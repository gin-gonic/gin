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
	`os`
	`runtime`
	`runtime/debug`
	`strings`

	`github.com/bytedance/sonic/internal/jit`
)


var (
    debugSyncGC  = os.Getenv("SONIC_SYNC_GC") != ""
    debugAsyncGC = os.Getenv("SONIC_NO_ASYNC_GC") == ""
)

var (
    _Instr_End _Instr = newInsOp(_OP_nil_1)

    _F_gc       = jit.Func(runtime.GC)
    _F_force_gc = jit.Func(debug.FreeOSMemory)
    _F_println  = jit.Func(println_wrapper)
    _F_print    = jit.Func(print)
)

func println_wrapper(i int, op1 int, op2 int){
    println(i, " Intrs ", op1, _OpNames[op1], "next: ", op2, _OpNames[op2])
}

func print(i int){
    println(i)
}

func (self *_Assembler) force_gc() {
    self.call_go(_F_gc)
    self.call_go(_F_force_gc)
}

func (self *_Assembler) debug_instr(i int, v *_Instr) {
    if debugSyncGC {
        if (i+1 == len(self.p)) {
            self.print_gc(i, v, &_Instr_End) 
        } else {
            next := &(self.p[i+1])
            self.print_gc(i, v, next)
            name := _OpNames[next.op()]
            if strings.Contains(name, "save") {
                return
            }
        }
        self.force_gc()
    }
}
