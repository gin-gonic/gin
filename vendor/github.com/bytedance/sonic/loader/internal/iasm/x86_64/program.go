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
	"fmt"
	"math"
	"math/bits"

	"github.com/bytedance/sonic/loader/internal/iasm/expr"
)

type (
	_PseudoType         int
	_InstructionEncoder func(*Program, ...interface{}) *Instruction
)

const (
	_PseudoNop _PseudoType = iota + 1
	_PseudoByte
	_PseudoWord
	_PseudoLong
	_PseudoQuad
	_PseudoData
	_PseudoAlign
)

func (self _PseudoType) String() string {
	switch self {
	case _PseudoNop:
		return ".nop"
	case _PseudoByte:
		return ".byte"
	case _PseudoWord:
		return ".word"
	case _PseudoLong:
		return ".long"
	case _PseudoQuad:
		return ".quad"
	case _PseudoData:
		return ".data"
	case _PseudoAlign:
		return ".align"
	default:
		panic("unreachable")
	}
}

type _Pseudo struct {
	kind _PseudoType
	data []byte
	uint uint64
	expr *expr.Expr
}

func (self *_Pseudo) free() {
	if self.expr != nil {
		self.expr.Free()
	}
}

func (self *_Pseudo) encode(m *[]byte, pc uintptr) int {
	switch self.kind {
	case _PseudoNop:
		return 0
	case _PseudoByte:
		self.encodeByte(m)
		return 1
	case _PseudoWord:
		self.encodeWord(m)
		return 2
	case _PseudoLong:
		self.encodeLong(m)
		return 4
	case _PseudoQuad:
		self.encodeQuad(m)
		return 8
	case _PseudoData:
		self.encodeData(m)
		return len(self.data)
	case _PseudoAlign:
		self.encodeAlign(m, pc)
		return self.alignSize(pc)
	default:
		panic("invalid pseudo instruction")
	}
}

func (self *_Pseudo) evalExpr(low int64, high int64) int64 {
	if v, err := self.expr.Evaluate(); err != nil {
		panic(err)
	} else if v < low || v > high {
		panic(fmt.Sprintf("expression out of range [%d, %d]: %d", low, high, v))
	} else {
		return v
	}
}

func (self *_Pseudo) alignSize(pc uintptr) int {
	if !ispow2(self.uint) {
		panic(fmt.Sprintf("alignment should be a power of 2, not %d", self.uint))
	} else {
		return align(int(pc), bits.TrailingZeros64(self.uint)) - int(pc)
	}
}

func (self *_Pseudo) encodeData(m *[]byte) {
	if m != nil {
		*m = append(*m, self.data...)
	}
}

func (self *_Pseudo) encodeByte(m *[]byte) {
	if m != nil {
		append8(m, byte(self.evalExpr(math.MinInt8, math.MaxUint8)))
	}
}

func (self *_Pseudo) encodeWord(m *[]byte) {
	if m != nil {
		append16(m, uint16(self.evalExpr(math.MinInt16, math.MaxUint16)))
	}
}

func (self *_Pseudo) encodeLong(m *[]byte) {
	if m != nil {
		append32(m, uint32(self.evalExpr(math.MinInt32, math.MaxUint32)))
	}
}

func (self *_Pseudo) encodeQuad(m *[]byte) {
	if m != nil {
		if v, err := self.expr.Evaluate(); err != nil {
			panic(err)
		} else {
			append64(m, uint64(v))
		}
	}
}

func (self *_Pseudo) encodeAlign(m *[]byte, pc uintptr) {
	if m != nil {
		if self.expr == nil {
			expandmm(m, self.alignSize(pc), 0)
		} else {
			expandmm(m, self.alignSize(pc), byte(self.evalExpr(math.MinInt8, math.MaxUint8)))
		}
	}
}

// Operands represents a sequence of operand required by an instruction.
type Operands [_N_args]interface{}

// InstructionDomain represents the domain of an instruction.
type InstructionDomain uint8

const (
	DomainGeneric InstructionDomain = iota
	DomainMMXSSE
	DomainAVX
	DomainFMA
	DomainCrypto
	DomainMask
	DomainAMDSpecific
	DomainMisc
	DomainPseudo
)

type (
	_BranchType uint8
)

const (
	_B_none _BranchType = iota
	_B_conditional
	_B_unconditional
)

// Instruction represents an unencoded instruction.
type Instruction struct {
	next   *Instruction
	pc     uintptr
	nb     int
	len    int
	argc   int
	name   string
	argv   Operands
	forms  [_N_forms]_Encoding
	pseudo _Pseudo
	branch _BranchType
	domain InstructionDomain
	prefix []byte
}

func (self *Instruction) add(flags int, encoder func(m *_Encoding, v []interface{})) {
	self.forms[self.len].flags = flags
	self.forms[self.len].encoder = encoder
	self.len++
}

func (self *Instruction) free() {
	self.clear()
	self.pseudo.free()
	//freeInstruction(self)
}

func (self *Instruction) clear() {
	for i := 0; i < self.argc; i++ {
		if v, ok := self.argv[i].(Disposable); ok {
			v.Free()
		}
	}
}

func (self *Instruction) check(e *_Encoding) bool {
	if (e.flags & _F_rel1) != 0 {
		return isRel8(self.argv[0])
	} else if (e.flags & _F_rel4) != 0 {
		return isRel32(self.argv[0]) || isLabel(self.argv[0])
	} else {
		return true
	}
}

func (self *Instruction) encode(m *[]byte) int {
	n := math.MaxInt64
	p := (*_Encoding)(nil)

	/* encode prefixes if any */
	if self.nb = len(self.prefix); m != nil {
		*m = append(*m, self.prefix...)
	}

	/* check for pseudo-instructions */
	if self.pseudo.kind != 0 {
		self.nb += self.pseudo.encode(m, self.pc)
		return self.nb
	}

	/* find the shortest encoding */
	for i := 0; i < self.len; i++ {
		if e := &self.forms[i]; self.check(e) {
			if v := e.encode(self.argv[:self.argc]); v < n {
				n = v
				p = e
			}
		}
	}

	/* add to buffer if needed */
	if m != nil {
		*m = append(*m, p.bytes[:n]...)
	}

	/* update the instruction length */
	self.nb += n
	return self.nb
}

/** Instruction Prefixes **/

const (
	_P_cs   = 0x2e
	_P_ds   = 0x3e
	_P_es   = 0x26
	_P_fs   = 0x64
	_P_gs   = 0x65
	_P_ss   = 0x36
	_P_lock = 0xf0
)

// CS overrides the memory operation of this instruction to CS.
func (self *Instruction) CS() *Instruction {
	self.prefix = append(self.prefix, _P_cs)
	return self
}

// DS overrides the memory operation of this instruction to DS,
// this is the default section for most instructions if not specified.
func (self *Instruction) DS() *Instruction {
	self.prefix = append(self.prefix, _P_ds)
	return self
}

// ES overrides the memory operation of this instruction to ES.
func (self *Instruction) ES() *Instruction {
	self.prefix = append(self.prefix, _P_es)
	return self
}

// FS overrides the memory operation of this instruction to FS.
func (self *Instruction) FS() *Instruction {
	self.prefix = append(self.prefix, _P_fs)
	return self
}

// GS overrides the memory operation of this instruction to GS.
func (self *Instruction) GS() *Instruction {
	self.prefix = append(self.prefix, _P_gs)
	return self
}

// SS overrides the memory operation of this instruction to SS.
func (self *Instruction) SS() *Instruction {
	self.prefix = append(self.prefix, _P_ss)
	return self
}

// LOCK causes the processor's LOCK# signal to be asserted during execution of
// the accompanying instruction (turns the instruction into an atomic instruction).
// In a multiprocessor environment, the LOCK# signal insures that the processor
// has exclusive use of any shared memory while the signal is asserted.
func (self *Instruction) LOCK() *Instruction {
	self.prefix = append(self.prefix, _P_lock)
	return self
}

/** Basic Instruction Properties **/

// Name returns the instruction name.
func (self *Instruction) Name() string {
	return self.name
}

// Domain returns the domain of this instruction.
func (self *Instruction) Domain() InstructionDomain {
	return self.domain
}

// Operands returns the operands of this instruction.
func (self *Instruction) Operands() []interface{} {
	return self.argv[:self.argc]
}

// Program represents a sequence of instructions.
type Program struct {
	arch *Arch
	head *Instruction
	tail *Instruction
}

const (
	_N_near       = 2 // near-branch (-128 ~ +127) takes 2 bytes to encode
	_N_far_cond   = 6 // conditional far-branch takes 6 bytes to encode
	_N_far_uncond = 5 // unconditional far-branch takes 5 bytes to encode
)

func (self *Program) clear() {
	for p, q := self.head, self.head; p != nil; p = q {
		q = p.next
		p.free()
	}
}

func (self *Program) alloc(name string, argc int, argv Operands) *Instruction {
	p := self.tail
	q := newInstruction(name, argc, argv)

	/* attach to tail if any */
	if p != nil {
		p.next = q
	} else {
		self.head = q
	}

	/* set the new tail */
	self.tail = q
	return q
}

func (self *Program) pseudo(kind _PseudoType) (p *Instruction) {
	p = self.alloc(kind.String(), 0, Operands{})
	p.domain = DomainPseudo
	p.pseudo.kind = kind
	return
}

func (self *Program) require(isa ISA) {
	if !self.arch.HasISA(isa) {
		panic("ISA '" + isa.String() + "' was not enabled")
	}
}

func (self *Program) branchSize(p *Instruction) int {
	switch p.branch {
	case _B_none:
		panic("p is not a branch")
	case _B_conditional:
		return _N_far_cond
	case _B_unconditional:
		return _N_far_uncond
	default:
		panic("invalid instruction")
	}
}

/** Pseudo-Instructions **/

// Byte is a pseudo-instruction to add raw byte to the assembled code.
func (self *Program) Byte(v *expr.Expr) (p *Instruction) {
	p = self.pseudo(_PseudoByte)
	p.pseudo.expr = v
	return
}

// Word is a pseudo-instruction to add raw uint16 as little-endian to the assembled code.
func (self *Program) Word(v *expr.Expr) (p *Instruction) {
	p = self.pseudo(_PseudoWord)
	p.pseudo.expr = v
	return
}

// Long is a pseudo-instruction to add raw uint32 as little-endian to the assembled code.
func (self *Program) Long(v *expr.Expr) (p *Instruction) {
	p = self.pseudo(_PseudoLong)
	p.pseudo.expr = v
	return
}

// Quad is a pseudo-instruction to add raw uint64 as little-endian to the assembled code.
func (self *Program) Quad(v *expr.Expr) (p *Instruction) {
	p = self.pseudo(_PseudoQuad)
	p.pseudo.expr = v
	return
}

// Data is a pseudo-instruction to add raw bytes to the assembled code.
func (self *Program) Data(v []byte) (p *Instruction) {
	p = self.pseudo(_PseudoData)
	p.pseudo.data = v
	return
}

// Align is a pseudo-instruction to ensure the PC is aligned to a certain value.
func (self *Program) Align(align uint64, padding *expr.Expr) (p *Instruction) {
	p = self.pseudo(_PseudoAlign)
	p.pseudo.uint = align
	p.pseudo.expr = padding
	return
}

/** Program Assembler **/

// Free returns the Program object into pool.
// Any operation performed after Free is undefined behavior.
//
// NOTE: This also frees all the instructions, labels, memory
//
//	operands and expressions associated with this program.
func (self *Program) Free() {
	self.clear()
	//freeProgram(self)
}

// Link pins a label at the current position.
func (self *Program) Link(p *Label) {
	if p.Dest != nil {
		panic("label was already linked")
	} else {
		p.Dest = self.pseudo(_PseudoNop)
	}
}

// Assemble assembles and links the entire program into machine code.
func (self *Program) Assemble(pc uintptr) (ret []byte) {
	orig := pc
	next := true
	offs := uintptr(0)

	/* Pass 0: PC-precompute, assume all labeled branches are far-branches. */
	for p := self.head; p != nil; p = p.next {
		if p.pc = pc; !isLabel(p.argv[0]) || p.branch == _B_none {
			pc += uintptr(p.encode(nil))
		} else {
			pc += uintptr(self.branchSize(p))
		}
	}

	/* allocate space for the machine code */
	nb := int(pc - orig)
	ret = make([]byte, 0, nb)

	/* Pass 1: adjust all the jumps */
	for next {
		next = false
		offs = uintptr(0)

		/* scan all the branches */
		for p := self.head; p != nil; p = p.next {
			var ok bool
			var lb *Label

			/* re-calculate the alignment here */
			if nb = p.nb; p.pseudo.kind == _PseudoAlign {
				p.pc -= offs
				offs += uintptr(nb - p.encode(nil))
				continue
			}

			/* adjust the program counter */
			p.pc -= offs
			lb, ok = p.argv[0].(*Label)

			/* only care about labeled far-branches */
			if !ok || p.nb == _N_near || p.branch == _B_none {
				continue
			}

			/* calculate the jump offset */
			size := self.branchSize(p)
			diff := lb.offset(p.pc, size)

			/* too far to be a near jump */
			if diff > 127 || diff < -128 {
				p.nb = size
				continue
			}

			/* a far jump becomes a near jump, calculate
			 * the PC adjustment value and assemble again */
			next = true
			p.nb = _N_near
			offs += uintptr(size - _N_near)
		}
	}

	/* Pass 3: link all the cross-references */
	for p := self.head; p != nil; p = p.next {
		for i := 0; i < p.argc; i++ {
			var ok bool
			var lb *Label
			var op *MemoryOperand

			/* resolve labels */
			if lb, ok = p.argv[i].(*Label); ok {
				p.argv[i] = lb.offset(p.pc, p.nb)
				continue
			}

			/* check for memory operands */
			if op, ok = p.argv[i].(*MemoryOperand); !ok {
				continue
			}

			/* check for label references */
			if op.Addr.Type != Reference {
				continue
			}

			/* replace the label with the real offset */
			op.Addr.Type = Offset
			op.Addr.Offset = op.Addr.Reference.offset(p.pc, p.nb)
		}
	}

	/* Pass 4: actually encode all the instructions */
	for p := self.head; p != nil; p = p.next {
		p.encode(&ret)
	}

	/* all done */
	return ret
}

// AssembleAndFree is like Assemble, but it frees the Program after assembling.
func (self *Program) AssembleAndFree(pc uintptr) (ret []byte) {
	ret = self.Assemble(pc)
	self.Free()
	return
}
