// Copyright 2022 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import "io"

// JsonApi the api for json codec.
type JsonApi interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
	NewEncoder(writer io.Writer) JsonEncoder
	NewDecoder(reader io.Reader) JsonDecoder
}

// A JsonEncoder interface writes JSON values to an output stream.
type JsonEncoder interface {
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
	Encode(v interface{}) error
}

// A JsonDecoder interface reads and decodes JSON values from an input stream.
type JsonDecoder interface {
	// UseNumber causes the Decoder to unmarshal a number into an interface{} as a
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
	Decode(v interface{}) error
}
