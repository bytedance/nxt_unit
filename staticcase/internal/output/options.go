// Copyright cweill/gotests authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package goparse contains logic for parsing Go files. Specifically it parses
// source and test files into domain models for generating tests.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package output

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"io"
	"io/ioutil"
	"os"

	"github.com/bytedance/nxt_unit/atgconstant"

	"github.com/bytedance/nxt_unit/staticcase/internal/models"
	"github.com/bytedance/nxt_unit/staticcase/internal/render"

	"golang.org/x/tools/imports"
)

type Options struct {
	PrintInputs    bool
	Subtests       bool
	Parallel       bool
	Named          bool
	Template       string
	TemplateDir    string
	Mocks          map[string]map[string]int
	Builders       map[string][]string
	MiddleBuilders []string
	GlobalInit     []string
	MiddleCodeInit []string
	TestCaseNum    int
	TemplateParams map[string]interface{}
	TemplateData   [][]byte
	Uid            string
	MockStatements map[string][]string
	TestMode       string
	FilePath       string
	render         *render.Render
	UseMockType    int
	Ctx            context.Context // Context of the smart unit
}

func (o *Options) Process(head *models.Header, funcs []*models.Function) ([]byte, error) {
	o.render = render.New()
	switch {
	case o.providesTemplateDir():
		if err := o.render.LoadCustomTemplates(o.TemplateDir); err != nil {
			return nil, fmt.Errorf("loading custom templates: %v", err)
		}
	case o.providesTemplate():
		if err := o.render.LoadCustomTemplatesName(o.Template); err != nil {
			return nil, fmt.Errorf("loading custom templates of name: %v", err)
		}
	case o.providesTemplateData():
		o.render.LoadFromData(o.TemplateData)
	}

	tf, err := ioutil.TempFile("", "gotests_")
	if err != nil {
		return nil, fmt.Errorf("ioutil.TempFile: %v", err)
	}
	defer tf.Close()
	defer os.Remove(tf.Name())

	// create physical copy of test
	b := &bytes.Buffer{}
	if err := o.writeTests(b, head, funcs); err != nil {
		return nil, err
	}
	out, err := imports.Process(tf.Name(), b.Bytes(), nil)
	if err != nil {
		fmt.Println("imports Process final case file fail,detail file:\n" + b.String())
		return nil, fmt.Errorf("imports.Process: %v", err)
	}
	return out, nil
}

func (o *Options) providesTemplateData() bool {
	return o != nil && len(o.TemplateData) > 0
}

func (o *Options) providesTemplateDir() bool {
	return o != nil && o.TemplateDir != ""
}

func (o *Options) providesTemplate() bool {
	return o != nil && o.Template != ""
}

func (o *Options) writeTests(w io.Writer, head *models.Header, funcs []*models.Function) error {
	if path, ok := importsMap[o.Template]; ok {
		head.Imports = append(head.Imports, &models.Import{
			Path: fmt.Sprintf(`"%s"`, path),
		})
	}
	if o.TestMode == atgconstant.MiddleCode {
		head.Imports = append(head.Imports, &models.Import{
			Path: fmt.Sprintf(`"%s"`, "sync"),
		})
	}

	b := bufio.NewWriter(w)
	var headerTemplate string = "header"
	if err := o.render.Header(b, head, headerTemplate); err != nil {
		return fmt.Errorf("render.Header: %v", err)
	}

	if o.TestMode == atgconstant.MiddleCode {
		wgStr := fmt.Sprintf("wg%s := sync.WaitGroup{}\n", o.Uid)
		// get context pkgname
		ctxPkgName, _ := duplicatepackagemanager.GetInstance(o.Ctx).PutAndGet("", "context")
		// define the render path of final suite
		declPath := fmt.Sprintf("t.Parallel();originPath := \"%s\" \n declLocker := sync.RWMutex{} \n declData := map[string][]string{}\n useMockMap := map[string]map[string]int{}\n type DeclResult struct {\n\t\tAvailableList []bool\n\t\tPathSync  sync.Map\n\t} \n declStatistics := map[string]DeclResult{} \n smartUnitCtx := duplicatepackagemanager.SetInstance(%s.Background())\n", o.FilePath, ctxPkgName)
		var orginalImportStr string
		for index := range head.OriginalImports {
			orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, head.OriginalImports[index].Name, head.OriginalImports[index].Path)
		}

		// Add system import to avoid the renaming for them
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"errors\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"context\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"fmt\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"strings\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"sync\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"testing\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"time\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"syscall\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"runtime/debug\"")
		orginalImportStr = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet(\"%s\",%s)\n", orginalImportStr, "", "\"runtime/stack\"")
		var globalInitSet string
		var globalInit string
		for _, statement := range o.GlobalInit {
			globalInit = fmt.Sprintf("%s%s\n", globalInit, statement)
		}
		for _, statement := range o.MiddleCodeInit {
			globalInitSet = fmt.Sprintf("%sduplicatepackagemanager.GetInstance(smartUnitCtx).SetInitBuilder(\"%s\")\n", globalInitSet, statement)
		}
		_, err := b.WriteString("func TestD" + o.Uid + "(t *testing.T)  { \n " + wgStr + declPath + globalInit + orginalImportStr + globalInitSet)
		if err != nil {
			return fmt.Errorf("render.TestFunction: %v", err)
		}
	}
	for index := range funcs {
		if o.TestMode == atgconstant.FinalTest && funcs[index].RowData == "" {
			continue
		}
		mocks := []string{}
		builders := []string{}
		// need fullName to get mocks because of the same functionName but different receiver
		switch o.UseMockType {
		case atgconstant.UseNoMock:
			mocks = []string{}
		case atgconstant.UseMockitoMock:
			if statementMap, ok := o.Mocks[funcs[index].Name]; ok {
				statements := make([]string, 0)
				for statement, _ := range statementMap {
					statements = append(statements, statement)
				}
				mocks = statements
			}
		case atgconstant.UseGoMonkeyMock:
			funcName := funcs[index].Name
			if o.TestMode == atgconstant.MiddleCode {
				if statementMap, ok := o.Mocks[funcName]; ok {
					statements := make([]string, 0)
					for statement, _ := range statementMap {
						statements = append(statements, statement)
					}
					mocks = statements
				}
			}
			if o.TestMode == atgconstant.FinalTest {
				if funcs[index].Receiver != nil && funcs[index].Receiver.Type != nil {
					funcName = funcs[index].Receiver.Type.String() + funcName
				}
				if mockFuncNames, ok := o.Mocks[funcName]; ok {
					for mockFuncName, outCount := range mockFuncNames {
						patchName := fmt.Sprintf("%sPatch", atghelper.RandStringBytes(5))
						totalOut := ""
						for index := 0; index < outCount; index++ {
							separate := ""
							if index > 0 {
								separate = ","
							}
							totalOut = totalOut + fmt.Sprintf("%stt.MonkeyOutputMap[\"%s\"][%d]", separate, mockFuncName, index)
						}
						mockStateMent := fmt.Sprintf("%s := gomonkeyv2.ApplyFuncReturn(%s,%s)\n\t\tdefer %s.Reset()", patchName, mockFuncName, totalOut, patchName)
						mocks = append(mocks, mockStateMent)
					}
				}
			}
		default:
			mocks = []string{}
		}
		if statement, ok := o.Builders[funcs[index].Name]; ok {
			builders = statement
		}
		err := o.render.TestFunction(b, funcs[index], o.PrintInputs, o.Subtests, o.Named, o.Parallel, o.TemplateParams, mocks, builders, o.MiddleBuilders, o.TestCaseNum, o.UseMockType, o.Uid, funcs[index].RowData, o.TestMode, o.FilePath, o.GlobalInit)
		if err != nil {
			return fmt.Errorf("render.TestFunction: %v", err)
		}
	}
	if o.TestMode == atgconstant.MiddleCode {
		// render the testSuite
		writer := `
// summary data info wait for the result of function
wg%s.Wait()
to := time.After(time.Millisecond * 3000)
Loop:
	for {
		select {
		case info := <-WorkPipe%s:
            // get declData KeyName
			declFuncName := info.FunctionName
			if info.ReceiverName != "" {
				declFuncName = info.ReceiverName + declFuncName
			}
			if info.IsStart != "" {
				declFuncName = "*" + declFuncName
			}
			if _, dataOk := declData[declFuncName]; dataOk {
				sResult, ok := declStatistics[declFuncName]
				if ok {
					sResult.AvailableList = append(sResult.AvailableList, false)
				} else {
					sResult.AvailableList = []bool{false}
					sResult.PathSync = sync.Map{}
				}
				declStatistics[declFuncName] = sResult
			}

			if info.Coverage >= 0 {
				// trim panic coverage == -1
				if sResult, ok := declStatistics[declFuncName]; ok {
					fmt.Println("pathid " + info.PathID)
					if _, exist := sResult.PathSync.Load(info.PathID); !exist {
						// record unique path
						sResult.PathSync.Store(info.PathID, struct{}{})
						// update available case
						sResult.AvailableList[len(sResult.AvailableList)-1] = true
						declStatistics[declFuncName] = sResult
					}
				}
			}
		case <-to:
			break Loop
		}
	}
var hitLine int
for _, hit := range HitSet%s{
	if hit > 0 {
		hitLine ++
	}
}
fmt.Println(fmt.Sprintf("coverage(%%v;%%v)-r \n", len(HitSet%s), hitLine))
res := map[string]string{}
	for k, v := range declData {
		sResult, ok := declStatistics[k]
		if ok {
			dataList := make([]string, 0)
			for index, available := range sResult.AvailableList {
				if available {
					dataList = append(dataList, v[index])
				}
			}
			if len(dataList) > 0 {
				res[k] = fmt.Sprintf("[]test{%%s}", strings.Join(dataList, ","))
			}
		}
	}
code, err := staticcase.RecordFinalSuite(smartUnitCtx, originPath, res,useMockMap, %v)
if err != nil{
    t.Fatal(err)
}
code,err = imports.Process("",code,nil)
if err != nil{
    t.Fatal(err)
}
testName := strings.ReplaceAll(filepath.Base(originPath),".go","_ATG_test.txt")
testFile :=  path.Join(filepath.Dir(originPath), testName)
if err := ioutil.WriteFile(testFile, code, atgconstant.NewFilePerm); err != nil {
    t.Fatalf("[render testsuite] has ioutil.WriteFile err: %%v", err)
}
`
		writer = fmt.Sprintf(writer, o.Uid, o.Uid, o.Uid, o.Uid, o.UseMockType)
		_, err := b.WriteString(writer + "}\n")
		if err != nil {
			return err
		}
	}
	return b.Flush()
}
