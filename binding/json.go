// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/codec/json"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// any as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	return decodeJSON(req.Body, obj, false)
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewReader(body), obj, false)
}

// strictJSONBinding is a JSON binding which always rejects unknown fields,
// regardless of the value of EnableDecoderDisallowUnknownFields.
type strictJSONBinding struct{}

func (strictJSONBinding) Name() string {
	return "strict-json"
}

func (strictJSONBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	return decodeJSON(req.Body, obj, true)
}

func (strictJSONBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewReader(body), obj, true)
}

func decodeJSON(r io.Reader, obj any, strict bool) error {
	decoder := json.API.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if strict || EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
