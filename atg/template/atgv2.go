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
	"math/rand"
	"time"
)

func GoodFuncv2(i, s int, bb string) (string, error, string) {
	if i*s > 9 {
		return bb, nil, "nil"
	}
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda()
	case gggaaaaa():
		bb = bb + "ss" + ddddda()
	}
	return bb + bb, nil, "nil"
}

func GoodFuncPanicv2(i, s int, bb string) (string, error, string) {
	if i*s > 9 {
		return bb, nil, "nil"
	}
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda()
	case gggaaaaa():
		bb = bb + "ss" + ddddda()
	}
	rand.Seed(time.Now().UnixNano())
	panic("err")
	// if rand.Float64() < 0.5
	// 	panic(" a panic massage ")
	// }
	return bb + bb, nil, "nil"
}

func DatePrinterv2(time time.Time, isAM bool, isUTC bool) string {
	timeStr := time.Format("15:04:05")
	if isAM {
		timeStr += " AM"
	} else {
		timeStr += " PM"
	}
	if isUTC {
		timeStr += " UTC"
	}
	return timeStr
}

func DatePrinterErrv2(time time.Time, isAM bool, isUTC bool) string {
	timeStr := time.Format("15:04:05")
	if isAM {
		timeStr += " AM"
	} else {
		timeStr += " PM"
		return timeStr
	}
	if isUTC {
		timeStr += " UTC"
	}
	return timeStr
}

func dddddav2() string {
	t := time.Now().String()
	return t + "ddd"
}

func dadadadav2(a int) int {
	t := time.Now().UnixNano()
	return int(t) + 3
}

func gggaaaaav2() int {
	t := time.Now().UnixNano()
	return int(t) + 3
}
