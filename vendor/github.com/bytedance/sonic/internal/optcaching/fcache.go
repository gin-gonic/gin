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
	"strings"
	"unicode"
	"unsafe"

	"github.com/bytedance/sonic/internal/envs"
	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/resolver"
	"github.com/bytedance/sonic/internal/rt"
)

const _AlignSize =  32
const _PaddingSize =  32

type FieldLookup interface {
	Set(fields []resolver.FieldMeta)
	Get(name string, caseSensitive bool) int
}

func isAscii(s string) bool {
	for i :=0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func NewFieldLookup(fields []resolver.FieldMeta) FieldLookup {
	var f FieldLookup
	isAsc := true
	n := len(fields)

	// when field name has non-ascii, use the fallback methods to use strings.ToLower
	for _, f := range fields {
		if !isAscii(f.Name) {
			isAsc = false
			break
		}
	}

	if n <= 8 {
		f =  NewSmallFieldMap(n)
	} else if envs.UseFastMap && n <= 128 && isAsc {
		f =   NewNormalFieldMap(n)
	} else {
		f =   NewFallbackFieldMap(n)
	}

	f.Set(fields)
	return f
}

// Map for keys nums max 8, idx is in [0, 8)
type SmallFieldMap struct {
	keys []string
	lowerKeys []string
}

func NewSmallFieldMap (hint int) *SmallFieldMap {
	return &SmallFieldMap{
		keys: make([]string, hint, hint),
		lowerKeys: make([]string, hint, hint),
	}
}

func (self *SmallFieldMap) Set(fields []resolver.FieldMeta) {
	if len(fields) > 8 {
		panic("small field map should use in small struct")
	}

	for i, f := range fields {
		self.keys[i] = f.Name
		self.lowerKeys[i] = strings.ToLower(f.Name)
	}
}

func (self *SmallFieldMap) Get(name string, caseSensitive bool) int {
	for i, k := range self.keys {
		if len(k) == len(name) && k == name {
			return i
		}
	}
	if caseSensitive {
		return -1
	}
	name = strings.ToLower(name)
	for i, k := range self.lowerKeys {
		if len(k) == len(name) && k == name {
			return i
		}
	}
	return -1
}


/*
1. select by the length: 0 ~ 32 and larger lengths
2. simd match the aligned prefix of the keys: 4/8/16/32 bytes or larger keys
3. check the key with strict match
4. check the key with case-insensitive match
5. find the index 

Mem Layout:
     fixed 33 * 5 bytes  165 bytes |||  variable keys  ||| variable lowerkeys
| length metadata array[33] ||| key0.0 | u8 | key0.1 | u8 | ...  || key1.0 | u8 | key1.1 | u8 | ...  ||| lowerkeys info ...

*/

// Map for keys nums max 255, idx is in [0, 255), idx 255 means not found.
// keysoffset
// | metadata | aligned key0 | aligned key1 | ... |
// 1 ~ 8
// 8 ~ 16
// 16 ~ 32
// > 32 keys use the long keys entry lists
// use bytes to reduce GC
type NormalFieldMap struct {
	keys  			[]byte
	longKeys		[]keyEntry
	// offset for lower
	lowOffset	    int
}

type keyEntry struct {
	key 		string
	lowerKey	string
	index		uint
}

func NewNormalFieldMap(n int) *NormalFieldMap {
	return &NormalFieldMap{
	}
}

const _HdrSlot = 33
const _HdrSize = _HdrSlot * 5

// use native SIMD to accelerate it
func (self *NormalFieldMap) Get(name string, caseSensitive bool) int {
	// small keys use native C
	if len(name) <= 32 {
		_ = native.LookupSmallKey
		lowOffset := self.lowOffset
		if caseSensitive {
			lowOffset = -1
		}
		return native.LookupSmallKey(&name, &self.keys, lowOffset);
	}
	return self.getLongKey(name, caseSensitive)
}

func (self *NormalFieldMap) getLongKey(name string, caseSensitive bool) int {
	for _, k := range self.longKeys {
		if len(k.key) != len(name) {
			continue;
		}
		if k.key == name {
			return int(k.index)
		}
	}

	if caseSensitive {
		return -1
	}

	lower := strings.ToLower(name)
	for _, k := range self.longKeys {
		if len(k.key) != len(name) {
			continue;
		}

		if k.lowerKey == lower {
			return int(k.index)
		}
	}
	return -1
}

func (self *NormalFieldMap) Getdouble(name string) int {
	if len(name) > 32 {
		for _, k := range self.longKeys {
			if len(k.key) != len(name) {
				continue;
			}
			if k.key == name {
				return int(k.index)
			}
		}
		return self.getCaseInsensitive(name)
	}

	// check the fixed length keys, not found the target length
	cnt := int(self.keys[5 * len(name)])
	if cnt == 0 {
		return -1
	}
	p := ((*rt.GoSlice)(unsafe.Pointer(&self.keys))).Ptr
	offset := int(*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(5 * len(name) + 1)))) + _HdrSize
	for i := 0; i < cnt; i++ {
		key := rt.Mem2Str(self.keys[offset: offset + len(name)])
		if key == name {
			return int(self.keys[offset + len(name)])
		}
		offset += len(name) + 1
	}

	return self.getCaseInsensitive(name)
}

func (self *NormalFieldMap) getCaseInsensitive(name string) int {
	lower := strings.ToLower(name)
	if len(name) > 32 {
		for _, k := range self.longKeys {
			if len(k.key) != len(name) {
				continue;
			}

			if k.lowerKey == lower {
				return int(k.index)
			}
		}
		return -1
	}

	cnt := int(self.keys[5 * len(name)])
	p := ((*rt.GoSlice)(unsafe.Pointer(&self.keys))).Ptr
	offset := int(*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(5 * len(name) + 1)))) + self.lowOffset
	for i := 0; i < cnt; i++ {
		key := rt.Mem2Str(self.keys[offset: offset + len(name)])
		if key == lower {
			return int(self.keys[offset + len(name)])
		}
		offset += len(name) + 1
	}

	return -1
}

type keysInfo struct {
	counts int
	lenSum int
	offset int
	cur    int
}

func (self *NormalFieldMap) Set(fields []resolver.FieldMeta) {
	if len(fields) <=8 || len(fields) > 128 {
		panic("normal field map should use in small struct")
	}

	// allocate the flat map in []byte
	var keyLenSum [_HdrSlot]keysInfo

	for i := 0; i < _HdrSlot; i++ {
		keyLenSum[i].offset = 0
		keyLenSum[i].counts = 0
		keyLenSum[i].lenSum = 0
		keyLenSum[i].cur = 0
	}

	kvLen := 0
	for _, f := range(fields) {
		len := len(f.Name)
		if len <= 32 {
			kvLen += len + 1 // key + index
			keyLenSum[len].counts++
			keyLenSum[len].lenSum += len + 1
		}

	}

	// add a padding size at last to make it friendly for SIMD.
	self.keys = make([]byte, _HdrSize + 2 * kvLen, _HdrSize + 2 * kvLen + _PaddingSize)
	self.lowOffset = _HdrSize + kvLen

	// initialize all keys offset
	self.keys[0] = byte(keyLenSum[0].counts)
	// self.keys[1:5] = 0 // offset is always zero here.
	i := 1
	p := ((*rt.GoSlice)(unsafe.Pointer(&self.keys))).Ptr
	for i < _HdrSlot {
		keyLenSum[i].offset = keyLenSum[i-1].offset + keyLenSum[i-1].lenSum
		self.keys[i * 5] = byte(keyLenSum[i].counts)
		// write the offset into []byte
		*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(i * 5 + 1))) = int32(keyLenSum[i].offset)
		i += 1

	}

	// fill the key into bytes
	for i, f := range(fields) {
		len := len(f.Name)
		if len <= 32 {
			offset := keyLenSum[len].offset +  keyLenSum[len].cur
			copy(self.keys[_HdrSize + offset: ], f.Name)
			copy(self.keys[self.lowOffset + offset: ], strings.ToLower(f.Name))
			self.keys[_HdrSize + offset + len] = byte(i)
			self.keys[self.lowOffset + offset + len] = byte(i)
			keyLenSum[len].cur += len + 1

		} else {
			self.longKeys = append(self.longKeys, keyEntry{f.Name, strings.ToLower(f.Name), uint(i)})
		}
	}

}

// use hashmap
type FallbackFieldMap struct {
	oders  []string
	inner  map[string]int
	backup map[string]int
}
 
 func NewFallbackFieldMap(n int) *FallbackFieldMap {
	 return &FallbackFieldMap{
		 oders:  make([]string, n, n),
		 inner:  make(map[string]int, n*2),
		 backup: make(map[string]int, n*2),
	 }
 }
 
 func (self *FallbackFieldMap) Get(name string, caseSensitive bool) int {
	 if i, ok := self.inner[name]; ok {
		 return i
	 } else if !caseSensitive {
		 return self.getCaseInsensitive(name)
	 } else {
		return -1
	 }
 }
 
 func (self *FallbackFieldMap) Set(fields []resolver.FieldMeta) {

	for i, f := range(fields) {
		name := f.Name
		self.oders[i] = name
		self.inner[name] = i
	
		/* add the case-insensitive version, prefer the one with smaller field ID */
		key := strings.ToLower(name)
		if v, ok := self.backup[key]; !ok || i < v {
			self.backup[key] = i
		}
	}
 }
 
 func (self *FallbackFieldMap) getCaseInsensitive(name string) int {
	 if i, ok := self.backup[strings.ToLower(name)]; ok {
		 return i
	 } else {
		 return -1
	 }
 }
 