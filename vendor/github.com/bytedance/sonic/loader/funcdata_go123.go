//go:build go1.23 && !go1.26
// +build go1.23,!go1.26

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
    `unsafe`
    `github.com/bytedance/sonic/loader/internal/rt`
)

const (
    _Magic uint32 = 0xFFFFFFF1
)

type moduledata struct {
    pcHeader     *pcHeader
    funcnametab  []byte
    cutab        []uint32
    filetab      []byte
    pctab        []byte
    pclntable    []byte
    ftab         []funcTab
    findfunctab  uintptr
    minpc, maxpc uintptr // first func address, last func address + last func size

    text, etext           uintptr // start/end of text, (etext-text) must be greater than MIN_FUNC
    noptrdata, enoptrdata uintptr
    data, edata           uintptr
    bss, ebss             uintptr
    noptrbss, enoptrbss   uintptr
    covctrs, ecovctrs     uintptr
    end, gcdata, gcbss    uintptr
    types, etypes         uintptr
    rodata                uintptr
    gofunc                uintptr // go.func.* is actual funcinfo object in image

    textsectmap []textSection // see runtime/symtab.go: textAddr()
    typelinks   []int32 // offsets from types
    itablinks   []*rt.GoItab

    ptab []ptabEntry

    pluginpath string
    pkghashes  []modulehash

    // This slice records the initializing tasks that need to be
	// done to start up the program. It is built by the linker.
	inittasks []unsafe.Pointer

    modulename   string
    modulehashes []modulehash

    hasmain uint8 // 1 if module contains the main function, 0 otherwise
    bad bool // module failed to load and should be ignored

    gcdatamask, gcbssmask bitVector

    typemap map[int32]*rt.GoType // offset to *_rtype in previous module

    next *moduledata
}

type _func struct {
    entryOff uint32 // start pc, as offset from moduledata.text/pcHeader.textStart
    nameOff  int32  // function name, as index into moduledata.funcnametab.

    args        int32  // in/out args size
    deferreturn uint32 // offset of start of a deferreturn call instruction from entry, if any.

    pcsp      uint32 
    pcfile    uint32
    pcln      uint32
    npcdata   uint32
    cuOffset  uint32 // runtime.cutab offset of this function's CU
    startLine int32  // line number of start of function (func keyword/TEXT directive)
    funcID    uint8 // set for certain special runtime functions
    flag      uint8
    _         [1]byte // pad
    nfuncdata uint8   // 
    
    // The end of the struct is followed immediately by two variable-length
    // arrays that reference the pcdata and funcdata locations for this
    // function.

    // pcdata contains the offset into moduledata.pctab for the start of
    // that index's table. e.g.,
    // &moduledata.pctab[_func.pcdata[_PCDATA_UnsafePoint]] is the start of
    // the unsafe point table.
    //
    // An offset of 0 indicates that there is no table.
    //
    // pcdata [npcdata]uint32

    // funcdata contains the offset past moduledata.gofunc which contains a
    // pointer to that index's funcdata. e.g.,
    // *(moduledata.gofunc +  _func.funcdata[_FUNCDATA_ArgsPointerMaps]) is
    // the argument pointer map.
    //
    // An offset of ^uint32(0) indicates that there is no entry.
    //
    // funcdata [nfuncdata]uint32
}
