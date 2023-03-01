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
	"io"
	"reflect"
)

var identifier Identifier

// GetIdentifier returns a new Identifier
func GetIdentifier() Identifier {
	if identifier == nil {
		identifier = &UUID{}
	}
	return identifier
}

// Identifier ...
type Identifier interface {
	Digit(v reflect.Value) (interface{}, error)
	Hyphenated(v reflect.Value) (interface{}, error)
}

// UUID struct
type UUID struct{}

// createUUID returns a 16 byte slice with random values
func createUUID() ([]byte, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(crypto, b)
	if err != nil {
		return b, err
	}
	// variant bits; see section 4.1.1
	b[8] = b[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	b[6] = b[6]&^0xf0 | 0x40
	return b, nil
}

func (u UUID) hyphenated() (string, error) {
	b, err := createUUID()
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, err
}

// Hyphenated returns a 36 byte hyphenated UUID
func (u UUID) Hyphenated(v reflect.Value) (interface{}, error) {
	return u.hyphenated()
}

// UUIDHyphenated get fake Hyphenated UUID
func UUIDHyphenated() string {
	u := UUID{}
	res, _ := u.hyphenated()
	return res
}

func (u UUID) digit() (string, error) {
	b, err := createUUID()
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x", b)
	return uuid, err
}

// Digit returns a 32 bytes UUID
func (u UUID) Digit(v reflect.Value) (interface{}, error) {
	return u.digit()
}

// UUIDDigit get fake Digit UUID
func UUIDDigit() string {
	u := UUID{}
	res, _ := u.digit()
	return res
}
