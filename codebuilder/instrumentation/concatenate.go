// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.

package instrumentation

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	xastutil "golang.org/x/tools/go/ast/astutil"

	"github.com/bytedance/nxt_unit/atghelper"

	"github.com/bytedance/nxt_unit/codebuilder/astutil"

	"golang.org/x/tools/imports"

	"github.com/bytedance/nxt_unit/atgconstant"
)

type TestedFunctionInfo struct {
	PredicatesPos []int
	FunctionNames []string
}

type Import struct {
	Name, Path string
}

// File is a wrapper for the state of a file used in the parser.
// The basic parse tree walker is a method of this type.
type File struct {
	fset             *token.FileSet
	name             string // Name of file.
	astFile          *ast.File
	content          []byte
	edit             *astutil.Buffer
	funcs            []Func
	branches         []Branch
	relatedFunctions []string
	count            int
	lastInit         bool
}

type SrcFileRecord struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	content []byte
	edit    *astutil.Buffer
	imports []*Import
}

type SrcFileModify struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	content []byte
	edit    *astutil.Buffer

	TargetFunctions map[string]struct{}
}

type TrcFile struct {
	fset          *token.FileSet
	name          string // Name of file.
	astFile       *ast.File
	content       []byte
	edit          *astutil.Buffer
	imports       []*Import // src file import
	TrcImports    []*Import // trc file import
	isWrong       bool
	FunctionNames map[string]struct{}
}

type Func struct {
	startByte token.Pos
	endByte   token.Pos
	Name      string
}

var once sync.Once

func Concatenate(src, trc string, tarImports []*Import) error {
	srcImports, err := GetImportsInfosFromFile(src)
	if err != nil {
		fmt.Printf("[Concatenate] err: %v", err)
		return fmt.Errorf("[Concatenate] err: %v", err)
	}

	// added the import to target source
	tgtWriter := &bytes.Buffer{}
	tgtFunction, err := AddedImportToTrc(tgtWriter, trc, srcImports, tarImports)
	if err != nil {
		fmt.Printf("[Concatenate] err: %v", err)
		return fmt.Errorf("[Concatenate] err: %v", err)
	}
	srcWriter := &bytes.Buffer{}
	err = ModifySrcBasedTarget(srcWriter, src, tgtFunction)
	if err != nil {
		fmt.Printf("[Concatenate] ModifySrcBasedTarget err: %v", err)
		return fmt.Errorf("[Concatenate] ModifySrcBasedTarget err: %v", err)
	}

	w := &bytes.Buffer{}
	// append the modified source file which removed the package name and imports to the target file
	w.Write(tgtWriter.Bytes())
	w.Write(srcWriter.Bytes())
	out, err := imports.Process(trc, w.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("[Concatenate] imports.Process err %v", err)
	}
	err = ioutil.WriteFile(trc, out, atgconstant.NewFilePerm)
	if err != nil {
		return fmt.Errorf("[Concatenate] imports.Process err %v", err)
	}

	return nil
}

func GetImportsInfosFromFile(filepath string) ([]*Import, error) {
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile has error: %v", err)
	}
	parsedFile, err := parser.ParseFile(fset, filepath, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ast parser.ParseFilee has error: %v", err)
	}
	file := &SrcFileRecord{
		fset:    fset,
		name:    filepath,
		content: content,
		edit:    astutil.NewBuffer(content),
		astFile: parsedFile,
		imports: make([]*Import, 0),
	}
	ast.Walk(file, file.astFile)
	return file.imports, nil
}

func (f *SrcFileRecord) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ImportSpec:
		if n.Path == nil {
			return nil
		}
		name := ""
		if n.Name != nil {
			name = n.Name.Name
		}
		f.imports = append(f.imports, &Import{
			Name: name,
			Path: n.Path.Value,
		})
		// The logic here is to remove the content in the file from the package xxx to the import )
		// case *ast.GenDecl:
		// 	// IMPORT Declarations
		// 	if n.Tok == token.IMPORT {
		// 		// Add the new import
		// 		f.edit.Delete(f.fset.Position(0).Offset, f.fset.Position(n.End()).Offset)
		// 	}
	}
	return f
}

func ModifySrcBasedTarget(w io.Writer, filepath string, targetFunction map[string]struct{}) error {
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("ParsedSrc ioutil.ReadFile has error: %v", err)
	}
	parsedFile, err := parser.ParseFile(fset, filepath, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ParsedSrc parser.ParseFile has error: %v", err)
	}
	file := &SrcFileModify{
		fset:            fset,
		name:            filepath,
		content:         content,
		edit:            astutil.NewBuffer(content),
		astFile:         parsedFile,
		TargetFunctions: targetFunction,
	}
	ast.Walk(file, file.astFile)
	newContent := file.edit.Bytes()
	w.Write(newContent)
	return nil
}

func (f *SrcFileModify) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	// The logic here is to remove the content in the file from the package xxx to the import )
	case *ast.GenDecl:
		// IMPORT Declarations
		if n.Tok == token.IMPORT {
			// Add the new import
			f.edit.Delete(f.fset.Position(0).Offset, f.fset.Position(n.End()).Offset)
		}
	case *ast.FuncDecl:
		// rename the same test function of target file
		_, exist := f.TargetFunctions[n.Name.Name]
		if exist {
			newName := n.Name.Name + atghelper.RandStringBytes(6)
			f.edit.Replace(int(n.Name.NamePos-1), int(n.Name.End()-1), newName)
		}
	}
	return f
}

func AddedImportToTrc(w io.Writer, filepath string, srcImports, tarImports []*Import) (map[string]struct{}, error) {
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return map[string]struct{}{}, fmt.Errorf("[AddedImportToTrc] ReadFile has error %v", err)
	}
	parsedFile, err := parser.ParseFile(fset, filepath, content, parser.ParseComments)
	if err != nil {
		return map[string]struct{}{}, fmt.Errorf("[AddedImportToTrc] ParseFile has error %v", err)
	}
	file := &TrcFile{
		fset:          fset,
		name:          filepath,
		content:       content,
		edit:          astutil.NewBuffer(content),
		astFile:       parsedFile,
		imports:       srcImports,
		TrcImports:    tarImports,
		FunctionNames: make(map[string]struct{}, 0),
	}
	ast.Walk(file, file.astFile)
	if len(tarImports) == 0 {
		tarImports = srcImports
	} else {
		// Add the src imports to the target imports.
		// src refers to the temporary file which is generated by the smart unit
		// tar refers to the existing _nxt_unit_test.go files.
		for _, srcImport := range srcImports {
			shouldAdded := true
			for _, tarImport := range tarImports {
				if srcImport.Path == tarImport.Path {
					shouldAdded = false
				}
			}
			if shouldAdded {
				srcPath := strings.ReplaceAll(srcImport.Path, "\"", "")
				name := strings.ReplaceAll(srcImport.Name, "\"", "")
				if atghelper.GetPkgName(srcPath) == name {
					xastutil.AddImport(fset, parsedFile, srcPath)
				} else {
					xastutil.AddNamedImport(fset, parsedFile, name, srcPath)
				}
			}
		}
	}
	if file.isWrong {
		fmt.Println("[AddedImportToTrc] you cannot generate the test for the same function ! Please remove the existing test function and regenerate it")
		return file.FunctionNames, fmt.Errorf("[AddedImportToTrc] smart unit cannot generate the same test for the same function")
	}
	if err := printer.Fprint(w, fset, parsedFile); err != nil {
		fmt.Println("[AddedImportToTrc] internal error")
		return nil, fmt.Errorf("[AddedImportToTrc] cannot transfer the code")
	}
	return file.FunctionNames, nil
}

func (f *TrcFile) Visit(node ast.Node) ast.Visitor {
	switch dd := node.(type) {
	case *ast.FuncDecl:
		f.FunctionNames[dd.Name.Name] = struct{}{}
	}
	return f
}
