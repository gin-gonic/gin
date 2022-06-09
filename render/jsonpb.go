// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONPB contains the given proto.Message object.
type JSONPB struct {
	Data   proto.Message
	datapb json.RawMessage
	Option protojson.MarshalOptions
}

// Render (JSONPB) writes data with custom ContentType.
func (r JSONPB) Render(w http.ResponseWriter) (err error) {
	r.datapb, err = r.Option.Marshal(r.Data)
	if err != nil {
		return err
	}
	writeContentType(w, jsonContentType)
	_, err = w.Write(r.datapb)
	return err
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSONPB) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
