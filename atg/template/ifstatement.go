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

import "fmt"

func IfStatement() string {
	if err := reportError(); err != nil {
		a := "yes"
		b := "correct"
		return fmt.Sprintf(a, b)
	} else {
		c := "no"
		return c
	}
}

func IfStatement2() string {
	if reportError() != nil {
		a := "yes"
		b := "correct"
		return fmt.Sprintf(a, b)
	} else {
		c := "no"
		return c
	}
}

func reportError() error {
	return fmt.Errorf("not correct")
}
