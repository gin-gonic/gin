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

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}
