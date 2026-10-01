// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin/codec/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblemMarshalJSON(t *testing.T) {
	p := Problem{
		Type:     "https://example.com/probs/out-of-credit",
		Title:    "You do not have enough credit.",
		Status:   403,
		Detail:   "Your current balance is 30, but that costs 50.",
		Instance: "/account/12345/msgs/abc",
	}

	jsonBytes, err := json.API.Marshal(p)

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"type": "https://example.com/probs/out-of-credit",
		"title": "You do not have enough credit.",
		"status": 403,
		"detail": "Your current balance is 30, but that costs 50.",
		"instance": "/account/12345/msgs/abc"
	}`, string(jsonBytes))
}

func TestProblemMarshalJSONOmitsEmptyMembers(t *testing.T) {
	jsonBytes, err := json.API.Marshal(Problem{Status: 404})

	require.NoError(t, err)
	assert.JSONEq(t, `{"status":404}`, string(jsonBytes))
}

func TestProblemMarshalJSONExtensions(t *testing.T) {
	p := Problem{
		Status:   403,
		Detail:   "Your current balance is 30, but that costs 50.",
		Instance: "/account/12345/msgs/abc",
		Extensions: map[string]any{
			"balance": 30,
			"status":  "extension members must not override standard members",
		},
	}

	jsonBytes, err := json.API.Marshal(p)

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"status": 403,
		"detail": "Your current balance is 30, but that costs 50.",
		"instance": "/account/12345/msgs/abc",
		"balance": 30
	}`, string(jsonBytes))
}

func TestProblemDetailsMiddleware(t *testing.T) {
	router := New()
	router.Use(ProblemDetails())
	router.GET("/error", func(c *Context) {
		c.Error(errors.New("boom")) //nolint:errcheck
	})

	w := PerformRequest(router, http.MethodGet, "/error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/problem+json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"title":"Internal Server Error","status":500,"detail":"boom"}`, w.Body.String())
}

func TestProblemDetailsMiddlewareKeepsErrorStatus(t *testing.T) {
	router := New()
	router.Use(ProblemDetails())
	router.GET("/conflict", func(c *Context) {
		c.Status(http.StatusConflict)
		c.Error(errors.New("already exists")) //nolint:errcheck
	})

	w := PerformRequest(router, http.MethodGet, "/conflict")

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "application/problem+json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"title":"Conflict","status":409,"detail":"already exists"}`, w.Body.String())
}

func TestProblemDetailsMiddlewareNoErrors(t *testing.T) {
	router := New()
	router.Use(ProblemDetails())
	router.GET("/ok", func(c *Context) {
		c.JSON(http.StatusOK, H{"foo": "bar"})
	})

	w := PerformRequest(router, http.MethodGet, "/ok")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"foo":"bar"}`, w.Body.String())
}

func TestProblemDetailsMiddlewareResponseAlreadyWritten(t *testing.T) {
	router := New()
	router.Use(ProblemDetails())
	router.GET("/written", func(c *Context) {
		c.JSON(http.StatusBadGateway, H{"error": "custom body"})
		c.Error(errors.New("boom")) //nolint:errcheck
	})

	w := PerformRequest(router, http.MethodGet, "/written")

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"error":"custom body"}`, w.Body.String())
}
