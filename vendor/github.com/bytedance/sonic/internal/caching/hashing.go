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

package caching

import (
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

var (
    V_strhash = rt.UnpackEface(rt.Strhash)
    S_strhash = *(*uintptr)(V_strhash.Value)
)

func StrHash(s string) uint64 {
    if v := rt.Strhash(unsafe.Pointer(&s), 0); v == 0 {
        return 1
    } else {
        return uint64(v)
    }
}
