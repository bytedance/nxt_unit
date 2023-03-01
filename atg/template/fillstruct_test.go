package template

import (
	"fmt"
	"testing"
)

func TestBasicType(t *testing.T) {
	type suBaseTest struct {
		A string
		B int
		C float64
	}
	tests := []suBaseTest{
		// add more test

		{
			A: "4hXvCIiwHH",
			B: 1636444755,
			C: 0.09208274320692282,
		}, {
			A: "4hXvCIiwHH",
			B: 1636444755,
			C: 0.09208274320692282,
		}, {
			A: "6Ud1EY1mdk",
			B: 461531886,
			C: 0.9352114064090074,
		}, {},
	}
	fmt.Println(tests)
}
