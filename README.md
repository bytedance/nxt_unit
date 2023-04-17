# NxtUnit
`NxtUnit` is an automatically unit test generation application for Go.\
You can compile it as the binary package and run it.

[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue)](https://github.com/bytedance/nxt_unit/blob/master/LICENSE)
[![Go](https://github.com/bytedance/nxt_unit/actions/workflows/go.yml/badge.svg)](https://github.com/bytedance/nxt_unit/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/bytedance/nxt_unit/branch/main/graph/badge.svg)](https://codecov.io/gh/bytedance/nxt_unit)
## Table of Contents

- [Introduction](#Introduction)
- [How To Use](#How-To-Use)
- [Generation Failure](#Generation-Failure)
- [License](#License)

## Introduction

Automated unit test generation has been studied for a long time and prior research has focused on dynamically compiled or 
dynamically typed programming languages such as Java and Python. However, few of the existing tools support Go, 
which is a popular statically compiled and typed programming language in the industry for server application development 
and used extensively in our production environment

`NxtUnit` is the tool that can automatically generate the unit test for Go. For example, given the original code
```Go
func Example (input1 int, input2 int) {
   if input1*input2 > 9 {
      return input1
   }
   switch input1 {
   case 20:
      input1 = +RPCCallee1(input2)
   case 40:
      input1 = +RPCCallee1(input2)
   }
   return input1
}
```

During the generation, you might see our intermediate code:

```Go
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	atgconstant "github.com/bytedance/nxt_unit/atgconstant"
	contexthelper "github.com/bytedance/nxt_unit/atghelper/contexthelper"
	mockfunc "github.com/bytedance/nxt_unit/codebuilder/mock"
	variablecard "github.com/bytedance/nxt_unit/codebuilder/variablecard"
	duplicatepackagemanager "github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	staticcase "github.com/bytedance/nxt_unit/staticcase"
	convey "github.com/smartystreets/goconvey/convey"
	imports "golang.org/x/tools/imports"
)

func TestDAUAMW(t *testing.T) {
	wgAUAMW := sync.WaitGroup{}
	t.Parallel()
	originPath := "/Users/siweiwang/go/src/github.com/nxt_unit/siwei.go"
	declLocker := sync.RWMutex{}
	declData := map[string][]string{}
	useMockMap := map[string]map[string]int{}
	type DeclResult struct {
		AvailableList []bool
		PathSync      sync.Map
	}
	declStatistics := map[string]DeclResult{}
	smartUnitCtx := duplicatepackagemanager.SetInstance(context.Background())
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "errors")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "context")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "fmt")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "strings")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "sync")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "testing")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "time")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "syscall")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "runtime/debug")
	duplicatepackagemanager.GetInstance(smartUnitCtx).PutAndGet("", "runtime/stack")

	wgAUAMW.Add(1)
	go func(t *testing.T) {
		type Args struct {
			Input1 int
			Input2 int
		}
		type test struct {
			Name  string
			Args  Args
			Want  int
			Mocks variablecard.MocksRecord
		}
		defer func() {
			wgAUAMW.Done()
		}()
		tt := test{}
		duplicatepackagemanager.GetInstance(smartUnitCtx).SetRelativePath(tt)
		var rowData []string
		useMock := make(map[string]int, 0)
		for i := 0; i < 4; i++ {
			convey.Convey(tt.Name, t, func() {
				mockRender := &mockfunc.StatementRender{
					MockStatement:   []string{},
					MonkeyOutputMap: make(variablecard.MonkeyOutputMap, 0),
					UsedMockFunc:    make(map[string]int, 0),
				}
				smartUnitCtx = contexthelper.SetVariableContext(smartUnitCtx, atgconstant.VariableContext{})
				tt = variablecard.VariableMutate(smartUnitCtx, reflect.TypeOf(tt), reflect.ValueOf(tt)).Interface().(test)
				defer func() {
				}()
				if got := ExampleAUAMW(tt.Args.Input1, tt.Args.Input2); got != tt.Want {
					tt.Want = got
				}
				tt.Mocks = mockRender.MockStatement
				useMock = mockRender.UsedMockFunc
				rowData = append(rowData, variablecard.ValueToString(smartUnitCtx, reflect.ValueOf(tt)))
			})
		}
		if len(rowData) <= 0 {
			return
		}
		declLocker.Lock()
		declData["Example"] = rowData
		useMockMap["Example"] = useMock
		declLocker.Unlock()
	}(t)

	// summary data info wait for the result of function
	wgAUAMW.Wait()
	to := time.After(time.Millisecond * 3000)
Loop:
	for {
		select {
		case info := <-WorkPipeAUAMW:
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
	for _, hit := range HitSetAUAMW {
		if hit > 0 {
			hitLine++
		}
	}
	fmt.Println(fmt.Sprintf("coverage(%v;%v)-r \n", len(HitSetAUAMW), hitLine))
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
				res[k] = fmt.Sprintf("[]test{%s}", strings.Join(dataList, ","))
			}
		}
	}
	code, err := staticcase.RecordFinalSuite(smartUnitCtx, originPath, res, useMockMap, 0)
	if err != nil {
		t.Fatal(err)
	}
	code, err = imports.Process("", code, nil)
	if err != nil {
		t.Fatal(err)
	}
	testName := strings.ReplaceAll(filepath.Base(originPath), ".go", "_ATG_test.txt")
	testFile := path.Join(filepath.Dir(originPath), testName)
	if err := ioutil.WriteFile(testFile, code, atgconstant.NewFilePerm); err != nil {
		t.Fatalf("[render testsuite] has ioutil.WriteFile err: %v", err)
	}
}
```

it can generate the unit test like below

```Go
import (
   testing "testing"
   gomonkey "github.com/agiledragon/gomonkey/v2"
   convey "github.com/smartystreets/goconvey/convey" 
)
func TestExampleFunction_URRDGU(t *testing.T) {
   type Args struct {
      Input1 int,  Input2 int
   }
   type test struct {
      Name            string
      Args            Args
      Want            int
      Mocks           func()
      MonkeyOutputMap map[string][]interface{}
   }
   tests := []test{test{
      Name: "Alice King",
      Args: Args{
         Input1: 20, Input2: 4,
      },
      Want:            20,
      Mocks:           func() {},
      MonkeyOutputMap: map[string][]interface{}{"RPCCallee1": []interface{}{10}},
   }, test{
      Name: "Eason King",
      Args: Args{
         Input1: 7, Input2: 1,
      },
      Want:            7,
      Mocks:           func() {},
      MonkeyOutputMap: map[string][]interface{}{"RPCCallee1": []interface{}{11}},
   }}
   for _, tt := range tests {
      convey.Convey(tt.Name, t, func() {
         tt.Mocks()
         PTNFTPatch := gomonkey.ApplyFuncReturn(RPCCallee1, tt.MonkeyOutputMap["RPCCallee1"][0])
         defer PTNFTPatch.Reset()
    if got := ExampleFunction(tt.Args.Input1, tt.Args.Input2); got != tt.Want {
            convey.So(got, convey.ShouldResemble, tt.Want)
         }
      })
   }
}
```


## How To Use
### Installation
```
go install github.com/bytedance/nxt_unit@latest
```

### Usage
```
-function_name(required)
    function name
-receiver_name(optional)
    the receiver name of your function
-receiver_is_star(optional) 
    whether your receiver is a pointer or not
-usage(required)
    option1: generate the unit test
    option2: generate the template
-go(optional) 
    your local go path
-file_name(required)
    absolute go path
```
###  Example
```
go build
./nxt_unit -file_path=[your path] -receiver_name=Decoder -receiver_is_star=true -function_name=Decode -usage=plugin
-go=/usr/local/go/bin/go
```
### Run generated unit test
```
go test xxxx_test.go -gcflags "all=-N -l"
```
`-gcflags "all=-N -l"` used for unblocking the inlining of the function
## Failure Scenarios
The failure might be caused by the following reasons:
1. The function is not exported
2. The fault of the gomonkey
3. You don't have permission to execute the file. Please see the [Solution](#Solution)

### Solution for the gomonkey
1 download the tool
```
cd `go env GOPATH`
git clone https://github.com/eisenxp/macos-golink-wrapper.git
```
2 rename the link to original_link
```
mv `go env GOTOOLDIR`/link `go env GOTOOLDIR`/original_link 
```
3 copy tool to GOTOOLDIR
```
cp `go env GOPATH`/macos-golink-wrapper/link  `go env GOTOOLDIR`/link 
```
4 authorize link
```
chmod +x `go env GOTOOLDIR`/link
```


## License

`NxtUnit` is licensed under the terms of the Apache license 2.0. See [LICENSE](LICENSE) for more information.
