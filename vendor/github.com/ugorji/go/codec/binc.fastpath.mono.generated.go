//go:build !notmono && !codec.notmono  && !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"slices"
	"sort"
)

type fastpathEBincBytes struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderBincBytes, *encFnInfo, reflect.Value)
}
type fastpathDBincBytes struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderBincBytes, *decFnInfo, reflect.Value)
}
type fastpathEsBincBytes [56]fastpathEBincBytes
type fastpathDsBincBytes [56]fastpathDBincBytes
type fastpathETBincBytes struct{}
type fastpathDTBincBytes struct{}

func (helperEncDriverBincBytes) fastpathEList() *fastpathEsBincBytes {
	var i uint = 0
	var s fastpathEsBincBytes
	fn := func(v interface{}, fe func(*encoderBincBytes, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathEBincBytes{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderBincBytes).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderBincBytes).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderBincBytes).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderBincBytes).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderBincBytes).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderBincBytes).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderBincBytes).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderBincBytes).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderBincBytes).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderBincBytes).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderBincBytes).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderBincBytes).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderBincBytes).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderBincBytes).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderBincBytes).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderBincBytes).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderBincBytes).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderBincBytes).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderBincBytes).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderBincBytes).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderBincBytes).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderBincBytes).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderBincBytes).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderBincBytes).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderBincBytes).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderBincBytes).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderBincBytes).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderBincBytes).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderBincBytes).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderBincBytes).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderBincBytes).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderBincBytes).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderBincBytes).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderBincBytes).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderBincBytes).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderBincBytes).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderBincBytes).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderBincBytes).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderBincBytes).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderBincBytes).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderBincBytes).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderBincBytes).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderBincBytes).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderBincBytes).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderBincBytes).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderBincBytes).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderBincBytes).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderBincBytes).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderBincBytes).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderBincBytes).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderBincBytes).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderBincBytes).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderBincBytes).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderBincBytes).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderBincBytes).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderBincBytes).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverBincBytes) fastpathDList() *fastpathDsBincBytes {
	var i uint = 0
	var s fastpathDsBincBytes
	fn := func(v interface{}, fd func(*decoderBincBytes, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDBincBytes{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderBincBytes).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderBincBytes).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderBincBytes).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderBincBytes).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderBincBytes).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderBincBytes).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderBincBytes).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderBincBytes).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderBincBytes).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderBincBytes).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderBincBytes).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderBincBytes).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderBincBytes).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderBincBytes).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderBincBytes).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderBincBytes).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderBincBytes).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderBincBytes).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderBincBytes).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderBincBytes).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderBincBytes).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderBincBytes).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderBincBytes).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderBincBytes).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderBincBytes).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderBincBytes).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderBincBytes).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderBincBytes).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderBincBytes).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderBincBytes).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderBincBytes).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderBincBytes).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderBincBytes).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderBincBytes).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderBincBytes).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderBincBytes).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderBincBytes).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderBincBytes).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderBincBytes).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderBincBytes).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderBincBytes).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderBincBytes).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderBincBytes).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderBincBytes).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderBincBytes).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderBincBytes).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderBincBytes).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderBincBytes).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderBincBytes).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderBincBytes).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderBincBytes).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderBincBytes).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderBincBytes).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderBincBytes).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderBincBytes).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderBincBytes).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverBincBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderBincBytes) bool {
	var ft fastpathETBincBytes
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
		_ = v
		return false
	}
	return true
}

func (e *encoderBincBytes) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceIntfV(v []interface{}, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceIntfV(v []interface{}, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceStringV(v []string, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceStringV(v []string, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceBytesV(v [][]byte, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceBytesV(v [][]byte, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceFloat32V(v []float32, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceFloat32V(v []float32, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceFloat64V(v []float64, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceFloat64V(v []float64, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceUint8V(v []uint8, e *encoderBincBytes) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETBincBytes) EncAsMapSliceUint8V(v []uint8, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceUint64V(v []uint64, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceUint64V(v []uint64, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceIntV(v []int, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceIntV(v []int, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceInt32V(v []int32, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceInt32V(v []int32, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceInt64V(v []int64, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceInt64V(v []int64, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincBytes
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
func (fastpathETBincBytes) EncSliceBoolV(v []bool, e *encoderBincBytes) {
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
func (fastpathETBincBytes) EncAsMapSliceBoolV(v []bool, e *encoderBincBytes) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincBytes) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETBincBytes) EncMapStringIntfV(v map[string]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETBincBytes) EncMapStringStringV(v map[string]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETBincBytes) EncMapStringBytesV(v map[string][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETBincBytes) EncMapStringUint8V(v map[string]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETBincBytes) EncMapStringUint64V(v map[string]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETBincBytes) EncMapStringIntV(v map[string]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETBincBytes) EncMapStringInt32V(v map[string]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETBincBytes) EncMapStringFloat64V(v map[string]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETBincBytes) EncMapStringBoolV(v map[string]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETBincBytes) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETBincBytes) EncMapUint8StringV(v map[uint8]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETBincBytes) EncMapUint8BytesV(v map[uint8][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETBincBytes) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETBincBytes) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETBincBytes) EncMapUint8IntV(v map[uint8]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETBincBytes) EncMapUint8Int32V(v map[uint8]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETBincBytes) EncMapUint8Float64V(v map[uint8]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETBincBytes) EncMapUint8BoolV(v map[uint8]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETBincBytes) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETBincBytes) EncMapUint64StringV(v map[uint64]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETBincBytes) EncMapUint64BytesV(v map[uint64][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETBincBytes) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETBincBytes) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETBincBytes) EncMapUint64IntV(v map[uint64]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETBincBytes) EncMapUint64Int32V(v map[uint64]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETBincBytes) EncMapUint64Float64V(v map[uint64]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETBincBytes) EncMapUint64BoolV(v map[uint64]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETBincBytes) EncMapIntIntfV(v map[int]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETBincBytes) EncMapIntStringV(v map[int]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETBincBytes) EncMapIntBytesV(v map[int][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETBincBytes) EncMapIntUint8V(v map[int]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETBincBytes) EncMapIntUint64V(v map[int]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETBincBytes) EncMapIntIntV(v map[int]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETBincBytes) EncMapIntInt32V(v map[int]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETBincBytes) EncMapIntFloat64V(v map[int]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETBincBytes) EncMapIntBoolV(v map[int]bool, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETBincBytes) EncMapInt32IntfV(v map[int32]interface{}, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETBincBytes) EncMapInt32StringV(v map[int32]string, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETBincBytes) EncMapInt32BytesV(v map[int32][]byte, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETBincBytes) EncMapInt32Uint8V(v map[int32]uint8, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETBincBytes) EncMapInt32Uint64V(v map[int32]uint64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETBincBytes) EncMapInt32IntV(v map[int32]int, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETBincBytes) EncMapInt32Int32V(v map[int32]int32, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETBincBytes) EncMapInt32Float64V(v map[int32]float64, e *encoderBincBytes) {
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
func (e *encoderBincBytes) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincBytes{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETBincBytes) EncMapInt32BoolV(v map[int32]bool, e *encoderBincBytes) {
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

func (helperDecDriverBincBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderBincBytes) bool {
	var ft fastpathDTBincBytes
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
		_ = v
		return false
	}
	return true
}

func (d *decoderBincBytes) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceIntfY(v []interface{}, d *decoderBincBytes) (v2 []interface{}, changed bool) {
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
func (fastpathDTBincBytes) DecSliceIntfN(v []interface{}, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceStringY(v []string, d *decoderBincBytes) (v2 []string, changed bool) {
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
func (fastpathDTBincBytes) DecSliceStringN(v []string, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceBytesY(v [][]byte, d *decoderBincBytes) (v2 [][]byte, changed bool) {
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
func (fastpathDTBincBytes) DecSliceBytesN(v [][]byte, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceFloat32Y(v []float32, d *decoderBincBytes) (v2 []float32, changed bool) {
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
func (fastpathDTBincBytes) DecSliceFloat32N(v []float32, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceFloat64Y(v []float64, d *decoderBincBytes) (v2 []float64, changed bool) {
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
func (fastpathDTBincBytes) DecSliceFloat64N(v []float64, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceUint8Y(v []uint8, d *decoderBincBytes) (v2 []uint8, changed bool) {
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
func (fastpathDTBincBytes) DecSliceUint8N(v []uint8, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceUint64Y(v []uint64, d *decoderBincBytes) (v2 []uint64, changed bool) {
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
func (fastpathDTBincBytes) DecSliceUint64N(v []uint64, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceIntY(v []int, d *decoderBincBytes) (v2 []int, changed bool) {
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
func (fastpathDTBincBytes) DecSliceIntN(v []int, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceInt32Y(v []int32, d *decoderBincBytes) (v2 []int32, changed bool) {
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
func (fastpathDTBincBytes) DecSliceInt32N(v []int32, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceInt64Y(v []int64, d *decoderBincBytes) (v2 []int64, changed bool) {
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
func (fastpathDTBincBytes) DecSliceInt64N(v []int64, d *decoderBincBytes) {
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

func (d *decoderBincBytes) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecSliceBoolY(v []bool, d *decoderBincBytes) (v2 []bool, changed bool) {
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
func (fastpathDTBincBytes) DecSliceBoolN(v []bool, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringStringL(v map[string]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringIntL(v map[string]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntStringL(v map[int]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntIntL(v map[int]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderBincBytes) {
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
func (d *decoderBincBytes) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincBytes
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
func (fastpathDTBincBytes) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderBincBytes) {
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

type fastpathEBincIO struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderBincIO, *encFnInfo, reflect.Value)
}
type fastpathDBincIO struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderBincIO, *decFnInfo, reflect.Value)
}
type fastpathEsBincIO [56]fastpathEBincIO
type fastpathDsBincIO [56]fastpathDBincIO
type fastpathETBincIO struct{}
type fastpathDTBincIO struct{}

func (helperEncDriverBincIO) fastpathEList() *fastpathEsBincIO {
	var i uint = 0
	var s fastpathEsBincIO
	fn := func(v interface{}, fe func(*encoderBincIO, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathEBincIO{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderBincIO).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderBincIO).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderBincIO).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderBincIO).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderBincIO).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderBincIO).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderBincIO).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderBincIO).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderBincIO).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderBincIO).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderBincIO).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderBincIO).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderBincIO).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderBincIO).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderBincIO).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderBincIO).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderBincIO).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderBincIO).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderBincIO).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderBincIO).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderBincIO).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderBincIO).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderBincIO).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderBincIO).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderBincIO).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderBincIO).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderBincIO).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderBincIO).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderBincIO).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderBincIO).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderBincIO).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderBincIO).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderBincIO).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderBincIO).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderBincIO).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderBincIO).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderBincIO).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderBincIO).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderBincIO).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderBincIO).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderBincIO).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderBincIO).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderBincIO).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderBincIO).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderBincIO).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderBincIO).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderBincIO).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderBincIO).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderBincIO).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderBincIO).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderBincIO).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderBincIO).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderBincIO).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderBincIO).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderBincIO).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderBincIO).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverBincIO) fastpathDList() *fastpathDsBincIO {
	var i uint = 0
	var s fastpathDsBincIO
	fn := func(v interface{}, fd func(*decoderBincIO, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDBincIO{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderBincIO).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderBincIO).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderBincIO).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderBincIO).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderBincIO).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderBincIO).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderBincIO).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderBincIO).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderBincIO).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderBincIO).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderBincIO).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderBincIO).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderBincIO).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderBincIO).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderBincIO).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderBincIO).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderBincIO).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderBincIO).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderBincIO).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderBincIO).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderBincIO).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderBincIO).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderBincIO).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderBincIO).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderBincIO).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderBincIO).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderBincIO).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderBincIO).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderBincIO).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderBincIO).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderBincIO).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderBincIO).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderBincIO).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderBincIO).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderBincIO).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderBincIO).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderBincIO).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderBincIO).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderBincIO).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderBincIO).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderBincIO).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderBincIO).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderBincIO).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderBincIO).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderBincIO).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderBincIO).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderBincIO).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderBincIO).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderBincIO).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderBincIO).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderBincIO).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderBincIO).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderBincIO).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderBincIO).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderBincIO).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderBincIO).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverBincIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderBincIO) bool {
	var ft fastpathETBincIO
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
		_ = v
		return false
	}
	return true
}

func (e *encoderBincIO) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceIntfV(v []interface{}, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceIntfV(v []interface{}, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceStringV(v []string, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceStringV(v []string, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceBytesV(v [][]byte, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceBytesV(v [][]byte, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceFloat32V(v []float32, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceFloat32V(v []float32, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceFloat64V(v []float64, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceFloat64V(v []float64, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceUint8V(v []uint8, e *encoderBincIO) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETBincIO) EncAsMapSliceUint8V(v []uint8, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceUint64V(v []uint64, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceUint64V(v []uint64, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceIntV(v []int, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceIntV(v []int, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceInt32V(v []int32, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceInt32V(v []int32, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceInt64V(v []int64, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceInt64V(v []int64, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETBincIO
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
func (fastpathETBincIO) EncSliceBoolV(v []bool, e *encoderBincIO) {
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
func (fastpathETBincIO) EncAsMapSliceBoolV(v []bool, e *encoderBincIO) {
	if len(v) == 0 {
		e.c = 0
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(len(v))
	e.mapStart(len(v) >> 1)
	for j := range v {
		if j&1 == 0 {
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

func (e *encoderBincIO) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETBincIO) EncMapStringIntfV(v map[string]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETBincIO) EncMapStringStringV(v map[string]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETBincIO) EncMapStringBytesV(v map[string][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETBincIO) EncMapStringUint8V(v map[string]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETBincIO) EncMapStringUint64V(v map[string]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETBincIO) EncMapStringIntV(v map[string]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETBincIO) EncMapStringInt32V(v map[string]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETBincIO) EncMapStringFloat64V(v map[string]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETBincIO) EncMapStringBoolV(v map[string]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETBincIO) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETBincIO) EncMapUint8StringV(v map[uint8]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETBincIO) EncMapUint8BytesV(v map[uint8][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETBincIO) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETBincIO) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETBincIO) EncMapUint8IntV(v map[uint8]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETBincIO) EncMapUint8Int32V(v map[uint8]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETBincIO) EncMapUint8Float64V(v map[uint8]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETBincIO) EncMapUint8BoolV(v map[uint8]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETBincIO) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETBincIO) EncMapUint64StringV(v map[uint64]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETBincIO) EncMapUint64BytesV(v map[uint64][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETBincIO) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETBincIO) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETBincIO) EncMapUint64IntV(v map[uint64]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETBincIO) EncMapUint64Int32V(v map[uint64]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETBincIO) EncMapUint64Float64V(v map[uint64]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETBincIO) EncMapUint64BoolV(v map[uint64]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETBincIO) EncMapIntIntfV(v map[int]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETBincIO) EncMapIntStringV(v map[int]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETBincIO) EncMapIntBytesV(v map[int][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETBincIO) EncMapIntUint8V(v map[int]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETBincIO) EncMapIntUint64V(v map[int]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETBincIO) EncMapIntIntV(v map[int]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETBincIO) EncMapIntInt32V(v map[int]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETBincIO) EncMapIntFloat64V(v map[int]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETBincIO) EncMapIntBoolV(v map[int]bool, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETBincIO) EncMapInt32IntfV(v map[int32]interface{}, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETBincIO) EncMapInt32StringV(v map[int32]string, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETBincIO) EncMapInt32BytesV(v map[int32][]byte, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETBincIO) EncMapInt32Uint8V(v map[int32]uint8, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETBincIO) EncMapInt32Uint64V(v map[int32]uint64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETBincIO) EncMapInt32IntV(v map[int32]int, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETBincIO) EncMapInt32Int32V(v map[int32]int32, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETBincIO) EncMapInt32Float64V(v map[int32]float64, e *encoderBincIO) {
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
func (e *encoderBincIO) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETBincIO{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETBincIO) EncMapInt32BoolV(v map[int32]bool, e *encoderBincIO) {
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

func (helperDecDriverBincIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderBincIO) bool {
	var ft fastpathDTBincIO
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
		_ = v
		return false
	}
	return true
}

func (d *decoderBincIO) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceIntfY(v []interface{}, d *decoderBincIO) (v2 []interface{}, changed bool) {
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
func (fastpathDTBincIO) DecSliceIntfN(v []interface{}, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceStringY(v []string, d *decoderBincIO) (v2 []string, changed bool) {
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
func (fastpathDTBincIO) DecSliceStringN(v []string, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceBytesY(v [][]byte, d *decoderBincIO) (v2 [][]byte, changed bool) {
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
func (fastpathDTBincIO) DecSliceBytesN(v [][]byte, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceFloat32Y(v []float32, d *decoderBincIO) (v2 []float32, changed bool) {
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
func (fastpathDTBincIO) DecSliceFloat32N(v []float32, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceFloat64Y(v []float64, d *decoderBincIO) (v2 []float64, changed bool) {
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
func (fastpathDTBincIO) DecSliceFloat64N(v []float64, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceUint8Y(v []uint8, d *decoderBincIO) (v2 []uint8, changed bool) {
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
func (fastpathDTBincIO) DecSliceUint8N(v []uint8, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceUint64Y(v []uint64, d *decoderBincIO) (v2 []uint64, changed bool) {
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
func (fastpathDTBincIO) DecSliceUint64N(v []uint64, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceIntY(v []int, d *decoderBincIO) (v2 []int, changed bool) {
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
func (fastpathDTBincIO) DecSliceIntN(v []int, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceInt32Y(v []int32, d *decoderBincIO) (v2 []int32, changed bool) {
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
func (fastpathDTBincIO) DecSliceInt32N(v []int32, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceInt64Y(v []int64, d *decoderBincIO) (v2 []int64, changed bool) {
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
func (fastpathDTBincIO) DecSliceInt64N(v []int64, d *decoderBincIO) {
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

func (d *decoderBincIO) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecSliceBoolY(v []bool, d *decoderBincIO) (v2 []bool, changed bool) {
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
func (fastpathDTBincIO) DecSliceBoolN(v []bool, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringStringL(v map[string]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringIntL(v map[string]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntStringL(v map[int]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntIntL(v map[int]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderBincIO) {
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
func (d *decoderBincIO) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTBincIO
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
func (fastpathDTBincIO) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderBincIO) {
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
