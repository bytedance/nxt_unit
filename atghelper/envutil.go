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

import (
	"bytes"
	"regexp"
	"strconv"

	"os/exec"

	"github.com/bytedance/nxt_unit/atgconstant"
)

type GoVersionInfo struct {
	FirstLevel  int
	SecondLevel int
	ThirdLevel  int
}

func GetGoVersion(path string) (GoVersionInfo, error) {
	var stdBuffer, stdErrBuff bytes.Buffer
	goVersionInfo := GoVersionInfo{}
	cmd := exec.Command(atgconstant.GoDirective, "version")
	cmd.Dir = path
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdErrBuff
	err := cmd.Run()
	if err != nil {
		return goVersionInfo, err
	}
	regexpStr := "go version go(\\d+).(\\d+).(\\d+)"
	// regexpStrBackup := "go version go(\\d+).(\\d+)"
	var match *regexp.Regexp
	match, err = regexp.Compile(regexpStr)
	if err != nil {
		match, err = regexp.Compile(regexpStr)
		if err != nil {
			return goVersionInfo, err
		}
	}
	versionInfo := match.FindStringSubmatch(stdBuffer.String())
	if len(versionInfo) > 0 {
		for index, _ := range versionInfo {

			switch index {
			case 1:
				firstLevel, _ := strconv.Atoi(versionInfo[index])
				goVersionInfo.FirstLevel = firstLevel
			case 2:
				secondLevel, _ := strconv.Atoi(versionInfo[index])
				goVersionInfo.SecondLevel = secondLevel
			case 3:
				thirdLevel, _ := strconv.Atoi(versionInfo[index])
				goVersionInfo.ThirdLevel = thirdLevel
			}
		}
	}
	return goVersionInfo, nil
}

func GetUseMockByVersion(path string) int {
	return atgconstant.UseMockitoMock
}
