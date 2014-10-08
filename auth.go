// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
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

// Implements a basic Basic HTTP Authorization. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
func BasicAuth(accounts Accounts) HandlerFunc {
	pairs, err := processAccounts(accounts)
	if err != nil {
		panic(err)
	}
	return func(c *Context) {
		// Search user in the slice of allowed credentials
		user, ok := searchCredential(pairs, c.Request.Header.Get("Authorization"))
		if !ok {
			// Credentials doesn't match, we return 401 Unauthorized and abort request.
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.Fail(401, errors.New("Unauthorized"))
		} else {
			// user is allowed, set UserId to key "user" in this context, the userId can be read later using
			// c.Get(gin.AuthUserKey)
			c.Set(AuthUserKey, user)
		}
	}
}

func processAccounts(accounts Accounts) (authPairs, error) {
	if len(accounts) == 0 {
		return nil, errors.New("Empty list of authorized credentials")
	}
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts {
		if len(user) == 0 {
			return nil, errors.New("User can not be empty")
		}
		base := user + ":" + password
		value := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
		pairs = append(pairs, authPair{
			Value: value,
			User:  user,
		})
	}
	// We have to sort the credentials in order to use bsearch later.
	sort.Sort(pairs)
	return pairs, nil
}

func searchCredential(pairs authPairs, auth string) (string, bool) {
	if len(auth) == 0 {
		return "", false
	}
	// Search user in the slice of allowed credentials
	r := sort.Search(len(pairs), func(i int) bool { return pairs[i].Value >= auth })
	if r < len(pairs) && secureCompare(pairs[r].Value, auth) {
		return pairs[r].User, true
	} else {
		return "", false
	}
}

func secureCompare(given, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		/* Securely compare actual to itself to keep constant time, but always return false */
		return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
	}
}
