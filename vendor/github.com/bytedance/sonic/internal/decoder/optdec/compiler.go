package optdec

import (
	"fmt"
	"reflect"

	"github.com/bytedance/sonic/option"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/internal/caching"
)

var (
	programCache = caching.CreateProgramCache()
)

func findOrCompile(vt *rt.GoType) (decFunc, error) {
	makeDecoder := func(vt *rt.GoType, _ ...interface{}) (interface{}, error) {
		ret, err := newCompiler().compileType(vt.Pack())
		return ret, err
	}
	if val := programCache.Get(vt); val != nil {
		return val.(decFunc), nil
	} else if ret, err := programCache.Compute(vt, makeDecoder); err == nil {
		return ret.(decFunc), nil
	} else {
		return nil, err
	}
}

type compiler struct {
	visited map[reflect.Type]bool
	depth   int
	counts  int
	opts 	option.CompileOptions
	namedPtr bool
}

func newCompiler() *compiler {
	return &compiler{
		visited: make(map[reflect.Type]bool),
		opts:  option.DefaultCompileOptions(),
	}
}

func (self *compiler) apply(opts option.CompileOptions) *compiler {
	self.opts = opts
	return self
}

const _CompileMaxDepth = 4096

func (c *compiler) enter(vt reflect.Type) {
	c.visited[vt] = true
	c.depth += 1

	if c.depth > _CompileMaxDepth {
		panic(*stackOverflow)
	}
}

func (c *compiler) exit(vt reflect.Type) {
	c.visited[vt] = false
	c.depth -= 1
}

func (c *compiler) compileInt(vt reflect.Type) decFunc {
	switch vt.Size() {
	case 4:
		switch vt.Kind() {
		case reflect.Uint:
			fallthrough
		case reflect.Uintptr:
			return &u32Decoder{}
		case reflect.Int:
			return &i32Decoder{}
		}
	case 8:
		switch vt.Kind() {
		case reflect.Uint:
			fallthrough
		case reflect.Uintptr:
			return &u64Decoder{}
		case reflect.Int:
			return &i64Decoder{}
		}
	default:
		panic("not supported pointer size: " + fmt.Sprint(vt.Size()))
	}
	panic("unreachable")
}

func (c *compiler) rescue(ep *error) {
	if val := recover(); val != nil {
		if err, ok := val.(error); ok {
			*ep = err
		} else {
			panic(val)
		}
	}
}

func (c *compiler) compileType(vt reflect.Type) (rt decFunc, err error) {
	defer c.rescue(&err)
	rt = c.compile(vt)
	return rt, err
}

func (c *compiler) compile(vt reflect.Type) decFunc {
	if c.visited[vt] {
		return &recuriveDecoder{
			typ: rt.UnpackType(vt),
		}
	}

	dec := c.tryCompilePtrUnmarshaler(vt, false)
	if dec != nil {
		return dec
	}

	return c.compileBasic(vt)
}

func (c *compiler) compileBasic(vt reflect.Type) decFunc {
	defer func() {
		c.counts += 1
	}()
	switch vt.Kind() {
	case reflect.Bool:
		return &boolDecoder{}
	case reflect.Int8:
		return &i8Decoder{}
	case reflect.Int16:
		return &i16Decoder{}
	case reflect.Int32:
		return &i32Decoder{}
	case reflect.Int64:
		return &i64Decoder{}
	case reflect.Uint8:
		return &u8Decoder{}
	case reflect.Uint16:
		return &u16Decoder{}
	case reflect.Uint32:
		return &u32Decoder{}
	case reflect.Uint64:
		return &u64Decoder{}
	case reflect.Float32:
		return &f32Decoder{}
	case reflect.Float64:
		return &f64Decoder{}
	case reflect.Uint:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Int:
		return c.compileInt(vt)
	case reflect.String:
		return c.compileString(vt)
	case reflect.Array:
		return c.compileArray(vt)
	case reflect.Interface:
		return c.compileInterface(vt)
	case reflect.Map:
		return c.compileMap(vt)
	case reflect.Ptr:
		return c.compilePtr(vt)
	case reflect.Slice:
		return c.compileSlice(vt)
	case reflect.Struct:
		return c.compileStruct(vt)
	default:
		return &unsupportedTypeDecoder{
			typ: rt.UnpackType(vt),
		}
	}
}

func (c *compiler) compilePtr(vt reflect.Type) decFunc {
	c.enter(vt)
	defer c.exit(vt)

	// special logic for Named Ptr, issue 379
	if reflect.PtrTo(vt.Elem()) != vt {
		c.namedPtr = true
		return &ptrDecoder{
			typ:   rt.UnpackType(vt.Elem()),
			deref: c.compileBasic(vt.Elem()),
		}
	}

	return &ptrDecoder{
		typ:   rt.UnpackType(vt.Elem()),
		deref: c.compile(vt.Elem()),
	}
}

func (c *compiler) compileArray(vt reflect.Type) decFunc {
	c.enter(vt)
	defer c.exit(vt)
	return &arrayDecoder{
		len:      vt.Len(),
		elemType: rt.UnpackType(vt.Elem()),
		elemDec:  c.compile(vt.Elem()),
		typ: vt,
	}
}

func (c *compiler) compileString(vt reflect.Type) decFunc {
	if vt == jsonNumberType {
		return &numberDecoder{}
	}
	return &stringDecoder{}

}

func (c *compiler) tryCompileSliceUnmarshaler(vt reflect.Type) decFunc {
	pt := reflect.PtrTo(vt.Elem())
	if pt.Implements(jsonUnmarshalerType) {
		return &sliceDecoder{
			elemType: rt.UnpackType(vt.Elem()),
			elemDec:  c.compile(vt.Elem()),
			typ: vt,
		}
	}

	if pt.Implements(encodingTextUnmarshalerType) {
		return &sliceDecoder{
			elemType: rt.UnpackType(vt.Elem()),
			elemDec:  c.compile(vt.Elem()),
			typ: vt,
		}
	}
	return nil
}

func (c *compiler) compileSlice(vt reflect.Type) decFunc {
	c.enter(vt)
	defer c.exit(vt)

	// Some common slice, use a decoder, to avoid function calls
	et := rt.UnpackType(vt.Elem())

	/* first checking `[]byte` */
	if et.Kind() == reflect.Uint8 /* []byte */ {
		return c.compileSliceBytes(vt)
	}

	dec := c.tryCompileSliceUnmarshaler(vt)
	if dec != nil {
		return dec
	}

	if vt == reflect.TypeOf([]interface{}{}) {
		return &sliceEfaceDecoder{}
	}
	if et.IsInt32() {
		return &sliceI32Decoder{}
	}
	if et.IsInt64() {
		return &sliceI64Decoder{}
	}
	if et.IsUint32() {
		return &sliceU32Decoder{}
	}
	if et.IsUint64() {
		return &sliceU64Decoder{}
	}
	if et.Kind() == reflect.String && et != rt.JsonNumberType {
		return &sliceStringDecoder{}
	}

	return &sliceDecoder{
		elemType: rt.UnpackType(vt.Elem()),
		elemDec:  c.compile(vt.Elem()),
		typ: vt,
	}
}

func (c *compiler) compileSliceBytes(vt reflect.Type) decFunc {
	ep := reflect.PtrTo(vt.Elem())

	if ep.Implements(jsonUnmarshalerType) {
		return &sliceBytesUnmarshalerDecoder{
			elemType: rt.UnpackType(vt.Elem()),
			elemDec:  c.compile(vt.Elem()),
			typ: vt,
		}
	}

	if ep.Implements(encodingTextUnmarshalerType) {
		return &sliceBytesUnmarshalerDecoder{
			elemType: rt.UnpackType(vt.Elem()),
			elemDec:  c.compile(vt.Elem()),
				typ: vt,
		}
	}

	return &sliceBytesDecoder{}
}

func (c *compiler) compileInterface(vt reflect.Type) decFunc {
	c.enter(vt)
	defer c.exit(vt)
	if vt.NumMethod() == 0 {
		return &efaceDecoder{}
	}

	if vt.Implements(jsonUnmarshalerType) {
		return &unmarshalJSONDecoder{
			typ: rt.UnpackType(vt),
		}
	}

	if vt.Implements(encodingTextUnmarshalerType) {
		return &unmarshalTextDecoder{
			typ: rt.UnpackType(vt),
		}
	}

	return &ifaceDecoder{
		typ: rt.UnpackType(vt),
	}
}

func (c *compiler) compileMap(vt reflect.Type) decFunc {
	c.enter(vt)
	defer c.exit(vt)
	// check the key unmarshaler at first
	decKey := tryCompileKeyUnmarshaler(vt)
	if decKey != nil {
		return &mapDecoder{
			mapType: rt.MapType(rt.UnpackType(vt)),
			keyDec:  decKey,
			elemDec: c.compile(vt.Elem()),
		}
	}

	// Most common map, use a decoder, to avoid function calls
	if vt == reflect.TypeOf(map[string]interface{}{}) {
		return &mapEfaceDecoder{}
	} else if vt == reflect.TypeOf(map[string]string{}) {
		return &mapStringDecoder{}
	}

	// Some common integer map later
	mt := rt.MapType(rt.UnpackType(vt))

	if mt.Key.Kind() == reflect.String && mt.Key != rt.JsonNumberType {
		return &mapStrKeyDecoder{
			mapType: mt,
			assign: rt.GetMapStrAssign(vt),
			elemDec: c.compile(vt.Elem()),
		}
	}

	if mt.Key.IsInt64() {
		return &mapI64KeyDecoder{
			mapType: mt,
			elemDec: c.compile(vt.Elem()),
			assign: rt.GetMap64Assign(vt),
		}
	}

	if mt.Key.IsInt32() {
		return &mapI32KeyDecoder{
			mapType: mt,
			elemDec: c.compile(vt.Elem()),
			assign: rt.GetMap32Assign(vt),
		}
	}

	if mt.Key.IsUint64() {
		return &mapU64KeyDecoder{
			mapType: mt,
			elemDec: c.compile(vt.Elem()),
			assign: rt.GetMap64Assign(vt),
		}
	}

	if mt.Key.IsUint32() {
		return &mapU32KeyDecoder{
			mapType: mt,
			elemDec: c.compile(vt.Elem()),
			assign: rt.GetMap32Assign(vt),
		}
	}

	// Generic map
	return &mapDecoder{
		mapType: mt,
		keyDec:  c.compileMapKey(vt),
		elemDec: c.compile(vt.Elem()),
	}
}

func tryCompileKeyUnmarshaler(vt reflect.Type) decKey {
	pt := reflect.PtrTo(vt.Key())

	/* check for `encoding.TextUnmarshaler` with pointer receiver */
	if pt.Implements(encodingTextUnmarshalerType) {
		return decodeKeyTextUnmarshaler
	}

	/* NOTE: encoding/json not support map key with `json.Unmarshaler` */
	return nil
}

func (c *compiler) compileMapKey(vt reflect.Type) decKey {
	switch vt.Key().Kind() {
	case reflect.Int8:
		return decodeKeyI8
	case reflect.Int16:
		return decodeKeyI16
	case reflect.Uint8:
		return decodeKeyU8
	case reflect.Uint16:
		return decodeKeyU16
	// NOTE: actually, encoding/json can't use float as map key
	case reflect.Float32:
		return decodeFloat32Key
	case reflect.Float64:
		return decodeFloat64Key
	case reflect.String:
		if rt.UnpackType(vt.Key()) == rt.JsonNumberType {
			return decodeJsonNumberKey
		}
		fallthrough
	default:
		return nil
	}
}

// maybe vt is a named type, and not a pointer receiver, see issue 379  
func (c *compiler) tryCompilePtrUnmarshaler(vt reflect.Type, strOpt bool) decFunc {
	pt := reflect.PtrTo(vt)

	/* check for `json.Unmarshaler` with pointer receiver */
	if pt.Implements(jsonUnmarshalerType) {
		return &unmarshalJSONDecoder{
			typ: rt.UnpackType(pt),
			strOpt: strOpt,
		}
	}

	/* check for `encoding.TextMarshaler` with pointer receiver */
	if pt.Implements(encodingTextUnmarshalerType) {
		/* TextUnmarshal not support, string tag */
		if strOpt {
			panicForInvalidStrType(vt)
		}
		return &unmarshalTextDecoder{
			typ: rt.UnpackType(pt),
		}
	}

	return nil
}

func panicForInvalidStrType(vt reflect.Type) {
	panic(error_type(rt.UnpackType(vt)))
}
