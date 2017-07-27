// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"io"
	"net/http"
)

type String struct {
	Format string
	Data   []interface{}
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func (r String) Render(w http.ResponseWriter) error {
	WriteString(w, r.Format, r.Data, false)
	return nil
}

func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}

func WriteString(w http.ResponseWriter, format string, data []interface{}, html bool) {
	if html {
		writeContentType(w, htmlContentType)
	} else {
		writeContentType(w, plainContentType)
	}
	if len(data) > 0 {
		fmt.Fprintf(w, format, data...)
	} else {
		io.WriteString(w, format)
	}
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
