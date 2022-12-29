// Copyright 2022 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProtoJSON contains the given interface object.
type ProtoJSON struct {
	Data any
}

// Render (ProtoJSON) marshals the given interface object and
// writes data with custom ContentType.
func (r ProtoJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := protojson.Marshal(r.Data.(protoreflect.ProtoMessage))
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// WriteContentType (ProtoBuf) writes ProtoBuf ContentType.
func (r ProtoJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
