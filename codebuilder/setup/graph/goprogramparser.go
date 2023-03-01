// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package graph

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bytedance/nxt_unit/manager/logextractor"

	"github.com/bytedance/nxt_unit/atghelper"

	"go/types"

	"github.com/bytedance/nxt_unit/atgconstant"

	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var RepoCalls = &sync.Map{}

// TODO: optimize this function
func ParsePackage(file string) (program *parsermodel.ProjectProgram, err error) {
	fileDir := path.Dir(file)
	projectStatus, ok := RepoCalls.Load(fileDir)
	if ok {
		return projectStatus.(*parsermodel.ProjectProgram), nil
	} else {
		defer func() {
			RepoCalls.Store(fileDir, program)
		}()
	}
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	cfg.Tests = false
	cfg.Dir = fileDir
	cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		// 如果是测试文件，解析忽略
		if strings.HasSuffix(filename, "_test.go") {
			return nil, nil
		}
		const mode = parser.DeclarationErrors
		return parser.ParseFile(fset, filename, src, mode)
	}

	initial, err := packages.Load(cfg, fileDir)
	if err != nil {
		return nil, fmt.Errorf("The error belongs to %w\n, the details is: %v", logextractor.ParseProgramError, err.Error())
	}

	// Create SSA packages for all well-typed packages.
	prog, pkgs := ssautil.AllPackages(initial, 0)
	prog.Build()

	allFuncs := ssautil.AllFunctions(prog)
	methodsByName := make(map[string][]*ssa.Function)
	for f := range allFuncs {
		if f.Signature.Recv() != nil {
			methodsByName[f.Name()] = append(methodsByName[f.Name()], f)
		}
	}

	// It means that we cannot find the method,
	// It might be this sectarian:
	// type mMT struct{}
	//
	// func (*mMT) GetSmartUnit(a int) string {
	//	fmt.Println("nono")
	//	return "ok"
	// }
	// We cannot find the smart unit because it is belongs to a lowercase receiver
	if len(pkgs) != 0 && pkgs[0] != nil {
		for _, m := range pkgs[0].Members {
			methods := getAllMethods(pkgs[0].Prog, m.Type())
			methods = append(methods, getAllMethods(pkgs[0].Prog, types.NewPointer(m.Type()))...)
			for _, method := range methods {
				if method.Recv() != nil {
					fn := pkgs[0].Prog.MethodValue(method)
					if fn == nil {
						continue
					}
					methodsByName[fn.Name()] = append(methodsByName[fn.Name()], fn)
				}
			}
		}
	}

	return &parsermodel.ProjectProgram{
		PkgName:       initial[0].Name,
		PkgPath:       initial[0].PkgPath,
		MethodsByName: methodsByName,
		Prog:          prog,
		Pkgs:          pkgs,
		AllFuncs:      allFuncs,
	}, nil
}

// Get all methods of a type.
func getAllMethods(prog *ssa.Program, typ types.Type) []*types.Selection {
	ms := prog.MethodSets.MethodSet(typ)
	methods := make([]*types.Selection, ms.Len())
	for i := 0; i < ms.Len(); i++ {
		methods[i] = ms.At(i)
	}
	return methods
}

// GetFunctionCallGraph TODO: GetFunctionCallGraph this function was called many times by testcode.go, parser.go.
// Need to reduce this functions call. Because it can save our running time
func GetFunctionCallGraph(option atgconstant.Options) (*parsermodel.ProjectFunction, *callgraph.Graph, error) {
	TestedFunctionAndCallees, err := GetPackageFunction(option)
	if err != nil {
		return TestedFunctionAndCallees, nil, err
	}

	if TestedFunctionAndCallees == nil {
		return TestedFunctionAndCallees, nil, fmt.Errorf("callees of Function not found")
	}

	return TestedFunctionAndCallees, BuildCallGraph(TestedFunctionAndCallees, option.FuncName), nil
}

// GetPackageFunction gets all functions from the package, the tested function, the constant in the packages
func GetPackageFunction(option atgconstant.Options) (*parsermodel.ProjectFunction, error) {
	program, err := ParsePackage(option.FilePath)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, fmt.Errorf("program is nil")
	}

	var projectFunction *parsermodel.ProjectFunction
	if option.ReceiverName == "" {
		for function := range program.AllFuncs {
			funcName := function.Name()
			if function.Pkg == nil {
				continue
			}
			funcFile := GetFuncFileName(function, program.Prog)
			if !isSameFile(funcFile, option.FilePath) {
				continue
			}

			if option.FuncName != funcName || strings.Contains(funcName, "$") {
				continue
			}

			// if option.ReceiverName != "" {
			//	functionReceive := atghelper.GetTheReceiveNameFromSSA(function.Signature.Recv().String())
			//	if functionReceive != option.ReceiverName {
			//		continue
			//	}
			// }

			projectFunction = &parsermodel.ProjectFunction{
				Function:                         function,
				Program:                          program,
				CalleeFunctionsForTargetFunction: map[string]*ssa.Function{},
			}
			break
		}
	}

	if option.ReceiverName != "" {
		var find bool
		for _, functions := range program.MethodsByName {
			find = false
			for _, function := range functions {
				funcName := function.Name()
				if function.Pkg == nil {
					continue
				}
				funcFile := GetFuncFileName(function, program.Prog)
				if !isSameFile(funcFile, option.FilePath) {
					continue
				}

				if option.FuncName != funcName || strings.Contains(funcName, "$") {
					continue
				}

				functionReceive := atghelper.GetTheReceiveNameFromSSA(function.Signature.Recv().String())
				if functionReceive != option.ReceiverName {
					continue
				}

				projectFunction = &parsermodel.ProjectFunction{
					Function:                         function,
					Program:                          program,
					CalleeFunctionsForTargetFunction: map[string]*ssa.Function{},
				}
				find = true
				break
			}
			if find {
				break
			}
		}
	}

	if projectFunction == nil {
		return projectFunction, fmt.Errorf("SSA Function Not found")
	}

	// Get all constants
	pkgs := program.Pkgs
	constants := make([]types.Object, 0)
	for _, pkg := range pkgs {
		if pkg == nil {
			continue
		}
		for _, member := range pkg.Members {
			if !strings.HasPrefix(member.Type().Underlying().String(), "func") && member.Object() != nil {
				constants = append(constants, member.Object())
			}
		}
	}
	projectFunction.PackageConstants = constants
	return projectFunction, nil
}

func GetFuncFileName(f *ssa.Function, prog *ssa.Program) string {
	filePos := prog.Fset.Position(f.Pos())
	return filePos.Filename
}

func isSameFile(a, b string) bool {
	a, err := filepath.Abs(a)
	if err != nil {
		return false
	}
	b, err = filepath.Abs(b)
	if err != nil {
		return false
	}
	return a == b
}
