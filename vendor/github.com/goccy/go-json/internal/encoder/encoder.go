package encoder

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

func (t OpType) IsMultipleOpHead() bool {
	switch t {
	case OpStructHead:
		return true
	case OpStructHeadSlice:
		return true
	case OpStructHeadArray:
		return true
	case OpStructHeadMap:
		return true
	case OpStructHeadStruct:
		return true
	case OpStructHeadOmitEmpty:
		return true
	case OpStructHeadOmitEmptySlice:
		return true
	case OpStructHeadOmitEmptyArray:
		return true
	case OpStructHeadOmitEmptyMap:
		return true
	case OpStructHeadOmitEmptyStruct:
		return true
	case OpStructHeadSlicePtr:
		return true
	case OpStructHeadOmitEmptySlicePtr:
		return true
	case OpStructHeadArrayPtr:
		return true
	case OpStructHeadOmitEmptyArrayPtr:
		return true
	case OpStructHeadMapPtr:
		return true
	case OpStructHeadOmitEmptyMapPtr:
		return true
	}
	return false
}

func (t OpType) IsMultipleOpField() bool {
	switch t {
	case OpStructField:
		return true
	case OpStructFieldSlice:
		return true
	case OpStructFieldArray:
		return true
	case OpStructFieldMap:
		return true
	case OpStructFieldStruct:
		return true
	case OpStructFieldOmitEmpty:
		return true
	case OpStructFieldOmitEmptySlice:
		return true
	case OpStructFieldOmitEmptyArray:
		return true
	case OpStructFieldOmitEmptyMap:
		return true
	case OpStructFieldOmitEmptyStruct:
		return true
	case OpStructFieldSlicePtr:
		return true
	case OpStructFieldOmitEmptySlicePtr:
		return true
	case OpStructFieldArrayPtr:
		return true
	case OpStructFieldOmitEmptyArrayPtr:
		return true
	case OpStructFieldMapPtr:
		return true
	case OpStructFieldOmitEmptyMapPtr:
		return true
	}
	return false
}

type OpcodeSet struct {
	Type                     *runtime.Type
	NoescapeKeyCode          *Opcode
	EscapeKeyCode            *Opcode
	InterfaceNoescapeKeyCode *Opcode
	InterfaceEscapeKeyCode   *Opcode
	CodeLength               int
	EndCode                  *Opcode
	Code                     Code
	QueryCache               map[string]*OpcodeSet
	cacheMu                  sync.RWMutex
}

func (s *OpcodeSet) getQueryCache(hash string) *OpcodeSet {
	s.cacheMu.RLock()
	codeSet := s.QueryCache[hash]
	s.cacheMu.RUnlock()
	return codeSet
}

func (s *OpcodeSet) setQueryCache(hash string, codeSet *OpcodeSet) {
	s.cacheMu.Lock()
	s.QueryCache[hash] = codeSet
	s.cacheMu.Unlock()
}

type CompiledCode struct {
	Code    *Opcode
	Linked  bool // whether recursive code already have linked
	CurLen  uintptr
	NextLen uintptr
}

const StartDetectingCyclesAfter = 1000

func Load(base uintptr, idx uintptr) uintptr {
	addr := base + idx
	return **(**uintptr)(unsafe.Pointer(&addr))
}

func Store(base uintptr, idx uintptr, p uintptr) {
	addr := base + idx
	**(**uintptr)(unsafe.Pointer(&addr)) = p
}

func LoadNPtr(base uintptr, idx uintptr, ptrNum int) uintptr {
	addr := base + idx
	p := **(**uintptr)(unsafe.Pointer(&addr))
	if p == 0 {
		return 0
	}
	return PtrToPtr(p)
	/*
		for i := 0; i < ptrNum; i++ {
			if p == 0 {
				return p
			}
			p = PtrToPtr(p)
		}
		return p
	*/
}

func PtrToUint64(p uintptr) uint64              { return **(**uint64)(unsafe.Pointer(&p)) }
func PtrToFloat32(p uintptr) float32            { return **(**float32)(unsafe.Pointer(&p)) }
func PtrToFloat64(p uintptr) float64            { return **(**float64)(unsafe.Pointer(&p)) }
func PtrToBool(p uintptr) bool                  { return **(**bool)(unsafe.Pointer(&p)) }
func PtrToBytes(p uintptr) []byte               { return **(**[]byte)(unsafe.Pointer(&p)) }
func PtrToNumber(p uintptr) json.Number         { return **(**json.Number)(unsafe.Pointer(&p)) }
func PtrToString(p uintptr) string              { return **(**string)(unsafe.Pointer(&p)) }
func PtrToSlice(p uintptr) *runtime.SliceHeader { return *(**runtime.SliceHeader)(unsafe.Pointer(&p)) }
func PtrToPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}
func PtrToNPtr(p uintptr, ptrNum int) uintptr {
	for i := 0; i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = PtrToPtr(p)
	}
	return p
}

func PtrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
func PtrToInterface(code *Opcode, p uintptr) interface{} {
	return *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: code.Type,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
}

func ErrUnsupportedValue(code *Opcode, ptr uintptr) *errors.UnsupportedValueError {
	v := *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: code.Type,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
	}))
	return &errors.UnsupportedValueError{
		Value: reflect.ValueOf(v),
		Str:   fmt.Sprintf("encountered a cycle via %s", code.Type),
	}
}

func ErrUnsupportedFloat(v float64) *errors.UnsupportedValueError {
	return &errors.UnsupportedValueError{
		Value: reflect.ValueOf(v),
		Str:   strconv.FormatFloat(v, 'g', -1, 64),
	}
}

func ErrMarshalerWithCode(code *Opcode, err error) *errors.MarshalerError {
	return &errors.MarshalerError{
		Type: runtime.RType2Type(code.Type),
		Err:  err,
	}
}

type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

type MapItem struct {
	Key   []byte
	Value []byte
}

type Mapslice struct {
	Items []MapItem
}

func (m *Mapslice) Len() int {
	return len(m.Items)
}

func (m *Mapslice) Less(i, j int) bool {
	return bytes.Compare(m.Items[i].Key, m.Items[j].Key) < 0
}

func (m *Mapslice) Swap(i, j int) {
	m.Items[i], m.Items[j] = m.Items[j], m.Items[i]
}

//nolint:structcheck,unused
type mapIter struct {
	key         unsafe.Pointer
	elem        unsafe.Pointer
	t           unsafe.Pointer
	h           unsafe.Pointer
	buckets     unsafe.Pointer
	bptr        unsafe.Pointer
	overflow    unsafe.Pointer
	oldoverflow unsafe.Pointer
	startBucket uintptr
	offset      uint8
	wrapped     bool
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

type MapContext struct {
	Start int
	First int
	Idx   int
	Slice *Mapslice
	Buf   []byte
	Len   int
	Iter  mapIter
}

var mapContextPool = sync.Pool{
	New: func() interface{} {
		return &MapContext{
			Slice: &Mapslice{},
		}
	},
}

func NewMapContext(mapLen int, unorderedMap bool) *MapContext {
	ctx := mapContextPool.Get().(*MapContext)
	if !unorderedMap {
		if len(ctx.Slice.Items) < mapLen {
			ctx.Slice.Items = make([]MapItem, mapLen)
		} else {
			ctx.Slice.Items = ctx.Slice.Items[:mapLen]
		}
	}
	ctx.Buf = ctx.Buf[:0]
	ctx.Iter = mapIter{}
	ctx.Idx = 0
	ctx.Len = mapLen
	return ctx
}

func ReleaseMapContext(c *MapContext) {
	mapContextPool.Put(c)
}

//go:linkname MapIterInit runtime.mapiterinit
//go:noescape
func MapIterInit(mapType *runtime.Type, m unsafe.Pointer, it *mapIter)

//go:linkname MapIterKey reflect.mapiterkey
//go:noescape
func MapIterKey(it *mapIter) unsafe.Pointer

//go:linkname MapIterNext reflect.mapiternext
//go:noescape
func MapIterNext(it *mapIter)

//go:linkname MapLen reflect.maplen
//go:noescape
func MapLen(m unsafe.Pointer) int

func AppendByteSlice(_ *RuntimeContext, b []byte, src []byte) []byte {
	if src == nil {
		return append(b, `null`...)
	}
	encodedLen := base64.StdEncoding.EncodedLen(len(src))
	b = append(b, '"')
	pos := len(b)
	remainLen := cap(b[pos:])
	var buf []byte
	if remainLen > encodedLen {
		buf = b[pos : pos+encodedLen]
	} else {
		buf = make([]byte, encodedLen)
	}
	base64.StdEncoding.Encode(buf, src)
	return append(append(b, buf...), '"')
}

func AppendFloat32(_ *RuntimeContext, b []byte, v float32) []byte {
	f64 := float64(v)
	abs := math.Abs(f64)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		f32 := float32(abs)
		if f32 < 1e-6 || f32 >= 1e21 {
			fmt = 'e'
		}
	}
	return strconv.AppendFloat(b, f64, fmt, -1, 32)
}

func AppendFloat64(_ *RuntimeContext, b []byte, v float64) []byte {
	abs := math.Abs(v)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}
	return strconv.AppendFloat(b, v, fmt, -1, 64)
}

func AppendBool(_ *RuntimeContext, b []byte, v bool) []byte {
	if v {
		return append(b, "true"...)
	}
	return append(b, "false"...)
}

var (
	floatTable = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'.': true,
		'e': true,
		'E': true,
		'+': true,
		'-': true,
	}
)

func AppendNumber(_ *RuntimeContext, b []byte, n json.Number) ([]byte, error) {
	if len(n) == 0 {
		return append(b, '0'), nil
	}
	for i := 0; i < len(n); i++ {
		if !floatTable[n[i]] {
			return nil, fmt.Errorf("json: invalid number literal %q", n)
		}
	}
	b = append(b, n...)
	return b, nil
}

func AppendMarshalJSON(ctx *RuntimeContext, code *Opcode, b []byte, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v) // convert by dynamic interface type
	if (code.Flags & AddrForMarshalerFlags) != 0 {
		if rv.CanAddr() {
			rv = rv.Addr()
		} else {
			newV := reflect.New(rv.Type())
			newV.Elem().Set(rv)
			rv = newV
		}
	}
	v = rv.Interface()
	var bb []byte
	if (code.Flags & MarshalerContextFlags) != 0 {
		marshaler, ok := v.(marshalerContext)
		if !ok {
			return AppendNull(ctx, b), nil
		}
		stdctx := ctx.Option.Context
		if ctx.Option.Flag&FieldQueryOption != 0 {
			stdctx = SetFieldQueryToContext(stdctx, code.FieldQuery)
		}
		b, err := marshaler.MarshalJSON(stdctx)
		if err != nil {
			return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
		}
		bb = b
	} else {
		marshaler, ok := v.(json.Marshaler)
		if !ok {
			return AppendNull(ctx, b), nil
		}
		b, err := marshaler.MarshalJSON()
		if err != nil {
			return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
		}
		bb = b
	}
	marshalBuf := ctx.MarshalBuf[:0]
	marshalBuf = append(append(marshalBuf, bb...), nul)
	compactedBuf, err := compact(b, marshalBuf, (ctx.Option.Flag&HTMLEscapeOption) != 0)
	if err != nil {
		return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
	}
	ctx.MarshalBuf = marshalBuf
	return compactedBuf, nil
}

func AppendMarshalJSONIndent(ctx *RuntimeContext, code *Opcode, b []byte, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v) // convert by dynamic interface type
	if (code.Flags & AddrForMarshalerFlags) != 0 {
		if rv.CanAddr() {
			rv = rv.Addr()
		} else {
			newV := reflect.New(rv.Type())
			newV.Elem().Set(rv)
			rv = newV
		}
	}
	v = rv.Interface()
	var bb []byte
	if (code.Flags & MarshalerContextFlags) != 0 {
		marshaler, ok := v.(marshalerContext)
		if !ok {
			return AppendNull(ctx, b), nil
		}
		b, err := marshaler.MarshalJSON(ctx.Option.Context)
		if err != nil {
			return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
		}
		bb = b
	} else {
		marshaler, ok := v.(json.Marshaler)
		if !ok {
			return AppendNull(ctx, b), nil
		}
		b, err := marshaler.MarshalJSON()
		if err != nil {
			return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
		}
		bb = b
	}
	marshalBuf := ctx.MarshalBuf[:0]
	marshalBuf = append(append(marshalBuf, bb...), nul)
	indentedBuf, err := doIndent(
		b,
		marshalBuf,
		string(ctx.Prefix)+strings.Repeat(string(ctx.IndentStr), int(ctx.BaseIndent+code.Indent)),
		string(ctx.IndentStr),
		(ctx.Option.Flag&HTMLEscapeOption) != 0,
	)
	if err != nil {
		return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
	}
	ctx.MarshalBuf = marshalBuf
	return indentedBuf, nil
}

func AppendMarshalText(ctx *RuntimeContext, code *Opcode, b []byte, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v) // convert by dynamic interface type
	if (code.Flags & AddrForMarshalerFlags) != 0 {
		if rv.CanAddr() {
			rv = rv.Addr()
		} else {
			newV := reflect.New(rv.Type())
			newV.Elem().Set(rv)
			rv = newV
		}
	}
	v = rv.Interface()
	marshaler, ok := v.(encoding.TextMarshaler)
	if !ok {
		return AppendNull(ctx, b), nil
	}
	bytes, err := marshaler.MarshalText()
	if err != nil {
		return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
	}
	return AppendString(ctx, b, *(*string)(unsafe.Pointer(&bytes))), nil
}

func AppendMarshalTextIndent(ctx *RuntimeContext, code *Opcode, b []byte, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v) // convert by dynamic interface type
	if (code.Flags & AddrForMarshalerFlags) != 0 {
		if rv.CanAddr() {
			rv = rv.Addr()
		} else {
			newV := reflect.New(rv.Type())
			newV.Elem().Set(rv)
			rv = newV
		}
	}
	v = rv.Interface()
	marshaler, ok := v.(encoding.TextMarshaler)
	if !ok {
		return AppendNull(ctx, b), nil
	}
	bytes, err := marshaler.MarshalText()
	if err != nil {
		return nil, &errors.MarshalerError{Type: reflect.TypeOf(v), Err: err}
	}
	return AppendString(ctx, b, *(*string)(unsafe.Pointer(&bytes))), nil
}

func AppendNull(_ *RuntimeContext, b []byte) []byte {
	return append(b, "null"...)
}

func AppendComma(_ *RuntimeContext, b []byte) []byte {
	return append(b, ',')
}

func AppendCommaIndent(_ *RuntimeContext, b []byte) []byte {
	return append(b, ',', '\n')
}

func AppendStructEnd(_ *RuntimeContext, b []byte) []byte {
	return append(b, '}', ',')
}

func AppendStructEndIndent(ctx *RuntimeContext, code *Opcode, b []byte) []byte {
	b = append(b, '\n')
	b = append(b, ctx.Prefix...)
	indentNum := ctx.BaseIndent + code.Indent - 1
	for i := uint32(0); i < indentNum; i++ {
		b = append(b, ctx.IndentStr...)
	}
	return append(b, '}', ',', '\n')
}

func AppendIndent(ctx *RuntimeContext, b []byte, indent uint32) []byte {
	b = append(b, ctx.Prefix...)
	indentNum := ctx.BaseIndent + indent
	for i := uint32(0); i < indentNum; i++ {
		b = append(b, ctx.IndentStr...)
	}
	return b
}

func IsNilForMarshaler(v interface{}) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(rv.Float()) == 0
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Func:
		return rv.IsNil()
	case reflect.Slice:
		return rv.IsNil() || rv.Len() == 0
	case reflect.String:
		return rv.Len() == 0
	}
	return false
}
