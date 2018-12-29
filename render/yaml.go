// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

func init() {
	Register(YAMLRenderType, YAMLFactory{})
}

// YAML contains the given interface object.
type YAML struct {
	Data interface{}
}

// YAMLFactory instance the YAML object.
type YAMLFactory struct{}

var yamlContentType = []string{"application/x-yaml; charset=utf-8"}

// Setup set data and opts
func (r *YAML) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *YAML) Reset() {
	r.Data = nil
}

// Render (YAML) marshals the given interface object and writes data with custom ContentType.
func (r *YAML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}

	w.Write(bytes)
	return nil
}

// WriteContentType (YAML) writes YAML ContentType for response.
func (r *YAML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, yamlContentType)
}

// Instance a new Render instance
func (YAMLFactory) Instance() RenderRecycler {
	return &YAML{}
}
