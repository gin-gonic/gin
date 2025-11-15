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
    `encoding/json`
    `fmt`
    `reflect`
    `sort`
    `strconv`
    `strings`
    `unsafe`

    `github.com/bytedance/sonic/internal/caching`
    `github.com/bytedance/sonic/internal/resolver`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/option`
)

type _Op uint8

const (
    _OP_any _Op = iota + 1
    _OP_dyn
    _OP_str
    _OP_bin
    _OP_bool
    _OP_num
    _OP_i8
    _OP_i16
    _OP_i32
    _OP_i64
    _OP_u8
    _OP_u16
    _OP_u32
    _OP_u64
    _OP_f32
    _OP_f64
    _OP_unquote
    _OP_nil_1
    _OP_nil_2
    _OP_nil_3
    _OP_empty_bytes
    _OP_deref
    _OP_index
    _OP_is_null
    _OP_is_null_quote
    _OP_map_init
    _OP_map_key_i8
    _OP_map_key_i16
    _OP_map_key_i32
    _OP_map_key_i64
    _OP_map_key_u8
    _OP_map_key_u16
    _OP_map_key_u32
    _OP_map_key_u64
    _OP_map_key_f32
    _OP_map_key_f64
    _OP_map_key_str
    _OP_map_key_utext
    _OP_map_key_utext_p
    _OP_array_skip
    _OP_array_clear
    _OP_array_clear_p
    _OP_slice_init
    _OP_slice_append
    _OP_object_next
    _OP_struct_field
    _OP_unmarshal
    _OP_unmarshal_p
    _OP_unmarshal_text
    _OP_unmarshal_text_p
    _OP_lspace
    _OP_match_char
    _OP_check_char
    _OP_load
    _OP_save
    _OP_drop
    _OP_drop_2
    _OP_recurse
    _OP_goto
    _OP_switch
    _OP_check_char_0
    _OP_dismatch_err
    _OP_go_skip
    _OP_skip_emtpy
    _OP_add
    _OP_check_empty
    _OP_unsupported
    _OP_debug
)

const (
    _INT_SIZE = 32 << (^uint(0) >> 63)
    _PTR_SIZE = 32 << (^uintptr(0) >> 63)
    _PTR_BYTE = unsafe.Sizeof(uintptr(0))
)

const (
    _MAX_ILBUF = 100000     // cutoff at 100k of IL instructions
    _MAX_FIELDS = 50        // cutoff at 50 fields struct
)

var _OpNames = [256]string {
    _OP_any              : "any",
    _OP_dyn              : "dyn",
    _OP_str              : "str",
    _OP_bin              : "bin",
    _OP_bool             : "bool",
    _OP_num              : "num",
    _OP_i8               : "i8",
    _OP_i16              : "i16",
    _OP_i32              : "i32",
    _OP_i64              : "i64",
    _OP_u8               : "u8",
    _OP_u16              : "u16",
    _OP_u32              : "u32",
    _OP_u64              : "u64",
    _OP_f32              : "f32",
    _OP_f64              : "f64",
    _OP_unquote          : "unquote",
    _OP_nil_1            : "nil_1",
    _OP_nil_2            : "nil_2",
    _OP_nil_3            : "nil_3",
    _OP_empty_bytes      : "empty bytes",
    _OP_deref            : "deref",
    _OP_index            : "index",
    _OP_is_null          : "is_null",
    _OP_is_null_quote    : "is_null_quote",
    _OP_map_init         : "map_init",
    _OP_map_key_i8       : "map_key_i8",
    _OP_map_key_i16      : "map_key_i16",
    _OP_map_key_i32      : "map_key_i32",
    _OP_map_key_i64      : "map_key_i64",
    _OP_map_key_u8       : "map_key_u8",
    _OP_map_key_u16      : "map_key_u16",
    _OP_map_key_u32      : "map_key_u32",
    _OP_map_key_u64      : "map_key_u64",
    _OP_map_key_f32      : "map_key_f32",
    _OP_map_key_f64      : "map_key_f64",
    _OP_map_key_str      : "map_key_str",
    _OP_map_key_utext    : "map_key_utext",
    _OP_map_key_utext_p  : "map_key_utext_p",
    _OP_array_skip       : "array_skip",
    _OP_slice_init       : "slice_init",
    _OP_slice_append     : "slice_append",
    _OP_object_next      : "object_next",
    _OP_struct_field     : "struct_field",
    _OP_unmarshal        : "unmarshal",
    _OP_unmarshal_p      : "unmarshal_p",
    _OP_unmarshal_text   : "unmarshal_text",
    _OP_unmarshal_text_p : "unmarshal_text_p",
    _OP_lspace           : "lspace",
    _OP_match_char       : "match_char",
    _OP_check_char       : "check_char",
    _OP_load             : "load",
    _OP_save             : "save",
    _OP_drop             : "drop",
    _OP_drop_2           : "drop_2",
    _OP_recurse          : "recurse",
    _OP_goto             : "goto",
    _OP_switch           : "switch",
    _OP_check_char_0     : "check_char_0",
    _OP_dismatch_err     : "dismatch_err",
    _OP_add              : "add",
    _OP_go_skip          : "go_skip",
    _OP_check_empty      : "check_empty",
    _OP_unsupported      : "unsupported type",
    _OP_debug            : "debug",
}

func (self _Op) String() string {
    if ret := _OpNames[self]; ret != "" {
        return ret
    } else {
        return "<invalid>"
    }
}

func _OP_int() _Op {
    switch _INT_SIZE {
        case 32: return _OP_i32
        case 64: return _OP_i64
        default: panic("unsupported int size")
    }
}

func _OP_uint() _Op {
    switch _INT_SIZE {
        case 32: return _OP_u32
        case 64: return _OP_u64
        default: panic("unsupported uint size")
    }
}

func _OP_uintptr() _Op {
    switch _PTR_SIZE {
        case 32: return _OP_u32
        case 64: return _OP_u64
        default: panic("unsupported pointer size")
    }
}

func _OP_map_key_int() _Op {
    switch _INT_SIZE {
        case 32: return _OP_map_key_i32
        case 64: return _OP_map_key_i64
        default: panic("unsupported int size")
    }
}

func _OP_map_key_uint() _Op {
    switch _INT_SIZE {
        case 32: return _OP_map_key_u32
        case 64: return _OP_map_key_u64
        default: panic("unsupported uint size")
    }
}

func _OP_map_key_uintptr() _Op {
    switch _PTR_SIZE {
        case 32: return _OP_map_key_u32
        case 64: return _OP_map_key_u64
        default: panic("unsupported pointer size")
    }
}

type _Instr struct {
    u uint64            // union {op: 8, vb: 8, vi: 48}, iv maybe int or len([]int)
    p unsafe.Pointer    // maybe GoSlice.Data, *GoType or *caching.FieldMap
}

func packOp(op _Op) uint64 {
    return uint64(op) << 56
}

func newInsOp(op _Op) _Instr {
    return _Instr{u: packOp(op)}
}

func newInsVi(op _Op, vi int) _Instr {
    return _Instr{u: packOp(op) | rt.PackInt(vi)}
}

func newInsVb(op _Op, vb byte) _Instr {
    return _Instr{u: packOp(op) | (uint64(vb) << 48)}
}

func newInsVs(op _Op, vs []int) _Instr {
    return _Instr {
        u: packOp(op) | rt.PackInt(len(vs)),
        p: (*rt.GoSlice)(unsafe.Pointer(&vs)).Ptr,
    }
}

func newInsVt(op _Op, vt reflect.Type) _Instr {
    return _Instr {
        u: packOp(op),
        p: unsafe.Pointer(rt.UnpackType(vt)),
    }
}

func newInsVtI(op _Op, vt reflect.Type, iv int) _Instr {
    return _Instr {
        u: packOp(op) | rt.PackInt(iv),
        p: unsafe.Pointer(rt.UnpackType(vt)),
    }
}

func newInsVf(op _Op, vf *caching.FieldMap) _Instr {
    return _Instr {
        u: packOp(op),
        p: unsafe.Pointer(vf),
    }
}

func (self _Instr) op() _Op {
    return _Op(self.u >> 56)
}

func (self _Instr) vi() int {
    return rt.UnpackInt(self.u)
}

func (self _Instr) vb() byte {
    return byte(self.u >> 48)
}

func (self _Instr) vs() (v []int) {
    (*rt.GoSlice)(unsafe.Pointer(&v)).Ptr = self.p
    (*rt.GoSlice)(unsafe.Pointer(&v)).Cap = self.vi()
    (*rt.GoSlice)(unsafe.Pointer(&v)).Len = self.vi()
    return
}

func (self _Instr) vf() *caching.FieldMap {
    return (*caching.FieldMap)(self.p)
}

func (self _Instr) vk() reflect.Kind {
    return (*rt.GoType)(self.p).Kind()
}

func (self _Instr) vt() reflect.Type {
    return (*rt.GoType)(self.p).Pack()
}

func (self _Instr) i64() int64 {
    return int64(self.vi())
}

func (self _Instr) vlen() int {
    return int((*rt.GoType)(self.p).Size)
}

func (self _Instr) isBranch() bool {
    switch self.op() {
        case _OP_goto          : fallthrough
        case _OP_switch        : fallthrough
        case _OP_is_null       : fallthrough
        case _OP_is_null_quote : fallthrough
        case _OP_check_char    : return true
        default                : return false
    }
}

func (self _Instr) disassemble() string {
    switch self.op() {
        case _OP_dyn              : fallthrough
        case _OP_deref            : fallthrough
        case _OP_map_key_i8       : fallthrough
        case _OP_map_key_i16      : fallthrough
        case _OP_map_key_i32      : fallthrough
        case _OP_map_key_i64      : fallthrough
        case _OP_map_key_u8       : fallthrough
        case _OP_map_key_u16      : fallthrough
        case _OP_map_key_u32      : fallthrough
        case _OP_map_key_u64      : fallthrough
        case _OP_map_key_f32      : fallthrough
        case _OP_map_key_f64      : fallthrough
        case _OP_map_key_str      : fallthrough
        case _OP_map_key_utext    : fallthrough
        case _OP_map_key_utext_p  : fallthrough
        case _OP_slice_init       : fallthrough
        case _OP_slice_append     : fallthrough
        case _OP_unmarshal        : fallthrough
        case _OP_unmarshal_p      : fallthrough
        case _OP_unmarshal_text   : fallthrough
        case _OP_unmarshal_text_p : fallthrough
        case _OP_recurse          : return fmt.Sprintf("%-18s%s", self.op(), self.vt())
        case _OP_goto             : fallthrough
        case _OP_is_null_quote    : fallthrough
        case _OP_is_null          : return fmt.Sprintf("%-18sL_%d", self.op(), self.vi())
        case _OP_index            : fallthrough
        case _OP_array_clear      : fallthrough
        case _OP_array_clear_p    : return fmt.Sprintf("%-18s%d", self.op(), self.vi())
        case _OP_switch           : return fmt.Sprintf("%-18s%s", self.op(), self.formatSwitchLabels())
        case _OP_struct_field     : return fmt.Sprintf("%-18s%s", self.op(), self.formatStructFields())
        case _OP_match_char       : return fmt.Sprintf("%-18s%s", self.op(), strconv.QuoteRune(rune(self.vb())))
        case _OP_check_char       : return fmt.Sprintf("%-18sL_%d, %s", self.op(), self.vi(), strconv.QuoteRune(rune(self.vb())))
        default                   : return self.op().String()
    }
}

func (self _Instr) formatSwitchLabels() string {
    var i int
    var v int
    var m []string

    /* format each label */
    for i, v = range self.vs() {
        m = append(m, fmt.Sprintf("%d=L_%d", i, v))
    }

    /* join them with "," */
    return strings.Join(m, ", ")
}

func (self _Instr) formatStructFields() string {
    var i uint64
    var r []string
    var m []struct{i int; n string}

    /* extract all the fields */
    for i = 0; i < self.vf().N; i++ {
        if v := self.vf().At(i); v.Hash != 0 {
            m = append(m, struct{i int; n string}{i: v.ID, n: v.Name})
        }
    }

    /* sort by field name */
    sort.Slice(m, func(i, j int) bool {
        return m[i].n < m[j].n
    })

    /* format each field */
    for _, v := range m {
        r = append(r, fmt.Sprintf("%s=%d", v.n, v.i))
    }

    /* join them with "," */
    return strings.Join(r, ", ")
}

type (
    _Program []_Instr
)

func (self _Program) pc() int {
    return len(self)
}

func (self _Program) tag(n int) {
    if n >= _MaxStack {
        panic("type nesting too deep")
    }
}

func (self _Program) pin(i int) {
    v := &self[i]
    v.u &= 0xffff000000000000
    v.u |= rt.PackInt(self.pc())
}

func (self _Program) rel(v []int) {
    for _, i := range v {
        self.pin(i)
    }
}

func (self *_Program) add(op _Op) {
    *self = append(*self, newInsOp(op))
}

func (self *_Program) int(op _Op, vi int) {
    *self = append(*self, newInsVi(op, vi))
}

func (self *_Program) chr(op _Op, vb byte) {
    *self = append(*self, newInsVb(op, vb))
}

func (self *_Program) tab(op _Op, vs []int) {
    *self = append(*self, newInsVs(op, vs))
}

func (self *_Program) rtt(op _Op, vt reflect.Type) {
    *self = append(*self, newInsVt(op, vt))
}

func (self *_Program) rtti(op _Op, vt reflect.Type, iv int) {
    *self = append(*self, newInsVtI(op, vt, iv))
}

func (self *_Program) fmv(op _Op, vf *caching.FieldMap) {
    *self = append(*self, newInsVf(op, vf))
}

func (self _Program) disassemble() string {
    nb  := len(self)
    tab := make([]bool, nb + 1)
    ret := make([]string, 0, nb + 1)

    /* prescan to get all the labels */
    for _, ins := range self {
        if ins.isBranch() {
            if ins.op() != _OP_switch {
                tab[ins.vi()] = true
            } else {
                for _, v := range ins.vs() {
                    tab[v] = true
                }
            }
        }
    }

    /* disassemble each instruction */
    for i, ins := range self {
        if !tab[i] {
            ret = append(ret, "\t" + ins.disassemble())
        } else {
            ret = append(ret, fmt.Sprintf("L_%d:\n\t%s", i, ins.disassemble()))
        }
    }

    /* add the last label, if needed */
    if tab[nb] {
        ret = append(ret, fmt.Sprintf("L_%d:", nb))
    }

    /* add an "end" indicator, and join all the strings */
    return strings.Join(append(ret, "\tend"), "\n")
}

type _Compiler struct {
    opts option.CompileOptions
    tab  map[reflect.Type]bool
    rec  map[reflect.Type]bool
}

func newCompiler() *_Compiler {
    return &_Compiler {
        opts: option.DefaultCompileOptions(),
        tab: map[reflect.Type]bool{},
        rec: map[reflect.Type]bool{},
    }
}

func (self *_Compiler) apply(opts option.CompileOptions) *_Compiler {
    self.opts = opts
    return self
}

func (self *_Compiler) rescue(ep *error) {
    if val := recover(); val != nil {
        if err, ok := val.(error); ok {
            *ep = err
        } else {
            panic(val)
        }
    }
}

func (self *_Compiler) compile(vt reflect.Type) (ret _Program, err error) {
    defer self.rescue(&err)
    self.compileOne(&ret, 0, vt)
    return
}

const (
    checkMarshalerFlags_quoted = 1
)

func (self *_Compiler) checkMarshaler(p *_Program, vt reflect.Type, flags int, exec bool) bool {
    pt := reflect.PtrTo(vt)

    /* check for `json.Unmarshaler` with pointer receiver */
    if pt.Implements(jsonUnmarshalerType) {
        if exec {
            p.add(_OP_lspace)
            p.rtti(_OP_unmarshal_p, pt, flags)
        }
        return true
    }

    /* check for `json.Unmarshaler` */
    if vt.Implements(jsonUnmarshalerType) {
        if exec {
            p.add(_OP_lspace)
            self.compileUnmarshalJson(p, vt, flags)
        }
        return true
    }

    if flags == checkMarshalerFlags_quoted {
        // text marshaler shouldn't be supported for quoted string
        return false
    }

    /* check for `encoding.TextMarshaler` with pointer receiver */
    if pt.Implements(encodingTextUnmarshalerType) {
        if exec {
            p.add(_OP_lspace)
            self.compileUnmarshalTextPtr(p, pt, flags)
        }
        return true
    }

    /* check for `encoding.TextUnmarshaler` */
    if vt.Implements(encodingTextUnmarshalerType) {
        if exec {
            p.add(_OP_lspace)
            self.compileUnmarshalText(p, vt, flags)
        }
        return true
    }

    return false
}

func (self *_Compiler) compileOne(p *_Program, sp int, vt reflect.Type) {
    /* check for recursive nesting */
    ok := self.tab[vt]
    if ok {
        p.rtt(_OP_recurse, vt)
        return
    }

    if self.checkMarshaler(p, vt, 0, true) {
        return
    }

    /* enter the recursion */
    p.add(_OP_lspace)
    self.tab[vt] = true
    self.compileOps(p, sp, vt)
    delete(self.tab, vt)
}

func (self *_Compiler) compileOps(p *_Program, sp int, vt reflect.Type) {
    switch vt.Kind() {
        case reflect.Bool      : self.compilePrimitive (vt, p, _OP_bool)
        case reflect.Int       : self.compilePrimitive (vt, p, _OP_int())
        case reflect.Int8      : self.compilePrimitive (vt, p, _OP_i8)
        case reflect.Int16     : self.compilePrimitive (vt, p, _OP_i16)
        case reflect.Int32     : self.compilePrimitive (vt, p, _OP_i32)
        case reflect.Int64     : self.compilePrimitive (vt, p, _OP_i64)
        case reflect.Uint      : self.compilePrimitive (vt, p, _OP_uint())
        case reflect.Uint8     : self.compilePrimitive (vt, p, _OP_u8)
        case reflect.Uint16    : self.compilePrimitive (vt, p, _OP_u16)
        case reflect.Uint32    : self.compilePrimitive (vt, p, _OP_u32)
        case reflect.Uint64    : self.compilePrimitive (vt, p, _OP_u64)
        case reflect.Uintptr   : self.compilePrimitive (vt, p, _OP_uintptr())
        case reflect.Float32   : self.compilePrimitive (vt, p, _OP_f32)
        case reflect.Float64   : self.compilePrimitive (vt, p, _OP_f64)
        case reflect.String    : self.compileString    (p, vt)
        case reflect.Array     : self.compileArray     (p, sp, vt)
        case reflect.Interface : self.compileInterface (p, vt)
        case reflect.Map       : self.compileMap       (p, sp, vt)
        case reflect.Ptr       : self.compilePtr       (p, sp, vt)
        case reflect.Slice     : self.compileSlice     (p, sp, vt)
        case reflect.Struct    : self.compileStruct    (p, sp, vt)
        default                : self.compileUnsupportedType      (p, vt)
    }
}

func (self *_Compiler) compileUnsupportedType(p *_Program, vt reflect.Type) {
    i := p.pc()
    p.add(_OP_is_null)
    p.rtt(_OP_unsupported, vt)
    p.pin(i)
}

func (self *_Compiler) compileMap(p *_Program, sp int, vt reflect.Type) {
    if vt.Key().Kind() != reflect.Interface && reflect.PtrTo(vt.Key()).Implements(encodingTextUnmarshalerType) {
        self.compileMapOp(p, sp, vt, _OP_map_key_utext_p)
    } else if vt.Key().Kind() != reflect.Interface && vt.Key().Implements(encodingTextUnmarshalerType) {
        self.compileMapOp(p, sp, vt, _OP_map_key_utext)
    } else {
        self.compileMapUt(p, sp, vt)
    }
}

func (self *_Compiler) compileMapUt(p *_Program, sp int, vt reflect.Type) {
    switch vt.Key().Kind() {
        case reflect.Int     : self.compileMapOp(p, sp, vt, _OP_map_key_int())
        case reflect.Int8    : self.compileMapOp(p, sp, vt, _OP_map_key_i8)
        case reflect.Int16   : self.compileMapOp(p, sp, vt, _OP_map_key_i16)
        case reflect.Int32   : self.compileMapOp(p, sp, vt, _OP_map_key_i32)
        case reflect.Int64   : self.compileMapOp(p, sp, vt, _OP_map_key_i64)
        case reflect.Uint    : self.compileMapOp(p, sp, vt, _OP_map_key_uint())
        case reflect.Uint8   : self.compileMapOp(p, sp, vt, _OP_map_key_u8)
        case reflect.Uint16  : self.compileMapOp(p, sp, vt, _OP_map_key_u16)
        case reflect.Uint32  : self.compileMapOp(p, sp, vt, _OP_map_key_u32)
        case reflect.Uint64  : self.compileMapOp(p, sp, vt, _OP_map_key_u64)
        case reflect.Uintptr : self.compileMapOp(p, sp, vt, _OP_map_key_uintptr())
        case reflect.Float32 : self.compileMapOp(p, sp, vt, _OP_map_key_f32)
        case reflect.Float64 : self.compileMapOp(p, sp, vt, _OP_map_key_f64)
        case reflect.String  : self.compileMapOp(p, sp, vt, _OP_map_key_str)
        default              : panic(&json.UnmarshalTypeError{Type: vt})
    }
}

func (self *_Compiler) compileMapOp(p *_Program, sp int, vt reflect.Type, op _Op) {
    i := p.pc()
    p.add(_OP_is_null)
    p.tag(sp + 1)
    skip := self.checkIfSkip(p, vt, '{')
    p.add(_OP_save)
    p.add(_OP_map_init)
    p.add(_OP_save)
    p.add(_OP_lspace)
    j := p.pc()
    p.chr(_OP_check_char, '}')
    p.chr(_OP_match_char, '"')
    skip2 := p.pc()
    p.rtt(op, vt)

    /* match the value separator */
    p.add(_OP_lspace)
    p.chr(_OP_match_char, ':')
    self.compileOne(p, sp + 2, vt.Elem())
    p.pin(skip2)
    p.add(_OP_load)
    k0 := p.pc()
    p.add(_OP_lspace)
    k1 := p.pc()
    p.chr(_OP_check_char, '}')
    p.chr(_OP_match_char, ',')
    p.add(_OP_lspace)
    p.chr(_OP_match_char, '"')
    skip3 := p.pc()
    p.rtt(op, vt)

    /* match the value separator */
    p.add(_OP_lspace)
    p.chr(_OP_match_char, ':')
    self.compileOne(p, sp + 2, vt.Elem())
    p.pin(skip3)
    p.add(_OP_load)
    p.int(_OP_goto, k0)
    p.pin(j)
    p.pin(k1)
    p.add(_OP_drop_2)
    x := p.pc()
    p.add(_OP_goto)
    p.pin(i)
    p.add(_OP_nil_1)
    p.pin(skip)
    p.pin(x)
}

func (self *_Compiler) compilePtr(p *_Program, sp int, et reflect.Type) {
    i := p.pc()
    p.add(_OP_is_null)

    /* dereference all the way down */
    for et.Kind() == reflect.Ptr {
        if self.checkMarshaler(p, et, 0, true) {
            return
        }
        et = et.Elem()
        p.rtt(_OP_deref, et)
    }

    /* check for recursive nesting */
    ok := self.tab[et]
    if ok {
        p.rtt(_OP_recurse, et)
    } else {
        /* enter the recursion */
        p.add(_OP_lspace)
        self.tab[et] = true

        /* not inline the pointer type
        * recursing the defined pointer type's elem will cause issue379.
        */
        self.compileOps(p, sp, et)
    }
    delete(self.tab, et)

    j := p.pc()
    p.add(_OP_goto)

    // set val pointer as nil
    p.pin(i)
    p.add(_OP_nil_1)

    // nothing todo
    p.pin(j)
}

func (self *_Compiler) compileArray(p *_Program, sp int, vt reflect.Type) {
    x := p.pc()
    p.add(_OP_is_null)
    p.tag(sp)
    skip := self.checkIfSkip(p, vt, '[')
    
    p.add(_OP_save)
    p.add(_OP_lspace)
    v := []int{p.pc()}
    p.chr(_OP_check_char, ']')

    /* decode every item */
    for i := 1; i <= vt.Len(); i++ {
        self.compileOne(p, sp + 1, vt.Elem())
        p.add(_OP_load)
        p.int(_OP_index, i * int(vt.Elem().Size()))
        p.add(_OP_lspace)
        v = append(v, p.pc())
        p.chr(_OP_check_char, ']')
        p.chr(_OP_match_char, ',')
    }

    /* drop rest of the array */
    p.add(_OP_array_skip)
    w := p.pc()
    p.add(_OP_goto)
    p.rel(v)

    /* check for pointer data */
    if rt.UnpackType(vt.Elem()).PtrData == 0 {
        p.int(_OP_array_clear, int(vt.Size()))
    } else {
        p.int(_OP_array_clear_p, int(vt.Size()))
    }

    /* restore the stack */
    p.pin(w)
    p.add(_OP_drop)

    p.pin(skip)
    p.pin(x)
}

func (self *_Compiler) compileSlice(p *_Program, sp int, vt reflect.Type) {
    if vt.Elem().Kind() == byteType.Kind() {
        self.compileSliceBin(p, sp, vt)
    } else {
        self.compileSliceList(p, sp, vt)
    }
}

func (self *_Compiler) compileSliceBin(p *_Program, sp int, vt reflect.Type) {
    i := p.pc()
    p.add(_OP_is_null)
    j := p.pc()
    p.chr(_OP_check_char, '[')
    skip := self.checkIfSkip(p, vt, '"')
    k := p.pc()
    p.chr(_OP_check_char, '"')
    p.add(_OP_bin)
    x := p.pc()
    p.add(_OP_goto)
    p.pin(j)
    self.compileSliceBody(p, sp, vt.Elem())
    y := p.pc()
    p.add(_OP_goto)

    // unmarshal `null` and `"` is different
    p.pin(i)
    p.add(_OP_nil_3)
    y2 := p.pc()
    p.add(_OP_goto)

    p.pin(k)
    p.add(_OP_empty_bytes)
    p.pin(x)
    p.pin(skip)
    p.pin(y)
    p.pin(y2)
}

func (self *_Compiler) compileSliceList(p *_Program, sp int, vt reflect.Type) {
    i := p.pc()
    p.add(_OP_is_null)
    p.tag(sp)
    skip := self.checkIfSkip(p, vt, '[')
    self.compileSliceBody(p, sp, vt.Elem())
    x := p.pc()
    p.add(_OP_goto)
    p.pin(i)
    p.add(_OP_nil_3)
    p.pin(x)
    p.pin(skip)
}

func (self *_Compiler) compileSliceBody(p *_Program, sp int, et reflect.Type) {
    p.add(_OP_lspace)
    j := p.pc()
    p.chr(_OP_check_empty, ']')
    p.rtt(_OP_slice_init, et)
    p.add(_OP_save)
    p.rtt(_OP_slice_append, et)
    self.compileOne(p, sp + 1, et)
    p.add(_OP_load)
    k0 := p.pc()
    p.add(_OP_lspace)
    k1 := p.pc()
    p.chr(_OP_check_char, ']')
    p.chr(_OP_match_char, ',')
    p.rtt(_OP_slice_append, et)
    self.compileOne(p, sp + 1, et)
    p.add(_OP_load)
    p.int(_OP_goto, k0)
    p.pin(k1)
    p.add(_OP_drop)
    p.pin(j)
}

func (self *_Compiler) compileString(p *_Program, vt reflect.Type) {
    if vt == jsonNumberType {
        self.compilePrimitive(vt, p, _OP_num)
    } else {
        self.compileStringBody(vt, p)
    }
}

func (self *_Compiler) compileStringBody(vt reflect.Type, p *_Program) {
    i := p.pc()
    p.add(_OP_is_null)
    skip := self.checkIfSkip(p, vt, '"')
    p.add(_OP_str)
    p.pin(i)
    p.pin(skip)
}

func (self *_Compiler) compileStruct(p *_Program, sp int, vt reflect.Type) {
    if sp >= self.opts.MaxInlineDepth || p.pc() >= _MAX_ILBUF || (sp > 0 && vt.NumField() >= _MAX_FIELDS) {
        p.rtt(_OP_recurse, vt)
        if self.opts.RecursiveDepth > 0 {
            self.rec[vt] = true
        }
    } else {
        self.compileStructBody(p, sp, vt)
    }
}

func (self *_Compiler) compileStructBody(p *_Program, sp int, vt reflect.Type) {
    fv := resolver.ResolveStruct(vt)
    fm, sw := caching.CreateFieldMap(len(fv)), make([]int, len(fv))

    /* start of object */
    p.tag(sp)
    n := p.pc()
    p.add(_OP_is_null)

    j := p.pc()
    p.chr(_OP_check_char_0, '{')
    p.rtt(_OP_dismatch_err, vt)

    /* special case for empty object */
    if len(fv) == 0 {
        p.pin(j)
        s := p.pc()
        p.add(_OP_skip_emtpy)
        p.pin(s)
        p.pin(n)
        return
    }

    skip := p.pc()
    p.add(_OP_go_skip)
    p.pin(j)
    p.int(_OP_add, 1)
    
    p.add(_OP_save)
    p.add(_OP_lspace)
    x := p.pc()
    p.chr(_OP_check_char, '}')
    p.chr(_OP_match_char, '"')
    p.fmv(_OP_struct_field, fm)
    p.add(_OP_lspace)
    p.chr(_OP_match_char, ':')
    p.tab(_OP_switch, sw)
    p.add(_OP_object_next)
    y0 := p.pc()
    p.add(_OP_lspace)
    y1 := p.pc()
    p.chr(_OP_check_char, '}')
    p.chr(_OP_match_char, ',')


    /* match the remaining fields */
    p.add(_OP_lspace)
    p.chr(_OP_match_char, '"')
    p.fmv(_OP_struct_field, fm)
    p.add(_OP_lspace)
    p.chr(_OP_match_char, ':')
    p.tab(_OP_switch, sw)
    p.add(_OP_object_next)
    p.int(_OP_goto, y0)

    /* process each field */
    for i, f := range fv {
        sw[i] = p.pc()
        fm.Set(f.Name, i)

        /* index to the field */
        for _, o := range f.Path {
            if p.int(_OP_index, int(o.Size)); o.Kind == resolver.F_deref {
                p.rtt(_OP_deref, o.Type)
            }
        }

        /* check for "stringnize" option */
        if (f.Opts & resolver.F_stringize) == 0 {
            self.compileOne(p, sp + 1, f.Type)
        } else {
            self.compileStructFieldStr(p, sp + 1, f.Type)
        }

        /* load the state, and try next field */
        p.add(_OP_load)
        p.int(_OP_goto, y0)
    }

    p.pin(x)
    p.pin(y1)
    p.add(_OP_drop)
    p.pin(n)
    p.pin(skip)
}

func (self *_Compiler) compileStructFieldStrUnmarshal(p *_Program, vt reflect.Type) {
    p.add(_OP_lspace)
    n0 := p.pc()
    p.add(_OP_is_null)
    self.checkMarshaler(p, vt, checkMarshalerFlags_quoted, true)
    p.pin(n0)
}

func (self *_Compiler) compileStructFieldStr(p *_Program, sp int, vt reflect.Type) {
    // according to std, json.Unmarshaler should be called before stringize
    // see https://github.com/bytedance/sonic/issues/670
    if self.checkMarshaler(p, vt, checkMarshalerFlags_quoted, false) {
        self.compileStructFieldStrUnmarshal(p, vt)
        return
    }

    n1 := -1
    ft := vt
    sv := false

    /* dereference the pointer if needed */
    if ft.Kind() == reflect.Ptr {
        ft = ft.Elem()
    }

    /* check if it can be stringized */
    switch ft.Kind() {
        case reflect.Bool    : sv = true
        case reflect.Int     : sv = true
        case reflect.Int8    : sv = true
        case reflect.Int16   : sv = true
        case reflect.Int32   : sv = true
        case reflect.Int64   : sv = true
        case reflect.Uint    : sv = true
        case reflect.Uint8   : sv = true
        case reflect.Uint16  : sv = true
        case reflect.Uint32  : sv = true
        case reflect.Uint64  : sv = true
        case reflect.Uintptr : sv = true
        case reflect.Float32 : sv = true
        case reflect.Float64 : sv = true
        case reflect.String  : sv = true
    }

    /* if it's not, ignore the "string" and follow the regular path */
    if !sv {
        self.compileOne(p, sp, vt)
        return
    }

    /* remove the leading space, and match the leading quote */
    vk := vt.Kind()
    p.add(_OP_lspace)
    n0 := p.pc()
    p.add(_OP_is_null)
    
    skip := self.checkIfSkip(p, stringType, '"')

    /* also check for inner "null" */
    n1 = p.pc()
    p.add(_OP_is_null_quote)

    /* dereference the pointer only when it is not null */
    if vk == reflect.Ptr {
        vt = vt.Elem()
        p.rtt(_OP_deref, vt)
    }

    n2 := p.pc()
    p.chr(_OP_check_char_0, '"')

    /* string opcode selector */
    _OP_string := func() _Op {
        if ft == jsonNumberType {
            return _OP_num
        } else {
            return _OP_unquote
        }
    }

    /* compile for each type */
    switch vt.Kind() {
        case reflect.Bool    : p.add(_OP_bool)
        case reflect.Int     : p.add(_OP_int())
        case reflect.Int8    : p.add(_OP_i8)
        case reflect.Int16   : p.add(_OP_i16)
        case reflect.Int32   : p.add(_OP_i32)
        case reflect.Int64   : p.add(_OP_i64)
        case reflect.Uint    : p.add(_OP_uint())
        case reflect.Uint8   : p.add(_OP_u8)
        case reflect.Uint16  : p.add(_OP_u16)
        case reflect.Uint32  : p.add(_OP_u32)
        case reflect.Uint64  : p.add(_OP_u64)
        case reflect.Uintptr : p.add(_OP_uintptr())
        case reflect.Float32 : p.add(_OP_f32)
        case reflect.Float64 : p.add(_OP_f64)
        case reflect.String  : p.add(_OP_string())
        default              : panic("not reachable")
    }

    /* the closing quote is not needed when parsing a pure string */
    if vt == jsonNumberType || vt.Kind() != reflect.String {
        p.chr(_OP_match_char, '"')
    }

    /* pin the `is_null_quote` jump location */
    if n1 != -1 && vk != reflect.Ptr {
        p.pin(n1)
    }

    /* "null" but not a pointer, act as if the field is not present */
    if vk != reflect.Ptr {
        pc2 := p.pc()
        p.add(_OP_goto)
        p.pin(n2)
        p.rtt(_OP_dismatch_err, vt)
        p.int(_OP_add, 1)
        p.pin(pc2)
        p.pin(n0)
        return
    }

    /* the "null" case of the pointer */
    pc := p.pc()
    p.add(_OP_goto)
    p.pin(n0) // `is_null` jump location
    p.pin(n1) // `is_null_quote` jump location
    p.add(_OP_nil_1)
    pc2 := p.pc()
    p.add(_OP_goto)
    p.pin(n2)
    p.rtt(_OP_dismatch_err, vt)
    p.int(_OP_add, 1)
    p.pin(pc)
    p.pin(pc2)
    p.pin(skip)
}

func (self *_Compiler) compileInterface(p *_Program, vt reflect.Type) {
    i := p.pc()
    p.add(_OP_is_null)

    /* check for empty interface */
    if vt.NumMethod() == 0 {
        p.add(_OP_any)
    } else {
        p.rtt(_OP_dyn, vt)
    }

    /* finish the OpCode */
    j := p.pc()
    p.add(_OP_goto)
    p.pin(i)
    p.add(_OP_nil_2)
    p.pin(j)
}

func (self *_Compiler) compilePrimitive(_ reflect.Type, p *_Program, op _Op) {
    i := p.pc()
    p.add(_OP_is_null)
    p.add(op)
    p.pin(i)
}

func (self *_Compiler) compileUnmarshalEnd(p *_Program, vt reflect.Type, i int) {
    j := p.pc()
    k := vt.Kind()

    /* not a pointer */
    if k != reflect.Ptr {
        p.pin(i)
        return
    }

    /* it seems that in Go JSON library, "null" takes priority over any kind of unmarshaler */
    p.add(_OP_goto)
    p.pin(i)
    p.add(_OP_nil_1)
    p.pin(j)
}

func (self *_Compiler) compileUnmarshalJson(p *_Program, vt reflect.Type, flags int) {
    i := p.pc()
    v := _OP_unmarshal
    p.add(_OP_is_null)

    /* check for dynamic interface */
    if vt.Kind() == reflect.Interface {
        v = _OP_dyn
    }

    /* call the unmarshaler */
    p.rtti(v, vt, flags)
    self.compileUnmarshalEnd(p, vt, i)
}

func (self *_Compiler) compileUnmarshalText(p *_Program, vt reflect.Type, iv int) {
    i := p.pc()
    v := _OP_unmarshal_text
    p.add(_OP_is_null)

    /* check for dynamic interface */
    if vt.Kind() == reflect.Interface {
        v = _OP_dyn
    } else {
        p.chr(_OP_match_char, '"')
    }

    /* call the unmarshaler */
    p.rtti(v, vt, iv)
    self.compileUnmarshalEnd(p, vt, i)
}

func (self *_Compiler) compileUnmarshalTextPtr(p *_Program, vt reflect.Type, iv int) {
    i := p.pc()
    p.add(_OP_is_null)
    p.chr(_OP_match_char, '"')
    p.rtti(_OP_unmarshal_text_p, vt, iv)
    p.pin(i)
}

func (self *_Compiler) checkIfSkip(p *_Program, vt reflect.Type, c byte) int {
    j := p.pc()
    p.chr(_OP_check_char_0, c)
    p.rtt(_OP_dismatch_err, vt)
    s := p.pc()
    p.add(_OP_go_skip)
    p.pin(j)
    p.int(_OP_add, 1)
    return s
}
