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
package instrumentation

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bytedance/nxt_unit/manager/logextractor"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"golang.org/x/tools/imports"
)

func NewFunctionBuilder(ctx context.Context) *functionBuilder {
	return &functionBuilder{}
}

type functionBuilder struct {
	filename  string
	id        string
	TotalLine int
	closeFunc func() error
}

func (f *functionBuilder) Offset() string {
	return f.id
}

func (f *functionBuilder) Build(ctx context.Context) (string, error) {
	opt, _ := contexthelper.GetOption(ctx)
	id, _ := contexthelper.GetBuilderVector(ctx)
	f.id = id
	newFile, total, err := NewInstrumentation(opt.FilePath, opt.FuncName, id)
	if err != nil {
		return "", err
	}
	f.TotalLine = total
	// fix for filepath has ".go"  string ,such as /a/c/username.go/e/f.go
	if strings.HasSuffix(opt.FilePath, ".go") {
		f.filename = strings.TrimSuffix(opt.FilePath, ".go") + id + "smartunit.txt"
	} else {
		f.filename = opt.FilePath
	}
	out, err := imports.Process(f.filename, newFile, nil)
	if err != nil {
		return "", fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.SourceCodePluginCodeError, err.Error())
	}
	if err := ioutil.WriteFile(f.filename, out, atgconstant.NewFilePerm); err != nil {
		return "", fmt.Errorf("[Sorry]: the bug comes from file write %v has ioutil.WriteFile err: %v", opt.FuncName, err)
	}
	return f.filename, nil
}

func (f *functionBuilder) Close() error {
	return f.closeFunc()
}
