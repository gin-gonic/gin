// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
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

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	err := decodeXML(req.Body, obj)
	if err != nil {
		return err
	}

	return validate(obj)
}

func (xmlBinding) BindOnly(req *http.Request, obj interface{}) error {
	return decodeXML(req.Body, obj)
}

func (b xmlBinding) BindBody(body []byte, obj interface{}) error {
	if err := b.BindBodyOnly(body, obj); err != nil {
		return err
	}

	return validate(obj)
}

func (xmlBinding) BindBodyOnly(body []byte, obj interface{}) error {
	return decodeXML(bytes.NewReader(body), obj)
}

func decodeXML(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return nil
}
