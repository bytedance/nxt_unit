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
	"context"
	"encoding/json"
	"fmt"
)

type SpecialContext struct {
	ItemID int
}

type SpecialConsumption struct {
	ItemID int
}

type AdBigStruct struct {
	ID           int
	Add          string
	IsSuccess    bool
	Cvr          float64
	Ssr          float64
	InitPctr     float64
	DpaCid       int64
	DpaProductId string
	DpaDirectId  int64
}

func (*SpecialContext) Consume(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("test")
	}
	return nil
}

func (SpecialConsumption) Consume(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("test")
	}
	return nil
}

func (SpecialConsumption) CompareInteger(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetAbParamsMap(str string) map[string]interface{} {
	abParams := make(map[string]interface{})
	_ = json.Unmarshal([]byte(str), &abParams)
	return abParams
}

func (ad *AdBigStruct) CheckBigStruct(param int) int {
	if ad.Add != "" {
		return ad.ID
	}
	if !ad.IsSuccess {
		return 1
	}
	return 0
}

type mMT struct {
}

func (*mMT) GetSmartUnit(a int) string {
	fmt.Println("nono")
	return "ok"
}
