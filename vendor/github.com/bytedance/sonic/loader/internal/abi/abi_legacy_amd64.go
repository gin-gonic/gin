//go:build !go1.17
// +build !go1.17

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

package abi

import (
	"fmt"
	"reflect"
	"runtime"
)

func ReservedRegs(callc bool) []Register {
	return nil
}

func salloc(p []Parameter, sp uint32, vt reflect.Type) (uint32, []Parameter) {
	switch vt.Kind() {
	case reflect.Bool:
		return sp + 8, append(p, mkStack(reflect.TypeOf(false), sp))
	case reflect.Int:
		return sp + 8, append(p, mkStack(intType, sp))
	case reflect.Int8:
		return sp + 8, append(p, mkStack(reflect.TypeOf(int8(0)), sp))
	case reflect.Int16:
		return sp + 8, append(p, mkStack(reflect.TypeOf(int16(0)), sp))
	case reflect.Int32:
		return sp + 8, append(p, mkStack(reflect.TypeOf(int32(0)), sp))
	case reflect.Int64:
		return sp + 8, append(p, mkStack(reflect.TypeOf(int64(0)), sp))
	case reflect.Uint:
		return sp + 8, append(p, mkStack(reflect.TypeOf(uint(0)), sp))
	case reflect.Uint8:
		return sp + 8, append(p, mkStack(reflect.TypeOf(uint8(0)), sp))
	case reflect.Uint16:
		return sp + 8, append(p, mkStack(reflect.TypeOf(uint16(0)), sp))
	case reflect.Uint32:
		return sp + 8, append(p, mkStack(reflect.TypeOf(uint32(0)), sp))
	case reflect.Uint64:
		return sp + 8, append(p, mkStack(reflect.TypeOf(uint64(0)), sp))
	case reflect.Uintptr:
		return sp + 8, append(p, mkStack(reflect.TypeOf(uintptr(0)), sp))
	case reflect.Float32:
		return sp + 8, append(p, mkStack(reflect.TypeOf(float32(0)), sp))
	case reflect.Float64:
		return sp + 8, append(p, mkStack(reflect.TypeOf(float64(0)), sp))
	case reflect.Complex64:
		panic("abi: go116: not implemented: complex64")
	case reflect.Complex128:
		panic("abi: go116: not implemented: complex128")
	case reflect.Array:
		panic("abi: go116: not implemented: arrays")
	case reflect.Chan:
		return sp + 8, append(p, mkStack(reflect.TypeOf((chan int)(nil)), sp))
	case reflect.Func:
		return sp + 8, append(p, mkStack(reflect.TypeOf((func())(nil)), sp))
	case reflect.Map:
		return sp + 8, append(p, mkStack(reflect.TypeOf((map[int]int)(nil)), sp))
	case reflect.Ptr:
		return sp + 8, append(p, mkStack(reflect.TypeOf((*int)(nil)), sp))
	case reflect.UnsafePointer:
		return sp + 8, append(p, mkStack(ptrType, sp))
	case reflect.Interface:
		return sp + 16, append(p, mkStack(ptrType, sp), mkStack(ptrType, sp+8))
	case reflect.Slice:
		return sp + 24, append(p, mkStack(ptrType, sp), mkStack(intType, sp+8), mkStack(intType, sp+16))
	case reflect.String:
		return sp + 16, append(p, mkStack(ptrType, sp), mkStack(intType, sp+8))
	case reflect.Struct:
		panic("abi: go116: not implemented: structs")
	default:
		panic("abi: invalid value type")
	}
}

func NewFunctionLayout(ft reflect.Type) FunctionLayout {
	var sp uint32
	var fn FunctionLayout

	/* assign every arguments */
	for i := 0; i < ft.NumIn(); i++ {
		sp, fn.Args = salloc(fn.Args, sp, ft.In(i))
	}

	/* assign every return value */
	for i := 0; i < ft.NumOut(); i++ {
		sp, fn.Rets = salloc(fn.Rets, sp, ft.Out(i))
	}

	/* update function ID and stack pointer */
	fn.FP = sp
	return fn
}

func (self *Frame) emitExchangeArgs(p *Program) {
	iregArgs, xregArgs := 0, 0
	for _, v := range self.desc.Args {
		if v.IsFloat != notFloatKind {
			xregArgs += 1
		} else {
			iregArgs += 1
		}
	}

	if iregArgs > len(iregOrderC) {
		panic("too many arguments, only support at most 6 integer arguments now")
	}
	if xregArgs > len(xregOrderC) {
		panic("too many arguments, only support at most 8 float arguments now")
	}

	ic, xc := iregArgs, xregArgs
	for i := 0; i < len(self.desc.Args); i++ {
		arg := self.desc.Args[i]
		if arg.IsFloat == floatKind64 {
			p.MOVSD(self.argv(i), xregOrderC[xregArgs-xc])
			xc -= 1
		} else if arg.IsFloat == floatKind32 {
			p.MOVSS(self.argv(i), xregOrderC[xregArgs-xc])
			xc -= 1
		} else {
			p.MOVQ(self.argv(i), iregOrderC[iregArgs-ic])
			ic -= 1
		}
	}
}

func (self *Frame) emitStackCheck(p *Program, to *Label, maxStack uintptr) {
	// get the current goroutine
	switch runtime.GOOS {
	case "linux":
		p.MOVQ(Abs(-8), R14).FS()
	case "darwin":
		p.MOVQ(Abs(0x30), R14).GS()
	case "windows":
		break // windows always stores G pointer at R14
	default:
		panic("unsupported operating system")
	}

	// check the stack guard
	p.LEAQ(Ptr(RSP, -int32(self.Size()+uint32(maxStack))), RAX)
	p.CMPQ(Ptr(R14, _G_stackguard0), RAX)
	p.JBE(to)
}

func (self *Frame) StackCheckTextSize() uint32 {
	p := DefaultArch.CreateProgram()

	// get the current goroutine
	switch runtime.GOOS {
	case "linux":
		p.MOVQ(Abs(-8), R14).FS()
	case "darwin":
		p.MOVQ(Abs(0x30), R14).GS()
	case "windows":
		break // windows always stores G pointer at R14
	default:
		panic("unsupported operating system")
	}

	// check the stack guard
	p.LEAQ(Ptr(RSP, -int32(self.Size())), RAX)
	p.CMPQ(Ptr(R14, _G_stackguard0), RAX)
	l := CreateLabel("")
	p.Link(l)
	p.JBE(l)

	return uint32(len(p.Assemble(0)))
}

func (self *Frame) emitExchangeRets(p *Program) {
	if len(self.desc.Rets) > 1 {
		panic("too many results, only support one result now")
	}
	// store result
	if len(self.desc.Rets) == 1 {
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
}
