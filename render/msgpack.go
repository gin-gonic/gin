// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"github.com/ugorji/go/codec"
)

func init() {
	Register(MsgPackRenderType, MsgPackFactory{})
}

// MsgPack contains the given interface object.
type MsgPack struct {
	Data interface{}
}

// MsgPackFactory instance the MsgPack object.
type MsgPackFactory struct{}

var msgpackContentType = []string{"application/msgpack; charset=utf-8"}

// Setup set data and opts
func (r *MsgPack) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *MsgPack) Reset() {
	r.Data = nil
}

// Render (MsgPack) encodes the given interface object and writes data with custom ContentType.
func (r *MsgPack) Render(w http.ResponseWriter) error {
	return WriteMsgPack(w, r.Data)
}

// WriteContentType (MsgPack) writes MsgPack ContentType.
func (r *MsgPack) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, msgpackContentType)
}

// WriteMsgPack writes MsgPack ContentType and encodes the given interface object.
func WriteMsgPack(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, msgpackContentType)
	var mh codec.MsgpackHandle
	return codec.NewEncoder(w, &mh).Encode(obj)
}

// Instance a new MsgPack object.
func (MsgPackFactory) Instance() RenderRecycler {
	return &MsgPack{}
}
