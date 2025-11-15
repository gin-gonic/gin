//go:build !notfastpath && !codec.notfastpath && (notmono || codec.notmono)

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from fastpath.notmono.go.tmpl - DO NOT EDIT.

package codec

import (
	"reflect"
	"slices"
	"sort"
)

type fastpathE[T encDriver] struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoder[T], *encFnInfo, reflect.Value)
}
type fastpathD[T decDriver] struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoder[T], *decFnInfo, reflect.Value)
}
type fastpathEs[T encDriver] [56]fastpathE[T]
type fastpathDs[T decDriver] [56]fastpathD[T]

type fastpathET[T encDriver] struct{}
type fastpathDT[T decDriver] struct{}

func (helperEncDriver[T]) fastpathEList() *fastpathEs[T] {
	var i uint = 0
	var s fastpathEs[T]
	fn := func(v interface{}, fe func(*encoder[T], *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathE[T]{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoder[T]).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoder[T]).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoder[T]).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoder[T]).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoder[T]).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoder[T]).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoder[T]).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoder[T]).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoder[T]).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoder[T]).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoder[T]).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoder[T]).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoder[T]).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoder[T]).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoder[T]).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoder[T]).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoder[T]).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoder[T]).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoder[T]).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoder[T]).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoder[T]).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoder[T]).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoder[T]).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoder[T]).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoder[T]).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoder[T]).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoder[T]).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoder[T]).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoder[T]).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoder[T]).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoder[T]).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoder[T]).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoder[T]).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoder[T]).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoder[T]).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoder[T]).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoder[T]).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoder[T]).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoder[T]).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoder[T]).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoder[T]).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoder[T]).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoder[T]).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoder[T]).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoder[T]).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoder[T]).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoder[T]).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoder[T]).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoder[T]).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoder[T]).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoder[T]).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoder[T]).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoder[T]).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoder[T]).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoder[T]).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoder[T]).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriver[T]) fastpathDList() *fastpathDs[T] {
	var i uint = 0
	var s fastpathDs[T]
	fn := func(v interface{}, fd func(*decoder[T], *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathD[T]{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoder[T]).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoder[T]).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoder[T]).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoder[T]).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoder[T]).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoder[T]).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoder[T]).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoder[T]).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoder[T]).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoder[T]).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoder[T]).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoder[T]).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoder[T]).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoder[T]).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoder[T]).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoder[T]).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoder[T]).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoder[T]).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoder[T]).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoder[T]).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoder[T]).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoder[T]).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoder[T]).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoder[T]).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoder[T]).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoder[T]).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoder[T]).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoder[T]).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoder[T]).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoder[T]).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoder[T]).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoder[T]).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoder[T]).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoder[T]).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoder[T]).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoder[T]).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoder[T]).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoder[T]).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoder[T]).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoder[T]).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoder[T]).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoder[T]).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoder[T]).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoder[T]).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoder[T]).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoder[T]).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoder[T]).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoder[T]).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoder[T]).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoder[T]).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoder[T]).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoder[T]).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoder[T]).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoder[T]).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoder[T]).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoder[T]).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

// -- encode

// -- -- fast path type switch
func (helperEncDriver[T]) fastpathEncodeTypeSwitch(iv interface{}, e *encoder[T]) bool {
	var ft fastpathET[T]
	switch v := iv.(type) {
	case []interface{}:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceIntfV(v, e)
		}
	case []string:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceStringV(v, e)
		}
	case [][]byte:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceBytesV(v, e)
		}
	case []float32:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceFloat32V(v, e)
		}
	case []float64:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceFloat64V(v, e)
		}
	case []uint8:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceUint8V(v, e)
		}
	case []uint64:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceUint64V(v, e)
		}
	case []int:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceIntV(v, e)
		}
	case []int32:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceInt32V(v, e)
		}
	case []int64:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceInt64V(v, e)
		}
	case []bool:
		if v == nil {
			e.e.writeNilArray()
		} else {
			ft.EncSliceBoolV(v, e)
		}
	case map[string]interface{}:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringIntfV(v, e)
		}
	case map[string]string:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringStringV(v, e)
		}
	case map[string][]byte:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringBytesV(v, e)
		}
	case map[string]uint8:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringUint8V(v, e)
		}
	case map[string]uint64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringUint64V(v, e)
		}
	case map[string]int:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringIntV(v, e)
		}
	case map[string]int32:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringInt32V(v, e)
		}
	case map[string]float64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringFloat64V(v, e)
		}
	case map[string]bool:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapStringBoolV(v, e)
		}
	case map[uint8]interface{}:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8IntfV(v, e)
		}
	case map[uint8]string:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8StringV(v, e)
		}
	case map[uint8][]byte:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8BytesV(v, e)
		}
	case map[uint8]uint8:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8Uint8V(v, e)
		}
	case map[uint8]uint64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8Uint64V(v, e)
		}
	case map[uint8]int:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8IntV(v, e)
		}
	case map[uint8]int32:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8Int32V(v, e)
		}
	case map[uint8]float64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8Float64V(v, e)
		}
	case map[uint8]bool:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint8BoolV(v, e)
		}
	case map[uint64]interface{}:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64IntfV(v, e)
		}
	case map[uint64]string:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64StringV(v, e)
		}
	case map[uint64][]byte:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64BytesV(v, e)
		}
	case map[uint64]uint8:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64Uint8V(v, e)
		}
	case map[uint64]uint64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64Uint64V(v, e)
		}
	case map[uint64]int:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64IntV(v, e)
		}
	case map[uint64]int32:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64Int32V(v, e)
		}
	case map[uint64]float64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64Float64V(v, e)
		}
	case map[uint64]bool:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapUint64BoolV(v, e)
		}
	case map[int]interface{}:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntIntfV(v, e)
		}
	case map[int]string:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntStringV(v, e)
		}
	case map[int][]byte:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntBytesV(v, e)
		}
	case map[int]uint8:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntUint8V(v, e)
		}
	case map[int]uint64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntUint64V(v, e)
		}
	case map[int]int:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntIntV(v, e)
		}
	case map[int]int32:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntInt32V(v, e)
		}
	case map[int]float64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntFloat64V(v, e)
		}
	case map[int]bool:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapIntBoolV(v, e)
		}
	case map[int32]interface{}:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32IntfV(v, e)
		}
	case map[int32]string:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32StringV(v, e)
		}
	case map[int32][]byte:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32BytesV(v, e)
		}
	case map[int32]uint8:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32Uint8V(v, e)
		}
	case map[int32]uint64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32Uint64V(v, e)
		}
	case map[int32]int:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32IntV(v, e)
		}
	case map[int32]int32:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32Int32V(v, e)
		}
	case map[int32]float64:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32Float64V(v, e)
		}
	case map[int32]bool:
		if v == nil {
			e.e.writeNilMap()
		} else {
			ft.EncMapInt32BoolV(v, e)
		}
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

// -- -- fast path functions
func (e *encoder[T]) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []interface{}
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]interface{})
	}
	if f.ti.mbs {
		ft.EncAsMapSliceIntfV(v, e)
		return
	}
	ft.EncSliceIntfV(v, e)
}
func (fastpathET[T]) EncSliceIntfV(v []interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		if !e.encodeBuiltin(v[j]) {
			e.encodeR(reflect.ValueOf(v[j]))
		}
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceIntfV(v []interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		if !e.encodeBuiltin(v[j]) {
			e.encodeR(reflect.ValueOf(v[j]))
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []string
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]string)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceStringV(v, e)
		return
	}
	ft.EncSliceStringV(v, e)
}
func (fastpathET[T]) EncSliceStringV(v []string, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeString(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceStringV(v []string, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeString(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v [][]byte
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([][]byte)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceBytesV(v, e)
		return
	}
	ft.EncSliceBytesV(v, e)
}
func (fastpathET[T]) EncSliceBytesV(v [][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeBytes(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceBytesV(v [][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeBytes(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []float32
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]float32)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceFloat32V(v, e)
		return
	}
	ft.EncSliceFloat32V(v, e)
}
func (fastpathET[T]) EncSliceFloat32V(v []float32, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeFloat32(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceFloat32V(v []float32, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat32(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []float64
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]float64)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceFloat64V(v, e)
		return
	}
	ft.EncSliceFloat64V(v, e)
}
func (fastpathET[T]) EncSliceFloat64V(v []float64, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeFloat64(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceFloat64V(v []float64, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeFloat64(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []uint8
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]uint8)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceUint8V(v, e)
		return
	}
	ft.EncSliceUint8V(v, e)
}
func (fastpathET[T]) EncSliceUint8V(v []uint8, e *encoder[T]) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathET[T]) EncAsMapSliceUint8V(v []uint8, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(uint64(v[j]))
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []uint64
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]uint64)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceUint64V(v, e)
		return
	}
	ft.EncSliceUint64V(v, e)
}
func (fastpathET[T]) EncSliceUint64V(v []uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeUint(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceUint64V(v []uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeUint(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []int
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]int)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceIntV(v, e)
		return
	}
	ft.EncSliceIntV(v, e)
}
func (fastpathET[T]) EncSliceIntV(v []int, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeInt(int64(v[j]))
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceIntV(v []int, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []int32
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]int32)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceInt32V(v, e)
		return
	}
	ft.EncSliceInt32V(v, e)
}
func (fastpathET[T]) EncSliceInt32V(v []int32, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeInt(int64(v[j]))
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceInt32V(v []int32, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(int64(v[j]))
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []int64
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]int64)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceInt64V(v, e)
		return
	}
	ft.EncSliceInt64V(v, e)
}
func (fastpathET[T]) EncSliceInt64V(v []int64, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeInt(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceInt64V(v []int64, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeInt(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathET[T]
	var v []bool
	if rv.Kind() == reflect.Array {
		rvGetSlice4Array(rv, &v)
	} else {
		v = rv2i(rv).([]bool)
	}
	if f.ti.mbs {
		ft.EncAsMapSliceBoolV(v, e)
		return
	}
	ft.EncSliceBoolV(v, e)
}
func (fastpathET[T]) EncSliceBoolV(v []bool, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(len(v))
	for j := range v {
		e.c = containerArrayElem
		e.e.WriteArrayElem(j == 0)
		e.e.EncodeBool(v[j])
	}
	e.c = 0
	e.e.WriteArrayEnd()
}
func (fastpathET[T]) EncAsMapSliceBoolV(v []bool, e *encoder[T]) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1) // e.mapStart(len(v) / 2)
	for j := range v {
		if j&1 == 0 { // if j%2 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
		} else {
			e.mapElemValue()
		}
		e.e.EncodeBool(v[j])
	}
	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoder[T]) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathET[T]) EncMapStringIntfV(v map[string]interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			if !e.encodeBuiltin(v[k2]) {
				e.encodeR(reflect.ValueOf(v[k2]))
			}
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			if !e.encodeBuiltin(v2) {
				e.encodeR(reflect.ValueOf(v2))
			}
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathET[T]) EncMapStringStringV(v map[string]string, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeString(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathET[T]) EncMapStringBytesV(v map[string][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeBytes(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeBytes(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathET[T]) EncMapStringUint8V(v map[string]uint8, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathET[T]) EncMapStringUint64V(v map[string]uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeUint(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathET[T]) EncMapStringIntV(v map[string]int, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathET[T]) EncMapStringInt32V(v map[string]int32, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathET[T]) EncMapStringFloat64V(v map[string]float64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathET[T]) EncMapStringBoolV(v map[string]bool, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]string, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeString(k2)
			e.mapElemValue()
			e.e.EncodeBool(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathET[T]) EncMapUint8IntfV(v map[uint8]interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			if !e.encodeBuiltin(v[k2]) {
				e.encodeR(reflect.ValueOf(v[k2]))
			}
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			if !e.encodeBuiltin(v2) {
				e.encodeR(reflect.ValueOf(v2))
			}
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathET[T]) EncMapUint8StringV(v map[uint8]string, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeString(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathET[T]) EncMapUint8BytesV(v map[uint8][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBytes(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBytes(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathET[T]) EncMapUint8Uint8V(v map[uint8]uint8, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathET[T]) EncMapUint8Uint64V(v map[uint8]uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathET[T]) EncMapUint8IntV(v map[uint8]int, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathET[T]) EncMapUint8Int32V(v map[uint8]int32, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathET[T]) EncMapUint8Float64V(v map[uint8]float64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathET[T]) EncMapUint8BoolV(v map[uint8]bool, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint8, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(uint64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathET[T]) EncMapUint64IntfV(v map[uint64]interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			if !e.encodeBuiltin(v[k2]) {
				e.encodeR(reflect.ValueOf(v[k2]))
			}
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			if !e.encodeBuiltin(v2) {
				e.encodeR(reflect.ValueOf(v2))
			}
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathET[T]) EncMapUint64StringV(v map[uint64]string, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeString(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathET[T]) EncMapUint64BytesV(v map[uint64][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBytes(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBytes(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathET[T]) EncMapUint64Uint8V(v map[uint64]uint8, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathET[T]) EncMapUint64Uint64V(v map[uint64]uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeUint(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathET[T]) EncMapUint64IntV(v map[uint64]int, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathET[T]) EncMapUint64Int32V(v map[uint64]int32, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathET[T]) EncMapUint64Float64V(v map[uint64]float64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathET[T]) EncMapUint64BoolV(v map[uint64]bool, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeUint(k2)
			e.mapElemValue()
			e.e.EncodeBool(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathET[T]) EncMapIntIntfV(v map[int]interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			if !e.encodeBuiltin(v[k2]) {
				e.encodeR(reflect.ValueOf(v[k2]))
			}
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			if !e.encodeBuiltin(v2) {
				e.encodeR(reflect.ValueOf(v2))
			}
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathET[T]) EncMapIntStringV(v map[int]string, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathET[T]) EncMapIntBytesV(v map[int][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBytes(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBytes(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathET[T]) EncMapIntUint8V(v map[int]uint8, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathET[T]) EncMapIntUint64V(v map[int]uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathET[T]) EncMapIntIntV(v map[int]int, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathET[T]) EncMapIntInt32V(v map[int]int32, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathET[T]) EncMapIntFloat64V(v map[int]float64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathET[T]) EncMapIntBoolV(v map[int]bool, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathET[T]) EncMapInt32IntfV(v map[int32]interface{}, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			if !e.encodeBuiltin(v[k2]) {
				e.encodeR(reflect.ValueOf(v[k2]))
			}
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			if !e.encodeBuiltin(v2) {
				e.encodeR(reflect.ValueOf(v2))
			}
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathET[T]) EncMapInt32StringV(v map[int32]string, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeString(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathET[T]) EncMapInt32BytesV(v map[int32][]byte, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBytes(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBytes(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathET[T]) EncMapInt32Uint8V(v map[int32]uint8, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(uint64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathET[T]) EncMapInt32Uint64V(v map[int32]uint64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeUint(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathET[T]) EncMapInt32IntV(v map[int32]int, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathET[T]) EncMapInt32Int32V(v map[int32]int32, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v[k2]))
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeInt(int64(v2))
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathET[T]) EncMapInt32Float64V(v map[int32]float64, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeFloat64(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}
func (e *encoder[T]) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathET[T]{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathET[T]) EncMapInt32BoolV(v map[int32]bool, e *encoder[T]) {
	if len(v) == 0 {
		e.e.WriteMapEmpty()
		return
	}
	var i uint
	e.mapStart(len(v))
	if e.h.Canonical {
		v2 := make([]int32, len(v))
		for k := range v {
			v2[i] = k
			i++
		}
		slices.Sort(v2)
		for i, k2 := range v2 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v[k2])
		}
	} else {
		i = 0
		for k2, v2 := range v {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			e.e.EncodeInt(int64(k2))
			e.mapElemValue()
			e.e.EncodeBool(v2)
			i++
		}
	}
	e.c = 0
	e.e.WriteMapEnd()
}

// -- decode

// -- -- fast path type switch
func (helperDecDriver[T]) fastpathDecodeTypeSwitch(iv interface{}, d *decoder[T]) bool {
	var ft fastpathDT[T]
	var changed bool
	var containerLen int
	switch v := iv.(type) {
	case []interface{}:
		ft.DecSliceIntfN(v, d)
	case *[]interface{}:
		var v2 []interface{}
		if v2, changed = ft.DecSliceIntfY(*v, d); changed {
			*v = v2
		}
	case []string:
		ft.DecSliceStringN(v, d)
	case *[]string:
		var v2 []string
		if v2, changed = ft.DecSliceStringY(*v, d); changed {
			*v = v2
		}
	case [][]byte:
		ft.DecSliceBytesN(v, d)
	case *[][]byte:
		var v2 [][]byte
		if v2, changed = ft.DecSliceBytesY(*v, d); changed {
			*v = v2
		}
	case []float32:
		ft.DecSliceFloat32N(v, d)
	case *[]float32:
		var v2 []float32
		if v2, changed = ft.DecSliceFloat32Y(*v, d); changed {
			*v = v2
		}
	case []float64:
		ft.DecSliceFloat64N(v, d)
	case *[]float64:
		var v2 []float64
		if v2, changed = ft.DecSliceFloat64Y(*v, d); changed {
			*v = v2
		}
	case []uint8:
		ft.DecSliceUint8N(v, d)
	case *[]uint8:
		var v2 []uint8
		if v2, changed = ft.DecSliceUint8Y(*v, d); changed {
			*v = v2
		}
	case []uint64:
		ft.DecSliceUint64N(v, d)
	case *[]uint64:
		var v2 []uint64
		if v2, changed = ft.DecSliceUint64Y(*v, d); changed {
			*v = v2
		}
	case []int:
		ft.DecSliceIntN(v, d)
	case *[]int:
		var v2 []int
		if v2, changed = ft.DecSliceIntY(*v, d); changed {
			*v = v2
		}
	case []int32:
		ft.DecSliceInt32N(v, d)
	case *[]int32:
		var v2 []int32
		if v2, changed = ft.DecSliceInt32Y(*v, d); changed {
			*v = v2
		}
	case []int64:
		ft.DecSliceInt64N(v, d)
	case *[]int64:
		var v2 []int64
		if v2, changed = ft.DecSliceInt64Y(*v, d); changed {
			*v = v2
		}
	case []bool:
		ft.DecSliceBoolN(v, d)
	case *[]bool:
		var v2 []bool
		if v2, changed = ft.DecSliceBoolY(*v, d); changed {
			*v = v2
		}
	case map[string]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringIntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]interface{}, decInferLen(containerLen, d.maxInitLen(), 32))
			}
			if containerLen != 0 {
				ft.DecMapStringIntfL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringStringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]string, decInferLen(containerLen, d.maxInitLen(), 32))
			}
			if containerLen != 0 {
				ft.DecMapStringStringL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringBytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string][]byte, decInferLen(containerLen, d.maxInitLen(), 40))
			}
			if containerLen != 0 {
				ft.DecMapStringBytesL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringUint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]uint8, decInferLen(containerLen, d.maxInitLen(), 17))
			}
			if containerLen != 0 {
				ft.DecMapStringUint8L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringUint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]uint64, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapStringUint64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringIntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]int, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapStringIntL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringInt32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]int32, decInferLen(containerLen, d.maxInitLen(), 20))
			}
			if containerLen != 0 {
				ft.DecMapStringInt32L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringFloat64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]float64, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapStringFloat64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[string]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapStringBoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[string]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[string]bool, decInferLen(containerLen, d.maxInitLen(), 17))
			}
			if containerLen != 0 {
				ft.DecMapStringBoolL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8IntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]interface{}, decInferLen(containerLen, d.maxInitLen(), 17))
			}
			if containerLen != 0 {
				ft.DecMapUint8IntfL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8StringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]string, decInferLen(containerLen, d.maxInitLen(), 17))
			}
			if containerLen != 0 {
				ft.DecMapUint8StringL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8BytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8][]byte, decInferLen(containerLen, d.maxInitLen(), 25))
			}
			if containerLen != 0 {
				ft.DecMapUint8BytesL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Uint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]uint8, decInferLen(containerLen, d.maxInitLen(), 2))
			}
			if containerLen != 0 {
				ft.DecMapUint8Uint8L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Uint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]uint64, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapUint8Uint64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8IntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]int, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapUint8IntL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Int32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]int32, decInferLen(containerLen, d.maxInitLen(), 5))
			}
			if containerLen != 0 {
				ft.DecMapUint8Int32L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8Float64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]float64, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapUint8Float64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint8]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint8BoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint8]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint8]bool, decInferLen(containerLen, d.maxInitLen(), 2))
			}
			if containerLen != 0 {
				ft.DecMapUint8BoolL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64IntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]interface{}, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapUint64IntfL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64StringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]string, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapUint64StringL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64BytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64][]byte, decInferLen(containerLen, d.maxInitLen(), 32))
			}
			if containerLen != 0 {
				ft.DecMapUint64BytesL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Uint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]uint8, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapUint64Uint8L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Uint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]uint64, decInferLen(containerLen, d.maxInitLen(), 16))
			}
			if containerLen != 0 {
				ft.DecMapUint64Uint64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64IntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]int, decInferLen(containerLen, d.maxInitLen(), 16))
			}
			if containerLen != 0 {
				ft.DecMapUint64IntL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Int32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]int32, decInferLen(containerLen, d.maxInitLen(), 12))
			}
			if containerLen != 0 {
				ft.DecMapUint64Int32L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64Float64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]float64, decInferLen(containerLen, d.maxInitLen(), 16))
			}
			if containerLen != 0 {
				ft.DecMapUint64Float64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[uint64]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapUint64BoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[uint64]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[uint64]bool, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapUint64BoolL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntIntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]interface{}, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapIntIntfL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntStringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]string, decInferLen(containerLen, d.maxInitLen(), 24))
			}
			if containerLen != 0 {
				ft.DecMapIntStringL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntBytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int][]byte, decInferLen(containerLen, d.maxInitLen(), 32))
			}
			if containerLen != 0 {
				ft.DecMapIntBytesL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntUint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]uint8, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapIntUint8L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntUint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]uint64, decInferLen(containerLen, d.maxInitLen(), 16))
			}
			if containerLen != 0 {
				ft.DecMapIntUint64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntIntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]int, decInferLen(containerLen, d.maxInitLen(), 16))
			}
			if containerLen != 0 {
				ft.DecMapIntIntL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntInt32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]int32, decInferLen(containerLen, d.maxInitLen(), 12))
			}
			if containerLen != 0 {
				ft.DecMapIntInt32L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntFloat64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]float64, decInferLen(containerLen, d.maxInitLen(), 16))
			}
			if containerLen != 0 {
				ft.DecMapIntFloat64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapIntBoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int]bool, decInferLen(containerLen, d.maxInitLen(), 9))
			}
			if containerLen != 0 {
				ft.DecMapIntBoolL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32IntfL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]interface{}:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]interface{}, decInferLen(containerLen, d.maxInitLen(), 20))
			}
			if containerLen != 0 {
				ft.DecMapInt32IntfL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32StringL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]string:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]string, decInferLen(containerLen, d.maxInitLen(), 20))
			}
			if containerLen != 0 {
				ft.DecMapInt32StringL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32BytesL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32][]byte:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32][]byte, decInferLen(containerLen, d.maxInitLen(), 28))
			}
			if containerLen != 0 {
				ft.DecMapInt32BytesL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Uint8L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]uint8:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]uint8, decInferLen(containerLen, d.maxInitLen(), 5))
			}
			if containerLen != 0 {
				ft.DecMapInt32Uint8L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Uint64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]uint64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]uint64, decInferLen(containerLen, d.maxInitLen(), 12))
			}
			if containerLen != 0 {
				ft.DecMapInt32Uint64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32IntL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]int:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]int, decInferLen(containerLen, d.maxInitLen(), 12))
			}
			if containerLen != 0 {
				ft.DecMapInt32IntL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Int32L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]int32:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]int32, decInferLen(containerLen, d.maxInitLen(), 8))
			}
			if containerLen != 0 {
				ft.DecMapInt32Int32L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32Float64L(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]float64:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]float64, decInferLen(containerLen, d.maxInitLen(), 12))
			}
			if containerLen != 0 {
				ft.DecMapInt32Float64L(*v, containerLen, d)
			}
			d.mapEnd()
		}
	case map[int32]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen != containerLenNil {
			if containerLen != 0 {
				ft.DecMapInt32BoolL(v, containerLen, d)
			}
			d.mapEnd()
		}
	case *map[int32]bool:
		if containerLen = d.mapStart(d.d.ReadMapStart()); containerLen == containerLenNil {
			*v = nil
		} else {
			if *v == nil {
				*v = make(map[int32]bool, decInferLen(containerLen, d.maxInitLen(), 5))
			}
			if containerLen != 0 {
				ft.DecMapInt32BoolL(*v, containerLen, d)
			}
			d.mapEnd()
		}
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}

// -- -- fast path functions

func (d *decoder[T]) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]interface{})
		if vv, changed := ft.DecSliceIntfY(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []interface{}
		rvGetSlice4Array(rv, &v)
		ft.DecSliceIntfN(v, d)
	default:
		ft.DecSliceIntfN(rv2i(rv).([]interface{}), d)
	}
}
func (fastpathDT[T]) DecSliceIntfY(v []interface{}, d *decoder[T]) (v2 []interface{}, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []interface{}) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 16)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]interface{}, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, nil))
		}
		d.decode(&v[uint(j)])
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]interface{}{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceIntfN(v []interface{}, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			d.decode(&v[uint(j)])
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]string)
		if vv, changed := ft.DecSliceStringY(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []string
		rvGetSlice4Array(rv, &v)
		ft.DecSliceStringN(v, d)
	default:
		ft.DecSliceStringN(rv2i(rv).([]string), d)
	}
}
func (fastpathDT[T]) DecSliceStringY(v []string, d *decoder[T]) (v2 []string, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []string) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 16)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]string, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, ""))
		}
		v[uint(j)] = d.detach2Str(d.d.DecodeStringAsBytes())
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]string{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceStringN(v []string, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = d.detach2Str(d.d.DecodeStringAsBytes())
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[][]byte)
		if vv, changed := ft.DecSliceBytesY(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v [][]byte
		rvGetSlice4Array(rv, &v)
		ft.DecSliceBytesN(v, d)
	default:
		ft.DecSliceBytesN(rv2i(rv).([][]byte), d)
	}
}
func (fastpathDT[T]) DecSliceBytesY(v [][]byte, d *decoder[T]) (v2 [][]byte, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst [][]byte) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 24)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([][]byte, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, nil))
		}
		v[uint(j)] = bytesOKdbi(d.decodeBytesInto(v[uint(j)], false))
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([][]byte{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceBytesN(v [][]byte, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = bytesOKdbi(d.decodeBytesInto(v[uint(j)], false))
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]float32)
		if vv, changed := ft.DecSliceFloat32Y(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []float32
		rvGetSlice4Array(rv, &v)
		ft.DecSliceFloat32N(v, d)
	default:
		ft.DecSliceFloat32N(rv2i(rv).([]float32), d)
	}
}
func (fastpathDT[T]) DecSliceFloat32Y(v []float32, d *decoder[T]) (v2 []float32, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []float32) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 4)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]float32, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = float32(d.d.DecodeFloat32())
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]float32{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceFloat32N(v []float32, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = float32(d.d.DecodeFloat32())
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]float64)
		if vv, changed := ft.DecSliceFloat64Y(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []float64
		rvGetSlice4Array(rv, &v)
		ft.DecSliceFloat64N(v, d)
	default:
		ft.DecSliceFloat64N(rv2i(rv).([]float64), d)
	}
}
func (fastpathDT[T]) DecSliceFloat64Y(v []float64, d *decoder[T]) (v2 []float64, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []float64) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 8)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]float64, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = d.d.DecodeFloat64()
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]float64{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceFloat64N(v []float64, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = d.d.DecodeFloat64()
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]uint8)
		if vv, changed := ft.DecSliceUint8Y(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []uint8
		rvGetSlice4Array(rv, &v)
		ft.DecSliceUint8N(v, d)
	default:
		ft.DecSliceUint8N(rv2i(rv).([]uint8), d)
	}
}
func (fastpathDT[T]) DecSliceUint8Y(v []uint8, d *decoder[T]) (v2 []uint8, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	if ctyp != valueTypeMap {
		var dbi dBytesIntoState
		v2, dbi = d.decodeBytesInto(v[:len(v):len(v)], false)
		return v2, dbi != dBytesIntoParamOut
	}
	containerLenS := d.mapStart(d.d.ReadMapStart()) * 2
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []uint8) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 1)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]uint8, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]uint8{})
	}
	d.mapEnd()
	return v, changed
}
func (fastpathDT[T]) DecSliceUint8N(v []uint8, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	if ctyp != valueTypeMap {
		d.decodeBytesInto(v[:len(v):len(v)], true)
		return
	}
	containerLenS := d.mapStart(d.d.ReadMapStart()) * 2
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	d.mapEnd()
}

func (d *decoder[T]) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]uint64)
		if vv, changed := ft.DecSliceUint64Y(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []uint64
		rvGetSlice4Array(rv, &v)
		ft.DecSliceUint64N(v, d)
	default:
		ft.DecSliceUint64N(rv2i(rv).([]uint64), d)
	}
}
func (fastpathDT[T]) DecSliceUint64Y(v []uint64, d *decoder[T]) (v2 []uint64, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []uint64) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 8)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]uint64, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = d.d.DecodeUint64()
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]uint64{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceUint64N(v []uint64, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = d.d.DecodeUint64()
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]int)
		if vv, changed := ft.DecSliceIntY(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []int
		rvGetSlice4Array(rv, &v)
		ft.DecSliceIntN(v, d)
	default:
		ft.DecSliceIntN(rv2i(rv).([]int), d)
	}
}
func (fastpathDT[T]) DecSliceIntY(v []int, d *decoder[T]) (v2 []int, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []int) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 8)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]int, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]int{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceIntN(v []int, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]int32)
		if vv, changed := ft.DecSliceInt32Y(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []int32
		rvGetSlice4Array(rv, &v)
		ft.DecSliceInt32N(v, d)
	default:
		ft.DecSliceInt32N(rv2i(rv).([]int32), d)
	}
}
func (fastpathDT[T]) DecSliceInt32Y(v []int32, d *decoder[T]) (v2 []int32, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []int32) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 4)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]int32, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]int32{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceInt32N(v []int32, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]int64)
		if vv, changed := ft.DecSliceInt64Y(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []int64
		rvGetSlice4Array(rv, &v)
		ft.DecSliceInt64N(v, d)
	default:
		ft.DecSliceInt64N(rv2i(rv).([]int64), d)
	}
}
func (fastpathDT[T]) DecSliceInt64Y(v []int64, d *decoder[T]) (v2 []int64, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []int64) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 8)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]int64, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, 0))
		}
		v[uint(j)] = d.d.DecodeInt64()
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]int64{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceInt64N(v []int64, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = d.d.DecodeInt64()
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoder[T]) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	switch rv.Kind() {
	case reflect.Ptr:
		v := rv2i(rv).(*[]bool)
		if vv, changed := ft.DecSliceBoolY(*v, d); changed {
			*v = vv
		}
	case reflect.Array:
		var v []bool
		rvGetSlice4Array(rv, &v)
		ft.DecSliceBoolN(v, d)
	default:
		ft.DecSliceBoolN(rv2i(rv).([]bool), d)
	}
}
func (fastpathDT[T]) DecSliceBoolY(v []bool, d *decoder[T]) (v2 []bool, changed bool) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return nil, v != nil
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	var j int
	fnv := func(dst []bool) { v, changed = dst, true }
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if containerLenS == len(v) {
			} else if containerLenS < 0 || containerLenS > cap(v) {
				if xlen := int(decInferLen(containerLenS, d.maxInitLen(), 1)); xlen <= cap(v) {
					fnv(v[:uint(xlen)])
				} else {
					v2 = make([]bool, uint(xlen))
					copy(v2, v)
					fnv(v2)
				}
			} else {
				fnv(v[:containerLenS])
			}
		}
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j >= len(v) {
			fnv(append(v, false))
		}
		v[uint(j)] = d.d.DecodeBool()
	}
	if j < len(v) {
		fnv(v[:uint(j)])
	} else if j == 0 && v == nil {
		fnv([]bool{})
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
	return v, changed
}
func (fastpathDT[T]) DecSliceBoolN(v []bool, d *decoder[T]) {
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil {
		return
	}
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}
	hasLen := containerLenS >= 0
	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if isArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}
		if j < len(v) {
			v[uint(j)] = d.d.DecodeBool()
		} else {
			d.arrayCannotExpand(len(v), j+1)
			d.swallow()
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}
func (d *decoder[T]) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]interface{})
		if *vp == nil {
			*vp = make(map[string]interface{}, decInferLen(containerLen, d.maxInitLen(), 32))
		}
		if containerLen != 0 {
			ft.DecMapStringIntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringIntfL(rv2i(rv).(map[string]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]interface{} given stream length: ", int64(containerLen))
	}
	var mv interface{}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]string)
		if *vp == nil {
			*vp = make(map[string]string, decInferLen(containerLen, d.maxInitLen(), 32))
		}
		if containerLen != 0 {
			ft.DecMapStringStringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringStringL(rv2i(rv).(map[string]string), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringStringL(v map[string]string, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]string given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = d.detach2Str(d.d.DecodeStringAsBytes())
	}
}
func (d *decoder[T]) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string][]byte)
		if *vp == nil {
			*vp = make(map[string][]byte, decInferLen(containerLen, d.maxInitLen(), 40))
		}
		if containerLen != 0 {
			ft.DecMapStringBytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringBytesL(rv2i(rv).(map[string][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string][]byte given stream length: ", int64(containerLen))
	}
	var mv []byte
	mapGet := !d.h.MapValueReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		v[mk], _ = d.decodeBytesInto(mv, false)
	}
}
func (d *decoder[T]) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]uint8)
		if *vp == nil {
			*vp = make(map[string]uint8, decInferLen(containerLen, d.maxInitLen(), 17))
		}
		if containerLen != 0 {
			ft.DecMapStringUint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringUint8L(rv2i(rv).(map[string]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]uint8 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
}
func (d *decoder[T]) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]uint64)
		if *vp == nil {
			*vp = make(map[string]uint64, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapStringUint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringUint64L(rv2i(rv).(map[string]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]uint64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = d.d.DecodeUint64()
	}
}
func (d *decoder[T]) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]int)
		if *vp == nil {
			*vp = make(map[string]int, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapStringIntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringIntL(rv2i(rv).(map[string]int), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringIntL(v map[string]int, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]int given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
}
func (d *decoder[T]) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]int32)
		if *vp == nil {
			*vp = make(map[string]int32, decInferLen(containerLen, d.maxInitLen(), 20))
		}
		if containerLen != 0 {
			ft.DecMapStringInt32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringInt32L(rv2i(rv).(map[string]int32), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]int32 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
}
func (d *decoder[T]) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]float64)
		if *vp == nil {
			*vp = make(map[string]float64, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapStringFloat64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringFloat64L(rv2i(rv).(map[string]float64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]float64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = d.d.DecodeFloat64()
	}
}
func (d *decoder[T]) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[string]bool)
		if *vp == nil {
			*vp = make(map[string]bool, decInferLen(containerLen, d.maxInitLen(), 17))
		}
		if containerLen != 0 {
			ft.DecMapStringBoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapStringBoolL(rv2i(rv).(map[string]bool), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[string]bool given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.detach2Str(d.d.DecodeStringAsBytes())
		d.mapElemValue()
		v[mk] = d.d.DecodeBool()
	}
}
func (d *decoder[T]) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]interface{})
		if *vp == nil {
			*vp = make(map[uint8]interface{}, decInferLen(containerLen, d.maxInitLen(), 17))
		}
		if containerLen != 0 {
			ft.DecMapUint8IntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8IntfL(rv2i(rv).(map[uint8]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]interface{} given stream length: ", int64(containerLen))
	}
	var mv interface{}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]string)
		if *vp == nil {
			*vp = make(map[uint8]string, decInferLen(containerLen, d.maxInitLen(), 17))
		}
		if containerLen != 0 {
			ft.DecMapUint8StringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8StringL(rv2i(rv).(map[uint8]string), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]string given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = d.detach2Str(d.d.DecodeStringAsBytes())
	}
}
func (d *decoder[T]) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8][]byte)
		if *vp == nil {
			*vp = make(map[uint8][]byte, decInferLen(containerLen, d.maxInitLen(), 25))
		}
		if containerLen != 0 {
			ft.DecMapUint8BytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8BytesL(rv2i(rv).(map[uint8][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8][]byte given stream length: ", int64(containerLen))
	}
	var mv []byte
	mapGet := !d.h.MapValueReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		v[mk], _ = d.decodeBytesInto(mv, false)
	}
}
func (d *decoder[T]) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]uint8)
		if *vp == nil {
			*vp = make(map[uint8]uint8, decInferLen(containerLen, d.maxInitLen(), 2))
		}
		if containerLen != 0 {
			ft.DecMapUint8Uint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Uint8L(rv2i(rv).(map[uint8]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]uint8 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
}
func (d *decoder[T]) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]uint64)
		if *vp == nil {
			*vp = make(map[uint8]uint64, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapUint8Uint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Uint64L(rv2i(rv).(map[uint8]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]uint64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = d.d.DecodeUint64()
	}
}
func (d *decoder[T]) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]int)
		if *vp == nil {
			*vp = make(map[uint8]int, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapUint8IntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8IntL(rv2i(rv).(map[uint8]int), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]int given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
}
func (d *decoder[T]) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]int32)
		if *vp == nil {
			*vp = make(map[uint8]int32, decInferLen(containerLen, d.maxInitLen(), 5))
		}
		if containerLen != 0 {
			ft.DecMapUint8Int32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Int32L(rv2i(rv).(map[uint8]int32), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]int32 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
}
func (d *decoder[T]) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]float64)
		if *vp == nil {
			*vp = make(map[uint8]float64, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapUint8Float64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8Float64L(rv2i(rv).(map[uint8]float64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]float64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = d.d.DecodeFloat64()
	}
}
func (d *decoder[T]) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint8]bool)
		if *vp == nil {
			*vp = make(map[uint8]bool, decInferLen(containerLen, d.maxInitLen(), 2))
		}
		if containerLen != 0 {
			ft.DecMapUint8BoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint8BoolL(rv2i(rv).(map[uint8]bool), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint8]bool given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
		d.mapElemValue()
		v[mk] = d.d.DecodeBool()
	}
}
func (d *decoder[T]) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]interface{})
		if *vp == nil {
			*vp = make(map[uint64]interface{}, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapUint64IntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64IntfL(rv2i(rv).(map[uint64]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]interface{} given stream length: ", int64(containerLen))
	}
	var mv interface{}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]string)
		if *vp == nil {
			*vp = make(map[uint64]string, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapUint64StringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64StringL(rv2i(rv).(map[uint64]string), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]string given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = d.detach2Str(d.d.DecodeStringAsBytes())
	}
}
func (d *decoder[T]) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64][]byte)
		if *vp == nil {
			*vp = make(map[uint64][]byte, decInferLen(containerLen, d.maxInitLen(), 32))
		}
		if containerLen != 0 {
			ft.DecMapUint64BytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64BytesL(rv2i(rv).(map[uint64][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64][]byte given stream length: ", int64(containerLen))
	}
	var mv []byte
	mapGet := !d.h.MapValueReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		v[mk], _ = d.decodeBytesInto(mv, false)
	}
}
func (d *decoder[T]) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]uint8)
		if *vp == nil {
			*vp = make(map[uint64]uint8, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapUint64Uint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Uint8L(rv2i(rv).(map[uint64]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]uint8 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
}
func (d *decoder[T]) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]uint64)
		if *vp == nil {
			*vp = make(map[uint64]uint64, decInferLen(containerLen, d.maxInitLen(), 16))
		}
		if containerLen != 0 {
			ft.DecMapUint64Uint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Uint64L(rv2i(rv).(map[uint64]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]uint64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = d.d.DecodeUint64()
	}
}
func (d *decoder[T]) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]int)
		if *vp == nil {
			*vp = make(map[uint64]int, decInferLen(containerLen, d.maxInitLen(), 16))
		}
		if containerLen != 0 {
			ft.DecMapUint64IntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64IntL(rv2i(rv).(map[uint64]int), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]int given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
}
func (d *decoder[T]) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]int32)
		if *vp == nil {
			*vp = make(map[uint64]int32, decInferLen(containerLen, d.maxInitLen(), 12))
		}
		if containerLen != 0 {
			ft.DecMapUint64Int32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Int32L(rv2i(rv).(map[uint64]int32), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]int32 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
}
func (d *decoder[T]) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]float64)
		if *vp == nil {
			*vp = make(map[uint64]float64, decInferLen(containerLen, d.maxInitLen(), 16))
		}
		if containerLen != 0 {
			ft.DecMapUint64Float64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64Float64L(rv2i(rv).(map[uint64]float64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]float64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = d.d.DecodeFloat64()
	}
}
func (d *decoder[T]) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[uint64]bool)
		if *vp == nil {
			*vp = make(map[uint64]bool, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapUint64BoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapUint64BoolL(rv2i(rv).(map[uint64]bool), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[uint64]bool given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := d.d.DecodeUint64()
		d.mapElemValue()
		v[mk] = d.d.DecodeBool()
	}
}
func (d *decoder[T]) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]interface{})
		if *vp == nil {
			*vp = make(map[int]interface{}, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapIntIntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntIntfL(rv2i(rv).(map[int]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]interface{} given stream length: ", int64(containerLen))
	}
	var mv interface{}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]string)
		if *vp == nil {
			*vp = make(map[int]string, decInferLen(containerLen, d.maxInitLen(), 24))
		}
		if containerLen != 0 {
			ft.DecMapIntStringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntStringL(rv2i(rv).(map[int]string), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntStringL(v map[int]string, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]string given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = d.detach2Str(d.d.DecodeStringAsBytes())
	}
}
func (d *decoder[T]) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int][]byte)
		if *vp == nil {
			*vp = make(map[int][]byte, decInferLen(containerLen, d.maxInitLen(), 32))
		}
		if containerLen != 0 {
			ft.DecMapIntBytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntBytesL(rv2i(rv).(map[int][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int][]byte given stream length: ", int64(containerLen))
	}
	var mv []byte
	mapGet := !d.h.MapValueReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		v[mk], _ = d.decodeBytesInto(mv, false)
	}
}
func (d *decoder[T]) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]uint8)
		if *vp == nil {
			*vp = make(map[int]uint8, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapIntUint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntUint8L(rv2i(rv).(map[int]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]uint8 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
}
func (d *decoder[T]) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]uint64)
		if *vp == nil {
			*vp = make(map[int]uint64, decInferLen(containerLen, d.maxInitLen(), 16))
		}
		if containerLen != 0 {
			ft.DecMapIntUint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntUint64L(rv2i(rv).(map[int]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]uint64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = d.d.DecodeUint64()
	}
}
func (d *decoder[T]) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]int)
		if *vp == nil {
			*vp = make(map[int]int, decInferLen(containerLen, d.maxInitLen(), 16))
		}
		if containerLen != 0 {
			ft.DecMapIntIntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntIntL(rv2i(rv).(map[int]int), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntIntL(v map[int]int, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]int given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
}
func (d *decoder[T]) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]int32)
		if *vp == nil {
			*vp = make(map[int]int32, decInferLen(containerLen, d.maxInitLen(), 12))
		}
		if containerLen != 0 {
			ft.DecMapIntInt32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntInt32L(rv2i(rv).(map[int]int32), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]int32 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
}
func (d *decoder[T]) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]float64)
		if *vp == nil {
			*vp = make(map[int]float64, decInferLen(containerLen, d.maxInitLen(), 16))
		}
		if containerLen != 0 {
			ft.DecMapIntFloat64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntFloat64L(rv2i(rv).(map[int]float64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]float64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = d.d.DecodeFloat64()
	}
}
func (d *decoder[T]) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int]bool)
		if *vp == nil {
			*vp = make(map[int]bool, decInferLen(containerLen, d.maxInitLen(), 9))
		}
		if containerLen != 0 {
			ft.DecMapIntBoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapIntBoolL(rv2i(rv).(map[int]bool), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int]bool given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
		d.mapElemValue()
		v[mk] = d.d.DecodeBool()
	}
}
func (d *decoder[T]) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]interface{})
		if *vp == nil {
			*vp = make(map[int32]interface{}, decInferLen(containerLen, d.maxInitLen(), 20))
		}
		if containerLen != 0 {
			ft.DecMapInt32IntfL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32IntfL(rv2i(rv).(map[int32]interface{}), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]interface{} given stream length: ", int64(containerLen))
	}
	var mv interface{}
	mapGet := !d.h.MapValueReset && !d.h.InterfaceReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		d.decode(&mv)
		v[mk] = mv
	}
}
func (d *decoder[T]) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]string)
		if *vp == nil {
			*vp = make(map[int32]string, decInferLen(containerLen, d.maxInitLen(), 20))
		}
		if containerLen != 0 {
			ft.DecMapInt32StringL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32StringL(rv2i(rv).(map[int32]string), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]string given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = d.detach2Str(d.d.DecodeStringAsBytes())
	}
}
func (d *decoder[T]) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32][]byte)
		if *vp == nil {
			*vp = make(map[int32][]byte, decInferLen(containerLen, d.maxInitLen(), 28))
		}
		if containerLen != 0 {
			ft.DecMapInt32BytesL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32BytesL(rv2i(rv).(map[int32][]byte), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32][]byte given stream length: ", int64(containerLen))
	}
	var mv []byte
	mapGet := !d.h.MapValueReset
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		if mapGet {
			mv = v[mk]
		} else {
			mv = nil
		}
		v[mk], _ = d.decodeBytesInto(mv, false)
	}
}
func (d *decoder[T]) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]uint8)
		if *vp == nil {
			*vp = make(map[int32]uint8, decInferLen(containerLen, d.maxInitLen(), 5))
		}
		if containerLen != 0 {
			ft.DecMapInt32Uint8L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Uint8L(rv2i(rv).(map[int32]uint8), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]uint8 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	}
}
func (d *decoder[T]) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]uint64)
		if *vp == nil {
			*vp = make(map[int32]uint64, decInferLen(containerLen, d.maxInitLen(), 12))
		}
		if containerLen != 0 {
			ft.DecMapInt32Uint64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Uint64L(rv2i(rv).(map[int32]uint64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]uint64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = d.d.DecodeUint64()
	}
}
func (d *decoder[T]) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]int)
		if *vp == nil {
			*vp = make(map[int32]int, decInferLen(containerLen, d.maxInitLen(), 12))
		}
		if containerLen != 0 {
			ft.DecMapInt32IntL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32IntL(rv2i(rv).(map[int32]int), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]int given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	}
}
func (d *decoder[T]) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]int32)
		if *vp == nil {
			*vp = make(map[int32]int32, decInferLen(containerLen, d.maxInitLen(), 8))
		}
		if containerLen != 0 {
			ft.DecMapInt32Int32L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Int32L(rv2i(rv).(map[int32]int32), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]int32 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	}
}
func (d *decoder[T]) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]float64)
		if *vp == nil {
			*vp = make(map[int32]float64, decInferLen(containerLen, d.maxInitLen(), 12))
		}
		if containerLen != 0 {
			ft.DecMapInt32Float64L(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32Float64L(rv2i(rv).(map[int32]float64), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]float64 given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = d.d.DecodeFloat64()
	}
}
func (d *decoder[T]) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDT[T]
	containerLen := d.mapStart(d.d.ReadMapStart())
	if rv.Kind() == reflect.Ptr {
		vp, _ := rv2i(rv).(*map[int32]bool)
		if *vp == nil {
			*vp = make(map[int32]bool, decInferLen(containerLen, d.maxInitLen(), 5))
		}
		if containerLen != 0 {
			ft.DecMapInt32BoolL(*vp, containerLen, d)
		}
	} else if containerLen != 0 {
		ft.DecMapInt32BoolL(rv2i(rv).(map[int32]bool), containerLen, d)
	}
	d.mapEnd()
}
func (fastpathDT[T]) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoder[T]) {
	if v == nil {
		halt.errorInt("cannot decode into nil map[int32]bool given stream length: ", int64(containerLen))
	}
	hasLen := containerLen >= 0
	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		d.mapElemKey(j == 0)
		mk := int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
		d.mapElemValue()
		v[mk] = d.d.DecodeBool()
	}
}
