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
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
)

func TestReflectString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"default",
			"\nfunc isNil(object interface{}) bool {\n\tif object == nil {\n\t\treturn true\n\t}\n\n\tvalue := reflect.ValueOf(object)\n\tkind := value.Kind()\n\tisNilableKind := containsKind(\n\t\t[]reflect.Kind{\n\t\t\treflect.Chan, reflect.Func,\n\t\t\treflect.Interface, reflect.Map,\n\t\t\treflect.Ptr, reflect.Slice},\n\t\tkind)\n\n\tif isNilableKind && value.IsNil() {\n\t\treturn true\n\t}\n\n\treturn false\n}\n\nfunc containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {\n\tfor i := 0; i < len(kinds); i++ {\n\t\tif kind == kinds[i] {\n\t\t\treturn true\n\t\t}\n\t}\n\n\treturn false\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReflectString(); got != tt.want {
				t.Errorf("ReflectString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReflectStringV2(t *testing.T) {
	type args struct {
		option atgconstant.Options
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test",
			args{
				option: atgconstant.Options{Uid: "test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReflectStringV2(tt.args.option)
		})
	}
}

func TestReflectStringV3(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"test",
			args{
				id: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReflectStringV3(tt.args.id)
		})
	}
}
