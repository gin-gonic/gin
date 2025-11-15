package optdec

import (
	"fmt"
	"reflect"

	caching "github.com/bytedance/sonic/internal/optcaching"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/internal/resolver"
)

const (
    _MAX_FIELDS = 50        // cutoff at 50 fields struct
)

func (c *compiler) compileIntStringOption(vt reflect.Type) decFunc {
	switch vt.Size() {
	case 4:
		switch vt.Kind() {
		case reflect.Uint:
			fallthrough
		case reflect.Uintptr:
			return &u32StringDecoder{}
		case reflect.Int:
			return &i32StringDecoder{}
		}
	case 8:
		switch vt.Kind() {
		case reflect.Uint:
			fallthrough
		case reflect.Uintptr:
			return &u64StringDecoder{}
		case reflect.Int:
			return &i64StringDecoder{}
		}
	default:
		panic("not supported pointer size: " + fmt.Sprint(vt.Size()))
	}
	panic("unreachable")
}

func isInteger(vt reflect.Type) bool {
	switch vt.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr, reflect.Int: return true
		default: return false
	}
}

func (c *compiler) assertStringOptTypes(vt reflect.Type) {
	if c.depth > _CompileMaxDepth {
		panic(*stackOverflow)
	}

	c.depth += 1
	defer func ()  {
		c.depth -= 1
	}()

	if isInteger(vt) {
		return
	}

	switch vt.Kind() {
	case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
		return
	case reflect.Ptr: c.assertStringOptTypes(vt.Elem())
	default:
		panicForInvalidStrType(vt)
	}
}

func (c *compiler) compileFieldStringOption(vt reflect.Type) decFunc {
	c.assertStringOptTypes(vt)
	unmDec := c.tryCompilePtrUnmarshaler(vt, true)
	if unmDec != nil { 
		return unmDec
	} 

	switch vt.Kind() {
	case reflect.String:
		if vt == jsonNumberType {
			return &numberStringDecoder{}
		}
		return &strStringDecoder{}
	case reflect.Bool:
		return &boolStringDecoder{}
	case reflect.Int8:
		return &i8StringDecoder{}
	case reflect.Int16:
		return &i16StringDecoder{}
	case reflect.Int32:
		return &i32StringDecoder{}
	case reflect.Int64:
		return &i64StringDecoder{}
	case reflect.Uint8:
		return &u8StringDecoder{}
	case reflect.Uint16:
		return &u16StringDecoder{}
	case reflect.Uint32:
		return &u32StringDecoder{}
	case reflect.Uint64:
		return &u64StringDecoder{}
	case reflect.Float32:
		return &f32StringDecoder{}
	case reflect.Float64:
		return &f64StringDecoder{}
	case reflect.Uint:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Int:
		return c.compileIntStringOption(vt)
	case reflect.Ptr:
		return &ptrStrDecoder{
			typ:   rt.UnpackType(vt.Elem()),
			deref: c.compileFieldStringOption(vt.Elem()),
		}
	default:
		panicForInvalidStrType(vt)
		return nil
	}
}

func (c *compiler) compileStruct(vt reflect.Type) decFunc {
	c.enter(vt)
	defer c.exit(vt)
	if c.namedPtr {
		c.namedPtr = false
		return c.compileStructBody(vt)
	}

	if c.depth >= c.opts.MaxInlineDepth + 1 || (c.counts > 0 &&  vt.NumField() >= _MAX_FIELDS) {
		return &recuriveDecoder{
			typ: rt.UnpackType(vt),
		}
	} else {
		return c.compileStructBody(vt)
	}
}

func (c *compiler) compileStructBody(vt reflect.Type) decFunc {
	fv := resolver.ResolveStruct(vt)
	entries := make([]fieldEntry, 0, len(fv))

	for _, f := range fv {
		var dec decFunc
		/* dealt with field tag options */
		if f.Opts&resolver.F_stringize != 0 {
			dec = c.compileFieldStringOption(f.Type)
		} else {
			dec = c.compile(f.Type)
		}

		/* deal with embedded pointer fields */
		if f.Path[0].Kind == resolver.F_deref {
			dec = &embeddedFieldPtrDecoder{
				field:    	f,
				fieldDec:   dec,
				fieldName:  f.Name,
			}
		}

		entries = append(entries, fieldEntry{
			FieldMeta: f,
			fieldDec:  dec,
		})
	}
	return &structDecoder{
		fieldMap:  	caching.NewFieldLookup(fv),
		fields:     entries,
		structName: vt.Name(),
		typ: 		vt,
	}
}
