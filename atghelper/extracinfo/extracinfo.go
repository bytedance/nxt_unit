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
package extracinfo

import (
	"context"
	"errors"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"go/types"
	"strings"

	"github.com/bytedance/nxt_unit/atghelper"

	"github.com/bytedance/nxt_unit/codebuilder/stratagy"

	"github.com/bytedance/nxt_unit/codebuilder/setup/parsermodel"
)

func ConvertToReceiver(ctx context.Context, signature *types.Signature, testedPackageName string) (*parsermodel.Field, error) {
	recv := signature.Recv()
	if recv == nil {
		return nil, nil
	}

	field := &parsermodel.Field{}
	field.Name = recv.Name()
	var expression = &parsermodel.Expression{}

	field.Type = expression
	underlyTypeStr := types.TypeString(recv.Type().Underlying(),
		func(p *types.Package) string {
			return qfByCtx(ctx, p)
		})
	if atghelper.IsBasicDataType(underlyTypeStr) {
		return nil, errors.New("not mock basic type receiver")
	}
	expression.Value = types.TypeString(recv.Type(),
		func(p *types.Package) string {
			return qfByCtx(ctx, p)
		})
	expression.PkgPath, expression.PkgName = getPkgPathAndPakName(recv.Type())

	// Hanlde the redis-v6 case
	if strings.Contains(expression.PkgName, "-") {
		expression.PkgName = strings.ReplaceAll(expression.PkgName, "-", "_")
	}

	// Special Logic for some receivers
	// The receiver here is used for the mock((receiver).function).Thus, the receiver should be like
	// the *kernel.AbContext
	if v, ok := stratagy.GetSpecialVariable(field, false, false); ok {
		expression.Value = v.GetConstructor()
		expression.PkgPath = v.PkgPath
		expression.PkgName = v.PkgName
		duplicatepackagemanager.GetInstance(ctx).PutAndGet(expression.PkgName, expression.PkgPath)
	}

	if !strings.Contains(expression.Value, ".") {
		expression.PkgName = ""
		expression.PkgPath = ""
	} else {
		tempInfo := duplicatepackagemanager.GetInstance(ctx).GetImportInfo()
		expression.PkgName = tempInfo.Name
		expression.PkgPath = tempInfo.PackagePath
	}
	if expression.Value == "io.Writer" {
		expression.IsWriter = true
	}

	return field, nil
}

// We used a ugly way to fetch the package name and path here.
// Also we store the name, create a new alias for that package.
func qfByCtx(ctx context.Context, pkg *types.Package) string {
	pkgName := pkg.Name()
	path := pkg.Path()
	// Special handling of go vendor
	paths := strings.SplitN(path, "vendor/", 2)
	if len(paths) == 2 {
		path = paths[1]
	}
	if len(path) == 0 {
		return pkgName
	}
	tempImportInfo := duplicatepackagemanager.GetInstance(ctx).GetImportInfo()
	tempImportInfo.Name, tempImportInfo.PackagePath =
		duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, path)
	duplicatepackagemanager.GetInstance(ctx).SetImportInfo(tempImportInfo)
	return tempImportInfo.Name
}

func getPkgPathAndPakName(t types.Type) (string, string) {

	if point, ok := t.Underlying().(*types.Pointer); ok {
		return getPkgPathAndPakName(point.Elem())
	}

	if slice, ok := t.Underlying().(*types.Slice); ok {
		return getPkgPathAndPakName(slice.Elem())
	}

	if _, ok := t.Underlying().(*types.Struct); ok {
		named, ok := t.(*types.Named)

		if !ok || named.Obj().Pkg() == nil {
			return "", ""
		}

		return named.Obj().Pkg().Path(), named.Obj().Pkg().Name()
	}

	if _, ok := t.Underlying().(*types.Named); ok {
		named := t.(*types.Named)

		if named.Obj().Pkg() == nil {
			return "", ""
		}
		return named.Obj().Pkg().Path(), named.Obj().Pkg().Name()
	}

	return "", ""
}
