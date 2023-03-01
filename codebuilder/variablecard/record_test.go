package variablecard

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/smartunitvariablebuild"
	"reflect"
	"strings"
	"testing"

	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"github.com/stretchr/testify/assert"
)

func TestErrorToString(t *testing.T) {
	typeA := reflect.TypeOf(errors.New).Out(0)
	atgconstant.PkgRelativePath = "smart unit"
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	for i := 0; i < 2; i++ {
		duplicatepackagemanager.Init()
		mutatedV := VariableMutate(ctx, typeA, reflect.New(typeA).Elem())
		res := ValueToString(ctx, mutatedV)
		assert.Contains(t, res, "errors")
	}
}

func TestTestPointer(t *testing.T) {
	var a *atgconstant.Options
	typeA := reflect.TypeOf(a)
	atgconstant.PkgRelativePath = "smart unit"
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	for i := 0; i < 100; i++ {
		duplicatepackagemanager.Init()
		mutatedV := VariableMutate(ctx, typeA, reflect.New(typeA.Elem()))
		res := ValueToString(ctx, mutatedV)
		assert.Contains(t, res, "&atgconstant.Options")
	}
}

// Use Go redis to test the result
// func TestMapValueToString2(t *testing.T) {
//	var a *goredis.Option
//	typeA := reflect.TypeOf(a)
//	atgconstant.PkgRelativePath = "smart unit"
//	vtx := atgconstant.VariableContext{}
//	ctx := context.Background()
//	ctx = contexthelper.SetVariableContext(ctx, vtx)
//	mutatedV := VariableMutate(ctx, typeA.Elem(), reflect.New(typeA.Elem()).Elem())
//	res := ValueToString(ctx, mutatedV)
//	fmt.Println(res)
// }

type ParamFunction func(a int, b ...string)

func TestFunctionArrayParamValueToString(t *testing.T) {
	var a ParamFunction
	typeA := reflect.TypeOf(a)
	atgconstant.PkgRelativePath = "smart unit"
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	mutatedV := VariableMutate(ctx, typeA, reflect.Zero(typeA))
	res := ValueToString(ctx, mutatedV)
	fmt.Println(res)
	assert.Contains(t, res, "...")
}

// DO NOT DELETE
// Test Locally because we cannot introduce the errorcode_i18n in go.mod
// func TestErrorCode(t *testing.T) {
//	s := aweme.SUCCESS
//	vtx := atgconstant.VariableContext{}
//	ctx := context.Background()
//	ctx = contexthelper.SetVariableContext(ctx, vtx)
//	atgconstant.PkgRelativePath = "smart unit"
//	duplicatepackagemanager.Init()
//	mutatedV := VariableMutate(ctx, reflect.TypeOf(s), reflect.ValueOf(s))
//	res := ValueToString(ctx, mutatedV)
//	assert.Contains(t, res, "aweme")
// }

func TestContext(t *testing.T) {
	atgconstant.PkgRelativePath = "smart unit"
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	s := reflect.TypeOf(func() context.Context { return nil }).Out(0)
	mutatedV := VariableMutate(ctx, s, reflect.New(s).Elem())
	res := ValueToString(ctx, mutatedV)
	assert.Equal(t, res, "context.Background()")
}

func TestInt64Slice(t *testing.T) {
	atgconstant.PkgRelativePath = "smart unit"
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	for i := 0; i < 100; i++ {
		var a []int64
		duplicatepackagemanager.Init()
		typeA := reflect.TypeOf(a)
		valueA := reflect.ValueOf(a)
		mutatedV := VariableMutate(ctx, typeA, valueA)
		res := ValueToString(ctx, mutatedV)
		assert.Contains(t, res, "[]int64")
	}
}

type Guard9 struct {
	Hello string
}

func TestPtrSpecial9String(t *testing.T) {
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	ctx = duplicatepackagemanager.SetInstance(ctx)
	a := &Guard9{}
	duplicatepackagemanager.GetInstance(ctx).SetRelativePath(*a)
	duplicatepackagemanager.Init()
	duplicatepackagemanager.Init()
	typeA := reflect.TypeOf(a)
	valueA := reflect.New(typeA.Elem()).Elem()
	mutatedV := VariableMutate(ctx, typeA.Elem(), valueA).Addr()
	res := ValueToString(ctx, mutatedV)
	assert.Contains(t, res, "&Guard9{")
}

func TestStructValueToString(t *testing.T) {
	atgconstant.PkgRelativePath = "smart unit"
	type args struct {
		ctx        context.Context
		vInterface reflect.Value
	}
	structOfInt := struct {
		Num int
	}{}
	type NamedStruct struct {
		Num int
	}
	var structOfNamed NamedStruct
	structOfSlice := struct {
		Path []string
	}{}
	structOfInterface := struct {
		Path []error
		sa   []int
		aass map[string]int
	}{}

	type StructOfPrivateValue struct {
		Path         []error
		privateValue int
	}
	structOfPrivateValue := StructOfPrivateValue{}
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "StructOfInt",
			args: struct {
				ctx        context.Context
				vInterface reflect.Value
			}{ctx: ctx,
				vInterface: reflect.ValueOf(structOfInt),
			},
			want: "struct { Num int }",
		},
		{
			name: "StructOfNamed",
			args: struct {
				ctx        context.Context
				vInterface reflect.Value
			}{ctx: ctx,
				vInterface: reflect.ValueOf(structOfNamed),
			},
			want: "variablecard.NamedStruct",
		},
		{
			name: "StructOfSlice",
			args: struct {
				ctx        context.Context
				vInterface reflect.Value
			}{ctx: ctx,
				vInterface: reflect.ValueOf(structOfSlice),
			},
			want: "struct { Path []string }",
		},
		{
			name: "StructOfSlice",
			args: struct {
				ctx        context.Context
				vInterface reflect.Value
			}{ctx: ctx,
				vInterface: reflect.ValueOf(structOfSlice),
			},
			want: "struct { Path []string }",
		},
		{
			name: "StructOfInterface",
			args: struct {
				ctx        context.Context
				vInterface reflect.Value
			}{ctx: ctx,
				vInterface: reflect.ValueOf(structOfInterface),
			},
			want: "struct { ",
		},
		{
			name: "StructOfPrivate",
			args: struct {
				ctx        context.Context
				vInterface reflect.Value
			}{ctx: ctx,
				vInterface: reflect.ValueOf(structOfPrivateValue),
			},
			want: "variablecard.StructOfPrivateValue{",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValueToString(tt.args.ctx, tt.args.vInterface); !strings.Contains(got, tt.want) {
				t.Errorf("ValueToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

type IntegerType int64

func TestInt64PtrToString(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	ctx = duplicatepackagemanager.SetInstance(ctx)
	duplicatepackagemanager.GetInstance(ctx).SetRelativePath(IntegerType(9))
	for i := 0; i < 10; i++ {
		var tmp *IntegerType
		duplicatepackagemanager.Init()
		typeA := reflect.TypeOf(tmp)
		valueA := reflect.ValueOf(tmp)
		mutatedV := VariableMutate(ctx, typeA, valueA)
		res := ValueToString(ctx, mutatedV)
		assert.Contains(t, res, "func() *IntegerType")
	}
}

func TestJsonStringToString(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	ctx = duplicatepackagemanager.SetInstance(ctx)
	duplicatepackagemanager.GetInstance(ctx).SetRelativePath(IntegerType(9))
	duplicatepackagemanager.Init()
	s := reflect.ValueOf("{\"Comment\":0,\"Download\":0,\"Duet\":0,\"Stitch\":0,\"React\":0}")
	res := ValueToString(ctx, s)
	assert.Contains(t, res, "\\\"Comment\\\"")
}

func TestReNameType(t *testing.T) {
	var tmp *IntegerType
	var tmpo *int64
	typeA := reflect.TypeOf(tmp)
	typeO := reflect.TypeOf(tmpo)
	// valueA := reflect.ValueOf(tmp)
	fmt.Println(typeA.Kind().String())
	if typeA.Kind() == reflect.Ptr {
		fmt.Println(typeA.Elem().Kind().String())
		fmt.Println(typeA.Elem().Name())
		tmp = func() *IntegerType { tt := IntegerType(0); return &tt }()
	}
	fmt.Println(typeO.Kind().String())
	if typeO.Kind() == reflect.Ptr {
		fmt.Println(typeO.Elem().Kind().String())
		fmt.Println(typeO.Elem().Name())
	}
}

type ArrayStringType [10]string

func TestStringArrayToString(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	atgconstant.PkgRelativePath = "github.com/bytedance/nxt_unit/codebuilder/variablecard"
	for i := 0; i < 10; i++ {
		var tmp ArrayStringType
		duplicatepackagemanager.Init()
		typeA := reflect.TypeOf(tmp)
		valueA := reflect.ValueOf(tmp)
		mutatedV := VariableMutate(ctx, typeA, valueA)
		res := ValueToString(ctx, mutatedV)
		assert.Contains(t, res, "ArrayStringType{")
	}
}

func TestBasicPtrToString(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	atgconstant.PkgRelativePath = "github.com/bytedance/nxt_unit/codebuilder/variablecard"
	value0 := int(23)
	value1 := int8(23)
	value2 := int16(23)
	value3 := int32(3)
	value4 := int64(23)
	value5 := uint(3)
	value6 := uint8(3)
	value7 := uint16(3)
	value8 := uint32(3)
	value9 := uint64(3)
	value10 := float32(10.05)
	value11 := float64(10.05)
	value12 := bool(true)
	value13 := string("xfx")
	basicPtrList := []interface{}{&value0, &value1, &value2, &value3, &value4, &value5, &value6, &value7, &value8, &value9, &value10, &value11, &value12, &value13}
	for i := 0; i < 1; i++ {
		for index := range basicPtrList {
			duplicatepackagemanager.Init()
			typeA := reflect.TypeOf(basicPtrList[index])
			valueA := reflect.ValueOf(basicPtrList)
			mutatedV := VariableMutate(ctx, typeA, valueA)
			res := ValueToString(ctx, mutatedV)
			assert.Contains(t, res, "atgconv.")
		}
	}
}

func TestRemoveSelfImported(t *testing.T) {
	a := "121321"
	assert.Equal(t, removeSelfImported(a), "121321")
	b := "121321.456"
	assert.Equal(t, removeSelfImported(b), "456")
}

func Test_removeSelfImported(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// {name: "case1", args: args{s: "abc2.hello"}, want: "hello"},
		// {name: "case2", args: args{s: "*abc2.hello"}, want: "*hello"},
		{name: "case3", args: args{s: "[]test1.SmartAtg"}, want: "[]SmartAtg"},
		// {name: "case4", args: args{s: "*abc2.hello"}, want: "*hello"},
		// {name: "case5", args: args{s: "map[abc2]atg.hello"}, want: "map[abc2]atg.hello"},
		// {name: "case6", args: args{s: "hello"}, want: "hello"},
		// {name: "case7", args: args{s: "[]*pack.Comment"}, want: "[]*Comment"},
		// {name: "case8", args: args{s: "[][]smartunit.hello"}, want: "[][]hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeSelfImported(tt.args.s); got != tt.want {
				t.Errorf("removeSelfImported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeSelfImportedForMap(t *testing.T) {
	type args struct {
		s       string
		isKey   bool
		isValue bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{s: "map[abc.edf]su", isKey: true, isValue: false}, want: "map[edf]su"},
		{name: "case2", args: args{s: "map[abc.edf]hello.su", isKey: true, isValue: true}, want: "map[edf]su"},
		{name: "case3", args: args{s: "map[edf]hello.su", isKey: false, isValue: true}, want: "map[edf]su"},
		{name: "case4", args: args{s: "map[*abc.edf]su", isKey: true, isValue: false}, want: "map[*edf]su"},
		{name: "case5", args: args{s: "map[*abc.edf]*hello.su", isKey: true, isValue: true}, want: "map[*edf]*su"},
		{name: "case5", args: args{s: "map[*abc.edf]*hello.su", isKey: false, isValue: false}, want: "map[*abc.edf]*hello.su"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeSelfImportedForMap(tt.args.s, tt.args.isKey, tt.args.isValue); got != tt.want {
				t.Errorf("removeSelfImportedForMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueToStringFunction(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	a := func(a int) (int, bool) {
		return 9, true
	}
	assert.Contains(t, ValueToString(ctx, reflect.ValueOf(a)), "func(int) (int")
	b := func() {
		return
	}
	assert.Contains(t, ValueToString(ctx, reflect.ValueOf(b)), "func()")
}

func TestVariableLogBuilder(t *testing.T) {
	type Info struct {
		Logo string
	}
	a := Info{}
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	s := smartunitvariablebuild.NewSpecialValueInjector()
	ctx = context.WithValue(ctx, "SpecialValueInjector", s)
	s.SetBuilder(SetName(), "SetName()")
	mutateV := VariableMutate(ctx, reflect.TypeOf(a), reflect.ValueOf(a))
	t.Log(mutateV.Interface().(Info).Logo)
	t.Log(ValueToString(ctx, reflect.ValueOf(a)))
}

// test map value is still innerPkg map
func TestDeepMapTrimName(t *testing.T) {
	type args struct {
		mapType interface{}
	}
	type ValueInfo struct {
		Key string
	}
	var va = ValueInfo{}
	type Certificate struct {
		Raw                     []byte // Complete ASN.1 DER content (certificate, signature algorithm and signature).
		RawTBSCertificate       []byte // Certificate part of raw ASN.1 DER content.
		RawSubjectPublicKeyInfo []byte // DER encoded SubjectPublicKeyInfo.
		RawSubject              []byte // DER encoded Subject
		RawIssuer               []byte // DER encoded Issuer
	}
	type x509 struct {
		Certificate Certificate
	}
	type XZMap map[string]ValueInfo
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	duplicatepackagemanager.GetInstance(ctx).SetRelativePath(va)
	inputVariadic := &InputVariadic{false, false}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{map[int]map[ValueInfo]ValueInfo{}}, want: "map[int]map[ValueInfo]ValueInfo"},
		{name: "case2", args: args{map[int]map[string]ValueInfo{}}, want: "map[int]map[string]ValueInfo"},
		{name: "case3", args: args{map[ValueInfo]map[string]ValueInfo{}}, want: "map[ValueInfo]map[string]ValueInfo"},
		{name: "case4", args: args{map[ValueInfo]map[ValueInfo]ValueInfo{}}, want: "map[ValueInfo]map[ValueInfo]ValueInfo"},
		{name: "case5", args: args{map[string]ValueInfo{}}, want: "map[string]ValueInfo"},
		{name: "case6", args: args{XZMap{}}, want: "XZMap"},
		{name: "case7", args: args{map[string][]*ValueInfo{}}, want: "map[string][]*ValueInfo"},
		{name: "case8", args: args{map[string][]*atgconstant.ImportInfo{}}, want: "map[string][]*atgconstant.ImportInfo"},
		{name: "case9", args: args{map[ValueInfo][]atgconstant.ImportInfo{}}, want: "map[ValueInfo][]atgconstant.ImportInfo"},
		{name: "case10", args: args{map[atgconstant.ImportInfo][]atgconstant.ImportInfo{}}, want: "map[atgconstant.ImportInfo][]atgconstant.ImportInfo"},
		{name: "case11", args: args{map[atgconstant.ImportInfo][]*ValueInfo{}}, want: "map[atgconstant.ImportInfo][]*ValueInfo"},
		{name: "case12", args: args{map[ValueInfo]map[ValueInfo]atgconstant.ImportInfo{}}, want: "map[ValueInfo]map[ValueInfo]atgconstant.ImportInfo"},
		{name: "case13", args: args{map[atgconstant.ImportInfo]map[ValueInfo]atgconstant.ImportInfo{}}, want: "map[atgconstant.ImportInfo]map[ValueInfo]atgconstant.ImportInfo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimName(ctx, reflect.TypeOf(tt.args.mapType), inputVariadic)
			fmt.Println("result :" + result)
			assert.Contains(t, result, tt.want)
		})
	}
}

// test slice element is not basic type
func TestDeepTSlicerimName(t *testing.T) {
	type args struct {
		sliceType interface{}
	}
	type ValueInfo struct {
		Key string
	}
	var va = ValueInfo{}
	type Certificate struct {
		Raw               []byte // Complete ASN.1 DER content (certificate, signature algorithm and signature).
		RawTBSCertificate []byte // Certificate part of raw ASN.1 DER content.
	}
	type x509 struct {
		Certificate Certificate
	}
	vtx := atgconstant.VariableContext{}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	duplicatepackagemanager.GetInstance(ctx).SetRelativePath(va)
	inputVariadic := &InputVariadic{false, false}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{[][]*x509{}}, want: "[][]*x509"},
		{name: "case2", args: args{[][]*atgconstant.ImportInfo{}}, want: "[][]*atgconstant.ImportInfo"},
		{name: "case3", args: args{[]*atgconstant.ImportInfo{}}, want: "[]*atgconstant.ImportInfo"},
		{name: "case4", args: args{[]*x509{}}, want: "[]*x509"},
		{name: "case5", args: args{[]map[string]ValueInfo{}}, want: "[]map[string]ValueInfo"},
		{name: "case6", args: args{[]map[string]*ValueInfo{}}, want: "[]map[string]*ValueInfo"},
		{name: "case7", args: args{[]map[string]*atgconstant.ImportInfo{}}, want: "[]map[string]*atgconstant.ImportInfo"},
		{name: "case8", args: args{[]map[ValueInfo]*atgconstant.ImportInfo{}}, want: "[]map[ValueInfo]*atgconstant.ImportInfo"},
		{name: "case9", args: args{[]chan atgconstant.ImportInfo{}}, want: "[]chan atgconstant.ImportInfo"},
		{name: "case10", args: args{[]chan ValueInfo{}}, want: "[]chan ValueInfo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimName(ctx, reflect.TypeOf(tt.args.sliceType), inputVariadic)
			fmt.Println("slice result :" + result)
			assert.Contains(t, result, tt.want)
		})
	}
}
