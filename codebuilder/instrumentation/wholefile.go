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
package instrumentation

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/nxt_unit/manager/logextractor"

	"github.com/bytedance/nxt_unit/codebuilder/astutil"
	xastutil "golang.org/x/tools/go/ast/astutil"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/google/uuid"
	util "github.com/typa01/go-utils"
)

func NewInstrumentation(file string, funcName, id string) ([]byte, int, error) {
	fset := token.NewFileSet()
	content, err := addNewImportToContent(file)
	if err != nil {
		return nil, 0, err
	}
	parsedFile, err := parser.ParseFile(fset, file, content, parser.ParseComments)
	if err != nil {
		return nil, 0, err
	}

	astAnalyzer := &CoverFile{
		fset:             fset,
		name:             file,
		content:          content,
		edit:             astutil.NewBuffer(content),
		astFile:          parsedFile,
		targetFunc:       funcName,
		relatedFunctions: []string{funcName},
		count:            1,
		lastInit:         false,
		resultLen:        3,
		Uid:              id,
	}
	ast.Walk(astAnalyzer, astAnalyzer.astFile)
	astAnalyzer.Rename()
	newContents := astAnalyzer.edit.String()
	newContents = strings.ReplaceAll(newContents, "astBranchTag", fmt.Sprintf("%v", len(astAnalyzer.branches)))
	newContents = strings.ReplaceAll(newContents, "\"astLineTag\"", fmt.Sprintf("%v", astAnalyzer.lines))
	newContents = newContents + astutil.ReflectStringV3(id)
	return []byte(newContents), astAnalyzer.lines, nil
}

// Add rename import to aviod duplicate with variable name
func addNewImportToContent(filepath string) ([]byte, error) {
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []byte{}, fmt.Errorf("AddImportToContent cannot parse File")
	}
	parsedFile, err := parser.ParseFile(fset, filepath, content, parser.ParseComments)
	if err != nil {
		return []byte{}, fmt.Errorf("AddImportToContent cannot parse File")
	}
	res := xastutil.AddNamedImport(fset, parsedFile, "smartunit_debug", "runtime/debug")
	res = res && xastutil.AddNamedImport(fset, parsedFile, "smartunit_strings", "strings")
	res = res && xastutil.AddNamedImport(fset, parsedFile, "smartunit_sort", "sort")
	res = res && xastutil.AddNamedImport(fset, parsedFile, "smartunit_runtime", "runtime")
	res = res && xastutil.AddNamedImport(fset, parsedFile, "smartunit_bytes", "bytes")
	if !res {
		return []byte{}, fmt.Errorf("AddImportToContent cannot add the import")
	}
	buffer := &bytes.Buffer{}
	if err := printer.Fprint(buffer, fset, parsedFile); err != nil {
		return nil, fmt.Errorf("AddImportToContent cannot transfer the code")
	}
	return buffer.Bytes(), nil
}

type Branch struct {
}

type CoverFile struct {
	fset             *token.FileSet
	name             string // Name of file.
	astFile          *ast.File
	content          []byte
	edit             *astutil.Buffer
	branches         []Branch
	defineCover      bool
	targetFunc       string
	lines            int
	relatedFunctions []string
	count            int
	lastInit         bool
	resultLen        int
	RepeatedName     bool
	Uid              string
	Pkg              string
}

// Visit implements the ast.Visitor interface.
// If the developer doesn't know how to write the code, please refer to
func (c *CoverFile) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BlockStmt:
		// If it's a switch or select, the body is a list of unittest clauses; don't tag the block itself.
		if len(n.List) > 0 {
			switch n.List[0].(type) {
			// 尝试确认行覆盖率
			default:
				for _, nLine := range n.List {
					clause := nLine.Pos()
					c.edit.Insert(c.offset(clause), c.commitTag())
				}
			case *ast.CaseClause: // switch
				for i, nLine := range n.List {
					clause := nLine.(*ast.CaseClause)
					// c.addBranchCounter(clause.Colon+1, clause.End())
					// 对其每一行进行追踪
					for _, line := range clause.Body {
						c.edit.Insert(c.offset(line.Pos()), c.commitTag())
					}
					c.edit.Insert(c.offset(clause.Colon+1), c.coverageTrace(fmt.Sprintf("%v - %v", len(n.List), i)))
				}
				return c
			case *ast.CommClause: // select
				for i, nLine := range n.List {
					clause := nLine.(*ast.CommClause)
					// 对其每一行进行追踪
					for _, line := range clause.Body {
						c.edit.Insert(c.offset(line.Pos()), c.commitTag())
					}
					// c.addBranchCounter(clause.Colon+1, clause.End())
					// branchPos := fmt.Sprintf("%v", clause.Colon+1)
					// c.edit.Insert(c.offset(clause.Colon+1), fmt.Sprint("fmt.Println(", "\"", "TrueDistances", "\"", ",", branchPos, ",", "0);"))
					// todo make it clear
					c.edit.Insert(c.offset(clause.Colon+1), c.coverageTrace(fmt.Sprintf("%v - %v", len(n.List), i)))
				}
				return c
			}
		}
	case *ast.IfStmt:
		if n.Init != nil {
			// c.lastInit = true
			ast.Walk(c, n.Init)
		}
		ast.Walk(c, n.Cond)
		// Below problem is used for the fitness function. However, in the smart unit 1.5 we don't
		// use the fitness function
		// We cannot process the if statement like this format:
		// 	if err := ajson.UnmarshalFromString(param.CommentInfo.Extra, &commentExtra); err == nil {
		//		syncCommenDigg(ctx, &commentExtra, param)
		//	}
		//  because the err variable cannot be access.
		// if c.lastInit {
		//	c.lastInit = false
		//	return nil
		// }
		// If statement InsertBranchDistanceStatement
		// insertPos := c.offset(n.If - 1)
		// branchPos := fmt.Sprintf("%v", n.Body.Lbrace+1)
		c.edit.Insert(c.offset(n.Body.Pos())+1, c.coverageTrace(time.Now().String()))
		// originExpr := GetOriginExpr(c, n)
		// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.TrueDistances, originExpr)
		// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.FalseDistances, []byte("!("+string(originExpr)+")"))
		// If statement addBranchCounter
		// c.addBranchCounter(n.Body.Lbrace+1, n.Body.Rbrace+1)
		// c.addBranchPredicateCounter(insertPos, branchPos)
		ast.Walk(c, n.Body)
		if n.Else == nil {
			return nil
		}
		switch n.Else.(type) {
		case *ast.BlockStmt:
			for _, line := range n.Else.(*ast.BlockStmt).List {
				c.edit.Insert(c.offset(line.Pos()), c.commitTag())
			}
			c.edit.Insert(c.offset(n.Else.Pos()+1), c.coverageTrace(time.Now().String()))
			return nil
		}
		// The elses are special, because if we have
		//	if x {
		//	} else if y {
		//	}
		// we want to cover the "if y". To do this, we need a place to drop the counter,
		// so we add a hidden block:
		//	if x {
		//	} else {
		//		if y {
		//		}
		//	}
		elseOffset := c.findText(n.Body.End(), "else")
		if elseOffset < 0 {
			panic("lost else")
		}
		c.edit.Insert(elseOffset+4, "{")
		c.edit.Insert(c.offset(n.Else.End()), "}")

		// We just created a block, now walk it.
		// Adjust the position of the new block to start after
		// the "else". That will cause it to follow the "{"
		// we inserted above.
		pos := c.fset.File(n.Body.End()).Pos(elseOffset + 4)
		switch stmt := n.Else.(type) {
		case *ast.IfStmt:
			block := &ast.BlockStmt{
				Lbrace: pos,
				List:   []ast.Stmt{stmt},
				Rbrace: stmt.End(),
			}
			n.Else = block
		case *ast.BlockStmt:
			// If else statement InsertBranchDistanceStatement
			// originExpr := GetOriginExpr(c, n)
			// branchPos := fmt.Sprintf("%v", pos)
			// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.TrueDistances, []byte("!("+string(originExpr)+")"))
			// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.FalseDistances, originExpr)
			// If else addBranchCounter
			// c.addBranchCounter(pos, n.Else.End())
			// c.addBranchPredicateCounter(insertPos, branchPos)
			stmt.Lbrace = pos
		default:
			panic("unexpected node type in if")
		}
		ast.Walk(c, n.Else)
		return nil
	case *ast.ReturnStmt:
		if len(n.Results) == 0 {
			return nil
		}

		// many to one relation.  It means it has multiple return signature but the return method only return one function
		if c.resultLen != len(n.Results) {
			// Temporally remove the logic, because it might invoke problem
			// _, ok := n.Results[0].(*ast.CallExpr)
			// if ok {
			// builder := util.NewStringBuilder()
			// for i := 0; i < c.resultLen; i++ {
			//	if i != 0 {
			//		builder.Append(",")
			//	}
			//	builder.Append(fmt.Sprint("smartUnitPluginRes", i))
			//	resultList = append(resultList, fmt.Sprint("isNil", c.Uid, "(", fmt.Sprint("smartUnitPluginRes", i), ") "))
			// }
			// originExpr := GetOriginExprV2(c, n)
			// res := decoratedReturnString(originExpr)
			// c.edit.Insert(insertPos-1, fmt.Sprint(";", builder.ToString(), ":=", res))
			// }
			return c
		}
	case *ast.SelectStmt:
		for _, v := range n.Body.List {
			switch v.(type) {
			case *ast.CommClause:
				// select statement
				// clause := v.(*ast.CommClause)
				// branchPos := fmt.Sprintf("%v", clause.Colon+1)
				// c.edit.Insert(c.offset(n.Pos()-1), atgconstant.TrueDistances+"["+branchPos+"]= 1;")
			}
		}
		// Don't annotate an empty select - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			return nil
		}
	case *ast.SwitchStmt:
		switchIdentExprAll := ""
		switchBinaryExprAll := ""
		for _, v := range n.Body.List {
			switchExpr := ""
			switchExprType := ""
			switch v.(type) {
			case *ast.CaseClause: // switch
				// switch : ast.Ident || ast.BinaryExpr
				switch n.Tag.(type) {
				case *ast.Ident:
					switchExpr = string(GetOriginExprV2(c, n.Tag))
					switchExprType = atgconstant.Ident
				case *ast.BinaryExpr:
					switchExpr = string(GetOriginExprV2(c, n.Tag))
					switchExprType = atgconstant.BinaryExpr
				}
				clause := v.(*ast.CaseClause)
				// insertPos := c.offset(n.Pos() - 1)
				for _, v := range clause.List {
					originExpr := GetOriginExprV2(c, v)
					if switchExprType == atgconstant.Ident {
						originExpr = []byte(switchExpr + " == " + string(originExpr))
						tempOriginExpr := "!(" + string(originExpr) + ")"
						if switchIdentExprAll == "" {
							switchIdentExprAll = tempOriginExpr
						} else {
							switchIdentExprAll = switchIdentExprAll + "&&" + tempOriginExpr
						}
					}
					if switchExprType == atgconstant.BinaryExpr {
						originExpr = []byte(switchExpr + " && " + string(originExpr))
						tempOriginExpr := "!(" + string(originExpr) + ")"
						if switchBinaryExprAll == "" {
							switchBinaryExprAll = tempOriginExpr
						} else {
							if !strings.Contains(switchBinaryExprAll, tempOriginExpr) {
								switchBinaryExprAll = switchBinaryExprAll + "&&" + tempOriginExpr
							}
						}
					}
					// branchPos := fmt.Sprintf("%v", clause.Colon+1)
					// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.TrueDistances, originExpr)
					// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.FalseDistances, []byte("!("+string(originExpr)+")"))
					// c.addBranchPredicateCounter(insertPos, branchPos)
				}
				// default statement
				// branchPos := fmt.Sprintf("%v", clause.Colon+1)
				if len(clause.List) == 0 {
					if switchIdentExprAll != "" {
						// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.TrueDistances, []byte(switchIdentExprAll))
						// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.FalseDistances, []byte("!("+switchIdentExprAll+")"))
						// c.addBranchPredicateCounter(insertPos, branchPos)
					}
					if switchBinaryExprAll != "" {
						if switchBinaryExprAll != "" {
							// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.TrueDistances, []byte(switchBinaryExprAll))
							// InsertBranchDistanceStatement(c, insertPos, branchPos, atgconstant.FalseDistances, []byte("!("+switchBinaryExprAll+")"))
							// c.addBranchPredicateCounter(insertPos, branchPos)
						}
					}
				}
			case *ast.CommClause: // select

			}
		}
		// Don't annotate an empty switch - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			if n.Init != nil {
				ast.Walk(c, n.Init)
			}
			if n.Tag != nil {
				ast.Walk(c, n.Tag)
			}
			return nil
		}
	case *ast.TypeSwitchStmt:
		// Don't annotate an empty type switch - creates a syntax error.
		if n.Body == nil || len(n.Body.List) == 0 {
			if n.Init != nil {
				ast.Walk(c, n.Init)
			}
			ast.Walk(c, n.Assign)
			return nil
		}
	case *ast.FuncLit:
		cwd, err := os.Getwd()
		relPath, err := filepath.Rel(cwd, c.name)
		if err != nil {
			relPath = c.name
		}
		c.addFuncCounter(n.Body.Lbrace+1, n.Body.Lbrace+2, "", relPath, "", "")
	case *ast.FuncDecl:
		// init  no need to instruction
		if n.Name.Name == "init" && n.Recv == nil {
			c.edit.Insert(c.offset(n.Body.Lbrace+1), "\n")
			c.edit.Insert(c.offset(n.Body.Lbrace+1), c.startInitCoverageTrace())
		} else {
			cwd, err := os.Getwd()
			relPath, err := filepath.Rel(cwd, c.name)
			if err != nil {
				relPath = c.name
			}
			// during instrumented file ,record tested function detail info
			isStart := ""
			recevierName := ""
			if n.Recv != nil && len(n.Recv.List) > 0 {
				switch expr := n.Recv.List[0].Type.(type) {
				case *ast.Ident:
					recevierName = expr.Name
				case *ast.StarExpr:
					isStart = "*"
					ident, ok1 := expr.X.(*ast.Ident)
					if ok1 {
						recevierName = ident.Name
					}
				}
			}
			c.addFuncCounter(n.Body.Lbrace+1, n.Body.Lbrace+2, n.Name.Name, relPath, isStart, recevierName)
		}
	}
	return c
}

// 放到最后
func (c *CoverFile) defineCoverageTrace() string {
	// 为了先满足语法树标准,astLineTag 在最终会被赋值为 行数量
	coverInfoV := fmt.Sprint("CoverInfoSU", c.Uid)
	workPipe := fmt.Sprint("WorkPipe", c.Uid)
	coverMap := fmt.Sprint("HitSet", c.Uid)
	linesV := fmt.Sprint("Lines", c.Uid)
	return fmt.Sprintf("\n const %v", fmt.Sprint(linesV, " = \"astLineTag\";var ", coverMap, "= [\"astLineTag\"]uint32{};", "type ", coverInfoV, " struct {PathID string;Coverage float64;ReturnString string;FunctionName string;Uid string;IsStart string;ReceiverName string};var ", workPipe, " = make(chan ", coverInfoV, ",10000);\n"))
}

func GetOriginExprV2(c *CoverFile, node ast.Node) []byte {
	start := node.Pos() - 1
	end := node.End() - 1
	switch n := node.(type) {
	case *ast.IfStmt:
		start = n.Cond.Pos() - 1
		end = n.Cond.End() - 1
	}
	return c.content[start:end]
}

// offset translates a token position into a 0-indexed byte offset.
func (c *CoverFile) offset(pos token.Pos) int {
	return c.fset.Position(pos).Offset
}

func (c *CoverFile) startCoverageTrace(funcName, fileName, isStart, receiverName string) string {
	branchVectorV := fmt.Sprint("branchVector")
	rV := fmt.Sprint("r", c.Uid)
	coverInfoV := fmt.Sprint("CoverInfoSU", c.Uid)
	coverInfoVStatment := fmt.Sprintf("CoverInfoSU%s{FunctionName: \"%s\", Uid: \"%s\",IsStart: \"%s\",ReceiverName: \"%s\"}", c.Uid, funcName, c.Uid, isStart, receiverName)
	workPipe := fmt.Sprint("WorkPipe", c.Uid)
	return fmt.Sprint(branchVectorV, " := map[string]int{};hitCommit := []func(){};var ", rV, " ", coverInfoV, " = ", coverInfoVStatment,
		";defer func() {if err := recover();err!= nil{", rV, ".Coverage = -1;",
		"fmt.Println(fmt.Sprintf(\"panic:(tested_func:"+funcName+"#"+fileName+"#%v#%v)-c\", err, string(smartunit_debug.Stack())));",
		"};", "var bvs []string;for k := range ", branchVectorV, " {bvs = append(bvs, k)};for _, h := range hitCommit{h()};smartunit_sort.Strings(bvs);",
		rV, ".PathID = smartunit_strings.Join(bvs,\"#\");", workPipe, " <- ", rV, "}();")
}

func (c *CoverFile) startInitCoverageTrace() string {
	branchVectorV := fmt.Sprint("branchVector")
	return fmt.Sprint(branchVectorV, " := map[string]int{};hitCommit := []func(){};fmt.Sprint(hitCommit,branchVector);")
}

func (c *CoverFile) coverageTrace(name string) string {
	id, err := uuid.NewUUID()
	if err != nil {
		name = strconv.FormatInt(time.Now().UnixNano(), 10)
	} else {
		name = id.String()
	}
	branchVectorV := fmt.Sprint("branchVector")
	builder := util.NewStringBuilder()
	builder.Append("if _,ok :=")
	builder.Append(branchVectorV)
	builder.Append("[\"")
	builder.Append(name)
	builder.Append("\"];ok{")
	builder.Append(branchVectorV)
	builder.Append("[\"")
	builder.Append(name)
	builder.Append("\"]++}else {")
	builder.Append(branchVectorV)
	builder.Append("[\"")
	builder.Append(name)
	builder.Append("\"]=1}")
	return builder.ToString()
}

func (c *CoverFile) commitTag() string {
	return fmt.Sprintf(`	hitCommit = append(hitCommit, func() {
		%v
	});`, c.tagLine())
}

func (c *CoverFile) tagLine() string {
	// 为了先满足语法树标准,astBranchTag 在最终会被赋值为 分支数量
	builder := util.NewStringBuilder()
	builder.Append(fmt.Sprintf("HitSet%s[%v]", c.Uid, c.lines))
	builder.Append(" ++;")
	c.lines++
	return builder.ToString()
}

// findText finds text in the original source, starting at pos.
// It correctly skips over comments and assumes it need not
// handle quoted strings.
// It returns a byte offset within f.src.
func (c *CoverFile) findText(pos token.Pos, text string) int {
	b := []byte(text)
	start := c.offset(pos)
	i := start
	s := c.content
	for i < len(s) {
		if bytes.HasPrefix(s[i:], b) {
			return i
		}
		if i+2 <= len(s) && s[i] == '/' && s[i+1] == '/' {
			for i < len(s) && s[i] != '\n' {
				i++
			}
			continue
		}
		if i+2 <= len(s) && s[i] == '/' && s[i+1] == '*' {
			for i += 2; ; i++ {
				if i+2 > len(s) {
					return 0
				}
				if s[i] == '*' && s[i+1] == '/' {
					i += 2
					break
				}
			}
			continue
		}
		i++
	}
	return -1
}

func (c *CoverFile) Rename() {
	for _, Del := range c.astFile.Decls {
		node, ok := Del.(ast.Node)
		if !ok {
			continue
		}
		switch n := node.(type) {
		case *ast.GenDecl:
			// rename the top del
			for _, del := range n.Specs {
				switch spec := del.(type) {
				case *ast.ValueSpec:
					for i, name := range spec.Names {
						if name.Name == "_" {
							continue
						}
						c.edit.Insert(c.offset(spec.Names[i].End()), c.Uid)
					}
				case *ast.TypeSpec:
					c.edit.Insert(c.offset(spec.Name.End()), c.Uid)
				}
			}
		case *ast.FuncLit:
			if !c.defineCover {
				c.edit.Insert(c.offset(n.Pos()), c.defineCoverageTrace())
				c.defineCover = true
			}
		case *ast.FuncDecl:
			// init also need record
			if n.Name.Name == "init" && n.Recv == nil {
				continue
			}
			if !c.defineCover {
				c.edit.Insert(c.offset(n.Pos()), c.defineCoverageTrace())
				c.defineCover = true
			}
			c.edit.Insert(c.offset(n.Name.End()), c.Uid)
		}
	}
}

func (c *CoverFile) addFuncCounter(insertPos, blockEnd token.Pos, funcName, fileName, isStart, receiverName string) {
	c.edit.Insert(c.offset(insertPos), "\n")
	c.edit.Insert(c.offset(insertPos), c.startCoverageTrace(funcName, fileName, isStart, receiverName))
	// f.edit.Insert(f.offset(insertPos), f.newFuncCounter(insertPos, blockEnd, name))
}

type functionCounter struct {
	fset      *token.FileSet
	astFile   *ast.File
	funcNames []string
}

func (f *functionCounter) Visit(node ast.Node) ast.Visitor {
	switch fn := node.(type) {
	case *ast.FuncDecl:
		// print actual function name
		if fn.Name.Name != "main" {
			f.funcNames = append(f.funcNames, fn.Name.Name)
		}
	}
	return f
}

func GetAllFunctionInFile(options atgconstant.Options) ([]string, error) {
	emptyArray := make([]string, 0)
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(options.FilePath)
	if err != nil {
		return emptyArray, fmt.Errorf("[getAllFunctionInFile] os cannot open file %v", err)
	}
	parsedFile, err := parser.ParseFile(fset, options.FilePath, content, parser.ParseComments)
	if err != nil {
		return emptyArray, fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.CannotParseTestedFunctionError, err.Error())
	}

	file := &functionCounter{
		fset:      fset,
		astFile:   parsedFile,
		funcNames: emptyArray,
	}
	ast.Walk(file, file.astFile)
	return file.funcNames, nil

}
