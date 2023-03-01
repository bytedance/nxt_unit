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
	"testing"
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

var testInfo = `# github/fakerinf/nxt_unit/atg/errinfo [github.com/bytedance/nxt_unit/atg/errinfo.test]
./atg_test.go:18:6: syntax error: unexpected newline, expecting type
./atg_test.go:52:6: ddddda redeclared in this block
	previous declaration at ./deadlock.go:48:6`

var ciInfo = `# /go/src/github/fakerinf/nxt_unit/atg/errinfo [/go/src/github/fakerinf/nxt_unit/atg/errinfo.test]
./atg_test.go:18:6: syntax error: unexpected newline, expecting type
./atg_test.go:52:6: ddddda redeclared in this block
	previous declaration at ./deadlock.go:48:6`
