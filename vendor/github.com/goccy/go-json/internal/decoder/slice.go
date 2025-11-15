package decoder

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

var (
	sliceType = runtime.Type2RType(
		reflect.TypeOf((*sliceHeader)(nil)).Elem(),
	)
	nilSlice = unsafe.Pointer(&sliceHeader{})
)

type sliceDecoder struct {
	elemType          *runtime.Type
	isElemPointerType bool
	valueDecoder      Decoder
	size              uintptr
	arrayPool         sync.Pool
	structName        string
	fieldName         string
}

// If use reflect.SliceHeader, data type is uintptr.
// In this case, Go compiler cannot trace reference created by newArray().
// So, define using unsafe.Pointer as data type
type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

const (
	defaultSliceCapacity = 2
)

func newSliceDecoder(dec Decoder, elemType *runtime.Type, size uintptr, structName, fieldName string) *sliceDecoder {
	return &sliceDecoder{
		valueDecoder:      dec,
		elemType:          elemType,
		isElemPointerType: elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Map,
		size:              size,
		arrayPool: sync.Pool{
			New: func() interface{} {
				return &sliceHeader{
					data: newArray(elemType, defaultSliceCapacity),
					len:  0,
					cap:  defaultSliceCapacity,
				}
			},
		},
		structName: structName,
		fieldName:  fieldName,
	}
}

func (d *sliceDecoder) newSlice(src *sliceHeader) *sliceHeader {
	slice := d.arrayPool.Get().(*sliceHeader)
	if src.len > 0 {
		// copy original elem
		if slice.cap < src.cap {
			data := newArray(d.elemType, src.cap)
			slice = &sliceHeader{data: data, len: src.len, cap: src.cap}
		} else {
			slice.len = src.len
		}
		copySlice(d.elemType, *slice, *src)
	} else {
		slice.len = 0
	}
	return slice
}

func (d *sliceDecoder) releaseSlice(p *sliceHeader) {
	d.arrayPool.Put(p)
}

//go:linkname copySlice reflect.typedslicecopy
func copySlice(elemType *runtime.Type, dst, src sliceHeader) int

//go:linkname newArray reflect.unsafe_NewArray
func newArray(*runtime.Type, int) unsafe.Pointer

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t *runtime.Type, dst, src unsafe.Pointer)

func (d *sliceDecoder) errNumber(offset int64) *errors.UnmarshalTypeError {
	return &errors.UnmarshalTypeError{
		Value:  "number",
		Type:   reflect.SliceOf(runtime.RType2Type(d.elemType)),
		Struct: d.structName,
		Field:  d.fieldName,
		Offset: offset,
	}
}

func (d *sliceDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	depth++
	if depth > maxDecodeNestingDepth {
		return errors.ErrExceededMaxDepth(s.char(), s.cursor)
	}

	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case 'n':
			if err := nullBytes(s); err != nil {
				return err
			}
			typedmemmove(sliceType, p, nilSlice)
			return nil
		case '[':
			s.cursor++
			if s.skipWhiteSpace() == ']' {
				dst := (*sliceHeader)(p)
				if dst.data == nil {
					dst.data = newArray(d.elemType, 0)
				} else {
					dst.len = 0
				}
				s.cursor++
				return nil
			}
			idx := 0
			slice := d.newSlice((*sliceHeader)(p))
			srcLen := slice.len
			capacity := slice.cap
			data := slice.data
			for {
				if capacity <= idx {
					src := sliceHeader{data: data, len: idx, cap: capacity}
					capacity *= 2
					data = newArray(d.elemType, capacity)
					dst := sliceHeader{data: data, len: idx, cap: capacity}
					copySlice(d.elemType, dst, src)
				}
				ep := unsafe.Pointer(uintptr(data) + uintptr(idx)*d.size)

				// if srcLen is greater than idx, keep the original reference
				if srcLen <= idx {
					if d.isElemPointerType {
						**(**unsafe.Pointer)(unsafe.Pointer(&ep)) = nil // initialize elem pointer
					} else {
						// assign new element to the slice
						typedmemmove(d.elemType, ep, unsafe_New(d.elemType))
					}
				}

				if err := d.valueDecoder.DecodeStream(s, depth, ep); err != nil {
					return err
				}
				s.skipWhiteSpace()
			RETRY:
				switch s.char() {
				case ']':
					slice.cap = capacity
					slice.len = idx + 1
					slice.data = data
					dst := (*sliceHeader)(p)
					dst.len = idx + 1
					if dst.len > dst.cap {
						dst.data = newArray(d.elemType, dst.len)
						dst.cap = dst.len
					}
					copySlice(d.elemType, *dst, *slice)
					d.releaseSlice(slice)
					s.cursor++
					return nil
				case ',':
					idx++
				case nul:
					if s.read() {
						goto RETRY
					}
					slice.cap = capacity
					slice.data = data
					d.releaseSlice(slice)
					goto ERROR
				default:
					slice.cap = capacity
					slice.data = data
					d.releaseSlice(slice)
					goto ERROR
				}
				s.cursor++
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return d.errNumber(s.totalOffset())
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
	return errors.ErrUnexpectedEndOfJSON("slice", s.totalOffset())
}

func (d *sliceDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
			continue
		case 'n':
			if err := validateNull(buf, cursor); err != nil {
				return 0, err
			}
			cursor += 4
			typedmemmove(sliceType, p, nilSlice)
			return cursor, nil
		case '[':
			cursor++
			cursor = skipWhiteSpace(buf, cursor)
			if buf[cursor] == ']' {
				dst := (*sliceHeader)(p)
				if dst.data == nil {
					dst.data = newArray(d.elemType, 0)
				} else {
					dst.len = 0
				}
				cursor++
				return cursor, nil
			}
			idx := 0
			slice := d.newSlice((*sliceHeader)(p))
			srcLen := slice.len
			capacity := slice.cap
			data := slice.data
			for {
				if capacity <= idx {
					src := sliceHeader{data: data, len: idx, cap: capacity}
					capacity *= 2
					data = newArray(d.elemType, capacity)
					dst := sliceHeader{data: data, len: idx, cap: capacity}
					copySlice(d.elemType, dst, src)
				}
				ep := unsafe.Pointer(uintptr(data) + uintptr(idx)*d.size)
				// if srcLen is greater than idx, keep the original reference
				if srcLen <= idx {
					if d.isElemPointerType {
						**(**unsafe.Pointer)(unsafe.Pointer(&ep)) = nil // initialize elem pointer
					} else {
						// assign new element to the slice
						typedmemmove(d.elemType, ep, unsafe_New(d.elemType))
					}
				}
				c, err := d.valueDecoder.Decode(ctx, cursor, depth, ep)
				if err != nil {
					return 0, err
				}
				cursor = c
				cursor = skipWhiteSpace(buf, cursor)
				switch buf[cursor] {
				case ']':
					slice.cap = capacity
					slice.len = idx + 1
					slice.data = data
					dst := (*sliceHeader)(p)
					dst.len = idx + 1
					if dst.len > dst.cap {
						dst.data = newArray(d.elemType, dst.len)
						dst.cap = dst.len
					}
					copySlice(d.elemType, *dst, *slice)
					d.releaseSlice(slice)
					cursor++
					return cursor, nil
				case ',':
					idx++
				default:
					slice.cap = capacity
					slice.data = data
					d.releaseSlice(slice)
					return 0, errors.ErrInvalidCharacter(buf[cursor], "slice", cursor)
				}
				cursor++
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return 0, d.errNumber(cursor)
		default:
			return 0, errors.ErrUnexpectedEndOfJSON("slice", cursor)
		}
	}
}

func (d *sliceDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	buf := ctx.Buf
	depth++
	if depth > maxDecodeNestingDepth {
		return nil, 0, errors.ErrExceededMaxDepth(buf[cursor], cursor)
	}

	ret := [][]byte{}
	for {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			cursor++
			continue
		case 'n':
			if err := validateNull(buf, cursor); err != nil {
				return nil, 0, err
			}
			cursor += 4
			return [][]byte{nullbytes}, cursor, nil
		case '[':
			cursor++
			cursor = skipWhiteSpace(buf, cursor)
			if buf[cursor] == ']' {
				cursor++
				return ret, cursor, nil
			}
			idx := 0
			for {
				child, found, err := ctx.Option.Path.node.Index(idx)
				if err != nil {
					return nil, 0, err
				}
				if found {
					if child != nil {
						oldPath := ctx.Option.Path.node
						ctx.Option.Path.node = child
						paths, c, err := d.valueDecoder.DecodePath(ctx, cursor, depth)
						if err != nil {
							return nil, 0, err
						}
						ctx.Option.Path.node = oldPath
						ret = append(ret, paths...)
						cursor = c
					} else {
						start := cursor
						end, err := skipValue(buf, cursor, depth)
						if err != nil {
							return nil, 0, err
						}
						ret = append(ret, buf[start:end])
						cursor = end
					}
				} else {
					c, err := skipValue(buf, cursor, depth)
					if err != nil {
						return nil, 0, err
					}
					cursor = c
				}
				cursor = skipWhiteSpace(buf, cursor)
				switch buf[cursor] {
				case ']':
					cursor++
					return ret, cursor, nil
				case ',':
					idx++
				default:
					return nil, 0, errors.ErrInvalidCharacter(buf[cursor], "slice", cursor)
				}
				cursor++
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return nil, 0, d.errNumber(cursor)
		default:
			return nil, 0, errors.ErrUnexpectedEndOfJSON("slice", cursor)
		}
	}
}
