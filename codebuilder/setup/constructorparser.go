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
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Get all constructor for each type by SSA
// We consider the constructor to satisfy such characteristics:
// 1. Return at least one value
// 2. Not a method
// 3. The return value contains a structure we need
func GetConstructors(pkg *packages.Package) (constructorMap map[types.Type][]*ssa.Function) {
	constructorMap = make(map[types.Type][]*ssa.Function)
	prog, _ := ssautil.Packages([]*packages.Package{pkg}, 0)
	allFuncs := ssautil.AllFunctions(prog)
	// find all constructor
	for function := range allFuncs {
		if function.Signature.Recv() != nil {
			continue
		}
		if function.Signature.Results().Len() == 0 {
			continue
		}
		rtnType := function.Signature.Results().At(0).Type()
		if _, ok := rtnType.Underlying().(*types.Struct); !ok {
			continue
		}
		constructorMap[rtnType] = append(constructorMap[rtnType], function)
	}
	return
}

// Get all constructor for certain type by SSA
// We consider the constructor to satisfy such characteristics:
// 1. Return at least one value
// 2. Not a method
// 3. The return value contains a structure we need
// TODO: support pointer
func GetConstructorsByType(pkg *packages.Package, typeList []types.Type) (constructorMap map[string][]*ssa.Function) {
	constructorMap = make(map[string][]*ssa.Function)
	prog, _ := ssautil.Packages([]*packages.Package{pkg}, 0)
	allFuncs := ssautil.AllFunctions(prog)
	// find all constructor
	for function := range allFuncs {
		if function.Signature.Recv() != nil {
			continue
		}
		// only support no param now
		if function.Signature.Params().Len() > 0 {
			continue
		}
		if function.Signature.Results().Len() == 0 {
			continue
		}
		rtnType := function.Signature.Results().At(0).Type()
		for _, t := range typeList {
			// if it is same type
			if rtnType.String() == t.String() {
				constructorMap[rtnType.String()] = append(constructorMap[rtnType.String()], function)
			}
		}
	}
	return
}

func GetConstructorsByFunc(function Functions) (constructorMap map[string][]*ssa.Function) {
	constructorMap = make(map[string][]*ssa.Function)
	pkgTypeMap := make(map[string][]types.Type)
	receiver := function.TestFunction.Function.Signature.Recv()
	if receiver != nil {
		pkgPath, err := getPackagePathofNamed(receiver.Type())
		if err == nil {
			pkgTypeMap[pkgPath] = []types.Type{receiver.Type()}
		}
	}
	params := function.TestFunction.Function.Signature.Params()
	for i := 0; i < params.Len(); i++ {
		pkgPath, err := getPackagePathofNamed(params.At(i).Type())
		if err != nil {
			continue
		}
		if typeList, ok := pkgTypeMap[pkgPath]; ok {
			pkgTypeMap[pkgPath] = append(typeList, params.At(i).Type())
		} else {
			pkgTypeMap[pkgPath] = []types.Type{params.At(i).Type()}
		}
	}
	for pkgPath, typeList := range pkgTypeMap {
		cfg := packages.Config{Mode: packages.LoadSyntax}
		initial, err := packages.Load(&cfg, pkgPath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if len(initial) < 1 {
			continue
		}
		for t, funcList := range GetConstructorsByType(initial[0], typeList) {
			constructorMap[t] = funcList
		}
	}
	return
}

func getPackagePathofNamed(typ types.Type) (string, error) {
	switch t := typ.(type) {
	case *types.Named:
		if t.Obj().Pkg() == nil {
			return "", errors.New("no pakcage")
		}
		return t.Obj().Pkg().Path(), nil
	case *types.Pointer:
		return getPackagePathofNamed(t.Elem())
	default:
		return "", errors.New("no pakcage")
	}
}
