package vm

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/goccy/go-json/internal/encoder"
	"github.com/goccy/go-json/internal/runtime"
)

const uintptrSize = 4 << (^uintptr(0) >> 63)

var (
	appendInt           = encoder.AppendInt
	appendUint          = encoder.AppendUint
	appendFloat32       = encoder.AppendFloat32
	appendFloat64       = encoder.AppendFloat64
	appendString        = encoder.AppendString
	appendByteSlice     = encoder.AppendByteSlice
	appendNumber        = encoder.AppendNumber
	errUnsupportedValue = encoder.ErrUnsupportedValue
	errUnsupportedFloat = encoder.ErrUnsupportedFloat
	mapiterinit         = encoder.MapIterInit
	mapiterkey          = encoder.MapIterKey
	mapitervalue        = encoder.MapIterValue
	mapiternext         = encoder.MapIterNext
	maplen              = encoder.MapLen
)

type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

type nonEmptyInterface struct {
	itab *struct {
		ityp *runtime.Type // static interface type
		typ  *runtime.Type // dynamic concrete type
		// unused fields...
	}
	ptr unsafe.Pointer
}

func errUnimplementedOp(op encoder.OpType) error {
	return fmt.Errorf("encoder: opcode %s has not been implemented", op)
}

func load(base uintptr, idx uint32) uintptr {
	addr := base + uintptr(idx)
	return **(**uintptr)(unsafe.Pointer(&addr))
}

func store(base uintptr, idx uint32, p uintptr) {
	addr := base + uintptr(idx)
	**(**uintptr)(unsafe.Pointer(&addr)) = p
}

func loadNPtr(base uintptr, idx uint32, ptrNum uint8) uintptr {
	addr := base + uintptr(idx)
	p := **(**uintptr)(unsafe.Pointer(&addr))
	for i := uint8(0); i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = ptrToPtr(p)
	}
	return p
}

func ptrToUint64(p uintptr, bitSize uint8) uint64 {
	switch bitSize {
	case 8:
		return (uint64)(**(**uint8)(unsafe.Pointer(&p)))
	case 16:
		return (uint64)(**(**uint16)(unsafe.Pointer(&p)))
	case 32:
		return (uint64)(**(**uint32)(unsafe.Pointer(&p)))
	case 64:
		return **(**uint64)(unsafe.Pointer(&p))
	}
	return 0
}
func ptrToFloat32(p uintptr) float32            { return **(**float32)(unsafe.Pointer(&p)) }
func ptrToFloat64(p uintptr) float64            { return **(**float64)(unsafe.Pointer(&p)) }
func ptrToBool(p uintptr) bool                  { return **(**bool)(unsafe.Pointer(&p)) }
func ptrToBytes(p uintptr) []byte               { return **(**[]byte)(unsafe.Pointer(&p)) }
func ptrToNumber(p uintptr) json.Number         { return **(**json.Number)(unsafe.Pointer(&p)) }
func ptrToString(p uintptr) string              { return **(**string)(unsafe.Pointer(&p)) }
func ptrToSlice(p uintptr) *runtime.SliceHeader { return *(**runtime.SliceHeader)(unsafe.Pointer(&p)) }
func ptrToPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}
func ptrToNPtr(p uintptr, ptrNum uint8) uintptr {
	for i := uint8(0); i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = ptrToPtr(p)
	}
	return p
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
func ptrToInterface(code *encoder.Opcode, p uintptr) interface{} {
	return *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: code.Type,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
}

func appendBool(_ *encoder.RuntimeContext, b []byte, v bool) []byte {
	if v {
		return append(b, "true"...)
	}
	return append(b, "false"...)
}

func appendNull(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, "null"...)
}

func appendComma(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, ',')
}

func appendNullComma(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, "null,"...)
}

func appendColon(_ *encoder.RuntimeContext, b []byte) []byte {
	last := len(b) - 1
	b[last] = ':'
	return b
}

func appendMapKeyValue(_ *encoder.RuntimeContext, _ *encoder.Opcode, b, key, value []byte) []byte {
	b = append(b, key...)
	b[len(b)-1] = ':'
	return append(b, value...)
}

func appendMapEnd(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte {
	b[len(b)-1] = '}'
	b = append(b, ',')
	return b
}

func appendMarshalJSON(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte, v interface{}) ([]byte, error) {
	return encoder.AppendMarshalJSON(ctx, code, b, v)
}

func appendMarshalText(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte, v interface{}) ([]byte, error) {
	return encoder.AppendMarshalText(ctx, code, b, v)
}

func appendArrayHead(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte {
	return append(b, '[')
}

func appendArrayEnd(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	b[last] = ']'
	return append(b, ',')
}

func appendEmptyArray(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '[', ']', ',')
}

func appendEmptyObject(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '{', '}', ',')
}

func appendObjectEnd(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	b[last] = '}'
	return append(b, ',')
}

func appendStructHead(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '{')
}

func appendStructKey(_ *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	return append(b, code.Key...)
}

func appendStructEnd(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte {
	return append(b, '}', ',')
}

func appendStructEndSkipLast(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	if b[last] == ',' {
		b[last] = '}'
		return appendComma(ctx, b)
	}
	return appendStructEnd(ctx, code, b)
}

func restoreIndent(_ *encoder.RuntimeContext, _ *encoder.Opcode, _ uintptr)               {}
func storeIndent(_ uintptr, _ *encoder.Opcode, _ uintptr)                                 {}
func appendMapKeyIndent(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte    { return b }
func appendArrayElemIndent(_ *encoder.RuntimeContext, _ *encoder.Opcode, b []byte) []byte { return b }
