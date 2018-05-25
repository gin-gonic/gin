// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin/json"
)

type JSON struct {
	Data interface{}
}

type IndentedJSON struct {
	Data interface{}
}

type SecureJSON struct {
	Prefix string
	Data   interface{}
}

type JsonpJSON struct {
	Callback string
	Data     interface{}
}

type SecureJSONPrefix string

var jsonContentType = []string{"application/json; charset=utf-8"}
var jsonpContentType = []string{"application/javascript; charset=utf-8"}

func (r JSON) Render(w http.ResponseWriter) (err error) {
	return WriteJSON(w, r.Data)
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func (r IndentedJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func (r IndentedJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r SecureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	// if the jsonBytes is array values
	if bytes.HasPrefix(jsonBytes, []byte("[")) && bytes.HasSuffix(jsonBytes, []byte("]")) {
		jsonBytes = append([]byte(r.Prefix), jsonBytes...)
	}
	_, err = w.Write(jsonBytes)
	return err
}

func (r SecureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r JsonpJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		_, err := w.Write(ret)
		return err
	}

	callback := template.JSEscapeString(r.Callback)

	tmp := []byte("")
	tmp = append(tmp, []byte(callback)...)
	tmp = append(tmp, []byte("(")...)
	tmp = append(tmp, ret...)
	tmp = append(tmp, []byte(")")...)

	_, err = w.Write(tmp)
	return err
}

func (r JsonpJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}
