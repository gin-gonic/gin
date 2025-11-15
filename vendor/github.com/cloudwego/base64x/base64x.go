/*
 * Copyright 2024 CloudWeGo Authors
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

package base64x

import (
    `encoding/base64`

    "github.com/cloudwego/base64x/internal/native"
)

// An Encoding is a radix 64 encoding/decoding scheme, defined by a
// 64-character alphabet. The most common encoding is the "base64"
// encoding defined in RFC 4648 and used in MIME (RFC 2045) and PEM
// (RFC 1421).  RFC 4648 also defines an alternate encoding, which is
// the standard encoding with - and _ substituted for + and /.
type Encoding int

const (
    _MODE_URL  = 1 << 0
    _MODE_RAW  = 1 << 1
    _MODE_AVX2 = 1 << 2
    _MODE_JSON = 1 << 3
)

// StdEncoding is the standard base64 encoding, as defined in
// RFC 4648.
const StdEncoding Encoding = 0

// URLEncoding is the alternate base64 encoding defined in RFC 4648.
// It is typically used in URLs and file names.
const URLEncoding Encoding = _MODE_URL

// RawStdEncoding is the standard raw, unpadded base64 encoding,
// as defined in RFC 4648 section 3.2.
//
// This is the same as StdEncoding but omits padding characters.
const RawStdEncoding Encoding = _MODE_RAW

// RawURLEncoding is the unpadded alternate base64 encoding defined in RFC 4648.
// It is typically used in URLs and file names.
//
// This is the same as URLEncoding but omits padding characters.
const RawURLEncoding Encoding = _MODE_RAW | _MODE_URL

// JSONStdEncoding is the StdEncoding and encoded as JSON string as RFC 8259.
const JSONStdEncoding Encoding = _MODE_JSON;

var (
    archFlags = 0
)

/** Encoder Functions **/

// Encode encodes src using the specified encoding, writing
// EncodedLen(len(src)) bytes to out.
//
// The encoding pads the output to a multiple of 4 bytes,
// so Encode is not appropriate for use on individual blocks
// of a large data stream.
//
// If out is not large enough to contain the encoded result,
// it will panic.
func (self Encoding) Encode(out []byte, src []byte) {
    if len(src) != 0 {
        if buf := out[:0:len(out)]; self.EncodedLen(len(src)) <= len(out) {
            self.EncodeUnsafe(&buf, src)
        } else {
            panic("encoder output buffer is too small")
        }
    }
}

// EncodeUnsafe behaves like Encode, except it does NOT check if
// out is large enough to contain the encoded result.
//
// It will also update the length of out.
func (self Encoding) EncodeUnsafe(out *[]byte, src []byte) {
    native.B64Encode(out, &src, int(self) | archFlags)
}

// EncodeToString returns the base64 encoding of src.
func (self Encoding) EncodeToString(src []byte) string {
    nbs := len(src)
    ret := make([]byte, 0, self.EncodedLen(nbs))

    /* encode in native code */
    self.EncodeUnsafe(&ret, src)
    return mem2str(ret)
}

// EncodedLen returns the length in bytes of the base64 encoding
// of an input buffer of length n.
func (self Encoding) EncodedLen(n int) int {
    if (self & _MODE_RAW) == 0 {
        return (n + 2) / 3 * 4
    } else {
        return (n * 8 + 5) / 6
    }
}

/** Decoder Functions **/

// Decode decodes src using the encoding enc. It writes at most
// DecodedLen(len(src)) bytes to out and returns the number of bytes
// written. If src contains invalid base64 data, it will return the
// number of bytes successfully written and base64.CorruptInputError.
//
// New line characters (\r and \n) are ignored.
//
// If out is not large enough to contain the encoded result,
// it will panic.
func (self Encoding) Decode(out []byte, src []byte) (int, error) {
    if len(src) == 0 {
        return 0, nil
    } else if buf := out[:0:len(out)]; self.DecodedLen(len(src)) <= len(out) {
        return self.DecodeUnsafe(&buf, src)
    } else {
        panic("decoder output buffer is too small")
    }
}

// DecodeUnsafe behaves like Decode, except it does NOT check if
// out is large enough to contain the decoded result.
//
// It will also update the length of out.
func (self Encoding) DecodeUnsafe(out *[]byte, src []byte) (int, error) {
    if n := native.B64Decode(out, mem2addr(src), len(src), int(self) | archFlags); n >= 0 {
        return n, nil
    } else {
        return 0, base64.CorruptInputError(-n - 1)
    }
}

// DecodeString returns the bytes represented by the base64 string s.
func (self Encoding) DecodeString(s string) ([]byte, error) {
    src := str2mem(s)
    ret := make([]byte, 0, self.DecodedLen(len(s)))

    /* decode into the allocated buffer */
    if _, err := self.DecodeUnsafe(&ret, src); err != nil {
        return nil, err
    } else {
        return ret, nil
    }
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base64-encoded data.
func (self Encoding) DecodedLen(n int) int {
    if (self & _MODE_RAW) == 0 {
        return n / 4 * 3
    } else {
        return n * 6 / 8
    }
}
