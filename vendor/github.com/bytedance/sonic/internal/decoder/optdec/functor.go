package optdec

import (
	"encoding/json"
	"math"
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/internal/resolver"
)

type decFunc interface {
	FromDom(vp unsafe.Pointer, node Node, ctx *context) error
}

type ptrDecoder struct {
	typ   *rt.GoType
	deref decFunc
}

// Pointer Value is allocated in the Caller
func (d *ptrDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	if *(*unsafe.Pointer)(vp) == nil {
		*(*unsafe.Pointer)(vp) = rt.Mallocgc(d.typ.Size, d.typ, true)
	}

	return d.deref.FromDom(*(*unsafe.Pointer)(vp), node, ctx)
}

type embeddedFieldPtrDecoder struct {
	field      resolver.FieldMeta
	fieldDec   decFunc
	fieldName  string
}

// Pointer Value is allocated in the Caller
func (d *embeddedFieldPtrDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	// seek into the pointer
	vp = unsafe.Pointer(uintptr(vp) - uintptr(d.field.Path[0].Size))
	for _, f := range d.field.Path {
		deref := rt.UnpackType(f.Type)
		vp = unsafe.Pointer(uintptr(vp) + f.Size)
		if f.Kind == resolver.F_deref {
			if  *(*unsafe.Pointer)(vp) == nil  {
				*(*unsafe.Pointer)(vp) = rt.Mallocgc(deref.Size, deref, true)
			}
			vp = *(*unsafe.Pointer)(vp)
		}
	}
	return d.fieldDec.FromDom(vp, node, ctx)
}

type i8Decoder struct{}

func (d *i8Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsI64(ctx)
	if !ok ||  ret > math.MaxInt8 || ret < math.MinInt8 {
		return error_mismatch(node, ctx, int8Type)
	}

	*(*int8)(vp) = int8(ret)
	return nil
}

type i16Decoder struct{}

func (d *i16Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsI64(ctx)
	if !ok || ret > math.MaxInt16 || ret < math.MinInt16 {
		return error_mismatch(node, ctx, int16Type)
	}

	*(*int16)(vp) = int16(ret)
	return nil
}

type i32Decoder struct{}

func (d *i32Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsI64(ctx)
	if !ok ||  ret > math.MaxInt32 || ret < math.MinInt32 {
		return error_mismatch(node, ctx, int32Type)
	}

	*(*int32)(vp) = int32(ret)
	return nil
}

type i64Decoder struct{}

func (d *i64Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsI64(ctx)
	if !ok  {
		return error_mismatch(node, ctx, int64Type)
	}

	*(*int64)(vp) = int64(ret)
	return nil
}

type u8Decoder struct{}

func (d *u8Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsU64(ctx)
	if !ok || ret > math.MaxUint8 {
		err := error_mismatch(node, ctx, uint8Type)
		return err
	}

	*(*uint8)(vp) = uint8(ret)
	return nil
}

type u16Decoder struct{}

func (d *u16Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsU64(ctx)
	if !ok || ret > math.MaxUint16 {
		return error_mismatch(node, ctx, uint16Type)
	}
	*(*uint16)(vp) = uint16(ret)
	return nil
}

type u32Decoder struct{}

func (d *u32Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsU64(ctx)
	if !ok || ret > math.MaxUint32 {
		return error_mismatch(node, ctx, uint32Type)
	}

	*(*uint32)(vp) = uint32(ret)
	return nil
}

type u64Decoder struct{}

func (d *u64Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsU64(ctx)
	if !ok {
		return error_mismatch(node, ctx, uint64Type)
	}

	*(*uint64)(vp) = uint64(ret)
	return nil
}

type f32Decoder struct{}

func (d *f32Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsF64(ctx)
	if !ok || ret > math.MaxFloat32 || ret < -math.MaxFloat32 {
		return error_mismatch(node, ctx, float32Type)
	}

	*(*float32)(vp) = float32(ret)
	return nil
}

type f64Decoder struct{}

func (d *f64Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsF64(ctx)
	if !ok {
		return  error_mismatch(node, ctx, float64Type)
	}

	*(*float64)(vp) = float64(ret)
	return nil
}

type boolDecoder struct {
}

func (d *boolDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsBool()
	if !ok {
		return error_mismatch(node, ctx, boolType)
	}

	*(*bool)(vp) = bool(ret)
	return nil
}

type stringDecoder struct {
}

func (d *stringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	ret, ok := node.AsStr(ctx)
	if !ok {
		return error_mismatch(node, ctx, stringType)
	}
	*(*string)(vp) = ret
	return nil
}

type numberDecoder struct {
}

func (d *numberDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	num, ok := node.AsNumber(ctx)
	if !ok {
		return error_mismatch(node, ctx, jsonNumberType)
	}
	*(*json.Number)(vp) = num
	return nil
}

type recuriveDecoder struct {
	typ *rt.GoType
}

func (d *recuriveDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	dec, err := findOrCompile(d.typ)
	if err != nil {
		return err
	}
	return dec.FromDom(vp, node, ctx)
}

type unsupportedTypeDecoder struct {
	typ *rt.GoType
}


func (d *unsupportedTypeDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}
	return error_unsuppoted(d.typ)
}

