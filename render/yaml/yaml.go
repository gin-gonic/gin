// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package yaml provides optional YAML rendering and binding for gin.
//
// YAML support is no longer compiled into the core gin module. Import this
// package to opt back in:
//
//	import (
//		"github.com/gin-gonic/gin"
//		"github.com/gin-gonic/gin/render/yaml"
//	)
//
//	yaml.Render(c, http.StatusOK, obj) // write a YAML response
//	yaml.ShouldBind(c, &obj)           // decode a YAML request body
//
// Importing the package registers the binding and renderer for the
// "application/x-yaml" and "application/yaml" content types so that the
// content-type negotiation done by c.ShouldBind and c.Negotiate keeps working.
package yaml

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/goccy/go-yaml"
)

// Content types handled by this package.
const (
	MIMEYAML  = binding.MIMEYAML
	MIMEYAML2 = binding.MIMEYAML2
)

var contentType = []string{"application/yaml; charset=utf-8"}

func init() {
	binding.Register(MIMEYAML, Binding)
	binding.Register(MIMEYAML2, Binding)
	factory := func(data any) render.Render { return renderer{Data: data} }
	render.Register(MIMEYAML, factory)
	render.Register(MIMEYAML2, factory)
}

// renderer implements render.Render for YAML responses.
type renderer struct {
	Data any
}

var _ render.Render = renderer{}

// Render marshals the data as YAML and writes it with the YAML content type.
func (r renderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

// WriteContentType writes the YAML content type.
func (r renderer) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, contentType)
}

// Render writes obj to the response as YAML with status code. It is the
// drop-in replacement for the former c.YAML(code, obj).
func Render(c *gin.Context, code int, obj any) {
	c.Render(code, renderer{Data: obj})
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Binding decodes YAML request bodies. It can be passed to
// c.ShouldBindWith / c.MustBindWith.
var Binding binding.BindingBody = yamlBinding{}

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (yamlBinding) Bind(req *http.Request, obj any) error {
	return decodeYAML(req.Body, obj)
}

func (yamlBinding) BindBody(body []byte, obj any) error {
	return decodeYAML(bytes.NewReader(body), obj)
}

func decodeYAML(r io.Reader, obj any) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return binding.Validate(obj)
}

// Bind binds the YAML request body to obj, aborting with HTTP 400 on error.
// It replaces the former c.BindYAML(obj).
func Bind(c *gin.Context, obj any) error {
	return c.MustBindWith(obj, Binding)
}

// ShouldBind binds the YAML request body to obj without aborting. It replaces
// the former c.ShouldBindYAML(obj).
func ShouldBind(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, Binding)
}
