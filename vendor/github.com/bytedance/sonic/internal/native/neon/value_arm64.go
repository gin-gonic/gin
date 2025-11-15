
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

package neon

import (
    `unsafe`
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func value(s unsafe.Pointer, n int, p int, v *types.JsonState, flags uint64) (ret int) {
    return __value(s, n, p, v, flags)
}

//go:nosplit
//go:noescape
//goland:noinspection GoUnusedParameter
func __value(s unsafe.Pointer, n int, p int, v *types.JsonState, flags uint64) (ret int)

