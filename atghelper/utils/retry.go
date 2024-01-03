/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package utils

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var Retry0TimeError = errors.New("retry time can't be 0")
var RetryTimeout0Error = errors.New("retry timeout can't be 0")
var TimeoutErr = errors.New("retry timeout")
var ShouldNotRetryErr = errors.New("Shouldn't retry")

type runAbortFunc func(interface{}) (bool, error)

func RetryDo(name string, times int, timeout time.Duration, run runAbortFunc, args interface{}) error {
	if times == 0 {
		return Retry0TimeError
	}
	if timeout == 0 {
		return RetryTimeout0Error
	}

	var err error
	timeoutAbort := int32(0)
	abort := false
	resultCh := make(chan error, 1)
	panicSignal := make(chan string)
	quit := make(chan bool)

	go func() {
		var runErr error
		defer func() {
			panicInfo := recover()
			if panicInfo == nil {
				return
			}
			runErr = errors.New(fmt.Sprint(panicInfo))
			select {
			//	received quit signal from the timeout
			case <-quit:
			case panicSignal <- fmt.Sprintf("[DoCanAbort-Recovery] name %s,err=%v smartunit_debug.Stack:\n%s", name, runErr, debug.Stack()):
			}
		}()
		for i := 0; i < times; i++ {
			if atomic.LoadInt32(&timeoutAbort) > 0 {
				fmt.Println("stop retry on timeout")
				return
			}
			// run!!!
			abort, runErr = run(args)

			if runErr == nil {
				resultCh <- nil
				return
			}
			// abort is meaningful only when err is not nil
			if abort {
				resultCh <- runErr
				return
			}
		}
		resultCh <- runErr
		return
	}()

	timer := time.NewTimer(timeout)
	select {
	case err = <-resultCh:
		// a read from ch has occurred
	case panicInfo := <-panicSignal:
		// panic!!!
		close(quit)
		panic(panicInfo)
	case <-timer.C:
		// the read from ch has timed out
		quit <- true
		err = TimeoutErr
		atomic.StoreInt32(&timeoutAbort, 1)
	}
	timer.Stop()

	return err
}
