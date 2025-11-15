// +build amd64,go1.17,!go1.26 arm64,go1.20,!go1.26

/*
 * Copyright 2023 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package encoder

import (
    `github.com/bytedance/sonic/internal/encoder`
)

// EnableFallback indicates if encoder use fallback
const EnableFallback = false

// Encoder represents a specific set of encoder configurations.
type Encoder = encoder.Encoder

// StreamEncoder uses io.Writer as input.
type StreamEncoder = encoder.StreamEncoder

// Options is a set of encoding options.
type Options = encoder.Options

const (
    // SortMapKeys indicates that the keys of a map needs to be sorted
    // before serializing into JSON.
    // WARNING: This hurts performance A LOT, USE WITH CARE.
    SortMapKeys Options = encoder.SortMapKeys

    // EscapeHTML indicates encoder to escape all HTML characters
    // after serializing into JSON (see https://pkg.go.dev/encoding/json#HTMLEscape).
    // WARNING: This hurts performance A LOT, USE WITH CARE.
    EscapeHTML Options = encoder.EscapeHTML

    // CompactMarshaler indicates that the output JSON from json.Marshaler
    // is always compact and needs no validation
    CompactMarshaler Options = encoder.CompactMarshaler

    // NoQuoteTextMarshaler indicates that the output text from encoding.TextMarshaler
    // is always escaped string and needs no quoting
    NoQuoteTextMarshaler Options = encoder.NoQuoteTextMarshaler

    // NoNullSliceOrMap indicates all empty Array or Object are encoded as '[]' or '{}',
    // instead of 'null'
    NoNullSliceOrMap Options = encoder.NoNullSliceOrMap

    // ValidateString indicates that encoder should validate the input string
    // before encoding it into JSON.
    ValidateString Options = encoder.ValidateString

    // NoValidateJSONMarshaler indicates that the encoder should not validate the output string
    // after encoding the JSONMarshaler to JSON.
    NoValidateJSONMarshaler Options = encoder.NoValidateJSONMarshaler

    // NoEncoderNewline indicates that the encoder should not add a newline after every message
    NoEncoderNewline Options = encoder.NoEncoderNewline

    // CompatibleWithStd is used to be compatible with std encoder.
    CompatibleWithStd Options = encoder.CompatibleWithStd

    // EncodeNullForInfOrNan encodes Infinity or NaN float values as 'null'
    // instead of returning an error.
    EncodeNullForInfOrNan Options = encoder.EncodeNullForInfOrNan
)


var (
    // Encode returns the JSON encoding of val, encoded with opts.
    Encode = encoder.Encode

    // EncodeIndented is like Encode but applies Indent to format the output.
    // Each JSON element in the output will begin on a new line beginning with prefix
    // followed by one or more copies of indent according to the indentation nesting.
    EncodeIndented = encoder.EncodeIndented

    // EncodeInto is like Encode but uses a user-supplied buffer instead of allocating a new one.
    EncodeInto = encoder.EncodeInto

    // HTMLEscape appends to dst the JSON-encoded src with <, >, &, U+2028 and U+2029
    // characters inside string literals changed to \u003c, \u003e, \u0026, \u2028, \u2029
    // so that the JSON will be safe to embed inside HTML <script> tags.
    // For historical reasons, web browsers don't honor standard HTML
    // escaping within <script> tags, so an alternative JSON encoding must
    // be used.
    HTMLEscape = encoder.HTMLEscape

    // Pretouch compiles vt ahead-of-time to avoid JIT compilation on-the-fly, in
    // order to reduce the first-hit latency.
    //
    // Opts are the compile options, for example, "option.WithCompileRecursiveDepth" is
    // a compile option to set the depth of recursive compile for the nested struct type.
    Pretouch = encoder.Pretouch

    // Quote returns the JSON-quoted version of s.
    Quote = encoder.Quote

    // Valid validates json and returns first non-blank character position,
    // if it is only one valid json value.
    // Otherwise returns invalid character position using start.
    //
    // Note: it does not check for the invalid UTF-8 characters.
    Valid = encoder.Valid

    // NewStreamEncoder adapts to encoding/json.NewDecoder API.
    //
    // NewStreamEncoder returns a new encoder that write to w.
    NewStreamEncoder = encoder.NewStreamEncoder
)
