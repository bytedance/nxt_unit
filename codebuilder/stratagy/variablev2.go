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
package stratagy

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
)

type SpecialVariableV2 struct {
	Type      string
	PkgName   string
	PkgPath   string
	Mutated   bool // we currently only support two types of mutation
	IsPointer bool
	CanBeNil  bool
}

func GetSpecialVariableV2(ctx context.Context, v reflect.Value) (string, bool) {
	vtx, ok := contexthelper.GetVariableContext(ctx)
	if !ok {
		return "", false
	}

	if !v.IsValid() {
		return "", false
	}

	t := v.Type()
	if t == nil {
		fmt.Printf("[GetSpecialVariableV2] v type is empty %v\n", v.String())
		return "", false
	}
	pkgPath := ""
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		pkgPath = t.Elem().PkgPath()
	default:
		pkgPath = t.PkgPath()
	}

	if pkgPath == "github.com/gin-gonic/gin" && strings.Contains(v.String(), "gin.Context") {
		if v.Kind() == reflect.Ptr {
			if vtx.CanBeNil && atghelper.RandomBool(atgconstant.SpecialValueBeNil) {
				return "nil", true
			}
			pkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet("gin", "github.com/gin-gonic/gin")
			return fmt.Sprint("&", pkgName, ".Context{Request: &http.Request{URL: &url.URL{Path: \"test_path\",}}}"), true
		}
		pkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet("gin", "github.com/gin-gonic/gin")
		return fmt.Sprint(pkgName, ".Context{Request: &http.Request{URL: &url.URL{Path: \"test_path\",}}}"), true
	}
	if strings.HasPrefix(v.String(), "<*errors.errorString") {
		if vtx.CanBeNil && atghelper.RandomBool(atgconstant.SpecialValueBeNil) {
			return "nil", true
		}
		duplicatepackagemanager.GetInstance(ctx).PutAndGet("errors", "")
		return "errors.New(\"smart unit\")", true
	}
	if strings.HasPrefix(v.String(), "<*context.") {
		duplicatepackagemanager.GetInstance(ctx).PutAndGet("context", "")
		return "context.Background()", true
	}
	return "", false
}
