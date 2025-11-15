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
    `unsafe`

    `github.com/bytedance/sonic/loader`
)

//go:nosplit
func pbool(v bool) uintptr {
    return freezeValue(unsafe.Pointer(&v))
}

//go:nosplit
func ptodec(p loader.Function) _Decoder {
    return *(*_Decoder)(unsafe.Pointer(&p))
}

func assert_eq(v int64, exp int64, msg string) {
    if v != exp {
        panic(msg)
    }
}
