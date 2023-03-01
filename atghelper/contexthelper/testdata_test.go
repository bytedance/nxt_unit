package contexthelper

// import (
// 	"context"
// 	"path"
// 	"testing"
//
// 	"gotest.tools/assert"
//
// 	"github.com/bytedance/nxt_unit/atgconstant"
// 	"github.com/bytedance/nxt_unit/codebuilder/setup"
// )
//
// func TestGetTestContext(t *testing.T) {
// 	option := atgconstant.Options{
// 		FuncName:     "DatePrinter",
// 		FilePath:     path.Join(atgconstant.GOPATHSRC, "github.com/bytedance/nxt_unit/atg/template/atg.go"),
// 		Level:        1,
// 		Maxtime:      4,
// 		GenerateType: atgconstant.GAMode,
// 		MinUnit:      "file",
// 		Uid:          "Vector",
// 	}
// 	sourceFunc, err := setup.GetFunctions(option)
// 	if err != nil {
// 		t.Log(err)
// 		return
// 	}
// 	vtx := atgconstant.VariableContext{}
// 	ctx := context.Background()
// 	ctx = SetOption(ctx, option)
// 	ctx = SetSetupFunc(ctx, sourceFunc)
// 	ctx = SetBuilderVector(ctx, option.Uid)
// 	ctx = SetVariableContext(ctx, vtx)
// 	_, b := GetOption(ctx)
// 	_, d := GetSetupFunc(ctx)
// 	_, f := GetBuilderVector(ctx)
// 	_, h := GetVariableContext(ctx)
// 	assert.Equal(t, b, true)
// 	assert.Equal(t, d, true)
// 	assert.Equal(t, f, true)
// 	assert.Equal(t, h, true)
// }
