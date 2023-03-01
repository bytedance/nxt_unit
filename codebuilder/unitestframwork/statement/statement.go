/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package statement

import (
	"context"
	"fmt"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytedance/nxt_unit/atghelper"

	"github.com/bytedance/nxt_unit/atgconstant"

	"github.com/bytedance/nxt_unit/atghelper/extracinfo"
	"golang.org/x/tools/go/ssa"
)

type Statement struct {
	Name             string
	OriginalFunction *ssa.Function // Means the original function
	PkgPath          string
	PkgName          string
	Expression       string // Used for the mockito mock. Mockito.call(xxxx, xxx, xxx)
	SpecialType      string // Type: "overpass".   //TODO: add this in the atg constant
	FunctionType     string // Used to force transform the  function
}

// Create Mocked Statement
func CreateMockedStatement(ctx context.Context, f *ssa.Function, testedPkgName string) (*Statement, error) {
	MockedStatement := &Statement{}
	pkgPath := f.Pkg.Pkg.Path()
	relativePath := duplicatepackagemanager.GetInstance(ctx).RelativePath()

	if pkgPath != relativePath && !CheckMockInternalFunc(pkgPath, relativePath) {
		return nil, fmt.Errorf("can't mock function from other internal package")
	}
	// trim function which can't ValueToString by unexported result field
	for i := 0; i < f.Signature.Results().Len(); i++ {
		if pkgPath != duplicatepackagemanager.GetInstance(ctx).RelativePath() && !IsExportedTypeStr(f.Signature.Results().At(i).Type().String()) {
			return nil, fmt.Errorf("can't mock unexported result of function from other package")
		}
	}
	// Special handling of go vendor
	paths := strings.SplitN(pkgPath, "vendor/", 2)
	if len(paths) == 2 {
		pkgPath = paths[1]
	}
	pkgName := f.Pkg.Pkg.Name()

	temPkgName, temPkgPath := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, pkgPath)
	MockedStatement.Name = f.Name()
	MockedStatement.OriginalFunction = f
	MockedStatement.PkgName = temPkgName
	MockedStatement.PkgPath = temPkgPath
	mockedReceiver, err := extracinfo.ConvertToReceiver(ctx, f.Signature, testedPkgName)
	if err != nil {
		return nil, err
	}
	if mockedReceiver != nil {
		// Special Logic for overpass, it used the lower case word: rawCallStruct. Thus, it is not unexported.
		// However, it provides a way to mock it.
		// Example input: iestns_judgement_user_security.RawCall.MGetUserSecurityForSearch(ctx, req)
		// Expected output:
		// iestns_judgement_user_security.SetMock.MGetUserSecurityForSearch(func(ctx context.Context, req *judgeUserSecurity.MGetUserSecurityForSearchReq) (*judgeUserSecurity.MGetUserSecurityForSearchResp, error) {
		//	return &judgeUserSecurity.MGetUserSecurityForSearchResp{
		//		BaseResp: &jbase.BaseResp{},
		//	}, nil
		// })

		// Condition for entering this: (1) the receiver string is rawCallStruct (2) the pkg path has "/overpass/".
		if strings.Contains(f.Signature.Recv().Type().String(), ".rawCallStruct") && strings.Contains(temPkgPath, "/overpass/") {
			MockedStatement.Expression = fmt.Sprintf("%v.SetMock.%v", temPkgName, f.Name())
			MockedStatement.SpecialType = "overpass"
			MockedStatement.Name = f.Name()
			var builder strings.Builder
			builder.WriteString("func(")
			// f.Signature.Params().Len() should greater or equal to 1
			for i := 0; i < f.Signature.Params().Len()-1; i++ {
				// for example, the s is var ctx context.context
				s := f.Signature.Params().At(i).String()
				_, _, typeName := GetPkgNameAndPathAndTypeFromOverPass(ctx, s)
				builder.WriteString(typeName)
				if i == f.Signature.Params().Len()-2 {
					continue
				}
				builder.WriteString(",")
			}
			builder.WriteString(")(")
			for j := 0; j < f.Signature.Results().Len(); j++ {
				s := f.Signature.Results().At(j).String()
				_, _, typeName := GetPkgNameAndPathAndTypeFromOverPass(ctx, s)
				builder.WriteString(typeName)
				if j == f.Signature.Results().Len()-1 {
					continue
				}
				builder.WriteString(",")
			}
			builder.WriteString(")")
			MockedStatement.FunctionType = builder.String()
			return MockedStatement, nil
		}

		// we can't mock unexported value from other package
		if pkgPath != duplicatepackagemanager.GetInstance(ctx).RelativePath() && !IsExportedTypeStr(f.Signature.Recv().Type().String()) {
			return nil, fmt.Errorf("can't mock unexported function's from other package")
		}
		MockedStatement.Expression = fmt.Sprintf("%s%s", mockedReceiver.InitiateReceiverForMock(), f.Name())

	} else {
		if MockedStatement.PkgName != "" {
			MockedStatement.Expression = fmt.Sprintf("%v.%v", temPkgName, f.Name())
		} else {
			MockedStatement.Expression = f.Name()
		}
	}
	if pkgPath == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
		pkg := pkgName + "."
		MockedStatement.Expression = strings.Replace(MockedStatement.Expression, pkg, "", -1)
	}

	return MockedStatement, nil
}

func IsExportedTypeStr(eleTypeStr string) bool {
	if strings.Contains(eleTypeStr, "/") {
		eleTypeStr = eleTypeStr[strings.LastIndex(eleTypeStr, "/")+1:]
	}

	if strings.Contains(eleTypeStr, ".") {
		eleTypeList := strings.Split(eleTypeStr, ".")
		if len(eleTypeList) > 1 {
			typeName := eleTypeList[len(eleTypeList)-1]
			ch, _ := utf8.DecodeRuneInString(typeName)
			return unicode.IsUpper(ch)
		} else {
			return false
		}
	} else {
		return atghelper.IsBasicDataType(eleTypeStr)
	}
}

// GetPkgNameAndPathAndTypeFromOverPass Need to avoid the import cycle
func GetPkgNameAndPathAndTypeFromOverPass(ctx context.Context, s string) (string, string, string) {
	// for example, the s is var ctx context.context
	// Overpass input should not have map or array. It will only have *xxx/xxx/xxxx/xxxx/abc.b  or xxx/xxx/xxx/abc.a
	sArray := strings.SplitN(s, " ", -1)
	if len(sArray) <= 0 {
		return "", "", ""
	}
	s = sArray[len(sArray)-1]

	// first, get the variable name
	stringArray := strings.SplitN(s, ".", -1)
	if len(stringArray) <= 1 {
		return s, "", s
	}
	variable := stringArray[len(stringArray)-1]

	// second get the pkg path
	var prefix strings.Builder
	for i := 0; i < len(stringArray)-1; i++ {
		prefix.WriteString(stringArray[i])
		if i == len(stringArray)-2 {
			break
		}
		prefix.WriteString(".")
	}
	if strings.Contains(prefix.String(), "/") {
		path := strings.ReplaceAll(prefix.String(), "*", "")
		pkgArray := strings.SplitN(path, "/", -1)
		// Ignore the pkgArray = 0, because it contains "/"
		pkgName, pkgPath := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgArray[len(pkgArray)-1], path)
		if stringArray[0][0] == '*' {
			variable = "*" + pkgName + "." + variable
		} else {
			variable = pkgName + "." + variable
		}
		return pkgName, pkgPath, variable
	}
	if stringArray[0][0] == '*' {
		variable = "*" + stringArray[0] + "." + variable
	} else {
		variable = stringArray[0] + "." + variable
	}
	return stringArray[0], "", variable
}

func CheckMockInternalFunc(pkgPath, relativePath string) bool {
	check := true
	if strings.Contains(pkgPath, atgconstant.InternalPkg) {
		// internal pkg is relative repo path
		if strings.Contains(pkgPath, atghelper.GetRepoByRelativePath(relativePath)) {
			// if tcc path need mock
			if strings.Contains(pkgPath, "tcc") {
				check = true
			} else {
				check = false
			}
		} else {
			// internal pkg is other repo pkg not mock
			check = false
		}
	}
	return check
}
