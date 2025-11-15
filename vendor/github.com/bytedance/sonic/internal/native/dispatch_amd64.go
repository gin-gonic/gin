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

package native

import (
    `unsafe`

    `github.com/bytedance/sonic/internal/cpu`
    `github.com/bytedance/sonic/internal/native/avx2`
    `github.com/bytedance/sonic/internal/native/sse`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
)

const MaxFrameSize   uintptr = 400

var (
    S_f64toa uintptr
    S_f32toa uintptr
    S_i64toa uintptr
    S_u64toa uintptr
    S_lspace uintptr
)

var (
    S_quote   uintptr
    S_unquote uintptr
)

var (
    S_value     uintptr
    S_vstring   uintptr
    S_vnumber   uintptr
    S_vsigned   uintptr
    S_vunsigned uintptr
)

var (
    S_skip_one    uintptr
    S_skip_one_fast    uintptr
    S_get_by_path    uintptr
    S_skip_array  uintptr
    S_skip_object uintptr
    S_skip_number uintptr
)

var (
    __Quote func(s unsafe.Pointer, nb int, dp unsafe.Pointer, dn unsafe.Pointer, flags uint64) int

    __Unquote func(s unsafe.Pointer, nb int, dp unsafe.Pointer, ep unsafe.Pointer, flags uint64) int

    __HTMLEscape func(s unsafe.Pointer, nb int, dp unsafe.Pointer, dn unsafe.Pointer) int

    __Value func(s unsafe.Pointer, n int, p int, v unsafe.Pointer, flags uint64) int

    __SkipOne func(s unsafe.Pointer, p unsafe.Pointer, m unsafe.Pointer, flags uint64) int

    __SkipOneFast func(s unsafe.Pointer, p unsafe.Pointer) int

    __GetByPath func(s unsafe.Pointer, p unsafe.Pointer, path unsafe.Pointer, m unsafe.Pointer) int

    __ValidateOne func(s unsafe.Pointer, p unsafe.Pointer, m unsafe.Pointer, flags uint64) int

    __I64toa func(out unsafe.Pointer, val int64) (ret int)

    __U64toa func(out unsafe.Pointer, val uint64) (ret int)

    __F64toa func(out unsafe.Pointer, val float64) (ret int)

    __F32toa func(out unsafe.Pointer, val float32) (ret int)

    __ValidateUTF8 func(s unsafe.Pointer, p unsafe.Pointer, m unsafe.Pointer) (ret int)

    __ValidateUTF8Fast func(s unsafe.Pointer) (ret int)

	__ParseWithPadding func(parser unsafe.Pointer) (ret int)

	__LookupSmallKey func(key  unsafe.Pointer, table  unsafe.Pointer, lowerOff int) (index int)
)

//go:nosplit
func Quote(s unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int, flags uint64) int {
    return __Quote(rt.NoEscape(unsafe.Pointer(s)), nb, rt.NoEscape(unsafe.Pointer(dp)), rt.NoEscape(unsafe.Pointer(dn)), flags)
}

//go:nosplit
func Unquote(s unsafe.Pointer, nb int, dp unsafe.Pointer, ep *int, flags uint64) int {
    return __Unquote(rt.NoEscape(unsafe.Pointer(s)), nb, rt.NoEscape(unsafe.Pointer(dp)), rt.NoEscape(unsafe.Pointer(ep)), flags)
}

//go:nosplit
func HTMLEscape(s unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int) int {
    return __HTMLEscape(rt.NoEscape(unsafe.Pointer(s)), nb, rt.NoEscape(unsafe.Pointer(dp)), rt.NoEscape(unsafe.Pointer(dn)))
}

//go:nosplit
func Value(s unsafe.Pointer, n int, p int, v *types.JsonState, flags uint64) int {
    return __Value(rt.NoEscape(unsafe.Pointer(s)), n, p, rt.NoEscape(unsafe.Pointer(v)), flags)
}

//go:nosplit
func SkipOne(s *string, p *int, m *types.StateMachine, flags uint64) int {
    return __SkipOne(rt.NoEscape(unsafe.Pointer(s)), rt.NoEscape(unsafe.Pointer(p)), rt.NoEscape(unsafe.Pointer(m)), flags)
}

//go:nosplit
func SkipOneFast(s *string, p *int) int {
    return __SkipOneFast(rt.NoEscape(unsafe.Pointer(s)), rt.NoEscape(unsafe.Pointer(p)))
}

//go:nosplit
func GetByPath(s *string, p *int, path *[]interface{}, m *types.StateMachine) int {
    return __GetByPath(rt.NoEscape(unsafe.Pointer(s)), rt.NoEscape(unsafe.Pointer(p)), rt.NoEscape(unsafe.Pointer(path)), rt.NoEscape(unsafe.Pointer(m)))
}

//go:nosplit
func ValidateOne(s *string, p *int, m *types.StateMachine, flags uint64) int {
    return __ValidateOne(rt.NoEscape(unsafe.Pointer(s)), rt.NoEscape(unsafe.Pointer(p)), rt.NoEscape(unsafe.Pointer(m)), flags)
}

//go:nosplit
func I64toa(out *byte, val int64) (ret int) {
    return __I64toa(rt.NoEscape(unsafe.Pointer(out)), val)
}

//go:nosplit
func U64toa(out *byte, val uint64) (ret int) {
    return __U64toa(rt.NoEscape(unsafe.Pointer(out)), val)
}

//go:nosplit
func F64toa(out *byte, val float64) (ret int) {
    return __F64toa(rt.NoEscape(unsafe.Pointer(out)), val)
}

//go:nosplit
func F32toa(out *byte, val float32) (ret int) {
    return __F32toa(rt.NoEscape(unsafe.Pointer(out)), val)
}

//go:nosplit
func ValidateUTF8(s *string, p *int, m *types.StateMachine) (ret int) {
    return __ValidateUTF8(rt.NoEscape(unsafe.Pointer(s)), rt.NoEscape(unsafe.Pointer(p)), rt.NoEscape(unsafe.Pointer(m)))
}

//go:nosplit
func ValidateUTF8Fast(s *string) (ret int) {
    return __ValidateUTF8Fast(rt.NoEscape(unsafe.Pointer(s)))
}

//go:nosplit
func ParseWithPadding(parser unsafe.Pointer) (ret int) {
    return __ParseWithPadding(rt.NoEscape(unsafe.Pointer(parser)))
}

//go:nosplit
func LookupSmallKey(key *string, table *[]byte, lowerOff int) (index int) {
    return __LookupSmallKey(rt.NoEscape(unsafe.Pointer(key)), rt.NoEscape(unsafe.Pointer(table)), lowerOff)
}

func useSSE() {
    sse.Use()
    S_f64toa      = sse.S_f64toa
    __F64toa      = sse.F_f64toa
    S_f32toa      = sse.S_f32toa
    __F32toa      = sse.F_f32toa
    S_i64toa      = sse.S_i64toa
    __I64toa      = sse.F_i64toa
    S_u64toa      = sse.S_u64toa
    __U64toa      = sse.F_u64toa
    S_lspace      = sse.S_lspace
    S_quote       = sse.S_quote
    __Quote       = sse.F_quote
    S_unquote     = sse.S_unquote
    __Unquote     = sse.F_unquote
    S_value       = sse.S_value
    __Value       = sse.F_value
    S_vstring     = sse.S_vstring
    S_vnumber     = sse.S_vnumber
    S_vsigned     = sse.S_vsigned
    S_vunsigned   = sse.S_vunsigned
    S_skip_one    = sse.S_skip_one
    __SkipOne     = sse.F_skip_one
    __SkipOneFast = sse.F_skip_one_fast
    S_skip_array  = sse.S_skip_array
    S_skip_object = sse.S_skip_object
    S_skip_number = sse.S_skip_number
    S_get_by_path = sse.S_get_by_path
    __GetByPath   = sse.F_get_by_path
    __HTMLEscape  = sse.F_html_escape
    __ValidateOne = sse.F_validate_one
    __ValidateUTF8= sse.F_validate_utf8
    __ValidateUTF8Fast = sse.F_validate_utf8_fast
	__ParseWithPadding = sse.F_parse_with_padding
	__LookupSmallKey = sse.F_lookup_small_key
}

func useAVX2() {
    avx2.Use()
    S_f64toa      = avx2.S_f64toa
    __F64toa      = avx2.F_f64toa
    S_f32toa      = avx2.S_f32toa
    __F32toa      = avx2.F_f32toa
    S_i64toa      = avx2.S_i64toa
    __I64toa      = avx2.F_i64toa
    S_u64toa      = avx2.S_u64toa
    __U64toa      = avx2.F_u64toa
    S_lspace      = avx2.S_lspace
    S_quote       = avx2.S_quote
    __Quote       = avx2.F_quote
    S_unquote     = avx2.S_unquote
    __Unquote     = avx2.F_unquote
    S_value       = avx2.S_value
    __Value       = avx2.F_value
    S_vstring     = avx2.S_vstring
    S_vnumber     = avx2.S_vnumber
    S_vsigned     = avx2.S_vsigned
    S_vunsigned   = avx2.S_vunsigned
    S_skip_one    = avx2.S_skip_one
    __SkipOne     = avx2.F_skip_one
    __SkipOneFast = avx2.F_skip_one_fast
    S_skip_array  = avx2.S_skip_array
    S_skip_object = avx2.S_skip_object
    S_skip_number = avx2.S_skip_number
    S_get_by_path = avx2.S_get_by_path
    __GetByPath   = avx2.F_get_by_path
    __HTMLEscape  = avx2.F_html_escape
    __ValidateOne = avx2.F_validate_one
    __ValidateUTF8= avx2.F_validate_utf8
    __ValidateUTF8Fast = avx2.F_validate_utf8_fast
	__ParseWithPadding = avx2.F_parse_with_padding
	__LookupSmallKey = avx2.F_lookup_small_key
}


func init() {
	if cpu.HasAVX2 {
		useAVX2()
	} else if cpu.HasSSE {
		useSSE()
	} else {
		panic("Unsupported CPU, lacks of AVX2 or SSE CPUID Flag. maybe it's too old to run Sonic.")
	}
}
