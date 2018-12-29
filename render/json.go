// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin/internal/json"
)

func init() {
	Register(JSONRenderType, JSONFactory{})
	Register(IntendedJSONRenderType, IndentedJSONFactory{})
	Register(JsonpJSONRenderType, JsonpJSONFactory{})
	Register(SecureJSONRenderType, SecureJSONFactory{})
	Register(AsciiJSONRenderType, AsciiJSONFactory{})
}

// JSON contains the given interface object.
type JSON struct {
	Data interface{}
}

// IndentedJSON contains the given interface object.
type IndentedJSON struct {
	Data interface{}
}

// SecureJSON contains the given interface object and its prefix.
type SecureJSON struct {
	Prefix string
	Data   interface{}
}

// JsonpJSON contains the given interface object its callback.
type JsonpJSON struct {
	Callback string
	Data     interface{}
}

// AsciiJSON contains the given interface object.
type AsciiJSON struct {
	Data interface{}
}

// JSONFactory instance the JSON object.
type JSONFactory struct{}

// IndentedJSONFactory instance the IndentedJSON object.
type IndentedJSONFactory struct{}

// SecureJSONFactory instance the SecureJSON object.
type SecureJSONFactory struct{}

// JsonpJSONFactory instance the JsonpJSON object.
type JsonpJSONFactory struct{}

// AsciiJSONFactory instance the AsciiJSON object.
type AsciiJSONFactory struct{}

// SecureJSONPrefix is a string which represents SecureJSON prefix.
type SecureJSONPrefix string

var jsonContentType = []string{"application/json; charset=utf-8"}
var jsonpContentType = []string{"application/javascript; charset=utf-8"}
var jsonAsciiContentType = []string{"application/json"}

// Setup set data and opts
func (r *JSON) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *JSON) Reset() {
	r.Data = nil
}

// Render (JSON) writes data with custom ContentType.
func (r *JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r *JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

// Setup set data and opts
func (r *IndentedJSON) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *IndentedJSON) Reset() {
	r.Data = nil
}

// Render (IndentedJSON) marshals the given interface object and writes it with custom ContentType.
func (r *IndentedJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

// WriteContentType (IndentedJSON) writes JSON ContentType.
func (r *IndentedJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Setup set data and opts
func (r *SecureJSON) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
	if len(opts) == 1 {
		r.Prefix, _ = opts[0].(string)
	}
}

// Reset clean data and opts
func (r *SecureJSON) Reset() {
	r.Data = nil
	r.Prefix = ""
}

// Render (SecureJSON) marshals the given interface object and writes it with custom ContentType.
func (r *SecureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	// if the jsonBytes is array values
	if bytes.HasPrefix(jsonBytes, []byte("[")) && bytes.HasSuffix(jsonBytes, []byte("]")) {
		w.Write([]byte(r.Prefix))
	}
	w.Write(jsonBytes)
	return nil
}

// WriteContentType (SecureJSON) writes JSON ContentType.
func (r *SecureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Setup set data and opts
func (r *JsonpJSON) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
	if len(opts) == 1 {
		if callback, ok := opts[0].(string); ok {
			r.Callback = callback
		} else {
			r.Callback = ""
		}
	}
}

// Reset clean data and opts
func (r *JsonpJSON) Reset() {
	r.Data = nil
	r.Callback = ""
}

// Render (JsonpJSON) marshals the given interface object and writes it and its callback with custom ContentType.
func (r *JsonpJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		w.Write(ret)
		return nil
	}

	callback := template.JSEscapeString(r.Callback)
	w.Write([]byte(callback))
	w.Write([]byte("("))
	w.Write(ret)
	w.Write([]byte(")"))

	return nil
}

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (r *JsonpJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}

// Setup set data and opts
func (r *AsciiJSON) Setup(data interface{}, opts ...interface{}) {
	r.Data = data
}

// Reset clean data and opts
func (r *AsciiJSON) Reset() {
	r.Data = nil
}

// Render (AsciiJSON) marshals the given interface object and writes it with custom ContentType.
func (r *AsciiJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, r := range string(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	w.Write(buffer.Bytes())
	return nil
}

// WriteContentType (AsciiJSON) writes JSON ContentType.
func (r *AsciiJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonAsciiContentType)
}

// Instance a new JSON object.
func (JSONFactory) Instance() RenderRecycler {
	return &JSON{}
}

// Instance a new IndentedJSON object.
func (IndentedJSONFactory) Instance() RenderRecycler {
	return &IndentedJSON{}
}

// Instance a new SecureJSON object.
func (SecureJSONFactory) Instance() RenderRecycler {
	return &SecureJSON{}
}

// Instance a new JsonpJSON object.
func (JsonpJSONFactory) Instance() RenderRecycler {
	return &JsonpJSON{}
}

// Instance a new AsciiJSON object.
func (AsciiJSONFactory) Instance() RenderRecycler {
	return &AsciiJSON{}
}
