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
package template

import (
	"testing"

	convey "github.com/smartystreets/goconvey/convey"
)

// we create the test template for the runnable function
// please fill the testcase and mock function
func TestGoodFunc_OOOWSU(t *testing.T) {
	type Args struct {
		I  int
		S  int
		Bb string
	}
	type test struct {
		Name    string
		Args    Args
		Want    string
		WantErr bool
	}
	tests := []test{
		// TODO: add the testcase
	}
	for _, tt := range tests {
		convey.Convey(tt.Name, t, func() {
			// TODO: add the return of mock functions
			got, err := GoodFunc(tt.Args.I, tt.Args.S, tt.Args.Bb)
			if (err != nil) != tt.WantErr {
				t.Errorf("%q. GoodFunc() error = %v, wantErr %v", tt.Name, err, tt.WantErr)
			}
			if got != tt.Want {
				t.Errorf("%q. GoodFunc() = %v, want %v", tt.Name, got, tt.Want)
			}
		})
	}
}
