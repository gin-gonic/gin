//go:build !notmono && !codec.notmono  && !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"slices"
	"sort"
)

type fastpathECborBytes struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderCborBytes, *encFnInfo, reflect.Value)
}
type fastpathDCborBytes struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderCborBytes, *decFnInfo, reflect.Value)
}
type fastpathEsCborBytes [56]fastpathECborBytes
type fastpathDsCborBytes [56]fastpathDCborBytes
type fastpathETCborBytes struct{}
type fastpathDTCborBytes struct{}

func (helperEncDriverCborBytes) fastpathEList() *fastpathEsCborBytes {
	var i uint = 0
	var s fastpathEsCborBytes
	fn := func(v interface{}, fe func(*encoderCborBytes, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathECborBytes{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderCborBytes).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderCborBytes).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderCborBytes).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderCborBytes).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderCborBytes).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderCborBytes).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderCborBytes).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderCborBytes).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderCborBytes).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderCborBytes).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderCborBytes).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderCborBytes).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderCborBytes).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderCborBytes).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderCborBytes).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderCborBytes).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderCborBytes).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderCborBytes).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderCborBytes).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderCborBytes).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderCborBytes).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderCborBytes).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderCborBytes).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderCborBytes).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderCborBytes).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderCborBytes).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderCborBytes).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderCborBytes).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderCborBytes).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderCborBytes).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderCborBytes).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderCborBytes).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderCborBytes).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderCborBytes).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderCborBytes).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderCborBytes).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderCborBytes).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderCborBytes).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderCborBytes).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderCborBytes).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderCborBytes).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderCborBytes).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderCborBytes).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderCborBytes).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderCborBytes).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderCborBytes).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderCborBytes).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderCborBytes).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderCborBytes).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderCborBytes).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderCborBytes).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderCborBytes).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderCborBytes).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderCborBytes).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderCborBytes).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderCborBytes).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverCborBytes) fastpathDList() *fastpathDsCborBytes {
	var i uint = 0
	var s fastpathDsCborBytes
	fn := func(v interface{}, fd func(*decoderCborBytes, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDCborBytes{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderCborBytes).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderCborBytes).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderCborBytes).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderCborBytes).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderCborBytes).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderCborBytes).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderCborBytes).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderCborBytes).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderCborBytes).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderCborBytes).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderCborBytes).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderCborBytes).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderCborBytes).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderCborBytes).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderCborBytes).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderCborBytes).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderCborBytes).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderCborBytes).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderCborBytes).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderCborBytes).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderCborBytes).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderCborBytes).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderCborBytes).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderCborBytes).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderCborBytes).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderCborBytes).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderCborBytes).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderCborBytes).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderCborBytes).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderCborBytes).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderCborBytes).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderCborBytes).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderCborBytes).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderCborBytes).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderCborBytes).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderCborBytes).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderCborBytes).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderCborBytes).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderCborBytes).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderCborBytes).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderCborBytes).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderCborBytes).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderCborBytes).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderCborBytes).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderCborBytes).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderCborBytes).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderCborBytes).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderCborBytes).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderCborBytes).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderCborBytes).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderCborBytes).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderCborBytes).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderCborBytes).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderCborBytes).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderCborBytes).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderCborBytes).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverCborBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderCborBytes) bool {
	var ft fastpathETCborBytes
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

func (e *encoderCborBytes) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceIntfV(v []interface{}, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceIntfV(v []interface{}, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceStringV(v []string, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceStringV(v []string, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceBytesV(v [][]byte, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceBytesV(v [][]byte, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceFloat32V(v []float32, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceFloat32V(v []float32, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceFloat64V(v []float64, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceFloat64V(v []float64, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceUint8V(v []uint8, e *encoderCborBytes) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETCborBytes) EncAsMapSliceUint8V(v []uint8, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceUint64V(v []uint64, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceUint64V(v []uint64, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceIntV(v []int, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceIntV(v []int, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceInt32V(v []int32, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceInt32V(v []int32, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceInt64V(v []int64, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceInt64V(v []int64, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborBytes
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
func (fastpathETCborBytes) EncSliceBoolV(v []bool, e *encoderCborBytes) {
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
func (fastpathETCborBytes) EncAsMapSliceBoolV(v []bool, e *encoderCborBytes) {
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

func (e *encoderCborBytes) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETCborBytes) EncMapStringIntfV(v map[string]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETCborBytes) EncMapStringStringV(v map[string]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETCborBytes) EncMapStringBytesV(v map[string][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETCborBytes) EncMapStringUint8V(v map[string]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETCborBytes) EncMapStringUint64V(v map[string]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETCborBytes) EncMapStringIntV(v map[string]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETCborBytes) EncMapStringInt32V(v map[string]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETCborBytes) EncMapStringFloat64V(v map[string]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETCborBytes) EncMapStringBoolV(v map[string]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETCborBytes) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETCborBytes) EncMapUint8StringV(v map[uint8]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETCborBytes) EncMapUint8BytesV(v map[uint8][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETCborBytes) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETCborBytes) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETCborBytes) EncMapUint8IntV(v map[uint8]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETCborBytes) EncMapUint8Int32V(v map[uint8]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETCborBytes) EncMapUint8Float64V(v map[uint8]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETCborBytes) EncMapUint8BoolV(v map[uint8]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETCborBytes) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETCborBytes) EncMapUint64StringV(v map[uint64]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETCborBytes) EncMapUint64BytesV(v map[uint64][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETCborBytes) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETCborBytes) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETCborBytes) EncMapUint64IntV(v map[uint64]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETCborBytes) EncMapUint64Int32V(v map[uint64]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETCborBytes) EncMapUint64Float64V(v map[uint64]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETCborBytes) EncMapUint64BoolV(v map[uint64]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETCborBytes) EncMapIntIntfV(v map[int]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETCborBytes) EncMapIntStringV(v map[int]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETCborBytes) EncMapIntBytesV(v map[int][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETCborBytes) EncMapIntUint8V(v map[int]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETCborBytes) EncMapIntUint64V(v map[int]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETCborBytes) EncMapIntIntV(v map[int]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETCborBytes) EncMapIntInt32V(v map[int]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETCborBytes) EncMapIntFloat64V(v map[int]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETCborBytes) EncMapIntBoolV(v map[int]bool, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETCborBytes) EncMapInt32IntfV(v map[int32]interface{}, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETCborBytes) EncMapInt32StringV(v map[int32]string, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETCborBytes) EncMapInt32BytesV(v map[int32][]byte, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETCborBytes) EncMapInt32Uint8V(v map[int32]uint8, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETCborBytes) EncMapInt32Uint64V(v map[int32]uint64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETCborBytes) EncMapInt32IntV(v map[int32]int, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETCborBytes) EncMapInt32Int32V(v map[int32]int32, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETCborBytes) EncMapInt32Float64V(v map[int32]float64, e *encoderCborBytes) {
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
func (e *encoderCborBytes) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborBytes{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETCborBytes) EncMapInt32BoolV(v map[int32]bool, e *encoderCborBytes) {
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

func (helperDecDriverCborBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderCborBytes) bool {
	var ft fastpathDTCborBytes
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

func (d *decoderCborBytes) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceIntfY(v []interface{}, d *decoderCborBytes) (v2 []interface{}, changed bool) {
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
func (fastpathDTCborBytes) DecSliceIntfN(v []interface{}, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceStringY(v []string, d *decoderCborBytes) (v2 []string, changed bool) {
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
func (fastpathDTCborBytes) DecSliceStringN(v []string, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceBytesY(v [][]byte, d *decoderCborBytes) (v2 [][]byte, changed bool) {
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
func (fastpathDTCborBytes) DecSliceBytesN(v [][]byte, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceFloat32Y(v []float32, d *decoderCborBytes) (v2 []float32, changed bool) {
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
func (fastpathDTCborBytes) DecSliceFloat32N(v []float32, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceFloat64Y(v []float64, d *decoderCborBytes) (v2 []float64, changed bool) {
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
func (fastpathDTCborBytes) DecSliceFloat64N(v []float64, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceUint8Y(v []uint8, d *decoderCborBytes) (v2 []uint8, changed bool) {
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
func (fastpathDTCborBytes) DecSliceUint8N(v []uint8, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceUint64Y(v []uint64, d *decoderCborBytes) (v2 []uint64, changed bool) {
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
func (fastpathDTCborBytes) DecSliceUint64N(v []uint64, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceIntY(v []int, d *decoderCborBytes) (v2 []int, changed bool) {
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
func (fastpathDTCborBytes) DecSliceIntN(v []int, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceInt32Y(v []int32, d *decoderCborBytes) (v2 []int32, changed bool) {
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
func (fastpathDTCborBytes) DecSliceInt32N(v []int32, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceInt64Y(v []int64, d *decoderCborBytes) (v2 []int64, changed bool) {
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
func (fastpathDTCborBytes) DecSliceInt64N(v []int64, d *decoderCborBytes) {
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

func (d *decoderCborBytes) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecSliceBoolY(v []bool, d *decoderCborBytes) (v2 []bool, changed bool) {
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
func (fastpathDTCborBytes) DecSliceBoolN(v []bool, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringStringL(v map[string]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringIntL(v map[string]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntStringL(v map[int]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntIntL(v map[int]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderCborBytes) {
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
func (d *decoderCborBytes) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborBytes
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
func (fastpathDTCborBytes) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderCborBytes) {
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

type fastpathECborIO struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderCborIO, *encFnInfo, reflect.Value)
}
type fastpathDCborIO struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderCborIO, *decFnInfo, reflect.Value)
}
type fastpathEsCborIO [56]fastpathECborIO
type fastpathDsCborIO [56]fastpathDCborIO
type fastpathETCborIO struct{}
type fastpathDTCborIO struct{}

func (helperEncDriverCborIO) fastpathEList() *fastpathEsCborIO {
	var i uint = 0
	var s fastpathEsCborIO
	fn := func(v interface{}, fe func(*encoderCborIO, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathECborIO{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderCborIO).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderCborIO).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderCborIO).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderCborIO).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderCborIO).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderCborIO).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderCborIO).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderCborIO).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderCborIO).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderCborIO).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderCborIO).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderCborIO).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderCborIO).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderCborIO).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderCborIO).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderCborIO).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderCborIO).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderCborIO).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderCborIO).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderCborIO).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderCborIO).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderCborIO).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderCborIO).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderCborIO).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderCborIO).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderCborIO).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderCborIO).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderCborIO).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderCborIO).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderCborIO).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderCborIO).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderCborIO).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderCborIO).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderCborIO).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderCborIO).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderCborIO).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderCborIO).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderCborIO).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderCborIO).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderCborIO).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderCborIO).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderCborIO).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderCborIO).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderCborIO).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderCborIO).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderCborIO).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderCborIO).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderCborIO).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderCborIO).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderCborIO).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderCborIO).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderCborIO).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderCborIO).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderCborIO).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderCborIO).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderCborIO).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverCborIO) fastpathDList() *fastpathDsCborIO {
	var i uint = 0
	var s fastpathDsCborIO
	fn := func(v interface{}, fd func(*decoderCborIO, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDCborIO{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderCborIO).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderCborIO).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderCborIO).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderCborIO).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderCborIO).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderCborIO).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderCborIO).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderCborIO).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderCborIO).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderCborIO).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderCborIO).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderCborIO).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderCborIO).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderCborIO).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderCborIO).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderCborIO).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderCborIO).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderCborIO).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderCborIO).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderCborIO).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderCborIO).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderCborIO).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderCborIO).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderCborIO).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderCborIO).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderCborIO).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderCborIO).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderCborIO).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderCborIO).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderCborIO).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderCborIO).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderCborIO).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderCborIO).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderCborIO).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderCborIO).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderCborIO).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderCborIO).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderCborIO).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderCborIO).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderCborIO).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderCborIO).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderCborIO).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderCborIO).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderCborIO).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderCborIO).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderCborIO).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderCborIO).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderCborIO).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderCborIO).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderCborIO).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderCborIO).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderCborIO).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderCborIO).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderCborIO).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderCborIO).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderCborIO).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverCborIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderCborIO) bool {
	var ft fastpathETCborIO
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

func (e *encoderCborIO) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceIntfV(v []interface{}, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceIntfV(v []interface{}, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceStringV(v []string, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceStringV(v []string, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceBytesV(v [][]byte, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceBytesV(v [][]byte, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceFloat32V(v []float32, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceFloat32V(v []float32, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceFloat64V(v []float64, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceFloat64V(v []float64, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceUint8V(v []uint8, e *encoderCborIO) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETCborIO) EncAsMapSliceUint8V(v []uint8, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceUint64V(v []uint64, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceUint64V(v []uint64, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceIntV(v []int, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceIntV(v []int, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceInt32V(v []int32, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceInt32V(v []int32, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceInt64V(v []int64, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceInt64V(v []int64, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETCborIO
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
func (fastpathETCborIO) EncSliceBoolV(v []bool, e *encoderCborIO) {
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
func (fastpathETCborIO) EncAsMapSliceBoolV(v []bool, e *encoderCborIO) {
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

func (e *encoderCborIO) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETCborIO) EncMapStringIntfV(v map[string]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETCborIO) EncMapStringStringV(v map[string]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETCborIO) EncMapStringBytesV(v map[string][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETCborIO) EncMapStringUint8V(v map[string]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETCborIO) EncMapStringUint64V(v map[string]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETCborIO) EncMapStringIntV(v map[string]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETCborIO) EncMapStringInt32V(v map[string]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETCborIO) EncMapStringFloat64V(v map[string]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETCborIO) EncMapStringBoolV(v map[string]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETCborIO) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETCborIO) EncMapUint8StringV(v map[uint8]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETCborIO) EncMapUint8BytesV(v map[uint8][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETCborIO) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETCborIO) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETCborIO) EncMapUint8IntV(v map[uint8]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETCborIO) EncMapUint8Int32V(v map[uint8]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETCborIO) EncMapUint8Float64V(v map[uint8]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETCborIO) EncMapUint8BoolV(v map[uint8]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETCborIO) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETCborIO) EncMapUint64StringV(v map[uint64]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETCborIO) EncMapUint64BytesV(v map[uint64][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETCborIO) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETCborIO) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETCborIO) EncMapUint64IntV(v map[uint64]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETCborIO) EncMapUint64Int32V(v map[uint64]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETCborIO) EncMapUint64Float64V(v map[uint64]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETCborIO) EncMapUint64BoolV(v map[uint64]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETCborIO) EncMapIntIntfV(v map[int]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETCborIO) EncMapIntStringV(v map[int]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETCborIO) EncMapIntBytesV(v map[int][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETCborIO) EncMapIntUint8V(v map[int]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETCborIO) EncMapIntUint64V(v map[int]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETCborIO) EncMapIntIntV(v map[int]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETCborIO) EncMapIntInt32V(v map[int]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETCborIO) EncMapIntFloat64V(v map[int]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETCborIO) EncMapIntBoolV(v map[int]bool, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETCborIO) EncMapInt32IntfV(v map[int32]interface{}, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETCborIO) EncMapInt32StringV(v map[int32]string, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETCborIO) EncMapInt32BytesV(v map[int32][]byte, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETCborIO) EncMapInt32Uint8V(v map[int32]uint8, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETCborIO) EncMapInt32Uint64V(v map[int32]uint64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETCborIO) EncMapInt32IntV(v map[int32]int, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETCborIO) EncMapInt32Int32V(v map[int32]int32, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETCborIO) EncMapInt32Float64V(v map[int32]float64, e *encoderCborIO) {
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
func (e *encoderCborIO) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETCborIO{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETCborIO) EncMapInt32BoolV(v map[int32]bool, e *encoderCborIO) {
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

func (helperDecDriverCborIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderCborIO) bool {
	var ft fastpathDTCborIO
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

func (d *decoderCborIO) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceIntfY(v []interface{}, d *decoderCborIO) (v2 []interface{}, changed bool) {
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
func (fastpathDTCborIO) DecSliceIntfN(v []interface{}, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceStringY(v []string, d *decoderCborIO) (v2 []string, changed bool) {
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
func (fastpathDTCborIO) DecSliceStringN(v []string, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceBytesY(v [][]byte, d *decoderCborIO) (v2 [][]byte, changed bool) {
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
func (fastpathDTCborIO) DecSliceBytesN(v [][]byte, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceFloat32Y(v []float32, d *decoderCborIO) (v2 []float32, changed bool) {
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
func (fastpathDTCborIO) DecSliceFloat32N(v []float32, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceFloat64Y(v []float64, d *decoderCborIO) (v2 []float64, changed bool) {
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
func (fastpathDTCborIO) DecSliceFloat64N(v []float64, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceUint8Y(v []uint8, d *decoderCborIO) (v2 []uint8, changed bool) {
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
func (fastpathDTCborIO) DecSliceUint8N(v []uint8, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceUint64Y(v []uint64, d *decoderCborIO) (v2 []uint64, changed bool) {
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
func (fastpathDTCborIO) DecSliceUint64N(v []uint64, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceIntY(v []int, d *decoderCborIO) (v2 []int, changed bool) {
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
func (fastpathDTCborIO) DecSliceIntN(v []int, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceInt32Y(v []int32, d *decoderCborIO) (v2 []int32, changed bool) {
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
func (fastpathDTCborIO) DecSliceInt32N(v []int32, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceInt64Y(v []int64, d *decoderCborIO) (v2 []int64, changed bool) {
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
func (fastpathDTCborIO) DecSliceInt64N(v []int64, d *decoderCborIO) {
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

func (d *decoderCborIO) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecSliceBoolY(v []bool, d *decoderCborIO) (v2 []bool, changed bool) {
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
func (fastpathDTCborIO) DecSliceBoolN(v []bool, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringStringL(v map[string]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringIntL(v map[string]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntStringL(v map[int]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntIntL(v map[int]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderCborIO) {
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
func (d *decoderCborIO) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTCborIO
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
func (fastpathDTCborIO) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderCborIO) {
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
