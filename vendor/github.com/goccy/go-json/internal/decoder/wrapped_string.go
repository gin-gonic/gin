package decoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goccy/go-json/internal/runtime"
)

type wrappedStringDecoder struct {
	typ           *runtime.Type
	dec           Decoder
	stringDecoder *stringDecoder
	structName    string
	fieldName     string
	isPtrType     bool
}

func newWrappedStringDecoder(typ *runtime.Type, dec Decoder, structName, fieldName string) *wrappedStringDecoder {
	return &wrappedStringDecoder{
		typ:           typ,
		dec:           dec,
		stringDecoder: newStringDecoder(structName, fieldName),
		structName:    structName,
		fieldName:     fieldName,
		isPtrType:     typ.Kind() == reflect.Ptr,
	}
}

func (d *wrappedStringDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	bytes, err := d.stringDecoder.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if bytes == nil {
		if d.isPtrType {
			*(*unsafe.Pointer)(p) = nil
		}
		return nil
	}
	b := make([]byte, len(bytes)+1)
	copy(b, bytes)
	if _, err := d.dec.Decode(&RuntimeContext{Buf: b}, 0, depth, p); err != nil {
		return err
	}
	return nil
}

func (d *wrappedStringDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.stringDecoder.decodeByte(ctx.Buf, cursor)
	if err != nil {
		return 0, err
	}
	if bytes == nil {
		if d.isPtrType {
			*(*unsafe.Pointer)(p) = nil
		}
		return c, nil
	}
	bytes = append(bytes, nul)
	oldBuf := ctx.Buf
	ctx.Buf = bytes
	if _, err := d.dec.Decode(ctx, 0, depth, p); err != nil {
		return 0, err
	}
	ctx.Buf = oldBuf
	return c, nil
}

func (d *wrappedStringDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return nil, 0, fmt.Errorf("json: wrapped string decoder does not support decode path")
}
