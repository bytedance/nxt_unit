package staticcase

import (
	"context"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytedance/nxt_unit/manager/lifemanager"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/codebuilder/instrumentation"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	atgconstant.PkgRelativePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template")
	os.Exit(exitVal)
}

func TestBackStage_RunFunction(t *testing.T) {
	convey.Convey("BackStage_RunFunction", t, func() {
		patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
		defer patch.Reset()
		patch = gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		patch = gomonkey.ApplyFuncReturn(os.Rename, nil)
		defer func() {
			lifemanager.Closer.Close()
		}()
		ctx := context.Background()
		opt, _ := contexthelper.GetOption(ctx)
		opt.MinUnit = "function"
		opt.FuncName = "Consume"
		opt.ReceiverName = "SpecialConsumption"
		opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/receiver.go")

		ctx, err := getContext(opt)
		if err != nil {
			t.Fatal(err)
		}
		tb := instrumentation.NewFunctionBuilder(ctx)
		// build 创建插桩文件
		_, err = tb.Build(ctx)
		if err != nil {
			t.Fatal(err)
		}
		te := NewTestsuiteEntry(ctx, opt.FilePath, atghelper.GlobalFileLock)
		_, err = te.Build(ctx)
		if err != nil {
			t.Fatal(err)
		}
		workDir := filepath.Dir(opt.FilePath)
		err = WorkToChangeGo(workDir, true)
		if err != nil {
			t.Fatal(err)
		}
		err = GenerateTest(workDir, true)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGenerateBaseTest(t *testing.T) {
	convey.Convey("GenerateBaseTest", t, func() {
		patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
		defer patch.Reset()
		patch1 := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch1.Reset()
		patch3 := gomonkey.ApplyFuncReturn(os.Rename, nil)
		defer patch3.Reset()
		ctx := contexthelper.GetTestContext()
		opt, _ := contexthelper.GetOption(ctx)
		opt.MinUnit = atgconstant.MinUnit
		opt.DirectoryPath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template")
		ctx, err := getContext(opt)
		if err != nil {
			t.Fatal(err)
		}
		err = GenerateBaseTest(ctx)
		if err != nil {
			t.Log(err)
		}
	})
}

func TestConstructorRun(t *testing.T) {
	convey.Convey("ConstructorRun", t, func() {
		patch1 := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch1.Reset()
		defer func() {
			lifemanager.Closer.Close()
		}()
		opt := atgconstant.Options{
			Level:         1,
			Maxtime:       4,
			GenerateType:  atgconstant.GAMode,
			MinUnit:       "function",
			FilePath:      path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/constructer_callers.go"),
			DirectoryPath: path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template"),
			FuncName:      "PrintName",
			Uid:           atghelper.RandStringBytes(5),
			Usage:         "plugin",
		}

		ctx, err := getContext(opt)
		if err != nil {
			t.Fatal(err)
		}
		tb := instrumentation.NewFunctionBuilder(ctx)
		// build 创建插桩文件
		_, err = tb.Build(ctx)
		if err != nil {
			t.Fatal(err)
		}
		te := NewTestsuiteEntry(ctx, opt.FilePath, atghelper.GlobalFileLock)
		_, err = te.Build(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "no need to geneate test") {
				t.Skip("file generated")
			}
			t.Fatal(err)
		}
		workDir := filepath.Dir(opt.FilePath)
		err = WorkToChangeGo(workDir, opt.DebugMode)
		if err != nil {
			t.Fatal(err)
		}
		_ = GenerateTestForPlugin(opt)
	})
}

func TestGenerateBaseTest_Function(t *testing.T) {
	convey.Convey("ConstructorRun", t, func() {
		patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
		defer patch.Reset()
		patch = gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		patch = gomonkey.ApplyFuncReturn(os.Rename, nil)
		defer patch.Reset()
		defer func() {
			lifemanager.Closer.Close()
		}()
		ctx := contexthelper.GetTestContext()
		opt, _ := contexthelper.GetOption(ctx)
		opt.DirectoryPath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template")
		opt.MinUnit = atgconstant.MinUnit
		ctx, err := getContext(opt)
		if err != nil {
			t.Fatal(err)
		}
		err = GenerateBaseTest(ctx)
		if err != nil {
			t.Log(err)
		}
	})
}

func TestGetDataStreamBuilder(t *testing.T) {
	opt := atgconstant.Options{
		FilePath:  path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/dataanalysis.go"),
		MinUnit:   "function",
		DebugMode: true,
		Usage:     "plugin",
		FuncName:  "QueryData",
	}
	opt.DirectoryPath = path.Dir(opt.FilePath)
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	mockStatements, builder := GetGlobalValueBuilder(ctx)
	t.Log(mockStatements, builder)
}

func TestGenerateBaseTest_Merge2(t *testing.T) {
	defer func() {
		lifemanager.Closer.Close()
	}()
	patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
	defer patch.Reset()
	patch1 := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
	defer patch1.Reset()
	patch2 := gomonkey.ApplyFuncReturn(os.Rename, nil)
	defer patch2.Reset()
	ctx := context.Background()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = atgconstant.MinUnit
	opt.MinUnit = "function"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/receiver.go")
	opt.DirectoryPath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template")
	opt.FuncName = "GetSmartUnit"
	opt.ReceiverName = "*mMT"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateBaseTest(ctx)
	if err != nil {
		t.Log(err)
	}
}

func TestGenerateBaseTest_Merge(t *testing.T) {
	defer func() {
		lifemanager.Closer.Close()
	}()
	patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
	defer patch.Reset()
	patch1 := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
	defer patch1.Reset()
	patch2 := gomonkey.ApplyFuncReturn(os.Rename, nil)
	defer patch2.Reset()
	ctx := context.Background()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = atgconstant.MinUnit
	opt.MinUnit = "function"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/receiver.go")
	opt.DirectoryPath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template")
	opt.FuncName = "CheckBigStruct"
	opt.ReceiverName = "*AdBigStruct"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateBaseTest(ctx)
	if err != nil {
		t.Log(err)
	}
}

func TestGenerateBaseTest_ReNameFunction(t *testing.T) {
	defer func() {
		lifemanager.Closer.Close()
	}()
	patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
	defer patch.Reset()
	patch1 := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
	defer patch1.Reset()
	patch2 := gomonkey.ApplyFuncReturn(os.Rename, nil)
	defer patch2.Reset()
	ctx := context.Background()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = atgconstant.MinUnit
	opt.MinUnit = "function"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/receiver.go")
	opt.DirectoryPath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template")
	opt.FuncName = "GetAbParamsMap"
	// opt.ReceiverName = "TikTokConsumption"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateBaseTest(ctx)
	if err != nil {
		t.Log(err)
	}
}
func TestGetAllMock(t *testing.T) {
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.DirectoryPath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template")
	opt.MinUnit = "function"
	opt.FuncName = "FunctionReturnError"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	mockStatement := GetAllMock(ctx, 2)
	if len(mockStatement) <= 0 {
		t.Fatal("mock function which return error failed")
	}
	t.Log(mockStatement)
}

func TestGetGlobalValueBuilder(t *testing.T) {
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = "function"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/mockpkg/dao/query.go")
	opt.DirectoryPath = path.Dir(opt.FilePath)
	opt.FuncName = "CallGlobalInterface"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	mockStatements, builder := GetGlobalValueBuilder(ctx)
	t.Log(mockStatements, builder)
}

func TestGetGlobalPointerBuilder(t *testing.T) {
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = "function"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/mockpkg/dao/query.go")
	opt.DirectoryPath = path.Dir(opt.FilePath)
	opt.FuncName = "CallGlobalPointer"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	mockStatements, builder := GetGlobalValueBuilder(ctx)
	t.Log(mockStatements, builder)
}

func TestPluginSDK_GlobalValue(t *testing.T) {
	convey.Convey("GlobalValue", t, func() {
		patch := gomonkey.ApplyFuncReturn(FixGoFile, nil)
		defer patch.Reset()
		patch = gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		patch = gomonkey.ApplyFuncReturn(os.Rename, nil)
		defer patch.Reset()
		defer func() {
			lifemanager.Closer.Close()
		}()
		ctx := contexthelper.GetTestContext()
		opt, _ := contexthelper.GetOption(ctx)
		opt.MinUnit = "file"
		opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/mockpkg/dao/query.go")
		opt.DirectoryPath = path.Dir(opt.FilePath)
		opt.FuncName = "CallGlobalInterface"
		opt.Uid = atghelper.RandStringBytes(5)
		opt.Usage = "plugin"
		ctx, err := getContext(opt)
		if err != nil {
			t.Fatal(err)
		}
		tb := instrumentation.NewFunctionBuilder(ctx)
		// build 创建插桩文件
		_, err = tb.Build(ctx)
		if err != nil {
			t.Fatal(err)
		}
		te := NewTestsuiteEntry(ctx, opt.FilePath, atghelper.GlobalFileLock)
		_, err = te.Build(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "no need to geneate test") {
				t.Skip("file generated")
			}
			t.Fatal(err)
		}
		workDir := filepath.Dir(opt.FilePath)
		err = WorkToChangeGo(workDir, opt.DebugMode)
		if err != nil {
			t.Fatal(err)
		}
		_ = GenerateTestForPlugin(opt)
	})
}

func TestGetPickStructFieldBuilder(t *testing.T) {
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = "function"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/dataanalysis.go")
	opt.DirectoryPath = path.Dir(opt.FilePath)
	opt.FuncName = "QueryData"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	builder := PickStructField(ctx)
	for _, code := range builder {
		t.Log(code)
	}
}

func TestGetPickStructFieldBuilderss(t *testing.T) {
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = "file"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/mockpkg/dao/query.go")
	opt.DirectoryPath = path.Dir(opt.FilePath)
	opt.FuncName = "CallGlobalInterface"
	opt.Uid = atghelper.RandStringBytes(5)
	opt.Usage = "plugin"
	ctx, err := getContext(opt)
	if err != nil {
		t.Fatal(err)
	}
	builder := PickStructField(ctx)
	for _, code := range builder {
		fmt.Println(code)
	}
}
