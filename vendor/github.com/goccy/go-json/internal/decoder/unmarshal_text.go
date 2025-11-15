package decoder

import (
	"bytes"
	"encoding"
	"fmt"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

type unmarshalTextDecoder struct {
	typ        *runtime.Type
	structName string
	fieldName  string
}

func newUnmarshalTextDecoder(typ *runtime.Type, structName, fieldName string) *unmarshalTextDecoder {
	return &unmarshalTextDecoder{
		typ:        typ,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *unmarshalTextDecoder) annotateError(cursor int64, err error) {
	switch e := err.(type) {
	case *errors.UnmarshalTypeError:
		e.Struct = d.structName
		e.Field = d.fieldName
	case *errors.SyntaxError:
		e.Offset = cursor
	}
}

var (
	nullbytes = []byte(`null`)
)

func (d *unmarshalTextDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	s.skipWhiteSpace()
	start := s.cursor
	if err := s.skipValue(depth); err != nil {
		return err
	}
	src := s.buf[start:s.cursor]
	if len(src) > 0 {
		switch src[0] {
		case '[':
			return &errors.UnmarshalTypeError{
				Value:  "array",
				Type:   runtime.RType2Type(d.typ),
				Offset: s.totalOffset(),
			}
		case '{':
			return &errors.UnmarshalTypeError{
				Value:  "object",
				Type:   runtime.RType2Type(d.typ),
				Offset: s.totalOffset(),
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return &errors.UnmarshalTypeError{
				Value:  "number",
				Type:   runtime.RType2Type(d.typ),
				Offset: s.totalOffset(),
			}
		case 'n':
			if bytes.Equal(src, nullbytes) {
				*(*unsafe.Pointer)(p) = nil
				return nil
			}
		}
	}
	dst := make([]byte, len(src))
	copy(dst, src)

	if b, ok := unquoteBytes(dst); ok {
		dst = b
	}
	v := *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: d.typ,
		ptr: p,
	}))
	if err := v.(encoding.TextUnmarshaler).UnmarshalText(dst); err != nil {
		d.annotateError(s.cursor, err)
		return err
	}
	return nil
}

func (d *unmarshalTextDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	cursor = skipWhiteSpace(buf, cursor)
	start := cursor
	end, err := skipValue(buf, cursor, depth)
	if err != nil {
		return 0, err
	}
	src := buf[start:end]
	if len(src) > 0 {
		switch src[0] {
		case '[':
			return 0, &errors.UnmarshalTypeError{
				Value:  "array",
				Type:   runtime.RType2Type(d.typ),
				Offset: start,
			}
		case '{':
			return 0, &errors.UnmarshalTypeError{
				Value:  "object",
				Type:   runtime.RType2Type(d.typ),
				Offset: start,
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, &errors.UnmarshalTypeError{
				Value:  "number",
				Type:   runtime.RType2Type(d.typ),
				Offset: start,
			}
		case 'n':
			if bytes.Equal(src, nullbytes) {
				*(*unsafe.Pointer)(p) = nil
				return end, nil
			}
		}
	}

	if s, ok := unquoteBytes(src); ok {
		src = s
	}
	v := *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: d.typ,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
	if err := v.(encoding.TextUnmarshaler).UnmarshalText(src); err != nil {
		d.annotateError(cursor, err)
		return 0, err
	}
	return end, nil
}

func (d *unmarshalTextDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return nil, 0, fmt.Errorf("json: unmarshal text decoder does not support decode path")
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	length := len(s)
	if length < 2 || s[0] != '"' || s[length-1] != '"' {
		return
	}
	s = s[1 : length-1]
	length -= 2

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < length {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == length {
		return s, true
	}

	b := make([]byte, length+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < length {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= length {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}
