// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type protoJSONBinding struct{}

func (protoJSONBinding) Name() string {
	return "protojson"
}

func (b protoJSONBinding) Bind(req *http.Request, obj any) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (protoJSONBinding) BindBody(body []byte, obj any) error {
	msg, ok := obj.(protoreflect.ProtoMessage)
	if !ok {
		return errors.New("obj is not ProtoMessage")
	}
	if err := protojson.Unmarshal(body, msg); err != nil {
		return err
	}
	// Here it's same to return validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
}
