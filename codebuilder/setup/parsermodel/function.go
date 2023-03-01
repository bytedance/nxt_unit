// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package parsermodel

import (
	"fmt"
	"strings"

	util "github.com/typa01/go-utils"
)

type Function struct {
	RootFunc     *Function
	PkgPath      string
	PkgName      string
	Name         string
	IsExported   bool
	Receiver     *Field
	Parameters   []*Field
	Results      []*Field
	ReturnsError bool
	NeedRandCase bool
	tags         []string // rpc
}

func (f *Function) OnlyReturnsOneValue() bool {
	return len(f.Results) == 1 && !f.ReturnsError
}

func (f *Function) OnlyReturnsError() bool {
	return len(f.Results) == 0 && f.ReturnsError
}

func (f *Function) ReturnsMultiple() bool {
	return len(f.Results) > 1
}

func (f *Function) IsNaked() bool {
	return f.Receiver == nil && len(f.Parameters) == 0 && len(f.Results) == 0
}

// TODO: 这里要改掉.
func (f *Function) MockFunc() string {
	if f.Receiver != nil {
		return f.mockObjFunc()
	}

	return f.mockCommonFunc()
}

func (f *Function) mockCommonFunc() string {
	builder := util.NewStringBuilder()

	mockName := "mock" + f.Name
	builder.Append(mockName)
	builder.Append(" := ")
	builder.Append("func(")

	params := getParams(f.RootFunc, f.Parameters)
	builder.Append(strings.Join(params, ", "))
	builder.Append(")")

	results := getParams(f.RootFunc, f.Results)
	result := strings.Join(results, ", ")
	if len(f.Results) > 1 {
		builder.Append("(").Append(result).Append(")")
	}

	if len(f.Results) == 1 {
		builder.Append(result)
	}

	builder.Append("{\n")
	f.GetResult(builder)
	builder.Append("\n}")

	builder.Append("\nMock(")
	builder.Append(f.MethodName())
	builder.Append(").To(")
	builder.Append(mockName)
	builder.Append(").Build()\n\n")

	return builder.ToString()
}

func (f *Function) mockObjFunc() string {
	builder := util.NewStringBuilder()

	mockName := "mock" + f.Name
	builder.Append(mockName)
	builder.Append(" := ")
	builder.Append("func(")

	params := getParams(f.RootFunc, f.Parameters)
	builder.Append(strings.Join(params, ", "))
	builder.Append(")")

	results := getParams(f.RootFunc, f.Results)
	result := strings.Join(results, ", ")
	if len(f.Results) > 1 {
		builder.Append("(").Append(result).Append(")")
	}

	if len(f.Results) == 1 {
		builder.Append(result)
	}

	builder.Append("{\n")
	f.GetResult(builder)
	builder.Append("\n}")
	builder.Append("\nMock((")
	builder.Append(f.GetReceiverString())
	builder.Append(").")
	builder.Append(f.Name)
	builder.Append(").To(")
	builder.Append(mockName)
	builder.Append(").Build()\n\n")

	return builder.ToString()
}

func getParams(RootFunc *Function, fields []*Field) []string {
	var rootPath = ""
	if RootFunc != nil {
		rootPath = RootFunc.PkgPath
	}

	if len(fields) == 0 {
		return nil
	}

	var params []string
	for _, parameter := range fields {
		params = append(params, parameter.ToString(rootPath))
	}
	return params
}

func (f *Function) GetResult(sb *util.StringBuilder) {
	if len(f.Results) == 0 {
		return
	}
	sb.Append("return ")
	var returns []string
	for _, result := range f.Results {
		value := result.GetRandValue(f.RootPath())
		returns = append(returns, value)
	}
	sb.Append(strings.Join(returns, ",\n"))
}

func (f *Function) RootPath() string {
	var rootPath = ""
	if f.RootFunc != nil {
		rootPath = f.RootFunc.PkgPath
	}
	return rootPath
}

func (f *Function) MethodName() string {
	if f.RootFunc != nil && f.PkgName == f.RootFunc.PkgName {
		return f.Name
	}

	return fmt.Sprintf("%v.%v", f.PkgName, f.Name)
}

func (f *Function) GetReceiverString() string {
	value := f.Receiver.Type.Value

	if f.RootFunc != nil && f.Receiver.Type.PkgPath == f.RootFunc.PkgPath {
		value = strings.Split(value, ".")[1]
		return fmt.Sprintf("*%v", value)
	}

	if strings.Contains(value, "*") {
		return value
	}
	return fmt.Sprintf("*%v", value)
}

func (f *Function) GetImports() []string {
	var importsMap = make(map[string]string)

	// Receiver
	f.Receiver.GetImports(importsMap)

	//param
	for _, parameter := range f.Parameters {
		parameter.GetImports(importsMap)
	}

	//param
	for _, result := range f.Results {
		result.GetImports(importsMap)
	}

	if len(importsMap) == 0 {
		return nil
	}

	// 移除根路径
	if f.RootFunc != nil {
		delete(importsMap, f.RootFunc.PkgPath)
	}

	var imports []string
	for key := range importsMap {
		imports = append(imports, key)
	}

	return imports
}
