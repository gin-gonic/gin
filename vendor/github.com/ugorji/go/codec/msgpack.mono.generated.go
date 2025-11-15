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

type helperEncDriverMsgpackBytes struct{}
type encFnMsgpackBytes struct {
	i  encFnInfo
	fe func(*encoderMsgpackBytes, *encFnInfo, reflect.Value)
}
type encRtidFnMsgpackBytes struct {
	rtid uintptr
	fn   *encFnMsgpackBytes
}
type encoderMsgpackBytes struct {
	dh helperEncDriverMsgpackBytes
	fp *fastpathEsMsgpackBytes
	e  msgpackEncDriverBytes
	encoderBase
}
type helperDecDriverMsgpackBytes struct{}
type decFnMsgpackBytes struct {
	i  decFnInfo
	fd func(*decoderMsgpackBytes, *decFnInfo, reflect.Value)
}
type decRtidFnMsgpackBytes struct {
	rtid uintptr
	fn   *decFnMsgpackBytes
}
type decoderMsgpackBytes struct {
	dh helperDecDriverMsgpackBytes
	fp *fastpathDsMsgpackBytes
	d  msgpackDecDriverBytes
	decoderBase
}
type msgpackEncDriverBytes struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *MsgpackHandle
	e *encoderBase
	w bytesEncAppender
}
type msgpackDecDriverBytes struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	decInit2er

	h *MsgpackHandle
	d *decoderBase
	r bytesDecReader

	bdAndBdread

	noBuiltInTypes
}

func (e *encoderMsgpackBytes) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderMsgpackBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderMsgpackBytes) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderMsgpackBytes) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderMsgpackBytes) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderMsgpackBytes) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderMsgpackBytes) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderMsgpackBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderMsgpackBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderMsgpackBytes) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderMsgpackBytes) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderMsgpackBytes) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderMsgpackBytes) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderMsgpackBytes) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderMsgpackBytes) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderMsgpackBytes) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderMsgpackBytes) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderMsgpackBytes) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderMsgpackBytes) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderMsgpackBytes) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderMsgpackBytes) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderMsgpackBytes) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderMsgpackBytes) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderMsgpackBytes) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderMsgpackBytes) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderMsgpackBytes) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderMsgpackBytes) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderMsgpackBytes) kSeqFn(rt reflect.Type) (fn *encFnMsgpackBytes) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderMsgpackBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnMsgpackBytes
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

func (e *encoderMsgpackBytes) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnMsgpackBytes
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

func (e *encoderMsgpackBytes) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderMsgpackBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderMsgpackBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderMsgpackBytes) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderMsgpackBytes) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderMsgpackBytes) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderMsgpackBytes) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderMsgpackBytes) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnMsgpackBytes

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

func (e *encoderMsgpackBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnMsgpackBytes) {
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

func (e *encoderMsgpackBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsMsgpackBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderMsgpackBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderMsgpackBytes) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderMsgpackBytes) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderMsgpackBytes) mustEncode(v interface{}) {
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

func (e *encoderMsgpackBytes) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderMsgpackBytes) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderMsgpackBytes) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderMsgpackBytes) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderMsgpackBytes) encodeValue(rv reflect.Value, fn *encFnMsgpackBytes) {

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

func (e *encoderMsgpackBytes) encodeValueNonNil(rv reflect.Value, fn *encFnMsgpackBytes) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderMsgpackBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderMsgpackBytes) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderMsgpackBytes) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderMsgpackBytes) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderMsgpackBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderMsgpackBytes) fn(t reflect.Type) *encFnMsgpackBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderMsgpackBytes) fnNoExt(t reflect.Type) *encFnMsgpackBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderMsgpackBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderMsgpackBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderMsgpackBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderMsgpackBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderMsgpackBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderMsgpackBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderMsgpackBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderMsgpackBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverMsgpackBytes) newEncoderBytes(out *[]byte, h Handle) *encoderMsgpackBytes {
	var c1 encoderMsgpackBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverMsgpackBytes) newEncoderIO(out io.Writer, h Handle) *encoderMsgpackBytes {
	var c1 encoderMsgpackBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverMsgpackBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsMsgpackBytes) (f *fastpathEMsgpackBytes, u reflect.Type) {
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

func (helperEncDriverMsgpackBytes) encFindRtidFn(s []encRtidFnMsgpackBytes, rtid uintptr) (i uint, fn *encFnMsgpackBytes) {

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

func (helperEncDriverMsgpackBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnMsgpackBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnMsgpackBytes](v))
	}
	return
}

func (dh helperEncDriverMsgpackBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsMsgpackBytes, checkExt bool) (fn *encFnMsgpackBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverMsgpackBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsMsgpackBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnMsgpackBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnMsgpackBytes = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverMsgpackBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsMsgpackBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnMsgpackBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnMsgpackBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnMsgpackBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnMsgpackBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnMsgpackBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverMsgpackBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsMsgpackBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnMsgpackBytes) {
	fn = new(encFnMsgpackBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderMsgpackBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderMsgpackBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderMsgpackBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderMsgpackBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderMsgpackBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderMsgpackBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderMsgpackBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderMsgpackBytes).textMarshal
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
					fn.fe = func(e *encoderMsgpackBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderMsgpackBytes).kBool
			case reflect.String:

				fn.fe = (*encoderMsgpackBytes).kString
			case reflect.Int:
				fn.fe = (*encoderMsgpackBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderMsgpackBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderMsgpackBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderMsgpackBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderMsgpackBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderMsgpackBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderMsgpackBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderMsgpackBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderMsgpackBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderMsgpackBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderMsgpackBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderMsgpackBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderMsgpackBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderMsgpackBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderMsgpackBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderMsgpackBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderMsgpackBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderMsgpackBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderMsgpackBytes).kStructSimple
				} else {
					fn.fe = (*encoderMsgpackBytes).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderMsgpackBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderMsgpackBytes).kErr
			default:

				fn.fe = (*encoderMsgpackBytes).kErr
			}
		}
	}
	return
}
func (d *decoderMsgpackBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderMsgpackBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderMsgpackBytes) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderMsgpackBytes) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderMsgpackBytes) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderMsgpackBytes) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderMsgpackBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderMsgpackBytes) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderMsgpackBytes) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderMsgpackBytes) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderMsgpackBytes) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderMsgpackBytes) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderMsgpackBytes) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderMsgpackBytes) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderMsgpackBytes) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderMsgpackBytes) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderMsgpackBytes) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderMsgpackBytes) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderMsgpackBytes) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderMsgpackBytes) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderMsgpackBytes) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderMsgpackBytes) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderMsgpackBytes) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderMsgpackBytes) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderMsgpackBytes) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderMsgpackBytes) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderMsgpackBytes) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderMsgpackBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderMsgpackBytes) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderMsgpackBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderMsgpackBytes) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderMsgpackBytes) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderMsgpackBytes) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnMsgpackBytes

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

func (d *decoderMsgpackBytes) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnMsgpackBytes
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

func (d *decoderMsgpackBytes) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnMsgpackBytes

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

func (d *decoderMsgpackBytes) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnMsgpackBytes
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

func (d *decoderMsgpackBytes) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsMsgpackBytes)

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

func (d *decoderMsgpackBytes) reset() {
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

func (d *decoderMsgpackBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderMsgpackBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderMsgpackBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderMsgpackBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderMsgpackBytes) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderMsgpackBytes) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderMsgpackBytes) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderMsgpackBytes) Release() {}

func (d *decoderMsgpackBytes) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderMsgpackBytes) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderMsgpackBytes) decode(iv interface{}) {
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

func (d *decoderMsgpackBytes) decodeValue(rv reflect.Value, fn *decFnMsgpackBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderMsgpackBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnMsgpackBytes) {

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

func (d *decoderMsgpackBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderMsgpackBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderMsgpackBytes) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderMsgpackBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderMsgpackBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderMsgpackBytes) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderMsgpackBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderMsgpackBytes) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderMsgpackBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderMsgpackBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderMsgpackBytes) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderMsgpackBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderMsgpackBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderMsgpackBytes) fn(t reflect.Type) *decFnMsgpackBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderMsgpackBytes) fnNoExt(t reflect.Type) *decFnMsgpackBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverMsgpackBytes) newDecoderBytes(in []byte, h Handle) *decoderMsgpackBytes {
	var c1 decoderMsgpackBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverMsgpackBytes) newDecoderIO(in io.Reader, h Handle) *decoderMsgpackBytes {
	var c1 decoderMsgpackBytes
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverMsgpackBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsMsgpackBytes) (f *fastpathDMsgpackBytes, u reflect.Type) {
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

func (helperDecDriverMsgpackBytes) decFindRtidFn(s []decRtidFnMsgpackBytes, rtid uintptr) (i uint, fn *decFnMsgpackBytes) {

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

func (helperDecDriverMsgpackBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnMsgpackBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnMsgpackBytes](v))
	}
	return
}

func (dh helperDecDriverMsgpackBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsMsgpackBytes,
	checkExt bool) (fn *decFnMsgpackBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverMsgpackBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsMsgpackBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnMsgpackBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnMsgpackBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverMsgpackBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsMsgpackBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnMsgpackBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnMsgpackBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnMsgpackBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnMsgpackBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnMsgpackBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverMsgpackBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsMsgpackBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnMsgpackBytes) {
	fn = new(decFnMsgpackBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderMsgpackBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderMsgpackBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderMsgpackBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderMsgpackBytes).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderMsgpackBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderMsgpackBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderMsgpackBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderMsgpackBytes).textUnmarshal
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
						fn.fd = func(d *decoderMsgpackBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderMsgpackBytes, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderMsgpackBytes).kBool
			case reflect.String:
				fn.fd = (*decoderMsgpackBytes).kString
			case reflect.Int:
				fn.fd = (*decoderMsgpackBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderMsgpackBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderMsgpackBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderMsgpackBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderMsgpackBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderMsgpackBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderMsgpackBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderMsgpackBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderMsgpackBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderMsgpackBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderMsgpackBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderMsgpackBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderMsgpackBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderMsgpackBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderMsgpackBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderMsgpackBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderMsgpackBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderMsgpackBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderMsgpackBytes).kStructSimple
				} else {
					fn.fd = (*decoderMsgpackBytes).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderMsgpackBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderMsgpackBytes).kInterface
			default:

				fn.fd = (*decoderMsgpackBytes).kErr
			}
		}
	}
	return
}
func (e *msgpackEncDriverBytes) EncodeNil() {
	e.w.writen1(mpNil)
}

func (e *msgpackEncDriverBytes) EncodeInt(i int64) {
	if e.h.PositiveIntUnsigned && i >= 0 {
		e.EncodeUint(uint64(i))
	} else if i > math.MaxInt8 {
		if i <= math.MaxInt16 {
			e.w.writen1(mpInt16)
			e.w.writen2(bigen.PutUint16(uint16(i)))
		} else if i <= math.MaxInt32 {
			e.w.writen1(mpInt32)
			e.w.writen4(bigen.PutUint32(uint32(i)))
		} else {
			e.w.writen1(mpInt64)
			e.w.writen8(bigen.PutUint64(uint64(i)))
		}
	} else if i >= -32 {
		if e.h.NoFixedNum {
			e.w.writen2(mpInt8, byte(i))
		} else {
			e.w.writen1(byte(i))
		}
	} else if i >= math.MinInt8 {
		e.w.writen2(mpInt8, byte(i))
	} else if i >= math.MinInt16 {
		e.w.writen1(mpInt16)
		e.w.writen2(bigen.PutUint16(uint16(i)))
	} else if i >= math.MinInt32 {
		e.w.writen1(mpInt32)
		e.w.writen4(bigen.PutUint32(uint32(i)))
	} else {
		e.w.writen1(mpInt64)
		e.w.writen8(bigen.PutUint64(uint64(i)))
	}
}

func (e *msgpackEncDriverBytes) EncodeUint(i uint64) {
	if i <= math.MaxInt8 {
		if e.h.NoFixedNum {
			e.w.writen2(mpUint8, byte(i))
		} else {
			e.w.writen1(byte(i))
		}
	} else if i <= math.MaxUint8 {
		e.w.writen2(mpUint8, byte(i))
	} else if i <= math.MaxUint16 {
		e.w.writen1(mpUint16)
		e.w.writen2(bigen.PutUint16(uint16(i)))
	} else if i <= math.MaxUint32 {
		e.w.writen1(mpUint32)
		e.w.writen4(bigen.PutUint32(uint32(i)))
	} else {
		e.w.writen1(mpUint64)
		e.w.writen8(bigen.PutUint64(uint64(i)))
	}
}

func (e *msgpackEncDriverBytes) EncodeBool(b bool) {
	if b {
		e.w.writen1(mpTrue)
	} else {
		e.w.writen1(mpFalse)
	}
}

func (e *msgpackEncDriverBytes) EncodeFloat32(f float32) {
	e.w.writen1(mpFloat)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *msgpackEncDriverBytes) EncodeFloat64(f float64) {
	e.w.writen1(mpDouble)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *msgpackEncDriverBytes) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
		return
	}
	t = t.UTC()
	sec, nsec := t.Unix(), uint64(t.Nanosecond())
	var data64 uint64
	var l = 4
	if sec >= 0 && sec>>34 == 0 {
		data64 = (nsec << 34) | uint64(sec)
		if data64&0xffffffff00000000 != 0 {
			l = 8
		}
	} else {
		l = 12
	}
	if e.h.WriteExt {
		e.encodeExtPreamble(mpTimeExtTagU, l)
	} else {
		e.writeContainerLen(msgpackContainerRawLegacy, l)
	}
	switch l {
	case 4:
		e.w.writen4(bigen.PutUint32(uint32(data64)))
	case 8:
		e.w.writen8(bigen.PutUint64(data64))
	case 12:
		e.w.writen4(bigen.PutUint32(uint32(nsec)))
		e.w.writen8(bigen.PutUint64(uint64(sec)))
	}
}

func (e *msgpackEncDriverBytes) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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
	if e.h.WriteExt {
		e.encodeExtPreamble(uint8(xtag), len(bs))
		e.w.writeb(bs)
	} else {
		e.EncodeBytes(bs)
	}
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *msgpackEncDriverBytes) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *msgpackEncDriverBytes) encodeExtPreamble(xtag byte, l int) {
	if l == 1 {
		e.w.writen2(mpFixExt1, xtag)
	} else if l == 2 {
		e.w.writen2(mpFixExt2, xtag)
	} else if l == 4 {
		e.w.writen2(mpFixExt4, xtag)
	} else if l == 8 {
		e.w.writen2(mpFixExt8, xtag)
	} else if l == 16 {
		e.w.writen2(mpFixExt16, xtag)
	} else if l < 256 {
		e.w.writen2(mpExt8, byte(l))
		e.w.writen1(xtag)
	} else if l < 65536 {
		e.w.writen1(mpExt16)
		e.w.writen2(bigen.PutUint16(uint16(l)))
		e.w.writen1(xtag)
	} else {
		e.w.writen1(mpExt32)
		e.w.writen4(bigen.PutUint32(uint32(l)))
		e.w.writen1(xtag)
	}
}

func (e *msgpackEncDriverBytes) WriteArrayStart(length int) {
	e.writeContainerLen(msgpackContainerList, length)
}

func (e *msgpackEncDriverBytes) WriteMapStart(length int) {
	e.writeContainerLen(msgpackContainerMap, length)
}

func (e *msgpackEncDriverBytes) WriteArrayEmpty() {

	e.w.writen1(mpFixArrayMin)
}

func (e *msgpackEncDriverBytes) WriteMapEmpty() {

	e.w.writen1(mpFixMapMin)
}

func (e *msgpackEncDriverBytes) EncodeString(s string) {
	var ct msgpackContainerType
	if e.h.WriteExt {
		if e.h.StringToRaw {
			ct = msgpackContainerBin
		} else {
			ct = msgpackContainerStr
		}
	} else {
		ct = msgpackContainerRawLegacy
	}
	e.writeContainerLen(ct, len(s))
	if len(s) > 0 {
		e.w.writestr(s)
	}
}

func (e *msgpackEncDriverBytes) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *msgpackEncDriverBytes) EncodeStringBytesRaw(bs []byte) {
	if e.h.WriteExt {
		e.writeContainerLen(msgpackContainerBin, len(bs))
	} else {
		e.writeContainerLen(msgpackContainerRawLegacy, len(bs))
	}
	if len(bs) > 0 {
		e.w.writeb(bs)
	}
}

func (e *msgpackEncDriverBytes) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *msgpackEncDriverBytes) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = mpNil
	}
	e.w.writen1(v)
}

func (e *msgpackEncDriverBytes) writeNilArray() {
	e.writeNilOr(mpFixArrayMin)
}

func (e *msgpackEncDriverBytes) writeNilMap() {
	e.writeNilOr(mpFixMapMin)
}

func (e *msgpackEncDriverBytes) writeNilBytes() {
	e.writeNilOr(mpFixStrMin)
}

func (e *msgpackEncDriverBytes) writeContainerLen(ct msgpackContainerType, l int) {
	if ct.fixCutoff > 0 && l < int(ct.fixCutoff) {
		e.w.writen1(ct.bFixMin | byte(l))
	} else if ct.b8 > 0 && l < 256 {
		e.w.writen2(ct.b8, uint8(l))
	} else if l < 65536 {
		e.w.writen1(ct.b16)
		e.w.writen2(bigen.PutUint16(uint16(l)))
	} else {
		e.w.writen1(ct.b32)
		e.w.writen4(bigen.PutUint32(uint32(l)))
	}
}

func (d *msgpackDecDriverBytes) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}
	bd := d.bd
	n := d.d.naked()
	var decodeFurther bool

	switch bd {
	case mpNil:
		n.v = valueTypeNil
		d.bdRead = false
	case mpFalse:
		n.v = valueTypeBool
		n.b = false
	case mpTrue:
		n.v = valueTypeBool
		n.b = true

	case mpFloat:
		n.v = valueTypeFloat
		n.f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case mpDouble:
		n.v = valueTypeFloat
		n.f = math.Float64frombits(bigen.Uint64(d.r.readn8()))

	case mpUint8:
		n.v = valueTypeUint
		n.u = uint64(d.r.readn1())
	case mpUint16:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint16(d.r.readn2()))
	case mpUint32:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint32(d.r.readn4()))
	case mpUint64:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint64(d.r.readn8()))

	case mpInt8:
		n.v = valueTypeInt
		n.i = int64(int8(d.r.readn1()))
	case mpInt16:
		n.v = valueTypeInt
		n.i = int64(int16(bigen.Uint16(d.r.readn2())))
	case mpInt32:
		n.v = valueTypeInt
		n.i = int64(int32(bigen.Uint32(d.r.readn4())))
	case mpInt64:
		n.v = valueTypeInt
		n.i = int64(int64(bigen.Uint64(d.r.readn8())))

	default:
		switch {
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:

			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:

			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd == mpStr8, bd == mpStr16, bd == mpStr32, bd >= mpFixStrMin && bd <= mpFixStrMax:
			d.d.fauxUnionReadRawBytes(d, d.h.WriteExt, d.h.RawToString)

		case bd == mpBin8, bd == mpBin16, bd == mpBin32:
			d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
		case bd == mpArray16, bd == mpArray32, bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			n.v = valueTypeArray
			decodeFurther = true
		case bd == mpMap16, bd == mpMap32, bd >= mpFixMapMin && bd <= mpFixMapMax:
			n.v = valueTypeMap
			decodeFurther = true
		case bd >= mpFixExt1 && bd <= mpFixExt16, bd >= mpExt8 && bd <= mpExt32:
			n.v = valueTypeExt
			clen := d.readExtLen()
			n.u = uint64(d.r.readn1())
			if n.u == uint64(mpTimeExtTagU) {
				n.v = valueTypeTime
				n.t = d.decodeTime(clen)
			} else {
				n.l = d.r.readx(uint(clen))
			}
		default:
			halt.errorf("cannot infer value: %s: Ox%x/%d/%s", msgBadDesc, bd, bd, mpdesc(bd))
		}
	}
	if !decodeFurther {
		d.bdRead = false
	}
	if n.v == valueTypeUint && d.h.SignedInteger {
		n.v = valueTypeInt
		n.i = int64(n.u)
	}
}

func (d *msgpackDecDriverBytes) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *msgpackDecDriverBytes) nextValueBytesBdReadR() {
	bd := d.bd

	var clen uint

	switch bd {
	case mpNil, mpFalse, mpTrue:
	case mpUint8, mpInt8:
		d.r.readn1()
	case mpUint16, mpInt16:
		d.r.skip(2)
	case mpFloat, mpUint32, mpInt32:
		d.r.skip(4)
	case mpDouble, mpUint64, mpInt64:
		d.r.skip(8)
	case mpStr8, mpBin8:
		clen = uint(d.r.readn1())
		d.r.skip(clen)
	case mpStr16, mpBin16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		d.r.skip(clen)
	case mpStr32, mpBin32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		d.r.skip(clen)
	case mpFixExt1:
		d.r.readn1()
		d.r.readn1()
	case mpFixExt2:
		d.r.readn1()
		d.r.skip(2)
	case mpFixExt4:
		d.r.readn1()
		d.r.skip(4)
	case mpFixExt8:
		d.r.readn1()
		d.r.skip(8)
	case mpFixExt16:
		d.r.readn1()
		d.r.skip(16)
	case mpExt8:
		clen = uint(d.r.readn1())
		d.r.readn1()
		d.r.skip(clen)
	case mpExt16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		d.r.readn1()
		d.r.skip(clen)
	case mpExt32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		d.r.readn1()
		d.r.skip(clen)
	case mpArray16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case mpArray32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case mpMap16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case mpMap32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	default:
		switch {
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
		case bd >= mpFixStrMin && bd <= mpFixStrMax:
			clen = uint(mpFixStrMin ^ bd)
			d.r.skip(clen)
		case bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			clen = uint(mpFixArrayMin ^ bd)
			for i := uint(0); i < clen; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		case bd >= mpFixMapMin && bd <= mpFixMapMax:
			clen = uint(mpFixMapMin ^ bd)
			for i := uint(0); i < clen; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		default:
			halt.errorf("nextValueBytes: cannot infer value: %s: Ox%x/%d/%s", msgBadDesc, bd, bd, mpdesc(bd))
		}
	}
	return
}

func (d *msgpackDecDriverBytes) decFloat4Int32() (f float32) {
	fbits := bigen.Uint32(d.r.readn4())
	f = math.Float32frombits(fbits)
	if !noFrac32(fbits) {
		halt.errorf("assigning integer value from float32 with a fraction: %v", f)
	}
	return
}

func (d *msgpackDecDriverBytes) decFloat4Int64() (f float64) {
	fbits := bigen.Uint64(d.r.readn8())
	f = math.Float64frombits(fbits)
	if !noFrac64(fbits) {
		halt.errorf("assigning integer value from float64 with a fraction: %v", f)
	}
	return
}

func (d *msgpackDecDriverBytes) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	switch d.bd {
	case mpUint8:
		i = int64(uint64(d.r.readn1()))
	case mpUint16:
		i = int64(uint64(bigen.Uint16(d.r.readn2())))
	case mpUint32:
		i = int64(uint64(bigen.Uint32(d.r.readn4())))
	case mpUint64:
		i = int64(bigen.Uint64(d.r.readn8()))
	case mpInt8:
		i = int64(int8(d.r.readn1()))
	case mpInt16:
		i = int64(int16(bigen.Uint16(d.r.readn2())))
	case mpInt32:
		i = int64(int32(bigen.Uint32(d.r.readn4())))
	case mpInt64:
		i = int64(bigen.Uint64(d.r.readn8()))
	case mpFloat:
		i = int64(d.decFloat4Int32())
	case mpDouble:
		i = int64(d.decFloat4Int64())
	default:
		switch {
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			i = int64(int8(d.bd))
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			i = int64(int8(d.bd))
		default:
			halt.errorf("cannot decode signed integer: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
		}
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverBytes) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	switch d.bd {
	case mpUint8:
		ui = uint64(d.r.readn1())
	case mpUint16:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case mpUint32:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case mpUint64:
		ui = bigen.Uint64(d.r.readn8())
	case mpInt8:
		if i := int64(int8(d.r.readn1())); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt16:
		if i := int64(int16(bigen.Uint16(d.r.readn2()))); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt32:
		if i := int64(int32(bigen.Uint32(d.r.readn4()))); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt64:
		if i := int64(bigen.Uint64(d.r.readn8())); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpFloat:
		if f := d.decFloat4Int32(); f >= 0 {
			ui = uint64(f)
		} else {
			halt.errorf("assigning negative float value: %v, to unsigned type", f)
		}
	case mpDouble:
		if f := d.decFloat4Int64(); f >= 0 {
			ui = uint64(f)
		} else {
			halt.errorf("assigning negative float value: %v, to unsigned type", f)
		}
	default:
		switch {
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			ui = uint64(d.bd)
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			halt.errorf("assigning negative signed value: %v, to unsigned type", int(d.bd))
		default:
			halt.errorf("cannot decode unsigned integer: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
		}
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverBytes) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.bd == mpFloat {
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	} else if d.bd == mpDouble {
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	} else {
		f = float64(d.DecodeInt64())
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverBytes) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == mpFalse || d.bd == 0 {

	} else if d.bd == mpTrue || d.bd == 1 {
		b = true
	} else {
		halt.errorf("cannot decode bool: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverBytes) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}

	var cond bool
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		clen = d.readContainerLen(msgpackContainerBin)
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		clen = d.readContainerLen(msgpackContainerStr)
	} else if bd == mpArray16 || bd == mpArray32 ||
		(bd >= mpFixArrayMin && bd <= mpFixArrayMax) {
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
	} else {
		halt.errorf("invalid byte descriptor for decoding bytes, got: 0x%x", d.bd)
	}

	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *msgpackDecDriverBytes) DecodeStringAsBytes() (out []byte, state dBytesAttachState) {
	out, state = d.DecodeBytes()
	if d.h.ValidateUnicode && !utf8.Valid(out) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", out)
	}
	return
}

func (d *msgpackDecDriverBytes) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *msgpackDecDriverBytes) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == mpNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *msgpackDecDriverBytes) TryNil() (v bool) {
	return d.advanceNil()
}

func (d *msgpackDecDriverBytes) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	bd := d.bd
	if bd == mpNil {
		d.bdRead = false
		return valueTypeNil
	} else if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		return valueTypeBytes
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		if d.h.WriteExt || d.h.RawToString {
			return valueTypeString
		}
		return valueTypeBytes
	} else if bd == mpArray16 || bd == mpArray32 || (bd >= mpFixArrayMin && bd <= mpFixArrayMax) {
		return valueTypeArray
	} else if bd == mpMap16 || bd == mpMap32 || (bd >= mpFixMapMin && bd <= mpFixMapMax) {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *msgpackDecDriverBytes) readContainerLen(ct msgpackContainerType) (clen int) {
	bd := d.bd
	if bd == ct.b8 {
		clen = int(d.r.readn1())
	} else if bd == ct.b16 {
		clen = int(bigen.Uint16(d.r.readn2()))
	} else if bd == ct.b32 {
		clen = int(bigen.Uint32(d.r.readn4()))
	} else if (ct.bFixMin & bd) == ct.bFixMin {
		clen = int(ct.bFixMin ^ bd)
	} else {
		halt.errorf("cannot read container length: %s: hex: %x, decimal: %d", msgBadDesc, bd, bd)
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverBytes) ReadMapStart() int {
	if d.advanceNil() {
		return containerLenNil
	}
	return d.readContainerLen(msgpackContainerMap)
}

func (d *msgpackDecDriverBytes) ReadArrayStart() int {
	if d.advanceNil() {
		return containerLenNil
	}
	return d.readContainerLen(msgpackContainerList)
}

func (d *msgpackDecDriverBytes) readExtLen() (clen int) {
	switch d.bd {
	case mpFixExt1:
		clen = 1
	case mpFixExt2:
		clen = 2
	case mpFixExt4:
		clen = 4
	case mpFixExt8:
		clen = 8
	case mpFixExt16:
		clen = 16
	case mpExt8:
		clen = int(d.r.readn1())
	case mpExt16:
		clen = int(bigen.Uint16(d.r.readn2()))
	case mpExt32:
		clen = int(bigen.Uint32(d.r.readn4()))
	default:
		halt.errorf("decoding ext bytes: found unexpected byte: %x", d.bd)
	}
	return
}

func (d *msgpackDecDriverBytes) DecodeTime() (t time.Time) {

	if d.advanceNil() {
		return
	}
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		clen = d.readContainerLen(msgpackContainerBin)
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		clen = d.readContainerLen(msgpackContainerStr)
	} else {

		d.bdRead = false
		b2 := d.r.readn1()
		if d.bd == mpFixExt4 && b2 == mpTimeExtTagU {
			clen = 4
		} else if d.bd == mpFixExt8 && b2 == mpTimeExtTagU {
			clen = 8
		} else if d.bd == mpExt8 && b2 == 12 && d.r.readn1() == mpTimeExtTagU {
			clen = 12
		} else {
			halt.errorf("invalid stream for decoding time as extension: got 0x%x, 0x%x", d.bd, b2)
		}
	}
	return d.decodeTime(clen)
}

func (d *msgpackDecDriverBytes) decodeTime(clen int) (t time.Time) {
	d.bdRead = false
	switch clen {
	case 4:
		t = time.Unix(int64(bigen.Uint32(d.r.readn4())), 0).UTC()
	case 8:
		tv := bigen.Uint64(d.r.readn8())
		t = time.Unix(int64(tv&0x00000003ffffffff), int64(tv>>34)).UTC()
	case 12:
		nsec := bigen.Uint32(d.r.readn4())
		sec := bigen.Uint64(d.r.readn8())
		t = time.Unix(int64(sec), int64(nsec)).UTC()
	default:
		halt.errorf("invalid length of bytes for decoding time - expecting 4 or 8 or 12, got %d", clen)
	}
	return
}

func (d *msgpackDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *msgpackDecDriverBytes) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *msgpackDecDriverBytes) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
	if xtagIn > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtagIn)
	}
	if d.advanceNil() {
		return
	}
	tag := uint8(xtagIn)
	xbd := d.bd
	if xbd == mpBin8 || xbd == mpBin16 || xbd == mpBin32 {
		xbs, bstate = d.DecodeBytes()
	} else if xbd == mpStr8 || xbd == mpStr16 || xbd == mpStr32 ||
		(xbd >= mpFixStrMin && xbd <= mpFixStrMax) {
		xbs, bstate = d.DecodeStringAsBytes()
	} else {
		clen := d.readExtLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag - got %b, expecting %v", xtag, tag)
		}
		xbs, ok = d.r.readxb(uint(clen))
		bstate = d.d.attachState(ok)

	}
	d.bdRead = false
	ok = true
	return
}

func (d *msgpackEncDriverBytes) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*MsgpackHandle)
	d.e = shared
	if shared.bytes {
		fp = msgpackFpEncBytes
	} else {
		fp = msgpackFpEncIO
	}

	d.init2(enc)
	return
}

func (e *msgpackEncDriverBytes) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *msgpackEncDriverBytes) writerEnd() { e.w.end() }

func (e *msgpackEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *msgpackEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *msgpackDecDriverBytes) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*MsgpackHandle)
	d.d = shared
	if shared.bytes {
		fp = msgpackFpDecBytes
	} else {
		fp = msgpackFpDecIO
	}

	d.init2(dec)
	return
}

func (d *msgpackDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *msgpackDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *msgpackDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *msgpackDecDriverBytes) descBd() string {
	return sprintf("%v (%s)", d.bd, mpdesc(d.bd))
}

func (d *msgpackDecDriverBytes) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}

type helperEncDriverMsgpackIO struct{}
type encFnMsgpackIO struct {
	i  encFnInfo
	fe func(*encoderMsgpackIO, *encFnInfo, reflect.Value)
}
type encRtidFnMsgpackIO struct {
	rtid uintptr
	fn   *encFnMsgpackIO
}
type encoderMsgpackIO struct {
	dh helperEncDriverMsgpackIO
	fp *fastpathEsMsgpackIO
	e  msgpackEncDriverIO
	encoderBase
}
type helperDecDriverMsgpackIO struct{}
type decFnMsgpackIO struct {
	i  decFnInfo
	fd func(*decoderMsgpackIO, *decFnInfo, reflect.Value)
}
type decRtidFnMsgpackIO struct {
	rtid uintptr
	fn   *decFnMsgpackIO
}
type decoderMsgpackIO struct {
	dh helperDecDriverMsgpackIO
	fp *fastpathDsMsgpackIO
	d  msgpackDecDriverIO
	decoderBase
}
type msgpackEncDriverIO struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *MsgpackHandle
	e *encoderBase
	w bufioEncWriter
}
type msgpackDecDriverIO struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	decInit2er

	h *MsgpackHandle
	d *decoderBase
	r ioDecReader

	bdAndBdread

	noBuiltInTypes
}

func (e *encoderMsgpackIO) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderMsgpackIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderMsgpackIO) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderMsgpackIO) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderMsgpackIO) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderMsgpackIO) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderMsgpackIO) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderMsgpackIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderMsgpackIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderMsgpackIO) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderMsgpackIO) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderMsgpackIO) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderMsgpackIO) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderMsgpackIO) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderMsgpackIO) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderMsgpackIO) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderMsgpackIO) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderMsgpackIO) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderMsgpackIO) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderMsgpackIO) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderMsgpackIO) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderMsgpackIO) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderMsgpackIO) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderMsgpackIO) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderMsgpackIO) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderMsgpackIO) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderMsgpackIO) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderMsgpackIO) kSeqFn(rt reflect.Type) (fn *encFnMsgpackIO) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderMsgpackIO) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnMsgpackIO
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

func (e *encoderMsgpackIO) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnMsgpackIO
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

func (e *encoderMsgpackIO) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderMsgpackIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderMsgpackIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderMsgpackIO) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderMsgpackIO) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderMsgpackIO) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderMsgpackIO) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderMsgpackIO) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnMsgpackIO

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

func (e *encoderMsgpackIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnMsgpackIO) {
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

func (e *encoderMsgpackIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsMsgpackIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderMsgpackIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderMsgpackIO) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderMsgpackIO) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderMsgpackIO) mustEncode(v interface{}) {
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

func (e *encoderMsgpackIO) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderMsgpackIO) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderMsgpackIO) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderMsgpackIO) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderMsgpackIO) encodeValue(rv reflect.Value, fn *encFnMsgpackIO) {

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

func (e *encoderMsgpackIO) encodeValueNonNil(rv reflect.Value, fn *encFnMsgpackIO) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderMsgpackIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderMsgpackIO) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderMsgpackIO) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderMsgpackIO) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderMsgpackIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderMsgpackIO) fn(t reflect.Type) *encFnMsgpackIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderMsgpackIO) fnNoExt(t reflect.Type) *encFnMsgpackIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderMsgpackIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderMsgpackIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderMsgpackIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderMsgpackIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderMsgpackIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderMsgpackIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderMsgpackIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderMsgpackIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverMsgpackIO) newEncoderBytes(out *[]byte, h Handle) *encoderMsgpackIO {
	var c1 encoderMsgpackIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverMsgpackIO) newEncoderIO(out io.Writer, h Handle) *encoderMsgpackIO {
	var c1 encoderMsgpackIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverMsgpackIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsMsgpackIO) (f *fastpathEMsgpackIO, u reflect.Type) {
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

func (helperEncDriverMsgpackIO) encFindRtidFn(s []encRtidFnMsgpackIO, rtid uintptr) (i uint, fn *encFnMsgpackIO) {

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

func (helperEncDriverMsgpackIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnMsgpackIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnMsgpackIO](v))
	}
	return
}

func (dh helperEncDriverMsgpackIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsMsgpackIO, checkExt bool) (fn *encFnMsgpackIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverMsgpackIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsMsgpackIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnMsgpackIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnMsgpackIO = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverMsgpackIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsMsgpackIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnMsgpackIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnMsgpackIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnMsgpackIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnMsgpackIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnMsgpackIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverMsgpackIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsMsgpackIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnMsgpackIO) {
	fn = new(encFnMsgpackIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderMsgpackIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderMsgpackIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderMsgpackIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderMsgpackIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderMsgpackIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderMsgpackIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderMsgpackIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderMsgpackIO).textMarshal
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
					fn.fe = func(e *encoderMsgpackIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderMsgpackIO).kBool
			case reflect.String:

				fn.fe = (*encoderMsgpackIO).kString
			case reflect.Int:
				fn.fe = (*encoderMsgpackIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderMsgpackIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderMsgpackIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderMsgpackIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderMsgpackIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderMsgpackIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderMsgpackIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderMsgpackIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderMsgpackIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderMsgpackIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderMsgpackIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderMsgpackIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderMsgpackIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderMsgpackIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderMsgpackIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderMsgpackIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderMsgpackIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderMsgpackIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderMsgpackIO).kStructSimple
				} else {
					fn.fe = (*encoderMsgpackIO).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderMsgpackIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderMsgpackIO).kErr
			default:

				fn.fe = (*encoderMsgpackIO).kErr
			}
		}
	}
	return
}
func (d *decoderMsgpackIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderMsgpackIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderMsgpackIO) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderMsgpackIO) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderMsgpackIO) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderMsgpackIO) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderMsgpackIO) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderMsgpackIO) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderMsgpackIO) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderMsgpackIO) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderMsgpackIO) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderMsgpackIO) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderMsgpackIO) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderMsgpackIO) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderMsgpackIO) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderMsgpackIO) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderMsgpackIO) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderMsgpackIO) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderMsgpackIO) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderMsgpackIO) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderMsgpackIO) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderMsgpackIO) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderMsgpackIO) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderMsgpackIO) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderMsgpackIO) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderMsgpackIO) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderMsgpackIO) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderMsgpackIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderMsgpackIO) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderMsgpackIO) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderMsgpackIO) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderMsgpackIO) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderMsgpackIO) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnMsgpackIO

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

func (d *decoderMsgpackIO) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnMsgpackIO
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

func (d *decoderMsgpackIO) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnMsgpackIO

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

func (d *decoderMsgpackIO) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnMsgpackIO
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

func (d *decoderMsgpackIO) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsMsgpackIO)

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

func (d *decoderMsgpackIO) reset() {
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

func (d *decoderMsgpackIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderMsgpackIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderMsgpackIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderMsgpackIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderMsgpackIO) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderMsgpackIO) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderMsgpackIO) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderMsgpackIO) Release() {}

func (d *decoderMsgpackIO) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderMsgpackIO) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderMsgpackIO) decode(iv interface{}) {
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

func (d *decoderMsgpackIO) decodeValue(rv reflect.Value, fn *decFnMsgpackIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderMsgpackIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnMsgpackIO) {

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

func (d *decoderMsgpackIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderMsgpackIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderMsgpackIO) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderMsgpackIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderMsgpackIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderMsgpackIO) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderMsgpackIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderMsgpackIO) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderMsgpackIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderMsgpackIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderMsgpackIO) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderMsgpackIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderMsgpackIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderMsgpackIO) fn(t reflect.Type) *decFnMsgpackIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderMsgpackIO) fnNoExt(t reflect.Type) *decFnMsgpackIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverMsgpackIO) newDecoderBytes(in []byte, h Handle) *decoderMsgpackIO {
	var c1 decoderMsgpackIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverMsgpackIO) newDecoderIO(in io.Reader, h Handle) *decoderMsgpackIO {
	var c1 decoderMsgpackIO
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverMsgpackIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsMsgpackIO) (f *fastpathDMsgpackIO, u reflect.Type) {
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

func (helperDecDriverMsgpackIO) decFindRtidFn(s []decRtidFnMsgpackIO, rtid uintptr) (i uint, fn *decFnMsgpackIO) {

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

func (helperDecDriverMsgpackIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnMsgpackIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnMsgpackIO](v))
	}
	return
}

func (dh helperDecDriverMsgpackIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsMsgpackIO,
	checkExt bool) (fn *decFnMsgpackIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverMsgpackIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsMsgpackIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnMsgpackIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnMsgpackIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverMsgpackIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsMsgpackIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnMsgpackIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnMsgpackIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnMsgpackIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnMsgpackIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnMsgpackIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverMsgpackIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsMsgpackIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnMsgpackIO) {
	fn = new(decFnMsgpackIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderMsgpackIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderMsgpackIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderMsgpackIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderMsgpackIO).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderMsgpackIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderMsgpackIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderMsgpackIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderMsgpackIO).textUnmarshal
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
						fn.fd = func(d *decoderMsgpackIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderMsgpackIO, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderMsgpackIO).kBool
			case reflect.String:
				fn.fd = (*decoderMsgpackIO).kString
			case reflect.Int:
				fn.fd = (*decoderMsgpackIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderMsgpackIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderMsgpackIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderMsgpackIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderMsgpackIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderMsgpackIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderMsgpackIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderMsgpackIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderMsgpackIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderMsgpackIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderMsgpackIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderMsgpackIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderMsgpackIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderMsgpackIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderMsgpackIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderMsgpackIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderMsgpackIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderMsgpackIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderMsgpackIO).kStructSimple
				} else {
					fn.fd = (*decoderMsgpackIO).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderMsgpackIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderMsgpackIO).kInterface
			default:

				fn.fd = (*decoderMsgpackIO).kErr
			}
		}
	}
	return
}
func (e *msgpackEncDriverIO) EncodeNil() {
	e.w.writen1(mpNil)
}

func (e *msgpackEncDriverIO) EncodeInt(i int64) {
	if e.h.PositiveIntUnsigned && i >= 0 {
		e.EncodeUint(uint64(i))
	} else if i > math.MaxInt8 {
		if i <= math.MaxInt16 {
			e.w.writen1(mpInt16)
			e.w.writen2(bigen.PutUint16(uint16(i)))
		} else if i <= math.MaxInt32 {
			e.w.writen1(mpInt32)
			e.w.writen4(bigen.PutUint32(uint32(i)))
		} else {
			e.w.writen1(mpInt64)
			e.w.writen8(bigen.PutUint64(uint64(i)))
		}
	} else if i >= -32 {
		if e.h.NoFixedNum {
			e.w.writen2(mpInt8, byte(i))
		} else {
			e.w.writen1(byte(i))
		}
	} else if i >= math.MinInt8 {
		e.w.writen2(mpInt8, byte(i))
	} else if i >= math.MinInt16 {
		e.w.writen1(mpInt16)
		e.w.writen2(bigen.PutUint16(uint16(i)))
	} else if i >= math.MinInt32 {
		e.w.writen1(mpInt32)
		e.w.writen4(bigen.PutUint32(uint32(i)))
	} else {
		e.w.writen1(mpInt64)
		e.w.writen8(bigen.PutUint64(uint64(i)))
	}
}

func (e *msgpackEncDriverIO) EncodeUint(i uint64) {
	if i <= math.MaxInt8 {
		if e.h.NoFixedNum {
			e.w.writen2(mpUint8, byte(i))
		} else {
			e.w.writen1(byte(i))
		}
	} else if i <= math.MaxUint8 {
		e.w.writen2(mpUint8, byte(i))
	} else if i <= math.MaxUint16 {
		e.w.writen1(mpUint16)
		e.w.writen2(bigen.PutUint16(uint16(i)))
	} else if i <= math.MaxUint32 {
		e.w.writen1(mpUint32)
		e.w.writen4(bigen.PutUint32(uint32(i)))
	} else {
		e.w.writen1(mpUint64)
		e.w.writen8(bigen.PutUint64(uint64(i)))
	}
}

func (e *msgpackEncDriverIO) EncodeBool(b bool) {
	if b {
		e.w.writen1(mpTrue)
	} else {
		e.w.writen1(mpFalse)
	}
}

func (e *msgpackEncDriverIO) EncodeFloat32(f float32) {
	e.w.writen1(mpFloat)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *msgpackEncDriverIO) EncodeFloat64(f float64) {
	e.w.writen1(mpDouble)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *msgpackEncDriverIO) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
		return
	}
	t = t.UTC()
	sec, nsec := t.Unix(), uint64(t.Nanosecond())
	var data64 uint64
	var l = 4
	if sec >= 0 && sec>>34 == 0 {
		data64 = (nsec << 34) | uint64(sec)
		if data64&0xffffffff00000000 != 0 {
			l = 8
		}
	} else {
		l = 12
	}
	if e.h.WriteExt {
		e.encodeExtPreamble(mpTimeExtTagU, l)
	} else {
		e.writeContainerLen(msgpackContainerRawLegacy, l)
	}
	switch l {
	case 4:
		e.w.writen4(bigen.PutUint32(uint32(data64)))
	case 8:
		e.w.writen8(bigen.PutUint64(data64))
	case 12:
		e.w.writen4(bigen.PutUint32(uint32(nsec)))
		e.w.writen8(bigen.PutUint64(uint64(sec)))
	}
}

func (e *msgpackEncDriverIO) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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
	if e.h.WriteExt {
		e.encodeExtPreamble(uint8(xtag), len(bs))
		e.w.writeb(bs)
	} else {
		e.EncodeBytes(bs)
	}
END:
	if ext == SelfExt {
		e.e.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			e.e.blist.put(bs0)
		}
	}
}

func (e *msgpackEncDriverIO) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *msgpackEncDriverIO) encodeExtPreamble(xtag byte, l int) {
	if l == 1 {
		e.w.writen2(mpFixExt1, xtag)
	} else if l == 2 {
		e.w.writen2(mpFixExt2, xtag)
	} else if l == 4 {
		e.w.writen2(mpFixExt4, xtag)
	} else if l == 8 {
		e.w.writen2(mpFixExt8, xtag)
	} else if l == 16 {
		e.w.writen2(mpFixExt16, xtag)
	} else if l < 256 {
		e.w.writen2(mpExt8, byte(l))
		e.w.writen1(xtag)
	} else if l < 65536 {
		e.w.writen1(mpExt16)
		e.w.writen2(bigen.PutUint16(uint16(l)))
		e.w.writen1(xtag)
	} else {
		e.w.writen1(mpExt32)
		e.w.writen4(bigen.PutUint32(uint32(l)))
		e.w.writen1(xtag)
	}
}

func (e *msgpackEncDriverIO) WriteArrayStart(length int) {
	e.writeContainerLen(msgpackContainerList, length)
}

func (e *msgpackEncDriverIO) WriteMapStart(length int) {
	e.writeContainerLen(msgpackContainerMap, length)
}

func (e *msgpackEncDriverIO) WriteArrayEmpty() {

	e.w.writen1(mpFixArrayMin)
}

func (e *msgpackEncDriverIO) WriteMapEmpty() {

	e.w.writen1(mpFixMapMin)
}

func (e *msgpackEncDriverIO) EncodeString(s string) {
	var ct msgpackContainerType
	if e.h.WriteExt {
		if e.h.StringToRaw {
			ct = msgpackContainerBin
		} else {
			ct = msgpackContainerStr
		}
	} else {
		ct = msgpackContainerRawLegacy
	}
	e.writeContainerLen(ct, len(s))
	if len(s) > 0 {
		e.w.writestr(s)
	}
}

func (e *msgpackEncDriverIO) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *msgpackEncDriverIO) EncodeStringBytesRaw(bs []byte) {
	if e.h.WriteExt {
		e.writeContainerLen(msgpackContainerBin, len(bs))
	} else {
		e.writeContainerLen(msgpackContainerRawLegacy, len(bs))
	}
	if len(bs) > 0 {
		e.w.writeb(bs)
	}
}

func (e *msgpackEncDriverIO) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *msgpackEncDriverIO) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = mpNil
	}
	e.w.writen1(v)
}

func (e *msgpackEncDriverIO) writeNilArray() {
	e.writeNilOr(mpFixArrayMin)
}

func (e *msgpackEncDriverIO) writeNilMap() {
	e.writeNilOr(mpFixMapMin)
}

func (e *msgpackEncDriverIO) writeNilBytes() {
	e.writeNilOr(mpFixStrMin)
}

func (e *msgpackEncDriverIO) writeContainerLen(ct msgpackContainerType, l int) {
	if ct.fixCutoff > 0 && l < int(ct.fixCutoff) {
		e.w.writen1(ct.bFixMin | byte(l))
	} else if ct.b8 > 0 && l < 256 {
		e.w.writen2(ct.b8, uint8(l))
	} else if l < 65536 {
		e.w.writen1(ct.b16)
		e.w.writen2(bigen.PutUint16(uint16(l)))
	} else {
		e.w.writen1(ct.b32)
		e.w.writen4(bigen.PutUint32(uint32(l)))
	}
}

func (d *msgpackDecDriverIO) DecodeNaked() {
	if !d.bdRead {
		d.readNextBd()
	}
	bd := d.bd
	n := d.d.naked()
	var decodeFurther bool

	switch bd {
	case mpNil:
		n.v = valueTypeNil
		d.bdRead = false
	case mpFalse:
		n.v = valueTypeBool
		n.b = false
	case mpTrue:
		n.v = valueTypeBool
		n.b = true

	case mpFloat:
		n.v = valueTypeFloat
		n.f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	case mpDouble:
		n.v = valueTypeFloat
		n.f = math.Float64frombits(bigen.Uint64(d.r.readn8()))

	case mpUint8:
		n.v = valueTypeUint
		n.u = uint64(d.r.readn1())
	case mpUint16:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint16(d.r.readn2()))
	case mpUint32:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint32(d.r.readn4()))
	case mpUint64:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint64(d.r.readn8()))

	case mpInt8:
		n.v = valueTypeInt
		n.i = int64(int8(d.r.readn1()))
	case mpInt16:
		n.v = valueTypeInt
		n.i = int64(int16(bigen.Uint16(d.r.readn2())))
	case mpInt32:
		n.v = valueTypeInt
		n.i = int64(int32(bigen.Uint32(d.r.readn4())))
	case mpInt64:
		n.v = valueTypeInt
		n.i = int64(int64(bigen.Uint64(d.r.readn8())))

	default:
		switch {
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:

			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:

			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd == mpStr8, bd == mpStr16, bd == mpStr32, bd >= mpFixStrMin && bd <= mpFixStrMax:
			d.d.fauxUnionReadRawBytes(d, d.h.WriteExt, d.h.RawToString)

		case bd == mpBin8, bd == mpBin16, bd == mpBin32:
			d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString)
		case bd == mpArray16, bd == mpArray32, bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			n.v = valueTypeArray
			decodeFurther = true
		case bd == mpMap16, bd == mpMap32, bd >= mpFixMapMin && bd <= mpFixMapMax:
			n.v = valueTypeMap
			decodeFurther = true
		case bd >= mpFixExt1 && bd <= mpFixExt16, bd >= mpExt8 && bd <= mpExt32:
			n.v = valueTypeExt
			clen := d.readExtLen()
			n.u = uint64(d.r.readn1())
			if n.u == uint64(mpTimeExtTagU) {
				n.v = valueTypeTime
				n.t = d.decodeTime(clen)
			} else {
				n.l = d.r.readx(uint(clen))
			}
		default:
			halt.errorf("cannot infer value: %s: Ox%x/%d/%s", msgBadDesc, bd, bd, mpdesc(bd))
		}
	}
	if !decodeFurther {
		d.bdRead = false
	}
	if n.v == valueTypeUint && d.h.SignedInteger {
		n.v = valueTypeInt
		n.i = int64(n.u)
	}
}

func (d *msgpackDecDriverIO) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *msgpackDecDriverIO) nextValueBytesBdReadR() {
	bd := d.bd

	var clen uint

	switch bd {
	case mpNil, mpFalse, mpTrue:
	case mpUint8, mpInt8:
		d.r.readn1()
	case mpUint16, mpInt16:
		d.r.skip(2)
	case mpFloat, mpUint32, mpInt32:
		d.r.skip(4)
	case mpDouble, mpUint64, mpInt64:
		d.r.skip(8)
	case mpStr8, mpBin8:
		clen = uint(d.r.readn1())
		d.r.skip(clen)
	case mpStr16, mpBin16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		d.r.skip(clen)
	case mpStr32, mpBin32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		d.r.skip(clen)
	case mpFixExt1:
		d.r.readn1()
		d.r.readn1()
	case mpFixExt2:
		d.r.readn1()
		d.r.skip(2)
	case mpFixExt4:
		d.r.readn1()
		d.r.skip(4)
	case mpFixExt8:
		d.r.readn1()
		d.r.skip(8)
	case mpFixExt16:
		d.r.readn1()
		d.r.skip(16)
	case mpExt8:
		clen = uint(d.r.readn1())
		d.r.readn1()
		d.r.skip(clen)
	case mpExt16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		d.r.readn1()
		d.r.skip(clen)
	case mpExt32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		d.r.readn1()
		d.r.skip(clen)
	case mpArray16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case mpArray32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case mpMap16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	case mpMap32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		for i := uint(0); i < clen; i++ {
			d.readNextBd()
			d.nextValueBytesBdReadR()
			d.readNextBd()
			d.nextValueBytesBdReadR()
		}
	default:
		switch {
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
		case bd >= mpFixStrMin && bd <= mpFixStrMax:
			clen = uint(mpFixStrMin ^ bd)
			d.r.skip(clen)
		case bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			clen = uint(mpFixArrayMin ^ bd)
			for i := uint(0); i < clen; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		case bd >= mpFixMapMin && bd <= mpFixMapMax:
			clen = uint(mpFixMapMin ^ bd)
			for i := uint(0); i < clen; i++ {
				d.readNextBd()
				d.nextValueBytesBdReadR()
				d.readNextBd()
				d.nextValueBytesBdReadR()
			}
		default:
			halt.errorf("nextValueBytes: cannot infer value: %s: Ox%x/%d/%s", msgBadDesc, bd, bd, mpdesc(bd))
		}
	}
	return
}

func (d *msgpackDecDriverIO) decFloat4Int32() (f float32) {
	fbits := bigen.Uint32(d.r.readn4())
	f = math.Float32frombits(fbits)
	if !noFrac32(fbits) {
		halt.errorf("assigning integer value from float32 with a fraction: %v", f)
	}
	return
}

func (d *msgpackDecDriverIO) decFloat4Int64() (f float64) {
	fbits := bigen.Uint64(d.r.readn8())
	f = math.Float64frombits(fbits)
	if !noFrac64(fbits) {
		halt.errorf("assigning integer value from float64 with a fraction: %v", f)
	}
	return
}

func (d *msgpackDecDriverIO) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	switch d.bd {
	case mpUint8:
		i = int64(uint64(d.r.readn1()))
	case mpUint16:
		i = int64(uint64(bigen.Uint16(d.r.readn2())))
	case mpUint32:
		i = int64(uint64(bigen.Uint32(d.r.readn4())))
	case mpUint64:
		i = int64(bigen.Uint64(d.r.readn8()))
	case mpInt8:
		i = int64(int8(d.r.readn1()))
	case mpInt16:
		i = int64(int16(bigen.Uint16(d.r.readn2())))
	case mpInt32:
		i = int64(int32(bigen.Uint32(d.r.readn4())))
	case mpInt64:
		i = int64(bigen.Uint64(d.r.readn8()))
	case mpFloat:
		i = int64(d.decFloat4Int32())
	case mpDouble:
		i = int64(d.decFloat4Int64())
	default:
		switch {
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			i = int64(int8(d.bd))
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			i = int64(int8(d.bd))
		default:
			halt.errorf("cannot decode signed integer: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
		}
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverIO) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	switch d.bd {
	case mpUint8:
		ui = uint64(d.r.readn1())
	case mpUint16:
		ui = uint64(bigen.Uint16(d.r.readn2()))
	case mpUint32:
		ui = uint64(bigen.Uint32(d.r.readn4()))
	case mpUint64:
		ui = bigen.Uint64(d.r.readn8())
	case mpInt8:
		if i := int64(int8(d.r.readn1())); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt16:
		if i := int64(int16(bigen.Uint16(d.r.readn2()))); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt32:
		if i := int64(int32(bigen.Uint32(d.r.readn4()))); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpInt64:
		if i := int64(bigen.Uint64(d.r.readn8())); i >= 0 {
			ui = uint64(i)
		} else {
			halt.errorf("assigning negative signed value: %v, to unsigned type", i)
		}
	case mpFloat:
		if f := d.decFloat4Int32(); f >= 0 {
			ui = uint64(f)
		} else {
			halt.errorf("assigning negative float value: %v, to unsigned type", f)
		}
	case mpDouble:
		if f := d.decFloat4Int64(); f >= 0 {
			ui = uint64(f)
		} else {
			halt.errorf("assigning negative float value: %v, to unsigned type", f)
		}
	default:
		switch {
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			ui = uint64(d.bd)
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			halt.errorf("assigning negative signed value: %v, to unsigned type", int(d.bd))
		default:
			halt.errorf("cannot decode unsigned integer: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
		}
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverIO) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	if d.bd == mpFloat {
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readn4())))
	} else if d.bd == mpDouble {
		f = math.Float64frombits(bigen.Uint64(d.r.readn8()))
	} else {
		f = float64(d.DecodeInt64())
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverIO) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == mpFalse || d.bd == 0 {

	} else if d.bd == mpTrue || d.bd == 1 {
		b = true
	} else {
		halt.errorf("cannot decode bool: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverIO) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}

	var cond bool
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		clen = d.readContainerLen(msgpackContainerBin)
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		clen = d.readContainerLen(msgpackContainerStr)
	} else if bd == mpArray16 || bd == mpArray32 ||
		(bd >= mpFixArrayMin && bd <= mpFixArrayMax) {
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
	} else {
		halt.errorf("invalid byte descriptor for decoding bytes, got: 0x%x", d.bd)
	}

	d.bdRead = false
	bs, cond = d.r.readxb(uint(clen))
	state = d.d.attachState(cond)
	return
}

func (d *msgpackDecDriverIO) DecodeStringAsBytes() (out []byte, state dBytesAttachState) {
	out, state = d.DecodeBytes()
	if d.h.ValidateUnicode && !utf8.Valid(out) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", out)
	}
	return
}

func (d *msgpackDecDriverIO) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *msgpackDecDriverIO) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == mpNil {
		d.bdRead = false
		return true
	}
	return
}

func (d *msgpackDecDriverIO) TryNil() (v bool) {
	return d.advanceNil()
}

func (d *msgpackDecDriverIO) ContainerType() (vt valueType) {
	if !d.bdRead {
		d.readNextBd()
	}
	bd := d.bd
	if bd == mpNil {
		d.bdRead = false
		return valueTypeNil
	} else if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		return valueTypeBytes
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		if d.h.WriteExt || d.h.RawToString {
			return valueTypeString
		}
		return valueTypeBytes
	} else if bd == mpArray16 || bd == mpArray32 || (bd >= mpFixArrayMin && bd <= mpFixArrayMax) {
		return valueTypeArray
	} else if bd == mpMap16 || bd == mpMap32 || (bd >= mpFixMapMin && bd <= mpFixMapMax) {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *msgpackDecDriverIO) readContainerLen(ct msgpackContainerType) (clen int) {
	bd := d.bd
	if bd == ct.b8 {
		clen = int(d.r.readn1())
	} else if bd == ct.b16 {
		clen = int(bigen.Uint16(d.r.readn2()))
	} else if bd == ct.b32 {
		clen = int(bigen.Uint32(d.r.readn4()))
	} else if (ct.bFixMin & bd) == ct.bFixMin {
		clen = int(ct.bFixMin ^ bd)
	} else {
		halt.errorf("cannot read container length: %s: hex: %x, decimal: %d", msgBadDesc, bd, bd)
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriverIO) ReadMapStart() int {
	if d.advanceNil() {
		return containerLenNil
	}
	return d.readContainerLen(msgpackContainerMap)
}

func (d *msgpackDecDriverIO) ReadArrayStart() int {
	if d.advanceNil() {
		return containerLenNil
	}
	return d.readContainerLen(msgpackContainerList)
}

func (d *msgpackDecDriverIO) readExtLen() (clen int) {
	switch d.bd {
	case mpFixExt1:
		clen = 1
	case mpFixExt2:
		clen = 2
	case mpFixExt4:
		clen = 4
	case mpFixExt8:
		clen = 8
	case mpFixExt16:
		clen = 16
	case mpExt8:
		clen = int(d.r.readn1())
	case mpExt16:
		clen = int(bigen.Uint16(d.r.readn2()))
	case mpExt32:
		clen = int(bigen.Uint32(d.r.readn4()))
	default:
		halt.errorf("decoding ext bytes: found unexpected byte: %x", d.bd)
	}
	return
}

func (d *msgpackDecDriverIO) DecodeTime() (t time.Time) {

	if d.advanceNil() {
		return
	}
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		clen = d.readContainerLen(msgpackContainerBin)
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		clen = d.readContainerLen(msgpackContainerStr)
	} else {

		d.bdRead = false
		b2 := d.r.readn1()
		if d.bd == mpFixExt4 && b2 == mpTimeExtTagU {
			clen = 4
		} else if d.bd == mpFixExt8 && b2 == mpTimeExtTagU {
			clen = 8
		} else if d.bd == mpExt8 && b2 == 12 && d.r.readn1() == mpTimeExtTagU {
			clen = 12
		} else {
			halt.errorf("invalid stream for decoding time as extension: got 0x%x, 0x%x", d.bd, b2)
		}
	}
	return d.decodeTime(clen)
}

func (d *msgpackDecDriverIO) decodeTime(clen int) (t time.Time) {
	d.bdRead = false
	switch clen {
	case 4:
		t = time.Unix(int64(bigen.Uint32(d.r.readn4())), 0).UTC()
	case 8:
		tv := bigen.Uint64(d.r.readn8())
		t = time.Unix(int64(tv&0x00000003ffffffff), int64(tv>>34)).UTC()
	case 12:
		nsec := bigen.Uint32(d.r.readn4())
		sec := bigen.Uint64(d.r.readn8())
		t = time.Unix(int64(sec), int64(nsec)).UTC()
	default:
		halt.errorf("invalid length of bytes for decoding time - expecting 4 or 8 or 12, got %d", clen)
	}
	return
}

func (d *msgpackDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *msgpackDecDriverIO) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *msgpackDecDriverIO) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
	if xtagIn > 0xff {
		halt.errorf("ext: tag must be <= 0xff; got: %v", xtagIn)
	}
	if d.advanceNil() {
		return
	}
	tag := uint8(xtagIn)
	xbd := d.bd
	if xbd == mpBin8 || xbd == mpBin16 || xbd == mpBin32 {
		xbs, bstate = d.DecodeBytes()
	} else if xbd == mpStr8 || xbd == mpStr16 || xbd == mpStr32 ||
		(xbd >= mpFixStrMin && xbd <= mpFixStrMax) {
		xbs, bstate = d.DecodeStringAsBytes()
	} else {
		clen := d.readExtLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag {
			halt.errorf("wrong extension tag - got %b, expecting %v", xtag, tag)
		}
		xbs, ok = d.r.readxb(uint(clen))
		bstate = d.d.attachState(ok)

	}
	d.bdRead = false
	ok = true
	return
}

func (d *msgpackEncDriverIO) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*MsgpackHandle)
	d.e = shared
	if shared.bytes {
		fp = msgpackFpEncBytes
	} else {
		fp = msgpackFpEncIO
	}

	d.init2(enc)
	return
}

func (e *msgpackEncDriverIO) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *msgpackEncDriverIO) writerEnd() { e.w.end() }

func (e *msgpackEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *msgpackEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *msgpackDecDriverIO) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*MsgpackHandle)
	d.d = shared
	if shared.bytes {
		fp = msgpackFpDecBytes
	} else {
		fp = msgpackFpDecIO
	}

	d.init2(dec)
	return
}

func (d *msgpackDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *msgpackDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *msgpackDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *msgpackDecDriverIO) descBd() string {
	return sprintf("%v (%s)", d.bd, mpdesc(d.bd))
}

func (d *msgpackDecDriverIO) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
