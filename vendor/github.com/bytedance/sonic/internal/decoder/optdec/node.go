package optdec

import (
	"encoding/json"
	"math"
	"unsafe"

	"github.com/bytedance/sonic/internal/envs"
	"github.com/bytedance/sonic/internal/rt"
)

type Context struct {
	Parser      *Parser
	efacePool   *efacePool
	Stack       boundedStack
	Utf8Inv     bool
}

func (ctx *Context) Options() uint64 {
	return ctx.Parser.options
}

/************************* Stack and Pool Helper *******************/

type parentStat struct {
	con 	unsafe.Pointer
	remain	uint64
}
type boundedStack struct {
	stack []parentStat
	index int
}

func newStack(size int) boundedStack {
	return boundedStack{
		stack: make([]parentStat, size + 2),
		index: 0,
	}
}

//go:nosplit
func (s *boundedStack) Pop() (unsafe.Pointer, int, bool){
	s.index--
	con := s.stack[s.index].con
	remain := s.stack[s.index].remain &^ (uint64(1) << 63)
	isObj := (s.stack[s.index].remain & (uint64(1) << 63)) != 0
	s.stack[s.index].con = nil
	s.stack[s.index].remain = 0
	return con, int(remain), isObj
}

//go:nosplit
func (s *boundedStack) Push(p unsafe.Pointer, remain int, isObj bool) {
	s.stack[s.index].con = p
	s.stack[s.index].remain = uint64(remain)
	if isObj {
		s.stack[s.index].remain |= (uint64(1) << 63)
	}
	s.index++
}

type efacePool struct{
	t64   		rt.T64Pool
	tslice 		rt.TslicePool
	tstring 	rt.TstringPool
	efaceSlice  rt.SlicePool
}

func newEfacePool(stat *jsonStat, useNumber bool) *efacePool {
	strs := int(stat.str)
	nums := 0
	if useNumber {
		strs += int(stat.number)
	} else {
		nums = int(stat.number)
	}

	return &efacePool{
		t64: rt.NewT64Pool(nums),
		tslice: rt.NewTslicePool(int(stat.array)),
		tstring: rt.NewTstringPool(strs),
		efaceSlice: rt.NewPool(rt.AnyType, int(stat.array_elems)),
	}
}

func (self *efacePool) GetMap(hint int) unsafe.Pointer {
	m := make(map[string]interface{}, hint)
	return *(*unsafe.Pointer)(unsafe.Pointer(&m))
}

func (self *efacePool) GetSlice(hint int) unsafe.Pointer {
	return unsafe.Pointer(self.efaceSlice.GetSlice(hint))
}

func (self *efacePool) ConvTSlice(val rt.GoSlice, typ *rt.GoType,  dst unsafe.Pointer) {
	self.tslice.Conv(val, typ, (*interface{})(dst))
}

func (self *efacePool) ConvF64(val float64, dst unsafe.Pointer) {
	self.t64.Conv(castU64(val), rt.Float64Type, (*interface{})(dst))
}

func (self *efacePool) ConvTstring(val string, dst unsafe.Pointer) {
	self.tstring.Conv(val, (*interface{})(dst))
}

func (self *efacePool) ConvTnum(val json.Number, dst unsafe.Pointer) {
	self.tstring.ConvNum(val, (*interface{})(dst))
}

/********************************************************/

func canUseFastMap( opts uint64, root *rt.GoType) bool {
	return envs.UseFastMap && (opts & (1 << _F_copy_string)) == 0 &&  (opts & (1 << _F_use_int64)) == 0  && (root == rt.AnyType || root == rt.MapEfaceType || root == rt.SliceEfaceType) 
}

func NewContext(json string, pos int, opts uint64, root *rt.GoType) (Context, error) {
	ctx := Context{
		Parser: newParser(json, pos, opts),
	}
	if root == rt.AnyType || root == rt.MapEfaceType || root == rt.SliceEfaceType {
		ctx.Parser.isEface = true
	}

	ecode := ctx.Parser.parse()

	if ecode != 0 {
		return ctx, ctx.Parser.fixError(ecode)
	}

	useNumber := (opts & (1 << _F_use_number )) != 0
	if canUseFastMap(opts, root) {
		ctx.efacePool = newEfacePool(&ctx.Parser.nbuf.stat, useNumber)
		ctx.Stack = newStack(int(ctx.Parser.nbuf.stat.max_depth))
	}

	return ctx, nil
}

func (ctx *Context) Delete() {
	ctx.Parser.free()
	ctx.Parser = nil
}

type Node struct {
	cptr uintptr
}

func NewNode(cptr uintptr) Node {
	return Node{cptr: cptr}
}

type Dom struct {
	cdom uintptr
}

func (ctx *Context) Root() Node {
	root := (uintptr)(((*rt.GoSlice)(unsafe.Pointer(&ctx.Parser.nodes))).Ptr)
	return Node{cptr: root}
}

type Array struct {
	cptr uintptr
}

type Object struct {
	cptr uintptr
}

func (obj Object) Len() int {
	cobj := ptrCast(obj.cptr)
	return int(uint64(cobj.val) & ConLenMask)
}

func (arr Array) Len() int {
	carr :=  ptrCast(arr.cptr)
	return int(uint64(carr.val) & ConLenMask)
}

// / Helper functions to eliminate CGO calls
func (val Node) Type() uint8 {
	ctype := ptrCast(val.cptr)
	return uint8(ctype.typ & TypeMask)
}

func (val Node) Next() uintptr {
	if val.Type() != KObject && val.Type() != KArray {
		return PtrOffset(val.cptr, 1)
	}
	cobj := ptrCast(val.cptr)
	offset := int64(uint64(cobj.val) >> ConLenBits)
	return PtrOffset(val.cptr, offset)
}

func (val *Node) next() {
	*val = NewNode(val.Next())
}

type NodeIter struct {
	next uintptr
}

func NewNodeIter(node Node) NodeIter {
	return NodeIter{next: node.cptr}
}

func (iter *NodeIter) Next() Node {
	ret := NewNode(iter.next)
	iter.next = PtrOffset(iter.next, 1)
	return ret
}


func (iter *NodeIter) Peek() Node {
	return NewNode(iter.next)
}

func (val Node) U64() uint64 {
	cnum := ptrCast(val.cptr)
	return *(*uint64)((unsafe.Pointer)(&(cnum.val)))
}

func (val Node) I64() int64 {
	cnum := ptrCast(val.cptr)
	return *(*int64)((unsafe.Pointer)(&(cnum.val)))
}

func (val Node) IsNull() bool {
	return val.Type() == KNull
}

func (val Node) IsNumber() bool {
	return val.Type() & KNumber != 0
}

func (val Node) F64() float64 {
	cnum := ptrCast(val.cptr)
	return *(*float64)((unsafe.Pointer)(&(cnum.val)))
}

func (val Node) Bool() bool {
	return val.Type() == KTrue
}

func (self Node) AsU64(ctx *Context) (uint64, bool) {
	if self.Type() == KUint {
		return self.U64(), true
	} else if self.Type() == KRawNumber {
		num, err := ParseU64(self.Raw(ctx))
		if err != nil {
			return 0, false
		}
		return num, true
	} else {
		return 0, false
	}
}

func (val *Node) AsObj() (Object, bool) {
	var ret Object
	if val.Type() != KObject {
		return ret, false
	}
	return Object{
		cptr: val.cptr,
	}, true
}

func (val Node) Obj() Object {
	return Object{cptr: val.cptr}
}

func (val Node) Arr() Array {
	return Array{cptr: val.cptr}
}

func (val *Node) AsArr() (Array, bool) {
	var ret Array
	if val.Type() != KArray {
		return ret, false
	}
	return Array{
		cptr: val.cptr,
	}, true
}

func (self Node) AsI64(ctx *Context) (int64, bool) {
	typ := self.Type()
	if typ == KUint && self.U64() <= math.MaxInt64 {
		return int64(self.U64()), true
	} else  if typ == KSint {
		return self.I64(), true
	} else if typ == KRawNumber {
		val, err := self.Number(ctx).Int64()
		if err != nil {
			return 0, false
		}
		return val, true
	} else {
		return 0, false
	}
}

func (self Node) AsByte(ctx *Context) (uint8, bool) {
	typ := self.Type()
	if typ == KUint && self.U64() <= math.MaxUint8 {
		return uint8(self.U64()), true
	} else if typ == KSint && self.I64() == 0 {
		return 0, true
	} else {
		return 0, false
	}
}

/********* Parse Node String into Value ***************/

func (val Node) ParseI64(ctx *Context) (int64, bool) {
	s, ok := val.AsStrRef(ctx)
	if !ok {
		return 0, false
	}

	if s == "null" {
		return 0, true
	}

	i, err := ParseI64(s)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (val Node) ParseBool(ctx *Context) (bool, bool) {
	s, ok := val.AsStrRef(ctx)
	if !ok {
		return false, false
	}

	if s == "null" {
		return false, true
	}

	b, err := ParseBool(s)
	if err != nil {
		return false, false
	}
	return b, true
}

func (val Node) ParseU64(ctx *Context) (uint64, bool) {
	s, ok := val.AsStrRef(ctx)
	if !ok {
		return 0, false
	}

	if s == "null" {
		return 0, true
	}

	i, err := ParseU64(s)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (val Node) ParseF64(ctx *Context) (float64, bool) {
	s, ok := val.AsStrRef(ctx)
	if !ok {
		return 0, false
	}

	if s == "null" {
		return 0, true
	}

	i, err := ParseF64(s)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (val Node) ParseString(ctx *Context) (string, bool) {
	// should not use AsStrRef
	s, ok := val.AsStr(ctx)
	if !ok {
		return "", false
	}

	if s == "null" {
		return "", true
	}

	s, err := Unquote(s)
	if err != nil {
		return "", false
	}
	return s, true
}


func (val Node) ParseNumber(ctx *Context) (json.Number, bool) {
	// should not use AsStrRef
	s, ok := val.AsStr(ctx)
	if !ok {
		return json.Number(""), false
	}

	if s == "null" {
		return json.Number(""), true
	}

	end, ok := SkipNumberFast(s, 0)
	// has error or trailing chars
	if !ok || end != len(s) {
		return json.Number(""),  false
	}
	return json.Number(s), true
}



func (val Node) AsF64(ctx *Context) (float64, bool) {
	switch val.Type() {
		case KUint: return float64(val.U64()), true
		case KSint: return float64(val.I64()), true
		case KReal: return float64(val.F64()), true
		case KRawNumber: f, err := val.Number(ctx).Float64(); return f, err == nil
		default: return 0, false
	}
}

func (val Node) AsBool() (bool, bool) {
	switch val.Type() {
		case KTrue: return true, true
		case KFalse: return false, true
		default: return false, false
	}
}

func (val Node) AsStr(ctx *Context) (string, bool) {
	switch val.Type() {
		case KStringCommon:
			s := val.StringRef(ctx)
			if (ctx.Options() & (1 << _F_copy_string) == 0) {
				return s, true
			}
			return string(rt.Str2Mem(s)), true
		case KStringEscaped:
			return val.StringCopyEsc(ctx), true
		default: return "", false
	}
}

func (val Node) AsStrRef(ctx *Context) (string, bool) {
	switch val.Type() {
	case KStringEscaped:
		node := ptrCast(val.cptr)
		offset := val.Position()
		len := int(node.val)
		return rt.Mem2Str(ctx.Parser.JsonBytes()[offset : offset + len]), true
	case KStringCommon:
		return val.StringRef(ctx), true
	default:
		return "", false
	}
}

func (val Node) AsStringText(ctx *Context) ([]byte, bool) {
	if !val.IsStr() {
		return nil, false
	}

	// clone to new bytes
	s, b := val.AsStrRef(ctx)
	return []byte(s), b
}

func (val Node) IsStr() bool {
	return (val.Type() == KStringCommon) || (val.Type() == KStringEscaped)
}

func (val Node) IsRawNumber() bool {
	return val.Type() == KRawNumber
}

func (val Node) Number(ctx *Context) json.Number {
	return json.Number(val.Raw(ctx))
}

func (val Node) Raw(ctx *Context) string {
	node := ptrCast(val.cptr)
	len := int(node.val)
	offset := val.Position()
	return ctx.Parser.Json[offset:int(offset+len)]
}

func (val Node) Position() int {
	node := ptrCast(val.cptr)
	return int(node.typ >> PosBits)
}

func (val Node) AsNumber(ctx *Context) (json.Number, bool) {
	// parse JSON string as number
	if val.IsStr() {
		s, _ := val.AsStr(ctx)
		if !ValidNumberFast(s) {
			return "", false
		} else {
			return json.Number(s), true
		}
	}

	return val.NonstrAsNumber(ctx)
}

func (val Node) NonstrAsNumber(ctx *Context) (json.Number, bool) {
	// deal with raw number
	if val.IsRawNumber() {
		return val.Number(ctx), true
	}

	// deal with parse number
	if !val.IsNumber() {
		return json.Number(""), false
	}

	start := val.Position()
	end, ok := SkipNumberFast(ctx.Parser.Json, start)
	if !ok {
		return "", false
	}
	return json.Number(ctx.Parser.Json[start:end]), true
}

func (val Node) AsRaw(ctx *Context) string {
	// fast path for unescaped strings
	switch val.Type() {
	case KNull:
		return "null"
	case KTrue:
		return "true"
	case KFalse:
		return "false"
	case KStringCommon:
		node := ptrCast(val.cptr)
		len := int(node.val)
		offset := val.Position()
		// add start and end quote
		ref := rt.Str2Mem(ctx.Parser.Json)[offset-1 : offset+len+1]
		return rt.Mem2Str(ref)
	case KRawNumber: fallthrough
	case KRaw: return val.Raw(ctx)
	case KStringEscaped:
		raw, _ := SkipOneFast(ctx.Parser.Json, val.Position() - 1)
		return raw
	default:
		raw, err := SkipOneFast(ctx.Parser.Json, val.Position())
		if err != nil {
			break
		}
		return raw
	}
	panic("should always be valid json here")
}

// reference from the input JSON as possible
func (val Node) StringRef(ctx *Context) string {
	return val.Raw(ctx)
}

//go:nocheckptr
func ptrCast(p uintptr) *node {
	return (*node)(unsafe.Pointer(p))
}

func (val Node) StringCopyEsc(ctx *Context) string {
	// check whether there are in padded
	node := ptrCast(val.cptr)
	len := int(node.val)
	offset := val.Position()
	return string(ctx.Parser.JsonBytes()[offset : offset + len])
}

func (val Node) Object() Object {
	return Object{cptr: val.cptr}
}

func (val Node) Array() Array {
	return Array{cptr: val.cptr}
}

func (val *Array) Children() uintptr {
	return PtrOffset(val.cptr, 1)
}

func (val *Object) Children() uintptr {
	return PtrOffset(val.cptr, 1)
}

func (val *Node) Equal(ctx *Context, lhs string) bool {
	// check whether escaped
	cstr := ptrCast(val.cptr)
	offset := int(val.Position())
	len := int(cstr.val)
	return lhs == ctx.Parser.Json[offset:offset+len]
}

func (node *Node) AsMapEface(ctx *Context, vp unsafe.Pointer) error {
	if node.IsNull() {
		return nil
	}

	obj, ok := node.AsObj()
	if !ok {
		return newUnmatched(node.Position(), rt.MapEfaceType)
	}

	var err, gerr error
	size := obj.Len()

	var m map[string]interface{}
	if *(*unsafe.Pointer)(vp) == nil {
		if ctx.efacePool != nil {
			p := ctx.efacePool.GetMap(size)
			m = *(*map[string]interface{})(unsafe.Pointer(&p))
		} else {
			m = make(map[string]interface{}, size)
		}
	} else {
		m = *(*map[string]interface{})(vp)
	}

	next := obj.Children()
	for i := 0; i < size; i++ {
		knode := NewNode(next)
		key, _ := knode.AsStr(ctx)
		val := NewNode(PtrOffset(next, 1))
		m[key], err = val.AsEface(ctx)
		next = val.cptr
		if gerr == nil && err != nil {
			gerr = err
		}
	}

	*(*map[string]interface{})(vp) = m
	return gerr
}

func (node *Node) AsMapString(ctx *Context, vp unsafe.Pointer) error {
	obj, ok := node.AsObj()
	if !ok {
		return newUnmatched(node.Position(), rt.MapStringType)
	}

	size := obj.Len()

	var m map[string]string
	if *(*unsafe.Pointer)(vp) == nil {
		m = make(map[string]string, size)
	} else {
		m = *(*map[string]string)(vp)
	}

	var gerr error
	next := obj.Children()
	for i := 0; i < size; i++ {
		knode := NewNode(next)
		key, _ := knode.AsStr(ctx)
		val := NewNode(PtrOffset(next, 1))
		m[key], ok = val.AsStr(ctx)
		if !ok {
			if gerr == nil {
				gerr = newUnmatched(val.Position(), rt.StringType)
			}
			next = val.Next()
		} else {
			next = PtrOffset(val.cptr, 1)
		}
	}

	*(*map[string]string)(vp) = m
	return gerr
}

func (node *Node) AsSliceEface(ctx *Context, vp unsafe.Pointer) error {
	arr, ok := node.AsArr()
	if !ok {
		return newUnmatched(node.Position(), rt.SliceEfaceType)
	}

	size := arr.Len()
	var s []interface{}
	if size != 0 && ctx.efacePool != nil {
		slice := rt.GoSlice {
			Ptr: ctx.efacePool.GetSlice(size),
			Len: size,
			Cap: size,
		}
		*(*rt.GoSlice)(unsafe.Pointer(&s)) = slice
	} else {
		s = *(*[]interface{})((unsafe.Pointer)(rt.MakeSlice(vp, rt.AnyType, size)))
	}

	*node = NewNode(arr.Children())

	var err, gerr error
	for i := 0; i < size; i++ {
		s[i], err = node.AsEface(ctx)
		if gerr == nil && err != nil {
			gerr = err
		}
	}

	*(*[]interface{})(vp) = s
	return nil
}

func (node *Node) AsSliceI32(ctx *Context, vp unsafe.Pointer) error {
	arr, ok := node.AsArr()
	if !ok {
		return newUnmatched(node.Position(), rt.SliceI32Type)
	}

	size := arr.Len()
	s := *(*[]int32)((unsafe.Pointer)(rt.MakeSlice(vp, rt.Int32Type, size)))
	next := arr.Children()

	var gerr error
	for i := 0; i < size; i++ {
		val := NewNode(next)
		ret, ok := val.AsI64(ctx)
		if !ok || ret > math.MaxInt32 || ret < math.MinInt32 {
			if gerr == nil {
				gerr = newUnmatched(val.Position(), rt.Int32Type)
			}
			next = val.Next()
		} else {
			s[i] = int32(ret)
			next = PtrOffset(val.cptr, 1)
		}
	}

	*(*[]int32)(vp) = s
	return gerr
}

func (node *Node) AsSliceI64(ctx *Context, vp unsafe.Pointer) error {
	arr, ok := node.AsArr()
	if !ok {
		return newUnmatched(node.Position(), rt.SliceI64Type)
	}

	size := arr.Len()
	s := *(*[]int64)((unsafe.Pointer)(rt.MakeSlice(vp, rt.Int64Type, size)))
	next := arr.Children()

	var gerr error
	for i := 0; i < size; i++ {
		val := NewNode(next)

		ret, ok := val.AsI64(ctx)
		if !ok {
			if gerr == nil {
				gerr = newUnmatched(val.Position(), rt.Int64Type)
			}
			next = val.Next()
		} else {
			s[i] = ret
			next = PtrOffset(val.cptr, 1)
		}
	}

	*(*[]int64)(vp) = s
	return gerr
}

func (node *Node) AsSliceU32(ctx *Context, vp unsafe.Pointer) error {
	arr, ok := node.AsArr()
	if !ok {
		return newUnmatched(node.Position(), rt.SliceU32Type)
	}

	size := arr.Len()
	next := arr.Children()
	s := *(*[]uint32)((unsafe.Pointer)(rt.MakeSlice(vp, rt.Uint32Type, size)))

	var gerr error
	for i := 0; i < size; i++ {
		val := NewNode(next)
		ret, ok := val.AsU64(ctx)
		if !ok ||  ret > math.MaxUint32 {
			if gerr == nil {
				gerr = newUnmatched(val.Position(), rt.Uint32Type)
			}
			next = val.Next()
		} else {
			s[i] = uint32(ret)
			next = PtrOffset(val.cptr, 1)
		}
	}

	*(*[]uint32)(vp) = s
	return gerr
}

func (node *Node) AsSliceU64(ctx *Context, vp unsafe.Pointer) error {
	arr, ok := node.AsArr()
	if !ok {
		return newUnmatched(node.Position(), rt.SliceU64Type)
	}

	size := arr.Len()
	next := arr.Children()

	s := *(*[]uint64)((unsafe.Pointer)(rt.MakeSlice(vp, rt.Uint64Type, size)))
	var gerr error
	for i := 0; i < size; i++ {
		val := NewNode(next)
		ret, ok := val.AsU64(ctx)
		if !ok {
			if gerr == nil {
				gerr = newUnmatched(val.Position(), rt.Uint64Type)
			}
			next = val.Next()
		} else {
			s[i] = ret
			next = PtrOffset(val.cptr, 1)
		}
	}

	*(*[]uint64)(vp) = s
	return gerr
}

func (node *Node) AsSliceString(ctx *Context, vp unsafe.Pointer) error {
	arr, ok := node.AsArr()
	if !ok {
		return newUnmatched(node.Position(), rt.SliceStringType)
	}

	size := arr.Len()
	next := arr.Children()
	s := *(*[]string)((unsafe.Pointer)(rt.MakeSlice(vp, rt.StringType, size)))

	var gerr error
	for i := 0; i < size; i++ {
		val := NewNode(next)
		ret, ok := val.AsStr(ctx)
		if !ok {
			if gerr == nil {
				gerr = newUnmatched(val.Position(), rt.StringType)
			}
			next = val.Next()
		} else {
			s[i] = ret
			next = PtrOffset(val.cptr, 1)
		}
	}

	*(*[]string)(vp) = s
	return gerr
}

func (val *Node) AsSliceBytes(ctx *Context) ([]byte, error) {
	var origin []byte
	switch val.Type() {
	case KStringEscaped:
		node := ptrCast(val.cptr)
		offset := val.Position()
		len := int(node.val)
		origin = ctx.Parser.JsonBytes()[offset : offset + len]
	case KStringCommon:
		origin = rt.Str2Mem(val.StringRef(ctx))
	case KArray:
		arr := val.Array()
		size := arr.Len()
		a := make([]byte, size)
		elem := NewNode(arr.Children())
		var gerr error
		var ok bool
		for i := 0; i < size; i++ {
			a[i], ok = elem.AsByte(ctx)
			if !ok && gerr == nil {
				gerr = newUnmatched(val.Position(), rt.BytesType)
			}
			elem = NewNode(PtrOffset(elem.cptr, 1))
		}
		return a, gerr
	default:
		return nil,  newUnmatched(val.Position(), rt.BytesType)
	}
	
	b64, err := rt.DecodeBase64(origin)
	if err != nil {
		return nil, newUnmatched(val.Position(), rt.BytesType)
	}
	return b64, nil
}

// AsEface will always ok, because we have parse in native.
func (node *Node) AsEface(ctx *Context) (interface{}, error) {
	if ctx.efacePool != nil {
		iter := NewNodeIter(*node)
		v := AsEfaceFast(&iter, ctx)
		*node = iter.Peek()
		return v, nil
	} else {
		return node.AsEfaceFallback(ctx)
	}
}

func parseSingleNode(node Node, ctx *Context) interface{} {
	var v interface{}
	switch node.Type() {
		case KObject: 			v = map[string]interface{}{}
		case KArray: 			v = []interface{}{}
		case KStringCommon: 	v = node.StringRef(ctx)
		case KStringEscaped:	v = node.StringCopyEsc(ctx)
		case KTrue:				v = true
		case KFalse:			v = false
		case KNull:				v = nil
		case KUint:				v = float64(node.U64())
		case KSint: 			v = float64(node.I64())
		case KReal:				v = float64(node.F64())
		case KRawNumber:		v = node.Number(ctx)
		default:				panic("unreachable for as eface")
	}
	return v
}

func castU64(val float64) uint64 {
	return *((*uint64)(unsafe.Pointer((&val))))
}

func AsEfaceFast(iter *NodeIter, ctx *Context) interface{} {
	var mp, sp, parent unsafe.Pointer // current container pointer
	var node Node
	var size int
	var isObj bool
	var slice rt.GoSlice
	var val unsafe.Pointer
	var vt **rt.GoType
	var vp *unsafe.Pointer
	var rootM unsafe.Pointer
	var rootS rt.GoSlice
	var root interface{}
	var key string

	node = iter.Next()

	switch node.Type() {
	case KObject: 
		size = node.Object().Len()
		if size != 0 {
			ctx.Stack.Push(nil, 0, true)
			mp = ctx.efacePool.GetMap(size)
			rootM = mp
			isObj = true
			goto _object_key
		} else {
			return rt.GoEface {
				Type: rt.MapEfaceType,
				Value: ctx.efacePool.GetMap(0),
			}.Pack()
		}
	case KArray:
		size = node.Array().Len()
		if size != 0 {
			ctx.Stack.Push(nil, 0, false)
			sp = ctx.efacePool.GetSlice(size)
			slice = rt.GoSlice {
				Ptr: sp,
				Len: size,
				Cap: size,
			}
			rootS = slice
			isObj = false
			val = sp
			goto _arr_val;
		} else {
			ctx.efacePool.ConvTSlice(rt.EmptySlice, rt.SliceEfaceType, unsafe.Pointer(&root))
		}
	case KStringCommon: 	ctx.efacePool.ConvTstring(node.StringRef(ctx), unsafe.Pointer(&root))
	case KStringEscaped:	ctx.efacePool.ConvTstring(node.StringCopyEsc(ctx), unsafe.Pointer(&root))  
	case KTrue:				root = true
	case KFalse:			root = false
	case KNull:				root = nil
	case KUint:				ctx.efacePool.ConvF64(float64(node.U64()), unsafe.Pointer(&root))  
	case KSint: 			ctx.efacePool.ConvF64(float64(node.I64()), unsafe.Pointer(&root))
	case KReal:				ctx.efacePool.ConvF64(node.F64(), unsafe.Pointer(&root))
	case KRawNumber:		ctx.efacePool.ConvTnum(node.Number(ctx), unsafe.Pointer(&root))
	default:				panic("unreachable for as eface")
	}
	return root

_object_key:
	node = iter.Next()
	if  node.Type() ==  KStringCommon {
		key = node.StringRef(ctx)
	} else {
		key = node.StringCopyEsc(ctx)
	}

	// interface{} slot in map bucket
	val = rt.Mapassign_faststr(rt.MapEfaceMapType, mp, key)
	vt = &(*rt.GoEface)(val).Type
	vp = &(*rt.GoEface)(val).Value

	// parse value node
	node = iter.Next()
	switch node.Type() {
		case KObject:
			newSize := node.Object().Len()
			newMp := ctx.efacePool.GetMap(newSize)
			*vt = rt.MapEfaceType
			*vp = newMp
			remain := size - 1
			isObj = true
			if newSize != 0 {
				if remain > 0 {
					ctx.Stack.Push(mp, remain, true)
				}
				mp = newMp
				size = newSize
				goto _object_key;
			}
		case KArray:
			newSize := node.Array().Len()
			if newSize == 0 {
				ctx.efacePool.ConvTSlice(rt.EmptySlice, rt.SliceEfaceType, val)
				break;
			}

			newSp := ctx.efacePool.GetSlice(newSize)
			// pack to []interface{}
			ctx.efacePool.ConvTSlice(rt.GoSlice{
				Ptr: newSp,
				Len: newSize,
				Cap: newSize,
			}, rt.SliceEfaceType, val)
			remain := size - 1
			if remain > 0 {
				ctx.Stack.Push(mp, remain, true)
			}
			val = newSp
			isObj = false
			size = newSize
			goto _arr_val;
		case KStringCommon:
			ctx.efacePool.ConvTstring(node.StringRef(ctx), val)
		case KStringEscaped:
			ctx.efacePool.ConvTstring(node.StringCopyEsc(ctx), val)
		case KTrue:
			rt.ConvTBool(true, (*interface{})(val))
		case KFalse:
			rt.ConvTBool(false, (*interface{})(val))
		case KNull: /* skip */
		case KUint:
			ctx.efacePool.ConvF64(float64(node.U64()), val)
		case KSint: 
			ctx.efacePool.ConvF64(float64(node.I64()), val)
		case KReal: 
			ctx.efacePool.ConvF64(node.F64(), val)
		case KRawNumber:
			ctx.efacePool.ConvTnum(node.Number(ctx), val)
		default: 
			panic("unreachable for as eface")
	}
	
	// check size 
	size -= 1
	if size != 0 {
		goto _object_key;
	}

	parent, size, isObj = ctx.Stack.Pop()

	// parent is empty
	if parent == nil {
		if isObj {
			return rt.GoEface {
				Type: rt.MapEfaceType,
				Value: rootM,
			}.Pack()
		} else {
			ctx.efacePool.ConvTSlice(rootS, rt.SliceEfaceType, (unsafe.Pointer)(&root))
			return root
		}
	}

	// continue to parse parent
	if isObj {
		mp = parent
		goto _object_key;
	} else {
		val = rt.PtrAdd(parent, rt.AnyType.Size)
		goto _arr_val;
	}

_arr_val:
	// interface{} slot in slice
	vt = &(*rt.GoEface)(val).Type
	vp = &(*rt.GoEface)(val).Value

	// parse value node
	node = iter.Next()
	switch node.Type() {
		case KObject:
			newSize := node.Object().Len()
			newMp := ctx.efacePool.GetMap(newSize)
			*vt = rt.MapEfaceType
			*vp = newMp
			remain := size - 1
			if newSize != 0 {
				// push next array elem into stack
				if remain > 0 {
					ctx.Stack.Push(val, remain, false)
				}
				mp = newMp
				size = newSize
				isObj = true
				goto _object_key;
			}
		case KArray:
			newSize := node.Array().Len()
			if newSize == 0 {
				ctx.efacePool.ConvTSlice(rt.EmptySlice, rt.SliceEfaceType, val)
				break;
			}
			
			newSp := ctx.efacePool.GetSlice(newSize)
			// pack to []interface{}
			ctx.efacePool.ConvTSlice(rt.GoSlice {
				Ptr: newSp,
				Len: newSize,
				Cap: newSize,
			}, rt.SliceEfaceType, val)

			remain := size - 1
			if remain > 0 {
				ctx.Stack.Push(val, remain, false)
			}

			val = newSp
			isObj = false
			size = newSize
			goto _arr_val;
		case KStringCommon:
			ctx.efacePool.ConvTstring(node.StringRef(ctx), val)
		case KStringEscaped:
			ctx.efacePool.ConvTstring(node.StringCopyEsc(ctx), val)
		case KTrue:
			rt.ConvTBool(true, (*interface{})(val))
		case KFalse:
			rt.ConvTBool(false, (*interface{})(val))
		case KNull: /* skip */
		case KUint:
			ctx.efacePool.ConvF64(float64(node.U64()), val)
		case KSint: 
			ctx.efacePool.ConvF64(float64(node.I64()), val)
		case KReal: 
			ctx.efacePool.ConvF64(node.F64(), val)
		case KRawNumber:
			ctx.efacePool.ConvTnum(node.Number(ctx), val)
		default: panic("unreachable for as eface")
	}

	// check size 
	size -= 1
	if size != 0 {
		val = rt.PtrAdd(val, rt.AnyType.Size)
		goto _arr_val;
	}


	parent, size, isObj = ctx.Stack.Pop()

	// parent is empty
	if parent == nil {
		if isObj {
			return rt.GoEface {
				Type: rt.MapEfaceType,
				Value: rootM,
			}.Pack()
		} else {
			ctx.efacePool.ConvTSlice(rootS, rt.SliceEfaceType, unsafe.Pointer(&root))
			return root
		}
	}

	// continue to parse parent
	if isObj {
		mp = parent
		goto _object_key;
	} else {
		val = rt.PtrAdd(parent, rt.AnyType.Size)
		goto _arr_val;
	}
}

func (node *Node) AsEfaceFallback(ctx *Context) (interface{}, error) {
	switch node.Type() {
	case KObject:
		obj := node.Object()
		size := obj.Len()
		m := make(map[string]interface{}, size)
		*node = NewNode(obj.Children())
		var gerr, err error
		for i := 0; i < size; i++ {
			key, _ := node.AsStr(ctx)
			*node = NewNode(PtrOffset(node.cptr, 1))
			m[key], err = node.AsEfaceFallback(ctx)
			if gerr == nil && err != nil {
				gerr = err
			}
		}
		return m, gerr
	case KArray:
		arr := node.Array()
		size := arr.Len()
		a := make([]interface{}, size)
		*node = NewNode(arr.Children())
		var gerr, err error
		for i := 0; i < size; i++ {
			a[i], err = node.AsEfaceFallback(ctx)
			if gerr == nil && err != nil {
				gerr = err
			}
		}
		return a, gerr
	case KStringCommon:
		str, _ := node.AsStr(ctx)
		*node = NewNode(PtrOffset(node.cptr, 1))
		return str, nil
	case KStringEscaped:
		str := node.StringCopyEsc(ctx)
		*node = NewNode(PtrOffset(node.cptr, 1))
		return str, nil
	case KTrue:
		*node = NewNode(PtrOffset(node.cptr, 1))
		return true, nil
	case KFalse:
		*node = NewNode(PtrOffset(node.cptr, 1))
		return false, nil
	case KNull:
		*node = NewNode(PtrOffset(node.cptr, 1))
		return nil, nil
	default:
		// use float64
		if ctx.Parser.options & (1 << _F_use_number) != 0 {
			num, ok := node.AsNumber(ctx)
			if !ok {
				// skip the unmatched type
				*node = NewNode(node.Next())
				return nil, newUnmatched(node.Position(), rt.JsonNumberType)
			} else {
				*node = NewNode(PtrOffset(node.cptr, 1))
				return num, nil
			}
		} else if  ctx.Parser.options & (1 << _F_use_int64) != 0 {
			// first try int64
			i, ok := node.AsI64(ctx)
			if ok {
				*node = NewNode(PtrOffset(node.cptr, 1))
				return i, nil
			}

			// is not integer, then use float64
			f, ok := node.AsF64(ctx)
			if ok {
				*node = NewNode(PtrOffset(node.cptr, 1))
				return f, nil
			}
		
			// skip the unmatched type
			*node = NewNode(node.Next())
			return nil, newUnmatched(node.Position(), rt.Int64Type)
		} else {
			num, ok := node.AsF64(ctx)
			if !ok {
				// skip the unmatched type
				*node = NewNode(node.Next())
				return nil, newUnmatched(node.Position(), rt.Float64Type)
			} else {
				*node = NewNode(PtrOffset(node.cptr, 1))
				return num, nil
			}
		}
	}
}

//go:nosplit
func PtrOffset(ptr uintptr, off int64) uintptr {
	return uintptr(int64(ptr) + off * int64(unsafe.Sizeof(node{})))
}
