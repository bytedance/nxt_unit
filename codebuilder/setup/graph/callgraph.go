// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package graph

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"

	"github.com/bytedance/nxt_unit/atgconstant"

	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

func BuildCallGraph(function *parsermodel.ProjectFunction, functionName string) *callgraph.Graph {
	cg := &callgraph.Graph{Nodes: make(map[*ssa.Function]*callgraph.Node)}
	cg.Root = cg.CreateNode(function.Function)
	createGraph(cg.Root, cg, function)
	// buildNode(cg.Root, cg, function)
	return cg
}

func createGraph(parentNode *callgraph.Node, cg *callgraph.Graph, projectFunction *parsermodel.ProjectFunction) *callgraph.Node {
	fNode := parentNode
	ScanFuncBlocks(parentNode.Func.Blocks, projectFunction, fNode, cg)
	return fNode
}

func getFunction(call *ssa.CallCommon, function *parsermodel.ProjectFunction) *ssa.Function {
	tiface := call.Value.Type().Underlying().(*types.Interface)
	methods := function.Program.LookupMethods(tiface, call.Method)

	var f *ssa.Function
	for _, method := range methods {
		if method.Name() == call.Method.Name() {
			f = method
		}
	}
	return f
}

func funcInStd(f *ssa.Function) bool {
	if f == nil || f.Pkg == nil || f.Pkg.Pkg == nil || f.Pkg.Pkg.Path() == "" {
		return true
	}
	funcPath := f.Pkg.Pkg.Path()
	if strings.HasPrefix(funcPath, "golang.org/x") {
		return true
	}

	if duplicatepackagemanager.IsPathShouldBeRemoved(funcPath) {
		return true
	}

	for name := range atgconstant.IgnoredPath {
		if strings.HasPrefix(funcPath, name) {
			return true
		}
	}

	// TODO: No idea why we need this
	path := joinPath(atgconstant.GOROOT, "src", funcPath)
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// hasCycleAndAdded checks whether the function is in the map
func hasCycleAndAdded(function *ssa.Function, projectFunction *parsermodel.ProjectFunction) bool {
	if function == nil || function.Pkg == nil {
		return true
	}
	funcFullName := fmt.Sprintf("%v/%v", function.Pkg.Pkg.Path(), function.Name())
	if _, ok := projectFunction.CalleeFunctionsForTargetFunction[funcFullName]; ok {
		return true
	}
	projectFunction.CalleeFunctionsForTargetFunction[funcFullName] = function
	return false
}

func joinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func ScanFuncBlocks(bacicBlocks []*ssa.BasicBlock, projectFunction *parsermodel.ProjectFunction, parentNode *callgraph.Node, cg *callgraph.Graph) {
	for _, b := range bacicBlocks {
		for _, instr := range b.Instrs {
			var gNode *callgraph.Node

			site, ok := instr.(ssa.CallInstruction)
			if !ok {
				continue
			}
			callCommon := site.Common()
			if callCommon == nil {
				continue
			}
			var downStreamFunc *ssa.Function
			if callCommon.IsInvoke() {
				invokeF := getFunction(callCommon, projectFunction)
				if invokeF != nil {
					downStreamFunc = invokeF
				}
			} else {
				calleeF := callCommon.StaticCallee()
				if calleeF != nil {
					downStreamFunc = calleeF
				}
			}
			if downStreamFunc == nil {
				continue
			}

			_, isDefer := site.(*ssa.Defer)
			_, isGoRoutine := site.(*ssa.Go)
			if strings.Contains(downStreamFunc.Name(), "$") && (isDefer || isGoRoutine) {
				// make callgraph: add downstream method to callGraph eg (*ssa.Go)(*ssa.Defer)
				// eg: go routine(*ssa.Defer)  eg:/atg/template/atg.go,funcName:ComplexDeferFunction
				// eg: go routine(*ssa.Go)  eg:/atg/template/atg.go,funcName:ComplexGoRoutineFunction
				// case use to get defer function downstream method
				// maybe missing corner cases ,waiting to supplement
				switch commonValue := callCommon.Value.(type) {
				case *ssa.Function:
					if commonValue != nil && commonValue.Blocks != nil && len(commonValue.Blocks) > 0 {
						ScanFuncBlocks(commonValue.Blocks, projectFunction, parentNode, cg)
					}
				case *ssa.MakeClosure:
					closureFunc, ok := commonValue.Fn.(*ssa.Function)
					if ok {
						if closureFunc != nil && closureFunc.Blocks != nil && len(closureFunc.Blocks) > 0 {
							ScanFuncBlocks(closureFunc.Blocks, projectFunction, parentNode, cg)
						}
					}
				}
			} else {
				gNode = createCalleeNode(downStreamFunc, cg, projectFunction)
			}

			if gNode == nil {
				continue
			}
			if atgconstant.GraphLevel == 1 {
				callgraph.AddEdge(parentNode, site, gNode)
				continue
			} else {
				callgraph.AddEdge(gNode, site, createGraph(gNode, cg, projectFunction))
			}
		}
	}
}

func createCalleeNode(downStreamFunc *ssa.Function, cg *callgraph.Graph, projectFunction *parsermodel.ProjectFunction) *callgraph.Node {
	if downStreamFunc == nil {
		return nil
	}

	if funcInStd(downStreamFunc) {
		return nil
	}

	if !strings.Contains(downStreamFunc.Name(), "$") {
		if hasCycleAndAdded(downStreamFunc, projectFunction) {
			return nil
		}
		return cg.CreateNode(downStreamFunc)
	}
	return nil
}
