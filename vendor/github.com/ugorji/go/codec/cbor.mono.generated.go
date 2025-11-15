//go:build !notmono && !codec.notmono 

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"

	"io"
	"math"
	"math/big"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

type helperEncDriverCborBytes struct{}
type encFnCborBytes struct {
	i  encFnInfo
	fe func(*encoderCborBytes, *encFnInfo, reflect.Value)
}
type encRtidFnCborBytes struct {
	rtid uintptr
	fn   *encFnCborBytes
}
type encoderCborBytes struct {
	dh helperEncDriverCborBytes
	fp *fastpathEsCborBytes
	e  cborEncDriverBytes
	encoderBase
}
type helperDecDriverCborBytes struct{}
type decFnCborBytes struct {
	i  decFnInfo
	fd func(*decoderCborBytes, *decFnInfo, reflect.Value)
}
type decRtidFnCborBytes struct {
	rtid uintptr
	fn   *decFnCborBytes
}
type decoderCborBytes struct {
	dh helperDecDriverCborBytes
	fp *fastpathDsCborBytes
	d  cborDecDriverBytes
	decoderBase
}
type cborEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoState
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT

	h   *CborHandle
	e   *encoderBase
	w   bytesEncAppender
	enc encoderI

	b [40]byte
}
type cborDecDriverBytes struct {
	decDriverNoopContainerReader

	noBuiltInTypes

	h   *CborHandle
	d   *decoderBase
	r   bytesDecReader
	dec decoderI
	bdAndBdread
}

func (e *encoderCborBytes) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderCborBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderCborBytes) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderCborBytes) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderCborBytes) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderCborBytes) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderCborBytes) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderCborBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderCborBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderCborBytes) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderCborBytes) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderCborBytes) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderCborBytes) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderCborBytes) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderCborBytes) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderCborBytes) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderCborBytes) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderCborBytes) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderCborBytes) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderCborBytes) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderCborBytes) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderCborBytes) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderCborBytes) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderCborBytes) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderCborBytes) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderCborBytes) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderCborBytes) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderCborBytes) kSeqFn(rt reflect.Type) (fn *encFnCborBytes) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderCborBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnCborBytes
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

func (e *encoderCborBytes) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnCborBytes
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

func (e *encoderCborBytes) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderCborBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderCborBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderCborBytes) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderCborBytes) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderCborBytes) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderCborBytes) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderCborBytes) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnCborBytes

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

func (e *encoderCborBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnCborBytes) {
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

func (e *encoderCborBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsCborBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderCborBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderCborBytes) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderCborBytes) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderCborBytes) mustEncode(v interface{}) {
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

func (e *encoderCborBytes) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderCborBytes) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderCborBytes) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderCborBytes) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderCborBytes) encodeValue(rv reflect.Value, fn *encFnCborBytes) {

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

func (e *encoderCborBytes) encodeValueNonNil(rv reflect.Value, fn *encFnCborBytes) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderCborBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderCborBytes) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderCborBytes) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderCborBytes) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderCborBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderCborBytes) fn(t reflect.Type) *encFnCborBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderCborBytes) fnNoExt(t reflect.Type) *encFnCborBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderCborBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderCborBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderCborBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderCborBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderCborBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderCborBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderCborBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderCborBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverCborBytes) newEncoderBytes(out *[]byte, h Handle) *encoderCborBytes {
	var c1 encoderCborBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverCborBytes) newEncoderIO(out io.Writer, h Handle) *encoderCborBytes {
	var c1 encoderCborBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverCborBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsCborBytes) (f *fastpathECborBytes, u reflect.Type) {
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

func (helperEncDriverCborBytes) encFindRtidFn(s []encRtidFnCborBytes, rtid uintptr) (i uint, fn *encFnCborBytes) {

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

func (helperEncDriverCborBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnCborBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnCborBytes](v))
	}
	return
}

func (dh helperEncDriverCborBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsCborBytes, checkExt bool) (fn *encFnCborBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverCborBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnCborBytes = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverCborBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnCborBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnCborBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnCborBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnCborBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverCborBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborBytes) {
	fn = new(encFnCborBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderCborBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderCborBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderCborBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderCborBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderCborBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderCborBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderCborBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderCborBytes).textMarshal
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
					fn.fe = func(e *encoderCborBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderCborBytes).kBool
			case reflect.String:

				fn.fe = (*encoderCborBytes).kString
			case reflect.Int:
				fn.fe = (*encoderCborBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderCborBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderCborBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderCborBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderCborBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderCborBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderCborBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderCborBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderCborBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderCborBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderCborBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderCborBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderCborBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderCborBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderCborBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderCborBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderCborBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderCborBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderCborBytes).kStructSimple
				} else {
					fn.fe = (*encoderCborBytes).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderCborBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderCborBytes).kErr
			default:

				fn.fe = (*encoderCborBytes).kErr
			}
		}
	}
	return
}
func (d *decoderCborBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderCborBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderCborBytes) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderCborBytes) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderCborBytes) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderCborBytes) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderCborBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderCborBytes) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderCborBytes) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderCborBytes) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderCborBytes) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderCborBytes) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderCborBytes) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderCborBytes) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderCborBytes) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderCborBytes) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderCborBytes) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderCborBytes) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderCborBytes) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderCborBytes) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderCborBytes) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderCborBytes) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborBytes) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborBytes) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderCborBytes) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderCborBytes) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderCborBytes) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderCborBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderCborBytes) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderCborBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderCborBytes) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderCborBytes) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderCborBytes) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnCborBytes

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

func (d *decoderCborBytes) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnCborBytes
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

func (d *decoderCborBytes) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnCborBytes

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

func (d *decoderCborBytes) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnCborBytes
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

func (d *decoderCborBytes) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsCborBytes)

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

func (d *decoderCborBytes) reset() {
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

func (d *decoderCborBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderCborBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderCborBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderCborBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderCborBytes) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderCborBytes) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderCborBytes) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderCborBytes) Release() {}

func (d *decoderCborBytes) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderCborBytes) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderCborBytes) decode(iv interface{}) {
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

func (d *decoderCborBytes) decodeValue(rv reflect.Value, fn *decFnCborBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderCborBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnCborBytes) {

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

func (d *decoderCborBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderCborBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderCborBytes) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderCborBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderCborBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderCborBytes) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderCborBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderCborBytes) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderCborBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderCborBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborBytes) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderCborBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderCborBytes) fn(t reflect.Type) *decFnCborBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderCborBytes) fnNoExt(t reflect.Type) *decFnCborBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverCborBytes) newDecoderBytes(in []byte, h Handle) *decoderCborBytes {
	var c1 decoderCborBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverCborBytes) newDecoderIO(in io.Reader, h Handle) *decoderCborBytes {
	var c1 decoderCborBytes
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverCborBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsCborBytes) (f *fastpathDCborBytes, u reflect.Type) {
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

func (helperDecDriverCborBytes) decFindRtidFn(s []decRtidFnCborBytes, rtid uintptr) (i uint, fn *decFnCborBytes) {

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

func (helperDecDriverCborBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnCborBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnCborBytes](v))
	}
	return
}

func (dh helperDecDriverCborBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsCborBytes,
	checkExt bool) (fn *decFnCborBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverCborBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnCborBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverCborBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnCborBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnCborBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnCborBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnCborBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverCborBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsCborBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborBytes) {
	fn = new(decFnCborBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderCborBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderCborBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderCborBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderCborBytes).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderCborBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderCborBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderCborBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderCborBytes).textUnmarshal
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
						fn.fd = func(d *decoderCborBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderCborBytes, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderCborBytes).kBool
			case reflect.String:
				fn.fd = (*decoderCborBytes).kString
			case reflect.Int:
				fn.fd = (*decoderCborBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderCborBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderCborBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderCborBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderCborBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderCborBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderCborBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderCborBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderCborBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderCborBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderCborBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderCborBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderCborBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderCborBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderCborBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderCborBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderCborBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderCborBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderCborBytes).kStructSimple
				} else {
					fn.fd = (*decoderCborBytes).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderCborBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderCborBytes).kInterface
			default:

				fn.fd = (*decoderCborBytes).kErr
			}
		}
	}
	return
}
func (e *cborEncDriverBytes) EncodeNil() {
	e.w.writen1(cborBdNil)
}

func (e *cborEncDriverBytes) EncodeBool(b bool) {
	if b {
		e.w.writen1(cborBdTrue)
	} else {
		e.w.writen1(cborBdFalse)
	}
}

func (e *cborEncDriverBytes) EncodeFloat32(f float32) {
	b := math.Float32bits(f)
	if e.h.OptimumSize {
		if h := floatToHalfFloatBits(b); halfFloatToFloatBits(h) == b {
			e.w.writen1(cborBdFloat16)
			e.w.writen2(bigen.PutUint16(h))
			return
		}
	}
	e.w.writen1(cborBdFloat32)
	e.w.writen4(bigen.PutUint32(b))
}

func (e *cborEncDriverBytes) EncodeFloat64(f float64) {
	if e.h.OptimumSize {
		if f32 := float32(f); float64(f32) == f {
			e.EncodeFloat32(f32)
			return
		}
	}
	e.w.writen1(cborBdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *cborEncDriverBytes) encUint(v uint64, bd byte) {
	if v <= 0x17 {
		e.w.writen1(byte(v) + bd)
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd+0x18, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 0x19)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 0x1a)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 0x1b)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *cborEncDriverBytes) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-1-v), cborBaseNegInt)
	} else {
		e.encUint(uint64(v), cborBaseUint)
	}
}

func (e *cborEncDriverBytes) EncodeUint(v uint64) {
	e.encUint(v, cborBaseUint)
}

func (e *cborEncDriverBytes) encLen(bd byte, length int) {
	e.encUint(uint64(length), bd)
}

func (e *cborEncDriverBytes) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else if e.h.TimeRFC3339 {
		e.encUint(0, cborBaseTag)
		e.encStringBytesS(cborBaseString, stringView(t.AppendFormat(e.b[:0], time.RFC3339Nano)))
	} else {
		e.encUint(1, cborBaseTag)
		t = t.UTC().Round(time.Microsecond)
		sec, nsec := t.Unix(), uint64(t.Nanosecond())
		if nsec == 0 {
			e.EncodeInt(sec)
		} else {
			e.EncodeFloat64(float64(sec) + float64(nsec)/1e9)
		}
	}
}

func (e *cborEncDriverBytes) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	e.encUint(uint64(xtag), cborBaseTag)
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.writeNilBytes()
	} else {
		e.enc.encodeI(v)
	}
}

func (e *cborEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encUint(uint64(re.Tag), cborBaseTag)
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encodeI(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *cborEncDriverBytes) WriteArrayEmpty() {
	if e.h.IndefiniteLength {
		e.w.writen2(cborBdIndefiniteArray, cborBdBreak)
	} else {
		e.w.writen1(cborBaseArray)

	}
}

func (e *cborEncDriverBytes) WriteMapEmpty() {
	if e.h.IndefiniteLength {
		e.w.writen2(cborBdIndefiniteMap, cborBdBreak)
	} else {
		e.w.writen1(cborBaseMap)

	}
}

func (e *cborEncDriverBytes) WriteArrayStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteArray)
	} else {
		e.encLen(cborBaseArray, length)
	}
}

func (e *cborEncDriverBytes) WriteMapStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteMap)
	} else {
		e.encLen(cborBaseMap, length)
	}
}

func (e *cborEncDriverBytes) WriteMapEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverBytes) WriteArrayEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverBytes) EncodeString(v string) {
	bb := cborBaseString
	if e.h.StringToRaw {
		bb = cborBaseBytes
	}
	e.encStringBytesS(bb, v)
}

func (e *cborEncDriverBytes) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *cborEncDriverBytes) EncodeStringBytesRaw(v []byte) {
	e.encStringBytesS(cborBaseBytes, stringView(v))
}

func (e *cborEncDriverBytes) encStringBytesS(bb byte, v string) {
	if e.h.IndefiniteLength {
		if bb == cborBaseBytes {
			e.w.writen1(cborBdIndefiniteBytes)
		} else {
			e.w.writen1(cborBdIndefiniteString)
		}
		vlen := uint(len(v))
		n := max(4, min(vlen/4, 1024))
		for i := uint(0); i < vlen; {
			i2 := i + n
			if i2 >= vlen {
				i2 = vlen
			}
			v2 := v[i:i2]
			e.encLen(bb, len(v2))
			e.w.writestr(v2)
			i = i2
		}
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(bb, len(v))
		e.w.writestr(v)
	}
}

func (e *cborEncDriverBytes) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *cborEncDriverBytes) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = cborBdNil
	}
	e.w.writen1(v)
}

func (e *cborEncDriverBytes) writeNilArray() {
	e.writeNilOr(cborBaseArray)
}

func (e *cborEncDriverBytes) writeNilMap() {
	e.writeNilOr(cborBaseMap)
}

func (e *cborEncDriverBytes) writeNilBytes() {
	e.writeNilOr(cborBaseBytes)
}

func (d *cborDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *cborDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdNil || d.bd == cborBdUndefined {
		d.bdRead = false
		return true
	}
	return
}

func (d *cborDecDriverBytes) TryNil() bool {
	return d.advanceNil()
}

func (d *cborDecDriverBytes) skipTags() {
	for d.bd>>5 == cborMajorTag {
		d.decUint()
		d.bd = d.r.readn1()
	}
}

func (d *cborDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	if d.bd == cborBdNil {
		d.bdRead = false
		return valueTypeNil
	}
	major := d.bd >> 5
	if major == cborMajorBytes {
		return valueTypeBytes
	} else if major == cborMajorString {
		return valueTypeString
	} else if major == cborMajorArray {
		return valueTypeArray
	} else if major == cborMajorMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *cborDecDriverBytes) CheckBreak() (v bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdBreak {
		d.bdRead = false
		v = true
	}
	return
}

func (d *cborDecDriverBytes) decUint() (ui uint64) {
	v := d.bd & 0x1f
	if v <= 0x17 {
		ui = uint64(v)
	} else if v == 0x18 {
		ui = uint64(d.r.readn1())
	} else if v == 0x19 {
		ui = uint64(bigen.Uint16(d.r.readn2()))
	} else if v == 0x1a {
		ui = uint64(bigen.Uint32(d.r.readn4()))
	} else if v == 0x1b {
		ui = uint64(bigen.Uint64(d.r.readn8()))
	} else {
		halt.errorf("invalid descriptor decoding uint: %x/%s (%x)", d.bd, cbordesc(d.bd), v)
	}
	return
}

func (d *cborDecDriverBytes) decLen() int {
	return int(d.decUint())
}

func (d *cborDecDriverBytes) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.r.readn2()))))
	case cborBdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case cborBdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		if d.bd>>5 == cborMajorTag {

			switch d.bd & 0x1f {
			case 2:
				f = d.decTagBigIntAsFloat(false)
			case 3:
				f = d.decTagBigIntAsFloat(true)
			case 4:
				f = d.decTagBigFloatAsFloat(true)
			case 5:
				f = d.decTagBigFloatAsFloat(false)
			default:
				ok = false
			}
		} else {
			ok = false
		}
	}
	return
}

func (d *cborDecDriverBytes) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd >> 5 {
	case cborMajorUint:
		ui = d.decUint()
	case cborMajorNegInt:
		ui = d.decUint()
		neg = true
	default:
		ok = false
	}
	return
}

func (d *cborDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	if d.bd == cborBdTrue {
		b = true
	} else if d.bd == cborBdFalse {
	} else {
		halt.errorf("not bool - %s %x/%s", msgBadDesc, d.bd, cbordesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorMap {
		halt.errorf("error reading map; got major type: %x, expected %x/%s", d.bd>>5, cborMajorMap, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverBytes) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorArray {
		halt.errorf("invalid array; got major type: %x, expect: %x/%s", d.bd>>5, cborMajorArray, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverBytes) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	fnEnsureNonNilBytes := func() {

		if bs == nil {
			bs = zeroByteSlice
			state = dBytesDetach
		}
	}
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		major := d.bd >> 5
		val4str := d.h.ValidateUnicode && major == cborMajorString
		bs = d.d.buf[:0]
		d.bdRead = false
		for !d.CheckBreak() {
			if d.bd>>5 != major {
				const msg = "malformed indefinite string/bytes %x (%s); " +
					"contains chunk with major type %v, expected %v"
				halt.errorf(msg, d.bd, cbordesc(d.bd), d.bd>>5, major)
			}
			n := uint(d.decLen())
			bs = append(bs, d.r.readx(n)...)
			d.bdRead = false
			if val4str && !utf8.Valid(bs[len(bs)-int(n):]) {
				const msg = "indefinite-length text string contains chunk " +
					"that is not a valid utf-8 sequence: 0x%x"
				halt.errorf(msg, bs[len(bs)-int(n):])
			}
		}
		d.bdRead = false
		d.d.buf = bs
		state = dBytesAttachBuffer
		fnEnsureNonNilBytes()
		return
	}
	if d.bd == cborBdIndefiniteArray {
		d.bdRead = false
		bs = d.d.buf[:0]
		for !d.CheckBreak() {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		d.d.buf = bs
		state = dBytesAttachBuffer
		fnEnsureNonNilBytes()
		return
	}
	var cond bool
	if d.bd>>5 == cborMajorArray {
		d.bdRead = false
		slen := d.decLen()
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
		fnEnsureNonNilBytes()
		return
	}
	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *cborDecDriverBytes) DecodeStringAsBytes() (out []byte, state dBytesAttachState) {
	out, state = d.DecodeBytes()
	if d.h.ValidateUnicode && !utf8.Valid(out) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", out)
	}
	return
}

func (d *cborDecDriverBytes) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
}

func (d *cborDecDriverBytes) decodeTime(xtag uint64) (t time.Time) {
	switch xtag {
	case 0:
		var err error
		t, err = time.Parse(time.RFC3339, stringView(bytesOKs(d.DecodeStringAsBytes())))
		halt.onerror(err)
	case 1:
		f1, f2 := math.Modf(d.DecodeFloat64())
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		halt.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	}
	t = t.UTC().Round(time.Microsecond)
	return
}

func (d *cborDecDriverBytes) preDecodeExt(checkTag bool, xtag uint64) (realxtag uint64, ok bool) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	realxtag = d.decUint()
	d.bdRead = false
	if checkTag && xtag != realxtag {
		halt.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
	}
	ok = true
	return
}

func (d *cborDecDriverBytes) DecodeRawExt(re *RawExt) {
	if realxtag, ok := d.preDecodeExt(false, 0); ok {
		re.Tag = realxtag
		d.dec.decode(&re.Value)
		d.bdRead = false
	}
}

func (d *cborDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if _, ok := d.preDecodeExt(true, xtag); ok {
		if ext == SelfExt {
			d.dec.decodeAs(rv, basetype, false)
		} else {
			d.dec.interfaceExtConvertAndDecode(rv, ext)
		}
		d.bdRead = false
	}
}

func (d *cborDecDriverBytes) decTagBigIntAsFloat(neg bool) (f float64) {
	bs, _ := d.DecodeBytes()
	bi := new(big.Int).SetBytes(bs)
	if neg {
		bi0 := bi
		bi = new(big.Int).Sub(big.NewInt(-1), bi0)
	}
	f, _ = bi.Float64()
	return
}

func (d *cborDecDriverBytes) decTagBigFloatAsFloat(decimal bool) (f float64) {
	if nn := d.r.readn1(); nn != 82 {
		halt.errorf("(%d) decoding decimal/big.Float: expected 2 numbers", nn)
	}
	exp := d.DecodeInt64()
	mant := d.DecodeInt64()
	if decimal {

		rf := readFloatResult{exp: int8(exp)}
		if mant >= 0 {
			rf.mantissa = uint64(mant)
		} else {
			rf.neg = true
			rf.mantissa = uint64(-mant)
		}
		f, _ = parseFloat64_reader(rf)

	} else {

		bfm := new(big.Float).SetPrec(64).SetInt64(mant)
		bf := new(big.Float).SetPrec(64).SetMantExp(bfm, int(exp))
		f, _ = bf.Float64()
	}
	return
}

func (d *cborDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool
	switch d.bd >> 5 {
	case cborMajorUint:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case cborMajorNegInt:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case cborMajorBytes:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case cborMajorString:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case cborMajorArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborMajorMap:
		n.v = valueTypeMap
		decodeFurther = true
	case cborMajorTag:
		n.v = valueTypeExt
		n.u = d.decUint()
		d.bdRead = false
		n.l = nil
		xx := d.h.getExtForTag(n.u)
		if xx == nil {
			switch n.u {
			case 0, 1:
				n.v = valueTypeTime
				n.t = d.decodeTime(n.u)
			case 2:
				n.f = d.decTagBigIntAsFloat(false)
				n.v = valueTypeFloat
			case 3:
				n.f = d.decTagBigIntAsFloat(true)
				n.v = valueTypeFloat
			case 4:
				n.f = d.decTagBigFloatAsFloat(true)
				n.v = valueTypeFloat
			case 5:
				n.f = d.decTagBigFloatAsFloat(false)
				n.v = valueTypeFloat
			case 55799:
				d.DecodeNaked()
			default:
				if d.h.SkipUnexpectedTags {
					d.DecodeNaked()
				}

			}
			return
		}

	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined:
			n.v = valueTypeNil
		case cborBdFalse:
			n.v = valueTypeBool
			n.b = false
		case cborBdTrue:
			n.v = valueTypeBool
			n.b = true
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			n.v = valueTypeFloat
			n.f = d.DecodeFloat64()
		default:
			halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	}
	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *cborDecDriverBytes) uintBytes() (v []byte, ui uint64) {

	switch vv := d.bd & 0x1f; vv {
	case 0x18:
		v = d.r.readx(1)
		ui = uint64(v[0])
	case 0x19:
		v = d.r.readx(2)
		ui = uint64(bigenstd.Uint16(v))
	case 0x1a:
		v = d.r.readx(4)
		ui = uint64(bigenstd.Uint32(v))
	case 0x1b:
		v = d.r.readx(8)
		ui = uint64(bigenstd.Uint64(v))
	default:
		if vv > 0x1b {
			halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
		}
		ui = uint64(vv)
	}
	return
}

func (d *cborDecDriverBytes) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *cborDecDriverBytes) nextValueBytesBdReadR() {

	var ui uint64

	switch d.bd >> 5 {
	case cborMajorUint, cborMajorNegInt:
		d.uintBytes()
	case cborMajorString, cborMajorBytes:
		if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				_, ui = d.uintBytes()
				d.r.skip(uint(ui))
			}
		} else {
			_, ui = d.uintBytes()
			d.r.skip(uint(ui))
		}
	case cborMajorArray:
		if d.bd == cborBdIndefiniteArray {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				d.nextValueBytesBdReadR()
			}
		} else {
			_, ui = d.uintBytes()
			for i := uint64(0); i < ui; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		}
	case cborMajorMap:
		if d.bd == cborBdIndefiniteMap {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			_, ui = d.uintBytes()
			for i := uint64(0); i < ui; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		}
	case cborMajorTag:
		d.uintBytes()
		d.readNextBd()
		d.nextValueBytesBdReadR()
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined, cborBdFalse, cborBdTrue:
		case cborBdFloat16:
			d.r.skip(2)
		case cborBdFloat32:
			d.r.skip(4)
		case cborBdFloat64:
			d.r.skip(8)
		default:
			halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
	}
	return
}

func (d *cborDecDriverBytes) reset() {
	d.bdAndBdread.reset()

}

func (d *cborEncDriverBytes) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*CborHandle)
	d.e = shared
	if shared.bytes {
		fp = cborFpEncBytes
	} else {
		fp = cborFpEncIO
	}

	d.init2(enc)
	return
}

func (e *cborEncDriverBytes) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *cborEncDriverBytes) writerEnd() { e.w.end() }

func (e *cborEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *cborEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *cborDecDriverBytes) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*CborHandle)
	d.d = shared
	if shared.bytes {
		fp = cborFpDecBytes
	} else {
		fp = cborFpDecIO
	}

	d.init2(dec)
	return
}

func (d *cborDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *cborDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *cborDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *cborDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, cbordesc(d.bd))
}

func (d *cborDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

func (d *cborEncDriverBytes) init2(enc encoderI) {
	d.enc = enc
}

func (d *cborDecDriverBytes) init2(dec decoderI) {
	d.dec = dec

}

type helperEncDriverCborIO struct{}
type encFnCborIO struct {
	i  encFnInfo
	fe func(*encoderCborIO, *encFnInfo, reflect.Value)
}
type encRtidFnCborIO struct {
	rtid uintptr
	fn   *encFnCborIO
}
type encoderCborIO struct {
	dh helperEncDriverCborIO
	fp *fastpathEsCborIO
	e  cborEncDriverIO
	encoderBase
}
type helperDecDriverCborIO struct{}
type decFnCborIO struct {
	i  decFnInfo
	fd func(*decoderCborIO, *decFnInfo, reflect.Value)
}
type decRtidFnCborIO struct {
	rtid uintptr
	fn   *decFnCborIO
}
type decoderCborIO struct {
	dh helperDecDriverCborIO
	fp *fastpathDsCborIO
	d  cborDecDriverIO
	decoderBase
}
type cborEncDriverIO struct {
	noBuiltInTypes
	encDriverNoState
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT

	h   *CborHandle
	e   *encoderBase
	w   bufioEncWriter
	enc encoderI

	b [40]byte
}
type cborDecDriverIO struct {
	decDriverNoopContainerReader

	noBuiltInTypes

	h   *CborHandle
	d   *decoderBase
	r   ioDecReader
	dec decoderI
	bdAndBdread
}

func (e *encoderCborIO) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderCborIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderCborIO) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderCborIO) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderCborIO) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderCborIO) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderCborIO) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderCborIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderCborIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderCborIO) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderCborIO) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderCborIO) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderCborIO) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderCborIO) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderCborIO) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderCborIO) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderCborIO) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderCborIO) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderCborIO) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderCborIO) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderCborIO) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderCborIO) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderCborIO) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderCborIO) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderCborIO) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderCborIO) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderCborIO) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderCborIO) kSeqFn(rt reflect.Type) (fn *encFnCborIO) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderCborIO) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnCborIO
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

func (e *encoderCborIO) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnCborIO
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

func (e *encoderCborIO) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderCborIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderCborIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderCborIO) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderCborIO) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderCborIO) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderCborIO) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderCborIO) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnCborIO

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

func (e *encoderCborIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnCborIO) {
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

func (e *encoderCborIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsCborIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderCborIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderCborIO) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderCborIO) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderCborIO) mustEncode(v interface{}) {
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

func (e *encoderCborIO) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderCborIO) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderCborIO) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderCborIO) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderCborIO) encodeValue(rv reflect.Value, fn *encFnCborIO) {

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

func (e *encoderCborIO) encodeValueNonNil(rv reflect.Value, fn *encFnCborIO) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderCborIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderCborIO) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderCborIO) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderCborIO) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderCborIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderCborIO) fn(t reflect.Type) *encFnCborIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderCborIO) fnNoExt(t reflect.Type) *encFnCborIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderCborIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderCborIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderCborIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderCborIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderCborIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderCborIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderCborIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderCborIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverCborIO) newEncoderBytes(out *[]byte, h Handle) *encoderCborIO {
	var c1 encoderCborIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverCborIO) newEncoderIO(out io.Writer, h Handle) *encoderCborIO {
	var c1 encoderCborIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverCborIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsCborIO) (f *fastpathECborIO, u reflect.Type) {
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

func (helperEncDriverCborIO) encFindRtidFn(s []encRtidFnCborIO, rtid uintptr) (i uint, fn *encFnCborIO) {

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

func (helperEncDriverCborIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnCborIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnCborIO](v))
	}
	return
}

func (dh helperEncDriverCborIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsCborIO, checkExt bool) (fn *encFnCborIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverCborIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnCborIO = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverCborIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnCborIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnCborIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnCborIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnCborIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverCborIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnCborIO) {
	fn = new(encFnCborIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderCborIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderCborIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderCborIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderCborIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderCborIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderCborIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderCborIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderCborIO).textMarshal
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
					fn.fe = func(e *encoderCborIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderCborIO).kBool
			case reflect.String:

				fn.fe = (*encoderCborIO).kString
			case reflect.Int:
				fn.fe = (*encoderCborIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderCborIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderCborIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderCborIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderCborIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderCborIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderCborIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderCborIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderCborIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderCborIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderCborIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderCborIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderCborIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderCborIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderCborIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderCborIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderCborIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderCborIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderCborIO).kStructSimple
				} else {
					fn.fe = (*encoderCborIO).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderCborIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderCborIO).kErr
			default:

				fn.fe = (*encoderCborIO).kErr
			}
		}
	}
	return
}
func (d *decoderCborIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderCborIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderCborIO) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderCborIO) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderCborIO) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderCborIO) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderCborIO) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderCborIO) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderCborIO) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderCborIO) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderCborIO) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderCborIO) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderCborIO) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderCborIO) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderCborIO) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderCborIO) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderCborIO) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderCborIO) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderCborIO) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderCborIO) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderCborIO) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderCborIO) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborIO) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderCborIO) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderCborIO) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderCborIO) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderCborIO) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderCborIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderCborIO) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderCborIO) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderCborIO) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderCborIO) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderCborIO) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnCborIO

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

func (d *decoderCborIO) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnCborIO
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

func (d *decoderCborIO) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnCborIO

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

func (d *decoderCborIO) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnCborIO
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

func (d *decoderCborIO) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsCborIO)

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

func (d *decoderCborIO) reset() {
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

func (d *decoderCborIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderCborIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderCborIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderCborIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderCborIO) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderCborIO) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderCborIO) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderCborIO) Release() {}

func (d *decoderCborIO) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderCborIO) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderCborIO) decode(iv interface{}) {
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

func (d *decoderCborIO) decodeValue(rv reflect.Value, fn *decFnCborIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderCborIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnCborIO) {

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

func (d *decoderCborIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderCborIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderCborIO) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderCborIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderCborIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderCborIO) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderCborIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderCborIO) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderCborIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderCborIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborIO) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderCborIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderCborIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderCborIO) fn(t reflect.Type) *decFnCborIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderCborIO) fnNoExt(t reflect.Type) *decFnCborIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverCborIO) newDecoderBytes(in []byte, h Handle) *decoderCborIO {
	var c1 decoderCborIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverCborIO) newDecoderIO(in io.Reader, h Handle) *decoderCborIO {
	var c1 decoderCborIO
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverCborIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsCborIO) (f *fastpathDCborIO, u reflect.Type) {
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

func (helperDecDriverCborIO) decFindRtidFn(s []decRtidFnCborIO, rtid uintptr) (i uint, fn *decFnCborIO) {

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

func (helperDecDriverCborIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnCborIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnCborIO](v))
	}
	return
}

func (dh helperDecDriverCborIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsCborIO,
	checkExt bool) (fn *decFnCborIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverCborIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnCborIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverCborIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnCborIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnCborIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnCborIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnCborIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverCborIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsCborIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnCborIO) {
	fn = new(decFnCborIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderCborIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderCborIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderCborIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderCborIO).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderCborIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderCborIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderCborIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderCborIO).textUnmarshal
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
						fn.fd = func(d *decoderCborIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderCborIO, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderCborIO).kBool
			case reflect.String:
				fn.fd = (*decoderCborIO).kString
			case reflect.Int:
				fn.fd = (*decoderCborIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderCborIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderCborIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderCborIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderCborIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderCborIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderCborIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderCborIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderCborIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderCborIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderCborIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderCborIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderCborIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderCborIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderCborIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderCborIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderCborIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderCborIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderCborIO).kStructSimple
				} else {
					fn.fd = (*decoderCborIO).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderCborIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderCborIO).kInterface
			default:

				fn.fd = (*decoderCborIO).kErr
			}
		}
	}
	return
}
func (e *cborEncDriverIO) EncodeNil() {
	e.w.writen1(cborBdNil)
}

func (e *cborEncDriverIO) EncodeBool(b bool) {
	if b {
		e.w.writen1(cborBdTrue)
	} else {
		e.w.writen1(cborBdFalse)
	}
}

func (e *cborEncDriverIO) EncodeFloat32(f float32) {
	b := math.Float32bits(f)
	if e.h.OptimumSize {
		if h := floatToHalfFloatBits(b); halfFloatToFloatBits(h) == b {
			e.w.writen1(cborBdFloat16)
			e.w.writen2(bigen.PutUint16(h))
			return
		}
	}
	e.w.writen1(cborBdFloat32)
	e.w.writen4(bigen.PutUint32(b))
}

func (e *cborEncDriverIO) EncodeFloat64(f float64) {
	if e.h.OptimumSize {
		if f32 := float32(f); float64(f32) == f {
			e.EncodeFloat32(f32)
			return
		}
	}
	e.w.writen1(cborBdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *cborEncDriverIO) encUint(v uint64, bd byte) {
	if v <= 0x17 {
		e.w.writen1(byte(v) + bd)
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd+0x18, uint8(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd + 0x19)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd + 0x1a)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd + 0x1b)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *cborEncDriverIO) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-1-v), cborBaseNegInt)
	} else {
		e.encUint(uint64(v), cborBaseUint)
	}
}

func (e *cborEncDriverIO) EncodeUint(v uint64) {
	e.encUint(v, cborBaseUint)
}

func (e *cborEncDriverIO) encLen(bd byte, length int) {
	e.encUint(uint64(length), bd)
}

func (e *cborEncDriverIO) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else if e.h.TimeRFC3339 {
		e.encUint(0, cborBaseTag)
		e.encStringBytesS(cborBaseString, stringView(t.AppendFormat(e.b[:0], time.RFC3339Nano)))
	} else {
		e.encUint(1, cborBaseTag)
		t = t.UTC().Round(time.Microsecond)
		sec, nsec := t.Unix(), uint64(t.Nanosecond())
		if nsec == 0 {
			e.EncodeInt(sec)
		} else {
			e.EncodeFloat64(float64(sec) + float64(nsec)/1e9)
		}
	}
}

func (e *cborEncDriverIO) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	e.encUint(uint64(xtag), cborBaseTag)
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.writeNilBytes()
	} else {
		e.enc.encodeI(v)
	}
}

func (e *cborEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encUint(uint64(re.Tag), cborBaseTag)
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encodeI(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *cborEncDriverIO) WriteArrayEmpty() {
	if e.h.IndefiniteLength {
		e.w.writen2(cborBdIndefiniteArray, cborBdBreak)
	} else {
		e.w.writen1(cborBaseArray)

	}
}

func (e *cborEncDriverIO) WriteMapEmpty() {
	if e.h.IndefiniteLength {
		e.w.writen2(cborBdIndefiniteMap, cborBdBreak)
	} else {
		e.w.writen1(cborBaseMap)

	}
}

func (e *cborEncDriverIO) WriteArrayStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteArray)
	} else {
		e.encLen(cborBaseArray, length)
	}
}

func (e *cborEncDriverIO) WriteMapStart(length int) {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdIndefiniteMap)
	} else {
		e.encLen(cborBaseMap, length)
	}
}

func (e *cborEncDriverIO) WriteMapEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverIO) WriteArrayEnd() {
	if e.h.IndefiniteLength {
		e.w.writen1(cborBdBreak)
	}
}

func (e *cborEncDriverIO) EncodeString(v string) {
	bb := cborBaseString
	if e.h.StringToRaw {
		bb = cborBaseBytes
	}
	e.encStringBytesS(bb, v)
}

func (e *cborEncDriverIO) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *cborEncDriverIO) EncodeStringBytesRaw(v []byte) {
	e.encStringBytesS(cborBaseBytes, stringView(v))
}

func (e *cborEncDriverIO) encStringBytesS(bb byte, v string) {
	if e.h.IndefiniteLength {
		if bb == cborBaseBytes {
			e.w.writen1(cborBdIndefiniteBytes)
		} else {
			e.w.writen1(cborBdIndefiniteString)
		}
		vlen := uint(len(v))
		n := max(4, min(vlen/4, 1024))
		for i := uint(0); i < vlen; {
			i2 := i + n
			if i2 >= vlen {
				i2 = vlen
			}
			v2 := v[i:i2]
			e.encLen(bb, len(v2))
			e.w.writestr(v2)
			i = i2
		}
		e.w.writen1(cborBdBreak)
	} else {
		e.encLen(bb, len(v))
		e.w.writestr(v)
	}
}

func (e *cborEncDriverIO) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *cborEncDriverIO) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = cborBdNil
	}
	e.w.writen1(v)
}

func (e *cborEncDriverIO) writeNilArray() {
	e.writeNilOr(cborBaseArray)
}

func (e *cborEncDriverIO) writeNilMap() {
	e.writeNilOr(cborBaseMap)
}

func (e *cborEncDriverIO) writeNilBytes() {
	e.writeNilOr(cborBaseBytes)
}

func (d *cborDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *cborDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdNil || d.bd == cborBdUndefined {
		d.bdRead = false
		return true
	}
	return
}

func (d *cborDecDriverIO) TryNil() bool {
	return d.advanceNil()
}

func (d *cborDecDriverIO) skipTags() {
	for d.bd>>5 == cborMajorTag {
		d.decUint()
		d.bd = d.r.readn1()
	}
}

func (d *cborDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	if d.bd == cborBdNil {
		d.bdRead = false
		return valueTypeNil
	}
	major := d.bd >> 5
	if major == cborMajorBytes {
		return valueTypeBytes
	} else if major == cborMajorString {
		return valueTypeString
	} else if major == cborMajorArray {
		return valueTypeArray
	} else if major == cborMajorMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *cborDecDriverIO) CheckBreak() (v bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == cborBdBreak {
		d.bdRead = false
		v = true
	}
	return
}

func (d *cborDecDriverIO) decUint() (ui uint64) {
	v := d.bd & 0x1f
	if v <= 0x17 {
		ui = uint64(v)
	} else if v == 0x18 {
		ui = uint64(d.r.readn1())
	} else if v == 0x19 {
		ui = uint64(bigen.Uint16(d.r.readn2()))
	} else if v == 0x1a {
		ui = uint64(bigen.Uint32(d.r.readn4()))
	} else if v == 0x1b {
		ui = uint64(bigen.Uint64(d.r.readn8()))
	} else {
		halt.errorf("invalid descriptor decoding uint: %x/%s (%x)", d.bd, cbordesc(d.bd), v)
	}
	return
}

func (d *cborDecDriverIO) decLen() int {
	return int(d.decUint())
}

func (d *cborDecDriverIO) decFloat() (f float64, ok bool) {
	ok = true
	switch d.bd {
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.r.readn2()))))
	case cborBdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case cborBdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	default:
		if d.bd>>5 == cborMajorTag {

			switch d.bd & 0x1f {
			case 2:
				f = d.decTagBigIntAsFloat(false)
			case 3:
				f = d.decTagBigIntAsFloat(true)
			case 4:
				f = d.decTagBigFloatAsFloat(true)
			case 5:
				f = d.decTagBigFloatAsFloat(false)
			default:
				ok = false
			}
		} else {
			ok = false
		}
	}
	return
}

func (d *cborDecDriverIO) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	switch d.bd >> 5 {
	case cborMajorUint:
		ui = d.decUint()
	case cborMajorNegInt:
		ui = d.decUint()
		neg = true
	default:
		ok = false
	}
	return
}

func (d *cborDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, true)
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	if d.bd == cborBdTrue {
		b = true
	} else if d.bd == cborBdFalse {
	} else {
		halt.errorf("not bool - %s %x/%s", msgBadDesc, d.bd, cbordesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorMap {
		halt.errorf("error reading map; got major type: %x, expected %x/%s", d.bd>>5, cborMajorMap, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverIO) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray {
		return containerLenUnknown
	}
	if d.bd>>5 != cborMajorArray {
		halt.errorf("invalid array; got major type: %x, expect: %x/%s", d.bd>>5, cborMajorArray, cbordesc(d.bd))
	}
	return d.decLen()
}

func (d *cborDecDriverIO) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	if d.h.SkipUnexpectedTags {
		d.skipTags()
	}
	fnEnsureNonNilBytes := func() {

		if bs == nil {
			bs = zeroByteSlice
			state = dBytesDetach
		}
	}
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
		major := d.bd >> 5
		val4str := d.h.ValidateUnicode && major == cborMajorString
		bs = d.d.buf[:0]
		d.bdRead = false
		for !d.CheckBreak() {
			if d.bd>>5 != major {
				const msg = "malformed indefinite string/bytes %x (%s); " +
					"contains chunk with major type %v, expected %v"
				halt.errorf(msg, d.bd, cbordesc(d.bd), d.bd>>5, major)
			}
			n := uint(d.decLen())
			bs = append(bs, d.r.readx(n)...)
			d.bdRead = false
			if val4str && !utf8.Valid(bs[len(bs)-int(n):]) {
				const msg = "indefinite-length text string contains chunk " +
					"that is not a valid utf-8 sequence: 0x%x"
				halt.errorf(msg, bs[len(bs)-int(n):])
			}
		}
		d.bdRead = false
		d.d.buf = bs
		state = dBytesAttachBuffer
		fnEnsureNonNilBytes()
		return
	}
	if d.bd == cborBdIndefiniteArray {
		d.bdRead = false
		bs = d.d.buf[:0]
		for !d.CheckBreak() {
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		}
		d.d.buf = bs
		state = dBytesAttachBuffer
		fnEnsureNonNilBytes()
		return
	}
	var cond bool
	if d.bd>>5 == cborMajorArray {
		d.bdRead = false
		slen := d.decLen()
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
		fnEnsureNonNilBytes()
		return
	}
	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *cborDecDriverIO) DecodeStringAsBytes() (out []byte, state dBytesAttachState) {
	out, state = d.DecodeBytes()
	if d.h.ValidateUnicode && !utf8.Valid(out) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", out)
	}
	return
}

func (d *cborDecDriverIO) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
}

func (d *cborDecDriverIO) decodeTime(xtag uint64) (t time.Time) {
	switch xtag {
	case 0:
		var err error
		t, err = time.Parse(time.RFC3339, stringView(bytesOKs(d.DecodeStringAsBytes())))
		halt.onerror(err)
	case 1:
		f1, f2 := math.Modf(d.DecodeFloat64())
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		halt.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	}
	t = t.UTC().Round(time.Microsecond)
	return
}

func (d *cborDecDriverIO) preDecodeExt(checkTag bool, xtag uint64) (realxtag uint64, ok bool) {
	if d.advanceNil() {
		return
	}
	if d.bd>>5 != cborMajorTag {
		halt.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	}
	realxtag = d.decUint()
	d.bdRead = false
	if checkTag && xtag != realxtag {
		halt.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
	}
	ok = true
	return
}

func (d *cborDecDriverIO) DecodeRawExt(re *RawExt) {
	if realxtag, ok := d.preDecodeExt(false, 0); ok {
		re.Tag = realxtag
		d.dec.decode(&re.Value)
		d.bdRead = false
	}
}

func (d *cborDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if _, ok := d.preDecodeExt(true, xtag); ok {
		if ext == SelfExt {
			d.dec.decodeAs(rv, basetype, false)
		} else {
			d.dec.interfaceExtConvertAndDecode(rv, ext)
		}
		d.bdRead = false
	}
}

func (d *cborDecDriverIO) decTagBigIntAsFloat(neg bool) (f float64) {
	bs, _ := d.DecodeBytes()
	bi := new(big.Int).SetBytes(bs)
	if neg {
		bi0 := bi
		bi = new(big.Int).Sub(big.NewInt(-1), bi0)
	}
	f, _ = bi.Float64()
	return
}

func (d *cborDecDriverIO) decTagBigFloatAsFloat(decimal bool) (f float64) {
	if nn := d.r.readn1(); nn != 82 {
		halt.errorf("(%d) decoding decimal/big.Float: expected 2 numbers", nn)
	}
	exp := d.DecodeInt64()
	mant := d.DecodeInt64()
	if decimal {

		rf := readFloatResult{exp: int8(exp)}
		if mant >= 0 {
			rf.mantissa = uint64(mant)
		} else {
			rf.neg = true
			rf.mantissa = uint64(-mant)
		}
		f, _ = parseFloat64_reader(rf)

	} else {

		bfm := new(big.Float).SetPrec(64).SetInt64(mant)
		bf := new(big.Float).SetPrec(64).SetMantExp(bfm, int(exp))
		f, _ = bf.Float64()
	}
	return
}

func (d *cborDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool
	switch d.bd >> 5 {
	case cborMajorUint:
		if d.h.SignedInteger {
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		} else {
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		}
	case cborMajorNegInt:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case cborMajorBytes:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case cborMajorString:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case cborMajorArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborMajorMap:
		n.v = valueTypeMap
		decodeFurther = true
	case cborMajorTag:
		n.v = valueTypeExt
		n.u = d.decUint()
		d.bdRead = false
		n.l = nil
		xx := d.h.getExtForTag(n.u)
		if xx == nil {
			switch n.u {
			case 0, 1:
				n.v = valueTypeTime
				n.t = d.decodeTime(n.u)
			case 2:
				n.f = d.decTagBigIntAsFloat(false)
				n.v = valueTypeFloat
			case 3:
				n.f = d.decTagBigIntAsFloat(true)
				n.v = valueTypeFloat
			case 4:
				n.f = d.decTagBigFloatAsFloat(true)
				n.v = valueTypeFloat
			case 5:
				n.f = d.decTagBigFloatAsFloat(false)
				n.v = valueTypeFloat
			case 55799:
				d.DecodeNaked()
			default:
				if d.h.SkipUnexpectedTags {
					d.DecodeNaked()
				}

			}
			return
		}

	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined:
			n.v = valueTypeNil
		case cborBdFalse:
			n.v = valueTypeBool
			n.b = false
		case cborBdTrue:
			n.v = valueTypeBool
			n.b = true
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			n.v = valueTypeFloat
			n.f = d.DecodeFloat64()
		default:
			halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	}
	if !decodeFurther {
		d.bdRead = false
	}
}

func (d *cborDecDriverIO) uintBytes() (v []byte, ui uint64) {

	switch vv := d.bd & 0x1f; vv {
	case 0x18:
		v = d.r.readx(1)
		ui = uint64(v[0])
	case 0x19:
		v = d.r.readx(2)
		ui = uint64(bigenstd.Uint16(v))
	case 0x1a:
		v = d.r.readx(4)
		ui = uint64(bigenstd.Uint32(v))
	case 0x1b:
		v = d.r.readx(8)
		ui = uint64(bigenstd.Uint64(v))
	default:
		if vv > 0x1b {
			halt.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
		}
		ui = uint64(vv)
	}
	return
}

func (d *cborDecDriverIO) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *cborDecDriverIO) nextValueBytesBdReadR() {

	var ui uint64

	switch d.bd >> 5 {
	case cborMajorUint, cborMajorNegInt:
		d.uintBytes()
	case cborMajorString, cborMajorBytes:
		if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				_, ui = d.uintBytes()
				d.r.skip(uint(ui))
			}
		} else {
			_, ui = d.uintBytes()
			d.r.skip(uint(ui))
		}
	case cborMajorArray:
		if d.bd == cborBdIndefiniteArray {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				d.nextValueBytesBdReadR()
			}
		} else {
			_, ui = d.uintBytes()
			for i := uint64(0); i < ui; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		}
	case cborMajorMap:
		if d.bd == cborBdIndefiniteMap {
			for {
				d.readNextBd()
				if d.bd == cborBdBreak {
					break
				}
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		} else {
			_, ui = d.uintBytes()
			for i := uint64(0); i < ui; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		}
	case cborMajorTag:
		d.uintBytes()
		d.readNextBd()
		d.nextValueBytesBdReadR()
	case cborMajorSimpleOrFloat:
		switch d.bd {
		case cborBdNil, cborBdUndefined, cborBdFalse, cborBdTrue:
		case cborBdFloat16:
			d.r.skip(2)
		case cborBdFloat32:
			d.r.skip(4)
		case cborBdFloat64:
			d.r.skip(8)
		default:
			halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
		}
	default:
		halt.errorf("nextValueBytes: Unrecognized d.bd: 0x%x", d.bd)
	}
	return
}

func (d *cborDecDriverIO) reset() {
	d.bdAndBdread.reset()

}

func (d *cborEncDriverIO) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*CborHandle)
	d.e = shared
	if shared.bytes {
		fp = cborFpEncBytes
	} else {
		fp = cborFpEncIO
	}

	d.init2(enc)
	return
}

func (e *cborEncDriverIO) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *cborEncDriverIO) writerEnd() { e.w.end() }

func (e *cborEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *cborEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *cborDecDriverIO) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*CborHandle)
	d.d = shared
	if shared.bytes {
		fp = cborFpDecBytes
	} else {
		fp = cborFpDecIO
	}

	d.init2(dec)
	return
}

func (d *cborDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *cborDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *cborDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *cborDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, cbordesc(d.bd))
}

func (d *cborDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

func (d *cborEncDriverIO) init2(enc encoderI) {
	d.enc = enc
}

func (d *cborDecDriverIO) init2(dec decoderI) {
	d.dec = dec

}
