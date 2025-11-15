package optdec

import (
	"encoding/json"
	"math"
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

type ptrStrDecoder struct {
	typ   *rt.GoType
	deref decFunc
}

// Pointer Value is allocated in the Caller
func (d *ptrStrDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return	error_mismatch(node, ctx, stringType)
	}

	if s == "null" {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	if *(*unsafe.Pointer)(vp) == nil {
		*(*unsafe.Pointer)(vp) = rt.Mallocgc(d.typ.Size, d.typ, true)
	}

	return d.deref.FromDom(*(*unsafe.Pointer)(vp), node, ctx)
}

type boolStringDecoder struct {
}

func (d *boolStringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return error_mismatch(node, ctx, stringType)
	}

	if s == "null" {
		return nil
	}

	b, err := ParseBool(s)
	if err != nil {
		return error_mismatch(node, ctx, boolType)
	}

	*(*bool)(vp) = b
	return nil
}

func parseI64(node Node, ctx *context) (int64, error, bool) {
	if node.IsNull() {
		return 0, nil, true
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return 0, error_mismatch(node, ctx, stringType), false
	}

	if s == "null" {
		return 0, nil, true
	}

	ret, err := ParseI64(s)
	return ret, err, false
}

type i8StringDecoder struct{}

func (d *i8StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseI64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	if ret > math.MaxInt8 || ret < math.MinInt8 {
		return error_mismatch(node, ctx, int8Type)
	}

	*(*int8)(vp) = int8(ret)
	return nil
}

type i16StringDecoder struct{}

func (d *i16StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseI64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	if ret > math.MaxInt16 || ret < math.MinInt16 {
		return error_mismatch(node, ctx, int16Type)
	}

	*(*int16)(vp) = int16(ret)
	return nil
}

type i32StringDecoder struct{}

func (d *i32StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseI64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	if ret > math.MaxInt32 || ret < math.MinInt32 {
		return error_mismatch(node, ctx, int32Type)
	}

	*(*int32)(vp) = int32(ret)
	return nil
}

type i64StringDecoder struct{}

func (d *i64StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseI64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	*(*int64)(vp) = int64(ret)
	return nil
}

func parseU64(node Node, ctx *context) (uint64, error, bool) {
	if node.IsNull() {
		return 0, nil, true
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return 0, error_mismatch(node, ctx, stringType), false
	}

	if s == "null" {
		return 0, nil, true
	}

	ret, err := ParseU64(s)
	return 	ret, err, false
}

type u8StringDecoder struct{}

func (d *u8StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseU64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	if ret > math.MaxUint8 {
		return error_mismatch(node, ctx, uint8Type)
	}

	*(*uint8)(vp) = uint8(ret)
	return nil
}

type u16StringDecoder struct{}

func (d *u16StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseU64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	if ret > math.MaxUint16 {
		return error_mismatch(node, ctx, uint16Type)
	}

	*(*uint16)(vp) = uint16(ret)
	return nil
}

type u32StringDecoder struct{}

func (d *u32StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseU64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	if ret > math.MaxUint32 {
		return error_mismatch(node, ctx, uint32Type)
	}

	*(*uint32)(vp) = uint32(ret)
	return nil
}


type u64StringDecoder struct{}

func (d *u64StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	ret, err, null := parseU64(node, ctx)
	if null {
		return nil
	}

	if err != nil {
		return err
	}

	*(*uint64)(vp) = uint64(ret)
	return nil
}

type f32StringDecoder struct{}

func (d *f32StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return error_mismatch(node, ctx, stringType)
	}

	if s == "null" {
		return nil
	}

	ret, err := ParseF64(s)
	if err != nil || ret > math.MaxFloat32 || ret < -math.MaxFloat32 {
		return error_mismatch(node, ctx, float32Type)
	}

	*(*float32)(vp) = float32(ret)
	return nil
}

type f64StringDecoder struct{}

func (d *f64StringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return error_mismatch(node, ctx, stringType)
	}

	if s == "null" {
		return nil
	}

	ret, err := ParseF64(s)
	if err != nil {
		return error_mismatch(node, ctx, float64Type)
	}

	*(*float64)(vp) = float64(ret)
	return nil
}

/* parse string field with string options */
type strStringDecoder struct{}

func (d *strStringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return error_mismatch(node, ctx, stringType)
	}

	if s == "null" {
		return nil
	}

	s, err := Unquote(s)
	if err != nil {
		return error_mismatch(node, ctx, stringType)
	}

	*(*string)(vp) = s
	return nil
}

type numberStringDecoder struct{}

func (d *numberStringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	s, ok := node.AsStrRef(ctx)
	if !ok {
		return error_mismatch(node, ctx, stringType)
	}

	if s == "null" {
		return nil
	}

	num, ok := node.ParseNumber(ctx)
	if !ok {
		return error_mismatch(node, ctx, jsonNumberType)
	}

	end, ok := SkipNumberFast(s, 0)
	// has error or trailing chars
	if !ok || end != len(s) {
		return error_mismatch(node, ctx, jsonNumberType)
	}

	*(*json.Number)(vp) = json.Number(num)
	return nil
}
