package errorparser

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
)

func TestParseValidatorError(t *testing.T) {

	_, ok := parseValidatorError(fmt.Errorf("not match"))
	assert.False(t, ok)

	_, ok = parseValidatorError(validator.ValidationErrors{})
	assert.True(t, ok)

}

func TestParseValidatorValidationErrors(t *testing.T) {

	jsonData := `{
		"text": "",
		"count": 1
	}`

	rbody := bytes.NewReader([]byte(jsonData))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", rbody)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	var obj struct {
		Text  string `json:"text" binding:"required"`
		Count int    `json:"count"`
	}

	err := c.Bind(&obj)
	require.Error(t, err)

	vErr, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	fErrs := []validator.FieldError(vErr)
	require.Equal(t, len(fErrs), 1)

	parseErrs := parseValidatorValidationErrors(vErr)
	require.Equal(t, len(parseErrs), 1)

	assert.Equal(t, parseErrs[0].ParamName, fErrs[0].Field())
	assert.Equal(t, parseErrs[0].ErrorType, ParseErrorTypeValidation)
	assert.Equal(t, parseErrs[0].InitialError, fErrs[0])

}
