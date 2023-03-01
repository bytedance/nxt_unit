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
	"fmt"
)

func init() {
	fmt.Println("hello")
}

var _ = "ddd"

func GoodFuncV3(i, s int, bb string) (string, error, string) {
	if i*s > 9 {
		return bb, nil, "nil"
	}
	switch i {
	case 67 + dadadada(3):
		bb = bb + "ss" + ddddda()
	case gggaaaaa():
		bb = bb + "ss" + ddddda()
	}
	return bb + bb, nil, "nil"
}
