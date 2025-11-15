package decoder

import (
	"fmt"
	"unsafe"

	"github.com/goccy/go-json/internal/runtime"
)

type ptrDecoder struct {
	dec        Decoder
	typ        *runtime.Type
	structName string
	fieldName  string
}

func newPtrDecoder(dec Decoder, typ *runtime.Type, structName, fieldName string) *ptrDecoder {
	return &ptrDecoder{
		dec:        dec,
		typ:        typ,
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *ptrDecoder) contentDecoder() Decoder {
	dec, ok := d.dec.(*ptrDecoder)
	if !ok {
		return d.dec
	}
	return dec.contentDecoder()
}

//nolint:golint
//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(*runtime.Type) unsafe.Pointer

func UnsafeNew(t *runtime.Type) unsafe.Pointer {
	return unsafe_New(t)
}

func (d *ptrDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	if s.skipWhiteSpace() == nul {
		s.read()
	}
	if s.char() == 'n' {
		if err := nullBytes(s); err != nil {
			return err
		}
		*(*unsafe.Pointer)(p) = nil
		return nil
	}
	var newptr unsafe.Pointer
	if *(*unsafe.Pointer)(p) == nil {
		newptr = unsafe_New(d.typ)
		*(*unsafe.Pointer)(p) = newptr
	} else {
		newptr = *(*unsafe.Pointer)(p)
	}
	if err := d.dec.DecodeStream(s, depth, newptr); err != nil {
		return err
	}
	return nil
}

func (d *ptrDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	cursor = skipWhiteSpace(buf, cursor)
	if buf[cursor] == 'n' {
		if err := validateNull(buf, cursor); err != nil {
			return 0, err
		}
		if p != nil {
			*(*unsafe.Pointer)(p) = nil
		}
		cursor += 4
		return cursor, nil
	}
	var newptr unsafe.Pointer
	if *(*unsafe.Pointer)(p) == nil {
		newptr = unsafe_New(d.typ)
		*(*unsafe.Pointer)(p) = newptr
	} else {
		newptr = *(*unsafe.Pointer)(p)
	}
	c, err := d.dec.Decode(ctx, cursor, depth, newptr)
	if err != nil {
		return 0, err
	}
	cursor = c
	return cursor, nil
}

func (d *ptrDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return nil, 0, fmt.Errorf("json: ptr decoder does not support decode path")
}
