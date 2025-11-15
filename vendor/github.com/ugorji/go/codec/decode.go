//go:build notmono || codec.notmono

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"io"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type helperDecDriver[T decDriver] struct{}

// decFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the
// code block that should be called (encapsulated in a function)
// instead of executing the checks every time.
type decFn[T decDriver] struct {
	i  decFnInfo
	fd func(*decoder[T], *decFnInfo, reflect.Value)
	// _  [1]uint64 // padding (cache-aligned)
}

type decRtidFn[T decDriver] struct {
	rtid uintptr
	fn   *decFn[T]
}

// ----

// Decoder reads and decodes an object from an input stream in a supported format.
//
// Decoder is NOT safe for concurrent use i.e. a Decoder cannot be used
// concurrently in multiple goroutines.
//
// However, as Decoder could be allocation heavy to initialize, a Reset method is provided
// so its state can be reused to decode new input streams repeatedly.
// This is the idiomatic way to use.
type decoder[T decDriver] struct {
	dh helperDecDriver[T]
	fp *fastpathDs[T]
	d  T
	decoderBase
}

func (d *decoder[T]) rawExt(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeRawExt(rv2i(rv).(*RawExt))
}

func (d *decoder[T]) ext(f *decFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *decoder[T]) selferUnmarshal(_ *decFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(&Decoder{d})
}

func (d *decoder[T]) binaryUnmarshal(_ *decFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs, _ := d.d.DecodeBytes()
	fnerr := bm.UnmarshalBinary(xbs)
	halt.onerror(fnerr)
}

func (d *decoder[T]) textUnmarshal(_ *decFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(bytesOKs(d.d.DecodeStringAsBytes()))
	halt.onerror(fnerr)
}

func (d *decoder[T]) jsonUnmarshal(_ *decFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *decoder[T]) jsonUnmarshalV(tm jsonUnmarshaler) {
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	halt.onerror(tm.UnmarshalJSON(d.d.nextValueBytes()))
}

func (d *decoder[T]) kErr(_ *decFnInfo, rv reflect.Value) {
	halt.errorf("unsupported decoding kind: %s, for %#v", rv.Kind(), rv)
	// halt.errorStr2("no decoding function defined for kind: ", rv.Kind().String())
}

func (d *decoder[T]) raw(_ *decFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *decoder[T]) kString(_ *decFnInfo, rv reflect.Value) {
	rvSetString(rv, d.detach2Str(d.d.DecodeStringAsBytes()))
}

func (d *decoder[T]) kBool(_ *decFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *decoder[T]) kTime(_ *decFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *decoder[T]) kFloat32(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.d.DecodeFloat32())
}

func (d *decoder[T]) kFloat64(_ *decFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *decoder[T]) kComplex64(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.d.DecodeFloat32(), 0))
}

func (d *decoder[T]) kComplex128(_ *decFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *decoder[T]) kInt(_ *decFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *decoder[T]) kInt8(_ *decFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *decoder[T]) kInt16(_ *decFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *decoder[T]) kInt32(_ *decFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *decoder[T]) kInt64(_ *decFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *decoder[T]) kUint(_ *decFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoder[T]) kUintptr(_ *decFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *decoder[T]) kUint8(_ *decFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *decoder[T]) kUint16(_ *decFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *decoder[T]) kUint32(_ *decFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *decoder[T]) kUint64(_ *decFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *decoder[T]) kInterfaceNaked(f *decFnInfo) (rvn reflect.Value) {
	// nil interface:
	// use some hieristics to decode it appropriately
	// based on the detected next value in the stream.
	n := d.naked()
	d.d.DecodeNaked()

	// We cannot decode non-nil stream value into nil interface with methods (e.g. io.Reader).
	// However, it is possible that the user has ways to pass in a type for a given interface
	//   - MapType
	//   - SliceType
	//   - Extensions
	//
	// Consequently, we should relax this. Put it behind a const flag for now.
	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		halt.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
	}

	// We generally make a pointer to the container here, and pass along,
	// so that they will be initialized later when we know the length of the collection.

	switch n.v {
	case valueTypeMap:
		mtid := d.mtid
		if mtid == 0 {
			if d.jsms { // if json, default to a map type with string keys
				mtid = mapStrIntfTypId // for json performance
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
			// // made map is fully initialized for direct modification.
			// // There's no need to make a pointer to it first.
			// rvn = makeMapReflect(d.h.MapType, 0)
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
		tag, bytes := n.u, n.l // calling decode below might taint the values
		bfn := d.h.getExtForTag(tag)
		var re = RawExt{Tag: tag}
		if bytes == nil {
			// one of the InterfaceExt ones: json and cbor.
			// (likely cbor, as json has no tagging support and won't reveal valueTypeExt)
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
			// one of the BytesExt ones: binc, msgpack, simple
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
		// if struct/array, directly store pointer into the interface
		if d.h.PreferPointerForStructOrArray && rvn.CanAddr() {
			if rk := rvn.Kind(); rk == reflect.Array || rk == reflect.Struct {
				rvn = rvn.Addr()
			}
		}
	case valueTypeNil:
		// rvn = reflect.Zero(f.ti.rt)
		// no-op
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

func (d *decoder[T]) kInterface(f *decFnInfo, rv reflect.Value) {
	// Note: A consequence of how kInterface works, is that
	// if an interface already contains something, we try
	// to decode into what was there before.
	// We do not replace with a generic value (as got from decodeNaked).
	//
	// every interface passed here MUST be settable.
	//
	// ensure you call rvSetIntf(...) before returning.

	isnilrv := rvIsNil(rv)

	var rvn reflect.Value

	if d.h.InterfaceReset {
		// check if mapping to a type: if so, initialize it and move on
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
		// check if mapping to a type: if so, initialize it and move on
		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			}
			return
		}
	} else {
		// now we have a non-nil interface value, meaning it contains a type
		rvn = rv.Elem()
	}

	// rvn is now a non-interface type

	canDecode, _ := isDecodeable(rvn)

	// Note: interface{} is settable, but underlying type may not be.
	// Consequently, we MAY have to allocate a value (containing the underlying value),
	// decode into it, and reset the interface to that new value.

	if !canDecode {
		rvn2 := d.oneShotAddrRV(rvn.Type(), rvn.Kind())
		rvSetDirect(rvn2, rvn)
		rvn = rvn2
	}

	d.decodeValue(rvn, nil)
	rvSetIntf(rv, rvn)
}

func (d *decoder[T]) kStructField(si *structFieldInfo, rv reflect.Value) {
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

func (d *decoder[T]) kStructSimple(f *decFnInfo, rv reflect.Value) {
	_ = d.d // early asserts d, d.d are not nil once
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
		// Not much gain from doing it two ways for array (used less frequently than structs).
		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		// iterate all the items in the stream.
		//   - if mapped elem-wise to a field, handle it
		//   - if more stream items than can be mapped, error it
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

func (d *decoder[T]) kStruct(f *decFnInfo, rv reflect.Value) {
	_ = d.d // early asserts d, d.d are not nil once
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
			// use if-else since <8 branches and we need good branch prediction for string
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
				// store rvkencname in new []byte, as it previously shares Decoder.b, which is used in decode
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
		// Not much gain from doing it two ways for array.
		// Arrays are not used as much for structs.
		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		// iterate all the items in the stream
		// if mapped elem-wise to a field, handle it
		// if more stream items than can be mapped, error it
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

func (d *decoder[T]) kSlice(f *decFnInfo, rv reflect.Value) {
	_ = d.d // early asserts d, d.d are not nil once
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).

	// Note: rv is a slice type here - guaranteed

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {
		// you can only decode bytes or string in the stream into a slice or array of bytes
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
			// not addressable byte slice, so do not decode into it past the length
			d.decodeBytesInto(rvbs[:len(rvbs):len(rvbs)], true)
		}
		return
	}

	// only expects valueType(Array|Map) - never Nil
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}

	// an array can never return a nil slice. so no need to check f.array here.
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

	var fn *decFn[T]

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
			} else if rvCanset { // rvlen1 > rvcap
				rvlen = rvlen1
				rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = !rvCanset
			} else { // rvlen1 > rvcap && !canSet
				halt.errorStr("cannot decode into non-settable slice")
			}
			if rvChanged && oldRvlenGtZero && rtelem0Mut {
				rvCopySlice(rv, rv0, rtelem) // only copy up to length NOT cap i.e. rv0.Slice(0, rvcap)
			}
		} else if containerLenS != rvlen {
			if rvCanset {
				rvlen = containerLenS
				rvSetSliceLen(rv, rvlen)
			}
		}
	}

	// consider creating new element once, and just decoding into it.
	var elemReset = d.h.SliceElementReset

	// when decoding into slices, there may be more values in the stream than the slice length.
	// decodeValue handles this better when coming from an addressable value (known to reflect.Value).
	// Consequently, builtin handling skips slices.

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
			if rvIsNil(rv) { // means hasLen = false
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

		// if indefinite, etc, then expand the slice if necessary
		if j >= rvlen {

			// expand the slice up to the cap.
			// Note that we did, so we have to reset it later.

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
				// note: 1 requested is hint/minimum - new capacity with more space
				rvlen = rvcap
				rvChanged = !rvCanset
			}
		}

		// we check if we can make this an addr, and do builtin
		// e.g. if []ints, then fastpath should handle it?
		// but if not, we should treat it as each element is *int, and decode into it.

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
				d.decode(rv2i(rvAddr(rv9, ti.tielem.ptr))) // d.decode(rv2i(rv9.Addr()))
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
		// rvlen = j
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

	if rvChanged { // infers rvCanset=true, so it can be reset
		rvSetDirect(rv0, rv)
	}
}

func (d *decoder[T]) kArray(f *decFnInfo, rv reflect.Value) {
	_ = d.d // early asserts d, d.d are not nil once
	// An array can be set from a map or array in stream.
	ti := f.ti
	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {
		// you can only decode bytes or string in the stream into a slice or array of bytes
		if ti.elemkind != uint8(reflect.Uint8) {
			halt.errorf("bytes/string in stream can decode into array of bytes, but not %v", ti.rt)
		}
		rvbs := rvGetArrayBytes(rv, nil)
		d.decodeBytesInto(rvbs, true)
		return
	}

	// only expects valueType(Array|Map) - never Nil
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}

	// an array can never return a nil slice. so no need to check f.array here.
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

	rvlen := rv.Len() // same as cap
	hasLen := containerLenS >= 0
	if hasLen && containerLenS > rvlen {
		halt.errorf("cannot decode into array with length: %v, less than container length: %v", any(rvlen), any(containerLenS))
	}

	// consider creating new element once, and just decoding into it.
	var elemReset = d.h.SliceElementReset

	var rtelemIsPtr bool
	var rtelemElem reflect.Type
	var fn *decFn[T]
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
		// note that you cannot expand the array if indefinite and we go past array length
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
				d.decode(rv2i(rvAddr(rv9, ti.tielem.ptr))) // d.decode(rv2i(rv9.Addr()))
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

func (d *decoder[T]) kChan(f *decFnInfo, rv reflect.Value) {
	_ = d.d // early asserts d, d.d are not nil once
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		halt.errorStr("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {
		// you can only decode bytes or string in the stream into a slice or array of bytes
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

	// only expects valueType(Array|Map) - never Nil
	var containerLenS int
	isArray := ctyp == valueTypeArray
	if isArray {
		containerLenS = d.arrayStart(d.d.ReadArrayStart())
	} else if ctyp == valueTypeMap {
		containerLenS = d.mapStart(d.d.ReadMapStart()) * 2
	} else {
		halt.errorStr2("decoding into a slice, expect map/array - got ", ctyp.String())
	}

	// an array can never return a nil slice. so no need to check f.array here.
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

	var fn *decFn[T]

	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int // = rv.Len()
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

	if rvChanged { // infers rvCanset=true, so it can be reset
		rvSetDirect(rv0, rv)
	}

}

func (d *decoder[T]) kMap(f *decFnInfo, rv reflect.Value) {
	_ = d.d // early asserts d, d.d are not nil once
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
	// kfast := mapKeyFastKindFor(ktypeKind)
	// visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	// visref := refBitset.isset(ti.elemkind)

	vtypePtr := vtypeKind == reflect.Ptr
	ktypePtr := ktypeKind == reflect.Ptr

	vTransient := decUseTransient && !vtypePtr && ti.tielem.flagCanTransient
	// keys are transient iff values are transient first
	kTransient := vTransient && !ktypePtr && ti.tikey.flagCanTransient

	var vtypeElem reflect.Type

	var keyFn, valFn *decFn[T]
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

	rvkMut := !scalarBitset.isset(ti.keykind) // if ktype is immutable, then re-use the same rvk.
	rvvMut := !scalarBitset.isset(ti.elemkind)
	rvvCanNil := isnilBitset.isset(ti.elemkind)

	// rvk: key
	// rvkn: if non-mutable, on each iteration of loop, set rvk to this
	// rvv: value
	// rvvn: if non-mutable, on each iteration of loop, set rvv to this
	//       if mutable, may be used as a temporary value for local-scoped operations
	// rvva: if mutable, used as transient value for use for key lookup
	// rvvz: zero value of map value type, used to do a map set when nil is found in stream
	var rvk, rvkn, rvv, rvvn, rvva, rvvz reflect.Value

	// we do a doMapGet if kind is mutable, and InterfaceReset=true if interface
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

	// Use a possibly transient (map) value (and key), to reduce allocation

	// when decoding into slices, there may be more values in the stream than the slice length.
	// decodeValue handles this better when coming from an addressable value (known to reflect.Value).
	// Consequently, builtin handling skips slices.

	var vElem, kElem reflect.Type
	kbuiltin := ti.tikey.flagDecBuiltin && ti.keykind != uint8(reflect.Slice)
	vbuiltin := ti.tielem.flagDecBuiltin // && ti.elemkind != uint8(reflect.Slice)
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
			// if vtypekind is a scalar and thus value will be decoded using TransientAddrK,
			// then it is ok to use TransientAddr2K for the map key.
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
			// special case if interface wrapping a byte slice
			if ktypeIsIntf {
				if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp {
					kstr2bs = rvGetBytes(rvk2)
					kstr, mapKeyStringSharesBytesBuf = d.bytes2Str(kstr2bs, dBytesAttachView)
					rvSetIntf(rvk, rv4istr(kstr))
				}
				// NOTE: consider failing early if map/slice/func
			}
		}

		// TryNil will try to read from the stream and check if a nil marker.
		//
		// When using ioDecReader (specifically in bufio mode), this TryNil call could
		// override part of the buffer used for the string key.
		//
		// To mitigate this, we do a special check for ioDecReader in bufio mode.
		if mapKeyStringSharesBytesBuf && d.bufio {
			if ktypeIsString {
				rvSetString(rvk, d.detach2Str(kstr2bs, att))
			} else { // ktypeIsIntf
				rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
			}
			mapKeyStringSharesBytesBuf = false
		}

		d.mapElemValue()

		if d.d.TryNil() {
			if mapKeyStringSharesBytesBuf {
				if ktypeIsString {
					rvSetString(rvk, d.detach2Str(kstr2bs, att))
				} else { // ktypeIsIntf
					rvSetIntf(rvk, rv4istr(d.detach2Str(kstr2bs, att)))
				}
			}
			// since a map, we have to set zero value if needed
			if !rvvz.IsValid() {
				rvvz = rvZeroK(vtype, vtypeKind)
			}
			mapSet(rv, rvk, rvvz, mparams)
			continue
		}

		// there is non-nil content in the stream to decode ...
		// consequently, it's ok to just directly create new value to the pointer (if vtypePtr)

		// set doMapSet to false iff u do a get, and the return value is a non-nil pointer
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
			case reflect.Ptr, reflect.Map: // ok to decode directly into map
				doMapSet = false
			case reflect.Interface:
				// if an interface{}, just decode into it iff a non-nil ptr/map, else allocate afresh
				rvvn = rvv.Elem()
				if k := rvvn.Kind(); (k == reflect.Ptr || k == reflect.Map) && !rvIsNil(rvvn) {
					d.decodeValueNoCheckNil(rvvn, nil) // valFn is incorrect here
					continue
				}
				// make addressable (so we can set the interface)
				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvSetIntf(rvvn, rvv)
				rvv = rvvn
			default:
				// make addressable (so you can set the slice/array elements, etc)
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
			rvv = reflect.New(vtypeElem) // non-nil in stream, so allocate value
		} else if vTransient {
			rvv = d.perType.TransientAddrK(vtype, vtypeKind)
		} else {
			rvv = rvZeroAddrK(vtype, vtypeKind)
		}

	DECODE_VALUE_NO_CHECK_NIL:
		if doMapSet && mapKeyStringSharesBytesBuf {
			if ktypeIsString {
				rvSetString(rvk, d.detach2Str(kstr2bs, att))
			} else { // ktypeIsIntf
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

func (d *decoder[T]) init(h Handle) {
	initHandle(h)
	callMake(&d.d)
	d.hh = h
	d.h = h.getBasicHandle()
	// d.zeroCopy = d.h.ZeroCopy
	// d.be = h.isBinary()
	d.err = errDecoderNotInitialized

	if d.h.InternString && d.is == nil {
		d.is.init()
	}

	// d.fp = fastpathDList[T]()
	d.fp = d.d.init(h, &d.decoderBase, d).(*fastpathDs[T]) // should set js, cbor, bytes, etc

	// d.cbreak = d.js || d.cbor

	if d.bytes {
		d.rtidFn = &d.h.rtidFnsDecBytes
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtBytes
	} else {
		d.bufio = d.h.ReaderBufferSize > 0
		d.rtidFn = &d.h.rtidFnsDecIO
		d.rtidFnNoExt = &d.h.rtidFnsDecNoExtIO
	}

	d.reset()
	// NOTE: do not initialize d.n here. It is lazily initialized in d.naked()
}

func (d *decoder[T]) reset() {
	d.d.reset()
	d.err = nil
	d.c = 0
	d.depth = 0
	d.calls = 0
	// reset all things which were cached from the Handle, but could change
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

// Reset the Decoder with a new Reader to decode from,
// clearing all state from last run(s).
func (d *decoder[T]) Reset(r io.Reader) {
	if d.bytes {
		halt.onerror(errDecNoResetBytesWithReader)
	}
	d.reset()
	if r == nil {
		r = &eofReader
	}
	d.d.resetInIO(r)
}

// ResetBytes resets the Decoder with a new []byte to decode from,
// clearing all state from last run(s).
func (d *decoder[T]) ResetBytes(in []byte) {
	if !d.bytes {
		halt.onerror(errDecNoResetReaderWithBytes)
	}
	d.resetBytes(in)
}

func (d *decoder[T]) resetBytes(in []byte) {
	d.reset()
	if in == nil {
		in = zeroByteSlice
	}
	d.d.resetInBytes(in)
}

// ResetString resets the Decoder with a new string to decode from,
// clearing all state from last run(s).
//
// It is a convenience function that calls ResetBytes with a
// []byte view into the string.
//
// This can be an efficient zero-copy if using default mode i.e. without codec.safe tag.
func (d *decoder[T]) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
//
// Sample usages:
//
//	// Decoding into a non-nil typed value
//	var f float32
//	err = codec.NewDecoder(r, handle).Decode(&f)
//
//	// Decoding into nil interface
//	var v interface{}
//	dec := codec.NewDecoder(r, handle)
//	err = dec.Decode(&v)
//
// When decoding into a nil interface{}, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64.
//   - Other values are decoded appropriately depending on the type:
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
//
// Configurations exist on the Handle to override defaults
// (e.g. for MapType, SliceType and how to decode raw bytes).
//
// When decoding into a non-nil interface{} value, the mode of encoding is based on the
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
//
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container i.e. the values
// in these containers will not be zero'ed before decoding.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// This in-place update maintains consistency in the decoding philosophy (i.e. we ALWAYS update
// in place by default). However, the consequence of this is that values in slices or maps
// which are not zero'ed before hand, will have part of the prior values in place after decode
// if the stream doesn't contain an update for those parts.
//
// This in-place update can be disabled by configuring the MapValueReset and SliceElementReset
// decode options available on every handle.
//
// Furthermore, when decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
// Note: we allow nil values in the stream anywhere except for map keys.
// A nil value in the encoded stream where a map key is expected is treated as an error.
//
// Note that an error from a Decode call will make the Decoder unusable moving forward.
// This is because the state of the Decoder, it's input stream, etc are no longer stable.
// Any subsequent calls to Decode will trigger the same error.
func (d *decoder[T]) Decode(v interface{}) (err error) {
	// tried to use closure, as runtime optimizes defer with no params.
	// This seemed to be causing weird issues (like circular reference found, unexpected panic, etc).
	// Also, see https://github.com/golang/go/issues/14939#issuecomment-417836139
	defer panicValToErr(d, callRecoverSentinel, &d.err, &err, debugging)
	d.mustDecode(v)
	return
}

// MustDecode is like Decode, but panics if unable to Decode.
//
// Note: This provides insight to the code location that triggered the error.
//
// Note that an error from a Decode call will make the Decoder unusable moving forward.
// This is because the state of the Decoder, it's input stream, etc are no longer stable.
// Any subsequent calls to Decode will trigger the same error.
func (d *decoder[T]) MustDecode(v interface{}) {
	defer panicValToErr(d, callRecoverSentinel, &d.err, nil, true)
	d.mustDecode(v)
	return
}

func (d *decoder[T]) mustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	// Top-level: v is a pointer and not nil.
	d.calls++
	d.decode(v)
	d.calls--
}

// Release is a no-op.
//
// Deprecated: Pooled resources are not used with a Decoder.
// This method is kept for compatibility reasons only.
func (d *decoder[T]) Release() {}

func (d *decoder[T]) swallow() {
	d.d.nextValueBytes()
}

func (d *decoder[T]) nextValueBytes() []byte {
	return d.d.nextValueBytes()
}

func (d *decoder[T]) decode(iv interface{}) {
	_ = d.d // early asserts d, d.d are not nil once
	// a switch with only concrete types can be optimized.
	// consequently, we deal with nil and interfaces outside the switch.

	rv, ok := isNil(iv, true) // handle nil pointers also
	if ok {
		halt.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {
	// case nil:
	// case Selfer:
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
		// not addressable byte slice, so do not decode into it past the length
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
		// we can't check non-predefined types, as they might be a Selfer or extension.
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

// decodeValue MUST be called by the actual value we want to decode into,
// not its addr or a reference to it.
//
// This way, we know if it is itself a pointer, and can handle nil in
// the stream effectively.
//
// Note that decodeValue will handle nil in the stream early, so that the
// subsequent calls i.e. kXXX methods, etc do not have to handle it themselves.
func (d *decoder[T]) decodeValue(rv reflect.Value, fn *decFn[T]) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
	} else {
		d.decodeValueNoCheckNil(rv, fn)
	}
}

func (d *decoder[T]) decodeValueNoCheckNil(rv reflect.Value, fn *decFn[T]) {
	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
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

func (d *decoder[T]) decodeAs(v interface{}, t reflect.Type, ext bool) {
	if ext {
		d.decodeValue(baseRV(v), d.fn(t))
	} else {
		d.decodeValue(baseRV(v), d.fnNoExt(t))
	}
}

func (d *decoder[T]) structFieldNotFound(index int, rvkencname string) {
	// Note: rvkencname is used only if there is an error, to pass into halt.errorf.
	// Consequently, it is ok to pass in a stringView
	// Since rvkencname may be a stringView, do NOT pass it to another function.
	if d.h.ErrorIfNoField {
		if index >= 0 {
			halt.errorInt("no matching struct field found when decoding stream array at index ", int64(index))
		} else if rvkencname != "" {
			halt.errorStr2("no matching struct field found when decoding stream map with key ", rvkencname)
		}
	}
	d.swallow()
}

// decodeBytesInto is a convenience delegate function to decDriver.DecodeBytes.
// It ensures that `in` is not a nil byte, before calling decDriver.DecodeBytes,
// as decDriver.DecodeBytes treats a nil as a hint to use its internal scratch buffer.
func (d *decoder[T]) decodeBytesInto(out []byte, mustFit bool) (v []byte, state dBytesIntoState) {
	v, att := d.d.DecodeBytes()
	if cap(v) == 0 || (att >= dBytesAttachViewZerocopy && !mustFit) {
		// no need to detach (since mustFit=false)
		// including v has no capacity (covers v == nil and []byte{})
		return
	}
	if len(v) == 0 {
		v = zeroByteSlice // cannot be re-sliced/appended to
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

func (d *decoder[T]) rawBytes() (v []byte) {
	// ensure that this is not a view into the bytes
	// i.e. if necessary, make new copy always.
	v = d.d.nextValueBytes()
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v) // using copy here triggers make+copy optimization eliding memclr
		v = vv
	}
	return
}

func (d *decoder[T]) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.d.NumBytesRead(), false)
}

// NumBytesRead returns the number of bytes read
func (d *decoder[T]) NumBytesRead() int {
	return d.d.NumBytesRead()
}

// ---- container tracking
// Note: We update the .c after calling the callback.
// This way, the callback can know what the last status was.

// MARKER: do not call mapEnd if mapStart returns containerLenNil.

// MARKER: optimize decoding since all formats do not truly support all decDriver'ish operations.
// - Read(Map|Array)Start is only supported by all formats.
// - CheckBreak is only supported by json and cbor.
// - Read(Map|Array)End is only supported by json.
// - Read(Map|Array)Elem(Kay|Value) is only supported by json.
// Honor these in the code, to reduce the number of interface calls (even if empty).

func (d *decoder[T]) containerNext(j, containerLen int, hasLen bool) bool {
	// return (hasLen && (j < containerLen)) || (!hasLen && !d.d.CheckBreak())
	if hasLen {
		return j < containerLen
	}
	return !d.d.CheckBreak()
}

func (d *decoder[T]) mapElemKey(firstTime bool) {
	d.d.ReadMapElemKey(firstTime)
	d.c = containerMapKey
}

func (d *decoder[T]) mapElemValue() {
	d.d.ReadMapElemValue()
	d.c = containerMapValue
}

func (d *decoder[T]) mapEnd() {
	d.d.ReadMapEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoder[T]) arrayElem(firstTime bool) {
	d.d.ReadArrayElem(firstTime)
	d.c = containerArrayElem
}

func (d *decoder[T]) arrayEnd() {
	d.d.ReadArrayEnd()
	d.depthDecr()
	d.c = 0
}

func (d *decoder[T]) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {
	// The ext may support different types for performance e.g. int if no fractions, else float64
	// Consequently, best mode is:
	// - decode next value into an interface{}
	// - pass it to the UpdateExt
	var vv interface{}
	d.decode(&vv)
	ext.UpdateExt(v, vv)
	// rv := d.interfaceExtConvertAndDecodeGetRV(v, ext)
	// d.decodeValue(rv, nil)
	// ext.UpdateExt(v, rv2i(rv))
}

func (d *decoder[T]) fn(t reflect.Type) *decFn[T] {
	return d.dh.decFnViaBH(t, d.rtidFn, d.h, d.fp, true)
}

func (d *decoder[T]) fnNoExt(t reflect.Type) *decFn[T] {
	return d.dh.decFnViaBH(t, d.rtidFnNoExt, d.h, d.fp, false)
}

// ----

func (helperDecDriver[T]) newDecoderBytes(in []byte, h Handle) *decoder[T] {
	var c1 decoder[T]
	c1.bytes = true
	c1.init(h)
	c1.ResetBytes(in) // MARKER check for error
	return &c1
}

func (helperDecDriver[T]) newDecoderIO(in io.Reader, h Handle) *decoder[T] {
	var c1 decoder[T]
	c1.init(h)
	c1.Reset(in)
	return &c1
}

// ----

func (helperDecDriver[T]) decFnloadFastpathUnderlying(ti *typeInfo, fp *fastpathDs[T]) (f *fastpathD[T], u reflect.Type) {
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

func (helperDecDriver[T]) decFindRtidFn(s []decRtidFn[T], rtid uintptr) (i uint, fn *decFn[T]) {
	// binary search. Adapted from sort/search.go. Use goto (not for loop) to allow inlining.
	var h uint // var h, i uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1 // avoid overflow when computing h // h = i + (j-i)/2
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

func (helperDecDriver[T]) decFromRtidFnSlice(fns *atomicRtidFnSlice) (s []decRtidFn[T]) {
	if v := fns.load(); v != nil {
		s = *(lowLevelToPtr[[]decRtidFn[T]](v))
	}
	return
}

func (dh helperDecDriver[T]) decFnViaBH(rt reflect.Type, fns *atomicRtidFnSlice, x *BasicHandle, fp *fastpathDs[T],
	checkExt bool) (fn *decFn[T]) {
	return dh.decFnVia(rt, fns, x.typeInfos(), &x.mu, x.extHandle, fp,
		checkExt, x.CheckCircularRef, x.timeBuiltin, x.binaryHandle, x.jsonHandle)
}

func (dh helperDecDriver[T]) decFnVia(rt reflect.Type, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDs[T],
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFn[T]) {
	rtid := rt2id(rt)
	var sp []decRtidFn[T] = dh.decFromRtidFnSlice(fns)
	if sp != nil {
		_, fn = dh.decFindRtidFn(sp, rtid)
	}
	if fn == nil {
		fn = dh.decFnViaLoader(rt, rtid, fns, tinfos, mu, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	}
	return
}

func (dh helperDecDriver[T]) decFnViaLoader(rt reflect.Type, rtid uintptr, fns *atomicRtidFnSlice,
	tinfos *TypeInfos, mu *sync.Mutex, exth extHandle, fp *fastpathDs[T],
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFn[T]) {

	fn = dh.decFnLoad(rt, rtid, tinfos, exth, fp, checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json)
	var sp []decRtidFn[T]
	mu.Lock()
	sp = dh.decFromRtidFnSlice(fns)
	// since this is an atomic load/store, we MUST use a different array each time,
	// else we have a data race when a store is happening simultaneously with a decFindRtidFn call.
	if sp == nil {
		sp = []decRtidFn[T]{{rtid, fn}}
		fns.store(ptrToLowLevel(&sp))
	} else {
		idx, fn2 := dh.decFindRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]decRtidFn[T], len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = decRtidFn[T]{rtid, fn}
			fns.store(ptrToLowLevel(&sp2))
		}
	}
	mu.Unlock()
	return
}

func (dh helperDecDriver[T]) decFnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos,
	exth extHandle, fp *fastpathDs[T],
	checkExt, checkCircularRef, timeBuiltin, binaryEncoding, json bool) (fn *decFn[T]) {
	fn = new(decFn[T])
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	// anything can be an extension except the built-in ones: time, raw and rawext.
	// ensure we check for these types, then if extension, before checking if
	// it implements one of the pre-declared interfaces.

	fi.addrDf = true

	if rtid == timeTypId && timeBuiltin {
		fn.fd = (*decoder[T]).kTime
	} else if rtid == rawTypId {
		fn.fd = (*decoder[T]).raw
	} else if rtid == rawExtTypId {
		fn.fd = (*decoder[T]).rawExt
		fi.addrD = true
	} else if xfFn := exth.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fd = (*decoder[T]).ext
		fi.addrD = true
	} else if ti.flagSelfer || ti.flagSelferPtr {
		fn.fd = (*decoder[T]).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
	} else if supportMarshalInterfaces && binaryEncoding &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fd = (*decoder[T]).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding && json &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {
		//If JSON, we should check JSONMarshal before textMarshal
		fn.fd = (*decoder[T]).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
	} else if supportMarshalInterfaces && !binaryEncoding &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fd = (*decoder[T]).textUnmarshal
		fi.addrD = ti.flagTextUnmarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			var rtid2 uintptr
			if !ti.flagHasPkgPath { // un-named type (slice or mpa or array)
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key) // ti.key for arrays = reflect.SliceOf(ti.elem)
				}
				if idx, ok := fastpathAvIndex(rtid2); ok {
					fn.fd = fp[idx].decfn
					fi.addrD = true
					fi.addrDf = false
					if rk == reflect.Array {
						fi.addrD = false // decode directly into array value (slice made from it)
					}
				}
			} else { // named type (with underlying type of map or slice or array)
				// try to use mapping for underlying type
				xfe, xrt := dh.decFnloadFastpathUnderlying(ti, fp)
				if xfe != nil {
					xfnf2 := xfe.decfn
					if rk == reflect.Array {
						fi.addrD = false // decode directly into array value (slice made from it)
						fn.fd = func(d *decoder[T], xf *decFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false // meaning it can be an address(ptr) or a value
						xptr2rt := reflect.PointerTo(xrt)
						fn.fd = func(d *decoder[T], xf *decFnInfo, xrv reflect.Value) {
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
				fn.fd = (*decoder[T]).kBool
			case reflect.String:
				fn.fd = (*decoder[T]).kString
			case reflect.Int:
				fn.fd = (*decoder[T]).kInt
			case reflect.Int8:
				fn.fd = (*decoder[T]).kInt8
			case reflect.Int16:
				fn.fd = (*decoder[T]).kInt16
			case reflect.Int32:
				fn.fd = (*decoder[T]).kInt32
			case reflect.Int64:
				fn.fd = (*decoder[T]).kInt64
			case reflect.Uint:
				fn.fd = (*decoder[T]).kUint
			case reflect.Uint8:
				fn.fd = (*decoder[T]).kUint8
			case reflect.Uint16:
				fn.fd = (*decoder[T]).kUint16
			case reflect.Uint32:
				fn.fd = (*decoder[T]).kUint32
			case reflect.Uint64:
				fn.fd = (*decoder[T]).kUint64
			case reflect.Uintptr:
				fn.fd = (*decoder[T]).kUintptr
			case reflect.Float32:
				fn.fd = (*decoder[T]).kFloat32
			case reflect.Float64:
				fn.fd = (*decoder[T]).kFloat64
			case reflect.Complex64:
				fn.fd = (*decoder[T]).kComplex64
			case reflect.Complex128:
				fn.fd = (*decoder[T]).kComplex128
			case reflect.Chan:
				fn.fd = (*decoder[T]).kChan
			case reflect.Slice:
				fn.fd = (*decoder[T]).kSlice
			case reflect.Array:
				fi.addrD = false // decode directly into array value (slice made from it)
				fn.fd = (*decoder[T]).kArray
			case reflect.Struct:
				if ti.simple {
					fn.fd = (*decoder[T]).kStructSimple
				} else {
					fn.fd = (*decoder[T]).kStruct
				}
			case reflect.Map:
				fn.fd = (*decoder[T]).kMap
			case reflect.Interface:
				// encode: reflect.Interface are handled already by preEncodeValue
				fn.fd = (*decoder[T]).kInterface
			default:
				// reflect.Ptr and reflect.Interface are handled already by preEncodeValue
				fn.fd = (*decoder[T]).kErr
			}
		}
	}
	return
}
