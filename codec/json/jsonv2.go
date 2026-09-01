// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build jsonv2 && go1.27 && goexperiment.jsonv2 && !jsoniter && !go_json && !(sonic && (linux || windows || darwin))

package json

import (
	"encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"io"
)

// Package indicates what library is being used for JSON encoding.
const Package = "encoding/json/v2"

func init() {
	API = jsonv2Api{}
}

// v1Compat matches the observable behavior of encoding/json v1: HTML escaping
// on, case-insensitive member matching, nil slices and maps as null, and the
// legacy omitempty and error semantics.
var v1Compat = json.DefaultOptionsV1()

// newline is written after each streamed value; a package-level slice keeps
// Encode from allocating one per call.
var newline = []byte{'\n'}

type jsonv2Api struct{}

func (jsonv2Api) Marshal(v any) ([]byte, error) {
	return jsonv2.Marshal(v, v1Compat)
}

func (jsonv2Api) Unmarshal(data []byte, v any) error {
	return jsonv2.Unmarshal(data, v, v1Compat)
}

func (jsonv2Api) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return jsonv2.Marshal(v, v1Compat,
		jsontext.WithIndentPrefix(prefix), jsontext.WithIndent(indent))
}

// NewEncoder wraps MarshalWrite rather than returning v1's encoder, which
// marshals into an intermediate buffer so that SetIndent can reformat the
// result. Gin never indents a stream, so writing straight to the io.Writer
// avoids that buffer. Re-measure BenchmarkEncoder before replacing this with
// json.NewEncoder.
func (jsonv2Api) NewEncoder(writer io.Writer) Encoder {
	// v1Compat already escapes HTML, matching v1's default.
	return &v2Encoder{writer: writer, opts: v1Compat}
}

// NewDecoder returns v1's decoder. It satisfies Decoder as-is, implements
// UseNumber natively (v2 exposes no option for it), and measures the same as a
// wrapper over jsontext.Decoder, so there is nothing to gain by wrapping.
func (jsonv2Api) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

type v2Encoder struct {
	writer io.Writer
	opts   jsonv2.Options
}

func (e *v2Encoder) SetEscapeHTML(on bool) {
	e.opts = jsonv2.JoinOptions(v1Compat, jsontext.EscapeForHTML(on))
}

func (e *v2Encoder) Encode(v any) error {
	// MarshalWrite omits the trailing newline that v1's Encoder.Encode writes.
	if err := jsonv2.MarshalWrite(e.writer, v, e.opts); err != nil {
		return err
	}
	_, err := e.writer.Write(newline)
	return err
}
