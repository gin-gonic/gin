// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package msgpack

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding/common"
	"github.com/ugorji/go/codec"
)

func init() {
	msgPack := msgpackBinding{}
	common.List[common.MIMEMSGPACK] = msgPack
	common.List[common.MIMEMSGPACK2] = msgPack
}

type msgpackBinding struct{}

func (msgpackBinding) Name() string {
	return "msgpack"
}

func (msgpackBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeMsgPack(req.Body, obj)
}

func (msgpackBinding) BindBody(body []byte, obj interface{}) error {
	return decodeMsgPack(bytes.NewReader(body), obj)
}

func decodeMsgPack(r io.Reader, obj interface{}) error {
	cdc := new(codec.MsgpackHandle)
	if err := codec.NewDecoder(r, cdc).Decode(&obj); err != nil {
		return err
	}
	return common.Validate(obj)
}
