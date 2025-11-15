package optdec

import (
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/decoder/consts"
	caching "github.com/bytedance/sonic/internal/optcaching"
	"github.com/bytedance/sonic/internal/resolver"
)

type fieldEntry struct {
	resolver.FieldMeta
	fieldDec decFunc
}

type structDecoder struct {
	fieldMap   caching.FieldLookup
	fields     []fieldEntry
	structName string
	typ        reflect.Type
}

func (d *structDecoder) FromDom(vp unsafe.Pointer, node Node, ctx *context) error {
	if node.IsNull() {
		return nil
	}

	var gerr error
	obj, ok := node.AsObj()
	if !ok {
		return error_mismatch(node, ctx, d.typ)
	}

	next := obj.Children()
	for i := 0; i < obj.Len(); i++ {
		key, _ := NewNode(next).AsStrRef(ctx)
		val := NewNode(PtrOffset(next, 1))
		next = val.Next()

		// find field idx
		idx := d.fieldMap.Get(key, ctx.Options()&uint64(consts.OptionCaseSensitive) != 0)
        if idx == -1 {
            if Options(ctx.Options())&OptionDisableUnknown != 0 {
                return error_field(key)
            }
            continue
        }

		offset := d.fields[idx].Path[0].Size
		elem := unsafe.Pointer(uintptr(vp) + offset)
		err := d.fields[idx].fieldDec.FromDom(elem, val, ctx)

		// deal with mismatch type errors
		if gerr == nil && err != nil {
			// TODO: better error info
			gerr = err
		}
	}
	return gerr
}

