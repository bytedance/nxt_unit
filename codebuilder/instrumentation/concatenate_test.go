package instrumentation

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConvertedImportsFromSrc(t *testing.T) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		assert.Equal(t, ok, true)
		return
	}
	filePath = path.Join(path.Dir(filePath), "../../atghelper/atgutil_test.go")
	convertedImports, err := GetImportsInfosFromFile(filePath)
	assert.Nil(t, err)
	assert.Equal(t, len(convertedImports), 3)
}
