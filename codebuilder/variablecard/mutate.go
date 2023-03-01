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
package variablecard

import (
	"context"
	"fmt"
	gomonkeyv2 "github.com/agiledragon/gomonkey/v2"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/faker"
	"github.com/bytedance/nxt_unit/smartunitvariablebuild"
	util "github.com/typa01/go-utils"
	"math/rand"
	"reflect"
	"strings"
)

// the context contains the follow information:
// canBeNil: Bool. If the variable is used for the mocked statement input, it could be nil.
// ID: ID tell us the relationship between the variable and statement.
// Level: determine how deeply that we do the recursion.
func VariableMutate(ctx context.Context, t reflect.Type, v reflect.Value) (mutate reflect.Value) {
	origin := v
	defer func() {
		if err := recover(); err != nil {
			// todo temp file manager
			// fmt.Printf("[VariableMutate] has error: %v\n", err)
			mutate = origin
		}
	}()
	if t == nil {
		return v
	}
	vtx, ok := contexthelper.GetVariableContext(ctx)
	if !ok {
		return v
	}
	newV, ok := smartunitvariablebuild.GetSpecialVariableV3(ctx, t)
	if ok {
		return newV
	}
	// TODO(siwei.wang): Assignable might panic, please avoid the panic here
	switch t.Kind() {
	case reflect.Bool:
		candidate := reflect.ValueOf(atghelper.RandomBool(0.5))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Int:
		candidate := reflect.ValueOf(int(v.Int() + int64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Int8:
		candidate := reflect.ValueOf(int8(v.Int() + int64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Int16:
		candidate := reflect.ValueOf(int16(v.Int() + int64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Int32:
		candidate := reflect.ValueOf(int32(v.Int() + int64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Int64:
		candidate := reflect.ValueOf(v.Int() + int64(atghelper.GetRandomFloat()*atgconstant.Delta))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Uint:
		candidate := reflect.ValueOf(uint(v.Uint() + uint64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Uint8:
		candidate := reflect.ValueOf(uint8(v.Uint() + uint64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Uint16:
		candidate := reflect.ValueOf(uint16(v.Uint() + uint64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Uint32:
		candidate := reflect.ValueOf(uint32(v.Uint() + uint64(atghelper.GetRandomFloat()*atgconstant.Delta)))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Uint64:
		candidate := reflect.ValueOf(v.Uint() + uint64(atghelper.GetRandomFloat()*atgconstant.Delta))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Float32:
		candidate := reflect.ValueOf(float32(v.Float() + atghelper.GetRandomFloat()*atgconstant.Delta))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Float64:
		candidate := reflect.ValueOf(v.Float() + atghelper.GetRandomFloat()*atgconstant.Delta)
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Complex64:
		tmp := complex(float32(real(v.Complex())+atghelper.GetRandomFloat()*atgconstant.Delta),
			float32(imag(v.Complex())+atghelper.GetRandomFloat()*atgconstant.Delta))
		candidate := reflect.ValueOf(tmp)
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.Complex128:
		tmp := complex(real(v.Complex())+atghelper.GetRandomFloat()*atgconstant.Delta,
			imag(v.Complex())+atghelper.GetRandomFloat()*atgconstant.Delta)
		candidate := reflect.ValueOf(tmp)
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate
	case reflect.String:
		mLen := rand.Intn(28) + 1
		builder := util.NewStringBuilder()
		original := v.String()
		for i := 0; i < mLen; i++ {
			builder.Append(atghelper.RandStringBytes(1))
		}
		candidate := reflect.ValueOf(fmt.Sprint(original, builder.ToString()))
		if !candidate.Type().AssignableTo(t) {
			candidate = candidate.Convert(t)
		}
		return candidate

	// TODO: Pointer change is hard. Need to implement it.
	// Caution: it is really easy to throw panic. Be careful when you implement the below part.
	case reflect.Ptr:
		if vtx.CanBeNil && atghelper.RandomBool(0.25) {
			var vNil interface{}
			vNil = nil
			return reflect.ValueOf(vNil)
		}
		if t.Elem().Kind() == reflect.Struct {
			_, exist := t.Elem().FieldByName("InterfaceTag")
			if exist {
				for i := 0; i < t.NumMethod(); i++ {
					f := t.Method(i)
					newF := reflect.MakeFunc(f.Type, func(args []reflect.Value) (results []reflect.Value) {
						// fmt.Println(fmt.Sprintf("Interface Call： %s , return defualt value", f.Name))
						for i := 0; i < f.Type.NumOut(); i++ {
							if f.Type.Out(i).Kind() == reflect.Ptr {
								v := reflect.New(f.Type.Out(i).Elem()).Elem()
								mutateElement := VariableMutate(ctx, f.Type.Out(i).Elem(), v)
								if mutateElement.CanAddr() {
									results = append(results, VariableMutate(ctx, f.Type.Out(i).Elem(), v).Addr())
								} else {
									var vNil interface{} = nil
									results = append(results, reflect.ValueOf(vNil))
								}
								continue
							}
							results = append(results, reflect.Zero(f.Type.Out(i)))
						}
						return
					})
					newFuncPatch := gomonkeyv2.ApplyFunc(f.Func.Interface(), newF.Interface())
					defer newFuncPatch.Reset()
				}
			}
		}
		// level means, how deep value we will generate for the struct. currently ,we only support level 3.
		// It means, we support maximum three levels.
		if !vtx.CanBeNil && atghelper.RandomBool(1.0) {
			switch t.Elem().Kind() {
			case reflect.Struct:
				return VariableMutate(ctx, t.Elem(), reflect.Zero(t.Elem())).Addr()
			default:
				newV, err := faker.GetValue(t, vtx.Level)
				if err == nil {
					return newV
				} else {
					if strings.Contains(err.Error(), "unexport") {
						var vNil interface{} = nil
						return reflect.ValueOf(vNil)
					}
				}
			}
		}
		return v
	case reflect.Struct:
		newV := reflect.New(t)
		newV = newV.Elem()
		vtx.Level += 1
		// to avoid circle link of struct
		if vtx.Level >= 10 {
			return reflect.Zero(t)
		}
		finalSet := false
		vtx.CanBeNil = false
		ctx = contexthelper.SetVariableContext(ctx, vtx)
		for i := 0; i < v.NumField(); i++ {
			f := newV.Field(i)
			// it's still no good enough to trim unuseful structField
			// todo: (liuguancheng.xiaohei) need to analysis compound variable
			// picker := ctx.Value("ValuePicker")
			// if picker != nil && !checkTestCaseStruct(t.Name()) {
			// 	fieldMap, ok := picker.(map[string]struct{})
			// 	if ok {
			// 		valueUnderlying := strings.Join([]string{t.PkgPath(), t.Name(), t.Field(i).Name}, ".")
			// 		if _, exist := fieldMap[valueUnderlying]; !exist {
			// 			continue
			// 		}
			// 	}
			//
			// }

			if f.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if f.CanSet() {
					canSet := atghelper.IsTypeExported(f.Type())
					if !canSet {
						break
					}
					// change value of field
					// FIXME: 需要其他类型的字面量列表，目前没存，暂时使用随机生成
					f.Set(VariableMutate(ctx, f.Type(), f).Convert(f.Type()))
					fieldName := strings.ToLower(t.Field(i).Name)
					var found bool
					for k, tagFunc := range faker.MapperTag {
						if found {
							break
						}
						if strings.Contains(fieldName, k) {
							fake, err := tagFunc(f)
							if err == nil {
								SafeSet(f, ctx, found, fake)
							}
						}
					}
					finalSet = true
				}
			}
		}
		if finalSet {
			return newV
		}
		return v
	case reflect.Interface:
		return v
	case reflect.Slice:
		sLen := v.Len()
		if sLen == 0 {
			sLen = rand.Intn(2)
			v = reflect.MakeSlice(t, 0, 0)
			for i := 0; i < sLen; i++ {
				mutatedV := VariableMutate(ctx, t.Elem(), reflect.Zero(t.Elem()))
				if t.Elem().AssignableTo(mutatedV.Type()) {
					v = reflect.Append(v, mutatedV)
				}
			}
		}
		// It will invoke panic, therefore, we add this rule.
		if v.Len() == 0 {
			return v
		}
		mutateTimes := rand.Intn(v.Len())
		for i := 0; i < mutateTimes; i++ {
			elem := v.Index(i)
			mutatedV := VariableMutate(ctx, t.Elem(), elem)
			if t.Elem().AssignableTo(mutatedV.Type()) {
				elem.Set(mutatedV)
			}
		}
		addTimes := rand.Intn(v.Len())
		for i := 0; i < addTimes; i++ {
			curI := rand.Intn(v.Len())
			mutatedV := VariableMutate(ctx, t.Elem(), v.Index(curI))
			if t.Elem().AssignableTo(mutatedV.Type()) {
				v = reflect.Append(v, mutatedV)
			}
		}
		return v
	case reflect.Map:
		mLen := v.Len()
		if mLen == 0 {
			// Init map with random key and value
			mLen = rand.Intn(2)
			v = reflect.MakeMap(t)
			for i := 0; i < mLen; i++ {
				mutatedV := VariableMutate(ctx, t.Elem(), reflect.Zero(t.Elem()))
				if t.Elem().AssignableTo(mutatedV.Type()) {
					newK := VariableMutate(ctx, t.Key(), reflect.Zero(t.Key()))
					v.SetMapIndex(newK.Convert(t.Key()), mutatedV)
				}
			}
		}
		// It will invoke panic, therefore, we add this rule.
		if mLen == 0 {
			return v
		}
		mutateTimes := rand.Intn(mLen)
		keys := v.MapKeys()
		for i := 0; i < mutateTimes; i++ {
			key := keys[rand.Intn(len(keys))]
			mutatedV := VariableMutate(ctx, t.Elem(), reflect.Zero(t.Elem()))
			if t.Elem().AssignableTo(mutatedV.Type()) {
				v.SetMapIndex(key, mutatedV)
			}
		}
		return v
	case reflect.Array:
		break
	case reflect.Chan:
		// v, err := GenerateVariableFromType(.Elem().Underlying(), tv)
		// if err != nil {
		//	panic("cannot transfer")
		// }
		// newV := reflect.ChanOf(reflect.ChanDir(v.Type.ChanDir()), v.Type)
		// v.Value = reflect.ValueOf(newV.Elem())
	case reflect.Func:
		// to avoid valuetostring we return default value
		value := reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			for i := 0; i < t.NumOut(); i++ {
				results = append(results, NewFunctionParam(t.Out(i)))
			}
			return results
		})
		return value
		// TODO: 是否需要处理？？？
		// v.Value = reflect.ValueOf(reflect.FuncOf([]{v.Type}, nil, false))
	}
	return v
}

func SafeSet(f reflect.Value, ctx context.Context, found bool, fake interface{}) {
	defer func() {
		if err := recover(); err != nil {
			f.Set(VariableMutate(ctx, f.Type(), f).Convert(f.Type()))
		} else {
			found = true
		}
	}()
	f.Set(reflect.ValueOf(fake))
}

func NewFunctionParam(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Ptr:
		// to make mockito has address to mock
		return reflect.New(t.Elem())
	// case reflect.Int, reflect.Bool, reflect.String, reflect.Float64, reflect.Float32, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//	return reflect.Zero(t)
	default:
		return reflect.Zero(t)
	}
}

func checkTestCaseStruct(name string) bool {
	var testCase = true
	switch strings.ToLower(name) {
	case "test":
	case "args":
	default:
		testCase = false
	}
	return testCase
}
