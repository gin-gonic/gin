//go:build go1.21 && !go1.26
// +build go1.21,!go1.26

// Copyright 2023 CloudWeGo Authors
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

package x86

import (
	"strconv"
	"unsafe"

	"github.com/bytedance/sonic/internal/jit"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/x86"
)

var (
    _V_writeBarrier = jit.Imm(int64(uintptr(unsafe.Pointer(&rt.RuntimeWriteBarrier))))

    _F_gcWriteBarrier2 = jit.Func(rt.GcWriteBarrier2)
)

func (self *Assembler) WritePtr(i int, ptr obj.Addr, old obj.Addr) {
    if old.Reg == x86.REG_AX || old.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _BX)
    self.Emit("CMPL", jit.Ptr(_BX, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.xsave(_SP_q)
    self.Emit("MOVQ", _F_gcWriteBarrier2, _BX)  // MOVQ ${fn}, AX
    self.Rjmp("CALL", _BX)  
    self.Emit("MOVQ", ptr, jit.Ptr(_SP_q, 0))
    self.Emit("MOVQ", old, _AX)
    self.Emit("MOVQ", _AX, jit.Ptr(_SP_q, 8))
    self.xload(_SP_q)  
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, old)
}
