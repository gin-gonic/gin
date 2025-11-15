// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !go1.9 || safe || codec.safe || appengine

package codec

import (
	// "hash/adler32"
	"math"
	"reflect"
	"sync/atomic"
	"time"
)

// This file has safe variants of some helper functions.
// MARKER: See helper_unsafe.go for the usage documentation.

const safeMode = true

func isTransientType4Size(size uint32) bool { return true }

type mapReqParams struct{}

func getMapReqParams(ti *typeInfo) (r mapReqParams) { return }

func byteAt(b []byte, index uint) byte {
	return b[index]
}

func setByteAt(b []byte, index uint, val byte) {
	b[index] = val
}

func stringView(v []byte) string {
	return string(v)
}

func bytesView(v string) []byte {
	return []byte(v)
}

func byteSliceSameData(v1 []byte, v2 []byte) bool {
	return cap(v1) != 0 && cap(v2) != 0 && &(v1[:1][0]) == &(v2[:1][0])
}

func isNil(v interface{}, checkPtr bool) (rv reflect.Value, b bool) {
	b = v == nil
	if b || !checkPtr {
		return
	}
	rv = reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		b = rv.IsNil()
	}
	return
}

func ptrToLowLevel(v interface{}) interface{} {
	return v
}

func lowLevelToPtr[T any](v interface{}) *T {
	return v.(*T)
}

func eq4i(i0, i1 interface{}) bool {
	return i0 == i1
}

func rv4iptr(i interface{}) reflect.Value { return reflect.ValueOf(i) }
func rv4istr(i interface{}) reflect.Value { return reflect.ValueOf(i) }

func rv2i(rv reflect.Value) interface{} {
	if rv.IsValid() {
		return rv.Interface()
	}
	return nil
}

func rvAddr(rv reflect.Value, ptrType reflect.Type) reflect.Value {
	return rv.Addr()
}

func rvPtrIsNil(rv reflect.Value) bool {
	return rv.IsNil()
}

func rvIsNil(rv reflect.Value) bool {
	return rv.IsNil()
}

func rvSetSliceLen(rv reflect.Value, length int) {
	rv.SetLen(length)
}

func rvZeroAddrK(t reflect.Type, k reflect.Kind) reflect.Value {
	return reflect.New(t).Elem()
}

func rvZeroK(t reflect.Type, k reflect.Kind) reflect.Value {
	return reflect.Zero(t)
}

func rvConvert(v reflect.Value, t reflect.Type) (rv reflect.Value) {
	// Note that reflect.Value.Convert(...) will make a copy if it is addressable.
	// Since we decode into the passed value, we must try to convert the addressable value..
	if v.CanAddr() {
		return v.Addr().Convert(reflect.PtrTo(t)).Elem()
	}
	return v.Convert(t)
}

func rt2id(rt reflect.Type) uintptr {
	return reflect.ValueOf(rt).Pointer()
}

func i2rtid(i interface{}) uintptr {
	return reflect.ValueOf(reflect.TypeOf(i)).Pointer()
}

// --------------------------

// is this an empty interface/ptr/struct/map/slice/chan/array
func isEmptyContainerValue(v reflect.Value, tinfos *TypeInfos, recursive bool) (empty bool) {
	switch v.Kind() {
	case reflect.Array:
		for i, vlen := 0, v.Len(); i < vlen; i++ {
			if !isEmptyValue(v.Index(i), tinfos, false) {
				return false
			}
		}
		return true
	case reflect.Map, reflect.Slice, reflect.Chan:
		return v.IsNil() || v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		empty = v.IsNil()
		if recursive && !empty {
			return isEmptyValue(v.Elem(), tinfos, recursive)
		}
		return empty
	case reflect.Struct:
		return isEmptyStruct(v, tinfos, recursive)
	}
	return false
}

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, recursive bool) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String:
		return v.Len() == 0
	case reflect.Array:
		// zero := reflect.Zero(v.Type().Elem())
		// can I just check if the whole value is equal to zeros? seems not.
		// can I just check if the whole value is equal to its zero value? no.
		// Well, then we check if each value is empty without recursive.
		for i, vlen := 0, v.Len(); i < vlen; i++ {
			if !isEmptyValue(v.Index(i), tinfos, false) {
				return false
			}
		}
		return true
	case reflect.Map, reflect.Slice, reflect.Chan:
		return v.IsNil() || v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Func, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.Interface, reflect.Ptr:
		isnil := v.IsNil()
		if recursive && !isnil {
			return isEmptyValue(v.Elem(), tinfos, recursive)
		}
		return isnil
	case reflect.Struct:
		return isEmptyStruct(v, tinfos, recursive)
	}
	return false
}

// isEmptyStruct is only called from isEmptyValue, and checks if a struct is empty:
//   - does it implement IsZero() bool
//   - is it comparable, and can i compare directly using ==
//   - if checkStruct, then walk through the encodable fields
//     and check if they are empty or not.
func isEmptyStruct(v reflect.Value, tinfos *TypeInfos, recursive bool) bool {
	// v is a struct kind - no need to check again.
	// We only check isZero on a struct kind, to reduce the amount of times
	// that we lookup the rtid and typeInfo for each type as we walk the tree.

	vt := v.Type()
	rtid := rt2id(vt)
	if tinfos == nil {
		tinfos = defTypeInfos
	}
	ti := tinfos.get(rtid, vt)
	if ti.rtid == timeTypId {
		return rv2i(v).(time.Time).IsZero()
	}
	if ti.flagIsZeroer {
		return rv2i(v).(isZeroer).IsZero()
	}
	if ti.flagIsZeroerPtr && v.CanAddr() {
		return rv2i(v.Addr()).(isZeroer).IsZero()
	}
	if ti.flagIsCodecEmptyer {
		return rv2i(v).(isCodecEmptyer).IsCodecEmpty()
	}
	if ti.flagIsCodecEmptyerPtr && v.CanAddr() {
		return rv2i(v.Addr()).(isCodecEmptyer).IsCodecEmpty()
	}
	if ti.flagComparable {
		return rv2i(v) == rv2i(rvZeroK(vt, reflect.Struct))
	}
	if !recursive {
		return false
	}
	// We only care about what we can encode/decode,
	// so that is what we use to check omitEmpty.
	for _, si := range ti.sfi.source() {
		sfv := si.fieldNoAlloc(v, true)
		if sfv.IsValid() && !isEmptyValue(sfv, tinfos, recursive) {
			return false
		}
	}
	return true
}

func makeMapReflect(t reflect.Type, size int) reflect.Value {
	return reflect.MakeMapWithSize(t, size)
}

// --------------------------

type perTypeElem struct {
	t    reflect.Type
	rtid uintptr
	zero reflect.Value
	addr [2]reflect.Value
}

func (x *perTypeElem) get(index uint8) (v reflect.Value) {
	v = x.addr[index%2]
	if v.IsValid() {
		v.Set(x.zero)
	} else {
		v = reflect.New(x.t).Elem()
		x.addr[index%2] = v
	}
	return
}

type perType struct {
	v []perTypeElem
}

type decPerType = perType

type encPerType = perType

func (x *perType) elem(t reflect.Type) *perTypeElem {
	rtid := rt2id(t)
	var h, i uint
	var j = uint(len(x.v))
LOOP:
	if i < j {
		h = (i + j) >> 1 // avoid overflow when computing h // h = i + (j-i)/2
		if x.v[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(x.v)) {
		if x.v[i].rtid != rtid {
			x.v = append(x.v, perTypeElem{})
			copy(x.v[i+1:], x.v[i:])
			x.v[i] = perTypeElem{t: t, rtid: rtid, zero: reflect.Zero(t)}
		}
	} else {
		x.v = append(x.v, perTypeElem{t: t, rtid: rtid, zero: reflect.Zero(t)})
	}
	return &x.v[i]
}

func (x *perType) TransientAddrK(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	return x.elem(t).get(0)
}

func (x *perType) TransientAddr2K(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	return x.elem(t).get(1)
}

func (x *perType) AddressableRO(v reflect.Value) (rv reflect.Value) {
	rv = x.elem(v.Type()).get(0)
	rvSetDirect(rv, v)
	return
}

// --------------------------
type mapIter struct {
	t      *reflect.MapIter
	m      reflect.Value
	values bool
}

func (t *mapIter) Next() (r bool) {
	return t.t.Next()
}

func (t *mapIter) Key() reflect.Value {
	return t.t.Key()
}

func (t *mapIter) Value() (r reflect.Value) {
	if t.values {
		return t.t.Value()
	}
	return
}

func (t *mapIter) Done() {}

func mapRange(t *mapIter, m, k, v reflect.Value, values bool) {
	*t = mapIter{
		m:      m,
		t:      m.MapRange(),
		values: values,
	}
}

// --------------------------
type structFieldInfos struct {
	c []*structFieldInfo
	s []*structFieldInfo
	t uint8To32TrieNode
	// byName map[string]*structFieldInfo // find sfi given a name
}

func (x *structFieldInfos) load(source, sorted []*structFieldInfo) {
	x.c = source
	x.s = sorted
}

// func (x *structFieldInfos) count() int                     { return len(x.c) }
func (x *structFieldInfos) source() (v []*structFieldInfo) { return x.c }
func (x *structFieldInfos) sorted() (v []*structFieldInfo) { return x.s }

// --------------------------

type uint8To32TrieNodeNoKids struct {
	key   uint8
	valid bool    // the value marks the end of a full stored string
	_     [2]byte // padding
	value uint32
}

type uint8To32TrieNodeKids = []uint8To32TrieNode

func (x *uint8To32TrieNode) setKids(kids []uint8To32TrieNode) { x.kids = kids }
func (x *uint8To32TrieNode) getKids() []uint8To32TrieNode     { return x.kids }
func (x *uint8To32TrieNode) truncKids()                       { x.kids = x.kids[:0] } // set len to 0

// --------------------------
func (n *fauxUnion) ru() reflect.Value {
	return reflect.ValueOf(&n.u).Elem()
}
func (n *fauxUnion) ri() reflect.Value {
	return reflect.ValueOf(&n.i).Elem()
}
func (n *fauxUnion) rf() reflect.Value {
	return reflect.ValueOf(&n.f).Elem()
}
func (n *fauxUnion) rl() reflect.Value {
	return reflect.ValueOf(&n.l).Elem()
}
func (n *fauxUnion) rs() reflect.Value {
	return reflect.ValueOf(&n.s).Elem()
}
func (n *fauxUnion) rt() reflect.Value {
	return reflect.ValueOf(&n.t).Elem()
}
func (n *fauxUnion) rb() reflect.Value {
	return reflect.ValueOf(&n.b).Elem()
}

// --------------------------
func rvSetBytes(rv reflect.Value, v []byte) {
	rv.SetBytes(v)
}

func rvSetString(rv reflect.Value, v string) {
	rv.SetString(v)
}

func rvSetBool(rv reflect.Value, v bool) {
	rv.SetBool(v)
}

func rvSetTime(rv reflect.Value, v time.Time) {
	rv.Set(reflect.ValueOf(v))
}

func rvSetFloat32(rv reflect.Value, v float32) {
	rv.SetFloat(float64(v))
}

func rvSetFloat64(rv reflect.Value, v float64) {
	rv.SetFloat(v)
}

func rvSetComplex64(rv reflect.Value, v complex64) {
	rv.SetComplex(complex128(v))
}

func rvSetComplex128(rv reflect.Value, v complex128) {
	rv.SetComplex(v)
}

func rvSetInt(rv reflect.Value, v int) {
	rv.SetInt(int64(v))
}

func rvSetInt8(rv reflect.Value, v int8) {
	rv.SetInt(int64(v))
}

func rvSetInt16(rv reflect.Value, v int16) {
	rv.SetInt(int64(v))
}

func rvSetInt32(rv reflect.Value, v int32) {
	rv.SetInt(int64(v))
}

func rvSetInt64(rv reflect.Value, v int64) {
	rv.SetInt(v)
}

func rvSetUint(rv reflect.Value, v uint) {
	rv.SetUint(uint64(v))
}

func rvSetUintptr(rv reflect.Value, v uintptr) {
	rv.SetUint(uint64(v))
}

func rvSetUint8(rv reflect.Value, v uint8) {
	rv.SetUint(uint64(v))
}

func rvSetUint16(rv reflect.Value, v uint16) {
	rv.SetUint(uint64(v))
}

func rvSetUint32(rv reflect.Value, v uint32) {
	rv.SetUint(uint64(v))
}

func rvSetUint64(rv reflect.Value, v uint64) {
	rv.SetUint(v)
}

// ----------------

func rvSetDirect(rv reflect.Value, v reflect.Value) {
	rv.Set(v)
}

func rvSetDirectZero(rv reflect.Value) {
	rv.Set(reflect.Zero(rv.Type()))
}

// func rvSet(rv reflect.Value, v reflect.Value) {
// 	rv.Set(v)
// }

func rvSetIntf(rv reflect.Value, v reflect.Value) {
	rv.Set(v)
}

func rvSetZero(rv reflect.Value) {
	rv.Set(reflect.Zero(rv.Type()))
}

func rvSlice(rv reflect.Value, length int) reflect.Value {
	return rv.Slice(0, length)
}

func rvMakeSlice(rv reflect.Value, ti *typeInfo, xlen, xcap int) (v reflect.Value, set bool) {
	v = reflect.MakeSlice(ti.rt, xlen, xcap)
	if rv.Len() > 0 {
		reflect.Copy(v, rv)
	}
	return
}

func rvGrowSlice(rv reflect.Value, ti *typeInfo, cap, incr int) (v reflect.Value, newcap int, set bool) {
	newcap = int(growCap(uint(cap), uint(ti.elemsize), uint(incr)))
	v = reflect.MakeSlice(ti.rt, newcap, newcap)
	if rv.Len() > 0 {
		reflect.Copy(v, rv)
	}
	return
}

// ----------------

func rvArrayIndex(rv reflect.Value, i int, _ *typeInfo, _ bool) reflect.Value {
	return rv.Index(i)
}

// func rvArrayIndex(rv reflect.Value, i int, ti *typeInfo) reflect.Value {
// 	return rv.Index(i)
// }

func rvSliceZeroCap(t reflect.Type) (v reflect.Value) {
	return reflect.MakeSlice(t, 0, 0)
}

func rvLenSlice(rv reflect.Value) int {
	return rv.Len()
}

func rvCapSlice(rv reflect.Value) int {
	return rv.Cap()
}

func rvGetArrayBytes(rv reflect.Value, scratch []byte) (bs []byte) {
	l := rv.Len()
	if scratch == nil && rv.CanAddr() {
		return rv.Slice(0, l).Bytes()
	}

	if l <= cap(scratch) {
		bs = scratch[:l]
	} else {
		bs = make([]byte, l)
	}
	reflect.Copy(reflect.ValueOf(bs), rv)
	return
}

func rvGetArray4Slice(rv reflect.Value) (v reflect.Value) {
	v = rvZeroAddrK(reflect.ArrayOf(rvLenSlice(rv), rv.Type().Elem()), reflect.Array)
	reflect.Copy(v, rv)
	return
}

func rvGetSlice4Array(rv reflect.Value, v interface{}) {
	// v is a pointer to a slice to be populated

	// rv.Slice fails if address is not addressable, which can occur during encoding.
	// Consequently, check if non-addressable, and if so, make new slice and copy into it first.
	// MARKER: this *may* cause allocation if non-addressable, unfortunately.

	rve := reflect.ValueOf(v).Elem()
	l := rv.Len()
	if rv.CanAddr() {
		rve.Set(rv.Slice(0, l))
	} else {
		rvs := reflect.MakeSlice(rve.Type(), l, l)
		reflect.Copy(rvs, rv)
		rve.Set(rvs)
	}
	// reflect.ValueOf(v).Elem().Set(rv.Slice(0, rv.Len()))
}

func rvCopySlice(dest, src reflect.Value, _ reflect.Type) {
	reflect.Copy(dest, src)
}

// ------------

func rvGetBool(rv reflect.Value) bool {
	return rv.Bool()
}

func rvGetBytes(rv reflect.Value) []byte {
	return rv.Bytes()
}

func rvGetTime(rv reflect.Value) time.Time {
	return rv2i(rv).(time.Time)
}

func rvGetString(rv reflect.Value) string {
	return rv.String()
}

func rvGetFloat64(rv reflect.Value) float64 {
	return rv.Float()
}

func rvGetFloat32(rv reflect.Value) float32 {
	return float32(rv.Float())
}

func rvGetComplex64(rv reflect.Value) complex64 {
	return complex64(rv.Complex())
}

func rvGetComplex128(rv reflect.Value) complex128 {
	return rv.Complex()
}

func rvGetInt(rv reflect.Value) int {
	return int(rv.Int())
}

func rvGetInt8(rv reflect.Value) int8 {
	return int8(rv.Int())
}

func rvGetInt16(rv reflect.Value) int16 {
	return int16(rv.Int())
}

func rvGetInt32(rv reflect.Value) int32 {
	return int32(rv.Int())
}

func rvGetInt64(rv reflect.Value) int64 {
	return rv.Int()
}

func rvGetUint(rv reflect.Value) uint {
	return uint(rv.Uint())
}

func rvGetUint8(rv reflect.Value) uint8 {
	return uint8(rv.Uint())
}

func rvGetUint16(rv reflect.Value) uint16 {
	return uint16(rv.Uint())
}

func rvGetUint32(rv reflect.Value) uint32 {
	return uint32(rv.Uint())
}

func rvGetUint64(rv reflect.Value) uint64 {
	return rv.Uint()
}

func rvGetUintptr(rv reflect.Value) uintptr {
	return uintptr(rv.Uint())
}

func rvLenMap(rv reflect.Value) int {
	return rv.Len()
}

// ------------ map range and map indexing ----------

func mapSet(m, k, v reflect.Value, _ mapReqParams) {
	m.SetMapIndex(k, v)
}

func mapGet(m, k, v reflect.Value, _ mapReqParams) (vv reflect.Value) {
	return m.MapIndex(k)
}

func mapAddrLoopvarRV(t reflect.Type, k reflect.Kind) (r reflect.Value) {
	return // reflect.New(t).Elem()
}

// ---------- ENCODER optimized ---------------

func (d *decoderBase) bytes2Str(in []byte, att dBytesAttachState) (s string, mutable bool) {
	return d.detach2Str(in, att), false
}

// ---------- structFieldInfo optimized ---------------

func (n *structFieldInfoNode) rvField(v reflect.Value) reflect.Value {
	return v.Field(int(n.index))
}

// ---------- others ---------------

// --------------------------
type atomicRtidFnSlice struct {
	v atomic.Value
}

func (x *atomicRtidFnSlice) load() interface{} {
	return x.v.Load()
}

func (x *atomicRtidFnSlice) store(p interface{}) {
	x.v.Store(p)
}
