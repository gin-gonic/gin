// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/internal/bytesconv"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"

// AuthProxyUserKey is the cookie name for proxy_user credential in basic auth for proxy.
const AuthProxyUserKey = "proxy_user"

// Accounts defines a key/value for user/pass list of authorized logins.
type Accounts map[string]string

type authPair struct {
	value string
	user  string
}

type authPairs []authPair

func (a authPairs) searchCredential(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if subtle.ConstantTimeCompare(bytesconv.StringToBytes(pair.value), bytesconv.StringToBytes(authValue)) == 1 {
			return pair.user, true
		}
	}
	return "", false
}

// BasicAuthWithRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the username and the value is the password, as well as the name of the Realm.
// If the realm is empty, "Authorization Required" will be used by default.
// In fact, 'realm' should contain at least the name of the host performing the authentication and might additionally
// indicate the collection of users who might have access. An example might be "registered_users@go.dev".
// (see http://tools.ietf.org/html/rfc2617#section-1.2 for more details)
func BasicAuthWithRealm(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		// Search user in the slice of allowed credentials
		user, found := pairs.searchCredential(c.requestHeader("Authorization"))
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(AuthUserKey, user)
	}
}

// BasicAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where
// the key is the username and the value is the password.
func BasicAuth(accounts Accounts) HandlerFunc {
	return BasicAuthWithRealm(accounts, "")
}

func processAccounts(accounts Accounts) authPairs {
	length := len(accounts)
	assert1(length > 0, "Empty list of authorized credentials")
	pairs := make(authPairs, 0, length)
	for user, password := range accounts {
		assert1(user != "", "User can not be empty")
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			value: value,
			user:  user,
		})
	}
	return pairs
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(bytesconv.StringToBytes(base))
}

// BasicAuthForProxyWithRealm returns a Basic HTTP Proxy-Authorization middleware.
// If the realm is empty, "Proxy Authorization Required" will be used by default.
func BasicAuthForProxyWithRealm(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Proxy Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		proxyUser, found := pairs.searchCredential(c.requestHeader("Proxy-Authorization"))
		if !found {
			// Credentials doesn't match, we return 407 and abort handlers chain.
			c.Header("Proxy-Authenticate", realm)
			c.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		}
		// The proxy_user credentials was found, set proxy_user's id to key AuthProxyUserKey in this context, the proxy_user's id can be read later using
		// c.MustGet(gin.AuthProxyUserKey).
		c.Set(AuthProxyUserKey, proxyUser)
	}
}

// BasicAuthForProxy returns a Basic HTTP Proxy-Authorization middleware.
func BasicAuthForProxy(accounts Accounts) HandlerFunc {
	return BasicAuthForProxyWithRealm(accounts, "")
}
