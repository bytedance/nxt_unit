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
package setup

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/codebuilder/setup/graph"
	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
	"golang.org/x/tools/go/callgraph"
)

// 函数描述
type FuncDesc struct {
	Path    string
	Package string
	Name    string
}

// 描述一个函数调用N个函数的一对多关系
type CallerRelation struct {
	Caller   FuncDesc
	Callees  []FuncDesc
	distance int
}

type CallGraph struct {
	RootFunc          *FuncDesc
	CallerRelationMap map[string]*CallerRelation
}

// 描述关键函数的一条反向调用关系
type CalleeRelation struct {
	Callee []FuncDesc
	CanFix bool // 该调用关系能反向找到gin.Context即可以自动修复
}

// Todo: init 在一个包或者一个项目里可能有多个，这里我们先不去管它，在之后要加上。
type Functions struct {
	TestFunction      *parsermodel.ProjectFunction
	DateSteam         map[interface{}]struct{}
	TestFuncCallGraph *callgraph.Graph
	SRCFile           string
	FunctionName      string
}

// GetFunctions gets the functions from tested function.
func GetFunctions(option atgconstant.Options) (Functions, error) {
	emptyFunction := Functions{}
	function, callGraph, err := graph.GetFunctionCallGraph(option)
	if err != nil {
		return emptyFunction, err
	}

	// analysis the date stream of testFunction
	// dataReference := map[interface{}]struct{}{}
	// for _, block := range function.Function.Blocks {
	//	for _, instr := range block.Instrs {
	//		GetTypeByAnalysisBlock(instr, dataReference)
	//	}
	// }
	dataReference := make(map[interface{}]struct{}, 0)
	return Functions{
		TestFunction:      function,
		TestFuncCallGraph: callGraph,
		DateSteam:         dataReference,
		SRCFile:           option.FilePath,
		FunctionName:      option.FuncName,
	}, nil
}

func GetTypeByAnalysisBlock(instr ssa.Instruction, reference map[interface{}]struct{}) {
	if reference == nil {
		return
	}
	typeList := []interface{}{}
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	switch n := instr.(type) {
	// The Field instruction yields the Field of struct X.
	//
	// The field is identified by its index within the field list of the
	// struct type of X; by using numeric indices we avoid ambiguity of
	// package-local identifiers and permit compact representations.
	//
	// Pos() returns the position of the ast.SelectorExpr.Sel for the
	// field, if explicit in the source.
	//
	// Example printed form:
	// 	t1 = t0.name [#1]
	case *ssa.FieldAddr:
		fieldType := n.X.Type().Underlying().(*types.Pointer).Elem().Underlying().(*types.Struct).Field(n.Field)
		structs := n.X.Type().Underlying().(*types.Pointer).Elem().(*types.Named).Obj()
		varUnderlying := strings.Join([]string{structs.Pkg().Path(), structs.Name(), fieldType.Name()}, ".")
		typeList = append(typeList, varUnderlying)
	case *ssa.Field:
		fieldType := n.X.Type().(*types.Struct).Field(n.Field)
		structs := n.X.Type().Underlying().(*types.Pointer).Elem().(*types.Named).Obj()
		varUnderlying := strings.Join([]string{structs.Pkg().Path(), structs.Name(), fieldType.Name()}, ".")
		typeList = append(typeList, varUnderlying)
	// like t0 > 9
	// it could set mode info to mutate t0 like t0 :=8.9 ~ 9.1
	case *ssa.BinOp:
		xt := n.X.Type()
		yt := n.Y.Type()
		typeList = append(typeList, xt.Underlying().String(), yt.Underlying().String())
	// like query(a,b,c)
	// (todo: liuguan.cheng ):  now we just pick structField more type is need to support
	case *ssa.Store:

	case *ssa.Call:
		// (todo: liuguan.cheng ):  we can get callee by analysis ssa
		// callee is n.Call.Method.Name()
		for _, arg := range n.Call.Args {
			typeList = append(typeList, arg.Type().Underlying().String())
		}
	case *ssa.Alloc:
		// if it's no alloc by heap, it means alloc local value but no param from function's arg
		// Alloc values are always addresses, and have pointer types, so the
		// type of the allocated variable is actually
		// Type().Underlying().(*types.Pointer).Elem().
		if n.Heap {
			typeList = append(typeList, n.Type().Underlying().(*types.Pointer).Elem().Underlying().String())
		}
	case *ssa.Return:
		for _, arg := range n.Results {
			typeList = append(typeList, arg.Type().Underlying().String())
		}
	case *ssa.MakeMap:
		typeList = append(typeList, n.Type().Underlying().String())
	case *ssa.Slice:
		typeList = append(typeList, n.Type().Underlying().String())
	}
	for _, t := range typeList {
		if _, exist := reference[t]; !exist {
			reference[t] = struct{}{}
		}
	}
}
