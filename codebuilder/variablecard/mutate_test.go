package variablecard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/smartunitvariablebuild"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StructTest1 struct {
	A string
	B int
}

func TestVariableMutate(t *testing.T) {
	testStruct := struct {
		A int
	}{1}
	type IntAlias int
	structWithTypeAlias := struct {
		A IntAlias
	}{1}
	structslice := make([]StructTest1, 0)
	structOfSlice := struct {
		Slice []string
	}{}
	structOfMap := struct {
		Map map[string]string
	}{}
	structOfInterface := struct {
		Map map[string]interface{}
	}{}
	type args struct {
		t           reflect.Type
		v           reflect.Value
		LiteralList []reflect.Value
	}
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)

	tests := []struct {
		name string
		args args
	}{
		{
			"Struct",
			args{
				reflect.TypeOf(testStruct),
				reflect.ValueOf(testStruct),
				nil,
			},
		},
		{
			"PtrToStruct",
			args{
				reflect.TypeOf(&testStruct),
				reflect.ValueOf(&testStruct),
				nil,
			},
		},
		{
			"StructWithTypeAlias",
			args{
				reflect.TypeOf(&structWithTypeAlias),
				reflect.ValueOf(&structWithTypeAlias),
				nil,
			},
		},
		{
			"StructOfSlice",
			args{
				reflect.TypeOf(structOfSlice),
				reflect.ValueOf(structOfSlice),
				nil,
			},
		},
		{
			"StructOfMap",
			args{
				reflect.TypeOf(structOfMap),
				reflect.ValueOf(structOfMap),
				nil,
			},
		},
		{
			"StructOfInterface",
			args{
				reflect.TypeOf(structOfInterface),
				reflect.ValueOf(structOfInterface),
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VariableMutate(ctx, tt.args.t, tt.args.v); reflect.DeepEqual(got, tt.args.v) {
				t.Errorf("No Mutation: VariableMutate() = %v, want %v", got, tt.args.v)
			}
		})
	}

	res := VariableMutate(ctx, reflect.TypeOf(structslice), reflect.ValueOf(structslice))
	assert.NotNil(t, res)
}

func TestMutatePtr(t *testing.T) {
	var guess *int
	guess = new(int)
	*guess = 1
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	res := VariableMutate(ctx, reflect.TypeOf(guess), reflect.ValueOf(guess))
	assert.Equal(t, res.Type(), reflect.TypeOf(guess))
}

func TestSpecialValue(t *testing.T) {
	guess := errors.New("hello")
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	res := VariableMutate(ctx, reflect.TypeOf(guess), reflect.ValueOf(guess))
	assert.Equal(t, res.Type(), reflect.TypeOf(guess))
}

func TestString(t *testing.T) {
	var a string
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	typeA := reflect.TypeOf(a)
	res := VariableMutate(ctx, reflect.TypeOf(a), reflect.ValueOf(typeA))
	assert.Equal(t, res.Type(), typeA)
}

func TestMutateComplex64(t *testing.T) {
	var guess complex64
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	guess = 1 + 5i
	res := VariableMutate(ctx, reflect.TypeOf(guess), reflect.ValueOf(guess))
	assert.Equal(t, res.Type(), reflect.TypeOf(guess))
}

func TestMutateComplex128(t *testing.T) {
	var guess complex128
	guess = 1 + 5i
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	res := VariableMutate(ctx, reflect.TypeOf(guess), reflect.ValueOf(guess))
	assert.Equal(t, res.Type(), reflect.TypeOf(guess))
}

func TestPlayground(t *testing.T) {
	type Bug int
	type test struct {
		A Bug
	}
	testStruct := test{1}
	type args struct {
		t           reflect.Type
		v           reflect.Value
		LiteralList []reflect.Value
	}
	vtx := atgconstant.VariableContext{
		Level:    0,
		ID:       0,
		CanBeNil: false,
	}
	ctx := context.Background()
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	tests := []struct {
		name string
		args args
	}{
		{
			"PtrToStruct",
			args{
				reflect.TypeOf(&testStruct),
				reflect.ValueOf(&testStruct),
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VariableMutate(ctx, tt.args.t, tt.args.v); reflect.DeepEqual(got, tt.args.v) {
				t.Errorf("No Mutation: VariableMutate() = %v, want %v", got, tt.args.v)
			}
		})
	}
}

// Test int,int8,int16,int32,int64
func TestMutliIntMutated(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	intList := []interface{}{int(23), int8(23), int16(23), int32(23), int64(23)}
	for _, tmp := range intList {
		mutateV := VariableMutate(ctx, reflect.TypeOf(tmp), reflect.ValueOf(tmp))
		assert.Equal(t, reflect.TypeOf(tmp), mutateV.Type())
	}
}

// Test Int Alias variable
func TestIntAliasMutation(t *testing.T) {
	type StrongInt int64
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	si := StrongInt(0)
	mutateV := VariableMutate(ctx, reflect.TypeOf(si), reflect.ValueOf(si))
	assert.Equal(t, reflect.TypeOf(si), mutateV.Type())
}

// Test String Alias variable
func TestStringAliasMutation(t *testing.T) {
	type StrongString string
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	si := StrongString("hello")
	mutateV := VariableMutate(ctx, reflect.TypeOf(si), reflect.ValueOf(si))
	assert.Equal(t, reflect.TypeOf(si), mutateV.Type())
}

// Test uint,uint8,uint16,uint32,uint64
func TestMutliUIntMutated(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	uintList := []interface{}{uint(23), uint8(23), uint16(23), uint32(23), uint64(23)}
	for _, tmp := range uintList {
		mutateV := VariableMutate(ctx, reflect.TypeOf(tmp), reflect.ValueOf(tmp))
		assert.Equal(t, reflect.TypeOf(tmp), mutateV.Type())
	}
}

// Test float32,float64,bool
func TestMutliFloatMutated(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	arrayList := []interface{}{float32(10.05), float64(10.05), bool(true)}
	for _, tmp := range arrayList {
		mutateV := VariableMutate(ctx, reflect.TypeOf(tmp), reflect.ValueOf(tmp))
		assert.Equal(t, reflect.TypeOf(tmp), mutateV.Type())
	}
}

type DB struct {
	Config string
	Map    map[string]int
}

func (d *DB) Find() *DB {
	return d
}

func (d *DB) Where() *DB {
	return d
}

// Test Pointer of struct
func TestPointerMutated(t *testing.T) {
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	newDB := &DB{}
	mutateV := VariableMutate(ctx, reflect.TypeOf(newDB), reflect.ValueOf(newDB))
	fmt.Println(ValueToString(ctx, reflect.ValueOf(mutateV.Interface())))
}

func TestMutateFakeData(t *testing.T) {
	// We have restriction for the struct field, thus we only keep 10's field here.
	type StructWithoutTags struct {
		// Latitude           float32
		// Longitude          float32
		CreditCardNumber string
		// CreditCardType     string
		Email string
		// DomainName         string
		IPV4 string
		// IPV6               string
		Password string
		// Jwt                string
		PhoneNumber string
		// MacAddress         string
		// URL                string
		UserName string
		// TollFreeNumber     string
		// E164PhoneNumber    string
		// TitleMale          string
		// TitleFemale        string
		// FirstName          string
		// FirstNameMale      string
		// FirstNameFemale    string
		// LastName           string
		// Name               string
		// UnixTime           int64
		Date string
		// Time               string
		// MonthName          string
		// Year               string
		// DayOfWeek          string
		// DayOfMonth         string
		// Timestamp          string
		// Century            string
		TimeZone string
		// TimePeriod         string
		// Word               string
		// Sentence           string
		// Paragraph          string
		// Currency           string
		// Amount             float64
		// AmountWithCurrency string
		// UUIDHypenated      string
		UUID string
	}
	arg := StructWithoutTags{}
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	mutateV := VariableMutate(ctx, reflect.TypeOf(arg), reflect.ValueOf(arg))
	fmt.Println(fmt.Sprint(mutateV.Interface()))

	// test the pointer
	mutateV = VariableMutate(ctx, reflect.TypeOf(&arg), reflect.ValueOf(arg))
	fmt.Println(fmt.Sprint(mutateV.Interface()))
}

type DBInterface interface {
	Find() *DB
	Where() *DB
}

func DemoFunc(db DBInterface) error {
	db.Where()
	return nil
}

func TestMutateByReflect(t *testing.T) {
	f := func(db DBInterface) *DB {
		return nil
	}
	iface := reflect.TypeOf(f).In(0)
	emptyStruct := reflect.StructOf([]reflect.StructField{
		{
			Name:      iface.Name(),
			Type:      iface,
			Anonymous: true,
		},
	})
	fmt.Println(emptyStruct.Method(0).Func.Interface())
}

func TestMutateByArg(t *testing.T) {
	type Arg struct {
		DB *struct {
			DBInterface
			InterfaceTag bool
		}
	}
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	arg := Arg{}
	VariableMutate(ctx, reflect.TypeOf(arg), reflect.ValueOf(arg))
}

func TestMutateMockFunc(t *testing.T) {
	a := func() context.Context { return nil }
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	arg := reflect.TypeOf(a).Out(0)
	mutateV := VariableMutate(ctx, arg, reflect.New(arg).Elem())
	_, ok := mutateV.Interface().(context.Context)
	if !ok {
		t.Fatal("no good at mutate")
	}
	fmt.Println(mutateV.Interface())
}

func TestMutatePointerDeepLevel(t *testing.T) {
	type One struct {
		A int
	}
	type Two struct {
		*One
	}
	type Three struct {
		*Two
	}
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	arg := &Three{}
	mutateV := VariableMutate(ctx, reflect.TypeOf(arg), reflect.ValueOf(arg))
	newV, ok := mutateV.Interface().(*Three)
	if !ok {
		t.Fatal("no good at mutate")
	}
	b, err := json.Marshal(newV)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestVariableMutateFunc(t *testing.T) {
	type Fn struct {
		F  func() (int, bool)
		Fp func() *Fn
	}
	a := Fn{}
	ctx := context.Background()
	vtx := atgconstant.VariableContext{Level: 0, ID: 0, CanBeNil: false}
	ctx = contexthelper.SetVariableContext(ctx, vtx)
	mutateV := VariableMutate(ctx, reflect.TypeOf(a), reflect.ValueOf(a))
	t.Log(mutateV.Interface().(Fn).F())
	t.Log(mutateV.Interface().(Fn).Fp())
}
func SetName() string {
	return "Smart Unit Name"
}

func TestVariableMutateBuilder(t *testing.T) {
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
}

type TikTokContext struct {
	ItemID int
}

type UserInfo struct {
	Name string
	Ctx  TikTokContext
	Bag  Bag
}

type Bag struct {
	P Pocket
	A int
	B int
}

type Pocket struct {
	Money  int
	UnUseA int32
	UnUseB int32
	UnUseC int32
}

func QueryData(u UserInfo) string {
	name := u.Name
	count := u.Bag.B + u.Bag.P.Money
	return fmt.Sprint(name, count)
}
