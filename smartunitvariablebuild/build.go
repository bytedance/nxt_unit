package smartunitvariablebuild

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
)

func NewErrcode_i18n() *status {
	if atghelper.RandomBool(atgconstant.SpecialValueBeNil) {
		return &status{
			code: 3001002,
		}
	}
	return &status{
		code: 3001210,
	}
}

type status struct {
	code int32
}

func (s *status) Code() int32 {
	return s.code
}
func (s *status) Msg(language string) string {
	return "test error"
}
func (s *status) StarlingKey() string {
	return "test error"
}

func (s *status) ValueToCode() string {
	switch s.code {
	case 3001002:
		return "aweme.SUCCESS"
	}
	return "aweme.ERR_SERVICE_INTERNAL"
}

type suErr struct {
	Ctx context.Context
	error
}

func (s *suErr) ValueToCode() string {
	pkgName, _ := duplicatepackagemanager.GetInstance(s.Ctx).PutAndGet("", "errors")
	return fmt.Sprintf("%s.New(\"smart unit\")", pkgName)
}

func NewErr(ctx context.Context) error {
	return &suErr{ctx, errors.New("smart unit")}
}

type suContext struct {
	context.Context
}

func (s *suContext) ValueToCode() string {
	return `context.Background()`
}

func NewContext() context.Context {
	return &suContext{
		context.Background(),
	}
}
