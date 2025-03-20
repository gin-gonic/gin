package gin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_Use(t *testing.T) {
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New(
		Use(middleware1, middleware2),
	)

	assert.Equal(t, 2, len(router.Handlers))
	compareFunc(t, middleware1, router.Handlers[0])
	compareFunc(t, middleware2, router.Handlers[1])
}

func TestOption_HttpMethod(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		optionFunc OptionFunc
		want       int
	}{
		{
			method: http.MethodGet,
			path:   "/get",
			optionFunc: GET("/get", func(c *Context) {
				assert.Equal(t, http.MethodGet, c.Request.Method)
				assert.Equal(t, "/get", c.Request.URL.Path)
			}),
		},
		{
			method: http.MethodPut,
			path:   "/put",
			optionFunc: PUT("/put", func(c *Context) {
				assert.Equal(t, http.MethodPut, c.Request.Method)
				assert.Equal(t, "/put", c.Request.URL.Path)
			}),
		},
		{
			method: http.MethodPost,
			path:   "/post",
			optionFunc: POST("/post", func(c *Context) {
				assert.Equal(t, http.MethodPost, c.Request.Method)
				assert.Equal(t, "/post", c.Request.URL.Path)
			}),
		},
		{
			method: http.MethodDelete,
			path:   "/delete",
			optionFunc: DELETE("/delete", func(c *Context) {
				assert.Equal(t, http.MethodDelete, c.Request.Method)
				assert.Equal(t, "/delete", c.Request.URL.Path)
			}),
		},
		{
			method: http.MethodPatch,
			path:   "/patch",
			optionFunc: PATCH("/patch", func(c *Context) {
				assert.Equal(t, http.MethodPatch, c.Request.Method)
				assert.Equal(t, "/patch", c.Request.URL.Path)
			}),
		},
		{
			method: http.MethodOptions,
			path:   "/options",
			optionFunc: OPTIONS("/options", func(c *Context) {
				assert.Equal(t, http.MethodOptions, c.Request.Method)
				assert.Equal(t, "/options", c.Request.URL.Path)
			}),
		},
		{
			method: http.MethodHead,
			path:   "/head",
			optionFunc: HEAD("/head", func(c *Context) {
				assert.Equal(t, http.MethodHead, c.Request.Method)
				assert.Equal(t, "/head", c.Request.URL.Path)
			}),
		},
		{
			method: "GET",
			path:   "/any",
			optionFunc: Any("/any", func(c *Context) {
				assert.Equal(t, http.MethodGet, c.Request.Method)
				assert.Equal(t, "/any", c.Request.URL.Path)
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			router := New(tt.optionFunc)
			w := PerformRequest(router, tt.method, tt.path)
			assert.Equal(t, 200, w.Code)
		})
	}
}

func TestOption_Any(t *testing.T) {
	method := make(chan string, 1)
	router := New(
		Any("/any", func(c *Context) {
			method <- c.Request.Method
			assert.Equal(t, "/any", c.Request.URL.Path)
		}),
	)

	tests := []struct {
		method string
	}{
		{http.MethodGet},
		{http.MethodPost},
		{http.MethodPut},
		{http.MethodPatch},
		{http.MethodDelete},
		{http.MethodHead},
		{http.MethodOptions},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			w := PerformRequest(router, tt.method, "/any")
			assert.Equal(t, 200, w.Code)
			assert.Equal(t, tt.method, <-method)
		})
	}
}

func TestOption_Group(t *testing.T) {
	router := New(
		Group("/v1", func(group *RouterGroup) {
			group.GET("/test", func(c *Context) {
				assert.Equal(t, http.MethodGet, c.Request.Method)
				assert.Equal(t, "/v1/test", c.Request.URL.Path)
			})
		}),
	)

	w := PerformRequest(router, http.MethodGet, "/v1/test")
	assert.Equal(t, 200, w.Code)
}

func TestOption_Route(t *testing.T) {
	router := New(
		Route(http.MethodGet, "/test", func(c *Context) {
			assert.Equal(t, http.MethodGet, c.Request.Method)
			assert.Equal(t, "/test", c.Request.URL.Path)
		}),
	)

	w := PerformRequest(router, http.MethodGet, "/test")
	assert.Equal(t, 200, w.Code)
}

func TestOption_StaticFS(t *testing.T) {
	router := New(
		StaticFS("/", http.Dir("./")),
	)

	w := PerformRequest(router, http.MethodGet, "/gin.go")
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "package gin")
}

func TestOption_StaticFile(t *testing.T) {
	router := New(
		StaticFile("/gin.go", "gin.go"),
	)

	w := PerformRequest(router, http.MethodGet, "/gin.go")
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "package gin")
}

func TestOption_Static(t *testing.T) {
	router := New(
		Static("/static", "./"),
	)

	w := PerformRequest(router, http.MethodGet, "/static/gin.go")
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "package gin")
}

func TestOption_NoRoute(t *testing.T) {
	router := New(
		NoRoute(func(c *Context) {
			c.String(http.StatusNotFound, "no route")
		}),
	)

	assert.Equal(t, 1, len(router.noRoute))

	w := PerformRequest(router, http.MethodGet, "/no-route")
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "no route", w.Body.String())
}
