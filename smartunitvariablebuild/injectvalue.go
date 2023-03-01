package smartunitvariablebuild

import "reflect"

func NewSpecialValueInjector() *SpecialValueInjector {
	return &SpecialValueInjector{
		ValueMap:   map[string]reflect.Value{},
		BuilderMap: map[string]string{},
	}
}

type SpecialValueInjector struct {
	ValueMap   map[string]reflect.Value
	BuilderMap map[string]string
}

func (s *SpecialValueInjector) Set(key string, value reflect.Value) {
	s.ValueMap[key] = value
}

func (s *SpecialValueInjector) SetBuilder(value interface{}, code string) {
	t := reflect.TypeOf(value)
	s.ValueMap[t.String()] = reflect.ValueOf(value)
	s.BuilderMap[t.String()] = code
}

func (s *SpecialValueInjector) GetCode(value reflect.Value) (string, bool) {
	v, exist := s.BuilderMap[value.Type().String()]
	return v, exist
}

func (s *SpecialValueInjector) Get(key string) (reflect.Value, bool) {
	v, exist := s.ValueMap[key]
	return v, exist
}
