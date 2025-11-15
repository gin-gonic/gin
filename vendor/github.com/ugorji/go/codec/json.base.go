// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode"
)

//--------------------------------

// jsonLits and jsonLitb are defined at the package level,
// so they are guaranteed to be stored efficiently, making
// for better append/string comparison/etc.
//
// (anecdotal evidence from some benchmarking on go 1.20 devel in 20220104)
const jsonLits = `"true"false"null"{}[]`

const (
	jsonLitT = 1
	jsonLitF = 6
	jsonLitN = 12
	jsonLitM = 17
	jsonLitA = 19
)

var jsonLitb = []byte(jsonLits)
var jsonNull = jsonLitb[jsonLitN : jsonLitN+4]
var jsonArrayEmpty = jsonLitb[jsonLitA : jsonLitA+2]
var jsonMapEmpty = jsonLitb[jsonLitM : jsonLitM+2]

const jsonEncodeUintSmallsString = "" +
	"00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

var jsonEncodeUintSmallsStringBytes = (*[len(jsonEncodeUintSmallsString)]byte)([]byte(jsonEncodeUintSmallsString))

const (
	jsonU4Chk2 = '0'
	jsonU4Chk1 = 'a' - 10
	jsonU4Chk0 = 'A' - 10
)

const (
	// If !jsonValidateSymbols, decoding will be faster, by skipping some checks:
	//   - If we see first character of null, false or true,
	//     do not validate subsequent characters.
	//   - e.g. if we see a n, assume null and skip next 3 characters,
	//     and do not validate they are ull.
	// P.S. Do not expect a significant decoding boost from this.
	jsonValidateSymbols = true

	// jsonEscapeMultiByteUnicodeSep controls whether some unicode characters
	// that are valid json but may bomb in some contexts are escaped during encoeing.
	//
	// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
	// Both technically valid JSON, but bomb on JSONP, so fix here unconditionally.
	jsonEscapeMultiByteUnicodeSep = true

	// jsonNakedBoolNumInQuotedStr is used during decoding into a blank interface{}
	// to control whether we detect quoted values of bools and null where a map key is expected,
	// and treat as nil, true or false.
	jsonNakedBoolNumInQuotedStr = true
)

var (
	// jsonTabs and jsonSpaces are used as caches for indents
	jsonTabs   [32]byte
	jsonSpaces [128]byte

	jsonHexEncoder hexEncoder
	// jsonTimeLayout is used to validate time layouts.
	// Unfortunately, we couldn't compare time.Time effectively, so punted.
	// jsonTimeLayout time.Time
)

func init() {
	for i := 0; i < len(jsonTabs); i++ {
		jsonTabs[i] = '\t'
	}
	for i := 0; i < len(jsonSpaces); i++ {
		jsonSpaces[i] = ' '
	}
	// jsonTimeLayout, err := time.Parse(time.Layout, time.Layout)
	// halt.onerror(err)
	// jsonTimeLayout = jsonTimeLayout.Round(time.Second).UTC()
}

// ----------------

type jsonBytesFmt uint8

const (
	jsonBytesFmtArray jsonBytesFmt = iota + 1
	jsonBytesFmtBase64
	jsonBytesFmtBase64url
	jsonBytesFmtBase32
	jsonBytesFmtBase32hex
	jsonBytesFmtBase16

	jsonBytesFmtHex = jsonBytesFmtBase16
)

type jsonTimeFmt uint8

const (
	jsonTimeFmtStringLayout jsonTimeFmt = iota + 1
	jsonTimeFmtUnix
	jsonTimeFmtUnixMilli
	jsonTimeFmtUnixMicro
	jsonTimeFmtUnixNano
)

type jsonBytesFmter = bytesEncoder

type jsonHandleOpts struct {
	rawext bool
	// bytesFmt used during encode to determine how to encode []byte
	bytesFmt jsonBytesFmt
	// timeFmt used during encode to determine how to encode a time.Time
	timeFmt jsonTimeFmt
	// timeFmtNum used during decode to decode a time.Time from an int64 in the stream
	timeFmtNum jsonTimeFmt
	// timeFmtLayouts used on decode, to try to parse time.Time until successful
	timeFmtLayouts []string
	// byteFmters used on decode, to try to parse []byte from a UTF-8 string encoding (e.g. base64)
	byteFmters []jsonBytesFmter
}

func jsonCheckTimeLayout(s string) (ok bool) {
	_, err := time.Parse(s, s)
	// t...Equal(jsonTimeLayout) always returns false - unsure why
	// return err == nil && t.Round(time.Second).UTC().Equal(jsonTimeLayout)
	return err == nil
}

func (x *jsonHandleOpts) reset(h *JsonHandle) {
	x.timeFmt = 0
	x.timeFmtNum = 0
	x.timeFmtLayouts = x.timeFmtLayouts[:0]
	if len(h.TimeFormat) != 0 {
		switch h.TimeFormat[0] {
		case "unix":
			x.timeFmt = jsonTimeFmtUnix
		case "unixmilli":
			x.timeFmt = jsonTimeFmtUnixMilli
		case "unixmicro":
			x.timeFmt = jsonTimeFmtUnixMicro
		case "unixnano":
			x.timeFmt = jsonTimeFmtUnixNano
		}
		x.timeFmtNum = x.timeFmt
		for _, v := range h.TimeFormat {
			if !strings.HasPrefix(v, "unix") && jsonCheckTimeLayout(v) {
				x.timeFmtLayouts = append(x.timeFmtLayouts, v)
			}
		}
	}
	if x.timeFmt == 0 { // both timeFmt and timeFmtNum are 0
		x.timeFmtNum = jsonTimeFmtUnix
		x.timeFmt = jsonTimeFmtStringLayout
		if len(x.timeFmtLayouts) == 0 {
			x.timeFmtLayouts = append(x.timeFmtLayouts, time.RFC3339Nano)
		}
	}

	x.bytesFmt = 0
	x.byteFmters = x.byteFmters[:0]
	var b64 bool
	if len(h.BytesFormat) != 0 {
		switch h.BytesFormat[0] {
		case "array":
			x.bytesFmt = jsonBytesFmtArray
		case "base64":
			x.bytesFmt = jsonBytesFmtBase64
		case "base64url":
			x.bytesFmt = jsonBytesFmtBase64url
		case "base32":
			x.bytesFmt = jsonBytesFmtBase32
		case "base32hex":
			x.bytesFmt = jsonBytesFmtBase32hex
		case "base16", "hex":
			x.bytesFmt = jsonBytesFmtBase16
		}
		for _, v := range h.BytesFormat {
			switch v {
			// case "array":
			case "base64":
				x.byteFmters = append(x.byteFmters, base64.StdEncoding)
				b64 = true
			case "base64url":
				x.byteFmters = append(x.byteFmters, base64.URLEncoding)
			case "base32":
				x.byteFmters = append(x.byteFmters, base32.StdEncoding)
			case "base32hex":
				x.byteFmters = append(x.byteFmters, base32.HexEncoding)
			case "base16", "hex":
				x.byteFmters = append(x.byteFmters, &jsonHexEncoder)
			}
		}
	}
	if x.bytesFmt == 0 {
		// either len==0 OR gibberish was in the first element; resolve here
		x.bytesFmt = jsonBytesFmtBase64
		if !b64 { // not present - so insert into pos 0
			x.byteFmters = append(x.byteFmters, nil)
			copy(x.byteFmters[1:], x.byteFmters[0:])
			x.byteFmters[0] = base64.StdEncoding
		}
	}
	// ----
	x.rawext = h.RawBytesExt != nil
}

var jsonEncBoolStrs = [2][2]string{
	{jsonLits[jsonLitF : jsonLitF+5], jsonLits[jsonLitT : jsonLitT+4]},
	{jsonLits[jsonLitF-1 : jsonLitF+6], jsonLits[jsonLitT-1 : jsonLitT+5]},
}

func jsonEncodeUint(neg, quotes bool, u uint64, b *[48]byte) []byte {
	// MARKER: use setByteAt/byteAt to elide the bounds-checks
	// when we are sure that we don't go beyond the bounds.

	// MARKER: copied mostly from std library: strconv/itoa.go
	// this should only be called on 64bit OS.

	var ss = jsonEncodeUintSmallsStringBytes[:]

	// typically, 19 or 20 bytes sufficient for decimal encoding a uint64
	var a = b[:24]
	var i = uint(len(a))

	if quotes {
		i--
		setByteAt(a, i, '"')
		// a[i] = '"'
	}
	var is, us uint // use uint, as those fit into a register on the platform
	if cpu32Bit {
		for u >= 1e9 {
			q := u / 1e9
			us = uint(u - q*1e9) // u % 1e9 fits into a uint
			for j := 4; j > 0; j-- {
				is = us % 100 * 2
				us /= 100
				i -= 2
				setByteAt(a, i+1, byteAt(ss, is+1))
				setByteAt(a, i, byteAt(ss, is))
			}
			i--
			setByteAt(a, i, byteAt(ss, us*2+1))
			u = q
		}
		// u is now < 1e9, so is guaranteed to fit into a uint
	}
	us = uint(u)
	for us >= 100 {
		is = us % 100 * 2
		us /= 100
		i -= 2
		setByteAt(a, i+1, byteAt(ss, is+1))
		setByteAt(a, i, byteAt(ss, is))
		// a[i+1], a[i] = ss[is+1], ss[is]
	}

	// us < 100
	is = us * 2
	i--
	setByteAt(a, i, byteAt(ss, is+1))
	// a[i] = ss[is+1]
	if us >= 10 {
		i--
		setByteAt(a, i, byteAt(ss, is))
		// a[i] = ss[is]
	}
	if neg {
		i--
		setByteAt(a, i, '-')
		// a[i] = '-'
	}
	if quotes {
		i--
		setByteAt(a, i, '"')
		// a[i] = '"'
	}
	return a[i:]
}

// MARKER: checkLitErr methods to prevent the got/expect parameters from escaping

//go:noinline
func jsonCheckLitErr3(got, expect [3]byte) {
	halt.errorf("expecting %s: got %s", expect, got)
}

//go:noinline
func jsonCheckLitErr4(got, expect [4]byte) {
	halt.errorf("expecting %s: got %s", expect, got)
}

func jsonSlashURune(cs [4]byte) (rr uint32) {
	for _, c := range cs {
		// best to use explicit if-else
		// - not a table, etc which involve memory loads, array lookup with bounds checks, etc
		if c >= '0' && c <= '9' {
			rr = rr*16 + uint32(c-jsonU4Chk2)
		} else if c >= 'a' && c <= 'f' {
			rr = rr*16 + uint32(c-jsonU4Chk1)
		} else if c >= 'A' && c <= 'F' {
			rr = rr*16 + uint32(c-jsonU4Chk0)
		} else {
			return unicode.ReplacementChar
		}
	}
	return
}

func jsonNakedNum(z *fauxUnion, bs []byte, preferFloat, signedInt bool) (err error) {
	// Note: jsonNakedNum is NEVER called with a zero-length []byte
	if preferFloat {
		z.v = valueTypeFloat
		z.f, err = parseFloat64(bs)
	} else {
		err = parseNumber(bs, z, signedInt)
	}
	return
}

//----------------------

// JsonHandle is a handle for JSON encoding format.
//
// Json is comprehensively supported:
//   - decodes numbers into interface{} as int, uint or float64
//     based on how the number looks and some config parameters e.g. PreferFloat, SignedInt, etc.
//   - decode integers from float formatted numbers e.g. 1.27e+8
//   - decode any json value (numbers, bool, etc) from quoted strings
//   - configurable way to encode/decode []byte .
//     by default, encodes and decodes []byte using base64 Std Encoding
//   - UTF-8 support for encoding and decoding
//
// It has better performance than the json library in the standard library,
// by leveraging the performance improvements of the codec library.
//
// In addition, it doesn't read more bytes than necessary during a decode, which allows
// reading multiple values from a stream containing json and non-json content.
// For example, a user can read a json value, then a cbor value, then a msgpack value,
// all from the same stream in sequence.
//
// Note that, when decoding quoted strings, invalid UTF-8 or invalid UTF-16 surrogate pairs are
// not treated as an error. Instead, they are replaced by the Unicode replacement character U+FFFD.
//
// Note also that the float values for NaN, +Inf or -Inf are encoded as null,
// as suggested by NOTE 4 of the ECMA-262 ECMAScript Language Specification 5.1 edition.
// see http://www.ecma-international.org/publications/files/ECMA-ST/Ecma-262.pdf .
//
// Note the following behaviour differences vs std-library encoding/json package:
//   - struct field names matched in case-sensitive manner
type JsonHandle struct {
	textEncodingType
	BasicHandle

	// Indent indicates how a value is encoded.
	//   - If positive, indent by that number of spaces.
	//   - If negative, indent by that number of tabs.
	Indent int8

	// IntegerAsString controls how integers (signed and unsigned) are encoded.
	//
	// Per the JSON Spec, JSON numbers are 64-bit floating point numbers.
	// Consequently, integers > 2^53 cannot be represented as a JSON number without losing precision.
	// This can be mitigated by configuring how to encode integers.
	//
	// IntegerAsString interpretes the following values:
	//   - if 'L', then encode integers > 2^53 as a json string.
	//   - if 'A', then encode all integers as a json string
	//             containing the exact integer representation as a decimal.
	//   - else    encode all integers as a json number (default)
	IntegerAsString byte

	// HTMLCharsAsIs controls how to encode some special characters to html: < > &
	//
	// By default, we encode them as \uXXX
	// to prevent security holes when served from some browsers.
	HTMLCharsAsIs bool

	// PreferFloat says that we will default to decoding a number as a float.
	// If not set, we will examine the characters of the number and decode as an
	// integer type if it doesn't have any of the characters [.eE].
	PreferFloat bool

	// TermWhitespace says that we add a whitespace character
	// at the end of an encoding.
	//
	// The whitespace is important, especially if using numbers in a context
	// where multiple items are written to a stream.
	TermWhitespace bool

	// MapKeyAsString says to encode all map keys as strings.
	//
	// Use this to enforce strict json output.
	// The only caveat is that nil value is ALWAYS written as null (never as "null")
	MapKeyAsString bool

	// _ uint64 // padding (cache line)

	// Note: below, we store hardly-used items e.g. RawBytesExt.
	// These values below may straddle a cache line, but they are hardly-used,
	// so shouldn't contribute to false-sharing except in rare cases.

	// RawBytesExt, if configured, is used to encode and decode raw bytes in a custom way.
	// If not configured, raw bytes are encoded to/from base64 text.
	RawBytesExt InterfaceExt

	// TimeFormat is an array of strings representing a time.Time format, with each one being either
	// a layout that honor the time.Time.Format specification.
	// In addition, at most one of the set below (unix, unixmilli, unixmicro, unixnana) can be specified
	// supporting encoding and decoding time as a number relative to the time epoch of Jan 1, 1970.
	//
	// During encode of a time.Time, the first entry in the array is used (defaults to RFC 3339).
	//
	// During decode,
	// - if a string, then each of the layout formats will be tried in order until a time.Time is decoded.
	// - if a number, then the sole unix entry is used.
	TimeFormat []string

	// BytesFormat is an array of strings representing how bytes are encoded.
	//
	// Supported values are base64 (default), base64url, base32, base32hex, base16 (synonymous with hex) and array.
	//
	// array is a special value configuring that bytes are encoded as a sequence of numbers.
	//
	// During encode of a []byte, the first entry is used (defaults to base64 if none specified).
	//
	// During decode
	// - if a string, then attempt decoding using each format in sequence until successful.
	// - if an array, then decode normally
	BytesFormat []string
}

func (h *JsonHandle) isJson() bool { return true }

// Name returns the name of the handle: json
func (h *JsonHandle) Name() string { return "json" }

// func (h *JsonHandle) desc(bd byte) string { return str4byte(bd) }
func (h *JsonHandle) desc(bd byte) string { return string(bd) }

func (h *JsonHandle) typical() bool {
	return h.Indent == 0 && !h.MapKeyAsString && h.IntegerAsString != 'A' && h.IntegerAsString != 'L'
}

// SetInterfaceExt sets an extension
func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}

func jsonFloatStrconvFmtPrec64(f float64) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	fbits := math.Float64bits(f)
	abs := math.Float64frombits(fbits &^ (1 << 63))
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac64(fbits) {
		prec = 1
	}
	return
}

func jsonFloatStrconvFmtPrec32(f float32) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	// directly handle Modf (to get fractions) and Abs (to get absolute)
	fbits := math.Float32bits(f)
	abs := math.Float32frombits(fbits &^ (1 << 31))
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac32(fbits) {
		prec = 1
	}
	return
}

var errJsonNoBd = errors.New("descBd unsupported in json")
