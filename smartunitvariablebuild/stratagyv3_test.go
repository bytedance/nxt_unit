package smartunitvariablebuild

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"reflect"
	"testing"
)

func TestGetSpecialVariableV3(t *testing.T) {
	type args struct {
		ctx context.Context
		v   reflect.Type
	}
	variableContext := atgconstant.VariableContext{}
	ctx := context.Background()
	vtx := contexthelper.SetVariableContext(ctx, variableContext)
	// TODO: below usage is related with redis initialization.
	tests := []struct {
		name  string
		args  args
		want  reflect.Value
		want1 bool
	}{
		{name: "case1",
			args:  args{ctx: vtx, v: reflect.TypeOf(func() error { return nil }).Out(0)},
			want:  reflect.ValueOf(errors.New("smart unit")),
			want1: true,
		},
		{name: "case2",
			args: args{ctx: vtx, v: reflect.TypeOf(func() context.Context {
				return nil
			}).Out(0)},
			want: reflect.ValueOf(&struct {
				context.Context
			}{context.Background()}),
			want1: true,
		},
		{name: "case3",
			args:  args{ctx: vtx, v: reflect.TypeOf(nil)},
			want:  reflect.ValueOf(""),
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := GetSpecialVariableV3(tt.args.ctx, tt.args.v)
			if !ok {
				t.Fatal("mutate err")
			}
		})
	}
}

func TestNewErr(t *testing.T) {
	fmt.Println(reflect.TypeOf(context.Background()).Kind())
	a := func() context.Context {
		return nil
	}
	fmt.Println(reflect.TypeOf(a).Out(0).Kind())
}

func TestGetSpecialVariableV32(t *testing.T) {
	s := NewSpecialValueInjector()
	s.Set(reflect.TypeOf(status{}).String(), reflect.ValueOf(status{code: 800}))
	variableContext := atgconstant.VariableContext{}
	ctx := context.Background()
	vtx := contexthelper.SetVariableContext(ctx, variableContext)
	vtx = context.WithValue(vtx, "SpecialValueInjector", s)
	value, ok := GetSpecialVariableV3(vtx, reflect.TypeOf(status{}))
	if !ok {
		t.Fatal("mutate error")
	}
	if value.Interface().(status).code != 800 {
		t.Fatal("mutate value error")
	}
}
