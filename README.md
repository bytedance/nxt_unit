# NxtUnit

`NxtUnit` is an automatically unit test generation application for Go, you can compile it as the binary package and run it

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
```
func Example (input1 int, input2 int) 
int {
   if input1*input2 > 9 
   {
      return input1
   }
   switch input1 
   {
   case 20:
      input1 = +RPCCallee1(input2)
   case 40:
      input1 = +RPCCallee1(input2)
   }
   return input1
}
```

it can generate the unit test like below
```
import (
   testing "testing"
   gomonkeyv2 "github.com/agiledragon/gomonkey/v2"
   convey "github.com/smartystreets/goconvey/convey" )
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
         PTNFTPatch := gomonkeyv2.ApplyFuncReturn(RPCCallee1, tt.MonkeyOutputMap["RPCCallee1"][0])
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
git install https://github.com/bytedance/nxt_unit@latest
```

### Explanation of the parameters
```
-function_name means your function name
-receiver_name means the receiver name
-receiver_is_star means whether your receiver is a pointer or not.
-usage has two version: plugin is to generate the unit test. -pluginq is to generate the template
-go is your local go path
-file_name is your absolute go path
```
###  Generate the unit test
```
nxt_unit -file_path=[your path] -receiver_name=Decoder -receiver_is_star=true -function_name=Decode -usage=plugin
-go=/usr/local/go/bin/go
```
### Run generated unit test
Remenber to add go build flag: `-gcflags "all=-N -l"`
## Generation-Failure
### solutions
【panic: permission denied】if your GoVersion>=1.17 on macOS Catalina 10.15.x
### 1.download the tool
cd `go env GOPATH`
git clone https://github.com/eisenxp/macos-golink-wrapper.git
### 2. rename the link to original_link
mv `go env GOTOOLDIR`/link `go env GOTOOLDIR`/original_link 
### 3. copy tool to GOTOOLDIR
cp `go env GOPATH`/macos-golink-wrapper/link  `go env GOTOOLDIR`/link 
### 4. authorize link
chmod +x `go env GOTOOLDIR`/link


## License

`NxtUnit` is licensed under the terms of the Apache license 2.0. See [LICENSE](LICENSE) for more information.