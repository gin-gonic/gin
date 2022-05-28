// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin/internal/bytesconv"
	"github.com/gin-gonic/gin/internal/json"
)

// JSON contains the given interface object.
type JSON struct {
	Data any
}

// IndentedJSON contains the given interface object.
type IndentedJSON struct {
	Data any
}

// SecureJSON contains the given interface object and its prefix.
type SecureJSON struct {
	Prefix string
	Data   any
}

// JsonpJSON contains the given interface object its callback.
type JsonpJSON struct {
	Callback string
	Data     any
}

// AsciiJSON contains the given interface object.
type AsciiJSON struct {
	Data any
}

// PureJSON contains the given interface object.
type PureJSON struct {
	Data any
}

var (
	jsonContentType      = []string{"application/json; charset=utf-8"}
	jsonpContentType     = []string{"application/javascript; charset=utf-8"}
	jsonASCIIContentType = []string{"application/json"}
)

// Render (JSON) writes data with custom ContentType.
func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj any) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

// Render (IndentedJSON) marshals the given interface object and writes it with custom ContentType.
func (r IndentedJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

// WriteContentType (IndentedJSON) writes JSON ContentType.
func (r IndentedJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Render (SecureJSON) marshals the given interface object and writes it with custom ContentType.
func (r SecureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	// if the jsonBytes is array values
	if bytes.HasPrefix(jsonBytes, bytesconv.StringToBytes("[")) && bytes.HasSuffix(jsonBytes,
		bytesconv.StringToBytes("]")) {
		if _, err = w.Write(bytesconv.StringToBytes(r.Prefix)); err != nil {
			return err
		}
	}
	_, err = w.Write(jsonBytes)
	return err
}

// WriteContentType (SecureJSON) writes JSON ContentType.
func (r SecureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Render (JsonpJSON) marshals the given interface object and writes it and its callback with custom ContentType.
func (r JsonpJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		_, err = w.Write(ret)
		return err
	}

	callback := template.JSEscapeString(r.Callback)
	if _, err = w.Write(bytesconv.StringToBytes(callback)); err != nil {
		return err
	}

	if _, err = w.Write(bytesconv.StringToBytes("(")); err != nil {
		return err
	}

	if _, err = w.Write(ret); err != nil {
		return err
	}

	if _, err = w.Write(bytesconv.StringToBytes(");")); err != nil {
		return err
	}

	return nil
}

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (r JsonpJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}

// Render (AsciiJSON) marshals the given interface object and writes it with custom ContentType.
func (r AsciiJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, r := range bytesconv.BytesToString(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	_, err = w.Write(buffer.Bytes())
	return err
}

// WriteContentType (AsciiJSON) writes JSON ContentType.
func (r AsciiJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonASCIIContentType)
}

// Render (PureJSON) writes custom ContentType and encodes the given interface object.
func (r PureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(r.Data)
}

// WriteContentType (PureJSON) writes custom ContentType.
func (r PureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
