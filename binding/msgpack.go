// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack
// +build !nomsgpack

package binding

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ugorji/go/codec"
)

type msgpackBinding struct{}

func (msgpackBinding) Name() string {
	return "msgpack"
}

func (b msgpackBinding) Bind(req *http.Request, obj interface{}) error {
	if err := b.BindOnly(req, obj); err != nil {
		return err
	}

	return validate(obj)
}

func (b msgpackBinding) BindOnly(req *http.Request, obj interface{}) error {
	return decodeMsgPack(req.Body, obj)
}

func (b msgpackBinding) BindBody(body []byte, obj interface{}) error {
	if err := b.BindBodyOnly(body, obj); err != nil {
		return err
	}

	return validate(obj)
}

func (b msgpackBinding) BindBodyOnly(body []byte, obj interface{}) error {
	return decodeMsgPack(bytes.NewReader(body), obj)
}

func decodeMsgPack(r io.Reader, obj interface{}) error {
	cdc := new(codec.MsgpackHandle)
	return codec.NewDecoder(r, cdc).Decode(&obj)
}
