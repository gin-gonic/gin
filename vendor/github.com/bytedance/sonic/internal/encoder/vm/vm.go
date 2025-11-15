// Copyright 2024 CloudWeGo Authors
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

package vm

import (
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/ir"
	"github.com/bytedance/sonic/internal/encoder/prim"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
)

const (
	_S_cond = iota
	_S_init
)

var (
	_T_json_Marshaler         = rt.UnpackType(vars.JsonMarshalerType)
	_T_encoding_TextMarshaler = rt.UnpackType(vars.EncodingTextMarshalerType)
)

func print_instr(buf []byte, pc int, op ir.Op, ins *ir.Instr, p unsafe.Pointer) {
	if len(buf) > 20 {
		fmt.Println(string(buf[len(buf)-20:]))
	} else {
		fmt.Println(string(buf))
	}
	fmt.Printf("pc %04d, op %v, ins %#v, ptr: %x\n", pc, op, ins.Disassemble(), p)
}

func Execute(b *[]byte, p unsafe.Pointer, s *vars.Stack, flags uint64, prog *ir.Program) (error) {
	pl := len(*prog)
	if pl <= 0 {
		return nil
	}

	var buf = *b
	var x int
	var q unsafe.Pointer
	var f uint64

	var pro = &(*prog)[0]
	for pc := 0; pc < pl; {
		ins := (*ir.Instr)(rt.Add(unsafe.Pointer(pro), ir.OpSize*uintptr(pc)))
		pc++
		op := ins.Op()

		switch op {
		case ir.OP_goto:
			pc = ins.Vi()
			continue
		case ir.OP_byte:
			v := ins.Byte()
			buf = append(buf, v)
		case ir.OP_text:
			v := ins.Vs()
			buf = append(buf, v...)
		case ir.OP_deref:
			p = *(*unsafe.Pointer)(p)
		case ir.OP_index:
			p = rt.Add(p, uintptr(ins.I64()))
		case ir.OP_load:
			// NOTICE: load CANNOT change f!
			x, _, p, q = s.Load() 
		case ir.OP_save:
			if !s.Save(x, f, p, q) {
				return vars.ERR_too_deep
			}
		case ir.OP_drop:
			x, f, p, q = s.Drop()
		case ir.OP_drop_2:
			s.Drop()
			x, f, p, q = s.Drop()
		case ir.OP_recurse:
			vt, pv := ins.Vp2()
			f := flags
			if pv {
				f |= (1 << alg.BitPointerValue)
			}
			*b = buf
			if vt.Indirect() {
				if err := EncodeTypedPointer(b, vt, (*unsafe.Pointer)(rt.NoEscape(unsafe.Pointer(&p))), s, f); err != nil {
					return err
				}
			} else {
				vp := (*unsafe.Pointer)(p)
				if err := EncodeTypedPointer(b, vt, vp, s, f); err != nil {
					return err
				}
			}
			buf = *b
		case ir.OP_is_nil:
			if is_nil(p) {
				pc = ins.Vi()
				continue
			}
		case ir.OP_is_nil_p1:
			if (*rt.GoEface)(p).Value == nil {
				pc = ins.Vi()
				continue
			}
		case ir.OP_null:
			buf = append(buf, 'n', 'u', 'l', 'l')
		case ir.OP_str:
			v := *(*string)(p)
			buf = alg.Quote(buf, v, false)
		case ir.OP_bool:
			if *(*bool)(p) {
				buf = append(buf, 't', 'r', 'u', 'e')
			} else {
				buf = append(buf, 'f', 'a', 'l', 's', 'e')
			}
		case ir.OP_i8:
			v := *(*int8)(p)
			buf = alg.I64toa(buf, int64(v))
		case ir.OP_i16:
			v := *(*int16)(p)
			buf = alg.I64toa(buf, int64(v))
		case ir.OP_i32:
			v := *(*int32)(p)
			buf = alg.I64toa(buf, int64(v))
		case ir.OP_i64:
			v := *(*int64)(p)
			buf = alg.I64toa(buf, int64(v))
		case ir.OP_u8:
			v := *(*uint8)(p)
			buf = alg.U64toa(buf, uint64(v))
		case ir.OP_u16:
			v := *(*uint16)(p)
			buf = alg.U64toa(buf, uint64(v))
		case ir.OP_u32:
			v := *(*uint32)(p)
			buf = alg.U64toa(buf, uint64(v))
		case ir.OP_u64:
			v := *(*uint64)(p)
			buf = alg.U64toa(buf, uint64(v))
		case ir.OP_f32:
			v := *(*float32)(p)
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				if flags&(1<<alg.BitEncodeNullForInfOrNan) != 0 {
					buf = append(buf, 'n', 'u', 'l', 'l')
					continue
				}
				return vars.ERR_nan_or_infinite
			}
			buf = alg.F32toa(buf, v)
		case ir.OP_f64:
			v := *(*float64)(p)
			if math.IsNaN(v) || math.IsInf(v, 0) {
				if flags&(1<<alg.BitEncodeNullForInfOrNan) != 0 {
					buf = append(buf, 'n', 'u', 'l', 'l')
					continue
				}
				return vars.ERR_nan_or_infinite
			}
			buf = alg.F64toa(buf, v)
		case ir.OP_bin:
			v := *(*[]byte)(p)
			buf = rt.EncodeBase64(buf, v)
		case ir.OP_quote:
			v := *(*string)(p)
			buf = alg.Quote(buf, v, true)
		case ir.OP_number:
			v := *(*json.Number)(p)
			if v == "" {
				buf = append(buf, '0')
			} else if !alg.IsValidNumber(string(v)) {
				return vars.Error_number(v)
			} else {
				buf = append(buf, v...)
			}
		case ir.OP_eface:
			*b = buf
			if err := EncodeTypedPointer(b, *(**rt.GoType)(p), (*unsafe.Pointer)(rt.Add(p, 8)), s, flags); err != nil {
				return err
			}
			buf = *b
		case ir.OP_iface:
			*b = buf
			if err := EncodeTypedPointer(b,  (*(**rt.GoItab)(p)).Vt, (*unsafe.Pointer)(rt.Add(p, 8)), s, flags); err != nil {
				return err
			}
			buf = *b
		case ir.OP_is_zero_map:
			v := *(*unsafe.Pointer)(p)
			if v == nil || rt.Maplen(v) == 0 {
				pc = ins.Vi()
				continue
			}
		case ir.OP_map_iter:
			v := *(*unsafe.Pointer)(p)
			vt := ins.Vr()
			it, err := alg.IteratorStart(rt.MapType(vt), v, flags)
			if err != nil {
				return err
			}
			q = unsafe.Pointer(it)
		case ir.OP_map_stop:
			it := (*alg.MapIterator)(q)
			alg.IteratorStop(it)
			q = nil
		case ir.OP_map_value_next:
			it := (*alg.MapIterator)(q)
			p = it.It.V
			alg.IteratorNext(it)
		case ir.OP_map_check_key:
			it := (*alg.MapIterator)(q)
			if it.It.K == nil {
				pc = ins.Vi()
				continue
			}
			p = it.It.K
		case ir.OP_marshal_text:
			vt, itab := ins.Vtab()
			var it rt.GoIface
			switch vt.Kind() {
				case reflect.Interface        : 
				if is_nil(p) {
					buf = append(buf, 'n', 'u', 'l', 'l')
					continue
				}
				it = rt.AssertI2I(_T_encoding_TextMarshaler, *(*rt.GoIface)(p))
				case reflect.Ptr, reflect.Map : it = convT2I(p, true, itab)
				default                       : it = convT2I(p, !vt.Indirect(), itab)
			}
			if err := prim.EncodeTextMarshaler(&buf, *(*encoding.TextMarshaler)(unsafe.Pointer(&it)), (flags)); err != nil {
				return err
			}
		case ir.OP_marshal_text_p:
			_, itab := ins.Vtab()
			it := convT2I(p, false, itab)
			if err := prim.EncodeTextMarshaler(&buf, *(*encoding.TextMarshaler)(unsafe.Pointer(&it)), (flags)); err != nil {
				return err
			}
		case ir.OP_map_write_key:
			if has_opts(flags, alg.BitSortMapKeys) {
				v := *(*string)(p)
				buf = alg.Quote(buf, v, false)
				pc = ins.Vi()
				continue
			}
		case ir.OP_slice_len:
			v := (*rt.GoSlice)(p)
			x = v.Len
			p = v.Ptr
			//TODO: why?
			f |= 1<<_S_init 
		case ir.OP_slice_next:
			if x == 0 {
				pc = ins.Vi()
				continue
			}
			x--
			if has_opts(f, _S_init) {
				f &= ^uint64(1 << _S_init)
			} else {
				p = rt.Add(p, uintptr(ins.Vlen()))
			}
		case ir.OP_cond_set:
			f |= 1<<_S_cond
		case ir.OP_cond_testc:
			if has_opts(f, _S_cond) {
				f &= ^uint64(1 << _S_cond)
				pc = ins.Vi()
				continue
			}
		case ir.OP_is_zero:
			fv := ins.VField()
			if prim.IsZero(p, fv) {
				pc = ins.Vi()
				continue
			}
		case ir.OP_is_zero_1:
			if *(*uint8)(p) == 0 {
				pc = ins.Vi()
				continue
			}
		case ir.OP_is_zero_2:
			if *(*uint16)(p) == 0 {
				pc = ins.Vi()
				continue
			}
		case ir.OP_is_zero_4:
			if *(*uint32)(p) == 0 {
				pc = ins.Vi()
				continue
			}
		case ir.OP_is_zero_8:
			if *(*uint64)(p) == 0 {
				pc = ins.Vi()
				continue
			}
		case ir.OP_empty_arr:
			if has_opts(flags, alg.BitNoNullSliceOrMap) {
				buf = append(buf, '[', ']')
			} else {
				buf = append(buf, 'n', 'u', 'l', 'l')
			}
		case ir.OP_empty_obj:
			if has_opts(flags, alg.BitNoNullSliceOrMap) {
				buf = append(buf, '{', '}')
			} else {
				buf = append(buf, 'n', 'u', 'l', 'l')
			}
		case ir.OP_marshal:
			vt, itab := ins.Vtab()
			var it rt.GoIface
			switch vt.Kind() {
				case reflect.Interface        : 
				if is_nil(p) {
					buf = append(buf, 'n', 'u', 'l', 'l')
					continue
				}
				it = rt.AssertI2I(_T_json_Marshaler, *(*rt.GoIface)(p))
				case reflect.Ptr, reflect.Map : it = convT2I(p, true, itab)
				default                       : it = convT2I(p, !vt.Indirect(), itab)
			}
			if err := prim.EncodeJsonMarshaler(&buf, *(*json.Marshaler)(unsafe.Pointer(&it)), (flags)); err != nil {
				return err
			}
		case ir.OP_marshal_p:
			_, itab := ins.Vtab()
			it := convT2I(p, false, itab)
			if err := prim.EncodeJsonMarshaler(&buf, *(*json.Marshaler)(unsafe.Pointer(&it)), (flags)); err != nil {
				return err
			}
		case ir.OP_unsupported:
			return vars.Error_unsuppoted(ins.GoType())
		default:
			panic(fmt.Sprintf("not implement %s at %d", ins.Op().String(), pc))
		}
	}

	*b = buf
	return nil
}


func has_opts(opts uint64, bit int) bool {
	return opts & (1<<bit) != 0
}

func is_nil(p unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(p) == nil
}

func convT2I(ptr unsafe.Pointer, deref bool, itab *rt.GoItab) (rt.GoIface) {
	if deref {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	return rt.GoIface{
		Itab:  itab,
		Value: ptr,
	}
}
