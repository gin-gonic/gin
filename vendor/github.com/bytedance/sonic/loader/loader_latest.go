
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

package loader

import (
    `github.com/bytedance/sonic/loader/internal/rt`
)

// LoadFuncs loads only one function as module, and returns the function pointer
//   - text: machine code
//   - funcName: function name
//   - frameSize: stack frame size. 
//   - argSize: argument total size (in bytes)
//   - argPtrs: indicates if a slot (8 Bytes) of arguments memory stores pointer, from low to high
//   - localPtrs: indicates if a slot (8 Bytes) of local variants memory stores pointer, from low to high
// 
// WARN: 
//   - the function MUST has fixed SP offset equaling to this, otherwise it go.gentraceback will fail
//   - the function MUST has only one stack map for all arguments and local variants
func (self Loader) LoadOne(text []byte, funcName string, frameSize int, argSize int, argPtrs []bool, localPtrs []bool, pcdata Pcdata) Function {
    size := uint32(len(text))

    fn := Func{
        Name: funcName,
        TextSize: size,
        ArgsSize: int32(argSize),
    }


    fn.Pcsp = &pcdata

    if self.NoPreempt {
        fn.PcUnsafePoint = &Pcdata{
            {PC: size, Val: PCDATA_UnsafePointUnsafe},
        }
    } else {
        fn.PcUnsafePoint = &Pcdata{
            {PC: size, Val: PCDATA_UnsafePointSafe},
        }
    }

    // NOTICE: suppose the function has only one stack map at index 0
    fn.PcStackMapIndex = &Pcdata{
        {PC: size, Val: 0},
    }

    if argPtrs != nil {
        args := rt.StackMapBuilder{}
        for _, b := range argPtrs {
            args.AddField(b)
        }
        fn.ArgsPointerMaps = args.Build()
    }
    
    if localPtrs != nil {
        locals := rt.StackMapBuilder{}
        for _, b := range localPtrs {
            locals.AddField(b)
        }
        fn.LocalsPointerMaps = locals.Build()
    }

    out := Load(text, []Func{fn}, self.Name + funcName, []string{self.File})
    return out[0]
}

// Load loads given machine codes and corresponding function information into go moduledata
// and returns runnable function pointer
// WARN: this API is experimental, use it carefully
func Load(text []byte, funcs []Func, modulename string, filenames []string) (out []Function) {
    ids := make([]string, len(funcs))
    for i, f := range funcs {
        ids[i] = f.Name
    }
    // generate module data and allocate memory address
    mod := makeModuledata(modulename, filenames, &funcs, text)

    // verify and register the new module
    moduledataverify1(mod)
    registerModule(mod)

    // 
    // encapsulate function address
    out = make([]Function, len(funcs))
    for i, s := range ids {
        for _, f := range funcs {
            if f.Name == s {
                m := uintptr(mod.text + uintptr(f.EntryOff))
                out[i] = Function(&m)
            }
        }
    }
    return 
}
