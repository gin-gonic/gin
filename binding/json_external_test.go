package binding_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJSONBindingEmptyBodyReturnsHelpfulError(t *testing.T) {
	type Req struct {
		Name string `json:"name" binding:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	var r Req
	err = c.ShouldBindJSON(&r)

	assert.Error(t, err)

	// Current behavior returns plain "EOF", which is not helpful.
	assert.NotEqual(t, "EOF", err.Error(), "error message should not be plain EOF")
	assert.Contains(t, err.Error(), "empty request body")
	assert.ErrorIs(t, err, io.EOF)

}
