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
package template

import (
	"fmt"

	"github.com/bytedance/nxt_unit/atg/mockpkg"
)

// DBInterface we can't get callee of call by interface
type DBInterface interface {
	Find() DBInterface
	Where() DBInterface
}

func Query(db DBInterface) error {
	type ss struct {
		arg struct{ DBInterface }
	}
	a := ss{}
	a.arg.Find()
	return fmt.Errorf("can,t mock ")
}

func DeepCopy(util mockpkg.Util) (err error) {
	util.DeepCopy("s")
	return nil
}

type A struct {
}

func (a *A) DeepCopyS(util mockpkg.Util) (err error) {
	util.DeepCopy("s")
	return nil
}
func StructToMap(value interface{}) map[interface{}]error {

	return map[interface{}]error{
		value: nil,
	}
}
