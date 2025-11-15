//go:build notmono || codec.notmono

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// By default, this json support uses base64 encoding for bytes, because you cannot
// store and read any arbitrary string in json (only unicode).
// However, the user can configre how to encode/decode bytes.
//
// This library specifically supports UTF-8 for encoding and decoding only.
//
// Note that the library will happily encode/decode things which are not valid
// json e.g. a map[int64]string. We do it for consistency. With valid json,
// we will encode and decode appropriately.
// Users can specify their map type if necessary to force it.
//
// We cannot use strconv.(Q|Unq)uote because json quotes/unquotes differently.

import (
	"encoding/base64"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

type jsonEncDriver[T encWriter] struct {
	noBuiltInTypes
	h *JsonHandle
	e *encoderBase
	s *bitset256 // safe set for characters (taking h.HTMLAsIs into consideration)

	w T
	// se interfaceExtWrapper

	enc encoderI

	timeFmtLayout string
	byteFmter     jsonBytesFmter
	// ---- cpu cache line boundary???

	// bytes2Arr bool
	// time2Num  bool
	timeFmt  jsonTimeFmt
	bytesFmt jsonBytesFmt

	di int8   // indent per: if negative, use tabs
	d  bool   // indenting?
	dl uint16 // indent level

	ks bool // map key as string
	is byte // integer as string

	typical bool

	rawext bool // rawext configured on the handle

	// buf *[]byte // used mostly for encoding []byte

	// scratch buffer for: encode time, numbers, etc
	//
	// RFC3339Nano uses 35 chars: 2006-01-02T15:04:05.999999999Z07:00
	// MaxUint64 uses 20 chars: 18446744073709551615
	// floats are encoded using: f/e fmt, and -1 precision, or 1 if no fractions.
	// This means we are limited by the number of characters for the
	// mantissa (up to 17), exponent (up to 3), signs (up to 3), dot (up to 1), E (up to 1)
	// for a total of 24 characters.
	//    -xxx.yyyyyyyyyyyye-zzz
	// Consequently, 35 characters should be sufficient for encoding time, integers or floats.
	// We use up all the remaining bytes to make this use full cache lines.
	b [48]byte
}

func (e *jsonEncDriver[T]) writeIndent() {
	e.w.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.di < 0 {
		x = -x
		for x > len(jsonTabs) {
			e.w.writeb(jsonTabs[:])
			x -= len(jsonTabs)
		}
		e.w.writeb(jsonTabs[:x])
	} else {
		for x > len(jsonSpaces) {
			e.w.writeb(jsonSpaces[:])
			x -= len(jsonSpaces)
		}
		e.w.writeb(jsonSpaces[:x])
	}
}

func (e *jsonEncDriver[T]) WriteArrayElem(firstTime bool) {
	if !firstTime {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver[T]) WriteMapElemKey(firstTime bool) {
	if !firstTime {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver[T]) WriteMapElemValue() {
	if e.d {
		e.w.writen2(':', ' ')
	} else {
		e.w.writen1(':')
	}
}

func (e *jsonEncDriver[T]) EncodeNil() {
	// We always encode nil as just null (never in quotes)
	// so we can easily decode if a nil in the json stream ie if initial token is n.

	// e.w.writestr(jsonLits[jsonLitN : jsonLitN+4])
	e.w.writeb(jsonNull)
}

func (e *jsonEncDriver[T]) encodeIntAsUint(v int64, quotes bool) {
	neg := v < 0
	if neg {
		v = -v
	}
	e.encodeUint(neg, quotes, uint64(v))
}

func (e *jsonEncDriver[T]) EncodeTime(t time.Time) {
	// Do NOT use MarshalJSON, as it allocates internally.
	// instead, we call AppendFormat directly, using our scratch buffer (e.b)

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	switch e.timeFmt {
	case jsonTimeFmtStringLayout:
		e.b[0] = '"'
		b := t.AppendFormat(e.b[1:1], e.timeFmtLayout)
		e.b[len(b)+1] = '"'
		e.w.writeb(e.b[:len(b)+2])
	case jsonTimeFmtUnix:
		e.encodeIntAsUint(t.Unix(), false)
	case jsonTimeFmtUnixMilli:
		e.encodeIntAsUint(t.UnixMilli(), false)
	case jsonTimeFmtUnixMicro:
		e.encodeIntAsUint(t.UnixMicro(), false)
	case jsonTimeFmtUnixNano:
		e.encodeIntAsUint(t.UnixNano(), false)
	}
}

func (e *jsonEncDriver[T]) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.writeNilBytes()
	} else {
		e.enc.encodeI(v)
	}
}

func (e *jsonEncDriver[T]) EncodeRawExt(re *RawExt) {
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encodeI(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *jsonEncDriver[T]) EncodeBool(b bool) {
	e.w.writestr(jsonEncBoolStrs[bool2int(e.ks && e.e.c == containerMapKey)%2][bool2int(b)%2])
}

func (e *jsonEncDriver[T]) encodeFloat(f float64, bitsize, fmt byte, prec int8) {
	var blen uint
	if e.ks && e.e.c == containerMapKey {
		blen = 2 + uint(len(strconv.AppendFloat(e.b[1:1], f, fmt, int(prec), int(bitsize))))
		// _ = e.b[:blen]
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
	} else {
		e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), int(bitsize)))
	}
}

func (e *jsonEncDriver[T]) EncodeFloat64(f float64) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.encodeFloat(f, 64, fmt, prec)
}

func (e *jsonEncDriver[T]) EncodeFloat32(f float32) {
	if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.encodeFloat(float64(f), 32, fmt, prec)
}

func (e *jsonEncDriver[T]) encodeUint(neg bool, quotes bool, u uint64) {
	e.w.writeb(jsonEncodeUint(neg, quotes, u, &e.b))
}

func (e *jsonEncDriver[T]) EncodeInt(v int64) {
	quotes := e.is == 'A' || e.is == 'L' && (v > 1<<53 || v < -(1<<53)) ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		if quotes {
			blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
		}
		return
	}

	if v < 0 {
		e.encodeUint(true, quotes, uint64(-v))
	} else {
		e.encodeUint(false, quotes, uint64(v))
	}
}

func (e *jsonEncDriver[T]) EncodeUint(v uint64) {
	quotes := e.is == 'A' || e.is == 'L' && v > 1<<53 ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		// use strconv directly, as optimized encodeUint only works on 64-bit alone
		if quotes {
			blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
		}
		return
	}

	e.encodeUint(false, quotes, v)
}

func (e *jsonEncDriver[T]) EncodeString(v string) {
	if e.h.StringToRaw {
		e.EncodeStringBytesRaw(bytesView(v))
		return
	}
	e.quoteStr(v)
}

func (e *jsonEncDriver[T]) EncodeStringNoEscape4Json(v string) { e.w.writeqstr(v) }

func (e *jsonEncDriver[T]) EncodeStringBytesRaw(v []byte) {
	if e.rawext {
		// explicitly convert v to interface{} so that v doesn't escape to heap
		iv := e.h.RawBytesExt.ConvertExt(any(v))
		if iv == nil {
			e.EncodeNil()
		} else {
			e.enc.encodeI(iv)
		}
		return
	}

	if e.bytesFmt == jsonBytesFmtArray {
		e.WriteArrayStart(len(v))
		for j := range v {
			e.WriteArrayElem(j == 0)
			e.encodeUint(false, false, uint64(v[j]))
		}
		e.WriteArrayEnd()
		return
	}

	// hardcode base64, so we call direct (not via interface) and hopefully inline
	var slen int
	if e.bytesFmt == jsonBytesFmtBase64 {
		slen = base64.StdEncoding.EncodedLen(len(v))
	} else {
		slen = e.byteFmter.EncodedLen(len(v))
	}
	slen += 2

	// bs := e.e.blist.check(*e.buf, n)[:slen]
	// *e.buf = bs

	bs := e.e.blist.peek(slen, false)[:slen]

	if e.bytesFmt == jsonBytesFmtBase64 {
		base64.StdEncoding.Encode(bs[1:], v)
	} else {
		e.byteFmter.Encode(bs[1:], v)
	}

	bs[len(bs)-1] = '"'
	bs[0] = '"'
	e.w.writeb(bs)
}

func (e *jsonEncDriver[T]) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *jsonEncDriver[T]) writeNilOr(v []byte) {
	if !e.h.NilCollectionToZeroLength {
		v = jsonNull
	}
	e.w.writeb(v)
}

func (e *jsonEncDriver[T]) writeNilBytes() {
	e.writeNilOr(jsonArrayEmpty)
}

func (e *jsonEncDriver[T]) writeNilArray() {
	e.writeNilOr(jsonArrayEmpty)
}

func (e *jsonEncDriver[T]) writeNilMap() {
	e.writeNilOr(jsonMapEmpty)
}

// indent is done as below:
//   - newline and indent are added before each mapKey or arrayElem
//   - newline and indent are added before each ending,
//     except there was no entry (so we can have {} or [])

func (e *jsonEncDriver[T]) WriteArrayEmpty() {
	e.w.writen2('[', ']')
}

func (e *jsonEncDriver[T]) WriteMapEmpty() {
	e.w.writen2('{', '}')
}

func (e *jsonEncDriver[T]) WriteArrayStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('[')
}

func (e *jsonEncDriver[T]) WriteArrayEnd() {
	if e.d {
		e.dl--
		// No need as encoder handles zero-len already
		// if e.e.c != containerArrayStart {
		e.writeIndent()
	}
	e.w.writen1(']')
}

func (e *jsonEncDriver[T]) WriteMapStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('{')
}

func (e *jsonEncDriver[T]) WriteMapEnd() {
	if e.d {
		e.dl--
		// No need as encoder handles zero-len already
		// if e.e.c != containerMapStart {
		e.writeIndent()
	}
	e.w.writen1('}')
}

func (e *jsonEncDriver[T]) quoteStr(s string) {
	// adapted from std pkg encoding/json
	const hex = "0123456789abcdef"
	e.w.writen1('"')
	var i, start uint
	for i < uint(len(s)) {
		// encode all bytes < 0x20 (except \r, \n).
		// also encode < > & to prevent security holes when served to some browsers.

		// We optimize for ascii, by assuming that most characters are in the BMP
		// and natively consumed by json without much computation.

		// if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
		// if (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b) {
		b := s[i]
		if e.s.isset(b) {
			i++
			continue
		}
		if b < utf8.RuneSelf {
			if start < i {
				e.w.writestr(s[start:i])
			}
			switch b {
			case '\\':
				e.w.writen2('\\', '\\')
			case '"':
				e.w.writen2('\\', '"')
			case '\n':
				e.w.writen2('\\', 'n')
			case '\t':
				e.w.writen2('\\', 't')
			case '\r':
				e.w.writen2('\\', 'r')
			case '\b':
				e.w.writen2('\\', 'b')
			case '\f':
				e.w.writen2('\\', 'f')
			default:
				e.w.writestr(`\u00`)
				e.w.writen2(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 { // meaning invalid encoding (so output as-is)
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\uFFFD`)
			i++
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
		// Both technically valid JSON, but bomb on JSONP, so fix here *unconditionally*.
		if jsonEscapeMultiByteUnicodeSep && (c == '\u2028' || c == '\u2029') {
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\u202`)
			e.w.writen1(hex[c&0xF])
			i += uint(size)
			start = i
			continue
		}
		i += uint(size)
	}
	if start < uint(len(s)) {
		e.w.writestr(s[start:])
	}
	e.w.writen1('"')
}

func (e *jsonEncDriver[T]) atEndOfEncode() {
	if e.h.TermWhitespace {
		var c byte = ' ' // default is that scalar is written, so output space
		if e.e.c != 0 {
			c = '\n' // for containers (map/list), output a newline
		}
		e.w.writen1(c)
	}
}

// ----------

type jsonDecDriver[T decReader] struct {
	noBuiltInTypes
	decDriverNoopNumberHelper
	h *JsonHandle
	d *decoderBase

	r T

	// scratch buffer used for base64 decoding (DecodeBytes in reuseBuf mode),
	// or reading doubleQuoted string (DecodeStringAsBytes, DecodeNaked)
	buf []byte

	tok  uint8   // used to store the token read right after skipWhiteSpace
	_    bool    // found null
	_    byte    // padding
	bstr [4]byte // scratch used for string \UXXX parsing

	jsonHandleOpts

	// se  interfaceExtWrapper

	// ---- cpu cache line boundary?

	// bytes bool

	dec decoderI
}

func (d *jsonDecDriver[T]) ReadMapStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '{' {
		halt.errorByte("read map - expect char '{' but got char: ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriver[T]) ReadArrayStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '[' {
		halt.errorByte("read array - expect char '[' but got char ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

// MARKER:
// We attempted making sure CheckBreak can be inlined, by moving the skipWhitespace
// call to an explicit (noinline) function call.
// However, this forces CheckBreak to always incur a function call if there was whitespace,
// with no clear benefit.

func (d *jsonDecDriver[T]) CheckBreak() bool {
	d.advance()
	return d.tok == '}' || d.tok == ']'
}

func (d *jsonDecDriver[T]) checkSep(xc byte) {
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriver[T]) ReadArrayElem(firstTime bool) {
	if !firstTime {
		d.checkSep(',')
	}
}

func (d *jsonDecDriver[T]) ReadArrayEnd() {
	d.checkSep(']')
}

func (d *jsonDecDriver[T]) ReadMapElemKey(firstTime bool) {
	d.ReadArrayElem(firstTime)
}

func (d *jsonDecDriver[T]) ReadMapElemValue() {
	d.checkSep(':')
}

func (d *jsonDecDriver[T]) ReadMapEnd() {
	d.checkSep('}')
}

//go:inline
func (d *jsonDecDriver[T]) readDelimError(xc uint8) {
	halt.errorf("read json delimiter - expect char '%c' but got char '%c'", xc, d.tok)
}

// MARKER: checkLit takes the readn(3|4) result as a parameter so they can be inlined.
// We pass the array directly to errorf, as passing slice pushes past inlining threshold,
// and passing slice also might cause allocation of the bs array on the heap.

func (d *jsonDecDriver[T]) checkLit3(got, expect [3]byte) {
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr3(got, expect)
	}
	d.tok = 0
}

func (d *jsonDecDriver[T]) checkLit4(got, expect [4]byte) {
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr4(got, expect)
	}
	d.tok = 0
}

func (d *jsonDecDriver[T]) skipWhitespace() {
	d.tok = d.r.skipWhitespace()
}

func (d *jsonDecDriver[T]) advance() {
	// handles jsonReadNum returning possibly non-printable value as tok
	if d.tok < 33 { // d.tok == 0 {
		d.skipWhitespace()
	}
}

func (d *jsonDecDriver[T]) nextValueBytes() []byte {
	consumeString := func() {
	TOP:
		_, c := d.r.jsonReadAsisChars()
		if c == '\\' { // consume next one and try again
			d.r.readn1()
			goto TOP
		}
	}

	d.advance() // ignore leading whitespace
	d.r.startRecording()

	// cursor = d.d.rb.c - 1 // cursor starts just before non-whitespace token
	switch d.tok {
	default:
		_, d.tok = d.r.jsonReadNum()
		// special case: trim last read token if a valid byte in stream
		if d.tok != 0 {
			vv := d.r.stopRecording()
			return vv[:len(vv)-1]
		}
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
	case '"':
		consumeString()
		d.tok = 0
	case '{', '[':
		var elem struct{}
		var stack []struct{}

		stack = append(stack, elem)

		for len(stack) != 0 {
			c := d.r.readn1()
			switch c {
			case '"':
				consumeString()
			case '{', '[':
				stack = append(stack, elem)
			case '}', ']':
				stack = stack[:len(stack)-1]
			}
		}
		d.tok = 0
	}
	return d.r.stopRecording()
}

func (d *jsonDecDriver[T]) TryNil() bool {
	d.advance()
	// we don't try to see if quoted "null" was here.
	// only the plain string: null denotes a nil (ie not quotes)
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriver[T]) DecodeBool() (v bool) {
	d.advance()
	// bool can be in quotes if and only if it's a map key
	fquot := d.d.c == containerMapKey && d.tok == '"'
	if fquot {
		d.tok = d.r.readn1()
	}
	switch d.tok {
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		// v = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		v = true
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		// v = false
	default:
		halt.errorByte("decode bool: got first char: ", d.tok)
		// v = false // "unreachable"
	}
	if fquot {
		d.r.readn1()
	}
	return
}

func (d *jsonDecDriver[T]) DecodeTime() (t time.Time) {
	// read string, and pass the string into json.unmarshal
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	var bs []byte
	// if a number, use the timeFmtNum
	if d.tok != '"' {
		bs, d.tok = d.r.jsonReadNum()
		i := d.parseInt64(bs)
		switch d.timeFmtNum {
		case jsonTimeFmtUnix:
			t = time.Unix(i, 0)
		case jsonTimeFmtUnixMilli:
			t = time.UnixMilli(i)
		case jsonTimeFmtUnixMicro:
			t = time.UnixMicro(i)
		case jsonTimeFmtUnixNano:
			t = time.Unix(0, i)
		default:
			halt.errorStr("invalid timeFmtNum")
		}
		return
	}

	// d.tok is now '"'
	// d.ensureReadingString()
	bs = d.readUnescapedString()
	var err error
	for _, v := range d.timeFmtLayouts {
		t, err = time.Parse(v, stringView(bs))
		if err == nil {
			return
		}
	}
	halt.errorStr("error decoding time")
	return
}

func (d *jsonDecDriver[T]) ContainerType() (vt valueType) {
	// check container type by checking the first char
	d.advance()

	// optimize this, so we don't do 4 checks but do one computation.
	// return jsonContainerSet[d.tok]

	// ContainerType is mostly called for Map and Array,
	// so this conditional is good enough (max 2 checks typically)
	if d.tok == '{' {
		return valueTypeMap
	} else if d.tok == '[' {
		return valueTypeArray
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return valueTypeNil
	} else if d.tok == '"' {
		return valueTypeString
	}
	return valueTypeUnset
}

func (d *jsonDecDriver[T]) decNumBytes() (bs []byte) {
	d.advance()
	if d.tok == '"' {
		bs = d.r.jsonReadUntilDblQuote()
		d.tok = 0
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	} else {
		bs, d.tok = d.r.jsonReadNum()
	}
	return
}

func (d *jsonDecDriver[T]) DecodeUint64() (u uint64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if neg {
		halt.errorf("negative number cannot be decoded as uint64: %s", any(b))
	}
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseUint"))
	}
	return
}

func (d *jsonDecDriver[T]) DecodeInt64() (v int64) {
	return d.parseInt64(d.decNumBytes())
}

func (d *jsonDecDriver[T]) parseInt64(b []byte) (v int64) {
	u, neg, ok := parseInteger_bytes(b)
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseInt"))
	}
	if chkOvf.Uint2Int(u, neg) {
		halt.errorBytes("overflow decoding number from ", b)
	}
	if neg {
		v = -int64(u)
	} else {
		v = int64(u)
	}
	return
}

func (d *jsonDecDriver[T]) DecodeFloat64() (f float64) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat64(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriver[T]) DecodeFloat32() (f float32) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat32(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriver[T]) advanceNil() (ok bool) {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if d.advanceNil() {
		return
	}
	if ext == SelfExt {
		d.dec.decodeAs(rv, basetype, false)
	} else {
		d.dec.interfaceExtConvertAndDecode(rv, ext)
	}
}

func (d *jsonDecDriver[T]) DecodeRawExt(re *RawExt) {
	if d.advanceNil() {
		return
	}
	d.dec.decode(&re.Value)
}

func (d *jsonDecDriver[T]) decBytesFromArray(bs []byte) []byte {
	d.advance()
	if d.tok != ']' {
		bs = append(bs, uint8(d.DecodeUint64()))
		d.advance()
	}
	for d.tok != ']' {
		if d.tok != ',' {
			halt.errorByte("read array element - expect char ',' but got char: ", d.tok)
		}
		d.tok = 0
		bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		d.advance()
	}
	d.tok = 0
	return bs
}

func (d *jsonDecDriver[T]) DecodeBytes() (bs []byte, state dBytesAttachState) {
	d.advance()
	state = dBytesDetach
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	state = dBytesAttachBuffer
	// if decoding into raw bytes, and the RawBytesExt is configured, use it to decode.
	if d.rawext {
		d.buf = d.buf[:0]
		d.dec.interfaceExtConvertAndDecode(&d.buf, d.h.RawBytesExt)
		bs = d.buf
		return
	}
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.tok == '[' {
		d.tok = 0
		// bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		bs = d.decBytesFromArray(d.buf[:0])
		d.buf = bs
		return
	}

	// base64 encodes []byte{} as "", and we encode nil []byte as null.
	// Consequently, base64 should decode null as a nil []byte, and "" as an empty []byte{}.

	d.ensureReadingString()
	bs1 := d.readUnescapedString()
	// base64 is most compact of supported formats; it's decodedlen is sufficient for all
	slen := base64.StdEncoding.DecodedLen(len(bs1))
	if slen == 0 {
		bs = zeroByteSlice
		state = dBytesDetach
	} else if slen <= cap(d.buf) {
		bs = d.buf[:slen]
	} else {
		d.buf = d.d.blist.putGet(d.buf, slen)[:slen]
		bs = d.buf
	}
	var err error
	for _, v := range d.byteFmters {
		// slen := v.DecodedLen(len(bs1))
		slen, err = v.Decode(bs, bs1)
		if err == nil {
			bs = bs[:slen]
			return
		}
	}
	halt.errorf("error decoding byte string '%s': %v", any(bs1), err)
	return
}

func (d *jsonDecDriver[T]) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
	d.advance()

	var cond bool
	// common case - hoist outside the switch statement
	if d.tok == '"' {
		d.tok = 0
		bs, cond = d.dblQuoteStringAsBytes()
		state = d.d.attachState(cond)
		return
	}

	state = dBytesDetach
	// handle non-string scalar: null, true, false or a number
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		// out = nil // []byte{}
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		bs = jsonLitb[jsonLitF : jsonLitF+5]
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		bs = jsonLitb[jsonLitT : jsonLitT+4]
	default:
		// try to parse a valid number
		bs, d.tok = d.r.jsonReadNum()
		state = d.d.attachState(!d.d.bytes)
	}
	return
}

func (d *jsonDecDriver[T]) ensureReadingString() {
	if d.tok != '"' {
		halt.errorByte("expecting string starting with '\"'; got ", d.tok)
	}
}

func (d *jsonDecDriver[T]) readUnescapedString() (bs []byte) {
	// d.ensureReadingString()
	bs = d.r.jsonReadUntilDblQuote()
	d.tok = 0
	return
}

func (d *jsonDecDriver[T]) dblQuoteStringAsBytes() (buf []byte, usingBuf bool) {
	bs, c := d.r.jsonReadAsisChars()
	if c == '"' {
		return bs, !d.d.bytes
	}
	buf = append(d.buf[:0], bs...)

	checkUtf8 := d.h.ValidateUnicode
	usingBuf = true

	for {
		// c is now '\'
		c = d.r.readn1()

		switch c {
		case '"', '\\', '/', '\'':
			buf = append(buf, c)
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'u':
			rr := d.appendStringAsBytesSlashU()
			if checkUtf8 && rr == unicode.ReplacementChar {
				d.buf = buf
				halt.errorBytes("invalid UTF-8 character found after: ", buf)
			}
			buf = append(buf, d.bstr[:utf8.EncodeRune(d.bstr[:], rr)]...)
		default:
			d.buf = buf
			halt.errorByte("unsupported escaped value: ", c)
		}

		bs, c = d.r.jsonReadAsisChars()
		buf = append(buf, bs...)
		if c == '"' {
			break
		}
	}
	d.buf = buf
	return
}

func (d *jsonDecDriver[T]) appendStringAsBytesSlashU() (r rune) {
	var rr uint32
	cs := d.r.readn4()
	if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
		return unicode.ReplacementChar
	}
	r = rune(rr)
	if utf16.IsSurrogate(r) {
		csu := d.r.readn2()
		cs = d.r.readn4()
		if csu[0] == '\\' && csu[1] == 'u' {
			if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
				return unicode.ReplacementChar
			}
			return utf16.DecodeRune(r, rune(rr))
		}
		return unicode.ReplacementChar
	}
	return
}

func (d *jsonDecDriver[T]) DecodeNaked() {
	z := d.d.naked()

	d.advance()
	var bs []byte
	var err error
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		z.v = valueTypeNil
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		z.v = valueTypeBool
		z.b = true
	case '{':
		z.v = valueTypeMap // don't consume. kInterfaceNaked will call ReadMapStart
	case '[':
		z.v = valueTypeArray // don't consume. kInterfaceNaked will call ReadArrayStart
	case '"':
		// if a string, and MapKeyAsString, then try to decode it as a bool or number first
		d.tok = 0
		bs, z.b = d.dblQuoteStringAsBytes()
		att := d.d.attachState(z.b)
		if jsonNakedBoolNumInQuotedStr &&
			d.h.MapKeyAsString && len(bs) > 0 && d.d.c == containerMapKey {
			switch string(bs) {
			// case "null": // nil is never quoted
			// 	z.v = valueTypeNil
			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default: // check if a number: float, int or uint
				if err = jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
					z.v = valueTypeString
					z.s = d.d.detach2Str(bs, att)
				}
			}
		} else {
			z.v = valueTypeString
			z.s = d.d.detach2Str(bs, att)
		}
	default: // number
		bs, d.tok = d.r.jsonReadNum()
		if len(bs) == 0 {
			halt.errorStr("decode number from empty string")
		}
		if err = jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
			halt.errorf("decode number from %s: %v", any(bs), err)
		}
	}
}

func (e *jsonEncDriver[T]) reset() {
	e.dl = 0
	// e.resetState()
	// (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b)
	// cache values from the handle
	e.typical = e.h.typical()
	if e.h.HTMLCharsAsIs {
		e.s = &jsonCharSafeBitset
	} else {
		e.s = &jsonCharHtmlSafeBitset
	}
	e.di = int8(e.h.Indent)
	e.d = e.h.Indent != 0
	e.ks = e.h.MapKeyAsString
	e.is = e.h.IntegerAsString

	var ho jsonHandleOpts
	ho.reset(e.h)
	e.timeFmt = ho.timeFmt
	e.bytesFmt = ho.bytesFmt
	e.timeFmtLayout = ""
	e.byteFmter = nil
	if len(ho.timeFmtLayouts) > 0 {
		e.timeFmtLayout = ho.timeFmtLayouts[0]
	}
	if len(ho.byteFmters) > 0 {
		e.byteFmter = ho.byteFmters[0]
	}
	e.rawext = ho.rawext
}

func (d *jsonDecDriver[T]) reset() {
	d.buf = d.d.blist.check(d.buf, 256)
	d.tok = 0
	// d.resetState()
	d.jsonHandleOpts.reset(d.h)
}

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *jsonEncDriver[T]) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*JsonHandle)
	d.e = shared
	if shared.bytes {
		fp = jsonFpEncBytes
	} else {
		fp = jsonFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *jsonEncDriver[T]) writeBytesAsis(b []byte) { e.w.writeb(b) }

// func (e *jsonEncDriver[T]) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }
func (e *jsonEncDriver[T]) writerEnd() { e.w.end() }

func (e *jsonEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *jsonEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *jsonDecDriver[T]) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*JsonHandle)
	d.d = shared
	if shared.bytes {
		fp = jsonFpDecBytes
	} else {
		fp = jsonFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *jsonDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *jsonDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *jsonDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

// ---- (custom stanza)

func (d *jsonDecDriver[T]) descBd() (s string) {
	halt.onerror(errJsonNoBd)
	return
}

func (d *jsonEncDriver[T]) init2(enc encoderI) {
	d.enc = enc
	// d.e.js = true
}

func (d *jsonDecDriver[T]) init2(dec decoderI) {
	d.dec = dec
	// var x []byte
	// d.buf = &x
	// d.buf = new([]byte)
	d.buf = d.buf[:0]
	// d.d.js = true
	d.d.jsms = d.h.MapKeyAsString
}
