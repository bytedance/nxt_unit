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
package lifemanager

var Closer = &closer{
	funcList: make([]func(), 0),
}

type closer struct {
	funcList []func()
}

func (c *closer) SetClose(f func()) {
	c.funcList = append(c.funcList, f)
}

func (c *closer) Close() {
	for _, closeFunc := range c.funcList {
		closeFunc()
	}
}

var SecondCloser = &closer{
	funcList: make([]func(), 0),
}
