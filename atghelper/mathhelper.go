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
	"math/rand"
	"time"
)

func RandomFloat(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	result := min + rand.Float64()*(max-min)
	return result
}

func RandomBool(rate float64) bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() < rate
}

func GetRandomFloat() float64 {
	return RandomFloat(0.0, 2.0)
}

var EPSILON = 0.00000001

func FloatEquals(a, b float64) bool {
	return (a-b) < EPSILON && (b-a) < EPSILON
}
