package optdec

import (
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

type sliceDecoder struct {
	elemType *rt.GoType
	elemDec  decFunc
	typ      reflect.Type
}

var (
	emptyPtr = &struct{}{}
)

func (d *sliceDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	arr, ok := node.AsArr()
	if !ok {
		return error_mismatch(node, ctx, d.typ)
	}

	slice := rt.MakeSlice(vp, d.elemType, arr.Len())
	elems := slice.Ptr
	next := arr.Children()

	var gerr error
	for i := 0; i < arr.Len(); i++ {
		val := NewNode(next)
		elem := unsafe.Pointer(uintptr(elems) + uintptr(i)*d.elemType.Size)
		err := d.elemDec.FromDom(elem, val, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = val.Next()
	}

	*(*rt.GoSlice)(vp) = *slice
	return gerr
}

type arrayDecoder struct {
	len      int
	elemType *rt.GoType
	elemDec  decFunc
	typ   	reflect.Type
}

//go:nocheckptr
func (d *arrayDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	arr, ok := node.AsArr()
	if !ok {
		return error_mismatch(node, ctx, d.typ)
	}

	next := arr.Children()
	i := 0

	var gerr error
	for ; i < d.len && i < arr.Len(); i++ {
		elem := unsafe.Pointer(uintptr(vp) + uintptr(i)*d.elemType.Size)
		val := NewNode(next)
		err := d.elemDec.FromDom(elem, val, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = val.Next()
	}

	/* zero rest of array */
	addr := uintptr(vp) + uintptr(i)*d.elemType.Size
	n := uintptr(d.len-i) * d.elemType.Size

	/* the boundary pointer may points to another unknown object, so we need to avoid using it */
	if n != 0 {
		rt.ClearMemory(d.elemType, unsafe.Pointer(addr), n)
	}
	return gerr
}

type sliceEfaceDecoder struct {
}

func (d *sliceEfaceDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	/* if slice is empty, just call `AsSliceEface` */
	if ((*rt.GoSlice)(vp)).Len == 0 {
		return node.AsSliceEface(ctx, vp)
	}
	
	decoder := sliceDecoder{
		elemType: rt.AnyType,
		elemDec:  &efaceDecoder{},
		typ:      rt.SliceEfaceType.Pack(),
	}

	return decoder.FromDom(vp, node, ctx)
}

type sliceI32Decoder struct {
}

func (d *sliceI32Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	return node.AsSliceI32(ctx, vp)
}

type sliceI64Decoder struct {
}

func (d *sliceI64Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	return node.AsSliceI64(ctx, vp)
}

type sliceU32Decoder struct {
}

func (d *sliceU32Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	return node.AsSliceU32(ctx, vp)
}

type sliceU64Decoder struct {
}

func (d *sliceU64Decoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	return node.AsSliceU64(ctx, vp)
}

type sliceStringDecoder struct {
}

func (d *sliceStringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	return node.AsSliceString(ctx, vp)
}

type sliceBytesDecoder struct {
}

func (d *sliceBytesDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	s, err := node.AsSliceBytes(ctx)
	*(*[]byte)(vp) = s
	return err
}

type sliceBytesUnmarshalerDecoder struct {
	elemType *rt.GoType
	elemDec  decFunc
	typ reflect.Type
}

func (d *sliceBytesUnmarshalerDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*rt.GoSlice)(vp) = rt.GoSlice{}
		return nil
	}

	/* parse JSON string into `[]byte` */
	if node.IsStr() {
		slice, err := node.AsSliceBytes(ctx)
		if err != nil {
			return err
		}
		*(*[]byte)(vp) = slice
		return nil
	}

	/* parse JSON array into `[]byte` */
	arr, ok := node.AsArr()
	if !ok {
		return error_mismatch(node, ctx, d.typ)
	}

	slice := rt.MakeSlice(vp, d.elemType, arr.Len())
	elems := slice.Ptr

	var gerr error
	next := arr.Children()
	for i := 0; i < arr.Len(); i++ {
		child := NewNode(next)
		elem := unsafe.Pointer(uintptr(elems) + uintptr(i)*d.elemType.Size)
		err := d.elemDec.FromDom(elem, child, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = child.Next()
	}

	*(*rt.GoSlice)(vp) = *slice
	return gerr
}
