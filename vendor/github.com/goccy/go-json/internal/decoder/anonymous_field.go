package decoder

import (
	"unsafe"

	"github.com/goccy/go-json/internal/runtime"
)

type anonymousFieldDecoder struct {
	structType *runtime.Type
	offset     uintptr
	dec        Decoder
}

func newAnonymousFieldDecoder(structType *runtime.Type, offset uintptr, dec Decoder) *anonymousFieldDecoder {
	return &anonymousFieldDecoder{
		structType: structType,
		offset:     offset,
		dec:        dec,
	}
}

func (d *anonymousFieldDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	if *(*unsafe.Pointer)(p) == nil {
		*(*unsafe.Pointer)(p) = unsafe_New(d.structType)
	}
	p = *(*unsafe.Pointer)(p)
	return d.dec.DecodeStream(s, depth, unsafe.Pointer(uintptr(p)+d.offset))
}

func (d *anonymousFieldDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	if *(*unsafe.Pointer)(p) == nil {
		*(*unsafe.Pointer)(p) = unsafe_New(d.structType)
	}
	p = *(*unsafe.Pointer)(p)
	return d.dec.Decode(ctx, cursor, depth, unsafe.Pointer(uintptr(p)+d.offset))
}

func (d *anonymousFieldDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return d.dec.DecodePath(ctx, cursor, depth)
}
