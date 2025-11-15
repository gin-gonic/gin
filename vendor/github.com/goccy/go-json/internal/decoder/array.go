package decoder

import (
	"fmt"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

type arrayDecoder struct {
	elemType     *runtime.Type
	size         uintptr
	valueDecoder Decoder
	alen         int
	structName   string
	fieldName    string
	zeroValue    unsafe.Pointer
}

func newArrayDecoder(dec Decoder, elemType *runtime.Type, alen int, structName, fieldName string) *arrayDecoder {
	// workaround to avoid checkptr errors. cannot use `*(*unsafe.Pointer)(unsafe_New(elemType))` directly.
	zeroValuePtr := unsafe_New(elemType)
	zeroValue := **(**unsafe.Pointer)(unsafe.Pointer(&zeroValuePtr))
	return &arrayDecoder{
		valueDecoder: dec,
		elemType:     elemType,
		size:         elemType.Size(),
		alen:         alen,
		structName:   structName,
		fieldName:    fieldName,
		zeroValue:    zeroValue,
	}
}

func (d *arrayDecoder) DecodeStream(s *Stream, depth int64, p unsafe.Pointer) error {
	depth++
	if depth > maxDecodeNestingDepth {
		return errors.ErrExceededMaxDepth(s.char(), s.cursor)
	}

	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
		case 'n':
			if err := nullBytes(s); err != nil {
				return err
			}
			return nil
		case '[':
			idx := 0
			s.cursor++
			if s.skipWhiteSpace() == ']' {
				for idx < d.alen {
					*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(idx)*d.size)) = d.zeroValue
					idx++
				}
				s.cursor++
				return nil
			}
			for {
				if idx < d.alen {
					if err := d.valueDecoder.DecodeStream(s, depth, unsafe.Pointer(uintptr(p)+uintptr(idx)*d.size)); err != nil {
						return err
					}
				} else {
					if err := s.skipValue(depth); err != nil {
						return err
					}
				}
				idx++
				switch s.skipWhiteSpace() {
				case ']':
					for idx < d.alen {
						*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(idx)*d.size)) = d.zeroValue
						idx++
					}
					s.cursor++
					return nil
				case ',':
					s.cursor++
					continue
				case nul:
					if s.read() {
						s.cursor++
						continue
					}
					goto ERROR
				default:
					goto ERROR
				}
			}
		case nul:
			if s.read() {
				continue
			}
			goto ERROR
		default:
			goto ERROR
		}
		s.cursor++
	}
ERROR:
	return errors.ErrUnexpectedEndOfJSON("array", s.totalOffset())
}

func (d *arrayDecoder) Decode(ctx *RuntimeContext, cursor, depth int64, p unsafe.Pointer) (int64, error) {
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
			return cursor, nil
		case '[':
			idx := 0
			cursor++
			cursor = skipWhiteSpace(buf, cursor)
			if buf[cursor] == ']' {
				for idx < d.alen {
					*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(idx)*d.size)) = d.zeroValue
					idx++
				}
				cursor++
				return cursor, nil
			}
			for {
				if idx < d.alen {
					c, err := d.valueDecoder.Decode(ctx, cursor, depth, unsafe.Pointer(uintptr(p)+uintptr(idx)*d.size))
					if err != nil {
						return 0, err
					}
					cursor = c
				} else {
					c, err := skipValue(buf, cursor, depth)
					if err != nil {
						return 0, err
					}
					cursor = c
				}
				idx++
				cursor = skipWhiteSpace(buf, cursor)
				switch buf[cursor] {
				case ']':
					for idx < d.alen {
						*(*unsafe.Pointer)(unsafe.Pointer(uintptr(p) + uintptr(idx)*d.size)) = d.zeroValue
						idx++
					}
					cursor++
					return cursor, nil
				case ',':
					cursor++
					continue
				default:
					return 0, errors.ErrInvalidCharacter(buf[cursor], "array", cursor)
				}
			}
		default:
			return 0, errors.ErrUnexpectedEndOfJSON("array", cursor)
		}
	}
}

func (d *arrayDecoder) DecodePath(ctx *RuntimeContext, cursor, depth int64) ([][]byte, int64, error) {
	return nil, 0, fmt.Errorf("json: array decoder does not support decode path")
}
