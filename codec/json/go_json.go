// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build go_json

package json

import (
	"io"

	"github.com/goccy/go-json"
)

// Package indicates what library is being used for JSON encoding.
const Package = "github.com/goccy/go-json"

func init() {
	API = gojsonApi{}
}

type gojsonApi struct{}

func (j gojsonApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j gojsonApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j gojsonApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j gojsonApi) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (j gojsonApi) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}
