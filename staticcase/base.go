package staticcase

import (
	"context"
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	logerror "github.com/bytedance/nxt_unit/manager/logextractor"

	"golang.org/x/tools/go/ssa"

	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"

	"github.com/bytedance/nxt_unit/manager/lifemanager"

	"github.com/bytedance/nxt_unit/codebuilder/instrumentation"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/codebuilder/unitestframwork/testcase"
	"golang.org/x/tools/imports"
)

func NewTestsuiteEntry(ctx context.Context, targetFile string, lock atghelper.RWLocker) *testsuiteEntry {
	opt, _ := contexthelper.GetOption(ctx)
	return &testsuiteEntry{
		name:       opt.Uid,
		filePath:   opt.FilePath,
		lock:       lock,
		targetFile: targetFile,
	}
}

type testsuiteEntry struct {
	name      string
	closeFunc func() error
	filePath  string
	lock      atghelper.RWLocker
	// mocks      []string
	targetFile string
	// asserts    []string
}

func (t *testsuiteEntry) Length() int {
	return 0
}

func (t *testsuiteEntry) Name() string {
	return t.name
}

func (t *testsuiteEntry) Build(ctx context.Context) (string, error) {
	code, err := t.Code(ctx)
	if err != nil {
		return "", err
	}
	t.lock.RLock()
	defer func() {
		t.lock.RUnlock()
	}()
	opt, _ := contexthelper.GetOption(ctx)
	f := path.Join(filepath.Dir(opt.FilePath), t.Name()+"_middle_code.txt")
	if err := ioutil.WriteFile(f, code, atgconstant.NewFilePerm); err != nil {
		return "", fmt.Errorf("[VariableCard] has ioutil.WriteFile err: %v", err)
	}
	return f, nil
}

func (t *testsuiteEntry) Close() error {
	t.lock.Lock()
	defer func() {
		t.lock.Unlock()
	}()
	return t.closeFunc()
}

func (t *testsuiteEntry) Code(ctx context.Context) ([]byte, error) {
	opt, _ := contexthelper.GetOption(ctx)
	gt, err := GenerateTests(t.targetFile, &Options{
		// Only:        regexp.MustCompile(fmt.Sprintf("^%s$", opt.FuncName+opt.Uid)),
		Ctx: ctx,
		TestCaseNum: func() int {
			switch opt.Usage {
			case atgconstant.PluginMode:
				return atgconstant.PluginTestSuiteMaxSize
			default:
				return atgconstant.TestSuiteMaxSize
			}
		}(),
		Uid:         opt.Uid,
		FilePath:    opt.FilePath,
		TestMode:    atgconstant.MiddleCode,
		UseMockType: opt.UseMockType,
	})
	if err != nil {
		return nil, err
	}
	if gt == nil {
		return nil, errors.New("can't create test")
	}
	code := gt.Output
	code, err = imports.Process("", code, nil)
	if err != nil {
		return nil, err
	}
	return code, nil
}

func GetAllMock(ctx context.Context, useMockType int) map[string]map[string]int {
	functionMockMap := make(map[string]map[string]int, 0)
	if useMockType == atgconstant.UseNoMock {
		return functionMockMap
	}
	functionMap, _ := contexthelper.GetSetupFuncMap(ctx)
	for funcName, functions := range functionMap {
		ts, err := testcase.CreateTestCase(ctx, functions.TestFunction)
		if err != nil {
			continue
		}
		mock := make(map[string]int, 0)
		for _, stat := range ts.Statements {
			randomMock := ""
			switch stat.SpecialType {
			case "overpass":
				mockStatement := fmt.Sprintf("%s(mockfunc.OverPassMakeCall(smartUnitCtx,\"%s\",mockRender,%s).(%s));", stat.Expression, stat.Expression, stat.Expression, stat.FunctionType)
				randomMock = mockStatement
			default:
				makeMockFunc := fmt.Sprintf("mockfunc.MakeCall(smartUnitCtx,\"%s\",mockRender,%s,%d)", stat.Expression, stat.Expression, useMockType)
				switch useMockType {
				case atgconstant.UseGoMonkeyMock:
					patchName := fmt.Sprintf("%sPatch", atghelper.RandStringBytes(5))
					randomMock = fmt.Sprintf("%s := gomonkeyv2.ApplyFunc(%s,%s)\n\t\tdefer %s.Reset()", patchName, stat.Expression, makeMockFunc, patchName)
				case atgconstant.UseMockitoMock:
					randomMock = fmt.Sprintf("mockito.Mock(%s).To(%s).Build();", stat.Expression, makeMockFunc)
				}

			}

			mock[randomMock] = 1
		}
		functionMockMap[funcName] = mock
	}
	return functionMockMap
}

func GetBaseMockByTempleType(ctx context.Context, useMockType int) map[string]map[string]int {
	functionMap, _ := contexthelper.GetSetupFuncMap(ctx)
	functionMockMap := map[string]map[string]int{}

	for funcName, functions := range functionMap {
		ts, err := testcase.CreateTestCase(ctx, functions.TestFunction)
		if err != nil {
			continue
		}
		mock := make(map[string]int, 0)
		for _, stat := range ts.Statements {
			randomMock := ""
			switch stat.SpecialType {
			case "overpass":
				randomMock = fmt.Sprint("// Please fill out the overpass mock yourself \n", "// ", stat.Expression, "(", stat.FunctionType, "{\n", "// \treturn [please fill the return here]\n", "//})")
			default:
				switch useMockType {
				case atgconstant.UseMockUnknown:
					randomMock = fmt.Sprintf("mockito.Mock(%s).Return().Build();", stat.Expression)
				case atgconstant.UseNoMock:
					// do nothing
				case atgconstant.UseMockitoMock:
					randomMock = fmt.Sprintf("mockito.Mock(%s).Return().Build();", stat.Expression)
				case atgconstant.UseGoMonkeyMock:
					patchName := fmt.Sprintf("%sPatch", atghelper.RandStringBytes(5))
					randomMock = fmt.Sprintf("%s := gomonkeyv2.ApplyFuncReturn(%s,\"%s\")\n\t\tdefer %s.Reset()", patchName, stat.Expression, "user fill function return output", patchName)
				}
			}
			if randomMock != "" {
				mock[randomMock] = 1
			}
		}
		functionMockMap[funcName] = mock
	}
	return functionMockMap
}

func GenerateBaseTest(ctx context.Context) error {
	opt, _ := contexthelper.GetOption(ctx)
	options := &Options{
		// Only:        regexp.MustCompile(fmt.Sprintf("^%s$", opt.FuncName+opt.Uid)),
		Ctx:             ctx,
		TestCaseNum:     atgconstant.TestSuiteMaxSize * 2,
		Uid:             opt.Uid,
		FilePath:        opt.FilePath,
		TestMode:        atgconstant.BaseTest,
		ExistTestImport: make([]*instrumentation.Import, 0),
		UseMockType:     opt.UseMockType,
	}

	// check _test.go exist,if exist options set imports from existing test file
	targetFilePath := getTempFilePath(opt, "_test.go")
	tgtExist := atghelper.IsFileExist(targetFilePath)
	var targetImports []*instrumentation.Import
	var err error
	if tgtExist {
		targetImports, err = instrumentation.GetImportsInfosFromFile(targetFilePath)
		if err != nil {
			return fmt.Errorf("the error belongs to %w, the detail is %v", logerror.GenerateTestTemplateCannotParseImportError, err.Error())
		}
		options.ExistTestImport = targetImports
	}

	// check _smart_unit_test.go exist, if exist options set imports from existing smart unit test file
	suFilePath := getTempFilePath(opt, "_nxt_unit_test.go")
	suExist := atghelper.IsFileExist(suFilePath)
	var suImports []*instrumentation.Import
	if suExist {
		suImports, err = instrumentation.GetImportsInfosFromFile(suFilePath)
		if err != nil {
			return fmt.Errorf("the error belongs to %w, the detail is %v", logerror.GenerateTestTemplateCannotParseImportError, err.Error())
		}
		targetImports = append(targetImports, suImports...)
		options.ExistTestImport = append(options.ExistTestImport, suImports...)
	}

	gt, err := GenerateTests(opt.FilePath, options)
	if err != nil {
		return fmt.Errorf("the error belongs to %w, the detail is %v", logerror.GenerateTestTemplateInternalError, err.Error())
	}
	code := gt.Output
	code, err = imports.Process("", code, nil)
	if err != nil {
		return fmt.Errorf("the error belongs to %w, the detail is %v", logerror.GenerateTemplateImportError, err.Error())
	}
	testName := strings.ReplaceAll(filepath.Base(opt.FilePath), ".go", "_base_test.go")
	testFile := path.Join(filepath.Dir(opt.FilePath), testName)
	if err = ioutil.WriteFile(testFile, code, atgconstant.NewFilePerm); err != nil {
		return err
	}
	// template rename based on naming conventions
	if !tgtExist {
		err = os.Rename(testFile, targetFilePath)
		if err != nil {
			return fmt.Errorf("the error belongs to %w, the detail is %v", logerror.GenerateTestTemplateRenameError, err.Error())
		}
		return nil
	}
	err = instrumentation.Concatenate(testFile, targetFilePath, targetImports)
	if err != nil {
		return fmt.Errorf("the error belongs to %w, the detail is %v", logerror.MergeTestConflictError, err.Error())
	}
	lifemanager.Closer.SetClose(func() {
		os.Remove(testFile)
	})
	return nil
}

func GetSpecialValueBuilder(ctx context.Context) map[string][]string {
	functionMap, _ := contexthelper.GetSetupFuncMap(ctx)
	// TODO: add constructor func map to context
	constructorMap, _ := contexthelper.GetConstructorFuncMap(ctx)
	functionBuilder := map[string][]string{}
	for funcName, functions := range functionMap {
		InjectorBuilder := sync.Once{}
		builders := make([]string, 0)
		initInjector := func() {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("smartunitvariablebuild", "github.com/bytedance/nxt_unit/smartunitvariablebuild")
			builders = append(builders, "injector := smartunitvariablebuild.NewSpecialValueInjector();")
			builders = append(builders, "smartUnitCtx = context.WithValue(smartUnitCtx,\"SpecialValueInjector\",injector)")
		}
		// add special func for ad/starry_assembler
		for i := 0; i < functions.TestFunction.Function.Signature.Params().Len(); i++ {
			in := functions.TestFunction.Function.Signature.Params().At(i)
			constructors, ok := constructorMap[in.Type().String()]
			switch {
			// if found constructors of param's type. Try to use constructor
			case ok:
				InjectorBuilder.Do(initInjector)
				constructor := constructors[rand.Intn(len(constructors))]
				storePkgName, _ := duplicatepackagemanager.GetInstance(ctx).Put(constructor.Pkg.Pkg.Name(), constructor.Pkg.Pkg.Path())
				// injector.SetBuilder(SetName(), "SetName()")
				// TODO: support muti params constructor
				if storePkgName == "" {
					// same pkg
					builders = append(builders, fmt.Sprintf("injector.SetBuilder(%s, \"%s()\")", fmt.Sprintf("%s()", constructor.Name()), constructor.Name()))
				} else {
					// other pkg
					builders = append(builders, fmt.Sprintf("injector.SetBuilder(%s, \"%s\")", fmt.Sprintf("%s.%s()", storePkgName, constructor.Name()), fmt.Sprintf("%s.%s()", storePkgName, constructor.Name())))
				}
			}
		}
		functionBuilder[funcName] = builders
	}
	return functionBuilder
}

func GetGlobalValueBuilder(ctx context.Context) ([]string, []string) {
	functionMap, _ := contexthelper.GetSetupFuncMap(ctx)
	initBuilder := make([]string, 0)
	middleCodeBuilder := make([]string, 0)
	var program *ssa.Program
	for _, functions := range functionMap {
		program = functions.TestFunction.Program.Prog
		for _, importsPkg := range program.AllPackages() {
			// gopkg has it's own init
			// if pkg has it's init functions, we don't mock global var
			if importsPkg.Func("init").Object() != nil {
				continue
			}
			for _, member := range importsPkg.Members {
				vt, ok := member.(*ssa.Global)
				if ok {
					if vt != nil && token.IsExported(vt.Name()) && CheckValueExis(ctx, functions.TestFunction.Function, vt.String(), vt.Name()) {
						entireDesc := vt.String()
						pointIndex := strings.LastIndex(entireDesc, ".")
						statement := path.Base(vt.String())
						if pointIndex != -1 {
							pkgPath := entireDesc[:pointIndex]
							stateIndex := strings.LastIndex(statement, ".")
							pkgName := statement[:stateIndex]
							varName := statement[stateIndex+1:]
							realPkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, pkgPath)
							if realPkgName != "" {
								statement = fmt.Sprintf("%s.%s", realPkgName, varName)
							} else {
								statement = varName
							}
						}

						if vt.Object() == nil {
							continue
						}
						switch v := vt.Object().Type().Underlying().(type) {
						case *types.Interface:

						case *types.Pointer:
							// get pkgPath & Name
							str := v.Elem().String()
							split := strings.LastIndex(str, ".")
							if split == -1 {
								builder := fmt.Sprintf("%s = new(%s);// mockito value of this type  ", statement, str)
								initBuilder = append(initBuilder, builder)
								middleCodeBuilder = append(middleCodeBuilder, builder)
							} else {
								pkgPath := str[:split]
								pkgName := atghelper.GetPkgName(pkgPath)
								pkgName, _ = duplicatepackagemanager.GetInstance(ctx).PutAndGet(pkgName, pkgPath)
								builder := ""
								if pkgName != "" {
									// check if it is exported
									if !token.IsExported(str[split+1:]) {
										continue
									}
									builder = fmt.Sprintf("%s = new(%s.%s);// mockito value of this type  ", statement, pkgName, str[split+1:])
								} else {
									builder = fmt.Sprintf("%s = new(%s);// mockito value of this type  ", statement, str[split+1:])
								}
								if builder != "" {
									initBuilder = append(initBuilder, builder)
								}
								middleCodeBuilder = append(middleCodeBuilder, builder)
							}
						default:
							builder := fmt.Sprintf("// unsupport init %v this of  type: %v  ", statement, vt.Type().String())
							initBuilder = append(initBuilder, builder)
							middleCodeBuilder = append(middleCodeBuilder, builder)
						}
					}
				}

			}
		}
	}
	if program == nil {
		return initBuilder, middleCodeBuilder
	}

	// It means, we get the source function.
	functions, ok := contexthelper.GetSetupFuncMap(ctx)
	if !ok {
		return initBuilder, middleCodeBuilder
	}
	optInfo, ok := contexthelper.GetOption(ctx)
	if !ok {
		return initBuilder, middleCodeBuilder
	}
	function, exist := functions[optInfo.FuncName]
	if !exist {
		return initBuilder, middleCodeBuilder
	}

	if function.TestFunction.Program == nil {
		return initBuilder, middleCodeBuilder
	}
	return initBuilder, middleCodeBuilder
}

func PickStructField(ctx context.Context) []string {
	pickBuilder := make([]string, 0)
	picks := sync.Once{}
	functionMap, _ := contexthelper.GetSetupFuncMap(ctx)
	for _, functions := range functionMap {
		if functions.DateSteam != nil {
			for key := range functions.DateSteam {
				picks.Do(func() {
					pickBuilder = append(
						pickBuilder,
						"picks := map[string]struct{}{}",
						"smartUnitCtx = context.WithValue(smartUnitCtx,\"ValuePicker\",picks)",
					)
				})
				pickBuilder = append(pickBuilder,
					fmt.Sprintf("picks[\"%s\"] = struct{}{}", key),
				)
			}
		}
	}
	return pickBuilder
}

func CheckValueExis(ctx context.Context, p *ssa.Function, globalValue, globalName string) bool {
	for _, block := range p.Blocks {
		for _, instr := range block.Instrs {
			if strings.Contains(instr.String(), globalValue) || strings.Contains(instr.String(), globalName) {
				return true
			}
		}
	}
	// for _, Atype := range p.RuntimeTypes(){
	// 	if strings.Contains(Atype.String(),"oneP"){
	//
	// 	}
	// 	if Atype.String() == globalValue{
	// 		fmt.Println(Atype.String(), globalValue)
	// 		return true
	// 	}
	// }
	return false
}
