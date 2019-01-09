// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/xml"
	"net/http"
)

func init() {
	Register(XMLRenderType, XMLFactory{})
}

// XML contains the given interface object.
type XML struct {
	Data interface{}
}

// XMLFactory instance the XML object.
type XMLFactory struct{}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Setup set data and opts
func (r *XML) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *XML) Reset() {
	r.Data = nil
}

// Render (XML) encodes the given interface object and writes data with custom ContentType.
func (r *XML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

// WriteContentType (XML) writes XML ContentType for response.
func (r *XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, xmlContentType)
}

// Instance a new XML object.
func (XMLFactory) Instance() RenderRecycler {
	return &XML{}
}
