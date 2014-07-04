package gin

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func runHandler(B *testing.B, handler HandlerFunc) {
	req, err := http.NewRequest("GET", "http://localhost/foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	c := &Context{
		Writer: &responseWriter{httptest.NewRecorder(), 0, false},
		Req:    req,
		index:  0,
	}

	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		c.index = 0
		handler(c)
	}
}

func runRequest(B *testing.B, r *Engine, path string) {
	// create fake request
	url := "http://localhost" + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// create fake writes
	w := httptest.NewRecorder()

	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkMiddlewareLogger(B *testing.B) {
	runHandler(B, Logger())
}

func BenchmarkDefaultOnlyPing(B *testing.B) {
	r := New()
	r.GET("/ping", func(c *Context) {
		c.String(200, "pong")
	})
	runRequest(B, r, "/ping")
}

func BenchmarkDefaultPing(B *testing.B) {
	r := Default()
	r.GET("/ping", func(c *Context) {
		c.String(200, "pong")
	})
	runRequest(B, r, "/ping")
}
