package instrumentation

import (
	"path"
	"runtime"
	"testing"

	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/imports"
)

func TestNewInstrumentation(t *testing.T) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../../atg/template/ifstatement.go")
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = "function"
	opt.FuncName = "CompareInteger"
	opt.ReceiverName = "TikTokConsumption"
	opt.FilePath = filePath

	o, _, err := NewInstrumentation(opt.FilePath, opt.FuncName, opt.Uid)
	assert.Equal(t, err, nil)
	_, err = imports.Process("", o, nil)
	assert.Equal(t, err, nil)
}
