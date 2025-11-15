/*
 * Copyright 2021 ByteDance Inc.
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

package sonic

import (
    `io`

    `github.com/bytedance/sonic/ast`
    `github.com/bytedance/sonic/internal/rt`
)

const (
    // UseStdJSON indicates you are using fallback implementation (encoding/json)
	UseStdJSON = iota
    // UseSonicJSON indicates you are using real sonic implementation
	UseSonicJSON
)

// APIKind is the kind of API, 0 is std json, 1 is sonic.
const APIKind = apiKind

// Config is a combination of sonic/encoder.Options and sonic/decoder.Options
type Config struct {
    // EscapeHTML indicates encoder to escape all HTML characters 
    // after serializing into JSON (see https://pkg.go.dev/encoding/json#HTMLEscape).
    // WARNING: This hurts performance A LOT, USE WITH CARE.
    EscapeHTML                    bool

    // SortMapKeys indicates encoder that the keys of a map needs to be sorted 
    // before serializing into JSON.
    // WARNING: This hurts performance A LOT, USE WITH CARE.
    SortMapKeys                   bool

    // CompactMarshaler indicates encoder that the output JSON from json.Marshaler 
    // is always compact and needs no validation 
    CompactMarshaler              bool

    // NoQuoteTextMarshaler indicates encoder that the output text from encoding.TextMarshaler 
    // is always escaped string and needs no quoting
    NoQuoteTextMarshaler          bool

    // NoNullSliceOrMap indicates encoder that all empty Array or Object are encoded as '[]' or '{}',
    // instead of 'null'
    NoNullSliceOrMap              bool

    // UseInt64 indicates decoder to unmarshal an integer into an interface{} as an
    // int64 instead of as a float64.
    UseInt64                      bool

    // UseNumber indicates decoder to unmarshal a number into an interface{} as a
    // json.Number instead of as a float64.
    UseNumber                     bool

    // UseUnicodeErrors indicates decoder to return an error when encounter invalid
    // UTF-8 escape sequences.
    UseUnicodeErrors              bool

    // DisallowUnknownFields indicates decoder to return an error when the destination
    // is a struct and the input contains object keys which do not match any
    // non-ignored, exported fields in the destination.
    DisallowUnknownFields         bool

    // CopyString indicates decoder to decode string values by copying instead of referring.
    CopyString                    bool

    // ValidateString indicates decoder and encoder to validate string values: decoder will return errors 
    // when unescaped control chars(\u0000-\u001f) in the string value of JSON.
    ValidateString                bool

    // NoValidateJSONMarshaler indicates that the encoder should not validate the output string
    // after encoding the JSONMarshaler to JSON.
    NoValidateJSONMarshaler       bool

    // NoValidateJSONSkip indicates the decoder should not validate the JSON value when skipping it,
    // such as unknown-fields, mismatched-type, redundant elements..
    NoValidateJSONSkip bool
    
    // NoEncoderNewline indicates that the encoder should not add a newline after every message
    NoEncoderNewline bool

    // Encode Infinity or Nan float into `null`, instead of returning an error.
    EncodeNullForInfOrNan bool

    // CaseSensitive indicates that the decoder should not ignore the case of object keys.
    CaseSensitive bool
}
 
var (
    // ConfigDefault is the default config of APIs, aiming at efficiency and safety.
    ConfigDefault = Config{}.Froze()
 
    // ConfigStd is the standard config of APIs, aiming at being compatible with encoding/json.
    ConfigStd = Config{
        EscapeHTML : true,
        SortMapKeys: true,
        CompactMarshaler: true,
        CopyString : true,
        ValidateString : true,
    }.Froze()
 
    // ConfigFastest is the fastest config of APIs, aiming at speed.
    ConfigFastest = Config{
        NoValidateJSONMarshaler: true,
        NoValidateJSONSkip: true,
    }.Froze()
)
 
 
// API is a binding of specific config.
// This interface is inspired by github.com/json-iterator/go,
// and has same behaviors under equivalent config.
type API interface {
    // MarshalToString returns the JSON encoding string of v
    MarshalToString(v interface{}) (string, error)
    // Marshal returns the JSON encoding bytes of v.
    Marshal(v interface{}) ([]byte, error)
    // MarshalIndent returns the JSON encoding bytes with indent and prefix.
    MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
    // UnmarshalFromString parses the JSON-encoded bytes and stores the result in the value pointed to by v.
    UnmarshalFromString(str string, v interface{}) error
    // Unmarshal parses the JSON-encoded string and stores the result in the value pointed to by v.
    Unmarshal(data []byte, v interface{}) error
    // NewEncoder create a Encoder holding writer
    NewEncoder(writer io.Writer) Encoder
    // NewDecoder create a Decoder holding reader
    NewDecoder(reader io.Reader) Decoder
    // Valid validates the JSON-encoded bytes and reports if it is valid
    Valid(data []byte) bool
}

// Encoder encodes JSON into io.Writer
type Encoder interface {
    // Encode writes the JSON encoding of v to the stream, followed by a newline character.
    Encode(val interface{}) error
    // SetEscapeHTML specifies whether problematic HTML characters 
    // should be escaped inside JSON quoted strings. 
    // The default behavior NOT ESCAPE 
    SetEscapeHTML(on bool)
    // SetIndent instructs the encoder to format each subsequent encoded value 
    // as if indented by the package-level function Indent(dst, src, prefix, indent).
    // Calling SetIndent("", "") disables indentation
    SetIndent(prefix, indent string)
}

// Decoder decodes JSON from io.Read
type Decoder interface {
    // Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
    Decode(val interface{}) error
    // Buffered returns a reader of the data remaining in the Decoder's buffer.
    // The reader is valid until the next call to Decode.
    Buffered() io.Reader
    // DisallowUnknownFields causes the Decoder to return an error when the destination is a struct 
    // and the input contains object keys which do not match any non-ignored, exported fields in the destination.
    DisallowUnknownFields()
    // More reports whether there is another element in the current array or object being parsed.
    More() bool
    // UseNumber causes the Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
    UseNumber()
}

// Marshal returns the JSON encoding bytes of v.
func Marshal(val interface{}) ([]byte, error) {
    return ConfigDefault.Marshal(val)
}

// MarshalIndent is like Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
    return ConfigDefault.MarshalIndent(v, prefix, indent)
}

// MarshalString returns the JSON encoding string of v.
func MarshalString(val interface{}) (string, error) {
    return ConfigDefault.MarshalToString(val)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// NOTICE: This API copies given buffer by default,
// if you want to pass JSON more efficiently, use UnmarshalString instead.
func Unmarshal(buf []byte, val interface{}) error {
    return ConfigDefault.Unmarshal(buf, val)
}

// UnmarshalString is like Unmarshal, except buf is a string.
func UnmarshalString(buf string, val interface{}) error {
    return ConfigDefault.UnmarshalFromString(buf, val)
}

// Get searches and locates the given path from src json,
// and returns a ast.Node representing the partially json.
//
// Each path arg must be integer or string:
//     - Integer is target index(>=0), means searching current node as array.
//     - String is target key, means searching current node as object.
//
// 
// Notice: It expects the src json is **Well-formed** and **Immutable** when calling,
// otherwise it may return unexpected result. 
// Considering memory safety, the returned JSON is **Copied** from the input
func Get(src []byte, path ...interface{}) (ast.Node, error) {
    return GetCopyFromString(rt.Mem2Str(src), path...)
}

//GetWithOptions searches and locates the given path from src json,
// with specific options of ast.Searcher
func GetWithOptions(src []byte, opts ast.SearchOptions, path ...interface{}) (ast.Node, error) {
    s := ast.NewSearcher(rt.Mem2Str(src))
    s.SearchOptions = opts
    return s.GetByPath(path...)
}

// GetFromString is same with Get except src is string.
//
// WARNING: The returned JSON is **Referenced** from the input. 
// Caching or long-time holding the returned node may cause OOM.
// If your src is big, consider use GetFromStringCopy().
func GetFromString(src string, path ...interface{}) (ast.Node, error) {
    return ast.NewSearcher(src).GetByPath(path...)
}

// GetCopyFromString is same with Get except src is string
func GetCopyFromString(src string, path ...interface{}) (ast.Node, error) {
    return ast.NewSearcher(src).GetByPathCopy(path...)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
    return ConfigDefault.Valid(data)
}

// Valid reports whether data is a valid JSON encoding.
func ValidString(data string) bool {
    return ConfigDefault.Valid(rt.Str2Mem(data))
}
