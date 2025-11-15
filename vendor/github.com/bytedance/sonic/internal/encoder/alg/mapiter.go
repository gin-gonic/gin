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

package alg

import (
	"encoding"
	"reflect"
	"strconv"
	"sync"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
)

type _MapPair struct {
    k string  // when the map key is integer, k is pointed to m
    v unsafe.Pointer
    m [32]byte
}

type MapIterator struct {
    It rt.GoMapIterator     // must be the first field
    kv rt.GoSlice           // slice of _MapPair
    ki int
}

var (
    iteratorPool = sync.Pool{}
    iteratorPair = rt.UnpackType(reflect.TypeOf(_MapPair{}))
)

func init() {
    if unsafe.Offsetof(MapIterator{}.It) != 0 {
        panic("_MapIterator.it is not the first field")
    }
}


func newIterator() *MapIterator {
    if v := iteratorPool.Get(); v == nil {
        return new(MapIterator)
    } else {
        return resetIterator(v.(*MapIterator))
    }
}

func resetIterator(p *MapIterator) *MapIterator {
    p.ki = 0
    p.It = rt.GoMapIterator{}
    p.kv.Len = 0
    return p
}

func (self *MapIterator) at(i int) *_MapPair {
    return (*_MapPair)(unsafe.Pointer(uintptr(self.kv.Ptr) + uintptr(i) * unsafe.Sizeof(_MapPair{})))
}

func (self *MapIterator) add() (p *_MapPair) {
    p = self.at(self.kv.Len)
    self.kv.Len++
    return
}

func (self *MapIterator) data() (p []_MapPair) {
    *(*rt.GoSlice)(unsafe.Pointer(&p)) = self.kv
    return
}

func (self *MapIterator) append(t *rt.GoType, k unsafe.Pointer, v unsafe.Pointer) (err error) {
    p := self.add()
    p.v = v
    tk := t.Kind()

    // followed as `encoding/json/emcode.go:resolveKeyName
    if tk == reflect.String {
        p.k = *(*string)(k)
        return nil
    }

    // check if the key implements the encoding.TextMarshaler interface
    if t.Pack().Implements(vars.EncodingTextMarshalerType) {
        if tk != reflect.Interface {
            return self.appendConcrete(p, t, k)
        } else {
            return self.appendInterface(p, t, k)
        }
    }

    return self.appendGeneric(p, t, tk, k)
}

func (self *MapIterator) appendGeneric(p *_MapPair, t *rt.GoType, v reflect.Kind, k unsafe.Pointer) error {
    switch v {
        case reflect.Int       : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int)(k)), 10))       ; return nil
        case reflect.Int8      : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int8)(k)), 10))      ; return nil
        case reflect.Int16     : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int16)(k)), 10))     ; return nil
        case reflect.Int32     : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int32)(k)), 10))     ; return nil
        case reflect.Int64     : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int64)(k)), 10))     ; return nil
        case reflect.Uint      : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint)(k)), 10))    ; return nil
        case reflect.Uint8     : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint8)(k)), 10))   ; return nil
        case reflect.Uint16    : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint16)(k)), 10))  ; return nil
        case reflect.Uint32    : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint32)(k)), 10))  ; return nil
        case reflect.Uint64    : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint64)(k)), 10))  ; return nil
        case reflect.Uintptr   : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uintptr)(k)), 10)) ; return nil
        case reflect.Bool      : if *(*bool)(k) { p.k = "true" } else { p.k = "false" }; return nil
        default                : return vars.Error_type(t.Pack())
    }
}

func (self *MapIterator) appendConcrete(p *_MapPair, t *rt.GoType, k unsafe.Pointer) error {
    // compiler has already checked that the type implements the encoding.MarshalText interface
    if !t.Indirect() {
        k = *(*unsafe.Pointer)(k)
    }

    // check the TextMarshaler interface
    eface := rt.GoEface{Value: k, Type: t}.Pack()
    e, ok := eface.(encoding.TextMarshaler)
    if !ok {
        return vars.Error_type(t.Pack())
    }

    // check for nil pointer
    if t.Kind() == reflect.Ptr && k == nil {
        p.k = ""
        return nil
    }
    
    out, err := e.MarshalText()
    if err != nil {
        return err
    }
    p.k = rt.Mem2Str(out)
    return nil
}

func (self *MapIterator) appendInterface(p *_MapPair, t *rt.GoType, k unsafe.Pointer) (err error) {
    if len(rt.IfaceType(t).Methods) == 0 {
        panic("unexpected map key type")
    }
    p.k, err = asText(k)
    return
}

func IteratorStop(p *MapIterator) {
    iteratorPool.Put(p)
}

func IteratorNext(p *MapIterator) {
    i := p.ki
    t := &p.It

    /* check for unordered iteration */
    if i < 0 {
        rt.Mapiternext(t)
        return
    }

    /* check for end of iteration */
    if p.ki >= p.kv.Len {
        t.K = nil
        t.V = nil
        return
    }

    /* update the key-value pair, and increase the pointer */
    t.K = unsafe.Pointer(&p.at(p.ki).k)
    t.V = p.at(p.ki).v
    p.ki++
}

func IteratorStart(t *rt.GoMapType, m unsafe.Pointer, fv uint64) (*MapIterator, error) {
    it := newIterator()
    rt.Mapiterinit(t, m, &it.It)
    count := rt.Maplen(m)

    /* check for key-sorting, empty map don't need sorting */
    if count == 0 || (fv & (1<<BitSortMapKeys)) == 0 {
        it.ki = -1
        return it, nil
    }

    /* pre-allocate space if needed */
    if count > it.kv.Cap {
        it.kv = rt.GrowSlice(iteratorPair, it.kv, count)
    }

    /* dump all the key-value pairs */
    for ; it.It.K != nil; rt.Mapiternext(&it.It) {
        if err := it.append(t.Key, it.It.K, it.It.V); err != nil {
            IteratorStop(it)
            return nil, err
        }
    }

    /* sort the keys, map with only 1 item don't need sorting */
    if it.ki = 1; count > 1 {
        radixQsort(it.data(), 0, maxDepth(it.kv.Len))
    }

    /* load the first pair into iterator */
    it.It.V = it.at(0).v
    it.It.K = unsafe.Pointer(&it.at(0).k)
    return it, nil
}

func asText(v unsafe.Pointer) (string, error) {
	text := rt.AssertI2I(rt.UnpackType(vars.EncodingTextMarshalerType), *(*rt.GoIface)(v))
	r, e := (*(*encoding.TextMarshaler)(unsafe.Pointer(&text))).MarshalText()
	return rt.Mem2Str(r), e
}

func IsValidNumber(s string) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and https://www.json.org/img/number.png

	if s == "" {
		return false
	}

	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if s == "" {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// Make sure we are at the end.
	return s == ""
}
