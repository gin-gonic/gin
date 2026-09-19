// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack

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

func (msgpackBinding) Bind(req *http.Request, obj any) error {
	return decodeMsgPack(req.Body, obj)
}

func (msgpackBinding) BindNoValidate(req *http.Request, obj any) error {
	return decodeMsgPackNoValidate(req.Body, obj)
}

func (msgpackBinding) BindBody(body []byte, obj any) error {
	return decodeMsgPack(bytes.NewReader(body), obj)
}

func (msgpackBinding) BindBodyNoValidate(body []byte, obj any) error {
	return decodeMsgPackNoValidate(bytes.NewReader(body), obj)
}

func decodeMsgPack(r io.Reader, obj any) error {
	if err := decodeMsgPackNoValidate(r, obj); err != nil {
		return err
	}
	return validate(obj)
}

func decodeMsgPackNoValidate(r io.Reader, obj any) error {
	cdc := new(codec.MsgpackHandle)
	return codec.NewDecoder(r, cdc).Decode(&obj)
}
