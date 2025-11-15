// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !safe && !codec.safe && !appengine && go1.21

// minimum of go 1.21 is needed, as that is the minimum for all features and linked functions we need
// - typedmemclr : go1.8
// - mapassign_fastXXX: go1.9
// - clear was added in go1.21
// - unsafe.String(Data): go1.20
// - unsafe.Add: go1.17
// - generics/any: go1.18
// etc

package codec

import (
	"reflect"
	_ "runtime" // needed for go linkname(s)
	"sync/atomic"
	"time"
	"unsafe"
)

// This file has unsafe variants of some helper functions.
// MARKER: See helper_unsafe.go for the usage documentation.
//
// There are a number of helper_*unsafe*.go files.
//
// - helper_unsafe
//   unsafe variants of dependent functions
// - helper_unsafe_compiler_gc (gc)
//   unsafe variants of dependent functions which cannot be shared with gollvm or gccgo
// - helper_not_unsafe_not_gc (gccgo/gollvm or safe)
//   safe variants of functions in helper_unsafe_compiler_gc
// - helper_not_unsafe (safe)
//   safe variants of functions in helper_unsafe
// - helper_unsafe_compiler_not_gc (gccgo, gollvm)
//   unsafe variants of functions/variables which non-standard compilers need
//
// This way, we can judiciously use build tags to include the right set of files
// for any compiler, and make it run optimally in unsafe mode.
//
// As of March 2021, we cannot differentiate whether running with gccgo or gollvm
// using a build constraint, as both satisfy 'gccgo' build tag.
// Consequently, we must use the lowest common denominator to support both.
//
// For reflect.Value code, we decided to do the following:
//    - if we know the kind, we can elide conditional checks for
//      - SetXXX (Int, Uint, String, Bool, etc)
//      - SetLen
//
// We can also optimize many others, incl IsNil, etc
//
// MARKER: Some functions here will not be hit during code coverage runs due to optimizations, e.g.
//   - rvCopySlice:      called by decode if rvGrowSlice did not set new slice into pointer to orig slice.
//                       however, helper_unsafe sets it, so no need to call rvCopySlice later
//   - rvSlice:          same as above
//
// MARKER: Handling flagIndir ----
//
// flagIndir means that the reflect.Value holds a pointer to the data itself.
//
// flagIndir can be set for:
// - references
//   Here, type.IfaceIndir() --> false
//   flagIndir is usually false (except when the value is addressable, where in flagIndir may be true)
// - everything else (numbers, bools, string, slice, struct, etc).
//   Here, type.IfaceIndir() --> true
//   flagIndir is always true
//
// This knowledge is used across this file, e.g. in rv2i and rvRefPtr

const safeMode = false

// helperUnsafeDirectAssignMapEntry says that we should not copy the pointer in the map
// to another value during mapRange/iteration and mapGet calls, but directly assign it.
//
// The only callers of mapRange/iteration is encode.
// Here, we just walk through the values and encode them
//
// The only caller of mapGet is decode.
// Here, it does a Get if the underlying value is a pointer, and decodes into that.
//
// For both users, we are very careful NOT to modify or keep the pointers around.
// Consequently, it is ok for take advantage of the performance that the map is not modified
// during an iteration and we can just "peek" at the internal value" in the map and use it.
const helperUnsafeDirectAssignMapEntry = true

// MARKER: keep in sync with GO_ROOT/src/reflect/value.go
const (
	unsafeFlagStickyRO = 1 << 5
	unsafeFlagEmbedRO  = 1 << 6
	unsafeFlagIndir    = 1 << 7
	unsafeFlagAddr     = 1 << 8
	unsafeFlagRO       = unsafeFlagStickyRO | unsafeFlagEmbedRO
	// unsafeFlagKindMask = (1 << 5) - 1 // 5 bits for 27 kinds (up to 31)
	// unsafeTypeKindDirectIface = 1 << 5
)

// transientSizeMax below is used in TransientAddr as the backing storage.
//
// Must be >= 16 as the maximum size is a complex128 (or string on 64-bit machines).
const transientSizeMax = 64

// should struct/array support internal strings and slices?
// const transientValueHasStringSlice = false

func isTransientType4Size(size uint32) bool { return size <= transientSizeMax }

type unsafeString struct {
	Data unsafe.Pointer
	Len  int
}

type unsafeSlice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type unsafeIntf struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

type unsafeReflectValue struct {
	unsafeIntf
	flag uintptr
}

// keep in sync with stdlib runtime/type.go
type unsafeRuntimeType struct {
	size uintptr
	// ... many other fields here
}

// unsafeZeroAddr and unsafeZeroSlice points to a read-only block of memory
// used for setting a zero value for most types or creating a read-only
// zero value for a given type.
var (
	unsafeZeroAddr  = unsafe.Pointer(&unsafeZeroArr[0])
	unsafeZeroSlice = unsafeSlice{unsafeZeroAddr, 0, 0}
)

// We use a scratch memory and an unsafeSlice for transient values:
//
// unsafeSlice is used for standalone strings and slices (outside an array or struct).
// scratch memory is used for other kinds, based on contract below:
// - numbers, bool are always transient
// - structs and arrays are transient iff they have no pointers i.e.
//   no string, slice, chan, func, interface, map, etc only numbers and bools.
// - slices and strings are transient (using the unsafeSlice)

type unsafePerTypeElem struct {
	arr   [transientSizeMax]byte // for bool, number, struct, array kinds
	slice unsafeSlice            // for string and slice kinds
}

func (x *unsafePerTypeElem) addrFor(k reflect.Kind) unsafe.Pointer {
	if k == reflect.String || k == reflect.Slice {
		x.slice = unsafeSlice{} // memclr
		return unsafe.Pointer(&x.slice)
	}
	clear(x.arr[:])
	// x.arr = [transientSizeMax]byte{} // memclr
	return unsafe.Pointer(&x.arr)
}

type perType struct {
	elems [2]unsafePerTypeElem
}

type decPerType = perType

type encPerType struct{}

// TransientAddrK is used for getting a *transient* value to be decoded into,
// which will right away be used for something else.
//
// See notes in helper.go about "Transient values during decoding"

func (x *perType) TransientAddrK(t reflect.Type, k reflect.Kind) reflect.Value {
	return rvZeroAddrTransientAnyK(t, k, x.elems[0].addrFor(k))
}

func (x *perType) TransientAddr2K(t reflect.Type, k reflect.Kind) reflect.Value {
	return rvZeroAddrTransientAnyK(t, k, x.elems[1].addrFor(k))
}

func (encPerType) AddressableRO(v reflect.Value) reflect.Value {
	return rvAddressableReadonly(v)
}

// byteAt returns the byte given an index which is guaranteed
// to be within the bounds of the slice i.e. we defensively
// already verified that the index is less than the length of the slice.
func byteAt(b []byte, index uint) byte {
	// return b[index]
	return *(*byte)(unsafe.Pointer(uintptr((*unsafeSlice)(unsafe.Pointer(&b)).Data) + uintptr(index)))
}

func setByteAt(b []byte, index uint, val byte) {
	// b[index] = val
	*(*byte)(unsafe.Pointer(uintptr((*unsafeSlice)(unsafe.Pointer(&b)).Data) + uintptr(index))) = val
}

// stringView returns a view of the []byte as a string.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
func stringView(v []byte) string {
	return *(*string)(unsafe.Pointer(&v))
}

// bytesView returns a view of the string as a []byte.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
func bytesView(v string) (b []byte) {
	sx := (*unsafeString)(unsafe.Pointer(&v))
	bx := (*unsafeSlice)(unsafe.Pointer(&b))
	bx.Data, bx.Len, bx.Cap = sx.Data, sx.Len, sx.Len
	return
}

func byteSliceSameData(v1 []byte, v2 []byte) bool {
	return (*unsafeSlice)(unsafe.Pointer(&v1)).Data == (*unsafeSlice)(unsafe.Pointer(&v2)).Data
}

// isNil checks - without much effort - if an interface is nil.
//
// returned rv is not guaranteed to be valid (e.g. if v == nil).
//
// Note that this will handle all pointer-sized types e.g.
// pointer, map, chan, func, etc.
func isNil(v interface{}, checkPtr bool) (rv reflect.Value, b bool) {
	b = ((*unsafeIntf)(unsafe.Pointer(&v))).ptr == nil
	return
}

func ptrToLowLevel[T any](ptr *T) unsafe.Pointer {
	return unsafe.Pointer(ptr)
}

func lowLevelToPtr[T any](v unsafe.Pointer) *T {
	return (*T)(v)
}

// Given that v is a reference (map/func/chan/ptr/unsafepointer) kind, return the pointer
func rvRefPtr(v *unsafeReflectValue) unsafe.Pointer {
	if v.flag&unsafeFlagIndir != 0 {
		return *(*unsafe.Pointer)(v.ptr)
	}
	return v.ptr
}

func eq4i(i0, i1 interface{}) bool {
	v0 := (*unsafeIntf)(unsafe.Pointer(&i0))
	v1 := (*unsafeIntf)(unsafe.Pointer(&i1))
	return v0.typ == v1.typ && v0.ptr == v1.ptr
}

func rv4iptr(i interface{}) (v reflect.Value) {
	// Main advantage here is that it is inlined, nothing escapes to heap, i is never nil
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.unsafeIntf = *(*unsafeIntf)(unsafe.Pointer(&i))
	uv.flag = uintptr(rkindPtr)
	return
}

func rv4istr(i interface{}) (v reflect.Value) {
	// Main advantage here is that it is inlined, nothing escapes to heap, i is never nil
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.unsafeIntf = *(*unsafeIntf)(unsafe.Pointer(&i))
	uv.flag = uintptr(rkindString) | unsafeFlagIndir
	return
}

func rv2i(rv reflect.Value) (i interface{}) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if refBitset.isset(byte(rv.Kind())) && urv.flag&unsafeFlagIndir != 0 {
		urv.ptr = *(*unsafe.Pointer)(urv.ptr)
	}
	return *(*interface{})(unsafe.Pointer(&urv.unsafeIntf))
}

func rvAddr(rv reflect.Value, ptrType reflect.Type) reflect.Value {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.flag = (urv.flag & unsafeFlagRO) | uintptr(reflect.Ptr)
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&ptrType))).ptr
	return rv
}

// return true if this rv - got from a pointer kind - is nil.
// For now, only use for struct fields of pointer types, as we're guaranteed
// that flagIndir will never be set.
func rvPtrIsNil(rv reflect.Value) bool {
	return rvIsNil(rv)
}

// checks if a nil'able value is nil
func rvIsNil(rv reflect.Value) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if urv.flag&unsafeFlagIndir == 0 {
		return urv.ptr == nil
	}
	// flagIndir is set for a reference (ptr/map/func/unsafepointer/chan)
	// OR kind is slice/interface
	return *(*unsafe.Pointer)(urv.ptr) == nil
}

func rvSetSliceLen(rv reflect.Value, length int) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	(*unsafeString)(urv.ptr).Len = length
}

func rvZeroAddrK(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
	urv.ptr = unsafeNew(urv.typ)
	return
}

func rvZeroAddrTransientAnyK(t reflect.Type, k reflect.Kind, addr unsafe.Pointer) (rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
	urv.ptr = addr
	return
}

func rvZeroK(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	if refBitset.isset(byte(k)) {
		urv.flag = uintptr(k)
	} else if rtsize2(urv.typ) <= uintptr(len(unsafeZeroArr)) {
		urv.flag = uintptr(k) | unsafeFlagIndir
		urv.ptr = unsafeZeroAddr
	} else { // meaning struct or array
		urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
		urv.ptr = unsafeNew(urv.typ)
	}
	return
}

// rvConvert will convert a value to a different type directly,
// ensuring that they still point to the same underlying value.
func rvConvert(v reflect.Value, t reflect.Type) reflect.Value {
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	return v
}

// rvAddressableReadonly returns an addressable reflect.Value.
//
// Use it within encode calls, when you just want to "read" the underlying ptr
// without modifying the value.
//
// Note that it cannot be used for r/w use, as those non-addressable values
// may have been stored in read-only memory, and trying to write the pointer
// may cause a segfault.
func rvAddressableReadonly(v reflect.Value) reflect.Value {
	// hack to make an addressable value out of a non-addressable one.
	// Assume folks calling it are passing a value that can be addressable, but isn't.
	// This assumes that the flagIndir is already set on it.
	// so we just set the flagAddr bit on the flag (and do not set the flagIndir).

	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.flag = uv.flag | unsafeFlagAddr // | unsafeFlagIndir

	return v
}

func rtsize2(rt unsafe.Pointer) uintptr {
	return ((*unsafeRuntimeType)(rt)).size
}

func rt2id(rt reflect.Type) uintptr {
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&rt))).ptr)
}

func i2rtid(i interface{}) uintptr {
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&i))).typ)
}

// --------------------------

func unsafeCmpZero(ptr unsafe.Pointer, size int) bool {
	// verified that size is always within right range, so no chance of OOM
	var s1 = unsafeString{ptr, size}
	var s2 = unsafeString{unsafeZeroAddr, size}
	if size > len(unsafeZeroArr) {
		arr := make([]byte, size)
		s2.Data = unsafe.Pointer(&arr[0])
	}
	return *(*string)(unsafe.Pointer(&s1)) == *(*string)(unsafe.Pointer(&s2)) // memcmp
}

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, recursive bool) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if urv.flag == 0 {
		return true
	}
	if recursive {
		return isEmptyValueFallbackRecur(urv, v, tinfos)
	}
	return unsafeCmpZero(urv.ptr, int(rtsize2(urv.typ)))
}

func isEmptyValueFallbackRecur(urv *unsafeReflectValue, v reflect.Value, tinfos *TypeInfos) bool {
	const recursive = true

	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String:
		return (*unsafeString)(urv.ptr).Len == 0
	case reflect.Slice:
		return (*unsafeSlice)(urv.ptr).Len == 0
	case reflect.Bool:
		return !*(*bool)(urv.ptr)
	case reflect.Int:
		return *(*int)(urv.ptr) == 0
	case reflect.Int8:
		return *(*int8)(urv.ptr) == 0
	case reflect.Int16:
		return *(*int16)(urv.ptr) == 0
	case reflect.Int32:
		return *(*int32)(urv.ptr) == 0
	case reflect.Int64:
		return *(*int64)(urv.ptr) == 0
	case reflect.Uint:
		return *(*uint)(urv.ptr) == 0
	case reflect.Uint8:
		return *(*uint8)(urv.ptr) == 0
	case reflect.Uint16:
		return *(*uint16)(urv.ptr) == 0
	case reflect.Uint32:
		return *(*uint32)(urv.ptr) == 0
	case reflect.Uint64:
		return *(*uint64)(urv.ptr) == 0
	case reflect.Uintptr:
		return *(*uintptr)(urv.ptr) == 0
	case reflect.Float32:
		return *(*float32)(urv.ptr) == 0
	case reflect.Float64:
		return *(*float64)(urv.ptr) == 0
	case reflect.Complex64:
		return unsafeCmpZero(urv.ptr, 8)
	case reflect.Complex128:
		return unsafeCmpZero(urv.ptr, 16)
	case reflect.Struct:
		// return isEmptyStruct(v, tinfos, recursive)
		if tinfos == nil {
			tinfos = defTypeInfos
		}
		ti := tinfos.find(uintptr(urv.typ))
		if ti == nil {
			ti = tinfos.load(v.Type())
		}
		return unsafeCmpZero(urv.ptr, int(ti.size))
	case reflect.Interface, reflect.Ptr:
		// isnil := urv.ptr == nil // (not sufficient, as a pointer value encodes the type)
		isnil := urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
		if recursive && !isnil {
			return isEmptyValue(v.Elem(), tinfos, recursive)
		}
		return isnil
	case reflect.UnsafePointer:
		return urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
	case reflect.Chan:
		return urv.ptr == nil || len_chan(rvRefPtr(urv)) == 0
	case reflect.Map:
		return urv.ptr == nil || len_map(rvRefPtr(urv)) == 0
	case reflect.Array:
		return v.Len() == 0 ||
			urv.ptr == nil ||
			urv.typ == nil ||
			rtsize2(urv.typ) == 0 ||
			unsafeCmpZero(urv.ptr, int(rtsize2(urv.typ)))
	}
	return false
}

// is this an empty interface/ptr/struct/map/slice/chan/array
func isEmptyContainerValue(v reflect.Value, tinfos *TypeInfos, recursive bool) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	switch v.Kind() {
	case reflect.Slice:
		return (*unsafeSlice)(urv.ptr).Len == 0
	case reflect.Struct:
		if tinfos == nil {
			tinfos = defTypeInfos
		}
		ti := tinfos.find(uintptr(urv.typ))
		if ti == nil {
			ti = tinfos.load(v.Type())
		}
		return unsafeCmpZero(urv.ptr, int(ti.size))
	case reflect.Interface, reflect.Ptr:
		// isnil := urv.ptr == nil // (not sufficient, as a pointer value encodes the type)
		isnil := urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
		if recursive && !isnil {
			return isEmptyValue(v.Elem(), tinfos, recursive)
		}
		return isnil
	case reflect.Chan:
		return urv.ptr == nil || len_chan(rvRefPtr(urv)) == 0
	case reflect.Map:
		return urv.ptr == nil || len_map(rvRefPtr(urv)) == 0
	case reflect.Array:
		return v.Len() == 0 ||
			urv.ptr == nil ||
			urv.typ == nil ||
			rtsize2(urv.typ) == 0 ||
			unsafeCmpZero(urv.ptr, int(rtsize2(urv.typ)))
	}
	return false
}

// --------------------------

type structFieldInfos struct {
	c unsafe.Pointer // source
	s unsafe.Pointer // sorted
	t uint8To32TrieNode

	length int

	// byName map[string]*structFieldInfo // find sfi given a name
}

// func (x *structFieldInfos) load(source, sorted []*structFieldInfo, sourceNames, sortedNames []string) {
func (x *structFieldInfos) load(source, sorted []*structFieldInfo) {
	var s *unsafeSlice
	s = (*unsafeSlice)(unsafe.Pointer(&source))
	x.c = s.Data
	x.length = s.Len
	s = (*unsafeSlice)(unsafe.Pointer(&sorted))
	x.s = s.Data
}

func (x *structFieldInfos) source() (v []*structFieldInfo) {
	*(*unsafeSlice)(unsafe.Pointer(&v)) = unsafeSlice{x.c, x.length, x.length}
	return
}

func (x *structFieldInfos) sorted() (v []*structFieldInfo) {
	*(*unsafeSlice)(unsafe.Pointer(&v)) = unsafeSlice{x.s, x.length, x.length}
	return
}

// --------------------------

type uint8To32TrieNodeNoKids struct {
	key     uint8
	valid   bool // the value marks the end of a full stored string
	numkids uint8
	_       byte // padding
	value   uint32
}

type uint8To32TrieNodeKids = *uint8To32TrieNode

func (x *uint8To32TrieNode) setKids(kids []uint8To32TrieNode) {
	x.numkids = uint8(len(kids))
	x.kids = &kids[0]
}
func (x *uint8To32TrieNode) getKids() (v []uint8To32TrieNode) {
	*(*unsafeSlice)(unsafe.Pointer(&v)) = unsafeSlice{unsafe.Pointer(x.kids), int(x.numkids), int(x.numkids)}
	return
}
func (x *uint8To32TrieNode) truncKids() { x.numkids = 0 }

// --------------------------

// Note that we do not atomically load/store length and data pointer separately,
// as this could lead to some races. Instead, we atomically load/store cappedSlice.

type atomicRtidFnSlice struct {
	v unsafe.Pointer // *[]codecRtidFn
}

func (x *atomicRtidFnSlice) load() (s unsafe.Pointer) {
	return atomic.LoadPointer(&x.v)
}

func (x *atomicRtidFnSlice) store(p unsafe.Pointer) {
	atomic.StorePointer(&x.v, p)
}

// --------------------------

// to create a reflect.Value for each member field of fauxUnion,
// we first create a global fauxUnion, and create reflect.Value
// for them all.
// This way, we have the flags and type in the reflect.Value.
// Then, when a reflect.Value is called, we just copy it,
// update the ptr to the fauxUnion's, and return it.

type unsafeDecNakedWrapper struct {
	fauxUnion
	ru, ri, rf, rl, rs, rb, rt reflect.Value // mapping to the primitives above
}

func (n *unsafeDecNakedWrapper) init() {
	n.ru = rv4iptr(&n.u).Elem()
	n.ri = rv4iptr(&n.i).Elem()
	n.rf = rv4iptr(&n.f).Elem()
	n.rl = rv4iptr(&n.l).Elem()
	n.rs = rv4iptr(&n.s).Elem()
	n.rt = rv4iptr(&n.t).Elem()
	n.rb = rv4iptr(&n.b).Elem()
	// n.rr[] = reflect.ValueOf(&n.)
}

var defUnsafeDecNakedWrapper unsafeDecNakedWrapper

func init() {
	defUnsafeDecNakedWrapper.init()
}

func (n *fauxUnion) ru() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.ru
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.u)
	return
}
func (n *fauxUnion) ri() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.ri
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.i)
	return
}
func (n *fauxUnion) rf() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rf
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.f)
	return
}
func (n *fauxUnion) rl() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rl
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.l)
	return
}
func (n *fauxUnion) rs() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rs
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.s)
	return
}
func (n *fauxUnion) rt() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rt
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.t)
	return
}
func (n *fauxUnion) rb() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rb
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.b)
	return
}

// --------------------------
func rvSetBytes(rv reflect.Value, v []byte) {
	*(*[]byte)(rvPtr(rv)) = v
}

func rvSetString(rv reflect.Value, v string) {
	*(*string)(rvPtr(rv)) = v
}

func rvSetBool(rv reflect.Value, v bool) {
	*(*bool)(rvPtr(rv)) = v
}

func rvSetTime(rv reflect.Value, v time.Time) {
	*(*time.Time)(rvPtr(rv)) = v
}

func rvSetFloat32(rv reflect.Value, v float32) {
	*(*float32)(rvPtr(rv)) = v
}

func rvSetFloat64(rv reflect.Value, v float64) {
	*(*float64)(rvPtr(rv)) = v
}

func rvSetComplex64(rv reflect.Value, v complex64) {
	*(*complex64)(rvPtr(rv)) = v
}

func rvSetComplex128(rv reflect.Value, v complex128) {
	*(*complex128)(rvPtr(rv)) = v
}

func rvSetInt(rv reflect.Value, v int) {
	*(*int)(rvPtr(rv)) = v
}

func rvSetInt8(rv reflect.Value, v int8) {
	*(*int8)(rvPtr(rv)) = v
}

func rvSetInt16(rv reflect.Value, v int16) {
	*(*int16)(rvPtr(rv)) = v
}

func rvSetInt32(rv reflect.Value, v int32) {
	*(*int32)(rvPtr(rv)) = v
}

func rvSetInt64(rv reflect.Value, v int64) {
	*(*int64)(rvPtr(rv)) = v
}

func rvSetUint(rv reflect.Value, v uint) {
	*(*uint)(rvPtr(rv)) = v
}

func rvSetUintptr(rv reflect.Value, v uintptr) {
	*(*uintptr)(rvPtr(rv)) = v
}

func rvSetUint8(rv reflect.Value, v uint8) {
	*(*uint8)(rvPtr(rv)) = v
}

func rvSetUint16(rv reflect.Value, v uint16) {
	*(*uint16)(rvPtr(rv)) = v
}

func rvSetUint32(rv reflect.Value, v uint32) {
	*(*uint32)(rvPtr(rv)) = v
}

func rvSetUint64(rv reflect.Value, v uint64) {
	*(*uint64)(rvPtr(rv)) = v
}

// ----------------

// rvSetZero is rv.Set(reflect.Zero(rv.Type()) for all kinds (including reflect.Interface).
func rvSetZero(rv reflect.Value) {
	rvSetDirectZero(rv)
}

func rvSetIntf(rv reflect.Value, v reflect.Value) {
	rv.Set(v)
}

// rvSetDirect is rv.Set for all kinds except reflect.Interface.
//
// Callers MUST not pass a value of kind reflect.Interface, as it may cause unexpected segfaults.
func rvSetDirect(rv reflect.Value, v reflect.Value) {
	// MARKER: rv.Set for kind reflect.Interface may do a separate allocation if a scalar value.
	// The book-keeping is onerous, so we just do the simple ones where a memmove is sufficient.
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if uv.flag&unsafeFlagIndir == 0 {
		*(*unsafe.Pointer)(urv.ptr) = uv.ptr
	} else if uv.ptr != unsafeZeroAddr {
		typedmemmove(urv.typ, urv.ptr, uv.ptr)
	} else if urv.ptr != unsafeZeroAddr {
		typedmemclr(urv.typ, urv.ptr)
	}
}

// rvSetDirectZero is rv.Set(reflect.Zero(rv.Type()) for all kinds except reflect.Interface.
func rvSetDirectZero(rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if urv.ptr != unsafeZeroAddr {
		typedmemclr(urv.typ, urv.ptr)
	}
}

// rvMakeSlice updates the slice to point to a new array.
// It copies data from old slice to new slice.
// It returns set=true iff it updates it, else it just returns a new slice pointing to a newly made array.
func rvMakeSlice(rv reflect.Value, ti *typeInfo, xlen, xcap int) (_ reflect.Value, set bool) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	ux := (*unsafeSlice)(urv.ptr)
	t := ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	s := unsafeSlice{newarray(t, xcap), xlen, xcap}
	if ux.Len > 0 {
		typedslicecopy(t, s, *ux)
	}
	*ux = s
	return rv, true
}

// rvSlice returns a sub-slice of the slice given new lenth,
// without modifying passed in value.
// It is typically called when we know that SetLen(...) cannot be done.
func rvSlice(rv reflect.Value, length int) reflect.Value {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	ux := *(*unsafeSlice)(urv.ptr) // copy slice header
	ux.Len = length
	urv.ptr = unsafe.Pointer(&ux)
	return rv
}

// rcGrowSlice updates the slice to point to a new array with the cap incremented, and len set to the new cap value.
// It copies data from old slice to new slice.
// It returns set=true iff it updates it, else it just returns a new slice pointing to a newly made array.
func rvGrowSlice(rv reflect.Value, ti *typeInfo, cap, incr int) (v reflect.Value, newcap int, set bool) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	ux := (*unsafeSlice)(urv.ptr)
	t := ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	*ux = unsafeGrowslice(t, *ux, cap, incr)
	ux.Len = ux.Cap
	return rv, ux.Cap, true
}

// ------------

func rvArrayIndex(rv reflect.Value, i int, ti *typeInfo, isSlice bool) (v reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if isSlice {
		uv.ptr = unsafe.Pointer(uintptr(((*unsafeSlice)(urv.ptr)).Data))
	} else {
		uv.ptr = unsafe.Pointer(uintptr(urv.ptr))
	}
	uv.ptr = unsafe.Add(uv.ptr, ti.elemsize*uint32(i))
	// uv.ptr = unsafe.Pointer(ptr + uintptr(int(ti.elemsize)*i))
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	uv.flag = uintptr(ti.elemkind) | unsafeFlagIndir | unsafeFlagAddr
	return
}

func rvSliceZeroCap(t reflect.Type) (v reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	urv.flag = uintptr(reflect.Slice) | unsafeFlagIndir
	urv.ptr = unsafe.Pointer(&unsafeZeroSlice)
	return
}

func rvLenSlice(rv reflect.Value) int {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return (*unsafeSlice)(urv.ptr).Len
}

func rvCapSlice(rv reflect.Value) int {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return (*unsafeSlice)(urv.ptr).Cap
}

// if scratch is nil, then return a writable view (assuming canAddr=true)
func rvGetArrayBytes(rv reflect.Value, _ []byte) (bs []byte) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	bx := (*unsafeSlice)(unsafe.Pointer(&bs))
	// bx.Data, bx.Len, bx.Cap = urv.ptr, rv.Len(), bx.Len
	bx.Data = urv.ptr
	bx.Len = rv.Len()
	bx.Cap = bx.Len
	return
}

func rvGetArray4Slice(rv reflect.Value) (v reflect.Value) {
	// It is possible that this slice is based off an array with a larger
	// len that we want (where array len == slice cap).
	// However, it is ok to create an array type that is a subset of the full
	// e.g. full slice is based off a *[16]byte, but we can create a *[4]byte
	// off of it. That is ok.
	//
	// Consequently, we use rvLenSlice, not rvCapSlice.

	t := reflect.ArrayOf(rvLenSlice(rv), rv.Type().Elem())
	// v = rvZeroAddrK(t, reflect.Array)

	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.flag = uintptr(reflect.Array) | unsafeFlagIndir | unsafeFlagAddr
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr

	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv.ptr = *(*unsafe.Pointer)(urv.ptr) // slice rv has a ptr to the slice.

	return
}

func rvGetSlice4Array(rv reflect.Value, v interface{}) {
	// v is a pointer to a slice to be populated
	uv := (*unsafeIntf)(unsafe.Pointer(&v))
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

	s := (*unsafeSlice)(uv.ptr)
	s.Data = urv.ptr
	s.Len = rv.Len()
	s.Cap = s.Len
}

func rvCopySlice(dest, src reflect.Value, elemType reflect.Type) {
	typedslicecopy((*unsafeIntf)(unsafe.Pointer(&elemType)).ptr,
		*(*unsafeSlice)((*unsafeReflectValue)(unsafe.Pointer(&dest)).ptr),
		*(*unsafeSlice)((*unsafeReflectValue)(unsafe.Pointer(&src)).ptr))
}

// ------------

func rvPtr(rv reflect.Value) unsafe.Pointer {
	return (*unsafeReflectValue)(unsafe.Pointer(&rv)).ptr
}

func rvGetBool(rv reflect.Value) bool {
	return *(*bool)(rvPtr(rv))
}

func rvGetBytes(rv reflect.Value) []byte {
	return *(*[]byte)(rvPtr(rv))
}

func rvGetTime(rv reflect.Value) time.Time {
	return *(*time.Time)(rvPtr(rv))
}

func rvGetString(rv reflect.Value) string {
	return *(*string)(rvPtr(rv))
}

func rvGetFloat64(rv reflect.Value) float64 {
	return *(*float64)(rvPtr(rv))
}

func rvGetFloat32(rv reflect.Value) float32 {
	return *(*float32)(rvPtr(rv))
}

func rvGetComplex64(rv reflect.Value) complex64 {
	return *(*complex64)(rvPtr(rv))
}

func rvGetComplex128(rv reflect.Value) complex128 {
	return *(*complex128)(rvPtr(rv))
}

func rvGetInt(rv reflect.Value) int {
	return *(*int)(rvPtr(rv))
}

func rvGetInt8(rv reflect.Value) int8 {
	return *(*int8)(rvPtr(rv))
}

func rvGetInt16(rv reflect.Value) int16 {
	return *(*int16)(rvPtr(rv))
}

func rvGetInt32(rv reflect.Value) int32 {
	return *(*int32)(rvPtr(rv))
}

func rvGetInt64(rv reflect.Value) int64 {
	return *(*int64)(rvPtr(rv))
}

func rvGetUint(rv reflect.Value) uint {
	return *(*uint)(rvPtr(rv))
}

func rvGetUint8(rv reflect.Value) uint8 {
	return *(*uint8)(rvPtr(rv))
}

func rvGetUint16(rv reflect.Value) uint16 {
	return *(*uint16)(rvPtr(rv))
}

func rvGetUint32(rv reflect.Value) uint32 {
	return *(*uint32)(rvPtr(rv))
}

func rvGetUint64(rv reflect.Value) uint64 {
	return *(*uint64)(rvPtr(rv))
}

func rvGetUintptr(rv reflect.Value) uintptr {
	return *(*uintptr)(rvPtr(rv))
}

func rvLenMap(rv reflect.Value) int {
	// maplen is not inlined, because as of go1.16beta, go:linkname's are not inlined.
	// thus, faster to call rv.Len() directly.
	//
	// MARKER: review after https://github.com/golang/go/issues/20019 fixed.

	// return rv.Len()

	return len_map(rvRefPtr((*unsafeReflectValue)(unsafe.Pointer(&rv))))
}

// Note: it is hard to find len(...) of an array type,
// as that is a field in the arrayType representing the array, and hard to introspect.
//
// func rvLenArray(rv reflect.Value) int {	return rv.Len() }

// ------------ map range and map indexing ----------

// regular calls to map via reflection: MapKeys, MapIndex, MapRange/MapIter etc
// will always allocate for each map key or value.
//
// It is more performant to provide a value that the map entry is set into,
// and that elides the allocation.
//
// go 1.4 through go 1.23 (in runtime/hashmap.go or runtime/map.go) has a hIter struct
// with the first 2 values being pointers for key and value of the current iteration.
// The next 6 values are pointers, followed by numeric types (uintptr, uint8, bool, etc).
// This *hIter is passed to mapiterinit, mapiternext, mapiterkey, mapiterelem.
//
// In go 1.24, swissmap was introduced, and it provides a compatibility layer
// for hIter (called linknameIter). This has only 2 pointer fields after the key and value pointers.
//
// Note: We bypass the reflect wrapper functions and just use the *hIter directly.
//
// When 'faking' these types with our own, we MUST ensure that the GC sees the pointers
// appropriately. These are reflected in goversion_(no)swissmap_unsafe.go files.
// In these files, we pad the extra spaces appropriately.
//
// Note: the faux hIter/linknameIter is directly embedded in unsafeMapIter below

type unsafeMapIter struct {
	mtyp, mptr unsafe.Pointer
	k, v       unsafeReflectValue
	kisref     bool
	visref     bool
	mapvalues  bool
	done       bool
	started    bool
	_          [3]byte // padding
	it         struct {
		key   unsafe.Pointer
		value unsafe.Pointer
		_     unsafeMapIterPadding
	}
}

func (t *unsafeMapIter) Next() (r bool) {
	if t == nil || t.done {
		return
	}
	if t.started {
		mapiternext((unsafe.Pointer)(&t.it))
	} else {
		t.started = true
	}

	t.done = t.it.key == nil
	if t.done {
		return
	}

	if helperUnsafeDirectAssignMapEntry || t.kisref {
		t.k.ptr = t.it.key
	} else {
		typedmemmove(t.k.typ, t.k.ptr, t.it.key)
	}

	if t.mapvalues {
		if helperUnsafeDirectAssignMapEntry || t.visref {
			t.v.ptr = t.it.value
		} else {
			typedmemmove(t.v.typ, t.v.ptr, t.it.value)
		}
	}

	return true
}

func (t *unsafeMapIter) Key() (r reflect.Value) {
	return *(*reflect.Value)(unsafe.Pointer(&t.k))
}

func (t *unsafeMapIter) Value() (r reflect.Value) {
	return *(*reflect.Value)(unsafe.Pointer(&t.v))
}

func (t *unsafeMapIter) Done() {}

type mapIter struct {
	unsafeMapIter
}

func mapRange(t *mapIter, m, k, v reflect.Value, mapvalues bool) {
	if rvIsNil(m) {
		t.done = true
		return
	}
	t.done = false
	t.started = false
	t.mapvalues = mapvalues

	// var urv *unsafeReflectValue

	urv := (*unsafeReflectValue)(unsafe.Pointer(&m))
	t.mtyp = urv.typ
	t.mptr = rvRefPtr(urv)

	// t.it = (*unsafeMapHashIter)(reflect_mapiterinit(t.mtyp, t.mptr))
	mapiterinit(t.mtyp, t.mptr, unsafe.Pointer(&t.it))

	t.k = *(*unsafeReflectValue)(unsafe.Pointer(&k))
	t.kisref = refBitset.isset(byte(k.Kind()))

	if mapvalues {
		t.v = *(*unsafeReflectValue)(unsafe.Pointer(&v))
		t.visref = refBitset.isset(byte(v.Kind()))
	} else {
		t.v = unsafeReflectValue{}
	}
}

// unsafeMapKVPtr returns the pointer if flagIndir, else it returns a pointer to the pointer.
// It is needed as maps always keep a reference to the underlying value.
func unsafeMapKVPtr(urv *unsafeReflectValue) unsafe.Pointer {
	if urv.flag&unsafeFlagIndir == 0 {
		return unsafe.Pointer(&urv.ptr)
	}
	return urv.ptr
}

// return an addressable reflect value that can be used in mapRange and mapGet operations.
//
// all calls to mapGet or mapRange will call here to get an addressable reflect.Value.
func mapAddrLoopvarRV(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	// return rvZeroAddrK(t, k)
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	// since we always set the ptr when helperUnsafeDirectAssignMapEntry=true,
	// we should only allocate if it is not true
	if !helperUnsafeDirectAssignMapEntry {
		urv.ptr = unsafeNew(urv.typ)
	}
	return
}

func makeMapReflect(typ reflect.Type, size int) (rv reflect.Value) {
	t := (*unsafeIntf)(unsafe.Pointer(&typ)).ptr
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = t
	urv.flag = uintptr(reflect.Map)
	urv.ptr = makemap(t, size, nil)
	return
}

func (d *decoderBase) bytes2Str(in []byte, state dBytesAttachState) (s string, mutable bool) {
	return stringView(in), state <= dBytesAttachBuffer
}

// ---------- structFieldInfo optimized ---------------

func (n *structFieldInfoNode) rvField(v reflect.Value) (rv reflect.Value) {
	// we already know this is exported, and maybe embedded (based on what si says)
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))

	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	// clear flagEmbedRO if necessary, and inherit permission bits from v
	urv.flag = uv.flag&(unsafeFlagStickyRO|unsafeFlagIndir|unsafeFlagAddr) | uintptr(n.kind)
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&n.typ))).ptr
	urv.ptr = unsafe.Pointer(uintptr(uv.ptr) + uintptr(n.offset))

	// *(*unsafeReflectValue)(unsafe.Pointer(&rv)) = unsafeReflectValue{
	// 	unsafeIntf: unsafeIntf{
	// 		typ: ((*unsafeIntf)(unsafe.Pointer(&n.typ))).ptr,
	// 		ptr: unsafe.Pointer(uintptr(uv.ptr) + uintptr(n.offset)),
	// 	},
	// 	flag: uv.flag&(unsafeFlagStickyRO|unsafeFlagIndir|unsafeFlagAddr) | uintptr(n.kind),
	// }

	return
}

// runtime chan and map are designed such that the first field is the count.
// len builtin uses this to get the length of a chan/map easily.
// leverage this knowledge, since maplen and chanlen functions from runtime package
// are go:linkname'd here, and thus not inlined as of go1.16beta

func len_map_chan(m unsafe.Pointer) int {
	if m == nil {
		return 0
	}
	return *((*int)(m))
}

func len_map(m unsafe.Pointer) int {
	// return maplen(m)
	return len_map_chan(m)
}
func len_chan(m unsafe.Pointer) int {
	// return chanlen(m)
	return len_map_chan(m)
}

func unsafeNew(typ unsafe.Pointer) unsafe.Pointer {
	return mallocgc(rtsize2(typ), typ, true)
}

// ---------- go linknames (LINKED to runtime/reflect) ---------------

// MARKER: always check that these linknames match subsequent versions of go
//
// Note that as of Jan 2021 (go 1.16 release), go:linkname(s) are not inlined
// outside of the standard library use (e.g. within sync, reflect, etc).
// If these link'ed functions were normally inlined, calling them here would
// not necessarily give a performance boost, due to function overhead.
//
// However, it seems most of these functions are not inlined anyway,
// as only maplen, chanlen and mapaccess are small enough to get inlined.
//
//   We checked this by going into $GOROOT/src/runtime and running:
//   $ go build -tags codec.notfastpath -gcflags "-m=2"

// reflect.{unsafe_New, unsafe_NewArray} are not supported in gollvm,
// failing with "error: undefined reference" error.
// however, runtime.{mallocgc, newarray} are supported, so use that instead.

//go:linkname mallocgc runtime.mallocgc
//go:noescape
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

//go:linkname newarray runtime.newarray
//go:noescape
func newarray(typ unsafe.Pointer, n int) unsafe.Pointer

//go:linkname mapiterinit runtime.mapiterinit
//go:noescape
func mapiterinit(typ unsafe.Pointer, m unsafe.Pointer, it unsafe.Pointer)

//go:linkname mapiternext runtime.mapiternext
//go:noescape
func mapiternext(it unsafe.Pointer) (key unsafe.Pointer)

//go:linkname mapassign runtime.mapassign
//go:noescape
func mapassign(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer

//go:linkname mapaccess2 runtime.mapaccess2
//go:noescape
func mapaccess2(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) (val unsafe.Pointer, ok bool)

//go:linkname makemap runtime.makemap
//go:noescape
func makemap(typ unsafe.Pointer, size int, h unsafe.Pointer) unsafe.Pointer

// reflect.typed{memmove, memclr, slicecopy} will handle checking if the type has pointers or not,
// and if a writeBarrier is needed, before delegating to the right method in the runtime.
//
// This is why we use the functions in reflect, and not the ones in runtime directly.
// Calling runtime.XXX here will lead to memory issues.

//go:linkname typedslicecopy reflect.typedslicecopy
//go:noescape
func typedslicecopy(elemType unsafe.Pointer, dst, src unsafeSlice) int

//go:linkname typedmemmove reflect.typedmemmove
//go:noescape
func typedmemmove(typ unsafe.Pointer, dst, src unsafe.Pointer)

//go:linkname typedmemclr reflect.typedmemclr
//go:noescape
func typedmemclr(typ unsafe.Pointer, dst unsafe.Pointer)
