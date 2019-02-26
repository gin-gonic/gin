// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"net/http"
)

// Redirect contains the http request reference and redirects status code and location.
type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

// Render (Redirect) redirects the http request to new location and writes redirect response.
func (r Redirect) Render(w http.ResponseWriter) error {
	// todo(thinkerou): go1.6 not support StatusPermanentRedirect(308)
	// when we upgrade go version we can use http.StatusPermanentRedirect
	if (r.Code < 300 || r.Code > 308) && r.Code != 201 {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

// WriteContentType (Redirect) don't write any ContentType.
func (r Redirect) WriteContentType(http.ResponseWriter) {}
