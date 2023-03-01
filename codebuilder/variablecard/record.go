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
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"github.com/bytedance/nxt_unit/smartunitvariablebuild"
	"reflect"
	"regexp"
	"strings"

	util "github.com/typa01/go-utils"
)

type MocksRecord []string
type MonkeyOutputMap map[string]string //key:functionName value is function return string

type SpecialValue interface {
	ValueToCode() string
}

// Because we cannot directly assign nil to the variable. Thus, we assign the nil to the variable if
// the value itself is nil.

// the context contains the follow information:
// ID: ID tell us the relationship between the variable and statement.
// Level refers to variable level. For example: tiktok{a,b,c}. tiktok's level is 0. a's level is 1.
// It is used for name triming.

// value refers to the reflect value
// But some times it might refers to the mocksRecord.
func ValueToString(ctx context.Context, v reflect.Value) (code string) {
	vtx, ok := contexthelper.GetVariableContext(ctx)
	if !ok {
		return "nil"
	}
	// it will stop render field
	// if len(vtx.MockedRecord) != 0 {
	// 	return fmt.Sprintf("func(){%s}", strings.Join(vtx.MockedRecord, ";\n"))
	// }
	// If the value cannot be nil, we don't set it to the nil.
	if !v.IsValid() || atghelper.IsValueNil(v) {
		return "nil"
	}

	if special, ok := v.Interface().(SpecialValue); ok {
		return special.ValueToCode()
	}

	if mocks, ok := v.Interface().(MocksRecord); ok {
		return fmt.Sprintf("func(){%s}", strings.Join(mocks, ";\n"))
	}
	if monkeyOutputMap, ok := v.Interface().(MonkeyOutputMap); ok {
		return getMonkeyOutputStrMap(monkeyOutputMap)
	}
	specialValue, ok := smartunitvariablebuild.RenderVariableV3(ctx, v)
	if ok {
		return specialValue
	}
	builder := util.NewStringBuilder()

	t := v.Type()
	switch t.Kind() {
	case reflect.Int, reflect.Bool, reflect.String, reflect.Float64, reflect.Float32, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		builder.Append(ParameterToString(ctx, t, v))
	case reflect.Ptr:
		newT := t.Elem()
		switch newT.Kind() {
		case reflect.Int, reflect.Bool, reflect.String, reflect.Float64, reflect.Float32, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if newT.PkgPath() != "" {
				typeName := newT.String()
				if newT.PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					typeName = removeSelfImported(typeName)
				} else {
					pkgName := atghelper.GetPkgName(newT.PkgPath())
					newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, newT.PkgPath())
					typeName = atghelper.ReplacePkgName(typeName, newPkgName, pkgName)
				}
				if newT.Kind().String() != newT.Name() {
					return BasicPtrReNameToString(ctx, newT, v.Elem(), typeName)
				}
				builder.Append("&")
				builder.Append(typeName)
				builder.Append("(")
				builder.Append(ParameterToString(ctx, newT, v.Elem()))
				builder.Append(")")
			} else {
				builder.Append(BasicPtrToString(ctx, newT, v.Elem()))
			}
		case reflect.Slice, reflect.Map:
			vtx.Level++
			ctx = contexthelper.SetVariableContext(ctx, vtx)
			builder.Append(fmt.Sprint("&", ValueToString(ctx, v.Elem())))
		case reflect.Struct:
			vtx.Level++
			ctx = contexthelper.SetVariableContext(ctx, vtx)
			builder.Append(fmt.Sprint("&", ValueToString(ctx, v.Elem())))
		case reflect.Interface:
			builder.Append("nil")
		}
	case reflect.Slice:
		builder.Append(fmt.Sprint(trimName(ctx, t, nil), "{", SliceElementToString(ctx, v), "}"))
	case reflect.Array:
		builder.Append(fmt.Sprint(trimName(ctx, t, nil), "{", SliceElementToString(ctx, v), "}"))
	case reflect.Map:
		builder.Append(fmt.Sprint(trimName(ctx, t, nil), "{\n", MapElementToString(ctx, v), "}"))
	case reflect.Struct:
		// Special logic, to remove the current package if the struct package is the same with the current one
		builder.Append(fmt.Sprint(trimName(ctx, t, nil), "{\n", StructFieldToString(ctx, v), "}"))
	case reflect.Interface:
		return "nil"
	case reflect.Func:
		builder.Append("func(")
		inputVariadic := &InputVariadic{t.IsVariadic(), false}
		for i := 0; i < t.NumIn(); i++ {
			if i == t.NumIn()-1 {
				inputVariadic.IsVariadicParam = true
			}
			if !atghelper.IsTypeExported(t.In(i)) {
				return "unexport variable"
			}
			builder.Append(trimName(ctx, t.In(i), inputVariadic))
			if i == t.NumIn()-1 {
				break
			}
			builder.Append(",")
		}
		builder.Append(")")
		builder.Append(" (")
		for j := 0; j < t.NumOut(); j++ {
			builder.Append(trimName(ctx, t.Out(j), nil))
			if j == t.NumOut()-1 {
				break
			}
			builder.Append(",")
		}
		builder.Append(")")
		// If it is overpass
		_, ok := contexthelper.GetBuilderVector(ctx)
		if ok {
			return builder.ToString()
		}
		var result []string
		var in []reflect.Value
		for i := 0; i < t.NumIn(); i++ {
			if t.IsVariadic() && i == t.NumIn()-1 {
				// If v is a variadic function, Call creates the variadic slice parameter
				// itself, copying in the corresponding values.
				in = append(in, reflect.Zero(t.In(i).Elem()))
			} else {
				in = append(in, reflect.Zero(t.In(i)))
			}
		}
		for _, value := range v.Call(in) {
			result = append(result, ValueToString(ctx, value))
		}
		r := strings.Join(result, ",")
		builder.Append("{ return ")
		builder.Append(r)
		builder.Append(" }")
	}
	return builder.ToString()
}

func BasicPtrToString(ctx context.Context, t reflect.Type, v reflect.Value) string {
	builder := util.NewStringBuilder()
	switch t.Kind() {
	case reflect.Int:
		builder.Append(fmt.Sprint("atgconv.IntPtr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Int8:
		builder.Append(fmt.Sprint("atgconv.Int8Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Int16:
		builder.Append(fmt.Sprint("atgconv.Int16Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Int32:
		builder.Append(fmt.Sprint("atgconv.Int32Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Int64:
		builder.Append(fmt.Sprint("atgconv.Int64Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Uint:
		builder.Append(fmt.Sprint("atgconv.UintPtr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Uint8:
		builder.Append(fmt.Sprint("atgconv.Uint8Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Uint16:
		builder.Append(fmt.Sprint("atgconv.Uint16Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Uint32:
		builder.Append(fmt.Sprint("atgconv.Uint32Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Uint64:
		builder.Append(fmt.Sprint("atgconv.Uint64Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Float32:
		builder.Append(fmt.Sprint("atgconv.Float32Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Float64:
		builder.Append(fmt.Sprint("atgconv.Float64Ptr(", ParameterToString(ctx, t, v), ")"))
	case reflect.Bool:
		builder.Append(fmt.Sprint("atgconv.BoolPtr(", ParameterToString(ctx, t, v), ")"))
	case reflect.String:
		builder.Append(fmt.Sprint("atgconv.StringPtr(", ParameterToString(ctx, t, v), ")"))
	}
	return builder.ToString()
}

func BasicPtrReNameToString(ctx context.Context, t reflect.Type, v reflect.Value, typeName string) string {
	// func string example: func() *IntegerType { tt := IntegerType(0); return &tt }()
	builder := util.NewStringBuilder()
	builder.Append("func() ")
	funReturnType := fmt.Sprintf("*%s", typeName)
	funReturn := "&tmp"
	builder.Append(funReturnType)
	builder.Append(fmt.Sprintf(" {tmp := %s;return %s}()", ParameterToString(ctx, t, v), funReturn))
	return builder.ToString()
}

// ParameterToString we need record basic value Corpus
func ParameterToString(ctx context.Context, t reflect.Type, v reflect.Value) string {
	builder := util.NewStringBuilder()
	switch t.Kind() {
	case reflect.Int:
		builder.Append(fmt.Sprint(v.Int()))
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		typeName := t.String()
		if t.PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
			typeName = removeSelfImported(typeName)
		} else {
			pkgPath := t.PkgPath()
			pkgName := atghelper.GetPkgName(pkgPath)
			newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, pkgPath)
			typeName = atghelper.ReplacePkgName(typeName, newPkgName, pkgName)
		}
		builder.Append(fmt.Sprint(typeName, "(", v.Int(), ")"))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typeName := t.String()
		if t.PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
			typeName = removeSelfImported(typeName)
		} else {
			pkgName := atghelper.GetPkgName(t.PkgPath())
			newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.PkgPath())
			typeName = atghelper.ReplacePkgName(typeName, newPkgName, pkgName)
		}
		builder.Append(fmt.Sprint(typeName, "(", v.Uint(), ")"))
	case reflect.Float64:
		builder.Append(fmt.Sprint(v.Float()))
	case reflect.Float32:
		typeName := t.String()
		if t.PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
			typeName = removeSelfImported(typeName)
		} else {
			pkgName := atghelper.GetPkgName(t.PkgPath())
			newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.PkgPath())
			typeName = atghelper.ReplacePkgName(typeName, newPkgName, pkgName)
		}
		builder.Append(fmt.Sprint(typeName, "(", v.Float(), ")"))
	case reflect.Bool:
		builder.Append(fmt.Sprint(v.Bool()))
	case reflect.String:
		content := v.String()
		if strings.Contains(content, "\"") {
			content = strings.ReplaceAll(content, "\"", "\\\"")
		}
		builder.Append(fmt.Sprint("\"", content, "\""))
	}
	return builder.ToString()
}

func SliceElementToString(ctx context.Context, v reflect.Value) string {
	if atghelper.IsValueNil(v) {
		return "nil"
	}
	builder := util.NewStringBuilder()
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		builder.Append(ValueToString(ctx, elem))
		builder.Append(",\n")
	}
	return builder.ToString()
}

func MapElementToString(ctx context.Context, v reflect.Value) string {
	if atghelper.IsValueNil(v) {
		return "nil"
	}
	builder := util.NewStringBuilder()
	for _, k := range v.MapKeys() {
		elem := v.MapIndex(k)
		builder.Append(ValueToString(ctx, k))
		builder.Append(":")
		builder.Append(ValueToString(ctx, elem))
		builder.Append(",\n")
	}
	return builder.ToString()
}

// TODO: use the json to do the transform
func StructFieldToString(ctx context.Context, v reflect.Value) string {
	if atghelper.IsValueNil(v) {
		return "nil"
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	builder := util.NewStringBuilder()
	// NumField will be panic if the value is a pointer.
	for i := 0; i < v.NumField(); i++ {
		// if !v.Field(i).CanSet() {
		// 	continue // to avoid panic to set on unexported field in struct
		// }
		// PkgPath is the package path that qualifies a lower case (unexported)
		// field name. It is empty for upper case (exported) field names.
		// See https://golang.org/ref/spec#Uniqueness_of_identifiers
		if !(v.Field(i).CanInterface()) {
			continue
		}
		res := func() bool {
			defer func() bool {
				return false
			}()
			if v.Field(i).IsValid() && v.Field(i).IsZero() {
				return true
			}
			return false
		}()
		if res {
			continue
		}

		if atghelper.IsTypeExported(v.Field(i).Type()) && (v.Field(i).IsValid() && !v.Field(i).IsZero()) {
			typeName := v.Type().Field(i).Name
			if v.Type().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
				typeName = removeSelfImported(typeName)
			} else {
				pkgPath := v.Type().PkgPath()
				pkgName := atghelper.GetPkgName(pkgPath)
				newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, pkgPath)
				typeName = atghelper.ReplacePkgName(typeName, newPkgName, pkgName)
			}
			builder.Append(typeName)
			builder.Append(":")
			builder.Append(ValueToString(ctx, v.Field(i)))
			builder.Append(",\n")
		}
	}
	return builder.ToString()
}

type InputVariadic struct {
	IsVariadic      bool
	IsVariadicParam bool
}

// Todo: we only handle the level 0's change, for other level's package replacement. We will do it.
func trimName(ctx context.Context, t reflect.Type, inputVariadic *InputVariadic) string {
	vtx, ok := contexthelper.GetVariableContext(ctx)
	if !ok {
		return "nil"
	}
	if duplicatepackagemanager.GetInstance(ctx).RelativePath() != "" {
		// it seems impossible to remove the package name inside struct
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Chan:
			switch t.Elem().Kind() {
			case reflect.Ptr:
				if t.Elem().Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImported(t.String())
				}
			default:
				if t.Elem().PkgPath() != "" && t.Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImported(t.String())
				}
			}
		case reflect.Map:
			// If the key and value are self-imported
			switch t.Elem().Kind() {
			case reflect.Ptr:
				if t.Elem().Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() && t.Key().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), true, true)
					// If the value are self-imported
				}
				switch t.Key().Kind() {
				case reflect.Ptr:
					if t.Elem().Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() && t.Key().Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
						return removeSelfImportedForMap(t.String(), true, true)
						// If the value are self-imported
					}
				}
			case reflect.Map:
				orignalValueName := t.Elem().String()
				realValueName := trimName(ctx, t.Elem(), &InputVariadic{false, false})
				str := strings.Replace(t.String(), orignalValueName, realValueName, 1)
				if t.Key().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(str, true, false)
				} else {
					return str
				}
			case reflect.Struct:
				if t.Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() && t.Key().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), true, true)
				}
			case reflect.Slice:
				orignalValueName := t.Elem().String()
				realValueName := trimName(ctx, t.Elem(), &InputVariadic{false, false})
				str := strings.Replace(t.String(), orignalValueName, realValueName, 1)
				if t.Key().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(str, true, false)
				} else {
					return str
				}

			default:
				if t.Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() && t.Key().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), true, true)
				}
			}
			// If the value are self-imported
			switch t.Elem().Kind() {
			case reflect.Ptr:
				if t.Elem().Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), false, true)
				}
			case reflect.Struct:
				if t.Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), false, true)
				}
			case reflect.Slice:
				orignalValueName := t.Elem().String()
				realValueName := trimName(ctx, t.Elem(), &InputVariadic{false, false})
				str := strings.Replace(t.String(), orignalValueName, realValueName, 1)
				return removeSelfImportedForMap(str, false, true)
			default:
				if t.Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), false, true)
				}
			}
			// If the key are self-imported
			switch t.Key().Kind() {
			case reflect.Ptr:
				if t.Key().Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), true, false)
				}
			default:
				if t.Key().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
					return removeSelfImportedForMap(t.String(), true, false)
				}
			}
		}
		if t.PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
			return removeSelfImported(t.String())
		}

		if vtx.Level > atgconstant.VariableMaxLevel {
			goto Breakthrough
		}
		// if p is a pointer. p.pkg path is not empty. However, t.pkgpath is empty.
		if t.Kind() == reflect.Ptr && t.Elem().PkgPath() == duplicatepackagemanager.GetInstance(ctx).RelativePath() {
			return removeSelfImported(t.String())
		}
	}

	// If the level exceed the maximum level, we jump to the next part. Otherwise, it will have the duplicated package
	// Problem
Breakthrough:
	// If type's package path is not empty
	if t.PkgPath() != "" {
		pkgName := atghelper.GetPkgName(t.PkgPath())
		newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.PkgPath())
		return atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
	}

	// If type's package path is empty
	switch t.Kind() {
	case reflect.Ptr:
		pkgName := atghelper.GetPkgName(t.Elem().PkgPath())
		newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().PkgPath())
		return atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
	case reflect.Slice:
		sliceName := t.String()
		pkgName := ""
		newPkgName := ""
		switch t.Elem().Kind() {
		case reflect.Ptr:
			pkgName = atghelper.GetPkgName(t.Elem().Elem().PkgPath())
			newPkgName, _ = duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().Elem().PkgPath())
		case reflect.Slice, reflect.Map, reflect.Chan, reflect.Array:
			orignalValueName := t.Elem().String()
			realValueName := trimName(ctx, t.Elem(), &InputVariadic{false, false})
			sliceName = strings.Replace(sliceName, orignalValueName, realValueName, 1)
			return sliceName
		default:
			pkgName = atghelper.GetPkgName(t.Elem().PkgPath())
			newPkgName, _ = duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().PkgPath())
		}
		sliceName = atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
		if inputVariadic != nil && inputVariadic.IsVariadic && inputVariadic.IsVariadicParam {
			sliceName = strings.Replace(sliceName, "[]", "...", 1)
		}
		return sliceName
	case reflect.Array, reflect.Chan:
		switch t.Elem().Kind() {
		case reflect.Ptr:
			pkgName := atghelper.GetPkgName(t.Elem().Elem().PkgPath())
			newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().Elem().PkgPath())
			return atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
		default:
			pkgName := atghelper.GetPkgName(t.Elem().PkgPath())
			newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().PkgPath())
			return atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
		}
	case reflect.Map:
		res := t.String()
		switch t.Key().Kind() {
		case reflect.Ptr:
			// If the key and value need to modify
			switch t.Elem().Kind() {
			case reflect.Ptr:
				if t.Elem().Elem().PkgPath() != "" {
					pkgName := atghelper.GetPkgName(t.Elem().Elem().PkgPath())
					newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().Elem().PkgPath())
					if pkgName != newPkgName {
						res = atghelper.ReplacePkgNameForMap(t.String(), newPkgName, pkgName, false)
					}
				}
			default:
				if t.Elem().PkgPath() != "" {
					pkgName := atghelper.GetPkgName(t.Elem().PkgPath())
					newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().PkgPath())
					if pkgName != newPkgName {
						res = atghelper.ReplacePkgNameForMap(t.String(), newPkgName, pkgName, false)
					}
				}
			}
			if t.Key().Elem().PkgPath() != "" {
				pkgName := atghelper.GetPkgName(t.Key().Elem().PkgPath())
				newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Key().Elem().PkgPath())
				if pkgName != newPkgName {
					return atghelper.ReplacePkgNameForMap(res, newPkgName, pkgName, true)
				}
			}
		default:
			if t.Key().PkgPath() != "" {
				pkgName := atghelper.GetPkgName(t.Key().PkgPath())
				newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Key().PkgPath())
				if pkgName != newPkgName {
					return atghelper.ReplacePkgNameForMap(t.String(), newPkgName, pkgName, true)
				}
			}
		}
		switch t.Elem().Kind() {
		case reflect.Ptr:
			if t.Elem().Elem().PkgPath() != "" {
				pkgName := atghelper.GetPkgName(t.Elem().Elem().PkgPath())
				newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().Elem().PkgPath())
				if pkgName != newPkgName {
					return atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
				}
			}
		default:
			if t.Elem().PkgPath() != "" {
				pkgName := atghelper.GetPkgName(t.Elem().PkgPath())
				newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, t.Elem().PkgPath())
				if pkgName != newPkgName {
					return atghelper.ReplacePkgName(t.String(), newPkgName, pkgName)
				}
			}
		}
	case reflect.Struct:
		// replace pkgName.Type in different defined pkgname condition
		replaceStructKeyMap := make(map[string]string, 0) // newStructKey:newPkgName.Type ,originalStructKey
		for i := 0; i < t.NumField(); i++ {
			name := t.Field(i).Type.Name()
			pkgPath := t.Field(i).Type.PkgPath()
			if pkgPath != "" {
				// XXX: @caoziguang Other place still using the old version GetPkgName so we just add a new one to replace.
				// It does work here but need to be update
				pkgName := atghelper.GetPkgNameV2(t.Field(i).Type.String())

				newPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, pkgPath)
				if pkgPath == duplicatepackagemanager.GetInstance(ctx).RelativePath() && newPkgName == "" {
					replaceStructKeyMap[name] = fmt.Sprintf("%s.%s", pkgName, name)
				} else {
					replaceStructKeyMap[fmt.Sprintf("%s.%s", newPkgName, name)] = fmt.Sprintf("%s.%s", pkgName, name)
				}
			}
		}
		resultName := t.String()
		for newFieldName, originalFieldName := range replaceStructKeyMap {
			resultName = strings.Replace(resultName, originalFieldName, newFieldName, 1)
		}
		return resultName
	}
	// No need to add the duplicated manager because type's pkg path is empty.
	return t.String()
}

// case5 map[a.b]c.d -> map[b]d or map[a.b]d or map[b]c.d
func removeSelfImportedForMap(s string, isKey bool, isValue bool) string {
	if !isKey && !isValue {
		return s
	}
	// case5 map[a.b]c -> map[b]c
	pathMatch, err := regexp.Compile(`^map\[(.*)\]`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS := pathMatch.FindAllString(s, 1)
	left := ""
	right := ""
	if len(matchedS) == 1 {
		// first let's handle left part: a.b
		leftArray := strings.SplitN(atghelper.DeepCopy(matchedS[0]), ".", 2)
		if len(leftArray) == 2 && isKey {
			if strings.Contains(leftArray[0], "*") {
				left = fmt.Sprint("map[*", leftArray[1])
			} else {
				left = fmt.Sprint("map[", leftArray[1])
			}
		}

		// If not "." in the left
		if len(leftArray) == 1 || (len(leftArray) == 2 && !isKey) {
			left = matchedS[0]
		}

		// right c.d
		rightRow := strings.Replace(s, atghelper.DeepCopy(matchedS[0]), "", 2)
		rightArray := strings.SplitN(rightRow, ".", 2)
		if len(rightArray) == 2 && isValue {
			if strings.Contains(rightArray[0], "*") {
				right = fmt.Sprint("*", rightArray[1])
			} else {
				right = rightArray[1]
			}
		}

		if len(rightArray) == 1 || (len(rightArray) == 2 && !isValue) {
			right = rightRow
		}
		return fmt.Sprint(left, right)
	} else {
		// for map type alias
		return removeSelfImported(s)
	}
}

// TODO(siwei.wang): need to apply the duplicated manager
func removeSelfImported(s string) string {
	// case1  [][]*a.b
	pathMatch, err := regexp.Compile(`^\[\]\[\]\*([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS := pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("[][]*", strings.Replace(s, matchedS[0], "", 1))
	}

	// case2  [][]a.b
	pathMatch, err = regexp.Compile(`^\[\]\[\]([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("[][]", strings.Replace(s, matchedS[0], "", 1))
	}

	// case3  []*a.b
	pathMatch, err = regexp.Compile(`^\[\]\*([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("[]*", strings.Replace(s, matchedS[0], "", 1))
	}

	// case4  []a.b
	pathMatch, err = regexp.Compile(`^\[\]([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("[]", strings.Replace(s, matchedS[0], "", 1))
	}

	// case5 chan *a.b
	pathMatch, err = regexp.Compile(`^chan\s\*([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("chan *", strings.Replace(s, matchedS[0], "", 1))
	}

	// case6 chan a.b
	pathMatch, err = regexp.Compile(`^chan\s([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("chan ", strings.Replace(s, matchedS[0], "", 1))
	}

	// case7 *a.b
	pathMatch, err = regexp.Compile(`^\*([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return fmt.Sprint("*", strings.Replace(s, matchedS[0], "", 1))
	}

	// case8 a.b
	pathMatch, err = regexp.Compile(`^([a-zA-Z0-9_-]+)\.`)
	if err != nil {
		fmt.Printf("[removeSelfImported] has error, the error is %v\n", err)
		return s
	}
	matchedS = pathMatch.FindAllString(s, 1)
	if len(matchedS) == 1 {
		return strings.Replace(s, matchedS[0], "", 1)
	}
	// // TBD more
	return s
}

func getMonkeyOutputStrMap(monkeyOutputMap MonkeyOutputMap) string {
	monkeyStr := ""
	headStr := "map[string][]interface{}{"
	tailStr := "}"
	cellHeadStr := "[]interface{}{"
	cellTailStr := "},"

	if len(monkeyOutputMap) > 0 {
		monkeyStr = monkeyStr + headStr
		for funcName, mockReturn := range monkeyOutputMap {
			if mockReturn != "" {
				monkeyStr = monkeyStr + fmt.Sprintf("\"%s\":%s%s%s", funcName, cellHeadStr, mockReturn, cellTailStr)
			}
		}
		monkeyStr = monkeyStr + tailStr
	}
	if monkeyStr != "" {
		return monkeyStr
	} else {
		return "map[string][]interface{}{}"
	}
}
