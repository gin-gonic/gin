// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

// zeroReader is an unbounded source of zero bytes. It models an attacker
// streaming an arbitrarily large request body without needing to
// materialize a large byte slice in the benchmark itself.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func BenchmarkBindJSONSmallBody_LimitEnabled(b *testing.B) {
	router := New()
	router.MaxRequestBodyBytes = defaultMaxRequestBodyBytes
	body := `{"foo":"bar","bar":"foo"}`
	c := CreateTestContextOnly(httptest.NewRecorder(), router)

	b.ReportAllocs()
	for b.Loop() {
		c.reset()
		c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		var obj struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}
		if err := c.BindJSON(&obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBindJSONSmallBody_LimitDisabled(b *testing.B) {
	router := New()
	router.MaxRequestBodyBytes = 0
	body := `{"foo":"bar","bar":"foo"}`
	c := CreateTestContextOnly(httptest.NewRecorder(), router)

	b.ReportAllocs()
	for b.Loop() {
		c.reset()
		c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		var obj struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}
		if err := c.BindJSON(&obj); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBindJSONLargeBody_LimitEnabled(b *testing.B) {
	router := New()
	router.MaxRequestBodyBytes = defaultMaxRequestBodyBytes
	body := `{"foo":"` + strings.Repeat("a", 1<<20) + `"}` // ~1 MiB value, under the limit
	c := CreateTestContextOnly(httptest.NewRecorder(), router)

	b.ReportAllocs()
	for b.Loop() {
		c.reset()
		c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		var obj struct {
			Foo string `json:"foo"`
		}
		if err := c.BindJSON(&obj); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBindRejectOversizedBody is the key regression benchmark: it
// proves rejection cost (and allocations) stay bounded by MaxRequestBodyBytes
// regardless of how much more data the client could have sent, instead of
// scaling with attacker-supplied body size as it did before this fix.
func BenchmarkBindRejectOversizedBody(b *testing.B) {
	router := New()
	router.MaxRequestBodyBytes = 4 << 10 // 4 KiB
	c := CreateTestContextOnly(httptest.NewRecorder(), router)

	b.ReportAllocs()
	for b.Loop() {
		c.reset()
		c.Request, _ = http.NewRequest(http.MethodPost, "/", io.NopCloser(zeroReader{}))
		var s string
		if err := c.ShouldBindWith(&s, binding.Plain); err == nil {
			b.Fatal("expected the oversized body to be rejected")
		}
	}
}

// BenchmarkLimitRequestBody isolates the per-request cost of the wrapper
// itself (no read), added to the top of every ShouldBindWith/ShouldBindBodyWith
// call regardless of body size.
func BenchmarkLimitRequestBody(b *testing.B) {
	router := New()
	router.MaxRequestBodyBytes = defaultMaxRequestBodyBytes
	c := CreateTestContextOnly(httptest.NewRecorder(), router)

	b.ReportAllocs()
	for b.Loop() {
		c.reset()
		c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader("x"))
		c.limitRequestBody()
	}
}
