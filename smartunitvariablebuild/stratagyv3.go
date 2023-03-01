package smartunitvariablebuild

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"github.com/bytedance/nxt_unit/atghelper/contexthelper"
	"github.com/bytedance/nxt_unit/manager/duplicatepackagemanager"
	"go/types"
	"path"
	"reflect"
	"strings"
)

func GetSpecialVariableV3(ctx context.Context, t reflect.Type) (reflect.Value, bool) {
	vtx, ok := contexthelper.GetVariableContext(ctx)
	if !ok {
		return reflect.ValueOf(""), false
	}
	if t == nil {
		return reflect.ValueOf(nil), true
	}
	switch t.Kind() {
	case reflect.Interface:
		pkgPath := ""
		switch t.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
			pkgPath = t.Elem().PkgPath()
		default:
			pkgPath = t.PkgPath()
		}
		fmt.Println("path ", pkgPath)
		if strings.Contains(t.String(), "io.Writer") {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("bytes", "")
			return reflect.ValueOf(bytes.NewBufferString("smartunit")), true
		}

		if strings.Contains(t.String(), "io.Reader") {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("strings", "")
			return reflect.ValueOf(strings.NewReader("smartunit")), true
		}

		if strings.HasPrefix(t.Name(), "error") {
			if vtx.CanBeNil && atghelper.RandomBool(atgconstant.SpecialValueBeNil) {
				return reflect.ValueOf(nil), true
			}
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("", "errors")
			return reflect.ValueOf(NewErr(ctx)), true
		}
		if strings.HasPrefix(t.String(), "context.Context") {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("context", "")
			return reflect.ValueOf(NewContext()), true
		}

	default:
		spv := ctx.Value("SpecialValueInjector")
		injector, ok := spv.(*SpecialValueInjector)
		if !ok {
			return reflect.ValueOf(""), false
		}
		value, exist := injector.Get(t.String())
		if exist {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet(path.Base(t.PkgPath()), t.PkgPath())
			return value, true
		}
		return reflect.ValueOf(""), false
	}
	return reflect.ValueOf(""), false
}

func RenderVariableV3(ctx context.Context, v reflect.Value) (string, bool) {
	t := v.Type()
	if t == nil {
		return "", false
	}
	pkgPath := ""
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		pkgPath = t.Elem().PkgPath()
	default:
		pkgPath = t.PkgPath()
	}

	switch t.Kind() {
	case reflect.Ptr:
		if pkgPath == "github.com/gin-gonic/gin" && v.String() == "<*gin.Context Value>" {
			pkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(path.Base("gin"), "github.com/gin-gonic/gin")
			return "&" + pkgName + ".Context{Request: &http.Request{URL: &url.URL{Path: \"test_path\",}}}", true
		}
	default:
		if pkgPath == "github.com/gin-gonic/gin" && v.String() == "<gin.Context Value>" {
			pkgName, _ := duplicatepackagemanager.GetInstance(ctx).PutAndGet(path.Base("gin"), "github.com/gin-gonic/gin")
			return pkgName + ".Context{Request: &http.Request{URL: &url.URL{Path: \"test_path\",}}}", true
		}
	}
	return "", false
}

type typeString string

const TypeString typeString = "TypeString"

// Check if there is a special variable builder for types.Type.
// Mainly used for special handling of interface{}.
// Return True means this type is ready to mock.
// It should update with GetSpecialVariableV3.
func SpecialVariableChecker(ctx context.Context, typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		return SpecialVariableChecker(ctx, t.Elem())
	case *types.Named:
		// Save type name to context. Cuz underlying type do not contain name
		ctx = context.WithValue(ctx, TypeString, t.String())
		return SpecialVariableChecker(ctx, t.Underlying())
	case *types.Interface:
		// Get type name from context
		str, _ := ctx.Value(TypeString).(string)
		if str == "io.Writer" {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("bytes", "")
			return true
		}
		if str == "io.Reader" {
			duplicatepackagemanager.GetInstance(ctx).PutAndGet("strings", "")
			return true
		}
		if str == "error" {
			return true
		}
		if str == "context.Context" {
			return true
		}
	default:
		return true
	}
	return false
}
