// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BSON contains the given interface object.
type BSON struct {
	Data any
}

var bsonContentType = []string{"application/bson"}

// Render (BSON) marshals the given interface object and writes data with custom ContentType.
func (r BSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := bson.Marshal(&r.Data)
	if err == nil {
		_, err = w.Write(bytes)
	}
	return err
}

// WriteContentType (BSONBuf) writes BSONBuf ContentType.
func (r BSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, bsonContentType)
}
