// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthSucceed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	r := New()
	accounts := Accounts{"admin": "password"}
	r.Use(BasicAuth(accounts))

	r.GET("/login", func(c *Context) {
		c.String(200, "autorized")
	})

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response code should be Ok, was: %s", w.Code)
	}
	bodyAsString := w.Body.String()

	if bodyAsString != "autorized" {
		t.Errorf("Response body should be `autorized`, was  %s", bodyAsString)
	}
}

func TestBasicAuth401(t *testing.T) {
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	r := New()
	accounts := Accounts{"foo": "bar"}
	r.Use(BasicAuth(accounts))

	r.GET("/login", func(c *Context) {
		c.String(200, "autorized")
	})

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Response code should be Not autorized, was: %s", w.Code)
	}

	if w.HeaderMap.Get("WWW-Authenticate") != "Basic realm=\"Authorization Required\"" {
		t.Errorf("WWW-Authenticate header is incorrect: %s", w.HeaderMap.Get("Content-Type"))
	}
}
