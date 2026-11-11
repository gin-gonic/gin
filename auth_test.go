// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/sha256"
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
		user:   "bar",
		value:  "Basic YmFyOmZvbw==",
		digest: sha256.Sum256([]byte("Basic YmFyOmZvbw==")),
	})
	assert.Contains(t, pairs, authPair{
		user:   "foo",
		value:  "Basic Zm9vOmJhcg==",
		digest: sha256.Sum256([]byte("Basic Zm9vOmJhcg==")),
	})
	assert.Contains(t, pairs, authPair{
		user:   "admin",
		value:  "Basic YWRtaW46cGFzc3dvcmQ=",
		digest: sha256.Sum256([]byte("Basic YWRtaW46cGFzc3dvcmQ=")),
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

func TestBasicAuthSearchCredentialLengthAndCollisionGuard(t *testing.T) {
	pairs := processAccounts(Accounts{
		"a":        "",
		"admin":    "password",
		"longuser": "very-long-secret-password-value",
	})

	for _, tc := range []struct {
		name      string
		header    string
		wantUser  string
		wantFound bool
	}{
		{
			name:      "empty password account matches",
			header:    authorizationHeader("a", ""),
			wantUser:  "a",
			wantFound: true,
		},
		{
			name:      "long account matches",
			header:    authorizationHeader("longuser", "very-long-secret-password-value"),
			wantUser:  "longuser",
			wantFound: true,
		},
		{
			name:   "shorter than any credential",
			header: "Basic YQ==",
		},
		{
			name:   "prefix of valid credential",
			header: authorizationHeader("admin", "password")[:10],
		},
		{
			name:   "longer than any credential",
			header: authorizationHeader("longuser", "very-long-secret-password-value-extra"),
		},
		{
			name:   "same length wrong password",
			header: authorizationHeader("admin", "passworx"),
		},
		{
			name:   "non-basic scheme",
			header: "Bearer some-token",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			user, found := pairs.searchCredential(tc.header)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantFound, found)
		})
	}

	// Verify exact value check prevents false positives even on matching digests.
	candidate := authorizationHeader("admin", "wrong")
	forged := authPairs{
		{
			user:   "admin",
			value:  authorizationHeader("admin", "password"),
			digest: sha256.Sum256([]byte(candidate)),
		},
	}
	user, found := forged.searchCredential(candidate)
	assert.Empty(t, user)
	assert.False(t, found)
}

func TestBasicAuthAuthorizationHeader(t *testing.T) {
	assert.Equal(t, "Basic YWRtaW46cGFzc3dvcmQ=", authorizationHeader("admin", "password"))
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
