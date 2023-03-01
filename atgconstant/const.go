/*
 * Copyright 2022 Bytedance Ltd. and/or its affiliates
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
package atgconstant

import (
	"go/build"
	"go/types"
	"math/rand"
	"os"
	"path"
	"time"
)

const (
	Execution          = "execution"
	CalcExpr           = "CalcExpr"
	CalcBinaryExpr     = "CalcBinaryExpr"
	BranchDistances    = "BranchDistances"
	TrueDistances      = "TrueDistances"
	ReturnType         = "ExecutedReturnType"
	FalseDistances     = "FalseDistances"
	Ident              = "Ident"
	BinaryExpr         = "BinaryExpr"
	ExecutedPredicates = "ExecutedPredicates"
	ExecutedFunctions  = "ExecutedFunctions"
	K                  = 1.0
	FilePath           = "/execution/testdata/ifstatementtestdata.go"
	IfStatement        = "_testdata/executiontest/ifstatementtestdata.go"
	SelectStatement    = "_testdata/executiontest/selectstatementtestdata.go"
	GoTypeConverter    = "/_testdata/gotypeconverter.go"

	StatementInsertion  float64 = 0.5
	TestCaseLength      int     = 20
	MaxAttempts         int     = 1000
	Delta               float64 = 20.0
	ChangeParameter     float64 = 0.5
	MutationDelete      float64 = 0.8
	MutationChange      float64 = 0.1
	MutationAdd         float64 = 100.0
	SpecialValueBeNil   float64 = 0.5
	MaxInt              int     = 2048
	TestCaseTrimAttempt int     = 2
	VariableMaxLevel    int     = 4

	// Initial probability of inserting a new test in a test suite
	TestInsertionProbability float64 = 0.1

	// Maximum number of test cases in a test suite
	TestSuiteMaxSize int = 2

	// Plugin maximum number of test cases in a test suite
	PluginTestSuiteMaxSize int = 4

	// Possibility that we shouldn't mock the statement
	MockStatementRatio = 0.2

	// Plugin timeout
	PluginTimeOut = 120

	// Without GA Timeout
	WithoutGATimeOut = 200

	// GA TIMEOUT
	GATimeOUT = 300

	// Population size of genetic algorithm
	Population int = 1

	// Maximum iterations
	AlgorithmIterations int = 1

	// Probability of crossover
	CrossoverRate float64 = 0.75

	// random string length
	RandomStringLen int = 5

	// Preprocess Test Suite. It determine how many test case that we want to exist in the test suite
	// Please take a look at the function: PreProcessTestSuite
	PreProcessTestSuite = 1

	// Receiver Struct Field Max limit
	ReceiverFieldMaxLimit = 66

	NewFilePerm os.FileMode = 0777

	GraphLevel        int = 1
	SmartUnitVersion      = "0.3.0"
	SimpleMode            = "simple"
	GAMode                = "ga"
	DebugMode             = false
	ReportMode            = false
	MinUnit               = "function"
	PluginMode            = "plugin"
	PluginQMode           = "pluginq"
	SplitFunctionMode     = "splitfunction"
	Backend               = "backstage"
	FileMode              = "file"
	InternalPkg           = "internal"
	FinalTest             = "finaltest"
	MiddleCode            = "middlecode"
	BaseTest              = "basetest"

	// plugin variables
	Case             = "case"
	Template         = "template"
	FunctionTemplate = "function_template"
)

// mock type define
const (
	UseMockUnknown int = iota
	UseNoMock
	UseMockitoMock
	UseGoMonkeyMock
)

type ImportInfo struct {
	Name        string
	PackagePath string
	Used        bool
}

// Set of options to use when generating tests.
type Options struct {
	// 被测函数名称
	FuncName string
	// 被测函数所在文件地址
	FilePath string
	// file_name
	FileName string
	// GA算法迭代次数：推荐value值大于2，小于10
	Level int
	// TestSuite 最大Case数量，推荐Value: 4或5
	Maxtime int
	// 最终生成testsuit的模式
	GenerateType       string
	DebugMode          bool
	ReportMode         bool
	MinUnit            string
	Uid                string
	DirectoryPath      string
	FunctionList       []string
	Usage              string
	ReceiverName       string
	FinalSuiteTestName string
	RunForFinalSuite   bool
	TemplateType       int
	UseMockType        int
}

// ExecutionValues is used for the test suite
type ExecutionValues struct {
	Coverage           float64             // test suite coverage
	TestCasePredicates map[string][]string // Key: test case name. Value: predicates chain.
	TestCaseReturnNils map[string][]bool   // Key: test case name. Value: return is nil or not. from left to right.
	PanicTcName        []string            // Panic test case names.
}

// ExecutionValue is used for the  test case
type ExecutionValue struct {
	Coverage          float64  // test case coverage
	TestCasePredicate []string // predicates chain.
	TestCaseReturnNil []bool   // return is nil or not. from left to right.
	PanicTcName       string   // Panic test case names.
}

type DynamicVariable struct {
	Var         *types.Var
	IsSignature bool
}

type VariableContext struct {
	Level        int
	ID           int
	CanBeNil     bool
	MockedRecord []string
}

var GOPATHSRC string
var GOROOT string
var GoDirective string
var IgnoredPath map[string]string
var ProjectPath = "/github.com/bytedance/nxt_unit"
var ProjectPrefix = "github.com"
var PkgRelativePath string
var TempImportInfo ImportInfo
var FinalSuiteTestName string

const Letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	// gopath
	gopath := build.Default.GOPATH
	GOPATHSRC = path.Join(gopath, "src")

	// goroot
	GOROOT = build.Default.GOROOT

	IgnoredPath = map[string]string{
		"github.com/golang/protobuf/proto": "proto",
		"github.com/gogo/protobuf/proto":   "proto",
	}
	rand.Seed(time.Now().UnixNano())
}
