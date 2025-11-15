//go:build !windows
// +build !windows

/**
 * Copyright 2023 ByteDance Inc.
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

package loader

import (
    `syscall`
)

const (
    _AP = syscall.MAP_ANON  | syscall.MAP_PRIVATE
    _RX = syscall.PROT_READ | syscall.PROT_EXEC
    _RW = syscall.PROT_READ | syscall.PROT_WRITE
)


func mmap(nb int) uintptr {
    if m, _, e := syscall.RawSyscall6(syscall.SYS_MMAP, 0, uintptr(nb), _RW, _AP, 0, 0); e != 0 {
        panic(e)
    } else {
        return m
    }
}

func mprotect(p uintptr, nb int) {
    if _, _, err := syscall.RawSyscall(syscall.SYS_MPROTECT, p, uintptr(nb), _RX); err != 0 {
        panic(err)
    }
}
