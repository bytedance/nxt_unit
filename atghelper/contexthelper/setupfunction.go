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
package contexthelper

import (
	"context"

	"github.com/bytedance/nxt_unit/codebuilder/setup"
	"golang.org/x/tools/go/ssa"
)

type functionsKey struct {
}

var FunctionsKey = functionsKey{}

type funcMapKey struct {
}

var FuncMapKey = funcMapKey{}

type constructorMapKey struct {
}

var ConstructorMapKey = constructorMapKey{}

func SetSetupFunc(ctx context.Context, st setup.Functions) context.Context {
	return context.WithValue(ctx, FunctionsKey, st)
}

func GetSetupFunc(ctx context.Context) (setup.Functions, bool) {
	value := ctx.Value(FunctionsKey)
	functions, ok := value.(setup.Functions)
	if !ok {
		return functions, false
	}
	return functions, true
}

func SetSetupFuncMap(ctx context.Context, st map[string]setup.Functions) context.Context {
	return context.WithValue(ctx, FuncMapKey, st)
}

func GetSetupFuncMap(ctx context.Context) (map[string]setup.Functions, bool) {
	value := ctx.Value(FuncMapKey)
	functions, ok := value.(map[string]setup.Functions)
	if !ok {
		return functions, false
	}
	return functions, true
}

func SetConstructorFuncMap(ctx context.Context, constructorMap map[string][]*ssa.Function) context.Context {
	value := ctx.Value(ConstructorMapKey)
	if functions, ok := value.(map[string][]*ssa.Function); ok {
		for k, m := range constructorMap {
			functions[k] = m
		}
		return context.WithValue(ctx, ConstructorMapKey, functions)
	}
	return context.WithValue(ctx, ConstructorMapKey, constructorMap)
}

func GetConstructorFuncMap(ctx context.Context) (map[string][]*ssa.Function, bool) {
	value := ctx.Value(ConstructorMapKey)
	functions, ok := value.(map[string][]*ssa.Function)
	if !ok {
		return functions, false
	}
	return functions, true
}
