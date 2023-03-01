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
package mock

import (
	"context"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/codebuilder/variablecard"
	"reflect"
	"strings"
)

type StatementRender struct {
	MockStatement   []string
	TestFuncCall    []string
	AssertStatement []string
	Imports         []string
	MonkeyOutputMap variablecard.MonkeyOutputMap
	UsedMockFunc    map[string]int
}

func OverPassMakeCall(ctx context.Context, funcName string, mockRender *StatementRender, m interface{}) (fun interface{}) {
	var out []reflect.Value
	function := reflect.TypeOf(m)
	if function.NumIn() == 0 {
		return nil
	}
	// Overpass get its input as our target function.
	function = function.In(0)
	if function.Kind() != reflect.Func {
		panic("it's no true")
	}
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	ctx = contexthelper.SetBuilderVector(ctx, "overpass")
	for i := 0; i < function.NumOut(); i++ {
		// 参数是 interface
		// function
		// 别名
		// 不行就放过，没有statement
		resultType := function.Out(i)
		switch resultType.Kind() {
		// reflect.Zero new the value of addressable nor settable
		// it make faker.GetValue(for pointer) fail
		// so we to Mutate struct for pointer by reflectValue.Addr()
		// see: https://halfrost.com/go_reflection/
		case reflect.Ptr:
			v := reflect.New(resultType.Elem()).Elem()
			mutationElem := variablecard.VariableMutate(ctx, resultType.Elem(), v)
			if mutationElem.CanAddr() {
				out = append(out, variablecard.VariableMutate(ctx, resultType.Elem(), v).Addr())
			} else {
				var vNil interface{} = nil
				out = append(out, reflect.ValueOf(vNil))
			}

		default:
			v := reflect.Zero(resultType)
			out = append(out, variablecard.VariableMutate(ctx, resultType, v))
		}
	}
	newFunc := reflect.MakeFunc(function, func(args []reflect.Value) (results []reflect.Value) {
		return out
	})
	var card []string
	for _, r := range out {
		card = append(card, variablecard.ValueToString(ctx, r))
	}
	mockStateMent := ""
	mockStateMent = fmt.Sprintf("%s(%s{\nreturn %s \n})",
		funcName,
		variablecard.ValueToString(ctx, newFunc),
		strings.Join(card, ", "))
	mockRender.MockStatement = append(mockRender.MockStatement, mockStateMent)
	return newFunc.Interface()
}

// TODO（siwei.wang）: alias will make MakeCall panic.
func MakeCall(ctx context.Context, funcName string, mockRender *StatementRender, m interface{}, useMockType int) (fun interface{}) {
	var out []reflect.Value
	function := reflect.TypeOf(m)
	if function.Kind() != reflect.Func {
		panic("it's no true")
	}
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	for i := 0; i < function.NumOut(); i++ {
		// 参数是 interface
		// function
		// 别名
		// 不行就放过，没有statement
		resultType := function.Out(i)
		switch resultType.Kind() {
		// reflect.Zero new the value of addressable nor settable
		// it make faker.GetValue(for pointer) fail
		// so we to Mutate struct for pointer by reflectValue.Addr()
		// see: https://halfrost.com/go_reflection/
		case reflect.Ptr:
			v := reflect.New(resultType.Elem()).Elem()
			mutationElement := variablecard.VariableMutate(ctx, resultType.Elem(), v)
			if mutationElement.CanAddr() {
				out = append(out, variablecard.VariableMutate(ctx, resultType.Elem(), v).Addr())
			} else {
				var vNil interface{} = nil
				out = append(out, reflect.ValueOf(vNil))
			}

		default:
			v := reflect.Zero(resultType)
			out = append(out, variablecard.VariableMutate(ctx, resultType, v))
		}
	}
	newFunc := reflect.MakeFunc(function, func(args []reflect.Value) (results []reflect.Value) {
		return out
	})
	ctx = contexthelper.SetVariableContext(ctx, atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false})
	var card []string
	for _, r := range out {
		card = append(card, variablecard.ValueToString(ctx, r))
	}
	for _, code := range card {
		if code == "unexport variable" {
			return newFunc.Interface()
		}
	}
	switch useMockType {
	case atgconstant.UseMockitoMock:
		mockStateMent := fmt.Sprintf("mockito.Mock(%s).Return(%s).Build()", funcName, strings.Join(card, ", "))
		mockRender.MockStatement = append(mockRender.MockStatement, mockStateMent)
	case atgconstant.UseGoMonkeyMock:
		mockRender.MonkeyOutputMap[funcName] = strings.Join(card, ", ")
		if len(card) != 0 {
			_, ok := mockRender.UsedMockFunc[funcName]
			if !ok {
				mockRender.UsedMockFunc[funcName] = len(card)
			}
		} else {
			if len(out) == 0 {
				_, ok := mockRender.UsedMockFunc[funcName]
				if !ok {
					mockRender.UsedMockFunc[funcName] = 0
				}
			}
		}
	}

	return newFunc.Interface()
}

func GetPrivateFunc(m interface{}) (r interface{}) {
	f := reflect.TypeOf(m)
	// it's no receiver func
	if f.NumIn() <= 0 {
		return m
	}
	var hit bool
	var name string
	for i := 0; i < f.In(0).NumMethod(); i++ {
		if f == f.In(0).Method(i).Type {
			name = f.In(0).Method(i).Name
			hit = true
			break
		}
	}

	if !hit {
		return m
	}
	// it means m it the receiver of
	method, ok := getNestedMethod(f.In(0), name)
	if ok {
		return method.Func.Interface()
	}
	return m
}

func GetNestedMethod(instance interface{}, methodName string) (interface{}, bool) {
	if typ := reflect.TypeOf(instance); typ != nil {
		if m, ok := getNestedMethod(typ, methodName); ok {
			return m.Func.Interface(), true
		}
	}
	return nil, false
}

func getNestedMethod(typ reflect.Type, methodName string) (reflect.Method, bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			if !typ.Field(i).Anonymous {
				// there is no need to acquire non-anonymous method
				continue
			}
			if m, ok := getNestedMethod(typ.Field(i).Type, methodName); ok {
				return m, true
			}
		}
	}
	// a struct receiver is prior to the corresponding pointer receiver
	if m, ok := typ.MethodByName(methodName); ok {
		return m, true
	}
	return reflect.PtrTo(typ).MethodByName(methodName)
}
