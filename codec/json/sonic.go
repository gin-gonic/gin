// Copyright 2022 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build sonic && avx && (linux || windows || darwin) && amd64

package json

import (
	"io"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin/codec/api"
)

func init() {
	Api = sonicApi{}
}

var json = sonic.ConfigStd

type sonicApi struct {
}

func (j sonicApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j sonicApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j sonicApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j sonicApi) NewEncoder(writer io.Writer) api.JsonEncoder {
	return json.NewEncoder(writer)
}

func (j sonicApi) NewDecoder(reader io.Reader) api.JsonDecoder {
	return json.NewDecoder(reader)
}
