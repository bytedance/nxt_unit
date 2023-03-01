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
package reporter

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/bytedance/nxt_unit/atgconstant"
)

const (
	Message  = "message"
	Issue    = "issue"
	Markdown = "markdown"
)

var BugReporter = bugReporter{
	bugs:   make([]bugInfo, 0),
	panics: make([]bugInfo, 0),
	mode:   Issue,
}

type bugReporter struct {
	sync.RWMutex
	mode   string
	bugs   []bugInfo
	panics []bugInfo
}

// ::add-issue rule=errcheck,severity=info,kind=bug,path=main.go,line=5,
// start_column=4::ineffectual assignment to `level`"

type bugInfo struct {
	line           string
	message        string
	detail         string
	targetFunc     string
	severity       string
	filePath       string
	fileName       string
	MirrorFileName string
}

func (b *bugReporter) Report(option atgconstant.Options) string {
	b.Lock()
	defer b.Unlock()
	b.trimBugs()
	switch b.mode {
	case Message:
		return b.addMessage(option)
	case Issue:
		return b.addIssue(option)
	case Markdown:
		return b.addMarkdown(option)
	default:
		return b.addIssue(option)
	}
}

func (b *bugReporter) Panics() int {
	return len(b.panics)
}

func (b *bugReporter) Bugs() int {
	return len(b.bugs)
}

func (b *bugReporter) ReportTemp(option atgconstant.Options) string {
	b.trimBugs()
	return b.addMarkdown(option)
}

// ::add-issue rule=errcheck,severity=info,kind=bug,path=main.go,line=5,start_column=4::ineffectual assignment to `level`
func (b *bugReporter) addIssue(option atgconstant.Options) string {
	var info string
	for _, bug := range b.bugs {
		if !option.ReportMode {
			if strings.Contains(bug.fileName, "_test") || strings.Contains(bug.fileName, "new.go") || strings.Contains(bug.fileName, "gotypeconverter") {
				continue
			}
		}
		// todo format filepath
		info = info + fmt.Sprintf("::add-issue rule=errcheck,severity=%s,kind=bug,path=%s,line=%s:: targetFunc-%s || message-%s \n", bug.severity,
			getIssuePath(bug.filePath), bug.line, option.FuncName, bug.message)
	}
	info = info + "\n"
	info = info + b.addMessage(option)
	return info
}

func getIssuePath(filePath string) string {
	filePath = strings.ReplaceAll(filePath, "/go/src/", "")
	dirs := strings.Split(filePath, "/")
	if len(dirs) <= 3 {
		return filePath
	}
	return path.Join(dirs[3:]...)
}

func (b *bugReporter) addMessage(option atgconstant.Options) string {
	var info string
	for _, bug := range b.panics {
		info = info + fmt.Sprintf("::add-message level=%s:: targetFunc-%s || %s \n", bug.severity, option.FuncName, markDown(bug))
	}
	return info
}

func (b *bugReporter) addMarkdown(option atgconstant.Options) string {
	var bugMessage string
	// panic info is no project's bug
	for _, bug := range b.panics {
		bugMessage = bugMessage + markdownPanic(bug)
	}
	filePath := option.FilePath
	filePath = strings.ReplaceAll(filePath, "/go/src/", "")
	// for _, bug := range b.bugs {
	// 	if !option.ReportMode {
	// 		if strings.Contains(bug.fileName, "_test") || strings.Contains(bug.fileName, "new.go") || strings.Contains(bug.fileName, "gotypeconverter") {
	// 			continue
	// 		}
	// 	}
	// 	bugMessage = bugMessage + markdownCompileErr(filePath, bug)
	// }
	if bugMessage == "" {
		return ""
	}
	info := fmt.Sprintf(`
## %s:%s bug report
 %s
`, filePath, option.FuncName, bugMessage)
	return info
}

func (b *bugReporter) GeneratePanicInfo() string {
	type issueDetail struct {
		LineNumber   int    `json:"line_number,omitempty"`
		ColumnNumber int    `json:"column_number,omitempty"`
		StackInfo    string `json:"stack_info,omitempty"`
	}
	type issue struct {
		IssueType   string      `json:"issue_type,omitempty"`
		FilePath    string      `json:"file_path,omitempty"`
		IssueTitle  string      `json:"issue_title,omitempty"`
		IssueDetail issueDetail `json:"issue_detail,omitempty"`
	}
	issues := make([]issue, 0)
	// panic info is no project's bug
	for _, bug := range b.panics {
		issues = append(issues, issue{
			IssueType:  "panic",
			FilePath:   bug.filePath,
			IssueTitle: bug.message,
			IssueDetail: issueDetail{
				// TODO: add line number
				StackInfo: bug.detail,
			},
		})
	}
	issue_byte, err := json.Marshal(issues)
	if err != nil {
		fmt.Println("JSON Marshal Err:", err.Error())
		return ""
	}
	return string(issue_byte)
}

func message(b bugInfo) string {
	return fmt.Sprintf("**%s**:**%s**: %s", b.filePath, b.line, b.message)
}

func markDown(b bugInfo) string {
	b.detail = strings.ReplaceAll(b.detail, "\n", "   ||  ")
	return fmt.Sprintf("panic **%s** ｜｜ panic info:  **%s** ｜｜ stack: ```%s``` ", b.filePath, b.message, b.detail)
}

func markdownPanic(b bugInfo) string {
	var panicInfo = `
### %s Panic

- Panic Message : %s

- Stack:

%s

`
	return fmt.Sprintf(panicInfo, b.filePath, b.message, " \n ``` \n "+b.detail+" \n ``` \n ")
}

func markdownCompileErr(funcName string, b bugInfo) string {
	filePath := b.filePath
	filePath = strings.ReplaceAll(filePath, "/go/src/", "")
	var errInfo = `
### %s ERR 

 - ERROR MESSAGE : %s

 - Code: %s:%s

`
	return fmt.Sprintf(errInfo, funcName, b.message, filePath, b.line)
}

func (b *bugReporter) trimBugs() {
	var left []bugInfo
	errMap := map[string]bool{}
	for _, bug := range b.bugs {
		ok, _ := errMap[bug.filePath+bug.line]
		if ok {
			continue
		}
		left = append(left, bug)
		errMap[bug.filePath+bug.line] = true
	}
	b.bugs = left
}

func (b *bugReporter) Analysis(contents string) {
	bug, err := AnalyserPanic(contents)
	if err != nil {
		return
	}
	b.panics = bug
}
