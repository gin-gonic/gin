// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj any) error {
	return decodeXML(req.Body, obj)
}

func (xmlBinding) BindNoValidate(req *http.Request, obj any) error {
	return decodeXMLNoValidate(req.Body, obj)
}

func (xmlBinding) BindBody(body []byte, obj any) error {
	return decodeXML(bytes.NewReader(body), obj)
}

func (xmlBinding) BindBodyNoValidate(body []byte, obj any) error {
	return decodeXMLNoValidate(bytes.NewReader(body), obj)
}

func decodeXML(r io.Reader, obj any) error {
	if err := decodeXMLNoValidate(r, obj); err != nil {
		return err
	}
	return validate(obj)
}

func decodeXMLNoValidate(r io.Reader, obj any) error {
	decoder := xml.NewDecoder(r)
	return decoder.Decode(obj)
}
