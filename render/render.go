// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

type Render interface {
	Render(http.ResponseWriter, int, ...interface{}) error
}

var (
	JSON         Render = jsonRender{}
	IndentedJSON Render = indentedJSON{}
	XML          Render = xmlRender{}
	HTMLPlain    Render = htmlPlainRender{}
	Plain        Render = plainTextRender{}
	Redirect     Render = redirectRender{}
	Data         Render = dataRender{}
	_            Render = HTMLRender{}
	_            Render = &HTMLDebugRender{}
)

func WriteHeader(w http.ResponseWriter, code int, contentType string) {
	contentType = joinStrings(contentType, "; charset=utf-8")
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}

func joinStrings(a ...string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := 0
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	n = 0
	for _, s := range a {
		n += copy(b[n:], s)
	}
	return string(b)
}
