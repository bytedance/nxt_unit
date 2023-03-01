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
package astutil

import (
	"fmt"
	"strings"

	"github.com/bytedance/nxt_unit/atgconstant"
)

// ReflectString returns the code in a string format
func ReflectString() string {
	return "\nfunc isNil(object interface{}) bool {\n\tif object == nil {\n\t\treturn true\n\t}\n\n\tvalue := reflect.ValueOf(object)\n\tkind := value.Kind()\n\tisNilableKind := containsKind(\n\t\t[]reflect.Kind{\n\t\t\treflect.Chan, reflect.Func,\n\t\t\treflect.Interface, reflect.Map,\n\t\t\treflect.Ptr, reflect.Slice},\n\t\tkind)\n\n\tif isNilableKind && value.IsNil() {\n\t\treturn true\n\t}\n\n\treturn false\n}\n\nfunc containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {\n\tfor i := 0; i < len(kinds); i++ {\n\t\tif kind == kinds[i] {\n\t\t\treturn true\n\t\t}\n\t}\n\n\treturn false\n}"
}

// Deprecated
func ReflectStringV2(option atgconstant.Options) string {
	isNilFunc := fmt.Sprintf("isNil%s(object interface{})", option.Uid)
	// PanicTrace(kb int)
	panicTraceFunc := fmt.Sprintf("PanicTrace%s(kb int)", option.Uid)
	originFunc := ReflectString() + `
	func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}
`
	originFunc = strings.ReplaceAll(originFunc, "isNil(object interface{})", isNilFunc)
	originFunc = strings.ReplaceAll(originFunc, "PanicTrace(kb int)", panicTraceFunc)
	originFunc = strings.ReplaceAll(originFunc, "containsKind", "containsKind"+option.Uid)
	return originFunc
}

func ReflectStringV3(id string) string {

	isNilFunc := fmt.Sprintf("isNil%s(object interface{})", id)
	// PanicTrace(kb int)
	panicTraceFunc := fmt.Sprintf("PanicTrace%s(kb int)", id)
	originFunc := ReflectString() + `
	func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := smartunit_runtime.Stack(stack, true)
	start := smartunit_bytes.Index(stack, s)
	stack = stack[start:length]
	start = smartunit_bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := smartunit_bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = smartunit_bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = smartunit_bytes.TrimRight(stack, "\n")
	return stack
}
`
	originFunc = strings.ReplaceAll(originFunc, "isNil(object interface{})", isNilFunc)
	originFunc = strings.ReplaceAll(originFunc, "PanicTrace(kb int)", panicTraceFunc)
	originFunc = strings.ReplaceAll(originFunc, "containsKind", "containsKind"+id)
	return originFunc
}
