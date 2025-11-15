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

package cpu

import (
    `fmt`
    `os`

    `github.com/klauspost/cpuid/v2`
)

var (
    HasAVX2 = cpuid.CPU.Has(cpuid.AVX2)
    HasSSE = cpuid.CPU.Has(cpuid.SSE)
)

func init() {
    switch v := os.Getenv("SONIC_MODE"); v {
        case ""       : break
        case "auto"   : break
        case "noavx"  : HasAVX2 = false
        // will also disable avx, act as `noavx`, we remain it to make sure forward compatibility
        case "noavx2" : HasAVX2 = false
        default       : panic(fmt.Sprintf("invalid mode: '%s', should be one of 'auto', 'noavx', 'noavx2'", v))
    }
}
