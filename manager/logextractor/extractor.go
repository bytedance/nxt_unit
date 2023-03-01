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
package logextractor

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

const (
	ModTidyLogType int8 = 1
	// error message
	ModTidyErrMessage     string = "missing go.sum entry"
	ModTidySuggestMessage string = "Sorry,suggest you execute go mod tidy,and generate again"
)

type LogExtractor interface {
	GenSuggestByErrLog() error
}

func GenSuggestByErrLog(e error, stdOutput, stdError string) error {
	return fmt.Errorf("the error belongs to %w\n, the details is %v\n, the stdError is: %v\n, the stdOutput is: %v. ", GoModTidyError, e.Error(), stdError, stdOutput)
}

// TODO: handle more scenarios.
func CommandAnalyze(log string) {
	// Scan the string from the top to the end. Because the log is sequential, if there is the
	stringArray := strings.SplitN(log, "\n", -1)
	indexInit := 0
	indexPointer := 0
	noChangeInit := true
	noChangePointer := true
	for i := 0; i < len(stringArray); i++ {
		if strings.Contains(stringArray[i], "init(") || strings.Contains(stringArray[i], "init.") {
			indexInit = i
			noChangeInit = false
			break
		}
		if strings.Contains(stringArray[i], "invalid memory address or nil pointer dereference") {
			indexPointer = i
			noChangePointer = false
			break
		}
	}
	if noChangeInit && noChangePointer {
		return
	}
	if !noChangeInit {
		left := math.Max(0.0, float64(indexInit-10))
		right := math.Min(float64(len(stringArray)), float64(indexInit+30))
		result := strings.Join(stringArray[int(left):int(right)], "\n")
		ExecutionLog.LogError(DependencyInitError.Error() + "\n The error is caused by:\n" + result)
	}
	if !noChangePointer {
		left := math.Max(0.0, float64(indexPointer-10))
		right := math.Min(float64(len(stringArray)), float64(indexPointer+30))
		result := strings.Join(stringArray[int(left):int(right)], "\n")
		ExecutionLog.LogError(NullPointerError.Error() + "\n The error is caused by:\n" + result)
	}
	return
}

var ExecutionLog = executionLog{
	buffer: bytes.Buffer{},
}

// All the log will be stored in the executionLog, and finally we will print it.
type executionLog struct {
	buffer      bytes.Buffer
	errorBuffer bytes.Buffer
	finalRes    bytes.Buffer
	debugBuffer bytes.Buffer
}

func (c *executionLog) Log(s string) {
	c.buffer.WriteString(s + "\n")
}

func (c *executionLog) LogFinalRes(s string) {
	c.finalRes.WriteString(s + "\n")
}

// Only supprt one error type now.
func (c *executionLog) LogError(s string) {
	c.errorBuffer.WriteString(s + "\n")
}

func (c *executionLog) DebugInfo(s string) {
	c.debugBuffer.WriteString(s + "\n")
}

//  1. Firstly print debug information. Because debug imformation sometimes are pretty long, it will sequeeze all space in the terminal.
//  2. Secondly print error information.
//  3. Finally print the final result.
func (c *executionLog) Print() {
	fmt.Println(c.buffer.String())
	if c.debugBuffer.Len() > 0 {
		fmt.Println("######################## debug info ########################\n" + c.debugBuffer.String() + "############################################################\n")
	}
	if c.errorBuffer.Len() > 0 {
		fmt.Println("######################## error code ########################\n" + c.errorBuffer.String() + "############################################################\n")
	}
	fmt.Println(c.finalRes.String())
}
