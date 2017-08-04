// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	pairs := processAccounts(Accounts{
		"admin": "password",
		"foo":   "bar",
		"bar":   "foo",
	})

	assert.Len(t, pairs, 3)
	assert.Contains(t, pairs, authPair{
		User:  "bar",
		Value: "Basic YmFyOmZvbw==",
	})
	assert.Contains(t, pairs, authPair{
		User:  "foo",
		Value: "Basic Zm9vOmJhcg==",
	})
	assert.Contains(t, pairs, authPair{
		User:  "admin",
		Value: "Basic YWRtaW46cGFzc3dvcmQ=",
	})
}

func TestBasicAuthFails(t *testing.T) {
	assert.Panics(t, func() { processAccounts(nil) })
	assert.Panics(t, func() {
		processAccounts(Accounts{
			"":    "password",
			"foo": "bar",
		})
	})
}

func TestBasicAuthSearchCredential(t *testing.T) {
	pairs := processAccounts(Accounts{
		"admin": "password",
		"foo":   "bar",
		"bar":   "foo",
	})

	user, found := pairs.searchCredential(authorizationHeader("admin", "password"))
	assert.Equal(t, "admin", user)
	assert.True(t, found)

	user, found = pairs.searchCredential(authorizationHeader("foo", "bar"))
	assert.Equal(t, "foo", user)
	assert.True(t, found)

	user, found = pairs.searchCredential(authorizationHeader("bar", "foo"))
	assert.Equal(t, "bar", user)
	assert.True(t, found)

	user, found = pairs.searchCredential(authorizationHeader("admins", "password"))
	assert.Empty(t, user)
	assert.False(t, found)

	user, found = pairs.searchCredential(authorizationHeader("foo", "bar "))
	assert.Empty(t, user)
	assert.False(t, found)

	user, found = pairs.searchCredential("")
	assert.Empty(t, user)
	assert.False(t, found)
}

func TestBasicAuthAuthorizationHeader(t *testing.T) {
	assert.Equal(t, "Basic YWRtaW46cGFzc3dvcmQ=", authorizationHeader("admin", "password"))
}

func TestBasicAuthSecureCompare(t *testing.T) {
	assert.True(t, secureCompare("1234567890", "1234567890"))
	assert.False(t, secureCompare("123456789", "1234567890"))
	assert.False(t, secureCompare("12345678900", "1234567890"))
	assert.False(t, secureCompare("1234567891", "1234567890"))
}

func TestBasicAuthSucceed(t *testing.T) {
	accounts := Accounts{"admin": "password"}
	router := New()
	router.Use(BasicAuth(accounts))
	router.GET("/login", func(c *Context) {
		c.String(200, c.MustGet(AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", authorizationHeader("admin", "password"))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "admin", w.Body.String())
}

func TestBasicAuth401(t *testing.T) {
	called := false
	accounts := Accounts{"foo": "bar"}
	router := New()
	router.Use(BasicAuth(accounts))
	router.GET("/login", func(c *Context) {
		called = true
		c.String(200, c.MustGet(AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "Basic realm=\"Authorization Required\"", w.HeaderMap.Get("WWW-Authenticate"))
}

func TestBasicAuth401WithCustomRealm(t *testing.T) {
	called := false
	accounts := Accounts{"foo": "bar"}
	router := New()
	router.Use(BasicAuthForRealm(accounts, "My Custom \"Realm\""))
	router.GET("/login", func(c *Context) {
		called = true
		c.String(200, c.MustGet(AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "Basic realm=\"My Custom \\\"Realm\\\"\"", w.HeaderMap.Get("WWW-Authenticate"))
}
