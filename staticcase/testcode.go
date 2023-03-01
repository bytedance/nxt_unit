// Package gotests contains the core logic for generating table-driven tests.
package staticcase

import (
	"context"
	"fmt"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"go/importer"
	"go/types"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/bytedance/nxt_unit/manager/logextractor"

	"github.com/bytedance/nxt_unit/codebuilder/instrumentation"

	"github.com/bytedance/nxt_unit/codebuilder/setup"

	"github.com/bytedance/nxt_unit/atgconstant"

	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"

	"github.com/bytedance/nxt_unit/staticcase/internal/goparser"
	"github.com/bytedance/nxt_unit/staticcase/internal/input"
	"github.com/bytedance/nxt_unit/staticcase/internal/models"
	"github.com/bytedance/nxt_unit/staticcase/internal/output"
)

// Options provides custom filters and parameters for generating tests.
type Options struct {
	Only            *regexp.Regexp            // Includes only functions that match.
	Exclude         *regexp.Regexp            // Excludes functions that match.
	Exported        bool                      // Include only exported methods
	PrintInputs     bool                      // Print function parameters in error messages
	Subtests        bool                      // Print tests using Go 1.7 subtests
	Ctx             context.Context           // Context of the smart unit
	TestCaseNum     int                       // Test case number
	Uid             string                    // uuid for each function
	TestMode        string                    // type of render test
	RecordRow       map[string]string         // rowData on each function
	Builder         []string                  // global value builder
	MiddleBuilder   []string                  // global value builder
	FilePath        string                    // realpath
	Parallel        bool                      // Print tests that runs the subtests in parallel.
	Named           bool                      // Create Map instead of slice
	Importer        func() types.Importer     // A custom importer.  // TODO: remoeve it because I don't know to use it
	Template        string                    // Name of custom template set
	TemplateDir     string                    // Path to custom template set
	TemplateParams  map[string]interface{}    // Custom external parameters
	TemplateData    [][]byte                  // Data slice for templates
	ExistTestImport []*instrumentation.Import // Imports from exist test file for template
	UseMockType     int                       // whether to user mock pkg in cases
	UseMockMap      map[string]map[string]int // Data for function for mock
}

// A GeneratedTest contains information about a test file with generated tests.
type GeneratedTest struct {
	Path      string             // The test file's absolute path.
	Functions []*models.Function // The functions with new test methods.
	Output    []byte             // The contents of the test file.
}

func RecordFinalSuite(ctx context.Context, filePath string, funcData map[string]string, useMockMap map[string]map[string]int, useMockType int) (code []byte, err error) {
	opt := &Options{}
	srcFiles, err := input.Files(filePath)
	if err != nil {
		return nil, fmt.Errorf("input.Files: %v", err)
	}
	if len(funcData) <= 0 {
		return nil, fmt.Errorf("[render final test]: no runnable code")
	}
	opt.UseMockType = useMockType
	opt.UseMockMap = useMockMap
	opt.Builder = duplicatepackagemanager.GetInstance(ctx).GetInitBuilder()
	opt.MiddleBuilder = duplicatepackagemanager.GetInstance(ctx).GetMiddleBuilder()
	duplicatepackagemanager.SetInstance(ctx)
	// duplicatepackagemanager.GetInstance(ctx).SetRelativePath(funcData)
	files, err := input.Files(path.Dir(filePath))
	if err != nil {
		return nil, fmt.Errorf("input.Files: %v", err)
	}
	if opt.Importer == nil || opt.Importer() == nil {
		opt.Importer = importer.Default
	}
	opt.RecordRow = funcData
	opt.TestMode = atgconstant.FinalTest
	if opt.RecordRow == nil {
		opt.RecordRow = map[string]string{}
	}
	opt.Ctx = ctx
	gt, err := generateTest(srcFiles[0], files, opt)
	if err != nil {
		return nil, err
	}
	if gt == nil {
		return nil, fmt.Errorf("generateTest is nil ")
	}
	return gt.Output, nil
}

// GenerateTests generates table-driven tests for the function and method
// signatures defined in the target source path file(s). The source path
// parameter can be either a Go source file or directory containing Go files.
func GenerateTests(srcPath string, opt *Options) (*GeneratedTest, error) {
	if opt == nil {
		opt = &Options{}
	}
	srcFiles, err := input.Files(srcPath)
	if err != nil {
		return nil, fmt.Errorf("input.Files: %v", err)
	}
	files, err := input.Files(path.Dir(srcPath))
	if err != nil {
		return nil, fmt.Errorf("input.Files: %v", err)
	}
	if opt.Importer == nil || opt.Importer() == nil {
		opt.Importer = importer.Default
	}

	// TODO(siwei.wang): No need the srcFiles. Because its length is always one.
	for _, srcP := range srcFiles {
		rs, err := generateTest(srcP, files, opt)
		if err != nil {
			return nil, err
		}
		return rs, nil
	}
	return nil, fmt.Errorf("[GenerateTests] doesn't srcFiles is not exceed one")
}

// result stores a generateTest result.
type result struct {
	gt  *GeneratedTest
	err error
}

// TODO: generateTest was called by middlecode only. Another generateTest is in the smartunit_mate
func generateTest(src models.Path, files []models.Path, opt *Options) (*GeneratedTest, error) {
	p := &goparser.Parser{Importer: opt.Importer()}
	sr, err := p.Parse(string(src), files)
	if err != nil {
		return nil, fmt.Errorf("Parser.Parse source file: %v", err)
	}
	h := sr.Header
	// TODO: below's logic is not related with the smart unit. They are reading the _test file of the original file
	// And then do something. However, I don't know what the funcs's doing. Just keep them.
	h.Code = nil // Code is only needed from parsed test files.
	testPath := models.Path(src).TestPath()
	h, tf, err := parseTestFile(p, testPath, h)
	if err != nil {
		return nil, err
	}
	h.Code = nil

	// Add the tested file imports and smart unit test file import
	// check _test.go exist,if exist options set imports from existing test file
	targetFilePath := strings.TrimSuffix(string(models.Path(src)), ".go") + "_test.go"
	tgtExist := atghelper.IsFileExist(targetFilePath)
	var targetImports []*instrumentation.Import
	if tgtExist {
		targetImports, err = instrumentation.GetImportsInfosFromFile(targetFilePath)
		if err == nil {
			opt.ExistTestImport = targetImports
		}
	}

	// check _smart_unit_test.go exist, if exist options set imports from existing smart unit test file
	suFilePath := strings.TrimSuffix(string(models.Path(src)), ".go") + "_nxt_unit_test.go"
	suExist := atghelper.IsFileExist(suFilePath)
	var suImports []*instrumentation.Import
	if suExist {
		suImports, err = instrumentation.GetImportsInfosFromFile(suFilePath)
		if err == nil {
			targetImports = append(targetImports, suImports...)
			opt.ExistTestImport = append(opt.ExistTestImport, suImports...)
		}
	}

	// All imports will be added to be duplicated manager. Inside the duplicated manager, it will remove the duplicate
	// packages.
	recordExistingImports(opt, h.Imports)
	functions, err := GetFunctions(sr, opt)
	if err != nil {
		return nil, fmt.Errorf("[generateTest] GetFunctions has error")
	}
	funcs := testableFuncs(functions, opt.Only, opt.Exclude, opt.Exported, tf)

	// TODO(siwei.wang): remove this global variable, because it is an ugly implementation
	if len(funcs) != 0 {
		temp := fmt.Sprint(funcs[0].TestName())
		// This will bring some problem. It will run the existing test. However, because our client does not
		// have the existing test. I assume it only run the test for the su.
		atgconstant.FinalSuiteTestName = temp[0 : len(temp)-7]
	}
	// Record its original imports
	h.OriginalImports = duplicatedManagerToImports(opt.Ctx)
	// Add the imports
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("gomonkeyv2", "github.com/agiledragon/gomonkey/v2")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("atgconstant", "github.com/bytedance/nxt_unit/atgconstant")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("atghelper", "github.com/bytedance/nxt_unit/atghelper")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("contexthelper", "github.com/bytedance/nxt_unit/atghelper/contexthelper")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("mockfunc", "github.com/bytedance/nxt_unit/codebuilder/mock")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "github.com/bytedance/nxt_unit/staticcase")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "golang.org/x/tools/imports")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("variablecard", "github.com/bytedance/nxt_unit/codebuilder/variablecard")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "testing")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("duplicatepackagemanager", "github.com/bytedance/nxt_unit/manager/duplicatepackagemanager")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "github.com/smartystreets/goconvey/convey")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "fmt")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "strings")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "syscall")
	duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet("", "context")

	// Convert the imports to its original imports
	ssaFunctionMap, exist := contexthelper.GetSetupFuncMap(opt.Ctx)
	if !exist {
		ssaFunctionMap = make(map[string]setup.Functions, 0)
	}
	h.Imports = duplicatedManagerToImports(opt.Ctx)
	for _, fun := range funcs {
		data, ok := opt.RecordRow[fun.FullName()]
		if ok {
			fun.RowData = data
		}
		ssaFunctionInfo, exist := ssaFunctionMap[fun.FullName()]
		if exist {
			fun.ContainAnonFuncs = len(ssaFunctionInfo.TestFunction.Function.AnonFuncs) + 1
		} else {
			fun.ContainAnonFuncs = 1
		}
	}
	if len(funcs) == 0 {
		return nil, fmt.Errorf("[generateTest] no need to geneate test")
	}
	mocks := map[string]map[string]int{}
	builder := map[string][]string{}
	initBuilder := make([]string, 0)
	middleCodeBuilder := make([]string, 0)
	// Create the mock statement and also record the package
	switch opt.TestMode {
	case atgconstant.MiddleCode:
		mocks = GetAllMock(opt.Ctx, opt.UseMockType)
		builder = GetSpecialValueBuilder(opt.Ctx)
		initBuilder, middleCodeBuilder = GetGlobalValueBuilder(opt.Ctx)
		// picks := PickStructField(opt.Ctx)
		// initBuilder = append(initBuilder, picks...)
	case atgconstant.BaseTest:
		mocks = GetBaseMockByTempleType(opt.Ctx, opt.UseMockType)
		initBuilder, _ = GetGlobalValueBuilder(opt.Ctx)
	case atgconstant.FinalTest:
		if opt.UseMockMap != nil && len(opt.UseMockMap) != 0 {
			mocks = opt.UseMockMap
		}
	}

	options := output.Options{
		PrintInputs:    opt.PrintInputs,
		Subtests:       opt.Subtests,
		Parallel:       opt.Parallel,
		Named:          opt.Named,
		Template:       opt.Template,
		TemplateDir:    opt.TemplateDir,
		TemplateParams: opt.TemplateParams,
		TemplateData:   opt.TemplateData,
		Mocks:          mocks,
		Builders:       builder,
		GlobalInit:     initBuilder,
		MiddleCodeInit: middleCodeBuilder,
		TestCaseNum:    opt.TestCaseNum,
		FilePath:       opt.FilePath,
		TestMode:       opt.TestMode,
		Uid:            opt.Uid,
		UseMockType:    opt.UseMockType,
		Ctx:            opt.Ctx,
	}
	b, err := options.Process(h, funcs)
	if err != nil {
		return nil, fmt.Errorf("the error belongs to %w, the detail is %v", logextractor.MiddleCodeGenerateError, err.Error())
	}
	return &GeneratedTest{
		Path:      testPath,
		Functions: funcs,
		Output:    b,
	}, nil
}

// We have two scenarios
// 1. generate the middle code: using the funcsMap
// 2. generate the final code: using the rowData
func GetFunctions(sr *goparser.Result, opt *Options) ([]*models.Function, error) {
	functions := make([]*models.Function, 0)
	if opt.RecordRow != nil {
		for key := range opt.RecordRow {
			for _, tempFunction := range sr.Funcs {
				if key == tempFunction.FullName() {
					functions = append(functions, tempFunction)
				}
			}
		}
		return functions, nil
	}

	functionsMap, ok := contexthelper.GetSetupFuncMap(opt.Ctx)
	if !ok {
		return nil, fmt.Errorf("context doesnt have source function")
	}
	for name, f := range functionsMap {
		for _, tempFunction := range sr.Funcs {
			if tempFunction.Name == name {
				// If both of them have not the receiver, it means they are the same functions
				if f.TestFunction.Function.Signature.Recv() == nil && tempFunction.Receiver == nil {
					functions = append(functions, tempFunction)
				}

				// If both of them have receiver and their receiver are the same
				// (TODO: siwei.wang) test the below logic because it might incur bugs
				if f.TestFunction.Function.Signature.Recv() != nil && tempFunction.Receiver != nil {
					testFunctionReceiver := atghelper.GetTheReceiveNameFromSSA(f.TestFunction.Function.Signature.Recv().Type().String())
					temFunctionReceiver := tempFunction.ReceiverName()
					if testFunctionReceiver == temFunctionReceiver {
						functions = append(functions, tempFunction)
					}
				}
			}
		}
	}
	return functions, nil
}

func duplicatedManagerToImports(ctx context.Context) []*models.Import {
	res := make([]*models.Import, 0)
	duplicatepackagemanager.GetInstance(ctx).UniquePkgMap.Range(func(key, value interface{}) bool {
		// They are system package
		if key == value {
			res = append(res, &models.Import{Path: fmt.Sprintf("\"%v\"", value)})
		} else {
			res = append(res, &models.Import{Name: fmt.Sprintf("%v", key), Path: fmt.Sprintf("\"%v\"", value)})
		}
		return true
	})
	return res
}

func recordExistingImports(opt *Options, originalImports []*models.Import) {
	// first set existTest file import then orignalImport
	existTestImports := opt.ExistTestImport
	for index := range existTestImports {
		if existTestImports[index] != nil && existTestImports[index].Name != "." {
			duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet(
				strings.ReplaceAll(existTestImports[index].Name, "\"", ""),
				strings.ReplaceAll(existTestImports[index].Path, "\"", ""))
		}
	}

	// second,add originalImport
	for _, singleImport := range originalImports {
		duplicatepackagemanager.GetInstance(opt.Ctx).PutAndGet(
			strings.ReplaceAll(singleImport.Name, "\"", ""),
			strings.ReplaceAll(singleImport.Path, "\"", ""))
	}
}

func parseTestFile(p *goparser.Parser, testPath string, h *models.Header) (*models.Header, []string, error) {
	if !atghelper.IsFileExist(testPath) {
		return h, nil, nil
	}
	tr, err := p.Parse(testPath, nil)
	if err != nil {
		if err == goparser.ErrEmptyFile {
			// Overwrite empty test files.
			return h, nil, nil
		}
		return nil, nil, fmt.Errorf("Parser.Parse test file: %v", err)
	}
	var testFuncs []string
	for _, fun := range tr.Funcs {
		testFuncs = append(testFuncs, fun.Name)
	}
	tr.Header.Imports = append(tr.Header.Imports, h.Imports...)
	return h, testFuncs, nil
}

func testableFuncs(funcs []*models.Function, only, excl *regexp.Regexp, exp bool, testFuncs []string) []*models.Function {
	sort.Strings(testFuncs)
	var fs []*models.Function
	for _, f := range funcs {
		if isExcluded(f, excl) || isUnexported(f, exp) || !isIncluded(f, only) || isInvalid(f) {
			continue
		}
		fs = append(fs, f)
	}
	return fs
}

func isInvalid(f *models.Function) bool {
	if f.Name == "init" && f.IsNaked() {
		return true
	}
	return false
}

func isExcluded(f *models.Function, excl *regexp.Regexp) bool {
	return excl != nil && (excl.MatchString(f.Name) || excl.MatchString(f.FullName()))
}

func isUnexported(f *models.Function, exp bool) bool {
	return exp && !f.IsExported
}

func isIncluded(f *models.Function, only *regexp.Regexp) bool {
	return only == nil || only.MatchString(f.Name) || only.MatchString(f.FullName())
}

func contains(ss []string, s string) bool {
	if i := sort.SearchStrings(ss, s); i < len(ss) && ss[i] == s {
		return true
	}
	return false
}
