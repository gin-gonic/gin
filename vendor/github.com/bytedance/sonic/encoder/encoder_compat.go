// +build !amd64,!arm64 go1.26 !go1.17 arm64,!go1.20

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
    `io`
    `bytes`
    `encoding/json`
    `reflect`

    `github.com/bytedance/sonic/option`
    `github.com/bytedance/sonic/internal/compat`
)

func init() {
    compat.Warn("sonic/encoder")
}

// EnableFallback indicates if encoder use fallback
const EnableFallback = true

// Options is a set of encoding options.
type Options uint64

const (
    bitSortMapKeys          = iota
    bitEscapeHTML          
    bitCompactMarshaler
    bitNoQuoteTextMarshaler
    bitNoNullSliceOrMap
    bitValidateString
    bitNoValidateJSONMarshaler
    bitNoEncoderNewline

    // used for recursive compile
    bitPointerValue = 63
)

const (
    // SortMapKeys indicates that the keys of a map needs to be sorted 
    // before serializing into JSON.
    // WARNING: This hurts performance A LOT, USE WITH CARE.
    SortMapKeys          Options = 1 << bitSortMapKeys

    // EscapeHTML indicates encoder to escape all HTML characters 
    // after serializing into JSON (see https://pkg.go.dev/encoding/json#HTMLEscape).
    // WARNING: This hurts performance A LOT, USE WITH CARE.
    EscapeHTML           Options = 1 << bitEscapeHTML

    // CompactMarshaler indicates that the output JSON from json.Marshaler 
    // is always compact and needs no validation 
    CompactMarshaler     Options = 1 << bitCompactMarshaler

    // NoQuoteTextMarshaler indicates that the output text from encoding.TextMarshaler 
    // is always escaped string and needs no quoting
    NoQuoteTextMarshaler Options = 1 << bitNoQuoteTextMarshaler

    // NoNullSliceOrMap indicates all empty Array or Object are encoded as '[]' or '{}',
    // instead of 'null'
    NoNullSliceOrMap     Options = 1 << bitNoNullSliceOrMap

    // ValidateString indicates that encoder should validate the input string
    // before encoding it into JSON.
    ValidateString       Options = 1 << bitValidateString

    // NoValidateJSONMarshaler indicates that the encoder should not validate the output string
    // after encoding the JSONMarshaler to JSON.
    NoValidateJSONMarshaler Options = 1 << bitNoValidateJSONMarshaler

    // NoEncoderNewline indicates that the encoder should not add a newline after every message
    NoEncoderNewline Options = 1 << bitNoEncoderNewline
  
    // CompatibleWithStd is used to be compatible with std encoder.
    CompatibleWithStd Options = SortMapKeys | EscapeHTML | CompactMarshaler
)

// Encoder represents a specific set of encoder configurations.
type Encoder struct {
    Opts Options
    prefix string
    indent string
}

// Encode returns the JSON encoding of v.
func (self *Encoder) Encode(v interface{}) ([]byte, error) {
    if self.indent != "" || self.prefix != "" { 
        return EncodeIndented(v, self.prefix, self.indent, self.Opts)
    }
    return Encode(v, self.Opts)
}

// SortKeys enables the SortMapKeys option.
func (self *Encoder) SortKeys() *Encoder {
    self.Opts |= SortMapKeys
    return self
}

// SetEscapeHTML specifies if option EscapeHTML opens
func (self *Encoder) SetEscapeHTML(f bool) {
    if f {
        self.Opts |= EscapeHTML
    } else {
        self.Opts &= ^EscapeHTML
    }
}

// SetValidateString specifies if option ValidateString opens
func (self *Encoder) SetValidateString(f bool) {
    if f {
        self.Opts |= ValidateString
    } else {
        self.Opts &= ^ValidateString
    }
}

// SetNoValidateJSONMarshaler specifies if option NoValidateJSONMarshaler opens
func (self *Encoder) SetNoValidateJSONMarshaler(f bool) {
    if f {
        self.Opts |= NoValidateJSONMarshaler
    } else {
        self.Opts &= ^NoValidateJSONMarshaler
    }
}

// SetNoEncoderNewline specifies if option NoEncoderNewline opens
func (self *Encoder) SetNoEncoderNewline(f bool) {
    if f {
        self.Opts |= NoEncoderNewline
    } else {
        self.Opts &= ^NoEncoderNewline
    }
}

// SetCompactMarshaler specifies if option CompactMarshaler opens
func (self *Encoder) SetCompactMarshaler(f bool) {
    if f {
        self.Opts |= CompactMarshaler
    } else {
        self.Opts &= ^CompactMarshaler
    }
}

// SetNoQuoteTextMarshaler specifies if option NoQuoteTextMarshaler opens
func (self *Encoder) SetNoQuoteTextMarshaler(f bool) {
    if f {
        self.Opts |= NoQuoteTextMarshaler
    } else {
        self.Opts &= ^NoQuoteTextMarshaler
    }
}

// SetIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function EncodeIndent().
// Calling SetIndent("", "") disables indentation.
func (enc *Encoder) SetIndent(prefix, indent string) {
    enc.prefix = prefix
    enc.indent = indent
}

// Quote returns the JSON-quoted version of s.
func Quote(s string) string {
    /* check for empty string */
    if s == "" {
        return `""`
    }

    out, _ := json.Marshal(s)
    return string(out)
}

// Encode returns the JSON encoding of val, encoded with opts.
func Encode(val interface{}, opts Options) ([]byte, error) {
   return json.Marshal(val)
}

// EncodeInto is like Encode but uses a user-supplied buffer instead of allocating
// a new one.
func EncodeInto(buf *[]byte, val interface{}, opts Options) error {
    if buf == nil {
        panic("user-supplied buffer buf is nil")
    }
    w := bytes.NewBuffer(*buf)
    enc := json.NewEncoder(w)
    enc.SetEscapeHTML((opts & EscapeHTML) != 0)
    err := enc.Encode(val)
    *buf = w.Bytes()
    l := len(*buf)
    if l > 0 && (opts & NoEncoderNewline != 0) && (*buf)[l-1] == '\n' {
        *buf = (*buf)[:l-1]
    }
    return err
}

// HTMLEscape appends to dst the JSON-encoded src with <, >, &, U+2028 and U+2029
// characters inside string literals changed to \u003c, \u003e, \u0026, \u2028, \u2029
// so that the JSON will be safe to embed inside HTML <script> tags.
// For historical reasons, web browsers don't honor standard HTML
// escaping within <script> tags, so an alternative JSON encoding must
// be used.
func HTMLEscape(dst []byte, src []byte) []byte {
   d := bytes.NewBuffer(dst)
   json.HTMLEscape(d, src)
   return d.Bytes()
}

// EncodeIndented is like Encode but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func EncodeIndented(val interface{}, prefix string, indent string, opts Options) ([]byte, error) {
   w := bytes.NewBuffer([]byte{})
   enc := json.NewEncoder(w)
   enc.SetEscapeHTML((opts & EscapeHTML) != 0)
   enc.SetIndent(prefix, indent)
   err := enc.Encode(val)
   out := w.Bytes()
   return out, err
}

// Pretouch compiles vt ahead-of-time to avoid JIT compilation on-the-fly, in
// order to reduce the first-hit latency.
//
// Opts are the compile options, for example, "option.WithCompileRecursiveDepth" is
// a compile option to set the depth of recursive compile for the nested struct type.
func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
   return nil
}

// Valid validates json and returns first non-blank character position,
// if it is only one valid json value.
// Otherwise returns invalid character position using start.
//
// Note: it does not check for the invalid UTF-8 characters.
func Valid(data []byte) (ok bool, start int) {
   return json.Valid(data), 0
}

// StreamEncoder uses io.Writer as 
type StreamEncoder = json.Encoder

// NewStreamEncoder adapts to encoding/json.NewDecoder API.
//
// NewStreamEncoder returns a new encoder that write to w.
func NewStreamEncoder(w io.Writer) *StreamEncoder {
   return json.NewEncoder(w)
}

