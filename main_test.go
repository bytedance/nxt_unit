package main

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/bytedance/nxt_unit/manager/lifemanager"
	"github.com/smartystreets/goconvey/convey"

	"github.com/bytedance/nxt_unit/atgconstant"
)

func TestPlugin_Function(t *testing.T) {
	convey.Convey("GlobalValue", t, func() {
		// defer func() {
		// 	lifemanager.Closer.Close()
		// }()
		*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
		*minUnit = atgconstant.MinUnit
		*funcName = "GoodFunc"
		*debugMode = true
		*usage = atgconstant.PluginMode
		*UseMockType = 3
		// patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		// defer patch.Reset()
		Plugin()
	})
}

func TestPlugin_InterfaceFunction(t *testing.T) {
	convey.Convey("GlobalValue", t, func() {
		defer func() {
			lifemanager.Closer.Close()
		}()
		*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
		*minUnit = atgconstant.MinUnit
		*funcName = "PrintInter"
		*debugMode = true
		*usage = atgconstant.PluginMode
		patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		Plugin()
	})
}

func TestPlugin_Functions(t *testing.T) {
	convey.Convey("GlobalValue", t, func() {
		defer func() {
			lifemanager.Closer.Close()
		}()
		*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
		*minUnit = atgconstant.MinUnit
		*funcName = "GoodFunc"
		*debugMode = true
		*usage = atgconstant.PluginMode
		// *ReceiverName = "Iss"
		// *ReceiverIsStar = true
		// patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		// defer patch.Reset()
		Plugin()
	})
}

func TestPlugin_File(t *testing.T) {
	convey.Convey("Plugin_File", t, func() {
		defer func() {
			lifemanager.Closer.Close()
		}()
		*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
		*minUnit = atgconstant.FileMode
		*debugMode = true
		*usage = atgconstant.PluginMode
		*UseMockType = 1
		patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		Plugin()
	})
}

func TestTemplate_atg(t *testing.T) {
	convey.Convey("TestTemplate_atg", t, func() {
		defer func() {
			lifemanager.Closer.Close()
		}()
		*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
		*minUnit = atgconstant.FileMode
		*funcName = "GoodFunc"
		*debugMode = true
		*usage = atgconstant.PluginQMode
		patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		Template()
	})
}

func TestTemplate_Fail(t *testing.T) {
	convey.Convey("TestTemplate_atg", t, func() {
		defer func() {
			lifemanager.Closer.Close()
		}()
		*filePath = "/Users/bytedance/workspace/hertz/pkg/protocol/http2/server.go"
		*minUnit = atgconstant.MinUnit
		*funcName = "registerConn"
		*ReceiverName = "serverInternalState"
		*ReceiverIsStar = true
		*debugMode = true
		*usage = atgconstant.PluginMode
		patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		Template()
	})
}

func TestTemplate_server(t *testing.T) {
	convey.Convey("TestTemplate_server", t, func() {
		defer func() {
			lifemanager.Closer.Close()
		}()
		*filePath = "/Users/bytedance/workspace/hertz/pkg/protocol/http2/server.go"
		*minUnit = atgconstant.FileMode
		*debugMode = true
		*usage = atgconstant.PluginMode
		patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
		defer patch.Reset()
		Template()
	})
}

func TestDeferCallGraph_Function(t *testing.T) {
	*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
	*minUnit = atgconstant.MinUnit
	*funcName = "ComplexDeferFunction"
	*debugMode = true
	*usage = atgconstant.PluginMode
	*UseMockType = 3
	patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
	defer patch.Reset()
	Plugin()
}

func TestGoRoutineCallGraph(t *testing.T) {
	*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
	*minUnit = atgconstant.MinUnit
	*funcName = "ComplexGoRoutineFunction"
	*debugMode = true
	*usage = atgconstant.PluginMode
	*UseMockType = 3
	patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
	defer patch.Reset()
	Plugin()
}

func TestDeferGoExistCallGraph(t *testing.T) {
	*filePath = path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "/atg/template/atg.go")
	*minUnit = atgconstant.MinUnit
	*funcName = "ComplexDeferExistFunction"
	*debugMode = true
	*usage = atgconstant.PluginMode
	*UseMockType = 3
	patch := gomonkey.ApplyFuncReturn(ioutil.WriteFile, nil)
	defer patch.Reset()
	Plugin()
}
