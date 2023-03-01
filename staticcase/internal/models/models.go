// Copyright cweill/gotests authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package goparse contains logic for parsing Go files. Specifically it parses
// source and test files into domain models for generating tests.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package models

import (
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	util "github.com/typa01/go-utils"
	"strings"
	"unicode"

	"github.com/bytedance/nxt_unit/atghelper"
)

type Expression struct {
	Value       string
	IsStar      bool
	IsVariadic  bool
	IsInterface bool
	IsWriter    bool
	Underlying  string
	PkgPath     string // 如果map结构，那么PkgPath表示的是key的包路径
	PkgName     string
}

func (e *Expression) String() string {
	value := e.Value
	if e.IsStar {
		value = "*" + value
	}
	if e.IsVariadic {
		return "[]" + value
	}
	return value
}

type Field struct {
	Name  string
	Type  *Expression
	Index int
}

func (f *Field) IsWriter() bool {
	return f.Type.IsWriter
}

func (f *Field) IsStruct() bool {
	return strings.HasPrefix(f.Type.Underlying, "struct")
}

func (f *Field) IsBasicType() bool {
	return isBasicType(f.Type.String()) || isBasicType(f.Type.Underlying)
}

func (f *Field) ToStringForStruct() string {
	return fmt.Sprintf("%s", f.Type.Value)
}

func (f *Field) FieldMaxIndex() int {
	return atgconstant.ReceiverFieldMaxLimit
}

func isBasicType(t string) bool {
	switch t {
	case "bool", "string", "int", "int8", "int16", "int32", "int64", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune",
		"float32", "float64", "complex64", "complex128":
		return true
	default:
		return false
	}
}

func (f *Field) IsNamed() bool {
	return f.Name != "" && f.Name != "_"
}

func (f *Field) ShortName() string {
	return strings.ToLower(string([]rune(f.Type.Value)[0]))
}

type Receiver struct {
	*Field
	Fields []*Field
}

func (r *Receiver) InitiateStruct() string {
	builder := util.NewStringBuilder()
	structS := r.ToStringForStruct()
	if strings.HasPrefix(structS, "*") {
		structS = strings.Replace(structS, "*", "&", 1)
	}
	result := fmt.Sprint(structS, "{}")
	builder.Append(result)
	return builder.ToString()
}

func (r *Receiver) InitiateReceiver() string {
	builder := util.NewStringBuilder()
	structS := r.ToStringForStruct()
	builder.Append(structS)
	return builder.ToString()
}

func (r *Receiver) InitiateReceiverForMock() string {
	builder := util.NewStringBuilder()
	structS := r.ToStringForStruct()
	builder.Append("(")
	builder.Append(structS)
	builder.Append(")")
	builder.Append(".")
	return builder.ToString()
}

// InitiateVariable used for gotype dynamic transformation
func (r *Receiver) InitiateVariable() string {
	builder := util.NewStringBuilder()
	structS := r.ToStringForStruct()
	structS = fmt.Sprint(" ", structS)
	builder.Append(structS)
	return builder.ToString()
}

type Function struct {
	Name             string
	IsExported       bool
	Receiver         *Receiver
	Parameters       []*Field
	Results          []*Field
	RowData          string
	ReturnsError     bool
	ContainAnonFuncs int
}

func (f *Function) TestParameters() []*Field {
	var ps []*Field
	for _, p := range f.Parameters {
		if p.IsWriter() {
			continue
		}
		ps = append(ps, p)
	}
	return ps
}

func (f *Function) TestResults() []*Field {
	var ps []*Field
	ps = append(ps, f.Results...)
	for _, p := range f.Parameters {
		if !p.IsWriter() {
			continue
		}
		ps = append(ps, &Field{
			Name: p.Name,
			Type: &Expression{
				Value:      "string",
				IsWriter:   true,
				Underlying: "string",
			},
			Index: len(ps),
		})
	}
	return ps
}

func (f *Function) ReturnsMultiple() bool {
	return len(f.Results) > 1
}

func (f *Function) OnlyReturnsOneValue() bool {
	return len(f.Results) == 1 && !f.ReturnsError
}

func (f *Function) OnlyReturnsError() bool {
	return len(f.Results) == 0 && f.ReturnsError
}

func (f *Function) FullName() string {
	var r string
	if f.Receiver != nil {
		if f.Receiver.Type.IsStar {
			r = fmt.Sprint("*", f.Receiver.Type.Value)
		} else {
			r = f.Receiver.Type.Value
		}
	}
	return r + f.Name
}

func (f *Function) ReceiverName() string {
	var r string
	if f.Receiver != nil {
		if f.Receiver.Type.IsStar {
			r = fmt.Sprint("*", f.Receiver.Type.Value)
		} else {
			r = f.Receiver.Type.Value
		}
	}
	return r
}

func (f *Function) TestName() string {
	if strings.HasPrefix(f.Name, "Test") {
		return f.Name
	}
	if f.Receiver != nil {
		receiverType := f.Receiver.Type.Value
		if unicode.IsLower([]rune(receiverType)[0]) {
			receiverType = "_" + receiverType
		}
		return "Test" + receiverType + "_" + f.Name + atghelper.RandStringBytes(4) + "SU"
	}
	if unicode.IsLower([]rune(f.Name)[0]) {
		return "Test_" + f.Name + "_" + atghelper.RandStringBytes(4) + "SU"
	}
	return "Test" + f.Name + "_" + atghelper.RandStringBytes(4) + "SU"
}

func (f *Function) IsNaked() bool {
	return f.Receiver == nil && len(f.Parameters) == 0 && len(f.Results) == 0
}

func (f *Function) AnonFuncsCount() int {
	return f.ContainAnonFuncs
}

type Import struct {
	Name, Path string
}

type Header struct {
	Comments        []string
	Package         string
	Imports         []*Import
	OriginalImports []*Import
	Code            []byte
}

type Path string

func (p Path) TestPath() string {
	if !p.IsTestPath() {
		return strings.TrimSuffix(string(p), ".go") + "_test.go"
	}
	return string(p)
}

func (p Path) IsTestPath() bool {
	return strings.HasSuffix(string(p), "_test.go")
}
