package utils

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetryDo(t *testing.T) {
	run := func(input interface{}) (bool, error) {
		total, ok := input.(int)
		if !ok {
			return false, errors.New("not Int type")
		}
		if total < 5 {
			return false, errors.New("not match condition")
		}
		return true, nil
	}
	err := RetryDo("test", 1, 10*time.Microsecond, run, 1)
	if err != nil {
		fmt.Println(err)
	}

}
