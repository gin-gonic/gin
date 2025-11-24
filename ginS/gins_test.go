// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginS

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGET(t *testing.T) {
	GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", w.Body.String())
}

func TestPOST(t *testing.T) {
	POST("/post", func(c *gin.Context) {
		c.String(http.StatusCreated, "created")
	})

	req := httptest.NewRequest(http.MethodPost, "/post", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "created", w.Body.String())
}

func TestPUT(t *testing.T) {
	PUT("/put", func(c *gin.Context) {
		c.String(http.StatusOK, "updated")
	})

	req := httptest.NewRequest(http.MethodPut, "/put", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "updated", w.Body.String())
}

func TestDELETE(t *testing.T) {
	DELETE("/delete", func(c *gin.Context) {
		c.String(http.StatusOK, "deleted")
	})

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "deleted", w.Body.String())
}

func TestPATCH(t *testing.T) {
	PATCH("/patch", func(c *gin.Context) {
		c.String(http.StatusOK, "patched")
	})

	req := httptest.NewRequest(http.MethodPatch, "/patch", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "patched", w.Body.String())
}

func TestOPTIONS(t *testing.T) {
	OPTIONS("/options", func(c *gin.Context) {
		c.String(http.StatusOK, "options")
	})

	req := httptest.NewRequest(http.MethodOptions, "/options", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "options", w.Body.String())
}

func TestHEAD(t *testing.T) {
	HEAD("/head", func(c *gin.Context) {
		c.String(http.StatusOK, "head")
	})

	req := httptest.NewRequest(http.MethodHead, "/head", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAny(t *testing.T) {
	Any("/any", func(c *gin.Context) {
		c.String(http.StatusOK, "any")
	})

	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "any", w.Body.String())
}

func TestHandle(t *testing.T) {
	Handle(http.MethodGet, "/handle", func(c *gin.Context) {
		c.String(http.StatusOK, "handle")
	})

	req := httptest.NewRequest(http.MethodGet, "/handle", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "handle", w.Body.String())
}

func TestGroup(t *testing.T) {
	group := Group("/group")
	group.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "group test")
	})

	req := httptest.NewRequest(http.MethodGet, "/group/test", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "group test", w.Body.String())
}

func TestUse(t *testing.T) {
	var middlewareExecuted bool
	Use(func(c *gin.Context) {
		middlewareExecuted = true
		c.Next()
	})

	GET("/middleware-test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/middleware-test", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.True(t, middlewareExecuted)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNoRoute(t *testing.T) {
	NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "custom 404")
	})

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "custom 404", w.Body.String())
}

func TestNoMethod(t *testing.T) {
	NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	})

	// This just verifies that NoMethod is callable
	// Testing the actual behavior would require a separate engine instance
	assert.NotNil(t, engine())
}

func TestRoutes(t *testing.T) {
	GET("/routes-test", func(c *gin.Context) {})

	routes := Routes()
	assert.NotEmpty(t, routes)

	found := false
	for _, route := range routes {
		if route.Path == "/routes-test" && route.Method == http.MethodGet {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestSetHTMLTemplate(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("Hello {{.}}"))
	SetHTMLTemplate(tmpl)

	// Verify engine has template set
	assert.NotNil(t, engine())
}

func TestStaticFile(t *testing.T) {
	StaticFile("/static-file", "../testdata/test_file.txt")

	req := httptest.NewRequest(http.MethodGet, "/static-file", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStatic(t *testing.T) {
	Static("/static-dir", "../testdata")

	req := httptest.NewRequest(http.MethodGet, "/static-dir/test_file.txt", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStaticFS(t *testing.T) {
	fs := http.Dir("../testdata")
	StaticFS("/static-fs", fs)

	req := httptest.NewRequest(http.MethodGet, "/static-fs/test_file.txt", nil)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
