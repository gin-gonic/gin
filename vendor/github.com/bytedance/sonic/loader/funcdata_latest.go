// go:build go1.18 && !go1.26
// +build go1.18,!go1.26

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
    `os`
    `sort`
    `unsafe`

    `github.com/bytedance/sonic/loader/internal/rt`
)

type funcTab struct {
    entry   uint32
    funcoff uint32
}

type pcHeader struct {
    magic          uint32  // 0xFFFFFFF0
    pad1, pad2     uint8   // 0,0
    minLC          uint8   // min instruction size
    ptrSize        uint8   // size of a ptr in bytes
    nfunc          int     // number of functions in the module
    nfiles         uint    // number of entries in the file tab
    textStart      uintptr // base for function entry PC offsets in this module, equal to moduledata.text
    funcnameOffset uintptr // offset to the funcnametab variable from pcHeader
    cuOffset       uintptr // offset to the cutab variable from pcHeader
    filetabOffset  uintptr // offset to the filetab variable from pcHeader
    pctabOffset    uintptr // offset to the pctab variable from pcHeader
    pclnOffset     uintptr // offset to the pclntab variable from pcHeader
}

type bitVector struct {
    n        int32 // # of bits
    bytedata *uint8
}

type ptabEntry struct {
    name int32
    typ  int32
}

type textSection struct {
    vaddr    uintptr // prelinked section vaddr
    end      uintptr // vaddr + section length
    baseaddr uintptr // relocated section address
}

type modulehash struct {
    modulename   string
    linktimehash string
    runtimehash  *string
}

// findfuncbucket is an array of these structures.
// Each bucket represents 4096 bytes of the text segment.
// Each subbucket represents 256 bytes of the text segment.
// To find a function given a pc, locate the bucket and subbucket for
// that pc. Add together the idx and subbucket value to obtain a
// function index. Then scan the functab array starting at that
// index to find the target function.
// This table uses 20 bytes for every 4096 bytes of code, or ~0.5% overhead.
type findfuncbucket struct {
    idx        uint32
    _SUBBUCKETS [16]byte
}

type compilationUnit struct {
    fileNames []string
}

func makeFtab(funcs []_func, maxpc uint32) (ftab []funcTab, pclntabSize int64, startLocations []uint32) {
    // Allocate space for the pc->func table. This structure consists of a pc offset
    // and an offset to the func structure. After that, we have a single pc
    // value that marks the end of the last function in the binary.
    pclntabSize = int64(len(funcs)*2*int(_PtrSize) + int(_PtrSize))
    startLocations = make([]uint32, len(funcs))
    for i, f := range funcs {
        pclntabSize = rnd(pclntabSize, int64(_PtrSize))
        //writePCToFunc
        startLocations[i] = uint32(pclntabSize)
        pclntabSize += int64(uint8(_FUNC_SIZE)+f.nfuncdata*4+uint8(f.npcdata)*4)
    }

    ftab = make([]funcTab, 0, len(funcs)+1)

    // write a map of pc->func info offsets 
    for i, f := range funcs {
        ftab = append(ftab, funcTab{uint32(f.entryOff), uint32(startLocations[i])})
    }

    // Final entry of table is just end pc offset.
    ftab = append(ftab, funcTab{maxpc, 0})
    return
}

// Pcln table format: [...]funcTab + [...]_Func
func makePclntable(size int64, startLocations []uint32, funcs []_func, maxpc uint32, pcdataOffs [][]uint32, funcdataOffs [][]uint32) (pclntab []byte) {
    // Allocate space for the pc->func table. This structure consists of a pc offset
    // and an offset to the func structure. After that, we have a single pc
    // value that marks the end of the last function in the binary.
    pclntab = make([]byte, size, size)

    // write a map of pc->func info offsets 
    offs := 0
    for i, f := range funcs {
        byteOrder.PutUint32(pclntab[offs:offs+4], uint32(f.entryOff))
        byteOrder.PutUint32(pclntab[offs+4:offs+8], uint32(startLocations[i]))
        offs += 8
    }
    // Final entry of table is just end pc offset.
    byteOrder.PutUint32(pclntab[offs:offs+4], maxpc)

    // write func info table
    for i := range funcs {
        off := startLocations[i]

        // write _func structure to pclntab
        fb := rt.BytesFrom(unsafe.Pointer(&funcs[i]), int(_FUNC_SIZE), int(_FUNC_SIZE))
        copy(pclntab[off:off+uint32(_FUNC_SIZE)], fb)
        off += uint32(_FUNC_SIZE)

        // NOTICE: _func.pcdata always starts from PcUnsafePoint, which is index 3
        for j := 3; j < len(pcdataOffs[i]); j++ {
            byteOrder.PutUint32(pclntab[off:off+4], uint32(pcdataOffs[i][j]))
            off += 4
        }

        // funcdata refs as offsets from gofunc
        for _, funcdata := range funcdataOffs[i] {
            byteOrder.PutUint32(pclntab[off:off+4], uint32(funcdata))
            off += 4
        }

    }

    return
}

// findfunc table used to map pc to belonging func, 
// returns the index in the func table.
//
// All text section are divided into buckets sized _BUCKETSIZE(4K):
//   every bucket is divided into _SUBBUCKETS sized _SUB_BUCKETSIZE(64),
//   and it has a base idx to plus the offset stored in jth subbucket.
// see findfunc() in runtime/symtab.go
func writeFindfunctab(out *[]byte, ftab []funcTab) (start int) {
    start = len(*out)

    max := ftab[len(ftab)-1].entry
    min := ftab[0].entry
    nbuckets := (max - min + _BUCKETSIZE - 1) / _BUCKETSIZE
    n := (max - min + _SUB_BUCKETSIZE - 1) / _SUB_BUCKETSIZE

    tab := make([]findfuncbucket, 0, nbuckets)
    var s, e = 0, 0
    for i := 0; i<int(nbuckets); i++ {
        // store the start s-th func of the bucket
        var fb = findfuncbucket{idx: uint32(s)}

        // find the last e-th func of the bucket
        var pc = min + uint32((i+1)*_BUCKETSIZE)
        for ; e < len(ftab)-1 && ftab[e+1].entry <= pc; e++ {}

        for j := 0; j<_SUBBUCKETS && (i*_SUBBUCKETS+j)<int(n); j++ {
            // store the start func of the subbucket
            fb._SUBBUCKETS[j] = byte(uint32(s) - fb.idx)

            // find the s-th end func of the subbucket
            pc = min + uint32(i*_BUCKETSIZE) + uint32((j+1)*_SUB_BUCKETSIZE)
            for ; s < len(ftab)-1 && ftab[s+1].entry <= pc; s++ {}            
        }

        s = e
        tab = append(tab, fb)
    }

    // write findfuncbucket
    if len(tab) > 0 {
        size := int(unsafe.Sizeof(findfuncbucket{}))*len(tab)
        *out = append(*out, rt.BytesFrom(unsafe.Pointer(&tab[0]), size, size)...)
    }
    return 
}

func makeModuledata(name string, filenames []string, funcsp *[]Func, text []byte) (mod *moduledata) {
    mod = new(moduledata)
    mod.modulename = name

    // sort funcs by entry
    funcs := *funcsp
    sort.Slice(funcs, func(i, j int) bool {
        return funcs[i].EntryOff < funcs[j].EntryOff
    })
    *funcsp = funcs

    // make filename table
    cu := make([]string, 0, len(filenames))
    cu = append(cu, filenames...)
    cutab, filetab, cuOffs := makeFilenametab([]compilationUnit{{cu}})
    mod.cutab = cutab
    mod.filetab = filetab

    // make funcname table
    funcnametab, nameOffs := makeFuncnameTab(funcs)
    mod.funcnametab = funcnametab

    // mmap() text and funcdata segments
    p := os.Getpagesize()
    size := int(rnd(int64(len(text)), int64(p)))
    addr := mmap(size)
    // copy the machine code
    s := rt.BytesFrom(unsafe.Pointer(addr), len(text), size)
    copy(s, text)
    // make it executable
    mprotect(addr, size)

    // assign addresses
    mod.text = addr
    mod.etext = addr + uintptr(size)
    mod.minpc = addr
    mod.maxpc = addr + uintptr(len(text))

    // make pcdata table
    // NOTICE: _func only use offset to index pcdata, thus no need mmap() pcdata 
    cuOff := cuOffs[0]
    pctab, pcdataOffs, _funcs := makePctab(funcs, cuOff, nameOffs)
    mod.pctab = pctab

    // write func data
    // NOTICE: _func use mod.gofunc+offset to directly point funcdata, thus need cache funcdata
    // TODO: estimate accurate capacity
    cache := make([]byte, 0, len(funcs)*int(_PtrSize)) 
    fstart, funcdataOffs := writeFuncdata(&cache, funcs)

    // make pc->func (binary search) func table
    ftab, pclntSize, startLocations := makeFtab(_funcs, uint32(len(text)))
    mod.ftab = ftab

    // write pc->func (modmap) findfunc table
    ffstart := writeFindfunctab(&cache, ftab)

    // cache funcdata and findfuncbucket
    moduleCache.Lock()
    moduleCache.m[mod] = cache
    moduleCache.Unlock()
    mod.gofunc = uintptr(unsafe.Pointer(&cache[fstart]))
    mod.findfunctab = uintptr(unsafe.Pointer(&cache[ffstart]))

    // make pclnt table
    pclntab := makePclntable(pclntSize, startLocations, _funcs, uint32(len(text)), pcdataOffs, funcdataOffs)
    mod.pclntable = pclntab

    // make pc header
    mod.pcHeader = &pcHeader {
        magic   : _Magic,
        minLC   : _MinLC,
        ptrSize : _PtrSize,
        nfunc   : len(funcs),
        nfiles: uint(len(cu)),
        textStart: mod.text,
        funcnameOffset: getOffsetOf(moduledata{}, "funcnametab"),
        cuOffset: getOffsetOf(moduledata{}, "cutab"),
        filetabOffset: getOffsetOf(moduledata{}, "filetab"),
        pctabOffset: getOffsetOf(moduledata{}, "pctab"),
        pclnOffset: getOffsetOf(moduledata{}, "pclntable"),
    }

    // special case: gcdata and gcbss must by non-empty
    mod.gcdata = uintptr(unsafe.Pointer(&emptyByte))
    mod.gcbss = uintptr(unsafe.Pointer(&emptyByte))

    return
}

// makePctab generates pcdelta->valuedelta tables for functions,
// and returns the table and the entry offset of every kind pcdata in the table.
func makePctab(funcs []Func, cuOffset uint32, nameOffset []int32) (pctab []byte, pcdataOffs [][]uint32, _funcs []_func) {
    _funcs = make([]_func, len(funcs))

    // Pctab offsets of 0 are considered invalid in the runtime. We respect
    // that by just padding a single byte at the beginning of runtime.pctab,
    // that way no real offsets can be zero.
    pctab = make([]byte, 1, 12*len(funcs)+1)
    pcdataOffs = make([][]uint32, len(funcs))

    for i, f := range funcs {
        _f := &_funcs[i]

        var writer = func(pc *Pcdata) {
            var ab []byte
            var err error
            if pc != nil {
                ab, err = pc.MarshalBinary()
                if err != nil {
                    panic(err)
                }
                pcdataOffs[i] = append(pcdataOffs[i], uint32(len(pctab)))
            } else {
                ab = []byte{0}
                pcdataOffs[i] = append(pcdataOffs[i], _PCDATA_INVALID_OFFSET)
            }
            pctab = append(pctab, ab...)
        }

        if f.Pcsp != nil {
            _f.pcsp = uint32(len(pctab))
        }
        writer(f.Pcsp)
        if f.Pcfile != nil {
            _f.pcfile = uint32(len(pctab))
        }
        writer(f.Pcfile)
        if f.Pcline != nil {
            _f.pcln = uint32(len(pctab))
        }
        writer(f.Pcline)
        writer(f.PcUnsafePoint)
        writer(f.PcStackMapIndex)
        writer(f.PcInlTreeIndex)
        writer(f.PcArgLiveIndex)
        
        _f.entryOff = f.EntryOff
        _f.nameOff = nameOffset[i]
        _f.args = f.ArgsSize
        _f.deferreturn = f.DeferReturn
        // NOTICE: _func.pcdata is always as [PCDATA_UnsafePoint(0) : PCDATA_ArgLiveIndex(3)]
        _f.npcdata = uint32(_N_PCDATA)
        _f.cuOffset = cuOffset
        _f.funcID = f.ID
        _f.flag = f.Flag
        _f.nfuncdata = uint8(_N_FUNCDATA)
    }

    return
}

func registerFunction(name string, pc uintptr, textSize uintptr, fp int, args int, size uintptr, argptrs uintptr, localptrs uintptr) {} 
