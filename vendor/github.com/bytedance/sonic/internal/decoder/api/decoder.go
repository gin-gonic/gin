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

package api

import (
    `reflect`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
	`github.com/bytedance/sonic/internal/decoder/consts`
	`github.com/bytedance/sonic/internal/decoder/errors`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/option`
)

const (
	_F_allow_control = consts.F_allow_control
	_F_copy_string = consts.F_copy_string
	_F_disable_unknown = consts.F_disable_unknown
	_F_disable_urc = consts.F_disable_urc
	_F_use_int64 = consts.F_use_int64
	_F_use_number = consts.F_use_number
	_F_validate_string = consts.F_validate_string
    _F_case_sensitive = consts.F_case_sensitive

	_MaxStack = consts.MaxStack

	OptionUseInt64 	       = consts.OptionUseInt64
	OptionUseNumber        = consts.OptionUseNumber
    OptionUseUnicodeErrors = consts.OptionUseUnicodeErrors
    OptionDisableUnknown   = consts.OptionDisableUnknown
    OptionCopyString       = consts.OptionCopyString
    OptionValidateString   = consts.OptionValidateString
    OptionNoValidateJSON   = consts.OptionNoValidateJSON
    OptionCaseSensitive    = consts.OptionCaseSensitive
)

type (
	Options = consts.Options
	MismatchTypeError = errors.MismatchTypeError
	SyntaxError = errors.SyntaxError
)

func (self *Decoder) SetOptions(opts Options) {
    if (opts & consts.OptionUseNumber != 0) && (opts & consts.OptionUseInt64 != 0) {
        panic("can't set OptionUseInt64 and OptionUseNumber both!")
    }
    self.f = uint64(opts)
}

// Decoder is the decoder context object
type Decoder struct {
    i int
    f uint64
    s string
}

// NewDecoder creates a new decoder instance.
func NewDecoder(s string) *Decoder {
    return &Decoder{s: s}
}

// Pos returns the current decoding position.
func (self *Decoder) Pos() int {
    return self.i
}

func (self *Decoder) Reset(s string) {
    self.s = s
    self.i = 0
    // self.f = 0
}

func (self *Decoder) CheckTrailings() error {
    pos := self.i
    buf := self.s
    /* skip all the trailing spaces */
    if pos != len(buf) {
        for pos < len(buf) && (types.SPACE_MASK & (1 << buf[pos])) != 0 {
            pos++
        }
    }

    /* then it must be at EOF */
    if pos == len(buf) {
        return nil
    }

    /* junk after JSON value */
    return SyntaxError {
        Src  : buf,
        Pos  : pos,
        Code : types.ERR_INVALID_CHAR,
    }
}


// Decode parses the JSON-encoded data from current position and stores the result
// in the value pointed to by val.
func (self *Decoder) Decode(val interface{}) error {
	return decodeImpl(&self.s, &self.i, self.f, val)
}

// UseInt64 indicates the Decoder to unmarshal an integer into an interface{} as an
// int64 instead of as a float64.
func (self *Decoder) UseInt64() {
    self.f  |= 1 << _F_use_int64
    self.f &^= 1 << _F_use_number
}

// UseNumber indicates the Decoder to unmarshal a number into an interface{} as a
// json.Number instead of as a float64.
func (self *Decoder) UseNumber() {
    self.f &^= 1 << _F_use_int64
    self.f  |= 1 << _F_use_number
}

// UseUnicodeErrors indicates the Decoder to return an error when encounter invalid
// UTF-8 escape sequences.
func (self *Decoder) UseUnicodeErrors() {
    self.f |= 1 << _F_disable_urc
}

// DisallowUnknownFields indicates the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func (self *Decoder) DisallowUnknownFields() {
    self.f |= 1 << _F_disable_unknown
}

// CopyString indicates the Decoder to decode string values by copying instead of referring.
func (self *Decoder) CopyString() {
    self.f |= 1 << _F_copy_string
}

// ValidateString causes the Decoder to validate string values when decoding string value 
// in JSON. Validation is that, returning error when unescaped control chars(0x00-0x1f) or
// invalid UTF-8 chars in the string value of JSON.
func (self *Decoder) ValidateString() {
    self.f |= 1 << _F_validate_string
}

// Pretouch compiles vt ahead-of-time to avoid JIT compilation on-the-fly, in
// order to reduce the first-hit latency.
//
// Opts are the compile options, for example, "option.WithCompileRecursiveDepth" is
// a compile option to set the depth of recursive compile for the nested struct type.
func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
	return pretouchImpl(vt, opts...)
}

// Skip skips only one json value, and returns first non-blank character position and its ending position if it is valid.
// Otherwise, returns negative error code using start and invalid character position using end
func Skip(data []byte) (start int, end int) {
    s := rt.Mem2Str(data)
    p := 0
    m := types.NewStateMachine()
    ret := native.SkipOne(&s, &p, m, uint64(0))
    types.FreeStateMachine(m) 
    return ret, p
}
