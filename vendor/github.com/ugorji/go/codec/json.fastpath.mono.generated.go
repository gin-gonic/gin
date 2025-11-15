//go:build !notmono && !codec.notmono  && !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"slices"
	"sort"
)

type fastpathEJsonBytes struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderJsonBytes, *encFnInfo, reflect.Value)
}
type fastpathDJsonBytes struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderJsonBytes, *decFnInfo, reflect.Value)
}
type fastpathEsJsonBytes [56]fastpathEJsonBytes
type fastpathDsJsonBytes [56]fastpathDJsonBytes
type fastpathETJsonBytes struct{}
type fastpathDTJsonBytes struct{}

func (helperEncDriverJsonBytes) fastpathEList() *fastpathEsJsonBytes {
	var i uint = 0
	var s fastpathEsJsonBytes
	fn := func(v interface{}, fe func(*encoderJsonBytes, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathEJsonBytes{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderJsonBytes).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderJsonBytes).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderJsonBytes).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderJsonBytes).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderJsonBytes).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderJsonBytes).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderJsonBytes).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderJsonBytes).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderJsonBytes).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderJsonBytes).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderJsonBytes).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderJsonBytes).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderJsonBytes).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderJsonBytes).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderJsonBytes).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderJsonBytes).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderJsonBytes).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderJsonBytes).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderJsonBytes).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderJsonBytes).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderJsonBytes).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderJsonBytes).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderJsonBytes).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderJsonBytes).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderJsonBytes).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderJsonBytes).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderJsonBytes).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderJsonBytes).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderJsonBytes).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderJsonBytes).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderJsonBytes).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderJsonBytes).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderJsonBytes).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderJsonBytes).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderJsonBytes).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderJsonBytes).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderJsonBytes).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderJsonBytes).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderJsonBytes).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderJsonBytes).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderJsonBytes).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderJsonBytes).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderJsonBytes).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderJsonBytes).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderJsonBytes).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderJsonBytes).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderJsonBytes).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderJsonBytes).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderJsonBytes).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderJsonBytes).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderJsonBytes).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderJsonBytes).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderJsonBytes).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderJsonBytes).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderJsonBytes).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderJsonBytes).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverJsonBytes) fastpathDList() *fastpathDsJsonBytes {
	var i uint = 0
	var s fastpathDsJsonBytes
	fn := func(v interface{}, fd func(*decoderJsonBytes, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDJsonBytes{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderJsonBytes).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderJsonBytes).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderJsonBytes).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderJsonBytes).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderJsonBytes).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderJsonBytes).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderJsonBytes).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderJsonBytes).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderJsonBytes).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderJsonBytes).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderJsonBytes).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderJsonBytes).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderJsonBytes).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderJsonBytes).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderJsonBytes).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderJsonBytes).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderJsonBytes).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderJsonBytes).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderJsonBytes).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderJsonBytes).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderJsonBytes).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderJsonBytes).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderJsonBytes).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderJsonBytes).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderJsonBytes).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderJsonBytes).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderJsonBytes).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderJsonBytes).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderJsonBytes).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderJsonBytes).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderJsonBytes).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderJsonBytes).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderJsonBytes).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderJsonBytes).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderJsonBytes).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderJsonBytes).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderJsonBytes).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderJsonBytes).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderJsonBytes).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderJsonBytes).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderJsonBytes).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderJsonBytes).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderJsonBytes).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderJsonBytes).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderJsonBytes).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderJsonBytes).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderJsonBytes).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderJsonBytes).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderJsonBytes).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderJsonBytes).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderJsonBytes).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderJsonBytes).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderJsonBytes).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderJsonBytes).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderJsonBytes).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderJsonBytes).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverJsonBytes) fastpathEncodeTypeSwitch(iv interface{}, e *encoderJsonBytes) bool {
	var ft fastpathETJsonBytes
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

func (e *encoderJsonBytes) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceIntfV(v []interface{}, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceIntfV(v []interface{}, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceStringV(v []string, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceStringV(v []string, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceBytesV(v [][]byte, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceBytesV(v [][]byte, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceFloat32V(v []float32, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceFloat32V(v []float32, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceFloat64V(v []float64, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceFloat64V(v []float64, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceUint8V(v []uint8, e *encoderJsonBytes) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETJsonBytes) EncAsMapSliceUint8V(v []uint8, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceUint64V(v []uint64, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceUint64V(v []uint64, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceIntV(v []int, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceIntV(v []int, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceInt32V(v []int32, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceInt32V(v []int32, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceInt64V(v []int64, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceInt64V(v []int64, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonBytes
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
func (fastpathETJsonBytes) EncSliceBoolV(v []bool, e *encoderJsonBytes) {
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
func (fastpathETJsonBytes) EncAsMapSliceBoolV(v []bool, e *encoderJsonBytes) {
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

func (e *encoderJsonBytes) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETJsonBytes) EncMapStringIntfV(v map[string]interface{}, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETJsonBytes) EncMapStringStringV(v map[string]string, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETJsonBytes) EncMapStringBytesV(v map[string][]byte, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETJsonBytes) EncMapStringUint8V(v map[string]uint8, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETJsonBytes) EncMapStringUint64V(v map[string]uint64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETJsonBytes) EncMapStringIntV(v map[string]int, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETJsonBytes) EncMapStringInt32V(v map[string]int32, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETJsonBytes) EncMapStringFloat64V(v map[string]float64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETJsonBytes) EncMapStringBoolV(v map[string]bool, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETJsonBytes) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETJsonBytes) EncMapUint8StringV(v map[uint8]string, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETJsonBytes) EncMapUint8BytesV(v map[uint8][]byte, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETJsonBytes) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETJsonBytes) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETJsonBytes) EncMapUint8IntV(v map[uint8]int, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETJsonBytes) EncMapUint8Int32V(v map[uint8]int32, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETJsonBytes) EncMapUint8Float64V(v map[uint8]float64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETJsonBytes) EncMapUint8BoolV(v map[uint8]bool, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETJsonBytes) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETJsonBytes) EncMapUint64StringV(v map[uint64]string, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETJsonBytes) EncMapUint64BytesV(v map[uint64][]byte, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETJsonBytes) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETJsonBytes) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETJsonBytes) EncMapUint64IntV(v map[uint64]int, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETJsonBytes) EncMapUint64Int32V(v map[uint64]int32, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETJsonBytes) EncMapUint64Float64V(v map[uint64]float64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETJsonBytes) EncMapUint64BoolV(v map[uint64]bool, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETJsonBytes) EncMapIntIntfV(v map[int]interface{}, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETJsonBytes) EncMapIntStringV(v map[int]string, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETJsonBytes) EncMapIntBytesV(v map[int][]byte, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETJsonBytes) EncMapIntUint8V(v map[int]uint8, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETJsonBytes) EncMapIntUint64V(v map[int]uint64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETJsonBytes) EncMapIntIntV(v map[int]int, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETJsonBytes) EncMapIntInt32V(v map[int]int32, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETJsonBytes) EncMapIntFloat64V(v map[int]float64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETJsonBytes) EncMapIntBoolV(v map[int]bool, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETJsonBytes) EncMapInt32IntfV(v map[int32]interface{}, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETJsonBytes) EncMapInt32StringV(v map[int32]string, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETJsonBytes) EncMapInt32BytesV(v map[int32][]byte, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETJsonBytes) EncMapInt32Uint8V(v map[int32]uint8, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETJsonBytes) EncMapInt32Uint64V(v map[int32]uint64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETJsonBytes) EncMapInt32IntV(v map[int32]int, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETJsonBytes) EncMapInt32Int32V(v map[int32]int32, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETJsonBytes) EncMapInt32Float64V(v map[int32]float64, e *encoderJsonBytes) {
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
func (e *encoderJsonBytes) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonBytes{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETJsonBytes) EncMapInt32BoolV(v map[int32]bool, e *encoderJsonBytes) {
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

func (helperDecDriverJsonBytes) fastpathDecodeTypeSwitch(iv interface{}, d *decoderJsonBytes) bool {
	var ft fastpathDTJsonBytes
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

func (d *decoderJsonBytes) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceIntfY(v []interface{}, d *decoderJsonBytes) (v2 []interface{}, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceIntfN(v []interface{}, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceStringY(v []string, d *decoderJsonBytes) (v2 []string, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceStringN(v []string, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceBytesY(v [][]byte, d *decoderJsonBytes) (v2 [][]byte, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceBytesN(v [][]byte, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceFloat32Y(v []float32, d *decoderJsonBytes) (v2 []float32, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceFloat32N(v []float32, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceFloat64Y(v []float64, d *decoderJsonBytes) (v2 []float64, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceFloat64N(v []float64, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceUint8Y(v []uint8, d *decoderJsonBytes) (v2 []uint8, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceUint8N(v []uint8, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceUint64Y(v []uint64, d *decoderJsonBytes) (v2 []uint64, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceUint64N(v []uint64, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceIntY(v []int, d *decoderJsonBytes) (v2 []int, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceIntN(v []int, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceInt32Y(v []int32, d *decoderJsonBytes) (v2 []int32, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceInt32N(v []int32, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceInt64Y(v []int64, d *decoderJsonBytes) (v2 []int64, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceInt64N(v []int64, d *decoderJsonBytes) {
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

func (d *decoderJsonBytes) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecSliceBoolY(v []bool, d *decoderJsonBytes) (v2 []bool, changed bool) {
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
func (fastpathDTJsonBytes) DecSliceBoolN(v []bool, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringStringL(v map[string]string, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringIntL(v map[string]int, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntStringL(v map[int]string, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntIntL(v map[int]int, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderJsonBytes) {
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
func (d *decoderJsonBytes) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonBytes
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
func (fastpathDTJsonBytes) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderJsonBytes) {
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

type fastpathEJsonIO struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encoderJsonIO, *encFnInfo, reflect.Value)
}
type fastpathDJsonIO struct {
	rtid  uintptr
	rt    reflect.Type
	decfn func(*decoderJsonIO, *decFnInfo, reflect.Value)
}
type fastpathEsJsonIO [56]fastpathEJsonIO
type fastpathDsJsonIO [56]fastpathDJsonIO
type fastpathETJsonIO struct{}
type fastpathDTJsonIO struct{}

func (helperEncDriverJsonIO) fastpathEList() *fastpathEsJsonIO {
	var i uint = 0
	var s fastpathEsJsonIO
	fn := func(v interface{}, fe func(*encoderJsonIO, *encFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathEJsonIO{rt2id(xrt), xrt, fe}
		i++
	}

	fn([]interface{}(nil), (*encoderJsonIO).fastpathEncSliceIntfR)
	fn([]string(nil), (*encoderJsonIO).fastpathEncSliceStringR)
	fn([][]byte(nil), (*encoderJsonIO).fastpathEncSliceBytesR)
	fn([]float32(nil), (*encoderJsonIO).fastpathEncSliceFloat32R)
	fn([]float64(nil), (*encoderJsonIO).fastpathEncSliceFloat64R)
	fn([]uint8(nil), (*encoderJsonIO).fastpathEncSliceUint8R)
	fn([]uint64(nil), (*encoderJsonIO).fastpathEncSliceUint64R)
	fn([]int(nil), (*encoderJsonIO).fastpathEncSliceIntR)
	fn([]int32(nil), (*encoderJsonIO).fastpathEncSliceInt32R)
	fn([]int64(nil), (*encoderJsonIO).fastpathEncSliceInt64R)
	fn([]bool(nil), (*encoderJsonIO).fastpathEncSliceBoolR)

	fn(map[string]interface{}(nil), (*encoderJsonIO).fastpathEncMapStringIntfR)
	fn(map[string]string(nil), (*encoderJsonIO).fastpathEncMapStringStringR)
	fn(map[string][]byte(nil), (*encoderJsonIO).fastpathEncMapStringBytesR)
	fn(map[string]uint8(nil), (*encoderJsonIO).fastpathEncMapStringUint8R)
	fn(map[string]uint64(nil), (*encoderJsonIO).fastpathEncMapStringUint64R)
	fn(map[string]int(nil), (*encoderJsonIO).fastpathEncMapStringIntR)
	fn(map[string]int32(nil), (*encoderJsonIO).fastpathEncMapStringInt32R)
	fn(map[string]float64(nil), (*encoderJsonIO).fastpathEncMapStringFloat64R)
	fn(map[string]bool(nil), (*encoderJsonIO).fastpathEncMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*encoderJsonIO).fastpathEncMapUint8IntfR)
	fn(map[uint8]string(nil), (*encoderJsonIO).fastpathEncMapUint8StringR)
	fn(map[uint8][]byte(nil), (*encoderJsonIO).fastpathEncMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*encoderJsonIO).fastpathEncMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*encoderJsonIO).fastpathEncMapUint8Uint64R)
	fn(map[uint8]int(nil), (*encoderJsonIO).fastpathEncMapUint8IntR)
	fn(map[uint8]int32(nil), (*encoderJsonIO).fastpathEncMapUint8Int32R)
	fn(map[uint8]float64(nil), (*encoderJsonIO).fastpathEncMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encoderJsonIO).fastpathEncMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*encoderJsonIO).fastpathEncMapUint64IntfR)
	fn(map[uint64]string(nil), (*encoderJsonIO).fastpathEncMapUint64StringR)
	fn(map[uint64][]byte(nil), (*encoderJsonIO).fastpathEncMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*encoderJsonIO).fastpathEncMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*encoderJsonIO).fastpathEncMapUint64Uint64R)
	fn(map[uint64]int(nil), (*encoderJsonIO).fastpathEncMapUint64IntR)
	fn(map[uint64]int32(nil), (*encoderJsonIO).fastpathEncMapUint64Int32R)
	fn(map[uint64]float64(nil), (*encoderJsonIO).fastpathEncMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encoderJsonIO).fastpathEncMapUint64BoolR)
	fn(map[int]interface{}(nil), (*encoderJsonIO).fastpathEncMapIntIntfR)
	fn(map[int]string(nil), (*encoderJsonIO).fastpathEncMapIntStringR)
	fn(map[int][]byte(nil), (*encoderJsonIO).fastpathEncMapIntBytesR)
	fn(map[int]uint8(nil), (*encoderJsonIO).fastpathEncMapIntUint8R)
	fn(map[int]uint64(nil), (*encoderJsonIO).fastpathEncMapIntUint64R)
	fn(map[int]int(nil), (*encoderJsonIO).fastpathEncMapIntIntR)
	fn(map[int]int32(nil), (*encoderJsonIO).fastpathEncMapIntInt32R)
	fn(map[int]float64(nil), (*encoderJsonIO).fastpathEncMapIntFloat64R)
	fn(map[int]bool(nil), (*encoderJsonIO).fastpathEncMapIntBoolR)
	fn(map[int32]interface{}(nil), (*encoderJsonIO).fastpathEncMapInt32IntfR)
	fn(map[int32]string(nil), (*encoderJsonIO).fastpathEncMapInt32StringR)
	fn(map[int32][]byte(nil), (*encoderJsonIO).fastpathEncMapInt32BytesR)
	fn(map[int32]uint8(nil), (*encoderJsonIO).fastpathEncMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*encoderJsonIO).fastpathEncMapInt32Uint64R)
	fn(map[int32]int(nil), (*encoderJsonIO).fastpathEncMapInt32IntR)
	fn(map[int32]int32(nil), (*encoderJsonIO).fastpathEncMapInt32Int32R)
	fn(map[int32]float64(nil), (*encoderJsonIO).fastpathEncMapInt32Float64R)
	fn(map[int32]bool(nil), (*encoderJsonIO).fastpathEncMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperDecDriverJsonIO) fastpathDList() *fastpathDsJsonIO {
	var i uint = 0
	var s fastpathDsJsonIO
	fn := func(v interface{}, fd func(*decoderJsonIO, *decFnInfo, reflect.Value)) {
		xrt := reflect.TypeOf(v)
		s[i] = fastpathDJsonIO{rt2id(xrt), xrt, fd}
		i++
	}

	fn([]interface{}(nil), (*decoderJsonIO).fastpathDecSliceIntfR)
	fn([]string(nil), (*decoderJsonIO).fastpathDecSliceStringR)
	fn([][]byte(nil), (*decoderJsonIO).fastpathDecSliceBytesR)
	fn([]float32(nil), (*decoderJsonIO).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*decoderJsonIO).fastpathDecSliceFloat64R)
	fn([]uint8(nil), (*decoderJsonIO).fastpathDecSliceUint8R)
	fn([]uint64(nil), (*decoderJsonIO).fastpathDecSliceUint64R)
	fn([]int(nil), (*decoderJsonIO).fastpathDecSliceIntR)
	fn([]int32(nil), (*decoderJsonIO).fastpathDecSliceInt32R)
	fn([]int64(nil), (*decoderJsonIO).fastpathDecSliceInt64R)
	fn([]bool(nil), (*decoderJsonIO).fastpathDecSliceBoolR)

	fn(map[string]interface{}(nil), (*decoderJsonIO).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*decoderJsonIO).fastpathDecMapStringStringR)
	fn(map[string][]byte(nil), (*decoderJsonIO).fastpathDecMapStringBytesR)
	fn(map[string]uint8(nil), (*decoderJsonIO).fastpathDecMapStringUint8R)
	fn(map[string]uint64(nil), (*decoderJsonIO).fastpathDecMapStringUint64R)
	fn(map[string]int(nil), (*decoderJsonIO).fastpathDecMapStringIntR)
	fn(map[string]int32(nil), (*decoderJsonIO).fastpathDecMapStringInt32R)
	fn(map[string]float64(nil), (*decoderJsonIO).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*decoderJsonIO).fastpathDecMapStringBoolR)
	fn(map[uint8]interface{}(nil), (*decoderJsonIO).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*decoderJsonIO).fastpathDecMapUint8StringR)
	fn(map[uint8][]byte(nil), (*decoderJsonIO).fastpathDecMapUint8BytesR)
	fn(map[uint8]uint8(nil), (*decoderJsonIO).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint64(nil), (*decoderJsonIO).fastpathDecMapUint8Uint64R)
	fn(map[uint8]int(nil), (*decoderJsonIO).fastpathDecMapUint8IntR)
	fn(map[uint8]int32(nil), (*decoderJsonIO).fastpathDecMapUint8Int32R)
	fn(map[uint8]float64(nil), (*decoderJsonIO).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*decoderJsonIO).fastpathDecMapUint8BoolR)
	fn(map[uint64]interface{}(nil), (*decoderJsonIO).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*decoderJsonIO).fastpathDecMapUint64StringR)
	fn(map[uint64][]byte(nil), (*decoderJsonIO).fastpathDecMapUint64BytesR)
	fn(map[uint64]uint8(nil), (*decoderJsonIO).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint64(nil), (*decoderJsonIO).fastpathDecMapUint64Uint64R)
	fn(map[uint64]int(nil), (*decoderJsonIO).fastpathDecMapUint64IntR)
	fn(map[uint64]int32(nil), (*decoderJsonIO).fastpathDecMapUint64Int32R)
	fn(map[uint64]float64(nil), (*decoderJsonIO).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*decoderJsonIO).fastpathDecMapUint64BoolR)
	fn(map[int]interface{}(nil), (*decoderJsonIO).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*decoderJsonIO).fastpathDecMapIntStringR)
	fn(map[int][]byte(nil), (*decoderJsonIO).fastpathDecMapIntBytesR)
	fn(map[int]uint8(nil), (*decoderJsonIO).fastpathDecMapIntUint8R)
	fn(map[int]uint64(nil), (*decoderJsonIO).fastpathDecMapIntUint64R)
	fn(map[int]int(nil), (*decoderJsonIO).fastpathDecMapIntIntR)
	fn(map[int]int32(nil), (*decoderJsonIO).fastpathDecMapIntInt32R)
	fn(map[int]float64(nil), (*decoderJsonIO).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*decoderJsonIO).fastpathDecMapIntBoolR)
	fn(map[int32]interface{}(nil), (*decoderJsonIO).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*decoderJsonIO).fastpathDecMapInt32StringR)
	fn(map[int32][]byte(nil), (*decoderJsonIO).fastpathDecMapInt32BytesR)
	fn(map[int32]uint8(nil), (*decoderJsonIO).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint64(nil), (*decoderJsonIO).fastpathDecMapInt32Uint64R)
	fn(map[int32]int(nil), (*decoderJsonIO).fastpathDecMapInt32IntR)
	fn(map[int32]int32(nil), (*decoderJsonIO).fastpathDecMapInt32Int32R)
	fn(map[int32]float64(nil), (*decoderJsonIO).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*decoderJsonIO).fastpathDecMapInt32BoolR)

	sort.Slice(s[:], func(i, j int) bool { return s[i].rtid < s[j].rtid })
	return &s
}

func (helperEncDriverJsonIO) fastpathEncodeTypeSwitch(iv interface{}, e *encoderJsonIO) bool {
	var ft fastpathETJsonIO
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

func (e *encoderJsonIO) fastpathEncSliceIntfR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceIntfV(v []interface{}, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceIntfV(v []interface{}, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceStringR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceStringV(v []string, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceStringV(v []string, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceBytesR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceBytesV(v [][]byte, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceBytesV(v [][]byte, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceFloat32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceFloat32V(v []float32, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceFloat32V(v []float32, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceFloat64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceFloat64V(v []float64, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceFloat64V(v []float64, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceUint8R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceUint8V(v []uint8, e *encoderJsonIO) {
	e.e.EncodeStringBytesRaw(v)
}
func (fastpathETJsonIO) EncAsMapSliceUint8V(v []uint8, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceUint64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceUint64V(v []uint64, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceUint64V(v []uint64, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceIntR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceIntV(v []int, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceIntV(v []int, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceInt32R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceInt32V(v []int32, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceInt32V(v []int32, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceInt64R(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceInt64V(v []int64, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceInt64V(v []int64, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncSliceBoolR(f *encFnInfo, rv reflect.Value) {
	var ft fastpathETJsonIO
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
func (fastpathETJsonIO) EncSliceBoolV(v []bool, e *encoderJsonIO) {
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
func (fastpathETJsonIO) EncAsMapSliceBoolV(v []bool, e *encoderJsonIO) {
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

func (e *encoderJsonIO) fastpathEncMapStringIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringIntfV(rv2i(rv).(map[string]interface{}), e)
}
func (fastpathETJsonIO) EncMapStringIntfV(v map[string]interface{}, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringStringV(rv2i(rv).(map[string]string), e)
}
func (fastpathETJsonIO) EncMapStringStringV(v map[string]string, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringBytesV(rv2i(rv).(map[string][]byte), e)
}
func (fastpathETJsonIO) EncMapStringBytesV(v map[string][]byte, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringUint8V(rv2i(rv).(map[string]uint8), e)
}
func (fastpathETJsonIO) EncMapStringUint8V(v map[string]uint8, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringUint64V(rv2i(rv).(map[string]uint64), e)
}
func (fastpathETJsonIO) EncMapStringUint64V(v map[string]uint64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringIntV(rv2i(rv).(map[string]int), e)
}
func (fastpathETJsonIO) EncMapStringIntV(v map[string]int, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringInt32V(rv2i(rv).(map[string]int32), e)
}
func (fastpathETJsonIO) EncMapStringInt32V(v map[string]int32, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringFloat64V(rv2i(rv).(map[string]float64), e)
}
func (fastpathETJsonIO) EncMapStringFloat64V(v map[string]float64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapStringBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapStringBoolV(rv2i(rv).(map[string]bool), e)
}
func (fastpathETJsonIO) EncMapStringBoolV(v map[string]bool, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8IntfV(rv2i(rv).(map[uint8]interface{}), e)
}
func (fastpathETJsonIO) EncMapUint8IntfV(v map[uint8]interface{}, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8StringV(rv2i(rv).(map[uint8]string), e)
}
func (fastpathETJsonIO) EncMapUint8StringV(v map[uint8]string, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8BytesV(rv2i(rv).(map[uint8][]byte), e)
}
func (fastpathETJsonIO) EncMapUint8BytesV(v map[uint8][]byte, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8Uint8V(rv2i(rv).(map[uint8]uint8), e)
}
func (fastpathETJsonIO) EncMapUint8Uint8V(v map[uint8]uint8, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8Uint64V(rv2i(rv).(map[uint8]uint64), e)
}
func (fastpathETJsonIO) EncMapUint8Uint64V(v map[uint8]uint64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8IntV(rv2i(rv).(map[uint8]int), e)
}
func (fastpathETJsonIO) EncMapUint8IntV(v map[uint8]int, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8Int32V(rv2i(rv).(map[uint8]int32), e)
}
func (fastpathETJsonIO) EncMapUint8Int32V(v map[uint8]int32, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8Float64V(rv2i(rv).(map[uint8]float64), e)
}
func (fastpathETJsonIO) EncMapUint8Float64V(v map[uint8]float64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint8BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint8BoolV(rv2i(rv).(map[uint8]bool), e)
}
func (fastpathETJsonIO) EncMapUint8BoolV(v map[uint8]bool, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64IntfV(rv2i(rv).(map[uint64]interface{}), e)
}
func (fastpathETJsonIO) EncMapUint64IntfV(v map[uint64]interface{}, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64StringV(rv2i(rv).(map[uint64]string), e)
}
func (fastpathETJsonIO) EncMapUint64StringV(v map[uint64]string, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64BytesV(rv2i(rv).(map[uint64][]byte), e)
}
func (fastpathETJsonIO) EncMapUint64BytesV(v map[uint64][]byte, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64Uint8V(rv2i(rv).(map[uint64]uint8), e)
}
func (fastpathETJsonIO) EncMapUint64Uint8V(v map[uint64]uint8, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64Uint64V(rv2i(rv).(map[uint64]uint64), e)
}
func (fastpathETJsonIO) EncMapUint64Uint64V(v map[uint64]uint64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64IntV(rv2i(rv).(map[uint64]int), e)
}
func (fastpathETJsonIO) EncMapUint64IntV(v map[uint64]int, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64Int32V(rv2i(rv).(map[uint64]int32), e)
}
func (fastpathETJsonIO) EncMapUint64Int32V(v map[uint64]int32, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64Float64V(rv2i(rv).(map[uint64]float64), e)
}
func (fastpathETJsonIO) EncMapUint64Float64V(v map[uint64]float64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapUint64BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapUint64BoolV(rv2i(rv).(map[uint64]bool), e)
}
func (fastpathETJsonIO) EncMapUint64BoolV(v map[uint64]bool, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntIntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntIntfV(rv2i(rv).(map[int]interface{}), e)
}
func (fastpathETJsonIO) EncMapIntIntfV(v map[int]interface{}, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntStringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntStringV(rv2i(rv).(map[int]string), e)
}
func (fastpathETJsonIO) EncMapIntStringV(v map[int]string, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntBytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntBytesV(rv2i(rv).(map[int][]byte), e)
}
func (fastpathETJsonIO) EncMapIntBytesV(v map[int][]byte, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntUint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntUint8V(rv2i(rv).(map[int]uint8), e)
}
func (fastpathETJsonIO) EncMapIntUint8V(v map[int]uint8, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntUint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntUint64V(rv2i(rv).(map[int]uint64), e)
}
func (fastpathETJsonIO) EncMapIntUint64V(v map[int]uint64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntIntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntIntV(rv2i(rv).(map[int]int), e)
}
func (fastpathETJsonIO) EncMapIntIntV(v map[int]int, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntInt32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntInt32V(rv2i(rv).(map[int]int32), e)
}
func (fastpathETJsonIO) EncMapIntInt32V(v map[int]int32, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntFloat64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntFloat64V(rv2i(rv).(map[int]float64), e)
}
func (fastpathETJsonIO) EncMapIntFloat64V(v map[int]float64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapIntBoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapIntBoolV(rv2i(rv).(map[int]bool), e)
}
func (fastpathETJsonIO) EncMapIntBoolV(v map[int]bool, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32IntfR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32IntfV(rv2i(rv).(map[int32]interface{}), e)
}
func (fastpathETJsonIO) EncMapInt32IntfV(v map[int32]interface{}, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32StringR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32StringV(rv2i(rv).(map[int32]string), e)
}
func (fastpathETJsonIO) EncMapInt32StringV(v map[int32]string, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32BytesR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32BytesV(rv2i(rv).(map[int32][]byte), e)
}
func (fastpathETJsonIO) EncMapInt32BytesV(v map[int32][]byte, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32Uint8R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32Uint8V(rv2i(rv).(map[int32]uint8), e)
}
func (fastpathETJsonIO) EncMapInt32Uint8V(v map[int32]uint8, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32Uint64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32Uint64V(rv2i(rv).(map[int32]uint64), e)
}
func (fastpathETJsonIO) EncMapInt32Uint64V(v map[int32]uint64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32IntR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32IntV(rv2i(rv).(map[int32]int), e)
}
func (fastpathETJsonIO) EncMapInt32IntV(v map[int32]int, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32Int32R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32Int32V(rv2i(rv).(map[int32]int32), e)
}
func (fastpathETJsonIO) EncMapInt32Int32V(v map[int32]int32, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32Float64R(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32Float64V(rv2i(rv).(map[int32]float64), e)
}
func (fastpathETJsonIO) EncMapInt32Float64V(v map[int32]float64, e *encoderJsonIO) {
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
func (e *encoderJsonIO) fastpathEncMapInt32BoolR(f *encFnInfo, rv reflect.Value) {
	fastpathETJsonIO{}.EncMapInt32BoolV(rv2i(rv).(map[int32]bool), e)
}
func (fastpathETJsonIO) EncMapInt32BoolV(v map[int32]bool, e *encoderJsonIO) {
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

func (helperDecDriverJsonIO) fastpathDecodeTypeSwitch(iv interface{}, d *decoderJsonIO) bool {
	var ft fastpathDTJsonIO
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

func (d *decoderJsonIO) fastpathDecSliceIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceIntfY(v []interface{}, d *decoderJsonIO) (v2 []interface{}, changed bool) {
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
func (fastpathDTJsonIO) DecSliceIntfN(v []interface{}, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceStringY(v []string, d *decoderJsonIO) (v2 []string, changed bool) {
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
func (fastpathDTJsonIO) DecSliceStringN(v []string, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceBytesY(v [][]byte, d *decoderJsonIO) (v2 [][]byte, changed bool) {
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
func (fastpathDTJsonIO) DecSliceBytesN(v [][]byte, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceFloat32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceFloat32Y(v []float32, d *decoderJsonIO) (v2 []float32, changed bool) {
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
func (fastpathDTJsonIO) DecSliceFloat32N(v []float32, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceFloat64Y(v []float64, d *decoderJsonIO) (v2 []float64, changed bool) {
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
func (fastpathDTJsonIO) DecSliceFloat64N(v []float64, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceUint8Y(v []uint8, d *decoderJsonIO) (v2 []uint8, changed bool) {
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
func (fastpathDTJsonIO) DecSliceUint8N(v []uint8, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceUint64Y(v []uint64, d *decoderJsonIO) (v2 []uint64, changed bool) {
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
func (fastpathDTJsonIO) DecSliceUint64N(v []uint64, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceIntY(v []int, d *decoderJsonIO) (v2 []int, changed bool) {
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
func (fastpathDTJsonIO) DecSliceIntN(v []int, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceInt32Y(v []int32, d *decoderJsonIO) (v2 []int32, changed bool) {
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
func (fastpathDTJsonIO) DecSliceInt32N(v []int32, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceInt64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceInt64Y(v []int64, d *decoderJsonIO) (v2 []int64, changed bool) {
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
func (fastpathDTJsonIO) DecSliceInt64N(v []int64, d *decoderJsonIO) {
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

func (d *decoderJsonIO) fastpathDecSliceBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecSliceBoolY(v []bool, d *decoderJsonIO) (v2 []bool, changed bool) {
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
func (fastpathDTJsonIO) DecSliceBoolN(v []bool, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringIntfL(v map[string]interface{}, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringStringL(v map[string]string, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringBytesL(v map[string][]byte, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringUint8L(v map[string]uint8, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringUint64L(v map[string]uint64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringIntL(v map[string]int, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringInt32L(v map[string]int32, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringFloat64L(v map[string]float64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapStringBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapStringBoolL(v map[string]bool, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8IntfL(v map[uint8]interface{}, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8StringL(v map[uint8]string, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8BytesL(v map[uint8][]byte, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8Uint8L(v map[uint8]uint8, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8Uint64L(v map[uint8]uint64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8IntL(v map[uint8]int, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8Int32L(v map[uint8]int32, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8Float64L(v map[uint8]float64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint8BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint8BoolL(v map[uint8]bool, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64IntfL(v map[uint64]interface{}, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64StringL(v map[uint64]string, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64BytesL(v map[uint64][]byte, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64Uint8L(v map[uint64]uint8, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64Uint64L(v map[uint64]uint64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64IntL(v map[uint64]int, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64Int32L(v map[uint64]int32, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64Float64L(v map[uint64]float64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapUint64BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapUint64BoolL(v map[uint64]bool, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntIntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntIntfL(v map[int]interface{}, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntStringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntStringL(v map[int]string, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntBytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntBytesL(v map[int][]byte, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntUint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntUint8L(v map[int]uint8, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntUint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntUint64L(v map[int]uint64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntIntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntIntL(v map[int]int, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntInt32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntInt32L(v map[int]int32, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntFloat64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntFloat64L(v map[int]float64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapIntBoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapIntBoolL(v map[int]bool, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32IntfR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32IntfL(v map[int32]interface{}, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32StringR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32StringL(v map[int32]string, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32BytesR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32BytesL(v map[int32][]byte, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32Uint8R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32Uint8L(v map[int32]uint8, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32Uint64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32Uint64L(v map[int32]uint64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32IntR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32IntL(v map[int32]int, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32Int32R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32Int32L(v map[int32]int32, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32Float64R(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32Float64L(v map[int32]float64, containerLen int, d *decoderJsonIO) {
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
func (d *decoderJsonIO) fastpathDecMapInt32BoolR(f *decFnInfo, rv reflect.Value) {
	var ft fastpathDTJsonIO
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
func (fastpathDTJsonIO) DecMapInt32BoolL(v map[int32]bool, containerLen int, d *decoderJsonIO) {
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
