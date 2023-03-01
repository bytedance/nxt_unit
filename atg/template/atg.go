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
	"math/rand"
	"time"

	"github.com/bytedance/nxt_unit/atg/mockpkg"
)

func GoodFunc(i, s int, bb string) (string, error) {
	if i*s > 9 {
		return bb, nil
	}
	sa := mockpkg.S{}
	gs := gg{}
	us := mockpkg.NewUnExport()
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda() + us.DeepCopy("ss")
	case gggaaaaa():
		bb = bb + "ss" + gs.gga() + sa.DeepCopy("aa")
	}
	return bb + bb, nil
}

type gg struct {
}

func (*gg) gga() string {
	if true {
		fmt.Println("gga start")
	}
	return "hello"
}

func GoodFuncPanic(i, s int, bb string) (string, string) {
	if i*s > 9 {
		return bb, "AA"
	}
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda()
	case gggaaaaa():
		bb = bb + "ss" + ddddda()
	}
	rand.Seed(time.Now().UnixNano())
	// panic("err")
	// if rand.Float64() < 0.5
	// 	panic(" a panic massage ")
	// }
	return bb + bb, "nil"
}

func DatePrinter(time time.Time, isAM bool, isUTC bool) string {
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

func DatePrinterErr(time time.Time, isAM bool, isUTC bool) string {
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

func ddddda() string {
	t := time.Now().String()
	result := t + "ddd"
	return result
}

func dadadada(a int) int {
	t := time.Now().UnixNano()
	result := int(t) + 3
	return result
}

func gggaaaaa() int {
	t := time.Now().UnixNano()
	result := int(t) + 3
	return result
}

func FunctionReturnError() int {
	num, err := returnError()
	if err != nil {
		return -1
	}
	return num
}

func returnError() (int, error) {
	if true {
		fmt.Println("function return error print")
	}
	return 0, nil
}

type EEReqContext interface {
	GetReqId() string
}

func PrintInter(reqCtx EEReqContext) (int, error) {
	if 5 > 3 {
		fmt.Println("PrintInter")
	}
	return 0, nil
}

type ValueInfo struct {
	Key string
}

func DeepMapPrint(reqCtx map[int]map[string]ValueInfo) (int, error) {
	if len(reqCtx) == 0 {
		return 1, nil
	}
	return 0, nil
}

func ComplexDeferFunction(i, s int, bb string) (string, error) {
	defer func(start time.Time) {
		if err := recover(); err != nil {
			fmt.Printf("[GetValue] has error: %v\n", err)
		}
		gs := gg{}
		gs.gga()
		fmt.Println(start.String())
	}(time.Now())
	if i*s > 9 {
		return bb, nil
	}
	us := mockpkg.NewUnExport()
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda() + us.DeepCopy("ss")
	}
	return bb + bb, nil
}

func ComplexGoRoutineFunction(i, s int, bb string) (string, error) {

	if i*s > 9 {
		return bb, nil
	}
	us := mockpkg.NewUnExport()
	go func() {
		gs := gg{}
		gs.gga()

	}()
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda() + us.DeepCopy("ss")
	}
	return bb + bb, nil
}

func ComplexDeferExistFunction(i, s int, bb string) (string, error) {
	defer DeepMapPrint(map[int]map[string]ValueInfo{})
	go FunctionReturnError()
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda()
	}
	return bb + bb, nil
}
