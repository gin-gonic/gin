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
	BasicAuthPair struct {
		Code string
		User string
	}
	Accounts map[string]string
	Pairs    []BasicAuthPair
)

func (a Pairs) Len() int           { return len(a) }
func (a Pairs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Pairs) Less(i, j int) bool { return a[i].Code < a[j].Code }

func processCredentials(accounts Accounts) (Pairs, error) {
	if len(accounts) == 0 {
		return nil, errors.New("Empty list of authorized credentials.")
	}
	pairs := make(Pairs, 0, len(accounts))
	for user, password := range accounts {
		if len(user) == 0 || len(password) == 0 {
			return nil, errors.New("User or password is empty")
		}
		base := user + ":" + password
		code := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
		pairs = append(pairs, BasicAuthPair{code, user})
	}
	// We have to sort the credentials in order to use bsearch later.
	sort.Sort(pairs)
	return pairs, nil
}

func secureCompare(given, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		/* Securely compare actual to itself to keep constant time, but always return false */
		return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
	}
}

func searchCredential(pairs Pairs, auth string) string {
	if len(auth) == 0 {
		return ""
	}
	// Search user in the slice of allowed credentials
	r := sort.Search(len(pairs), func(i int) bool { return pairs[i].Code >= auth })
	if r < len(pairs) && secureCompare(pairs[r].Code, auth) {
		return pairs[r].User
	} else {
		return ""
	}
}

// Implements a basic Basic HTTP Authorization. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
func BasicAuth(accounts Accounts) HandlerFunc {

	pairs, err := processCredentials(accounts)
	if err != nil {
		panic(err)
	}
	return func(c *Context) {
		// Search user in the slice of allowed credentials
		user := searchCredential(pairs, c.Request.Header.Get("Authorization"))
		if len(user) == 0 {
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
