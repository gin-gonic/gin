package decoder

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

type funcDecoder struct {
	typ        *runtime.Type
	structName string
	fieldName  string
}

func newFuncDecoder(typ *runtime.Type, structName, fieldName string) *funcDecoder {
	fnDecoder := &funcDecoder{typ, structName, fieldName}
	return fnDecoder
}

func (d *funcDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	s.skipWhiteSpace()
	start := s.cursor
	if err := s.skipValue(depth); err != nil {
		return err
	}
	src := s.buf[start:s.cursor]
	if len(src) > 0 {
		switch src[0] {
		case '"':
			return &errors.UnmarshalTypeError{
				Value:  "string",
				Type:   runtime.RType2Type(d.typ),
				Offset: s.totalOffset(),
			}
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
			if err := nullBytes(s); err != nil {
				return err
			}
			*(*unsafe.Pointer)(p) = nil
			return nil
		case 't':
			if err := trueBytes(s); err == nil {
				return &errors.UnmarshalTypeError{
					Value:  "boolean",
					Type:   runtime.RType2Type(d.typ),
					Offset: s.totalOffset(),
				}
			}
		case 'f':
			if err := falseBytes(s); err == nil {
				return &errors.UnmarshalTypeError{
					Value:  "boolean",
					Type:   runtime.RType2Type(d.typ),
					Offset: s.totalOffset(),
				}
			}
		}
	}
	return errors.ErrInvalidBeginningOfValue(s.buf[s.cursor], s.totalOffset())
}

func (d *funcDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
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
		case '"':
			return 0, &errors.UnmarshalTypeError{
				Value:  "string",
				Type:   runtime.RType2Type(d.typ),
				Offset: start,
			}
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
		case 't':
			if err := validateTrue(buf, start); err == nil {
				return 0, &errors.UnmarshalTypeError{
					Value:  "boolean",
					Type:   runtime.RType2Type(d.typ),
					Offset: start,
				}
			}
		case 'f':
			if err := validateFalse(buf, start); err == nil {
				return 0, &errors.UnmarshalTypeError{
					Value:  "boolean",
					Type:   runtime.RType2Type(d.typ),
					Offset: start,
				}
			}
		}
	}
	return cursor, errors.ErrInvalidBeginningOfValue(buf[cursor], cursor)
}

func (d *funcDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return nil, 0, fmt.Errorf("json: func decoder does not support decode path")
}
