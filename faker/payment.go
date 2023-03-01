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
	"strconv"
	"strings"
)

// creditCard struct
type creditCard struct {
	ccType   string
	length   int
	prefixes []int
}

var creditCards = map[string]creditCard{
	"visa":             {"VISA", 16, []int{4539, 4556, 4916, 4532, 4929, 40240071, 4485, 4716, 4}},
	"mastercard":       {"MasterCard", 16, []int{51, 52, 53, 54, 55}},
	"american express": {"American Express", 15, []int{34, 37}},
	"discover":         {"Discover", 16, []int{6011}},
	"jcb":              {"JCB", 16, []int{3528, 3538, 3548, 3558, 3568, 3578, 3588}},
	"diners club":      {"Diners Club", 14, []int{36, 38, 39}},
}

var pay Render

var cacheCreditCard string

// GetPayment returns a new Render interface of Payment struct
func GetPayment() Render {
	if pay == nil {
		pay = &Payment{}
	}
	return pay
}

// SetPayment set custom Network
func SetPayment(p Render) {
	pay = p
}

// Render contains Whole Random Credit Card Generators with their types
type Render interface {
	CreditCardType(v reflect.Value) (interface{}, error)
	CreditCardNumber(v reflect.Value) (interface{}, error)
}

// Payment struct
type Payment struct{}

func (p Payment) cctype() string {
	n := len(creditCards)
	if cacheCreditCard != "" {
		return cacheCreditCard
	}
	var ccTypes []string

	for _, cc := range creditCards {
		ccTypes = append(ccTypes, cc.ccType)
	}
	cacheCreditCard = ccTypes[rand.Intn(n)]
	return cacheCreditCard
}

// CreditCardType returns one of the following credit values:
// VISA, MasterCard, American Express, Discover, JCB and Diners Club
func (p Payment) CreditCardType(v reflect.Value) (interface{}, error) {
	return p.cctype(), nil
}

// CCType get a credit card type randomly in string (VISA, MasterCard, etc)
func CCType() string {
	p := Payment{}
	return p.cctype()
}

func (p Payment) ccnumber() string {
	ccType := p.cctype()
	cacheCreditCard = ccType
	card := creditCards[strings.ToLower(ccType)]
	prefix := strconv.Itoa(card.prefixes[rand.Intn(len(card.prefixes))])

	num := prefix
	digit := randomStringNumber(card.length - len(prefix))

	num += digit
	return num
}

// CreditCardNumber generated credit card number according to the card number rules
func (p Payment) CreditCardNumber(v reflect.Value) (interface{}, error) {
	return p.ccnumber(), nil
}

// CCNumber get a credit card number randomly in string (VISA, MasterCard, etc)
func CCNumber() string {
	p := Payment{}
	return p.ccnumber()
}
