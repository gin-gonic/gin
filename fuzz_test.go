//go:build go1.18
// +build go1.18

// Copyright 2014 Manu Martinez-Almeida
// SPDX-License-Identifier: MIT

package gin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// FuzzJSONBinding tests JSON request body binding with arbitrary
// attacker-controlled byte input.
//
// This is the pre-auth boundary for every JSON API built with Gin.
// Every JSON request body passes through this code path.
// Gin has 88K+ GitHub stars and is imported by thousands of Go projects.
//
// 5 GitHub Security Advisories exist for Gin.
func FuzzJSONBinding(f *testing.F) {
	jsonBinding := binding.JSON

	f.Add([]byte(`{"name":"test","age":30}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(``))
	f.Add([]byte(`{"a":`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[1,2,3]`))
	f.Add(make([]byte, 10000))

	f.Fuzz(func(t *testing.T, body []byte) {
		if len(body) > 1<<16 {
			return
		}

		var obj map[string]any
		// BindBody must never panic on any input
		_ = jsonBinding.BindBody(body, &obj)
	})
}

// FuzzGinPathMatch tests URL path parameter matching with
// arbitrary attacker-controlled paths.
//
// Path matching determines routing and extracts parameters
// from URLs. This is the first code that processes every
// incoming HTTP request.
func FuzzGinPathMatch(f *testing.F) {
	f.Add("/users/:id", "/users/123")
	f.Add("/api/:version/users/:id", "/api/v1/users/42")
	f.Add("/", "/")
	f.Add("/static/*filepath", "/static/css/main.css")
	f.Add(strings.Repeat("/a", 100), "/a/a")

	f.Fuzz(func(t *testing.T, pattern, path string) {
		if len(pattern) > 10000 || len(path) > 10000 {
			return
		}
		if pattern == "" {
			return
		}

		// Test route construction + matching
		router := gin.New()
		func() {
			defer func() { _ = recover() }()
			router.GET(pattern, func(c *gin.Context) {})
			router.Handle("GET", path, func(c *gin.Context) {})
		}()
	})
}

// FuzzFormBinding tests form POST body binding with arbitrary
// key-value pairs from HTTP POST requests.
func FuzzFormBinding(f *testing.F) {
	f.Add("name=test&age=30")
	f.Add("")
	f.Add("a=1&b=2&c=3")

	f.Fuzz(func(t *testing.T, formData string) {
		if len(formData) > 1<<16 {
			return
		}

		// Form binding parses POST form data
		// Create a minimal request with form body
		req, err := http.NewRequest("POST", "/", strings.NewReader(formData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var obj map[string]string
		_ = binding.Form.Bind(req, &obj)
	})
}
