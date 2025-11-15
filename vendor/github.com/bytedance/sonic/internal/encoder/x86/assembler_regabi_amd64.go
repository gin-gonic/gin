//go:build go1.17 && !go1.26
// +build go1.17,!go1.26

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

package x86

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/bytedance/sonic/internal/cpu"
	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/ir"
	"github.com/bytedance/sonic/internal/encoder/prim"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/jit"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/x86"

	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/rt"
)

/** Register Allocations
 *
 *  State Registers:
 *
 *      %rbx : stack base
 *      %rdi : result pointer
 *      %rsi : result length
 *      %rdx : result capacity
 *      %r12 : sp->p
 *      %r13 : sp->q
 *      %r14 : sp->x
 *      %r15 : sp->f
 *
 *  Error Registers:
 *
 *      %r10 : error type register
 *      %r11 : error pointer register
 */

/** Function Prototype & Stack Map
 *
 *  func (buf *[]byte, p unsafe.Pointer, sb *_Stack, fv uint64) (err error)
 *
 *  buf    :   (FP)
 *  p      :  8(FP)
 *  sb     : 16(FP)
 *  fv     : 24(FP)
 *  err.vt : 32(FP)
 *  err.vp : 40(FP)
 */

const (
	_S_cond = iota
	_S_init
)

const (
	_FP_args   = 32 // 32 bytes for spill registers of arguments
	_FP_fargs  = 40 // 40 bytes for passing arguments to other Go functions
	_FP_saves  = 64 // 64 bytes for saving the registers before CALL instructions
	_FP_locals = 24 // 24 bytes for local variables
)

const (
	_FP_loffs = _FP_fargs + _FP_saves
	FP_offs  = _FP_loffs + _FP_locals
	// _FP_offs  = _FP_loffs + _FP_locals + _FP_debug
	_FP_size = FP_offs + 8 // 8 bytes for the parent frame pointer
	_FP_base = _FP_size + 8 // 8 bytes for the return address
)

const (
	_FM_exp32 = 0x7f800000
	_FM_exp64 = 0x7ff0000000000000
)

const (
	_IM_null   = 0x6c6c756e // 'null'
	_IM_true   = 0x65757274 // 'true'
	_IM_fals   = 0x736c6166 // 'fals' ('false' without the 'e')
	_IM_open   = 0x00225c22 // '"\"∅'
	_IM_array  = 0x5d5b     // '[]'
	_IM_object = 0x7d7b     // '{}'
	_IM_mulv   = -0x5555555555555555
)

const (
	_LB_more_space        = "_more_space"
	_LB_more_space_return = "_more_space_return_"
)

const (
	_LB_error                 = "_error"
	_LB_error_too_deep        = "_error_too_deep"
	_LB_error_invalid_number  = "_error_invalid_number"
	_LB_error_nan_or_infinite = "_error_nan_or_infinite"
	_LB_panic                 = "_panic"
)

var (
	_AX = jit.Reg("AX")
	_BX = jit.Reg("BX")
	_CX = jit.Reg("CX")
	_DX = jit.Reg("DX")
	_DI = jit.Reg("DI")
	_SI = jit.Reg("SI")
	_BP = jit.Reg("BP")
	_SP = jit.Reg("SP")
	_R8 = jit.Reg("R8")
	_R9 = jit.Reg("R9")
)

var (
	_X0 = jit.Reg("X0")
	_Y0 = jit.Reg("Y0")
)

var (
	_ST = jit.Reg("R15") // can't use R14 since it's always scratched by Go...
	_RP = jit.Reg("DI")
	_RL = jit.Reg("SI")
	_RC = jit.Reg("DX")
)

var (
	_LR = jit.Reg("R9")
	_ET = jit.Reg("AX")
	_EP = jit.Reg("BX")
)

var (
	_SP_p = jit.Reg("R10") // saved on BX when call_c
	_SP_q = jit.Reg("R11") // saved on BP when call_c
	_SP_x = jit.Reg("R12")
	_SP_f = jit.Reg("R13")
)

var (
	_ARG_rb = jit.Ptr(_SP, _FP_base)
	_ARG_vp = jit.Ptr(_SP, _FP_base+8)
	_ARG_sb = jit.Ptr(_SP, _FP_base+16)
	_ARG_fv = jit.Ptr(_SP, _FP_base+24)
)

var (
	_RET_et = _ET
	_RET_ep = _EP
)

var (
	_VAR_sp = jit.Ptr(_SP, _FP_fargs+_FP_saves)
	_VAR_dn = jit.Ptr(_SP, _FP_fargs+_FP_saves+8)
	_VAR_vp = jit.Ptr(_SP, _FP_fargs+_FP_saves+16)
)

var (
	_REG_ffi = []obj.Addr{_RP, _RL, _RC, _SP_q}
	_REG_b64 = []obj.Addr{_SP_p, _SP_q}

	_REG_all = []obj.Addr{_ST, _SP_x, _SP_f, _SP_p, _SP_q, _RP, _RL, _RC}
	_REG_ms  = []obj.Addr{_ST, _SP_x, _SP_f, _SP_p, _SP_q, _LR}
	_REG_enc = []obj.Addr{_ST, _SP_x, _SP_f, _SP_p, _SP_q, _RL}
)

type Assembler struct {
	Name string
	jit.BaseAssembler
	p    ir.Program
	x    int
}

func NewAssembler(p ir.Program) *Assembler {
	return new(Assembler).Init(p)
}

/** Assembler Interface **/

func (self *Assembler) Load() vars.Encoder {
	return ptoenc(self.BaseAssembler.Load("encode_"+self.Name, _FP_size, _FP_args, vars.ArgPtrs, vars.LocalPtrs))
}

func (self *Assembler) Init(p ir.Program) *Assembler {
	self.p = p
	self.BaseAssembler.Init(self.compile)
	return self
}

func (self *Assembler) compile() {
	self.prologue()
	self.instrs()
	self.epilogue()
	self.builtins()
}

/** Assembler Stages **/

var _OpFuncTab = [256]func(*Assembler, *ir.Instr){
	ir.OP_null:           (*Assembler)._asm_OP_null,
	ir.OP_empty_arr:      (*Assembler)._asm_OP_empty_arr,
	ir.OP_empty_obj:      (*Assembler)._asm_OP_empty_obj,
	ir.OP_bool:           (*Assembler)._asm_OP_bool,
	ir.OP_i8:             (*Assembler)._asm_OP_i8,
	ir.OP_i16:            (*Assembler)._asm_OP_i16,
	ir.OP_i32:            (*Assembler)._asm_OP_i32,
	ir.OP_i64:            (*Assembler)._asm_OP_i64,
	ir.OP_u8:             (*Assembler)._asm_OP_u8,
	ir.OP_u16:            (*Assembler)._asm_OP_u16,
	ir.OP_u32:            (*Assembler)._asm_OP_u32,
	ir.OP_u64:            (*Assembler)._asm_OP_u64,
	ir.OP_f32:            (*Assembler)._asm_OP_f32,
	ir.OP_f64:            (*Assembler)._asm_OP_f64,
	ir.OP_str:            (*Assembler)._asm_OP_str,
	ir.OP_bin:            (*Assembler)._asm_OP_bin,
	ir.OP_quote:          (*Assembler)._asm_OP_quote,
	ir.OP_number:         (*Assembler)._asm_OP_number,
	ir.OP_eface:          (*Assembler)._asm_OP_eface,
	ir.OP_iface:          (*Assembler)._asm_OP_iface,
	ir.OP_byte:           (*Assembler)._asm_OP_byte,
	ir.OP_text:           (*Assembler)._asm_OP_text,
	ir.OP_deref:          (*Assembler)._asm_OP_deref,
	ir.OP_index:          (*Assembler)._asm_OP_index,
	ir.OP_load:           (*Assembler)._asm_OP_load,
	ir.OP_save:           (*Assembler)._asm_OP_save,
	ir.OP_drop:           (*Assembler)._asm_OP_drop,
	ir.OP_drop_2:         (*Assembler)._asm_OP_drop_2,
	ir.OP_recurse:        (*Assembler)._asm_OP_recurse,
	ir.OP_is_nil:         (*Assembler)._asm_OP_is_nil,
	ir.OP_is_nil_p1:      (*Assembler)._asm_OP_is_nil_p1,
	ir.OP_is_zero_1:      (*Assembler)._asm_OP_is_zero_1,
	ir.OP_is_zero_2:      (*Assembler)._asm_OP_is_zero_2,
	ir.OP_is_zero_4:      (*Assembler)._asm_OP_is_zero_4,
	ir.OP_is_zero_8:      (*Assembler)._asm_OP_is_zero_8,
	ir.OP_is_zero_map:    (*Assembler)._asm_OP_is_zero_map,
	ir.OP_goto:           (*Assembler)._asm_OP_goto,
	ir.OP_map_iter:       (*Assembler)._asm_OP_map_iter,
	ir.OP_map_stop:       (*Assembler)._asm_OP_map_stop,
	ir.OP_map_check_key:  (*Assembler)._asm_OP_map_check_key,
	ir.OP_map_write_key:  (*Assembler)._asm_OP_map_write_key,
	ir.OP_map_value_next: (*Assembler)._asm_OP_map_value_next,
	ir.OP_slice_len:      (*Assembler)._asm_OP_slice_len,
	ir.OP_slice_next:     (*Assembler)._asm_OP_slice_next,
	ir.OP_marshal:        (*Assembler)._asm_OP_marshal,
	ir.OP_marshal_p:      (*Assembler)._asm_OP_marshal_p,
	ir.OP_marshal_text:   (*Assembler)._asm_OP_marshal_text,
	ir.OP_marshal_text_p: (*Assembler)._asm_OP_marshal_text_p,
	ir.OP_cond_set:       (*Assembler)._asm_OP_cond_set,
	ir.OP_cond_testc:     (*Assembler)._asm_OP_cond_testc,
	ir.OP_unsupported:    (*Assembler)._asm_OP_unsupported,
	ir.OP_is_zero:        (*Assembler)._asm_OP_is_zero,
}

func (self *Assembler) instr(v *ir.Instr) {
	if fn := _OpFuncTab[v.Op()]; fn != nil {
		fn(self, v)
	} else {
		panic(fmt.Sprintf("invalid opcode: %d", v.Op()))
	}
}

func (self *Assembler) instrs() {
	for i, v := range self.p {
		self.Mark(i)
		self.instr(&v)
		self.debug_instr(i, &v)
	}
}

func (self *Assembler) builtins() {
	self.more_space()
	self.error_too_deep()
	self.error_invalid_number()
	self.error_nan_or_infinite()
	self.go_panic()
}

func (self *Assembler) epilogue() {
	self.Mark(len(self.p))
	self.Emit("XORL", _ET, _ET)
	self.Emit("XORL", _EP, _EP)
	self.Link(_LB_error)
	self.Emit("MOVQ", _ARG_rb, _CX)                // MOVQ rb<>+0(FP), CX
	self.Emit("MOVQ", _RL, jit.Ptr(_CX, 8))        // MOVQ RL, 8(CX)
	self.Emit("MOVQ", jit.Imm(0), _ARG_rb)         // MOVQ AX, rb<>+0(FP)
	self.Emit("MOVQ", jit.Imm(0), _ARG_vp)         // MOVQ BX, vp<>+8(FP)
	self.Emit("MOVQ", jit.Imm(0), _ARG_sb)         // MOVQ CX, sb<>+16(FP)
	self.Emit("MOVQ", jit.Ptr(_SP, FP_offs), _BP) // MOVQ _FP_offs(SP), BP
	self.Emit("ADDQ", jit.Imm(_FP_size), _SP)      // ADDQ $_FP_size, SP
	self.Emit("RET")                               // RET
}

func (self *Assembler) prologue() {
	self.Emit("SUBQ", jit.Imm(_FP_size), _SP)      // SUBQ $_FP_size, SP
	self.Emit("MOVQ", _BP, jit.Ptr(_SP, FP_offs)) // MOVQ BP, _FP_offs(SP)
	self.Emit("LEAQ", jit.Ptr(_SP, FP_offs), _BP) // LEAQ _FP_offs(SP), BP
	self.Emit("MOVQ", _AX, _ARG_rb)                // MOVQ AX, rb<>+0(FP)
	self.Emit("MOVQ", _BX, _ARG_vp)                // MOVQ BX, vp<>+8(FP)
	self.Emit("MOVQ", _CX, _ARG_sb)                // MOVQ CX, sb<>+16(FP)
	self.Emit("MOVQ", _DI, _ARG_fv)                // MOVQ DI, rb<>+24(FP)
	self.Emit("MOVQ", jit.Ptr(_AX, 0), _RP)        // MOVQ (AX)  , DI
	self.Emit("MOVQ", jit.Ptr(_AX, 8), _RL)        // MOVQ 8(AX) , SI
	self.Emit("MOVQ", jit.Ptr(_AX, 16), _RC)       // MOVQ 16(AX), DX
	self.Emit("MOVQ", _BX, _SP_p)                  // MOVQ BX, R10
	self.Emit("MOVQ", _CX, _ST)                    // MOVQ CX, R8
	self.Emit("XORL", _SP_x, _SP_x)                // XORL R10, R12
	self.Emit("XORL", _SP_f, _SP_f)                // XORL R11, R13
	self.Emit("XORL", _SP_q, _SP_q)                // XORL R13, R11
}

/** Assembler Inline Functions **/

func (self *Assembler) xsave(reg ...obj.Addr) {
	for i, v := range reg {
		if i > _FP_saves/8-1 {
			panic("too many registers to save")
		} else {
			self.Emit("MOVQ", v, jit.Ptr(_SP, _FP_fargs+int64(i)*8))
		}
	}
}

func (self *Assembler) xload(reg ...obj.Addr) {
	for i, v := range reg {
		if i > _FP_saves/8-1 {
			panic("too many registers to load")
		} else {
			self.Emit("MOVQ", jit.Ptr(_SP, _FP_fargs+int64(i)*8), v)
		}
	}
}

func (self *Assembler) rbuf_di() {
	if _RP.Reg != x86.REG_DI {
		panic("register allocation messed up: RP != DI")
	} else {
		self.Emit("ADDQ", _RL, _RP)
	}
}

func (self *Assembler) store_int(nd int, fn obj.Addr, ins string) {
	self.check_size(nd)
	self.save_c()                          // SAVE   $C_regs
	self.rbuf_di()                         // MOVQ   RP, DI
	self.Emit(ins, jit.Ptr(_SP_p, 0), _SI) // $ins   (SP.p), SI
	self.call_c(fn)                        // CALL_C $fn
	self.Emit("ADDQ", _AX, _RL)            // ADDQ   AX, RL
}

func (self *Assembler) store_str(s string) {
	i := 0
	m := rt.Str2Mem(s)

	/* 8-byte stores */
	for i <= len(m)-8 {
		self.Emit("MOVQ", jit.Imm(rt.Get64(m[i:])), _AX)       // MOVQ $s[i:], AX
		self.Emit("MOVQ", _AX, jit.Sib(_RP, _RL, 1, int64(i))) // MOVQ AX, i(RP)(RL)
		i += 8
	}

	/* 4-byte stores */
	if i <= len(m)-4 {
		self.Emit("MOVL", jit.Imm(int64(rt.Get32(m[i:]))), jit.Sib(_RP, _RL, 1, int64(i))) // MOVL $s[i:], i(RP)(RL)
		i += 4
	}

	/* 2-byte stores */
	if i <= len(m)-2 {
		self.Emit("MOVW", jit.Imm(int64(rt.Get16(m[i:]))), jit.Sib(_RP, _RL, 1, int64(i))) // MOVW $s[i:], i(RP)(RL)
		i += 2
	}

	/* last byte */
	if i < len(m) {
		self.Emit("MOVB", jit.Imm(int64(m[i])), jit.Sib(_RP, _RL, 1, int64(i))) // MOVB $s[i:], i(RP)(RL)
	}
}

func (self *Assembler) check_size(n int) {
	self.check_size_rl(jit.Ptr(_RL, int64(n)))
}

func (self *Assembler) check_size_r(r obj.Addr, d int) {
	self.check_size_rl(jit.Sib(_RL, r, 1, int64(d)))
}

func (self *Assembler) check_size_rl(v obj.Addr) {
	idx := self.x
	key := _LB_more_space_return + strconv.Itoa(idx)

	/* the following code relies on LR == R9 to work */
	if _LR.Reg != x86.REG_R9 {
		panic("register allocation messed up: LR != R9")
	}

	/* check for buffer capacity */
	self.x++
	self.Emit("LEAQ", v, _AX)   // LEAQ $v, AX
	self.Emit("CMPQ", _AX, _RC) // CMPQ AX, RC
	self.Sjmp("JBE", key)       // JBE  _more_space_return_{n}
	self.slice_grow_ax(key)     // GROW $key
	self.Link(key)              // _more_space_return_{n}:
}

func (self *Assembler) slice_grow_ax(ret string) {
	self.Byte(0x4c, 0x8d, 0x0d)      // LEAQ ?(PC), R9
	self.Sref(ret, 4)                // .... &ret
	self.Sjmp("JMP", _LB_more_space) // JMP  _more_space
}

/** State Stack Helpers **/



func (self *Assembler) save_state() {
	self.Emit("MOVQ", jit.Ptr(_ST, 0), _CX)            // MOVQ (ST), CX
	self.Emit("LEAQ", jit.Ptr(_CX, vars.StateSize), _R9)   // LEAQ vars.StateSize(CX), R9
	self.Emit("CMPQ", _R9, jit.Imm(vars.StackLimit))       // CMPQ R9, $vars.StackLimit
	self.Sjmp("JAE", _LB_error_too_deep)               // JA   _error_too_deep
	self.Emit("MOVQ", _SP_x, jit.Sib(_ST, _CX, 1, 8))  // MOVQ SP.x, 8(ST)(CX)
	self.Emit("MOVQ", _SP_f, jit.Sib(_ST, _CX, 1, 16)) // MOVQ SP.f, 16(ST)(CX)
	self.WritePtr(0, _SP_p, jit.Sib(_ST, _CX, 1, 24))  // MOVQ SP.p, 24(ST)(CX)
	self.WritePtr(1, _SP_q, jit.Sib(_ST, _CX, 1, 32))  // MOVQ SP.q, 32(ST)(CX)
	self.Emit("MOVQ", _R9, jit.Ptr(_ST, 0))            // MOVQ R9, (ST)
}

func (self *Assembler) drop_state(decr int64) {
	self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)            // MOVQ  (ST), AX
	self.Emit("SUBQ", jit.Imm(decr), _AX)              // SUBQ  $decr, AX
	self.Emit("MOVQ", _AX, jit.Ptr(_ST, 0))            // MOVQ  AX, (ST)
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 8), _SP_x)  // MOVQ  8(ST)(AX), SP.x
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 16), _SP_f) // MOVQ  16(ST)(AX), SP.f
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 24), _SP_p) // MOVQ  24(ST)(AX), SP.p
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 32), _SP_q) // MOVQ  32(ST)(AX), SP.q
	self.Emit("PXOR", _X0, _X0)                        // PXOR  X0, X0
	self.Emit("MOVOU", _X0, jit.Sib(_ST, _AX, 1, 8))   // MOVOU X0, 8(ST)(AX)
	self.Emit("MOVOU", _X0, jit.Sib(_ST, _AX, 1, 24))  // MOVOU X0, 24(ST)(AX)
}

/** Buffer Helpers **/

func (self *Assembler) add_char(ch byte) {
	self.Emit("MOVB", jit.Imm(int64(ch)), jit.Sib(_RP, _RL, 1, 0)) // MOVB $ch, (RP)(RL)
	self.Emit("ADDQ", jit.Imm(1), _RL)                             // ADDQ $1, RL
}

func (self *Assembler) add_long(ch uint32, n int64) {
	self.Emit("MOVL", jit.Imm(int64(ch)), jit.Sib(_RP, _RL, 1, 0)) // MOVL $ch, (RP)(RL)
	self.Emit("ADDQ", jit.Imm(n), _RL)                             // ADDQ $n, RL
}

func (self *Assembler) add_text(ss string) {
	self.store_str(ss)                              // TEXT $ss
	self.Emit("ADDQ", jit.Imm(int64(len(ss))), _RL) // ADDQ ${len(ss)}, RL
}

// get *buf at AX
func (self *Assembler) prep_buffer_AX() {
	self.Emit("MOVQ", _ARG_rb, _AX)         // MOVQ rb<>+0(FP), AX
	self.Emit("MOVQ", _RL, jit.Ptr(_AX, 8)) // MOVQ RL, 8(AX)
}

func (self *Assembler) save_buffer() {
	self.Emit("MOVQ", _ARG_rb, _CX)          // MOVQ rb<>+0(FP), CX
	self.Emit("MOVQ", _RP, jit.Ptr(_CX, 0))  // MOVQ RP, (CX)
	self.Emit("MOVQ", _RL, jit.Ptr(_CX, 8))  // MOVQ RL, 8(CX)
	self.Emit("MOVQ", _RC, jit.Ptr(_CX, 16)) // MOVQ RC, 16(CX)
}

// get *buf at AX
func (self *Assembler) load_buffer_AX() {
	self.Emit("MOVQ", _ARG_rb, _AX)          // MOVQ rb<>+0(FP), AX
	self.Emit("MOVQ", jit.Ptr(_AX, 0), _RP)  // MOVQ (AX), RP
	self.Emit("MOVQ", jit.Ptr(_AX, 8), _RL)  // MOVQ 8(AX), RL
	self.Emit("MOVQ", jit.Ptr(_AX, 16), _RC) // MOVQ 16(AX), RC
}

/** Function Interface Helpers **/

func (self *Assembler) call(pc obj.Addr) {
	self.Emit("MOVQ", pc, _LR) // MOVQ $pc, AX
	self.Rjmp("CALL", _LR)     // CALL AX
}

func (self *Assembler) save_c() {
	self.xsave(_REG_ffi...) // SAVE $REG_ffi
}

func (self *Assembler) call_b64(pc obj.Addr) {
	self.xsave(_REG_b64...) // SAVE $REG_all
	self.call(pc)           // CALL $pc
	self.xload(_REG_b64...) // LOAD $REG_ffi
}

func (self *Assembler) call_c(pc obj.Addr) {
	self.Emit("XCHGQ", _SP_p, _BX)
	self.call(pc)           // CALL $pc
	self.xload(_REG_ffi...) // LOAD $REG_ffi
	self.Emit("XCHGQ", _SP_p, _BX)
}

func (self *Assembler) call_go(pc obj.Addr) {
	self.xsave(_REG_all...) // SAVE $REG_all
	self.call(pc)           // CALL $pc
	self.xload(_REG_all...) // LOAD $REG_all
}

func (self *Assembler) call_more_space(pc obj.Addr) {
	self.xsave(_REG_ms...) // SAVE $REG_all
	self.call(pc)          // CALL $pc
	self.xload(_REG_ms...) // LOAD $REG_all
}

func (self *Assembler) call_encoder(pc obj.Addr) {
	self.xsave(_REG_enc...) // SAVE $REG_all
	self.call(pc)           // CALL $pc
	self.xload(_REG_enc...) // LOAD $REG_all
}

func (self *Assembler) call_marshaler(fn obj.Addr, it *rt.GoType, vt reflect.Type) {
	switch vt.Kind() {
	case reflect.Interface:
		self.call_marshaler_i(fn, it)
	case reflect.Ptr, reflect.Map:
		self.call_marshaler_v(fn, it, vt, true)
	// struct/array of 1 direct iface type can be direct
	default:
		self.call_marshaler_v(fn, it, vt, !rt.UnpackType(vt).Indirect())
	}
}

var (
	_F_assertI2I = jit.Func(rt.AssertI2I)
)

func (self *Assembler) call_marshaler_i(fn obj.Addr, it *rt.GoType) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _AX) // MOVQ    (SP.p), AX
	self.Emit("TESTQ", _AX, _AX)              // TESTQ   AX, AX
	self.Sjmp("JZ", "_null_{n}")              // JZ      _null_{n}
	self.Emit("MOVQ", _AX, _BX)               // MOVQ    AX, BX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _CX) // MOVQ    8(SP.p), CX
	self.Emit("MOVQ", jit.Gtype(it), _AX)     // MOVQ    $it, AX
	self.call_go(_F_assertI2I)                // CALL_GO assertI2I
	self.Emit("TESTQ", _AX, _AX)              // TESTQ   AX, AX
	self.Sjmp("JZ", "_null_{n}")              // JZ      _null_{n}
	self.Emit("MOVQ", _BX, _CX)               // MOVQ   BX, CX
	self.Emit("MOVQ", _AX, _BX)               // MOVQ   AX, BX
	self.prep_buffer_AX()
	self.Emit("MOVQ", _ARG_fv, _DI) // MOVQ   ARG.fv, DI
	self.call_go(fn)                // CALL    $fn
	self.Emit("TESTQ", _ET, _ET)    // TESTQ ET, ET
	self.Sjmp("JNZ", _LB_error)     // JNZ   _error
	self.load_buffer_AX()
	self.Sjmp("JMP", "_done_{n}")                                 // JMP     _done_{n}
	self.Link("_null_{n}")                                        // _null_{n}:
	self.check_size(4)                                            // SIZE    $4
	self.Emit("MOVL", jit.Imm(_IM_null), jit.Sib(_RP, _RL, 1, 0)) // MOVL    $'null', (RP)(RL*1)
	self.Emit("ADDQ", jit.Imm(4), _RL)                            // ADDQ    $4, RL
	self.Link("_done_{n}")                                        // _done_{n}:
}

func (self *Assembler) call_marshaler_v(fn obj.Addr, it *rt.GoType, vt reflect.Type, deref bool) {
	self.prep_buffer_AX()                    // MOVE {buf}, (SP)
	self.Emit("MOVQ", jit.Itab(it, vt), _BX) // MOVQ $(itab(it, vt)), BX

	/* dereference the pointer if needed */
	if !deref {
		self.Emit("MOVQ", _SP_p, _CX) // MOVQ SP.p, CX
	} else {
		self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _CX) // MOVQ 0(SP.p), CX
	}

	/* call the encoder, and perform error checks */
	self.Emit("MOVQ", _ARG_fv, _DI) // MOVQ   ARG.fv, DI
	self.call_go(fn)                // CALL  $fn
	self.Emit("TESTQ", _ET, _ET)    // TESTQ ET, ET
	self.Sjmp("JNZ", _LB_error)     // JNZ   _error
	self.load_buffer_AX()
}

/** Builtin: _more_space **/

var (
	_T_byte      = jit.Type(vars.ByteType)
	_F_growslice = jit.Func(rt.GrowSlice)

	_T_json_Marshaler         = rt.UnpackType(vars.JsonMarshalerType)
	_T_encoding_TextMarshaler = rt.UnpackType(vars.EncodingTextMarshalerType)
)

// AX must saving n
func (self *Assembler) more_space() {
	self.Link(_LB_more_space)
	self.Emit("MOVQ", _RP, _BX)        // MOVQ DI, BX
	self.Emit("MOVQ", _RL, _CX)        // MOVQ SI, CX
	self.Emit("MOVQ", _RC, _DI)        // MOVQ DX, DI
	self.Emit("MOVQ", _AX, _SI)        // MOVQ AX, SI
	self.Emit("MOVQ", _T_byte, _AX)    // MOVQ $_T_byte, AX
	self.call_more_space(_F_growslice) // CALL $pc
	self.Emit("MOVQ", _AX, _RP)        // MOVQ AX, DI
	self.Emit("MOVQ", _BX, _RL)        // MOVQ BX, SI
	self.Emit("MOVQ", _CX, _RC)        // MOVQ CX, DX
	self.save_buffer()                 // SAVE {buf}
	self.Rjmp("JMP", _LR)              // JMP  LR
}

/** Builtin Errors **/

var (
	_V_ERR_too_deep               = jit.Imm(int64(uintptr(unsafe.Pointer(vars.ERR_too_deep))))
	_V_ERR_nan_or_infinite        = jit.Imm(int64(uintptr(unsafe.Pointer(vars.ERR_nan_or_infinite))))
	_I_json_UnsupportedValueError = jit.Itab(rt.UnpackType(vars.ErrorType), vars.JsonUnsupportedValueType)
)

func (self *Assembler) error_too_deep() {
	self.Link(_LB_error_too_deep)
	self.Emit("MOVQ", _V_ERR_too_deep, _EP)               // MOVQ $_V_ERR_too_deep, EP
	self.Emit("MOVQ", _I_json_UnsupportedValueError, _ET) // MOVQ $_I_json_UnsupportedValuError, ET
	self.Sjmp("JMP", _LB_error)                           // JMP  _error
}

func (self *Assembler) error_invalid_number() {
	self.Link(_LB_error_invalid_number)
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _AX) // MOVQ    0(SP), AX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _BX) // MOVQ    8(SP), BX
	self.call_go(_F_error_number)             // CALL_GO error_number
	self.Sjmp("JMP", _LB_error)               // JMP     _error
}

func (self *Assembler) error_nan_or_infinite() {
	self.Link(_LB_error_nan_or_infinite)
	self.Emit("MOVQ", _V_ERR_nan_or_infinite, _EP)        // MOVQ $_V_ERR_nan_or_infinite, EP
	self.Emit("MOVQ", _I_json_UnsupportedValueError, _ET) // MOVQ $_I_json_UnsupportedValuError, ET
	self.Sjmp("JMP", _LB_error)                           // JMP  _error
}

/** String Encoding Routine **/

var (
	_F_quote = jit.Imm(int64(native.S_quote))
	_F_panic = jit.Func(vars.GoPanic)
)

func (self *Assembler) go_panic() {
	self.Link(_LB_panic)
	self.Emit("MOVQ", _SP_p, _BX)
	self.Emit("MOVQ", _RP, _CX)
	self.Emit("MOVQ", _RL, _DI)
	self.call_go(_F_panic)
}

func (self *Assembler) encode_string(doubleQuote bool) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _AX) // MOVQ  8(SP.p), AX
	self.Emit("TESTQ", _AX, _AX)              // TESTQ AX, AX
	self.Sjmp("JZ", "_str_empty_{n}")         // JZ    _str_empty_{n}
	self.Emit("CMPQ", jit.Ptr(_SP_p, 0), jit.Imm(0))
	self.Sjmp("JNE", "_str_next_{n}")
	self.Emit("MOVQ", jit.Imm(int64(vars.PanicNilPointerOfNonEmptyString)), _AX)
	self.Sjmp("JMP", _LB_panic)
	self.Link("_str_next_{n}")

	/* opening quote, check for double quote */
	if !doubleQuote {
		self.check_size_r(_AX, 2) // SIZE $2
		self.add_char('"')        // CHAR $'"'
	} else {
		self.check_size_r(_AX, 6)  // SIZE $6
		self.add_long(_IM_open, 3) // TEXT $`"\"`
	}

	/* quoting loop */
	self.Emit("XORL", _AX, _AX)     // XORL AX, AX
	self.Emit("MOVQ", _AX, _VAR_sp) // MOVQ AX, sp
	self.Link("_str_loop_{n}")      // _str_loop_{n}:
	self.save_c()                   // SAVE $REG_ffi

	/* load the output buffer first, and then input buffer,
	 * because the parameter registers collide with RP / RL / RC */
	self.Emit("MOVQ", _RC, _CX)                     // MOVQ RC, CX
	self.Emit("SUBQ", _RL, _CX)                     // SUBQ RL, CX
	self.Emit("MOVQ", _CX, _VAR_dn)                 // MOVQ CX, dn
	self.Emit("LEAQ", jit.Sib(_RP, _RL, 1, 0), _DX) // LEAQ (RP)(RL), DX
	self.Emit("LEAQ", _VAR_dn, _CX)                 // LEAQ dn, CX
	self.Emit("MOVQ", _VAR_sp, _AX)                 // MOVQ sp, AX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _DI)       // MOVQ (SP.p), DI
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _SI)       // MOVQ 8(SP.p), SI
	self.Emit("ADDQ", _AX, _DI)                     // ADDQ AX, DI
	self.Emit("SUBQ", _AX, _SI)                     // SUBQ AX, SI

	/* set the flags based on `doubleQuote` */
	if !doubleQuote {
		self.Emit("XORL", _R8, _R8) // XORL R8, R8
	} else {
		self.Emit("MOVL", jit.Imm(types.F_DOUBLE_UNQUOTE), _R8) // MOVL ${types.F_DOUBLE_UNQUOTE}, R8
	}

	/* call the native quoter */
	self.call_c(_F_quote)           // CALL  quote
	self.Emit("ADDQ", _VAR_dn, _RL) // ADDQ  dn, RL

	self.Emit("TESTQ", _AX, _AX)      // TESTQ AX, AX
	self.Sjmp("JS", "_str_space_{n}") // JS    _str_space_{n}

	/* close the string, check for double quote */
	if !doubleQuote {
		self.check_size(1)               // SIZE $1
		self.add_char('"')               // CHAR $'"'
		self.Sjmp("JMP", "_str_end_{n}") // JMP  _str_end_{n}
	} else {
		self.check_size(3)               // SIZE $3
		self.add_text("\\\"\"")          // TEXT $'\""'
		self.Sjmp("JMP", "_str_end_{n}") // JMP  _str_end_{n}
	}

	/* not enough space to contain the quoted string */
	self.Link("_str_space_{n}")                     // _str_space_{n}:
	self.Emit("NOTQ", _AX)                          // NOTQ AX
	self.Emit("ADDQ", _AX, _VAR_sp)                 // ADDQ AX, sp
	self.Emit("LEAQ", jit.Sib(_RC, _RC, 1, 0), _AX) // LEAQ (RC)(RC), AX
	self.slice_grow_ax("_str_loop_{n}")             // GROW _str_loop_{n}

	/* empty string, check for double quote */
	if !doubleQuote {
		self.Link("_str_empty_{n}") // _str_empty_{n}:
		self.check_size(2)          // SIZE $2
		self.add_text("\"\"")       // TEXT $'""'
		self.Link("_str_end_{n}")   // _str_end_{n}:
	} else {
		self.Link("_str_empty_{n}")   // _str_empty_{n}:
		self.check_size(6)            // SIZE $6
		self.add_text("\"\\\"\\\"\"") // TEXT $'"\"\""'
		self.Link("_str_end_{n}")     // _str_end_{n}:
	}
}

/** OpCode Assembler Functions **/



var (
	_F_f64toa    = jit.Imm(int64(native.S_f64toa))
	_F_f32toa    = jit.Imm(int64(native.S_f32toa))
	_F_i64toa    = jit.Imm(int64(native.S_i64toa))
	_F_u64toa    = jit.Imm(int64(native.S_u64toa))
	_F_b64encode = jit.Imm(int64(rt.SubrB64Encode))
)

var (
	_F_memmove       = jit.Func(rt.Memmove)
	_F_error_number  = jit.Func(vars.Error_number)
	_F_isValidNumber = jit.Func(alg.IsValidNumber)
)

var (
	_F_iteratorStop  = jit.Func(alg.IteratorStop)
	_F_iteratorNext  = jit.Func(alg.IteratorNext)
	_F_iteratorStart = jit.Func(alg.IteratorStart)
)

var (
	_F_encodeTypedPointer  obj.Addr
	_F_encodeJsonMarshaler obj.Addr
	_F_encodeTextMarshaler obj.Addr
)

const (
	_MODE_AVX2 = 1 << 2
)

func init() {
	_F_encodeJsonMarshaler = jit.Func(prim.EncodeJsonMarshaler)
	_F_encodeTextMarshaler = jit.Func(prim.EncodeTextMarshaler)
	_F_encodeTypedPointer  = jit.Func(EncodeTypedPointer)
}

func (self *Assembler) _asm_OP_null(_ *ir.Instr) {
	self.check_size(4)
	self.Emit("MOVL", jit.Imm(_IM_null), jit.Sib(_RP, _RL, 1, 0)) // MOVL $'null', (RP)(RL*1)
	self.Emit("ADDQ", jit.Imm(4), _RL)                            // ADDQ $4, RL
}

func (self *Assembler) _asm_OP_empty_arr(_ *ir.Instr) {
	self.Emit("BTQ", jit.Imm(int64(alg.BitNoNullSliceOrMap)), _ARG_fv)
	self.Sjmp("JC", "_empty_arr_{n}")
	self._asm_OP_null(nil)
	self.Sjmp("JMP", "_empty_arr_end_{n}")
	self.Link("_empty_arr_{n}")
	self.check_size(2)
	self.Emit("MOVW", jit.Imm(_IM_array), jit.Sib(_RP, _RL, 1, 0))
	self.Emit("ADDQ", jit.Imm(2), _RL)
	self.Link("_empty_arr_end_{n}")
}

func (self *Assembler) _asm_OP_empty_obj(_ *ir.Instr) {
	self.Emit("BTQ", jit.Imm(int64(alg.BitNoNullSliceOrMap)), _ARG_fv)
	self.Sjmp("JC", "_empty_obj_{n}")
	self._asm_OP_null(nil)
	self.Sjmp("JMP", "_empty_obj_end_{n}")
	self.Link("_empty_obj_{n}")
	self.check_size(2)
	self.Emit("MOVW", jit.Imm(_IM_object), jit.Sib(_RP, _RL, 1, 0))
	self.Emit("ADDQ", jit.Imm(2), _RL)
	self.Link("_empty_obj_end_{n}")
}

func (self *Assembler) _asm_OP_bool(_ *ir.Instr) {
	self.Emit("CMPB", jit.Ptr(_SP_p, 0), jit.Imm(0))              // CMPB (SP.p), $0
	self.Sjmp("JE", "_false_{n}")                                 // JE   _false_{n}
	self.check_size(4)                                            // SIZE $4
	self.Emit("MOVL", jit.Imm(_IM_true), jit.Sib(_RP, _RL, 1, 0)) // MOVL $'true', (RP)(RL*1)
	self.Emit("ADDQ", jit.Imm(4), _RL)                            // ADDQ $4, RL
	self.Sjmp("JMP", "_end_{n}")                                  // JMP  _end_{n}
	self.Link("_false_{n}")                                       // _false_{n}:
	self.check_size(5)                                            // SIZE $5
	self.Emit("MOVL", jit.Imm(_IM_fals), jit.Sib(_RP, _RL, 1, 0)) // MOVL $'fals', (RP)(RL*1)
	self.Emit("MOVB", jit.Imm('e'), jit.Sib(_RP, _RL, 1, 4))      // MOVB $'e', 4(RP)(RL*1)
	self.Emit("ADDQ", jit.Imm(5), _RL)                            // ADDQ $5, RL
	self.Link("_end_{n}")                                         // _end_{n}:
}

func (self *Assembler) _asm_OP_i8(_ *ir.Instr) {
	self.store_int(4, _F_i64toa, "MOVBQSX")
}

func (self *Assembler) _asm_OP_i16(_ *ir.Instr) {
	self.store_int(6, _F_i64toa, "MOVWQSX")
}

func (self *Assembler) _asm_OP_i32(_ *ir.Instr) {
	self.store_int(17, _F_i64toa, "MOVLQSX")
}

func (self *Assembler) _asm_OP_i64(_ *ir.Instr) {
	self.store_int(21, _F_i64toa, "MOVQ")
}

func (self *Assembler) _asm_OP_u8(_ *ir.Instr) {
	self.store_int(3, _F_u64toa, "MOVBQZX")
}

func (self *Assembler) _asm_OP_u16(_ *ir.Instr) {
	self.store_int(5, _F_u64toa, "MOVWQZX")
}

func (self *Assembler) _asm_OP_u32(_ *ir.Instr) {
	self.store_int(16, _F_u64toa, "MOVLQZX")
}

func (self *Assembler) _asm_OP_u64(_ *ir.Instr) {
	self.store_int(20, _F_u64toa, "MOVQ")
}

func (self *Assembler) _asm_OP_f32(_ *ir.Instr) {
	self.check_size(32)
	self.Emit("MOVL", jit.Ptr(_SP_p, 0), _AX)  // MOVL     (SP.p), AX
	self.Emit("ANDL", jit.Imm(_FM_exp32), _AX) // ANDL     $_FM_exp32, AX
	self.Emit("XORL", jit.Imm(_FM_exp32), _AX) // XORL     $_FM_exp32, AX
	self.Sjmp("JNZ",  "_encode_normal_f32_{n}")// JNZ      _encode_normal_f32_{n}
	self.Emit("BTQ", jit.Imm(alg.BitEncodeNullForInfOrNan), _ARG_fv) // BTQ ${BitEncodeNullForInfOrNan}, fv
	self.Sjmp("JNC", _LB_error_nan_or_infinite) // JNC     _error_nan_or_infinite
	self._asm_OP_null(nil)
	self.Sjmp("JMP", "_encode_f32_end_{n}")    // JMP      _encode_f32_end_{n}
	self.Link("_encode_normal_f32_{n}")
	self.save_c()                              // SAVE     $C_regs
	self.rbuf_di()                             // MOVQ     RP, DI
	self.Emit("MOVSS", jit.Ptr(_SP_p, 0), _X0) // MOVSS    (SP.p), X0
	self.call_c(_F_f32toa)                     // CALL_C   f32toa
	self.Emit("ADDQ", _AX, _RL)                // ADDQ     AX, RL
	self.Link("_encode_f32_end_{n}")
}

func (self *Assembler) _asm_OP_f64(_ *ir.Instr) {
	self.check_size(32)
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _AX)  // MOVQ   (SP.p), AX
	self.Emit("MOVQ", jit.Imm(_FM_exp64), _CX) // MOVQ   $_FM_exp64, CX
	self.Emit("ANDQ", _CX, _AX)                // ANDQ   CX, AX
	self.Emit("XORQ", _CX, _AX)                // XORQ   CX, AX
	self.Sjmp("JNZ",  "_encode_normal_f64_{n}")// JNZ    _encode_normal_f64_{n}
	self.Emit("BTQ", jit.Imm(alg.BitEncodeNullForInfOrNan), _ARG_fv) // BTQ ${BitEncodeNullForInfOrNan}, fv
	self.Sjmp("JNC", _LB_error_nan_or_infinite)// JNC    _error_nan_or_infinite
	self._asm_OP_null(nil)
	self.Sjmp("JMP", "_encode_f64_end_{n}")    // JMP    _encode_f64_end_{n}
	self.Link("_encode_normal_f64_{n}")
	self.save_c()                              // SAVE   $C_regs
	self.rbuf_di()                             // MOVQ   RP, DI
	self.Emit("MOVSD", jit.Ptr(_SP_p, 0), _X0) // MOVSD  (SP.p), X0
	self.call_c(_F_f64toa)                     // CALL_C f64toa
	self.Emit("ADDQ", _AX, _RL)                // ADDQ   AX, RL
	self.Link("_encode_f64_end_{n}")
}

func (self *Assembler) _asm_OP_str(_ *ir.Instr) {
	self.encode_string(false)
}

func (self *Assembler) _asm_OP_bin(_ *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _AX)       // MOVQ 8(SP.p), AX
	self.Emit("ADDQ", jit.Imm(2), _AX)              // ADDQ $2, AX
	self.Emit("MOVQ", jit.Imm(_IM_mulv), _CX)       // MOVQ $_MF_mulv, CX
	self.Emit("MOVQ", _DX, _BX)                     // MOVQ DX, BX
	self.From("MULQ", _CX)                          // MULQ CX
	self.Emit("LEAQ", jit.Sib(_DX, _DX, 1, 1), _AX) // LEAQ 1(DX)(DX), AX
	self.Emit("ORQ", jit.Imm(2), _AX)               // ORQ  $2, AX
	self.Emit("MOVQ", _BX, _DX)                     // MOVQ BX, DX
	self.check_size_r(_AX, 0)                       // SIZE AX
	self.add_char('"')                              // CHAR $'"'
	self.Emit("MOVQ", _ARG_rb, _DI)                 // MOVQ rb<>+0(FP), DI
	self.Emit("MOVQ", _RL, jit.Ptr(_DI, 8))         // MOVQ SI, 8(DI)
	self.Emit("MOVQ", _SP_p, _SI)                   // MOVQ SP.p, SI

	/* check for AVX2 support */
	if !cpu.HasAVX2 {
		self.Emit("XORL", _DX, _DX) // XORL DX, DX
	} else {
		self.Emit("MOVL", jit.Imm(_MODE_AVX2), _DX) // MOVL $_MODE_AVX2, DX
	}

	/* call the encoder */
	self.call_b64(_F_b64encode) // CALL b64encode
	self.load_buffer_AX()       // LOAD {buf}
	self.add_char('"')          // CHAR $'"'
}

func (self *Assembler) _asm_OP_quote(_ *ir.Instr) {
	self.encode_string(true)
}

func (self *Assembler) _asm_OP_number(_ *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _BX) // MOVQ    (SP.p), BX
	self.Emit("TESTQ", _BX, _BX)              // TESTQ   BX, BX
	self.Sjmp("JZ", "_empty_{n}")
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _AX) // MOVQ    (SP.p), AX
	self.Emit("TESTQ", _AX, _AX)              // TESTQ   AX, AX
	self.Sjmp("JNZ", "_number_next_{n}")
	self.Emit("MOVQ", jit.Imm(int64(vars.PanicNilPointerOfNonEmptyString)), _AX)
	self.Sjmp("JMP", _LB_panic)
	self.Link("_number_next_{n}")
	self.call_go(_F_isValidNumber)                  // CALL_GO isValidNumber
	self.Emit("CMPB", _AX, jit.Imm(0))              // CMPB    AX, $0
	self.Sjmp("JE", _LB_error_invalid_number)       // JE      _error_invalid_number
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _BX)       // MOVQ    (SP.p), BX
	self.check_size_r(_BX, 0)                       // SIZE    BX
	self.Emit("LEAQ", jit.Sib(_RP, _RL, 1, 0), _AX) // LEAQ    (RP)(RL), AX
	self.Emit("ADDQ", jit.Ptr(_SP_p, 8), _RL)       // ADDQ    8(SP.p), RL
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _BX)       // MOVOU   (SP.p), BX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _CX)       // MOVOU   X0, 8(SP)
	self.call_go(_F_memmove)                        // CALL_GO memmove
	self.Emit("MOVQ", _ARG_rb, _AX)                 // MOVQ rb<>+0(FP), AX
	self.Emit("MOVQ", _RL, jit.Ptr(_AX, 8))         // MOVQ RL, 8(AX)
	self.Sjmp("JMP", "_done_{n}")                   // JMP     _done_{n}
	self.Link("_empty_{n}")                         // _empty_{n}
	self.check_size(1)                              // SIZE    $1
	self.add_char('0')                              // CHAR    $'0'
	self.Link("_done_{n}")                          // _done_{n}:
}

func (self *Assembler) _asm_OP_eface(_ *ir.Instr) {
	self.prep_buffer_AX()                     // MOVE  {buf}, AX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _BX) // MOVQ  (SP.p), BX
	self.Emit("LEAQ", jit.Ptr(_SP_p, 8), _CX) // LEAQ  8(SP.p), CX
	self.Emit("MOVQ", _ST, _DI)               // MOVQ  ST, DI
	self.Emit("MOVQ", _ARG_fv, _SI)           // MOVQ  fv, AX
	self.call_encoder(_F_encodeTypedPointer)  // CALL  encodeTypedPointer
	self.Emit("TESTQ", _ET, _ET)              // TESTQ ET, ET
	self.Sjmp("JNZ", _LB_error)               // JNZ   _error
	self.load_buffer_AX()
}

func (self *Assembler) _asm_OP_iface(_ *ir.Instr) {
	self.prep_buffer_AX()                     // MOVE  {buf}, AX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _CX) // MOVQ  (SP.p), CX
	self.Emit("MOVQ", jit.Ptr(_CX, 8), _BX)   // MOVQ  8(CX), BX
	self.Emit("LEAQ", jit.Ptr(_SP_p, 8), _CX) // LEAQ  8(SP.p), CX
	self.Emit("MOVQ", _ST, _DI)               // MOVQ  ST, DI
	self.Emit("MOVQ", _ARG_fv, _SI)           // MOVQ  fv, AX
	self.call_encoder(_F_encodeTypedPointer)  // CALL  encodeTypedPointer
	self.Emit("TESTQ", _ET, _ET)              // TESTQ ET, ET
	self.Sjmp("JNZ", _LB_error)               // JNZ   _error
	self.load_buffer_AX()
}

func (self *Assembler) _asm_OP_byte(p *ir.Instr) {
	self.check_size(1)
	self.Emit("MOVB", jit.Imm(p.I64()), jit.Sib(_RP, _RL, 1, 0)) // MOVL p.Vi(), (RP)(RL*1)
	self.Emit("ADDQ", jit.Imm(1), _RL)                           // ADDQ $1, RL
}

func (self *Assembler) _asm_OP_text(p *ir.Instr) {
	self.check_size(len(p.Vs())) // SIZE ${len(p.Vs())}
	self.add_text(p.Vs())        // TEXT ${p.Vs()}
}

func (self *Assembler) _asm_OP_deref(_ *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _SP_p) // MOVQ (SP.p), SP.p
}

func (self *Assembler) _asm_OP_index(p *ir.Instr) {
	self.Emit("MOVQ", jit.Imm(p.I64()), _AX) // MOVQ $p.Vi(), AX
	self.Emit("ADDQ", _AX, _SP_p)            // ADDQ AX, SP.p
}

func (self *Assembler) _asm_OP_load(_ *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             // MOVQ (ST), AX
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, -24), _SP_x) // MOVQ -24(ST)(AX), SP.x
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, -8), _SP_p)  // MOVQ -8(ST)(AX), SP.p
	self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 0), _SP_q)   // MOVQ (ST)(AX), SP.q
}

func (self *Assembler) _asm_OP_save(_ *ir.Instr) {
	self.save_state()
}

func (self *Assembler) _asm_OP_drop(_ *ir.Instr) {
	self.drop_state(vars.StateSize)
}

func (self *Assembler) _asm_OP_drop_2(_ *ir.Instr) {
	self.drop_state(vars.StateSize * 2)                   // DROP  $(vars.StateSize * 2)
	self.Emit("MOVOU", _X0, jit.Sib(_ST, _AX, 1, 56)) // MOVOU X0, 56(ST)(AX)
}

func (self *Assembler) _asm_OP_recurse(p *ir.Instr) {
	self.prep_buffer_AX() // MOVE {buf}, (SP)
	vt, pv := p.Vp()
	self.Emit("MOVQ", jit.Type(vt), _BX) // MOVQ $(type(p.Vt())), BX

	/* check for indirection */
	if !rt.UnpackType(vt).Indirect() {
		self.Emit("MOVQ", _SP_p, _CX) // MOVQ SP.p, CX
	} else {
		self.Emit("MOVQ", _SP_p, _VAR_vp) // MOVQ SP.p, VAR.vp
		self.Emit("LEAQ", _VAR_vp, _CX)   // LEAQ VAR.vp, CX
	}

	/* call the encoder */
	self.Emit("MOVQ", _ST, _DI)     // MOVQ  ST, DI
	self.Emit("MOVQ", _ARG_fv, _SI) // MOVQ  $fv, SI
	if pv {
		self.Emit("BTSQ", jit.Imm(alg.BitPointerValue), _SI) // BTSQ $1, SI
	}

	self.call_encoder(_F_encodeTypedPointer) // CALL  encodeTypedPointer
	self.Emit("TESTQ", _ET, _ET)             // TESTQ ET, ET
	self.Sjmp("JNZ", _LB_error)              // JNZ   _error
	self.load_buffer_AX()
}

func (self *Assembler) _asm_OP_is_nil(p *ir.Instr) {
	self.Emit("CMPQ", jit.Ptr(_SP_p, 0), jit.Imm(0)) // CMPQ (SP.p), $0
	self.Xjmp("JE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_is_nil_p1(p *ir.Instr) {
	self.Emit("CMPQ", jit.Ptr(_SP_p, 8), jit.Imm(0)) // CMPQ 8(SP.p), $0
	self.Xjmp("JE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_is_zero_1(p *ir.Instr) {
	self.Emit("CMPB", jit.Ptr(_SP_p, 0), jit.Imm(0)) // CMPB (SP.p), $0
	self.Xjmp("JE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_is_zero_2(p *ir.Instr) {
	self.Emit("CMPW", jit.Ptr(_SP_p, 0), jit.Imm(0)) // CMPW (SP.p), $0
	self.Xjmp("JE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_is_zero_4(p *ir.Instr) {
	self.Emit("CMPL", jit.Ptr(_SP_p, 0), jit.Imm(0)) // CMPL (SP.p), $0
	self.Xjmp("JE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_is_zero_8(p *ir.Instr) {
	self.Emit("CMPQ", jit.Ptr(_SP_p, 0), jit.Imm(0)) // CMPQ (SP.p), $0
	self.Xjmp("JE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_is_zero_map(p *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _AX)      // MOVQ  (SP.p), AX
	self.Emit("TESTQ", _AX, _AX)                   // TESTQ AX, AX
	self.Xjmp("JZ", p.Vi())                        // JZ    p.Vi()
	self.Emit("CMPQ", jit.Ptr(_AX, 0), jit.Imm(0)) // CMPQ  (AX), $0
	self.Xjmp("JE", p.Vi())                        // JE    p.Vi()
}

var (
	_F_is_zero = jit.Func(prim.IsZero)
	_T_reflect_Type = rt.UnpackIface(reflect.Type(nil))
)

func (self *Assembler) _asm_OP_is_zero(p *ir.Instr) {
	fv := p.VField()
	self.Emit("MOVQ", _SP_p, _AX) // ptr
	self.Emit("MOVQ", jit.ImmPtr(unsafe.Pointer(fv)), _BX) // fv
	self.call_go(_F_is_zero) // CALL  $fn
	self.Emit("CMPB", _AX, jit.Imm(0)) // CMPB (SP.p), $0
	self.Xjmp("JNE", p.Vi())                          // JE   p.Vi()
}

func (self *Assembler) _asm_OP_goto(p *ir.Instr) {
	self.Xjmp("JMP", p.Vi())
}

func (self *Assembler) _asm_OP_map_iter(p *ir.Instr) {
	self.Emit("MOVQ", jit.Type(p.Vt()), _AX)  // MOVQ    $p.Vt(), AX
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _BX) // MOVQ    (SP.p), BX
	self.Emit("MOVQ", _ARG_fv, _CX)           // MOVQ    fv, CX
	self.call_go(_F_iteratorStart)            // CALL_GO iteratorStart
	self.Emit("MOVQ", _AX, _SP_q)             // MOVQ    AX, SP.q
	self.Emit("MOVQ", _BX, _ET)               // MOVQ    32(SP), ET
	self.Emit("MOVQ", _CX, _EP)               // MOVQ    40(SP), EP
	self.Emit("TESTQ", _ET, _ET)              // TESTQ   ET, ET
	self.Sjmp("JNZ", _LB_error)               // JNZ     _error
}

func (self *Assembler) _asm_OP_map_stop(_ *ir.Instr) {
	self.Emit("MOVQ", _SP_q, _AX)   // MOVQ    SP.q, AX
	self.call_go(_F_iteratorStop)   // CALL_GO iteratorStop
	self.Emit("XORL", _SP_q, _SP_q) // XORL    SP.q, SP.q
}

func (self *Assembler) _asm_OP_map_check_key(p *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_q, 0), _SP_p) // MOVQ    (SP.q), SP.p
	self.Emit("TESTQ", _SP_p, _SP_p)            // TESTQ   SP.p, SP.p
	self.Xjmp("JZ", p.Vi())                     // JNZ     p.Vi()
}

func (self *Assembler) _asm_OP_map_write_key(p *ir.Instr) {
	self.Emit("BTQ", jit.Imm(alg.BitSortMapKeys), _ARG_fv) // BTQ ${SortMapKeys}, fv
	self.Sjmp("JNC", "_unordered_key_{n}")             // JNC _unordered_key_{n}
	self.encode_string(false)                          // STR $false
	self.Xjmp("JMP", p.Vi())                           // JMP ${p.Vi()}
	self.Link("_unordered_key_{n}")                    // _unordered_key_{n}:
}

func (self *Assembler) _asm_OP_map_value_next(_ *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_q, 8), _SP_p) // MOVQ    8(SP.q), SP.p
	self.Emit("MOVQ", _SP_q, _AX)               // MOVQ    SP.q, AX
	self.call_go(_F_iteratorNext)               // CALL_GO iteratorNext
}

func (self *Assembler) _asm_OP_slice_len(_ *ir.Instr) {
	self.Emit("MOVQ", jit.Ptr(_SP_p, 8), _SP_x)  // MOVQ  8(SP.p), SP.x
	self.Emit("MOVQ", jit.Ptr(_SP_p, 0), _SP_p)  // MOVQ  (SP.p), SP.p
	self.Emit("ORQ", jit.Imm(1<<_S_init), _SP_f) // ORQ   $(1<<_S_init), SP.f
}

func (self *Assembler) _asm_OP_slice_next(p *ir.Instr) {
	self.Emit("TESTQ", _SP_x, _SP_x)                        // TESTQ   SP.x, SP.x
	self.Xjmp("JZ", p.Vi())                                 // JZ      p.Vi()
	self.Emit("SUBQ", jit.Imm(1), _SP_x)                    // SUBQ    $1, SP.x
	self.Emit("BTRQ", jit.Imm(_S_init), _SP_f)              // BTRQ    $_S_init, SP.f
	self.Emit("LEAQ", jit.Ptr(_SP_p, int64(p.Vlen())), _AX) // LEAQ    $(p.vlen())(SP.p), AX
	self.Emit("CMOVQCC", _AX, _SP_p)                        // CMOVQNC AX, SP.p
}

func (self *Assembler) _asm_OP_marshal(p *ir.Instr) {
	self.call_marshaler(_F_encodeJsonMarshaler, _T_json_Marshaler, p.Vt())
}

func (self *Assembler) _asm_OP_marshal_p(p *ir.Instr) {
	if p.Vk() != reflect.Ptr {
		panic("marshal_p: invalid type")
	} else {
		self.call_marshaler_v(_F_encodeJsonMarshaler, _T_json_Marshaler, p.Vt(), false)
	}
}

func (self *Assembler) _asm_OP_marshal_text(p *ir.Instr) {
	self.call_marshaler(_F_encodeTextMarshaler, _T_encoding_TextMarshaler, p.Vt())
}

func (self *Assembler) _asm_OP_marshal_text_p(p *ir.Instr) {
	if p.Vk() != reflect.Ptr {
		panic("marshal_text_p: invalid type")
	} else {
		self.call_marshaler_v(_F_encodeTextMarshaler, _T_encoding_TextMarshaler, p.Vt(), false)
	}
}

func (self *Assembler) _asm_OP_cond_set(_ *ir.Instr) {
	self.Emit("ORQ", jit.Imm(1<<_S_cond), _SP_f) // ORQ $(1<<_S_cond), SP.f
}

func (self *Assembler) _asm_OP_cond_testc(p *ir.Instr) {
	self.Emit("BTRQ", jit.Imm(_S_cond), _SP_f) // BTRQ $_S_cond, SP.f
	self.Xjmp("JC", p.Vi())
}

var _F_error_unsupported = jit.Func(vars.Error_unsuppoted)

func (self *Assembler) _asm_OP_unsupported(i *ir.Instr) {
	typ := int64(uintptr(unsafe.Pointer(i.GoType())))
	self.Emit("MOVQ", jit.Imm(typ), _AX)
	self.call_go(_F_error_unsupported)
	self.Sjmp("JMP", _LB_error)
}

func (self *Assembler) print_gc(i int, p1 *ir.Instr, p2 *ir.Instr) {
	self.Emit("MOVQ", jit.Imm(int64(p2.Op())), _CX) // MOVQ $(p2.Op()), AX
	self.Emit("MOVQ", jit.Imm(int64(p1.Op())), _BX) // MOVQ $(p1.Op()), BX
	self.Emit("MOVQ", jit.Imm(int64(i)), _AX)       // MOVQ $(i), CX
	self.call_go(_F_println)
}
