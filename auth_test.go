// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
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
		user:  "bar",
		value: "Basic YmFyOmZvbw==",
	})
	assert.Contains(t, pairs, authPair{
		user:  "foo",
		value: "Basic Zm9vOmJhcg==",
	})
	assert.Contains(t, pairs, authPair{
		user:  "admin",
		value: "Basic YWRtaW46cGFzc3dvcmQ=",
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

func TestBasicAuthSchemeCase(t *testing.T) {
	accounts := Accounts{"admin": "password"}
	credentials := strings.TrimPrefix(authorizationHeader("admin", "password"), "Basic ")
	tests := []struct {
		name   string
		header string
		valid  bool
	}{
		{name: "canonical", header: "Basic " + credentials, valid: true},
		{name: "lowercase", header: "basic " + credentials, valid: true},
		{name: "uppercase", header: "BASIC " + credentials, valid: true},
		{name: "mixed case", header: "bAsIc " + credentials, valid: true},
		{name: "empty"},
		{name: "missing credentials", header: "Basic "},
		{name: "missing space", header: "Basic" + credentials},
		{name: "tab separator", header: "Basic\t" + credentials},
		{name: "other scheme", header: "Bearer " + credentials},
		{name: "longer scheme", header: "BasicX " + credentials},
		{name: "non-ASCII scheme", header: "ba\u017fic " + credentials},
		{name: "credential case", header: "basic " + strings.ToLower(credentials)},
		{name: "username case", header: authorizationHeader("Admin", "password")},
		{name: "password case", header: authorizationHeader("admin", "Password")},
		{name: "malformed credentials", header: "basic !invalid!"},
	}
	middlewares := []struct {
		name      string
		handler   HandlerFunc
		header    string
		userKey   string
		challenge string
		realm     string
		status    int
	}{
		{
			name: "origin", handler: BasicAuth(accounts),
			header: "Authorization", userKey: AuthUserKey,
			challenge: "WWW-Authenticate", realm: `Basic realm="Authorization Required"`,
			status: http.StatusUnauthorized,
		},
		{
			name: "proxy", handler: BasicAuthForProxy(accounts, ""),
			header: "Proxy-Authorization", userKey: AuthProxyUserKey,
			challenge: "Proxy-Authenticate", realm: `Basic realm="Proxy Authorization Required"`,
			status: http.StatusProxyAuthRequired,
		},
	}
	for _, middleware := range middlewares {
		t.Run(middleware.name, func(t *testing.T) {
			router := New()
			router.Use(middleware.handler)
			router.GET("/", func(c *Context) {
				c.String(http.StatusOK, c.MustGet(middleware.userKey).(string))
			})
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					w := httptest.NewRecorder()
					req := httptest.NewRequest(http.MethodGet, "/", nil)
					req.Header.Set(middleware.header, tt.header)
					router.ServeHTTP(w, req)
					if tt.valid {
						assert.Equal(t, http.StatusOK, w.Code)
						assert.Equal(t, "admin", w.Body.String())
						assert.Empty(t, w.Header().Get(middleware.challenge))
					} else {
						assert.Equal(t, middleware.status, w.Code)
						assert.Empty(t, w.Body.String())
						assert.Equal(t, middleware.realm, w.Header().Get(middleware.challenge))
					}
				})
			}
		})
	}
}

func TestBasicAuthSucceed(t *testing.T) {
	accounts := Accounts{"admin": "password"}
	router := New()
	router.Use(BasicAuth(accounts))
	router.GET("/login", func(c *Context) {
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", authorizationHeader("admin", "password"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin", w.Body.String())
}

func TestBasicAuth401(t *testing.T) {
	called := false
	accounts := Accounts{"foo": "bar"}
	router := New()
	router.Use(BasicAuth(accounts))
	router.GET("/login", func(c *Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Basic realm=\"Authorization Required\"", w.Header().Get("WWW-Authenticate"))
}

func TestBasicAuth401WithCustomRealm(t *testing.T) {
	called := false
	accounts := Accounts{"foo": "bar"}
	router := New()
	router.Use(BasicAuthForRealm(accounts, "My Custom \"Realm\""))
	router.GET("/login", func(c *Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Basic realm=\"My Custom \\\"Realm\\\"\"", w.Header().Get("WWW-Authenticate"))
}

func TestBasicAuthForProxySucceed(t *testing.T) {
	accounts := Accounts{"admin": "password"}
	router := New()
	router.Use(BasicAuthForProxy(accounts, ""))
	router.Any("/*proxyPath", func(c *Context) {
		c.String(http.StatusOK, c.MustGet(AuthProxyUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Proxy-Authorization", authorizationHeader("admin", "password"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin", w.Body.String())
}

func TestBasicAuthForProxy407(t *testing.T) {
	called := false
	accounts := Accounts{"foo": "bar"}
	router := New()
	router.Use(BasicAuthForProxy(accounts, ""))
	router.Any("/*proxyPath", func(c *Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthProxyUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusProxyAuthRequired, w.Code)
	assert.Equal(t, "Basic realm=\"Proxy Authorization Required\"", w.Header().Get("Proxy-Authenticate"))
}
