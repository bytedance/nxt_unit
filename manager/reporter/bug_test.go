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
	"fmt"
	"path"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
)

func TestBugReporter_Analysis(t *testing.T) {
	t.Skip("test on local")
	pathMatch, err := regexp.Compile(`# (.*?) \[`)
	if err != nil {
		fmt.Println("[Analysis Code err] :" + err.Error())
		return
	}
	match, err := regexp.Compile(`./(.*?):(.*?):(.*?): (.*?)xxxx`)
	if err != nil {
		fmt.Println("[Analysis Code err] :" + err.Error())
		return
	}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line := line + "xxxx"
		results := match.FindStringSubmatch(line)
		if len(results) != 5 {
			paths := pathMatch.FindStringSubmatch(line)
			if len(paths) == 2 {
				t.Log(paths[1])
			}
			continue
		}
		t.Log(results[1])
		t.Log(path.Join("ggg/dasdjh", results[1]))
	}
}

func TestBugReporter_getpath(t *testing.T) {
	var demoLine = `# /go/src/github/fakerinf/nxt_unit/atg/errinfo [/go/src/github/fakerinf/nxt_unit/atg/errinfo.test]`
	ok, paths := getPath(demoLine)
	if !ok {
		t.Fatal("can't no get path")
	}
	t.Log(paths)
}

var s = `# github/fakerinf/nxt_unit/atg/errinfo [github/fakerinf/nxt_unit/atg/errinfo.test]
./deadlock.go:18:6: syntax error: unexpected newline, expecting type
./deadlock.go:52:6: ddddda redeclared in this block
	previous declaration at ./deadlock.go:48:6`

var ciInfo = `# /go/src/github/fakerinf/nxt_unit/atg/errinfo [/go/src/github/fakerinf/nxt_unit/atg/errinfo.test]
./atg_test.go:18:6: syntax error: unexpected newline, expecting type
./atg_test.go:52:6: ddddda redeclared in this block
	previous declaration at ./deadlock.go:48:6`

func Test_bugReporter_Report(t *testing.T) {
	type fields struct {
		mode   string
		bugs   []bugInfo
		panics []bugInfo
	}
	type args struct {
		option atgconstant.Options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"case1",
			fields{
				Issue,
				[]bugInfo{
					{
						"1", "bug", "bug detail", "bugFunc", "warning", "/go/src/example.com/test/code/bug.go", "bug.go", ""},
				},
				[]bugInfo{
					{
						"1", "panic", "panic detail", "panicFunc", "error", "/go/src/example.com/test/code/panic.go", "panic.go", "",
					},
				},
			},
			args{
				atgconstant.Options{FilePath: "/go/src/example.com/test/code/panic.go", FuncName: "func"},
			},
			"::add-issue rule=errcheck,severity=warning,kind=bug,path=bug.go,line=1:: targetFunc-func || message-bug \n\n::add-message level=error:: targetFunc-func || panic **/go/src/example.com/test/code/panic.go** || panic info:  **panic** || stack: ```panic detail```  \n",
		},
		{
			"case2",
			fields{
				Message,
				[]bugInfo{
					{
						"1", "bug", "bug detail", "bugFunc", "warning", "/go/src/example.com/test/code/bug.go", "bug.go", ""},
				},
				[]bugInfo{
					{
						"1", "panic", "panic detail", "panicFunc", "error", "/go/src/example.com/test/code/panic.go", "panic.go", "",
					},
				},
			},
			args{
				atgconstant.Options{FilePath: "/go/src/example.com/test/code/panic.go", FuncName: "func"},
			},
			"::add-message level=error:: targetFunc-func || panic **/go/src/example.com/test/code/panic.go** || panic info:  **panic** || stack: ```panic detail```  \n",
		},
		{
			"case3",
			fields{
				Markdown,
				[]bugInfo{
					{
						"1", "bug", "bug detail", "bugFunc", "warning", "/go/src/example.com/test/code/bug.go", "bug.go", ""},
				},
				[]bugInfo{
					{
						"1", "panic", "panic detail", "panicFunc", "error", "/go/src/example.com/test/code/panic.go", "panic.go", "",
					},
				},
			},
			args{
				atgconstant.Options{FilePath: "/go/src/example.com/test/code/panic.go", FuncName: "func"},
			},
			"\n## example.com/test/code/panic.go:func bug report\n \n### /go/src/example.com/test/code/panic.go Panic\n\n- Panic Message : panic\n\n- Stack:\n\n \n ``` \n panic detail \n ``` \n \n\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bugReporter{
				RWMutex: sync.RWMutex{},
				mode:    tt.fields.mode,
				bugs:    tt.fields.bugs,
				panics:  tt.fields.panics,
			}
			if got := b.Report(tt.args.option); got != tt.want {
				t.Errorf("bugReporter.Report() = %v, want %v", got, tt.want)
			}
		})
	}
}
