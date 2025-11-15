//go:build !notfastpath && !codec.notfastpath

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from fastpath.go.tmpl - DO NOT EDIT.

package codec

// Fast path functions try to create a fast path encode or decode implementation
// for common maps and slices.
//
// We define the functions and register them in this single file
// so as not to pollute the encode.go and decode.go, and create a dependency in there.
// This file can be omitted without causing a build failure.
//
// The advantage of fast paths is:
//	  - Many calls bypass reflection altogether
//
// Currently support
//	  - slice of all builtin types (numeric, bool, string, []byte)
//    - maps of builtin types to builtin or interface{} type, EXCEPT FOR
//      keys of type uintptr, int8/16/32, uint16/32, float32/64, bool, interface{}
//      AND values of type type int8/16/32, uint16/32
// This should provide adequate "typical" implementations.
//
// Note that fast track decode functions must handle values for which an address cannot be obtained.
// For example:
//	 m2 := map[string]int{}
//	 p2 := []interface{}{m2}
//	 // decoding into p2 will bomb if fast track functions do not treat like unaddressable.
//

import (
	"reflect"
	"slices"
	"sort"
)

const fastpathEnabled = true

type fastpathARtid [56]uintptr

type fastpathRtRtid struct {
	rtid uintptr
	rt   reflect.Type
}
type fastpathARtRtid [56]fastpathRtRtid

var (
	fastpathAvRtidArr   fastpathARtid
	fastpathAvRtRtidArr fastpathARtRtid
	fastpathAvRtid      = fastpathAvRtidArr[:]
	fastpathAvRtRtid    = fastpathAvRtRtidArr[:]
)

func fastpathAvIndex(rtid uintptr) (i uint, ok bool) {
	return searchRtids(fastpathAvRtid, rtid)
}

func init() {
	var i uint = 0
	fn := func(v interface{}) {
		xrt := reflect.TypeOf(v)
		xrtid := rt2id(xrt)
		xptrtid := rt2id(reflect.PointerTo(xrt))
		fastpathAvRtid[i] = xrtid
		fastpathAvRtRtid[i] = fastpathRtRtid{rtid: xrtid, rt: xrt}
		encBuiltinRtids = append(encBuiltinRtids, xrtid, xptrtid)
		decBuiltinRtids = append(decBuiltinRtids, xrtid, xptrtid)
		i++
	}

	fn([]interface{}(nil))
	fn([]string(nil))
	fn([][]byte(nil))
	fn([]float32(nil))
	fn([]float64(nil))
	fn([]uint8(nil))
	fn([]uint64(nil))
	fn([]int(nil))
	fn([]int32(nil))
	fn([]int64(nil))
	fn([]bool(nil))

	fn(map[string]interface{}(nil))
	fn(map[string]string(nil))
	fn(map[string][]byte(nil))
	fn(map[string]uint8(nil))
	fn(map[string]uint64(nil))
	fn(map[string]int(nil))
	fn(map[string]int32(nil))
	fn(map[string]float64(nil))
	fn(map[string]bool(nil))
	fn(map[uint8]interface{}(nil))
	fn(map[uint8]string(nil))
	fn(map[uint8][]byte(nil))
	fn(map[uint8]uint8(nil))
	fn(map[uint8]uint64(nil))
	fn(map[uint8]int(nil))
	fn(map[uint8]int32(nil))
	fn(map[uint8]float64(nil))
	fn(map[uint8]bool(nil))
	fn(map[uint64]interface{}(nil))
	fn(map[uint64]string(nil))
	fn(map[uint64][]byte(nil))
	fn(map[uint64]uint8(nil))
	fn(map[uint64]uint64(nil))
	fn(map[uint64]int(nil))
	fn(map[uint64]int32(nil))
	fn(map[uint64]float64(nil))
	fn(map[uint64]bool(nil))
	fn(map[int]interface{}(nil))
	fn(map[int]string(nil))
	fn(map[int][]byte(nil))
	fn(map[int]uint8(nil))
	fn(map[int]uint64(nil))
	fn(map[int]int(nil))
	fn(map[int]int32(nil))
	fn(map[int]float64(nil))
	fn(map[int]bool(nil))
	fn(map[int32]interface{}(nil))
	fn(map[int32]string(nil))
	fn(map[int32][]byte(nil))
	fn(map[int32]uint8(nil))
	fn(map[int32]uint64(nil))
	fn(map[int32]int(nil))
	fn(map[int32]int32(nil))
	fn(map[int32]float64(nil))
	fn(map[int32]bool(nil))

	sort.Slice(fastpathAvRtid, func(i, j int) bool { return fastpathAvRtid[i] < fastpathAvRtid[j] })
	sort.Slice(fastpathAvRtRtid, func(i, j int) bool { return fastpathAvRtRtid[i].rtid < fastpathAvRtRtid[j].rtid })
	slices.Sort(encBuiltinRtids)
	slices.Sort(decBuiltinRtids)
}

func fastpathDecodeSetZeroTypeSwitch(iv interface{}) bool {
	switch v := iv.(type) {
	case *[]interface{}:
		*v = nil
	case *[]string:
		*v = nil
	case *[][]byte:
		*v = nil
	case *[]float32:
		*v = nil
	case *[]float64:
		*v = nil
	case *[]uint8:
		*v = nil
	case *[]uint64:
		*v = nil
	case *[]int:
		*v = nil
	case *[]int32:
		*v = nil
	case *[]int64:
		*v = nil
	case *[]bool:
		*v = nil

	case *map[string]interface{}:
		*v = nil
	case *map[string]string:
		*v = nil
	case *map[string][]byte:
		*v = nil
	case *map[string]uint8:
		*v = nil
	case *map[string]uint64:
		*v = nil
	case *map[string]int:
		*v = nil
	case *map[string]int32:
		*v = nil
	case *map[string]float64:
		*v = nil
	case *map[string]bool:
		*v = nil
	case *map[uint8]interface{}:
		*v = nil
	case *map[uint8]string:
		*v = nil
	case *map[uint8][]byte:
		*v = nil
	case *map[uint8]uint8:
		*v = nil
	case *map[uint8]uint64:
		*v = nil
	case *map[uint8]int:
		*v = nil
	case *map[uint8]int32:
		*v = nil
	case *map[uint8]float64:
		*v = nil
	case *map[uint8]bool:
		*v = nil
	case *map[uint64]interface{}:
		*v = nil
	case *map[uint64]string:
		*v = nil
	case *map[uint64][]byte:
		*v = nil
	case *map[uint64]uint8:
		*v = nil
	case *map[uint64]uint64:
		*v = nil
	case *map[uint64]int:
		*v = nil
	case *map[uint64]int32:
		*v = nil
	case *map[uint64]float64:
		*v = nil
	case *map[uint64]bool:
		*v = nil
	case *map[int]interface{}:
		*v = nil
	case *map[int]string:
		*v = nil
	case *map[int][]byte:
		*v = nil
	case *map[int]uint8:
		*v = nil
	case *map[int]uint64:
		*v = nil
	case *map[int]int:
		*v = nil
	case *map[int]int32:
		*v = nil
	case *map[int]float64:
		*v = nil
	case *map[int]bool:
		*v = nil
	case *map[int32]interface{}:
		*v = nil
	case *map[int32]string:
		*v = nil
	case *map[int32][]byte:
		*v = nil
	case *map[int32]uint8:
		*v = nil
	case *map[int32]uint64:
		*v = nil
	case *map[int32]int:
		*v = nil
	case *map[int32]int32:
		*v = nil
	case *map[int32]float64:
		*v = nil
	case *map[int32]bool:
		*v = nil

	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}
