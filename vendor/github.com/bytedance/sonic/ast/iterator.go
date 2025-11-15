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

package ast

import (
	"fmt"

	"github.com/bytedance/sonic/internal/caching"
	"github.com/bytedance/sonic/internal/native/types"
)

type Pair struct {
    hash  uint64
    Key   string
    Value Node
}

func NewPair(key string, val Node) Pair {
    return Pair{
        hash: caching.StrHash(key),
        Key: key,
        Value: val,
    }
}

// Values returns iterator for array's children traversal
func (self *Node) Values() (ListIterator, error) {
    if err := self.should(types.V_ARRAY); err != nil {
        return ListIterator{}, err
    }
    return self.values(), nil
}

func (self *Node) values() ListIterator {
    return ListIterator{Iterator{p: self}}
}

// Properties returns iterator for object's children traversal
func (self *Node) Properties() (ObjectIterator, error) {
    if err := self.should(types.V_OBJECT); err != nil {
        return ObjectIterator{}, err
    }
    return self.properties(), nil
}

func (self *Node) properties() ObjectIterator {
    return ObjectIterator{Iterator{p: self}}
}

type Iterator struct {
    i int
    p *Node
}

func (self *Iterator) Pos() int {
    return self.i
}

func (self *Iterator) Len() int {
    return self.p.len()
}

// HasNext reports if it is the end of iteration or has error.
func (self *Iterator) HasNext() bool {
    if !self.p.isLazy() {
        return self.p.Valid() && self.i < self.p.len()
    } else if self.p.t == _V_ARRAY_LAZY {
        return self.p.skipNextNode().Valid()
    } else if self.p.t == _V_OBJECT_LAZY {
        pair := self.p.skipNextPair()
        if pair == nil {
            return false
        }
        return pair.Value.Valid()
    }
    return false
}

// ListIterator is specialized iterator for V_ARRAY
type ListIterator struct {
    Iterator
}

// ObjectIterator is specialized iterator for V_ARRAY
type ObjectIterator struct {
    Iterator
}

func (self *ListIterator) next() *Node {
next_start:
    if !self.HasNext() {
        return nil
    } else {
        n := self.p.nodeAt(self.i)
        self.i++
        if !n.Exists() {
            goto next_start
        }
        return n
    }
}

// Next scans through children of underlying V_ARRAY, 
// copies each child to v, and returns .HasNext().
func (self *ListIterator) Next(v *Node) bool {
    n := self.next()
    if n == nil {
        return false
    }
    *v = *n
    return true
}

func (self *ObjectIterator) next() *Pair {
next_start:
    if !self.HasNext() {
        return nil
    } else {
        n := self.p.pairAt(self.i)
        self.i++
        if n == nil || !n.Value.Exists() {
            goto next_start
        }
        return n
    }
}

// Next scans through children of underlying V_OBJECT, 
// copies each child to v, and returns .HasNext().
func (self *ObjectIterator) Next(p *Pair) bool {
    n := self.next()
    if n == nil {
        return false
    }
    *p = *n
    return true
}

// Sequence represents scanning path of single-layer nodes.
// Index indicates the value's order in both V_ARRAY and V_OBJECT json.
// Key is the value's key (for V_OBJECT json only, otherwise it will be nil).
type Sequence struct {
    Index int 
    Key *string
    // Level int
}

// String is string representation of one Sequence
func (s Sequence) String() string {
    k := ""
    if s.Key != nil {
        k = *s.Key
    }
    return fmt.Sprintf("Sequence(%d, %q)", s.Index, k)
}

type Scanner func(path Sequence, node *Node) bool

// ForEach scans one V_OBJECT node's children from JSON head to tail, 
// and pass the Sequence and Node of corresponding JSON value.
//
// Especially, if the node is not V_ARRAY or V_OBJECT,
// the node itself will be returned and Sequence.Index == -1.
// 
// NOTICE: An unset node WON'T trigger sc, but its index still counts into Path.Index
func (self *Node) ForEach(sc Scanner) error {
    if err := self.checkRaw(); err != nil {
        return err
    }
    switch self.itype() {
    case types.V_ARRAY:
        iter, err := self.Values()
        if err != nil {
            return err
        }
        v := iter.next()
        for v != nil {
            if !sc(Sequence{iter.i-1, nil}, v) {
                return nil
            }
            v = iter.next()
        }
    case types.V_OBJECT:
        iter, err := self.Properties()
        if err != nil {
            return err
        }
        v := iter.next()
        for v != nil {
            if !sc(Sequence{iter.i-1, &v.Key}, &v.Value) {
                return nil
            }
            v = iter.next()
        }
    default:
        if self.Check() != nil {
            return self
        }
        sc(Sequence{-1, nil}, self)
    }
    return nil
}
