package optdec

import (
	"encoding"
	"encoding/json"
	"unsafe"
	"reflect"

	"github.com/bytedance/sonic/internal/rt"
)

type efaceDecoder struct {
}

func (d *efaceDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	/* check the defined pointer type for issue 379 */
	eface := (*rt.GoEface)(vp)

	/*
	 not pointer type, or nil pointer, or self-pointed interface{}, such as 
		```go
		var v interface{}
		v = &v
		return v
		``` see `issue758_test.go`.
	*/
	if eface.Value == nil || eface.Type.Kind() != reflect.Ptr || eface.Value == vp {
		ret, err := node.AsEface(ctx)
		if err != nil {
			return err
		}
		*(*interface{})(vp) = ret
		return nil
	}

	if node.IsNull() {
		if eface.Type.Indirect() || (!eface.Type.Indirect() &&  eface.Type.Pack().Elem().Kind() != reflect.Ptr) {
			*(*interface{})(vp) = nil
			return nil
		}
	}

	etp := rt.PtrElem(eface.Type)
	vp = eface.Value

	if eface.Type.IsNamed() {
		// check named pointer type, avoid call its `Unmarshaler`
		newp := vp
		etp = eface.Type
		vp = unsafe.Pointer(&newp)
	} else if !eface.Type.Indirect() {
		// check direct value
		etp = rt.UnpackType(eface.Type.Pack().Elem())
	}

	dec, err := findOrCompile(etp)
	if err != nil {
		return err
	}

	return dec.FromDom(vp, node, ctx)
}

type ifaceDecoder struct {
	typ *rt.GoType
}

func (d *ifaceDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	iface := *(*rt.GoIface)(vp)
	if iface.Itab == nil {
		return error_type(d.typ)
	}

	vt := iface.Itab.Vt
	if vt.Kind() != reflect.Ptr || iface.Value == nil {
		return error_type(d.typ)
	}

	etp := rt.PtrElem(vt)
	vp = iface.Value

	/* check the defined pointer type for issue 379 */
	if vt.IsNamed() {
		newp := vp
		etp = vt
		vp = unsafe.Pointer(&newp)
	}

	dec, err := findOrCompile(etp)
	if err != nil {
		return err
	}

	return dec.FromDom(vp, node, ctx)
}

type unmarshalTextDecoder struct {
	typ *rt.GoType
}

func (d *unmarshalTextDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		*(*unsafe.Pointer)(vp) = nil
		return nil
	}

	txt, ok := node.AsStringText(ctx)
	if !ok {
		return error_mismatch(node, ctx, d.typ.Pack())
	}

	v := *(*interface{})(unsafe.Pointer(&rt.GoEface{
		Type:  d.typ,
		Value: vp,
	}))

	// fast path
	if u, ok :=  v.(encoding.TextUnmarshaler); ok {
		return u.UnmarshalText(txt)
	}

	// slow path
	rv := reflect.ValueOf(v)
	if u, ok := rv.Interface().(encoding.TextUnmarshaler); ok {
		return u.UnmarshalText(txt)
	}

	return error_type(d.typ)
}

type unmarshalJSONDecoder struct {
	typ 	*rt.GoType
	strOpt	bool
}

func (d *unmarshalJSONDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	v := *(*interface{})(unsafe.Pointer(&rt.GoEface{
		Type: d.typ,
		Value: vp,
	}))

	var input []byte
	if d.strOpt && node.IsNull() {
		input = []byte("null")
	} else if d.strOpt {
		s, ok := node.AsStringText(ctx)
		if !ok {
			return error_mismatch(node, ctx, d.typ.Pack())
		}
		input = s
	} else {
		input = []byte(node.AsRaw(ctx))
	}

	// fast path
	if u, ok :=  v.(json.Unmarshaler); ok {
		return u.UnmarshalJSON((input))
	}

	// slow path
	rv := reflect.ValueOf(v)
	if u, ok := rv.Interface().(json.Unmarshaler); ok {
		return u.UnmarshalJSON(input)
	}

	return error_type(d.typ)
}
