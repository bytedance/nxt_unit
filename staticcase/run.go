package staticcase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/nxt_unit/manager/logextractor"

	"github.com/bytedance/nxt_unit/codebuilder/instrumentation"
	"github.com/bytedance/nxt_unit/manager/reporter"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/atghelper/utils"
	"github.com/bytedance/nxt_unit/manager/lifemanager"
)

func Run(done func(), option atgconstant.Options) int {
	var err error
	var totalLine int
	var panicInfo string
	go func() {
		<-time.After(time.Second * atgconstant.WithoutGATimeOut)
		done()
	}()
	defer func() {
		if err := recover(); err != nil {
			done()
			// todo temp file manager
			logextractor.ExecutionLog.Log(fmt.Sprintf("[Run] encounter the panic, the panic is %v\n", err))
			reporter.InternelErrorReporter.AddErrorFunction(option, fmt.Errorf("%v", err))
		}
		done()
		ReportInternalError(option, err, panicInfo)
	}()

	ctx, err := Init(&option)
	if err != nil {

		return totalLine
	}
	totalLine, err = CreatMiddleCode(ctx)
	if err != nil {
		fmt.Printf("Sorry, we cannot add the logging to your original file, the error is %v\n", err.Error())
	}
	return totalLine
}
func CreatMiddleCode(ctx context.Context) (int, error) {
	tb := instrumentation.NewFunctionBuilder(ctx)
	// build instructionFile
	instructionFile, err := tb.Build(ctx)
	// if middle code error we need remove instructionFile
	lifemanager.Closer.SetClose(func() {
		_ = os.Remove(instructionFile)
	})
	if err != nil {
		return tb.TotalLine, err
	}
	opt, _ := contexthelper.GetOption(ctx)
	te := NewTestsuiteEntry(ctx, opt.FilePath, atghelper.GlobalFileLock)
	testFile, err := te.Build(ctx)
	lifemanager.Closer.SetClose(func() {
		_ = os.Remove(testFile)
	})
	if err != nil {
		return tb.TotalLine, err
	}
	// for plugin mod we change text to code immediately
	if opt.Usage == atgconstant.PluginMode || opt.Usage == atgconstant.SplitFunctionMode {
		atghelper.PluginSDK = atghelper.NewPluginSDK(instructionFile, testFile)
		err := os.Rename(instructionFile, atghelper.PluginSDK.InstructionFile())
		if err != nil {
			return tb.TotalLine, err
		}
		lifemanager.Closer.SetClose(func() {
			os.Remove(atghelper.PluginSDK.InstructionFile())
		})
		err = os.Rename(testFile, atghelper.PluginSDK.AtgTestFile())
		if err != nil {
			return tb.TotalLine, err
		}
		lifemanager.Closer.SetClose(func() {
			os.Remove(atghelper.PluginSDK.AtgTestFile())
		})
	}
	return tb.TotalLine, err
}

func RunForPlugin(option atgconstant.Options) error {
	var err error
	var panicInfo string
	defer func() {
		if err := recover(); err != nil {
			logextractor.ExecutionLog.Log(fmt.Sprintf("[Run] encounter the panic, the panic is %v\n", err))
			reporter.InternelErrorReporter.AddErrorFunction(option, fmt.Errorf("%v", err))
		}
		ReportInternalError(option, err, panicInfo)
	}()

	ctx, err := Init(&option)
	if err != nil {
		return err
	}
	_, err = CreatMiddleCode(ctx)
	if err != nil {
		return fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.MiddleCodeGenerateError, err.Error())
	}
	return nil
}

func WorkForPlugin(opt atgconstant.Options) error {
	if !checkFileName(opt.FilePath) {
		return logextractor.LocalFileNotSupportedError
	}
	return RunForPlugin(opt)
}

func Work(path string, functionList string, useMockType int) error {
	wg := sync.WaitGroup{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !checkFileName(path) {
			return nil
		}
		option := atgconstant.Options{
			FilePath:     path,
			Level:        1,
			Maxtime:      4,
			MinUnit:      "file",
			Uid:          atghelper.RandStringBytes(10),
			FunctionList: strings.SplitN(functionList, ",", -1),
			UseMockType:  useMockType,
		}
		wg.Add(1)
		done := sync.Once{}
		Run(func() {
			done.Do(wg.Done)
		}, option)
		return nil
	})
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func WorkToChangeGo(dir string, shouldPrintLog bool) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		ReNameFileToGo(path, "middle_code.txt", "middle_code_test.go", true, shouldPrintLog)
		ReNameFileToGo(path, "smartunit.txt", ".go", true, shouldPrintLog)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ReNameFileToGo(path, oldSuffix, newSuffix string, isTemp bool, shouldPrintLog bool) error {
	if strings.Contains(path, oldSuffix) {
		realName := strings.ReplaceAll(path, oldSuffix, newSuffix)
		err := os.Rename(path, realName)
		if err != nil && shouldPrintLog {
			fmt.Println(err)
		}
		if isTemp {
			lifemanager.Closer.SetClose(
				func() {
					os.Remove(realName)
				})
		}
		return nil
	}
	return nil
}

func FixGoFile(path string) error {
	var stdBuffer, stdErrBuff bytes.Buffer
	cmd := exec.Command(atgconstant.GoDirective, "mod", "tidy")
	cmd.Dir = path
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdErrBuff
	err := cmd.Run()
	if err != nil {
		return logextractor.GenSuggestByErrLog(err, stdBuffer.String(), stdErrBuff.String())
	}
	return nil
}

// GenerateTestForPlugin the test for the plugin. It is firstly generate the test with the _smart_unit_test_[random_number].go
// If there is the _smart_unit_test.go file. We merge these two file. If there is no _smart_unit_test.go file, we remove the
// random number.
func GenerateTestForPlugin(opt atgconstant.Options) error {
	_ = PluginCmd(opt)
	// We don't reply on the error from the command because it is not accurate. Some go test will run if there are
	// panics happened in the code.
	// We only monitor if there is the file path exist. If it is existed. We think it is succeeded.
	testText := strings.ReplaceAll(opt.FilePath, ".go", "_ATG_test.txt")
	if !atghelper.IsFileExist(testText) {
		return fmt.Errorf("the error belongs to %w, please check out the debug info", logextractor.MiddleCodeCannotGenerateFinalCodeError)
	}
	targetFilePath := getTempFilePath(opt, "_nxt_unit_test.go")
	tempFileName := strings.ReplaceAll(testText, "_ATG_test.txt", fmt.Sprint("_", opt.Uid, "_nxt_unit_test", ".go"))
	lifemanager.Closer.SetClose(
		func() {
			_ = os.Remove(tempFileName)
			_ = os.Remove(testText)
		})

	switch {
	// if GenerateTest at once, rename tempFile to target name
	case !atghelper.IsFileExist(targetFilePath):
		err := os.Rename(testText, targetFilePath)
		if err != nil {
			return fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.GenerateTestRenameError, err.Error())
		}
		// Execute temp file for the backend script. If it is not runnable, we return the error.
		// opt.Usage == atgconstant.PluginMode is here because generate the for file is also going
		// through this logic. Therefore, we should keep it otherwise generate for file will fail by one failed case.
		if opt.Usage == atgconstant.SplitFunctionMode || opt.Usage == atgconstant.PluginMode {
			_ = FixGoFile(opt.DirectoryPath)
			opt.RunForFinalSuite = true // Only works once
			opt.FinalSuiteTestName = atgconstant.FinalSuiteTestName
			// TestCommonPidBiddingHandler_TransPidBid_ATG
			err = PluginCmd(opt)
			if err != nil {
				// lifemanager.Closer.SetClose(
				// 	func() {
				// 		os.Remove(targetFilePath)
				// 	})
				// Doesn't need to prompt the error because it is not accurate.
				return fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.GenerateTestNotRunnableError, err.Error())
			}
			opt.RunForFinalSuite = false
		}
		return nil
	}
	err := os.Rename(testText, tempFileName)

	if err != nil {
		return fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.GenerateTestRenameError, err.Error())
	}
	// Execute temp file for the backend script. If it is not runnable, we return the error.
	if opt.Usage == atgconstant.SplitFunctionMode || opt.Usage == atgconstant.PluginMode {
		_ = FixGoFile(opt.DirectoryPath)
		opt.RunForFinalSuite = true // Only works once
		opt.FinalSuiteTestName = atgconstant.FinalSuiteTestName
		err = PluginCmd(opt)
		if err != nil {
			return fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.GenerateTestNotRunnableError, err.Error())
		}
		opt.RunForFinalSuite = false
	}
	targetImports, err := instrumentation.GetImportsInfosFromFile(targetFilePath)
	if err != nil {
		targetImports = make([]*instrumentation.Import, 0)
	}
	err = instrumentation.Concatenate(tempFileName, targetFilePath, targetImports)
	if err != nil {
		return fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.MergeTestConflictError, err.Error())
	}
	return nil
}

func PluginCmd(opt atgconstant.Options) error {
	var stdBuffer bytes.Buffer
	var stdErrBuff bytes.Buffer
	var cmd *exec.Cmd
	defer func() {
		if opt.Usage == atgconstant.PluginMode {
			if stdBuffer.Len() > 0 {
				logextractor.CommandAnalyze(stdBuffer.String())
			}
			if stdErrBuff.Len() > 0 {
				logextractor.CommandAnalyze(stdErrBuff.String())
			}
		}
		logextractor.ExecutionLog.DebugInfo(stdBuffer.String())
		logextractor.ExecutionLog.DebugInfo(stdErrBuff.String())
	}()
	if opt.RunForFinalSuite {
		cmd = exec.Command(atgconstant.GoDirective, "test", "-gcflags=all=-N -l", "-v", "-vet=off", "-count=1", "-timeout=80s", fmt.Sprintf("-test.run=%s", opt.FinalSuiteTestName))
	} else {
		switch opt.MinUnit {
		case atgconstant.MinUnit, atgconstant.FileMode:
			cmd = exec.Command(atgconstant.GoDirective, "test", "-gcflags=all=-N -l", "-v", "-vet=off", "-count=1", "-timeout=100s", fmt.Sprintf("-test.run=%s", opt.Uid))
		default:
			cmd = exec.Command(atgconstant.GoDirective, "test", "-gcflags=all=-N -l", "-v", "-vet=off", "-count=1", "-timeout=80s", "./...")
		}
	}
	cmd.Env = os.Environ()
	cmd.Dir = opt.DirectoryPath
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdErrBuff
	err := cmd.Run()
	return err
}

func GenerateTest(dir string, shouldPrintLog bool) error {
	var stdBuffer, stdErrBuff bytes.Buffer
	defer func() {
		if !shouldPrintLog {
			return
		}
		fmt.Println(stdBuffer.String())
		fmt.Println(stdErrBuff.String())
	}()
	ChangeFailFileSuffix(dir, "middle_code_test.go", "middle_code_fail_test.txt")
	cmd := exec.Command(atgconstant.GoDirective, "test", "-gcflags=\"all=-N -l\"", "-v", "-vet=off", "-timeout=40s", "./...")
	cmd.Dir = dir
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdErrBuff
	err := cmd.Run()
	if err != nil && shouldPrintLog {
		fmt.Println(err)
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		_ = ReNameFileToGo(path, "_ATG_test.txt", "_nxt_unit_test.go", false, shouldPrintLog)
		return nil
	})
	reporter.BugReporter.Analysis(stdBuffer.String())
	panicInfo := strings.ReplaceAll(reporter.BugReporter.GeneratePanicInfo(), "\n", "super&&world")
	fmt.Printf("panic_info(%s)-r\n", panicInfo)

	err = ChangeFailFileSuffix(dir, "_nxt_unit_test.go", "_final_code_fail_test.txt")
	return err
}

func ChangeFailFileSuffix(dir string, oldSuffix, newSuffix string) error {
	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	trimCheck := func(input interface{}) (bool, error) {
		cmdDir, ok := input.(string)
		if !ok {
			return true, errors.New("input is not string")
		}
		var stdBuffer, stdErrBuff bytes.Buffer
		cmd := exec.Command(atgconstant.GoDirective, "vet", "-composites=false", "-tests=false", "-unreachable=false", "-unusedresult=false", "-assign=false", "-structtag=false", "./...")
		cmd.Dir = cmdDir
		cmd.Stdout = &stdBuffer
		cmd.Stderr = &stdErrBuff
		err := cmd.Run()
		if (stdErrBuff.String() == "" && stdBuffer.String() == "") || err == nil {
			return true, err
		}

		if stdErrBuff.String() != "" {
			infos := strings.Split(stdErrBuff.String(), "\n")
			for _, info := range infos {
				fmt.Println(info)
				if strings.Contains(info, ".go") {
					match, err := regexp.Compile(`vet: (.*?):`)
					if err == nil {
						strs := match.FindAllStringSubmatch(info, 1)
						for index := range strs {
							if len(strs[index]) >= 2 {
								fileName := strings.TrimLeft(strs[index][1], ".")
								filePath := path.Join(dir, fileName)
								if strings.Contains(filePath, oldSuffix) {
									err = os.Rename(filePath, strings.ReplaceAll(filePath, oldSuffix, newSuffix))
									if err != nil {
										fmt.Println(err)
									}
								}
							}
						}
					}
				}
			}
			return false, err
		}
		return true, err
	}
	utils.RetryDo("trimCheck", 100, 10*time.Minute, trimCheck, dir)

	return nil
}

// we are going to remove Vector_middle_code_test.go, [path]_BridgeVector.go, middle_code_fail_test.txt
// _final_code_fail_test.txt.
func RemoveGeneratedFile(opt atgconstant.Options) {
	fileName := strings.ReplaceAll(path.Base(opt.FilePath), ".go", "")
	fileDir := path.Dir(opt.FilePath)
	_ = os.Remove(path.Join(fileDir, fmt.Sprint(fileName, "_final_code_fail_test.txt")))
	_ = os.Remove(path.Join(fileDir, fmt.Sprint(fileName, "Vector.go")))
	_ = os.Remove(path.Join(fileDir, "middle_code_fail_test.txt"))
	_ = os.Remove(path.Join(fileDir, "Vector_middle_code_test.go"))
}

func CheckEmptyFileAndThenRemove(filePath string, shouldPrintLog bool) error {
	fileDir := path.Dir(filePath)
	fileName := strings.ReplaceAll(path.Base(filePath), ".go", "")
	file, err := os.Open(path.Join(fileDir, fmt.Sprint(fileName, "_ATG_test.txt")))
	if err != nil && shouldPrintLog {
		fmt.Println(err)
	}
	b, _ := ioutil.ReadAll(file)
	if strings.Count(string(b), "\n") <= 2 {
		os.Remove(path.Join(fileDir, fmt.Sprint(fileName, "_ATG_test.txt")))
		// TODO(siwei.wang) add disable error here
		return fmt.Errorf("[CheckEmptyFileAndThenRemove] cannot generate the test for file: %v", filePath)
	}
	return nil
}

func getTempFilePath(opt atgconstant.Options, suffix string) string {
	fileName := path.Base(opt.FilePath)
	newFileName := strings.ReplaceAll(fileName, ".go", "")
	return path.Join(opt.DirectoryPath, fmt.Sprint(newFileName, suffix))
}

func WorkOnFile(option atgconstant.Options) error {
	if !checkFileName(option.FilePath) {
		return logextractor.LocalFileNotSupportedError
	}
	wg := sync.WaitGroup{}
	funcNameArray, err := instrumentation.GetAllFunctionInFile(option)
	if err != nil {
		return err
	}
	option.FunctionList = funcNameArray
	wg.Add(1)
	done := sync.Once{}
	RunSplitFunction(func() {
		done.Do(wg.Done)
	}, option)
	wg.Wait()
	return nil
}

func WorkForSplitFunction(path string, useMockType int) error {
	wg := sync.WaitGroup{}
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if !checkFileName(filePath) {
			return nil
		}
		funcNameArray, errcheck := instrumentation.GetAllFunctionInFile(atgconstant.Options{FilePath: filePath})
		if errcheck != nil {
			return errcheck
		}

		option := atgconstant.Options{
			FilePath:     filePath,
			Level:        1,
			Maxtime:      4,
			Uid:          atghelper.RandStringBytes(10),
			FunctionList: funcNameArray,
			UseMockType:  useMockType,
		}

		wg.Add(1)
		done := sync.Once{}
		RunSplitFunction(func() {
			done.Do(wg.Done)
		}, option)

		return nil
	})
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func RunSplitFunction(done func(), option atgconstant.Options) {
	var err error
	var panicInfo string
	go func() {
		<-time.After(time.Second * atgconstant.WithoutGATimeOut)
		done()
		defer func() {
			if err := recover(); err != nil {
				done()
				fmt.Printf("[Run] encounter the panic, the panic is %v\n", err)
				reporter.InternelErrorReporter.AddErrorFunction(option, fmt.Errorf("%v", err))
			}
			done()
			ReportInternalError(option, err, panicInfo)
		}()
	}()

	for _, functionName := range option.FunctionList {
		RunSplitFunctionTask(option.FilePath, functionName, option.UseMockType)
	}
}

func RunSplitFunctionTask(filePath string, functionName string, useMockType int) {
	defer func() {
		// clean all generated files to avoid conflicts
		lifemanager.Closer.Close()
	}()

	dir := path.Dir(filePath)
	option := atgconstant.Options{
		FilePath:      filePath,
		Level:         1,
		Maxtime:       4,
		MinUnit:       atgconstant.MinUnit,
		Uid:           atghelper.RandStringBytes(10),
		FuncName:      functionName,
		Usage:         atgconstant.SplitFunctionMode,
		DirectoryPath: dir,
		UseMockType:   useMockType,
	}
	// warning :not delete println,plugin get necessary msg
	fmt.Println("plugin sdk use UID is:", option.Uid)
	err := WorkForPlugin(option)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = FixGoFile(dir)
	if err != nil {
		fmt.Println(err)
	}
	err = GenerateTestForPlugin(option)
	if err != nil {
		fmt.Printf("Sorry, we cannot generate the test for function: %s, the error is %v\n", functionName, err.Error())
		fmt.Println("splitfunction_failure(1)-r")
	} else {
		fmt.Println("splitfunction_success(1)-r")
	}
}

func UpdateSmartUnit(path string) error {
	var stdBuffer, stdErrBuff bytes.Buffer
	cmd := exec.Command(atgconstant.GoDirective, "mod", "download")
	cmd.Dir = path
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdErrBuff
	err := cmd.Run()
	if err != nil {
		return logextractor.GenSuggestByErrLog(err, stdBuffer.String(), stdErrBuff.String())
	}
	return nil
}
