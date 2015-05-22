// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/json"
	"net/http"
)

type (
	JSON struct {
		Data interface{}
	}

	IndentedJSON struct {
		Data interface{}
	}
)

const jsonContentType = "application/json; charset=utf-8"

func (r JSON) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", jsonContentType)
	return json.NewEncoder(w).Encode(r.Data)
}

func (r IndentedJSON) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", jsonContentType)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}
