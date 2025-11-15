package optdec

import (
	"encoding"
	"encoding/json"
	"math"
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

/** Decoder for most common map types: map[string]interface{}, map[string]string **/

type mapEfaceDecoder struct {
}

func (d *mapEfaceDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*map[string]interface{})(vp) = nil
		return nil
	}

	return node.AsMapEface(ctx, vp)
}

type mapStringDecoder struct {
}

func (d *mapStringDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*map[string]string)(vp) = nil
		return nil
	}

	return node.AsMapString(ctx, vp)
}

/** Decoder for map with string key **/

type mapStrKeyDecoder struct {
	mapType *rt.GoMapType
	elemDec decFunc
	assign  rt.MapStrAssign
	typ 	reflect.Type
}

func (d *mapStrKeyDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	obj, ok := node.AsObj()
	if !ok {
		return error_mismatch(node, ctx, d.mapType.Pack())
	}

	// allocate map
	m := *(*unsafe.Pointer)(vp)
	if m == nil {
		m = rt.Makemap(&d.mapType.GoType, obj.Len())
	}

	var gerr error
	next := obj.Children()
	for i := 0; i < obj.Len(); i++ {
		keyn := NewNode(next)
		key, _ := keyn.AsStr(ctx)

		valn := NewNode(PtrOffset(next, 1))
		valp := d.assign(d.mapType, m, key)
		err := d.elemDec.FromDom(valp, valn, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = valn.Next()
	}

	*(*unsafe.Pointer)(vp) = m
	return gerr
}

/** Decoder for map with int32 or int64 key **/

type mapI32KeyDecoder struct {
	mapType *rt.GoMapType
	elemDec decFunc
	assign rt.Map32Assign
}

func (d *mapI32KeyDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	obj, ok := node.AsObj()
	if !ok {
		return error_mismatch(node, ctx, d.mapType.Pack())
	}

	// allocate map
	m := *(*unsafe.Pointer)(vp)
	if m == nil {
		m = rt.Makemap(&d.mapType.GoType, obj.Len())
	}

	next := obj.Children()
	var gerr error
	for i := 0; i < obj.Len(); i++ {
		keyn := NewNode(next)
		k, ok := keyn.ParseI64(ctx)
		if !ok || k > math.MaxInt32 || k < math.MinInt32 {
			if gerr == nil {
				gerr = error_mismatch(keyn, ctx, d.mapType.Pack())
			}
			valn := NewNode(PtrOffset(next, 1))
			next = valn.Next()
			continue
		}

		key := int32(k)
		ku32 := *(*uint32)(unsafe.Pointer(&key))
		valn := NewNode(PtrOffset(next, 1))
		valp := d.assign(d.mapType, m, ku32)
		err := d.elemDec.FromDom(valp, valn, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}

		next = valn.Next()
	}

	*(*unsafe.Pointer)(vp) = m
	return gerr
}

type mapI64KeyDecoder struct {
	mapType *rt.GoMapType
	elemDec decFunc
	assign rt.Map64Assign
}

func (d *mapI64KeyDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	obj, ok := node.AsObj()
	if !ok {
		return error_mismatch(node, ctx, d.mapType.Pack())
	}

	// allocate map
	m := *(*unsafe.Pointer)(vp)
	if m == nil {
		m = rt.Makemap(&d.mapType.GoType, obj.Len())
	}

	var gerr error
	next := obj.Children()
	for i := 0; i < obj.Len(); i++ {
		keyn := NewNode(next)
		key, ok := keyn.ParseI64(ctx)

		if !ok {
			if gerr == nil {
				gerr = error_mismatch(keyn, ctx, d.mapType.Pack())
			}
			valn := NewNode(PtrOffset(next, 1))
			next = valn.Next()
			continue
		}

		ku64 := *(*uint64)(unsafe.Pointer(&key))
		valn := NewNode(PtrOffset(next, 1))
		valp := d.assign(d.mapType, m, ku64)
		err := d.elemDec.FromDom(valp, valn, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = valn.Next()
	}

	*(*unsafe.Pointer)(vp) = m
	return gerr
}

/** Decoder for map with unt32 or uint64 key **/

type mapU32KeyDecoder struct {
	mapType *rt.GoMapType
	elemDec decFunc
	assign  rt.Map32Assign
}

func (d *mapU32KeyDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	obj, ok := node.AsObj()
	if !ok {
		return error_mismatch(node, ctx, d.mapType.Pack())
	}

	// allocate map
	m := *(*unsafe.Pointer)(vp)
	if m == nil {
		m = rt.Makemap(&d.mapType.GoType, obj.Len())
	}

	var gerr error
	next := obj.Children()
	for i := 0; i < obj.Len(); i++ {
		keyn := NewNode(next)
		k, ok := keyn.ParseU64(ctx)
		if !ok || k > math.MaxUint32 {
			if gerr == nil {
				gerr = error_mismatch(keyn, ctx, d.mapType.Pack())
			}
			valn := NewNode(PtrOffset(next, 1))
			next = valn.Next()
			continue
		}

		key := uint32(k)
		valn := NewNode(PtrOffset(next, 1))
		valp := d.assign(d.mapType, m, key)
		err := d.elemDec.FromDom(valp, valn, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = valn.Next()
	}

	*(*unsafe.Pointer)(vp) = m
	return gerr
}

type mapU64KeyDecoder struct {
	mapType *rt.GoMapType
	elemDec decFunc
	assign rt.Map64Assign
}

func (d *mapU64KeyDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	obj, ok := node.AsObj()
	if !ok {
		return  error_mismatch(node, ctx, d.mapType.Pack())
	}
	// allocate map
	m := *(*unsafe.Pointer)(vp)
	if m == nil {
		m = rt.Makemap(&d.mapType.GoType, obj.Len())
	}

	var gerr error
	next := obj.Children()
	for i := 0; i < obj.Len(); i++ {
		keyn := NewNode(next)
		key, ok := keyn.ParseU64(ctx)
		if !ok {
			if gerr == nil {
				gerr = error_mismatch(keyn, ctx, d.mapType.Pack())
			}
			valn := NewNode(PtrOffset(next, 1))
			next = valn.Next()
			continue
		}

		valn := NewNode(PtrOffset(next, 1))
		valp := d.assign(d.mapType, m, key)
		err := d.elemDec.FromDom(valp, valn, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
		next = valn.Next()
	}

	*(*unsafe.Pointer)(vp) = m
	return gerr
}

/** Decoder for generic cases */

type decKey func(dec *mapDecoder, raw string) (interface{}, error)

func decodeKeyU8(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	ret, err := ParseU64(key)
	if err != nil {
		return nil, err
	}
	if ret > math.MaxUint8 {
		return nil, error_value(key, dec.mapType.Key.Pack())
	}
	return uint8(ret), nil
}

func decodeKeyU16(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	ret, err := ParseU64(key)
	if err != nil {
		return nil, err
	}
	if ret > math.MaxUint16 {
		return nil, error_value(key, dec.mapType.Key.Pack())
	}
	return uint16(ret), nil
}

func decodeKeyI8(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	ret, err := ParseI64(key)
	if err != nil {
		return nil, err
	}
	if ret > math.MaxInt8 || ret < math.MinInt8 {
		return nil, error_value(key, dec.mapType.Key.Pack())
	}
	return int8(ret), nil
}

func decodeKeyI16(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	ret, err := ParseI64(key)
	if err != nil {
		return nil, err
	}
	if ret > math.MaxInt16 || ret < math.MinInt16 {
		return nil, error_value(key, dec.mapType.Key.Pack())
	}
	return int16(ret), nil
}

func decodeKeyTextUnmarshaler(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	ret := reflect.New(dec.mapType.Key.Pack()).Interface()
	err = ret.(encoding.TextUnmarshaler).UnmarshalText(rt.Str2Mem(key))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func decodeFloat32Key(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	ret, err := ParseF64(key)
	if err != nil {
		return nil, err
	}
	if ret > math.MaxFloat32 || ret < -math.MaxFloat32 {
		return nil, error_value(key, dec.mapType.Key.Pack())
	}
	return float32(ret), nil
}

func decodeFloat64Key(dec *mapDecoder, raw string) (interface{}, error) {
	key, err := Unquote(raw)
	if err != nil {
		return nil, err
	}
	return ParseF64(key)
}

func decodeJsonNumberKey(dec *mapDecoder, raw string) (interface{}, error) {
	// skip the quote
	raw = raw[1:len(raw)-1]
	end, ok := SkipNumberFast(raw, 0)

	// check trailing chars
	if !ok || end != len(raw) {
		return nil, error_value(raw, rt.JsonNumberType.Pack())
	}
	
	return json.Number(raw[0:end]), nil
}

type mapDecoder struct {
	mapType *rt.GoMapType
	keyDec  decKey
	elemDec decFunc
}

func (d *mapDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	obj, ok := node.AsObj()
	if !ok || d.keyDec == nil {
		return error_mismatch(node, ctx, d.mapType.Pack())
	}

	// allocate map
	m := *(*unsafe.Pointer)(vp)
	if m == nil {
		m = rt.Makemap(&d.mapType.GoType, obj.Len())
	}

	next := obj.Children()
	var gerr error
	for i := 0; i < obj.Len(); i++ {
		keyn := NewNode(next)
		raw := keyn.AsRaw(ctx)

		key, err := d.keyDec(d, raw)
		if err != nil {
			if gerr == nil {
				gerr = error_mismatch(keyn, ctx, d.mapType.Pack())
			}
			valn := NewNode(PtrOffset(next, 1))
			next = valn.Next()
			continue
		}

		valn := NewNode(PtrOffset(next, 1))
		keyp := rt.UnpackEface(key).Value
		valp := rt.Mapassign(d.mapType, m, keyp)
		err = d.elemDec.FromDom(valp, valn, ctx)
		if gerr == nil && err != nil {
			gerr = err
		}

		next = valn.Next()
	}

	*(*unsafe.Pointer)(vp) = m
	return gerr
}
