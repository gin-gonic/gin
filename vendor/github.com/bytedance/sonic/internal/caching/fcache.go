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
    `strings`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

type FieldMap struct {
    N uint64
    b unsafe.Pointer
    m map[string]int
}

type FieldEntry struct {
    ID   int
    Name string
    Hash uint64
}

const (
    FieldMap_N     = int64(unsafe.Offsetof(FieldMap{}.N))
    FieldMap_b     = int64(unsafe.Offsetof(FieldMap{}.b))
	FieldEntrySize = int64(unsafe.Sizeof(FieldEntry{}))
)

func newBucket(n int) unsafe.Pointer {
    v := make([]FieldEntry, n)
    return (*rt.GoSlice)(unsafe.Pointer(&v)).Ptr
}

func CreateFieldMap(n int) *FieldMap {
    return &FieldMap {
        N: uint64(n * 2),
        b: newBucket(n * 2),    // LoadFactor = 0.5
        m: make(map[string]int, n * 2),
    }
}

func (self *FieldMap) At(p uint64) *FieldEntry {
    off := uintptr(p) * uintptr(FieldEntrySize)
    return (*FieldEntry)(unsafe.Pointer(uintptr(self.b) + off))
}

// Get searches FieldMap by name. JIT generated assembly does NOT call this
// function, rather it implements its own version directly in assembly. So
// we must ensure this function stays in sync with the JIT generated one.
func (self *FieldMap) Get(name string) int {
    h := StrHash(name)
    p := h % self.N
    s := self.At(p)

    /* find the element;
     * the hash map is never full, so the loop will always terminate */
    for s.Hash != 0 {
        if s.Hash == h && s.Name == name {
            return s.ID
        } else {
            p = (p + 1) % self.N
            s = self.At(p)
        }
    }

    /* not found */
    return -1
}

func (self *FieldMap) Set(name string, i int) {
    h := StrHash(name)
    p := h % self.N
    s := self.At(p)

    /* searching for an empty slot;
     * the hash map is never full, so the loop will always terminate */
    for s.Hash != 0 {
        p = (p + 1) % self.N
        s = self.At(p)
    }

    /* set the value */
    s.ID   = i
    s.Hash = h
    s.Name = name

    /* add the case-insensitive version, prefer the one with smaller field ID */
    key := strings.ToLower(name)
    if v, ok := self.m[key]; !ok || i < v {
        self.m[key] = i
    }
}

func (self *FieldMap) GetCaseInsensitive(name string) int {
    if i, ok := self.m[strings.ToLower(name)]; ok {
        return i
    } else {
        return -1
    }
}
