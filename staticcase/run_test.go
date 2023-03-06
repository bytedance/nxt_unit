package staticcase

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/manager/lifemanager"
	"github.com/stretchr/testify/assert"
)

// run in any path, it will Generate  middle test code in parallel by rendering to text file
func TestWork(t *testing.T) {
	defer func() {
		lifemanager.Closer.Close()
	}()
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../atg/template")
	err := Work(filePath, "", 2)
	if err != nil {
		t.Fatal(err)
	}
	err = WorkToChangeGo(filePath, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckEmptyFileAndThenRemove(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	value := "hello\n"
	err = ioutil.WriteFile(path.Join(dir, "testcase1"), []byte(value), 0644)
	defer os.Remove(path.Join(dir, "testcase1"))
	res := CheckEmptyFileAndThenRemove(path.Join(dir, "testcase1"), false)
	assert.NotNil(t, res)
}

func Test_getTempFilePath(t *testing.T) {
	type args struct {
		opt    atgconstant.Options
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				opt:    atgconstant.Options{DirectoryPath: "test1", FilePath: "test1/file1.go"},
				suffix: fmt.Sprint("_nxt_unit_test_", "abc", ".go"),
			},
			want: "test1/file1_nxt_unit_test_abc.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTempFilePath(tt.args.opt, tt.args.suffix); got != tt.want {
				t.Errorf("getTempFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTestOnPlugin(t *testing.T) {
	defer func() {
		lifemanager.Closer.Close()
	}()
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.Usage = atgconstant.PluginMode
	opt.MinUnit = atgconstant.FileMode
	opt.DirectoryPath = filepath.Dir(opt.FilePath)
	err := WorkForPlugin(opt)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = GenerateTestForPlugin(opt)
	if err != nil {
		ctx, err := Init(&opt)
		if err != nil {
			fmt.Printf("Sorry, we cannot generate the test for you, the error is %v\n", err.Error())
		}
		err = GenerateBaseTest(ctx)
		if err != nil {
			fmt.Printf("Sorry, we cannot generate base test %v\n", err.Error())
		}
		fmt.Println("Successfully generate the unit test template!")
		return
	}
}

func TestGenerateTestOnSplitFunctionMode(t *testing.T) {
	defer func() {
		lifemanager.Closer.Close()
	}()
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.Usage = atgconstant.SplitFunctionMode
	opt.MinUnit = atgconstant.FileMode
	opt.DirectoryPath = filepath.Dir(opt.FilePath)
	err := WorkForSplitFunction(opt.DirectoryPath, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
}
