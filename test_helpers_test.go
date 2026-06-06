// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTestHandlerFlushesStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"201 Created", http.StatusCreated},
		{"204 No Content", http.StatusNoContent},
		{"400 Bad Request", http.StatusBadRequest},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			code := tt.statusCode

			c := RunTestHandler(w, req, func(c *Context) {
				c.Status(code)
			})

			assert.Equal(t, code, w.Code)
			assert.Equal(t, code, c.Writer.Status())
		})
	}
}

func TestRunTestHandlerMultipleHandlers(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/resource", nil)

	callOrder := []string{}

	middleware := func(c *Context) {
		callOrder = append(callOrder, "middleware")
		c.Next()
	}

	handler := func(c *Context) {
		callOrder = append(callOrder, "handler")
		c.Status(http.StatusCreated)
	}

	c := RunTestHandler(w, req, middleware, handler)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusCreated, c.Writer.Status())
	assert.Equal(t, []string{"middleware", "handler"}, callOrder)
}

func TestRunTestHandlerDefaultStatus(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	RunTestHandler(w, req, func(c *Context) {
		// Handler that does not set a status explicitly.
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRunTestHandlerSetsRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/items/42", nil)

	c := RunTestHandler(w, req, func(c *Context) {
		c.Status(http.StatusOK)
	})

	assert.Equal(t, req, c.Request)
}

func TestCreateTestContextBackwardCompatible(t *testing.T) {
	// Verify that CreateTestContext still behaves the same way:
	// status is stored internally but NOT flushed to the ResponseRecorder.
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	c.Status(http.StatusCreated)

	// The internal status should be set.
	assert.Equal(t, http.StatusCreated, c.Writer.Status())

	// But w.Code should still be the default 200 because WriteHeaderNow
	// was never called. This confirms backward compatibility.
	assert.Equal(t, http.StatusOK, w.Code)
}
