package gin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeWriter struct{}

func (_ FakeWriter) Write(d []byte) (int, error) {
	return 0, nil
}

func (_ FakeWriter) WriteString(d string) (int, error) {
	return 0, nil
}

func runRequest(B *testing.B, r *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := httptest.NewRecorder()

	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkOneRoute(B *testing.B) {
	router := New()
	router.GET("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkManyHandlers(B *testing.B) {
	DefaultWriter = FakeWriter{}
	//router := Default()
	router := New()
	router.Use(Recovery(), Logger())
	router.Use(func(c *Context) {})
	router.GET("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark5Params(B *testing.B) {
	DefaultWriter = new(bytes.Buffer)
	router := New()
	router.Use(func(c *Context) {})
	router.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {})
	runRequest(B, router, "GET", "/param/path/to/parameter/john/12345")
}

func BenchmarkOneRouteJSON(B *testing.B) {
	router := New()
	data := H{
		"status": "ok",
	}
	router.GET("/json", func(c *Context) {
		//c.JSON(200, data)
		c.Writer.WriteHeader(200)
		json.NewEncoder(c.Writer).Encode(data)
	})
	runRequest(B, router, "GET", "/json")
}

func BenchmarkOneRouteString(B *testing.B) {
	router := New()
	router.GET("/text", func(c *Context) {
		c.String(200, "this is a plain text")
	})
	runRequest(B, router, "GET", "/text")
}

func BenchmarkManyRoutes(B *testing.B) {
	router := New()
	router.Any("/ping", func(c *Context) {})
	runRequest(B, router, "PUT", "/ping")
}

func Benchmark404(B *testing.B) {
	router := New()
	router.Any("/something", func(c *Context) {})
	router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark404Many(B *testing.B) {
	router := New()
	router.GET("/", func(c *Context) {})
	router.GET("/path/to/something", func(c *Context) {})
	router.GET("/post/:id", func(c *Context) {})
	router.GET("/view/:id", func(c *Context) {})
	router.GET("/favicon.ico", func(c *Context) {})
	router.GET("/robots.txt", func(c *Context) {})
	router.GET("/delete/:id", func(c *Context) {})
	router.GET("/user/:id/:mode", func(c *Context) {})

	router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/viewfake")
}
