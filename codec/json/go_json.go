// Copyright 2017 Bo-Yi Wu. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build go_json

package json

import (
	"io"

	"github.com/gin-gonic/gin/codec/api"
	"github.com/goccy/go-json"
)

func init() {
	Api = gojsonApi{}
}

type gojsonApi struct {
}

func (j gojsonApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j gojsonApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j gojsonApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j gojsonApi) NewEncoder(writer io.Writer) api.JsonEncoder {
	return json.NewEncoder(writer)
}

func (j gojsonApi) NewDecoder(reader io.Reader) api.JsonDecoder {
	return json.NewDecoder(reader)
}
