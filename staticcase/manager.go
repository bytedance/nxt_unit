package staticcase

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/bytedance/nxt_unit/manager/logextractor"

	"github.com/bytedance/nxt_unit/codebuilder/instrumentation"

	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"golang.org/x/tools/go/ssa"

	"github.com/bytedance/nxt_unit/manager/reporter"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/codebuilder/setup"
)

func Init(option *atgconstant.Options) (context.Context, error) {
	if option.FileName == "" {
		option.FileName = path.Base(option.FilePath)
	}
	return getContext(*option)
}

func getContext(option atgconstant.Options) (ctx context.Context, err error) {
	ctx = context.Background()
	ctx = contexthelper.SetOption(ctx, option)
	ctx = contexthelper.SetBuilderVector(ctx, option.Uid)
	ctx = duplicatepackagemanager.SetInstance(ctx)
	switch option.MinUnit {
	case "file":
		// Get all function in the file
		funcNameArray, err := instrumentation.GetAllFunctionInFile(option)
		if err != nil {
			return nil, fmt.Errorf("the error belongs to %w, the detail: %v", logextractor.CannotParseTestedFunctionError, err.Error())
		}
		funcInfo := &struct {
			pkgPath string
			exist   bool
		}{}
		funcsMap := make(map[string]setup.Functions, 0)
		constructorMap := make(map[string][]*ssa.Function)
		for _, funcName := range funcNameArray {
			option.FuncName = funcName
			sourceFunc, err := setup.GetFunctions(option)
			if err == nil {
				if !funcInfo.exist {
					ctx = contexthelper.SetSetupFunc(ctx, sourceFunc)
					funcInfo.exist = true
					funcInfo.pkgPath = sourceFunc.TestFunction.Program.PkgPath
				}
				funcsMap[funcName] = sourceFunc
				for t, funcList := range setup.GetConstructorsByFunc(sourceFunc) {
					constructorMap[t] = funcList
				}
			} else {
				fmt.Printf("GetFunctions  FilePath %v,err %v\n", option.FilePath, err)
			}
		}
		if !funcInfo.exist {
			return nil, logextractor.CannotFindTestedFunctionError
		}
		duplicatepackagemanager.GetInstance(ctx).SetRelativeString(funcInfo.pkgPath)
		ctx = contexthelper.SetSetupFuncMap(ctx, funcsMap)
		ctx = contexthelper.SetConstructorFuncMap(ctx, constructorMap)
	case "function":
		funcInfo := &struct {
			pkgPath string
			exist   bool
		}{}
		funcMap := map[string]setup.Functions{}
		sourceFunc, err := setup.GetFunctions(option)
		if err != nil {
			return nil, fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.CannotParseTestedFunctionError, err.Error())
		}
		if !funcInfo.exist {
			ctx = contexthelper.SetSetupFunc(ctx, sourceFunc)
			funcInfo.exist = true
			funcInfo.pkgPath = sourceFunc.TestFunction.Program.PkgPath
		}
		funcMap[sourceFunc.TestFunction.Function.Name()] = sourceFunc
		constructorMap := setup.GetConstructorsByFunc(sourceFunc)
		duplicatepackagemanager.GetInstance(ctx).SetRelativeString(sourceFunc.TestFunction.Program.PkgPath)
		ctx = contexthelper.SetSetupFuncMap(ctx, funcMap)
		ctx = contexthelper.SetConstructorFuncMap(ctx, constructorMap)
	case "multifunction":
		funcInfo := &struct {
			pkgPath string
			exist   bool
		}{}
		funcsMap := map[string]setup.Functions{}
		constructorMap := make(map[string][]*ssa.Function)
		for _, funcName := range option.FunctionList {
			option.FuncName = funcName
			sourceFunc, err := setup.GetFunctions(option)
			if err == nil {
				if !funcInfo.exist {
					ctx = contexthelper.SetSetupFunc(ctx, sourceFunc)
					funcInfo.exist = true
					funcInfo.pkgPath = sourceFunc.TestFunction.Program.PkgPath
				}
				funcsMap[funcName] = sourceFunc
				for t, funcList := range setup.GetConstructorsByFunc(sourceFunc) {
					constructorMap[t] = funcList
				}
			} else {
				fmt.Printf("GetFunctions  FilePath %v,funcArray %v,err %v\n", option.FilePath, option.FunctionList, err)
			}
		}
		if !funcInfo.exist {
			return nil, logextractor.CannotFindTestedFunctionError
		}
		ctx = contexthelper.SetSetupFuncMap(ctx, funcsMap)
		ctx = contexthelper.SetConstructorFuncMap(ctx, constructorMap)
	default:
		sourceFunc, err := setup.GetFunctions(option)
		if err != nil {
			return nil, fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.CannotParseTestedFunctionError, err.Error())
		}
		fmt.Println("tested function parsed")
		ctx = contexthelper.SetSetupFunc(ctx, sourceFunc)
		constructorMap := setup.GetConstructorsByFunc(sourceFunc)
		ctx = contexthelper.SetConstructorFuncMap(ctx, constructorMap)
	}
	return ctx, nil
}

// Plugin doesn't need to report the message. Only the debug mode, we can see the message.
func ReportInternalError(option atgconstant.Options, err error, panicInfo string) {
	if option.Usage == "plugin" && !option.DebugMode {
		return
	}
	if err != nil {
		reporter.InternelErrorReporter.AddErrorFunction(option, err)
	}
	// report ci message
	fmt.Println(reporter.BugReporter.Report(option))
	fmt.Println(reporter.InternelErrorReporter.Report())
	if panicInfo != "" {
		fmt.Printf("panic_info(%s)-r\n", panicInfo)
	}
	time.Sleep(time.Second)
}

func checkFileName(fileName string) bool {
	switch {
	case strings.Contains(fileName, "kitex_gen"):
		return false
	case strings.Contains(fileName, "thrift_gen"):
		return false
	case strings.Contains(fileName, "/vendor/"):
		return false
	case strings.Contains(fileName, "_test/"):
		return false
	case strings.Contains(fileName, "_test.go"):
		return false
	case strings.HasSuffix(fileName, ".go"):
		return true
	}
	return false
}
