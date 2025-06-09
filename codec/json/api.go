// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import "io"

// API the json codec in use.
var API Core

// Core the api for json codec.
type Core interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
	NewEncoder(writer io.Writer) Encoder
	NewDecoder(reader io.Reader) Decoder
}

// Encoder an interface writes JSON values to an output stream.
type Encoder interface {
	// SetEscapeHTML specifies whether problematic HTML characters
	// should be escaped inside JSON quoted strings.
	// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
	// to avoid certain safety problems that can arise when embedding JSON in HTML.
	//
	// In non-HTML settings where the escaping interferes with the readability
	// of the output, SetEscapeHTML(false) disables this behavior.
	SetEscapeHTML(on bool)

	// Encode writes the JSON encoding of v to the stream,
	// followed by a newline character.
	//
	// See the documentation for Marshal for details about the
	// conversion of Go values to JSON.
	Encode(v any) error
}

// Decoder an interface reads and decodes JSON values from an input stream.
type Decoder interface {
	// UseNumber causes the Decoder to unmarshal a number into an any as a
	// Number instead of as a float64.
	UseNumber()

	// DisallowUnknownFields causes the Decoder to return an error when the destination
	// is a struct and the input contains object keys which do not match any
	// non-ignored, exported fields in the destination.
	DisallowUnknownFields()

	// Decode reads the next JSON-encoded value from its
	// input and stores it in the value pointed to by v.
	//
	// See the documentation for Unmarshal for details about
	// the conversion of JSON into a Go value.
	Decode(v any) error
}
