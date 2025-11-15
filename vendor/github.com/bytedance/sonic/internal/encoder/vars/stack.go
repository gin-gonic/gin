/**
 * Copyright 2024 ByteDance Inc.
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

package vars

import (
	"bytes"
	"sync"
	"unsafe"

	"github.com/bytedance/sonic/internal/caching"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/option"
)

type State struct {
	x int
	f uint64
	p unsafe.Pointer
	q unsafe.Pointer
}

type Stack struct {
	sp uintptr
	sb [MaxStack]State
}

var (
	bytesPool    = sync.Pool{}
	stackPool    = sync.Pool{
		New: func() interface{} {
			return &Stack{}
		},
	}
	bufferPool   = sync.Pool{}
	programCache = caching.CreateProgramCache()
)

func ResetProgramCache() {
	programCache.Reset()
}

func NewBytes() *[]byte {
	if ret := bytesPool.Get(); ret != nil {
		return ret.(*[]byte)
	} else {
		ret := make([]byte, 0, option.DefaultEncoderBufferSize)
		return &ret
	}
}

func NewStack() *Stack {
	ret :=  stackPool.Get().(*Stack)
	ret.sp = 0
	return ret
}

func ResetStack(p *Stack) {
	rt.MemclrNoHeapPointers(unsafe.Pointer(p), StackSize)
}

func (s *Stack) Top() *State {
	return (*State)(rt.Add(unsafe.Pointer(&s.sb[0]), s.sp))
}

func (s *Stack) Cur() *State {
	return (*State)(rt.Add(unsafe.Pointer(&s.sb[0]), s.sp - uintptr(StateSize)))
}

const _MaxStackSP = uintptr(MaxStack * StateSize)

func (s *Stack) Push(v State) bool {
	if uintptr(s.sp) >= _MaxStackSP {
		return false
	}
	st := s.Top()
	*st = v
	s.sp += uintptr(StateSize)
	return true
}

func (s *Stack) Pop() State {
	s.sp -= uintptr(StateSize)
	st := s.Top()
	ret := *st
	*st = State{}
	return ret
}

func (s *Stack) Load() (int, uint64, unsafe.Pointer, unsafe.Pointer) {
	st := s.Cur()
	return st.x, st.f, st.p, st.q
}

func (s *Stack) Save(x int, f uint64, p unsafe.Pointer, q unsafe.Pointer) bool {
	return s.Push(State{x: x, f:f, p: p, q: q})
}

func (s *Stack) Drop() (int, uint64, unsafe.Pointer, unsafe.Pointer) {
	st := s.Pop()
	return st.x, st.f, st.p, st.q
}

func NewBuffer() *bytes.Buffer {
	if ret := bufferPool.Get(); ret != nil {
		return ret.(*bytes.Buffer)
	} else {
		return bytes.NewBuffer(make([]byte, 0, option.DefaultEncoderBufferSize))
	}
}

func FreeBytes(p *[]byte) {
	if rt.CanSizeResue(cap(*p)) {
		(*p) = (*p)[:0]
		bytesPool.Put(p)
	}
}

func FreeStack(p *Stack) {
	p.sp = 0
	stackPool.Put(p)
}

func FreeBuffer(p *bytes.Buffer) {
	if rt.CanSizeResue(cap(p.Bytes())) {
		p.Reset()
		bufferPool.Put(p)
	}
}

var (
	ArgPtrs   = []bool{true, true, true, false}
	LocalPtrs = []bool{}
	
    ArgPtrs_generic   = []bool{true}
    LocalPtrs_generic = []bool{}
)