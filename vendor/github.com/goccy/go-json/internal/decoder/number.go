package decoder

import (
	"encoding/json"
	"strconv"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
)

type numberDecoder struct {
	stringDecoder *stringDecoder
	op            func(unsafe.Pointer, json.Number)
	structName    string
	fieldName     string
}

func newNumberDecoder(structName, fieldName string, op func(unsafe.Pointer, json.Number)) *numberDecoder {
	return &numberDecoder{
		stringDecoder: newStringDecoder(structName, fieldName),
		op:            op,
		structName:    structName,
		fieldName:     fieldName,
	}
}

func (d *numberDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if _, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&bytes)), 64); err != nil {
		return errors.ErrSyntax(err.Error(), s.totalOffset())
	}
	d.op(p, json.Number(string(bytes)))
	s.reset()
	return nil
}

func (d *numberDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if _, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&bytes)), 64); err != nil {
		return 0, errors.ErrSyntax(err.Error(), c)
	}
	cursor = c
	s := *(*string)(unsafe.Pointer(&bytes))
	d.op(p, json.Number(s))
	return cursor, nil
}

func (d *numberDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	bytes, c, err := d.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return nil, 0, err
	}
	if bytes == nil {
		return [][]byte{nullbytes}, c, nil
	}
	return [][]byte{bytes}, c, nil
}

func (d *numberDecoder) decodeStreamByte(s *Stream) ([]byte, error) {
	start := s.cursor
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
		case '"':
			return d.stringDecoder.decodeStreamByte(s)
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
	if s.cursor == start {
		return nil, errors.ErrInvalidBeginningOfValue(s.char(), s.totalOffset())
	}
	return nil, errors.ErrUnexpectedEndOfJSON("json.Number", s.totalOffset())
}

func (d *numberDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
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
		case '"':
			return d.stringDecoder.decodeByte(buf, cursor)
		default:
			return nil, 0, errors.ErrUnexpectedEndOfJSON("json.Number", cursor)
		}
	}
}
