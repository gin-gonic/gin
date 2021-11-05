// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (b xmlBinding) Bind(req *http.Request, obj interface{}) error {
	return b.BindContext(context.Background(), req, obj)
}

func (xmlBinding) BindContext(ctx context.Context, req *http.Request, obj interface{}) error {
	return decodeXML(ctx, req.Body, obj)
}

func (b xmlBinding) BindBody(body []byte, obj interface{}) error {
	return b.BindBodyContext(context.Background(), body, obj)
}

func (xmlBinding) BindBodyContext(ctx context.Context, body []byte, obj interface{}) error {
	return decodeXML(ctx, bytes.NewReader(body), obj)
}

func decodeXML(ctx context.Context, r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validateContext(ctx, obj)
}
