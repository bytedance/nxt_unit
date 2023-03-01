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
	"sync"

	util "github.com/typa01/go-utils"

	"github.com/bytedance/nxt_unit/atgconstant"
)

var InternelErrorReporter = internalError{}

type internalError struct {
	sync.RWMutex
	errorNumber          int
	disableNumber        int
	errorFuncLocations   []functionLocation // it is caused by the internal error, like package / import
	disableFuncLocations []functionLocation // it is caused by complex struct or interface
}

type functionLocation struct {
	functionName string
	functionPath string
	errorInfo    string
}

func (c *internalError) Report() string {
	c.Lock()
	defer c.Unlock()
	defer func() {
		c.errorNumber = 0
		c.disableNumber = 0
		c.errorFuncLocations = make([]functionLocation, 0)
		c.disableFuncLocations = make([]functionLocation, 0)
	}()
	return c.addMessage()
}

func (c *internalError) AddErrorFunction(options atgconstant.Options, err error) {
	c.Lock()
	defer c.Unlock()
	c.errorNumber += 1
	c.errorFuncLocations = append(c.errorFuncLocations, functionLocation{
		functionPath: options.FilePath,
		functionName: options.FuncName,
		errorInfo:    fmt.Sprint(err),
	})
}

func (c *internalError) AddDisableFunction(options atgconstant.Options) {
	c.disableNumber += 1
	c.disableFuncLocations = append(c.disableFuncLocations, functionLocation{
		functionPath: options.FilePath,
		functionName: options.FuncName,
	})
}

func (c *internalError) addMessage() string {
	if len(c.errorFuncLocations) == 0 && len(c.disableFuncLocations) == 0 {
		return ""
	}
	builder := util.NewStringBuilder()
	builder.Append(fmt.Sprintf("errornumber(%v)-r\n", c.errorNumber))
	for _, v := range c.errorFuncLocations {
		builder.Append(fmt.Sprintf("internalerror(%s;%s;)-r\n", v.functionPath, v.functionName))
		builder.Append(fmt.Sprintf("internalerrorstring(%s;)-r\n", v.errorInfo))
	}
	builder.Append(fmt.Sprintf("disablenumber(%v)-r\n", c.disableNumber))
	for _, v := range c.disableFuncLocations {
		builder.Append(fmt.Sprintf("disableerror(%s;%s;)-r\n", v.functionPath, v.functionName))
	}
	return builder.ToString()
}
