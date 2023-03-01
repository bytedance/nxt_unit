package instrumentation

import (
	"fmt"
	"path"
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"golang.org/x/tools/imports"
)

func TestNewInstrumentation(t *testing.T) {
	ctx := contexthelper.GetTestContext()
	opt, _ := contexthelper.GetOption(ctx)
	opt.MinUnit = "function"
	opt.FuncName = "CompareInteger"
	opt.ReceiverName = "TikTokConsumption"
	opt.FilePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atg/template/ifstatement.go")

	o, _, err := NewInstrumentation(opt.FilePath, opt.FuncName, opt.Uid)
	if err != nil {
		t.Fatal(err)
	}
	code, err := imports.Process("", o, nil)
	fmt.Println(string(code))

}
