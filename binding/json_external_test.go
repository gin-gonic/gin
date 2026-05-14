package binding_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONBindingEmptyBodyReturnsHelpfulError(t *testing.T) {
	type Req struct {
		Name string `json:"name" binding:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	var r Req
	err = c.ShouldBindJSON(&r)

	require.Error(t, err)

	// Error message should be more descriptive than plain EOF,
	// while still preserving io.EOF via wrapping.
	assert.NotEqual(t, "EOF", err.Error())
	assert.Contains(t, err.Error(), "empty request body")
	assert.ErrorIs(t, err, io.EOF)
}
