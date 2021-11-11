// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack
// +build !nomsgpack

package binding

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/ugorji/go/codec"
)

type msgpackBinding struct{}

func (msgpackBinding) Name() string {
	return "msgpack"
}

func (b msgpackBinding) Bind(req *http.Request, obj interface{}) error {
	return b.BindContext(context.Background(), req, obj)
}

func (msgpackBinding) BindContext(ctx context.Context, req *http.Request, obj interface{}) error {
	return decodeMsgPack(ctx, req.Body, obj)
}

func (b msgpackBinding) BindBody(body []byte, obj interface{}) error {
	return b.BindBodyContext(context.Background(), body, obj)
}

func (msgpackBinding) BindBodyContext(ctx context.Context, body []byte, obj interface{}) error {
	return decodeMsgPack(ctx, bytes.NewReader(body), obj)
}

func decodeMsgPack(ctx context.Context, r io.Reader, obj interface{}) error {
	cdc := new(codec.MsgpackHandle)
	if err := codec.NewDecoder(r, cdc).Decode(&obj); err != nil {
		return err
	}
	return validateContext(ctx, obj)
}
