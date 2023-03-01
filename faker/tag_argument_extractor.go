// MIT License
//
// Copyright (c) 2017 Iman Tumorang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package faker

import (
	"fmt"
	"strconv"
	"strings"
)

func extractFloat64FromTagArgs(args []string) (interface{}, error) {
	bytes := 64
	var floatValues []float64
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseFloat(k, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, j)
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractFloat32FromTagArgs(args []string) (interface{}, error) {
	bytes := 32
	var floatValues []float32
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseFloat(k, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, float32(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt64FromTagArgs(args []string) (interface{}, error) {
	bytes := 64
	var floatValues []int64
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, j)
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt32FromTagArgs(args []string) (interface{}, error) {
	bytes := 32
	var floatValues []int32
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int32(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt16FromTagArgs(args []string) (interface{}, error) {
	bytes := 16
	var floatValues []int16
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int16(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt8FromTagArgs(args []string) (interface{}, error) {
	bytes := 8
	var floatValues []int8
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int8(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractIntFromTagArgs(args []string) (interface{}, error) {
	bytes := 0
	var floatValues []int
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint64FromTagArgs(args []string) (interface{}, error) {
	bytes := 64
	var floatValues []uint64
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, j)
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint32FromTagArgs(args []string) (interface{}, error) {
	bytes := 32
	var floatValues []uint32
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint32(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint16FromTagArgs(args []string) (interface{}, error) {
	bytes := 16
	var floatValues []uint16
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint16(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint8FromTagArgs(args []string) (interface{}, error) {
	bytes := 8
	var floatValues []uint8
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint8(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUintFromTagArgs(args []string) (interface{}, error) {
	bytes := 0
	var floatValues []uint
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}
