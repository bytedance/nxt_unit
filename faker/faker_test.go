package faker

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
	A    error
	B    string
	C    *types.Var
	Skip string `faker:"-"`
	Mail string `faker:"email"`
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
	var a ErrorTest
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
	assert.Nil(t, err)
}

func TestGetValue11(t *testing.T) {
	var a *ErrorTest
	v, err := GetValue(reflect.TypeOf(a), 3)
	fmt.Println(v)
	assert.Nil(t, err)
}

func Test_isZero(t *testing.T) {
	type args struct {
		field reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"map",
			args{
				field: reflect.ValueOf(make(map[string]string)),
			},
			true,
			false,
		},
		{
			"int",
			args{
				field: reflect.ValueOf(0),
			},
			true,
			false,
		},
		{
			"error_struct",
			args{
				field: reflect.ValueOf(make([]int64, 0)),
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isZero(tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("isZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
