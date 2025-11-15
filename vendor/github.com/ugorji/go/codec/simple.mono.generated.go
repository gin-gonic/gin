//go:build !notmono && !codec.notmono 

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"

	"io"
	"math"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
)

type helperEncDriverSimpleBytes struct{}
type encFnSimpleBytes struct {
	i  encFnInfo
	fe func(*encoderSimpleBytes, *encFnInfo, reflect.Value)
}
type encRtidFnSimpleBytes struct {
	rtid uintptr
	fn   *encFnSimpleBytes
}
type encoderSimpleBytes struct {
	dh helperEncDriverSimpleBytes
	fp *fastpathEsSimpleBytes
	e  simpleEncDriverBytes
	encoderBase
}
type helperDecDriverSimpleBytes struct{}
type decFnSimpleBytes struct {
	i  decFnInfo
	fd func(*decoderSimpleBytes, *decFnInfo, reflect.Value)
}
type decRtidFnSimpleBytes struct {
	rtid uintptr
	fn   *decFnSimpleBytes
}
type decoderSimpleBytes struct {
	dh helperDecDriverSimpleBytes
	fp *fastpathDsSimpleBytes
	d  simpleDecDriverBytes
	decoderBase
}
type simpleEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *SimpleHandle
	e *encoderBase

	w bytesEncAppender
}
type simpleDecDriverBytes struct {
	h *SimpleHandle
	d *decoderBase
	r bytesDecReader

	bdAndBdread

	noBuiltInTypes

	decDriverNoopContainerReader
	decInit2er
}

func (e *encoderSimpleBytes) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderSimpleBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderSimpleBytes) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderSimpleBytes) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderSimpleBytes) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderSimpleBytes) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderSimpleBytes) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderSimpleBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderSimpleBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderSimpleBytes) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderSimpleBytes) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderSimpleBytes) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderSimpleBytes) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderSimpleBytes) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderSimpleBytes) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderSimpleBytes) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderSimpleBytes) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderSimpleBytes) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderSimpleBytes) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderSimpleBytes) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderSimpleBytes) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderSimpleBytes) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderSimpleBytes) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderSimpleBytes) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderSimpleBytes) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderSimpleBytes) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderSimpleBytes) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderSimpleBytes) kSeqFn(rt reflect.Type) (fn *encFnSimpleBytes) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderSimpleBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
	var l int
	if isSlice {
		l = rvLenSlice(rv)
	} else {
		l = rv.Len()
	}
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(l)
	e.mapStart(l >> 1)

	var fn *encFnSimpleBytes
	builtin := ti.tielem.flagEncBuiltin
	if !builtin {
		fn = e.kSeqFn(ti.elem)
	}

	j := 0
	e.c = containerMapKey
	e.e.WriteMapElemKey(true)
	for {
		rvv := rvArrayIndex(rv, j, ti, isSlice)
		if builtin {
			e.encodeIB(rv2i(baseRVRV(rvv)))
		} else {
			e.encodeValue(rvv, fn)
		}
		j++
		if j == l {
			break
		}
		if j&1 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(false)
		} else {
			e.mapElemValue()
		}
	}
	e.c = 0
	e.e.WriteMapEnd()

}

func (e *encoderSimpleBytes) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
	var l int
	if isSlice {
		l = rvLenSlice(rv)
	} else {
		l = rv.Len()
	}
	if l <= 0 {
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(l)

	var fn *encFnSimpleBytes
	if !ti.tielem.flagEncBuiltin {
		fn = e.kSeqFn(ti.elem)
	}

	j := 0
	e.c = containerArrayElem
	e.e.WriteArrayElem(true)
	builtin := ti.tielem.flagEncBuiltin
	for {
		rvv := rvArrayIndex(rv, j, ti, isSlice)
		if builtin {
			e.encodeIB(rv2i(baseRVRV(rvv)))
		} else {
			e.encodeValue(rvv, fn)
		}
		j++
		if j == l {
			break
		}
		e.c = containerArrayElem
		e.e.WriteArrayElem(false)
	}

	e.c = 0
	e.e.WriteArrayEnd()
}

func (e *encoderSimpleBytes) kChan(f *encFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		halt.errorStr("send-only channel cannot be encoded")
	}
	if !f.ti.mbs && uint8TypId == rt2id(f.ti.elem) {
		e.kSliceBytesChan(rv)
		return
	}
	rtslice := reflect.SliceOf(f.ti.elem)
	rv = chanToSlice(rv, rtslice, e.h.ChanRecvTimeout)
	ti := e.h.getTypeInfo(rt2id(rtslice), rtslice)
	if f.ti.mbs {
		e.kArrayWMbs(rv, ti, true)
	} else {
		e.kArrayW(rv, ti, true)
	}
}

func (e *encoderSimpleBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderSimpleBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderSimpleBytes) kSliceBytesChan(rv reflect.Value) {

	bs0 := e.blist.peek(32, true)
	bs := bs0

	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok {
		ch = irv.(chan byte)
	}

L1:
	switch timeout := e.h.ChanRecvTimeout; {
	case timeout == 0:
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			}
		}
	case timeout > 0:
		tt := time.NewTimer(timeout)
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:

				break L1
			}
		}
	default:
		for b := range ch {
			bs = append(bs, b)
		}
	}

	e.e.EncodeBytes(bs)
	e.blist.put(bs)
	if !byteSliceSameData(bs0, bs) {
		e.blist.put(bs0)
	}
}

func (e *encoderSimpleBytes) kStructFieldKey(keyType valueType, encName string) {

	if keyType == valueTypeString {
		e.e.EncodeString(encName)
	} else if keyType == valueTypeInt {
		e.e.EncodeInt(must.Int(strconv.ParseInt(encName, 10, 64)))
	} else if keyType == valueTypeUint {
		e.e.EncodeUint(must.Uint(strconv.ParseUint(encName, 10, 64)))
	} else if keyType == valueTypeFloat {
		e.e.EncodeFloat64(must.Float(strconv.ParseFloat(encName, 64)))
	} else {
		halt.errorStr2("invalid struct key type: ", keyType.String())
	}

}

func (e *encoderSimpleBytes) kStructSimple(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	tisfi := f.ti.sfi.source()

	chkCirRef := e.h.CheckCircularRef
	var si *structFieldInfo
	var j int

	if f.ti.toArray || e.h.StructToArray {
		if len(tisfi) == 0 {
			e.e.WriteArrayEmpty()
			return
		}
		e.arrayStart(len(tisfi))
		for j, si = range tisfi {
			e.c = containerArrayElem
			e.e.WriteArrayElem(j == 0)
			if si.encBuiltin {
				e.encodeIB(rv2i(si.fieldNoAlloc(rv, true)))
			} else {
				e.encodeValue(si.fieldNoAlloc(rv, !chkCirRef), nil)
			}
		}
		e.c = 0
		e.e.WriteArrayEnd()
	} else {
		if len(tisfi) == 0 {
			e.e.WriteMapEmpty()
			return
		}
		if e.h.Canonical {
			tisfi = f.ti.sfi.sorted()
		}
		e.mapStart(len(tisfi))
		for j, si = range tisfi {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
			e.e.EncodeStringNoEscape4Json(si.encName)
			e.mapElemValue()
			if si.encBuiltin {
				e.encodeIB(rv2i(si.fieldNoAlloc(rv, true)))
			} else {
				e.encodeValue(si.fieldNoAlloc(rv, !chkCirRef), nil)
			}
		}
		e.c = 0
		e.e.WriteMapEnd()
	}
}

func (e *encoderSimpleBytes) kStruct(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	ti := f.ti
	toMap := !(ti.toArray || e.h.StructToArray)
	var mf map[string]interface{}
	if ti.flagMissingFielder {
		toMap = true
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
	} else if ti.flagMissingFielderPtr {
		toMap = true
		if rv.CanAddr() {
			mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder).CodecMissingFields()
		} else {
			mf = rv2i(e.addrRV(rv, ti.rt, ti.ptr)).(MissingFielder).CodecMissingFields()
		}
	}
	newlen := len(mf)
	tisfi := ti.sfi.source()
	newlen += len(tisfi)

	var fkvs = e.slist.get(newlen)[:newlen]

	recur := e.h.RecursiveEmptyCheck
	chkCirRef := e.h.CheckCircularRef

	var xlen int

	var kv sfiRv
	var j int
	var sf encStructFieldObj
	if toMap {
		newlen = 0
		if e.h.Canonical {
			tisfi = f.ti.sfi.sorted()
		}
		for _, si := range tisfi {

			if si.omitEmpty {
				kv.r = si.fieldNoAlloc(rv, false)
				if isEmptyValue(kv.r, e.h.TypeInfos, recur) {
					continue
				}
			} else {
				kv.r = si.fieldNoAlloc(rv, si.encBuiltin || !chkCirRef)
			}
			kv.v = si
			fkvs[newlen] = kv
			newlen++
		}

		var mf2s []stringIntf
		if len(mf) != 0 {
			mf2s = make([]stringIntf, 0, len(mf))
			for k, v := range mf {
				if k == "" {
					continue
				}
				if ti.infoFieldOmitempty && isEmptyValue(reflect.ValueOf(v), e.h.TypeInfos, recur) {
					continue
				}
				mf2s = append(mf2s, stringIntf{k, v})
			}
		}

		xlen = newlen + len(mf2s)
		if xlen == 0 {
			e.e.WriteMapEmpty()
			goto END
		}

		e.mapStart(xlen)

		if len(mf2s) != 0 && e.h.Canonical {
			mf2w := make([]encStructFieldObj, newlen+len(mf2s))
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				mf2w[j] = encStructFieldObj{kv.v.encName, kv.r, nil, true,
					!kv.v.encNameEscape4Json, kv.v.encBuiltin}
			}
			for _, v := range mf2s {
				mf2w[j] = encStructFieldObj{v.v, reflect.Value{}, v.i, false, false, false}
				j++
			}
			sort.Sort((encStructFieldObjSlice)(mf2w))
			for j, sf = range mf2w {
				e.c = containerMapKey
				e.e.WriteMapElemKey(j == 0)
				if ti.keyType == valueTypeString && sf.noEsc4json {
					e.e.EncodeStringNoEscape4Json(sf.key)
				} else {
					e.kStructFieldKey(ti.keyType, sf.key)
				}
				e.mapElemValue()
				if sf.isRv {
					if sf.builtin {
						e.encodeIB(rv2i(baseRVRV(sf.rv)))
					} else {
						e.encodeValue(sf.rv, nil)
					}
				} else {
					if !e.encodeBuiltin(sf.intf) {
						e.encodeR(reflect.ValueOf(sf.intf))
					}

				}
			}
		} else {
			keytyp := ti.keyType
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				e.c = containerMapKey
				e.e.WriteMapElemKey(j == 0)
				if ti.keyType == valueTypeString && !kv.v.encNameEscape4Json {
					e.e.EncodeStringNoEscape4Json(kv.v.encName)
				} else {
					e.kStructFieldKey(keytyp, kv.v.encName)
				}
				e.mapElemValue()
				if kv.v.encBuiltin {
					e.encodeIB(rv2i(baseRVRV(kv.r)))
				} else {
					e.encodeValue(kv.r, nil)
				}
			}
			for _, v := range mf2s {
				e.c = containerMapKey
				e.e.WriteMapElemKey(j == 0)
				e.kStructFieldKey(keytyp, v.v)
				e.mapElemValue()
				if !e.encodeBuiltin(v.i) {
					e.encodeR(reflect.ValueOf(v.i))
				}

				j++
			}
		}

		e.c = 0
		e.e.WriteMapEnd()
	} else {
		newlen = len(tisfi)
		for i, si := range tisfi {

			if si.omitEmpty {

				kv.r = si.fieldNoAlloc(rv, false)
				if isEmptyContainerValue(kv.r, e.h.TypeInfos, recur) {
					kv.r = reflect.Value{}
				}
			} else {
				kv.r = si.fieldNoAlloc(rv, si.encBuiltin || !chkCirRef)
			}
			kv.v = si
			fkvs[i] = kv
		}

		if newlen == 0 {
			e.e.WriteArrayEmpty()
			goto END
		}

		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ {
			e.c = containerArrayElem
			e.e.WriteArrayElem(j == 0)
			kv = fkvs[j]
			if !kv.r.IsValid() {
				e.e.EncodeNil()
			} else if kv.v.encBuiltin {
				e.encodeIB(rv2i(baseRVRV(kv.r)))
			} else {
				e.encodeValue(kv.r, nil)
			}
		}
		e.c = 0
		e.e.WriteArrayEnd()
	}

END:

	e.slist.put(fkvs)
}

func (e *encoderSimpleBytes) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnSimpleBytes

	ktypeKind := reflect.Kind(f.ti.keykind)
	vtypeKind := reflect.Kind(f.ti.elemkind)

	rtval := f.ti.elem
	rtvalkind := vtypeKind
	for rtvalkind == reflect.Ptr {
		rtval = rtval.Elem()
		rtvalkind = rtval.Kind()
	}
	if rtvalkind != reflect.Interface {
		valFn = e.fn(rtval)
	}

	var rvv = mapAddrLoopvarRV(f.ti.elem, vtypeKind)

	rtkey := f.ti.key
	var keyTypeIsString = stringTypId == rt2id(rtkey)
	if keyTypeIsString {
		keyFn = e.fn(rtkey)
	} else {
		for rtkey.Kind() == reflect.Ptr {
			rtkey = rtkey.Elem()
		}
		if rtkey.Kind() != reflect.Interface {
			keyFn = e.fn(rtkey)
		}
	}

	if e.h.Canonical {
		e.kMapCanonical(f.ti, rv, rvv, keyFn, valFn)
		e.c = 0
		e.e.WriteMapEnd()
		return
	}

	var rvk = mapAddrLoopvarRV(f.ti.key, ktypeKind)

	var it mapIter
	mapRange(&it, rv, rvk, rvv, true)

	kbuiltin := f.ti.tikey.flagEncBuiltin
	vbuiltin := f.ti.tielem.flagEncBuiltin
	for j := 0; it.Next(); j++ {
		rv = it.Key()
		e.c = containerMapKey
		e.e.WriteMapElemKey(j == 0)
		if keyTypeIsString {
			e.e.EncodeString(rvGetString(rv))
		} else if kbuiltin {
			e.encodeIB(rv2i(baseRVRV(rv)))
		} else {
			e.encodeValue(rv, keyFn)
		}
		e.mapElemValue()
		rv = it.Value()
		if vbuiltin {
			e.encodeIB(rv2i(baseRVRV(rv)))
		} else {
			e.encodeValue(it.Value(), valFn)
		}
	}
	it.Done()

	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoderSimpleBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnSimpleBytes) {
	_ = e.e

	rtkey := ti.key
	rtkeydecl := rtkey.PkgPath() == "" && rtkey.Name() != ""

	mks := rv.MapKeys()
	rtkeyKind := rtkey.Kind()
	mparams := getMapReqParams(ti)

	switch rtkeyKind {
	case reflect.Bool:

		if len(mks) == 2 && mks[0].Bool() {
			mks[0], mks[1] = mks[1], mks[0]
		}
		for i := range mks {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeBool(mks[i].Bool())
			} else {
				e.encodeValueNonNil(mks[i], keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mks[i], rvv, mparams), valFn)
		}
	case reflect.String:
		mksv := make([]orderedRv[string], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = rvGetString(k)
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeString(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		mksv := make([]orderedRv[uint64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Uint()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeUint(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		mksv := make([]orderedRv[int64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Int()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeInt(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Float32:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeFloat32(float32(mksv[i].v))
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Float64:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeFloat64(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	default:
		if rtkey == timeTyp {
			mksv := make([]timeRv, len(mks))
			for i, k := range mks {
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			}
			slices.SortFunc(mksv, cmpTimeRv)
			for i := range mksv {
				e.c = containerMapKey
				e.e.WriteMapElemKey(i == 0)
				e.e.EncodeTime(mksv[i].v)
				e.mapElemValue()
				e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
			}
			break
		}

		bs0 := e.blist.get(len(mks) * 16)
		mksv := bs0
		mksbv := make([]bytesRv, len(mks))

		sideEncode(e.hh, &e.h.sideEncPool, func(se encoderI) {
			se.ResetBytes(&mksv)
			for i, k := range mks {
				v := &mksbv[i]
				l := len(mksv)
				se.setContainerState(containerMapKey)
				se.encodeR(baseRVRV(k))
				se.atEndOfEncode()
				se.writerEnd()
				v.r = k
				v.v = mksv[l:]
			}
		})

		slices.SortFunc(mksbv, cmpBytesRv)
		for j := range mksbv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
			e.e.writeBytesAsis(mksbv[j].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksbv[j].r, rvv, mparams), valFn)
		}
		e.blist.put(mksv)
		if !byteSliceSameData(bs0, mksv) {
			e.blist.put(bs0)
		}
	}
}

func (e *encoderSimpleBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsSimpleBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderSimpleBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderSimpleBytes) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderSimpleBytes) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderSimpleBytes) mustEncode(v interface{}) {
	halt.onerror(e.err)
	if e.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	e.calls++
	if !e.encodeBuiltin(v) {
		e.encodeR(reflect.ValueOf(v))
	}

	e.calls--
	if e.calls == 0 {
		e.e.atEndOfEncode()
		e.e.writerEnd()
	}
}

func (e *encoderSimpleBytes) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderSimpleBytes) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderSimpleBytes) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderSimpleBytes) encodeBuiltin(iv interface{}) (ok bool) {
	ok = true
	switch v := iv.(type) {
	case nil:
		e.e.EncodeNil()

	case Raw:
		e.rawBytes(v)
	case string:
		e.e.EncodeString(v)
	case bool:
		e.e.EncodeBool(v)
	case int:
		e.e.EncodeInt(int64(v))
	case int8:
		e.e.EncodeInt(int64(v))
	case int16:
		e.e.EncodeInt(int64(v))
	case int32:
		e.e.EncodeInt(int64(v))
	case int64:
		e.e.EncodeInt(v)
	case uint:
		e.e.EncodeUint(uint64(v))
	case uint8:
		e.e.EncodeUint(uint64(v))
	case uint16:
		e.e.EncodeUint(uint64(v))
	case uint32:
		e.e.EncodeUint(uint64(v))
	case uint64:
		e.e.EncodeUint(v)
	case uintptr:
		e.e.EncodeUint(uint64(v))
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)
	case complex64:
		e.encodeComplex64(v)
	case complex128:
		e.encodeComplex128(v)
	case time.Time:
		e.e.EncodeTime(v)
	case []byte:
		e.e.EncodeBytes(v)
	default:

		ok = !skipFastpathTypeSwitchInDirectCall && e.dh.fastpathEncodeTypeSwitch(iv, e)
	}
	return
}

func (e *encoderSimpleBytes) encodeValue(rv reflect.Value, fn *encFnSimpleBytes) {

	var ciPushes int

	var rvp reflect.Value
	var rvpValid bool

RV:
	switch rv.Kind() {
	case reflect.Ptr:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
		rvpValid = true
		rvp = rv
		rv = rv.Elem()

		if e.h.CheckCircularRef && e.ci.canPushElemKind(rv.Kind()) {
			e.ci.push(rv2i(rvp))
			ciPushes++
		}
		goto RV
	case reflect.Interface:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
		rvpValid = false
		rvp = reflect.Value{}
		rv = rv.Elem()
		fn = nil
		goto RV
	case reflect.Map:
		if rvIsNil(rv) {
			if e.h.NilCollectionToZeroLength {
				e.e.WriteMapEmpty()
			} else {
				e.e.EncodeNil()
			}
			goto END
		}
	case reflect.Slice, reflect.Chan:
		if rvIsNil(rv) {
			if e.h.NilCollectionToZeroLength {
				e.e.WriteArrayEmpty()
			} else {
				e.e.EncodeNil()
			}
			goto END
		}
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		goto END
	}

	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if !fn.i.addrE {

	} else if rvpValid {
		rv = rvp
	} else if rv.CanAddr() {
		rv = rvAddr(rv, fn.i.ti.ptr)
	} else {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)

END:
	if ciPushes > 0 {
		e.ci.pop(ciPushes)
	}
}

func (e *encoderSimpleBytes) encodeValueNonNil(rv reflect.Value, fn *encFnSimpleBytes) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderSimpleBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderSimpleBytes) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderSimpleBytes) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderSimpleBytes) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderSimpleBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderSimpleBytes) fn(t reflect.Type) *encFnSimpleBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderSimpleBytes) fnNoExt(t reflect.Type) *encFnSimpleBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderSimpleBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderSimpleBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderSimpleBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderSimpleBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderSimpleBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderSimpleBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderSimpleBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderSimpleBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverSimpleBytes) newEncoderBytes(out *[]byte, h Handle) *encoderSimpleBytes {
	var c1 encoderSimpleBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverSimpleBytes) newEncoderIO(out io.Writer, h Handle) *encoderSimpleBytes {
	var c1 encoderSimpleBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverSimpleBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsSimpleBytes) (f *fastpathESimpleBytes, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

func (helperEncDriverSimpleBytes) encFindRtidFn(s []encRtidFnSimpleBytes, rtid uintptr) (i uint, fn *encFnSimpleBytes) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperEncDriverSimpleBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnSimpleBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnSimpleBytes](v))
	}
	return
}

func (dh helperEncDriverSimpleBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsSimpleBytes, checkExt bool) (fn *encFnSimpleBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverSimpleBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnSimpleBytes = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverSimpleBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnSimpleBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnSimpleBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnSimpleBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnSimpleBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverSimpleBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleBytes) {
	fn = new(encFnSimpleBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderSimpleBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderSimpleBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderSimpleBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderSimpleBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderSimpleBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderSimpleBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderSimpleBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderSimpleBytes).textMarshal
		fi.addrE = ti.flagTextMarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {

			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fe = fp[idx].encfn
				}
			} else {

				xfe, xrt := dh.encFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf := xfe.encfn
					fn.fe = func(e *encoderSimpleBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderSimpleBytes).kBool
			case reflect.String:

				fn.fe = (*encoderSimpleBytes).kString
			case reflect.Int:
				fn.fe = (*encoderSimpleBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderSimpleBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderSimpleBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderSimpleBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderSimpleBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderSimpleBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderSimpleBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderSimpleBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderSimpleBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderSimpleBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderSimpleBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderSimpleBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderSimpleBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderSimpleBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderSimpleBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderSimpleBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderSimpleBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderSimpleBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderSimpleBytes).kStructSimple
				} else {
					fn.fe = (*encoderSimpleBytes).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderSimpleBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderSimpleBytes).kErr
			default:

				fn.fe = (*encoderSimpleBytes).kErr
			}
		}
	}
	return
}
func (d *decoderSimpleBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderSimpleBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderSimpleBytes) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderSimpleBytes) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderSimpleBytes) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderSimpleBytes) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderSimpleBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderSimpleBytes) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderSimpleBytes) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderSimpleBytes) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderSimpleBytes) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderSimpleBytes) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderSimpleBytes) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderSimpleBytes) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderSimpleBytes) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderSimpleBytes) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderSimpleBytes) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderSimpleBytes) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderSimpleBytes) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderSimpleBytes) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderSimpleBytes) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderSimpleBytes) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleBytes) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleBytes) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderSimpleBytes) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderSimpleBytes) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderSimpleBytes) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderSimpleBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

	n := d.naked()
	d.d.DecodeNaked()

	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		halt.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
	}

	switch n.v {
	case valueTypeMap:
		mtid := d.mtid
		if mtid == 0 {
			if d.jsms {
				mtid = mapStrIntfTypId
			} else {
				mtid = mapIntfIntfTypId
			}
		}
		if mtid == mapStrIntfTypId {
			var v2 map[string]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if mtid == mapIntfIntfTypId {
			var v2 map[interface{}]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.mtr {
			rvn = reflect.New(d.h.MapType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {

			rvn = rvZeroAddrK(d.h.MapType, reflect.Map)
			d.decodeValue(rvn, nil)
		}
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId {
			var v2 []interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.str {
			rvn = reflect.New(d.h.SliceType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.SliceType, reflect.Slice)
			d.decodeValue(rvn, nil)
		}
		if d.h.PreferArrayOverSlice {
			rvn = rvGetArray4Slice(rvn)
		}
	case valueTypeExt:
		tag, bytes := n.u, n.l
		bfn := d.h.getExtForTag(tag)
		var re = RawExt{Tag: tag}
		if bytes == nil {

			if bfn == nil {
				d.decode(&re.Value)
				rvn = rv4iptr(&re).Elem()
			} else if bfn.ext == SelfExt {
				rvn = rvZeroAddrK(bfn.rt, bfn.rt.Kind())
				d.decodeValue(rvn, d.fnNoExt(bfn.rt))
			} else {
				rvn = reflect.New(bfn.rt)
				d.interfaceExtConvertAndDecode(rv2i(rvn), bfn.ext)
				rvn = rvn.Elem()
			}
		} else {

			if bfn == nil {
				re.setData(bytes, false)
				rvn = rv4iptr(&re).Elem()
			} else {
				rvn = reflect.New(bfn.rt)
				if bfn.ext == SelfExt {
					sideDecode(d.hh, &d.h.sideDecPool, func(sd decoderI) { oneOffDecode(sd, rv2i(rvn), bytes, bfn.rt, false) })
				} else {
					bfn.ext.ReadExt(rv2i(rvn), bytes)
				}
				rvn = rvn.Elem()
			}
		}

		if d.h.PreferPointerForStructOrArray && rvn.CanAddr() {
			if rk := rvn.Kind(); rk == reflect.Array || rk == reflect.Struct {
				rvn = rvn.Addr()
			}
		}
	case valueTypeNil:

	case valueTypeInt:
		rvn = n.ri()
	case valueTypeUint:
		rvn = n.ru()
	case valueTypeFloat:
		rvn = n.rf()
	case valueTypeBool:
		rvn = n.rb()
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs()
	case valueTypeBytes:
		rvn = n.rl()
	case valueTypeTime:
		rvn = n.rt()
	default:
		halt.errorStr2("kInterfaceNaked: unexpected valueType: ", n.v.String())
	}
	return
}

func (d *decoderSimpleBytes) kInterface(f *decFnInfo, rv reflect.Value) {

	isnilrv := rvIsNil(rv)

	var rvn reflect.Value

	if d.h.InterfaceReset {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			} else if !isnilrv {
				decSetNonNilRV2Zero4Intf(rv)
			}
			return
		}
	} else if isnilrv {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			}
			return
		}
	} else {

		rvn = rv.Elem()
	}

	canDecode, _ := isDecodeable(rvn)

	if !canDecode {
		rvn2 := d.oneShotAddrRV(rvn.Type(), rvn.Kind())
		rvSetDirect(rvn2, rvn)
		rvn = rvn2
	}

	d.decodeValue(rvn, nil)
	rvSetIntf(rv, rvn)
}

func (d *decoderSimpleBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		rv = si.fieldNoAlloc(rv, true)
		if rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
	} else if si.decBuiltin {
		rv = rvAddr(si.fieldAlloc(rv), si.ptrTyp)
		d.decode(rv2i(rv))
	} else {
		fn := d.fn(si.baseTyp)
		rv = si.fieldAlloc(rv)
		if fn.i.addrD {
			rv = rvAddr(rv, si.ptrTyp)
		}
		fn.fd(d, &fn.i, rv)
	}
}

func (d *decoderSimpleBytes) kStructSimple(f *decFnInfo, rv reflect.Value) {
	_ = d.d
	ctyp := d.d.ContainerType()
	ti := f.ti
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var rvkencname []byte
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey(j == 0)
			sab, att := d.d.DecodeStringAsBytes()
			rvkencname = d.usableStructFieldNameBytes(rvkencname, sab, att)
			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}

		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem(j == 0)
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}
		d.arrayEnd()
	} else {
		halt.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderSimpleBytes) kStruct(f *decFnInfo, rv reflect.Value) {
	_ = d.d
	ctyp := d.d.ContainerType()
	ti := f.ti
	var mf MissingFielder
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder)
	} else if ti.flagMissingFielderPtr {
		mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder)
	}
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var name2 []byte
		var rvkencname []byte
		tkt := ti.keyType
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey(j == 0)

			if tkt == valueTypeString {
				sab, att := d.d.DecodeStringAsBytes()
				rvkencname = d.usableStructFieldNameBytes(rvkencname, sab, att)
			} else if tkt == valueTypeInt {
				rvkencname = strconv.AppendInt(d.b[:0], d.d.DecodeInt64(), 10)
			} else if tkt == valueTypeUint {
				rvkencname = strconv.AppendUint(d.b[:0], d.d.DecodeUint64(), 10)
			} else if tkt == valueTypeFloat {
				rvkencname = strconv.AppendFloat(d.b[:0], d.d.DecodeFloat64(), 'f', -1, 64)
			} else {
				halt.errorStr2("invalid struct key type: ", ti.keyType.String())
			}

			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else if mf != nil {

				name2 = append(name2[:0], rvkencname...)
				var f interface{}
				d.decode(&f)
				if !mf.CodecMissingField(name2, f) && d.h.ErrorIfNoField {
					halt.errorStr2("no matching struct field when decoding stream map with key: ", stringView(name2))
				}
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}

		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem(j == 0)
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}

		d.arrayEnd()
	} else {
		halt.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderSimpleBytes) kSlice(f *decFnInfo, rv reflect.Value) {
	_ = d.d

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		rvbs := rvGetBytes(rv)
		if rvCanset {
			bs2, bst := d.decodeBytesInto(rvbs, false)
			if bst != dBytesIntoParamOut {
				rvSetBytes(rv, bs2)
			}
		} else {

			d.decodeBytesInto(rvbs[:len(rvbs):len(rvbs)], true)
		}
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

	if containerLenS == 0 {
		if rvCanset {
			if rvIsNil(rv) {
				rvSetDirect(rv, rvSliceZeroCap(ti.rt))
			} else {
				rvSetSliceLen(rv, 0)
			}
		}
		if isArray {
			d.arrayEnd()
		} else {
			d.mapEnd()
		}
		return
	}

	rtelem0Mut := !scalarBitset.isset(ti.elemkind)
	rtelem := ti.elem

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnSimpleBytes

	var rvChanged bool

	var rv0 = rv
	var rv9 reflect.Value

	rvlen := rvLenSlice(rv)
	rvcap := rvCapSlice(rv)
	maxInitLen := d.maxInitLen()
	hasLen := containerLenS >= 0
	if hasLen {
		if containerLenS > rvcap {
			oldRvlenGtZero := rvlen > 0
			rvlen1 := int(decInferLen(containerLenS, maxInitLen, uint(ti.elemsize)))
			if rvlen1 == rvlen {
			} else if rvlen1 <= rvcap {
				if rvCanset {
					rvlen = rvlen1
					rvSetSliceLen(rv, rvlen)
				}
			} else if rvCanset {
				rvlen = rvlen1
				rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = !rvCanset
			} else {
				halt.errorStr("cannot decode into non-settable slice")
			}
			if rvChanged && oldRvlenGtZero && rtelem0Mut {
				rvCopySlice(rv, rv0, rtelem)
			}
		} else if containerLenS != rvlen {
			if rvCanset {
				rvlen = containerLenS
				rvSetSliceLen(rv, rvlen)
			}
		}
	}

	var elemReset = d.h.SliceElementReset

	var rtelemIsPtr bool
	var rtelemElem reflect.Type
	builtin := ti.tielem.flagDecBuiltin
	if builtin {
		rtelemIsPtr = ti.elemkind == uint8(reflect.Ptr)
		if rtelemIsPtr {
			rtelemElem = ti.elem.Elem()
		}
	}

	var j int
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if rvCanset {
					rvlen = int(decInferLen(containerLenS, maxInitLen, uint(ti.elemsize)))
					rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
					rvcap = rvlen
					rvChanged = !rvCanset
				} else {
					halt.errorStr("cannot decode into non-settable slice")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}

		if ctyp == valueTypeArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}

		if j >= rvlen {

			if rvlen < rvcap {
				rvlen = rvcap
				if rvCanset {
					rvSetSliceLen(rv, rvlen)
				} else if rvChanged {
					rv = rvSlice(rv, rvlen)
				} else {
					halt.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					halt.onerror(errExpandSliceCannotChange)
				}
				rv, rvcap, rvCanset = rvGrowSlice(rv, f.ti, rvcap, 1)

				rvlen = rvcap
				rvChanged = !rvCanset
			}
		}

		rv9 = rvArrayIndex(rv, j, f.ti, true)
		if elemReset {
			rvSetZero(rv9)
		}
		if d.d.TryNil() {
			rvSetZero(rv9)
		} else if builtin {
			if rtelemIsPtr {
				if rvIsNil(rv9) {
					rvSetDirect(rv9, reflect.New(rtelemElem))
				}
				d.decode(rv2i(rv9))
			} else {
				d.decode(rv2i(rvAddr(rv9, ti.tielem.ptr)))
			}
		} else {
			d.decodeValueNoCheckNil(rv9, fn)
		}
	}
	if j < rvlen {
		if rvCanset {
			rvSetSliceLen(rv, j)
		} else if rvChanged {
			rv = rvSlice(rv, j)
		}

	} else if j == 0 && rvIsNil(rv) {
		if rvCanset {
			rv = rvSliceZeroCap(ti.rt)
			rvCanset = false
			rvChanged = true
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}

	if rvChanged {
		rvSetDirect(rv0, rv)
	}
}

func (d *decoderSimpleBytes) kArray(f *decFnInfo, rv reflect.Value) {
	_ = d.d

	ti := f.ti
	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {

		if ti.elemkind != uint8(reflect.Uint8) {
			halt.errorf("bytes/string in stream can decode into array of bytes, but not %v", ti.rt)
		}
		rvbs := rvGetArrayBytes(rv, nil)
		d.decodeBytesInto(rvbs, true)
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

	if containerLenS == 0 {
		if isArray {
			d.arrayEnd()
		} else {
			d.mapEnd()
		}
		return
	}

	rtelem := ti.elem
	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var rv9 reflect.Value

	rvlen := rv.Len()
	hasLen := containerLenS >= 0
	if hasLen && containerLenS > rvlen {
		halt.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
	}

	var elemReset = d.h.SliceElementReset

	var rtelemIsPtr bool
	var rtelemElem reflect.Type
	var fn *decFnSimpleBytes
	builtin := ti.tielem.flagDecBuiltin
	if builtin {
		rtelemIsPtr = ti.elemkind == uint8(reflect.Ptr)
		if rtelemIsPtr {
			rtelemElem = ti.elem.Elem()
		}
	} else {
		fn = d.fn(rtelem)
	}

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if ctyp == valueTypeArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}

		if j >= rvlen {
			d.arrayCannotExpand(rvlen, j+1)
			d.swallow()
			continue
		}

		rv9 = rvArrayIndex(rv, j, f.ti, false)
		if elemReset {
			rvSetZero(rv9)
		}
		if d.d.TryNil() {
			rvSetZero(rv9)
		} else if builtin {
			if rtelemIsPtr {
				if rvIsNil(rv9) {
					rvSetDirect(rv9, reflect.New(rtelemElem))
				}
				d.decode(rv2i(rv9))
			} else {
				d.decode(rv2i(rvAddr(rv9, ti.tielem.ptr)))
			}
		} else {
			d.decodeValueNoCheckNil(rv9, fn)
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoderSimpleBytes) kChan(f *decFnInfo, rv reflect.Value) {
	_ = d.d

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		halt.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2, _ := d.d.DecodeBytes()
		irv := rv2i(rv)
		ch, ok := irv.(chan<- byte)
		if !ok {
			ch = irv.(chan byte)
		}
		for _, b := range bs2 {
			ch <- b
		}
		return
	}

	var rvCanset = rv.CanSet()

	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}

	if containerLenS == 0 {
		if rvCanset && rvIsNil(rv) {
			rvSetDirect(rv, reflect.MakeChan(ti.rt, 0))
		}
		if isArray {
			d.arrayEnd()
		} else {
			d.mapEnd()
		}
		return
	}

	rtelem := ti.elem
	useTransient := decUseTransient && ti.elemkind != byte(reflect.Ptr) && ti.tielem.flagCanTransient

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnSimpleBytes

	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int
	hasLen := containerLenS >= 0
	maxInitLen := d.maxInitLen()

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if hasLen {
					rvlen = int(decInferLen(containerLenS, maxInitLen, uint(ti.elemsize)))
				} else {
					rvlen = decDefChanCap
				}
				if rvCanset {
					rv = reflect.MakeChan(ti.rt, rvlen)
					rvChanged = true
				} else {
					halt.errorStr("cannot decode into non-settable chan")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}

		if ctyp == valueTypeArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}

		if rv9.IsValid() {
			rvSetZero(rv9)
		} else if useTransient {
			rv9 = d.perType.TransientAddrK(ti.elem, reflect.Kind(ti.elemkind))
		} else {
			rv9 = rvZeroAddrK(ti.elem, reflect.Kind(ti.elemkind))
		}
		if !d.d.TryNil() {
			d.decodeValueNoCheckNil(rv9, fn)
		}
		rv.Send(rv9)
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}

	if rvChanged {
		rvSetDirect(rv0, rv)
	}

}

func (d *decoderSimpleBytes) kMap(f *decFnInfo, rv reflect.Value) {
	_ = d.d
	containerLen := d.mapStart(d.d.ReadMapStart())
	ti := f.ti
	if rvIsNil(rv) {
		rvlen := int(decInferLen(containerLen, d.maxInitLen(), uint(ti.keysize+ti.elemsize)))
		rvSetDirect(rv, makeMapReflect(ti.rt, rvlen))
	}

	if containerLen == 0 {
		d.mapEnd()
		return
	}

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := reflect.Kind(ti.elemkind)
	ktypeKind := reflect.Kind(ti.keykind)
	mparams := getMapReqParams(ti)

	vtypePtr := vtypeKind == reflect.Ptr
	ktypePtr := ktypeKind == reflect.Ptr

	vTransient := decUseTransient && !vtypePtr && ti.tielem.flagCanTransient

	kTransient := vTransient && !ktypePtr && ti.tikey.flagCanTransient

	var vtypeElem reflect.Type

	var keyFn, valFn *decFnSimpleBytes
	var ktypeLo, vtypeLo = ktype, vtype

	if ktypeKind == reflect.Ptr {
		for ktypeLo = ktype.Elem(); ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() {
		}
	}

	if vtypePtr {
		vtypeElem = vtype.Elem()
		for vtypeLo = vtypeElem; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() {
		}
	}

	rvkMut := !scalarBitset.isset(ti.keykind)
	rvvMut := !scalarBitset.isset(ti.elemkind)
	rvvCanNil := isnilBitset.isset(ti.elemkind)

	var rvk, rvkn, rvv, rvvn, rvva, rvvz reflect.Value

	var doMapGet, doMapSet bool

	if !d.h.MapValueReset {
		if rvvMut && (vtypeKind != reflect.Interface || !d.h.InterfaceReset) {
			doMapGet = true
			rvva = mapAddrLoopvarRV(vtype, vtypeKind)
		}
	}

	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId
	hasLen := containerLen >= 0

	var kstr2bs []byte
	var kstr string

	var mapKeyStringSharesBytesBuf bool
	var att dBytesAttachState

	var vElem, kElem reflect.Type
	kbuiltin := ti.tikey.flagDecBuiltin && ti.keykind != uint8(reflect.Slice)
	vbuiltin := ti.tielem.flagDecBuiltin
	if kbuiltin && ktypePtr {
		kElem = ti.key.Elem()
	}
	if vbuiltin && vtypePtr {
		vElem = ti.elem.Elem()
	}

	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		mapKeyStringSharesBytesBuf = false
		kstr = ""
		if j == 0 {

			if kTransient {
				rvk = d.perType.TransientAddr2K(ktype, ktypeKind)
			} else {
				rvk = rvZeroAddrK(ktype, ktypeKind)
			}
			if !rvkMut {
				rvkn = rvk
			}
			if !rvvMut {
				if vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
			}
			if !ktypeIsString && keyFn == nil {
				keyFn = d.fn(ktypeLo)
			}
			if valFn == nil {
				valFn = d.fn(vtypeLo)
			}
		} else if rvkMut {
			rvSetZero(rvk)
		} else {
			rvk = rvkn
		}

		d.mapElemKey(j == 0)

		if d.d.TryNil() {
			rvSetZero(rvk)
		} else if ktypeIsString {
			kstr2bs, att = d.d.DecodeStringAsBytes()
			kstr, mapKeyStringSharesBytesBuf = d.bytes2Str(kstr2bs, att)
			rvSetString(rvk, kstr)
		} else {
			if kbuiltin {
				if ktypePtr {
					if rvIsNil(rvk) {
						rvSetDirect(rvk, reflect.New(kElem))
					}
					d.decode(rv2i(rvk))
				} else {
					d.decode(rv2i(rvAddr(rvk, ti.tikey.ptr)))
				}
			} else {
				d.decodeValueNoCheckNil(rvk, keyFn)
			}

			if ktypeIsIntf {
				if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp {
					kstr2bs = rvGetBytes(rvk2)
					kstr, mapKeyStringSharesBytesBuf = d.bytes2Str(kstr2bs, dBytesAttachView)
					rvSetIntf(rvk, rv4istr(kstr))
				}

			}
		}

		if mapKeyStringSharesBytesBuf && d.bufio {
			if ktypeIsString {
				rvSetString(rvk, d.detach2Str(kstr2bs, att))
			} else {
				rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
			}
			mapKeyStringSharesBytesBuf = false
		}

		d.mapElemValue()

		if d.d.TryNil() {
			if mapKeyStringSharesBytesBuf {
				if ktypeIsString {
					rvSetString(rvk, d.detach2Str(kstr2bs, att))
				} else {
					rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
				}
			}

			if !rvvz.IsValid() {
				rvvz = rvZeroK(vtype, vtypeKind)
			}
			mapSet(rv, rvk, rvvz, mparams)
			continue
		}

		doMapSet = true

		if !rvvMut {
			rvv = rvvn
		} else if !doMapGet {
			goto NEW_RVV
		} else {
			rvv = mapGet(rv, rvk, rvva, mparams)
			if !rvv.IsValid() || (rvvCanNil && rvIsNil(rvv)) {
				goto NEW_RVV
			}
			switch vtypeKind {
			case reflect.Ptr, reflect.Map:
				doMapSet = false
			case reflect.Interface:

				rvvn = rvv.Elem()
				if k := rvvn.Kind(); (k == reflect.Ptr || k == reflect.Map) && !rvIsNil(rvvn) {
					d.decodeValueNoCheckNil(rvvn, nil)
					continue
				}

				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvSetIntf(rvvn, rvv)
				rvv = rvvn
			default:

				if vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
				rvSetDirect(rvvn, rvv)
				rvv = rvvn
			}
		}
		goto DECODE_VALUE_NO_CHECK_NIL

	NEW_RVV:
		if vtypePtr {
			rvv = reflect.New(vtypeElem)
		} else if vTransient {
			rvv = d.perType.TransientAddrK(vtype, vtypeKind)
		} else {
			rvv = rvZeroAddrK(vtype, vtypeKind)
		}

	DECODE_VALUE_NO_CHECK_NIL:
		if doMapSet && mapKeyStringSharesBytesBuf {
			if ktypeIsString {
				rvSetString(rvk, d.detach2Str(kstr2bs, att))
			} else {
				rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
			}
		}
		if vbuiltin {
			if vtypePtr {
				if rvIsNil(rvv) {
					rvSetDirect(rvv, reflect.New(vElem))
				}
				d.decode(rv2i(rvv))
			} else {
				d.decode(rv2i(rvAddr(rvv, ti.tielem.ptr)))
			}
		} else {
			d.decodeValueNoCheckNil(rvv, valFn)
		}
		if doMapSet {
			mapSet(rv, rvk, rvv, mparams)
		}
	}

	d.mapEnd()
}

func (d *decoderSimpleBytes) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsSimpleBytes)

	if d.bytes {
		d.rtidFn = &d.h.rtidFnsDecBytes
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtBytes
	} else {
		d.bufio = d.h.ReaderBufferSize > 0
		d.rtidFn = &d.h.rtidFnsDecIO
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtIO
	}

	d.reset()

}

func (d *decoderSimpleBytes) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
	d.depth = 0
	d.calls = 0

	d.maxdepth = decDefMaxDepth
	if d.h.MaxDepth > 0 {
		d.maxdepth = d.h.MaxDepth
	}
	d.mtid = 0
	d.stid = 0
	d.mtr = false
	d.str = false
	if d.h.MapType != nil {
		d.mtid = rt2id(d.h.MapType)
		_, d.mtr = fastpathAvIndex(d.mtid)
	}
	if d.h.SliceType != nil {
		d.stid = rt2id(d.h.SliceType)
		_, d.str = fastpathAvIndex(d.stid)
	}
}

func (d *decoderSimpleBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderSimpleBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderSimpleBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderSimpleBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderSimpleBytes) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderSimpleBytes) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderSimpleBytes) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderSimpleBytes) Release() {}

func (d *decoderSimpleBytes) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderSimpleBytes) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderSimpleBytes) decode(iv interface{}) {
	_ = d.d

	rv, ok := isNil(iv, true)
	if ok {
		halt.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {

	case *string:
		*v = d.detach2Str(d.d.DecodeStringAsBytes())
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *uintptr:
		*v = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *float32:
		*v = d.d.DecodeFloat32()
	case *float64:
		*v = d.d.DecodeFloat64()
	case *complex64:
		*v = complex(d.d.DecodeFloat32(), 0)
	case *complex128:
		*v = complex(d.d.DecodeFloat64(), 0)
	case *[]byte:
		*v, _ = d.decodeBytesInto(*v, false)
	case []byte:

		d.decodeBytesInto(v[:len(v):len(v)], true)
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface{}:
		d.decodeValue(rv4iptr(v), nil)

	case reflect.Value:
		if ok, _ = isDecodeable(v); !ok {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)

	default:

		if skipFastpathTypeSwitchInDirectCall || !d.dh.fastpathDecodeTypeSwitch(iv, d) {
			if !rv.IsValid() {
				rv = reflect.ValueOf(iv)
			}
			if ok, _ = isDecodeable(rv); !ok {
				d.haltAsNotDecodeable(rv)
			}
			d.decodeValue(rv, nil)
		}
	}
}

func (d *decoderSimpleBytes) decodeValue(rv reflect.Value, fn *decFnSimpleBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderSimpleBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnSimpleBytes) {

	var rvp reflect.Value
	var rvpValid bool
PTR:
	if rv.Kind() == reflect.Ptr {
		rvpValid = true
		if rvIsNil(rv) {
			rvSetDirect(rv, reflect.New(rv.Type().Elem()))
		}
		rvp = rv
		rv = rv.Elem()
		goto PTR
	}

	if fn == nil {
		fn = d.fn(rv.Type())
	}
	if fn.i.addrD {
		if rvpValid {
			rv = rvp
		} else if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else if fn.i.addrDf {
			halt.errorStr("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *decoderSimpleBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderSimpleBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderSimpleBytes) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
	v, att := d.d.DecodeBytes()
	if cap(v) == 0 || (att >= dBytesAttachViewZerocopy && !mustFit) {

		return
	}
	if len(v) == 0 {
		v = zeroByteSlice
		return
	}
	if len(out) == len(v) {
		state = dBytesIntoParamOut
	} else if cap(out) >= len(v) {
		out = out[:len(v)]
		state = dBytesIntoParamOutSlice
	} else if mustFit {
		halt.errorf("bytes capacity insufficient for decoded bytes: got/expected: %d/%d", len(v), len(out))
	} else {
		out = make([]byte, len(v))
		state = dBytesIntoNew
	}
	copy(out, v)
	v = out
	return
}

func (d *decoderSimpleBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderSimpleBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderSimpleBytes) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderSimpleBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderSimpleBytes) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderSimpleBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderSimpleBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleBytes) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderSimpleBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderSimpleBytes) fn(t reflect.Type) *decFnSimpleBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderSimpleBytes) fnNoExt(t reflect.Type) *decFnSimpleBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverSimpleBytes) newDecoderBytes(in []byte, h Handle) *decoderSimpleBytes {
	var c1 decoderSimpleBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverSimpleBytes) newDecoderIO(in io.Reader, h Handle) *decoderSimpleBytes {
	var c1 decoderSimpleBytes
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverSimpleBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsSimpleBytes) (f *fastpathDSimpleBytes, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

func (helperDecDriverSimpleBytes) decFindRtidFn(s []decRtidFnSimpleBytes, rtid uintptr) (i uint, fn *decFnSimpleBytes) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperDecDriverSimpleBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnSimpleBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnSimpleBytes](v))
	}
	return
}

func (dh helperDecDriverSimpleBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsSimpleBytes,
	checkExt bool) (fn *decFnSimpleBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverSimpleBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnSimpleBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverSimpleBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnSimpleBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnSimpleBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnSimpleBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnSimpleBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverSimpleBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsSimpleBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleBytes) {
	fn = new(decFnSimpleBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderSimpleBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderSimpleBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderSimpleBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderSimpleBytes).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderSimpleBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderSimpleBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderSimpleBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderSimpleBytes).textUnmarshal
		fi.addrD = ti.flagTextUnmarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fd = fp[idx].decfn
					fi.addrD = true
					fi.addrDf = false
					if rk == reflect.Array {
						fi.addrD = false
					}
				}
			} else {

				xfe, xrt := dh.decFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf2 := xfe.decfn
					if rk == reflect.Array {
						fi.addrD = false
						fn.fd = func(d *decoderSimpleBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderSimpleBytes, xf *decFnInfo, xrv reflect.Value) {
							if xrv.Kind() == reflect.Ptr {
								xfnf2(d, xf, rvConvert(xrv, xptr2rt))
							} else {
								xfnf2(d, xf, rvConvert(xrv, xrt))
							}
						}
					}
				}
			}
		}
		if fn.fd == nil {
			switch rk {
			case reflect.Bool:
				fn.fd = (*decoderSimpleBytes).kBool
			case reflect.String:
				fn.fd = (*decoderSimpleBytes).kString
			case reflect.Int:
				fn.fd = (*decoderSimpleBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderSimpleBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderSimpleBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderSimpleBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderSimpleBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderSimpleBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderSimpleBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderSimpleBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderSimpleBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderSimpleBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderSimpleBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderSimpleBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderSimpleBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderSimpleBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderSimpleBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderSimpleBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderSimpleBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderSimpleBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderSimpleBytes).kStructSimple
				} else {
					fn.fd = (*decoderSimpleBytes).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderSimpleBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderSimpleBytes).kInterface
			default:

				fn.fd = (*decoderSimpleBytes).kErr
			}
		}
	}
	return
}
func (e *simpleEncDriverBytes) EncodeNil() {
	e.w.writen1(simpleVdNil)
}

func (e *simpleEncDriverBytes) EncodeBool(b bool) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && !b {
		e.EncodeNil()
		return
	}
	if b {
		e.w.writen1(simpleVdTrue)
	} else {
		e.w.writen1(simpleVdFalse)
	}
}

func (e *simpleEncDriverBytes) EncodeFloat32(f float32) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat32)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *simpleEncDriverBytes) EncodeFloat64(f float64) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *simpleEncDriverBytes) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), simpleVdNegInt)
	} else {
		e.encUint(uint64(v), simpleVdPosInt)
	}
}

func (e *simpleEncDriverBytes) EncodeUint(v uint64) {
	e.encUint(v, simpleVdPosInt)
}

func (e *simpleEncDriverBytes) encUint(v uint64, bd uint8) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == 0 {
		e.EncodeNil()
		return
	}
	if v <= math.MaxUint8 {
		e.w.writen2(bd, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 1)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 2)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 3)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *simpleEncDriverBytes) encLen(bd byte, length int) {
	if length == 0 {
		e.w.writen1(bd)
	} else if length <= math.MaxUint8 {
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	} else if length <= math.MaxUint16 {
		e.w.writen1(bd + 2)
		e.w.writen2(bigen.PutUint16(uint16(length)))
	} else if int64(length) <= math.MaxUint32 {
		e.w.writen1(bd + 3)
		e.w.writen4(bigen.PutUint32(uint32(length)))
	} else {
		e.w.writen1(bd + 4)
		e.w.writen8(bigen.PutUint64(uint64(length)))
	}
}

func (e *simpleEncDriverBytes) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	var bs0, bs []byte
	if ext == SelfExt {
		bs0 = e.e.blist.get(1024)
		bs = bs0
		sideEncode(e.h, &e.h.sideEncPool, func(se encoderI) { oneOffEncode(se, v, &bs, basetype, false) })
	} else {
		bs = ext.WriteExt(v)
	}
	if bs == nil {
		e.writeNilBytes()
		goto END
	}
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *simpleEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *simpleEncDriverBytes) encodeExtPreamble(xtag byte, length int) {
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
}

func (e *simpleEncDriverBytes) WriteArrayStart(length int) {
	e.encLen(simpleVdArray, length)
}

func (e *simpleEncDriverBytes) WriteMapStart(length int) {
	e.encLen(simpleVdMap, length)
}

func (e *simpleEncDriverBytes) WriteArrayEmpty() {

	e.w.writen1(simpleVdArray)
}

func (e *simpleEncDriverBytes) WriteMapEmpty() {

	e.w.writen1(simpleVdMap)
}

func (e *simpleEncDriverBytes) EncodeString(v string) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == "" {
		e.EncodeNil()
		return
	}
	if e.h.StringToRaw {
		e.encLen(simpleVdByteArray, len(v))
	} else {
		e.encLen(simpleVdString, len(v))
	}
	e.w.writestr(v)
}

func (e *simpleEncDriverBytes) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *simpleEncDriverBytes) EncodeStringBytesRaw(v []byte) {

	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
}

func (e *simpleEncDriverBytes) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *simpleEncDriverBytes) encodeNilBytes() {
	b := byte(simpleVdNil)
	if e.h.NilCollectionToZeroLength {
		b = simpleVdArray
	}
	e.w.writen1(b)
}

func (e *simpleEncDriverBytes) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = simpleVdNil
	}
	e.w.writen1(v)
}

func (e *simpleEncDriverBytes) writeNilArray() {
	e.writeNilOr(simpleVdArray)
}

func (e *simpleEncDriverBytes) writeNilMap() {
	e.writeNilOr(simpleVdMap)
}

func (e *simpleEncDriverBytes) writeNilBytes() {
	e.writeNilOr(simpleVdByteArray)
}

func (e *simpleEncDriverBytes) EncodeTime(t time.Time) {

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	v, err := t.MarshalBinary()
	halt.onerror(err)
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
}

func (d *simpleDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *simpleDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *simpleDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	switch d.bd {
	case simpleVdNil:
		d.bdRead = false
		return valueTypeNil
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		return valueTypeBytes
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		return valueTypeString
	case simpleVdArray, simpleVdArray + 1,
		simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		return valueTypeArray
	case simpleVdMap, simpleVdMap + 1,
		simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *simpleDecDriverBytes) TryNil() bool {
	return d.advanceNil()
}

func (d *simpleDecDriverBytes) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case simpleVdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *simpleDecDriverBytes) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case simpleVdPosInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case simpleVdPosInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
		neg = true
	default:
		ok = false

	}

	return
}

func (d *simpleDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == simpleVdFalse {
	} else if d.bd == simpleVdTrue {
		b = true
	} else {
		halt.errorf("cannot decode bool - %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverBytes) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverBytes) uint2Len(ui uint64) int {
	if chkOvf.Uint(ui, intBitsize) {
		halt.errorf("overflow integer: %v", ui)
	}
	return int(ui)
}

func (d *simpleDecDriverBytes) decLen() int {
	switch d.bd & 7 {
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(bigen.Uint16(d.r.readn2()))
	case 3:
		return d.uint2Len(uint64(bigen.Uint32(d.r.readn4())))
	case 4:
		return d.uint2Len(bigen.Uint64(d.r.readn8()))
	}
	halt.errorf("cannot read length: bd%%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
}

func (d *simpleDecDriverBytes) DecodeStringAsBytes() ([]byte, dBytesAttachState) {
	return d.DecodeBytes()
}

func (d *simpleDecDriverBytes) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool

	if d.bd >= simpleVdArray && d.bd <= simpleVdArray+4 {
		slen := d.ReadArrayStart()
		bs, cond = usableByteSlice(d.d.buf, slen)
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		if cond {
			d.d.buf = bs
		}
		state = dBytesAttachBuffer

		return
	}

	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *simpleDecDriverBytes) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd != simpleVdTime {
		halt.errorf("invalid descriptor for time.Time - expect 0x%x, received 0x%x", simpleVdTime, d.bd)
	}
	d.bdRead = false
	clen := uint(d.r.readn1())
	b := d.r.readx(clen)
	halt.onerror((&t).UnmarshalBinary(b))
	return
}

func (d *simpleDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	xbs, _, _, ok := d.decodeExtV(ext != nil, xtag)
	if !ok {
		return
	}
	if ext == SelfExt {
		sideDecode(d.h, &d.h.sideDecPool, func(sd decoderI) { oneOffDecode(sd, rv, xbs, basetype, false) })
	} else {
		ext.ReadExt(rv, xbs)
	}
}

func (d *simpleDecDriverBytes) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *simpleDecDriverBytes) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
	if xtagIn > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtagIn)
	}
	if d.advanceNil() {
		return
	}
	tag := uint8(xtagIn)
	switch d.bd {
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		xbs, ok = d.r.readxb(uint(l))
		bstate = d.d.attachState(ok)
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs, bstate = d.DecodeBytes()
	default:
		halt.errorf("ext - %s - expecting extensions/bytearray, got: 0x%x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	ok = true
	return
}

func (d *simpleDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd {
	case simpleVdNil:
		n.v = valueTypeNil
	case simpleVdFalse:
		n.v = valueTypeBool
		n.b = false
	case simpleVdTrue:
		n.v = valueTypeBool
		n.b = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case simpleVdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdTime:
		n.v = valueTypeTime
		n.t = d.DecodeTime()
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(uint(l))

	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2,
		simpleVdArray + 3, simpleVdArray + 4:
		n.v = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, d.bd)
	}

	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *simpleDecDriverBytes) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *simpleDecDriverBytes) nextValueBytesBdReadR() {
	c := d.bd

	var length uint

	switch c {
	case simpleVdNil, simpleVdFalse, simpleVdTrue, simpleVdString, simpleVdByteArray:

	case simpleVdPosInt, simpleVdNegInt:
		d.r.readn1()
	case simpleVdPosInt + 1, simpleVdNegInt + 1:
		d.r.skip(2)
	case simpleVdPosInt + 2, simpleVdNegInt + 2, simpleVdFloat32:
		d.r.skip(4)
	case simpleVdPosInt + 3, simpleVdNegInt + 3, simpleVdFloat64:
		d.r.skip(8)
	case simpleVdTime:
		c = d.r.readn1()
		d.r.skip(uint(c))

	default:
		switch c & 7 {
		case 0:
			length = 0
		case 1:
			b := d.r.readn1()
			length = uint(b)
		case 2:
			x := d.r.readn2()
			length = uint(bigen.Uint16(x))
		case 3:
			x := d.r.readn4()
			length = uint(bigen.Uint32(x))
		case 4:
			x := d.r.readn8()
			length = uint(bigen.Uint64(x))
		}

		bExt := c >= simpleVdExt && c <= simpleVdExt+7
		bStr := c >= simpleVdString && c <= simpleVdString+7
		bByteArray := c >= simpleVdByteArray && c <= simpleVdByteArray+7
		bArray := c >= simpleVdArray && c <= simpleVdArray+7
		bMap := c >= simpleVdMap && c <= simpleVdMap+7

		if !(bExt || bStr || bByteArray || bArray || bMap) {
			halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, c)
		}

		if bExt {
			d.r.readn1()
		}

		if length == 0 {
			break
		}

		if bArray {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else if bMap {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			d.r.skip(length)
		}
	}
	return
}

func (d *simpleEncDriverBytes) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*SimpleHandle)
	d.e = shared
	if shared.bytes {
		fp = simpleFpEncBytes
	} else {
		fp = simpleFpEncIO
	}

	d.init2(enc)
	return
}

func (e *simpleEncDriverBytes) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *simpleEncDriverBytes) writerEnd() { e.w.end() }

func (e *simpleEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *simpleEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *simpleDecDriverBytes) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*SimpleHandle)
	d.d = shared
	if shared.bytes {
		fp = simpleFpDecBytes
	} else {
		fp = simpleFpDecIO
	}

	d.init2(dec)
	return
}

func (d *simpleDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *simpleDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *simpleDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *simpleDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, simpledesc(d.bd))
}

func (d *simpleDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

type helperEncDriverSimpleIO struct{}
type encFnSimpleIO struct {
	i  encFnInfo
	fe func(*encoderSimpleIO, *encFnInfo, reflect.Value)
}
type encRtidFnSimpleIO struct {
	rtid uintptr
	fn   *encFnSimpleIO
}
type encoderSimpleIO struct {
	dh helperEncDriverSimpleIO
	fp *fastpathEsSimpleIO
	e  simpleEncDriverIO
	encoderBase
}
type helperDecDriverSimpleIO struct{}
type decFnSimpleIO struct {
	i  decFnInfo
	fd func(*decoderSimpleIO, *decFnInfo, reflect.Value)
}
type decRtidFnSimpleIO struct {
	rtid uintptr
	fn   *decFnSimpleIO
}
type decoderSimpleIO struct {
	dh helperDecDriverSimpleIO
	fp *fastpathDsSimpleIO
	d  simpleDecDriverIO
	decoderBase
}
type simpleEncDriverIO struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *SimpleHandle
	e *encoderBase

	w bufioEncWriter
}
type simpleDecDriverIO struct {
	h *SimpleHandle
	d *decoderBase
	r ioDecReader

	bdAndBdread

	noBuiltInTypes

	decDriverNoopContainerReader
	decInit2er
}

func (e *encoderSimpleIO) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderSimpleIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderSimpleIO) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderSimpleIO) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderSimpleIO) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderSimpleIO) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderSimpleIO) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderSimpleIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderSimpleIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderSimpleIO) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderSimpleIO) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderSimpleIO) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderSimpleIO) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderSimpleIO) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderSimpleIO) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderSimpleIO) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderSimpleIO) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderSimpleIO) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderSimpleIO) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderSimpleIO) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderSimpleIO) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderSimpleIO) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderSimpleIO) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderSimpleIO) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderSimpleIO) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderSimpleIO) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderSimpleIO) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderSimpleIO) kSeqFn(rt reflect.Type) (fn *encFnSimpleIO) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderSimpleIO) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
	var l int
	if isSlice {
		l = rvLenSlice(rv)
	} else {
		l = rv.Len()
	}
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.haltOnMbsOddLen(l)
	e.mapStart(l >> 1)

	var fn *encFnSimpleIO
	builtin := ti.tielem.flagEncBuiltin
	if !builtin {
		fn = e.kSeqFn(ti.elem)
	}

	j := 0
	e.c = containerMapKey
	e.e.WriteMapElemKey(true)
	for {
		rvv := rvArrayIndex(rv, j, ti, isSlice)
		if builtin {
			e.encodeIB(rv2i(baseRVRV(rvv)))
		} else {
			e.encodeValue(rvv, fn)
		}
		j++
		if j == l {
			break
		}
		if j&1 == 0 {
			e.c = containerMapKey
			e.e.WriteMapElemKey(false)
		} else {
			e.mapElemValue()
		}
	}
	e.c = 0
	e.e.WriteMapEnd()

}

func (e *encoderSimpleIO) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
	var l int
	if isSlice {
		l = rvLenSlice(rv)
	} else {
		l = rv.Len()
	}
	if l <= 0 {
		e.e.WriteArrayEmpty()
		return
	}
	e.arrayStart(l)

	var fn *encFnSimpleIO
	if !ti.tielem.flagEncBuiltin {
		fn = e.kSeqFn(ti.elem)
	}

	j := 0
	e.c = containerArrayElem
	e.e.WriteArrayElem(true)
	builtin := ti.tielem.flagEncBuiltin
	for {
		rvv := rvArrayIndex(rv, j, ti, isSlice)
		if builtin {
			e.encodeIB(rv2i(baseRVRV(rvv)))
		} else {
			e.encodeValue(rvv, fn)
		}
		j++
		if j == l {
			break
		}
		e.c = containerArrayElem
		e.e.WriteArrayElem(false)
	}

	e.c = 0
	e.e.WriteArrayEnd()
}

func (e *encoderSimpleIO) kChan(f *encFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		halt.errorStr("send-only channel cannot be encoded")
	}
	if !f.ti.mbs && uint8TypId == rt2id(f.ti.elem) {
		e.kSliceBytesChan(rv)
		return
	}
	rtslice := reflect.SliceOf(f.ti.elem)
	rv = chanToSlice(rv, rtslice, e.h.ChanRecvTimeout)
	ti := e.h.getTypeInfo(rt2id(rtslice), rtslice)
	if f.ti.mbs {
		e.kArrayWMbs(rv, ti, true)
	} else {
		e.kArrayW(rv, ti, true)
	}
}

func (e *encoderSimpleIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderSimpleIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderSimpleIO) kSliceBytesChan(rv reflect.Value) {

	bs0 := e.blist.peek(32, true)
	bs := bs0

	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok {
		ch = irv.(chan byte)
	}

L1:
	switch timeout := e.h.ChanRecvTimeout; {
	case timeout == 0:
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			}
		}
	case timeout > 0:
		tt := time.NewTimer(timeout)
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:

				break L1
			}
		}
	default:
		for b := range ch {
			bs = append(bs, b)
		}
	}

	e.e.EncodeBytes(bs)
	e.blist.put(bs)
	if !byteSliceSameData(bs0, bs) {
		e.blist.put(bs0)
	}
}

func (e *encoderSimpleIO) kStructFieldKey(keyType valueType, encName string) {

	if keyType == valueTypeString {
		e.e.EncodeString(encName)
	} else if keyType == valueTypeInt {
		e.e.EncodeInt(must.Int(strconv.ParseInt(encName, 10, 64)))
	} else if keyType == valueTypeUint {
		e.e.EncodeUint(must.Uint(strconv.ParseUint(encName, 10, 64)))
	} else if keyType == valueTypeFloat {
		e.e.EncodeFloat64(must.Float(strconv.ParseFloat(encName, 64)))
	} else {
		halt.errorStr2("invalid struct key type: ", keyType.String())
	}

}

func (e *encoderSimpleIO) kStructSimple(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	tisfi := f.ti.sfi.source()

	chkCirRef := e.h.CheckCircularRef
	var si *structFieldInfo
	var j int

	if f.ti.toArray || e.h.StructToArray {
		if len(tisfi) == 0 {
			e.e.WriteArrayEmpty()
			return
		}
		e.arrayStart(len(tisfi))
		for j, si = range tisfi {
			e.c = containerArrayElem
			e.e.WriteArrayElem(j == 0)
			if si.encBuiltin {
				e.encodeIB(rv2i(si.fieldNoAlloc(rv, true)))
			} else {
				e.encodeValue(si.fieldNoAlloc(rv, !chkCirRef), nil)
			}
		}
		e.c = 0
		e.e.WriteArrayEnd()
	} else {
		if len(tisfi) == 0 {
			e.e.WriteMapEmpty()
			return
		}
		if e.h.Canonical {
			tisfi = f.ti.sfi.sorted()
		}
		e.mapStart(len(tisfi))
		for j, si = range tisfi {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
			e.e.EncodeStringNoEscape4Json(si.encName)
			e.mapElemValue()
			if si.encBuiltin {
				e.encodeIB(rv2i(si.fieldNoAlloc(rv, true)))
			} else {
				e.encodeValue(si.fieldNoAlloc(rv, !chkCirRef), nil)
			}
		}
		e.c = 0
		e.e.WriteMapEnd()
	}
}

func (e *encoderSimpleIO) kStruct(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	ti := f.ti
	toMap := !(ti.toArray || e.h.StructToArray)
	var mf map[string]interface{}
	if ti.flagMissingFielder {
		toMap = true
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
	} else if ti.flagMissingFielderPtr {
		toMap = true
		if rv.CanAddr() {
			mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder).CodecMissingFields()
		} else {
			mf = rv2i(e.addrRV(rv, ti.rt, ti.ptr)).(MissingFielder).CodecMissingFields()
		}
	}
	newlen := len(mf)
	tisfi := ti.sfi.source()
	newlen += len(tisfi)

	var fkvs = e.slist.get(newlen)[:newlen]

	recur := e.h.RecursiveEmptyCheck
	chkCirRef := e.h.CheckCircularRef

	var xlen int

	var kv sfiRv
	var j int
	var sf encStructFieldObj
	if toMap {
		newlen = 0
		if e.h.Canonical {
			tisfi = f.ti.sfi.sorted()
		}
		for _, si := range tisfi {

			if si.omitEmpty {
				kv.r = si.fieldNoAlloc(rv, false)
				if isEmptyValue(kv.r, e.h.TypeInfos, recur) {
					continue
				}
			} else {
				kv.r = si.fieldNoAlloc(rv, si.encBuiltin || !chkCirRef)
			}
			kv.v = si
			fkvs[newlen] = kv
			newlen++
		}

		var mf2s []stringIntf
		if len(mf) != 0 {
			mf2s = make([]stringIntf, 0, len(mf))
			for k, v := range mf {
				if k == "" {
					continue
				}
				if ti.infoFieldOmitempty && isEmptyValue(reflect.ValueOf(v), e.h.TypeInfos, recur) {
					continue
				}
				mf2s = append(mf2s, stringIntf{k, v})
			}
		}

		xlen = newlen + len(mf2s)
		if xlen == 0 {
			e.e.WriteMapEmpty()
			goto END
		}

		e.mapStart(xlen)

		if len(mf2s) != 0 && e.h.Canonical {
			mf2w := make([]encStructFieldObj, newlen+len(mf2s))
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				mf2w[j] = encStructFieldObj{kv.v.encName, kv.r, nil, true,
					!kv.v.encNameEscape4Json, kv.v.encBuiltin}
			}
			for _, v := range mf2s {
				mf2w[j] = encStructFieldObj{v.v, reflect.Value{}, v.i, false, false, false}
				j++
			}
			sort.Sort((encStructFieldObjSlice)(mf2w))
			for j, sf = range mf2w {
				e.c = containerMapKey
				e.e.WriteMapElemKey(j == 0)
				if ti.keyType == valueTypeString && sf.noEsc4json {
					e.e.EncodeStringNoEscape4Json(sf.key)
				} else {
					e.kStructFieldKey(ti.keyType, sf.key)
				}
				e.mapElemValue()
				if sf.isRv {
					if sf.builtin {
						e.encodeIB(rv2i(baseRVRV(sf.rv)))
					} else {
						e.encodeValue(sf.rv, nil)
					}
				} else {
					if !e.encodeBuiltin(sf.intf) {
						e.encodeR(reflect.ValueOf(sf.intf))
					}

				}
			}
		} else {
			keytyp := ti.keyType
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				e.c = containerMapKey
				e.e.WriteMapElemKey(j == 0)
				if ti.keyType == valueTypeString && !kv.v.encNameEscape4Json {
					e.e.EncodeStringNoEscape4Json(kv.v.encName)
				} else {
					e.kStructFieldKey(keytyp, kv.v.encName)
				}
				e.mapElemValue()
				if kv.v.encBuiltin {
					e.encodeIB(rv2i(baseRVRV(kv.r)))
				} else {
					e.encodeValue(kv.r, nil)
				}
			}
			for _, v := range mf2s {
				e.c = containerMapKey
				e.e.WriteMapElemKey(j == 0)
				e.kStructFieldKey(keytyp, v.v)
				e.mapElemValue()
				if !e.encodeBuiltin(v.i) {
					e.encodeR(reflect.ValueOf(v.i))
				}

				j++
			}
		}

		e.c = 0
		e.e.WriteMapEnd()
	} else {
		newlen = len(tisfi)
		for i, si := range tisfi {

			if si.omitEmpty {

				kv.r = si.fieldNoAlloc(rv, false)
				if isEmptyContainerValue(kv.r, e.h.TypeInfos, recur) {
					kv.r = reflect.Value{}
				}
			} else {
				kv.r = si.fieldNoAlloc(rv, si.encBuiltin || !chkCirRef)
			}
			kv.v = si
			fkvs[i] = kv
		}

		if newlen == 0 {
			e.e.WriteArrayEmpty()
			goto END
		}

		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ {
			e.c = containerArrayElem
			e.e.WriteArrayElem(j == 0)
			kv = fkvs[j]
			if !kv.r.IsValid() {
				e.e.EncodeNil()
			} else if kv.v.encBuiltin {
				e.encodeIB(rv2i(baseRVRV(kv.r)))
			} else {
				e.encodeValue(kv.r, nil)
			}
		}
		e.c = 0
		e.e.WriteArrayEnd()
	}

END:

	e.slist.put(fkvs)
}

func (e *encoderSimpleIO) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnSimpleIO

	ktypeKind := reflect.Kind(f.ti.keykind)
	vtypeKind := reflect.Kind(f.ti.elemkind)

	rtval := f.ti.elem
	rtvalkind := vtypeKind
	for rtvalkind == reflect.Ptr {
		rtval = rtval.Elem()
		rtvalkind = rtval.Kind()
	}
	if rtvalkind != reflect.Interface {
		valFn = e.fn(rtval)
	}

	var rvv = mapAddrLoopvarRV(f.ti.elem, vtypeKind)

	rtkey := f.ti.key
	var keyTypeIsString = stringTypId == rt2id(rtkey)
	if keyTypeIsString {
		keyFn = e.fn(rtkey)
	} else {
		for rtkey.Kind() == reflect.Ptr {
			rtkey = rtkey.Elem()
		}
		if rtkey.Kind() != reflect.Interface {
			keyFn = e.fn(rtkey)
		}
	}

	if e.h.Canonical {
		e.kMapCanonical(f.ti, rv, rvv, keyFn, valFn)
		e.c = 0
		e.e.WriteMapEnd()
		return
	}

	var rvk = mapAddrLoopvarRV(f.ti.key, ktypeKind)

	var it mapIter
	mapRange(&it, rv, rvk, rvv, true)

	kbuiltin := f.ti.tikey.flagEncBuiltin
	vbuiltin := f.ti.tielem.flagEncBuiltin
	for j := 0; it.Next(); j++ {
		rv = it.Key()
		e.c = containerMapKey
		e.e.WriteMapElemKey(j == 0)
		if keyTypeIsString {
			e.e.EncodeString(rvGetString(rv))
		} else if kbuiltin {
			e.encodeIB(rv2i(baseRVRV(rv)))
		} else {
			e.encodeValue(rv, keyFn)
		}
		e.mapElemValue()
		rv = it.Value()
		if vbuiltin {
			e.encodeIB(rv2i(baseRVRV(rv)))
		} else {
			e.encodeValue(it.Value(), valFn)
		}
	}
	it.Done()

	e.c = 0
	e.e.WriteMapEnd()
}

func (e *encoderSimpleIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnSimpleIO) {
	_ = e.e

	rtkey := ti.key
	rtkeydecl := rtkey.PkgPath() == "" && rtkey.Name() != ""

	mks := rv.MapKeys()
	rtkeyKind := rtkey.Kind()
	mparams := getMapReqParams(ti)

	switch rtkeyKind {
	case reflect.Bool:

		if len(mks) == 2 && mks[0].Bool() {
			mks[0], mks[1] = mks[1], mks[0]
		}
		for i := range mks {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeBool(mks[i].Bool())
			} else {
				e.encodeValueNonNil(mks[i], keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mks[i], rvv, mparams), valFn)
		}
	case reflect.String:
		mksv := make([]orderedRv[string], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = rvGetString(k)
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeString(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		mksv := make([]orderedRv[uint64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Uint()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeUint(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		mksv := make([]orderedRv[int64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Int()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeInt(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Float32:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeFloat32(float32(mksv[i].v))
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	case reflect.Float64:
		mksv := make([]orderedRv[float64], len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		slices.SortFunc(mksv, cmpOrderedRv)
		for i := range mksv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(i == 0)
			if rtkeydecl {
				e.e.EncodeFloat64(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
		}
	default:
		if rtkey == timeTyp {
			mksv := make([]timeRv, len(mks))
			for i, k := range mks {
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			}
			slices.SortFunc(mksv, cmpTimeRv)
			for i := range mksv {
				e.c = containerMapKey
				e.e.WriteMapElemKey(i == 0)
				e.e.EncodeTime(mksv[i].v)
				e.mapElemValue()
				e.encodeValue(mapGet(rv, mksv[i].r, rvv, mparams), valFn)
			}
			break
		}

		bs0 := e.blist.get(len(mks) * 16)
		mksv := bs0
		mksbv := make([]bytesRv, len(mks))

		sideEncode(e.hh, &e.h.sideEncPool, func(se encoderI) {
			se.ResetBytes(&mksv)
			for i, k := range mks {
				v := &mksbv[i]
				l := len(mksv)
				se.setContainerState(containerMapKey)
				se.encodeR(baseRVRV(k))
				se.atEndOfEncode()
				se.writerEnd()
				v.r = k
				v.v = mksv[l:]
			}
		})

		slices.SortFunc(mksbv, cmpBytesRv)
		for j := range mksbv {
			e.c = containerMapKey
			e.e.WriteMapElemKey(j == 0)
			e.e.writeBytesAsis(mksbv[j].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksbv[j].r, rvv, mparams), valFn)
		}
		e.blist.put(mksv)
		if !byteSliceSameData(bs0, mksv) {
			e.blist.put(bs0)
		}
	}
}

func (e *encoderSimpleIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsSimpleIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderSimpleIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderSimpleIO) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderSimpleIO) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderSimpleIO) mustEncode(v interface{}) {
	halt.onerror(e.err)
	if e.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	e.calls++
	if !e.encodeBuiltin(v) {
		e.encodeR(reflect.ValueOf(v))
	}

	e.calls--
	if e.calls == 0 {
		e.e.atEndOfEncode()
		e.e.writerEnd()
	}
}

func (e *encoderSimpleIO) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderSimpleIO) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderSimpleIO) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderSimpleIO) encodeBuiltin(iv interface{}) (ok bool) {
	ok = true
	switch v := iv.(type) {
	case nil:
		e.e.EncodeNil()

	case Raw:
		e.rawBytes(v)
	case string:
		e.e.EncodeString(v)
	case bool:
		e.e.EncodeBool(v)
	case int:
		e.e.EncodeInt(int64(v))
	case int8:
		e.e.EncodeInt(int64(v))
	case int16:
		e.e.EncodeInt(int64(v))
	case int32:
		e.e.EncodeInt(int64(v))
	case int64:
		e.e.EncodeInt(v)
	case uint:
		e.e.EncodeUint(uint64(v))
	case uint8:
		e.e.EncodeUint(uint64(v))
	case uint16:
		e.e.EncodeUint(uint64(v))
	case uint32:
		e.e.EncodeUint(uint64(v))
	case uint64:
		e.e.EncodeUint(v)
	case uintptr:
		e.e.EncodeUint(uint64(v))
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)
	case complex64:
		e.encodeComplex64(v)
	case complex128:
		e.encodeComplex128(v)
	case time.Time:
		e.e.EncodeTime(v)
	case []byte:
		e.e.EncodeBytes(v)
	default:

		ok = !skipFastpathTypeSwitchInDirectCall && e.dh.fastpathEncodeTypeSwitch(iv, e)
	}
	return
}

func (e *encoderSimpleIO) encodeValue(rv reflect.Value, fn *encFnSimpleIO) {

	var ciPushes int

	var rvp reflect.Value
	var rvpValid bool

RV:
	switch rv.Kind() {
	case reflect.Ptr:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
		rvpValid = true
		rvp = rv
		rv = rv.Elem()

		if e.h.CheckCircularRef && e.ci.canPushElemKind(rv.Kind()) {
			e.ci.push(rv2i(rvp))
			ciPushes++
		}
		goto RV
	case reflect.Interface:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			goto END
		}
		rvpValid = false
		rvp = reflect.Value{}
		rv = rv.Elem()
		fn = nil
		goto RV
	case reflect.Map:
		if rvIsNil(rv) {
			if e.h.NilCollectionToZeroLength {
				e.e.WriteMapEmpty()
			} else {
				e.e.EncodeNil()
			}
			goto END
		}
	case reflect.Slice, reflect.Chan:
		if rvIsNil(rv) {
			if e.h.NilCollectionToZeroLength {
				e.e.WriteArrayEmpty()
			} else {
				e.e.EncodeNil()
			}
			goto END
		}
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		goto END
	}

	if fn == nil {
		fn = e.fn(rv.Type())
	}

	if !fn.i.addrE {

	} else if rvpValid {
		rv = rvp
	} else if rv.CanAddr() {
		rv = rvAddr(rv, fn.i.ti.ptr)
	} else {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)

END:
	if ciPushes > 0 {
		e.ci.pop(ciPushes)
	}
}

func (e *encoderSimpleIO) encodeValueNonNil(rv reflect.Value, fn *encFnSimpleIO) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderSimpleIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderSimpleIO) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderSimpleIO) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderSimpleIO) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderSimpleIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderSimpleIO) fn(t reflect.Type) *encFnSimpleIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderSimpleIO) fnNoExt(t reflect.Type) *encFnSimpleIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderSimpleIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderSimpleIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderSimpleIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderSimpleIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderSimpleIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderSimpleIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderSimpleIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderSimpleIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverSimpleIO) newEncoderBytes(out *[]byte, h Handle) *encoderSimpleIO {
	var c1 encoderSimpleIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverSimpleIO) newEncoderIO(out io.Writer, h Handle) *encoderSimpleIO {
	var c1 encoderSimpleIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverSimpleIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsSimpleIO) (f *fastpathESimpleIO, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

func (helperEncDriverSimpleIO) encFindRtidFn(s []encRtidFnSimpleIO, rtid uintptr) (i uint, fn *encFnSimpleIO) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperEncDriverSimpleIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnSimpleIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnSimpleIO](v))
	}
	return
}

func (dh helperEncDriverSimpleIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsSimpleIO, checkExt bool) (fn *encFnSimpleIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverSimpleIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnSimpleIO = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverSimpleIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnSimpleIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnSimpleIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnSimpleIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnSimpleIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverSimpleIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnSimpleIO) {
	fn = new(encFnSimpleIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderSimpleIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderSimpleIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderSimpleIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderSimpleIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderSimpleIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderSimpleIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderSimpleIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderSimpleIO).textMarshal
		fi.addrE = ti.flagTextMarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {

			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fe = fp[idx].encfn
				}
			} else {

				xfe, xrt := dh.encFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf := xfe.encfn
					fn.fe = func(e *encoderSimpleIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderSimpleIO).kBool
			case reflect.String:

				fn.fe = (*encoderSimpleIO).kString
			case reflect.Int:
				fn.fe = (*encoderSimpleIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderSimpleIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderSimpleIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderSimpleIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderSimpleIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderSimpleIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderSimpleIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderSimpleIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderSimpleIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderSimpleIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderSimpleIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderSimpleIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderSimpleIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderSimpleIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderSimpleIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderSimpleIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderSimpleIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderSimpleIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderSimpleIO).kStructSimple
				} else {
					fn.fe = (*encoderSimpleIO).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderSimpleIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderSimpleIO).kErr
			default:

				fn.fe = (*encoderSimpleIO).kErr
			}
		}
	}
	return
}
func (d *decoderSimpleIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderSimpleIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderSimpleIO) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderSimpleIO) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderSimpleIO) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderSimpleIO) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderSimpleIO) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderSimpleIO) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderSimpleIO) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderSimpleIO) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderSimpleIO) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderSimpleIO) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderSimpleIO) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderSimpleIO) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderSimpleIO) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderSimpleIO) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderSimpleIO) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderSimpleIO) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderSimpleIO) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderSimpleIO) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderSimpleIO) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderSimpleIO) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleIO) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderSimpleIO) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderSimpleIO) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderSimpleIO) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderSimpleIO) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderSimpleIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

	n := d.naked()
	d.d.DecodeNaked()

	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		halt.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
	}

	switch n.v {
	case valueTypeMap:
		mtid := d.mtid
		if mtid == 0 {
			if d.jsms {
				mtid = mapStrIntfTypId
			} else {
				mtid = mapIntfIntfTypId
			}
		}
		if mtid == mapStrIntfTypId {
			var v2 map[string]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if mtid == mapIntfIntfTypId {
			var v2 map[interface{}]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.mtr {
			rvn = reflect.New(d.h.MapType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {

			rvn = rvZeroAddrK(d.h.MapType, reflect.Map)
			d.decodeValue(rvn, nil)
		}
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId {
			var v2 []interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.str {
			rvn = reflect.New(d.h.SliceType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.SliceType, reflect.Slice)
			d.decodeValue(rvn, nil)
		}
		if d.h.PreferArrayOverSlice {
			rvn = rvGetArray4Slice(rvn)
		}
	case valueTypeExt:
		tag, bytes := n.u, n.l
		bfn := d.h.getExtForTag(tag)
		var re = RawExt{Tag: tag}
		if bytes == nil {

			if bfn == nil {
				d.decode(&re.Value)
				rvn = rv4iptr(&re).Elem()
			} else if bfn.ext == SelfExt {
				rvn = rvZeroAddrK(bfn.rt, bfn.rt.Kind())
				d.decodeValue(rvn, d.fnNoExt(bfn.rt))
			} else {
				rvn = reflect.New(bfn.rt)
				d.interfaceExtConvertAndDecode(rv2i(rvn), bfn.ext)
				rvn = rvn.Elem()
			}
		} else {

			if bfn == nil {
				re.setData(bytes, false)
				rvn = rv4iptr(&re).Elem()
			} else {
				rvn = reflect.New(bfn.rt)
				if bfn.ext == SelfExt {
					sideDecode(d.hh, &d.h.sideDecPool, func(sd decoderI) { oneOffDecode(sd, rv2i(rvn), bytes, bfn.rt, false) })
				} else {
					bfn.ext.ReadExt(rv2i(rvn), bytes)
				}
				rvn = rvn.Elem()
			}
		}

		if d.h.PreferPointerForStructOrArray && rvn.CanAddr() {
			if rk := rvn.Kind(); rk == reflect.Array || rk == reflect.Struct {
				rvn = rvn.Addr()
			}
		}
	case valueTypeNil:

	case valueTypeInt:
		rvn = n.ri()
	case valueTypeUint:
		rvn = n.ru()
	case valueTypeFloat:
		rvn = n.rf()
	case valueTypeBool:
		rvn = n.rb()
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs()
	case valueTypeBytes:
		rvn = n.rl()
	case valueTypeTime:
		rvn = n.rt()
	default:
		halt.errorStr2("kInterfaceNaked: unexpected valueType: ", n.v.String())
	}
	return
}

func (d *decoderSimpleIO) kInterface(f *decFnInfo, rv reflect.Value) {

	isnilrv := rvIsNil(rv)

	var rvn reflect.Value

	if d.h.InterfaceReset {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			} else if !isnilrv {
				decSetNonNilRV2Zero4Intf(rv)
			}
			return
		}
	} else if isnilrv {

		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			}
			return
		}
	} else {

		rvn = rv.Elem()
	}

	canDecode, _ := isDecodeable(rvn)

	if !canDecode {
		rvn2 := d.oneShotAddrRV(rvn.Type(), rvn.Kind())
		rvSetDirect(rvn2, rvn)
		rvn = rvn2
	}

	d.decodeValue(rvn, nil)
	rvSetIntf(rv, rvn)
}

func (d *decoderSimpleIO) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		rv = si.fieldNoAlloc(rv, true)
		if rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
	} else if si.decBuiltin {
		rv = rvAddr(si.fieldAlloc(rv), si.ptrTyp)
		d.decode(rv2i(rv))
	} else {
		fn := d.fn(si.baseTyp)
		rv = si.fieldAlloc(rv)
		if fn.i.addrD {
			rv = rvAddr(rv, si.ptrTyp)
		}
		fn.fd(d, &fn.i, rv)
	}
}

func (d *decoderSimpleIO) kStructSimple(f *decFnInfo, rv reflect.Value) {
	_ = d.d
	ctyp := d.d.ContainerType()
	ti := f.ti
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var rvkencname []byte
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey(j == 0)
			sab, att := d.d.DecodeStringAsBytes()
			rvkencname = d.usableStructFieldNameBytes(rvkencname, sab, att)
			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}

		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem(j == 0)
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}
		d.arrayEnd()
	} else {
		halt.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderSimpleIO) kStruct(f *decFnInfo, rv reflect.Value) {
	_ = d.d
	ctyp := d.d.ContainerType()
	ti := f.ti
	var mf MissingFielder
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder)
	} else if ti.flagMissingFielderPtr {
		mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder)
	}
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var name2 []byte
		var rvkencname []byte
		tkt := ti.keyType
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey(j == 0)

			if tkt == valueTypeString {
				sab, att := d.d.DecodeStringAsBytes()
				rvkencname = d.usableStructFieldNameBytes(rvkencname, sab, att)
			} else if tkt == valueTypeInt {
				rvkencname = strconv.AppendInt(d.b[:0], d.d.DecodeInt64(), 10)
			} else if tkt == valueTypeUint {
				rvkencname = strconv.AppendUint(d.b[:0], d.d.DecodeUint64(), 10)
			} else if tkt == valueTypeFloat {
				rvkencname = strconv.AppendFloat(d.b[:0], d.d.DecodeFloat64(), 'f', -1, 64)
			} else {
				halt.errorStr2("invalid struct key type: ", ti.keyType.String())
			}

			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else if mf != nil {

				name2 = append(name2[:0], rvkencname...)
				var f interface{}
				d.decode(&f)
				if !mf.CodecMissingField(name2, f) && d.h.ErrorIfNoField {
					halt.errorStr2("no matching struct field when decoding stream map with key: ", stringView(name2))
				}
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}

		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem(j == 0)
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}

		d.arrayEnd()
	} else {
		halt.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *decoderSimpleIO) kSlice(f *decFnInfo, rv reflect.Value) {
	_ = d.d

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		rvbs := rvGetBytes(rv)
		if rvCanset {
			bs2, bst := d.decodeBytesInto(rvbs, false)
			if bst != dBytesIntoParamOut {
				rvSetBytes(rv, bs2)
			}
		} else {

			d.decodeBytesInto(rvbs[:len(rvbs):len(rvbs)], true)
		}
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

	if containerLenS == 0 {
		if rvCanset {
			if rvIsNil(rv) {
				rvSetDirect(rv, rvSliceZeroCap(ti.rt))
			} else {
				rvSetSliceLen(rv, 0)
			}
		}
		if isArray {
			d.arrayEnd()
		} else {
			d.mapEnd()
		}
		return
	}

	rtelem0Mut := !scalarBitset.isset(ti.elemkind)
	rtelem := ti.elem

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnSimpleIO

	var rvChanged bool

	var rv0 = rv
	var rv9 reflect.Value

	rvlen := rvLenSlice(rv)
	rvcap := rvCapSlice(rv)
	maxInitLen := d.maxInitLen()
	hasLen := containerLenS >= 0
	if hasLen {
		if containerLenS > rvcap {
			oldRvlenGtZero := rvlen > 0
			rvlen1 := int(decInferLen(containerLenS, maxInitLen, uint(ti.elemsize)))
			if rvlen1 == rvlen {
			} else if rvlen1 <= rvcap {
				if rvCanset {
					rvlen = rvlen1
					rvSetSliceLen(rv, rvlen)
				}
			} else if rvCanset {
				rvlen = rvlen1
				rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = !rvCanset
			} else {
				halt.errorStr("cannot decode into non-settable slice")
			}
			if rvChanged && oldRvlenGtZero && rtelem0Mut {
				rvCopySlice(rv, rv0, rtelem)
			}
		} else if containerLenS != rvlen {
			if rvCanset {
				rvlen = containerLenS
				rvSetSliceLen(rv, rvlen)
			}
		}
	}

	var elemReset = d.h.SliceElementReset

	var rtelemIsPtr bool
	var rtelemElem reflect.Type
	builtin := ti.tielem.flagDecBuiltin
	if builtin {
		rtelemIsPtr = ti.elemkind == uint8(reflect.Ptr)
		if rtelemIsPtr {
			rtelemElem = ti.elem.Elem()
		}
	}

	var j int
	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if rvCanset {
					rvlen = int(decInferLen(containerLenS, maxInitLen, uint(ti.elemsize)))
					rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
					rvcap = rvlen
					rvChanged = !rvCanset
				} else {
					halt.errorStr("cannot decode into non-settable slice")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}

		if ctyp == valueTypeArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}

		if j >= rvlen {

			if rvlen < rvcap {
				rvlen = rvcap
				if rvCanset {
					rvSetSliceLen(rv, rvlen)
				} else if rvChanged {
					rv = rvSlice(rv, rvlen)
				} else {
					halt.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					halt.onerror(errExpandSliceCannotChange)
				}
				rv, rvcap, rvCanset = rvGrowSlice(rv, f.ti, rvcap, 1)

				rvlen = rvcap
				rvChanged = !rvCanset
			}
		}

		rv9 = rvArrayIndex(rv, j, f.ti, true)
		if elemReset {
			rvSetZero(rv9)
		}
		if d.d.TryNil() {
			rvSetZero(rv9)
		} else if builtin {
			if rtelemIsPtr {
				if rvIsNil(rv9) {
					rvSetDirect(rv9, reflect.New(rtelemElem))
				}
				d.decode(rv2i(rv9))
			} else {
				d.decode(rv2i(rvAddr(rv9, ti.tielem.ptr)))
			}
		} else {
			d.decodeValueNoCheckNil(rv9, fn)
		}
	}
	if j < rvlen {
		if rvCanset {
			rvSetSliceLen(rv, j)
		} else if rvChanged {
			rv = rvSlice(rv, j)
		}

	} else if j == 0 && rvIsNil(rv) {
		if rvCanset {
			rv = rvSliceZeroCap(ti.rt)
			rvCanset = false
			rvChanged = true
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}

	if rvChanged {
		rvSetDirect(rv0, rv)
	}
}

func (d *decoderSimpleIO) kArray(f *decFnInfo, rv reflect.Value) {
	_ = d.d

	ti := f.ti
	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {

		if ti.elemkind != uint8(reflect.Uint8) {
			halt.errorf("bytes/string in stream can decode into array of bytes, but not %v", ti.rt)
		}
		rvbs := rvGetArrayBytes(rv, nil)
		d.decodeBytesInto(rvbs, true)
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

	if containerLenS == 0 {
		if isArray {
			d.arrayEnd()
		} else {
			d.mapEnd()
		}
		return
	}

	rtelem := ti.elem
	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var rv9 reflect.Value

	rvlen := rv.Len()
	hasLen := containerLenS >= 0
	if hasLen && containerLenS > rvlen {
		halt.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
	}

	var elemReset = d.h.SliceElementReset

	var rtelemIsPtr bool
	var rtelemElem reflect.Type
	var fn *decFnSimpleIO
	builtin := ti.tielem.flagDecBuiltin
	if builtin {
		rtelemIsPtr = ti.elemkind == uint8(reflect.Ptr)
		if rtelemIsPtr {
			rtelemElem = ti.elem.Elem()
		}
	} else {
		fn = d.fn(rtelem)
	}

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if ctyp == valueTypeArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}

		if j >= rvlen {
			d.arrayCannotExpand(rvlen, j+1)
			d.swallow()
			continue
		}

		rv9 = rvArrayIndex(rv, j, f.ti, false)
		if elemReset {
			rvSetZero(rv9)
		}
		if d.d.TryNil() {
			rvSetZero(rv9)
		} else if builtin {
			if rtelemIsPtr {
				if rvIsNil(rv9) {
					rvSetDirect(rv9, reflect.New(rtelemElem))
				}
				d.decode(rv2i(rv9))
			} else {
				d.decode(rv2i(rvAddr(rv9, ti.tielem.ptr)))
			}
		} else {
			d.decodeValueNoCheckNil(rv9, fn)
		}
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}
}

func (d *decoderSimpleIO) kChan(f *decFnInfo, rv reflect.Value) {
	_ = d.d

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		halt.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {

		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			halt.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2, _ := d.d.DecodeBytes()
		irv := rv2i(rv)
		ch, ok := irv.(chan<- byte)
		if !ok {
			ch = irv.(chan byte)
		}
		for _, b := range bs2 {
			ch <- b
		}
		return
	}

	var rvCanset = rv.CanSet()

	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}

	if containerLenS == 0 {
		if rvCanset && rvIsNil(rv) {
			rvSetDirect(rv, reflect.MakeChan(ti.rt, 0))
		}
		if isArray {
			d.arrayEnd()
		} else {
			d.mapEnd()
		}
		return
	}

	rtelem := ti.elem
	useTransient := decUseTransient && ti.elemkind != byte(reflect.Ptr) && ti.tielem.flagCanTransient

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *decFnSimpleIO

	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int
	hasLen := containerLenS >= 0
	maxInitLen := d.maxInitLen()

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if hasLen {
					rvlen = int(decInferLen(containerLenS, maxInitLen, uint(ti.elemsize)))
				} else {
					rvlen = decDefChanCap
				}
				if rvCanset {
					rv = reflect.MakeChan(ti.rt, rvlen)
					rvChanged = true
				} else {
					halt.errorStr("cannot decode into non-settable chan")
				}
			}
			if fn == nil {
				fn = d.fn(rtelem)
			}
		}

		if ctyp == valueTypeArray {
			d.arrayElem(j == 0)
		} else if j&1 == 0 {
			d.mapElemKey(j == 0)
		} else {
			d.mapElemValue()
		}

		if rv9.IsValid() {
			rvSetZero(rv9)
		} else if useTransient {
			rv9 = d.perType.TransientAddrK(ti.elem, reflect.Kind(ti.elemkind))
		} else {
			rv9 = rvZeroAddrK(ti.elem, reflect.Kind(ti.elemkind))
		}
		if !d.d.TryNil() {
			d.decodeValueNoCheckNil(rv9, fn)
		}
		rv.Send(rv9)
	}
	if isArray {
		d.arrayEnd()
	} else {
		d.mapEnd()
	}

	if rvChanged {
		rvSetDirect(rv0, rv)
	}

}

func (d *decoderSimpleIO) kMap(f *decFnInfo, rv reflect.Value) {
	_ = d.d
	containerLen := d.mapStart(d.d.ReadMapStart())
	ti := f.ti
	if rvIsNil(rv) {
		rvlen := int(decInferLen(containerLen, d.maxInitLen(), uint(ti.keysize+ti.elemsize)))
		rvSetDirect(rv, makeMapReflect(ti.rt, rvlen))
	}

	if containerLen == 0 {
		d.mapEnd()
		return
	}

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := reflect.Kind(ti.elemkind)
	ktypeKind := reflect.Kind(ti.keykind)
	mparams := getMapReqParams(ti)

	vtypePtr := vtypeKind == reflect.Ptr
	ktypePtr := ktypeKind == reflect.Ptr

	vTransient := decUseTransient && !vtypePtr && ti.tielem.flagCanTransient

	kTransient := vTransient && !ktypePtr && ti.tikey.flagCanTransient

	var vtypeElem reflect.Type

	var keyFn, valFn *decFnSimpleIO
	var ktypeLo, vtypeLo = ktype, vtype

	if ktypeKind == reflect.Ptr {
		for ktypeLo = ktype.Elem(); ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() {
		}
	}

	if vtypePtr {
		vtypeElem = vtype.Elem()
		for vtypeLo = vtypeElem; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() {
		}
	}

	rvkMut := !scalarBitset.isset(ti.keykind)
	rvvMut := !scalarBitset.isset(ti.elemkind)
	rvvCanNil := isnilBitset.isset(ti.elemkind)

	var rvk, rvkn, rvv, rvvn, rvva, rvvz reflect.Value

	var doMapGet, doMapSet bool

	if !d.h.MapValueReset {
		if rvvMut && (vtypeKind != reflect.Interface || !d.h.InterfaceReset) {
			doMapGet = true
			rvva = mapAddrLoopvarRV(vtype, vtypeKind)
		}
	}

	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId
	hasLen := containerLen >= 0

	var kstr2bs []byte
	var kstr string

	var mapKeyStringSharesBytesBuf bool
	var att dBytesAttachState

	var vElem, kElem reflect.Type
	kbuiltin := ti.tikey.flagDecBuiltin && ti.keykind != uint8(reflect.Slice)
	vbuiltin := ti.tielem.flagDecBuiltin
	if kbuiltin && ktypePtr {
		kElem = ti.key.Elem()
	}
	if vbuiltin && vtypePtr {
		vElem = ti.elem.Elem()
	}

	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		mapKeyStringSharesBytesBuf = false
		kstr = ""
		if j == 0 {

			if kTransient {
				rvk = d.perType.TransientAddr2K(ktype, ktypeKind)
			} else {
				rvk = rvZeroAddrK(ktype, ktypeKind)
			}
			if !rvkMut {
				rvkn = rvk
			}
			if !rvvMut {
				if vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
			}
			if !ktypeIsString && keyFn == nil {
				keyFn = d.fn(ktypeLo)
			}
			if valFn == nil {
				valFn = d.fn(vtypeLo)
			}
		} else if rvkMut {
			rvSetZero(rvk)
		} else {
			rvk = rvkn
		}

		d.mapElemKey(j == 0)

		if d.d.TryNil() {
			rvSetZero(rvk)
		} else if ktypeIsString {
			kstr2bs, att = d.d.DecodeStringAsBytes()
			kstr, mapKeyStringSharesBytesBuf = d.bytes2Str(kstr2bs, att)
			rvSetString(rvk, kstr)
		} else {
			if kbuiltin {
				if ktypePtr {
					if rvIsNil(rvk) {
						rvSetDirect(rvk, reflect.New(kElem))
					}
					d.decode(rv2i(rvk))
				} else {
					d.decode(rv2i(rvAddr(rvk, ti.tikey.ptr)))
				}
			} else {
				d.decodeValueNoCheckNil(rvk, keyFn)
			}

			if ktypeIsIntf {
				if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp {
					kstr2bs = rvGetBytes(rvk2)
					kstr, mapKeyStringSharesBytesBuf = d.bytes2Str(kstr2bs, dBytesAttachView)
					rvSetIntf(rvk, rv4istr(kstr))
				}

			}
		}

		if mapKeyStringSharesBytesBuf && d.bufio {
			if ktypeIsString {
				rvSetString(rvk, d.detach2Str(kstr2bs, att))
			} else {
				rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
			}
			mapKeyStringSharesBytesBuf = false
		}

		d.mapElemValue()

		if d.d.TryNil() {
			if mapKeyStringSharesBytesBuf {
				if ktypeIsString {
					rvSetString(rvk, d.detach2Str(kstr2bs, att))
				} else {
					rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
				}
			}

			if !rvvz.IsValid() {
				rvvz = rvZeroK(vtype, vtypeKind)
			}
			mapSet(rv, rvk, rvvz, mparams)
			continue
		}

		doMapSet = true

		if !rvvMut {
			rvv = rvvn
		} else if !doMapGet {
			goto NEW_RVV
		} else {
			rvv = mapGet(rv, rvk, rvva, mparams)
			if !rvv.IsValid() || (rvvCanNil && rvIsNil(rvv)) {
				goto NEW_RVV
			}
			switch vtypeKind {
			case reflect.Ptr, reflect.Map:
				doMapSet = false
			case reflect.Interface:

				rvvn = rvv.Elem()
				if k := rvvn.Kind(); (k == reflect.Ptr || k == reflect.Map) && !rvIsNil(rvvn) {
					d.decodeValueNoCheckNil(rvvn, nil)
					continue
				}

				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvSetIntf(rvvn, rvv)
				rvv = rvvn
			default:

				if vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
				rvSetDirect(rvvn, rvv)
				rvv = rvvn
			}
		}
		goto DECODE_VALUE_NO_CHECK_NIL

	NEW_RVV:
		if vtypePtr {
			rvv = reflect.New(vtypeElem)
		} else if vTransient {
			rvv = d.perType.TransientAddrK(vtype, vtypeKind)
		} else {
			rvv = rvZeroAddrK(vtype, vtypeKind)
		}

	DECODE_VALUE_NO_CHECK_NIL:
		if doMapSet && mapKeyStringSharesBytesBuf {
			if ktypeIsString {
				rvSetString(rvk, d.detach2Str(kstr2bs, att))
			} else {
				rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
			}
		}
		if vbuiltin {
			if vtypePtr {
				if rvIsNil(rvv) {
					rvSetDirect(rvv, reflect.New(vElem))
				}
				d.decode(rv2i(rvv))
			} else {
				d.decode(rv2i(rvAddr(rvv, ti.tielem.ptr)))
			}
		} else {
			d.decodeValueNoCheckNil(rvv, valFn)
		}
		if doMapSet {
			mapSet(rv, rvk, rvv, mparams)
		}
	}

	d.mapEnd()
}

func (d *decoderSimpleIO) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsSimpleIO)

	if d.bytes {
		d.rtidFn = &d.h.rtidFnsDecBytes
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtBytes
	} else {
		d.bufio = d.h.ReaderBufferSize > 0
		d.rtidFn = &d.h.rtidFnsDecIO
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtIO
	}

	d.reset()

}

func (d *decoderSimpleIO) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
	d.depth = 0
	d.calls = 0

	d.maxdepth = decDefMaxDepth
	if d.h.MaxDepth > 0 {
		d.maxdepth = d.h.MaxDepth
	}
	d.mtid = 0
	d.stid = 0
	d.mtr = false
	d.str = false
	if d.h.MapType != nil {
		d.mtid = rt2id(d.h.MapType)
		_, d.mtr = fastpathAvIndex(d.mtid)
	}
	if d.h.SliceType != nil {
		d.stid = rt2id(d.h.SliceType)
		_, d.str = fastpathAvIndex(d.stid)
	}
}

func (d *decoderSimpleIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderSimpleIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderSimpleIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderSimpleIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderSimpleIO) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderSimpleIO) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderSimpleIO) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderSimpleIO) Release() {}

func (d *decoderSimpleIO) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderSimpleIO) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderSimpleIO) decode(iv interface{}) {
	_ = d.d

	rv, ok := isNil(iv, true)
	if ok {
		halt.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {

	case *string:
		*v = d.detach2Str(d.d.DecodeStringAsBytes())
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *uintptr:
		*v = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *float32:
		*v = d.d.DecodeFloat32()
	case *float64:
		*v = d.d.DecodeFloat64()
	case *complex64:
		*v = complex(d.d.DecodeFloat32(), 0)
	case *complex128:
		*v = complex(d.d.DecodeFloat64(), 0)
	case *[]byte:
		*v, _ = d.decodeBytesInto(*v, false)
	case []byte:

		d.decodeBytesInto(v[:len(v):len(v)], true)
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface{}:
		d.decodeValue(rv4iptr(v), nil)

	case reflect.Value:
		if ok, _ = isDecodeable(v); !ok {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)

	default:

		if skipFastpathTypeSwitchInDirectCall || !d.dh.fastpathDecodeTypeSwitch(iv, d) {
			if !rv.IsValid() {
				rv = reflect.ValueOf(iv)
			}
			if ok, _ = isDecodeable(rv); !ok {
				d.haltAsNotDecodeable(rv)
			}
			d.decodeValue(rv, nil)
		}
	}
}

func (d *decoderSimpleIO) decodeValue(rv reflect.Value, fn *decFnSimpleIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderSimpleIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnSimpleIO) {

	var rvp reflect.Value
	var rvpValid bool
PTR:
	if rv.Kind() == reflect.Ptr {
		rvpValid = true
		if rvIsNil(rv) {
			rvSetDirect(rv, reflect.New(rv.Type().Elem()))
		}
		rvp = rv
		rv = rv.Elem()
		goto PTR
	}

	if fn == nil {
		fn = d.fn(rv.Type())
	}
	if fn.i.addrD {
		if rvpValid {
			rv = rvp
		} else if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else if fn.i.addrDf {
			halt.errorStr("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *decoderSimpleIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderSimpleIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderSimpleIO) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
	v, att := d.d.DecodeBytes()
	if cap(v) == 0 || (att >= dBytesAttachViewZerocopy && !mustFit) {

		return
	}
	if len(v) == 0 {
		v = zeroByteSlice
		return
	}
	if len(out) == len(v) {
		state = dBytesIntoParamOut
	} else if cap(out) >= len(v) {
		out = out[:len(v)]
		state = dBytesIntoParamOutSlice
	} else if mustFit {
		halt.errorf("bytes capacity insufficient for decoded bytes: got/expected: %d/%d", len(v), len(out))
	} else {
		out = make([]byte, len(v))
		state = dBytesIntoNew
	}
	copy(out, v)
	v = out
	return
}

func (d *decoderSimpleIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderSimpleIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderSimpleIO) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderSimpleIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderSimpleIO) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderSimpleIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderSimpleIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleIO) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderSimpleIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderSimpleIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderSimpleIO) fn(t reflect.Type) *decFnSimpleIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderSimpleIO) fnNoExt(t reflect.Type) *decFnSimpleIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverSimpleIO) newDecoderBytes(in []byte, h Handle) *decoderSimpleIO {
	var c1 decoderSimpleIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverSimpleIO) newDecoderIO(in io.Reader, h Handle) *decoderSimpleIO {
	var c1 decoderSimpleIO
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverSimpleIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsSimpleIO) (f *fastpathDSimpleIO, u reflect.Type) {
	rtid := rt2id(ti.fastpathUnderlying)
	idx, ok := fastpathAvIndex(rtid)
	if !ok {
		return
	}
	f = &fp[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflect.ArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

func (helperDecDriverSimpleIO) decFindRtidFn(s []decRtidFnSimpleIO, rtid uintptr) (i uint, fn *decFnSimpleIO) {

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (helperDecDriverSimpleIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnSimpleIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnSimpleIO](v))
	}
	return
}

func (dh helperDecDriverSimpleIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsSimpleIO,
	checkExt bool) (fn *decFnSimpleIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverSimpleIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnSimpleIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverSimpleIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnSimpleIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnSimpleIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnSimpleIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnSimpleIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverSimpleIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsSimpleIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnSimpleIO) {
	fn = new(decFnSimpleIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderSimpleIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderSimpleIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderSimpleIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderSimpleIO).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderSimpleIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderSimpleIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderSimpleIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderSimpleIO).textUnmarshal
		fi.addrD = ti.flagTextUnmarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			var rtid2 uintptr
			if !ti.flagHasPkgPath {
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fd = fp[idx].decfn
					fi.addrD = true
					fi.addrDf = false
					if rk == reflect.Array {
						fi.addrD = false
					}
				}
			} else {

				xfe, xrt := dh.decFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf2 := xfe.decfn
					if rk == reflect.Array {
						fi.addrD = false
						fn.fd = func(d *decoderSimpleIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderSimpleIO, xf *decFnInfo, xrv reflect.Value) {
							if xrv.Kind() == reflect.Ptr {
								xfnf2(d, xf, rvConvert(xrv, xptr2rt))
							} else {
								xfnf2(d, xf, rvConvert(xrv, xrt))
							}
						}
					}
				}
			}
		}
		if fn.fd == nil {
			switch rk {
			case reflect.Bool:
				fn.fd = (*decoderSimpleIO).kBool
			case reflect.String:
				fn.fd = (*decoderSimpleIO).kString
			case reflect.Int:
				fn.fd = (*decoderSimpleIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderSimpleIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderSimpleIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderSimpleIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderSimpleIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderSimpleIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderSimpleIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderSimpleIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderSimpleIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderSimpleIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderSimpleIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderSimpleIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderSimpleIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderSimpleIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderSimpleIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderSimpleIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderSimpleIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderSimpleIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderSimpleIO).kStructSimple
				} else {
					fn.fd = (*decoderSimpleIO).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderSimpleIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderSimpleIO).kInterface
			default:

				fn.fd = (*decoderSimpleIO).kErr
			}
		}
	}
	return
}
func (e *simpleEncDriverIO) EncodeNil() {
	e.w.writen1(simpleVdNil)
}

func (e *simpleEncDriverIO) EncodeBool(b bool) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && !b {
		e.EncodeNil()
		return
	}
	if b {
		e.w.writen1(simpleVdTrue)
	} else {
		e.w.writen1(simpleVdFalse)
	}
}

func (e *simpleEncDriverIO) EncodeFloat32(f float32) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat32)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *simpleEncDriverIO) EncodeFloat64(f float64) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *simpleEncDriverIO) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), simpleVdNegInt)
	} else {
		e.encUint(uint64(v), simpleVdPosInt)
	}
}

func (e *simpleEncDriverIO) EncodeUint(v uint64) {
	e.encUint(v, simpleVdPosInt)
}

func (e *simpleEncDriverIO) encUint(v uint64, bd uint8) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == 0 {
		e.EncodeNil()
		return
	}
	if v <= math.MaxUint8 {
		e.w.writen2(bd, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 1)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 2)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 3)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *simpleEncDriverIO) encLen(bd byte, length int) {
	if length == 0 {
		e.w.writen1(bd)
	} else if length <= math.MaxUint8 {
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	} else if length <= math.MaxUint16 {
		e.w.writen1(bd + 2)
		e.w.writen2(bigen.PutUint16(uint16(length)))
	} else if int64(length) <= math.MaxUint32 {
		e.w.writen1(bd + 3)
		e.w.writen4(bigen.PutUint32(uint32(length)))
	} else {
		e.w.writen1(bd + 4)
		e.w.writen8(bigen.PutUint64(uint64(length)))
	}
}

func (e *simpleEncDriverIO) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	var bs0, bs []byte
	if ext == SelfExt {
		bs0 = e.e.blist.get(1024)
		bs = bs0
		sideEncode(e.h, &e.h.sideEncPool, func(se encoderI) { oneOffEncode(se, v, &bs, basetype, false) })
	} else {
		bs = ext.WriteExt(v)
	}
	if bs == nil {
		e.writeNilBytes()
		goto END
	}
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *simpleEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *simpleEncDriverIO) encodeExtPreamble(xtag byte, length int) {
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
}

func (e *simpleEncDriverIO) WriteArrayStart(length int) {
	e.encLen(simpleVdArray, length)
}

func (e *simpleEncDriverIO) WriteMapStart(length int) {
	e.encLen(simpleVdMap, length)
}

func (e *simpleEncDriverIO) WriteArrayEmpty() {

	e.w.writen1(simpleVdArray)
}

func (e *simpleEncDriverIO) WriteMapEmpty() {

	e.w.writen1(simpleVdMap)
}

func (e *simpleEncDriverIO) EncodeString(v string) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && v == "" {
		e.EncodeNil()
		return
	}
	if e.h.StringToRaw {
		e.encLen(simpleVdByteArray, len(v))
	} else {
		e.encLen(simpleVdString, len(v))
	}
	e.w.writestr(v)
}

func (e *simpleEncDriverIO) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *simpleEncDriverIO) EncodeStringBytesRaw(v []byte) {

	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
}

func (e *simpleEncDriverIO) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *simpleEncDriverIO) encodeNilBytes() {
	b := byte(simpleVdNil)
	if e.h.NilCollectionToZeroLength {
		b = simpleVdArray
	}
	e.w.writen1(b)
}

func (e *simpleEncDriverIO) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = simpleVdNil
	}
	e.w.writen1(v)
}

func (e *simpleEncDriverIO) writeNilArray() {
	e.writeNilOr(simpleVdArray)
}

func (e *simpleEncDriverIO) writeNilMap() {
	e.writeNilOr(simpleVdMap)
}

func (e *simpleEncDriverIO) writeNilBytes() {
	e.writeNilOr(simpleVdByteArray)
}

func (e *simpleEncDriverIO) EncodeTime(t time.Time) {

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	v, err := t.MarshalBinary()
	halt.onerror(err)
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
}

func (d *simpleDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *simpleDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *simpleDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	switch d.bd {
	case simpleVdNil:
		d.bdRead = false
		return valueTypeNil
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		return valueTypeBytes
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		return valueTypeString
	case simpleVdArray, simpleVdArray + 1,
		simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		return valueTypeArray
	case simpleVdMap, simpleVdMap + 1,
		simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *simpleDecDriverIO) TryNil() bool {
	return d.advanceNil()
}

func (d *simpleDecDriverIO) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case simpleVdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		ok = false
	}
	return
}

func (d *simpleDecDriverIO) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd {
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case simpleVdPosInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case simpleVdPosInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(bigen.Uint16(d.r.readn2()))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(bigen.Uint32(d.r.readn4()))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(bigen.Uint64(d.r.readn8()))
		neg = true
	default:
		ok = false

	}

	return
}

func (d *simpleDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == simpleVdFalse {
	} else if d.bd == simpleVdTrue {
		b = true
	} else {
		halt.errorf("cannot decode bool - %s: %x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverIO) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriverIO) uint2Len(ui uint64) int {
	if chkOvf.Uint(ui, intBitsize) {
		halt.errorf("overflow integer: %v", ui)
	}
	return int(ui)
}

func (d *simpleDecDriverIO) decLen() int {
	switch d.bd & 7 {
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(bigen.Uint16(d.r.readn2()))
	case 3:
		return d.uint2Len(uint64(bigen.Uint32(d.r.readn4())))
	case 4:
		return d.uint2Len(bigen.Uint64(d.r.readn8()))
	}
	halt.errorf("cannot read length: bd%%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
}

func (d *simpleDecDriverIO) DecodeStringAsBytes() ([]byte, dBytesAttachState) {
	return d.DecodeBytes()
}

func (d *simpleDecDriverIO) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool

	if d.bd >= simpleVdArray && d.bd <= simpleVdArray+4 {
		slen := d.ReadArrayStart()
		bs, cond = usableByteSlice(d.d.buf, slen)
		for i := 0; i < len(bs); i++ {
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		}
		for i := len(bs); i < slen; i++ {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		if cond {
			d.d.buf = bs
		}
		state = dBytesAttachBuffer

		return
	}

	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *simpleDecDriverIO) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd != simpleVdTime {
		halt.errorf("invalid descriptor for time.Time - expect 0x%x, received 0x%x", simpleVdTime, d.bd)
	}
	d.bdRead = false
	clen := uint(d.r.readn1())
	b := d.r.readx(clen)
	halt.onerror((&t).UnmarshalBinary(b))
	return
}

func (d *simpleDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	xbs, _, _, ok := d.decodeExtV(ext != nil, xtag)
	if !ok {
		return
	}
	if ext == SelfExt {
		sideDecode(d.h, &d.h.sideDecPool, func(sd decoderI) { oneOffDecode(sd, rv, xbs, basetype, false) })
	} else {
		ext.ReadExt(rv, xbs)
	}
}

func (d *simpleDecDriverIO) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *simpleDecDriverIO) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
	if xtagIn > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtagIn)
	}
	if d.advanceNil() {
		return
	}
	tag := uint8(xtagIn)
	switch d.bd {
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag. Got %b. Expecting: %v", xtag, tag)
		}
		xbs, ok = d.r.readxb(uint(l))
		bstate = d.d.attachState(ok)
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs, bstate = d.DecodeBytes()
	default:
		halt.errorf("ext - %s - expecting extensions/bytearray, got: 0x%x", msgBadDesc, d.bd)
	}
	d.bdRead = false
	ok = true
	return
}

func (d *simpleDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.bd {
	case simpleVdNil:
		n.v = valueTypeNil
	case simpleVdFalse:
		n.v = valueTypeBool
		n.b = false
	case simpleVdTrue:
		n.v = valueTypeBool
		n.b = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case simpleVdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdTime:
		n.v = valueTypeTime
		n.t = d.DecodeTime()
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(uint(l))

	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2,
		simpleVdArray + 3, simpleVdArray + 4:
		n.v = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, d.bd)
	}

	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *simpleDecDriverIO) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *simpleDecDriverIO) nextValueBytesBdReadR() {
	c := d.bd

	var length uint

	switch c {
	case simpleVdNil, simpleVdFalse, simpleVdTrue, simpleVdString, simpleVdByteArray:

	case simpleVdPosInt, simpleVdNegInt:
		d.r.readn1()
	case simpleVdPosInt + 1, simpleVdNegInt + 1:
		d.r.skip(2)
	case simpleVdPosInt + 2, simpleVdNegInt + 2, simpleVdFloat32:
		d.r.skip(4)
	case simpleVdPosInt + 3, simpleVdNegInt + 3, simpleVdFloat64:
		d.r.skip(8)
	case simpleVdTime:
		c = d.r.readn1()
		d.r.skip(uint(c))

	default:
		switch c & 7 {
		case 0:
			length = 0
		case 1:
			b := d.r.readn1()
			length = uint(b)
		case 2:
			x := d.r.readn2()
			length = uint(bigen.Uint16(x))
		case 3:
			x := d.r.readn4()
			length = uint(bigen.Uint32(x))
		case 4:
			x := d.r.readn8()
			length = uint(bigen.Uint64(x))
		}

		bExt := c >= simpleVdExt && c <= simpleVdExt+7
		bStr := c >= simpleVdString && c <= simpleVdString+7
		bByteArray := c >= simpleVdByteArray && c <= simpleVdByteArray+7
		bArray := c >= simpleVdArray && c <= simpleVdArray+7
		bMap := c >= simpleVdMap && c <= simpleVdMap+7

		if !(bExt || bStr || bByteArray || bArray || bMap) {
			halt.errorf("cannot infer value - %s 0x%x", msgBadDesc, c)
		}

		if bExt {
			d.r.readn1()
		}

		if length == 0 {
			break
		}

		if bArray {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else if bMap {
			for i := uint(0); i < length; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			d.r.skip(length)
		}
	}
	return
}

func (d *simpleEncDriverIO) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*SimpleHandle)
	d.e = shared
	if shared.bytes {
		fp = simpleFpEncBytes
	} else {
		fp = simpleFpEncIO
	}

	d.init2(enc)
	return
}

func (e *simpleEncDriverIO) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *simpleEncDriverIO) writerEnd() { e.w.end() }

func (e *simpleEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *simpleEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *simpleDecDriverIO) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*SimpleHandle)
	d.d = shared
	if shared.bytes {
		fp = simpleFpDecBytes
	} else {
		fp = simpleFpDecIO
	}

	d.init2(dec)
	return
}

func (d *simpleDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *simpleDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *simpleDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *simpleDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, simpledesc(d.bd))
}

func (d *simpleDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
