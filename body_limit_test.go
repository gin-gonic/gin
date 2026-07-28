// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/codec/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// expectedTooLargeStatus mirrors the go-json caveat documented in
// TestContextBindRequestTooLarge: go-json does not propagate
// *http.MaxBytesError, so the response falls back to a generic 400.
func expectedTooLargeStatus() int {
	if json.Package == "github.com/goccy/go-json" {
		return http.StatusBadRequest
	}
	return http.StatusRequestEntityTooLarge
}

// --- unit-level: limitRequestBody branch coverage ---

func TestLimitRequestBodyNilEngine(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("hello"))
	c := &Context{Request: req}
	body := c.Request.Body

	assert.NotPanics(t, func() { c.limitRequestBody() })
	assert.Equal(t, body, c.Request.Body, "body must be untouched when c.engine is nil")
}

func TestLimitRequestBodyDisabled(t *testing.T) {
	for _, limit := range []int64{0, -1, -100} {
		t.Run("", func(t *testing.T) {
			c, _ := CreateTestContext(httptest.NewRecorder())
			c.engine.MaxRequestBodyBytes = limit
			c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader("hello"))
			body := c.Request.Body

			c.limitRequestBody()

			assert.Equal(t, body, c.Request.Body, "body must be untouched when the limit is disabled")
		})
	}
}

func TestLimitRequestBodyNilRequest(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request = nil

	assert.NotPanics(t, func() { c.limitRequestBody() })
}

func TestLimitRequestBodyNilBody(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	require.Nil(t, c.Request.Body)

	assert.NotPanics(t, func() { c.limitRequestBody() })
	assert.Nil(t, c.Request.Body)
}

func TestLimitRequestBodyWraps(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.engine.MaxRequestBodyBytes = 4
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader("hello world"))

	c.limitRequestBody()

	_, err := io.ReadAll(c.Request.Body)
	var maxBytesErr *http.MaxBytesError
	require.ErrorAs(t, err, &maxBytesErr)
	assert.Equal(t, int64(4), maxBytesErr.Limit)
}

// --- integration: ShouldBindWith / Bind / MustBindWith ---

func TestBindWithinBodyLimit(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.engine.MaxRequestBodyBytes = 1024
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar"}`))

	var obj struct {
		Foo string `json:"foo"`
	}
	require.NoError(t, c.BindJSON(&obj))
	assert.Equal(t, "bar", obj.Foo)
	assert.False(t, c.IsAborted())
}

func TestBindExceedsBodyLimit(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.engine.MaxRequestBodyBytes = 10
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar", "bar":"foo"}`))

	var obj struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	err := c.BindJSON(&obj)
	require.Error(t, err)
	c.Writer.WriteHeaderNow()

	if json.Package != "github.com/goccy/go-json" {
		var maxBytesErr *http.MaxBytesError
		require.ErrorAs(t, err, &maxBytesErr)
	}
	assert.Empty(t, obj.Foo)
	assert.Empty(t, obj.Bar)
	assert.Equal(t, expectedTooLargeStatus(), w.Code)
	assert.True(t, c.IsAborted())
}

func TestBindBodyLimitBoundary(t *testing.T) {
	tests := []struct {
		name    string
		bodyLen int
		limit   int64
		wantErr bool
	}{
		{"at limit succeeds", 10, 10, false},
		{"over limit by one fails", 11, 10, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := CreateTestContext(httptest.NewRecorder())
			c.engine.MaxRequestBodyBytes = tt.limit
			body := strings.Repeat("a", tt.bodyLen)
			c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

			var s string
			err := c.ShouldBindWith(&s, binding.Plain)

			if tt.wantErr {
				require.Error(t, err)
				var maxBytesErr *http.MaxBytesError
				require.ErrorAs(t, err, &maxBytesErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, body, s)
			}
		})
	}
}

func TestBindBodyLimitDisabled(t *testing.T) {
	for _, limit := range []int64{0, -1} {
		t.Run("", func(t *testing.T) {
			c, _ := CreateTestContext(httptest.NewRecorder())
			c.engine.MaxRequestBodyBytes = limit
			body := strings.Repeat("a", 5000) // larger than any limit used elsewhere in this file
			c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

			var s string
			require.NoError(t, c.ShouldBindWith(&s, binding.Plain))
			assert.Equal(t, body, s)
		})
	}
}

func TestBindBodyLimitAcrossFormats(t *testing.T) {
	type obj struct {
		Foo string `json:"foo" xml:"foo" bson:"foo"`
	}

	oversizedJSON := `{"foo":"` + strings.Repeat("x", 100) + `"}`
	oversizedXML := `<root><foo>` + strings.Repeat("x", 100) + `</foo></root>`
	bsonBody, err := bson.Marshal(&obj{Foo: strings.Repeat("x", 100)})
	require.NoError(t, err)
	oversizedPlain := strings.Repeat("x", 100)

	tests := []struct {
		name   string
		b      binding.Binding
		body   string
		target func() any
	}{
		{"json", binding.JSON, oversizedJSON, func() any { return &obj{} }},
		{"xml", binding.XML, oversizedXML, func() any { return &obj{} }},
		{"bson", binding.BSON, string(bsonBody), func() any { return &obj{} }},
		{"plain", binding.Plain, oversizedPlain, func() any { var s string; return &s }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := CreateTestContext(httptest.NewRecorder())
			c.engine.MaxRequestBodyBytes = 10 // well under every body above
			c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))

			err := c.ShouldBindWith(tt.target(), tt.b)
			require.Error(t, err)
			var maxBytesErr *http.MaxBytesError
			require.ErrorAs(t, err, &maxBytesErr)
		})
	}
}

func TestBindNonBodyBindingsUnaffectedByLimit(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.engine.MaxRequestBodyBytes = 1 // tiny enough to reject almost any body
	c.Request, _ = http.NewRequest(http.MethodPost, "/?foo=bar", nil)
	c.Request.Header.Add("rate", "8000")

	var q struct {
		Foo string `form:"foo"`
	}
	require.NoError(t, c.BindQuery(&q))
	assert.Equal(t, "bar", q.Foo)

	var h struct {
		Rate int `header:"rate"`
	}
	require.NoError(t, c.ShouldBindHeader(&h))
	assert.Equal(t, 8000, h.Rate)
}

// --- integration: ShouldBindBodyWith ---

func TestShouldBindBodyWithExceedsBodyLimit(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.engine.MaxRequestBodyBytes = 5
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar"}`))

	var obj struct {
		Foo string `json:"foo"`
	}
	err := c.ShouldBindBodyWith(&obj, binding.JSON)
	require.Error(t, err)
	var maxBytesErr *http.MaxBytesError
	require.ErrorAs(t, err, &maxBytesErr)
}

func TestShouldBindBodyWithWithinBodyLimit(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.engine.MaxRequestBodyBytes = 1024
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar"}`))

	var obj struct {
		Foo string `json:"foo"`
	}
	require.NoError(t, c.ShouldBindBodyWith(&obj, binding.JSON))
	assert.Equal(t, "bar", obj.Foo)

	cached, ok := c.Get(BodyBytesKey)
	require.True(t, ok)
	assert.JSONEq(t, `{"foo":"bar"}`, string(cached.([]byte)))
}

func TestShouldBindBodyWithCachedReadBypassesReWrap(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.engine.MaxRequestBodyBytes = 1024
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar"}`))

	var obj1 struct {
		Foo string `json:"foo"`
	}
	require.NoError(t, c.ShouldBindBodyWith(&obj1, binding.JSON))
	assert.Equal(t, "bar", obj1.Foo)

	// Second call must succeed from the BodyBytesKey cache even though
	// c.Request.Body has already been drained and re-wrapped by the first
	// call's limitRequestBody().
	var obj2 struct {
		Foo string `json:"foo"`
	}
	require.NoError(t, c.ShouldBindBodyWith(&obj2, binding.JSON))
	assert.Equal(t, "bar", obj2.Foo)
}
