/**
* Copyright 2023 ByteDance Inc.
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

package loader

import (
	`reflect`
	`unsafe`

	`github.com/bytedance/sonic/loader/internal/abi`
	`github.com/bytedance/sonic/loader/internal/rt`
)

var _C_Redzone = []bool{false, false, false, false}

// CFunc is a function information for C func
type CFunc struct {
	// C function name
	Name     string

	// entry pc relative to entire text segment
	EntryOff uint32

	// function text size in bytes
	TextSize uint32

	// maximum stack depth of the function
	MaxStack uintptr

	// PC->SP delta lists of the function
	Pcsp     [][2]uint32
}

// GoC is the wrapper for Go calls to C
type GoC struct {
	// CName is the name of corresponding C function
	CName     string

	// CEntry points out where to store the entry address of corresponding C function.
	// It won't be set if nil
	CEntry   *uintptr

	// GoFunc is the POINTER of corresponding go stub function. 
	// It is used to generate Go-C ABI conversion wrapper and receive the wrapper's address 
	//   eg. &func(a int, b int) int 
	//     FOR 
	//     int add(int a, int b)
	// It won't be set if nil
	GoFunc   interface{} 
}

// WrapGoC wraps C functions and loader it into Go stubs
func WrapGoC(text []byte, natives []CFunc, stubs []GoC, modulename string, filename string) {
	funcs := make([]Func, len(natives))
	
	// register C funcs
	for i, f := range natives {
		fn := Func{
			Flag: FuncFlag_ASM,
			EntryOff: f.EntryOff,
			TextSize: f.TextSize,
			Name: f.Name,
		}
		if len(f.Pcsp) != 0 {
			fn.Pcsp = (*Pcdata)(unsafe.Pointer(&natives[i].Pcsp))
		}
		// NOTICE: always forbid async preempt
		fn.PcUnsafePoint = &Pcdata{
			{PC: f.TextSize, Val: PCDATA_UnsafePointUnsafe},
		}
		// NOTICE: always refer to first file
		fn.Pcfile = &Pcdata{
			{PC: f.TextSize, Val: 0},
		}
		// NOTICE: always refer to first line
		fn.Pcline = &Pcdata{
			{PC: f.TextSize, Val: 1},
		}
		// NOTICE: copystack need locals stackmap
		fn.PcStackMapIndex = &Pcdata{
			{PC: f.TextSize, Val: 0},
		}
		sm := rt.StackMapBuilder{}
		sm.AddField(false)
		fn.ArgsPointerMaps = sm.Build()
		fn.LocalsPointerMaps = sm.Build()
		funcs[i] = fn
	}
	rets := Load(text, funcs, modulename, []string{filename})

	// got absolute entry address
	native_entry := **(**uintptr)(unsafe.Pointer(&rets[0]))
	// println("native_entry: ", native_entry)

	wraps := make([]Func, 0, len(stubs))
	wrapIds := make([]int, 0, len(stubs))
	code := make([]byte, 0, len(wraps))
	entryOff := uint32(0)

	// register go wrappers
	for i := range stubs {
		for j := range natives {
			if stubs[i].CName != natives[j].Name {
				continue
			}
			
			// calculate corresponding C entry
			pc := uintptr(native_entry + uintptr(natives[j].EntryOff))
			if stubs[i].CEntry != nil {
				*stubs[i].CEntry = pc
			}

			// no need to generate wrapper, continue next
			if stubs[i].GoFunc == nil {
				continue
			}

			// assemble wrapper codes
			layout := abi.NewFunctionLayout(reflect.TypeOf(stubs[i].GoFunc).Elem())
			frame := abi.NewFrame(&layout, _C_Redzone, true) 
			tcode := abi.CallC(pc, frame, natives[j].MaxStack)
			code = append(code, tcode...)
			size := uint32(len(tcode))
		
			fn := Func{
				Flag: FuncFlag_ASM,
				ArgsSize: int32(layout.ArgSize()),
				EntryOff: entryOff,
				TextSize: size,
				Name: stubs[i].CName + "_go",
			}

			// add check-stack and grow-stack texts' pcsp
			fn.Pcsp = &Pcdata{
				{PC: uint32(frame.StackCheckTextSize()), Val: 0},
				{PC: size - uint32(frame.GrowStackTextSize()), Val: int32(frame.Size())},
				{PC: size, Val: 0},
			}
			// NOTICE: always refer to first file
			fn.Pcfile = &Pcdata{
				{PC: size, Val: 0},
			}
			// NOTICE: always refer to first line
			fn.Pcline = &Pcdata{
				{PC: size, Val: 1},
			}
			// NOTICE: always forbid async preempt
			fn.PcUnsafePoint = &Pcdata{
				{PC: size, Val: PCDATA_UnsafePointUnsafe},
			}

			// register pointer stackmaps
			fn.PcStackMapIndex = &Pcdata{
				{PC: size, Val: 0},
			}
			fn.ArgsPointerMaps = frame.ArgPtrs()
			fn.LocalsPointerMaps = frame.LocalPtrs()

			entryOff += size
			wraps = append(wraps, fn)
			wrapIds = append(wrapIds, i)
		}
	}
	gofuncs := Load(code, wraps, modulename+"/go", []string{filename+".go"})

	// set go func value 
	for i := range gofuncs {
		idx := wrapIds[i]
		w := rt.UnpackEface(stubs[idx].GoFunc)
		*(*Function)(w.Value) = gofuncs[i]
	}
}
