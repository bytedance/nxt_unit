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
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
)

const (
	panicType   = 1
	compileType = 2
)

func AnalyserPanic(contents string) (bug []bugInfo, err error) {
	errPlace := make(map[string]struct{})
	match, err := regexp.Compile(`panic:\(tested_func:(.*?)#(.*?)#(.*?)#(?is)(.*?)\)-c`)
	if err != nil {
		return bug, err
	}
	reports := match.FindAllStringSubmatch(contents, -1)
	if len(reports) <= 0 {
		return bug, errors.New("can't analysis panic")
	}
	for _, panicInfo := range reports {
		var oneBug bugInfo
		oneBug.targetFunc = panicInfo[1]
		oneBug.filePath = panicInfo[2]
		oneBug.message = panicInfo[3]
		oneBug.detail = "panic:" + oneBug.message + "\n" + panicInfo[4]
		oneBug.severity = "warning"
		oneBug.fileName = path.Base(oneBug.filePath)
		_, exist := errPlace[oneBug.filePath+oneBug.message]
		if exist {
			continue
		}
		errPlace[oneBug.filePath+oneBug.message] = struct{}{}
		bug = append(bug, oneBug)
	}
	return
}

func AnalyserInfo(contents, filePath, mirrorFileName string) ([]bugInfo, error, int) {
	p := panicErr{}
	bug, err := p.Analysis(context.Background(), contents, filePath, mirrorFileName)
	if err == nil {
		return bug, nil, panicType
	}
	c := compileErr{}
	bug, err = c.Analysis(context.Background(), contents, filePath)
	if err != nil {
		return bug, err, 0
	}
	return bug, nil, compileType
}

type Analyser interface {
	Analysis(ctx context.Context, contents, filePath string) bugInfo
}

type panicErr struct {
}

func (p *panicErr) Analysis(ctx context.Context, contents, filePath, mirrorFileName string) (bug []bugInfo, err error) {
	errPlace := make(map[string]struct{}, 0)
	match, err := regexp.Compile(`panic:\((?is)(.*?)#(?is)(.*?)\)-c`)
	if err != nil {
		return bug, err
	}
	reports := match.FindAllStringSubmatch(contents, -1)
	if len(reports) <= 0 {
		return bug, errors.New("can't analysis panic")
	}
	for _, panicInfo := range reports {
		var oneBug bugInfo
		oneBug.message = panicInfo[1]
		oneBug.detail = "panic:" + oneBug.message + "\n" + panicInfo[2]
		oneBug.severity = "warning"
		oneBug.filePath = getPackRelativePath(filePath)
		oneBug.MirrorFileName = mirrorFileName
		oneBug.fileName = path.Base(filePath)
		if !p.pass(oneBug.detail) {
			continue
		}

		if mirrorFileName != "" {
			lineNumReg := fmt.Sprintf("%s:([0-9]*) +", mirrorFileName)
			lineMatch, err := regexp.Compile(lineNumReg)
			if err != nil {
				fmt.Printf("line regexp mirrorFileName: %v,Err: %v\n", mirrorFileName, err)
			} else {
				bugLines := lineMatch.FindAllStringSubmatch(panicInfo[2], 1)
				for index, bugLine := range bugLines {
					if index == 0 {
						if len(bugLine) > 1 {
							oneBug.line = bugLine[1]
						}
						break
					}
				}
			}
			// file path and line unique
			if oneBug.line != "" {
				_, exist := errPlace[oneBug.filePath+mirrorFileName+oneBug.line]
				if exist {
					continue
				}
			}
			errPlace[oneBug.filePath+mirrorFileName+oneBug.line] = struct{}{}
		}
		bug = append(bug, oneBug)
	}
	return
}
func getPackRelativePath(filePath string) string {
	reg, err := regexp.Compile(`github\.com/([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)/(.*)`)
	if err != nil {
		fmt.Println(err)
	}
	result := reg.FindAllStringSubmatch(filePath, 1)
	if len(result) > 0 {
		if len(result[0]) == 4 {
			return result[0][3]
		}
	}
	return filePath
}

func (p *panicErr) pass(content string) bool {
	// if strings.Contains(content,"deadlock!") || strings.Contains(content,"nterface conversion!"){
	// 	return true
	// }
	if strings.Contains(content, "invalid memory address or nil pointer dereference") {
		return false
	}
	return true
}

type compileErr struct {
}

func (c *compileErr) Analysis(ctx context.Context, contents, filePath string) (bug []bugInfo, err error) {
	match, err := regexp.Compile(`./(.*?):(.*?):(.*?): (.*?)xxxx`)
	if err != nil {
		return bug, err
	}
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		line := line + "xxxx"
		results := match.FindStringSubmatch(line)
		if len(results) != 5 {
			ok, paths := getPath(line)
			if ok {
				filePath = paths
			}
			continue
		}
		file := results[1]
		lineNum := results[2]
		message := results[4]
		bug = append(bug, bugInfo{
			filePath: path.Join(path.Dir(filePath), file),
			message:  message,
			line:     lineNum,
			severity: "error",
			fileName: file,
		})
	}
	return bug, nil
}

func getPath(line string) (bool, string) {
	pathMatch, err := regexp.Compile(`# (.*?) \[`)
	if err != nil {
		fmt.Println("[Analysis Code err] :" + err.Error())
		return false, ""
	}
	paths := pathMatch.FindStringSubmatch(line)
	if len(paths) == 2 {
		source := paths[1]
		source = strings.Replace(source, "/go/src/", "", 1)
		return true, source
	}
	return false, ""
}
