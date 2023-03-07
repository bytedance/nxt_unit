package setup

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/stretchr/testify/assert"
)

var baseDir string

func TestMain(m *testing.M) {
	_, baseDir, _, _ = runtime.Caller(0)
	m.Run()
}

// This case will find two functions of consume
func TestGetFunctions3(t *testing.T) {
	filePath := path.Join(path.Dir(baseDir), "../../atg/template/receiver.go")
	opt := atgconstant.Options{
		FilePath:     filePath,
		MinUnit:      "function",
		DebugMode:    true,
		Usage:        "plugin",
		ReceiverName: "",
		FuncName:     "Consume",
	}
	_, err := GetFunctions(opt)
	assert.Nil(t, err)
}

func TestGetFunctions4(t *testing.T) {
	filePath := path.Join(path.Dir(baseDir), "../../atg/template/receiver.go")
	opt := atgconstant.Options{
		FilePath:     filePath,
		MinUnit:      "function",
		DebugMode:    true,
		Usage:        "plugin",
		ReceiverName: "Hello",
		FuncName:     "Consume",
	}
	_, err := GetFunctions(opt)
	assert.NotNil(t, err)
}

func TestGetReferences(t *testing.T) {
	filePath := path.Join(path.Dir(baseDir), "../../atg/template/dataanalysis.go")
	opt := atgconstant.Options{
		FilePath:  filePath,
		MinUnit:   "function",
		DebugMode: true,
		Usage:     "plugin",
		FuncName:  "QueryData",
	}
	f, err := GetFunctions(opt)
	if err != nil {
		t.Fatal(err)
	}

	for key := range f.DateSteam {
		t.Log(key)
	}
}

func TestGetReferencesPointer(t *testing.T) {
	filePath := path.Join(path.Dir(baseDir), "../../atg/template/dataanalysis.go")
	opt := atgconstant.Options{
		FilePath:  filePath,
		MinUnit:   "function",
		DebugMode: true,
		Usage:     "plugin",
		FuncName:  "QueryDataPointer",
	}
	f, err := GetFunctions(opt)
	if err != nil {
		t.Fatal(err)
	}

	for key := range f.DateSteam {
		if strings.Contains(fmt.Sprint(key), "*") {
			t.Fatal(err)
		}
		t.Log(key)
	}
}
