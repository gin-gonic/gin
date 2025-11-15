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

package jit

import (
    `reflect`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
    `github.com/twitchyliquid64/golang-asm/obj`
)

func Func(f interface{}) obj.Addr {
    if p := rt.UnpackEface(f); p.Type.Kind() != reflect.Func {
        panic("f is not a function")
    } else {
        return Imm(*(*int64)(p.Value))
    }
}

func Type(t reflect.Type) obj.Addr {
    return Gtype(rt.UnpackType(t))
}

func Itab(i *rt.GoType, t reflect.Type) obj.Addr {
    return Imm(int64(uintptr(unsafe.Pointer(rt.GetItab(rt.IfaceType(i), rt.UnpackType(t), false)))))
}

func Gitab(i *rt.GoItab) obj.Addr {
    return Imm(int64(uintptr(unsafe.Pointer(i))))
}

func Gtype(t *rt.GoType) obj.Addr {
    return Imm(int64(uintptr(unsafe.Pointer(t))))
}
