// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build go1.7

package render

import (
	"net/http"

	"github.com/gin-gonic/gin/internal/json"
)

func init() {
	Register(PureJSONRenderType, &PureJsonFactory{})
}

// PureJSON contains the given interface object.
type PureJSON struct {
	Data interface{}
}

// Setup set data and opts
func (r *PureJSON) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *PureJSON) Reset() {
	r.Data = nil
}

// JSONFactory instance the PureJson object.
type PureJsonFactory struct{}

// Render (PureJSON) writes custom ContentType and encodes the given interface object.
func (r PureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(r.Data)
}

// WriteContentType (PureJSON) writes custom ContentType.
func (r PureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Instance a new Render instance
func (PureJsonFactory) Instance() RenderRecycler {
	return &PureJSON{}
}
