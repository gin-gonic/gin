// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !safe && !codec.safe && !appengine && go1.9 && !gc

package codec

import (
	"reflect"
	_ "runtime" // needed for go linkname(s)
	"unsafe"
)

var unsafeZeroArr [1024]byte

type mapReqParams struct {
	ref bool
}

func getMapReqParams(ti *typeInfo) (r mapReqParams) {
	r.ref = refBitset.isset(ti.elemkind)
	return
}

// runtime.growslice does not work with gccgo, failing with "growslice: cap out of range" error.
// consequently, we just call newarray followed by typedslicecopy directly.

func unsafeGrowslice(typ unsafe.Pointer, old unsafeSlice, cap, incr int) (v unsafeSlice) {
	size := rtsize2(typ)
	if size == 0 {
		return unsafeSlice{unsafe.Pointer(&unsafeZeroArr[0]), old.Len, cap + incr}
	}
	newcap := int(growCap(uint(cap), uint(size), uint(incr)))
	v = unsafeSlice{Data: newarray(typ, newcap), Len: old.Len, Cap: newcap}
	if old.Len > 0 {
		typedslicecopy(typ, v, old)
	}
	// memmove(v.Data, old.Data, size*uintptr(old.Len))
	return
}

// runtime.{mapassign_fastXXX, mapaccess2_fastXXX} are not supported in gollvm,
// failing with "error: undefined reference" error.
// so we just use runtime.{mapassign, mapaccess2} directly

func mapSet(m, k, v reflect.Value, p mapReqParams) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))
	var vtyp = urv.typ
	var vptr = unsafeMapKVPtr(urv)

	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mptr := rvRefPtr(urv)

	vvptr := mapassign(urv.typ, mptr, kptr)
	typedmemmove(vtyp, vvptr, vptr)
}

func mapGet(m, k, v reflect.Value, p mapReqParams) (_ reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mptr := rvRefPtr(urv)

	vvptr, ok := mapaccess2(urv.typ, mptr, kptr)

	if !ok {
		return
	}

	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))

	if helperUnsafeDirectAssignMapEntry || p.ref {
		urv.ptr = vvptr
	} else {
		typedmemmove(urv.typ, urv.ptr, vvptr)
	}

	return v
}
