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

package resolver

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
    _ "unsafe"
)

type FieldOpts int
type OffsetType int

const (
    F_omitempty FieldOpts = 1 << iota
    F_stringize
    F_omitzero
)

const (
    F_offset OffsetType = iota
    F_deref
)

type Offset struct {
    Size uintptr
    Kind OffsetType
    Type reflect.Type
}

type FieldMeta struct {
    Name string
    Path []Offset
    Opts FieldOpts
    Type reflect.Type
    IsZero func(reflect.Value) bool
}

func (self *FieldMeta) String() string {
    var path []string
    var opts []string

    /* dump the field path */
    for _, off := range self.Path {
        if off.Kind == F_offset {
            path = append(path, fmt.Sprintf("%d", off.Size))
        } else {
            path = append(path, fmt.Sprintf("%d.(*%s)", off.Size, off.Type))
        }
    }

    /* check for "string" */
    if (self.Opts & F_stringize) != 0 {
        opts = append(opts, "string")
    }

    /* check for "omitempty" */
    if (self.Opts & F_omitempty) != 0 {
        opts = append(opts, "omitempty")
    }

    /* format the field */
    return fmt.Sprintf(
        "{Field \"%s\" @ %s, opts=%s, type=%s}",
        self.Name,
        strings.Join(path, "."),
        strings.Join(opts, ","),
        self.Type,
    )
}

func (self *FieldMeta) optimize() {
    var n int
    var v uintptr

    /* merge adjacent offsets */
    for _, o := range self.Path {
        if v += o.Size; o.Kind == F_deref {
            self.Path[n].Size    = v
            self.Path[n].Type, v = o.Type, 0
            self.Path[n].Kind, n = F_deref, n + 1
        }
    }

    /* last offset value */
    if v != 0 {
        self.Path[n].Size = v
        self.Path[n].Type = nil
        self.Path[n].Kind = F_offset
        n++
    }

    /* must be at least 1 offset */
    if n != 0 {
        self.Path = self.Path[:n]
    } else {
        self.Path = []Offset{{Kind: F_offset}}
    }
}

func resolveFields(vt reflect.Type) []FieldMeta {
    tfv := typeFields(vt)
    ret := []FieldMeta(nil)

    /* convert each field */
    for _, fv := range tfv.list {
        /* add to result */
        ret = append(ret, FieldMeta{})
        fm := &ret[len(ret)-1]

        item := vt
        path := []Offset(nil)

        /* check for "string" */
        if fv.quoted {
            fm.Opts |= F_stringize
        }

        /* check for "omitempty" */
        if fv.omitEmpty {
            fm.Opts |= F_omitempty
        }

        /* handle the "omitzero" */
        handleOmitZero(fv, fm)

        /* dump the field path */
        for _, i := range fv.index {
            kind := F_offset
            fval := item.Field(i)
            item  = fval.Type

            /* deref the pointer if needed */
            if item.Kind() == reflect.Ptr {
                kind = F_deref
                item = item.Elem()
            }

            /* add to path */
            path = append(path, Offset {
                Kind: kind,
                Type: item,
                Size: fval.Offset,
            })
        }

        /* get the index to the last offset */
        idx := len(path) - 1
        fvt := path[idx].Type

        /* do not dereference into fields */
        if path[idx].Kind == F_deref {
            fvt = reflect.PtrTo(fvt)
            path[idx].Kind = F_offset
        }

        fm.Type = fvt
        fm.Path = path
        fm.Name = fv.name
    }

    /* optimize the offsets */
    for i := range ret {
        ret[i].optimize()
    }

    /* all done */
    return ret
}

var (
    fieldLock  = sync.RWMutex{}
    fieldCache = map[reflect.Type][]FieldMeta{}
)

func ResolveStruct(vt reflect.Type) []FieldMeta {
    var ok bool
    var fm []FieldMeta

    /* attempt to read from cache */
    fieldLock.RLock()
    fm, ok = fieldCache[vt]
    fieldLock.RUnlock()

    /* check if it was cached */
    if ok {
        return fm
    }

    /* otherwise use write-lock */
    fieldLock.Lock()
    defer fieldLock.Unlock()

    /* double check */
    if fm, ok = fieldCache[vt]; ok {
        return fm
    }

    /* resolve the field */
    fm = resolveFields(vt)
    fieldCache[vt] = fm
    return fm
}

func handleOmitZero(fv StdField, fm *FieldMeta) {
    if fv.omitZero {
        fm.Opts |= F_omitzero
        fm.IsZero = fv.isZero
    }
}
