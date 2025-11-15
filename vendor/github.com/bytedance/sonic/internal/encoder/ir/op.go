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

package ir

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/resolver"
	"github.com/bytedance/sonic/internal/rt"
)

type Op uint8

const (
	OP_null Op = iota + 1
	OP_empty_arr
	OP_empty_obj
	OP_bool
	OP_i8
	OP_i16
	OP_i32
	OP_i64
	OP_u8
	OP_u16
	OP_u32
	OP_u64
	OP_f32
	OP_f64
	OP_str
	OP_bin
	OP_quote
	OP_number
	OP_eface
	OP_iface
	OP_byte
	OP_text
	OP_deref
	OP_index
	OP_load
	OP_save
	OP_drop
	OP_drop_2
	OP_recurse
	OP_is_nil
	OP_is_nil_p1
	OP_is_zero_1
	OP_is_zero_2
	OP_is_zero_4
	OP_is_zero_8
	OP_is_zero_map
	OP_goto
	OP_map_iter
	OP_map_stop
	OP_map_check_key
	OP_map_write_key
	OP_map_value_next
	OP_slice_len
	OP_slice_next
	OP_marshal
	OP_marshal_p
	OP_marshal_text
	OP_marshal_text_p
	OP_cond_set
	OP_cond_testc
	OP_unsupported
	OP_is_zero
)

const (
	_INT_SIZE = 32 << (^uint(0) >> 63)
	_PTR_SIZE = 32 << (^uintptr(0) >> 63)
	_PTR_BYTE = unsafe.Sizeof(uintptr(0))
)

const OpSize = unsafe.Sizeof(NewInsOp(0))

var OpNames = [256]string{
	OP_null:           "null",
	OP_empty_arr:      "empty_arr",
	OP_empty_obj:      "empty_obj",
	OP_bool:           "bool",
	OP_i8:             "i8",
	OP_i16:            "i16",
	OP_i32:            "i32",
	OP_i64:            "i64",
	OP_u8:             "u8",
	OP_u16:            "u16",
	OP_u32:            "u32",
	OP_u64:            "u64",
	OP_f32:            "f32",
	OP_f64:            "f64",
	OP_str:            "str",
	OP_bin:            "bin",
	OP_quote:          "quote",
	OP_number:         "number",
	OP_eface:          "eface",
	OP_iface:          "iface",
	OP_byte:           "byte",
	OP_text:           "text",
	OP_deref:          "deref",
	OP_index:          "index",
	OP_load:           "load",
	OP_save:           "save",
	OP_drop:           "drop",
	OP_drop_2:         "drop_2",
	OP_recurse:        "recurse",
	OP_is_nil:         "is_nil",
	OP_is_nil_p1:      "is_nil_p1",
	OP_is_zero_1:      "is_zero_1",
	OP_is_zero_2:      "is_zero_2",
	OP_is_zero_4:      "is_zero_4",
	OP_is_zero_8:      "is_zero_8",
	OP_is_zero_map:    "is_zero_map",
	OP_goto:           "goto",
	OP_map_iter:       "map_iter",
	OP_map_stop:       "map_stop",
	OP_map_check_key:  "map_check_key",
	OP_map_write_key:  "map_write_key",
	OP_map_value_next: "map_value_next",
	OP_slice_len:      "slice_len",
	OP_slice_next:     "slice_next",
	OP_marshal:        "marshal",
	OP_marshal_p:      "marshal_p",
	OP_marshal_text:   "marshal_text",
	OP_marshal_text_p: "marshal_text_p",
	OP_cond_set:       "cond_set",
	OP_cond_testc:     "cond_testc",
	OP_unsupported:    "unsupported type",
}

func (self Op) String() string {
	if ret := OpNames[self]; ret != "" {
		return ret
	} else {
		return "<invalid>"
	}
}

func OP_int() Op {
	switch _INT_SIZE {
	case 32:
		return OP_i32
	case 64:
		return OP_i64
	default:
		panic("unsupported int size")
	}
}

func OP_uint() Op {
	switch _INT_SIZE {
	case 32:
		return OP_u32
	case 64:
		return OP_u64
	default:
		panic("unsupported uint size")
	}
}

func OP_uintptr() Op {
	switch _PTR_SIZE {
	case 32:
		return OP_u32
	case 64:
		return OP_u64
	default:
		panic("unsupported pointer size")
	}
}

func OP_is_zero_ints() Op {
	switch _INT_SIZE {
	case 32:
		return OP_is_zero_4
	case 64:
		return OP_is_zero_8
	default:
		panic("unsupported integer size")
	}
}

type Instr struct {
	o Op
	u int            // union {op: 8, _: 8, vi: 48}, vi maybe int or len(str)
	p unsafe.Pointer // maybe GoString.Ptr, or *GoType
}

func NewInsOp(op Op) Instr {
	return Instr{o: op}
}

func NewInsVi(op Op, vi int) Instr {
	return Instr{o: op, u: vi}
}

func NewInsVs(op Op, vs string) Instr {
	return Instr{
		o: op,
		u: len(vs),
		p: (*rt.GoString)(unsafe.Pointer(&vs)).Ptr,
	}
}

func NewInsVt(op Op, vt reflect.Type) Instr {
	return Instr{
		o: op,
		p: unsafe.Pointer(rt.UnpackType(vt)),
	}
}

type typAndTab struct {
	vt *rt.GoType
	itab *rt.GoItab
}

type typAndField struct {
	vt reflect.Type
	fv *resolver.FieldMeta
}

func NewInsVtab(op Op, vt reflect.Type, itab *rt.GoItab) Instr {
	return Instr{
		o: op,
		p: unsafe.Pointer(&typAndTab{
			vt: rt.UnpackType(vt),
			itab: itab,
		}),
	}
}

func NewInsField(op Op, fv *resolver.FieldMeta) Instr {
	return Instr{
		o: op,
		p: unsafe.Pointer(fv),
	}
}

func NewInsVp(op Op, vt reflect.Type, pv bool) Instr {
	i := 0
	if pv {
		i = 1
	}
	return Instr{
		o: op,
		u: i,
		p: unsafe.Pointer(rt.UnpackType(vt)),
	}
}

func (self Instr) Op() Op {
	return Op(self.o)
}

func (self Instr) Vi() int {
	return self.u
}

func (self Instr) Vf() uint8 {
	return (*rt.GoType)(self.p).KindFlags
}

func (self Instr) VField() (*resolver.FieldMeta) {
	return (*resolver.FieldMeta)(self.p)
}

func (self Instr) Vs() (v string) {
	(*rt.GoString)(unsafe.Pointer(&v)).Ptr = self.p
	(*rt.GoString)(unsafe.Pointer(&v)).Len = self.Vi()
	return
}

func (self Instr) Vk() reflect.Kind {
	return (*rt.GoType)(self.p).Kind()
}

func (self Instr) GoType() *rt.GoType {
	return (*rt.GoType)(self.p)
}

func (self Instr) Vt() reflect.Type {
	return (*rt.GoType)(self.p).Pack()
}

func (self Instr) Vr() *rt.GoType {
	return (*rt.GoType)(self.p)
}

func (self Instr) Vp() (vt reflect.Type, pv bool) {
	return (*rt.GoType)(self.p).Pack(), self.u == 1
}

func (self Instr) Vtab() (vt *rt.GoType, itab *rt.GoItab) {
	tt := (*typAndTab)(self.p)
	return tt.vt, tt.itab
}

func (self Instr) Vp2() (vt *rt.GoType, pv bool) {
	return (*rt.GoType)(self.p), self.u == 1
}

func (self Instr) I64() int64 {
	return int64(self.Vi())
}

func (self Instr) Byte() byte {
	return byte(self.Vi())
}

func (self Instr) Vlen() int {
	return int((*rt.GoType)(self.p).Size)
}

func (self Instr) isBranch() bool {
	switch self.Op() {
	case OP_goto:
		fallthrough
	case OP_is_nil:
		fallthrough
	case OP_is_nil_p1:
		fallthrough
	case OP_is_zero_1:
		fallthrough
	case OP_is_zero_2:
		fallthrough
	case OP_is_zero_4:
		fallthrough
	case OP_is_zero_8:
		fallthrough
	case OP_map_check_key:
		fallthrough
	case OP_map_write_key:
		fallthrough
	case OP_slice_next:
		fallthrough
	case OP_cond_testc:
		return true
	default:
		return false
	}
}

func (self Instr) Disassemble() string {
	switch self.Op() {
	case OP_byte:
		return fmt.Sprintf("%-18s%s", self.Op().String(), strconv.QuoteRune(rune(self.Vi())))
	case OP_text:
		return fmt.Sprintf("%-18s%s", self.Op().String(), strconv.Quote(self.Vs()))
	case OP_index:
		return fmt.Sprintf("%-18s%d", self.Op().String(), self.Vi())
	case OP_recurse:
		fallthrough
	case OP_map_iter:
		return fmt.Sprintf("%-18s%s", self.Op().String(), self.Vt())
	case OP_marshal:
		fallthrough
	case OP_marshal_p:
		fallthrough
	case OP_marshal_text:
		fallthrough
	case OP_marshal_text_p:
		vt, _ := self.Vtab()
		return fmt.Sprintf("%-18s%s", self.Op().String(), vt.Pack())
	case OP_goto:
		fallthrough
	case OP_is_nil:
		fallthrough
	case OP_is_nil_p1:
		fallthrough
	case OP_is_zero_1:
		fallthrough
	case OP_is_zero_2:
		fallthrough
	case OP_is_zero_4:
		fallthrough
	case OP_is_zero_8:
		fallthrough
	case OP_is_zero_map:
		fallthrough
	case OP_cond_testc:
		fallthrough
	case OP_map_check_key:
		fallthrough
	case OP_map_write_key:
		return fmt.Sprintf("%-18sL_%d", self.Op().String(), self.Vi())
	case OP_slice_next:
		return fmt.Sprintf("%-18sL_%d, %s", self.Op().String(), self.Vi(), self.Vt())
	default:
		return fmt.Sprintf("%#v", self) 
	}
}

type (
	Program []Instr
)

func (self Program) PC() int {
	return len(self)
}

func (self Program) Tag(n int) {
	if n >= vars.MaxStack {
		panic("type nesting too deep")
	}
}

func (self Program) Pin(i int) {
	v := &self[i]
	v.u = self.PC()
}

func (self Program) Rel(v []int) {
	for _, i := range v {
		self.Pin(i)
	}
}

func (self *Program) Add(op Op) {
	*self = append(*self, NewInsOp(op))
}

func (self *Program) Key(op Op) {
	*self = append(*self,
		NewInsVi(OP_byte, '"'),
		NewInsOp(op),
		NewInsVi(OP_byte, '"'),
	)
}

func (self *Program) Int(op Op, vi int) {
	*self = append(*self, NewInsVi(op, vi))
}

func (self *Program) Str(op Op, vs string) {
	*self = append(*self, NewInsVs(op, vs))
}

func (self *Program) Rtt(op Op, vt reflect.Type) {
	*self = append(*self, NewInsVt(op, vt))
}

func (self *Program) Vp(op Op, vt reflect.Type, pv bool) {
	*self = append(*self, NewInsVp(op, vt, pv))
}

func (self *Program) Vtab(op Op, vt reflect.Type, itab *rt.GoItab) {
	*self = append(*self, NewInsVtab(op, vt, itab))
}

func (self *Program) VField(op Op, fv *resolver.FieldMeta) {
	*self = append(*self, NewInsField(op, fv))
}

func (self Program) Disassemble() string {
	nb := len(self)
	tab := make([]bool, nb+1)
	ret := make([]string, 0, nb+1)

	/* prescan to get all the labels */
	for _, ins := range self {
		if ins.isBranch() {
			tab[ins.Vi()] = true
		}
	}

	/* disassemble each instruction */
	for i, ins := range self {
		if !tab[i] {
			ret = append(ret, "\t"+ins.Disassemble())
		} else {
			ret = append(ret, fmt.Sprintf("L_%d:\n\t%s", i, ins.Disassemble()))
		}
	}

	/* add the last label, if needed */
	if tab[nb] {
		ret = append(ret, fmt.Sprintf("L_%d:", nb))
	}

	/* add an "end" indicator, and join all the strings */
	return strings.Join(append(ret, "\tend"), "\n")
}
