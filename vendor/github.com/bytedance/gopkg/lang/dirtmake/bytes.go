// Copyright 2024 ByteDance Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dirtmake

import (
	"unsafe"
)

type slice struct {
	data unsafe.Pointer
	len  int
	cap  int
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

// Bytes allocates a byte slice but does not clean up the memory it references.
// Throw a fatal error instead of panic if cap is greater than runtime.maxAlloc.
// NOTE: MUST set any byte element before it's read.
func Bytes(len, cap int) (b []byte) {
	if len < 0 || len > cap {
		panic("dirtmake.Bytes: len out of range")
	}
	p := mallocgc(uintptr(cap), nil, false)
	sh := (*slice)(unsafe.Pointer(&b))
	sh.data = p
	sh.len = len
	sh.cap = cap
	return
}
