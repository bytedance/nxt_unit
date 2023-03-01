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
package main

import (
	"flag"
	"fmt"
	mateAtgconstant "github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/manager/logextractor"
	"os"
	"os/exec"
	"path"
	"runtime/debug"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	matePkgManager "github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"github.com/bytedance/nxt_unit/manager/lifemanager"
	"github.com/bytedance/nxt_unit/staticcase"
)

var (
	funcName       = flag.String("function_name", "", "tested function name")
	filePath       = flag.String("file_path", "", `tested function file path`)
	debugMode      = flag.Bool("debug_mode", atgconstant.DebugMode, `used for debug`)
	minUnit        = flag.String("min_unit", atgconstant.MinUnit, `generate the unit tests for function or file`)
	usage          = flag.String("usage", "", "plugin: used for engineer")
	directoryPath  = flag.String("directory_path", "", "tested repository root directory")
	functionList   = flag.String("function_list", "", "passed in a list of function names")
	ReceiverName   = flag.String("receiver_name", "", "used to receive the receiver name")
	ReceiverIsStar = flag.Bool("receiver_is_start", false, "used to know the receiver has a pointer")
	templateType   = flag.Int("template_type", 0, "special template type")
	UseMockType    = flag.Int("use_mock_type", atgconstant.UseMockUnknown, "default is mockito. use nomock=1,mockito=2, gomonkey=3. gomonkey support go>=1.17")
	versionFlag    = flag.Bool("v", false, "Print the current version and exit")
	currentTag     = "unknown"
)

var help string

func Init() {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		currentTag = bi.Main.Version
	} else {
		logextractor.ExecutionLog.Log("Failed to get the current version")
	}
	flag.StringVar(&atgconstant.GoDirective, "go", "go", "Path or command of go directive")
	flag.StringVar(&help, "help", "help", "help message for nxt_unit flags")
	matePkgManager.Init()
}

func main() {
	Init()
	flag.Parse()
	defer func() {
		lifemanager.Closer.Close()
		logextractor.ExecutionLog.Print()
	}()

	if *versionFlag {
		logextractor.ExecutionLog.Log(currentTag)
		return
	}

	if *debugMode == true {
		cmd := exec.Command("bash", "-c", "go env")
		output, _ := cmd.CombinedOutput()
		goEnv := string(output)
		logextractor.ExecutionLog.Log("######################### go env ##########################")
		logextractor.ExecutionLog.Log(goEnv)
		logextractor.ExecutionLog.Log("####################################################################")
	}

	switch *usage {
	case atgconstant.PluginMode:
		err := Plugin()
		logextractor.ExecutionLog.LogFinalRes("######################### Conclusion ##########################")
		if err != nil {
			logextractor.ExecutionLog.LogFinalRes("Sorry, we cannot generate the test for you, Please check the error code above")
			logextractor.ExecutionLog.LogError(err.Error())
			return
		}
		logextractor.ExecutionLog.LogFinalRes("Successfully generate the unit test!")
		return
	case atgconstant.PluginQMode:
		err := Template()
		logextractor.ExecutionLog.LogFinalRes("######################### Conclusion ##########################")
		if err != nil {
			logextractor.ExecutionLog.LogFinalRes("Sorry, we cannot generate the template for you, Please check the error code above")
			logextractor.ExecutionLog.LogError(err.Error())
			return
		}
		logextractor.ExecutionLog.LogFinalRes("Successfully generate the unit test template!")
		return
	case atgconstant.Backend:
		fmt.Println("back stage task is deprecated")
		// BackStageTask()
	case atgconstant.SplitFunctionMode:
		SplitFunctionTask()
	default:
		fmt.Println("Please give your usage")
	}
}

// BackStageTask Deprecated.
// If you Need the backstage task restart, we need to fix the file mode.
func BackStageTask() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if *directoryPath != "" {
		dir = *directoryPath
	}
	err = staticcase.Work(dir, *functionList, GetUseMockType(dir))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("generating middle code")
	err = staticcase.WorkToChangeGo(dir, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("fix go file")
	err = staticcase.FixGoFile(dir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("execution middle code")
	err = staticcase.GenerateTest(dir, true)
	if err != nil {
		fmt.Println(err)
	}
}

func Plugin() error {
	var dir string
	if *directoryPath != "" {
		dir = *directoryPath
	} else {
		dir = path.Dir(*filePath)
	}
	if *ReceiverIsStar && *ReceiverName != "" {
		*ReceiverName = fmt.Sprint("*", *ReceiverName)
	}
	option := atgconstant.Options{
		FilePath:      *filePath,
		Level:         1,
		Maxtime:       4,
		MinUnit:       *minUnit,
		Uid:           atghelper.RandStringBytes(10),
		FuncName:      *funcName,
		DebugMode:     *debugMode,
		Usage:         *usage,
		DirectoryPath: dir,
		ReceiverName:  *ReceiverName,
		UseMockType:   GetUseMockType(dir),
	}
	var err error
	// warning :not delete println,plugin get necessary msg
	// logextractor.ExecutionLog.Log(fmt.Sprintf("plugin sdk use UID is %v\n", option.Uid))
	err = staticcase.UpdateSmartUnit(dir)
	if err != nil {
		return err
	}
	// Generate For file
	if option.MinUnit == "file" {
		err := staticcase.WorkOnFile(option)
		if err != nil {
			logextractor.ExecutionLog.Log(err.Error())
		}
		return nil
	}

	// Generate for function
	err = staticcase.WorkForPlugin(option)
	if err != nil {
		return err
	}
	err = staticcase.FixGoFile(dir)
	if err != nil {
		logextractor.ExecutionLog.Log(err.Error())
	}
	err = staticcase.GenerateTestForPlugin(option)
	if err != nil {
		return err
	}
	// fix go file again
	err = staticcase.FixGoFile(dir)
	if err != nil {
		logextractor.ExecutionLog.Log(err.Error())
	}
	return nil
}

func Template() error {
	var dir string
	if *directoryPath != "" {
		dir = *directoryPath
	} else {
		dir = path.Dir(*filePath)
	}
	if *ReceiverIsStar && *ReceiverName != "" {
		*ReceiverName = fmt.Sprint("*", *ReceiverName)
	}
	option := atgconstant.Options{
		FilePath:      *filePath,
		Level:         1,
		Maxtime:       4,
		MinUnit:       *minUnit,
		Uid:           atghelper.RandStringBytes(10),
		FuncName:      *funcName,
		DebugMode:     *debugMode,
		Usage:         *usage,
		DirectoryPath: dir,
		ReceiverName:  *ReceiverName,
		TemplateType:  *templateType,
		UseMockType:   GetUseMockType(dir),
	}
	var err error
	// fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.MiddleCodeGenerateError, err.Error())
	// warning :not delete println,plugin get necessary msg
	logextractor.ExecutionLog.Log(fmt.Sprintf("plugin sdk use UID is %v\n", option.Uid))

	// Generate For file
	// Generate For file
	if option.MinUnit == "file" {
		err := staticcase.WorkOnFile(option)
		if err != nil {
			logextractor.ExecutionLog.Log(err.Error())
		}
		return nil
	}

	// Generate for function
	ctx, err := staticcase.Init(&option)
	if err != nil {
		return err
	}
	err = staticcase.GenerateBaseTest(ctx)
	if err != nil {
		return err
	}
	// fix go file again
	err = staticcase.FixGoFile(dir)
	if err != nil {
		logextractor.ExecutionLog.Log(err.Error())
	}
	return nil
}

func SplitFunctionTask() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if *directoryPath != "" {
		dir = *directoryPath
	}
	err = staticcase.WorkForSplitFunction(dir, GetUseMockType(dir))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generate the unit test!")
}

func GetUseMockType(dir string) int {
	// switch *UseMockType {
	// case mateAtgconstant.UseMockUnknown:
	// 	return mateAtgconstant.UseGoMonkeyMock
	// }
	return mateAtgconstant.UseGoMonkeyMock
}
