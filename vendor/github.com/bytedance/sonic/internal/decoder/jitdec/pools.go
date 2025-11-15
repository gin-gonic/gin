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

package jitdec

import (
    `sync`
    `unsafe`

    `github.com/bytedance/sonic/internal/caching`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
)

const (
    _MinSlice = 2
    _MaxStack = 4096 // 4k slots
    _MaxStackBytes = _MaxStack * _PtrBytes
    _MaxDigitNums = types.MaxDigitNums  // used in atof fallback algorithm
)

const (
    _PtrBytes   = _PTR_SIZE / 8
    _FsmOffset  = (_MaxStack + 1) * _PtrBytes
    _DbufOffset = _FsmOffset + int64(unsafe.Sizeof(types.StateMachine{})) + types.MAX_RECURSE * _PtrBytes
    _EpOffset   = _DbufOffset + _MaxDigitNums
    _StackSize  = unsafe.Sizeof(_Stack{})
)

var (
    stackPool     = sync.Pool{}
    valueCache    = []unsafe.Pointer(nil)
    fieldCache    = []*caching.FieldMap(nil)
    fieldCacheMux = sync.Mutex{}
    programCache  = caching.CreateProgramCache()
)

type _Stack struct {
    sp uintptr
    sb [_MaxStack]unsafe.Pointer
    mm types.StateMachine
    vp [types.MAX_RECURSE]unsafe.Pointer
    dp [_MaxDigitNums]byte
    ep unsafe.Pointer
}

type _Decoder func(
    s  string,
    i  int,
    vp unsafe.Pointer,
    sb *_Stack,
    fv uint64,
    sv string, // DO NOT pass value to this argument, since it is only used for local _VAR_sv
    vk unsafe.Pointer, // DO NOT pass value to this argument, since it is only used for local _VAR_vk
) (int, error)

var _KeepAlive struct {
    s string
    i int
    vp unsafe.Pointer
    sb *_Stack
    fv uint64
    sv string
    vk unsafe.Pointer

    ret int
    err error

    frame_decoder [_FP_offs]byte
    frame_generic [_VD_offs]byte
}

var (
    argPtrs   = []bool{true, false, false, true, true, false, true, false, true}
    localPtrs = []bool{}
)

var (
    argPtrs_generic   = []bool{true}
    localPtrs_generic = []bool{}
)

func newStack() *_Stack {
    if ret := stackPool.Get(); ret == nil {
        return new(_Stack)
    } else {
        return ret.(*_Stack)
    }
}

func resetStack(p *_Stack) {
    rt.MemclrNoHeapPointers(unsafe.Pointer(p), _StackSize)
}

func freeStack(p *_Stack) {
    p.sp = 0
    stackPool.Put(p)
}

func freezeValue(v unsafe.Pointer) uintptr {
    valueCache = append(valueCache, v)
    return uintptr(v)
}

func freezeFields(v *caching.FieldMap) int64 {
    fieldCacheMux.Lock()
    fieldCache = append(fieldCache, v)
    fieldCacheMux.Unlock()
    return referenceFields(v)
}

func referenceFields(v *caching.FieldMap) int64 {
    return int64(uintptr(unsafe.Pointer(v)))
}

func makeDecoder(vt *rt.GoType, _ ...interface{}) (interface{}, error) {
    if pp, err := newCompiler().compile(vt.Pack()); err != nil {
        return nil, err
    } else {
        return newAssembler(pp).Load(), nil
    }
}

func findOrCompile(vt *rt.GoType) (_Decoder, error) {
    if val := programCache.Get(vt); val != nil {
        return val.(_Decoder), nil
    } else if ret, err := programCache.Compute(vt, makeDecoder); err == nil {
        return ret.(_Decoder), nil
    } else {
        return nil, err
    }
}
