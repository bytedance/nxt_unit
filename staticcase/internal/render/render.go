// Copyright cweill/gotests authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package goparse contains logic for parsing Go files. Specifically it parses
// source and test files into domain models for generating tests.

//This file may have been modified by Bytedance Ltd. and/or its affiliates (“Bytedance's Modifications”).
// All Bytedance's Modifications are Copyright (2022) Bytedance Ltd. and/or its affiliates.
package render

import (
	"fmt"
	"github.com/bytedance/nxt_unit/staticcase/internal/render/bindata"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/bytedance/nxt_unit/atgconstant"

	"github.com/bytedance/nxt_unit/staticcase/internal/models"
	"github.com/bytedance/nxt_unit/staticcase/templates"
)

type Render struct {
	tmpls *template.Template
}

func New() *Render {
	r := Render{
		tmpls: template.New("render").Funcs(map[string]interface{}{
			"Field":    fieldName,
			"Receiver": receiverName,
			"Param":    parameterName,
			"Want":     wantName,
			"Got":      gotName,
		}),
	}

	// default templates first
	for _, name := range bindata.AssetNames() {
		r.tmpls = template.Must(r.tmpls.Parse(string(bindata.MustAsset(name))))
	}
	return &r
}

// LoadCustomTemplates allows to load in custom templates from a specified path.
func (r *Render) LoadCustomTemplates(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir: %v", err)
	}

	templateFiles := []string{}
	for _, f := range files {
		templateFiles = append(templateFiles, path.Join(dir, f.Name()))
	}
	r.tmpls, err = r.tmpls.ParseFiles(templateFiles...)
	if err != nil {
		return fmt.Errorf("tmpls.ParseFiles: %v", err)
	}

	return nil
}

// LoadCustomTemplatesName allows to load in custom templates of a specified name from the templates directory.
func (r *Render) LoadCustomTemplatesName(name string) error {
	f, err := templates.Dir(false, "/").Open(name)
	if err != nil {
		return fmt.Errorf("templates.Open: %v", err)
	}

	files, err := f.Readdir(nFile)
	if err != nil {
		return fmt.Errorf("f.Readdir: %v", err)
	}
	for _, f := range files {
		text, err := templates.FSString(false, path.Join("/", name, f.Name()))
		if err != nil {
			return fmt.Errorf("templates.FSString: %v", err)
		}

		if tmpls, err := r.tmpls.Parse(text); err != nil {
			return fmt.Errorf("tmpls.Parse: %v", err)
		} else {
			r.tmpls = tmpls
		}
	}

	return nil
}

// LoadFromData allows to load from a data slice
func (r *Render) LoadFromData(templateData [][]byte) {
	for _, d := range templateData {
		r.tmpls = template.Must(r.tmpls.Parse(string(d)))
	}
}

func (r *Render) Header(w io.Writer, h *models.Header, headerTemplate string) error {
	if err := r.tmpls.ExecuteTemplate(w, headerTemplate, h); err != nil {
		return err
	}
	_, err := w.Write(h.Code)
	return err
}

func (r *Render) TestFunction(
	w io.Writer,
	f *models.Function,
	printInputs bool,
	subtests bool,
	named bool,
	parallel bool,
	params map[string]interface{},
	mock []string,
	builder []string,
	middleBuilder []string,
	testCaseNum int,
	useMockType int,
	Uid string,
	rowData string,
	testMode string,
	filePath string,
	globalInit []string,
) error {
	switch testMode {
	case atgconstant.FinalTest:
		return r.tmpls.ExecuteTemplate(w, "finalsuite", struct {
			*models.Function
			PrintInputs    bool
			Subtests       bool
			Parallel       bool
			Named          bool
			UseMockType    int
			RowData        string
			Builders       []string
			MiddleBuilders []string
			MockStateMents []string
			Uid            string
			TemplateParams map[string]interface{}
		}{
			Function:       f,
			PrintInputs:    printInputs,
			Subtests:       subtests,
			Parallel:       parallel,
			Named:          named,
			UseMockType:    useMockType,
			RowData:        rowData,
			Builders:       builder,
			MiddleBuilders: middleBuilder,
			MockStateMents: mock,
			Uid:            "",
			TemplateParams: params,
		})
	case atgconstant.MiddleCode:
		return r.tmpls.ExecuteTemplate(w, "function", struct {
			*models.Function
			PrintInputs    bool
			Subtests       bool
			Parallel       bool
			Named          bool
			UseMockType    int
			Mocks          []string
			Builders       []string
			TestCaseNum    int
			MaxTestCaseN   []string
			Uid            string
			FilePath       string
			RealName       string
			TemplateParams map[string]interface{}
		}{
			Function:       f,
			PrintInputs:    printInputs,
			Subtests:       subtests,
			Parallel:       parallel,
			Named:          named,
			UseMockType:    useMockType,
			Uid:            Uid,
			FilePath:       filePath,
			Mocks:          mock,
			Builders:       builder,
			TestCaseNum:    testCaseNum,
			RealName:       strings.Replace(f.Name, Uid, "", -1),
			TemplateParams: params,
		})
	case atgconstant.BaseTest:
		var baseTemp string = "basefunc"

		return r.tmpls.ExecuteTemplate(w, baseTemp, struct {
			*models.Function
			PrintInputs    bool
			Subtests       bool
			Mocks          []string
			Parallel       bool
			Named          bool
			RowData        string
			Uid            string
			TemplateParams map[string]interface{}
			GlobalInit     []string
			UseMockType    int
		}{
			Function:       f,
			PrintInputs:    printInputs,
			Subtests:       subtests,
			Mocks:          mock,
			Parallel:       parallel,
			Named:          named,
			TemplateParams: params,
			GlobalInit:     globalInit,
			UseMockType:    useMockType,
		})
	default:
		return nil
	}
}
