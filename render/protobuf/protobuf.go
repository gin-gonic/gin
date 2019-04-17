// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protobuf

import (
	"net/http"

	"github.com/gin-gonic/gin/render/common"
	"github.com/golang/protobuf/proto"
)

func init() {
	common.List["ProtoBuf"] = NewProtoBuf
}

// ProtoBuf contains the given interface object.
type ProtoBuf struct {
	Data interface{}
}

var protobufContentType = []string{"application/x-protobuf"}

// Render (ProtoBuf) marshals the given interface object and writes data with custom ContentType.
func (r ProtoBuf) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := proto.Marshal(r.Data.(proto.Message))
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// WriteContentType (ProtoBuf) writes ProtoBuf ContentType.
func (r ProtoBuf) WriteContentType(w http.ResponseWriter) {
	common.WriteContentType(w, protobufContentType)
}

//NewProtoBuf build a new ProtoBuf render
func NewProtoBuf(obj interface{}, _ map[string]string) common.Render {
	return ProtoBuf{Data: obj}
}
