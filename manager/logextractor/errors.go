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
package logextractor

import (
	"errors"
)

// InspiredCloudCreatedError Inspired Cloud Error
var InspiredCloudCreatedError = errors.New("Error Code is N3042, inspired cloud create\n")
var InspiredCloudUpdatedError = errors.New("Error Code is N3043, inspired cloud update.\n")

// LocalFileNotSupported LocalFileError Local File is not supported
var LocalFileNotSupportedError = errors.New("Error Code is P3044, local file is not supported.We cannot generate the unit test for this file.For example, this file belongs to kite_x, govender\n")

// Cannot parse the test function
var CannotParseTestedFunctionError = errors.New("Error Code is P3045, cannot find tested function.\n")

// ParseProgramError Source Code Add Code Error
var ParseProgramError = errors.New("Error Code is P3046, cannot find go executable file.\n")

// SourceCodePluginCodeError Source Code Add Code Error
var SourceCodePluginCodeError = errors.New("Error Code is I3046, cannot add plugin to source code.\n")

// Middle Code generate error
var MiddleCodeGenerateError = errors.New("Error Code is I3047, middle code cannot be generated.\n")

// Go mod tidy in Fix Go File
var GoModTidyError = errors.New("Error Code is P3048, go mod tidy error.\n")

// Middle code cannot generate the final code
var MiddleCodeCannotGenerateFinalCodeError = errors.New("Error Code is I3049, cannot generate the final code for middle code.\n")

// Generate test is not runnable
var GenerateTestNotRunnableError = errors.New("Error Code is I3050, generate code cannot be runnable.\n")

// Merge test has a conflict
var MergeTestConflictError = errors.New("Error Code is I3051, merge code conflict.\n")

// Generate test template fail after generate test failed
var GenerateTestTemplateError = errors.New("Error Code is I3052, cannot generate template.\nThis error is normally go with other errors. It is not able to be fixed.\nPlease fix other errors first.\n")

// Generate template import process error
var GenerateTemplateImportError = errors.New("Error Code is I3053, template import merge fail.\n")

// Cannot find the tested function
var CannotFindTestedFunctionError = errors.New("Error Code is I3058, cannot find tested function.\nYou can run go vet [my/project/...] or go build to resolve the compiling issue.")

// GoGetError Go Get Error
var GoGetError = errors.New("Error Code is P3054, go get error.\nYou can run go vet [my/project/...] or go build to resolve the compiling issue.")

// Generate test template cannot parse import error
var GenerateTestTemplateCannotParseImportError = errors.New("Error Code is I3055, template import cannot be parsed.\n")

// Generate test template rename error
var GenerateTestTemplateRenameError = errors.New("Error Code is I3056, template rename error.\n")

// Generate test rename error
var GenerateTestRenameError = errors.New("Error Code is I3057, test rename error.\n")

// Generate test template internal error
var GenerateTestTemplateInternalError = errors.New("Error Code is I3059, template internal error.\n")

// DependencyInitError DependencyError DependencyInitError Dependency Error
var DependencyInitError = errors.New("error Code is P3023.\n")
var NullPointerError = errors.New("error Code is P3055.\n")
