// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func setupHTMLFiles(t *testing.T, mode string, tls bool, loadMethod func(*Engine)) *httptest.Server {
	SetMode(mode)
	defer SetMode(TestMode)

	var router *Engine
	captureOutput(t, func() {
		router = New()
		router.Delims("{[{", "}]}")
		router.SetFuncMap(template.FuncMap{
			"formatAsDate": formatAsDate,
		})
		loadMethod(router)
		router.GET("/test", func(c *Context) {
			c.HTML(http.StatusOK, "hello.tmpl", map[string]string{"name": "world"})
		})
		router.GET("/raw", func(c *Context) {
			c.HTML(http.StatusOK, "raw.tmpl", map[string]any{
				"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC), //nolint:gofumpt
			})
		})
	})

	var ts *httptest.Server

	if tls {
		ts = httptest.NewTLSServer(router)
	} else {
		ts = httptest.NewServer(router)
	}

	return ts
}

func TestLoadHTMLGlobDebugMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		DebugMode,
		false,
		func(router *Engine) {
			router.LoadHTMLGlob("./testdata/template/*")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestH2c(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	r := Default()
	r.UseH2C = true
	r.GET("/", func(c *Context) {
		c.String(200, "<h1>Hello world</h1>")
	})
	go func() {
		err := http.Serve(ln, r.Handler())
		if err != nil {
			t.Log(err)
		}
	}()
	defer ln.Close()

	url := "http://" + ln.Addr().String() + "/"

	httpClient := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}

	res, err := httpClient.Get(url)
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobTestMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		false,
		func(router *Engine) {
			router.LoadHTMLGlob("./testdata/template/*")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobReleaseMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		ReleaseMode,
		false,
		func(router *Engine) {
			router.LoadHTMLGlob("./testdata/template/*")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobUsingTLS(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		DebugMode,
		true,
		func(router *Engine) {
			router.LoadHTMLGlob("./testdata/template/*")
		},
	)
	defer ts.Close()

	// Use InsecureSkipVerify for avoiding `x509: certificate signed by unknown authority` error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobFromFuncMap(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		DebugMode,
		false,
		func(router *Engine) {
			router.LoadHTMLGlob("./testdata/template/*")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/raw")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Date: 2017/07/01", string(resp))
}

func init() {
	SetMode(TestMode)
}

func TestCreateEngine(t *testing.T) {
	router := New()
	assert.Equal(t, "/", router.basePath)
	assert.Equal(t, router.engine, router)
	assert.Empty(t, router.Handlers)
}

func TestLoadHTMLFilesTestMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFiles("./testdata/template/hello.tmpl", "./testdata/template/raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFilesDebugMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		DebugMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFiles("./testdata/template/hello.tmpl", "./testdata/template/raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFilesReleaseMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		ReleaseMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFiles("./testdata/template/hello.tmpl", "./testdata/template/raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFilesUsingTLS(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		true,
		func(router *Engine) {
			router.LoadHTMLFiles("./testdata/template/hello.tmpl", "./testdata/template/raw.tmpl")
		},
	)
	defer ts.Close()

	// Use InsecureSkipVerify for avoiding `x509: certificate signed by unknown authority` error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFilesFuncMap(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFiles("./testdata/template/hello.tmpl", "./testdata/template/raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/raw")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Date: 2017/07/01", string(resp))
}

var tmplFS = http.Dir("testdata/template")

func TestLoadHTMLFSTestMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFS(tmplFS, "hello.tmpl", "raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFSDebugMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		DebugMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFS(tmplFS, "hello.tmpl", "raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFSReleaseMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		ReleaseMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFS(tmplFS, "hello.tmpl", "raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFSUsingTLS(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		true,
		func(router *Engine) {
			router.LoadHTMLFS(tmplFS, "hello.tmpl", "raw.tmpl")
		},
	)
	defer ts.Close()

	// Use InsecureSkipVerify for avoiding `x509: certificate signed by unknown authority` error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLFSFuncMap(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		TestMode,
		false,
		func(router *Engine) {
			router.LoadHTMLFS(tmplFS, "hello.tmpl", "raw.tmpl")
		},
	)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/raw")
	if err != nil {
		t.Error(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Date: 2017/07/01", string(resp))
}

func TestAddRoute(t *testing.T) {
	router := New()
	router.addRoute(http.MethodGet, "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 1)
	assert.NotNil(t, router.trees.get(http.MethodGet))
	assert.Nil(t, router.trees.get(http.MethodPost))

	router.addRoute(http.MethodPost, "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 2)
	assert.NotNil(t, router.trees.get(http.MethodGet))
	assert.NotNil(t, router.trees.get(http.MethodPost))

	router.addRoute(http.MethodPost, "/post", HandlersChain{func(_ *Context) {}})
	assert.Len(t, router.trees, 2)
}

func TestAddRouteFails(t *testing.T) {
	router := New()
	assert.Panics(t, func() { router.addRoute("", "/", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute(http.MethodGet, "a", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute(http.MethodGet, "/", HandlersChain{}) })

	router.addRoute(http.MethodPost, "/post", HandlersChain{func(_ *Context) {}})
	assert.Panics(t, func() {
		router.addRoute(http.MethodPost, "/post", HandlersChain{func(_ *Context) {}})
	})
}

func TestCreateDefaultRouter(t *testing.T) {
	router := Default()
	assert.Len(t, router.Handlers, 2)
}

func TestNoRouteWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoRoute(middleware0)
	assert.Nil(t, router.Handlers)
	assert.Len(t, router.noRoute, 1)
	assert.Len(t, router.allNoRoute, 1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware0)

	router.NoRoute(middleware1, middleware0)
	assert.Len(t, router.noRoute, 2)
	assert.Len(t, router.allNoRoute, 2)
	compareFunc(t, router.noRoute[0], middleware1)
	compareFunc(t, router.allNoRoute[0], middleware1)
	compareFunc(t, router.noRoute[1], middleware0)
	compareFunc(t, router.allNoRoute[1], middleware0)
}

func TestNoRouteWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoRoute(middleware0)
	assert.Len(t, router.allNoRoute, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoRoute, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware1)
	compareFunc(t, router.allNoRoute[2], middleware0)
}

func TestNoMethodWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoMethod(middleware0)
	assert.Empty(t, router.Handlers)
	assert.Len(t, router.noMethod, 1)
	assert.Len(t, router.allNoMethod, 1)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware0)

	router.NoMethod(middleware1, middleware0)
	assert.Len(t, router.noMethod, 2)
	assert.Len(t, router.allNoMethod, 2)
	compareFunc(t, router.noMethod[0], middleware1)
	compareFunc(t, router.allNoMethod[0], middleware1)
	compareFunc(t, router.noMethod[1], middleware0)
	compareFunc(t, router.allNoMethod[1], middleware0)
}

func TestRebuild404Handlers(t *testing.T) {
}

func TestNoMethodWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoMethod(middleware0)
	assert.Len(t, router.allNoMethod, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noMethod, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware2)
	compareFunc(t, router.allNoMethod[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoMethod, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noMethod, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware2)
	compareFunc(t, router.allNoMethod[1], middleware1)
	compareFunc(t, router.allNoMethod[2], middleware0)
}

func compareFunc(t *testing.T, a, b any) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func TestListOfRoutes(t *testing.T) {
	router := New()
	router.GET("/favicon.ico", handlerTest1)
	router.GET("/", handlerTest1)
	group := router.Group("/users")
	{
		group.GET("/", handlerTest2)
		group.GET("/:id", handlerTest1)
		group.POST("/:id", handlerTest2)
	}
	router.Static("/static", ".")

	list := router.Routes()

	assert.Len(t, list, 7)
	assertRoutePresent(t, list, RouteInfo{
		Method:  http.MethodGet,
		Path:    "/favicon.ico",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  http.MethodGet,
		Path:    "/users/",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest2$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  http.MethodGet,
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  http.MethodPost,
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest2$",
	})
}

func TestEngineHandleContext(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.Request.URL.Path = "/v2"
		r.HandleContext(c)
	})
	v2 := r.Group("/v2")
	{
		v2.GET("/", func(c *Context) {})
	}

	assert.NotPanics(t, func() {
		w := PerformRequest(r, http.MethodGet, "/")
		assert.Equal(t, 301, w.Code)
	})
}

func TestEngineHandleContextManyReEntries(t *testing.T) {
	expectValue := 10000

	var handlerCounter, middlewareCounter int64

	r := New()
	r.Use(func(c *Context) {
		atomic.AddInt64(&middlewareCounter, 1)
	})
	r.GET("/:count", func(c *Context) {
		countStr := c.Param("count")
		count, err := strconv.Atoi(countStr)
		require.NoError(t, err)

		n, err := c.Writer.Write([]byte("."))
		require.NoError(t, err)
		assert.Equal(t, 1, n)

		switch {
		case count > 0:
			c.Request.URL.Path = "/" + strconv.Itoa(count-1)
			r.HandleContext(c)
		}
	}, func(c *Context) {
		atomic.AddInt64(&handlerCounter, 1)
	})

	assert.NotPanics(t, func() {
		w := PerformRequest(r, http.MethodGet, "/"+strconv.Itoa(expectValue-1)) // include 0 value
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expectValue, w.Body.Len())
	})

	assert.Equal(t, int64(expectValue), handlerCounter)
	assert.Equal(t, int64(expectValue), middlewareCounter)
}

func TestEngineHandleContextPreventsMiddlewareReEntry(t *testing.T) {
	// given
	var handlerCounterV1, handlerCounterV2, middlewareCounterV1 int64

	r := New()
	v1 := r.Group("/v1")
	{
		v1.Use(func(c *Context) {
			atomic.AddInt64(&middlewareCounterV1, 1)
		})
		v1.GET("/test", func(c *Context) {
			atomic.AddInt64(&handlerCounterV1, 1)
			c.Status(http.StatusOK)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/test", func(c *Context) {
			c.Request.URL.Path = "/v1/test"
			r.HandleContext(c)
		}, func(c *Context) {
			atomic.AddInt64(&handlerCounterV2, 1)
		})
	}

	// when
	responseV1 := PerformRequest(r, "GET", "/v1/test")
	responseV2 := PerformRequest(r, "GET", "/v2/test")

	// then
	assert.Equal(t, 200, responseV1.Code)
	assert.Equal(t, 200, responseV2.Code)
	assert.Equal(t, int64(2), handlerCounterV1)
	assert.Equal(t, int64(2), middlewareCounterV1)
	assert.Equal(t, int64(1), handlerCounterV2)
}

func TestPrepareTrustedCIRDsWith(t *testing.T) {
	r := New()

	// valid ipv4 cidr
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("0.0.0.0/0")}
		err := r.SetTrustedProxies([]string{"0.0.0.0/0"})

		require.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedCIDRs)
	}

	// invalid ipv4 cidr
	{
		err := r.SetTrustedProxies([]string{"192.168.1.33/33"})

		require.Error(t, err)
	}

	// valid ipv4 address
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("192.168.1.33/32")}

		err := r.SetTrustedProxies([]string{"192.168.1.33"})

		require.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedCIDRs)
	}

	// invalid ipv4 address
	{
		err := r.SetTrustedProxies([]string{"192.168.1.256"})

		require.Error(t, err)
	}

	// valid ipv6 address
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("2002:0000:0000:1234:abcd:ffff:c0a8:0101/128")}
		err := r.SetTrustedProxies([]string{"2002:0000:0000:1234:abcd:ffff:c0a8:0101"})

		require.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedCIDRs)
	}

	// invalid ipv6 address
	{
		err := r.SetTrustedProxies([]string{"gggg:0000:0000:1234:abcd:ffff:c0a8:0101"})

		require.Error(t, err)
	}

	// valid ipv6 cidr
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("::/0")}
		err := r.SetTrustedProxies([]string{"::/0"})

		require.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedCIDRs)
	}

	// invalid ipv6 cidr
	{
		err := r.SetTrustedProxies([]string{"gggg:0000:0000:1234:abcd:ffff:c0a8:0101/129"})

		require.Error(t, err)
	}

	// valid combination
	{
		expectedTrustedCIDRs := []*net.IPNet{
			parseCIDR("::/0"),
			parseCIDR("192.168.0.0/16"),
			parseCIDR("172.16.0.1/32"),
		}
		err := r.SetTrustedProxies([]string{
			"::/0",
			"192.168.0.0/16",
			"172.16.0.1",
		})

		require.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedCIDRs)
	}

	// invalid combination
	{
		err := r.SetTrustedProxies([]string{
			"::/0",
			"192.168.0.0/16",
			"172.16.0.256",
		})

		require.Error(t, err)
	}

	// nil value
	{
		err := r.SetTrustedProxies(nil)

		assert.Nil(t, r.trustedCIDRs)
		require.NoError(t, err)
	}
}

func parseCIDR(cidr string) *net.IPNet {
	_, parsedCIDR, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println(err)
	}
	return parsedCIDR
}

func assertRoutePresent(t *testing.T, gotRoutes RoutesInfo, wantRoute RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			assert.Regexp(t, wantRoute.Handler, gotRoute.Handler)
			return
		}
	}
	t.Errorf("route not found: %v", wantRoute)
}

func handlerTest1(c *Context) {}
func handlerTest2(c *Context) {}

func TestNewOptionFunc(t *testing.T) {
	fc := func(e *Engine) {
		e.GET("/test1", handlerTest1)
		e.GET("/test2", handlerTest2)

		e.Use(func(c *Context) {
			c.Next()
		})
	}

	r := New(fc)

	routes := r.Routes()
	assertRoutePresent(t, routes, RouteInfo{Path: "/test1", Method: http.MethodGet, Handler: "github.com/gin-gonic/gin.handlerTest1"})
	assertRoutePresent(t, routes, RouteInfo{Path: "/test2", Method: http.MethodGet, Handler: "github.com/gin-gonic/gin.handlerTest2"})
}

func TestWithOptionFunc(t *testing.T) {
	r := New()

	r.With(func(e *Engine) {
		e.GET("/test1", handlerTest1)
		e.GET("/test2", handlerTest2)

		e.Use(func(c *Context) {
			c.Next()
		})
	})

	routes := r.Routes()
	assertRoutePresent(t, routes, RouteInfo{Path: "/test1", Method: http.MethodGet, Handler: "github.com/gin-gonic/gin.handlerTest1"})
	assertRoutePresent(t, routes, RouteInfo{Path: "/test2", Method: http.MethodGet, Handler: "github.com/gin-gonic/gin.handlerTest2"})
}

type Birthday string

func (b *Birthday) UnmarshalParam(param string) error {
	*b = Birthday(strings.ReplaceAll(param, "-", "/"))
	return nil
}

func TestCustomUnmarshalStruct(t *testing.T) {
	route := Default()
	var request struct {
		Birthday Birthday `form:"birthday"`
	}
	route.GET("/test", func(ctx *Context) {
		_ = ctx.BindQuery(&request)
		ctx.JSON(200, request.Birthday)
	})
	req := httptest.NewRequest(http.MethodGet, "/test?birthday=2000-01-01", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `"2000/01/01"`, w.Body.String())
}

// Test the fix for https://github.com/gin-gonic/gin/issues/4002
func TestMethodNotAllowedNoRoute(t *testing.T) {
	g := New()
	g.HandleMethodNotAllowed = true

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	assert.NotPanics(t, func() { g.ServeHTTP(resp, req) })
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// TestTreesMapInitialization tests that treesMap is properly initialized
func TestTreesMapInitialization(t *testing.T) {
	router := New()

	// Verify treesMap is initialized as an empty map
	assert.NotNil(t, router.treesMap)
	assert.Empty(t, router.treesMap)

	// Verify trees slice is also initialized
	assert.NotNil(t, router.trees)
	assert.Empty(t, router.trees)
}

// TestTreesMapSynchronization tests that treesMap stays in sync with trees slice
func TestTreesMapSynchronization(t *testing.T) {
	router := New()

	// Add a GET route
	router.addRoute(http.MethodGet, "/", HandlersChain{func(_ *Context) {}})

	// Verify both trees and treesMap are updated
	assert.Len(t, router.trees, 1)
	assert.Len(t, router.treesMap, 1)

	// Verify treesMap contains the correct method
	root, exists := router.treesMap[http.MethodGet]
	assert.True(t, exists)
	assert.NotNil(t, root)

	// Verify trees slice also contains the same method
	treeRoot := router.trees.get(http.MethodGet)
	assert.Equal(t, root, treeRoot)

	// Add a POST route
	router.addRoute(http.MethodPost, "/post", HandlersChain{func(_ *Context) {}})

	// Verify both are updated again
	assert.Len(t, router.trees, 2)
	assert.Len(t, router.treesMap, 2)

	// Verify both methods exist in treesMap
	_, getExists := router.treesMap[http.MethodGet]
	_, postExists := router.treesMap[http.MethodPost]
	assert.True(t, getExists)
	assert.True(t, postExists)

	// Verify trees slice also has both methods
	assert.NotNil(t, router.trees.get(http.MethodGet))
	assert.NotNil(t, router.trees.get(http.MethodPost))
}

// TestTreesMapLookupPerformance tests the O(1) lookup performance
func TestTreesMapLookupPerformance(t *testing.T) {
	router := New()

	// Add multiple routes with different methods
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	// Test that all methods can be found via treesMap (O(1) lookup)
	for _, method := range methods {
		root, exists := router.treesMap[method]
		assert.True(t, exists, "Method %s should exist in treesMap", method)
		assert.NotNil(t, root, "Root for method %s should not be nil", method)
	}

	// Test that non-existent method returns false
	_, exists := router.treesMap["NONEXISTENT"]
	assert.False(t, exists, "Non-existent method should not exist in treesMap")
}

// TestTreesMapWithMultipleRoutesPerMethod tests treesMap with multiple routes per HTTP method
func TestTreesMapWithMultipleRoutesPerMethod(t *testing.T) {
	router := New()

	// Add multiple routes for the same method
	router.addRoute(http.MethodGet, "/", HandlersChain{func(_ *Context) {}})
	router.addRoute(http.MethodGet, "/users", HandlersChain{func(_ *Context) {}})
	router.addRoute(http.MethodGet, "/users/:id", HandlersChain{func(_ *Context) {}})

	// Verify treesMap still has only one entry for GET method
	assert.Len(t, router.treesMap, 1)

	// Verify the GET method exists and points to the root node
	root, exists := router.treesMap[http.MethodGet]
	assert.True(t, exists)
	assert.NotNil(t, root)

	// Verify trees slice also has only one entry for GET
	assert.Len(t, router.trees, 1)
	assert.Equal(t, root, router.trees.get(http.MethodGet))
}

// TestTreesMapConcurrentAccess tests that treesMap works correctly with sequential access
// Note: treesMap is not designed for concurrent writes, as route registration is typically
// done during application initialization, not at runtime
func TestTreesMapSequentialAccess(t *testing.T) {
	router := New()

	// Add routes sequentially (which is the typical use case)
	for i := 0; i < 10; i++ {
		method := http.MethodGet + strconv.Itoa(i)
		router.addRoute(method, "/route"+strconv.Itoa(i), HandlersChain{func(_ *Context) {}})
	}

	// Verify all routes were added successfully
	assert.Len(t, router.trees, 10)
	assert.Len(t, router.treesMap, 10)

	// Verify all methods exist in treesMap
	for i := 0; i < 10; i++ {
		method := http.MethodGet + strconv.Itoa(i)
		_, exists := router.treesMap[method]
		assert.True(t, exists, "Method %s should exist in treesMap", method)
	}
}

// TestTreesMapWithSpecialMethods tests treesMap with special HTTP methods
func TestTreesMapWithSpecialMethods(t *testing.T) {
	router := New()

	// Test with various HTTP methods including less common ones
	specialMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect,
		http.MethodTrace,
	}

	for _, method := range specialMethods {
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	// Verify all methods exist in treesMap
	assert.Len(t, router.treesMap, len(specialMethods))
	assert.Len(t, router.trees, len(specialMethods))

	for _, method := range specialMethods {
		root, exists := router.treesMap[method]
		assert.True(t, exists, "Method %s should exist in treesMap", method)
		assert.NotNil(t, root, "Root for method %s should not be nil", method)

		// Verify consistency with trees slice
		treeRoot := router.trees.get(method)
		assert.Equal(t, root, treeRoot, "treesMap and trees should point to same root for method %s", method)
	}
}

// TestTreesMapEmptyLookup tests lookup behavior with empty treesMap
func TestTreesMapEmptyLookup(t *testing.T) {
	router := New()

	// Verify empty treesMap behavior
	assert.Empty(t, router.treesMap)

	// Test lookup of non-existent method
	_, exists := router.treesMap[http.MethodGet]
	assert.False(t, exists)

	// Test lookup of empty string method
	_, exists = router.treesMap[""]
	assert.False(t, exists)
}

// TestTreesMapConsistencyWithTrees tests that treesMap and trees remain consistent
func TestTreesMapConsistencyWithTrees(t *testing.T) {
	router := New()

	// Add routes and verify consistency at each step
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut}

	for i, method := range methods {
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})

		// After each addition, verify consistency
		assert.Len(t, router.trees, i+1)
		assert.Len(t, router.treesMap, i+1)

		// Verify each method exists in both structures
		for j := 0; j <= i; j++ {
			currentMethod := methods[j]

			// Check treesMap
			mapRoot, mapExists := router.treesMap[currentMethod]
			assert.True(t, mapExists, "Method %s should exist in treesMap after %d additions", currentMethod, i+1)
			assert.NotNil(t, mapRoot)

			// Check trees slice
			treeRoot := router.trees.get(currentMethod)
			assert.NotNil(t, treeRoot, "Method %s should exist in trees after %d additions", currentMethod, i+1)

			// Verify they point to the same root
			assert.Equal(t, mapRoot, treeRoot, "treesMap and trees should point to same root for method %s", currentMethod)
		}
	}
}

// BenchmarkTreesMapLookup benchmarks the O(1) treesMap lookup performance
func BenchmarkTreesMapLookup(b *testing.B) {
	router := New()

	// Add multiple HTTP methods to test lookup performance
	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect,
		http.MethodTrace, "CUSTOM1", "CUSTOM2", "CUSTOM3",
	}

	for _, method := range methods {
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Test O(1) lookup via treesMap
			_, exists := router.treesMap[http.MethodGet]
			if !exists {
				b.Fatal("Expected method to exist")
			}
		}
	})
}

// BenchmarkTreesSliceLookup benchmarks the O(n) trees slice lookup for comparison
func BenchmarkTreesSliceLookup(b *testing.B) {
	router := New()

	// Add multiple HTTP methods to test lookup performance
	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect,
		http.MethodTrace, "CUSTOM1", "CUSTOM2", "CUSTOM3",
	}

	for _, method := range methods {
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Test O(n) lookup via trees slice
			root := router.trees.get(http.MethodGet)
			if root == nil {
				b.Fatal("Expected method to exist")
			}
		}
	})
}

// BenchmarkTreesMapLookupWithManyMethods benchmarks treesMap lookup with many HTTP methods
func BenchmarkTreesMapLookupWithManyMethods(b *testing.B) {
	router := New()

	// Add many HTTP methods to test scalability
	for i := 0; i < 50; i++ {
		method := "METHOD" + strconv.Itoa(i)
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	// Add the method we'll be looking up
	router.addRoute(http.MethodGet, "/get", HandlersChain{func(_ *Context) {}})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Test O(1) lookup via treesMap (should be constant time regardless of total methods)
			_, exists := router.treesMap[http.MethodGet]
			if !exists {
				b.Fatal("Expected method to exist")
			}
		}
	})
}

// BenchmarkTreesSliceLookupWithManyMethods benchmarks trees slice lookup with many HTTP methods
func BenchmarkTreesSliceLookupWithManyMethods(b *testing.B) {
	router := New()

	// Add many HTTP methods to test scalability
	for i := 0; i < 50; i++ {
		method := "METHOD" + strconv.Itoa(i)
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	// Add the method we'll be looking up
	router.addRoute(http.MethodGet, "/get", HandlersChain{func(_ *Context) {}})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Test O(n) lookup via trees slice (should be slower with more methods)
			root := router.trees.get(http.MethodGet)
			if root == nil {
				b.Fatal("Expected method to exist")
			}
		}
	})
}

// TestTreesMapIntegrationWithHTTPRequests tests treesMap integration with actual HTTP requests
func TestTreesMapIntegrationWithHTTPRequests(t *testing.T) {
	router := New()

	// Add routes with different methods
	router.GET("/get", func(c *Context) {
		c.String(200, "GET response")
	})
	router.POST("/post", func(c *Context) {
		c.String(200, "POST response")
	})
	router.PUT("/put", func(c *Context) {
		c.String(200, "PUT response")
	})

	// Test that HTTP requests work correctly with treesMap optimization
	testCases := []struct {
		method           string
		path             string
		expectedResponse string
	}{
		{http.MethodGet, "/get", "GET response"},
		{http.MethodPost, "/post", "POST response"},
		{http.MethodPut, "/put", "PUT response"},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, tc.expectedResponse, w.Body.String())
	}

	// Verify treesMap was populated correctly
	assert.Len(t, router.treesMap, 3)
	assert.Contains(t, router.treesMap, http.MethodGet)
	assert.Contains(t, router.treesMap, http.MethodPost)
	assert.Contains(t, router.treesMap, http.MethodPut)
}

// TestTreesMapWithDuplicateMethods tests behavior when adding routes with duplicate methods
func TestTreesMapWithDuplicateMethods(t *testing.T) {
	router := New()

	// Add first route with GET method
	router.addRoute(http.MethodGet, "/first", HandlersChain{func(_ *Context) {}})

	// Verify initial state
	assert.Len(t, router.treesMap, 1)
	assert.Len(t, router.trees, 1)

	// Add second route with same GET method (should not create new entry)
	router.addRoute(http.MethodGet, "/second", HandlersChain{func(_ *Context) {}})

	// Verify treesMap still has only one entry for GET
	assert.Len(t, router.treesMap, 1)
	assert.Len(t, router.trees, 1)

	// Verify the GET method still exists and points to the same root
	root, exists := router.treesMap[http.MethodGet]
	assert.True(t, exists)
	assert.NotNil(t, root)
	assert.Equal(t, root, router.trees.get(http.MethodGet))
}

// TestTreesMapWithEmptyMethod tests edge case with empty method string
func TestTreesMapWithEmptyMethod(t *testing.T) {
	router := New()

	// This should panic due to validation in addRoute
	assert.Panics(t, func() {
		router.addRoute("", "/path", HandlersChain{func(_ *Context) {}})
	})

	// Verify treesMap remains empty
	assert.Empty(t, router.treesMap)
	assert.Empty(t, router.trees)
}

// TestTreesMapWithNilHandlers tests edge case with nil handlers
func TestTreesMapWithNilHandlers(t *testing.T) {
	router := New()

	// This should panic due to validation in addRoute
	assert.Panics(t, func() {
		router.addRoute(http.MethodGet, "/path", nil)
	})

	// Verify treesMap remains empty
	assert.Empty(t, router.treesMap)
	assert.Empty(t, router.trees)
}

// TestTreesMapWithInvalidPath tests edge case with invalid path
func TestTreesMapWithInvalidPath(t *testing.T) {
	router := New()

	// This should panic due to validation in addRoute
	assert.Panics(t, func() {
		router.addRoute(http.MethodGet, "invalid-path", HandlersChain{func(_ *Context) {}})
	})

	// Verify treesMap remains empty
	assert.Empty(t, router.treesMap)
	assert.Empty(t, router.trees)
}

// TestTreesMapMemoryUsage tests that treesMap doesn't cause memory leaks
func TestTreesMapMemoryUsage(t *testing.T) {
	router := New()

	// Add and remove routes multiple times
	for i := 0; i < 100; i++ {
		method := "METHOD" + strconv.Itoa(i)
		router.addRoute(method, "/path"+strconv.Itoa(i), HandlersChain{func(_ *Context) {}})
	}

	// Verify all routes were added
	assert.Len(t, router.treesMap, 100)
	assert.Len(t, router.trees, 100)

	// Verify memory usage is reasonable (treesMap should have 100 entries)
	assert.Equal(t, 100, len(router.treesMap))

	// Test that lookups still work
	for i := 0; i < 100; i++ {
		method := "METHOD" + strconv.Itoa(i)
		_, exists := router.treesMap[method]
		assert.True(t, exists, "Method %s should exist", method)
	}
}

// TestTreesMapWithVeryLongMethodNames tests treesMap with very long method names
func TestTreesMapWithVeryLongMethodNames(t *testing.T) {
	router := New()

	// Create a very long method name
	longMethod := strings.Repeat("VERY_LONG_METHOD_NAME_", 10)

	router.addRoute(longMethod, "/long", HandlersChain{func(_ *Context) {}})

	// Verify the long method was added correctly
	assert.Len(t, router.treesMap, 1)
	assert.Len(t, router.trees, 1)

	// Verify lookup works with long method name
	root, exists := router.treesMap[longMethod]
	assert.True(t, exists)
	assert.NotNil(t, root)
	assert.Equal(t, root, router.trees.get(longMethod))
}

// TestTreesMapWithSpecialCharacters tests treesMap with method names containing special characters
func TestTreesMapWithSpecialCharacters(t *testing.T) {
	router := New()

	// Test with method names containing special characters
	specialMethods := []string{
		"METHOD-WITH-DASH",
		"METHOD_WITH_UNDERSCORE",
		"METHOD.WITH.DOTS",
		"METHOD123WITH456NUMBERS",
		"method-with-lowercase",
	}

	for _, method := range specialMethods {
		router.addRoute(method, "/"+strings.ToLower(method), HandlersChain{func(_ *Context) {}})
	}

	// Verify all special methods were added
	assert.Len(t, router.treesMap, len(specialMethods))
	assert.Len(t, router.trees, len(specialMethods))

	// Verify all special methods can be looked up
	for _, method := range specialMethods {
		root, exists := router.treesMap[method]
		assert.True(t, exists, "Method %s should exist in treesMap", method)
		assert.NotNil(t, root, "Root for method %s should not be nil", method)
	}
}
