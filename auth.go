// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
)

const (
	AuthUserKey = "user"
)

type (
	Accounts map[string]string
	authPair struct {
		Value string
		User  string
	}
	authPairs []authPair
)

func (a authPairs) Len() int           { return len(a) }
func (a authPairs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a authPairs) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (a authPairs) searchCredential(auth string) (string, bool) {
	if len(auth) == 0 {
		return "", false
	}
	// Search user in the slice of allowed credentials
	r := sort.Search(len(a), func(i int) bool { return a[i].Value >= auth })
	if r < len(a) && secureCompare(a[r].Value, auth) {
		return a[r].User, true
	} else {
		return "", false
	}
}

// Implements a basic Basic HTTP Authorization. It takes as arguments a map[string]string where
// the key is the user name and the value is the password, as well as the name of the Realm
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Authorization Required"
	}
	realm = fmt.Sprintf("Basic realm=\"%s\"", realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		// Search user in the slice of allowed credentials
		user, ok := pairs.searchCredential(c.Request.Header.Get("Authorization"))
		if !ok {
			// Credentials doesn't match, we return 401 Unauthorized and abort request.
			c.Writer.Header().Set("WWW-Authenticate", realm)
			c.Fail(401, errors.New("Unauthorized"))
		} else {
			// user is allowed, set UserId to key "user" in this context, the userId can be read later using
			// c.Get(gin.AuthUserKey)
			c.Set(AuthUserKey, user)
		}
	}
}

// Implements a basic Basic HTTP Authorization. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
func BasicAuth(accounts Accounts) HandlerFunc {
	return BasicAuthForRealm(accounts, "")
}

func processAccounts(accounts Accounts) authPairs {
	if len(accounts) == 0 {
		panic("Empty list of authorized credentials")
	}
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts {
		if len(user) == 0 {
			panic("User can not be empty")
		}
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			Value: value,
			User:  user,
		})
	}
	// We have to sort the credentials in order to use bsearch later.
	sort.Sort(pairs)
	return pairs
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}

func secureCompare(given, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		/* Securely compare actual to itself to keep constant time, but always return false */
		return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
	}
}
