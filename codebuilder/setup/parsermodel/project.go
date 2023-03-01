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
package parsermodel

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

type Imethod struct {
	I  *types.Interface
	Id string
}

type ProjectProgram struct {
	PkgPath       string
	PkgName       string
	Prog          *ssa.Program
	Pkgs          []*ssa.Package
	AllFuncs      map[*ssa.Function]bool
	MethodsByName map[string][]*ssa.Function
	MethodsMemo   map[Imethod][]*ssa.Function
}

type ProjectFunction struct {
	ModuleName string

	Function *ssa.Function
	Program  *ProjectProgram

	CalleeFunctionsForTargetFunction map[string]*ssa.Function
	PackageConstants                 []types.Object
}

func (prog *ProjectProgram) LookupMethods(I *types.Interface, m *types.Func) []*ssa.Function {

	if prog.MethodsMemo == nil {
		prog.MethodsMemo = make(map[Imethod][]*ssa.Function)
	}

	id := m.Id()
	methods, ok := prog.MethodsMemo[Imethod{I, id}]
	if !ok {
		for _, f := range prog.MethodsByName[m.Name()] {
			if f.Signature.Recv() == nil {
				continue
			}
			C := f.Signature.Recv().Type() // named or *named
			if types.Implements(C, I) {
				methods = append(methods, f)
			}
		}
		prog.MethodsMemo[Imethod{I, id}] = methods
	}
	return methods
}
