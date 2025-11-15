//
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
//

package x86_64

import (
	"encoding/binary"
	"math"
)

/** Operand Encoding Helpers **/

func imml(v interface{}) byte {
	return byte(toImmAny(v) & 0x0f)
}

func relv(v interface{}) int64 {
	switch r := v.(type) {
	case *Label:
		return 0
	case RelativeOffset:
		return int64(r)
	default:
		panic("invalid relative offset")
	}
}

func addr(v interface{}) interface{} {
	switch a := v.(*MemoryOperand).Addr; a.Type {
	case Memory:
		return a.Memory
	case Offset:
		return a.Offset
	case Reference:
		return a.Reference
	default:
		panic("invalid memory operand type")
	}
}

func bcode(v interface{}) byte {
	if m, ok := v.(*MemoryOperand); !ok {
		panic("v is not a memory operand")
	} else if m.Broadcast == 0 {
		return 0
	} else {
		return 1
	}
}

func vcode(v interface{}) byte {
	switch r := v.(type) {
	case XMMRegister:
		return byte(r)
	case YMMRegister:
		return byte(r)
	case ZMMRegister:
		return byte(r)
	case MaskedRegister:
		return vcode(r.Reg)
	default:
		panic("v is not a vector register")
	}
}

func kcode(v interface{}) byte {
	switch r := v.(type) {
	case KRegister:
		return byte(r)
	case XMMRegister:
		return 0
	case YMMRegister:
		return 0
	case ZMMRegister:
		return 0
	case RegisterMask:
		return byte(r.K)
	case MaskedRegister:
		return byte(r.Mask.K)
	case *MemoryOperand:
		return toKcodeMem(r)
	default:
		panic("v is not a maskable operand")
	}
}

func zcode(v interface{}) byte {
	switch r := v.(type) {
	case KRegister:
		return 0
	case XMMRegister:
		return 0
	case YMMRegister:
		return 0
	case ZMMRegister:
		return 0
	case RegisterMask:
		return toZcodeRegM(r)
	case MaskedRegister:
		return toZcodeRegM(r.Mask)
	case *MemoryOperand:
		return toZcodeMem(r)
	default:
		panic("v is not a maskable operand")
	}
}

func lcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r & 0x07)
	case Register16:
		return byte(r & 0x07)
	case Register32:
		return byte(r & 0x07)
	case Register64:
		return byte(r & 0x07)
	case KRegister:
		return byte(r & 0x07)
	case MMRegister:
		return byte(r & 0x07)
	case XMMRegister:
		return byte(r & 0x07)
	case YMMRegister:
		return byte(r & 0x07)
	case ZMMRegister:
		return byte(r & 0x07)
	case MaskedRegister:
		return lcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func hcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r>>3) & 1
	case Register16:
		return byte(r>>3) & 1
	case Register32:
		return byte(r>>3) & 1
	case Register64:
		return byte(r>>3) & 1
	case KRegister:
		return byte(r>>3) & 1
	case MMRegister:
		return byte(r>>3) & 1
	case XMMRegister:
		return byte(r>>3) & 1
	case YMMRegister:
		return byte(r>>3) & 1
	case ZMMRegister:
		return byte(r>>3) & 1
	case MaskedRegister:
		return hcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func ecode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r>>4) & 1
	case Register16:
		return byte(r>>4) & 1
	case Register32:
		return byte(r>>4) & 1
	case Register64:
		return byte(r>>4) & 1
	case KRegister:
		return byte(r>>4) & 1
	case MMRegister:
		return byte(r>>4) & 1
	case XMMRegister:
		return byte(r>>4) & 1
	case YMMRegister:
		return byte(r>>4) & 1
	case ZMMRegister:
		return byte(r>>4) & 1
	case MaskedRegister:
		return ecode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func hlcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return toHLcodeReg8(r)
	case Register16:
		return byte(r & 0x0f)
	case Register32:
		return byte(r & 0x0f)
	case Register64:
		return byte(r & 0x0f)
	case KRegister:
		return byte(r & 0x0f)
	case MMRegister:
		return byte(r & 0x0f)
	case XMMRegister:
		return byte(r & 0x0f)
	case YMMRegister:
		return byte(r & 0x0f)
	case ZMMRegister:
		return byte(r & 0x0f)
	case MaskedRegister:
		return hlcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func ehcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r>>3) & 0x03
	case Register16:
		return byte(r>>3) & 0x03
	case Register32:
		return byte(r>>3) & 0x03
	case Register64:
		return byte(r>>3) & 0x03
	case KRegister:
		return byte(r>>3) & 0x03
	case MMRegister:
		return byte(r>>3) & 0x03
	case XMMRegister:
		return byte(r>>3) & 0x03
	case YMMRegister:
		return byte(r>>3) & 0x03
	case ZMMRegister:
		return byte(r>>3) & 0x03
	case MaskedRegister:
		return ehcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func toImmAny(v interface{}) int64 {
	if x, ok := asInt64(v); ok {
		return x
	} else {
		panic("value is not an integer")
	}
}

func toHcodeOpt(v interface{}) byte {
	if v == nil {
		return 0
	} else {
		return hcode(v)
	}
}

func toEcodeVMM(v interface{}, x byte) byte {
	switch r := v.(type) {
	case XMMRegister:
		return ecode(r)
	case YMMRegister:
		return ecode(r)
	case ZMMRegister:
		return ecode(r)
	default:
		return x
	}
}

func toKcodeMem(v *MemoryOperand) byte {
	if !v.Masked {
		return 0
	} else {
		return byte(v.Mask.K)
	}
}

func toZcodeMem(v *MemoryOperand) byte {
	if !v.Masked || v.Mask.Z {
		return 0
	} else {
		return 1
	}
}

func toZcodeRegM(v RegisterMask) byte {
	if v.Z {
		return 1
	} else {
		return 0
	}
}

func toHLcodeReg8(v Register8) byte {
	switch v {
	case AH:
		fallthrough
	case BH:
		fallthrough
	case CH:
		fallthrough
	case DH:
		panic("ah/bh/ch/dh registers never use 4-bit encoding")
	default:
		return byte(v & 0x0f)
	}
}

/** Instruction Encoding Helpers **/

const (
	_N_inst = 16
)

const (
	_F_rel1 = 1 << iota
	_F_rel4
)

type _Encoding struct {
	len     int
	flags   int
	bytes   [_N_inst]byte
	encoder func(m *_Encoding, v []interface{})
}

// buf ensures len + n <= len(bytes).
func (self *_Encoding) buf(n int) []byte {
	if i := self.len; i+n > _N_inst {
		panic("instruction too long")
	} else {
		return self.bytes[i:]
	}
}

// emit encodes a single byte.
func (self *_Encoding) emit(v byte) {
	self.buf(1)[0] = v
	self.len++
}

// imm1 encodes a single byte immediate value.
func (self *_Encoding) imm1(v int64) {
	self.emit(byte(v))
}

// imm2 encodes a two-byte immediate value in little-endian.
func (self *_Encoding) imm2(v int64) {
	binary.LittleEndian.PutUint16(self.buf(2), uint16(v))
	self.len += 2
}

// imm4 encodes a 4-byte immediate value in little-endian.
func (self *_Encoding) imm4(v int64) {
	binary.LittleEndian.PutUint32(self.buf(4), uint32(v))
	self.len += 4
}

// imm8 encodes an 8-byte immediate value in little-endian.
func (self *_Encoding) imm8(v int64) {
	binary.LittleEndian.PutUint64(self.buf(8), uint64(v))
	self.len += 8
}

// vex2 encodes a 2-byte or 3-byte VEX prefix.
//
//	2-byte VEX prefix:
//
// Requires: VEX.W = 0, VEX.mmmmm = 0b00001 and VEX.B = VEX.X = 0
//
//	+----------------+
//
// Byte 0: | Bits 0-7: 0xc5 |
//
//	+----------------+
//
//	+-----------+----------------+----------+--------------+
//
// Byte 1: | Bit 7: ~R | Bits 3-6 ~vvvv | Bit 2: L | Bits 0-1: pp |
//
//	+-----------+----------------+----------+--------------+
//
//	                 3-byte VEX prefix:
//	+----------------+
//
// Byte 0: | Bits 0-7: 0xc4 |
//
//	+----------------+
//
//	+-----------+-----------+-----------+-------------------+
//
// Byte 1: | Bit 7: ~R | Bit 6: ~X | Bit 5: ~B | Bits 0-4: 0b00001 |
//
//	+-----------+-----------+-----------+-------------------+
//
//	+----------+-----------------+----------+--------------+
//
// Byte 2: | Bit 7: 0 | Bits 3-6: ~vvvv | Bit 2: L | Bits 0-1: pp |
//
//	+----------+-----------------+----------+--------------+
func (self *_Encoding) vex2(lpp byte, r byte, rm interface{}, vvvv byte) {
	var b byte
	var x byte

	/* VEX.R must be a single-bit mask */
	if r > 1 {
		panic("VEX.R must be a 1-bit mask")
	}

	/* VEX.Lpp must be a 3-bit mask */
	if lpp&^0b111 != 0 {
		panic("VEX.Lpp must be a 3-bit mask")
	}

	/* VEX.vvvv must be a 4-bit mask */
	if vvvv&^0b1111 != 0 {
		panic("VEX.vvvv must be a 4-bit mask")
	}

	/* encode the RM bits if any */
	if rm != nil {
		switch v := rm.(type) {
		case *Label:
			break
		case Register:
			b = hcode(v)
		case MemoryAddress:
			b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
		case RelativeOffset:
			break
		default:
			panic("rm is expected to be a register or a memory address")
		}
	}

	/* if VEX.B and VEX.X are zeroes, 2-byte VEX prefix can be used */
	if x == 0 && b == 0 {
		self.emit(0xc5)
		self.emit(0xf8 ^ (r << 7) ^ (vvvv << 3) ^ lpp)
	} else {
		self.emit(0xc4)
		self.emit(0xe1 ^ (r << 7) ^ (x << 6) ^ (b << 5))
		self.emit(0x78 ^ (vvvv << 3) ^ lpp)
	}
}

// vex3 encodes a 3-byte VEX or XOP prefix.
//
//	                3-byte VEX/XOP prefix
//	+-----------------------------------+
//
// Byte 0: | Bits 0-7: 0xc4 (VEX) / 0x8f (XOP) |
//
//	+-----------------------------------+
//
//	+-----------+-----------+-----------+-----------------+
//
// Byte 1: | Bit 7: ~R | Bit 6: ~X | Bit 5: ~B | Bits 0-4: mmmmm |
//
//	+-----------+-----------+-----------+-----------------+
//
//	+----------+-----------------+----------+--------------+
//
// Byte 2: | Bit 7: W | Bits 3-6: ~vvvv | Bit 2: L | Bits 0-1: pp |
//
//	+----------+-----------------+----------+--------------+
func (self *_Encoding) vex3(esc byte, mmmmm byte, wlpp byte, r byte, rm interface{}, vvvv byte) {
	var b byte
	var x byte

	/* VEX.R must be a single-bit mask */
	if r > 1 {
		panic("VEX.R must be a 1-bit mask")
	}

	/* VEX.vvvv must be a 4-bit mask */
	if vvvv&^0b1111 != 0 {
		panic("VEX.vvvv must be a 4-bit mask")
	}

	/* escape must be a 3-byte VEX (0xc4) or XOP (0x8f) prefix */
	if esc != 0xc4 && esc != 0x8f {
		panic("escape must be a 3-byte VEX (0xc4) or XOP (0x8f) prefix")
	}

	/* VEX.W____Lpp is expected to have no bits set except 0, 1, 2 and 7 */
	if wlpp&^0b10000111 != 0 {
		panic("VEX.W____Lpp is expected to have no bits set except 0, 1, 2 and 7")
	}

	/* VEX.m-mmmm is expected to be a 5-bit mask */
	if mmmmm&^0b11111 != 0 {
		panic("VEX.m-mmmm is expected to be a 5-bit mask")
	}

	/* encode the RM bits */
	switch v := rm.(type) {
	case *Label:
		break
	case MemoryAddress:
		b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
	case RelativeOffset:
		break
	default:
		panic("rm is expected to be a register or a memory address")
	}

	/* encode the 3-byte VEX or XOP prefix */
	self.emit(esc)
	self.emit(0xe0 ^ (r << 7) ^ (x << 6) ^ (b << 5) ^ mmmmm)
	self.emit(0x78 ^ (vvvv << 3) ^ wlpp)
}

// evex encodes a 4-byte EVEX prefix.
func (self *_Encoding) evex(mm byte, w1pp byte, ll byte, rr byte, rm interface{}, vvvvv byte, aaa byte, zz byte, bb byte) {
	var b byte
	var x byte

	/* EVEX.b must be a single-bit mask */
	if bb > 1 {
		panic("EVEX.b must be a 1-bit mask")
	}

	/* EVEX.z must be a single-bit mask */
	if zz > 1 {
		panic("EVEX.z must be a 1-bit mask")
	}

	/* EVEX.mm must be a 2-bit mask */
	if mm&^0b11 != 0 {
		panic("EVEX.mm must be a 2-bit mask")
	}

	/* EVEX.L'L must be a 2-bit mask */
	if ll&^0b11 != 0 {
		panic("EVEX.L'L must be a 2-bit mask")
	}

	/* EVEX.R'R must be a 2-bit mask */
	if rr&^0b11 != 0 {
		panic("EVEX.R'R must be a 2-bit mask")
	}

	/* EVEX.aaa must be a 3-bit mask */
	if aaa&^0b111 != 0 {
		panic("EVEX.aaa must be a 3-bit mask")
	}

	/* EVEX.v'vvvv must be a 5-bit mask */
	if vvvvv&^0b11111 != 0 {
		panic("EVEX.v'vvvv must be a 5-bit mask")
	}

	/* EVEX.W____1pp is expected to have no bits set except 0, 1, 2, and 7 */
	if w1pp&^0b10000011 != 0b100 {
		panic("EVEX.W____1pp is expected to have no bits set except 0, 1, 2, and 7")
	}

	/* extract bits from EVEX.R'R and EVEX.v'vvvv */
	r1, r0 := rr>>1, rr&1
	v1, v0 := vvvvv>>4, vvvvv&0b1111

	/* encode the RM bits if any */
	if rm != nil {
		switch m := rm.(type) {
		case *Label:
			break
		case Register:
			b, x = hcode(m), ecode(m)
		case MemoryAddress:
			b, x, v1 = toHcodeOpt(m.Base), toHcodeOpt(m.Index), toEcodeVMM(m.Index, v1)
		case RelativeOffset:
			break
		default:
			panic("rm is expected to be a register or a memory address")
		}
	}

	/* EVEX prefix bytes */
	p0 := (r0 << 7) | (x << 6) | (b << 5) | (r1 << 4) | mm
	p1 := (v0 << 3) | w1pp
	p2 := (zz << 7) | (ll << 5) | (b << 4) | (v1 << 3) | aaa

	/* p0: invert RXBR' (bits 4-7)
	 * p1: invert vvvv  (bits 3-6)
	 * p2: invert V'    (bit  3) */
	self.emit(0x62)
	self.emit(p0 ^ 0xf0)
	self.emit(p1 ^ 0x78)
	self.emit(p2 ^ 0x08)
}

// rexm encodes a mandatory REX prefix.
func (self *_Encoding) rexm(w byte, r byte, rm interface{}) {
	var b byte
	var x byte

	/* REX.R must be 0 or 1 */
	if r != 0 && r != 1 {
		panic("REX.R must be 0 or 1")
	}

	/* REX.W must be 0 or 1 */
	if w != 0 && w != 1 {
		panic("REX.W must be 0 or 1")
	}

	/* encode the RM bits */
	switch v := rm.(type) {
	case *Label:
		break
	case MemoryAddress:
		b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
	case RelativeOffset:
		break
	default:
		panic("rm is expected to be a register or a memory address")
	}

	/* encode the REX prefix */
	self.emit(0x40 | (w << 3) | (r << 2) | (x << 1) | b)
}

// rexo encodes an optional REX prefix.
func (self *_Encoding) rexo(r byte, rm interface{}, force bool) {
	var b byte
	var x byte

	/* REX.R must be 0 or 1 */
	if r != 0 && r != 1 {
		panic("REX.R must be 0 or 1")
	}

	/* encode the RM bits */
	switch v := rm.(type) {
	case *Label:
		break
	case Register:
		b = hcode(v)
	case MemoryAddress:
		b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
	case RelativeOffset:
		break
	default:
		panic("rm is expected to be a register or a memory address")
	}

	/* if REX.R, REX.X, and REX.B are all zeroes, REX prefix can be omitted */
	if force || r != 0 || x != 0 || b != 0 {
		self.emit(0x40 | (r << 2) | (x << 1) | b)
	}
}

// mrsd encodes ModR/M, SIB and Displacement.
//
//	ModR/M byte
//
// +----------------+---------------+---------------+
// | Bits 6-7: Mode | Bits 3-5: Reg | Bits 0-2: R/M |
// +----------------+---------------+---------------+
//
//	SIB byte
//
// +-----------------+-----------------+----------------+
// | Bits 6-7: Scale | Bits 3-5: Index | Bits 0-2: Base |
// +-----------------+-----------------+----------------+
func (self *_Encoding) mrsd(reg byte, rm interface{}, disp8v int32) {
	var ok bool
	var mm MemoryAddress
	var ro RelativeOffset

	/* ModRM encodes the lower 3-bit of the register */
	if reg > 7 {
		panic("invalid register bits")
	}

	/* check the displacement scale */
	switch disp8v {
	case 1:
		break
	case 2:
		break
	case 4:
		break
	case 8:
		break
	case 16:
		break
	case 32:
		break
	case 64:
		break
	default:
		panic("invalid displacement size")
	}

	/* special case: unresolved labels, assuming a zero offset */
	if _, ok = rm.(*Label); ok {
		self.emit(0x05 | (reg << 3))
		self.imm4(0)
		return
	}

	/* special case: RIP-relative offset
	 * ModRM.Mode == 0 and ModeRM.R/M == 5 indicates (rip + disp32) addressing */
	if ro, ok = rm.(RelativeOffset); ok {
		self.emit(0x05 | (reg << 3))
		self.imm4(int64(ro))
		return
	}

	/* must be a generic memory address */
	if mm, ok = rm.(MemoryAddress); !ok {
		panic("rm must be a memory address")
	}

	/* absolute addressing, encoded as disp(%rbp,%rsp,1) */
	if mm.Base == nil && mm.Index == nil {
		self.emit(0x04 | (reg << 3))
		self.emit(0x25)
		self.imm4(int64(mm.Displacement))
		return
	}

	/* no SIB byte */
	if mm.Index == nil && lcode(mm.Base) != 0b100 {
		cc := lcode(mm.Base)
		dv := mm.Displacement

		/* ModRM.Mode == 0 (no displacement) */
		if dv == 0 && mm.Base != RBP && mm.Base != R13 {
			if cc == 0b101 {
				panic("rbp/r13 is not encodable as a base register (interpreted as disp32 address)")
			} else {
				self.emit((reg << 3) | cc)
				return
			}
		}

		/* ModRM.Mode == 1 (8-bit displacement) */
		if dq := dv / disp8v; dq >= math.MinInt8 && dq <= math.MaxInt8 && dv%disp8v == 0 {
			self.emit(0x40 | (reg << 3) | cc)
			self.imm1(int64(dq))
			return
		}

		/* ModRM.Mode == 2 (32-bit displacement) */
		self.emit(0x80 | (reg << 3) | cc)
		self.imm4(int64(mm.Displacement))
		return
	}

	/* all encodings below use ModRM.R/M = 4 (0b100) to indicate the presence of SIB */
	if mm.Index == RSP {
		panic("rsp is not encodable as an index register (interpreted as no index)")
	}

	/* index = 4 (0b100) denotes no-index encoding */
	var scale byte
	var index byte = 0x04

	/* encode the scale byte */
	if mm.Scale != 0 {
		switch mm.Scale {
		case 1:
			scale = 0
		case 2:
			scale = 1
		case 4:
			scale = 2
		case 8:
			scale = 3
		default:
			panic("invalid scale value")
		}
	}

	/* encode the index byte */
	if mm.Index != nil {
		index = lcode(mm.Index)
	}

	/* SIB.Base = 5 (0b101) and ModRM.Mode = 0 indicates no-base encoding with disp32 */
	if mm.Base == nil {
		self.emit((reg << 3) | 0b100)
		self.emit((scale << 6) | (index << 3) | 0b101)
		self.imm4(int64(mm.Displacement))
		return
	}

	/* base L-code & displacement value */
	cc := lcode(mm.Base)
	dv := mm.Displacement

	/* ModRM.Mode == 0 (no displacement) */
	if dv == 0 && cc != 0b101 {
		self.emit((reg << 3) | 0b100)
		self.emit((scale << 6) | (index << 3) | cc)
		return
	}

	/* ModRM.Mode == 1 (8-bit displacement) */
	if dq := dv / disp8v; dq >= math.MinInt8 && dq <= math.MaxInt8 && dv%disp8v == 0 {
		self.emit(0x44 | (reg << 3))
		self.emit((scale << 6) | (index << 3) | cc)
		self.imm1(int64(dq))
		return
	}

	/* ModRM.Mode == 2 (32-bit displacement) */
	self.emit(0x84 | (reg << 3))
	self.emit((scale << 6) | (index << 3) | cc)
	self.imm4(int64(mm.Displacement))
}

// encode invokes the encoder to encode this instruction.
func (self *_Encoding) encode(v []interface{}) int {
	self.len = 0
	self.encoder(self, v)
	return self.len
}
