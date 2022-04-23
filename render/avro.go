// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"github.com/hamba/avro"
)

// AVRO contains the given interface object.
type AVRO struct {
	Data   interface{}
	schema avro.Schema
}

var avroContentType = []string{"application/x-avro; charset=utf-8"}

// Render (AVRO) marshals the given interface object and writes data with custom ContentType.
func (r AVRO) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := avro.Marshal(r.schema, r.Data)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// WriteContentType (AVRO) writes AVRO ContentType for response.
func (r AVRO) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, avroContentType)
}
