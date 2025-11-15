package decoder

import (
	"strconv"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
)

type floatDecoder struct {
	op         func(unsafe.Pointer, float64)
	structName string
	fieldName  string
}

func newFloatDecoder(structName, fieldName string, op func(unsafe.Pointer, float64)) *floatDecoder {
	return &floatDecoder{op: op, structName: structName, fieldName: fieldName}
}

var (
	floatTable = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'.': true,
		'e': true,
		'E': true,
		'+': true,
		'-': true,
	}

	validEndNumberChar = [256]bool{
		nul:  true,
		' ':  true,
		'\t': true,
		'\r': true,
		'\n': true,
		',':  true,
		':':  true,
		'}':  true,
		']':  true,
	}
)

func floatBytes(s *Stream) []byte {
	start := s.cursor
	for {
		s.cursor++
		if floatTable[s.char()] {
			continue
		} else if s.char() == nul {
			if s.read() {
				s.cursor-- // for retry current character
				continue
			}
		}
		break
	}
	return s.buf[start:s.cursor]
}

func (d *floatDecoder) decodeStreamByte(s *Stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return floatBytes(s), nil
		case 'n':
			if err := nullBytes(s); err != nil {
				return nil, err
			}
			return nil, nil
		case nul:
			if s.read() {
				continue
			}
			goto ERROR
		default:
			goto ERROR
		}
	}
ERROR:
	return nil, errors.ErrUnexpectedEndOfJSON("float", s.totalOffset())
}

func (d *floatDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := cursor
			cursor++
			for floatTable[buf[cursor]] {
				cursor++
			}
			num := buf[start:cursor]
			return num, cursor, nil
		case 'n':
			if err := validateNull(buf, cursor); err != nil {
				return nil, 0, err
			}
			cursor += 4
			return nil, cursor, nil
		default:
			return nil, 0, errors.ErrUnexpectedEndOfJSON("float", cursor)
		}
	}
}

func (d *floatDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	str := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return errors.ErrSyntax(err.Error(), s.totalOffset())
	}
	d.op(p, f64)
	return nil
}

func (d *floatDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	if !validEndNumberChar[buf[cursor]] {
		return 0, errors.ErrUnexpectedEndOfJSON("float", cursor)
	}
	s := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.ErrSyntax(err.Error(), cursor)
	}
	d.op(p, f64)
	return cursor, nil
}

func (d *floatDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	buf := ctx.Buf
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return nil, 0, err
	}
	if bytes == nil {
		return [][]byte{nullbytes}, c, nil
	}
	return [][]byte{bytes}, c, nil
}
