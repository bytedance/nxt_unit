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
package setup

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
)

type LiteralVisitor struct {
	LiteralMap map[reflect.Type][]reflect.Value
}

type HeaderVisitor struct {
	LiteralMap map[reflect.Type][]reflect.Value
}

type FileVisitor struct {
	LiteralMap   map[reflect.Type][]reflect.Value
	functionName string
}

func NewLiteralVisitor() (v *LiteralVisitor) {
	v = new(LiteralVisitor)
	v.LiteralMap = make(map[reflect.Type][]reflect.Value)
	return
}

func NewHeaderVisitor() (v *HeaderVisitor) {
	v = new(HeaderVisitor)
	v.LiteralMap = make(map[reflect.Type][]reflect.Value)
	return
}

func NewFileVisitor(functionName string) (v *FileVisitor) {
	v = new(FileVisitor)
	v.LiteralMap = make(map[reflect.Type][]reflect.Value)
	v.functionName = functionName
	return
}

// literalNewName 记录复杂数据
func (lv *LiteralVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BasicLit:
		switch n.Kind {
		case token.INT:
			i, err := strconv.Atoi(n.Value)
			if err != nil {
				return nil
			}
			if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(i)]; ok {
				lv.LiteralMap[reflect.TypeOf(i)] = append(LiteralList, reflect.ValueOf(i))
			} else {
				lv.LiteralMap[reflect.TypeOf(i)] = []reflect.Value{reflect.ValueOf(i)}
			}
		case token.FLOAT:
			f, err := strconv.ParseFloat(n.Value, 64)
			// refer the line 50.
			if err == nil {
				if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(f)]; ok {
					lv.LiteralMap[reflect.TypeOf(f)] = append(LiteralList, reflect.ValueOf(f))
				} else {
					lv.LiteralMap[reflect.TypeOf(f)] = []reflect.Value{reflect.ValueOf(f)}
				}
			}
		case token.IMAG:
			// TODO: Complex Number not support yet
		case token.CHAR:
			str, _ := strconv.Unquote(n.Value)
			if str != "" {
				if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(str[0])]; ok {
					lv.LiteralMap[reflect.TypeOf(str[0])] = append(LiteralList, reflect.ValueOf(str[0]))
				} else {
					lv.LiteralMap[reflect.TypeOf(str[0])] = []reflect.Value{reflect.ValueOf(str[0])}
				}
			}
		case token.STRING:
			str, err := strconv.Unquote(n.Value)
			if err == nil {
				if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(str)]; ok {
					lv.LiteralMap[reflect.TypeOf(str)] = append(LiteralList, reflect.ValueOf(str))
				} else {
					lv.LiteralMap[reflect.TypeOf(str)] = []reflect.Value{reflect.ValueOf(str)}
				}
			}
		}
	case *ast.StructType:
		// Rethink about the struct type. because it does not provide us the information. For example,
		// a = 5. 5 is meanning.
		if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(n.Fields)]; ok {
			lv.LiteralMap[reflect.TypeOf(n.Fields)] = append(LiteralList, reflect.ValueOf(n.Fields))
		} else {
			lv.LiteralMap[reflect.TypeOf(n.Fields)] = []reflect.Value{reflect.ValueOf(n.Fields)}
		}
	case *ast.ArrayType:
		if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(n.Elt)]; ok {
			lv.LiteralMap[reflect.TypeOf(n.Elt)] = append(LiteralList, reflect.ValueOf(n.Elt))
		} else {
			lv.LiteralMap[reflect.TypeOf(n.Elt)] = []reflect.Value{reflect.ValueOf(n.Elt)}
		}
	case *ast.SliceExpr:
		if LiteralList, ok := lv.LiteralMap[reflect.TypeOf(n.X)]; ok {
			lv.LiteralMap[reflect.TypeOf(n.X)] = append(LiteralList, reflect.ValueOf(n.X))
		} else {
			lv.LiteralMap[reflect.TypeOf(n.X)] = []reflect.Value{reflect.ValueOf(n.X)}
		}
	case *ast.InterfaceType:
		initValue := node.(*ast.InterfaceType)
		//bes, err := json.Marshal(initValue)
		addContext(initValue)
	case *ast.ChanType:
		newV := reflect.ChanOf(reflect.ChanDir(n.Value.Pos()), reflect.TypeOf(n.Value))
		if LiteralList, ok := lv.LiteralMap[newV]; ok {
			lv.LiteralMap[reflect.TypeOf(newV)] = append(LiteralList, reflect.ValueOf(newV))
		} else {
			lv.LiteralMap[reflect.TypeOf(newV)] = []reflect.Value{reflect.ValueOf(newV)}
		}
	}
	return lv
}

// This is used to parse the content in the target function.
func (v *FileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.IfStmt:
		literalVisitor := NewLiteralVisitor()
		ast.Walk(literalVisitor, n.Cond)
		ast.Walk(v, n.Cond)
		for key, value := range v.LiteralMap {
			if LiteralList, ok := v.LiteralMap[key]; ok {
				v.LiteralMap[key] = append(LiteralList, value...)
			} else {
				v.LiteralMap[key] = value
			}
		}
	case *ast.FuncDecl:
		curName := n.Name.Name
		if n.Recv != nil && len(n.Recv.List) > 0 {
			curName = fmt.Sprintf("%v", reflect.ValueOf(n.Recv.List[0].Type).Elem().FieldByIndex([]int{1})) + "." + curName
		}
		if curName == v.functionName {
			ast.Walk(v, n.Body)
		}
	}
	return v
}

// This is used to parse the struct, interface, constant in the head of the file.
func (v *HeaderVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		headerVisitor := NewLiteralVisitor()
		for _, spec := range n.Specs {
			// TODO: consider more case, and test it.
			if n, ok := spec.(*ast.TypeSpec); ok {
				ast.Walk(headerVisitor, n.Type)
			}
			if n, ok := spec.(*ast.ValueSpec); ok {
				for _, ele := range n.Values {
					ast.Walk(headerVisitor, ele)
				}
			}
		}

		for key, value := range v.LiteralMap {
			if LiteralList, ok := v.LiteralMap[key]; ok {
				v.LiteralMap[key] = append(LiteralList, value...)
			} else {
				v.LiteralMap[key] = value
			}
		}
	}
	return v
}

// addContext 添加context参数
func addContext(iface *ast.InterfaceType) {
	// 接口方法不为空时，遍历接口方法
	if iface.Methods != nil || iface.Methods.List != nil {
		for _, v := range iface.Methods.List {
			ft := v.Type.(*ast.FuncType)
			hasContext := false
			// 判断参数中是否包含context.Context类型
			for _, v := range ft.Params.List {
				if expr, ok := v.Type.(*ast.SelectorExpr); ok {
					if ident, ok := expr.X.(*ast.Ident); ok {
						if ident.Name == "context" {
							hasContext = true
						}
					}
				}
			}
			// 为没有context参数的方法添加context参数
			if !hasContext {
				ctxField := &ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent("ctx"),
					},
					// Notice: 没有考虑import别名的情况
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("context"),
						Sel: ast.NewIdent("Context"),
					},
				}
				list := []*ast.Field{
					ctxField,
				}
				ft.Params.List = append(list, ft.Params.List...)
			}
		}
	}
}

func GetLiteralFromFile(SRCFile string, FunctionName string) (map[reflect.Type][]reflect.Value, error) {
	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs(SRCFile)
	astFile, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileVisitor := NewFileVisitor(FunctionName)
	ast.Walk(fileVisitor, astFile)
	headerVisitor := NewHeaderVisitor()
	ast.Walk(headerVisitor, astFile)

	// Combine two visitors.
	for k, v := range headerVisitor.LiteralMap {
		fileVisitor.LiteralMap[k] = v
	}
	return fileVisitor.LiteralMap, nil
}
