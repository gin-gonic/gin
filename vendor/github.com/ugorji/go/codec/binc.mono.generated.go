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
	"unicode/utf8"
)

type helperEncDriverBincBytes struct{}
type encFnBincBytes struct {
	i  encFnInfo
	fe func(*encoderBincBytes, *encFnInfo, reflect.Value)
}
type encRtidFnBincBytes struct {
	rtid uintptr
	fn   *encFnBincBytes
}
type encoderBincBytes struct {
	dh helperEncDriverBincBytes
	fp *fastpathEsBincBytes
	e  bincEncDriverBytes
	encoderBase
}
type helperDecDriverBincBytes struct{}
type decFnBincBytes struct {
	i  decFnInfo
	fd func(*decoderBincBytes, *decFnInfo, reflect.Value)
}
type decRtidFnBincBytes struct {
	rtid uintptr
	fn   *decFnBincBytes
}
type decoderBincBytes struct {
	dh helperDecDriverBincBytes
	fp *fastpathDsBincBytes
	d  bincDecDriverBytes
	decoderBase
}
type bincEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT
	encInit2er

	h *BincHandle
	e *encoderBase
	w bytesEncAppender
	bincEncState
}
type bincDecDriverBytes struct {
	decDriverNoopContainerReader

	decInit2er
	noBuiltInTypes

	h *BincHandle
	d *decoderBase
	r bytesDecReader

	bincDecState
}

func (e *encoderBincBytes) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderBincBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderBincBytes) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderBincBytes) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderBincBytes) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderBincBytes) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderBincBytes) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderBincBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderBincBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderBincBytes) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderBincBytes) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderBincBytes) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderBincBytes) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderBincBytes) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderBincBytes) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderBincBytes) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderBincBytes) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderBincBytes) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderBincBytes) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderBincBytes) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderBincBytes) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderBincBytes) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderBincBytes) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderBincBytes) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderBincBytes) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderBincBytes) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderBincBytes) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderBincBytes) kSeqFn(rt reflect.Type) (fn *encFnBincBytes) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderBincBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnBincBytes
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

func (e *encoderBincBytes) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnBincBytes
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

func (e *encoderBincBytes) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderBincBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderBincBytes) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderBincBytes) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderBincBytes) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincBytes) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincBytes) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnBincBytes

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

func (e *encoderBincBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnBincBytes) {
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

func (e *encoderBincBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsBincBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderBincBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderBincBytes) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderBincBytes) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderBincBytes) mustEncode(v interface{}) {
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

func (e *encoderBincBytes) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderBincBytes) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderBincBytes) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderBincBytes) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderBincBytes) encodeValue(rv reflect.Value, fn *encFnBincBytes) {

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

func (e *encoderBincBytes) encodeValueNonNil(rv reflect.Value, fn *encFnBincBytes) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderBincBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderBincBytes) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderBincBytes) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderBincBytes) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderBincBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderBincBytes) fn(t reflect.Type) *encFnBincBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderBincBytes) fnNoExt(t reflect.Type) *encFnBincBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderBincBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderBincBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderBincBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderBincBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderBincBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderBincBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderBincBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderBincBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverBincBytes) newEncoderBytes(out *[]byte, h Handle) *encoderBincBytes {
	var c1 encoderBincBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverBincBytes) newEncoderIO(out io.Writer, h Handle) *encoderBincBytes {
	var c1 encoderBincBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverBincBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsBincBytes) (f *fastpathEBincBytes, u reflect.Type) {
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

func (helperEncDriverBincBytes) encFindRtidFn(s []encRtidFnBincBytes, rtid uintptr) (i uint, fn *encFnBincBytes) {

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

func (helperEncDriverBincBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnBincBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnBincBytes](v))
	}
	return
}

func (dh helperEncDriverBincBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsBincBytes, checkExt bool) (fn *encFnBincBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverBincBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnBincBytes = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverBincBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnBincBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnBincBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnBincBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnBincBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverBincBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincBytes) {
	fn = new(encFnBincBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderBincBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderBincBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderBincBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderBincBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderBincBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderBincBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderBincBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderBincBytes).textMarshal
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
					fn.fe = func(e *encoderBincBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderBincBytes).kBool
			case reflect.String:

				fn.fe = (*encoderBincBytes).kString
			case reflect.Int:
				fn.fe = (*encoderBincBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderBincBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderBincBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderBincBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderBincBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderBincBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderBincBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderBincBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderBincBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderBincBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderBincBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderBincBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderBincBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderBincBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderBincBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderBincBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderBincBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderBincBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderBincBytes).kStructSimple
				} else {
					fn.fe = (*encoderBincBytes).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderBincBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderBincBytes).kErr
			default:

				fn.fe = (*encoderBincBytes).kErr
			}
		}
	}
	return
}
func (d *decoderBincBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderBincBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderBincBytes) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderBincBytes) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderBincBytes) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderBincBytes) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderBincBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderBincBytes) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderBincBytes) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderBincBytes) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderBincBytes) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderBincBytes) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderBincBytes) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderBincBytes) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderBincBytes) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderBincBytes) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderBincBytes) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderBincBytes) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderBincBytes) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderBincBytes) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderBincBytes) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderBincBytes) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincBytes) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincBytes) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderBincBytes) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderBincBytes) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderBincBytes) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderBincBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderBincBytes) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderBincBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderBincBytes) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderBincBytes) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderBincBytes) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnBincBytes

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

func (d *decoderBincBytes) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnBincBytes
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

func (d *decoderBincBytes) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnBincBytes

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

func (d *decoderBincBytes) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnBincBytes
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

func (d *decoderBincBytes) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsBincBytes)

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

func (d *decoderBincBytes) reset() {
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

func (d *decoderBincBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderBincBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderBincBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderBincBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderBincBytes) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderBincBytes) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderBincBytes) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderBincBytes) Release() {}

func (d *decoderBincBytes) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderBincBytes) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderBincBytes) decode(iv interface{}) {
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

func (d *decoderBincBytes) decodeValue(rv reflect.Value, fn *decFnBincBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderBincBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnBincBytes) {

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

func (d *decoderBincBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderBincBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderBincBytes) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderBincBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderBincBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderBincBytes) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderBincBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderBincBytes) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderBincBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderBincBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincBytes) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderBincBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderBincBytes) fn(t reflect.Type) *decFnBincBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderBincBytes) fnNoExt(t reflect.Type) *decFnBincBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverBincBytes) newDecoderBytes(in []byte, h Handle) *decoderBincBytes {
	var c1 decoderBincBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverBincBytes) newDecoderIO(in io.Reader, h Handle) *decoderBincBytes {
	var c1 decoderBincBytes
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverBincBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsBincBytes) (f *fastpathDBincBytes, u reflect.Type) {
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

func (helperDecDriverBincBytes) decFindRtidFn(s []decRtidFnBincBytes, rtid uintptr) (i uint, fn *decFnBincBytes) {

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

func (helperDecDriverBincBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnBincBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnBincBytes](v))
	}
	return
}

func (dh helperDecDriverBincBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsBincBytes,
	checkExt bool) (fn *decFnBincBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverBincBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnBincBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverBincBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnBincBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnBincBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnBincBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnBincBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverBincBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsBincBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincBytes) {
	fn = new(decFnBincBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderBincBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderBincBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderBincBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderBincBytes).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderBincBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderBincBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderBincBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderBincBytes).textUnmarshal
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
						fn.fd = func(d *decoderBincBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderBincBytes, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderBincBytes).kBool
			case reflect.String:
				fn.fd = (*decoderBincBytes).kString
			case reflect.Int:
				fn.fd = (*decoderBincBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderBincBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderBincBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderBincBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderBincBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderBincBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderBincBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderBincBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderBincBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderBincBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderBincBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderBincBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderBincBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderBincBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderBincBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderBincBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderBincBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderBincBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderBincBytes).kStructSimple
				} else {
					fn.fd = (*decoderBincBytes).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderBincBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderBincBytes).kInterface
			default:

				fn.fd = (*decoderBincBytes).kErr
			}
		}
	}
	return
}
func (e *bincEncDriverBytes) EncodeNil() {
	e.w.writen1(bincBdNil)
}

func (e *bincEncDriverBytes) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else {
		bs := bincEncodeTime(t)
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	}
}

func (e *bincEncDriverBytes) EncodeBool(b bool) {
	if b {
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	}
}

func (e *bincEncDriverBytes) encSpFloat(f float64) (done bool) {
	if f == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
	} else if math.IsNaN(float64(f)) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNan)
	} else if math.IsInf(float64(f), +1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpPosInf)
	} else if math.IsInf(float64(f), -1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegInf)
	} else {
		return
	}
	return true
}

func (e *bincEncDriverBytes) EncodeFloat32(f float32) {
	if !e.encSpFloat(float64(f)) {
		e.w.writen1(bincVdFloat<<4 | bincFlBin32)
		e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
	}
}

func (e *bincEncDriverBytes) EncodeFloat64(f float64) {
	if e.encSpFloat(f) {
		return
	}
	b := bigen.PutUint64(math.Float64bits(f))
	if bincDoPrune {
		i := 7
		for ; i >= 0 && (b[i] == 0); i-- {
		}
		i++
		if i <= 6 {
			e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.w.writen1(byte(i))
			e.w.writeb(b[:i])
			return
		}
	}
	e.w.writen1(bincVdFloat<<4 | bincFlBin64)
	e.w.writen8(b)
}

func (e *bincEncDriverBytes) encIntegerPrune32(bd byte, pos bool, v uint64) {
	b := bigen.PutUint32(uint32(v))
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 3 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 3)
		e.w.writen4(b)
	}
}

func (e *bincEncDriverBytes) encIntegerPrune64(bd byte, pos bool, v uint64) {
	b := bigen.PutUint64(v)
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 7 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 7)
		e.w.writen8(b)
	}
}

func (e *bincEncDriverBytes) EncodeInt(v int64) {
	if v >= 0 {
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	} else if v == -1 {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	} else {
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	}
}

func (e *bincEncDriverBytes) EncodeUint(v uint64) {
	e.encUint(bincVdPosInt<<4, true, v)
}

func (e *bincEncDriverBytes) encUint(bd byte, pos bool, v uint64) {
	if v == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	} else if pos && v >= 1 && v <= 16 {
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.encIntegerPrune32(bd, pos, v)
	} else {
		e.encIntegerPrune64(bd, pos, v)
	}
}

func (e *bincEncDriverBytes) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *bincEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *bincEncDriverBytes) encodeExtPreamble(xtag byte, length int) {
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncDriverBytes) WriteArrayStart(length int) {
	e.encLen(bincVdArray<<4, uint64(length))
}

func (e *bincEncDriverBytes) WriteMapStart(length int) {
	e.encLen(bincVdMap<<4, uint64(length))
}

func (e *bincEncDriverBytes) WriteArrayEmpty() {

	e.w.writen1(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriverBytes) WriteMapEmpty() {

	e.w.writen1(bincVdMap<<4 | uint8(0+4))
}

func (e *bincEncDriverBytes) EncodeSymbol(v string) {

	l := len(v)
	if l == 0 {
		e.encBytesLen(cUTF8, 0)
		return
	} else if l == 1 {
		e.encBytesLen(cUTF8, 1)
		e.w.writen1(v[0])
		return
	}
	if e.m == nil {
		e.m = make(map[string]uint16, 16)
	}
	ui, ok := e.m[v]
	if ok {
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			e.w.writen2(bigen.PutUint16(ui))
		}
	} else {
		e.e.seq++
		ui = e.e.seq
		e.m[v] = ui
		var lenprec uint8
		if l <= math.MaxUint8 {

		} else if l <= math.MaxUint16 {
			lenprec = 1
		} else if int64(l) <= math.MaxUint32 {
			lenprec = 2
		} else {
			lenprec = 3
		}
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4|0x4|lenprec, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			e.w.writen2(bigen.PutUint16(ui))
		}
		if lenprec == 0 {
			e.w.writen1(byte(l))
		} else if lenprec == 1 {
			e.w.writen2(bigen.PutUint16(uint16(l)))
		} else if lenprec == 2 {
			e.w.writen4(bigen.PutUint32(uint32(l)))
		} else {
			e.w.writen8(bigen.PutUint64(uint64(l)))
		}
		e.w.writestr(v)
	}
}

func (e *bincEncDriverBytes) EncodeString(v string) {
	if e.h.StringToRaw {
		e.encLen(bincVdByteArray<<4, uint64(len(v)))
		if len(v) > 0 {
			e.w.writestr(v)
		}
		return
	}
	e.EncodeStringEnc(cUTF8, v)
}

func (e *bincEncDriverBytes) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *bincEncDriverBytes) EncodeStringEnc(c charEncoding, v string) {
	if e.e.c == containerMapKey && c == cUTF8 && (e.h.AsSymbols == 1) {
		e.EncodeSymbol(v)
		return
	}
	e.encLen(bincVdString<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writestr(v)
	}
}

func (e *bincEncDriverBytes) EncodeStringBytesRaw(v []byte) {
	e.encLen(bincVdByteArray<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writeb(v)
	}
}

func (e *bincEncDriverBytes) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *bincEncDriverBytes) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = bincBdNil
	}
	e.w.writen1(v)
}

func (e *bincEncDriverBytes) writeNilArray() {
	e.writeNilOr(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriverBytes) writeNilMap() {
	e.writeNilOr(bincVdMap<<4 | uint8(0+4))
}

func (e *bincEncDriverBytes) writeNilBytes() {
	e.writeNilOr(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriverBytes) encBytesLen(c charEncoding, length uint64) {

	if c == cRAW {
		e.encLen(bincVdByteArray<<4, length)
	} else {
		e.encLen(bincVdString<<4, length)
	}
}

func (e *bincEncDriverBytes) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l+4))
	} else {
		e.encLenNumber(bd, l)
	}
}

func (e *bincEncDriverBytes) encLenNumber(bd byte, v uint64) {
	if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd | 0x02)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd | 0x03)
		e.w.writen8(bigen.PutUint64(uint64(v)))
	}
}

func (d *bincDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
}

func (d *bincDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *bincDecDriverBytes) TryNil() bool {
	return d.advanceNil()
}

func (d *bincDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return valueTypeNil
	} else if d.vd == bincVdByteArray {
		return valueTypeBytes
	} else if d.vd == bincVdString {
		return valueTypeString
	} else if d.vd == bincVdArray {
		return valueTypeArray
	} else if d.vd == bincVdMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *bincDecDriverBytes) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.vd != bincVdTimestamp {
		halt.errorf("cannot decode time - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	t, err := bincDecodeTime(d.r.readx(uint(d.vs)))
	halt.onerror(err)
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) decFloatPruned(maxlen uint8) {
	l := d.r.readn1()
	if l > maxlen {
		halt.errorf("cannot read float - at most %v bytes used to represent float - received %v bytes", maxlen, l)
	}
	for i := l; i < maxlen; i++ {
		d.d.b[i] = 0
	}
	d.r.readb(d.d.b[0:l])
}

func (d *bincDecDriverBytes) decFloatPre32() (b [4]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn4()
	} else {
		d.decFloatPruned(4)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverBytes) decFloatPre64() (b [8]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn8()
	} else {
		d.decFloatPruned(8)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverBytes) decFloatVal() (f float64) {
	switch d.vs & 0x7 {
	case bincFlBin32:
		f = float64(math.Float32frombits(bigen.Uint32(d.decFloatPre32())))
	case bincFlBin64:
		f = math.Float64frombits(bigen.Uint64(d.decFloatPre64()))
	default:

		halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincDecDriverBytes) decUint() (v uint64) {
	switch d.vs {
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		v = uint64(bigen.Uint16(d.r.readn2()))
	case 2:
		b3 := d.r.readn3()
		var b [4]byte
		copy(b[1:], b3[:])
		v = uint64(bigen.Uint32(b))
	case 3:
		v = uint64(bigen.Uint32(d.r.readn4()))
	case 4, 5, 6:

		bs := d.d.b[:8]
		clear(bs)
		d.r.readb(bs[(7 - d.vs):])
		v = bigen.Uint64(*(*[8]byte)(bs))
	case 7:
		v = bigen.Uint64(d.r.readn8())
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	return
}

func (d *bincDecDriverBytes) uintBytes() (bs []byte) {
	switch d.vs {
	case 0:
		bs = d.d.b[:1]
		bs[0] = d.r.readn1()
		return
	case 1:
		bs = d.d.b[:2]
	case 2:
		bs = d.d.b[:3]
	case 3:
		bs = d.d.b[:4]
	case 4, 5, 6:
		lim := 7 - d.vs
		bs = d.d.b[lim:8]
	case 7:
		bs = d.d.b[:8]
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	d.r.readb(bs)
	return
}

func (d *bincDecDriverBytes) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdPosInt {
		ui = d.decUint()
	} else if vd == bincVdNegInt {
		ui = d.decUint()
		neg = true
	} else if vd == bincVdSmallInt {
		ui = uint64(d.vs) + 1
	} else if vd == bincVdSpecial {
		if vs == bincSpZero {

		} else if vs == bincSpNegOne {
			neg = true
			ui = 1
		} else {
			ok = false

		}
	} else {
		ok = false

	}
	return
}

func (d *bincDecDriverBytes) decFloat() (f float64, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdSpecial {
		if vs == bincSpNan {
			f = math.NaN()
		} else if vs == bincSpPosInf {
			f = math.Inf(1)
		} else if vs == bincSpZeroFloat || vs == bincSpZero {

		} else if vs == bincSpNegInf {
			f = math.Inf(-1)
		} else {
			ok = false

		}
	} else if vd == bincVdFloat {
		f = d.decFloatVal()
	} else {
		ok = false
	}
	return
}

func (d *bincDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == (bincVdSpecial | bincSpFalse) {

	} else if d.bd == (bincVdSpecial | bincSpTrue) {
		b = true
	} else {
		halt.errorf("bool - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdMap {
		halt.errorf("map - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdArray {
		halt.errorf("array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) decLen() int {
	if d.vs > 3 {
		return int(d.vs - 4)
	}
	return int(d.decLenNumber())
}

func (d *bincDecDriverBytes) decLenNumber() (v uint64) {
	if x := d.vs; x == 0 {
		v = uint64(d.r.readn1())
	} else if x == 1 {
		v = uint64(bigen.Uint16(d.r.readn2()))
	} else if x == 2 {
		v = uint64(bigen.Uint32(d.r.readn4()))
	} else {
		v = bigen.Uint64(d.r.readn8())
	}
	return
}

func (d *bincDecDriverBytes) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool
	var slen = -1
	switch d.vd {
	case bincVdString, bincVdByteArray:
		slen = d.decLen()
		bs, cond = d.r.readxb(uint(slen))
		state = d.d.attachState(cond)
	case bincVdSymbol:

		var symbol uint16
		vs := d.vs
		if vs&0x8 == 0 {
			symbol = uint16(d.r.readn1())
		} else {
			symbol = uint16(bigen.Uint16(d.r.readn2()))
		}
		if d.s == nil {
			d.s = make(map[uint16][]byte, 16)
		}

		if vs&0x4 == 0 {
			bs = d.s[symbol]
		} else {
			switch vs & 0x3 {
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(bigen.Uint16(d.r.readn2()))
			case 2:
				slen = int(bigen.Uint32(d.r.readn4()))
			case 3:
				slen = int(bigen.Uint64(d.r.readn8()))
			}

			bs, cond = d.r.readxb(uint(slen))
			bs = d.d.detach2Bytes(bs, d.d.attachState(cond))
			d.s[symbol] = bs
		}
		state = dBytesDetach
	default:
		halt.errorf("string/bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if d.h.ValidateUnicode && !utf8.Valid(bs) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", bs)
	}

	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool
	if d.vd == bincVdArray {
		slen := d.ReadArrayStart()
		bs, cond = usableByteSlice(d.d.buf, slen)
		for i := 0; i < slen; i++ {
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
	if !(d.vd == bincVdString || d.vd == bincVdByteArray) {
		halt.errorf("bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *bincDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *bincDecDriverBytes) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *bincDecDriverBytes) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
	if xtagIn > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtagIn)
	}
	if d.advanceNil() {
		return
	}
	tag := uint8(xtagIn)
	if d.vd == bincVdCustomExt {
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag - got %b, expecting: %v", xtag, tag)
		}
		xbs, ok = d.r.readxb(uint(l))
		bstate = d.d.attachState(ok)

	} else if d.vd == bincVdByteArray {
		xbs, bstate = d.DecodeBytes()
	} else {
		halt.errorf("ext expects extensions or byte array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	ok = true
	return
}

func (d *bincDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil:
			n.v = valueTypeNil
		case bincSpFalse:
			n.v = valueTypeBool
			n.b = false
		case bincSpTrue:
			n.v = valueTypeBool
			n.b = true
		case bincSpNan:
			n.v = valueTypeFloat
			n.f = math.NaN()
		case bincSpPosInf:
			n.v = valueTypeFloat
			n.f = math.Inf(1)
		case bincSpNegInf:
			n.v = valueTypeFloat
			n.f = math.Inf(-1)
		case bincSpZeroFloat:
			n.v = valueTypeFloat
			n.f = float64(0)
		case bincSpZero:
			n.v = valueTypeUint
			n.u = uint64(0)
		case bincSpNegOne:
			n.v = valueTypeInt
			n.i = int64(-1)
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
		n.v = valueTypeUint
		n.u = uint64(int8(d.vs)) + 1
	case bincVdPosInt:
		n.v = valueTypeUint
		n.u = d.decUint()
	case bincVdNegInt:
		n.v = valueTypeInt
		n.i = -(int64(d.decUint()))
	case bincVdFloat:
		n.v = valueTypeFloat
		n.f = d.decFloatVal()
	case bincVdString:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case bincVdByteArray:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case bincVdSymbol:
		n.v = valueTypeSymbol
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case bincVdTimestamp:
		n.v = valueTypeTime
		tt, err := bincDecodeTime(d.r.readx(uint(d.vs)))
		halt.onerror(err)
		n.t = tt
	case bincVdCustomExt:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(uint(l))
	case bincVdArray:
		n.v = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if !decodeFurther {
		d.bdRead = false
	}
	if n.v == valueTypeUint && d.h.SignedInteger {
		n.v = valueTypeInt
		n.i = int64(n.u)
	}
}

func (d *bincDecDriverBytes) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *bincDecDriverBytes) nextValueBytesBdReadR() {
	fnLen := func(vs byte) uint {
		switch vs {
		case 0:
			x := d.r.readn1()
			return uint(x)
		case 1:
			x := d.r.readn2()
			return uint(bigen.Uint16(x))
		case 2:
			x := d.r.readn4()
			return uint(bigen.Uint32(x))
		case 3:
			x := d.r.readn8()
			return uint(bigen.Uint64(x))
		default:
			return uint(vs - 4)
		}
	}

	var clen uint

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil, bincSpFalse, bincSpTrue, bincSpNan, bincSpPosInf:
		case bincSpNegInf, bincSpZeroFloat, bincSpZero, bincSpNegOne:
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
	case bincVdPosInt, bincVdNegInt:
		d.uintBytes()
	case bincVdFloat:
		fn := func(xlen byte) {
			if d.vs&0x8 != 0 {
				xlen = d.r.readn1()
				if xlen > 8 {
					halt.errorf("cannot read float - at most 8 bytes used to represent float - received %v bytes", xlen)
				}
			}
			d.r.readb(d.d.b[:xlen])
		}
		switch d.vs & 0x7 {
		case bincFlBin32:
			fn(4)
		case bincFlBin64:
			fn(8)
		default:
			halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdString, bincVdByteArray:
		clen = fnLen(d.vs)
		d.r.skip(clen)
	case bincVdSymbol:
		if d.vs&0x8 == 0 {
			d.r.readn1()
		} else {
			d.r.skip(2)
		}
		if d.vs&0x4 != 0 {
			clen = fnLen(d.vs & 0x3)
			d.r.skip(clen)
		}
	case bincVdTimestamp:
		d.r.skip(uint(d.vs))
	case bincVdCustomExt:
		clen = fnLen(d.vs)
		d.r.readn1()
		d.r.skip(clen)
	case bincVdArray:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case bincVdMap:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincEncDriverBytes) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*BincHandle)
	d.e = shared
	if shared.bytes {
		fp = bincFpEncBytes
	} else {
		fp = bincFpEncIO
	}

	d.init2(enc)
	return
}

func (e *bincEncDriverBytes) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *bincEncDriverBytes) writerEnd() { e.w.end() }

func (e *bincEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *bincEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *bincDecDriverBytes) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*BincHandle)
	d.d = shared
	if shared.bytes {
		fp = bincFpDecBytes
	} else {
		fp = bincFpDecIO
	}

	d.init2(dec)
	return
}

func (d *bincDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *bincDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *bincDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *bincDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, bincdescbd(d.bd))
}

func (d *bincDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

type helperEncDriverBincIO struct{}
type encFnBincIO struct {
	i  encFnInfo
	fe func(*encoderBincIO, *encFnInfo, reflect.Value)
}
type encRtidFnBincIO struct {
	rtid uintptr
	fn   *encFnBincIO
}
type encoderBincIO struct {
	dh helperEncDriverBincIO
	fp *fastpathEsBincIO
	e  bincEncDriverIO
	encoderBase
}
type helperDecDriverBincIO struct{}
type decFnBincIO struct {
	i  decFnInfo
	fd func(*decoderBincIO, *decFnInfo, reflect.Value)
}
type decRtidFnBincIO struct {
	rtid uintptr
	fn   *decFnBincIO
}
type decoderBincIO struct {
	dh helperDecDriverBincIO
	fp *fastpathDsBincIO
	d  bincDecDriverIO
	decoderBase
}
type bincEncDriverIO struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT
	encInit2er

	h *BincHandle
	e *encoderBase
	w bufioEncWriter
	bincEncState
}
type bincDecDriverIO struct {
	decDriverNoopContainerReader

	decInit2er
	noBuiltInTypes

	h *BincHandle
	d *decoderBase
	r ioDecReader

	bincDecState
}

func (e *encoderBincIO) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderBincIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderBincIO) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderBincIO) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderBincIO) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderBincIO) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderBincIO) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderBincIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderBincIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderBincIO) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderBincIO) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderBincIO) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderBincIO) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderBincIO) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderBincIO) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderBincIO) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderBincIO) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderBincIO) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderBincIO) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderBincIO) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderBincIO) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderBincIO) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderBincIO) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderBincIO) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderBincIO) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderBincIO) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderBincIO) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderBincIO) kSeqFn(rt reflect.Type) (fn *encFnBincIO) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderBincIO) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnBincIO
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

func (e *encoderBincIO) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnBincIO
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

func (e *encoderBincIO) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderBincIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderBincIO) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderBincIO) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderBincIO) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincIO) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderBincIO) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnBincIO

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

func (e *encoderBincIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnBincIO) {
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

func (e *encoderBincIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsBincIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderBincIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderBincIO) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderBincIO) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderBincIO) mustEncode(v interface{}) {
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

func (e *encoderBincIO) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderBincIO) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderBincIO) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderBincIO) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderBincIO) encodeValue(rv reflect.Value, fn *encFnBincIO) {

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

func (e *encoderBincIO) encodeValueNonNil(rv reflect.Value, fn *encFnBincIO) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderBincIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderBincIO) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderBincIO) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderBincIO) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderBincIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderBincIO) fn(t reflect.Type) *encFnBincIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderBincIO) fnNoExt(t reflect.Type) *encFnBincIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderBincIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderBincIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderBincIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderBincIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderBincIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderBincIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderBincIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderBincIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverBincIO) newEncoderBytes(out *[]byte, h Handle) *encoderBincIO {
	var c1 encoderBincIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverBincIO) newEncoderIO(out io.Writer, h Handle) *encoderBincIO {
	var c1 encoderBincIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverBincIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsBincIO) (f *fastpathEBincIO, u reflect.Type) {
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

func (helperEncDriverBincIO) encFindRtidFn(s []encRtidFnBincIO, rtid uintptr) (i uint, fn *encFnBincIO) {

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

func (helperEncDriverBincIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnBincIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnBincIO](v))
	}
	return
}

func (dh helperEncDriverBincIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsBincIO, checkExt bool) (fn *encFnBincIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverBincIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnBincIO = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverBincIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnBincIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnBincIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnBincIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnBincIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverBincIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnBincIO) {
	fn = new(encFnBincIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderBincIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderBincIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderBincIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderBincIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderBincIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderBincIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderBincIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderBincIO).textMarshal
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
					fn.fe = func(e *encoderBincIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderBincIO).kBool
			case reflect.String:

				fn.fe = (*encoderBincIO).kString
			case reflect.Int:
				fn.fe = (*encoderBincIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderBincIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderBincIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderBincIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderBincIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderBincIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderBincIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderBincIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderBincIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderBincIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderBincIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderBincIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderBincIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderBincIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderBincIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderBincIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderBincIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderBincIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderBincIO).kStructSimple
				} else {
					fn.fe = (*encoderBincIO).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderBincIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderBincIO).kErr
			default:

				fn.fe = (*encoderBincIO).kErr
			}
		}
	}
	return
}
func (d *decoderBincIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderBincIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderBincIO) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderBincIO) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderBincIO) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderBincIO) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderBincIO) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderBincIO) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderBincIO) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderBincIO) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderBincIO) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderBincIO) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderBincIO) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderBincIO) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderBincIO) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderBincIO) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderBincIO) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderBincIO) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderBincIO) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderBincIO) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderBincIO) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderBincIO) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincIO) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderBincIO) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderBincIO) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderBincIO) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderBincIO) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderBincIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderBincIO) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderBincIO) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderBincIO) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderBincIO) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderBincIO) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnBincIO

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

func (d *decoderBincIO) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnBincIO
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

func (d *decoderBincIO) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnBincIO

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

func (d *decoderBincIO) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnBincIO
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

func (d *decoderBincIO) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsBincIO)

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

func (d *decoderBincIO) reset() {
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

func (d *decoderBincIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderBincIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderBincIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderBincIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderBincIO) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderBincIO) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderBincIO) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderBincIO) Release() {}

func (d *decoderBincIO) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderBincIO) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderBincIO) decode(iv interface{}) {
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

func (d *decoderBincIO) decodeValue(rv reflect.Value, fn *decFnBincIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderBincIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnBincIO) {

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

func (d *decoderBincIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderBincIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderBincIO) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderBincIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderBincIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderBincIO) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderBincIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderBincIO) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderBincIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderBincIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincIO) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderBincIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderBincIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderBincIO) fn(t reflect.Type) *decFnBincIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderBincIO) fnNoExt(t reflect.Type) *decFnBincIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverBincIO) newDecoderBytes(in []byte, h Handle) *decoderBincIO {
	var c1 decoderBincIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverBincIO) newDecoderIO(in io.Reader, h Handle) *decoderBincIO {
	var c1 decoderBincIO
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverBincIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsBincIO) (f *fastpathDBincIO, u reflect.Type) {
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

func (helperDecDriverBincIO) decFindRtidFn(s []decRtidFnBincIO, rtid uintptr) (i uint, fn *decFnBincIO) {

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

func (helperDecDriverBincIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnBincIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnBincIO](v))
	}
	return
}

func (dh helperDecDriverBincIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsBincIO,
	checkExt bool) (fn *decFnBincIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverBincIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnBincIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverBincIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnBincIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnBincIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnBincIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnBincIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverBincIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsBincIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnBincIO) {
	fn = new(decFnBincIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderBincIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderBincIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderBincIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderBincIO).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderBincIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderBincIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderBincIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderBincIO).textUnmarshal
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
						fn.fd = func(d *decoderBincIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderBincIO, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderBincIO).kBool
			case reflect.String:
				fn.fd = (*decoderBincIO).kString
			case reflect.Int:
				fn.fd = (*decoderBincIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderBincIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderBincIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderBincIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderBincIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderBincIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderBincIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderBincIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderBincIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderBincIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderBincIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderBincIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderBincIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderBincIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderBincIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderBincIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderBincIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderBincIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderBincIO).kStructSimple
				} else {
					fn.fd = (*decoderBincIO).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderBincIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderBincIO).kInterface
			default:

				fn.fd = (*decoderBincIO).kErr
			}
		}
	}
	return
}
func (e *bincEncDriverIO) EncodeNil() {
	e.w.writen1(bincBdNil)
}

func (e *bincEncDriverIO) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else {
		bs := bincEncodeTime(t)
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	}
}

func (e *bincEncDriverIO) EncodeBool(b bool) {
	if b {
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	}
}

func (e *bincEncDriverIO) encSpFloat(f float64) (done bool) {
	if f == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
	} else if math.IsNaN(float64(f)) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNan)
	} else if math.IsInf(float64(f), +1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpPosInf)
	} else if math.IsInf(float64(f), -1) {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegInf)
	} else {
		return
	}
	return true
}

func (e *bincEncDriverIO) EncodeFloat32(f float32) {
	if !e.encSpFloat(float64(f)) {
		e.w.writen1(bincVdFloat<<4 | bincFlBin32)
		e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
	}
}

func (e *bincEncDriverIO) EncodeFloat64(f float64) {
	if e.encSpFloat(f) {
		return
	}
	b := bigen.PutUint64(math.Float64bits(f))
	if bincDoPrune {
		i := 7
		for ; i >= 0 && (b[i] == 0); i-- {
		}
		i++
		if i <= 6 {
			e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.w.writen1(byte(i))
			e.w.writeb(b[:i])
			return
		}
	}
	e.w.writen1(bincVdFloat<<4 | bincFlBin64)
	e.w.writen8(b)
}

func (e *bincEncDriverIO) encIntegerPrune32(bd byte, pos bool, v uint64) {
	b := bigen.PutUint32(uint32(v))
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 3 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 3)
		e.w.writen4(b)
	}
}

func (e *bincEncDriverIO) encIntegerPrune64(bd byte, pos bool, v uint64) {
	b := bigen.PutUint64(v)
	if bincDoPrune {
		i := byte(pruneSignExt(b[:], pos))
		e.w.writen1(bd | 7 - i)
		e.w.writeb(b[i:])
	} else {
		e.w.writen1(bd | 7)
		e.w.writen8(b)
	}
}

func (e *bincEncDriverIO) EncodeInt(v int64) {
	if v >= 0 {
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	} else if v == -1 {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	} else {
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	}
}

func (e *bincEncDriverIO) EncodeUint(v uint64) {
	e.encUint(bincVdPosInt<<4, true, v)
}

func (e *bincEncDriverIO) encUint(bd byte, pos bool, v uint64) {
	if v == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	} else if pos && v >= 1 && v <= 16 {
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.encIntegerPrune32(bd, pos, v)
	} else {
		e.encIntegerPrune64(bd, pos, v)
	}
}

func (e *bincEncDriverIO) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *bincEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *bincEncDriverIO) encodeExtPreamble(xtag byte, length int) {
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncDriverIO) WriteArrayStart(length int) {
	e.encLen(bincVdArray<<4, uint64(length))
}

func (e *bincEncDriverIO) WriteMapStart(length int) {
	e.encLen(bincVdMap<<4, uint64(length))
}

func (e *bincEncDriverIO) WriteArrayEmpty() {

	e.w.writen1(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriverIO) WriteMapEmpty() {

	e.w.writen1(bincVdMap<<4 | uint8(0+4))
}

func (e *bincEncDriverIO) EncodeSymbol(v string) {

	l := len(v)
	if l == 0 {
		e.encBytesLen(cUTF8, 0)
		return
	} else if l == 1 {
		e.encBytesLen(cUTF8, 1)
		e.w.writen1(v[0])
		return
	}
	if e.m == nil {
		e.m = make(map[string]uint16, 16)
	}
	ui, ok := e.m[v]
	if ok {
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			e.w.writen2(bigen.PutUint16(ui))
		}
	} else {
		e.e.seq++
		ui = e.e.seq
		e.m[v] = ui
		var lenprec uint8
		if l <= math.MaxUint8 {

		} else if l <= math.MaxUint16 {
			lenprec = 1
		} else if int64(l) <= math.MaxUint32 {
			lenprec = 2
		} else {
			lenprec = 3
		}
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4|0x4|lenprec, byte(ui))
		} else {
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			e.w.writen2(bigen.PutUint16(ui))
		}
		if lenprec == 0 {
			e.w.writen1(byte(l))
		} else if lenprec == 1 {
			e.w.writen2(bigen.PutUint16(uint16(l)))
		} else if lenprec == 2 {
			e.w.writen4(bigen.PutUint32(uint32(l)))
		} else {
			e.w.writen8(bigen.PutUint64(uint64(l)))
		}
		e.w.writestr(v)
	}
}

func (e *bincEncDriverIO) EncodeString(v string) {
	if e.h.StringToRaw {
		e.encLen(bincVdByteArray<<4, uint64(len(v)))
		if len(v) > 0 {
			e.w.writestr(v)
		}
		return
	}
	e.EncodeStringEnc(cUTF8, v)
}

func (e *bincEncDriverIO) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *bincEncDriverIO) EncodeStringEnc(c charEncoding, v string) {
	if e.e.c == containerMapKey && c == cUTF8 && (e.h.AsSymbols == 1) {
		e.EncodeSymbol(v)
		return
	}
	e.encLen(bincVdString<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writestr(v)
	}
}

func (e *bincEncDriverIO) EncodeStringBytesRaw(v []byte) {
	e.encLen(bincVdByteArray<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writeb(v)
	}
}

func (e *bincEncDriverIO) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *bincEncDriverIO) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = bincBdNil
	}
	e.w.writen1(v)
}

func (e *bincEncDriverIO) writeNilArray() {
	e.writeNilOr(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriverIO) writeNilMap() {
	e.writeNilOr(bincVdMap<<4 | uint8(0+4))
}

func (e *bincEncDriverIO) writeNilBytes() {
	e.writeNilOr(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriverIO) encBytesLen(c charEncoding, length uint64) {

	if c == cRAW {
		e.encLen(bincVdByteArray<<4, length)
	} else {
		e.encLen(bincVdString<<4, length)
	}
}

func (e *bincEncDriverIO) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l+4))
	} else {
		e.encLenNumber(bd, l)
	}
}

func (e *bincEncDriverIO) encLenNumber(bd byte, v uint64) {
	if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v))
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.w.writen1(bd | 0x02)
		e.w.writen4(bigen.PutUint32(uint32(v)))
	} else {
		e.w.writen1(bd | 0x03)
		e.w.writen8(bigen.PutUint64(uint64(v)))
	}
}

func (d *bincDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
}

func (d *bincDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *bincDecDriverIO) TryNil() bool {
	return d.advanceNil()
}

func (d *bincDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return valueTypeNil
	} else if d.vd == bincVdByteArray {
		return valueTypeBytes
	} else if d.vd == bincVdString {
		return valueTypeString
	} else if d.vd == bincVdArray {
		return valueTypeArray
	} else if d.vd == bincVdMap {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *bincDecDriverIO) DecodeTime() (t time.Time) {
	if d.advanceNil() {
		return
	}
	if d.vd != bincVdTimestamp {
		halt.errorf("cannot decode time - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	t, err := bincDecodeTime(d.r.readx(uint(d.vs)))
	halt.onerror(err)
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) decFloatPruned(maxlen uint8) {
	l := d.r.readn1()
	if l > maxlen {
		halt.errorf("cannot read float - at most %v bytes used to represent float - received %v bytes", maxlen, l)
	}
	for i := l; i < maxlen; i++ {
		d.d.b[i] = 0
	}
	d.r.readb(d.d.b[0:l])
}

func (d *bincDecDriverIO) decFloatPre32() (b [4]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn4()
	} else {
		d.decFloatPruned(4)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverIO) decFloatPre64() (b [8]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn8()
	} else {
		d.decFloatPruned(8)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriverIO) decFloatVal() (f float64) {
	switch d.vs & 0x7 {
	case bincFlBin32:
		f = float64(math.Float32frombits(bigen.Uint32(d.decFloatPre32())))
	case bincFlBin64:
		f = math.Float64frombits(bigen.Uint64(d.decFloatPre64()))
	default:

		halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincDecDriverIO) decUint() (v uint64) {
	switch d.vs {
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		v = uint64(bigen.Uint16(d.r.readn2()))
	case 2:
		b3 := d.r.readn3()
		var b [4]byte
		copy(b[1:], b3[:])
		v = uint64(bigen.Uint32(b))
	case 3:
		v = uint64(bigen.Uint32(d.r.readn4()))
	case 4, 5, 6:

		bs := d.d.b[:8]
		clear(bs)
		d.r.readb(bs[(7 - d.vs):])
		v = bigen.Uint64(*(*[8]byte)(bs))
	case 7:
		v = bigen.Uint64(d.r.readn8())
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	return
}

func (d *bincDecDriverIO) uintBytes() (bs []byte) {
	switch d.vs {
	case 0:
		bs = d.d.b[:1]
		bs[0] = d.r.readn1()
		return
	case 1:
		bs = d.d.b[:2]
	case 2:
		bs = d.d.b[:3]
	case 3:
		bs = d.d.b[:4]
	case 4, 5, 6:
		lim := 7 - d.vs
		bs = d.d.b[lim:8]
	case 7:
		bs = d.d.b[:8]
	default:
		halt.errorf("unsigned integers with greater than 64 bits of precision not supported: d.vs: %v %x", d.vs, d.vs)
	}
	d.r.readb(bs)
	return
}

func (d *bincDecDriverIO) decInteger() (ui uint64, neg, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdPosInt {
		ui = d.decUint()
	} else if vd == bincVdNegInt {
		ui = d.decUint()
		neg = true
	} else if vd == bincVdSmallInt {
		ui = uint64(d.vs) + 1
	} else if vd == bincVdSpecial {
		if vs == bincSpZero {

		} else if vs == bincSpNegOne {
			neg = true
			ui = 1
		} else {
			ok = false

		}
	} else {
		ok = false

	}
	return
}

func (d *bincDecDriverIO) decFloat() (f float64, ok bool) {
	ok = true
	vd, vs := d.vd, d.vs
	if vd == bincVdSpecial {
		if vs == bincSpNan {
			f = math.NaN()
		} else if vs == bincSpPosInf {
			f = math.Inf(1)
		} else if vs == bincSpZeroFloat || vs == bincSpZero {

		} else if vs == bincSpNegInf {
			f = math.Inf(-1)
		} else {
			ok = false

		}
	} else if vd == bincVdFloat {
		f = d.decFloatVal()
	} else {
		ok = false
	}
	return
}

func (d *bincDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == (bincVdSpecial | bincSpFalse) {

	} else if d.bd == (bincVdSpecial | bincSpTrue) {
		b = true
	} else {
		halt.errorf("bool - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdMap {
		halt.errorf("map - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	if d.vd != bincVdArray {
		halt.errorf("array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	length = d.decLen()
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) decLen() int {
	if d.vs > 3 {
		return int(d.vs - 4)
	}
	return int(d.decLenNumber())
}

func (d *bincDecDriverIO) decLenNumber() (v uint64) {
	if x := d.vs; x == 0 {
		v = uint64(d.r.readn1())
	} else if x == 1 {
		v = uint64(bigen.Uint16(d.r.readn2()))
	} else if x == 2 {
		v = uint64(bigen.Uint32(d.r.readn4()))
	} else {
		v = bigen.Uint64(d.r.readn8())
	}
	return
}

func (d *bincDecDriverIO) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool
	var slen = -1
	switch d.vd {
	case bincVdString, bincVdByteArray:
		slen = d.decLen()
		bs, cond = d.r.readxb(uint(slen))
		state = d.d.attachState(cond)
	case bincVdSymbol:

		var symbol uint16
		vs := d.vs
		if vs&0x8 == 0 {
			symbol = uint16(d.r.readn1())
		} else {
			symbol = uint16(bigen.Uint16(d.r.readn2()))
		}
		if d.s == nil {
			d.s = make(map[uint16][]byte, 16)
		}

		if vs&0x4 == 0 {
			bs = d.s[symbol]
		} else {
			switch vs & 0x3 {
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(bigen.Uint16(d.r.readn2()))
			case 2:
				slen = int(bigen.Uint32(d.r.readn4()))
			case 3:
				slen = int(bigen.Uint64(d.r.readn8()))
			}

			bs, cond = d.r.readxb(uint(slen))
			bs = d.d.detach2Bytes(bs, d.d.attachState(cond))
			d.s[symbol] = bs
		}
		state = dBytesDetach
	default:
		halt.errorf("string/bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if d.h.ValidateUnicode && !utf8.Valid(bs) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", bs)
	}

	d.bdRead = false
	return
}

func (d *bincDecDriverIO) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool
	if d.vd == bincVdArray {
		slen := d.ReadArrayStart()
		bs, cond = usableByteSlice(d.d.buf, slen)
		for i := 0; i < slen; i++ {
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
	if !(d.vd == bincVdString || d.vd == bincVdByteArray) {
		halt.errorf("bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	clen := d.decLen()
	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *bincDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *bincDecDriverIO) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *bincDecDriverIO) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
	if xtagIn > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtagIn)
	}
	if d.advanceNil() {
		return
	}
	tag := uint8(xtagIn)
	if d.vd == bincVdCustomExt {
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag - got %b, expecting: %v", xtag, tag)
		}
		xbs, ok = d.r.readxb(uint(l))
		bstate = d.d.attachState(ok)

	} else if d.vd == bincVdByteArray {
		xbs, bstate = d.DecodeBytes()
	} else {
		halt.errorf("ext expects extensions or byte array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	ok = true
	return
}

func (d *bincDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}

	n := d.d.naked()
	var decodeFurther bool

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil:
			n.v = valueTypeNil
		case bincSpFalse:
			n.v = valueTypeBool
			n.b = false
		case bincSpTrue:
			n.v = valueTypeBool
			n.b = true
		case bincSpNan:
			n.v = valueTypeFloat
			n.f = math.NaN()
		case bincSpPosInf:
			n.v = valueTypeFloat
			n.f = math.Inf(1)
		case bincSpNegInf:
			n.v = valueTypeFloat
			n.f = math.Inf(-1)
		case bincSpZeroFloat:
			n.v = valueTypeFloat
			n.f = float64(0)
		case bincSpZero:
			n.v = valueTypeUint
			n.u = uint64(0)
		case bincSpNegOne:
			n.v = valueTypeInt
			n.i = int64(-1)
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
		n.v = valueTypeUint
		n.u = uint64(int8(d.vs)) + 1
	case bincVdPosInt:
		n.v = valueTypeUint
		n.u = d.decUint()
	case bincVdNegInt:
		n.v = valueTypeInt
		n.i = -(int64(d.decUint()))
	case bincVdFloat:
		n.v = valueTypeFloat
		n.f = d.decFloatVal()
	case bincVdString:
		n.v = valueTypeString
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case bincVdByteArray:
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
	case bincVdSymbol:
		n.v = valueTypeSymbol
		n.s = d.d.detach2Str(d.DecodeStringAsBytes())
	case bincVdTimestamp:
		n.v = valueTypeTime
		tt, err := bincDecodeTime(d.r.readx(uint(d.vs)))
		halt.onerror(err)
		n.t = tt
	case bincVdCustomExt:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(uint(l))
	case bincVdArray:
		n.v = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}

	if !decodeFurther {
		d.bdRead = false
	}
	if n.v == valueTypeUint && d.h.SignedInteger {
		n.v = valueTypeInt
		n.i = int64(n.u)
	}
}

func (d *bincDecDriverIO) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *bincDecDriverIO) nextValueBytesBdReadR() {
	fnLen := func(vs byte) uint {
		switch vs {
		case 0:
			x := d.r.readn1()
			return uint(x)
		case 1:
			x := d.r.readn2()
			return uint(bigen.Uint16(x))
		case 2:
			x := d.r.readn4()
			return uint(bigen.Uint32(x))
		case 3:
			x := d.r.readn8()
			return uint(bigen.Uint64(x))
		default:
			return uint(vs - 4)
		}
	}

	var clen uint

	switch d.vd {
	case bincVdSpecial:
		switch d.vs {
		case bincSpNil, bincSpFalse, bincSpTrue, bincSpNan, bincSpPosInf:
		case bincSpNegInf, bincSpZeroFloat, bincSpZero, bincSpNegOne:
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
	case bincVdPosInt, bincVdNegInt:
		d.uintBytes()
	case bincVdFloat:
		fn := func(xlen byte) {
			if d.vs&0x8 != 0 {
				xlen = d.r.readn1()
				if xlen > 8 {
					halt.errorf("cannot read float - at most 8 bytes used to represent float - received %v bytes", xlen)
				}
			}
			d.r.readb(d.d.b[:xlen])
		}
		switch d.vs & 0x7 {
		case bincFlBin32:
			fn(4)
		case bincFlBin64:
			fn(8)
		default:
			halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdString, bincVdByteArray:
		clen = fnLen(d.vs)
		d.r.skip(clen)
	case bincVdSymbol:
		if d.vs&0x8 == 0 {
			d.r.readn1()
		} else {
			d.r.skip(2)
		}
		if d.vs&0x4 != 0 {
			clen = fnLen(d.vs & 0x3)
			d.r.skip(clen)
		}
	case bincVdTimestamp:
		d.r.skip(uint(d.vs))
	case bincVdCustomExt:
		clen = fnLen(d.vs)
		d.r.readn1()
		d.r.skip(clen)
	case bincVdArray:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case bincVdMap:
		clen = fnLen(d.vs)
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	default:
		halt.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincEncDriverIO) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*BincHandle)
	d.e = shared
	if shared.bytes {
		fp = bincFpEncBytes
	} else {
		fp = bincFpEncIO
	}

	d.init2(enc)
	return
}

func (e *bincEncDriverIO) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *bincEncDriverIO) writerEnd() { e.w.end() }

func (e *bincEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *bincEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *bincDecDriverIO) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*BincHandle)
	d.d = shared
	if shared.bytes {
		fp = bincFpDecBytes
	} else {
		fp = bincFpDecIO
	}

	d.init2(dec)
	return
}

func (d *bincDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *bincDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *bincDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *bincDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, bincdescbd(d.bd))
}

func (d *bincDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
