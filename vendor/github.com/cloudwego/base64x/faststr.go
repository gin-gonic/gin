/*
 * Copyright 2024 CloudWeGo Authors
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

package base64x

import (
    `reflect`
    `unsafe`
)

func mem2str(v []byte) (s string) {
    (*reflect.StringHeader)(unsafe.Pointer(&s)).Len  = (*reflect.SliceHeader)(unsafe.Pointer(&v)).Len
    (*reflect.StringHeader)(unsafe.Pointer(&s)).Data = (*reflect.SliceHeader)(unsafe.Pointer(&v)).Data
    return
}

func str2mem(s string) (v []byte) {
    (*reflect.SliceHeader)(unsafe.Pointer(&v)).Cap  = (*reflect.StringHeader)(unsafe.Pointer(&s)).Len
    (*reflect.SliceHeader)(unsafe.Pointer(&v)).Len  = (*reflect.StringHeader)(unsafe.Pointer(&s)).Len
    (*reflect.SliceHeader)(unsafe.Pointer(&v)).Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
    return
}

func mem2addr(v []byte) unsafe.Pointer {
    return *(*unsafe.Pointer)(unsafe.Pointer(&v))
}

// NoEscape hides a pointer from escape analysis. NoEscape is
// the identity function but escape analysis doesn't think the
// output depends on the input. NoEscape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//go:nosplit
//goland:noinspection GoVetUnsafePointer
func noEscape(p unsafe.Pointer) unsafe.Pointer {
    x := uintptr(p)
    return unsafe.Pointer(x ^ 0)
}
