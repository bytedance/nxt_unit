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
package mock

import (
	"reflect"
)

func MakeAOP(funcName string, mockRender *StatementRender, m interface{}) (fun interface{}) {
	function := reflect.TypeOf(m)
	if function.Kind() != reflect.Func {
		panic("it's no true")
	}

	newFunc := reflect.MakeFunc(function, func(args []reflect.Value) (results []reflect.Value) {
		// TODO: line 22 is never used.
		// record in
		// var card []string
		// vtx := atgconstant.VariableContext{MockedRecord: card}
		// ctx := context.Background()
		// ctx = contexthelper.SetVariableContext(ctx, vtx)
		// for _, r := range outs {
		// 	card = append(card, variablecard.ValueToString(ctx, r))
		// }
		outs := reflect.ValueOf(m).Call(args)
		// assert outs
		return outs
	})

	// mockStateMent := fmt.Sprintf("mock%v := mockito.Mock(%s).Return(%s).Build(); defer mock%v.UnPatch();",
	// 	len(mockRender.MockStatement),
	// 	funcName,
	// 	strings.Join(card, ", "),
	// 	len(mockRender.MockStatement))
	// mockRender.MockStatement = append(mockRender.MockStatement, mockStateMent)
	return newFunc.Interface()
}
