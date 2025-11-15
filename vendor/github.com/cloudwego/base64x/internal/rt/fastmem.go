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
	"unsafe"
)


// NoEscape hides a pointer from escape analysis. NoEscape is
// the identity function but escape analysis doesn't think the
// output depends on the input. NoEscape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//go:nosplit
//goland:noinspection GoVetUnsafePointer
func NoEscape(p unsafe.Pointer) unsafe.Pointer {
    x := uintptr(p)
    return unsafe.Pointer(x ^ 0)
}

//go:nosplit
func MoreStack(size uintptr)

//go:nosplit
func Add(ptr unsafe.Pointer, off uintptr) unsafe.Pointer {
    return unsafe.Pointer(uintptr(ptr) + off)
}
