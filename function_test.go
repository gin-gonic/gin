package gin

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func someFunction() {
	// this empty function is used by TestFunctionName()
}

func TestParseFunction(t *testing.T) {
	f := parseFunction(someFunction)
	assert.Equal(t, "github.com/gin-gonic/gin.someFunction", f.Name)
	assert.Equal(t, "function_test.go", filepath.Base(f.File))
	assert.Equal(t, 12, f.Line)
}

func TestFunctionName(t *testing.T) {
	assert.Equal(t, nameOfFunction(someFunction), "github.com/gin-gonic/gin.someFunction")
}
