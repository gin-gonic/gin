// Copyright 2017 Bo-Yi Wu. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build jsoniter

package json

import (
	"io"

	"github.com/gin-gonic/gin/codec/api"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	Api = jsoniterApi{}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsoniterApi struct {
}

func (j jsoniterApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsoniterApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j jsoniterApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j jsoniterApi) NewEncoder(writer io.Writer) api.JsonEncoder {
	return json.NewEncoder(writer)
}

func (j jsoniterApi) NewDecoder(reader io.Reader) api.JsonDecoder {
	return json.NewDecoder(reader)
}
