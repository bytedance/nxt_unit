package reporter

import (
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
)

func TestCoverReporter_Analysis(t *testing.T) {
	c := coverReporter{
		hitLine: map[string]map[string]int{},
	}
	c.Analysis("funcCover(GoodFuncPanic;8;5)-r",
		atgconstant.Options{FilePath: "test123", Uid: "ABCDEFGH", FuncName: "123"}, "tee")
	t.Log(c.totalLine, "--", c.hitLine)
}

func TestAnalyserAddTotal(t *testing.T) {
	c := coverReporter{}
	c.AddTotal(8)
	t.Log(c.Report())
}

func TestCoverReporter_BestCoverage(t *testing.T) {
	c := coverReporter{
		hitLine: map[string]map[string]int{},
	}
	c.Analysis("funcCover(GoodFuncPanic;8;5)-r",
		atgconstant.Options{FilePath: "test123", Uid: "ABCDEFGH", FuncName: "123"}, "tee")
	c.Analysis("funcCover(GoodFuncPanic;8;6)-r",
		atgconstant.Options{FilePath: "test123", Uid: "ABCDEFGH", FuncName: "123"}, "bbee")
	c.SetBestIndividual("bbee")
	c.AddTotal(8)
	t.Log(c.totalLine, "--", c.hitCount())
}
