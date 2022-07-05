// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/internal/bytesconv"
)

// String contains the given interface object slice and its format.
type String struct {
	Format string
	Data   []any
}

var plainContentType = []string{"text/plain; charset=utf-8"}

// Render (String) writes data with custom ContentType.
func (r String) Render(w http.ResponseWriter) error {
	return WriteString(w, r.Format, r.Data, false)
}

// WriteContentType (String) writes Plain ContentType.
func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}

// WriteString writes data according to its format and write custom ContentType.
func WriteString(w http.ResponseWriter, format string, data []any, html bool) (err error) {
	if html {
		writeContentType(w, htmlContentType)
	} else {
		writeContentType(w, plainContentType)
	}
	if len(data) > 0 {
		_, err = fmt.Fprintf(w, format, data...)
		return
	}
	_, err = w.Write(bytesconv.StringToBytes(format))
	return
}

// StringHTML will function exactly the same as the String struct
// but it will inject an html ContentType to the response
type StringHTML struct {
	Format string
	Data   []interface{}
}

func (r StringHTML) Render(w http.ResponseWriter) error {
	WriteString(w, r.Format, r.Data, true)
	return nil
}

func (r StringHTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
