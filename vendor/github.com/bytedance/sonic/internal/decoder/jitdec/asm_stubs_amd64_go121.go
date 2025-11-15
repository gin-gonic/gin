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

package jitdec

import (
    `strconv`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/internal/jit`
    `github.com/twitchyliquid64/golang-asm/obj`
    `github.com/twitchyliquid64/golang-asm/obj/x86`
)

// Notice: gcWriteBarrier must use R11 register!!
var _R11 = _IC

var (
    _V_writeBarrier = jit.Imm(int64(uintptr(unsafe.Pointer(&rt.RuntimeWriteBarrier))))

    _F_gcWriteBarrier2 = jit.Func(rt.GcWriteBarrier2)
)

func (self *_Assembler) WritePtrAX(i int, rec obj.Addr, saveDI bool) {
    self.Emit("MOVQ", _V_writeBarrier, _R9)
    self.Emit("CMPL", jit.Ptr(_R9, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    if saveDI {
        self.save(_DI, _R11)
    } else {
        self.save(_R11)
    }
    self.Emit("MOVQ", _F_gcWriteBarrier2, _R11)  
    self.Rjmp("CALL", _R11)  
    self.Emit("MOVQ", _AX, jit.Ptr(_R11, 0))
    self.Emit("MOVQ", rec, _DI)
    self.Emit("MOVQ", _DI, jit.Ptr(_R11, 8))
    if saveDI {
        self.load(_DI, _R11)
    } else {
        self.load(_R11)
    }   
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", _AX, rec)
}

func (self *_Assembler) WriteRecNotAX(i int, ptr obj.Addr, rec obj.Addr, saveDI bool, saveAX bool) {
    if rec.Reg == x86.REG_AX || rec.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _R9)
    self.Emit("CMPL", jit.Ptr(_R9, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    if saveAX {
        self.save(_AX, _R11)
    } else {
        self.save(_R11)
    }
    self.Emit("MOVQ", _F_gcWriteBarrier2, _R11)  
    self.Rjmp("CALL", _R11)  
    self.Emit("MOVQ", ptr, jit.Ptr(_R11, 0))
    self.Emit("MOVQ", rec, _AX)
    self.Emit("MOVQ", _AX, jit.Ptr(_R11, 8))
    if saveAX {
        self.load(_AX, _R11)
    } else {
        self.load(_R11)
    }   
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, rec)
}

func (self *_ValueDecoder) WritePtrAX(i int, rec obj.Addr, saveDI bool) {
    self.Emit("MOVQ", _V_writeBarrier, _R9)
    self.Emit("CMPL", jit.Ptr(_R9, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    if saveDI {
        self.save(_DI, _R11)
    } else {
        self.save(_R11)
    }
    self.Emit("MOVQ", _F_gcWriteBarrier2, _R11)  
    self.Rjmp("CALL", _R11)   
    self.Emit("MOVQ", _AX, jit.Ptr(_R11, 0))
    self.Emit("MOVQ", rec, _DI)
    self.Emit("MOVQ", _DI, jit.Ptr(_R11, 8))
    if saveDI {
        self.load(_DI, _R11)
    } else {
        self.load(_R11)
    }   
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", _AX, rec)
}

func (self *_ValueDecoder) WriteRecNotAX(i int, ptr obj.Addr, rec obj.Addr, saveDI bool) {
    if rec.Reg == x86.REG_AX || rec.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _AX)
    self.Emit("CMPL", jit.Ptr(_AX, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.save(_R11)
    self.Emit("MOVQ", _F_gcWriteBarrier2, _R11)  
    self.Rjmp("CALL", _R11)     
    self.Emit("MOVQ", ptr, jit.Ptr(_R11, 0))
    self.Emit("MOVQ", rec, _AX)
    self.Emit("MOVQ", _AX, jit.Ptr(_R11, 8))
    self.load(_R11)  
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, rec)
}