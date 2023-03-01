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

	redisv2 "github.com/bytedance/nxt_unit/atg/mockpkg/duplicateredis/redis"
	"github.com/bytedance/nxt_unit/atg/mockpkg/redis"
)

type UserInfo struct {
	Name    string
	Ctx     SpecialContext
	Bag     Bag
	ClientA redis.Client
	ClientB redisv2.Client
}

type Bag struct {
	P Pocket
	A int
	B int
}

type Pocket struct {
	Money  int
	UnUseA int32
	UnUseB int32
	UnUseC int32
}

func QueryData(u UserInfo) string {
	name := u.Name + u.ClientA.Address + u.ClientB.Address
	count := u.Bag.B + u.Bag.P.Money
	return fmt.Sprint(name, count)
}

type UserInfoPointer struct {
	Name    string
	Ctx     *SpecialContext
	Bag     *Bag
	ClientA *redis.Client
	ClientB *redisv2.Client
}

func QueryDataPointer(u *UserInfo) string {
	name := u.Name + u.ClientA.Address + u.ClientB.Address
	count := u.Bag.B + u.Bag.P.Money
	return fmt.Sprint(name, count)
}
