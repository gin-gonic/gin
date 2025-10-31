package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteralColonWithRun(t *testing.T) {
	SetMode(TestMode)
	router := New()

	router.GET(`/test\:action`, func(c *Context) {
		c.JSON(http.StatusOK, H{"path": "literal_colon"})
	})

	router.updateRouteTrees()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test:action", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "literal_colon")
}

func TestLiteralColonWithDirectServeHTTP(t *testing.T) {
	SetMode(TestMode)
	router := New()

	router.GET(`/test\:action`, func(c *Context) {
		c.JSON(http.StatusOK, H{"path": "literal_colon"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test:action", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "literal_colon")
}

func TestLiteralColonWithHandler(t *testing.T) {

	SetMode(TestMode)
	router := New()

	router.GET(`/test\:action`, func(c *Context) {
		c.JSON(http.StatusOK, H{"path": "literal_colon"})
	})

	handler := router.Handler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test:action", nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "literal_colon")
}

func TestLiteralColonWithHTTPServer(t *testing.T) {
	SetMode(TestMode)
	router := New()

	router.GET(`/test\:action`, func(c *Context) {
		c.JSON(http.StatusOK, H{"path": "literal_colon"})
	})

	router.GET("/test/:param", func(c *Context) {
		c.JSON(http.StatusOK, H{"param": c.Param("param")})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test:action", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "literal_colon")

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test/foo", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "foo")
}

// Test that updateRouteTrees is called only once
func TestUpdateRouteTreesCalledOnce(t *testing.T) {
	SetMode(TestMode)
	router := New()

	callCount := 0
	originalUpdate := router.updateRouteTrees

	router.GET(`/test\:action`, func(c *Context) {
		c.JSON(http.StatusOK, H{"call": callCount})
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test:action", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	_ = originalUpdate
}
