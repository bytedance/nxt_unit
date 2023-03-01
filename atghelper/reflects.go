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
package atghelper

import (
	"context"
	"reflect"
)

var (
	reflectTypeSetForNilValidation = func() map[reflect.Kind]struct{} {
		emptyStruct := struct{}{}
		kindSet := make(map[reflect.Kind]struct{})
		kindSet[reflect.Chan] = emptyStruct
		kindSet[reflect.Func] = emptyStruct
		kindSet[reflect.Interface] = emptyStruct
		kindSet[reflect.Ptr] = emptyStruct
		kindSet[reflect.UnsafePointer] = emptyStruct
		return kindSet
	}()
)

// IsValueNil returns if reflect.Value is nil.
// panics when value is invalid.
func IsValueNil(value reflect.Value) bool {
	if !value.IsValid() {
		panic("value is invalid")
	}
	kind := value.Kind()
	if _, exists := reflectTypeSetForNilValidation[kind]; exists {
		return value.IsNil()
	}
	return false
}

// IsTypeFunction returns if the argument is reflect.Func.
func IsTypeFunction(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	return typ.Kind() == reflect.Func
}

// IsTypeIteratee returns if the argument is an iteratee.
func IsTypeIteratee(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	kind := typ.Kind()
	return kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}

// IsTypeCollection returns if the argument is a slice/array.
func IsTypeCollection(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	kind := typ.Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

// IsTypeSlice returns if argument is slice type
func IsTypeSlice(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	return typ.Kind() == reflect.Slice
}

// IsTypeArray returns if argument is array type
func IsTypeArray(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	return typ.Kind() == reflect.Array
}

// IsTypeMap returns if the argument is a map
func IsTypeMap(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	return typ.Kind() == reflect.Map
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// IsTypeImplementsError returns if the argument is impl of error
func IsTypeImplementsError(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	return typ.Implements(errorType)
}

var ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()

// IsTypeImplementsContext returns if the argument is impl of context.Context
func IsTypeImplementsContext(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	return typ.Implements(ctxType)
}
