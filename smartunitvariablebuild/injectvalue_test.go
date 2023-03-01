package smartunitvariablebuild

import (
	"reflect"
	"testing"
)

func NewName() string {
	return "smartunit name"
}

func TestSpecialValueInjector_Get(t *testing.T) {
	s := NewSpecialValueInjector()
	s.Set("int", reflect.ValueOf(2))
	v, ok := s.Get("int")
	if !ok {
		t.Fatal("test error")
	}
	t.Log(v.Interface())
}

func TestSpecialValueInjector_Bulider(t *testing.T) {
	s := NewSpecialValueInjector()
	s.SetBuilder(NewName(), "NewName()")
	v, ok := s.Get("string")
	if !ok {
		t.Fatal("test error")
	}
	t.Log(v.Interface())
	code, ok := s.GetCode(reflect.ValueOf(""))
	if !ok {
		t.Fatal("test error")
	}
	t.Log(code)
}
