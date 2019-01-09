// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"github.com/golang/protobuf/proto"
)

func init() {
	Register(ProtoBufRenderType, ProtoBufFactory{})
}

// ProtoBuf contains the given interface object.
type ProtoBuf struct {
	Data interface{}
}

// ProtoBufFactory instance the ProtoBuf object.
type ProtoBufFactory struct{}

var protobufContentType = []string{"application/x-protobuf"}

// Setup set data and opts
func (r *ProtoBuf) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *ProtoBuf) Reset() {
	r.Data = nil
}

// Render (ProtoBuf) marshals the given interface object and writes data with custom ContentType.
func (r *ProtoBuf) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := proto.Marshal(r.Data.(proto.Message))
	if err != nil {
		return err
	}

	w.Write(bytes)
	return nil
}

// WriteContentType (ProtoBuf) writes ProtoBuf ContentType.
func (r *ProtoBuf) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, protobufContentType)
}

// Instance a new ProtoBuf object.
func (ProtoBufFactory) Instance() RenderRecycler {
	return &ProtoBuf{}
}
