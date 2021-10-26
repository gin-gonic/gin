// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func testRouteOK(method string, t *testing.T) {
	passed := false
	passedAny := false
	r := New()
	r.Any("/test2", func(c *Context) {
		passedAny = true
	})
	r.Handle(method, "/test", func(c *Context) {
		passed = true
	})

	w := performRequest(r, method, "/test")
	assert.True(t, passed)
	assert.Equal(t, http.StatusOK, w.Code)

	performRequest(r, method, "/test2")
	assert.True(t, passedAny)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK(method string, t *testing.T) {
	passed := false
	router := New()
	router.Handle(method, "/test_2", func(c *Context) {
		passed = true
	})

	w := performRequest(router, method, "/test")

	assert.False(t, passed)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK2(method string, t *testing.T) {
	passed := false
	router := New()
	router.HandleMethodNotAllowed = true
	var methodRoute string
	if method == http.MethodPost {
		methodRoute = http.MethodGet
	} else {
		methodRoute = http.MethodPost
	}
	router.Handle(methodRoute, "/test", func(c *Context) {
		passed = true
	})

	w := performRequest(router, method, "/test")

	assert.False(t, passed)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestRouterMethod(t *testing.T) {
	router := New()
	router.PUT("/hey2", func(c *Context) {
		c.String(http.StatusOK, "sup2")
	})

	router.PUT("/hey", func(c *Context) {
		c.String(http.StatusOK, "called")
	})

	router.PUT("/hey3", func(c *Context) {
		c.String(http.StatusOK, "sup3")
	})

	w := performRequest(router, http.MethodPut, "/hey")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "called", w.Body.String())
}

func TestRouterGroupRouteOK(t *testing.T) {
	testRouteOK(http.MethodGet, t)
	testRouteOK(http.MethodPost, t)
	testRouteOK(http.MethodPut, t)
	testRouteOK(http.MethodPatch, t)
	testRouteOK(http.MethodHead, t)
	testRouteOK(http.MethodOptions, t)
	testRouteOK(http.MethodDelete, t)
	testRouteOK(http.MethodConnect, t)
	testRouteOK(http.MethodTrace, t)
}

func TestRouteNotOK(t *testing.T) {
	testRouteNotOK(http.MethodGet, t)
	testRouteNotOK(http.MethodPost, t)
	testRouteNotOK(http.MethodPut, t)
	testRouteNotOK(http.MethodPatch, t)
	testRouteNotOK(http.MethodHead, t)
	testRouteNotOK(http.MethodOptions, t)
	testRouteNotOK(http.MethodDelete, t)
	testRouteNotOK(http.MethodConnect, t)
	testRouteNotOK(http.MethodTrace, t)
}

func TestRouteNotOK2(t *testing.T) {
	testRouteNotOK2(http.MethodGet, t)
	testRouteNotOK2(http.MethodPost, t)
	testRouteNotOK2(http.MethodPut, t)
	testRouteNotOK2(http.MethodPatch, t)
	testRouteNotOK2(http.MethodHead, t)
	testRouteNotOK2(http.MethodOptions, t)
	testRouteNotOK2(http.MethodDelete, t)
	testRouteNotOK2(http.MethodConnect, t)
	testRouteNotOK2(http.MethodTrace, t)
}

func TestRouteRedirectTrailingSlash(t *testing.T) {
	router := New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = true
	router.GET("/path", func(c *Context) {})
	router.GET("/path2/", func(c *Context) {})
	router.POST("/path3", func(c *Context) {})
	router.PUT("/path4/", func(c *Context) {})

	w := performRequest(router, http.MethodGet, "/path/")
	assert.Equal(t, "/path", w.Header().Get("Location"))
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(router, http.MethodGet, "/path2")
	assert.Equal(t, "/path2/", w.Header().Get("Location"))
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(router, http.MethodPost, "/path3/")
	assert.Equal(t, "/path3", w.Header().Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	w = performRequest(router, http.MethodPut, "/path4")
	assert.Equal(t, "/path4/", w.Header().Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodGet, "/path2/")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodPost, "/path3")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodPut, "/path4/")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodGet, "/path2", header{Key: "X-Forwarded-Prefix", Value: "/api"})
	assert.Equal(t, "/api/path2/", w.Header().Get("Location"))
	assert.Equal(t, 301, w.Code)

	w = performRequest(router, http.MethodGet, "/path2/", header{Key: "X-Forwarded-Prefix", Value: "/api/"})
	assert.Equal(t, 200, w.Code)

	router.RedirectTrailingSlash = false

	w = performRequest(router, http.MethodGet, "/path/")
	assert.Equal(t, http.StatusNotFound, w.Code)
	w = performRequest(router, http.MethodGet, "/path2")
	assert.Equal(t, http.StatusNotFound, w.Code)
	w = performRequest(router, http.MethodPost, "/path3/")
	assert.Equal(t, http.StatusNotFound, w.Code)
	w = performRequest(router, http.MethodPut, "/path4")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouteRedirectFixedPath(t *testing.T) {
	router := New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = false

	router.GET("/path", func(c *Context) {})
	router.GET("/Path2", func(c *Context) {})
	router.POST("/PATH3", func(c *Context) {})
	router.POST("/Path4/", func(c *Context) {})

	w := performRequest(router, http.MethodGet, "/PATH")
	assert.Equal(t, "/path", w.Header().Get("Location"))
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(router, http.MethodGet, "/path2")
	assert.Equal(t, "/Path2", w.Header().Get("Location"))
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(router, http.MethodPost, "/path3")
	assert.Equal(t, "/PATH3", w.Header().Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	w = performRequest(router, http.MethodPost, "/path4")
	assert.Equal(t, "/Path4/", w.Header().Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

// TestContextParamsGet tests that a parameter can be parsed from the URL.
func TestRouteParamsByName(t *testing.T) {
	name := ""
	lastName := ""
	wild := ""
	router := New()
	router.GET("/test/:name/:last_name/*wild", func(c *Context) {
		name = c.Params.ByName("name")
		lastName = c.Params.ByName("last_name")
		var ok bool
		wild, ok = c.Params.Get("wild")

		assert.True(t, ok)
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, lastName, c.Param("last_name"))

		assert.Empty(t, c.Param("wtf"))
		assert.Empty(t, c.Params.ByName("wtf"))

		wtf, ok := c.Params.Get("wtf")
		assert.Empty(t, wtf)
		assert.False(t, ok)
	})

	w := performRequest(router, http.MethodGet, "/test/john/smith/is/super/great")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "john", name)
	assert.Equal(t, "smith", lastName)
	assert.Equal(t, "/is/super/great", wild)
}

// TestContextParamsGet tests that a parameter can be parsed from the URL even with extra slashes.
func TestRouteParamsByNameWithExtraSlash(t *testing.T) {
	name := ""
	lastName := ""
	wild := ""
	router := New()
	router.RemoveExtraSlash = true
	router.GET("/test/:name/:last_name/*wild", func(c *Context) {
		name = c.Params.ByName("name")
		lastName = c.Params.ByName("last_name")
		var ok bool
		wild, ok = c.Params.Get("wild")

		assert.True(t, ok)
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, lastName, c.Param("last_name"))

		assert.Empty(t, c.Param("wtf"))
		assert.Empty(t, c.Params.ByName("wtf"))

		wtf, ok := c.Params.Get("wtf")
		assert.Empty(t, wtf)
		assert.False(t, ok)
	})

	w := performRequest(router, http.MethodGet, "//test//john//smith//is//super//great")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "john", name)
	assert.Equal(t, "smith", lastName)
	assert.Equal(t, "/is/super/great", wild)
}

// TestHandleStaticFile - ensure the static file handles properly
func TestRouteStaticFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := ioutil.TempFile(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString("Gin Web Framework")
	assert.NoError(t, err)
	f.Close()

	dir, filename := filepath.Split(f.Name())

	// SETUP gin
	router := New()
	router.Static("/using_static", dir)
	router.StaticFile("/result", f.Name())

	w := performRequest(router, http.MethodGet, "/using_static/"+filename)
	w2 := performRequest(router, http.MethodGet, "/result")

	assert.Equal(t, w, w2)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Gin Web Framework", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	w3 := performRequest(router, http.MethodHead, "/using_static/"+filename)
	w4 := performRequest(router, http.MethodHead, "/result")

	assert.Equal(t, w3, w4)
	assert.Equal(t, http.StatusOK, w3.Code)
}

// TestHandleStaticDir - ensure the root/sub dir handles properly
func TestRouteStaticListingDir(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir("./", true))

	w := performRequest(router, http.MethodGet, "/")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "gin.go")
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

// TestHandleHeadToDir - ensure the root/sub dir handles properly
func TestRouteStaticNoListing(t *testing.T) {
	router := New()
	router.Static("/", "./")

	w := performRequest(router, http.MethodGet, "/")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NotContains(t, w.Body.String(), "gin.go")
}

func TestRouterMiddlewareAndStatic(t *testing.T) {
	router := New()
	static := router.Group("/", func(c *Context) {
		c.Writer.Header().Add("Last-Modified", "Mon, 02 Jan 2006 15:04:05 MST")
		c.Writer.Header().Add("Expires", "Mon, 02 Jan 2006 15:04:05 MST")
		c.Writer.Header().Add("X-GIN", "Gin Framework")
	})
	static.Static("/", "./")

	w := performRequest(router, http.MethodGet, "/gin.go")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "package gin")
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.NotEqual(t, w.Header().Get("Last-Modified"), "Mon, 02 Jan 2006 15:04:05 MST")
	assert.Equal(t, "Mon, 02 Jan 2006 15:04:05 MST", w.Header().Get("Expires"))
	assert.Equal(t, "Gin Framework", w.Header().Get("x-GIN"))
}

func TestRouteNotAllowedEnabled(t *testing.T) {
	router := New()
	router.HandleMethodNotAllowed = true
	router.POST("/path", func(c *Context) {})
	w := performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	router.NoMethod(func(c *Context) {
		c.String(http.StatusTeapot, "responseText")
	})
	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, "responseText", w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestRouteNotAllowedEnabled2(t *testing.T) {
	router := New()
	router.HandleMethodNotAllowed = true
	// add one methodTree to trees
	router.addRoute(http.MethodPost, "/", HandlersChain{func(_ *Context) {}})
	router.GET("/path2", func(c *Context) {})
	w := performRequest(router, http.MethodPost, "/path2")
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestRouteNotAllowedDisabled(t *testing.T) {
	router := New()
	router.HandleMethodNotAllowed = false
	router.POST("/path", func(c *Context) {})
	w := performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusNotFound, w.Code)

	router.NoMethod(func(c *Context) {
		c.String(http.StatusTeapot, "responseText")
	})
	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterNotFoundWithRemoveExtraSlash(t *testing.T) {
	router := New()
	router.RemoveExtraSlash = true
	router.GET("/path", func(c *Context) {})
	router.GET("/", func(c *Context) {})

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/../path", http.StatusOK, ""},    // CleanPath
		{"/nope", http.StatusNotFound, ""}, // NotFound
	}
	for _, tr := range testRoutes {
		w := performRequest(router, "GET", tr.route)
		assert.Equal(t, tr.code, w.Code)
		if w.Code != http.StatusNotFound {
			assert.Equal(t, tr.location, fmt.Sprint(w.Header().Get("Location")))
		}
	}
}

func TestRouterNotFound(t *testing.T) {
	router := New()
	router.RedirectFixedPath = true
	router.GET("/path", func(c *Context) {})
	router.GET("/dir/", func(c *Context) {})
	router.GET("/", func(c *Context) {})

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/path/", http.StatusMovedPermanently, "/path"},   // TSR -/
		{"/dir", http.StatusMovedPermanently, "/dir/"},     // TSR +/
		{"/PATH", http.StatusMovedPermanently, "/path"},    // Fixed Case
		{"/DIR/", http.StatusMovedPermanently, "/dir/"},    // Fixed Case
		{"/PATH/", http.StatusMovedPermanently, "/path"},   // Fixed Case -/
		{"/DIR", http.StatusMovedPermanently, "/dir/"},     // Fixed Case +/
		{"/../path", http.StatusMovedPermanently, "/path"}, // Without CleanPath
		{"/nope", http.StatusNotFound, ""},                 // NotFound
	}
	for _, tr := range testRoutes {
		w := performRequest(router, http.MethodGet, tr.route)
		assert.Equal(t, tr.code, w.Code)
		if w.Code != http.StatusNotFound {
			assert.Equal(t, tr.location, fmt.Sprint(w.Header().Get("Location")))
		}
	}

	// Test custom not found handler
	var notFound bool
	router.NoRoute(func(c *Context) {
		c.AbortWithStatus(http.StatusNotFound)
		notFound = true
	})
	w := performRequest(router, http.MethodGet, "/nope")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.True(t, notFound)

	// Test other method than GET (want 307 instead of 301)
	router.PATCH("/path", func(c *Context) {})
	w = performRequest(router, http.MethodPatch, "/path/")
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "map[Location:[/path]]", fmt.Sprint(w.Header()))

	// Test special case where no node for the prefix "/" exists
	router = New()
	router.GET("/a", func(c *Context) {})
	w = performRequest(router, http.MethodGet, "/")
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Reproduction test for the bug of issue #2843
	router = New()
	router.NoRoute(func(c *Context) {
		if c.Request.RequestURI == "/login" {
			c.String(200, "login")
		}
	})
	router.GET("/logout", func(c *Context) {
		c.String(200, "logout")
	})
	w = performRequest(router, http.MethodGet, "/login")
	assert.Equal(t, "login", w.Body.String())
	w = performRequest(router, http.MethodGet, "/logout")
	assert.Equal(t, "logout", w.Body.String())
}

func TestRouterStaticFSNotFound(t *testing.T) {
	router := New()
	router.StaticFS("/", http.FileSystem(http.Dir("/thisreallydoesntexist/")))
	router.NoRoute(func(c *Context) {
		c.String(404, "non existent")
	})

	w := performRequest(router, http.MethodGet, "/nonexistent")
	assert.Equal(t, "non existent", w.Body.String())

	w = performRequest(router, http.MethodHead, "/nonexistent")
	assert.Equal(t, "non existent", w.Body.String())
}

func TestRouterStaticFSFileNotFound(t *testing.T) {
	router := New()

	router.StaticFS("/", http.FileSystem(http.Dir(".")))

	assert.NotPanics(t, func() {
		performRequest(router, http.MethodGet, "/nonexistent")
	})
}

// Reproduction test for the bug of issue #1805
func TestMiddlewareCalledOnceByRouterStaticFSNotFound(t *testing.T) {
	router := New()

	// Middleware must be called just only once by per request.
	middlewareCalledNum := 0
	router.Use(func(c *Context) {
		middlewareCalledNum++
	})

	router.StaticFS("/", http.FileSystem(http.Dir("/thisreallydoesntexist/")))

	// First access
	performRequest(router, http.MethodGet, "/nonexistent")
	assert.Equal(t, 1, middlewareCalledNum)

	// Second access
	performRequest(router, http.MethodHead, "/nonexistent")
	assert.Equal(t, 2, middlewareCalledNum)
}

func TestRouteRawPath(t *testing.T) {
	route := New()
	route.UseRawPath = true

	route.POST("/project/:name/build/:num", func(c *Context) {
		name := c.Params.ByName("name")
		num := c.Params.ByName("num")

		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, num, c.Param("num"))

		assert.Equal(t, "Some/Other/Project", name)
		assert.Equal(t, "222", num)
	})

	w := performRequest(route, http.MethodPost, "/project/Some%2FOther%2FProject/build/222")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouteRawPathNoUnescape(t *testing.T) {
	route := New()
	route.UseRawPath = true
	route.UnescapePathValues = false

	route.POST("/project/:name/build/:num", func(c *Context) {
		name := c.Params.ByName("name")
		num := c.Params.ByName("num")

		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, num, c.Param("num"))

		assert.Equal(t, "Some%2FOther%2FProject", name)
		assert.Equal(t, "333", num)
	})

	w := performRequest(route, http.MethodPost, "/project/Some%2FOther%2FProject/build/333")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouteServeErrorWithWriteHeader(t *testing.T) {
	route := New()
	route.Use(func(c *Context) {
		c.Status(421)
		c.Next()
	})

	w := performRequest(route, http.MethodGet, "/NotFound")
	assert.Equal(t, 421, w.Code)
	assert.Equal(t, 0, w.Body.Len())
}

func TestRouteContextHoldsFullPath(t *testing.T) {
	router := New()

	// Test routes
	routes := []string{
		"/simple",
		"/project/:name",
		"/",
		"/news/home",
		"/news",
		"/simple-two/one",
		"/simple-two/one-two",
		"/project/:name/build/*params",
		"/project/:name/bui",
		"/user/:id/status",
		"/user/:id",
		"/user/:id/profile",
	}

	for _, route := range routes {
		actualRoute := route
		router.GET(route, func(c *Context) {
			// For each defined route context should contain its full path
			assert.Equal(t, actualRoute, c.FullPath())
			c.AbortWithStatus(http.StatusOK)
		})
	}

	for _, route := range routes {
		w := performRequest(router, http.MethodGet, route)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Test not found
	router.Use(func(c *Context) {
		// For not found routes full path is empty
		assert.Equal(t, "", c.FullPath())
	})

	w := performRequest(router, http.MethodGet, "/not-found")
	assert.Equal(t, http.StatusNotFound, w.Code)
}
