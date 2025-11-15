//
// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package x86_64

import (
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

const (
	_CC_digit = 1 << iota
	_CC_ident
	_CC_ident0
	_CC_number
)

func ispow2(v uint64) bool {
	return (v & (v - 1)) == 0
}

func isdigit(cc rune) bool {
	return '0' <= cc && cc <= '9'
}

func isalpha(cc rune) bool {
	return (cc >= 'a' && cc <= 'z') || (cc >= 'A' && cc <= 'Z')
}

func isident(cc rune) bool {
	return cc == '_' || isalpha(cc) || isdigit(cc)
}

func isident0(cc rune) bool {
	return cc == '_' || isalpha(cc)
}

func isnumber(cc rune) bool {
	return (cc == 'b' || cc == 'B') ||
		(cc == 'o' || cc == 'O') ||
		(cc == 'x' || cc == 'X') ||
		(cc >= '0' && cc <= '9') ||
		(cc >= 'a' && cc <= 'f') ||
		(cc >= 'A' && cc <= 'F')
}

func align(v int, n int) int {
	return (((v - 1) >> n) + 1) << n
}

func append8(m *[]byte, v byte) {
	*m = append(*m, v)
}

func append16(m *[]byte, v uint16) {
	p := len(*m)
	*m = append(*m, 0, 0)
	binary.LittleEndian.PutUint16((*m)[p:], v)
}

func append32(m *[]byte, v uint32) {
	p := len(*m)
	*m = append(*m, 0, 0, 0, 0)
	binary.LittleEndian.PutUint32((*m)[p:], v)
}

func append64(m *[]byte, v uint64) {
	p := len(*m)
	*m = append(*m, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.LittleEndian.PutUint64((*m)[p:], v)
}

func expandmm(m *[]byte, n int, v byte) {
	sl := (*_GoSlice)(unsafe.Pointer(m))
	nb := sl.len + n

	/* grow as needed */
	if nb > cap(*m) {
		*m = growslice(byteType, *m, nb)
	}

	/* fill the new area */
	memset(unsafe.Pointer(uintptr(sl.ptr)+uintptr(sl.len)), v, uintptr(n))
	sl.len = nb
}

func memset(p unsafe.Pointer, c byte, n uintptr) {
	if c != 0 {
		memsetv(p, c, n)
	} else {
		memclrNoHeapPointers(p, n)
	}
}

func memsetv(p unsafe.Pointer, c byte, n uintptr) {
	for i := uintptr(0); i < n; i++ {
		*(*byte)(unsafe.Pointer(uintptr(p) + i)) = c
	}
}

func literal64(v string) (uint64, error) {
	var nb int
	var ch rune
	var ex error
	var mm [12]byte

	/* unquote the runes */
	for v != "" {
		if ch, _, v, ex = strconv.UnquoteChar(v, '\''); ex != nil {
			return 0, ex
		} else if nb += utf8.EncodeRune(mm[nb:], ch); nb > 8 {
			return 0, errors.New("multi-char constant too large")
		}
	}

	/* convert to uint64 */
	return *(*uint64)(unsafe.Pointer(&mm)), nil
}

var (
	byteWrap = reflect.TypeOf(byte(0))
	byteType = (*_GoType)(efaceOf(byteWrap).ptr)
)

//go:linkname growslice runtime.growslice
func growslice(_ *_GoType, _ []byte, _ int) []byte

//go:noescape
//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(_ unsafe.Pointer, _ uintptr)
