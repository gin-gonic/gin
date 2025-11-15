//go:build go1.17
// +build go1.17

/*
 * Copyright 2022 ByteDance Inc.
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

/** Go Internal ABI implementation
 *
 *  This module implements the function layout algorithm described by the Go internal ABI.
 *  See https://github.com/golang/go/blob/master/src/cmd/compile/abi-internal.md for more info.
 */

package abi

import (
	"fmt"
	"reflect"

	x64 "github.com/bytedance/sonic/loader/internal/iasm/x86_64"
)

/** Frame Structure of the Generated Function
    FP  +------------------------------+
        |             . . .            |
        | 2nd reg argument spill space |
        + 1st reg argument spill space |
        | <pointer-sized alignment>    |
        |             . . .            |
        | 2nd stack-assigned result    |
        + 1st stack-assigned result    |
        | <pointer-sized alignment>    |
        |             . . .            |
        | 2nd stack-assigned argument  |
        | 1st stack-assigned argument  |
        | stack-assigned receiver      |
prev()  +------------------------------+ (Previous Frame)
                Return PC              |
size()  -------------------------------|
               Saved RBP               |
offs()  -------------------------------|
           1th Reserved Registers      |
        -------------------------------|
           2th Reserved Registers      |
        -------------------------------|
           Local Variables             |
    RSP -------------------------------|↓ lower addresses
*/

const zeroRegGo = x64.XMM15

var iregOrderGo = [...]Register64{
	x64.RAX, // RDI
	x64.RBX, // RSI
	x64.RCX, // RDX
	x64.RDI, // RCX
	x64.RSI, // R8
	x64.R8,  // R9
	x64.R9,
	x64.R10,
	x64.R11,
}

var xregOrderGo = [...]XMMRegister{
	x64.XMM0,
	x64.XMM1,
	x64.XMM2,
	x64.XMM3,
	x64.XMM4,
	x64.XMM5,
	x64.XMM6,
	x64.XMM7,
	x64.XMM8,
	x64.XMM9,
	x64.XMM10,
	x64.XMM11,
	x64.XMM12,
	x64.XMM13,
	x64.XMM14,
}

func ReservedRegs(callc bool) []Register {
	if callc {
		return nil
	}
	return []Register{
		R14, // current goroutine
		R15, // GOT reference
	}
}

type stackAlloc struct {
	s uint32
	i int
	x int
}

func (self *stackAlloc) reset() {
	self.i, self.x = 0, 0
}

func (self *stackAlloc) ireg(vt reflect.Type) (p Parameter) {
	p = mkIReg(vt, iregOrderGo[self.i])
	self.i++
	return
}

func (self *stackAlloc) xreg(vt reflect.Type) (p Parameter) {
	p = mkXReg(vt, xregOrderGo[self.x])
	self.x++
	return
}

func (self *stackAlloc) stack(vt reflect.Type) (p Parameter) {
	p = mkStack(vt, self.s)
	self.s += uint32(vt.Size())
	return
}

func (self *stackAlloc) spill(n uint32, a int) uint32 {
	self.s = alignUp(self.s, a) + n
	return self.s
}

func (self *stackAlloc) alloc(p []Parameter, vt reflect.Type) []Parameter {
	nb := vt.Size()
	vk := vt.Kind()

	/* zero-sized objects are allocated on stack */
	if nb == 0 {
		return append(p, mkStack(intType, self.s))
	}

	/* check for value type */
	switch vk {
	case reflect.Bool:
		return self.valloc(p, reflect.TypeOf(false))
	case reflect.Int:
		return self.valloc(p, intType)
	case reflect.Int8:
		return self.valloc(p, reflect.TypeOf(int8(0)))
	case reflect.Int16:
		return self.valloc(p, reflect.TypeOf(int16(0)))
	case reflect.Int32:
		return self.valloc(p, reflect.TypeOf(uint32(0)))
	case reflect.Int64:
		return self.valloc(p, reflect.TypeOf(int64(0)))
	case reflect.Uint:
		return self.valloc(p, reflect.TypeOf(uint(0)))
	case reflect.Uint8:
		return self.valloc(p, reflect.TypeOf(uint8(0)))
	case reflect.Uint16:
		return self.valloc(p, reflect.TypeOf(uint16(0)))
	case reflect.Uint32:
		return self.valloc(p, reflect.TypeOf(uint32(0)))
	case reflect.Uint64:
		return self.valloc(p, reflect.TypeOf(uint64(0)))
	case reflect.Uintptr:
		return self.valloc(p, reflect.TypeOf(uintptr(0)))
	case reflect.Float32:
		return self.valloc(p, reflect.TypeOf(float32(0)))
	case reflect.Float64:
		return self.valloc(p, reflect.TypeOf(float64(0)))
	case reflect.Complex64:
		panic("abi: go117: not implemented: complex64")
	case reflect.Complex128:
		panic("abi: go117: not implemented: complex128")
	case reflect.Array:
		panic("abi: go117: not implemented: arrays")
	case reflect.Chan:
		return self.valloc(p, reflect.TypeOf((chan int)(nil)))
	case reflect.Func:
		return self.valloc(p, reflect.TypeOf((func())(nil)))
	case reflect.Map:
		return self.valloc(p, reflect.TypeOf((map[int]int)(nil)))
	case reflect.Ptr:
		return self.valloc(p, reflect.TypeOf((*int)(nil)))
	case reflect.UnsafePointer:
		return self.valloc(p, ptrType)
	case reflect.Interface:
		return self.valloc(p, ptrType, ptrType)
	case reflect.Slice:
		return self.valloc(p, ptrType, intType, intType)
	case reflect.String:
		return self.valloc(p, ptrType, intType)
	case reflect.Struct:
		panic("abi: go117: not implemented: structs")
	default:
		panic("abi: invalid value type")
	}
}

func (self *stackAlloc) valloc(p []Parameter, vts ...reflect.Type) []Parameter {
	for _, vt := range vts {
		enum := isFloat(vt)
		if enum != notFloatKind && self.x < len(xregOrderGo) {
			p = append(p, self.xreg(vt))
		} else if enum == notFloatKind && self.i < len(iregOrderGo) {
			p = append(p, self.ireg(vt))
		} else {
			p = append(p, self.stack(vt))
		}
	}
	return p
}

func NewFunctionLayout(ft reflect.Type) FunctionLayout {
	var sa stackAlloc
	var fn FunctionLayout

	/* assign every arguments */
	for i := 0; i < ft.NumIn(); i++ {
		fn.Args = sa.alloc(fn.Args, ft.In(i))
	}

	/* reset the register counter, and add a pointer alignment field */
	sa.reset()

	/* assign every return value */
	for i := 0; i < ft.NumOut(); i++ {
		fn.Rets = sa.alloc(fn.Rets, ft.Out(i))
	}

	sa.spill(0, PtrAlign)

	/* assign spill slots */
	for i := 0; i < len(fn.Args); i++ {
		if fn.Args[i].InRegister {
			fn.Args[i].Mem = sa.spill(PtrSize, PtrAlign) - PtrSize
		}
	}

	/* add the final pointer alignment field */
	fn.FP = sa.spill(0, PtrAlign)
	return fn
}

func (self *Frame) emitExchangeArgs(p *Program) {
	iregArgs := make([]Parameter, 0, len(self.desc.Args))
	xregArgs := 0
	for _, v := range self.desc.Args {
		if v.InRegister {
			if v.IsFloat != notFloatKind {
				xregArgs += 1
			} else {
				iregArgs = append(iregArgs, v)
			}
		} else {
			panic("not support stack-assigned arguments now")
		}
	}
	if xregArgs > len(xregOrderC) {
		panic("too many arguments, only support at most 8 integer register arguments now")
	}

	switch len(iregArgs) {
	case 0, 1, 2, 3:
		{
			//Fast-Path: when arguments count are less than four, just exchange the registers
			for i := 0; i < len(iregArgs); i++ {
				p.MOVQ(iregOrderGo[i], iregOrderC[i])
			}
		}
	case 4, 5, 6:
		{
			// need to spill 3th ~ regArgs registers before exchange
			for i := 3; i < len(iregArgs); i++ {
				arg := iregArgs[i]
				// pointer args have already been spilled
				if !arg.IsPointer {
					p.MOVQ(iregOrderGo[i], Ptr(RSP, int32(self.Prev()+arg.Mem)))
				}
			}
			p.MOVQ(iregOrderGo[0], iregOrderC[0])
			p.MOVQ(iregOrderGo[1], iregOrderC[1])
			p.MOVQ(iregOrderGo[2], iregOrderC[2])
			for i := 3; i < len(iregArgs); i++ {
				arg := iregArgs[i]
				p.MOVQ(Ptr(RSP, int32(self.Prev()+arg.Mem)), iregOrderC[i])
			}
		}
	default:
		panic("too many arguments, only support at most 6 integer register arguments now")
	}
}

func (self *Frame) emitStackCheck(p *Program, to *Label, maxStack uintptr) {
	p.LEAQ(Ptr(RSP, int32(-(self.Size()+uint32(maxStack)))), R12)
	p.CMPQ(Ptr(R14, _G_stackguard0), R12)
	p.JBE(to)
}

func (self *Frame) StackCheckTextSize() uint32 {
	p := DefaultArch.CreateProgram()
	p.LEAQ(Ptr(RSP, int32(-(self.Size()))), R12)
	p.CMPQ(Ptr(R14, _G_stackguard0), R12)
	to := CreateLabel("")
	p.Link(to)
	p.JBE(to)
	return uint32(len(p.Assemble(0)))
}

func (self *Frame) emitExchangeRets(p *Program) {
	if len(self.desc.Rets) > 1 {
		panic("too many results, only support one result now")
	}
	// store result
	if len(self.desc.Rets) == 1 && !self.desc.Rets[0].InRegister {
		if self.desc.Rets[0].IsFloat == floatKind64 {
			p.MOVSD(xregOrderC[0], self.retv(0))
		} else if self.desc.Rets[0].IsFloat == floatKind32 {
			p.MOVSS(xregOrderC[0], self.retv(0))
		} else {
			p.MOVQ(RAX, self.retv(0))
		}
	}
}

func (self *Frame) emitRestoreRegs(p *Program) {
	// load reserved registers
	for i, r := range ReservedRegs(self.ccall) {
		switch r.(type) {
		case Register64:
			p.MOVQ(self.resv(i), r)
		case XMMRegister:
			p.MOVSD(self.resv(i), r)
		default:
			panic(fmt.Sprintf("unsupported register type %t to reserve", r))
		}
	}
	// zero xmm15 for go abi
	p.XORPS(zeroRegGo, zeroRegGo)
}
