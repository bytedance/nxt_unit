/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package duplicatepackagemanager

import (
	"context"
	"fmt"
	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/bytedance/nxt_unit/atghelper"
	"reflect"
	"strings"
	"sync"
)

func NewPackageManager() *DuplicatePackageManager {
	return &DuplicatePackageManager{
		CheckPkgMap:   make(map[string]string, 0),
		GlobalBuilder: []string{},
		MiddleBuilder: []string{},
	}
}

var packageManagerKey struct{}

func SetInstance(ctx context.Context) context.Context {
	return context.WithValue(ctx, packageManagerKey, NewPackageManager())
}

func GetInstance(ctx context.Context) *DuplicatePackageManager {
	value := ctx.Value(packageManagerKey)
	manager, ok := value.(*DuplicatePackageManager)
	if !ok {
		return GetInstanceA()
	}
	return manager
}

type DuplicatePackageManager struct {
	CurrentPkgName string
	CurrentPkgPath string
	relativePath   string
	tempImports    atgconstant.ImportInfo
	GlobalBuilder  []string
	MiddleBuilder  []string
	UniquePkgMap   sync.Map          //key:pkg_name,key:pkg_path
	CheckPkgMap    map[string]string //key:pkg_path,key:pkg_name
}

var d *DuplicatePackageManager

func init() {
	Init()
}

func Init() {
	d = &DuplicatePackageManager{CheckPkgMap: make(map[string]string, 0)}
}

// GetInstanceA - get singleton instance pre-initialized
func GetInstanceA() *DuplicatePackageManager {
	return d
}

func (d *DuplicatePackageManager) GetPkgName(pkgPath string) string {
	var pkgName string
	d.UniquePkgMap.Range(func(key, value interface{}) bool {
		valueStr, ok := value.(string)
		if ok {
			if valueStr == pkgPath {
				pkgName, ok = key.(string)
				if ok {
					return false
				}
				return true
			}
		}
		return true
	})
	return pkgName
}

func (d *DuplicatePackageManager) Put(pkgName string, pkgPath string) (string, bool) {
	if pkgName == "" {
		pkgName = atghelper.GetPkgName(pkgPath)
	}
	storePkgName, success := d.LoadOrStoreImportPkg(pkgName, pkgPath)
	return storePkgName, success
}

func (d *DuplicatePackageManager) SetInitBuilder(builder string) {
	d.GlobalBuilder = append(d.GlobalBuilder, builder)
}

func (d *DuplicatePackageManager) GetInitBuilder() []string {
	return d.GlobalBuilder
}

func (d *DuplicatePackageManager) SetMiddleBuilder(builder string) {
	d.MiddleBuilder = append(d.MiddleBuilder, builder)
}

func (d *DuplicatePackageManager) GetMiddleBuilder() []string {
	return d.MiddleBuilder
}

// PutAndGet PkgName are supposed to be not empty. if that's empty, let's parse its pkg name.
func (d *DuplicatePackageManager) PutAndGet(pkgName string, pkgPath string) (string, string) {
	if pkgPath == "" {
		return "", ""
	}

	// Self Import
	if pkgPath == d.relativePath {
		return "", ""
	}
	if pkgName == "" {
		pkgName = atghelper.GetPkgName(pkgPath)
	}

	// pkgName should not be in a format of "a-b"
	if strings.Contains(pkgName, "-") {
		pkgName = strings.Replace(pkgName, "-", "_", -1)
	}

	storePkgName, success := d.LoadOrStoreImportPkg(pkgName, pkgPath)
	if !success {
		fmt.Printf(" Duplicated manager [PutAndGet] LoadOrStoreImportPkg is not correct, the pkg name is %v,storePkgName %v, the pkg path is %v\n", pkgName, pkgPath, storePkgName)
	}
	// Special Logic
	if IsPathShouldBeRemoved(pkgPath) {
		return pkgName, ""
	}
	return storePkgName, pkgPath
}

func (d *DuplicatePackageManager) SetRelativePath(value interface{}) {
	d.relativePath = reflect.TypeOf(value).PkgPath()
}

func (d *DuplicatePackageManager) SetRelativeString(path string) {
	d.relativePath = path
}

func (d *DuplicatePackageManager) SetImportInfo(temp atgconstant.ImportInfo) {
	d.tempImports = temp
}

func (d *DuplicatePackageManager) GetImportInfo() atgconstant.ImportInfo {
	return d.tempImports
}

func (d *DuplicatePackageManager) RelativePath() string {
	return d.relativePath
}

// The bool result is true if the value load or store success, false if fail
func (d *DuplicatePackageManager) LoadOrStoreImportPkg(pkgName, pkgPath string) (string, bool) {
	// check pkg path is exist
	existPkgName, exist := d.IsPkgPathExit(pkgPath)
	// PkgMap exist renamed pkgName by smart unit
	if exist {
		orignPkgName := atghelper.GetPkgName(pkgPath)
		if strings.Contains(existPkgName, "SmartU") || (existPkgName != pkgName && existPkgName != orignPkgName) {
			loadValue, found := d.UniquePkgMap.Load(existPkgName)
			loadPath, ok := loadValue.(string)
			if found && ok && loadPath == pkgPath {
				return existPkgName, true
			}
		}

	}
	storeValue, loaded := d.UniquePkgMap.LoadOrStore(pkgName, pkgPath)
	storePath, ok := storeValue.(string)
	if !ok {
		fmt.Printf("pkgName %v,storeValue: %v is not string", pkgName, storeValue)
		storePath = ""
	}
	// exist pkgName
	if loaded {
		// check same pkgName and pkgPath
		if storePath == pkgPath {
			return pkgName, true
		} else {
			newPkgName := pkgName + "SmartU" + atghelper.RandStringBytes(5)
			_, success := d.LoadOrStoreImportPkg(newPkgName, pkgPath)
			if success {
				d.CheckPkgMap[pkgPath] = newPkgName
			}
			return newPkgName, success
		}
	} else {
		d.CheckPkgMap[pkgPath] = pkgName
	}
	return pkgName, true
}

// The bool result is true if the value load or store success, false if fail
func (d *DuplicatePackageManager) UniquePkgMapLen() int {
	total := 0
	d.UniquePkgMap.Range(func(key, value interface{}) bool {
		total++
		return true
	})
	return total
}

// The bool result is true if the value load or store success, false if fail
func (d *DuplicatePackageManager) IsPkgPathExit(pkgPath string) (string, bool) {
	pkgName, ok := d.CheckPkgMap[pkgPath]
	return pkgName, ok
}

// warning:  Use with caution
func (d *DuplicatePackageManager) ClearUniquePkgMap() {
	d.UniquePkgMap.Range(func(key, value interface{}) bool {
		d.UniquePkgMap.Delete(key)
		return true
	})
}

func IsPathShouldBeRemoved(path string) bool {
	if strings.Contains(path, "github.com/pkg/errors") {
		return true
	}
	return false
}
