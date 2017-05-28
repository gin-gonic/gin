// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"

	"github.com/casbin/casbin"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer) HandlerFunc {
	return func(c *Context) {
		a := &BasicAuthorizer{enforcer: e}

		if !a.CheckPermission(c.Request) {
			a.RequirePermission(c.Writer)
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) bool {
	user := a.GetUserName(r)
	method := r.Method
	path := r.URL.Path
	return a.enforcer.Enforce(user, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(w http.ResponseWriter) {
		w.WriteHeader(403)
	w.Write([]byte("403 Forbidden\n"))
}
