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
    `unsafe`
    `reflect`
)

//go:nosplit
func Mem2Str(v []byte) (s string) {
    (*GoString)(unsafe.Pointer(&s)).Len = (*GoSlice)(unsafe.Pointer(&v)).Len
    (*GoString)(unsafe.Pointer(&s)).Ptr = (*GoSlice)(unsafe.Pointer(&v)).Ptr
    return
}

//go:nosplit
func Str2Mem(s string) (v []byte) {
    (*GoSlice)(unsafe.Pointer(&v)).Cap = (*GoString)(unsafe.Pointer(&s)).Len
    (*GoSlice)(unsafe.Pointer(&v)).Len = (*GoString)(unsafe.Pointer(&s)).Len
    (*GoSlice)(unsafe.Pointer(&v)).Ptr = (*GoString)(unsafe.Pointer(&s)).Ptr
    return
}

func BytesFrom(p unsafe.Pointer, n int, c int) (r []byte) {
    (*GoSlice)(unsafe.Pointer(&r)).Ptr = p
    (*GoSlice)(unsafe.Pointer(&r)).Len = n
    (*GoSlice)(unsafe.Pointer(&r)).Cap = c
    return
}

func FuncAddr(f interface{}) unsafe.Pointer {
    if vv := UnpackEface(f); vv.Type.Kind() != reflect.Func {
        panic("f is not a function")
    } else {
        return *(*unsafe.Pointer)(vv.Value)
    }
}

//go:nocheckptr
func IndexChar(src string, index int) unsafe.Pointer {
	return unsafe.Pointer(uintptr((*GoString)(unsafe.Pointer(&src)).Ptr) + uintptr(index))
}

//go:nocheckptr
func IndexByte(ptr []byte, index int) unsafe.Pointer {
	return unsafe.Pointer(uintptr((*GoSlice)(unsafe.Pointer(&ptr)).Ptr) + uintptr(index))
}
