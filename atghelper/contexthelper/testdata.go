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
package contexthelper

import (
	"context"
	"path"
	"runtime"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/codebuilder/setup"
)

func GetTestContext() context.Context {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	filePath = path.Join(path.Dir(filePath), "../../atg/template/atg.go")
	option := atgconstant.Options{
		FuncName:     "GoodFunc",
		FilePath:     filePath,
		Level:        1,
		Maxtime:      4,
		GenerateType: atgconstant.GAMode,
		Uid:          "Vector",
	}
	sourceFunc, err := setup.GetFunctions(option)
	if err != nil {
		panic(err)
	}
	atgconstant.PkgRelativePath = sourceFunc.TestFunction.Program.PkgPath
	ctx := context.Background()
	ctx = SetOption(ctx, option)
	ctx = SetSetupFunc(ctx, sourceFunc)
	ctx = SetBuilderVector(ctx, option.Uid)
	return ctx
}

func GetTestContextV2() context.Context {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	filePath = path.Dir(filePath)
	directoryPath := path.Join(filePath, "../../atg/template")
	filePath = path.Join(filePath, "../../atg/template/atg.go")
	option := atgconstant.Options{
		FuncName:      "siwei123",
		FilePath:      filePath,
		Level:         1,
		Maxtime:       4,
		GenerateType:  atgconstant.GAMode,
		MinUnit:       "file",
		Uid:           "Vector",
		DirectoryPath: directoryPath,
	}
	sourceFunc, err := setup.GetFunctions(option)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = SetOption(ctx, option)
	ctx = SetSetupFunc(ctx, sourceFunc)
	ctx = SetBuilderVector(ctx, option.Uid)
	return ctx
}
