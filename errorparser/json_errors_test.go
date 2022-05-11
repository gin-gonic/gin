package errorparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
)

func TestParseJsonDecodeError(t *testing.T) {

	_, ok := parseJsonDecodeError(fmt.Errorf("not match"))
	assert.False(t, ok)

	_, ok = parseJsonDecodeError(&json.UnmarshalTypeError{})
	assert.True(t, ok)

	_, ok = parseJsonDecodeError(&json.SyntaxError{})
	assert.True(t, ok)

}

func TestParseJsonUnmarshalTypeError(t *testing.T) {

	jsonData := `{
		"text": "text",
		"count": "1"
	}`

	rbody := bytes.NewReader([]byte(jsonData))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", rbody)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	var obj struct {
		Text  string `json:"text"`
		Count int    `json:"count"`
	}

	err := c.Bind(&obj)
	require.Error(t, err)

	typeErr, ok := err.(*json.UnmarshalTypeError)
	require.True(t, ok)

	parseErrs := parseJsonUnmarshalTypeError(typeErr)
	require.Equal(t, len(parseErrs), 1)

	assert.Equal(t, parseErrs[0].ParamName, "count")
	assert.Equal(t, parseErrs[0].ErrorType, ParseErrorTypeMismatch)
	assert.Equal(t, parseErrs[0].InitialError, err)

}

func TestParseJsonSyntaxError(t *testing.T) {

	jsonData := `{
		"text": "text"
		"count": 1
	}`

	rbody := bytes.NewReader([]byte(jsonData))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", rbody)
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)

	var obj struct {
		Text  string `json:"text"`
		Count int    `json:"count"`
	}

	err := c.Bind(&obj)
	require.Error(t, err)

	typeErr, ok := err.(*json.SyntaxError)
	require.True(t, ok)

	parseErrs := parseJsonSyntaxError(typeErr)
	require.Equal(t, len(parseErrs), 1)

	assert.Equal(t, parseErrs[0].ParamName, "")
	assert.Equal(t, parseErrs[0].ErrorType, ParseErrorTypeBadInput)
	assert.Equal(t, parseErrs[0].InitialError, err)

}
