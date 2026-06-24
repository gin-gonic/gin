// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package msgpack provides optional MessagePack rendering and binding for gin.
//
// MessagePack support is no longer compiled into the core gin module. Import
// this package (for its helpers, or with a blank identifier purely for the
// init-time registration) to opt back in:
//
//	import (
//		"github.com/gin-gonic/gin"
//		"github.com/gin-gonic/gin/render/msgpack"
//	)
//
//	msgpack.Render(c, http.StatusOK, obj) // write a MessagePack response
//	msgpack.ShouldBind(c, &obj)           // decode a MessagePack request body
//
// Importing the package also registers the binding and renderer for the
// "application/x-msgpack" and "application/msgpack" content types, so the
// content-type negotiation done by c.ShouldBind and c.Negotiate keeps working.
package msgpack

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/ugorji/go/codec"
)

// Content types handled by this package.
const (
	MIMEMSGPACK  = binding.MIMEMSGPACK
	MIMEMSGPACK2 = binding.MIMEMSGPACK2
)

var contentType = []string{"application/msgpack; charset=utf-8"}

func init() {
	binding.Register(MIMEMSGPACK, Binding)
	binding.Register(MIMEMSGPACK2, Binding)
	factory := func(data any) render.Render { return renderer{Data: data} }
	render.Register(MIMEMSGPACK, factory)
	render.Register(MIMEMSGPACK2, factory)
}

// renderer implements render.Render for MessagePack responses.
type renderer struct {
	Data any
}

var _ render.Render = renderer{}

// Render encodes the data as MessagePack and writes it with the MessagePack
// content type.
func (r renderer) Render(w http.ResponseWriter) error {
	return WriteMsgPack(w, r.Data)
}

// WriteContentType writes the MessagePack content type.
func (r renderer) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, contentType)
}

// Render writes obj to the response as MessagePack with status code. It is the
// drop-in replacement for the former c.MsgPack(code, obj).
func Render(c *gin.Context, code int, obj any) {
	c.Render(code, renderer{Data: obj})
}

// WriteMsgPack writes the MessagePack content type and encodes obj to w.
func WriteMsgPack(w http.ResponseWriter, obj any) error {
	writeContentType(w, contentType)
	var mh codec.MsgpackHandle
	return codec.NewEncoder(w, &mh).Encode(obj)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Binding decodes MessagePack request bodies. It can be passed to
// c.ShouldBindWith / c.MustBindWith.
var Binding binding.BindingBody = msgpackBinding{}

type msgpackBinding struct{}

func (msgpackBinding) Name() string {
	return "msgpack"
}

func (msgpackBinding) Bind(req *http.Request, obj any) error {
	return decodeMsgPack(req.Body, obj)
}

func (msgpackBinding) BindBody(body []byte, obj any) error {
	return decodeMsgPack(bytes.NewReader(body), obj)
}

func decodeMsgPack(r io.Reader, obj any) error {
	cdc := new(codec.MsgpackHandle)
	if err := codec.NewDecoder(r, cdc).Decode(&obj); err != nil {
		return err
	}
	return binding.Validate(obj)
}

// Bind binds the MessagePack request body to obj, aborting with HTTP 400 on
// error. It replaces the former c.BindWith(obj, binding.MsgPack).
func Bind(c *gin.Context, obj any) error {
	return c.MustBindWith(obj, Binding)
}

// ShouldBind binds the MessagePack request body to obj without aborting.
func ShouldBind(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, Binding)
}
