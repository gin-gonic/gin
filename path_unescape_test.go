// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutePathUnescape(t *testing.T) {
	for _, mode := range []struct {
		name                 string
		raw, escaped, decode bool
	}{
		{name: "default", decode: true},
		{name: "default_no_unescape"},
		{name: "raw", raw: true, decode: true},
		{name: "raw_no_unescape", raw: true},
		{name: "escaped", escaped: true, decode: true},
		{name: "escaped_no_unescape", escaped: true},
		{name: "escaped_over_raw", raw: true, escaped: true, decode: true},
		{name: "escaped_over_raw_no_unescape", raw: true, escaped: true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			router := New()
			router.UseRawPath = mode.raw
			router.UseEscapedPath = mode.escaped
			router.UnescapePathValues = mode.decode
			router.Use(func(c *Context) {
				c.Header("Middleware-Param", c.Param("value"))
				c.Next()
			})
			for _, route := range []string{"/single/:value/end", "/catch/*value"} {
				router.GET(route, func(c *Context) {
					c.Header("Matched-Route", c.FullPath())
					c.Header("Query-Value", c.Query("q"))
					c.String(http.StatusOK, "%s", c.Param("value"))
				})
			}
			for _, tc := range []struct {
				name, encoded, decoded, raw string
			}{
				{name: "literal_plus", encoded: "a+b", decoded: "a+b"},
				{name: "encoded_plus", encoded: "a%2Bb", decoded: "a+b", raw: "a%2Bb"},
				{name: "lowercase_escape", encoded: "a%2bb", decoded: "a+b", raw: "a%2bb"},
				{name: "space", encoded: "a%20b", decoded: "a b"},
				{name: "plus_and_space", encoded: "a+b%20c", decoded: "a+b c"},
				{name: "plus_and_slash", encoded: "a+%2Fb", decoded: "a+/b", raw: "a+%2Fb"},
				{name: "encoded_once", encoded: "a+%252Fb", decoded: "a+%2Fb"},
				{name: "percent", encoded: "a+%25b", decoded: "a+%b"},
				{name: "reserved", encoded: "a+%3Fb%23c%26d%3De", decoded: "a+?b#c&d=e", raw: "a+%3Fb%23c%26d%3De"},
				{name: "unicode", encoded: "%E2%98%83+%F0%9F%98%80", decoded: "☃+😀"},
			} {
				for _, catchAll := range []bool{false, true} {
					name := tc.name + "/single"
					if catchAll {
						name = tc.name + "/catch_all"
					}
					t.Run(name, func(t *testing.T) {
						prefix, suffix, route := "/single/", "/end", "/single/:value/end"
						if catchAll {
							prefix, suffix, route = "/catch/", "", "/catch/*value"
						}
						req := httptest.NewRequest(http.MethodGet, prefix+tc.encoded+suffix+"?q=a+b", nil)
						originalURL := *req.URL
						w := httptest.NewRecorder()
						router.ServeHTTP(w, req)
						assert.Equal(t, originalURL, *req.URL)

						// The default route matcher sees decoded slashes as separators.
						if !mode.raw && !mode.escaped && !catchAll && strings.Contains(tc.decoded, "/") {
							assert.Equal(t, http.StatusNotFound, w.Code)
							return
						}
						want := tc.decoded
						if !mode.decode {
							if mode.escaped {
								want = tc.encoded
							} else if mode.raw && tc.raw != "" {
								want = tc.raw
							}
						}
						if catchAll {
							want = "/" + want
						}
						assert.Equal(t, http.StatusOK, w.Code)
						assert.Equal(t, want, w.Body.String())
						assert.Equal(t, want, w.Header().Get("Middleware-Param"))
						assert.Equal(t, route, w.Header().Get("Matched-Route"))
						assert.Equal(t, "a b", w.Header().Get("Query-Value"))
					})
				}
			}
		})
	}
}

func TestRoutePathUnescapeHTTP(t *testing.T) {
	for _, escaped := range []bool{false, true} {
		for _, http2 := range []bool{false, true} {
			name := "raw/http1"
			if escaped {
				name = "escaped/http1"
			}
			if http2 {
				name = strings.TrimSuffix(name, "http1") + "http2"
			}
			t.Run(name, func(t *testing.T) {
				router := New()
				router.UseRawPath = !escaped
				router.UseEscapedPath = escaped
				router.GET("/item/:value", func(c *Context) {
					c.String(http.StatusOK, "%s|%s", c.Param("value"), c.Query("q"))
				})
				server := httptest.NewUnstartedServer(router)
				server.EnableHTTP2 = http2
				if http2 {
					server.StartTLS()
				} else {
					server.Start()
				}
				defer server.Close()
				client := server.Client()
				client.Timeout = 5 * time.Second
				resp, err := client.Get(server.URL + "/item/a+%2Fb?q=a+b")
				require.NoError(t, err)
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "a+/b|a b", string(body))
				if http2 {
					assert.Equal(t, 2, resp.ProtoMajor)
				} else {
					assert.Equal(t, 1, resp.ProtoMajor)
				}
			})
		}
	}
}

func FuzzRoutePathUnescape(f *testing.F) {
	for _, value := range []string{"", "+", " ", "%2F", "/", "a+%2Fb", "☃+😀", "\x00\xff", "a?b#c&d=e"} {
		f.Add(value)
	}
	router := New()
	router.UseEscapedPath = true
	router.GET("/single/:value/end", func(c *Context) {
		c.String(http.StatusOK, "%s", c.Param("value"))
	})
	router.GET("/catch/*value", func(c *Context) {
		c.String(http.StatusOK, "%s", c.Param("value"))
	})
	f.Fuzz(func(t *testing.T, value string) {
		// Keep the parameter nonempty, including when the generated value is empty.
		value = "v" + value
		escaped := url.PathEscape(value)
		for _, tc := range []struct{ path, want string }{
			{path: "/single/" + escaped + "/end", want: value},
			{path: "/catch/" + escaped, want: "/" + value},
		} {
			w := PerformRequest(router, http.MethodGet, tc.path)
			require.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tc.want, w.Body.String())
		}
	})
}

func BenchmarkRoutePathUnescape(b *testing.B) {
	for _, path := range []string{"plain", "a+b", "a%2Bb", "a%20b", "a+%2Fb"} {
		b.Run(path, func(b *testing.B) {
			tree := &node{}
			tree.addRoute("/single/:value/end", HandlersChain{func(*Context) {}})
			params, skipped := getParams(), getSkippedNodes()
			requestPath := "/single/" + path + "/end"
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				*params = (*params)[:0]
				*skipped = (*skipped)[:0]
				value := tree.getValue(requestPath, params, skipped, true)
				if value.handlers == nil {
					b.Fatal("route did not match")
				}
			}
		})
	}
}
