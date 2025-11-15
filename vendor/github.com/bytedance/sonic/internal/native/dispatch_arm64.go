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
	"unsafe"

	neon "github.com/bytedance/sonic/internal/native/neon"
	"github.com/bytedance/sonic/internal/native/types"
)

const (
	MaxFrameSize   uintptr = 200
	BufPaddingSize int     = 64
)

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
	S_parse_with_padding uintptr
	S_lookup_small_key uintptr
)

//go:nosplit
//go:noescape
//go:linkname Quote github.com/bytedance/sonic/internal/native/neon.__quote
func Quote(s unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int, flags uint64) int

//go:nosplit
//go:noescape
//go:linkname Unquote github.com/bytedance/sonic/internal/native/neon.__unquote
func Unquote(s unsafe.Pointer, nb int, dp unsafe.Pointer, ep *int, flags uint64) int

//go:nosplit
//go:noescape
//go:linkname HTMLEscape github.com/bytedance/sonic/internal/native/neon.__html_escape
func HTMLEscape(s unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int) int

//go:nosplit
//go:noescape
//go:linkname Value github.com/bytedance/sonic/internal/native/neon.__value
func Value(s unsafe.Pointer, n int, p int, v *types.JsonState, flags uint64) int

//go:nosplit
//go:noescape
//go:linkname SkipOne github.com/bytedance/sonic/internal/native/neon.__skip_one
func SkipOne(s *string, p *int, m *types.StateMachine, flags uint64) int

//go:nosplit
//go:noescape
//go:linkname SkipOneFast github.com/bytedance/sonic/internal/native/neon.__skip_one_fast
func SkipOneFast(s *string, p *int) int

//go:nosplit
//go:noescape
//go:linkname GetByPath github.com/bytedance/sonic/internal/native/neon.__get_by_path
func GetByPath(s *string, p *int, path *[]interface{}, m *types.StateMachine) int

//go:nosplit
//go:noescape
//go:linkname ValidateOne github.com/bytedance/sonic/internal/native/neon.__validate_one
func ValidateOne(s *string, p *int, m *types.StateMachine, flags uint64) int

//go:nosplit
//go:noescape
//go:linkname I64toa github.com/bytedance/sonic/internal/native/neon.__i64toa
func I64toa(out *byte, val int64) (ret int)

//go:nosplit
//go:noescape
//go:linkname U64toa github.com/bytedance/sonic/internal/native/neon.__u64toa
func U64toa(out *byte, val uint64) (ret int)

//go:nosplit
//go:noescape
//go:linkname F64toa github.com/bytedance/sonic/internal/native/neon.__f64toa
func F64toa(out *byte, val float64) (ret int)

//go:nosplit
//go:noescape
//go:linkname F32toa github.com/bytedance/sonic/internal/native/neon.__f32toa
func F32toa(out *byte, val float32) (ret int)

//go:nosplit
//go:noescape
//go:linkname ValidateUTF8 github.com/bytedance/sonic/internal/native/neon.__validate_utf8
func ValidateUTF8(s *string, p *int, m *types.StateMachine) (ret int)

//go:nosplit
//go:noescape
//go:linkname ValidateUTF8Fast github.com/bytedance/sonic/internal/native/neon.__validate_utf8_fast
func ValidateUTF8Fast(s *string) (ret int)

//go:nosplit
//go:noescape
//go:linkname ParseWithPadding github.com/bytedance/sonic/internal/native/neon.__parse_with_padding
func ParseWithPadding(parser unsafe.Pointer) (ret int) 

//go:nosplit
//go:noescape
//go:linkname LookupSmallKey github.com/bytedance/sonic/internal/native/neon.__lookup_small_key
func LookupSmallKey(key *string, table *[]byte, lowerOff int) (index int)


func useNeon() {
	S_f64toa = neon.S_f64toa
	S_f32toa = neon.S_f32toa
	S_i64toa = neon.S_i64toa
	S_u64toa = neon.S_u64toa
	S_lspace = neon.S_lspace
	S_quote = neon.S_quote
	S_unquote = neon.S_unquote
	S_value = neon.S_value
	S_vstring = neon.S_vstring
	S_vnumber = neon.S_vnumber
	S_vsigned = neon.S_vsigned
	S_vunsigned = neon.S_vunsigned
	S_skip_one = neon.S_skip_one
	S_skip_one_fast = neon.S_skip_one_fast
	S_skip_array = neon.S_skip_array
	S_skip_object = neon.S_skip_object
	S_skip_number = neon.S_skip_number
	S_get_by_path = neon.S_get_by_path
	S_parse_with_padding = neon.S_parse_with_padding
	S_lookup_small_key = neon.S_lookup_small_key
}

func init() {
	useNeon()
}
