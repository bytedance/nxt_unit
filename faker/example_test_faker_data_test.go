package faker

import (
	"reflect"
)

// Gondoruwo ...
type Gondoruwo struct {
	Name       string
	Locatadata int
}

// custom type that aliases over slice of byte
type CustomUUID []byte

// Sample ...
type Sample struct {
	ID        int64      `faker:"customIdFaker"`
	Gondoruwo Gondoruwo  `faker:"gondoruwo"`
	Danger    string     `faker:"danger"`
	UUID      CustomUUID `faker:"customUUID"`
}

// CustomGenerator ...
func CustomGenerator() {
	_ = AddProvider("customIdFaker", func(v reflect.Value) (interface{}, error) {
		return int64(43), nil
	})
	_ = AddProvider("danger", func(v reflect.Value) (interface{}, error) {
		return "danger-ranger", nil
	})

	_ = AddProvider("gondoruwo", func(v reflect.Value) (interface{}, error) {
		obj := Gondoruwo{
			Name:       "Power",
			Locatadata: 324,
		}
		return obj, nil
	})

	_ = AddProvider("customUUID", func(v reflect.Value) (interface{}, error) {
		s := CustomUUID{
			0, 8, 7, 2, 3,
		}
		return s, nil
	})
}
