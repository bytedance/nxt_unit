package goparser

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bytedance/nxt_unit/codebuilder/setup/graph"
	"gotest.tools/assert"

	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/staticcase/internal/models"

	"github.com/bytedance/nxt_unit/atgconstant"
)

func TestParser_Parse_Interface(t *testing.T) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../../../atg/template/interface.go")
	p := &Parser{}
	sr, err := p.Parse(filePath, nil)
	if err != nil {
		t.Fatalf("Parser.Parse source file: %v", err)
	}
	for _, field := range sr.Funcs[1].Parameters {
		fmt.Println(field.Type.String())
	}
}

func TestParser_Parse(t *testing.T) {
	opt := atgconstant.Options{}
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/interface.go")
	opt.FuncName = "DeepCopy"
	// check interface from other pkg by SSA
	funcInfoSSA, _, err := graph.GetFunctionCallGraph(opt)
	if err == nil {
		for name, f := range funcInfoSSA.Program.MethodsByName {
			fmt.Println(name, f)
		}
	}
}

func TestParser_Parse_Interface_OtherPkg(t *testing.T) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../../../atg/template/interface.go")
	p := &Parser{}
	sr, err := p.Parse(filePath, nil)
	if err != nil {
		t.Fatalf("Parser.Parse source file: %v", err)
	}
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = atgconstant.MinUnit
	opt.FilePath = filePath
	opt.DirectoryPath = filepath.Dir(opt.FilePath)
	opt.FuncName = "DeepCopy"
	// funcs, err := setup.GetFunctions(opt)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// params := funcs.TestFunction.Function.Params[0].Type()
	var target *models.Function
	for _, f := range sr.Funcs {
		if f.Name == "DeepCopy" {
			target = f
		}
	}
	if target == nil {
		t.Fatal("no function")
	}
	// switch params.Underlying().(type) {
	// case *types.Interface:
	// 	target.Parameters[0].Type.IsInterface = true
	// 	target.Parameters[0].Type.Value = fmt.Sprintf("struct{%v}", target.Parameters[0].Type.Value)
	// }
	for _, field := range target.Parameters {
		fmt.Println(field.Type.String())
	}
}

func TestParser_Parse_Interface_OtherPkg_Reciver(t *testing.T) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../../../atg/template/interface.go")
	p := &Parser{}
	sr, err := p.Parse(filePath, nil)
	if err != nil {
		t.Fatalf("Parser.Parse source file: %v", err)
	}
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = atgconstant.MinUnit
	opt.FilePath = filePath
	opt.DirectoryPath = filepath.Dir(opt.FilePath)
	opt.FuncName = "DeepCopyS"
	var target *models.Function
	for _, f := range sr.Funcs {
		if f.Name == "DeepCopyS" {
			target = f
		}
	}
	if target == nil {
		t.Fatal("no function")
	}
	for _, field := range target.Parameters {
		fmt.Println(field.Type.String())
	}
}

func TestParser_Parse_EmptyInterface(t *testing.T) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../../../atg/template/interface.go")
	p := &Parser{}
	sr, err := p.Parse(filePath, nil)
	if err != nil {
		t.Fatalf("Parser.Parse source file: %v", err)
	}
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = atgconstant.MinUnit
	opt.FilePath = filePath
	opt.DirectoryPath = filepath.Dir(opt.FilePath)
	opt.FuncName = "StructToMap"
	var target *models.Function
	for _, f := range sr.Funcs {
		if f.Name == "StructToMap" {
			target = f
		}
	}
	if target == nil {
		t.Fatal("no function")
	}
	for _, field := range target.Parameters {
		fmt.Println(field.Type.String())
	}
}
