//go:build !notmono && !codec.notmono 

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"

	"encoding/base64"
	"io"
	"math"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

type helperEncDriverJsonBytes struct{}
type encFnJsonBytes struct {
	i  encFnInfo
	fe func(*encoderJsonBytes, *encFnInfo, reflect.Value)
}
type encRtidFnJsonBytes struct {
	rtid uintptr
	fn   *encFnJsonBytes
}
type encoderJsonBytes struct {
	dh helperEncDriverJsonBytes
	fp *fastpathEsJsonBytes
	e  jsonEncDriverBytes
	encoderBase
}
type helperDecDriverJsonBytes struct{}
type decFnJsonBytes struct {
	i  decFnInfo
	fd func(*decoderJsonBytes, *decFnInfo, reflect.Value)
}
type decRtidFnJsonBytes struct {
	rtid uintptr
	fn   *decFnJsonBytes
}
type decoderJsonBytes struct {
	dh helperDecDriverJsonBytes
	fp *fastpathDsJsonBytes
	d  jsonDecDriverBytes
	decoderBase
}
type jsonEncDriverBytes struct {
	noBuiltInTypes
	h *JsonHandle
	e *encoderBase
	s *bitset256

	w bytesEncAppender

	enc encoderI

	timeFmtLayout string
	byteFmter     jsonBytesFmter

	timeFmt  jsonTimeFmt
	bytesFmt jsonBytesFmt

	di int8
	d  bool
	dl uint16

	ks bool
	is byte

	typical bool

	rawext bool

	b [48]byte
}
type jsonDecDriverBytes struct {
	noBuiltInTypes
	decDriverNoopNumberHelper
	h *JsonHandle
	d *decoderBase

	r bytesDecReader

	buf []byte

	tok  uint8
	_    bool
	_    byte
	bstr [4]byte

	jsonHandleOpts

	dec decoderI
}

func (e *encoderJsonBytes) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderJsonBytes) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderJsonBytes) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderJsonBytes) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderJsonBytes) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderJsonBytes) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderJsonBytes) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderJsonBytes) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderJsonBytes) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderJsonBytes) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderJsonBytes) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderJsonBytes) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderJsonBytes) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderJsonBytes) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderJsonBytes) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderJsonBytes) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderJsonBytes) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderJsonBytes) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderJsonBytes) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderJsonBytes) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderJsonBytes) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderJsonBytes) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderJsonBytes) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderJsonBytes) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderJsonBytes) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderJsonBytes) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderJsonBytes) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderJsonBytes) kSeqFn(rt reflect.Type) (fn *encFnJsonBytes) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderJsonBytes) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnJsonBytes
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

func (e *encoderJsonBytes) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnJsonBytes
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

func (e *encoderJsonBytes) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderJsonBytes) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderJsonBytes) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderJsonBytes) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderJsonBytes) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderJsonBytes) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderJsonBytes) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderJsonBytes) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnJsonBytes

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

func (e *encoderJsonBytes) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnJsonBytes) {
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

func (e *encoderJsonBytes) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsJsonBytes)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderJsonBytes) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderJsonBytes) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderJsonBytes) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderJsonBytes) mustEncode(v interface{}) {
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

func (e *encoderJsonBytes) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderJsonBytes) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderJsonBytes) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderJsonBytes) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderJsonBytes) encodeValue(rv reflect.Value, fn *encFnJsonBytes) {

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

func (e *encoderJsonBytes) encodeValueNonNil(rv reflect.Value, fn *encFnJsonBytes) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderJsonBytes) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderJsonBytes) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderJsonBytes) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderJsonBytes) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderJsonBytes) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderJsonBytes) fn(t reflect.Type) *encFnJsonBytes {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderJsonBytes) fnNoExt(t reflect.Type) *encFnJsonBytes {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderJsonBytes) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderJsonBytes) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderJsonBytes) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderJsonBytes) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderJsonBytes) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderJsonBytes) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderJsonBytes) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderJsonBytes) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverJsonBytes) newEncoderBytes(out *[]byte, h Handle) *encoderJsonBytes {
	var c1 encoderJsonBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverJsonBytes) newEncoderIO(out io.Writer, h Handle) *encoderJsonBytes {
	var c1 encoderJsonBytes
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverJsonBytes) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsJsonBytes) (f *fastpathEJsonBytes, u reflect.Type) {
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

func (helperEncDriverJsonBytes) encFindRtidFn(s []encRtidFnJsonBytes, rtid uintptr) (i uint, fn *encFnJsonBytes) {

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

func (helperEncDriverJsonBytes) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnJsonBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnJsonBytes](v))
	}
	return
}

func (dh helperEncDriverJsonBytes) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsJsonBytes, checkExt bool) (fn *encFnJsonBytes) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverJsonBytes) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsJsonBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnJsonBytes) {
	rtid := rt2id(rt)
	var sp []encRtidFnJsonBytes = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverJsonBytes) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsJsonBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnJsonBytes) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnJsonBytes
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnJsonBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnJsonBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnJsonBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverJsonBytes) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsJsonBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnJsonBytes) {
	fn = new(encFnJsonBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderJsonBytes).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderJsonBytes).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderJsonBytes).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderJsonBytes).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderJsonBytes).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderJsonBytes).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderJsonBytes).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderJsonBytes).textMarshal
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
					fn.fe = func(e *encoderJsonBytes, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderJsonBytes).kBool
			case reflect.String:

				fn.fe = (*encoderJsonBytes).kString
			case reflect.Int:
				fn.fe = (*encoderJsonBytes).kInt
			case reflect.Int8:
				fn.fe = (*encoderJsonBytes).kInt8
			case reflect.Int16:
				fn.fe = (*encoderJsonBytes).kInt16
			case reflect.Int32:
				fn.fe = (*encoderJsonBytes).kInt32
			case reflect.Int64:
				fn.fe = (*encoderJsonBytes).kInt64
			case reflect.Uint:
				fn.fe = (*encoderJsonBytes).kUint
			case reflect.Uint8:
				fn.fe = (*encoderJsonBytes).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderJsonBytes).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderJsonBytes).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderJsonBytes).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderJsonBytes).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderJsonBytes).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderJsonBytes).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderJsonBytes).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderJsonBytes).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderJsonBytes).kChan
			case reflect.Slice:
				fn.fe = (*encoderJsonBytes).kSlice
			case reflect.Array:
				fn.fe = (*encoderJsonBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderJsonBytes).kStructSimple
				} else {
					fn.fe = (*encoderJsonBytes).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderJsonBytes).kMap
			case reflect.Interface:

				fn.fe = (*encoderJsonBytes).kErr
			default:

				fn.fe = (*encoderJsonBytes).kErr
			}
		}
	}
	return
}
func (d *decoderJsonBytes) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderJsonBytes) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderJsonBytes) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderJsonBytes) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderJsonBytes) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderJsonBytes) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderJsonBytes) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderJsonBytes) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderJsonBytes) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderJsonBytes) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderJsonBytes) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderJsonBytes) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderJsonBytes) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderJsonBytes) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderJsonBytes) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderJsonBytes) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderJsonBytes) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderJsonBytes) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderJsonBytes) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderJsonBytes) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderJsonBytes) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderJsonBytes) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderJsonBytes) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderJsonBytes) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderJsonBytes) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderJsonBytes) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderJsonBytes) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderJsonBytes) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderJsonBytes) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderJsonBytes) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderJsonBytes) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderJsonBytes) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderJsonBytes) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnJsonBytes

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

func (d *decoderJsonBytes) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnJsonBytes
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

func (d *decoderJsonBytes) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnJsonBytes

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

func (d *decoderJsonBytes) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnJsonBytes
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

func (d *decoderJsonBytes) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsJsonBytes)

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

func (d *decoderJsonBytes) reset() {
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

func (d *decoderJsonBytes) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderJsonBytes) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderJsonBytes) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderJsonBytes) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderJsonBytes) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderJsonBytes) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderJsonBytes) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderJsonBytes) Release() {}

func (d *decoderJsonBytes) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderJsonBytes) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderJsonBytes) decode(iv interface{}) {
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

func (d *decoderJsonBytes) decodeValue(rv reflect.Value, fn *decFnJsonBytes) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderJsonBytes) decodeValueNoCheckNil(rv reflect.Value, fn *decFnJsonBytes) {

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

func (d *decoderJsonBytes) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderJsonBytes) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderJsonBytes) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderJsonBytes) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderJsonBytes) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderJsonBytes) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderJsonBytes) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderJsonBytes) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderJsonBytes) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderJsonBytes) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderJsonBytes) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderJsonBytes) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderJsonBytes) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderJsonBytes) fn(t reflect.Type) *decFnJsonBytes {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderJsonBytes) fnNoExt(t reflect.Type) *decFnJsonBytes {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverJsonBytes) newDecoderBytes(in []byte, h Handle) *decoderJsonBytes {
	var c1 decoderJsonBytes
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverJsonBytes) newDecoderIO(in io.Reader, h Handle) *decoderJsonBytes {
	var c1 decoderJsonBytes
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverJsonBytes) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsJsonBytes) (f *fastpathDJsonBytes, u reflect.Type) {
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

func (helperDecDriverJsonBytes) decFindRtidFn(s []decRtidFnJsonBytes, rtid uintptr) (i uint, fn *decFnJsonBytes) {

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

func (helperDecDriverJsonBytes) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnJsonBytes) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnJsonBytes](v))
	}
	return
}

func (dh helperDecDriverJsonBytes) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsJsonBytes,
	checkExt bool) (fn *decFnJsonBytes) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverJsonBytes) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsJsonBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnJsonBytes) {
	rtid := rt2id(rt)
	var sp []decRtidFnJsonBytes = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverJsonBytes) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsJsonBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnJsonBytes) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnJsonBytes
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnJsonBytes{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnJsonBytes, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnJsonBytes{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverJsonBytes) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsJsonBytes,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnJsonBytes) {
	fn = new(decFnJsonBytes)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderJsonBytes).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderJsonBytes).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderJsonBytes).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderJsonBytes).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderJsonBytes).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderJsonBytes).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderJsonBytes).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderJsonBytes).textUnmarshal
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
						fn.fd = func(d *decoderJsonBytes, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderJsonBytes, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderJsonBytes).kBool
			case reflect.String:
				fn.fd = (*decoderJsonBytes).kString
			case reflect.Int:
				fn.fd = (*decoderJsonBytes).kInt
			case reflect.Int8:
				fn.fd = (*decoderJsonBytes).kInt8
			case reflect.Int16:
				fn.fd = (*decoderJsonBytes).kInt16
			case reflect.Int32:
				fn.fd = (*decoderJsonBytes).kInt32
			case reflect.Int64:
				fn.fd = (*decoderJsonBytes).kInt64
			case reflect.Uint:
				fn.fd = (*decoderJsonBytes).kUint
			case reflect.Uint8:
				fn.fd = (*decoderJsonBytes).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderJsonBytes).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderJsonBytes).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderJsonBytes).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderJsonBytes).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderJsonBytes).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderJsonBytes).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderJsonBytes).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderJsonBytes).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderJsonBytes).kChan
			case reflect.Slice:
				fn.fd = (*decoderJsonBytes).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderJsonBytes).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderJsonBytes).kStructSimple
				} else {
					fn.fd = (*decoderJsonBytes).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderJsonBytes).kMap
			case reflect.Interface:

				fn.fd = (*decoderJsonBytes).kInterface
			default:

				fn.fd = (*decoderJsonBytes).kErr
			}
		}
	}
	return
}
func (e *jsonEncDriverBytes) writeIndent() {
	e.w.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.di < 0 {
		x = -x
		for x > len(jsonTabs) {
			e.w.writeb(jsonTabs[:])
			x -= len(jsonTabs)
		}
		e.w.writeb(jsonTabs[:x])
	} else {
		for x > len(jsonSpaces) {
			e.w.writeb(jsonSpaces[:])
			x -= len(jsonSpaces)
		}
		e.w.writeb(jsonSpaces[:x])
	}
}

func (e *jsonEncDriverBytes) WriteArrayElem(firstTime bool) {
	if !firstTime {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriverBytes) WriteMapElemKey(firstTime bool) {
	if !firstTime {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriverBytes) WriteMapElemValue() {
	if e.d {
		e.w.writen2(':', ' ')
	} else {
		e.w.writen1(':')
	}
}

func (e *jsonEncDriverBytes) EncodeNil() {

	e.w.writeb(jsonNull)
}

func (e *jsonEncDriverBytes) encodeIntAsUint(v int64, quotes bool) {
	neg := v < 0
	if neg {
		v = -v
	}
	e.encodeUint(neg, quotes, uint64(v))
}

func (e *jsonEncDriverBytes) EncodeTime(t time.Time) {

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	switch e.timeFmt {
	case jsonTimeFmtStringLayout:
		e.b[0] = '"'
		b := t.AppendFormat(e.b[1:1], e.timeFmtLayout)
		e.b[len(b)+1] = '"'
		e.w.writeb(e.b[:len(b)+2])
	case jsonTimeFmtUnix:
		e.encodeIntAsUint(t.Unix(), false)
	case jsonTimeFmtUnixMilli:
		e.encodeIntAsUint(t.UnixMilli(), false)
	case jsonTimeFmtUnixMicro:
		e.encodeIntAsUint(t.UnixMicro(), false)
	case jsonTimeFmtUnixNano:
		e.encodeIntAsUint(t.UnixNano(), false)
	}
}

func (e *jsonEncDriverBytes) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.writeNilBytes()
	} else {
		e.enc.encodeI(v)
	}
}

func (e *jsonEncDriverBytes) EncodeRawExt(re *RawExt) {
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encodeI(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *jsonEncDriverBytes) EncodeBool(b bool) {
	e.w.writestr(jsonEncBoolStrs[bool2int(e.ks && e.e.c == containerMapKey)%2][bool2int(b)%2])
}

func (e *jsonEncDriverBytes) encodeFloat(f float64, bitsize, fmt byte, prec int8) {
	var blen uint
	if e.ks && e.e.c == containerMapKey {
		blen = 2 + uint(len(strconv.AppendFloat(e.b[1:1], f, fmt, int(prec), int(bitsize))))

		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
	} else {
		e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), int(bitsize)))
	}
}

func (e *jsonEncDriverBytes) EncodeFloat64(f float64) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.encodeFloat(f, 64, fmt, prec)
}

func (e *jsonEncDriverBytes) EncodeFloat32(f float32) {
	if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.encodeFloat(float64(f), 32, fmt, prec)
}

func (e *jsonEncDriverBytes) encodeUint(neg bool, quotes bool, u uint64) {
	e.w.writeb(jsonEncodeUint(neg, quotes, u, &e.b))
}

func (e *jsonEncDriverBytes) EncodeInt(v int64) {
	quotes := e.is == 'A' || e.is == 'L' && (v > 1<<53 || v < -(1<<53)) ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		if quotes {
			blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
		}
		return
	}

	if v < 0 {
		e.encodeUint(true, quotes, uint64(-v))
	} else {
		e.encodeUint(false, quotes, uint64(v))
	}
}

func (e *jsonEncDriverBytes) EncodeUint(v uint64) {
	quotes := e.is == 'A' || e.is == 'L' && v > 1<<53 ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {

		if quotes {
			blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
		}
		return
	}

	e.encodeUint(false, quotes, v)
}

func (e *jsonEncDriverBytes) EncodeString(v string) {
	if e.h.StringToRaw {
		e.EncodeStringBytesRaw(bytesView(v))
		return
	}
	e.quoteStr(v)
}

func (e *jsonEncDriverBytes) EncodeStringNoEscape4Json(v string) { e.w.writeqstr(v) }

func (e *jsonEncDriverBytes) EncodeStringBytesRaw(v []byte) {
	if e.rawext {

		iv := e.h.RawBytesExt.ConvertExt(any(v))
		if iv == nil {
			e.EncodeNil()
		} else {
			e.enc.encodeI(iv)
		}
		return
	}

	if e.bytesFmt == jsonBytesFmtArray {
		e.WriteArrayStart(len(v))
		for j := range v {
			e.WriteArrayElem(j == 0)
			e.encodeUint(false, false, uint64(v[j]))
		}
		e.WriteArrayEnd()
		return
	}

	var slen int
	if e.bytesFmt == jsonBytesFmtBase64 {
		slen = base64.StdEncoding.EncodedLen(len(v))
	} else {
		slen = e.byteFmter.EncodedLen(len(v))
	}
	slen += 2

	bs := e.e.blist.peek(slen, false)[:slen]

	if e.bytesFmt == jsonBytesFmtBase64 {
		base64.StdEncoding.Encode(bs[1:], v)
	} else {
		e.byteFmter.Encode(bs[1:], v)
	}

	bs[len(bs)-1] = '"'
	bs[0] = '"'
	e.w.writeb(bs)
}

func (e *jsonEncDriverBytes) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *jsonEncDriverBytes) writeNilOr(v []byte) {
	if !e.h.NilCollectionToZeroLength {
		v = jsonNull
	}
	e.w.writeb(v)
}

func (e *jsonEncDriverBytes) writeNilBytes() {
	e.writeNilOr(jsonArrayEmpty)
}

func (e *jsonEncDriverBytes) writeNilArray() {
	e.writeNilOr(jsonArrayEmpty)
}

func (e *jsonEncDriverBytes) writeNilMap() {
	e.writeNilOr(jsonMapEmpty)
}

func (e *jsonEncDriverBytes) WriteArrayEmpty() {
	e.w.writen2('[', ']')
}

func (e *jsonEncDriverBytes) WriteMapEmpty() {
	e.w.writen2('{', '}')
}

func (e *jsonEncDriverBytes) WriteArrayStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('[')
}

func (e *jsonEncDriverBytes) WriteArrayEnd() {
	if e.d {
		e.dl--

		e.writeIndent()
	}
	e.w.writen1(']')
}

func (e *jsonEncDriverBytes) WriteMapStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('{')
}

func (e *jsonEncDriverBytes) WriteMapEnd() {
	if e.d {
		e.dl--

		e.writeIndent()
	}
	e.w.writen1('}')
}

func (e *jsonEncDriverBytes) quoteStr(s string) {

	const hex = "0123456789abcdef"
	e.w.writen1('"')
	var i, start uint
	for i < uint(len(s)) {

		b := s[i]
		if e.s.isset(b) {
			i++
			continue
		}
		if b < utf8.RuneSelf {
			if start < i {
				e.w.writestr(s[start:i])
			}
			switch b {
			case '\\':
				e.w.writen2('\\', '\\')
			case '"':
				e.w.writen2('\\', '"')
			case '\n':
				e.w.writen2('\\', 'n')
			case '\t':
				e.w.writen2('\\', 't')
			case '\r':
				e.w.writen2('\\', 'r')
			case '\b':
				e.w.writen2('\\', 'b')
			case '\f':
				e.w.writen2('\\', 'f')
			default:
				e.w.writestr(`\u00`)
				e.w.writen2(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\uFFFD`)
			i++
			start = i
			continue
		}

		if jsonEscapeMultiByteUnicodeSep && (c == '\u2028' || c == '\u2029') {
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\u202`)
			e.w.writen1(hex[c&0xF])
			i += uint(size)
			start = i
			continue
		}
		i += uint(size)
	}
	if start < uint(len(s)) {
		e.w.writestr(s[start:])
	}
	e.w.writen1('"')
}

func (e *jsonEncDriverBytes) atEndOfEncode() {
	if e.h.TermWhitespace {
		var c byte = ' '
		if e.e.c != 0 {
			c = '\n'
		}
		e.w.writen1(c)
	}
}

func (d *jsonDecDriverBytes) ReadMapStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '{' {
		halt.errorByte("read map - expect char '{' but got char: ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriverBytes) ReadArrayStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '[' {
		halt.errorByte("read array - expect char '[' but got char ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriverBytes) CheckBreak() bool {
	d.advance()
	return d.tok == '}' || d.tok == ']'
}

func (d *jsonDecDriverBytes) checkSep(xc byte) {
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriverBytes) ReadArrayElem(firstTime bool) {
	if !firstTime {
		d.checkSep(',')
	}
}

func (d *jsonDecDriverBytes) ReadArrayEnd() {
	d.checkSep(']')
}

func (d *jsonDecDriverBytes) ReadMapElemKey(firstTime bool) {
	d.ReadArrayElem(firstTime)
}

func (d *jsonDecDriverBytes) ReadMapElemValue() {
	d.checkSep(':')
}

func (d *jsonDecDriverBytes) ReadMapEnd() {
	d.checkSep('}')
}

func (d *jsonDecDriverBytes) readDelimError(xc uint8) {
	halt.errorf("read json delimiter - expect char '%c' but got char '%c'", xc, d.tok)
}

func (d *jsonDecDriverBytes) checkLit3(got, expect [3]byte) {
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr3(got, expect)
	}
	d.tok = 0
}

func (d *jsonDecDriverBytes) checkLit4(got, expect [4]byte) {
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr4(got, expect)
	}
	d.tok = 0
}

func (d *jsonDecDriverBytes) skipWhitespace() {
	d.tok = d.r.skipWhitespace()
}

func (d *jsonDecDriverBytes) advance() {

	if d.tok < 33 {
		d.skipWhitespace()
	}
}

func (d *jsonDecDriverBytes) nextValueBytes() []byte {
	consumeString := func() {
	TOP:
		_, c := d.r.jsonReadAsisChars()
		if c == '\\' {
			d.r.readn1()
			goto TOP
		}
	}

	d.advance()
	d.r.startRecording()

	switch d.tok {
	default:
		_, d.tok = d.r.jsonReadNum()

		if d.tok != 0 {
			vv := d.r.stopRecording()
			return vv[:len(vv)-1]
		}
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
	case '"':
		consumeString()
		d.tok = 0
	case '{', '[':
		var elem struct{}
		var stack []struct{}

		stack = append(stack, elem)

		for len(stack) != 0 {
			c := d.r.readn1()
			switch c {
			case '"':
				consumeString()
			case '{', '[':
				stack = append(stack, elem)
			case '}', ']':
				stack = stack[:len(stack)-1]
			}
		}
		d.tok = 0
	}
	return d.r.stopRecording()
}

func (d *jsonDecDriverBytes) TryNil() bool {
	d.advance()

	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriverBytes) DecodeBool() (v bool) {
	d.advance()

	fquot := d.d.c == containerMapKey && d.tok == '"'
	if fquot {
		d.tok = d.r.readn1()
	}
	switch d.tok {
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())

	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		v = true
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())

	default:
		halt.errorByte("decode bool: got first char: ", d.tok)

	}
	if fquot {
		d.r.readn1()
	}
	return
}

func (d *jsonDecDriverBytes) DecodeTime() (t time.Time) {

	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	var bs []byte

	if d.tok != '"' {
		bs, d.tok = d.r.jsonReadNum()
		i := d.parseInt64(bs)
		switch d.timeFmtNum {
		case jsonTimeFmtUnix:
			t = time.Unix(i, 0)
		case jsonTimeFmtUnixMilli:
			t = time.UnixMilli(i)
		case jsonTimeFmtUnixMicro:
			t = time.UnixMicro(i)
		case jsonTimeFmtUnixNano:
			t = time.Unix(0, i)
		default:
			halt.errorStr("invalid timeFmtNum")
		}
		return
	}

	bs = d.readUnescapedString()
	var err error
	for _, v := range d.timeFmtLayouts {
		t, err = time.Parse(v, stringView(bs))
		if err == nil {
			return
		}
	}
	halt.errorStr("error decoding time")
	return
}

func (d *jsonDecDriverBytes) ContainerType() (vt valueType) {

	d.advance()

	if d.tok == '{' {
		return valueTypeMap
	} else if d.tok == '[' {
		return valueTypeArray
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return valueTypeNil
	} else if d.tok == '"' {
		return valueTypeString
	}
	return valueTypeUnset
}

func (d *jsonDecDriverBytes) decNumBytes() (bs []byte) {
	d.advance()
	if d.tok == '"' {
		bs = d.r.jsonReadUntilDblQuote()
		d.tok = 0
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	} else {
		bs, d.tok = d.r.jsonReadNum()
	}
	return
}

func (d *jsonDecDriverBytes) DecodeUint64() (u uint64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if neg {
		halt.errorf("negative number cannot be decoded as uint64: %s", any(b))
	}
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseUint"))
	}
	return
}

func (d *jsonDecDriverBytes) DecodeInt64() (v int64) {
	return d.parseInt64(d.decNumBytes())
}

func (d *jsonDecDriverBytes) parseInt64(b []byte) (v int64) {
	u, neg, ok := parseInteger_bytes(b)
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseInt"))
	}
	if chkOvf.Uint2Int(u, neg) {
		halt.errorBytes("overflow decoding number from ", b)
	}
	if neg {
		v = -int64(u)
	} else {
		v = int64(u)
	}
	return
}

func (d *jsonDecDriverBytes) DecodeFloat64() (f float64) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat64(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriverBytes) DecodeFloat32() (f float32) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat32(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriverBytes) advanceNil() (ok bool) {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriverBytes) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if d.advanceNil() {
		return
	}
	if ext == SelfExt {
		d.dec.decodeAs(rv, basetype, false)
	} else {
		d.dec.interfaceExtConvertAndDecode(rv, ext)
	}
}

func (d *jsonDecDriverBytes) DecodeRawExt(re *RawExt) {
	if d.advanceNil() {
		return
	}
	d.dec.decode(&re.Value)
}

func (d *jsonDecDriverBytes) decBytesFromArray(bs []byte) []byte {
	d.advance()
	if d.tok != ']' {
		bs = append(bs, uint8(d.DecodeUint64()))
		d.advance()
	}
	for d.tok != ']' {
		if d.tok != ',' {
			halt.errorByte("read array element - expect char ',' but got char: ", d.tok)
		}
		d.tok = 0
		bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		d.advance()
	}
	d.tok = 0
	return bs
}

func (d *jsonDecDriverBytes) DecodeBytes() (bs []byte, state dBytesAttachState) {
	d.advance()
	state = dBytesDetach
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	state = dBytesAttachBuffer

	if d.rawext {
		d.buf = d.buf[:0]
		d.dec.interfaceExtConvertAndDecode(&d.buf, d.h.RawBytesExt)
		bs = d.buf
		return
	}

	if d.tok == '[' {
		d.tok = 0

		bs = d.decBytesFromArray(d.buf[:0])
		d.buf = bs
		return
	}

	d.ensureReadingString()
	bs1 := d.readUnescapedString()

	slen := base64.StdEncoding.DecodedLen(len(bs1))
	if slen == 0 {
		bs = zeroByteSlice
		state = dBytesDetach
	} else if slen <= cap(d.buf) {
		bs = d.buf[:slen]
	} else {
		d.buf = d.d.blist.putGet(d.buf, slen)[:slen]
		bs = d.buf
	}
	var err error
	for _, v := range d.byteFmters {

		slen, err = v.Decode(bs, bs1)
		if err == nil {
			bs = bs[:slen]
			return
		}
	}
	halt.errorf("error decoding byte string '%s': %v", any(bs1), err)
	return
}

func (d *jsonDecDriverBytes) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
	d.advance()

	var cond bool

	if d.tok == '"' {
		d.tok = 0
		bs, cond = d.dblQuoteStringAsBytes()
		state = d.d.attachState(cond)
		return
	}

	state = dBytesDetach

	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())

	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		bs = jsonLitb[jsonLitF : jsonLitF+5]
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		bs = jsonLitb[jsonLitT : jsonLitT+4]
	default:

		bs, d.tok = d.r.jsonReadNum()
		state = d.d.attachState(!d.d.bytes)
	}
	return
}

func (d *jsonDecDriverBytes) ensureReadingString() {
	if d.tok != '"' {
		halt.errorByte("expecting string starting with '\"'; got ", d.tok)
	}
}

func (d *jsonDecDriverBytes) readUnescapedString() (bs []byte) {

	bs = d.r.jsonReadUntilDblQuote()
	d.tok = 0
	return
}

func (d *jsonDecDriverBytes) dblQuoteStringAsBytes() (buf []byte, usingBuf bool) {
	bs, c := d.r.jsonReadAsisChars()
	if c == '"' {
		return bs, !d.d.bytes
	}
	buf = append(d.buf[:0], bs...)

	checkUtf8 := d.h.ValidateUnicode
	usingBuf = true

	for {

		c = d.r.readn1()

		switch c {
		case '"', '\\', '/', '\'':
			buf = append(buf, c)
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'u':
			rr := d.appendStringAsBytesSlashU()
			if checkUtf8 && rr == unicode.ReplacementChar {
				d.buf = buf
				halt.errorBytes("invalid UTF-8 character found after: ", buf)
			}
			buf = append(buf, d.bstr[:utf8.EncodeRune(d.bstr[:], rr)]...)
		default:
			d.buf = buf
			halt.errorByte("unsupported escaped value: ", c)
		}

		bs, c = d.r.jsonReadAsisChars()
		buf = append(buf, bs...)
		if c == '"' {
			break
		}
	}
	d.buf = buf
	return
}

func (d *jsonDecDriverBytes) appendStringAsBytesSlashU() (r rune) {
	var rr uint32
	cs := d.r.readn4()
	if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
		return unicode.ReplacementChar
	}
	r = rune(rr)
	if utf16.IsSurrogate(r) {
		csu := d.r.readn2()
		cs = d.r.readn4()
		if csu[0] == '\\' && csu[1] == 'u' {
			if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
				return unicode.ReplacementChar
			}
			return utf16.DecodeRune(r, rune(rr))
		}
		return unicode.ReplacementChar
	}
	return
}

func (d *jsonDecDriverBytes) DecodeNaked() {
	z := d.d.naked()

	d.advance()
	var bs []byte
	var err error
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		z.v = valueTypeNil
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		z.v = valueTypeBool
		z.b = true
	case '{':
		z.v = valueTypeMap
	case '[':
		z.v = valueTypeArray
	case '"':

		d.tok = 0
		bs, z.b = d.dblQuoteStringAsBytes()
		att := d.d.attachState(z.b)
		if jsonNakedBoolNumInQuotedStr &&
			d.h.MapKeyAsString && len(bs) > 0 && d.d.c == containerMapKey {
			switch string(bs) {

			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default:
				if err = jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
					z.v = valueTypeString
					z.s = d.d.detach2Str(bs, att)
				}
			}
		} else {
			z.v = valueTypeString
			z.s = d.d.detach2Str(bs, att)
		}
	default:
		bs, d.tok = d.r.jsonReadNum()
		if len(bs) == 0 {
			halt.errorStr("decode number from empty string")
		}
		if err = jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
			halt.errorf("decode number from %s: %v", any(bs), err)
		}
	}
}

func (e *jsonEncDriverBytes) reset() {
	e.dl = 0

	e.typical = e.h.typical()
	if e.h.HTMLCharsAsIs {
		e.s = &jsonCharSafeBitset
	} else {
		e.s = &jsonCharHtmlSafeBitset
	}
	e.di = int8(e.h.Indent)
	e.d = e.h.Indent != 0
	e.ks = e.h.MapKeyAsString
	e.is = e.h.IntegerAsString

	var ho jsonHandleOpts
	ho.reset(e.h)
	e.timeFmt = ho.timeFmt
	e.bytesFmt = ho.bytesFmt
	e.timeFmtLayout = ""
	e.byteFmter = nil
	if len(ho.timeFmtLayouts) > 0 {
		e.timeFmtLayout = ho.timeFmtLayouts[0]
	}
	if len(ho.byteFmters) > 0 {
		e.byteFmter = ho.byteFmters[0]
	}
	e.rawext = ho.rawext
}

func (d *jsonDecDriverBytes) reset() {
	d.buf = d.d.blist.check(d.buf, 256)
	d.tok = 0

	d.jsonHandleOpts.reset(d.h)
}

func (d *jsonEncDriverBytes) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*JsonHandle)
	d.e = shared
	if shared.bytes {
		fp = jsonFpEncBytes
	} else {
		fp = jsonFpEncIO
	}

	d.init2(enc)
	return
}

func (e *jsonEncDriverBytes) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *jsonEncDriverBytes) writerEnd() { e.w.end() }

func (e *jsonEncDriverBytes) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *jsonEncDriverBytes) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *jsonDecDriverBytes) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*JsonHandle)
	d.d = shared
	if shared.bytes {
		fp = jsonFpDecBytes
	} else {
		fp = jsonFpDecIO
	}

	d.init2(dec)
	return
}

func (d *jsonDecDriverBytes) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *jsonDecDriverBytes) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *jsonDecDriverBytes) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *jsonDecDriverBytes) descBd() (s string) {
	halt.onerror(errJsonNoBd)
	return
}

func (d *jsonEncDriverBytes) init2(enc encoderI) {
	d.enc = enc

}

func (d *jsonDecDriverBytes) init2(dec decoderI) {
	d.dec = dec

	d.buf = d.buf[:0]

	d.d.jsms = d.h.MapKeyAsString
}

type helperEncDriverJsonIO struct{}
type encFnJsonIO struct {
	i  encFnInfo
	fe func(*encoderJsonIO, *encFnInfo, reflect.Value)
}
type encRtidFnJsonIO struct {
	rtid uintptr
	fn   *encFnJsonIO
}
type encoderJsonIO struct {
	dh helperEncDriverJsonIO
	fp *fastpathEsJsonIO
	e  jsonEncDriverIO
	encoderBase
}
type helperDecDriverJsonIO struct{}
type decFnJsonIO struct {
	i  decFnInfo
	fd func(*decoderJsonIO, *decFnInfo, reflect.Value)
}
type decRtidFnJsonIO struct {
	rtid uintptr
	fn   *decFnJsonIO
}
type decoderJsonIO struct {
	dh helperDecDriverJsonIO
	fp *fastpathDsJsonIO
	d  jsonDecDriverIO
	decoderBase
}
type jsonEncDriverIO struct {
	noBuiltInTypes
	h *JsonHandle
	e *encoderBase
	s *bitset256

	w bufioEncWriter

	enc encoderI

	timeFmtLayout string
	byteFmter     jsonBytesFmter

	timeFmt  jsonTimeFmt
	bytesFmt jsonBytesFmt

	di int8
	d  bool
	dl uint16

	ks bool
	is byte

	typical bool

	rawext bool

	b [48]byte
}
type jsonDecDriverIO struct {
	noBuiltInTypes
	decDriverNoopNumberHelper
	h *JsonHandle
	d *decoderBase

	r ioDecReader

	buf []byte

	tok  uint8
	_    bool
	_    byte
	bstr [4]byte

	jsonHandleOpts

	dec decoderI
}

func (e *encoderJsonIO) rawExt(_ *encFnInfo, rv reflect.Value) {
	if re := rv2i(rv).(*RawExt); re == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeRawExt(re)
	}
}

func (e *encoderJsonIO) ext(f *encFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *encoderJsonIO) selferMarshal(_ *encFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(&Encoder{e})
}

func (e *encoderJsonIO) binaryMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *encoderJsonIO) textMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *encoderJsonIO) jsonMarshal(_ *encFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *encoderJsonIO) raw(_ *encFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *encoderJsonIO) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *encoderJsonIO) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		halt.errorf("cannot encode complex number: %v, with imaginary values: %v", any(v), any(imag(v)))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *encoderJsonIO) kBool(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *encoderJsonIO) kTime(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *encoderJsonIO) kString(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *encoderJsonIO) kFloat32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *encoderJsonIO) kFloat64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *encoderJsonIO) kComplex64(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *encoderJsonIO) kComplex128(_ *encFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *encoderJsonIO) kInt(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *encoderJsonIO) kInt8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *encoderJsonIO) kInt16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *encoderJsonIO) kInt32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *encoderJsonIO) kInt64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *encoderJsonIO) kUint(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *encoderJsonIO) kUint8(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *encoderJsonIO) kUint16(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *encoderJsonIO) kUint32(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *encoderJsonIO) kUint64(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *encoderJsonIO) kUintptr(_ *encFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *encoderJsonIO) kSeqFn(rt reflect.Type) (fn *encFnJsonIO) {

	if rt = baseRT(rt); rt.Kind() != reflect.Interface {
		fn = e.fn(rt)
	}
	return
}

func (e *encoderJsonIO) kArrayWMbs(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnJsonIO
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

func (e *encoderJsonIO) kArrayW(rv reflect.Value, ti *typeInfo, isSlice bool) {
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

	var fn *encFnJsonIO
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

func (e *encoderJsonIO) kChan(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderJsonIO) kSlice(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, true)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {

		e.e.EncodeBytes(rvGetBytes(rv))
	} else {
		e.kArrayW(rv, f.ti, true)
	}
}

func (e *encoderJsonIO) kArray(f *encFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti, false)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, nil))
	} else {
		e.kArrayW(rv, f.ti, false)
	}
}

func (e *encoderJsonIO) kSliceBytesChan(rv reflect.Value) {

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

func (e *encoderJsonIO) kStructFieldKey(keyType valueType, encName string) {

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

func (e *encoderJsonIO) kStructSimple(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderJsonIO) kStruct(f *encFnInfo, rv reflect.Value) {
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

func (e *encoderJsonIO) kMap(f *encFnInfo, rv reflect.Value) {
	_ = e.e
	l := rvLenMap(rv)
	if l == 0 {
		e.e.WriteMapEmpty()
		return
	}
	e.mapStart(l)

	var keyFn, valFn *encFnJsonIO

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

func (e *encoderJsonIO) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *encFnJsonIO) {
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

func (e *encoderJsonIO) init(h Handle) {
	initHandle(h)
	callMake(&e.e)
	e.hh = h
	e.h = h.getBasicHandle()

	e.err = errEncoderNotInitialized

	e.fp = e.e.init(h, &e.encoderBase, e).(*fastpathEsJsonIO)

	if e.bytes {
		e.rtidFn = &e.h.rtidFnsEncBytes
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtBytes
	} else {
		e.rtidFn = &e.h.rtidFnsEncIO
		e.rtidFnNoExt = &e.h.rtidFnsEncNoExtIO
	}

	e.reset()
}

func (e *encoderJsonIO) reset() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}

func (e *encoderJsonIO) Encode(v interface{}) (err error) {

	defer panicValToErr(e, callRecoverSentinel, &e.err, &err, debugging)
	e.mustEncode(v)
	return
}

func (e *encoderJsonIO) MustEncode(v interface{}) {
	defer panicValToErr(e, callRecoverSentinel, &e.err, nil, true)
	e.mustEncode(v)
	return
}

func (e *encoderJsonIO) mustEncode(v interface{}) {
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

func (e *encoderJsonIO) encodeI(iv interface{}) {
	if !e.encodeBuiltin(iv) {
		e.encodeR(reflect.ValueOf(iv))
	}
}

func (e *encoderJsonIO) encodeIB(iv interface{}) {
	if !e.encodeBuiltin(iv) {

		halt.errorStr("[should not happen] invalid type passed to encodeBuiltin")
	}
}

func (e *encoderJsonIO) encodeR(base reflect.Value) {
	e.encodeValue(base, nil)
}

func (e *encoderJsonIO) encodeBuiltin(iv interface{}) (ok bool) {
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

func (e *encoderJsonIO) encodeValue(rv reflect.Value, fn *encFnJsonIO) {

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

func (e *encoderJsonIO) encodeValueNonNil(rv reflect.Value, fn *encFnJsonIO) {

	if fn.i.addrE {
		if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else {
			rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
		}
	}
	fn.fe(e, &fn.i, rv)
}

func (e *encoderJsonIO) encodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		e.encodeValue(baseRV(v), e.fn(t))
	} else {
		e.encodeValue(baseRV(v), e.fnNoExt(t))
	}
}

func (e *encoderJsonIO) marshalUtf8(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *encoderJsonIO) marshalAsis(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.writeBytesAsis(bs)
	}
}

func (e *encoderJsonIO) marshalRaw(bs []byte, fnerr error) {
	halt.onerror(fnerr)
	e.e.EncodeBytes(bs)
}

func (e *encoderJsonIO) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		halt.errorBytes("Raw values cannot be encoded: ", v)
	}
	e.e.writeBytesAsis(v)
}

func (e *encoderJsonIO) fn(t reflect.Type) *encFnJsonIO {
	return e.dh.encFnViaBH(t, e.rtidFn, e.h, e.fp, true)
}

func (e *encoderJsonIO) fnNoExt(t reflect.Type) *encFnJsonIO {
	return e.dh.encFnViaBH(t, e.rtidFnNoExt, e.h, e.fp, false)
}

func (e *encoderJsonIO) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *encoderJsonIO) mapElemValue() {
	e.e.WriteMapElemValue()
	e.c = containerMapValue
}

func (e *encoderJsonIO) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *encoderJsonIO) writerEnd() {
	e.e.writerEnd()
}

func (e *encoderJsonIO) atEndOfEncode() {
	e.e.atEndOfEncode()
}

func (e *encoderJsonIO) Reset(w io.Writer) {
	if e.bytes {
		halt.onerror(errEncNoResetBytesWithWriter)
	}
	e.reset()
	if w == nil {
		w = io.Discard
	}
	e.e.resetOutIO(w)
}

func (e *encoderJsonIO) ResetBytes(out *[]byte) {
	if !e.bytes {
		halt.onerror(errEncNoResetWriterWithBytes)
	}
	e.resetBytes(out)
}

func (e *encoderJsonIO) resetBytes(out *[]byte) {
	e.reset()
	if out == nil {
		out = &bytesEncAppenderDefOut
	}
	e.e.resetOutBytes(out)
}

func (helperEncDriverJsonIO) newEncoderBytes(out *[]byte, h Handle) *encoderJsonIO {
	var c1 encoderJsonIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(out)
	return &c1
}

func (helperEncDriverJsonIO) newEncoderIO(out io.Writer, h Handle) *encoderJsonIO {
	var c1 encoderJsonIO
	c1.bytes = false
	c1.init(h)
	c1.Reset(out)
	return &c1
}

func (helperEncDriverJsonIO) encFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathEsJsonIO) (f *fastpathEJsonIO, u reflect.Type) {
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

func (helperEncDriverJsonIO) encFindRtidFn(s []encRtidFnJsonIO, rtid uintptr) (i uint, fn *encFnJsonIO) {

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

func (helperEncDriverJsonIO) encFromRtidFnSlice(fns *atomicRtidFnSlice) (s []encRtidFnJsonIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]encRtidFnJsonIO](v))
	}
	return
}

func (dh helperEncDriverJsonIO) encFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice,
	x *BasicHandle, fp *fastpathEsJsonIO, checkExt bool) (fn *encFnJsonIO) {
	return dh.encFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperEncDriverJsonIO) encFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsJsonIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnJsonIO) {
	rtid := rt2id(rt)
	var sp []encRtidFnJsonIO = dh.encFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.encFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.encFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperEncDriverJsonIO) encFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathEsJsonIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnJsonIO) {

	fn = dh.encFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []encRtidFnJsonIO
	mu.Lock()
	sp = dh.encFromRtidFnSlice(fns)

	if sp == nil {
		sp = []encRtidFnJsonIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.encFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]encRtidFnJsonIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = encRtidFnJsonIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperEncDriverJsonIO) encFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathEsJsonIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *encFnJsonIO) {
	fn = new(encFnJsonIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	if rtid == timeTypId && timeBuiltin {
		fn.fe = (*encoderJsonIO).kTime
	} else if rtid == rawTypId {
		fn.fe = (*encoderJsonIO).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*encoderJsonIO).rawExt
		fi.addrE = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*encoderJsonIO).ext
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fe = (*encoderJsonIO).selferMarshal
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*encoderJsonIO).binaryMarshal
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fe = (*encoderJsonIO).jsonMarshal
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*encoderJsonIO).textMarshal
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
					fn.fe = func(e *encoderJsonIO, xf *encFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*encoderJsonIO).kBool
			case reflect.String:

				fn.fe = (*encoderJsonIO).kString
			case reflect.Int:
				fn.fe = (*encoderJsonIO).kInt
			case reflect.Int8:
				fn.fe = (*encoderJsonIO).kInt8
			case reflect.Int16:
				fn.fe = (*encoderJsonIO).kInt16
			case reflect.Int32:
				fn.fe = (*encoderJsonIO).kInt32
			case reflect.Int64:
				fn.fe = (*encoderJsonIO).kInt64
			case reflect.Uint:
				fn.fe = (*encoderJsonIO).kUint
			case reflect.Uint8:
				fn.fe = (*encoderJsonIO).kUint8
			case reflect.Uint16:
				fn.fe = (*encoderJsonIO).kUint16
			case reflect.Uint32:
				fn.fe = (*encoderJsonIO).kUint32
			case reflect.Uint64:
				fn.fe = (*encoderJsonIO).kUint64
			case reflect.Uintptr:
				fn.fe = (*encoderJsonIO).kUintptr
			case reflect.Float32:
				fn.fe = (*encoderJsonIO).kFloat32
			case reflect.Float64:
				fn.fe = (*encoderJsonIO).kFloat64
			case reflect.Complex64:
				fn.fe = (*encoderJsonIO).kComplex64
			case reflect.Complex128:
				fn.fe = (*encoderJsonIO).kComplex128
			case reflect.Chan:
				fn.fe = (*encoderJsonIO).kChan
			case reflect.Slice:
				fn.fe = (*encoderJsonIO).kSlice
			case reflect.Array:
				fn.fe = (*encoderJsonIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fe = (*encoderJsonIO).kStructSimple
				} else {
					fn.fe = (*encoderJsonIO).kStruct
				}
			case reflect.Map:
				fn.fe = (*encoderJsonIO).kMap
			case reflect.Interface:

				fn.fe = (*encoderJsonIO).kErr
			default:

				fn.fe = (*encoderJsonIO).kErr
			}
		}
	}
	return
}
func (d *decoderJsonIO) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoderJsonIO) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoderJsonIO) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoderJsonIO) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoderJsonIO) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoderJsonIO) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoderJsonIO) jsonUnmarshalV(tm jsonUnmarshaler) {

	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoderJsonIO) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)

}

func (d *decoderJsonIO) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoderJsonIO) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoderJsonIO) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoderJsonIO) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoderJsonIO) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoderJsonIO) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoderJsonIO) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoderJsonIO) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoderJsonIO) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoderJsonIO) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoderJsonIO) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoderJsonIO) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoderJsonIO) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoderJsonIO) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderJsonIO) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoderJsonIO) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoderJsonIO) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoderJsonIO) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoderJsonIO) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoderJsonIO) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {

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

func (d *decoderJsonIO) kInterface(f *decFnInfo, rv reflect.Value) {

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

func (d *decoderJsonIO) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoderJsonIO) kStructSimple(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderJsonIO) kStruct(f *decFnInfo, rv reflect.Value) {
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

func (d *decoderJsonIO) kSlice(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnJsonIO

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

func (d *decoderJsonIO) kArray(f *decFnInfo, rv reflect.Value) {
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
	var fn *decFnJsonIO
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

func (d *decoderJsonIO) kChan(f *decFnInfo, rv reflect.Value) {
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

	var fn *decFnJsonIO

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

func (d *decoderJsonIO) kMap(f *decFnInfo, rv reflect.Value) {
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

	var keyFn, valFn *decFnJsonIO
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

func (d *decoderJsonIO) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()

	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDsJsonIO)

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

func (d *decoderJsonIO) reset() {
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

func (d *decoderJsonIO) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

func (d *decoderJsonIO) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoderJsonIO) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

func (d *decoderJsonIO) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *decoderJsonIO) Decode(v interface{}) (err error) {

	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

func (d *decoderJsonIO) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoderJsonIO) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	d.calls++
	d.decode(v)
	d.calls--
}

func (d *decoderJsonIO) Release() {}

func (d *decoderJsonIO) swallow() {
	d.d.nextValueBytes()
}

func (d *decoderJsonIO) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoderJsonIO) decode(iv interface{}) {
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

func (d *decoderJsonIO) decodeValue(rv reflect.Value, fn *decFnJsonIO) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoderJsonIO) decodeValueNoCheckNil(rv reflect.Value, fn *decFnJsonIO) {

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

func (d *decoderJsonIO) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoderJsonIO) structFieldNotFound(index int, rvkencname string) {

	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

func (d *decoderJsonIO) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
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

func (d *decoderJsonIO) rawBytes() (v []byte) {

	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v)
		v = vv
	}
	return
}

func (d *decoderJsonIO) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

func (d *decoderJsonIO) NumBytesRead() int {
	return d.d.NumBytesRead()
}

func (d *decoderJsonIO) containerNext(j, containerLen int, hasLen bool) bool {

	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoderJsonIO) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoderJsonIO) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoderJsonIO) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderJsonIO) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoderJsonIO) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoderJsonIO) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {

	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)

}

func (d *decoderJsonIO) fn(t reflect.Type) *decFnJsonIO {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoderJsonIO) fnNoExt(t reflect.Type) *decFnJsonIO {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

func (helperDecDriverJsonIO) newDecoderBytes(in []byte, h Handle) *decoderJsonIO {
	var c1 decoderJsonIO
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in)
	return &c1
}

func (helperDecDriverJsonIO) newDecoderIO(in io.Reader, h Handle) *decoderJsonIO {
	var c1 decoderJsonIO
	c1.init(h)
	c1.Reset(in)
	return &c1
}

func (helperDecDriverJsonIO) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDsJsonIO) (f *fastpathDJsonIO, u reflect.Type) {
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

func (helperDecDriverJsonIO) decFindRtidFn(s []decRtidFnJsonIO, rtid uintptr) (i uint, fn *decFnJsonIO) {

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

func (helperDecDriverJsonIO) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFnJsonIO) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFnJsonIO](v))
	}
	return
}

func (dh helperDecDriverJsonIO) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDsJsonIO,
	checkExt bool) (fn *decFnJsonIO) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriverJsonIO) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsJsonIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnJsonIO) {
	rtid := rt2id(rt)
	var sp []decRtidFnJsonIO = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriverJsonIO) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDsJsonIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnJsonIO) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFnJsonIO
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)

	if sp == nil {
		sp = []decRtidFnJsonIO{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFnJsonIO, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFnJsonIO{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriverJsonIO) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDsJsonIO,
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFnJsonIO) {
	fn = new(decFnJsonIO)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoderJsonIO).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoderJsonIO).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoderJsonIO).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoderJsonIO).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoderJsonIO).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoderJsonIO).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {

		fn.fd = (*decoderJsonIO).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoderJsonIO).textUnmarshal
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
						fn.fd = func(d *decoderJsonIO, xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoderJsonIO, xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoderJsonIO).kBool
			case reflect.String:
				fn.fd = (*decoderJsonIO).kString
			case reflect.Int:
				fn.fd = (*decoderJsonIO).kInt
			case reflect.Int8:
				fn.fd = (*decoderJsonIO).kInt8
			case reflect.Int16:
				fn.fd = (*decoderJsonIO).kInt16
			case reflect.Int32:
				fn.fd = (*decoderJsonIO).kInt32
			case reflect.Int64:
				fn.fd = (*decoderJsonIO).kInt64
			case reflect.Uint:
				fn.fd = (*decoderJsonIO).kUint
			case reflect.Uint8:
				fn.fd = (*decoderJsonIO).kUint8
			case reflect.Uint16:
				fn.fd = (*decoderJsonIO).kUint16
			case reflect.Uint32:
				fn.fd = (*decoderJsonIO).kUint32
			case reflect.Uint64:
				fn.fd = (*decoderJsonIO).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoderJsonIO).kUintptr
			case reflect.Float32:
				fn.fd = (*decoderJsonIO).kFloat32
			case reflect.Float64:
				fn.fd = (*decoderJsonIO).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoderJsonIO).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoderJsonIO).kComplex128
			case reflect.Chan:
				fn.fd = (*decoderJsonIO).kChan
			case reflect.Slice:
				fn.fd = (*decoderJsonIO).kSlice
			case reflect.Array:
				fi.addrD = false
				fn.fd = (*decoderJsonIO).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoderJsonIO).kStructSimple
				} else {
					fn.fd = (*decoderJsonIO).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoderJsonIO).kMap
			case reflect.Interface:

				fn.fd = (*decoderJsonIO).kInterface
			default:

				fn.fd = (*decoderJsonIO).kErr
			}
		}
	}
	return
}
func (e *jsonEncDriverIO) writeIndent() {
	e.w.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.di < 0 {
		x = -x
		for x > len(jsonTabs) {
			e.w.writeb(jsonTabs[:])
			x -= len(jsonTabs)
		}
		e.w.writeb(jsonTabs[:x])
	} else {
		for x > len(jsonSpaces) {
			e.w.writeb(jsonSpaces[:])
			x -= len(jsonSpaces)
		}
		e.w.writeb(jsonSpaces[:x])
	}
}

func (e *jsonEncDriverIO) WriteArrayElem(firstTime bool) {
	if !firstTime {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriverIO) WriteMapElemKey(firstTime bool) {
	if !firstTime {
		e.w.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriverIO) WriteMapElemValue() {
	if e.d {
		e.w.writen2(':', ' ')
	} else {
		e.w.writen1(':')
	}
}

func (e *jsonEncDriverIO) EncodeNil() {

	e.w.writeb(jsonNull)
}

func (e *jsonEncDriverIO) encodeIntAsUint(v int64, quotes bool) {
	neg := v < 0
	if neg {
		v = -v
	}
	e.encodeUint(neg, quotes, uint64(v))
}

func (e *jsonEncDriverIO) EncodeTime(t time.Time) {

	if t.IsZero() {
		e.EncodeNil()
		return
	}
	switch e.timeFmt {
	case jsonTimeFmtStringLayout:
		e.b[0] = '"'
		b := t.AppendFormat(e.b[1:1], e.timeFmtLayout)
		e.b[len(b)+1] = '"'
		e.w.writeb(e.b[:len(b)+2])
	case jsonTimeFmtUnix:
		e.encodeIntAsUint(t.Unix(), false)
	case jsonTimeFmtUnixMilli:
		e.encodeIntAsUint(t.UnixMilli(), false)
	case jsonTimeFmtUnixMicro:
		e.encodeIntAsUint(t.UnixMicro(), false)
	case jsonTimeFmtUnixNano:
		e.encodeIntAsUint(t.UnixNano(), false)
	}
}

func (e *jsonEncDriverIO) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if ext == SelfExt {
		e.enc.encodeAs(rv, basetype, false)
	} else if v := ext.ConvertExt(rv); v == nil {
		e.writeNilBytes()
	} else {
		e.enc.encodeI(v)
	}
}

func (e *jsonEncDriverIO) EncodeRawExt(re *RawExt) {
	if re.Data != nil {
		e.w.writeb(re.Data)
	} else if re.Value != nil {
		e.enc.encodeI(re.Value)
	} else {
		e.EncodeNil()
	}
}

func (e *jsonEncDriverIO) EncodeBool(b bool) {
	e.w.writestr(jsonEncBoolStrs[bool2int(e.ks && e.e.c == containerMapKey)%2][bool2int(b)%2])
}

func (e *jsonEncDriverIO) encodeFloat(f float64, bitsize, fmt byte, prec int8) {
	var blen uint
	if e.ks && e.e.c == containerMapKey {
		blen = 2 + uint(len(strconv.AppendFloat(e.b[1:1], f, fmt, int(prec), int(bitsize))))

		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
	} else {
		e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), int(bitsize)))
	}
}

func (e *jsonEncDriverIO) EncodeFloat64(f float64) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.encodeFloat(f, 64, fmt, prec)
}

func (e *jsonEncDriverIO) EncodeFloat32(f float32) {
	if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.encodeFloat(float64(f), 32, fmt, prec)
}

func (e *jsonEncDriverIO) encodeUint(neg bool, quotes bool, u uint64) {
	e.w.writeb(jsonEncodeUint(neg, quotes, u, &e.b))
}

func (e *jsonEncDriverIO) EncodeInt(v int64) {
	quotes := e.is == 'A' || e.is == 'L' && (v > 1<<53 || v < -(1<<53)) ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		if quotes {
			blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
		}
		return
	}

	if v < 0 {
		e.encodeUint(true, quotes, uint64(-v))
	} else {
		e.encodeUint(false, quotes, uint64(v))
	}
}

func (e *jsonEncDriverIO) EncodeUint(v uint64) {
	quotes := e.is == 'A' || e.is == 'L' && v > 1<<53 ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {

		if quotes {
			blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.w.writeb(e.b[:blen])
		} else {
			e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
		}
		return
	}

	e.encodeUint(false, quotes, v)
}

func (e *jsonEncDriverIO) EncodeString(v string) {
	if e.h.StringToRaw {
		e.EncodeStringBytesRaw(bytesView(v))
		return
	}
	e.quoteStr(v)
}

func (e *jsonEncDriverIO) EncodeStringNoEscape4Json(v string) { e.w.writeqstr(v) }

func (e *jsonEncDriverIO) EncodeStringBytesRaw(v []byte) {
	if e.rawext {

		iv := e.h.RawBytesExt.ConvertExt(any(v))
		if iv == nil {
			e.EncodeNil()
		} else {
			e.enc.encodeI(iv)
		}
		return
	}

	if e.bytesFmt == jsonBytesFmtArray {
		e.WriteArrayStart(len(v))
		for j := range v {
			e.WriteArrayElem(j == 0)
			e.encodeUint(false, false, uint64(v[j]))
		}
		e.WriteArrayEnd()
		return
	}

	var slen int
	if e.bytesFmt == jsonBytesFmtBase64 {
		slen = base64.StdEncoding.EncodedLen(len(v))
	} else {
		slen = e.byteFmter.EncodedLen(len(v))
	}
	slen += 2

	bs := e.e.blist.peek(slen, false)[:slen]

	if e.bytesFmt == jsonBytesFmtBase64 {
		base64.StdEncoding.Encode(bs[1:], v)
	} else {
		e.byteFmter.Encode(bs[1:], v)
	}

	bs[len(bs)-1] = '"'
	bs[0] = '"'
	e.w.writeb(bs)
}

func (e *jsonEncDriverIO) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *jsonEncDriverIO) writeNilOr(v []byte) {
	if !e.h.NilCollectionToZeroLength {
		v = jsonNull
	}
	e.w.writeb(v)
}

func (e *jsonEncDriverIO) writeNilBytes() {
	e.writeNilOr(jsonArrayEmpty)
}

func (e *jsonEncDriverIO) writeNilArray() {
	e.writeNilOr(jsonArrayEmpty)
}

func (e *jsonEncDriverIO) writeNilMap() {
	e.writeNilOr(jsonMapEmpty)
}

func (e *jsonEncDriverIO) WriteArrayEmpty() {
	e.w.writen2('[', ']')
}

func (e *jsonEncDriverIO) WriteMapEmpty() {
	e.w.writen2('{', '}')
}

func (e *jsonEncDriverIO) WriteArrayStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('[')
}

func (e *jsonEncDriverIO) WriteArrayEnd() {
	if e.d {
		e.dl--

		e.writeIndent()
	}
	e.w.writen1(']')
}

func (e *jsonEncDriverIO) WriteMapStart(length int) {
	if e.d {
		e.dl++
	}
	e.w.writen1('{')
}

func (e *jsonEncDriverIO) WriteMapEnd() {
	if e.d {
		e.dl--

		e.writeIndent()
	}
	e.w.writen1('}')
}

func (e *jsonEncDriverIO) quoteStr(s string) {

	const hex = "0123456789abcdef"
	e.w.writen1('"')
	var i, start uint
	for i < uint(len(s)) {

		b := s[i]
		if e.s.isset(b) {
			i++
			continue
		}
		if b < utf8.RuneSelf {
			if start < i {
				e.w.writestr(s[start:i])
			}
			switch b {
			case '\\':
				e.w.writen2('\\', '\\')
			case '"':
				e.w.writen2('\\', '"')
			case '\n':
				e.w.writen2('\\', 'n')
			case '\t':
				e.w.writen2('\\', 't')
			case '\r':
				e.w.writen2('\\', 'r')
			case '\b':
				e.w.writen2('\\', 'b')
			case '\f':
				e.w.writen2('\\', 'f')
			default:
				e.w.writestr(`\u00`)
				e.w.writen2(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\uFFFD`)
			i++
			start = i
			continue
		}

		if jsonEscapeMultiByteUnicodeSep && (c == '\u2028' || c == '\u2029') {
			if start < i {
				e.w.writestr(s[start:i])
			}
			e.w.writestr(`\u202`)
			e.w.writen1(hex[c&0xF])
			i += uint(size)
			start = i
			continue
		}
		i += uint(size)
	}
	if start < uint(len(s)) {
		e.w.writestr(s[start:])
	}
	e.w.writen1('"')
}

func (e *jsonEncDriverIO) atEndOfEncode() {
	if e.h.TermWhitespace {
		var c byte = ' '
		if e.e.c != 0 {
			c = '\n'
		}
		e.w.writen1(c)
	}
}

func (d *jsonDecDriverIO) ReadMapStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '{' {
		halt.errorByte("read map - expect char '{' but got char: ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriverIO) ReadArrayStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return containerLenNil
	}
	if d.tok != '[' {
		halt.errorByte("read array - expect char '[' but got char ", d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriverIO) CheckBreak() bool {
	d.advance()
	return d.tok == '}' || d.tok == ']'
}

func (d *jsonDecDriverIO) checkSep(xc byte) {
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriverIO) ReadArrayElem(firstTime bool) {
	if !firstTime {
		d.checkSep(',')
	}
}

func (d *jsonDecDriverIO) ReadArrayEnd() {
	d.checkSep(']')
}

func (d *jsonDecDriverIO) ReadMapElemKey(firstTime bool) {
	d.ReadArrayElem(firstTime)
}

func (d *jsonDecDriverIO) ReadMapElemValue() {
	d.checkSep(':')
}

func (d *jsonDecDriverIO) ReadMapEnd() {
	d.checkSep('}')
}

func (d *jsonDecDriverIO) readDelimError(xc uint8) {
	halt.errorf("read json delimiter - expect char '%c' but got char '%c'", xc, d.tok)
}

func (d *jsonDecDriverIO) checkLit3(got, expect [3]byte) {
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr3(got, expect)
	}
	d.tok = 0
}

func (d *jsonDecDriverIO) checkLit4(got, expect [4]byte) {
	if jsonValidateSymbols && got != expect {
		jsonCheckLitErr4(got, expect)
	}
	d.tok = 0
}

func (d *jsonDecDriverIO) skipWhitespace() {
	d.tok = d.r.skipWhitespace()
}

func (d *jsonDecDriverIO) advance() {

	if d.tok < 33 {
		d.skipWhitespace()
	}
}

func (d *jsonDecDriverIO) nextValueBytes() []byte {
	consumeString := func() {
	TOP:
		_, c := d.r.jsonReadAsisChars()
		if c == '\\' {
			d.r.readn1()
			goto TOP
		}
	}

	d.advance()
	d.r.startRecording()

	switch d.tok {
	default:
		_, d.tok = d.r.jsonReadNum()

		if d.tok != 0 {
			vv := d.r.stopRecording()
			return vv[:len(vv)-1]
		}
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
	case '"':
		consumeString()
		d.tok = 0
	case '{', '[':
		var elem struct{}
		var stack []struct{}

		stack = append(stack, elem)

		for len(stack) != 0 {
			c := d.r.readn1()
			switch c {
			case '"':
				consumeString()
			case '{', '[':
				stack = append(stack, elem)
			case '}', ']':
				stack = stack[:len(stack)-1]
			}
		}
		d.tok = 0
	}
	return d.r.stopRecording()
}

func (d *jsonDecDriverIO) TryNil() bool {
	d.advance()

	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriverIO) DecodeBool() (v bool) {
	d.advance()

	fquot := d.d.c == containerMapKey && d.tok == '"'
	if fquot {
		d.tok = d.r.readn1()
	}
	switch d.tok {
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())

	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		v = true
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())

	default:
		halt.errorByte("decode bool: got first char: ", d.tok)

	}
	if fquot {
		d.r.readn1()
	}
	return
}

func (d *jsonDecDriverIO) DecodeTime() (t time.Time) {

	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	var bs []byte

	if d.tok != '"' {
		bs, d.tok = d.r.jsonReadNum()
		i := d.parseInt64(bs)
		switch d.timeFmtNum {
		case jsonTimeFmtUnix:
			t = time.Unix(i, 0)
		case jsonTimeFmtUnixMilli:
			t = time.UnixMilli(i)
		case jsonTimeFmtUnixMicro:
			t = time.UnixMicro(i)
		case jsonTimeFmtUnixNano:
			t = time.Unix(0, i)
		default:
			halt.errorStr("invalid timeFmtNum")
		}
		return
	}

	bs = d.readUnescapedString()
	var err error
	for _, v := range d.timeFmtLayouts {
		t, err = time.Parse(v, stringView(bs))
		if err == nil {
			return
		}
	}
	halt.errorStr("error decoding time")
	return
}

func (d *jsonDecDriverIO) ContainerType() (vt valueType) {

	d.advance()

	if d.tok == '{' {
		return valueTypeMap
	} else if d.tok == '[' {
		return valueTypeArray
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return valueTypeNil
	} else if d.tok == '"' {
		return valueTypeString
	}
	return valueTypeUnset
}

func (d *jsonDecDriverIO) decNumBytes() (bs []byte) {
	d.advance()
	if d.tok == '"' {
		bs = d.r.jsonReadUntilDblQuote()
		d.tok = 0
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
	} else {
		bs, d.tok = d.r.jsonReadNum()
	}
	return
}

func (d *jsonDecDriverIO) DecodeUint64() (u uint64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if neg {
		halt.errorf("negative number cannot be decoded as uint64: %s", any(b))
	}
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseUint"))
	}
	return
}

func (d *jsonDecDriverIO) DecodeInt64() (v int64) {
	return d.parseInt64(d.decNumBytes())
}

func (d *jsonDecDriverIO) parseInt64(b []byte) (v int64) {
	u, neg, ok := parseInteger_bytes(b)
	if !ok {
		halt.onerror(strconvParseErr(b, "ParseInt"))
	}
	if chkOvf.Uint2Int(u, neg) {
		halt.errorBytes("overflow decoding number from ", b)
	}
	if neg {
		v = -int64(u)
	} else {
		v = int64(u)
	}
	return
}

func (d *jsonDecDriverIO) DecodeFloat64() (f float64) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat64(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriverIO) DecodeFloat32() (f float32) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat32(bs)
	halt.onerror(err)
	return
}

func (d *jsonDecDriverIO) advanceNil() (ok bool) {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriverIO) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if d.advanceNil() {
		return
	}
	if ext == SelfExt {
		d.dec.decodeAs(rv, basetype, false)
	} else {
		d.dec.interfaceExtConvertAndDecode(rv, ext)
	}
}

func (d *jsonDecDriverIO) DecodeRawExt(re *RawExt) {
	if d.advanceNil() {
		return
	}
	d.dec.decode(&re.Value)
}

func (d *jsonDecDriverIO) decBytesFromArray(bs []byte) []byte {
	d.advance()
	if d.tok != ']' {
		bs = append(bs, uint8(d.DecodeUint64()))
		d.advance()
	}
	for d.tok != ']' {
		if d.tok != ',' {
			halt.errorByte("read array element - expect char ',' but got char: ", d.tok)
		}
		d.tok = 0
		bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		d.advance()
	}
	d.tok = 0
	return bs
}

func (d *jsonDecDriverIO) DecodeBytes() (bs []byte, state dBytesAttachState) {
	d.advance()
	state = dBytesDetach
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		return
	}
	state = dBytesAttachBuffer

	if d.rawext {
		d.buf = d.buf[:0]
		d.dec.interfaceExtConvertAndDecode(&d.buf, d.h.RawBytesExt)
		bs = d.buf
		return
	}

	if d.tok == '[' {
		d.tok = 0

		bs = d.decBytesFromArray(d.buf[:0])
		d.buf = bs
		return
	}

	d.ensureReadingString()
	bs1 := d.readUnescapedString()

	slen := base64.StdEncoding.DecodedLen(len(bs1))
	if slen == 0 {
		bs = zeroByteSlice
		state = dBytesDetach
	} else if slen <= cap(d.buf) {
		bs = d.buf[:slen]
	} else {
		d.buf = d.d.blist.putGet(d.buf, slen)[:slen]
		bs = d.buf
	}
	var err error
	for _, v := range d.byteFmters {

		slen, err = v.Decode(bs, bs1)
		if err == nil {
			bs = bs[:slen]
			return
		}
	}
	halt.errorf("error decoding byte string '%s': %v", any(bs1), err)
	return
}

func (d *jsonDecDriverIO) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
	d.advance()

	var cond bool

	if d.tok == '"' {
		d.tok = 0
		bs, cond = d.dblQuoteStringAsBytes()
		state = d.d.attachState(cond)
		return
	}

	state = dBytesDetach

	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())

	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		bs = jsonLitb[jsonLitF : jsonLitF+5]
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		bs = jsonLitb[jsonLitT : jsonLitT+4]
	default:

		bs, d.tok = d.r.jsonReadNum()
		state = d.d.attachState(!d.d.bytes)
	}
	return
}

func (d *jsonDecDriverIO) ensureReadingString() {
	if d.tok != '"' {
		halt.errorByte("expecting string starting with '\"'; got ", d.tok)
	}
}

func (d *jsonDecDriverIO) readUnescapedString() (bs []byte) {

	bs = d.r.jsonReadUntilDblQuote()
	d.tok = 0
	return
}

func (d *jsonDecDriverIO) dblQuoteStringAsBytes() (buf []byte, usingBuf bool) {
	bs, c := d.r.jsonReadAsisChars()
	if c == '"' {
		return bs, !d.d.bytes
	}
	buf = append(d.buf[:0], bs...)

	checkUtf8 := d.h.ValidateUnicode
	usingBuf = true

	for {

		c = d.r.readn1()

		switch c {
		case '"', '\\', '/', '\'':
			buf = append(buf, c)
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'u':
			rr := d.appendStringAsBytesSlashU()
			if checkUtf8 && rr == unicode.ReplacementChar {
				d.buf = buf
				halt.errorBytes("invalid UTF-8 character found after: ", buf)
			}
			buf = append(buf, d.bstr[:utf8.EncodeRune(d.bstr[:], rr)]...)
		default:
			d.buf = buf
			halt.errorByte("unsupported escaped value: ", c)
		}

		bs, c = d.r.jsonReadAsisChars()
		buf = append(buf, bs...)
		if c == '"' {
			break
		}
	}
	d.buf = buf
	return
}

func (d *jsonDecDriverIO) appendStringAsBytesSlashU() (r rune) {
	var rr uint32
	cs := d.r.readn4()
	if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
		return unicode.ReplacementChar
	}
	r = rune(rr)
	if utf16.IsSurrogate(r) {
		csu := d.r.readn2()
		cs = d.r.readn4()
		if csu[0] == '\\' && csu[1] == 'u' {
			if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
				return unicode.ReplacementChar
			}
			return utf16.DecodeRune(r, rune(rr))
		}
		return unicode.ReplacementChar
	}
	return
}

func (d *jsonDecDriverIO) DecodeNaked() {
	z := d.d.naked()

	d.advance()
	var bs []byte
	var err error
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.r.readn3())
		z.v = valueTypeNil
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.r.readn4())
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.r.readn3())
		z.v = valueTypeBool
		z.b = true
	case '{':
		z.v = valueTypeMap
	case '[':
		z.v = valueTypeArray
	case '"':

		d.tok = 0
		bs, z.b = d.dblQuoteStringAsBytes()
		att := d.d.attachState(z.b)
		if jsonNakedBoolNumInQuotedStr &&
			d.h.MapKeyAsString && len(bs) > 0 && d.d.c == containerMapKey {
			switch string(bs) {

			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default:
				if err = jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
					z.v = valueTypeString
					z.s = d.d.detach2Str(bs, att)
				}
			}
		} else {
			z.v = valueTypeString
			z.s = d.d.detach2Str(bs, att)
		}
	default:
		bs, d.tok = d.r.jsonReadNum()
		if len(bs) == 0 {
			halt.errorStr("decode number from empty string")
		}
		if err = jsonNakedNum(z, bs, d.h.PreferFloat, d.h.SignedInteger); err != nil {
			halt.errorf("decode number from %s: %v", any(bs), err)
		}
	}
}

func (e *jsonEncDriverIO) reset() {
	e.dl = 0

	e.typical = e.h.typical()
	if e.h.HTMLCharsAsIs {
		e.s = &jsonCharSafeBitset
	} else {
		e.s = &jsonCharHtmlSafeBitset
	}
	e.di = int8(e.h.Indent)
	e.d = e.h.Indent != 0
	e.ks = e.h.MapKeyAsString
	e.is = e.h.IntegerAsString

	var ho jsonHandleOpts
	ho.reset(e.h)
	e.timeFmt = ho.timeFmt
	e.bytesFmt = ho.bytesFmt
	e.timeFmtLayout = ""
	e.byteFmter = nil
	if len(ho.timeFmtLayouts) > 0 {
		e.timeFmtLayout = ho.timeFmtLayouts[0]
	}
	if len(ho.byteFmters) > 0 {
		e.byteFmter = ho.byteFmters[0]
	}
	e.rawext = ho.rawext
}

func (d *jsonDecDriverIO) reset() {
	d.buf = d.d.blist.check(d.buf, 256)
	d.tok = 0

	d.jsonHandleOpts.reset(d.h)
}

func (d *jsonEncDriverIO) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*JsonHandle)
	d.e = shared
	if shared.bytes {
		fp = jsonFpEncBytes
	} else {
		fp = jsonFpEncIO
	}

	d.init2(enc)
	return
}

func (e *jsonEncDriverIO) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *jsonEncDriverIO) writerEnd() { e.w.end() }

func (e *jsonEncDriverIO) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *jsonEncDriverIO) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

func (d *jsonDecDriverIO) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*JsonHandle)
	d.d = shared
	if shared.bytes {
		fp = jsonFpDecBytes
	} else {
		fp = jsonFpDecIO
	}

	d.init2(dec)
	return
}

func (d *jsonDecDriverIO) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *jsonDecDriverIO) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *jsonDecDriverIO) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

func (d *jsonDecDriverIO) descBd() (s string) {
	halt.onerror(errJsonNoBd)
	return
}

func (d *jsonEncDriverIO) init2(enc encoderI) {
	d.enc = enc

}

func (d *jsonDecDriverIO) init2(dec decoderI) {
	d.dec = dec

	d.buf = d.buf[:0]

	d.d.jsms = d.h.MapKeyAsString
}
