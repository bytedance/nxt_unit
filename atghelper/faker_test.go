package atghelper

import (
	"context"
	"fmt"
	"go/types"
	"reflect"
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/stretchr/testify/assert"
)

type ErrorTest struct {
	A error
	b string
	C *types.Var
}

func TestGetValue1(t *testing.T) {
	var a *atgconstant.Options
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v)
	assert.NotNil(t, err)
}

func TestGetValue2(t *testing.T) {
	var a *atgconstant.ExecutionValues
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v)
	assert.Nil(t, err)
}

func TestGetValue3(t *testing.T) {
	var a *map[string]*atgconstant.ImportInfo
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v)
	assert.Nil(t, err)
}

func TestGetValue4(t *testing.T) {
	var a interface{}
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v)
	assert.NotNil(t, err)
}

func TestGetValue5(t *testing.T) {
	var a context.Context
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v)
	assert.NotNil(t, err)
}

// This one is correct, but because we don't import the tiktok combo so we just left it here.
//func TestGetValue6(t *testing.T) {
//	var a *kernel.RequestContext
//	v, err := GetValue(reflect.TypeOf(a), 0)
//	fmt.Println(v.String())
//	assert.Nil(t, err)
//}

func TestGetValue7(t *testing.T) {
	var a *types.Var
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v.Elem())
	assert.Nil(t, err)
}

func TestGetValue8(t *testing.T) {
	var a *ErrorTest
	v, err := GetValue(reflect.TypeOf(a), 1)
	fmt.Println(v)
	assert.Nil(t, err)
}

func TestGetValue9(t *testing.T) {
	var a *[][]string
	_, err := GetValue(reflect.TypeOf(a), 0)
	assert.Nil(t, err)
}

func TestGetValue10(t *testing.T) {
	var a *ErrorTest
	v, err := GetValue(reflect.TypeOf(a), 0)
	fmt.Println(v)
	assert.NotNil(t, err)
}

func TestGetValue11(t *testing.T) {
	var a *ErrorTest
	v, err := GetValue(reflect.TypeOf(a), 3)
	fmt.Println(v)
	assert.Nil(t, err)
}
