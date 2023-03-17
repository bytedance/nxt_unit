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
package atghelper

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/bytedance/nxt_unit/atgconstant"
)

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func GetRepoByRelativePath(relativePath string) string {
	results := strings.Split(relativePath, "/")
	if len(results) >= 3 {
		return strings.Join(results[0:3], "/")
	} else {
		return strings.Join(results[0:], "/")
	}
}

func RemoveTestCaseFile(dirName, fileName string) {
	GlobalFileLock.Lock()
	defer func() {
		GlobalFileLock.Unlock()
	}()
	path := path.Join(atgconstant.ProjectPath, dirName, fileName)
	os.Remove(path)
}

func RemoveDirectory(dirName string) {
	GlobalFileLock.Lock()
	defer func() {
		GlobalFileLock.Unlock()
	}()
	path := path.Join(atgconstant.ProjectPath, dirName)
	os.RemoveAll(path)
}

// UpperCaseFirstLetter won't change the input string.
func UpperCaseFirstLetter(s string) string {
	news := DeepCopy(s)
	if len(s) == 0 {
		return ""
	}
	return fmt.Sprint(strings.ToUpper(news[0:1]) + news[1:len(s)])
}

func DeepCopy(s string) string {
	var sb strings.Builder
	sb.WriteString(s)
	return sb.String()
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = atgconstant.Letters[rand.Intn(len(atgconstant.Letters))]
	}
	return string(b)
}

func DeepCopySlice(bytes []byte) []byte {
	cpy := make([]byte, len(bytes))
	copy(cpy, bytes)
	return cpy
}

func RemoveGoStyleComments(content []byte, comments []string) []byte {
	for _, cmt := range comments {
		matched, err := regexp.Compile(`//.*` + regexp.QuoteMeta(cmt))
		if err != nil {
			continue
		}
		content = matched.ReplaceAll(content, []byte(""))
	}
	return content
}

func RemoveCStyleComments(content []byte) []byte {
	// http://blog.ostermiller.org/find-comment
	cCmt, err := regexp.Compile(`/\*([^*]|[\r\n]|(\*+([^*/]|[\r\n])))*\*+/`)
	if err != nil {
		return content
	}
	return cCmt.ReplaceAll(content, []byte(""))
}

func RemoveGoComments(path string) ([]byte, error) {
	origin, err := ioutil.ReadFile(path)
	if err != nil {
		return origin, fmt.Errorf("[RemoveGoComments] read file err: %v", err)
	}
	cpy := DeepCopySlice(origin)
	RemoveCStyleComments(cpy)
	comments, err := LintPackageComment(path)
	if err != nil {
		return origin, fmt.Errorf("[RemoveGoComments] LintPackageComment err: %v", err)
	}
	content := RemoveGoStyleComments(cpy, comments)
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return origin, fmt.Errorf("[RemoveGoComments] cannot write the file err: %v", err)
	}
	return origin, nil
}

func RecoverOriginalFile(content []byte, path string) error {
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return fmt.Errorf("[RecoverOriginalFile] cannot write the file err: %v", err)
	}
	return nil
}

// get all comments
func LintPackageComment(path string) ([]string, error) {
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[LintPackageComment] cannot read the file")
	}
	parsedFile, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("[LintPackageComment] cannot parse the file")
	}
	res := make([]string, 0)
	for _, cg := range parsedFile.Comments {
		temp := strings.ReplaceAll(cg.Text(), "\n", "")
		temp = strings.ReplaceAll(temp, "\t", "")
		temp = strings.ReplaceAll(temp, "\a", "")
		res = append(res, temp)
	}
	return res, nil
}

var GoKeywords = []string{
	"interface", "break", "default", "func", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range ", "type",
	"continue", "for", "import", "return", "var"}

func GetPkgName(path string) string {
	if strings.Contains(path, "/") {
		stringArray := strings.SplitN(path, "/", -1)
		res := stringArray[len(stringArray)-1]
		if strings.Contains(res, "-") {
			res = strings.ReplaceAll(res, "-", "_")
		} else if Contains(GoKeywords, res) {
			res = "_" + res
		}
		return res
	} else {
		return path
	}
}

func ReplacePkgName(s string, pkgName string, originalPkgName string) string {
	pathMatch, err := regexp.Compile(`^(.*)([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS := pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		nameMath, err := regexp.Compile(`([a-zA-Z0-9_-]+)\.`)
		if err == nil {
			realNames := nameMath.FindAllString(matchedS[0], 1)
			if len(realNames) == 1 {
				realPkgName := strings.Trim(realNames[0], ".")
				// use real pkg name not cut from path
				if realPkgName != originalPkgName && pkgName != "" {
					originalPkgName = realPkgName
				} else if realPkgName != originalPkgName && pkgName == "" {
					originalPkgName = realNames[0]
				}
			}
		}
		left := strings.Replace(matchedS[0], originalPkgName, pkgName, 1)
		return fmt.Sprint(left, s[len(matchedS[0]):])
	}
	return s
}

func ReplacePkgNameForMap(s string, pkgName string, originalPkgName string, fromLeft bool) string {
	switch fromLeft {
	case true:
		pathMatch, err := regexp.Compile(`^(.*)([a-zA-Z0-9_-]+)\.`)
		if err != nil {
			fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
			return s
		}
		matchedS := pathMatch.FindAllString(s, 1)
		if len(matchedS) == 1 {
			left := strings.Replace(matchedS[0], originalPkgName, pkgName, 1)
			return fmt.Sprint(left, s[len(matchedS[0]):])
		}
		return s
	default:
		pathMatch, err := regexp.Compile(`\](.*)([a-zA-Z0-9_-]+)$`)
		if err != nil {
			fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
			return s
		}
		matchedS := pathMatch.FindAllString(s, 1)
		if len(matchedS) == 1 {
			right := strings.Replace(matchedS[0], originalPkgName, pkgName, 1)
			return fmt.Sprint(s[0:len(matchedS[0])], right)
		}
		return s
	}
}

func IsSystemPkg(path string) bool {
	return !strings.Contains(path, "/")
}

// The receive name like:
//  var github.com/bytedance/nxt_unit/unittesttemplate.TikTokConsumption
//  var *github.com/bytedance/nxt_unit/unittesttemplate.TikTokConsumption
func GetTheReceiveNameFromSSA(s string) string {
	s = strings.TrimPrefix(s, "var")
	s = strings.ReplaceAll(s, " ", "")
	lastElement := path.Base(s)
	res := ""
	if strings.Contains(lastElement, ".") {
		resArray := strings.SplitN(lastElement, ".", -1)
		res = resArray[len(resArray)-1]
	}
	if strings.Contains(s, "*") {
		return fmt.Sprint("*", res)
	}
	return res
}
func IsTypeExported(t reflect.Type) bool {
	if !strings.Contains(t.String(), ".") {
		return true
	}
	nArray := strings.SplitN(t.String(), ".", -1)
	if len(nArray) > 0 {
		for index, letter := range nArray[1] {
			if index >= 1 {
				break
			}
			return unicode.IsUpper(letter)
		}
		return false
	}
	return true
}

// TODO: @caoziguang GetPkgName does not work properly in some scenarios.
// It should be replace by GetPkgNameV2

// get pkgname of reflect.Type
func GetPkgNameV2(strExp string) string {
	if !strings.Contains(strExp, ".") {
		return strExp
	}
	nameList := strings.SplitN(strExp, ".", 2)
	return nameList[0]
}

func IsBasicDataType(t string) bool {
	switch t {
	case "bool", "string", "int", "int8", "int16", "int32", "int64", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune",
		"float32", "float64", "complex64", "complex128", "error":
		return true
	default:
		return false
	}
}
