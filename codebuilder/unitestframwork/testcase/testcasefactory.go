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
package testcase

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/bytedance/nxt_unit/atghelper"

	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
	"github.com/bytedance/nxt_unit/codebuilder/unitestframwork/statement"
	"golang.org/x/tools/go/ssa"
)

type StatementUsageValue struct {
	Statement statement.Statement
	Used      bool
}

type TestCase struct {
	Name              string
	TestedPackageName string
	Statements        []statement.Statement
	StatementUsageMap map[string]*StatementUsageValue
	ID                int64
	IsChanged         bool
}

// Test case only contains the mocked function
func CreateTestCase(ctx context.Context, f *parsermodel.ProjectFunction) (*TestCase, error) {
	testCase := &TestCase{
		Name:              fmt.Sprint("Test", atghelper.UpperCaseFirstLetter(f.Function.Name()), atghelper.RandStringBytes(8)),
		TestedPackageName: f.Function.Pkg.Pkg.Name(),
		Statements:        make([]statement.Statement, 0),
		StatementUsageMap: make(map[string]*StatementUsageValue, 0),
		ID:                rand.Int63(),
	}
	for _, function := range f.CalleeFunctionsForTargetFunction {
		InsertRandomCall(ctx, function, testCase, len(testCase.Statements), f.Program.PkgPath)
	}

	return testCase, nil
}

func InsertRandomCall(ctx context.Context, function *ssa.Function, testCase *TestCase, position int, pkgPath string) bool {
	if function == nil {
		return false
	}
	st, _ := statement.CreateMockedStatement(ctx, function, pkgPath)
	if st != nil {
		testCase.Statements = append(testCase.Statements, *st)
	}
	return true
}
