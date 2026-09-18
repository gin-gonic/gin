// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package bson provides optional BSON rendering and binding for gin.
//
// BSON support is no longer compiled into the core gin module. Import this
// package to opt back in:
//
//	import (
//		"github.com/gin-gonic/gin"
//		"github.com/gin-gonic/gin/render/bson"
//	)
//
//	bson.Render(c, http.StatusOK, obj) // write a BSON response
//	bson.ShouldBind(c, &obj)           // decode a BSON request body
//
// Importing the package registers the binding and renderer for the
// "application/bson" content type so that the content-type negotiation done by
// c.ShouldBind and c.Negotiate keeps working.
package bson

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MIMEBSON is the content type handled by this package.
const MIMEBSON = binding.MIMEBSON

var contentType = []string{"application/bson"}

func init() {
	binding.Register(MIMEBSON, Binding)
	render.Register(MIMEBSON, func(data any) render.Render { return renderer{Data: data} })
}

// renderer implements render.Render for BSON responses.
type renderer struct {
	Data any
}

var _ render.Render = renderer{}

// Render marshals the data as BSON and writes it with the BSON content type.
func (r renderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	bytes, err := bson.Marshal(&r.Data)
	if err == nil {
		_, err = w.Write(bytes)
	}
	return err
}

// WriteContentType writes the BSON content type.
func (r renderer) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, contentType)
}

// Render writes obj to the response as BSON with status code. It is the
// drop-in replacement for the former c.BSON(code, obj).
func Render(c *gin.Context, code int, obj any) {
	c.Render(code, renderer{Data: obj})
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Binding decodes BSON request bodies. It can be passed to
// c.ShouldBindWith / c.MustBindWith.
var Binding binding.BindingBody = bsonBinding{}

type bsonBinding struct{}

func (bsonBinding) Name() string {
	return "bson"
}

func (b bsonBinding) Bind(req *http.Request, obj any) error {
	buf, err := io.ReadAll(req.Body)
	if err == nil {
		err = b.BindBody(buf, obj)
	}
	return err
}

func (bsonBinding) BindBody(body []byte, obj any) error {
	return bson.Unmarshal(body, obj)
}

// Bind binds the BSON request body to obj, aborting with HTTP 400 on error.
func Bind(c *gin.Context, obj any) error {
	return c.MustBindWith(obj, Binding)
}

// ShouldBind binds the BSON request body to obj without aborting.
func ShouldBind(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, Binding)
}
