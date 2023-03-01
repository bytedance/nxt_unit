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
package stratagy

import (
	"fmt"
	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
)

type SpecialVariable struct {
	Type      string
	PkgName   string
	PkgPath   string
	Mutated   bool // we currently only support two types of mutation
	IsPointer bool
	CanBeNil  bool
}

func GetSpecialVariable(newV *parsermodel.Field, isPointer bool, canBeNil bool) (*SpecialVariable, bool) {
	if newV.Type == nil {
		fmt.Printf("[GetSpecialVariable] extracinfo.ConvertToModelField's type is nil, the name is %v\n", newV.Name)
		return nil, false
	}
	return nil, false
}

func (s *SpecialVariable) GetConstructor() string {
	switch s.Type {
	case "*redis.cmdable":
		return redisCmdableConstructor()
	default:
		fmt.Println("[GetConstructor] has error, the special variable type is not correct")
		return ""
	}
}

func redisCmdableConstructor() string {
	return "*goredis.Client"
}
