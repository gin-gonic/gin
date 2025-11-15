/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rt

import (
	"reflect"
	"unsafe"
)

//go:noescape
//go:linkname Memmove runtime.memmove
func Memmove(to unsafe.Pointer, from unsafe.Pointer, n uintptr)
//go:noescape
//go:linkname MemEqual runtime.memequal
//goland:noinspection GoUnusedParameter
func MemEqual(a unsafe.Pointer, b unsafe.Pointer, size uintptr) bool

//go:linkname Mapiternext runtime.mapiternext
func Mapiternext(it *GoMapIterator)

//go:linkname Mapiterinit runtime.mapiterinit
func Mapiterinit(t *GoMapType, m unsafe.Pointer, it *GoMapIterator)

//go:linkname Maplen reflect.maplen
func Maplen(h unsafe.Pointer) int 

//go:nosplit
//go:linkname MemclrHasPointers runtime.memclrHasPointers
//goland:noinspection GoUnusedParameter
func MemclrHasPointers(ptr unsafe.Pointer, n uintptr)

//go:linkname MemclrNoHeapPointers runtime.memclrNoHeapPointers
//goland:noinspection GoUnusedParameter
func MemclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

//go:linkname newarray runtime.newarray
func newarray(typ *GoType, n int) unsafe.Pointer 

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func ClearMemory(et *GoType, ptr unsafe.Pointer, size uintptr) {
	if et.PtrData == 0 {
		MemclrNoHeapPointers(ptr, size)
	} else {
		MemclrHasPointers(ptr, size)
	}
}

// runtime.maxElementSize
const _max_map_element_size uintptr = 128

func IsMapfast(vt reflect.Type) bool {
	return vt.Elem().Size() <= _max_map_element_size
}

//go:linkname Mallocgc runtime.mallocgc
//goland:noinspection GoUnusedParameter
func Mallocgc(size uintptr, typ *GoType, needzero bool) unsafe.Pointer

//go:linkname Makemap reflect.makemap
func Makemap(*GoType, int) unsafe.Pointer

//go:linkname MakemapSmall runtime.makemap_small
func MakemapSmall() unsafe.Pointer

//go:linkname Mapassign runtime.mapassign
//goland:noinspection GoUnusedParameter
func Mapassign(t *GoMapType, h unsafe.Pointer, k unsafe.Pointer) unsafe.Pointer

//go:linkname Mapassign_fast32 runtime.mapassign_fast32
//goland:noinspection GoUnusedParameter
func Mapassign_fast32(t *GoMapType, h unsafe.Pointer, k uint32) unsafe.Pointer

//go:linkname Mapassign_fast64 runtime.mapassign_fast64
//goland:noinspection GoUnusedParameter
func Mapassign_fast64(t *GoMapType, h unsafe.Pointer, k uint64) unsafe.Pointer

//go:linkname Mapassign_faststr runtime.mapassign_faststr
//goland:noinspection GoUnusedParameter
func Mapassign_faststr(t *GoMapType, h unsafe.Pointer, s string) unsafe.Pointer

type MapStrAssign func (t *GoMapType, h unsafe.Pointer, s string) unsafe.Pointer

func GetMapStrAssign(vt reflect.Type) MapStrAssign {
	if IsMapfast(vt) {
		return Mapassign_faststr
	} else {
		return func (t *GoMapType, h unsafe.Pointer, s string) unsafe.Pointer {
			return Mapassign(t, h, unsafe.Pointer(&s))
		}
	}
}

type Map32Assign func(t *GoMapType, h unsafe.Pointer, k uint32) unsafe.Pointer

func GetMap32Assign(vt reflect.Type) Map32Assign {
	if IsMapfast(vt) {
		return Mapassign_fast32
	} else {
		return func (t *GoMapType, h unsafe.Pointer, s uint32) unsafe.Pointer {
			return Mapassign(t, h, unsafe.Pointer(&s))
		}
	}
}

type Map64Assign func(t *GoMapType, h unsafe.Pointer, k uint64) unsafe.Pointer

func GetMap64Assign(vt reflect.Type) Map64Assign {
	if IsMapfast(vt) {
		return Mapassign_fast64
	} else {
		return func (t *GoMapType, h unsafe.Pointer, s uint64) unsafe.Pointer {
			return Mapassign(t, h, unsafe.Pointer(&s))
		}
	}
}


var emptyBytes = make([]byte, 0, 0)
var EmptySlice = *(*GoSlice)(unsafe.Pointer(&emptyBytes))

//go:linkname MakeSliceStd runtime.makeslice
//goland:noinspection GoUnusedParameter
func MakeSliceStd(et *GoType, len int, cap int) unsafe.Pointer

func MakeSlice(oldPtr unsafe.Pointer, et *GoType, newLen int) *GoSlice {
	if newLen == 0 {
		return &EmptySlice
	}

	if *(*unsafe.Pointer)(oldPtr) == nil {
		return &GoSlice{
			Ptr: MakeSliceStd(et, newLen, newLen),
			Len: newLen,
			Cap: newLen,
		}
	}

	old := (*GoSlice)(oldPtr)
	if old.Cap >= newLen {
		old.Len = newLen
		return old
	}

	new := GrowSlice(et, *old, newLen)

	// we should clear the memory from [oldLen:newLen]
	if et.PtrData == 0 {
		oldlenmem := uintptr(old.Len) * et.Size
		newlenmem := uintptr(newLen) * et.Size
		MemclrNoHeapPointers(add(new.Ptr, oldlenmem), newlenmem-oldlenmem)
	}

	new.Len = newLen
	return &new
}

//go:nosplit
//go:linkname Throw runtime.throw
//goland:noinspection GoUnusedParameter
func Throw(s string)

//go:linkname ConvT64 runtime.convT64
//goland:noinspection GoUnusedParameter
func ConvT64(v uint64) unsafe.Pointer

//go:linkname ConvTslice runtime.convTslice
//goland:noinspection GoUnusedParameter
func ConvTslice(v []byte) unsafe.Pointer

//go:linkname ConvTstring runtime.convTstring
//goland:noinspection GoUnusedParameter
func ConvTstring(v string) unsafe.Pointer

//go:linkname Mapassign_fast64ptr runtime.mapassign_fast64ptr
//goland:noinspection GoUnusedParameter
func Mapassign_fast64ptr(t *GoMapType, h unsafe.Pointer, k unsafe.Pointer) unsafe.Pointer

//go:noescape
//go:linkname Strhash runtime.strhash
func Strhash(_ unsafe.Pointer, _ uintptr) uintptr
