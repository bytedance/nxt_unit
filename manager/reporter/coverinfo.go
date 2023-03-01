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
	"regexp"
	"strconv"
	"time"

	"github.com/bytedance/nxt_unit/atgconstant"
)

var CoverReporter = coverReporter{
	start: time.Now(),
	//  testsuite-function-hit_line
	hitLine:        make(map[string]map[string]int),
	bestIndividual: []string{},
}

type coverReporter struct {
	mode      string
	success   int
	fail      int
	start     time.Time
	totalLine int
	fileTotal int
	// record the finalsuite
	bestIndividual []string
	// record every testsuite in hitLine
	hitLine map[string]map[string]int
}

func (c *coverReporter) SetBestIndividual(name string) {
	c.bestIndividual = append(c.bestIndividual, name)
}

func (c *coverReporter) AddTotal(fileLine int) {
	c.fileTotal = fileLine
}

func (c *coverReporter) Report() string {
	defer func() {
		c.success = 0
		c.fail = 0
		c.start = time.Now()
	}()
	switch c.mode {
	case Message:
		return c.addMessage()
	default:
		return c.addMessage()
	}
}

func (c *coverReporter) Analysis(content string, options atgconstant.Options, testsuiteName string) {
	// fmt.Println("覆盖率提取：",options.FuncName)
	match, err := regexp.Compile(`funcCover\((.*?);(.*?);(.*?)\)-r`)
	if err != nil {
		return
	}
	cover := match.FindStringSubmatch(content)
	if len(cover) != 4 {
		return
	}
	total, err := strconv.Atoi(cover[2])
	if err != nil {
		return
	}
	hit, err := strconv.Atoi(cover[3])
	if err != nil {
		return
	}
	c.totalLine = c.totalLine + total
	_, ok := c.hitLine[testsuiteName]
	// init testsuiteCover map
	if !ok {
		c.hitLine[testsuiteName] = map[string]int{}
	}
	count, ok := c.hitLine[testsuiteName][options.FuncName]
	if !ok {
		if c.hitLine[testsuiteName] == nil {
			c.hitLine[testsuiteName] = map[string]int{}
		}
		c.hitLine[testsuiteName][options.FuncName] = hit
		return
	}
	// update hits in ga mode
	if count < hit {
		c.hitLine[testsuiteName][options.FuncName] = hit
	}
}

func (c *coverReporter) RecordGoodTestCase(v int) {
	c.success = c.success + v
}

func (c *coverReporter) RecordFailTestCase(v int) {
	c.fail = c.fail + v
}

func (c *coverReporter) stop() string {
	rt := time.Since(c.start)
	return strconv.FormatFloat(rt.Seconds(), 'f', -1, 64)
}

func (c *coverReporter) addIssue() string {
	return ""
}

func (c *coverReporter) hitCount() int {
	var total int
	for _, bestTestSuite := range c.bestIndividual {
		bestCoverage, ok := c.hitLine[bestTestSuite]
		if !ok {
			continue
		}
		for _, hit := range bestCoverage {
			total = total + hit
		}
	}
	return total
}

func (c *coverReporter) addMessage() string {
	realTotal := c.totalLine
	if c.fileTotal > 0 {
		realTotal = c.fileTotal
	}
	coverage := fmt.Sprintf("coverage(%v;%v)-r \n", realTotal, c.hitCount())
	coverage = coverage + fmt.Sprintf("coverresult(%v;%v;%s)-r", c.success, c.fail, c.stop())
	return coverage
}
