package duplicatepackagemanager

import (
	"github.com/bytedance/nxt_unit/atgconstant"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, d)
}

// Test Self Import with empty pkg list
func TestDuplicatePackageManager_PutAndGet_Case1(t *testing.T) {
	Init()
	atgconstant.PkgRelativePath = "github/pkg/env"
	GetInstanceA().SetRelativeString(atgconstant.PkgRelativePath)
	pkgName, pkgPath := GetInstanceA().PutAndGet("env", "github/pkg/env")
	assert.Empty(t, pkgName)
	assert.Empty(t, pkgPath)
	atgconstant.PkgRelativePath = ""
	d.ClearUniquePkgMap()
}

// Test Self Import with the same element in the pkg list
func TestDuplicatePackageManager_PutAndGet_Case2(t *testing.T) {
	Init()
	GetInstanceA().Put("env", "github/pkg/env")
	atgconstant.PkgRelativePath = "github/pkg/env"
	GetInstanceA().SetRelativeString(atgconstant.PkgRelativePath)
	pkgName, pkgPath := d.PutAndGet("env", "github/pkg/env")
	assert.Equal(t, d.UniquePkgMapLen(), 1)
	assert.Empty(t, pkgName)
	assert.Empty(t, pkgPath)
	d.ClearUniquePkgMap()
	atgconstant.PkgRelativePath = ""
}

// Test Self Import with the different element in the pkg list
func TestDuplicatePackageManager_PutAndGet_Case3(t *testing.T) {
	Init()
	GetInstanceA().Put("hello", "github/pkg/env")
	atgconstant.PkgRelativePath = "github/pkg/env"
	GetInstanceA().SetRelativeString(atgconstant.PkgRelativePath)
	pkgName, pkgPath := d.PutAndGet("env", "github/pkg/env")
	assert.Equal(t, d.UniquePkgMapLen(), 1)
	assert.Empty(t, pkgName)
	assert.Empty(t, pkgPath)
	d.ClearUniquePkgMap()
	atgconstant.PkgRelativePath = ""
}

// Test Special logic
func TestDuplicatePackageManager_PutAndGet_Case4(t *testing.T) {
	Init()
	pkgName, pkgPath := d.PutAndGet("env", "github.com/pkg/errors")
	assert.Equal(t, pkgName, "env")
	assert.Empty(t, pkgPath)
	assert.Equal(t, 1, d.UniquePkgMapLen())
	d.ClearUniquePkgMap()
}

// Test different package in the pkg list
func TestDuplicatePackageManager_PutAndGet_Case6(t *testing.T) {
	Init()
	GetInstanceA().Put("env", "github/pkg/env")
	pkgName, pkgPath := d.PutAndGet("env", "github/pkg2/env")
	assert.Equal(t, d.UniquePkgMapLen(), 2)
	assert.Contains(t, pkgName, "env")
	assert.Equal(t, pkgPath, "github/pkg2/env")
	d.ClearUniquePkgMap()
}

// Test same path in the pkg PUR
func TestDuplicatePackageManager_Put_Case(t *testing.T) {
	Init()
	pkgName, success := GetInstanceA().Put("utils", "github/pkg/utils")
	getName := d.GetPkgName("github/pkg/utils")
	assert.Equal(t, true, success)
	assert.Equal(t, pkgName, getName)
}
