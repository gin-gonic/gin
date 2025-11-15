package decoder

import (
	"bytes"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
)

type stringDecoder struct {
	structName string
	fieldName  string
}

func newStringDecoder(structName, fieldName string) *stringDecoder {
	return &stringDecoder{
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *stringDecoder) errUnmarshalType(typeName string, offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  typeName,
		Type:   reflect.TypeOf(""),
		Offset: offset,
		Struct: d.structName,
		Field:  d.fieldName,
	}
}

func (d *stringDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	**(**string)(unsafe.Pointer(&p)) = *(*string)(unsafe.Pointer(&bytes))
	s.reset()
	return nil
}

func (d *stringDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	**(**string)(unsafe.Pointer(&p)) = *(*string)(unsafe.Pointer(&bytes))
	return cursor, nil
}

func (d *stringDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return nil, 0, err
	}
	if bytes == nil {
		return [][]byte{nullbytes}, c, nil
	}
	return [][]byte{bytes}, c, nil
}

var (
	hexToInt = [256]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
	}
)

func unicodeToRune(code []byte) rune {
	var r rune
	for i := 0; i < len(code); i++ {
		r = r*16 + rune(hexToInt[code[i]])
	}
	return r
}

func readAtLeast(s *Stream, n int64, p *unsafe.Pointer) bool {
	for s.cursor+n >= s.length {
		if !s.read() {
			return false
		}
		*p = s.bufptr()
	}
	return true
}

func decodeUnicodeRune(s *Stream, p unsafe.Pointer) (rune, int64, unsafe.Pointer, error) {
	const defaultOffset = 5
	const surrogateOffset = 11

	if !readAtLeast(s, defaultOffset, &p) {
		return rune(0), 0, nil, errors.ErrInvalidCharacter(s.char(), "escaped string", s.totalOffset())
	}

	r := unicodeToRune(s.buf[s.cursor+1 : s.cursor+defaultOffset])
	if utf16.IsSurrogate(r) {
		if !readAtLeast(s, surrogateOffset, &p) {
			return unicode.ReplacementChar, defaultOffset, p, nil
		}
		if s.buf[s.cursor+defaultOffset] != '\\' || s.buf[s.cursor+defaultOffset+1] != 'u' {
			return unicode.ReplacementChar, defaultOffset, p, nil
		}
		r2 := unicodeToRune(s.buf[s.cursor+defaultOffset+2 : s.cursor+surrogateOffset])
		if r := utf16.DecodeRune(r, r2); r != unicode.ReplacementChar {
			return r, surrogateOffset, p, nil
		}
	}
	return r, defaultOffset, p, nil
}

func decodeUnicode(s *Stream, p unsafe.Pointer) (unsafe.Pointer, error) {
	const backSlashAndULen = 2 // length of \u

	r, offset, pp, err := decodeUnicodeRune(s, p)
	if err != nil {
		return nil, err
	}
	unicode := []byte(string(r))
	unicodeLen := int64(len(unicode))
	s.buf = append(append(s.buf[:s.cursor-1], unicode...), s.buf[s.cursor+offset:]...)
	unicodeOrgLen := offset - 1
	s.length = s.length - (backSlashAndULen + (unicodeOrgLen - unicodeLen))
	s.cursor = s.cursor - backSlashAndULen + unicodeLen
	return pp, nil
}

func decodeEscapeString(s *Stream, p unsafe.Pointer) (unsafe.Pointer, error) {
	s.cursor++
RETRY:
	switch s.buf[s.cursor] {
	case '"':
		s.buf[s.cursor] = '"'
	case '\\':
		s.buf[s.cursor] = '\\'
	case '/':
		s.buf[s.cursor] = '/'
	case 'b':
		s.buf[s.cursor] = '\b'
	case 'f':
		s.buf[s.cursor] = '\f'
	case 'n':
		s.buf[s.cursor] = '\n'
	case 'r':
		s.buf[s.cursor] = '\r'
	case 't':
		s.buf[s.cursor] = '\t'
	case 'u':
		return decodeUnicode(s, p)
	case nul:
		if !s.read() {
			return nil, errors.ErrInvalidCharacter(s.char(), "escaped string", s.totalOffset())
		}
		p = s.bufptr()
		goto RETRY
	default:
		return nil, errors.ErrUnexpectedEndOfJSON("string", s.totalOffset())
	}
	s.buf = append(s.buf[:s.cursor-1], s.buf[s.cursor:]...)
	s.length--
	s.cursor--
	p = s.bufptr()
	return p, nil
}

var (
	runeErrBytes    = []byte(string(utf8.RuneError))
	runeErrBytesLen = int64(len(runeErrBytes))
)

func stringBytes(s *Stream) ([]byte, error) {
	_, cursor, p := s.stat()
	cursor++ // skip double quote char
	start := cursor
	for {
		switch char(p, cursor) {
		case '\\':
			s.cursor = cursor
			pp, err := decodeEscapeString(s, p)
			if err != nil {
				return nil, err
			}
			p = pp
			cursor = s.cursor
		case '"':
			literal := s.buf[start:cursor]
			cursor++
			s.cursor = cursor
			return literal, nil
		case
			// 0x00 is nul, 0x5c is '\\', 0x22 is '"' .
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, // 0x00-0x0F
			0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, // 0x10-0x1F
			0x20, 0x21 /*0x22,*/, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, // 0x20-0x2F
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, // 0x30-0x3F
			0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, // 0x40-0x4F
			0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B /*0x5C,*/, 0x5D, 0x5E, 0x5F, // 0x50-0x5F
			0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, // 0x60-0x6F
			0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F: // 0x70-0x7F
			// character is ASCII. skip to next char
		case
			0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, // 0x80-0x8F
			0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F, // 0x90-0x9F
			0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, // 0xA0-0xAF
			0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF, // 0xB0-0xBF
			0xC0, 0xC1, // 0xC0-0xC1
			0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF: // 0xF5-0xFE
			// character is invalid
			s.buf = append(append(append([]byte{}, s.buf[:cursor]...), runeErrBytes...), s.buf[cursor+1:]...)
			_, _, p = s.stat()
			cursor += runeErrBytesLen
			s.length += runeErrBytesLen
			continue
		case nul:
			s.cursor = cursor
			if s.read() {
				_, cursor, p = s.stat()
				continue
			}
			goto ERROR
		case 0xEF:
			// RuneError is {0xEF, 0xBF, 0xBD}
			if s.buf[cursor+1] == 0xBF && s.buf[cursor+2] == 0xBD {
				// found RuneError: skip
				cursor += 2
				break
			}
			fallthrough
		default:
			// multi bytes character
			if !utf8.FullRune(s.buf[cursor : len(s.buf)-1]) {
				s.cursor = cursor
				if s.read() {
					_, cursor, p = s.stat()
					continue
				}
				goto ERROR
			}
			r, size := utf8.DecodeRune(s.buf[cursor:])
			if r == utf8.RuneError {
				s.buf = append(append(append([]byte{}, s.buf[:cursor]...), runeErrBytes...), s.buf[cursor+1:]...)
				cursor += runeErrBytesLen
				s.length += runeErrBytesLen
				_, _, p = s.stat()
			} else {
				cursor += int64(size)
			}
			continue
		}
		cursor++
	}
ERROR:
	return nil, errors.ErrUnexpectedEndOfJSON("string", s.totalOffset())
}

func (d *stringDecoder) decodeStreamByte(s *Stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '[':
			return nil, d.errUnmarshalType("array", s.totalOffset())
		case '{':
			return nil, d.errUnmarshalType("object", s.totalOffset())
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return nil, d.errUnmarshalType("number", s.totalOffset())
		case '"':
			return stringBytes(s)
		case 'n':
			if err := nullBytes(s); err != nil {
				return nil, err
			}
			return nil, nil
		case nul:
			if s.read() {
				continue
			}
		}
		break
	}
	return nil, errors.ErrInvalidBeginningOfValue(s.char(), s.totalOffset())
}

func (d *stringDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
		case '[':
			return nil, 0, d.errUnmarshalType("array", cursor)
		case '{':
			return nil, 0, d.errUnmarshalType("object", cursor)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return nil, 0, d.errUnmarshalType("number", cursor)
		case '"':
			cursor++
			start := cursor
			b := (*sliceHeader)(unsafe.Pointer(&buf)).data
			escaped := 0
			for {
				switch char(b, cursor) {
				case '\\':
					escaped++
					cursor++
					switch char(b, cursor) {
					case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
						cursor++
					case 'u':
						buflen := int64(len(buf))
						if cursor+5 >= buflen {
							return nil, 0, errors.ErrUnexpectedEndOfJSON("escaped string", cursor)
						}
						for i := int64(1); i <= 4; i++ {
							c := char(b, cursor+i)
							if !(('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')) {
								return nil, 0, errors.ErrSyntax(fmt.Sprintf("json: invalid character %c in \\u hexadecimal character escape", c), cursor+i)
							}
						}
						cursor += 5
					default:
						return nil, 0, errors.ErrUnexpectedEndOfJSON("escaped string", cursor)
					}
					continue
				case '"':
					literal := buf[start:cursor]
					if escaped > 0 {
						literal = literal[:unescapeString(literal)]
					}
					cursor++
					return literal, cursor, nil
				case nul:
					return nil, 0, errors.ErrUnexpectedEndOfJSON("string", cursor)
				}
				cursor++
			}
		case 'n':
			if err := validateNull(buf, cursor); err != nil {
				return nil, 0, err
			}
			cursor += 4
			return nil, cursor, nil
		default:
			return nil, 0, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
		}
	}
}

var unescapeMap = [256]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func unsafeAdd(ptr unsafe.Pointer, offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + uintptr(offset))
}

func unescapeString(buf []byte) int {
	p := (*sliceHeader)(unsafe.Pointer(&buf)).data
	end := unsafeAdd(p, len(buf))
	src := unsafeAdd(p, bytes.IndexByte(buf, '\\'))
	dst := src
	for src != end {
		c := char(src, 0)
		if c == '\\' {
			escapeChar := char(src, 1)
			if escapeChar != 'u' {
				*(*byte)(dst) = unescapeMap[escapeChar]
				src = unsafeAdd(src, 2)
				dst = unsafeAdd(dst, 1)
			} else {
				v1 := hexToInt[char(src, 2)]
				v2 := hexToInt[char(src, 3)]
				v3 := hexToInt[char(src, 4)]
				v4 := hexToInt[char(src, 5)]
				code := rune((v1 << 12) | (v2 << 8) | (v3 << 4) | v4)
				if code >= 0xd800 && code < 0xdc00 && uintptr(unsafeAdd(src, 11)) < uintptr(end) {
					if char(src, 6) == '\\' && char(src, 7) == 'u' {
						v1 := hexToInt[char(src, 8)]
						v2 := hexToInt[char(src, 9)]
						v3 := hexToInt[char(src, 10)]
						v4 := hexToInt[char(src, 11)]
						lo := rune((v1 << 12) | (v2 << 8) | (v3 << 4) | v4)
						if lo >= 0xdc00 && lo < 0xe000 {
							code = (code-0xd800)<<10 | (lo - 0xdc00) + 0x10000
							src = unsafeAdd(src, 6)
						}
					}
				}
				var b [utf8.UTFMax]byte
				n := utf8.EncodeRune(b[:], code)
				switch n {
				case 4:
					*(*byte)(unsafeAdd(dst, 3)) = b[3]
					fallthrough
				case 3:
					*(*byte)(unsafeAdd(dst, 2)) = b[2]
					fallthrough
				case 2:
					*(*byte)(unsafeAdd(dst, 1)) = b[1]
					fallthrough
				case 1:
					*(*byte)(unsafeAdd(dst, 0)) = b[0]
				}
				src = unsafeAdd(src, 6)
				dst = unsafeAdd(dst, n)
			}
		} else {
			*(*byte)(dst) = c
			src = unsafeAdd(src, 1)
			dst = unsafeAdd(dst, 1)
		}
	}
	return int(uintptr(dst) - uintptr(p))
}
