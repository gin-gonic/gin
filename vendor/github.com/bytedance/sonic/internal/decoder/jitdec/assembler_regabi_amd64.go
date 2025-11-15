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

package jitdec

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"github.com/bytedance/sonic/internal/caching"
	"github.com/bytedance/sonic/internal/jit"
	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/twitchyliquid64/golang-asm/obj"
)

/** Register Allocations
 *
 *  State Registers:
 *
 *      %r13 : stack base
 *      %r10 : input pointer
 *      %r12 : input length
 *      %r11 : input cursor
 *      %r15 : value pointer
 *
 *  Error Registers:
 *
 *      %rax : error type register
 *      %rbx : error pointer register
 */

/** Function Prototype & Stack Map
 *
 *  func (s string, ic int, vp unsafe.Pointer, sb *_Stack, fv uint64, sv string) (rc int, err error)
 *
 *  s.buf  :   (FP)
 *  s.len  :  8(FP)
 *  ic     : 16(FP)
 *  vp     : 24(FP)
 *  sb     : 32(FP)
 *  fv     : 40(FP)
 *  sv     : 56(FP)
 *  err.vt : 72(FP)
 *  err.vp : 80(FP)
 */

const (
    _FP_args   = 72     // 72 bytes to pass and spill register arguments
    _FP_fargs  = 80     // 80 bytes for passing arguments to other Go functions
    _FP_saves  = 48     // 48 bytes for saving the registers before CALL instructions
    _FP_locals = 144    // 144 bytes for local variables
)

const (
    _FP_offs = _FP_fargs + _FP_saves + _FP_locals
    _FP_size = _FP_offs + 8     // 8 bytes for the parent frame pointer
    _FP_base = _FP_size + 8     // 8 bytes for the return address
)

const (
    _IM_null = 0x6c6c756e   // 'null'
    _IM_true = 0x65757274   // 'true'
    _IM_alse = 0x65736c61   // 'alse' ('false' without the 'f')
)

const (
    _BM_space = (1 << ' ') | (1 << '\t') | (1 << '\r') | (1 << '\n')
)

const (
    _MODE_JSON = 1 << 3 // base64 mode
)

const (
    _LB_error           = "_error"
    _LB_im_error        = "_im_error"
    _LB_eof_error       = "_eof_error"
    _LB_type_error      = "_type_error"
    _LB_field_error     = "_field_error"
    _LB_range_error     = "_range_error"
    _LB_stack_error     = "_stack_error"
    _LB_base64_error    = "_base64_error"
    _LB_unquote_error   = "_unquote_error"
    _LB_parsing_error   = "_parsing_error"
    _LB_parsing_error_v = "_parsing_error_v"
    _LB_mismatch_error   = "_mismatch_error"
)

const (
    _LB_char_0_error  = "_char_0_error"
    _LB_char_1_error  = "_char_1_error"
    _LB_char_2_error  = "_char_2_error"
    _LB_char_3_error  = "_char_3_error"
    _LB_char_4_error  = "_char_4_error"
    _LB_char_m2_error = "_char_m2_error"
    _LB_char_m3_error = "_char_m3_error"
)

const (
    _LB_skip_one = "_skip_one"
    _LB_skip_key_value = "_skip_key_value"
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
    _X0 = jit.Reg("X0")
    _X1 = jit.Reg("X1")
    _X15 = jit.Reg("X15")
)

var (
    _IP = jit.Reg("R10")  // saved on BP when callc
    _IC = jit.Reg("R11")  // saved on BX when call_c
    _IL = jit.Reg("R12")
    _ST = jit.Reg("R13")
    _VP = jit.Reg("R15")
)

var (
    _DF = jit.Reg("AX")    // reuse AX in generic decoder for flags
    _ET = jit.Reg("AX")
    _EP = jit.Reg("BX")
)



var (
    _ARG_s  = _ARG_sp
    _ARG_sp = jit.Ptr(_SP, _FP_base + 0)
    _ARG_sl = jit.Ptr(_SP, _FP_base + 8)
    _ARG_ic = jit.Ptr(_SP, _FP_base + 16)
    _ARG_vp = jit.Ptr(_SP, _FP_base + 24)
    _ARG_sb = jit.Ptr(_SP, _FP_base + 32)
    _ARG_fv = jit.Ptr(_SP, _FP_base + 40)
)

var (
    _ARG_sv   = _ARG_sv_p
    _ARG_sv_p = jit.Ptr(_SP, _FP_base + 48)
    _ARG_sv_n = jit.Ptr(_SP, _FP_base + 56)
    _ARG_vk   = jit.Ptr(_SP, _FP_base + 64)
)

var (
    _VAR_st = _VAR_st_Vt
    _VAR_sr = jit.Ptr(_SP, _FP_fargs + _FP_saves)
)

var (
    _VAR_st_Vt = jit.Ptr(_SP, _FP_fargs + _FP_saves + 0)
    _VAR_st_Dv = jit.Ptr(_SP, _FP_fargs + _FP_saves + 8)
    _VAR_st_Iv = jit.Ptr(_SP, _FP_fargs + _FP_saves + 16)
    _VAR_st_Ep = jit.Ptr(_SP, _FP_fargs + _FP_saves + 24)
    _VAR_st_Db = jit.Ptr(_SP, _FP_fargs + _FP_saves + 32)
    _VAR_st_Dc = jit.Ptr(_SP, _FP_fargs + _FP_saves + 40)
)

var (
    _VAR_ss_AX = jit.Ptr(_SP, _FP_fargs + _FP_saves + 48)
    _VAR_ss_CX = jit.Ptr(_SP, _FP_fargs + _FP_saves + 56)
    _VAR_ss_SI = jit.Ptr(_SP, _FP_fargs + _FP_saves + 64)
    _VAR_ss_R8 = jit.Ptr(_SP, _FP_fargs + _FP_saves + 72)
    _VAR_ss_R9 = jit.Ptr(_SP, _FP_fargs + _FP_saves + 80)
)

var (
    _VAR_bs_p = jit.Ptr(_SP, _FP_fargs + _FP_saves + 88)
    _VAR_bs_n = jit.Ptr(_SP, _FP_fargs + _FP_saves + 96)
    _VAR_bs_LR = jit.Ptr(_SP, _FP_fargs + _FP_saves + 104)
)

var _VAR_fl = jit.Ptr(_SP, _FP_fargs + _FP_saves + 112)

var (
    _VAR_et = jit.Ptr(_SP, _FP_fargs + _FP_saves + 120) // save mismatched type
    _VAR_pc = jit.Ptr(_SP, _FP_fargs + _FP_saves + 128) // save skip return pc
    _VAR_ic = jit.Ptr(_SP, _FP_fargs + _FP_saves + 136) // save mismatched position
)

type _Assembler struct {
    jit.BaseAssembler
    p _Program
    name string
}

func newAssembler(p _Program) *_Assembler {
    return new(_Assembler).Init(p)
}

/** Assembler Interface **/

func (self *_Assembler) Load() _Decoder {
    return ptodec(self.BaseAssembler.Load("decode_"+self.name, _FP_size, _FP_args, argPtrs, localPtrs))
}

func (self *_Assembler) Init(p _Program) *_Assembler {
    self.p = p
    self.BaseAssembler.Init(self.compile)
    return self
}

func (self *_Assembler) compile() {
    self.prologue()
    self.instrs()
    self.epilogue()
    self.copy_string()
    self.escape_string()
    self.escape_string_twice()
    self.skip_one()
    self.skip_key_value()
    self.type_error()
    self.mismatch_error()
    self.field_error()
    self.range_error()
    self.stack_error()
    self.base64_error()
    self.parsing_error()
}

/** Assembler Stages **/

var _OpFuncTab = [256]func(*_Assembler, *_Instr) {
    _OP_any              : (*_Assembler)._asm_OP_any,
    _OP_dyn              : (*_Assembler)._asm_OP_dyn,
    _OP_str              : (*_Assembler)._asm_OP_str,
    _OP_bin              : (*_Assembler)._asm_OP_bin,
    _OP_bool             : (*_Assembler)._asm_OP_bool,
    _OP_num              : (*_Assembler)._asm_OP_num,
    _OP_i8               : (*_Assembler)._asm_OP_i8,
    _OP_i16              : (*_Assembler)._asm_OP_i16,
    _OP_i32              : (*_Assembler)._asm_OP_i32,
    _OP_i64              : (*_Assembler)._asm_OP_i64,
    _OP_u8               : (*_Assembler)._asm_OP_u8,
    _OP_u16              : (*_Assembler)._asm_OP_u16,
    _OP_u32              : (*_Assembler)._asm_OP_u32,
    _OP_u64              : (*_Assembler)._asm_OP_u64,
    _OP_f32              : (*_Assembler)._asm_OP_f32,
    _OP_f64              : (*_Assembler)._asm_OP_f64,
    _OP_unquote          : (*_Assembler)._asm_OP_unquote,
    _OP_nil_1            : (*_Assembler)._asm_OP_nil_1,
    _OP_nil_2            : (*_Assembler)._asm_OP_nil_2,
    _OP_nil_3            : (*_Assembler)._asm_OP_nil_3,
    _OP_empty_bytes      : (*_Assembler)._asm_OP_empty_bytes,
    _OP_deref            : (*_Assembler)._asm_OP_deref,
    _OP_index            : (*_Assembler)._asm_OP_index,
    _OP_is_null          : (*_Assembler)._asm_OP_is_null,
    _OP_is_null_quote    : (*_Assembler)._asm_OP_is_null_quote,
    _OP_map_init         : (*_Assembler)._asm_OP_map_init,
    _OP_map_key_i8       : (*_Assembler)._asm_OP_map_key_i8,
    _OP_map_key_i16      : (*_Assembler)._asm_OP_map_key_i16,
    _OP_map_key_i32      : (*_Assembler)._asm_OP_map_key_i32,
    _OP_map_key_i64      : (*_Assembler)._asm_OP_map_key_i64,
    _OP_map_key_u8       : (*_Assembler)._asm_OP_map_key_u8,
    _OP_map_key_u16      : (*_Assembler)._asm_OP_map_key_u16,
    _OP_map_key_u32      : (*_Assembler)._asm_OP_map_key_u32,
    _OP_map_key_u64      : (*_Assembler)._asm_OP_map_key_u64,
    _OP_map_key_f32      : (*_Assembler)._asm_OP_map_key_f32,
    _OP_map_key_f64      : (*_Assembler)._asm_OP_map_key_f64,
    _OP_map_key_str      : (*_Assembler)._asm_OP_map_key_str,
    _OP_map_key_utext    : (*_Assembler)._asm_OP_map_key_utext,
    _OP_map_key_utext_p  : (*_Assembler)._asm_OP_map_key_utext_p,
    _OP_array_skip       : (*_Assembler)._asm_OP_array_skip,
    _OP_array_clear      : (*_Assembler)._asm_OP_array_clear,
    _OP_array_clear_p    : (*_Assembler)._asm_OP_array_clear_p,
    _OP_slice_init       : (*_Assembler)._asm_OP_slice_init,
    _OP_slice_append     : (*_Assembler)._asm_OP_slice_append,
    _OP_object_next      : (*_Assembler)._asm_OP_object_next,
    _OP_struct_field     : (*_Assembler)._asm_OP_struct_field,
    _OP_unmarshal        : (*_Assembler)._asm_OP_unmarshal,
    _OP_unmarshal_p      : (*_Assembler)._asm_OP_unmarshal_p,
    _OP_unmarshal_text   : (*_Assembler)._asm_OP_unmarshal_text,
    _OP_unmarshal_text_p : (*_Assembler)._asm_OP_unmarshal_text_p,
    _OP_lspace           : (*_Assembler)._asm_OP_lspace,
    _OP_match_char       : (*_Assembler)._asm_OP_match_char,
    _OP_check_char       : (*_Assembler)._asm_OP_check_char,
    _OP_load             : (*_Assembler)._asm_OP_load,
    _OP_save             : (*_Assembler)._asm_OP_save,
    _OP_drop             : (*_Assembler)._asm_OP_drop,
    _OP_drop_2           : (*_Assembler)._asm_OP_drop_2,
    _OP_recurse          : (*_Assembler)._asm_OP_recurse,
    _OP_goto             : (*_Assembler)._asm_OP_goto,
    _OP_switch           : (*_Assembler)._asm_OP_switch,
    _OP_check_char_0     : (*_Assembler)._asm_OP_check_char_0,
    _OP_dismatch_err     : (*_Assembler)._asm_OP_dismatch_err,
    _OP_go_skip          : (*_Assembler)._asm_OP_go_skip,
    _OP_skip_emtpy       : (*_Assembler)._asm_OP_skip_empty,
    _OP_add              : (*_Assembler)._asm_OP_add,
    _OP_check_empty      : (*_Assembler)._asm_OP_check_empty,
    _OP_unsupported      : (*_Assembler)._asm_OP_unsupported,
    _OP_debug            : (*_Assembler)._asm_OP_debug,
}

func (self *_Assembler) _asm_OP_debug(_ *_Instr) {
    self.Byte(0xcc)
}

func (self *_Assembler) instr(v *_Instr) {
    if fn := _OpFuncTab[v.op()]; fn != nil {
        fn(self, v)
    } else {
        panic(fmt.Sprintf("invalid opcode: %d", v.op()))
    }
}

func (self *_Assembler) instrs() {
    for i, v := range self.p {
        self.Mark(i)
        self.instr(&v)
        self.debug_instr(i, &v)
    }
}

func (self *_Assembler) epilogue() {
    self.Mark(len(self.p))
    self.Emit("XORL", _EP, _EP)                     // XORL EP, EP
    self.Emit("MOVQ", _VAR_et, _ET)                 // MOVQ VAR_et, ET
    self.Emit("TESTQ", _ET, _ET)                    // TESTQ ET, ET
    self.Sjmp("JNZ", _LB_mismatch_error)            // JNZ _LB_mismatch_error
    self.Link(_LB_error)                            // _error:
    self.Emit("MOVQ", _EP, _CX)                     // MOVQ BX, CX
    self.Emit("MOVQ", _ET, _BX)                     // MOVQ AX, BX
    self.Emit("MOVQ", _IC, _AX)                     // MOVQ IC, AX
    self.Emit("MOVQ", jit.Imm(0), _ARG_sp)          // MOVQ $0, sv.p<>+48(FP)
    self.Emit("MOVQ", jit.Imm(0), _ARG_vp)          // MOVQ $0, sv.p<>+48(FP)
    self.Emit("MOVQ", jit.Imm(0), _ARG_sv_p)        // MOVQ $0, sv.p<>+48(FP)
    self.Emit("MOVQ", jit.Imm(0), _ARG_vk)          // MOVQ $0, vk<>+64(FP)
    self.Emit("MOVQ", jit.Ptr(_SP, _FP_offs), _BP)  // MOVQ _FP_offs(SP), BP
    self.Emit("ADDQ", jit.Imm(_FP_size), _SP)       // ADDQ $_FP_size, SP
    self.Emit("RET")                                // RET
}

func (self *_Assembler) prologue() {
    self.Emit("SUBQ", jit.Imm(_FP_size), _SP)       // SUBQ $_FP_size, SP
    self.Emit("MOVQ", _BP, jit.Ptr(_SP, _FP_offs))  // MOVQ BP, _FP_offs(SP)
    self.Emit("LEAQ", jit.Ptr(_SP, _FP_offs), _BP)  // LEAQ _FP_offs(SP), BP
    self.Emit("MOVQ", _AX, _ARG_sp)                 // MOVQ AX, s.p<>+0(FP)
    self.Emit("MOVQ", _AX, _IP)                     // MOVQ AX, IP
    self.Emit("MOVQ", _BX, _ARG_sl)                 // MOVQ BX, s.l<>+8(FP)
    self.Emit("MOVQ", _BX, _IL)                     // MOVQ BX, IL
    self.Emit("MOVQ", _CX, _ARG_ic)                 // MOVQ CX, ic<>+16(FP)
    self.Emit("MOVQ", _CX, _IC)                     // MOVQ CX, IC
    self.Emit("MOVQ", _DI, _ARG_vp)                 // MOVQ DI, vp<>+24(FP)
    self.Emit("MOVQ", _DI, _VP)                     // MOVQ DI, VP
    self.Emit("MOVQ", _SI, _ARG_sb)                 // MOVQ SI, sb<>+32(FP)
    self.Emit("MOVQ", _SI, _ST)                     // MOVQ SI, ST
    self.Emit("MOVQ", _R8, _ARG_fv)                 // MOVQ R8, fv<>+40(FP)
    self.Emit("MOVQ", jit.Imm(0), _ARG_sv_p)        // MOVQ $0, sv.p<>+48(FP)
    self.Emit("MOVQ", jit.Imm(0), _ARG_sv_n)        // MOVQ $0, sv.n<>+56(FP)
    self.Emit("MOVQ", jit.Imm(0), _ARG_vk)          // MOVQ $0, vk<>+64(FP)
    self.Emit("MOVQ", jit.Imm(0), _VAR_et)          // MOVQ $0, et<>+120(FP)
    // initialize digital buffer first
    self.Emit("MOVQ", jit.Imm(_MaxDigitNums), _VAR_st_Dc)    // MOVQ $_MaxDigitNums, ss.Dcap
    self.Emit("LEAQ", jit.Ptr(_ST, _DbufOffset), _AX)        // LEAQ _DbufOffset(ST), AX
    self.Emit("MOVQ", _AX, _VAR_st_Db)                       // MOVQ AX, ss.Dbuf
}

/** Function Calling Helpers **/

var (
    _REG_go = []obj.Addr { _ST, _VP, _IP, _IL, _IC }
    _REG_rt = []obj.Addr { _ST, _VP, _IP, _IL, _IC }
)

func (self *_Assembler) save(r ...obj.Addr) {
    for i, v := range r {
        if i > _FP_saves / 8 - 1 {
            panic("too many registers to save")
        } else {
            self.Emit("MOVQ", v, jit.Ptr(_SP, _FP_fargs + int64(i) * 8))
        }
    }
}

func (self *_Assembler) load(r ...obj.Addr) {
    for i, v := range r {
        if i > _FP_saves / 8 - 1 {
            panic("too many registers to load")
        } else {
            self.Emit("MOVQ", jit.Ptr(_SP, _FP_fargs + int64(i) * 8), v)
        }
    }
}

func (self *_Assembler) call(fn obj.Addr) {
    self.Emit("MOVQ", fn, _R9)  // MOVQ ${fn}, R11
    self.Rjmp("CALL", _R9)      // CALL R11
}

func (self *_Assembler) call_go(fn obj.Addr) {
    self.save(_REG_go...)   // SAVE $REG_go
    self.call(fn)
    self.load(_REG_go...)   // LOAD $REG_go
}

func (self *_Assembler) callc(fn obj.Addr) {
    self.save(_IP)
    self.call(fn)
    self.Emit("XORPS", _X15, _X15)
    self.load(_IP)
}

func (self *_Assembler) call_c(fn obj.Addr) {
    self.Emit("XCHGQ", _IC, _BX)
    self.callc(fn)
    self.Emit("XCHGQ", _IC, _BX)
}

func (self *_Assembler) call_sf(fn obj.Addr) {
    self.Emit("LEAQ", _ARG_s, _DI)                      // LEAQ s<>+0(FP), DI
    self.Emit("MOVQ", _IC, _ARG_ic)                     // MOVQ IC, ic<>+16(FP)
    self.Emit("LEAQ", _ARG_ic, _SI)                     // LEAQ ic<>+16(FP), SI
    self.Emit("LEAQ", jit.Ptr(_ST, _FsmOffset), _DX)    // LEAQ _FsmOffset(ST), DX
    self.Emit("MOVQ", _ARG_fv, _CX)
    self.callc(fn)
    self.Emit("MOVQ", _ARG_ic, _IC)                     // MOVQ ic<>+16(FP), IC
}

func (self *_Assembler) call_vf(fn obj.Addr) {
    self.Emit("LEAQ", _ARG_s, _DI)      // LEAQ s<>+0(FP), DI
    self.Emit("MOVQ", _IC, _ARG_ic)     // MOVQ IC, ic<>+16(FP)
    self.Emit("LEAQ", _ARG_ic, _SI)     // LEAQ ic<>+16(FP), SI
    self.Emit("LEAQ", _VAR_st, _DX)     // LEAQ st, DX
    self.callc(fn)
    self.Emit("MOVQ", _ARG_ic, _IC)     // MOVQ ic<>+16(FP), IC
}

/** Assembler Error Handlers **/

var (
    _F_convT64        = jit.Func(rt.ConvT64)
    _F_error_wrap     = jit.Func(error_wrap)
    _F_error_type     = jit.Func(error_type)
    _F_error_field    = jit.Func(error_field)
    _F_error_value    = jit.Func(error_value)
    _F_error_mismatch = jit.Func(error_mismatch)
)

var (
    _I_int8    , _T_int8    = rtype(reflect.TypeOf(int8(0)))
    _I_int16   , _T_int16   = rtype(reflect.TypeOf(int16(0)))
    _I_int32   , _T_int32   = rtype(reflect.TypeOf(int32(0)))
    _I_uint8   , _T_uint8   = rtype(reflect.TypeOf(uint8(0)))
    _I_uint16  , _T_uint16  = rtype(reflect.TypeOf(uint16(0)))
    _I_uint32  , _T_uint32  = rtype(reflect.TypeOf(uint32(0)))
    _I_float32 , _T_float32 = rtype(reflect.TypeOf(float32(0)))
)

var (
    _T_error                    = rt.UnpackType(errorType)
    _I_base64_CorruptInputError = jit.Itab(_T_error, base64CorruptInputError)
)

var (
    _V_stackOverflow              = jit.Imm(int64(uintptr(unsafe.Pointer(&stackOverflow))))
    _I_json_UnsupportedValueError = jit.Itab(_T_error, reflect.TypeOf(new(json.UnsupportedValueError)))
    _I_json_MismatchTypeError     = jit.Itab(_T_error, reflect.TypeOf(new(MismatchTypeError)))
    _I_json_MismatchQuotedError   = jit.Itab(_T_error, reflect.TypeOf(new(MismatchQuotedError)))
)

func (self *_Assembler) type_error() {
    self.Link(_LB_type_error)                   // _type_error:
    self.call_go(_F_error_type)                 // CALL_GO error_type
    self.Sjmp("JMP" , _LB_error)                // JMP     _error
}

func (self *_Assembler) mismatch_error() {
    self.Link(_LB_mismatch_error)                     // _type_error:
    self.Emit("MOVQ", _VAR_et, _ET)                   // MOVQ _VAR_et, ET
    self.Emit("MOVQ", _I_json_MismatchTypeError, _CX) // MOVQ _I_json_MismatchType, CX
    self.Emit("CMPQ", _ET, _CX)                       // CMPQ ET, CX
    self.Emit("MOVQ", jit.Ptr(_ST, _EpOffset), _EP)   // MOVQ stack.Ep, EP
    self.Sjmp("JE"  , _LB_error)                      // JE _LB_error
    self.Emit("MOVQ", _ARG_sp, _AX)
    self.Emit("MOVQ", _ARG_sl, _BX)
    self.Emit("MOVQ", _VAR_ic, _CX)
    self.Emit("MOVQ", _VAR_et, _DI)
    self.call_go(_F_error_mismatch)             // CALL_GO error_type
    self.Sjmp("JMP" , _LB_error)                // JMP     _error
}

func (self *_Assembler) field_error() {
    self.Link(_LB_field_error)                  // _field_error:
    self.Emit("MOVQ", _ARG_sv_p, _AX)           // MOVQ   sv.p, AX
    self.Emit("MOVQ", _ARG_sv_n, _BX)           // MOVQ   sv.n, BX
    self.call_go(_F_error_field)                // CALL_GO error_field
    self.Sjmp("JMP" , _LB_error)                // JMP     _error
}

func (self *_Assembler) range_error() {
    self.Link(_LB_range_error)                  // _range_error:
    self.Emit("MOVQ", _ET, _CX)                 // MOVQ    ET, CX
    self.slice_from(_VAR_st_Ep, 0)              // SLICE   st.Ep, $0
    self.Emit("MOVQ", _DI, _AX)                 // MOVQ    DI, AX
    self.Emit("MOVQ", _EP, _DI)                 // MOVQ    EP, DI
    self.Emit("MOVQ", _SI, _BX)                 // MOVQ    SI, BX
    self.call_go(_F_error_value)                // CALL_GO error_value
    self.Sjmp("JMP" , _LB_error)                // JMP     _error
}

func (self *_Assembler) stack_error() {
    self.Link(_LB_stack_error)                              // _stack_error:
    self.Emit("MOVQ", _V_stackOverflow, _EP)                // MOVQ ${_V_stackOverflow}, EP
    self.Emit("MOVQ", _I_json_UnsupportedValueError, _ET)   // MOVQ ${_I_json_UnsupportedValueError}, ET
    self.Sjmp("JMP" , _LB_error)                            // JMP  _error
}

func (self *_Assembler) base64_error() {
    self.Link(_LB_base64_error)
    self.Emit("NEGQ", _AX)                                  // NEGQ    AX
    self.Emit("SUBQ", jit.Imm(1), _AX)                      // SUBQ    $1, AX
    self.call_go(_F_convT64)                                // CALL_GO convT64
    self.Emit("MOVQ", _AX, _EP)                             // MOVQ    AX, EP
    self.Emit("MOVQ", _I_base64_CorruptInputError, _ET)     // MOVQ    ${itab(base64.CorruptInputError)}, ET
    self.Sjmp("JMP" , _LB_error)                            // JMP     _error
}

func (self *_Assembler) parsing_error() {
    self.Link(_LB_eof_error)                                            // _eof_error:
    self.Emit("MOVQ" , _IL, _IC)                                        // MOVQ    IL, IC
    self.Emit("MOVL" , jit.Imm(int64(types.ERR_EOF)), _EP)              // MOVL    ${types.ERR_EOF}, EP
    self.Sjmp("JMP"  , _LB_parsing_error)                               // JMP     _parsing_error
    self.Link(_LB_unquote_error)                                        // _unquote_error:
    self.Emit("SUBQ" , _VAR_sr, _SI)                                    // SUBQ    sr, SI
    self.Emit("SUBQ" , _SI, _IC)                                        // SUBQ    IL, IC
    self.Link(_LB_parsing_error_v)                                      // _parsing_error_v:
    self.Emit("MOVQ" , _AX, _EP)                                        // MOVQ    AX, EP
    self.Emit("NEGQ" , _EP)                                             // NEGQ    EP
    self.Sjmp("JMP"  , _LB_parsing_error)                               // JMP     _parsing_error
    self.Link(_LB_char_m3_error)                                        // _char_m3_error:
    self.Emit("SUBQ" , jit.Imm(1), _IC)                                 // SUBQ    $1, IC
    self.Link(_LB_char_m2_error)                                        // _char_m2_error:
    self.Emit("SUBQ" , jit.Imm(2), _IC)                                 // SUBQ    $2, IC
    self.Sjmp("JMP"  , _LB_char_0_error)                                // JMP     _char_0_error
    self.Link(_LB_im_error)                                             // _im_error:
    self.Emit("CMPB" , _CX, jit.Sib(_IP, _IC, 1, 0))                    // CMPB    CX, (IP)(IC)
    self.Sjmp("JNE"  , _LB_char_0_error)                                // JNE     _char_0_error
    self.Emit("SHRL" , jit.Imm(8), _CX)                                 // SHRL    $8, CX
    self.Emit("CMPB" , _CX, jit.Sib(_IP, _IC, 1, 1))                    // CMPB    CX, 1(IP)(IC)
    self.Sjmp("JNE"  , _LB_char_1_error)                                // JNE     _char_1_error
    self.Emit("SHRL" , jit.Imm(8), _CX)                                 // SHRL    $8, CX
    self.Emit("CMPB" , _CX, jit.Sib(_IP, _IC, 1, 2))                    // CMPB    CX, 2(IP)(IC)
    self.Sjmp("JNE"  , _LB_char_2_error)                                // JNE     _char_2_error
    self.Sjmp("JMP"  , _LB_char_3_error)                                // JNE     _char_3_error
    self.Link(_LB_char_4_error)                                         // _char_4_error:
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 // ADDQ    $1, IC
    self.Link(_LB_char_3_error)                                         // _char_3_error:
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 // ADDQ    $1, IC
    self.Link(_LB_char_2_error)                                         // _char_2_error:
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 // ADDQ    $1, IC
    self.Link(_LB_char_1_error)                                         // _char_1_error:
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 // ADDQ    $1, IC
    self.Link(_LB_char_0_error)                                         // _char_0_error:
    self.Emit("MOVL" , jit.Imm(int64(types.ERR_INVALID_CHAR)), _EP)     // MOVL    ${types.ERR_INVALID_CHAR}, EP
    self.Link(_LB_parsing_error)                                        // _parsing_error:
    self.Emit("MOVQ" , _EP, _DI)                                        // MOVQ    EP, DI
    self.Emit("MOVQ",  _ARG_sp, _AX)                                     // MOVQ  sp, AX
    self.Emit("MOVQ",  _ARG_sl, _BX)                                     // MOVQ  sl, BX
    self.Emit("MOVQ" , _IC, _CX)                                        // MOVQ    IC, CX
    self.call_go(_F_error_wrap)                                         // CALL_GO error_wrap
    self.Sjmp("JMP"  , _LB_error)                                       // JMP     _error
}

func (self *_Assembler) _asm_OP_dismatch_err(p *_Instr) {
    self.Emit("MOVQ", _IC, _VAR_ic)      
    self.Emit("MOVQ", jit.Type(p.vt()), _ET)     
    self.Emit("MOVQ", _ET, _VAR_et)
}

func (self *_Assembler) _asm_OP_go_skip(p *_Instr) {
    self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
    self.Xref(p.vi(), 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)            // JMP     _skip_one
}

var _F_IndexByte = jit.Func(strings.IndexByte)

func (self *_Assembler) _asm_OP_skip_empty(p *_Instr) {
    self.call_sf(_F_skip_one)                   // CALL_SF skip_one
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)     // JS      _parse_error_v
    self.Emit("BTQ", jit.Imm(_F_disable_unknown), _ARG_fv) 
    self.Xjmp("JNC", p.vi())
    self.Emit("LEAQ", jit.Sib(_IC, _AX, 1, 0), _BX)
    self.Emit("MOVQ", _BX, _ARG_sv_n)
    self.Emit("LEAQ", jit.Sib(_IP, _AX, 1, 0), _AX)
    self.Emit("MOVQ", _AX, _ARG_sv_p)
    self.Emit("MOVQ", jit.Imm(':'), _CX)
    self.call_go(_F_IndexByte)
    self.Emit("TESTQ", _AX, _AX)
    // disallow unknown field
    self.Sjmp("JNS", _LB_field_error)
}

func (self *_Assembler) skip_one() {
    self.Link(_LB_skip_one)                     // _skip:
    self.Emit("MOVQ", _VAR_ic, _IC)             // MOVQ    _VAR_ic, IC
    self.call_sf(_F_skip_one)                   // CALL_SF skip_one
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)     // JS      _parse_error_v
    self.Emit("MOVQ" , _VAR_pc, _R9)            // MOVQ    pc, R9
    self.Rjmp("JMP"  , _R9)                     // JMP     (R9)
}

func (self *_Assembler) skip_key_value() {
    self.Link(_LB_skip_key_value)               // _skip:
    // skip the key
    self.Emit("MOVQ", _VAR_ic, _IC)             // MOVQ    _VAR_ic, IC
    self.call_sf(_F_skip_one)                   // CALL_SF skip_one
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)     // JS      _parse_error_v
    // match char ':'
    self.lspace("_global_1")
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(':'))
    self.Sjmp("JNE"  , _LB_parsing_error_v)     // JNE     _parse_error_v
    self.Emit("ADDQ", jit.Imm(1), _IC)          // ADDQ    $1, IC
    self.lspace("_global_2")
    // skip the value
    self.call_sf(_F_skip_one)                   // CALL_SF skip_one
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)     // JS      _parse_error_v
    // jump back to specified address
    self.Emit("MOVQ" , _VAR_pc, _R9)            // MOVQ    pc, R9
    self.Rjmp("JMP"  , _R9)                     // JMP     (R9)
}


/** Memory Management Routines **/

var (
    _T_byte     = jit.Type(byteType)
    _F_mallocgc = jit.Func(rt.Mallocgc)
)

func (self *_Assembler) malloc_AX(nb obj.Addr, ret obj.Addr) {
    self.Emit("MOVQ", nb, _AX)                  // MOVQ    ${nb}, AX
    self.Emit("MOVQ", _T_byte, _BX)             // MOVQ    ${type(byte)}, BX
    self.Emit("XORL", _CX, _CX)                 // XORL    CX, CX
    self.call_go(_F_mallocgc)                   // CALL_GO mallocgc
    self.Emit("MOVQ", _AX, ret)                 // MOVQ    AX, ${ret}
}

func (self *_Assembler) valloc(vt reflect.Type, ret obj.Addr) {
    self.Emit("MOVQ", jit.Imm(int64(vt.Size())), _AX)   // MOVQ    ${vt.Size()}, AX
    self.Emit("MOVQ", jit.Type(vt), _BX)                // MOVQ    ${vt}, BX
    self.Emit("MOVB", jit.Imm(1), _CX)                  // MOVB    $1, CX
    self.call_go(_F_mallocgc)                           // CALL_GO mallocgc
    self.Emit("MOVQ", _AX, ret)                         // MOVQ    AX, ${ret}
}

func (self *_Assembler) valloc_AX(vt reflect.Type) {
    self.Emit("MOVQ", jit.Imm(int64(vt.Size())), _AX)   // MOVQ    ${vt.Size()}, AX
    self.Emit("MOVQ", jit.Type(vt), _BX)                // MOVQ    ${vt}, BX
    self.Emit("MOVB", jit.Imm(1), _CX)                  // MOVB    $1, CX
    self.call_go(_F_mallocgc)                           // CALL_GO mallocgc
}

func (self *_Assembler) vfollow(vt reflect.Type) {
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _AX)    // MOVQ   (VP), AX
    self.Emit("TESTQ", _AX, _AX)                // TESTQ  AX, AX
    self.Sjmp("JNZ"  , "_end_{n}")              // JNZ    _end_{n}
    self.valloc_AX(vt)                          // VALLOC ${vt}, AX
    self.WritePtrAX(1, jit.Ptr(_VP, 0), true)   // MOVQ   AX, (VP)
    self.Link("_end_{n}")                       // _end_{n}:
    self.Emit("MOVQ" , _AX, _VP)                // MOVQ   AX, VP
}

/** Value Parsing Routines **/

var (
    _F_vstring   = jit.Imm(int64(native.S_vstring))
    _F_vnumber   = jit.Imm(int64(native.S_vnumber))
    _F_vsigned   = jit.Imm(int64(native.S_vsigned))
    _F_vunsigned = jit.Imm(int64(native.S_vunsigned))
)

func (self *_Assembler) check_err(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ" , _VAR_st_Vt, _AX)         // MOVQ st.Vt, AX
    self.Emit("TESTQ", _AX, _AX)                // CMPQ AX, ${native.V_STRING}
    // try to skip the value
    if vt != nil {
        self.Sjmp("JNS" , "_check_err_{n}")        // JNE  _parsing_error_v
        self.Emit("MOVQ", jit.Type(vt), _ET)         
        self.Emit("MOVQ", _ET, _VAR_et)
        if pin2 != -1 {
            self.Emit("SUBQ", jit.Imm(1), _BX)
            self.Emit("MOVQ", _BX, _VAR_ic)
            self.Byte(0x4c  , 0x8d, 0x0d)         // LEAQ (PC), R9
            self.Xref(pin2, 4)
            self.Emit("MOVQ", _R9, _VAR_pc)
            self.Sjmp("JMP" , _LB_skip_key_value)
        } else {
            self.Emit("MOVQ", _BX, _VAR_ic)
            self.Byte(0x4c  , 0x8d, 0x0d)         // LEAQ (PC), R9
            self.Sref(pin, 4)
            self.Emit("MOVQ", _R9, _VAR_pc)
            self.Sjmp("JMP" , _LB_skip_one)
        }
        self.Link("_check_err_{n}")
    } else {
        self.Sjmp("JS"   , _LB_parsing_error_v)     // JNE  _parsing_error_v
    }
}

func (self *_Assembler) check_eof(d int64) {
    if d == 1 {
        self.Emit("CMPQ", _IC, _IL)         // CMPQ IC, IL
        self.Sjmp("JAE" , _LB_eof_error)    // JAE  _eof_error
    } else {
        self.Emit("LEAQ", jit.Ptr(_IC, d), _AX)     // LEAQ ${d}(IC), AX
        self.Emit("CMPQ", _AX, _IL)                 // CMPQ AX, IL
        self.Sjmp("JA"  , _LB_eof_error)            // JA   _eof_error
    }
}


func (self *_Assembler) parse_string() {
    self.Emit("MOVQ", _ARG_fv, _CX)
    self.call_vf(_F_vstring)
    self.check_err(nil, "", -1)
}

func (self *_Assembler) parse_number(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ", _IC, _BX)       // save ic when call native func    
    self.call_vf(_F_vnumber)
    self.check_err(vt, pin, pin2)
}

func (self *_Assembler) parse_signed(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ", _IC, _BX)       // save ic when call native func    
    self.call_vf(_F_vsigned)
    self.check_err(vt, pin, pin2)
}

func (self *_Assembler) parse_unsigned(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ", _IC, _BX)       // save ic when call native func    
    self.call_vf(_F_vunsigned)
    self.check_err(vt, pin, pin2)
}

// Pointer: DI, Size: SI, Return: R9  
func (self *_Assembler) copy_string() {
    self.Link("_copy_string")
    self.Emit("MOVQ", _DI, _VAR_bs_p)
    self.Emit("MOVQ", _SI, _VAR_bs_n)
    self.Emit("MOVQ", _R9, _VAR_bs_LR)
    self.malloc_AX(_SI, _ARG_sv_p)                              
    self.Emit("MOVQ", _VAR_bs_p, _BX)
    self.Emit("MOVQ", _VAR_bs_n, _CX)
    self.call_go(_F_memmove)
    self.Emit("MOVQ", _ARG_sv_p, _DI)
    self.Emit("MOVQ", _VAR_bs_n, _SI)
    self.Emit("MOVQ", _VAR_bs_LR, _R9)
    self.Rjmp("JMP", _R9)
}

// Pointer: DI, Size: SI, Return: R9
func (self *_Assembler) escape_string() {
    self.Link("_escape_string")
    self.Emit("MOVQ" , _DI, _VAR_bs_p)
    self.Emit("MOVQ" , _SI, _VAR_bs_n)
    self.Emit("MOVQ" , _R9, _VAR_bs_LR)
    self.malloc_AX(_SI, _DX)                                    // MALLOC SI, DX
    self.Emit("MOVQ" , _DX, _ARG_sv_p)
    self.Emit("MOVQ" , _VAR_bs_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_n, _SI)                                  
    self.Emit("LEAQ" , _VAR_sr, _CX)                            // LEAQ   sr, CX
    self.Emit("XORL" , _R8, _R8)                                // XORL   R8, R8
    self.Emit("BTQ"  , jit.Imm(_F_disable_urc), _ARG_fv)        // BTQ    ${_F_disable_urc}, fv
    self.Emit("SETCC", _R8)                                     // SETCC  R8
    self.Emit("SHLQ" , jit.Imm(types.B_UNICODE_REPLACE), _R8)   // SHLQ   ${types.B_UNICODE_REPLACE}, R8
    self.call_c(_F_unquote)                                       // CALL   unquote
    self.Emit("MOVQ" , _VAR_bs_n, _SI)                                  // MOVQ   ${n}, SI
    self.Emit("ADDQ" , jit.Imm(1), _SI)                         // ADDQ   $1, SI
    self.Emit("TESTQ", _AX, _AX)                                // TESTQ  AX, AX
    self.Sjmp("JS"   , _LB_unquote_error)                       // JS     _unquote_error
    self.Emit("MOVQ" , _AX, _SI)
    self.Emit("MOVQ" , _ARG_sv_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_LR, _R9)
    self.Rjmp("JMP", _R9)
}

func (self *_Assembler) escape_string_twice() {
    self.Link("_escape_string_twice")
    self.Emit("MOVQ" , _DI, _VAR_bs_p)
    self.Emit("MOVQ" , _SI, _VAR_bs_n)
    self.Emit("MOVQ" , _R9, _VAR_bs_LR)
    self.malloc_AX(_SI, _DX)                                        // MALLOC SI, DX
    self.Emit("MOVQ" , _DX, _ARG_sv_p)
    self.Emit("MOVQ" , _VAR_bs_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_n, _SI)        
    self.Emit("LEAQ" , _VAR_sr, _CX)                                // LEAQ   sr, CX
    self.Emit("MOVL" , jit.Imm(types.F_DOUBLE_UNQUOTE), _R8)        // MOVL   ${types.F_DOUBLE_UNQUOTE}, R8
    self.Emit("BTQ"  , jit.Imm(_F_disable_urc), _ARG_fv)            // BTQ    ${_F_disable_urc}, AX
    self.Emit("XORL" , _AX, _AX)                                    // XORL   AX, AX
    self.Emit("SETCC", _AX)                                         // SETCC  AX
    self.Emit("SHLQ" , jit.Imm(types.B_UNICODE_REPLACE), _AX)       // SHLQ   ${types.B_UNICODE_REPLACE}, AX
    self.Emit("ORQ"  , _AX, _R8)                                    // ORQ    AX, R8
    self.call_c(_F_unquote)                                         // CALL   unquote
    self.Emit("MOVQ" , _VAR_bs_n, _SI)                              // MOVQ   ${n}, SI
    self.Emit("ADDQ" , jit.Imm(3), _SI)                             // ADDQ   $3, SI
    self.Emit("TESTQ", _AX, _AX)                                    // TESTQ  AX, AX
    self.Sjmp("JS"   , _LB_unquote_error)                           // JS     _unquote_error
    self.Emit("MOVQ" , _AX, _SI)
    self.Emit("MOVQ" , _ARG_sv_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_LR, _R9)
    self.Rjmp("JMP", _R9)
}

/** Range Checking Routines **/

var (
    _V_max_f32 = jit.Imm(int64(uintptr(unsafe.Pointer(_Vp_max_f32))))
    _V_min_f32 = jit.Imm(int64(uintptr(unsafe.Pointer(_Vp_min_f32))))
)

var (
    _Vp_max_f32 = new(float32)
    _Vp_min_f32 = new(float32)
)

func init() {
    *_Vp_max_f32 = math.MaxFloat32
    *_Vp_min_f32 = -math.MaxFloat32
}

func (self *_Assembler) range_single_X0() {
    self.Emit("CVTSD2SS", _VAR_st_Dv, _X0)              // CVTSD2SS _VAR_st_Dv, X0
    self.Emit("MOVQ"    , _V_max_f32, _CX)              // MOVQ     _max_f32, CX
    self.Emit("MOVQ"    , jit.Gitab(_I_float32), _ET)   // MOVQ     ${itab(float32)}, ET
    self.Emit("MOVQ"    , jit.Gtype(_T_float32), _EP)   // MOVQ     ${type(float32)}, EP
    self.Emit("UCOMISS" , jit.Ptr(_CX, 0), _X0)         // UCOMISS  (CX), X0
    self.Sjmp("JA"      , _LB_range_error)              // JA       _range_error
    self.Emit("MOVQ"    , _V_min_f32, _CX)              // MOVQ     _min_f32, CX
    self.Emit("UCOMISS" , jit.Ptr(_CX, 0), _X0)         // UCOMISS  (CX), X0
    self.Sjmp("JB"      , _LB_range_error)              // JB      _range_error
}

func (self *_Assembler) range_signed_CX(i *rt.GoItab, t *rt.GoType, a int64, b int64) {
    self.Emit("MOVQ", _VAR_st_Iv, _CX)      // MOVQ st.Iv, CX
    self.Emit("MOVQ", jit.Gitab(i), _ET)    // MOVQ ${i}, ET
    self.Emit("MOVQ", jit.Gtype(t), _EP)    // MOVQ ${t}, EP
    self.Emit("CMPQ", _CX, jit.Imm(a))      // CMPQ CX, ${a}
    self.Sjmp("JL"  , _LB_range_error)      // JL   _range_error
    self.Emit("CMPQ", _CX, jit.Imm(b))      // CMPQ CX, ${B}
    self.Sjmp("JG"  , _LB_range_error)      // JG   _range_error
}

func (self *_Assembler) range_unsigned_CX(i *rt.GoItab, t *rt.GoType, v uint64) {
    self.Emit("MOVQ" , _VAR_st_Iv, _CX)         // MOVQ  st.Iv, CX
    self.Emit("MOVQ" , jit.Gitab(i), _ET)       // MOVQ  ${i}, ET
    self.Emit("MOVQ" , jit.Gtype(t), _EP)       // MOVQ  ${t}, EP
    self.Emit("TESTQ", _CX, _CX)                // TESTQ CX, CX
    self.Sjmp("JS"   , _LB_range_error)         // JS    _range_error
    self.Emit("CMPQ" , _CX, jit.Imm(int64(v)))  // CMPQ  CX, ${a}
    self.Sjmp("JA"   , _LB_range_error)         // JA    _range_error
}

func (self *_Assembler) range_uint32_CX(i *rt.GoItab, t *rt.GoType) {
    self.Emit("MOVQ" , _VAR_st_Iv, _CX)         // MOVQ  st.Iv, CX
    self.Emit("MOVQ" , jit.Gitab(i), _ET)       // MOVQ  ${i}, ET
    self.Emit("MOVQ" , jit.Gtype(t), _EP)       // MOVQ  ${t}, EP
    self.Emit("TESTQ", _CX, _CX)                // TESTQ CX, CX
    self.Sjmp("JS"   , _LB_range_error)         // JS    _range_error
    self.Emit("MOVL" , _CX, _DX)                // MOVL  CX, DX
    self.Emit("CMPQ" , _CX, _DX)                // CMPQ  CX, DX
    self.Sjmp("JNE"  , _LB_range_error)         // JNZ   _range_error
}

/** String Manipulating Routines **/

var (
    _F_unquote = jit.Imm(int64(native.S_unquote))
)

func (self *_Assembler) slice_from(p obj.Addr, d int64) {
    self.Emit("MOVQ", p, _SI)   // MOVQ    ${p}, SI
    self.slice_from_r(_SI, d)   // SLICE_R SI, ${d}
}

func (self *_Assembler) slice_from_r(p obj.Addr, d int64) {
    self.Emit("LEAQ", jit.Sib(_IP, p, 1, 0), _DI)   // LEAQ (IP)(${p}), DI
    self.Emit("NEGQ", p)                            // NEGQ ${p}
    self.Emit("LEAQ", jit.Sib(_IC, p, 1, d), _SI)   // LEAQ d(IC)(${p}), SI
}

func (self *_Assembler) unquote_once(p obj.Addr, n obj.Addr, stack bool, copy bool) {
    self.slice_from(_VAR_st_Iv, -1)                             // SLICE  st.Iv, $-1
    self.Emit("CMPQ", _VAR_st_Ep, jit.Imm(-1))                 // CMPQ   st.Ep, $-1
    self.Sjmp("JE"  , "_noescape_{n}")                         // JE     _escape_{n}
    self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
    self.Sref("_unquote_once_write_{n}", 4)
    self.Sjmp("JMP" , "_escape_string")
    self.Link("_noescape_{n}")
    if copy {
        self.Emit("BTQ" , jit.Imm(_F_copy_string), _ARG_fv)    
        self.Sjmp("JNC", "_unquote_once_write_{n}")
        self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
        self.Sref("_unquote_once_write_{n}", 4)
        self.Sjmp("JMP", "_copy_string")
    }
    self.Link("_unquote_once_write_{n}")
    self.Emit("MOVQ", _SI, n)                                  // MOVQ   SI, ${n}
    if stack {
        self.Emit("MOVQ", _DI, p) 
    } else {
        self.WriteRecNotAX(10, _DI, p, false, false)
    }
}

func (self *_Assembler) unquote_twice(p obj.Addr, n obj.Addr, stack bool) {
    self.Emit("CMPQ" , _VAR_st_Ep, jit.Imm(-1))                     // CMPQ   st.Ep, $-1
    self.Sjmp("JE"   , _LB_eof_error)                               // JE     _eof_error
    self.Emit("CMPB" , jit.Sib(_IP, _IC, 1, -3), jit.Imm('\\'))     // CMPB   -3(IP)(IC), $'\\'
    self.Sjmp("JNE"  , _LB_char_m3_error)                           // JNE    _char_m3_error
    self.Emit("CMPB" , jit.Sib(_IP, _IC, 1, -2), jit.Imm('"'))      // CMPB   -2(IP)(IC), $'"'
    self.Sjmp("JNE"  , _LB_char_m2_error)                           // JNE    _char_m2_error
    self.slice_from(_VAR_st_Iv, -3)                                 // SLICE  st.Iv, $-3
    self.Emit("MOVQ" , _SI, _AX)                                    // MOVQ   SI, AX
    self.Emit("ADDQ" , _VAR_st_Iv, _AX)                             // ADDQ   st.Iv, AX
    self.Emit("CMPQ" , _VAR_st_Ep, _AX)                             // CMPQ   st.Ep, AX
    self.Sjmp("JE"   , "_noescape_{n}")                             // JE     _noescape_{n}
    self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
    self.Sref("_unquote_twice_write_{n}", 4)
    self.Sjmp("JMP" , "_escape_string_twice")
    self.Link("_noescape_{n}")                                      // _noescape_{n}:
    self.Emit("BTQ"  , jit.Imm(_F_copy_string), _ARG_fv)    
    self.Sjmp("JNC", "_unquote_twice_write_{n}") 
    self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
    self.Sref("_unquote_twice_write_{n}", 4)
    self.Sjmp("JMP", "_copy_string")
    self.Link("_unquote_twice_write_{n}")
    self.Emit("MOVQ" , _SI, n)                                      // MOVQ   SI, ${n}
    if stack {
        self.Emit("MOVQ", _DI, p) 
    } else {
        self.WriteRecNotAX(12, _DI, p, false, false)
    }
    self.Link("_unquote_twice_end_{n}")
}

/** Memory Clearing Routines **/

var (
    _F_memclrHasPointers    = jit.Func(rt.MemclrHasPointers)
    _F_memclrNoHeapPointers = jit.Func(rt.MemclrNoHeapPointers)
)

func (self *_Assembler) mem_clear_fn(ptrfree bool) {
    if !ptrfree {
        self.call_go(_F_memclrHasPointers)
    } else {
        self.call_go(_F_memclrNoHeapPointers)
    }
}

func (self *_Assembler) mem_clear_rem(size int64, ptrfree bool) {
    self.Emit("MOVQ", jit.Imm(size), _BX)               // MOVQ    ${size}, BX
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             // MOVQ    (ST), AX
    self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 0), _AX)     // MOVQ    (ST)(AX), AX
    self.Emit("SUBQ", _VP, _AX)                         // SUBQ    VP, AX
    self.Emit("ADDQ", _AX, _BX)                         // ADDQ    AX, BX
    self.Emit("MOVQ", _VP, _AX)                         // MOVQ    VP, (SP)
    self.mem_clear_fn(ptrfree)                          // CALL_GO memclr{Has,NoHeap}Pointers
}

/** Map Assigning Routines **/

var (
    _F_mapassign           = jit.Func(rt.Mapassign)
    _F_mapassign_fast32    = jit.Func(rt.Mapassign_fast32)
    _F_mapassign_faststr   = jit.Func(rt.Mapassign_faststr)
    _F_mapassign_fast64ptr = jit.Func(rt.Mapassign_fast64ptr)
)

var (
    _F_decodeJsonUnmarshaler obj.Addr
    _F_decodeJsonUnmarshalerQuoted obj.Addr
    _F_decodeTextUnmarshaler obj.Addr
)

func init() {
    _F_decodeJsonUnmarshaler = jit.Func(decodeJsonUnmarshaler)
    _F_decodeJsonUnmarshalerQuoted = jit.Func(decodeJsonUnmarshalerQuoted)
    _F_decodeTextUnmarshaler = jit.Func(decodeTextUnmarshaler)
}

func (self *_Assembler) mapaccess_ptr(t reflect.Type) {
    if rt.MapType(rt.UnpackType(t)).IndirectElem() {
        self.vfollow(t.Elem())
    }
}

func (self *_Assembler) mapassign_std(t reflect.Type, v obj.Addr) {
    self.Emit("LEAQ", v, _AX)               // LEAQ      ${v}, AX
    self.mapassign_call_from_AX(t, _F_mapassign)    // MAPASSIGN ${t}, mapassign
}

func (self *_Assembler) mapassign_str_fast(t reflect.Type, p obj.Addr, n obj.Addr) {
    self.Emit("MOVQ", jit.Type(t), _AX)         // MOVQ    ${t}, AX
    self.Emit("MOVQ", _VP, _BX)                 // MOVQ    VP, BX
    self.Emit("MOVQ", p, _CX)                   // MOVQ    ${p}, CX
    self.Emit("MOVQ", n, _DI)                   // MOVQ    ${n}, DI
    self.call_go(_F_mapassign_faststr)          // CALL_GO ${fn}
    self.Emit("MOVQ", _AX, _VP)                 // MOVQ    AX, VP
    self.mapaccess_ptr(t)
}

func (self *_Assembler) mapassign_call_from_AX(t reflect.Type, fn obj.Addr) {
    self.Emit("MOVQ", _AX, _CX)
    self.Emit("MOVQ", jit.Type(t), _AX)         // MOVQ    ${t}, AX
    self.Emit("MOVQ", _VP, _BX)                 // MOVQ    VP, _BX
    self.call_go(fn)                            // CALL_GO ${fn}
    self.Emit("MOVQ", _AX, _VP)                 // MOVQ    AX, VP
}

func (self *_Assembler) mapassign_fastx(t reflect.Type, fn obj.Addr) {
    self.mapassign_call_from_AX(t, fn)
    self.mapaccess_ptr(t)
}

func (self *_Assembler) mapassign_utext(t reflect.Type, addressable bool) {
    pv := false
    vk := t.Key()
    tk := t.Key()

    /* deref pointer if needed */
    if vk.Kind() == reflect.Ptr {
        pv = true
        vk = vk.Elem()
    }

    /* addressable value with pointer receiver */
    if addressable {
        pv = false
        tk = reflect.PtrTo(tk)
    }

    /* allocate the key, and call the unmarshaler */
    self.valloc(vk, _BX)                        // VALLOC  ${vk}, BX
    // must spill vk pointer since next call_go may invoke GC
    self.Emit("MOVQ" , _BX, _ARG_vk)
    self.Emit("MOVQ" , jit.Type(tk), _AX)       // MOVQ    ${tk}, AX
    self.Emit("MOVQ" , _ARG_sv_p, _CX)          // MOVQ    sv.p, CX
    self.Emit("MOVQ" , _ARG_sv_n, _DI)          // MOVQ    sv.n, DI
    self.call_go(_F_decodeTextUnmarshaler)      // CALL_GO decodeTextUnmarshaler
    self.Emit("TESTQ", _ET, _ET)                // TESTQ   ET, ET
    self.Sjmp("JNZ"  , _LB_error)               // JNZ     _error
    self.Emit("MOVQ" , _ARG_vk, _AX)            // MOVQ    VAR.vk, AX
    self.Emit("MOVQ", jit.Imm(0), _ARG_vk)

    /* select the correct assignment function */
    if !pv {
        self.mapassign_call_from_AX(t, _F_mapassign)
    } else {
        self.mapassign_fastx(t, _F_mapassign_fast64ptr)
    }
}

/** External Unmarshaler Routines **/

var (
    _F_skip_one = jit.Imm(int64(native.S_skip_one))
    _F_skip_array  = jit.Imm(int64(native.S_skip_array))
    _F_skip_number = jit.Imm(int64(native.S_skip_number))
)

func (self *_Assembler) unmarshal_json(t reflect.Type, deref bool, f obj.Addr) {
    self.call_sf(_F_skip_one)                                   // CALL_SF   skip_one
    self.Emit("TESTQ", _AX, _AX)                                // TESTQ     AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)                     // JS        _parse_error_v
    self.Emit("MOVQ", _IC, _VAR_ic)                             // store for mismatche error skip
    self.slice_from_r(_AX, 0)                                   // SLICE_R   AX, $0
    self.Emit("MOVQ" , _DI, _ARG_sv_p)                          // MOVQ      DI, sv.p
    self.Emit("MOVQ" , _SI, _ARG_sv_n)                          // MOVQ      SI, sv.n
    self.unmarshal_func(t, f, deref)     // UNMARSHAL json, ${t}, ${deref}
}

func (self *_Assembler) unmarshal_text(t reflect.Type, deref bool) {
    self.parse_string()                                         // PARSE     STRING
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)        // UNQUOTE   once, sv.p, sv.n
    self.unmarshal_func(t, _F_decodeTextUnmarshaler, deref)     // UNMARSHAL text, ${t}, ${deref}
}

func (self *_Assembler) unmarshal_func(t reflect.Type, fn obj.Addr, deref bool) {
    pt := t
    vk := t.Kind()

    /* allocate the field if needed */
    if deref && vk == reflect.Ptr {
        self.Emit("MOVQ" , _VP, _BX)                // MOVQ   VP, BX
        self.Emit("MOVQ" , jit.Ptr(_BX, 0), _BX)    // MOVQ   (BX), BX
        self.Emit("TESTQ", _BX, _BX)                // TESTQ  BX, BX
        self.Sjmp("JNZ"  , "_deref_{n}")            // JNZ    _deref_{n}
        self.valloc(t.Elem(), _BX)                  // VALLOC ${t.Elem()}, BX
        self.WriteRecNotAX(3, _BX, jit.Ptr(_VP, 0), false, false)    // MOVQ   BX, (VP)
        self.Link("_deref_{n}")                     // _deref_{n}:
    } else {
        /* set value pointer */
        self.Emit("MOVQ", _VP, _BX)                 // MOVQ   (VP), BX
    }

    /* set value type */
    self.Emit("MOVQ", jit.Type(pt), _AX)        // MOVQ ${pt}, AX

    /* set the source string and call the unmarshaler */
    self.Emit("MOVQ" , _ARG_sv_p, _CX)          // MOVQ    sv.p, CX
    self.Emit("MOVQ" , _ARG_sv_n, _DI)          // MOVQ    sv.n, DI
    self.call_go(fn)                            // CALL_GO ${fn}
    self.Emit("TESTQ", _ET, _ET)                // TESTQ   ET, ET
    if fn == _F_decodeJsonUnmarshalerQuoted {
        self.Sjmp("JZ"  , "_unmarshal_func_end_{n}")            // JZ   _unmarshal_func_end_{n}
        self.Emit("MOVQ", _I_json_MismatchQuotedError, _CX)     // MOVQ _I_json_MismatchQuotedError, CX
        self.Emit("CMPQ", _ET, _CX)            // check if MismatchQuotedError
        self.Sjmp("JNE" , _LB_error)           // JNE     _error
        self.Emit("MOVQ", jit.Type(t), _CX)    // store current type 
        self.Emit("MOVQ", _CX, _VAR_et)        // store current type as mismatched type
        self.Emit("MOVQ", _VAR_ic, _IC)        // recover the pos at mismatched, continue to parse
        self.Emit("XORL", _ET, _ET)            // clear ET
        self.Link("_unmarshal_func_end_{n}")
    } else {
        self.Sjmp("JNE" , _LB_error)           // JNE     _error
    }
}

/** Dynamic Decoding Routine **/

var (
    _F_decodeTypedPointer obj.Addr
)

func init() {
    _F_decodeTypedPointer = jit.Func(decodeTypedPointer)
}

func (self *_Assembler) decode_dynamic(vt obj.Addr, vp obj.Addr) {
    self.Emit("MOVQ" , vp, _SI)    // MOVQ    ${vp}, SI
    self.Emit("MOVQ" , vt, _DI)    // MOVQ    ${vt}, DI
    self.Emit("MOVQ", _ARG_sp, _AX)            // MOVQ    sp, AX
    self.Emit("MOVQ", _ARG_sl, _BX)            // MOVQ    sp, BX
    self.Emit("MOVQ" , _IC, _CX)                // MOVQ    IC, CX
    self.Emit("MOVQ" , _ST, _R8)                // MOVQ    ST, R8 
    self.Emit("MOVQ" , _ARG_fv, _R9)            // MOVQ    fv, R9
    self.save(_REG_rt...)
    self.Emit("MOVQ", _F_decodeTypedPointer, _IL)  // MOVQ ${fn}, R11
    self.Rjmp("CALL", _IL)      // CALL R11
    self.load(_REG_rt...)
    self.Emit("MOVQ" , _AX, _IC)                // MOVQ    AX, IC
    self.Emit("MOVQ" , _BX, _ET)                // MOVQ    BX, ET
    self.Emit("MOVQ" , _CX, _EP)                // MOVQ    CX, EP
    self.Emit("TESTQ", _ET, _ET)                // TESTQ   ET, ET
    self.Sjmp("JE", "_decode_dynamic_end_{n}")  // JE, _decode_dynamic_end_{n}
    self.Emit("MOVQ", _I_json_MismatchTypeError, _CX) // MOVQ _I_json_MismatchTypeError, CX
    self.Emit("CMPQ", _ET, _CX)                 // CMPQ ET, CX
    self.Sjmp("JNE",  _LB_error)                // JNE  LB_error
    self.Emit("MOVQ", _ET, _VAR_et)             // MOVQ ET, VAR_et
    self.WriteRecNotAX(14, _EP, jit.Ptr(_ST, _EpOffset), false, false) // MOVQ EP, stack.Ep
    self.Link("_decode_dynamic_end_{n}")
}

/** OpCode Assembler Functions **/

var (
    _F_memequal         = jit.Func(rt.MemEqual)
    _F_memmove          = jit.Func(rt.Memmove)
    _F_growslice        = jit.Func(rt.GrowSlice)
    _F_makeslice        = jit.Func(rt.MakeSliceStd)
    _F_makemap_small    = jit.Func(rt.MakemapSmall)
    _F_mapassign_fast64 = jit.Func(rt.Mapassign_fast64)
)

var (
    _F_lspace  = jit.Imm(int64(native.S_lspace))
    _F_strhash = jit.Imm(int64(caching.S_strhash))
)

var (
    _F_b64decode   = jit.Imm(int64(rt.SubrB64Decode))
    _F_decodeValue = jit.Imm(int64(_subr_decode_value))
)

var (
    _F_FieldMap_GetCaseInsensitive obj.Addr
    _Empty_Slice = []byte{}
    _Zero_Base = int64(uintptr(((*rt.GoSlice)(unsafe.Pointer(&_Empty_Slice))).Ptr))
)

const (
    _MODE_AVX2 = 1 << 2
)

const (
    _Fe_ID   = int64(unsafe.Offsetof(caching.FieldEntry{}.ID))
    _Fe_Name = int64(unsafe.Offsetof(caching.FieldEntry{}.Name))
    _Fe_Hash = int64(unsafe.Offsetof(caching.FieldEntry{}.Hash))
)

const (
    _Vk_Ptr       = int64(reflect.Ptr)
    _Gt_KindFlags = int64(unsafe.Offsetof(rt.GoType{}.KindFlags))
)

func init() {
    _F_FieldMap_GetCaseInsensitive = jit.Func((*caching.FieldMap).GetCaseInsensitive)
}

func (self *_Assembler) _asm_OP_any(_ *_Instr) {
    self.Emit("MOVQ"   , jit.Ptr(_VP, 8), _CX)              // MOVQ    8(VP), CX
    self.Emit("TESTQ"  , _CX, _CX)                          // TESTQ   CX, CX
    self.Sjmp("JZ"     , "_decode_{n}")                     // JZ      _decode_{n}
    self.Emit("CMPQ"   , _CX, _VP)                          // CMPQ    CX, VP
    self.Sjmp("JE"     , "_decode_{n}")                     // JE      _decode_{n}
    self.Emit("MOVQ"   , jit.Ptr(_VP, 0), _AX)              // MOVQ    (VP), AX
    self.Emit("MOVBLZX", jit.Ptr(_AX, _Gt_KindFlags), _DX)  // MOVBLZX _Gt_KindFlags(AX), DX
    self.Emit("ANDL"   , jit.Imm(rt.F_kind_mask), _DX)      // ANDL    ${F_kind_mask}, DX
    self.Emit("CMPL"   , _DX, jit.Imm(_Vk_Ptr))             // CMPL    DX, ${reflect.Ptr}
    self.Sjmp("JNE"    , "_decode_{n}")                     // JNE     _decode_{n}
    self.Emit("LEAQ"   , jit.Ptr(_VP, 8), _DI)              // LEAQ    8(VP), DI
    self.decode_dynamic(_AX, _DI)                           // DECODE  AX, DI
    self.Sjmp("JMP"    , "_decode_end_{n}")                 // JMP     _decode_end_{n}
    self.Link("_decode_{n}")                                // _decode_{n}:
    self.Emit("MOVQ"   , _ARG_fv, _DF)                      // MOVQ    fv, DF
    self.Emit("MOVQ"   , _ST, jit.Ptr(_SP, 0))              // MOVQ    _ST, (SP)
    self.call(_F_decodeValue)                               // CALL    decodeValue
    self.Emit("MOVQ"   , jit.Imm(0), jit.Ptr(_SP, 0))              // MOVQ    _ST, (SP)
    self.Emit("TESTQ"  , _EP, _EP)                          // TESTQ   EP, EP
    self.Sjmp("JNZ"    , _LB_parsing_error)                 // JNZ     _parsing_error
    self.Link("_decode_end_{n}")                            // _decode_end_{n}:
}

func (self *_Assembler) _asm_OP_dyn(p *_Instr) {
    self.Emit("MOVQ"   , jit.Type(p.vt()), _ET)             // MOVQ    ${p.vt()}, ET
    self.Emit("CMPQ"   , jit.Ptr(_VP, 8), jit.Imm(0))       // CMPQ    8(VP), $0
    self.Sjmp("JNE"     , "_decode_dyn_non_nil_{n}")                    // JE      _type_error

    /* if nil iface, call skip one */
    self.Emit("MOVQ", _IC, _VAR_ic)
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)       
    self.Sref("_decode_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)

    self.Link("_decode_dyn_non_nil_{n}")                    // _decode_dyn_non_nil_{n}:
    self.Emit("MOVQ"   , jit.Ptr(_VP, 0), _CX)              // MOVQ    (VP), CX
    self.Emit("MOVQ"   , jit.Ptr(_CX, 8), _CX)              // MOVQ    8(CX), CX
    self.Emit("MOVBLZX", jit.Ptr(_CX, _Gt_KindFlags), _DX)  // MOVBLZX _Gt_KindFlags(CX), DX
    self.Emit("ANDL"   , jit.Imm(rt.F_kind_mask), _DX)      // ANDL    ${F_kind_mask}, DX
    self.Emit("CMPL"   , _DX, jit.Imm(_Vk_Ptr))             // CMPL    DX, ${reflect.Ptr}
    self.Sjmp("JE"    , "_decode_dyn_ptr_{n}")              // JNE     _type_error

    self.Emit("MOVQ", _IC, _VAR_ic)
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)       
    self.Sref("_decode_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)

    self.Link("_decode_dyn_ptr_{n}")                        // _decode_dyn_ptr_{n}:
    self.Emit("LEAQ"   , jit.Ptr(_VP, 8), _DI)              // LEAQ    8(VP), DI
    self.decode_dynamic(_CX, _DI)                           // DECODE  CX, DI
    self.Link("_decode_end_{n}")                            // _decode_end_{n}:
}

func (self *_Assembler) _asm_OP_unsupported(p *_Instr) {
    self.Emit("MOVQ", jit.Type(p.vt()), _ET)               // MOVQ    ${p.vt()}, ET
    self.Sjmp("JMP" , _LB_type_error)                      // JMP     _LB_type_error
}

func (self *_Assembler) _asm_OP_str(_ *_Instr) {
    self.parse_string()                                     // PARSE   STRING
    self.unquote_once(jit.Ptr(_VP, 0), jit.Ptr(_VP, 8), false, true)     // UNQUOTE once, (VP), 8(VP)
}

func (self *_Assembler) _asm_OP_bin(_ *_Instr) {
    self.parse_string()                                 // PARSE  STRING
    self.slice_from(_VAR_st_Iv, -1)                     // SLICE  st.Iv, $-1
    self.Emit("MOVQ" , _DI, jit.Ptr(_VP, 0))            // MOVQ   DI, (VP)
    self.Emit("MOVQ" , _SI, jit.Ptr(_VP, 8))            // MOVQ   SI, 8(VP)
    self.Emit("SHRQ" , jit.Imm(2), _SI)                 // SHRQ   $2, SI
    self.Emit("LEAQ" , jit.Sib(_SI, _SI, 2, 0), _SI)    // LEAQ   (SI)(SI*2), SI
    self.Emit("MOVQ" , _SI, jit.Ptr(_VP, 16))           // MOVQ   SI, 16(VP)
    self.malloc_AX(_SI, _SI)                               // MALLOC SI, SI

    // TODO: due to base64x's bug, only use AVX mode now
    self.Emit("MOVL", jit.Imm(_MODE_JSON), _CX)          //  MOVL $_MODE_JSON, CX

    /* call the decoder */
    self.Emit("XORL" , _DX, _DX)                // XORL  DX, DX
    self.Emit("MOVQ" , _VP, _DI)                // MOVQ  VP, DI

    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _R8)    // MOVQ SI, （VP)
    self.WriteRecNotAX(4, _SI, jit.Ptr(_VP, 0), true, false)    // XCHGQ SI, (VP) 
    self.Emit("MOVQ" , _R8, _SI)

    self.Emit("XCHGQ", _DX, jit.Ptr(_VP, 8))    // XCHGQ DX, 8(VP)
    self.call_c(_F_b64decode)                     // CALL  b64decode
    self.Emit("TESTQ", _AX, _AX)                // TESTQ AX, AX
    self.Sjmp("JS"   , _LB_base64_error)        // JS    _base64_error
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 8))    // MOVQ  AX, 8(VP)
}

func (self *_Assembler) _asm_OP_bool(_ *_Instr) {
    self.Emit("LEAQ", jit.Ptr(_IC, 4), _AX)                     // LEAQ 4(IC), AX
    self.Emit("CMPQ", _AX, _IL)                                 // CMPQ AX, IL
    self.Sjmp("JA"  , _LB_eof_error)                            // JA   _eof_error
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('f'))    // CMPB (IP)(IC), $'f'
    self.Sjmp("JE"  , "_false_{n}")                             // JE   _false_{n}
    self.Emit("MOVL", jit.Imm(_IM_true), _CX)                   // MOVL $"true", CX
    self.Emit("CMPL", _CX, jit.Sib(_IP, _IC, 1, 0))             // CMPL CX, (IP)(IC)
    self.Sjmp("JE" , "_bool_true_{n}")          
    // try to skip the value
    self.Emit("MOVQ", _IC, _VAR_ic)           
    self.Emit("MOVQ", _T_bool, _ET)     
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
    self.Sref("_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one) 

    self.Link("_bool_true_{n}")
    self.Emit("MOVQ", _AX, _IC)                                 // MOVQ AX, IC
    self.Emit("MOVB", jit.Imm(1), jit.Ptr(_VP, 0))              // MOVB $1, (VP)
    self.Sjmp("JMP" , "_end_{n}")                               // JMP  _end_{n}
    self.Link("_false_{n}")                                     // _false_{n}:
    self.Emit("ADDQ", jit.Imm(1), _AX)                          // ADDQ $1, AX
    self.Emit("ADDQ", jit.Imm(1), _IC)                          // ADDQ $1, IC
    self.Emit("CMPQ", _AX, _IL)                                 // CMPQ AX, IL
    self.Sjmp("JA"  , _LB_eof_error)                            // JA   _eof_error
    self.Emit("MOVL", jit.Imm(_IM_alse), _CX)                   // MOVL $"alse", CX
    self.Emit("CMPL", _CX, jit.Sib(_IP, _IC, 1, 0))             // CMPL CX, (IP)(IC)
    self.Sjmp("JNE" , _LB_im_error)                             // JNE  _im_error
    self.Emit("MOVQ", _AX, _IC)                                 // MOVQ AX, IC
    self.Emit("XORL", _AX, _AX)                                 // XORL AX, AX
    self.Emit("MOVB", _AX, jit.Ptr(_VP, 0))                     // MOVB AX, (VP)
    self.Link("_end_{n}")                                       // _end_{n}:
}

func (self *_Assembler) _asm_OP_num(_ *_Instr) {
    self.Emit("MOVQ", jit.Imm(0), _VAR_fl)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('"'))
    self.Emit("MOVQ", _IC, _BX)
    self.Sjmp("JNE", "_skip_number_{n}")
    self.Emit("MOVQ", jit.Imm(1), _VAR_fl)
    self.Emit("ADDQ", jit.Imm(1), _IC)
    self.Link("_skip_number_{n}")

    /* call skip_number */
    self.Emit("LEAQ", _ARG_s, _DI)                      // LEAQ  s<>+0(FP), DI
    self.Emit("MOVQ", _IC, _ARG_ic)                     // MOVQ  IC, ic<>+16(FP)
    self.Emit("LEAQ", _ARG_ic, _SI)                     // LEAQ  ic<>+16(FP), SI
    self.callc(_F_skip_number)                          // CALL  _F_skip_number
    self.Emit("MOVQ", _ARG_ic, _IC)                     // MOVQ  ic<>+16(FP), IC
    self.Emit("TESTQ", _AX, _AX)                        // TESTQ AX, AX
    self.Sjmp("JNS"   , "_num_next_{n}")

    /* call skip one */
    self.Emit("MOVQ", _BX, _VAR_ic)           
    self.Emit("MOVQ", _T_number, _ET)     
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)       
    self.Sref("_num_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)

    /* assign string */
    self.Link("_num_next_{n}")
    self.slice_from_r(_AX, 0)
    self.Emit("BTQ", jit.Imm(_F_copy_string), _ARG_fv)
    self.Sjmp("JNC", "_num_write_{n}")
    self.Byte(0x4c, 0x8d, 0x0d)         // LEAQ (PC), R9
    self.Sref("_num_write_{n}", 4)
    self.Sjmp("JMP", "_copy_string")
    self.Link("_num_write_{n}")
    self.Emit("MOVQ", _SI, jit.Ptr(_VP, 8))     // MOVQ  SI, 8(VP)
    self.WriteRecNotAX(13, _DI, jit.Ptr(_VP, 0), false, false)
    self.Emit("CMPQ", _VAR_fl, jit.Imm(1))
    self.Sjmp("JNE", "_num_end_{n}")
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('"'))
    self.Sjmp("JNE", _LB_char_0_error)
    self.Emit("ADDQ", jit.Imm(1), _IC)
    self.Link("_num_end_{n}")
}

func (self *_Assembler) _asm_OP_i8(_ *_Instr) {
    var pin = "_i8_end_{n}"
    self.parse_signed(int8Type, pin, -1)                                                 // PARSE int8
    self.range_signed_CX(_I_int8, _T_int8, math.MinInt8, math.MaxInt8)     // RANGE int8
    self.Emit("MOVB", _CX, jit.Ptr(_VP, 0))                             // MOVB  CX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_i16(_ *_Instr) {
    var pin = "_i16_end_{n}"
    self.parse_signed(int16Type, pin, -1)                                                     // PARSE int16
    self.range_signed_CX(_I_int16, _T_int16, math.MinInt16, math.MaxInt16)     // RANGE int16
    self.Emit("MOVW", _CX, jit.Ptr(_VP, 0))                                 // MOVW  CX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_i32(_ *_Instr) {
    var pin = "_i32_end_{n}"
    self.parse_signed(int32Type, pin, -1)                                                     // PARSE int32
    self.range_signed_CX(_I_int32, _T_int32, math.MinInt32, math.MaxInt32)     // RANGE int32
    self.Emit("MOVL", _CX, jit.Ptr(_VP, 0))                                 // MOVL  CX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_i64(_ *_Instr) {
    var pin = "_i64_end_{n}"
    self.parse_signed(int64Type, pin, -1)                         // PARSE int64
    self.Emit("MOVQ", _VAR_st_Iv, _AX)          // MOVQ  st.Iv, AX
    self.Emit("MOVQ", _AX, jit.Ptr(_VP, 0))     // MOVQ  AX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u8(_ *_Instr) {
    var pin = "_u8_end_{n}"
    self.parse_unsigned(uint8Type, pin, -1)                                   // PARSE uint8
    self.range_unsigned_CX(_I_uint8, _T_uint8, math.MaxUint8)  // RANGE uint8
    self.Emit("MOVB", _CX, jit.Ptr(_VP, 0))                 // MOVB  CX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u16(_ *_Instr) {
    var pin = "_u16_end_{n}"
    self.parse_unsigned(uint16Type, pin, -1)                                       // PARSE uint16
    self.range_unsigned_CX(_I_uint16, _T_uint16, math.MaxUint16)   // RANGE uint16
    self.Emit("MOVW", _CX, jit.Ptr(_VP, 0))                     // MOVW  CX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u32(_ *_Instr) {
    var pin = "_u32_end_{n}"
    self.parse_unsigned(uint32Type, pin, -1)                                       // PARSE uint32
    self.range_uint32_CX(_I_uint32, _T_uint32)   // RANGE uint32
    self.Emit("MOVL", _CX, jit.Ptr(_VP, 0))                     // MOVL  CX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u64(_ *_Instr) {
    var pin = "_u64_end_{n}"
    self.parse_unsigned(uint64Type, pin, -1)                       // PARSE uint64
    self.Emit("MOVQ", _VAR_st_Iv, _AX)          // MOVQ  st.Iv, AX
    self.Emit("MOVQ", _AX, jit.Ptr(_VP, 0))     // MOVQ  AX, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_f32(_ *_Instr) {
    var pin = "_f32_end_{n}"
    self.parse_number(float32Type, pin, -1)                         // PARSE NUMBER
    self.range_single_X0()                         // RANGE float32
    self.Emit("MOVSS", _X0, jit.Ptr(_VP, 0))    // MOVSS X0, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_f64(_ *_Instr) {
    var pin = "_f64_end_{n}"
    self.parse_number(float64Type, pin, -1)                         // PARSE NUMBER
    self.Emit("MOVSD", _VAR_st_Dv, _X0)         // MOVSD st.Dv, X0
    self.Emit("MOVSD", _X0, jit.Ptr(_VP, 0))    // MOVSD X0, (VP)
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_unquote(_ *_Instr) {
    self.check_eof(2)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('\\'))   // CMPB    (IP)(IC), $'\\'
    self.Sjmp("JNE" , _LB_char_0_error)                         // JNE     _char_0_error
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 1), jit.Imm('"'))    // CMPB    1(IP)(IC), $'"'
    self.Sjmp("JNE" , _LB_char_1_error)                         // JNE     _char_1_error
    self.Emit("ADDQ", jit.Imm(2), _IC)                          // ADDQ    $2, IC
    self.parse_string()                                         // PARSE   STRING
    self.unquote_twice(jit.Ptr(_VP, 0), jit.Ptr(_VP, 8), false)        // UNQUOTE twice, (VP), 8(VP)
}

func (self *_Assembler) _asm_OP_nil_1(_ *_Instr) {
    self.Emit("XORL", _AX, _AX)                 // XORL AX, AX
    self.Emit("MOVQ", _AX, jit.Ptr(_VP, 0))     // MOVQ AX, (VP)
}

func (self *_Assembler) _asm_OP_nil_2(_ *_Instr) {
    self.Emit("PXOR" , _X0, _X0)                // PXOR  X0, X0
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 0))    // MOVOU X0, (VP)
}

func (self *_Assembler) _asm_OP_nil_3(_ *_Instr) {
    self.Emit("XORL" , _AX, _AX)                // XORL  AX, AX
    self.Emit("PXOR" , _X0, _X0)                // PXOR  X0, X0
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 0))    // MOVOU X0, (VP)
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 16))   // MOVOU AX, 16(VP)
}

var (
    bytes []byte = make([]byte, 0)
    zerobytes = (*rt.GoSlice)(unsafe.Pointer(&bytes)).Ptr
    _ZERO_PTR = jit.Imm(int64(uintptr(zerobytes)))
)

func (self *_Assembler) _asm_OP_empty_bytes(_ *_Instr) {
    self.Emit("MOVQ", _ZERO_PTR, _AX)
    self.Emit("PXOR" , _X0, _X0)
    self.Emit("MOVQ", _AX,  jit.Ptr(_VP, 0))
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 8))
}

func (self *_Assembler) _asm_OP_deref(p *_Instr) {
    self.vfollow(p.vt())
}

func (self *_Assembler) _asm_OP_index(p *_Instr) {
    self.Emit("MOVQ", jit.Imm(p.i64()), _AX)    // MOVQ ${p.vi()}, AX
    self.Emit("ADDQ", _AX, _VP)                 // ADDQ _AX, _VP
}

func (self *_Assembler) _asm_OP_is_null(p *_Instr) {
    self.Emit("LEAQ"   , jit.Ptr(_IC, 4), _AX)                          // LEAQ    4(IC), AX
    self.Emit("CMPQ"   , _AX, _IL)                                      // CMPQ    AX, IL
    self.Sjmp("JA"     , "_not_null_{n}")                               // JA      _not_null_{n}
    self.Emit("CMPL"   , jit.Sib(_IP, _IC, 1, 0), jit.Imm(_IM_null))    // CMPL    (IP)(IC), $"null"
    self.Emit("CMOVQEQ", _AX, _IC)                                      // CMOVQEQ AX, IC
    self.Xjmp("JE"     , p.vi())                                        // JE      {p.vi()}
    self.Link("_not_null_{n}")                                          // _not_null_{n}:
}

func (self *_Assembler) _asm_OP_is_null_quote(p *_Instr) {
    self.Emit("LEAQ"   , jit.Ptr(_IC, 5), _AX)                          // LEAQ    4(IC), AX
    self.Emit("CMPQ"   , _AX, _IL)                                      // CMPQ    AX, IL
    self.Sjmp("JA"     , "_not_null_quote_{n}")                         // JA      _not_null_quote_{n}
    self.Emit("CMPL"   , jit.Sib(_IP, _IC, 1, 0), jit.Imm(_IM_null))    // CMPL    (IP)(IC), $"null"
    self.Sjmp("JNE"    , "_not_null_quote_{n}")                         // JNE     _not_null_quote_{n}
    self.Emit("CMPB"   , jit.Sib(_IP, _IC, 1, 4), jit.Imm('"'))         // CMPB    4(IP)(IC), $'"'
    self.Emit("CMOVQEQ", _AX, _IC)                                      // CMOVQEQ AX, IC
    self.Xjmp("JE"     , p.vi())                                        // JE      {p.vi()}
    self.Link("_not_null_quote_{n}")                                    // _not_null_quote_{n}:
}

func (self *_Assembler) _asm_OP_map_init(_ *_Instr) {
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _AX)    // MOVQ    (VP), AX
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JNZ"  , "_end_{n}")              // JNZ     _end_{n}
    self.call_go(_F_makemap_small)              // CALL_GO makemap_small
    self.WritePtrAX(6, jit.Ptr(_VP, 0), false)    // MOVQ    AX, (VP)
    self.Link("_end_{n}")                       // _end_{n}:
    self.Emit("MOVQ" , _AX, _VP)                // MOVQ    AX, VP
}

func (self *_Assembler) _asm_OP_map_key_i8(p *_Instr) {
    self.parse_signed(int8Type, "", p.vi())                                                 // PARSE     int8
    self.range_signed_CX(_I_int8, _T_int8, math.MinInt8, math.MaxInt8)     // RANGE     int8
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                              // MAPASSIGN int8, mapassign, st.Iv
}

func (self *_Assembler) _asm_OP_map_key_i16(p *_Instr) {
    self.parse_signed(int16Type, "", p.vi())                                                     // PARSE     int16
    self.range_signed_CX(_I_int16, _T_int16, math.MinInt16, math.MaxInt16)     // RANGE     int16
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                                  // MAPASSIGN int16, mapassign, st.Iv
}

func (self *_Assembler) _asm_OP_map_key_i32(p *_Instr) {
    self.parse_signed(int32Type, "", p.vi())                                                     // PARSE     int32
    self.range_signed_CX(_I_int32, _T_int32, math.MinInt32, math.MaxInt32)     // RANGE     int32
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)                                  // MAPASSIGN int32, mapassign, st.Iv
    } else {
        self.Emit("MOVQ", _CX, _AX)                                         // MOVQ CX, AX
        self.mapassign_fastx(vt, _F_mapassign_fast32)                       // MAPASSIGN int32, mapassign_fast32
    }
}

func (self *_Assembler) _asm_OP_map_key_i64(p *_Instr) {
    self.parse_signed(int64Type, "", p.vi())                                 // PARSE     int64
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)              // MAPASSIGN int64, mapassign, st.Iv
    } else {
        self.Emit("MOVQ", _VAR_st_Iv, _AX)              // MOVQ      st.Iv, AX
        self.mapassign_fastx(vt, _F_mapassign_fast64)   // MAPASSIGN int64, mapassign_fast64
    }
}

func (self *_Assembler) _asm_OP_map_key_u8(p *_Instr) {
    self.parse_unsigned(uint8Type, "", p.vi())                                   // PARSE     uint8
    self.range_unsigned_CX(_I_uint8, _T_uint8, math.MaxUint8)  // RANGE     uint8
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                    // MAPASSIGN uint8, vt.Iv
}

func (self *_Assembler) _asm_OP_map_key_u16(p *_Instr) {
    self.parse_unsigned(uint16Type, "", p.vi())                                       // PARSE     uint16
    self.range_unsigned_CX(_I_uint16, _T_uint16, math.MaxUint16)   // RANGE     uint16
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                      // MAPASSIGN uint16, vt.Iv
}

func (self *_Assembler) _asm_OP_map_key_u32(p *_Instr) {
    self.parse_unsigned(uint32Type, "", p.vi())                                       // PARSE     uint32
    self.range_unsigned_CX(_I_uint32, _T_uint32, math.MaxUint32)   // RANGE     uint32
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)                      // MAPASSIGN uint32, vt.Iv
    } else {
        self.Emit("MOVQ", _CX, _AX)                             // MOVQ CX, AX
        self.mapassign_fastx(vt, _F_mapassign_fast32)           // MAPASSIGN uint32, mapassign_fast32
    }
}

func (self *_Assembler) _asm_OP_map_key_u64(p *_Instr) {
    self.parse_unsigned(uint64Type, "", p.vi())                                       // PARSE     uint64
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)                      // MAPASSIGN uint64, vt.Iv
    } else {
        self.Emit("MOVQ", _VAR_st_Iv, _AX)                      // MOVQ      st.Iv, AX
        self.mapassign_fastx(vt, _F_mapassign_fast64)           // MAPASSIGN uint64, mapassign_fast64
    }
}

func (self *_Assembler) _asm_OP_map_key_f32(p *_Instr) {
    self.parse_number(float32Type, "", p.vi())                     // PARSE     NUMBER
    self.range_single_X0()                     // RANGE     float32
    self.Emit("MOVSS", _X0, _VAR_st_Dv)     // MOVSS     X0, st.Dv
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Dv)  // MAPASSIGN ${p.vt()}, mapassign, st.Dv
}

func (self *_Assembler) _asm_OP_map_key_f64(p *_Instr) {
    self.parse_number(float64Type, "", p.vi())                     // PARSE     NUMBER
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Dv)  // MAPASSIGN ${p.vt()}, mapassign, st.Dv
}

func (self *_Assembler) _asm_OP_map_key_str(p *_Instr) {
    self.parse_string()                          // PARSE     STRING
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)      // UNQUOTE   once, sv.p, sv.n
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.valloc(vt.Key(), _DI)
        self.Emit("MOVOU", _ARG_sv, _X0)
        self.Emit("MOVOU", _X0, jit.Ptr(_DI, 0))
        self.mapassign_std(vt, jit.Ptr(_DI, 0))        // MAPASSIGN string, DI, SI
    } else {
        self.mapassign_str_fast(vt, _ARG_sv_p, _ARG_sv_n)    // MAPASSIGN string, DI, SI
    }
}

func (self *_Assembler) _asm_OP_map_key_utext(p *_Instr) {
    self.parse_string()                         // PARSE     STRING
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)     // UNQUOTE   once, sv.p, sv.n
    self.mapassign_utext(p.vt(), false)         // MAPASSIGN utext, ${p.vt()}, false
}

func (self *_Assembler) _asm_OP_map_key_utext_p(p *_Instr) {
    self.parse_string()                         // PARSE     STRING
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)     // UNQUOTE   once, sv.p, sv.n
    self.mapassign_utext(p.vt(), true)          // MAPASSIGN utext, ${p.vt()}, true
}

func (self *_Assembler) _asm_OP_array_skip(_ *_Instr) {
    self.call_sf(_F_skip_array)                 // CALL_SF skip_array
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)     // JS      _parse_error_v
}

func (self *_Assembler) _asm_OP_array_clear(p *_Instr) {
    self.mem_clear_rem(p.i64(), true)
}

func (self *_Assembler) _asm_OP_array_clear_p(p *_Instr) {
    self.mem_clear_rem(p.i64(), false)
}

func (self *_Assembler) _asm_OP_slice_init(p *_Instr) {
    self.Emit("XORL" , _AX, _AX)                    // XORL    AX, AX
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 8))        // MOVQ    AX, 8(VP)
    self.Emit("MOVQ" , jit.Ptr(_VP, 16), _BX)       // MOVQ    16(VP), BX
    self.Emit("TESTQ", _BX, _BX)                    // TESTQ   BX, BX
    self.Sjmp("JNZ"  , "_done_{n}")                 // JNZ     _done_{n}
    self.Emit("MOVQ" , jit.Imm(_MinSlice), _CX)     // MOVQ    ${_MinSlice}, CX
    self.Emit("MOVQ" , _CX, jit.Ptr(_VP, 16))       // MOVQ    CX, 16(VP)
    self.Emit("MOVQ" , jit.Type(p.vt()), _AX)       // MOVQ    ${p.vt()}, DX
    self.call_go(_F_makeslice)                      // CALL_GO makeslice
    self.WritePtrAX(7, jit.Ptr(_VP, 0), false)      // MOVQ    AX, (VP)
    self.Emit("XORL" , _AX, _AX)                    // XORL    AX, AX
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 8))        // MOVQ    AX, 8(VP)
    self.Link("_done_{n}")                          // _done_{n}
}

func (self *_Assembler) _asm_OP_check_empty(p *_Instr) {
    rbracket := p.vb()
    if rbracket == ']' {
        self.check_eof(1)
        self.Emit("LEAQ", jit.Ptr(_IC, 1), _AX)                              // LEAQ    1(IC), AX
        self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(rbracket))) // CMPB    (IP)(IC), ']'
        self.Sjmp("JNE" , "_not_empty_array_{n}")                            // JNE     _not_empty_array_{n}
        self.Emit("MOVQ", _AX, _IC)                                          // MOVQ    AX, IC
        self.Emit("MOVQ", jit.Imm(_Zero_Base), _AX)
        self.WritePtrAX(9, jit.Ptr(_VP, 0), false)
        self.Emit("PXOR", _X0, _X0)                                          // PXOR    X0, X0
        self.Emit("MOVOU", _X0, jit.Ptr(_VP, 8))                             // MOVOU   X0, 8(VP)
        self.Xjmp("JMP" , p.vi())                                            // JMP     {p.vi()}
        self.Link("_not_empty_array_{n}")
    } else {
        panic("only implement check empty array here!")
    }
}

func (self *_Assembler) _asm_OP_slice_append(p *_Instr) {
    self.Emit("MOVQ" , jit.Ptr(_VP, 8), _AX)            // MOVQ    8(VP), AX
    self.Emit("CMPQ" , _AX, jit.Ptr(_VP, 16))           // CMPQ    AX, 16(VP)
    self.Sjmp("JB"   , "_index_{n}")                    // JB      _index_{n}
    self.Emit("MOVQ" , _AX, _SI)                        // MOVQ    AX, SI
    self.Emit("SHLQ" , jit.Imm(1), _SI)                 // SHLQ    $1, SI
    self.Emit("MOVQ" , jit.Type(p.vt()), _AX)           // MOVQ    ${p.vt()}, AX
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _BX)            // MOVQ   (VP), BX
    self.Emit("MOVQ" , jit.Ptr(_VP, 8), _CX)            // MOVQ    8(VP), CX
    self.Emit("MOVQ" , jit.Ptr(_VP, 16), _DI)           // MOVQ    16(VP), DI
    self.call_go(_F_growslice)                          // CALL_GO growslice
    self.WritePtrAX(8, jit.Ptr(_VP, 0), false)          // MOVQ    AX, (VP)
    self.Emit("MOVQ" , _BX, jit.Ptr(_VP, 8))            // MOVQ    BX, 8(VP)
    self.Emit("MOVQ" , _CX, jit.Ptr(_VP, 16))           // MOVQ    CX, 16(VP)

    // because growslice not zero memory {oldcap, newlen} when append et not has ptrdata.
    // but we should zero it, avoid decode it as random values.
    if rt.UnpackType(p.vt()).PtrData == 0 {
        self.Emit("MOVQ" , _CX, _DI)                        // MOVQ    CX, DI
        self.Emit("SUBQ" , _BX, _DI)                        // MOVQ    BX, DI
    
        self.Emit("ADDQ" , jit.Imm(1), jit.Ptr(_VP, 8))     // ADDQ    $1, 8(VP)
        self.Emit("MOVQ" , _AX, _VP)                        // MOVQ    AX, VP
        self.Emit("MOVQ" , jit.Imm(int64(p.vlen())), _CX)   // MOVQ    ${p.vlen()}, CX
        self.Emit("MOVQ" , _BX, _AX)                        // MOVQ    BX, AX 
        self.From("MULQ" , _CX)                             // MULQ    CX
        self.Emit("ADDQ" , _AX, _VP)                        // ADDQ    AX, VP

        self.Emit("MOVQ" , _DI, _AX)                        // MOVQ    SI, AX
        self.From("MULQ" , _CX)                             // MULQ    BX
        self.Emit("MOVQ" , _AX, _BX)                        // ADDQ    AX, BX
        self.Emit("MOVQ" , _VP, _AX)                        // MOVQ    VP, AX
        self.mem_clear_fn(true)                             // CALL_GO memclr{Has,NoHeap}
        self.Sjmp("JMP", "_append_slice_end_{n}")
    }

    self.Emit("MOVQ" , _BX, _AX)                        // MOVQ    BX, AX
    self.Link("_index_{n}")                             // _index_{n}:
    self.Emit("ADDQ" , jit.Imm(1), jit.Ptr(_VP, 8))     // ADDQ    $1, 8(VP)
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _VP)            // MOVQ    (VP), VP
    self.Emit("MOVQ" , jit.Imm(int64(p.vlen())), _CX)   // MOVQ    ${p.vlen()}, CX
    self.From("MULQ" , _CX)                             // MULQ    CX
    self.Emit("ADDQ" , _AX, _VP)                        // ADDQ    AX, VP
    self.Link("_append_slice_end_{n}")
}

func (self *_Assembler) _asm_OP_object_next(_ *_Instr) {
    self.call_sf(_F_skip_one)                   // CALL_SF skip_one
    self.Emit("TESTQ", _AX, _AX)                // TESTQ   AX, AX
    self.Sjmp("JS"   , _LB_parsing_error_v)     // JS      _parse_error_v
}

func (self *_Assembler) _asm_OP_struct_field(p *_Instr) {
    assert_eq(caching.FieldEntrySize, 32, "invalid field entry size")
    self.Emit("MOVQ" , jit.Imm(-1), _AX)                        // MOVQ    $-1, AX
    self.Emit("MOVQ" , _AX, _VAR_sr)                            // MOVQ    AX, sr
    self.parse_string()                                         // PARSE   STRING
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, false)                     // UNQUOTE once, sv.p, sv.n
    self.Emit("LEAQ" , _ARG_sv, _AX)                            // LEAQ    sv, AX
    self.Emit("XORL" , _BX, _BX)                                // XORL    BX, BX
    self.call_go(_F_strhash)                                    // CALL_GO strhash
    self.Emit("MOVQ" , _AX, _R9)                                // MOVQ    AX, R9
    self.Emit("MOVQ" , jit.Imm(freezeFields(p.vf())), _CX)      // MOVQ    ${p.vf()}, CX
    self.Emit("MOVQ" , jit.Ptr(_CX, caching.FieldMap_b), _SI)   // MOVQ    FieldMap.b(CX), SI
    self.Emit("MOVQ" , jit.Ptr(_CX, caching.FieldMap_N), _CX)   // MOVQ    FieldMap.N(CX), CX
    self.Emit("TESTQ", _CX, _CX)                                // TESTQ   CX, CX
    self.Sjmp("JZ"   , "_try_lowercase_{n}")                    // JZ      _try_lowercase_{n}
    self.Link("_loop_{n}")                                      // _loop_{n}:
    self.Emit("XORL" , _DX, _DX)                                // XORL    DX, DX
    self.From("DIVQ" , _CX)                                     // DIVQ    CX
    self.Emit("LEAQ" , jit.Ptr(_DX, 1), _AX)                    // LEAQ    1(DX), AX
    self.Emit("SHLQ" , jit.Imm(5), _DX)                         // SHLQ    $5, DX
    self.Emit("LEAQ" , jit.Sib(_SI, _DX, 1, 0), _DI)            // LEAQ    (SI)(DX), DI
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_Hash), _R8)             // MOVQ    FieldEntry.Hash(DI), R8
    self.Emit("TESTQ", _R8, _R8)                                // TESTQ   R8, R8
    self.Sjmp("JZ"   , "_try_lowercase_{n}")                    // JZ      _try_lowercase_{n}
    self.Emit("CMPQ" , _R8, _R9)                                // CMPQ    R8, R9
    self.Sjmp("JNE"  , "_loop_{n}")                             // JNE     _loop_{n}
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_Name + 8), _DX)         // MOVQ    FieldEntry.Name+8(DI), DX
    self.Emit("CMPQ" , _DX, _ARG_sv_n)                          // CMPQ    DX, sv.n
    self.Sjmp("JNE"  , "_loop_{n}")                             // JNE     _loop_{n}
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_ID), _R8)               // MOVQ    FieldEntry.ID(DI), R8
    self.Emit("MOVQ" , _AX, _VAR_ss_AX)                         // MOVQ    AX, ss.AX
    self.Emit("MOVQ" , _CX, _VAR_ss_CX)                         // MOVQ    CX, ss.CX
    self.Emit("MOVQ" , _SI, _VAR_ss_SI)                         // MOVQ    SI, ss.SI
    self.Emit("MOVQ" , _R8, _VAR_ss_R8)                         // MOVQ    R8, ss.R8
    self.Emit("MOVQ" , _R9, _VAR_ss_R9)                         // MOVQ    R9, ss.R9
    self.Emit("MOVQ" , _ARG_sv_p, _AX)                          // MOVQ    _VAR_sv_p, AX
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_Name), _CX)             // MOVQ    FieldEntry.Name(DI), CX
    self.Emit("MOVQ" , _CX, _BX)                                // MOVQ    CX, 8(SP)
    self.Emit("MOVQ" , _DX, _CX)                                // MOVQ    DX, 16(SP)
    self.call_go(_F_memequal)                                   // CALL_GO memequal
    self.Emit("MOVB" , _AX, _DX)                                // MOVB    24(SP), DX
    self.Emit("MOVQ" , _VAR_ss_AX, _AX)                         // MOVQ    ss.AX, AX
    self.Emit("MOVQ" , _VAR_ss_CX, _CX)                         // MOVQ    ss.CX, CX
    self.Emit("MOVQ" , _VAR_ss_SI, _SI)                         // MOVQ    ss.SI, SI
    self.Emit("MOVQ" , _VAR_ss_R9, _R9)                         // MOVQ    ss.R9, R9
    self.Emit("TESTB", _DX, _DX)                                // TESTB   DX, DX
    self.Sjmp("JZ"   , "_loop_{n}")                             // JZ      _loop_{n}
    self.Emit("MOVQ" , _VAR_ss_R8, _R8)                         // MOVQ    ss.R8, R8
    self.Emit("MOVQ" , _R8, _VAR_sr)                            // MOVQ    R8, sr
    self.Sjmp("JMP"  , "_end_{n}")                              // JMP     _end_{n}
    self.Link("_try_lowercase_{n}")                             // _try_lowercase_{n}:
    self.Emit("BTQ"  , jit.Imm(_F_case_sensitive), _ARG_fv)     // check if enable option CaseSensitive
    self.Sjmp("JC"   , "_unknown_{n}")                         
    self.Emit("MOVQ" , jit.Imm(referenceFields(p.vf())), _AX)   // MOVQ    ${p.vf()}, AX
    self.Emit("MOVQ", _ARG_sv_p, _BX)                            // MOVQ   sv, BX
    self.Emit("MOVQ", _ARG_sv_n, _CX)                            // MOVQ   sv, CX
    self.call_go(_F_FieldMap_GetCaseInsensitive)                // CALL_GO FieldMap::GetCaseInsensitive
    self.Emit("MOVQ" , _AX, _VAR_sr)                            // MOVQ    AX, _VAR_sr
    self.Emit("TESTQ", _AX, _AX)                                // TESTQ   AX, AX
    self.Sjmp("JNS"  , "_end_{n}")                              // JNS     _end_{n}
    self.Link("_unknown_{n}")
    // HACK: because `_VAR_sr` maybe used in `F_vstring`, so we should clear here again for `_OP_switch`.
    self.Emit("MOVQ" , jit.Imm(-1), _AX)                        // MOVQ    $-1, AX
    self.Emit("MOVQ" , _AX, _VAR_sr)                            // MOVQ    AX, sr
    self.Emit("BTQ"  , jit.Imm(_F_disable_unknown), _ARG_fv)    // BTQ     ${_F_disable_unknown}, fv
    self.Sjmp("JC"   , _LB_field_error)                         // JC      _field_error
    self.Link("_end_{n}")                                       // _end_{n}:
}

func (self *_Assembler) _asm_OP_unmarshal(p *_Instr) {
    if iv := p.i64(); iv != 0 {
        self.unmarshal_json(p.vt(), true, _F_decodeJsonUnmarshalerQuoted)
    } else {
        self.unmarshal_json(p.vt(), true, _F_decodeJsonUnmarshaler)
    }
}

func (self *_Assembler) _asm_OP_unmarshal_p(p *_Instr) {
    if iv := p.i64(); iv != 0 {
        self.unmarshal_json(p.vt(), false, _F_decodeJsonUnmarshalerQuoted)
    } else {
        self.unmarshal_json(p.vt(), false, _F_decodeJsonUnmarshaler)
    }
}

func (self *_Assembler) _asm_OP_unmarshal_text(p *_Instr) {
    self.unmarshal_text(p.vt(), true)
}

func (self *_Assembler) _asm_OP_unmarshal_text_p(p *_Instr) {
    self.unmarshal_text(p.vt(), false)
}

func (self *_Assembler) _asm_OP_lspace(_ *_Instr) {
    self.lspace("_{n}")
}

func (self *_Assembler) lspace(subfix string) {
    var label = "_lspace" + subfix
    self.Emit("CMPQ"   , _IC, _IL)                      // CMPQ    IC, IL
    self.Sjmp("JAE"    , _LB_eof_error)                 // JAE     _eof_error
    self.Emit("MOVQ"   , jit.Imm(_BM_space), _DX)       // MOVQ    _BM_space, DX
    self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  // MOVBQZX (IP)(IC), AX
    self.Emit("CMPQ"   , _AX, jit.Imm(' '))             // CMPQ    AX, $' '
    self.Sjmp("JA"     , label)                // JA      _nospace_{n}
    self.Emit("BTQ"    , _AX, _DX)                      // BTQ     AX, DX
    self.Sjmp("JNC"    , label)                // JNC     _nospace_{n}

    /* test up to 4 characters */
    for i := 0; i < 3; i++ {
        self.Emit("ADDQ"   , jit.Imm(1), _IC)               // ADDQ    $1, IC
        self.Emit("CMPQ"   , _IC, _IL)                      // CMPQ    IC, IL
        self.Sjmp("JAE"    , _LB_eof_error)                 // JAE     _eof_error
        self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  // MOVBQZX (IP)(IC), AX
        self.Emit("CMPQ"   , _AX, jit.Imm(' '))             // CMPQ    AX, $' '
        self.Sjmp("JA"     , label)                // JA      _nospace_{n}
        self.Emit("BTQ"    , _AX, _DX)                      // BTQ     AX, DX
        self.Sjmp("JNC"    , label)                // JNC     _nospace_{n}
    }

    /* handle over to the native function */
    self.Emit("MOVQ"   , _IP, _DI)                      // MOVQ    IP, DI
    self.Emit("MOVQ"   , _IL, _SI)                      // MOVQ    IL, SI
    self.Emit("MOVQ"   , _IC, _DX)                      // MOVQ    IC, DX
    self.callc(_F_lspace)                                // CALL    lspace
    self.Emit("TESTQ"  , _AX, _AX)                      // TESTQ   AX, AX
    self.Sjmp("JS"     , _LB_parsing_error_v)           // JS      _parsing_error_v
    self.Emit("CMPQ"   , _AX, _IL)                      // CMPQ    AX, IL
    self.Sjmp("JAE"    , _LB_eof_error)                 // JAE     _eof_error
    self.Emit("MOVQ"   , _AX, _IC)                      // MOVQ    AX, IC
    self.Link(label)                           // _nospace_{n}:
}

func (self *_Assembler) _asm_OP_match_char(p *_Instr) {
    self.match_char(p.vb())
}

func (self *_Assembler) match_char(char byte) {
    self.check_eof(1)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(char)))  // CMPB (IP)(IC), ${p.vb()}
    self.Sjmp("JNE" , _LB_char_0_error)                                 // JNE  _char_0_error
    self.Emit("ADDQ", jit.Imm(1), _IC)                                  // ADDQ $1, IC
}

func (self *_Assembler) _asm_OP_check_char(p *_Instr) {
    self.check_eof(1)
    self.Emit("LEAQ"   , jit.Ptr(_IC, 1), _AX)                              // LEAQ    1(IC), AX
    self.Emit("CMPB"   , jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(p.vb())))   // CMPB    (IP)(IC), ${p.vb()}
    self.Emit("CMOVQEQ", _AX, _IC)                                          // CMOVQEQ AX, IC
    self.Xjmp("JE"     , p.vi())                                            // JE      {p.vi()}
}

func (self *_Assembler) _asm_OP_check_char_0(p *_Instr) {
    self.check_eof(1)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(p.vb())))   // CMPB    (IP)(IC), ${p.vb()}
    self.Xjmp("JE"  , p.vi())                                            // JE      {p.vi()}
}

func (self *_Assembler) _asm_OP_add(p *_Instr) {
    self.Emit("ADDQ", jit.Imm(int64(p.vi())), _IC)  // ADDQ ${p.vi()}, IC
}

func (self *_Assembler) _asm_OP_load(_ *_Instr) {
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             // MOVQ (ST), AX
    self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 0), _VP)     // MOVQ (ST)(AX), VP
}

func (self *_Assembler) _asm_OP_save(_ *_Instr) {
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _CX)             // MOVQ (ST), CX
    self.Emit("CMPQ", _CX, jit.Imm(_MaxStackBytes))     // CMPQ CX, ${_MaxStackBytes}
    self.Sjmp("JAE"  , _LB_stack_error)                  // JA   _stack_error
    self.WriteRecNotAX(0 , _VP, jit.Sib(_ST, _CX, 1, 8), false, false) // MOVQ VP, 8(ST)(CX)
    self.Emit("ADDQ", jit.Imm(8), _CX)                  // ADDQ $8, CX
    self.Emit("MOVQ", _CX, jit.Ptr(_ST, 0))             // MOVQ CX, (ST)
}

func (self *_Assembler) _asm_OP_drop(_ *_Instr) {
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             // MOVQ (ST), AX
    self.Emit("SUBQ", jit.Imm(8), _AX)                  // SUBQ $8, AX
    self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 8), _VP)     // MOVQ 8(ST)(AX), VP
    self.Emit("MOVQ", _AX, jit.Ptr(_ST, 0))             // MOVQ AX, (ST)
    self.Emit("XORL", _BX, _BX)                         // XORL BX, BX
    self.Emit("MOVQ", _BX, jit.Sib(_ST, _AX, 1, 8))     // MOVQ BX, 8(ST)(AX)
}

func (self *_Assembler) _asm_OP_drop_2(_ *_Instr) {
    self.Emit("MOVQ" , jit.Ptr(_ST, 0), _AX)            // MOVQ  (ST), AX
    self.Emit("SUBQ" , jit.Imm(16), _AX)                // SUBQ  $16, AX
    self.Emit("MOVQ" , jit.Sib(_ST, _AX, 1, 8), _VP)    // MOVQ  8(ST)(AX), VP
    self.Emit("MOVQ" , _AX, jit.Ptr(_ST, 0))            // MOVQ  AX, (ST)
    self.Emit("PXOR" , _X0, _X0)                        // PXOR  X0, X0
    self.Emit("MOVOU", _X0, jit.Sib(_ST, _AX, 1, 8))    // MOVOU X0, 8(ST)(AX)
}

func (self *_Assembler) _asm_OP_recurse(p *_Instr) {
    self.Emit("MOVQ", jit.Type(p.vt()), _AX)    // MOVQ   ${p.vt()}, AX
    self.decode_dynamic(_AX, _VP)               // DECODE AX, VP
}

func (self *_Assembler) _asm_OP_goto(p *_Instr) {
    self.Xjmp("JMP", p.vi())
}

func (self *_Assembler) _asm_OP_switch(p *_Instr) {
    self.Emit("MOVQ", _VAR_sr, _AX)             // MOVQ sr, AX
    self.Emit("CMPQ", _AX, jit.Imm(p.i64()))    // CMPQ AX, ${len(p.vs())}
    self.Sjmp("JAE" , "_default_{n}")           // JAE  _default_{n}

    /* jump table selector */
    self.Byte(0x48, 0x8d, 0x3d)                         // LEAQ    ?(PC), DI
    self.Sref("_switch_table_{n}", 4)                   // ....    &_switch_table_{n}
    self.Emit("MOVLQSX", jit.Sib(_DI, _AX, 4, 0), _AX)  // MOVLQSX (DI)(AX*4), AX
    self.Emit("ADDQ"   , _DI, _AX)                      // ADDQ    DI, AX
    self.Rjmp("JMP"    , _AX)                           // JMP     AX
    self.Link("_switch_table_{n}")                      // _switch_table_{n}:

    /* generate the jump table */
    for i, v := range p.vs() {
        self.Xref(v, int64(-i) * 4)
    }

    /* default case */
    self.Link("_default_{n}")
    self.NOP()
}

func (self *_Assembler) print_gc(i int, p1 *_Instr, p2 *_Instr) {
    self.Emit("MOVQ", jit.Imm(int64(p2.op())),  _CX)// MOVQ $(p2.op()), 16(SP)
    self.Emit("MOVQ", jit.Imm(int64(p1.op())),  _BX) // MOVQ $(p1.op()), 8(SP)
    self.Emit("MOVQ", jit.Imm(int64(i)),  _AX)       // MOVQ $(i), (SP)
    self.call_go(_F_println)
}
