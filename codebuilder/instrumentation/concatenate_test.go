package instrumentation

import (
	"path"
	"testing"

	"github.com/bytedance/nxt_unit/atgconstant"
	"github.com/stretchr/testify/assert"
)

func TestGetConvertedImportsFromSrc(t *testing.T) {
	path := path.Join(atgconstant.GOPATHSRC, atgconstant.ProjectPath, "atghelper/atgutil_test.go")
	convertedImports, err := GetImportsInfosFromFile(path)
	assert.Nil(t, err)
	assert.Equal(t, len(convertedImports), 3)
}
