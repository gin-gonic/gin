/**
 * Copyright 2023 ByteDance Inc.
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

package ast

import (
	"sort"
	"unsafe"

	"github.com/bytedance/sonic/internal/caching"
)

type nodeChunk [_DEFAULT_NODE_CAP]Node

type linkedNodes struct {
    head   nodeChunk
    tail   []*nodeChunk
    size   int
}

func (self *linkedNodes) Cap() int {
    if self == nil {
        return 0
    }
    return (len(self.tail)+1)*_DEFAULT_NODE_CAP 
}

func (self *linkedNodes) Len() int {
    if self == nil {
        return 0
    }
    return self.size 
}

func (self *linkedNodes) At(i int) (*Node) {
    if self == nil {
        return nil
    }
    if i >= 0 && i<self.size && i < _DEFAULT_NODE_CAP {
        return &self.head[i]
    } else if i >= _DEFAULT_NODE_CAP && i<self.size  {
        a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
        if a < len(self.tail) {
            return &self.tail[a][b]
        }
    }
    return nil
}

func (self *linkedNodes) MoveOne(source int,  target int) {
    if source == target {
        return
    }
    if source < 0 || source >= self.size || target < 0 || target >= self.size {
        return
    }
    // reserve source
    n := *self.At(source)
    if source < target {
        // move every element (source,target] one step back
        for i:=source; i<target; i++ {
            *self.At(i) = *self.At(i+1)
        } 
    } else {
        // move every element [target,source) one step forward
        for i:=source; i>target; i-- {
            *self.At(i) = *self.At(i-1)
        }
    } 
    // set target
    *self.At(target) = n
}

func (self *linkedNodes) Pop() {
    if self == nil || self.size == 0 {
        return
    }
    self.Set(self.size-1, Node{})
    self.size--
}

func (self *linkedNodes) Push(v Node) {
    self.Set(self.size, v)
}


func (self *linkedNodes) Set(i int, v Node) {
    if i < _DEFAULT_NODE_CAP {
        self.head[i] = v
        if self.size <= i {
            self.size = i+1
        }
        return
    }
    a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
    if a < 0 {
        self.head[b] = v
    } else {
        self.growTailLength(a+1)
        var n = &self.tail[a]
        if *n == nil {
            *n = new(nodeChunk)
        }
        (*n)[b] = v
    }
    if self.size <= i {
        self.size = i+1
    }
}

func (self *linkedNodes) growTailLength(l int) {
    if l <= len(self.tail) {
        return
    }
    c := cap(self.tail)
    for c < l {
        c += 1 + c>>_APPEND_GROW_SHIFT
    }
    if c == cap(self.tail) {
        self.tail = self.tail[:l]
        return
    }
    tmp := make([]*nodeChunk, l, c)
    copy(tmp, self.tail)
    self.tail = tmp
}

func (self *linkedNodes) ToSlice(con []Node) {
    if len(con) < self.size {
        return
    }
    i := (self.size-1)
    a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
    if a < 0 {
        copy(con, self.head[:b+1])
        return
    } else {
        copy(con, self.head[:])
        con = con[_DEFAULT_NODE_CAP:]
    }

    for i:=0; i<a; i++ {
        copy(con, self.tail[i][:])
        con = con[_DEFAULT_NODE_CAP:]
    }
    copy(con, self.tail[a][:b+1])
}

func (self *linkedNodes) FromSlice(con []Node) {
    self.size = len(con)
    i := self.size-1
    a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
    if a < 0 {
        copy(self.head[:b+1], con)
        return
    } else {
        copy(self.head[:], con)
        con = con[_DEFAULT_NODE_CAP:]
    }

    if cap(self.tail) <= a {
        c := (a+1) + (a+1)>>_APPEND_GROW_SHIFT
        self.tail = make([]*nodeChunk, a+1, c)
    }
    self.tail = self.tail[:a+1]

    for i:=0; i<a; i++ {
        self.tail[i] = new(nodeChunk)
        copy(self.tail[i][:], con)
        con = con[_DEFAULT_NODE_CAP:]
    }

    self.tail[a] = new(nodeChunk)
    copy(self.tail[a][:b+1], con)
}

type pairChunk [_DEFAULT_NODE_CAP]Pair

type linkedPairs struct {
    index map[uint64]int
    head pairChunk
    tail []*pairChunk
    size int
}

func (self *linkedPairs) BuildIndex() {
    if self.index == nil {
        self.index = make(map[uint64]int, self.size)
    }
    for i:=0; i<self.size; i++ {
        p := self.At(i)
        self.index[p.hash] = i
    }
}

func (self *linkedPairs) Cap() int {
    if self == nil {
        return 0
    }
    return (len(self.tail)+1)*_DEFAULT_NODE_CAP 
}

func (self *linkedPairs) Len() int {
    if self == nil {
        return 0
    }
    return self.size 
}

func (self *linkedPairs) At(i int) *Pair {
    if self == nil {
        return nil
    }
    if i >= 0 && i < _DEFAULT_NODE_CAP && i<self.size {
        return &self.head[i]
    } else if i >= _DEFAULT_NODE_CAP && i<self.size {
        a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
        if a < len(self.tail) {
            return &self.tail[a][b]
        }
    }
    return nil
}

func (self *linkedPairs) Push(v Pair) {
    self.Set(self.size, v)
}

func (self *linkedPairs) Pop() {
    if self == nil || self.size == 0 {
        return
    }
    self.Unset(self.size-1)
    self.size--
}

func (self *linkedPairs) Unset(i int) {
    if self.index != nil {
        p := self.At(i)
        delete(self.index, p.hash)
    }
    self.set(i, Pair{}) 
}

func (self *linkedPairs) Set(i int, v Pair) {
    if self.index != nil {
        h := v.hash
        self.index[h] = i
    }
    self.set(i, v)
}

func (self *linkedPairs) set(i int, v Pair) {
    if i < _DEFAULT_NODE_CAP {
        self.head[i] = v
        if self.size <= i {
            self.size = i+1
        }
        return
    }
    a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
    if a < 0 {
        self.head[b] = v
    } else {
        self.growTailLength(a+1)
        var n = &self.tail[a]
        if *n == nil {
            *n = new(pairChunk)
        }
        (*n)[b] = v
    }
    if self.size <= i {
        self.size = i+1
    }
}

func (self *linkedPairs) growTailLength(l int) {
    if l <= len(self.tail) {
        return
    }
    c := cap(self.tail)
    for c < l {
        c += 1 + c>>_APPEND_GROW_SHIFT
    }
    if c == cap(self.tail) {
        self.tail = self.tail[:l]
        return
    }
    tmp := make([]*pairChunk, l, c)
    copy(tmp, self.tail)
    self.tail = tmp
}

// linear search
func (self *linkedPairs) Get(key string) (*Pair, int) {
    if self.index != nil {
        // fast-path
        i, ok := self.index[caching.StrHash(key)]
        if ok {
            n := self.At(i)
            if n.Key == key {
                return n, i
            }
            // hash conflicts
            goto linear_search
        } else {
            return nil, -1
        }
    }
linear_search:
    for i:=0; i<self.size; i++ {
        if n := self.At(i); n.Key == key {
            return n, i
        }
    }
    return nil, -1
}

func (self *linkedPairs) ToSlice(con []Pair) {
    if len(con) < self.size {
        return
    }
    i := self.size-1
    a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP

    if a < 0 {
        copy(con, self.head[:b+1])
        return
    } else {
        copy(con, self.head[:])
        con = con[_DEFAULT_NODE_CAP:]
    }

    for i:=0; i<a; i++ {
        copy(con, self.tail[i][:])
        con = con[_DEFAULT_NODE_CAP:]
    }
    copy(con, self.tail[a][:b+1])
}

func (self *linkedPairs) ToMap(con map[string]Node) {
    for i:=0; i<self.size; i++ {
        n := self.At(i)
        con[n.Key] = n.Value
    }
}

func (self *linkedPairs) copyPairs(to []Pair, from []Pair, l int) {
    copy(to, from)
    if self.index != nil {
        for i:=0; i<l; i++ {
            // NOTICE: in case of user not pass hash, just cal it
            h := caching.StrHash(from[i].Key)
            from[i].hash = h
            self.index[h] = i
        }
    }
}

func (self *linkedPairs) FromSlice(con []Pair) {
    self.size = len(con)
    i := self.size-1
    a, b := i/_DEFAULT_NODE_CAP-1, i%_DEFAULT_NODE_CAP
    if a < 0 {
        self.copyPairs(self.head[:b+1], con, b+1)
        return
    } else {
        self.copyPairs(self.head[:], con, len(self.head))
        con = con[_DEFAULT_NODE_CAP:]
    }

    if cap(self.tail) <= a {
        c := (a+1) + (a+1)>>_APPEND_GROW_SHIFT
        self.tail = make([]*pairChunk, a+1, c)
    }
    self.tail = self.tail[:a+1]

    for i:=0; i<a; i++ {
        self.tail[i] = new(pairChunk)
        self.copyPairs(self.tail[i][:], con, len(self.tail[i]))
        con = con[_DEFAULT_NODE_CAP:]
    }

    self.tail[a] = new(pairChunk)
    self.copyPairs(self.tail[a][:b+1], con, b+1)
}

func (self *linkedPairs) Less(i, j int) bool {
    return lessFrom(self.At(i).Key, self.At(j).Key, 0)
}

func (self *linkedPairs) Swap(i, j int) {
    a, b := self.At(i), self.At(j)
    if self.index != nil {
        self.index[a.hash] = j
        self.index[b.hash] = i
    }
    *a, *b = *b, *a
}

func (self *linkedPairs) Sort() {
    sort.Stable(self)
}

// Compare two strings from the pos d.
func lessFrom(a, b string, d int) bool {
    l := len(a)
    if l > len(b) {
        l = len(b)
    }
    for i := d; i < l; i++ {
        if a[i] == b[i] {
            continue
        }
        return a[i] < b[i]
    }
    return len(a) < len(b)
}

type parseObjectStack struct {
    parser Parser
    v      linkedPairs
}

type parseArrayStack struct {
    parser Parser
    v      linkedNodes
}

func newLazyArray(p *Parser) Node {
    s := new(parseArrayStack)
    s.parser = *p
    return Node{
        t: _V_ARRAY_LAZY,
        p: unsafe.Pointer(s),
    }
}

func newLazyObject(p *Parser) Node {
    s := new(parseObjectStack)
    s.parser = *p
    return Node{
        t: _V_OBJECT_LAZY,
        p: unsafe.Pointer(s),
    }
}

func (self *Node) getParserAndArrayStack() (*Parser, *parseArrayStack) {
    stack := (*parseArrayStack)(self.p)
    return &stack.parser, stack
}

func (self *Node) getParserAndObjectStack() (*Parser, *parseObjectStack) {
    stack := (*parseObjectStack)(self.p)
    return &stack.parser, stack
}

