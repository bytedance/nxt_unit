package setup

import (
	"fmt"
	"go/types"
	"path"
	"testing"

	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func TestGetAllConstructor(t *testing.T) {
	fileDir := path.Join(baseDir, "../../../atg/template/constructors/constructors.go")
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	cfg.Tests = false

	initial, err := packages.Load(cfg, fileDir)
	if err != nil {
		t.Fatal(err)
		return
	}
	type args struct {
		pkg *packages.Package
	}
	tests := []struct {
		name               string
		args               args
		wantConstructorMap map[types.Type][]*ssa.Function
	}{
		{
			name: "simple",
			args: args{
				pkg: initial[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConstructorMap := GetConstructors(tt.args.pkg)
			if len(gotConstructorMap) != 2 {
				t.Errorf("want len 2, got len %v", len(gotConstructorMap))
			}
		})
	}
}

func TestGetConstructorsByFunc(t *testing.T) {
	fileDir := path.Join(baseDir, "../../../atg/template/constructer_callers.go")
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	cfg.Tests = false

	initial, err := packages.Load(cfg, fileDir)
	if err != nil {
		t.Fatal(err)
		return
	}
	prog, _ := ssautil.AllPackages(initial, 0)
	prog.Build()

	functions := ssautil.AllFunctions(prog)
	var function Functions
	for f := range functions {
		if f.Name() == "PrintName" {
			function.TestFunction = new(parsermodel.ProjectFunction)
			function.TestFunction.Function = f
		}
	}
	type args struct {
		function Functions
	}
	tests := []struct {
		name               string
		args               args
		wantConstructorMap map[types.Type][]*ssa.Function
	}{
		{
			name: "simple",
			args: args{
				function: function,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConstructorMap := GetConstructorsByFunc(tt.args.function)
			fmt.Println(gotConstructorMap)
		})
	}
}
