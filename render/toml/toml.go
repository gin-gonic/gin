// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package toml provides optional TOML rendering and binding for gin.
//
// TOML support is no longer compiled into the core gin module. Import this
// package to opt back in:
//
//	import (
//		"github.com/gin-gonic/gin"
//		"github.com/gin-gonic/gin/render/toml"
//	)
//
//	toml.Render(c, http.StatusOK, obj) // write a TOML response
//	toml.ShouldBind(c, &obj)           // decode a TOML request body
//
// Importing the package registers the binding and renderer for the
// "application/toml" content type so that the content-type negotiation done by
// c.ShouldBind and c.Negotiate keeps working.
package toml

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/pelletier/go-toml/v2"
)

// MIMETOML is the content type handled by this package.
const MIMETOML = binding.MIMETOML

var contentType = []string{"application/toml; charset=utf-8"}

func init() {
	binding.Register(MIMETOML, Binding)
	render.Register(MIMETOML, func(data any) render.Render { return renderer{Data: data} })
}

// renderer implements render.Render for TOML responses.
type renderer struct {
	Data any
}

var _ render.Render = renderer{}

// Render marshals the data as TOML and writes it with the TOML content type.
func (r renderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	bytes, err := toml.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

// WriteContentType writes the TOML content type.
func (r renderer) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, contentType)
}

// Render writes obj to the response as TOML with status code. It is the
// drop-in replacement for the former c.TOML(code, obj).
func Render(c *gin.Context, code int, obj any) {
	c.Render(code, renderer{Data: obj})
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Binding decodes TOML request bodies. It can be passed to
// c.ShouldBindWith / c.MustBindWith.
var Binding binding.BindingBody = tomlBinding{}

type tomlBinding struct{}

func (tomlBinding) Name() string {
	return "toml"
}

func (tomlBinding) Bind(req *http.Request, obj any) error {
	return decodeToml(req.Body, obj)
}

func (tomlBinding) BindBody(body []byte, obj any) error {
	return decodeToml(bytes.NewReader(body), obj)
}

func decodeToml(r io.Reader, obj any) error {
	decoder := toml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return binding.Validate(obj)
}

// Bind binds the TOML request body to obj, aborting with HTTP 400 on error.
// It replaces the former c.BindTOML(obj).
func Bind(c *gin.Context, obj any) error {
	return c.MustBindWith(obj, Binding)
}

// ShouldBind binds the TOML request body to obj without aborting. It replaces
// the former c.ShouldBindTOML(obj).
func ShouldBind(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, Binding)
}
