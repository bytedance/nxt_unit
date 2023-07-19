// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package parsermodel

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	util "github.com/typa01/go-utils"
)

type Expression struct {
	Value         string //没有[]以及*的类型
	PkgPath       string // 如果map结构，那么PkgPath表示的是key的包路径
	PkgName       string
	IsStar        bool //是否指针
	IsStruct      bool
	IsVariadic    bool //是否...模拟数组的形式
	IsWriter      bool
	Underlying    string //golang类型
	ExtraPkgPaths []string
	IsSignature   bool
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

func (e *Expression) ToString(rootPath string) string {
	log.Printf("rootPath = %v,field path = %v", rootPath, e.PkgPath)
	value := e.Value
	if rootPath != "" && rootPath == e.PkgPath {
		value = strings.Replace(value, e.PkgName+".", "", 1)
	}

	if e.IsVariadic {
		return fmt.Sprintf("...%s", value[2:])
	}
	return fmt.Sprintf("%s", value)
}

type Field struct {
	Name   string
	Type   *Expression
	Index  int
	Fields []*Field
}

func (f *Field) InitiateStruct() string {
	builder := util.NewStringBuilder()
	structS := f.ToStringForStruct()
	if strings.HasPrefix(structS, "*") {
		structS = strings.Replace(structS, "*", "&", 1)
	}
	result := fmt.Sprint(structS, "{}")
	builder.Append(result)
	return builder.ToString()
}

func (f *Field) InitiateReceiver() string {
	builder := util.NewStringBuilder()
	structS := f.ToStringForStruct()
	builder.Append(structS)
	return builder.ToString()
}

func (f *Field) InitiateReceiverForMock() string {
	builder := util.NewStringBuilder()
	structS := f.ToStringForStruct()
	builder.Append("(")
	builder.Append(structS)
	builder.Append(")")
	builder.Append(".")
	return builder.ToString()
}

// InitiateVariable used for gotype dynamic transformation
func (f *Field) InitiateVariable() string {
	builder := util.NewStringBuilder()
	structS := f.ToStringForStruct()
	structS = fmt.Sprint(" ", structS)
	builder.Append(structS)
	return builder.ToString()
}

func (f *Field) IsBasicType() bool {
	return isBasicType(f.Type.String()) || isBasicType(f.Type.Underlying)
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

func (f *Field) ShortName() string {
	return strings.ToLower(string([]rune(f.Type.Value)[0]))
}

func (f *Field) IsWriter() bool {
	return f.Type.IsWriter
}

func (f *Field) ToStringForStruct() string {
	return fmt.Sprintf("%s", f.Type.Value)
}

func (f *Field) ToString(rootPath string) string {
	return fmt.Sprintf("%s %s", f.Name, f.Type.ToString(rootPath))
}

// ToStringForInitialVariable only used for the signature
func (f *Field) ToStringForInitialVariable() string {
	builder := util.NewStringBuilder()
	if f.Type.IsStar {
		builder.Append("*")
	}
	if f.Type.IsVariadic {
		builder.Append("[]")
	}
	builder.Append(f.Type.Value)
	return builder.ToString()
}

func (f *Field) IsStruct() bool {
	return strings.HasPrefix(f.Type.Underlying, "struct")
}

func (f *Field) IsStar() bool {
	return f.Type.IsStar
}

func (f *Field) IsInterface() bool {
	return strings.HasPrefix(f.Type.Underlying, "interface")
}
func (f *Field) IsList() bool {
	return strings.HasPrefix(f.Type.Underlying, "[]")
}
func (f *Field) IsMap() bool {
	return strings.HasPrefix(f.Type.Underlying, "map")
}

func (f *Field) IsNamed() bool {
	return f.Name != "" && f.Name != "_"
}

func (f *Field) GetRandValue(rootPath string) (res string) {
	if f.Type.IsVariadic {
		return ""
	}
	if f.IsStruct() {
		if len(f.Fields) > 0 {
			builder := util.NewStringBuilder()
			for _, field := range f.Fields {
				value := field.GetStructFieldRandValue(rootPath)
				if len(value) > 0 {
					builder.Append(field.Name).Append(":").Append(value).Append(",\n")
				} else {
					builder.Append(field.Name).Append(":").Append("nil,\n")
				}
			}
			if f.Type.IsStar {
				return fmt.Sprintf("&%+v{\n %v}", f.Type.ToString(rootPath)[1:], builder.ToString())
			}
			return fmt.Sprintf("%+v{\n %v}", f.Type.ToString(rootPath), builder.ToString())
		}
		if f.Type.IsStar {
			return fmt.Sprintf("&%+v{}", f.Type.ToString(rootPath))
		}
		return f.Type.ToString(rootPath) + "{}"
	}

	if f.IsList() || f.IsMap() {
		if f.Type.IsStar {
			return "&" + f.Type.ToString(rootPath) + "{}"
		}
		return f.Type.ToString(rootPath) + "{}"
	}
	if res = GetRandValueByType(f.Type.Value, f.Type.IsStar); len(res) > 0 {
		return res
	}
	return ""
}

func (f *Field) GetStructFieldRandValue(rootPath string) (res string) {
	if f.Type.IsVariadic {
		return ""
	}
	if f.IsStruct() {
		if f.Type.IsStar {
			return fmt.Sprintf("&%+v{}", f.Type.ToString(rootPath)[1:])
		}
		return f.Type.ToString(rootPath) + "{}"
	}
	if f.IsList() || f.IsMap() {
		if f.Type.IsStar {
			return "&" + f.Type.Value + "{}"
		}
		return f.Type.Value + "{}"
	}
	if res = GetRandValueByType(f.Type.Value, f.Type.IsStar); len(res) > 0 {
		return res
	}
	return ""
}

func (f *Field) GetImports(importsMap map[string]string) {
	if f == nil || f.Type.PkgPath == "" {
		return
	}

	if len(f.Fields) > 0 {
		for _, field := range f.Fields {
			field.GetImports(importsMap)
		}
	}

	path := f.Type.PkgPath
	importsMap[path] = path
}

func GetRandValueByType(typeValue string, isStars ...bool) string {
	isStar := false
	if len(isStars) > 0 && isStars[0] {
		isStar = true
	}
	switch typeValue {
	case "context.Context":
		return "context.Background()"
	case "bool":
		if isStar {
			return fmt.Sprintf(`thrift.BoolPtr(%+v)`, GetRandValueByType(typeValue))
		}
		return RandChoice("true", "false")
	case "int", "int8", "int16", "uint",
		"uint8", "uint16":
		if isStar {
			return fmt.Sprintf(`thrift.IntPtr(int(%+v))`, GetRandValueByType(typeValue))
		}
		return ToStr(rand.Int31())
	case "int32":
		if isStar {
			return fmt.Sprintf(`thrift.Int32Ptr(int32(%+v))`, GetRandValueByType(typeValue))
		}
		return ToStr(rand.Int31())
	case "int64", "uint32", "uint64":
		if isStar {
			return fmt.Sprintf(`thrift.Int64Ptr(int64(%+v))`, GetRandValueByType(typeValue))
		}
		return ToStr(rand.Int63())
	case "string":
		if isStar {
			return fmt.Sprintf(`thrift.StringPtr(%+v)`, GetRandValueByType(typeValue))
		}
		return fmt.Sprintf("\"%+v\"", RandomString(0))
	case "float32":
		if isStar {
			return fmt.Sprintf(`thrift.Float32Ptr(%+v)`, GetRandValueByType(typeValue))
		}
		return fmt.Sprintf("%+v", rand.Float32())
	case "float64":
		if isStar {
			return fmt.Sprintf(`thrift.Float64Ptr(%+v)`, GetRandValueByType(typeValue))
		}
		return fmt.Sprintf("%+v", rand.Float64())
	case "error":
		return "nil"
	}

	return "nil"
}

//time 2 10位时间戳
func MakeTimeStampMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Second)
}
func ToStr(i interface{}) string {
	if i == nil {
		return ""
	}
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case int8:
		return strconv.Itoa(int(i.(int8)))
	case int16:
		return strconv.Itoa(int(i.(int16)))
	case int32:
		return strconv.Itoa(int(i.(int32)))
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	case time.Time:
		return ToStr(MakeTimeStampMs(i.(time.Time)))
	}
	return ""
}

func RandChoice(choices ...string) string {
	if len(choices) == 0 {
		return ""
	}
	i := rand.Intn(len(choices))
	return choices[i]
}

func RandomString(n int, allowedChars ...[]rune) string {
	if n == 0 {
		n = rand.Intn(15)
	}
	var letters []rune

	if len(allowedChars) == 0 {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
