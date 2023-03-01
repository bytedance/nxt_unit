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
package atghelper

import "strings"

var PluginSDK Plugin

type Plugin interface {
	InstructionFile() string
	AtgTestFile() string
}

type pluginSDK struct {
	instructionFile string
	atgTestFile     string
}

func NewPluginSDK(instructionFile, atgTestFile string) Plugin {
	return &pluginSDK{
		instructionFile: instructionFile,
		atgTestFile:     atgTestFile,
	}
}

func (p *pluginSDK) InstructionFile() string {
	newName := strings.ReplaceAll(p.instructionFile, "smartunit.txt", ".go")
	return newName
}

func (p *pluginSDK) AtgTestFile() string {
	newName := strings.ReplaceAll(p.atgTestFile, "middle_code.txt", "middle_code_test.go")
	return newName
}
