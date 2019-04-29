// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package msgpack

import (
	"net/http"

	"github.com/gin-gonic/gin/render/common"
	"github.com/ugorji/go/codec"
)

func init() {
	common.List["MsgPack"] = NewMsgPack
}

// MsgPack contains the given interface object.
type MsgPack struct {
	Data interface{}
}

var msgpackContentType = []string{"application/msgpack; charset=utf-8"}

// WriteContentType (MsgPack) writes MsgPack ContentType.
func (r MsgPack) WriteContentType(w http.ResponseWriter) {
	common.WriteContentType(w, msgpackContentType)
}

// Render (MsgPack) encodes the given interface object and writes data with custom ContentType.
func (r MsgPack) Render(w http.ResponseWriter) error {
	return WriteMsgPack(w, r.Data)
}

// WriteMsgPack writes MsgPack ContentType and encodes the given interface object.
func WriteMsgPack(w http.ResponseWriter, obj interface{}) error {
	common.WriteContentType(w, msgpackContentType)
	var mh codec.MsgpackHandle
	return codec.NewEncoder(w, &mh).Encode(obj)
}

//NewMsgPack build a new MsgPack render
func NewMsgPack(obj interface{}, _ map[string]string) common.Render {
	return MsgPack{Data: obj}
}
