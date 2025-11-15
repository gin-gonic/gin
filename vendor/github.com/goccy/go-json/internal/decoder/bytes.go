package decoder

import (
	"encoding/base64"
	"fmt"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

type bytesDecoder struct {
	typ           *runtime.Type
	sliceDecoder  Decoder
	stringDecoder *stringDecoder
	structName    string
	fieldName     string
}

func byteUnmarshalerSliceDecoder(typ *runtime.Type, structName string, fieldName string) Decoder {
	var unmarshalDecoder Decoder
	switch {
	case runtime.PtrTo(typ).Implements(unmarshalJSONType):
		unmarshalDecoder = newUnmarshalJSONDecoder(runtime.PtrTo(typ), structName, fieldName)
	case runtime.PtrTo(typ).Implements(unmarshalTextType):
		unmarshalDecoder = newUnmarshalTextDecoder(runtime.PtrTo(typ), structName, fieldName)
	default:
		unmarshalDecoder, _ = compileUint8(typ, structName, fieldName)
	}
	return newSliceDecoder(unmarshalDecoder, typ, 1, structName, fieldName)
}

func newBytesDecoder(typ *runtime.Type, structName string, fieldName string) *bytesDecoder {
	return &bytesDecoder{
		typ:           typ,
		sliceDecoder:  byteUnmarshalerSliceDecoder(typ, structName, fieldName),
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
	}
}

func (d *bytesDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamBinary(s, depth, p)
	if err != nil {
		return err
	}
	if bytes == nil {
		s.reset()
		return nil
	}
	decodedLen := base64.StdEncoding.DecodedLen(len(bytes))
	buf := make([]byte, decodedLen)
	n, err := base64.StdEncoding.Decode(buf, bytes)
	if err != nil {
		return err
	}
	*(*[]byte)(p) = buf[:n]
	s.reset()
	return nil
}

func (d *bytesDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeBinary(ctx, cursor, depth, p)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		return c, nil
	}
	cursor = c
	decodedLen := base64.StdEncoding.DecodedLen(len(bytes))
	b := make([]byte, decodedLen)
	n, err := base64.StdEncoding.Decode(b, bytes)
	if err != nil {
		return 0, err
	}
	*(*[]byte)(p) = b[:n]
	return cursor, nil
}

func (d *bytesDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return nil, 0, fmt.Errorf("json: []byte decoder does not support decode path")
}

func (d *bytesDecoder) decodeStreamBinary(s *Stream, depth int64, p unsafe.Pointer) ([]byte, error) {
	c := s.skipWhiteSpace()
	if c == '[' {
		if d.sliceDecoder == nil {
			return nil, &errors.UnmarshalTypeError{
				Type:   runtime.RType2Type(d.typ),
				Offset: s.totalOffset(),
			}
		}
		err := d.sliceDecoder.DecodeStream(s, depth, p)
		return nil, err
	}
	return d.stringDecoder.decodeStreamByte(s)
}

func (d *bytesDecoder) decodeBinary(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) ([]byte, int64, error) {
	buf := ctx.Buf
	cursor = skipWhiteSpace(buf, cursor)
	if buf[cursor] == '[' {
		if d.sliceDecoder == nil {
			return nil, 0, &errors.UnmarshalTypeError{
				Type:   runtime.RType2Type(d.typ),
				Offset: cursor,
			}
		}
		c, err := d.sliceDecoder.Decode(ctx, cursor, depth, p)
		if err != nil {
			return nil, 0, err
		}
		return nil, c, nil
	}
	return d.stringDecoder.decodeByte(buf, cursor)
}
