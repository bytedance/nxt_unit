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
	"reflect"
)

var address Addresser

// GetAddress returns a new Addresser interface of Address
func GetAddress() Addresser {
	if address == nil {
		address = &Address{}
	}
	return address
}

// SetAddress sets custom Address
func SetAddress(net Addresser) {
	address = net
}

// Addresser is logical layer for Address
type Addresser interface {
	Latitude(v reflect.Value) (interface{}, error)
	Longitude(v reflect.Value) (interface{}, error)
}

// Address struct
type Address struct{}

func (i Address) latitude() float32 {
	return (rand.Float32() * 180) - 90
}

// Latitude sets latitude of the address
func (i Address) Latitude(v reflect.Value) (interface{}, error) {
	kind := v.Kind()
	val := i.latitude()
	if kind == reflect.Float32 {
		return val, nil
	}
	return float64(val), nil
}

func (i Address) longitude() float32 {
	return (rand.Float32() * 360) - 180
}

// Longitude sets longitude of the address
func (i Address) Longitude(v reflect.Value) (interface{}, error) {
	kind := v.Kind()
	val := i.longitude()
	if kind == reflect.Float32 {
		return val, nil
	}
	return float64(val), nil
}

// Longitude get fake longitude randomly
func Longitude() float64 {
	address := Address{}
	return float64(address.longitude())
}

// Latitude get fake latitude randomly
func Latitude() float64 {
	address := Address{}
	return float64(address.latitude())
}
