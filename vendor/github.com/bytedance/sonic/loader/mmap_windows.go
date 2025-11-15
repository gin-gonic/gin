//go:build windows
// +build windows

// build

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

package loader

import (
    `syscall`
    `unsafe`
)

const (
    MEM_COMMIT  = 0x00001000
    MEM_RESERVE = 0x00002000
)

var (
    libKernel32                = syscall.NewLazyDLL("KERNEL32.DLL")
    libKernel32_VirtualAlloc   = libKernel32.NewProc("VirtualAlloc")
    libKernel32_VirtualProtect = libKernel32.NewProc("VirtualProtect")
)

func mmap(nb int) uintptr {
    addr, err := winapi_VirtualAlloc(0, nb, MEM_COMMIT|MEM_RESERVE, syscall.PAGE_READWRITE)
    if err != nil {
        panic(err)
    }
    return addr
}

func mprotect(p uintptr, nb int) (oldProtect int) {
    err := winapi_VirtualProtect(p, nb, syscall.PAGE_EXECUTE_READ, &oldProtect)
    if err != nil {
        panic(err)
    }
    return
}

// winapi_VirtualAlloc allocate memory
// Doc: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
func winapi_VirtualAlloc(lpAddr uintptr, dwSize int, flAllocationType int, flProtect int) (uintptr, error) {
    r1, _, err := libKernel32_VirtualAlloc.Call(
        lpAddr,
        uintptr(dwSize),
        uintptr(flAllocationType),
        uintptr(flProtect),
    )
    if r1 == 0 {
        return 0, err
    }
    return r1, nil
}

// winapi_VirtualProtect change memory protection
// Doc: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualprotect
func winapi_VirtualProtect(lpAddr uintptr, dwSize int, flNewProtect int, lpflOldProtect *int) error {
    r1, _, err := libKernel32_VirtualProtect.Call(
        lpAddr,
        uintptr(dwSize),
        uintptr(flNewProtect),
        uintptr(unsafe.Pointer(lpflOldProtect)),
    )
    if r1 == 0 {
        return err
    }
    return nil
}
