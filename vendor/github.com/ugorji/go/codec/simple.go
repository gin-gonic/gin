//go:build notmono || codec.notmono

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
	"math"
	"reflect"
	"time"
)

type simpleEncDriver[T encWriter] struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *SimpleHandle
	e *encoderBase
	// b [8]byte
	w T
}

func (e *simpleEncDriver[T]) EncodeNil() {
	e.w.writen1(simpleVdNil)
}

func (e *simpleEncDriver[T]) EncodeBool(b bool) {
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

func (e *simpleEncDriver[T]) EncodeFloat32(f float32) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat32)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *simpleEncDriver[T]) EncodeFloat64(f float64) {
	if e.h.EncZeroValuesAsNil && e.e.c != containerMapKey && f == 0.0 {
		e.EncodeNil()
		return
	}
	e.w.writen1(simpleVdFloat64)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *simpleEncDriver[T]) EncodeInt(v int64) {
	if v < 0 {
		e.encUint(uint64(-v), simpleVdNegInt)
	} else {
		e.encUint(uint64(v), simpleVdPosInt)
	}
}

func (e *simpleEncDriver[T]) EncodeUint(v uint64) {
	e.encUint(v, simpleVdPosInt)
}

func (e *simpleEncDriver[T]) encUint(v uint64, bd uint8) {
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
	} else { // if v <= math.MaxUint64 {
		e.w.writen1(bd + 3)
		e.w.writen8(bigen.PutUint64(v))
	}
}

func (e *simpleEncDriver[T]) encLen(bd byte, length int) {
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

func (e *simpleEncDriver[T]) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *simpleEncDriver[T]) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *simpleEncDriver[T]) encodeExtPreamble(xtag byte, length int) {
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
}

func (e *simpleEncDriver[T]) WriteArrayStart(length int) {
	e.encLen(simpleVdArray, length)
}

func (e *simpleEncDriver[T]) WriteMapStart(length int) {
	e.encLen(simpleVdMap, length)
}

func (e *simpleEncDriver[T]) WriteArrayEmpty() {
	// e.WriteArrayStart(0) = e.encLen(simpleVdArray, 0)
	e.w.writen1(simpleVdArray)
}

func (e *simpleEncDriver[T]) WriteMapEmpty() {
	// e.WriteMapStart(0) = e.encLen(simpleVdMap, 0)
	e.w.writen1(simpleVdMap)
}

func (e *simpleEncDriver[T]) EncodeString(v string) {
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

func (e *simpleEncDriver[T]) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *simpleEncDriver[T]) EncodeStringBytesRaw(v []byte) {
	// if e.h.EncZeroValuesAsNil && e.c != containerMapKey && v == nil {
	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
}

func (e *simpleEncDriver[T]) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *simpleEncDriver[T]) encodeNilBytes() {
	b := byte(simpleVdNil)
	if e.h.NilCollectionToZeroLength {
		b = simpleVdArray
	}
	e.w.writen1(b)
}

func (e *simpleEncDriver[T]) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = simpleVdNil
	}
	e.w.writen1(v)
}

func (e *simpleEncDriver[T]) writeNilArray() {
	e.writeNilOr(simpleVdArray)
}

func (e *simpleEncDriver[T]) writeNilMap() {
	e.writeNilOr(simpleVdMap)
}

func (e *simpleEncDriver[T]) writeNilBytes() {
	e.writeNilOr(simpleVdByteArray)
}

func (e *simpleEncDriver[T]) EncodeTime(t time.Time) {
	// if e.h.EncZeroValuesAsNil && e.c != containerMapKey && t.IsZero() {
	if t.IsZero() {
		e.EncodeNil()
		return
	}
	v, err := t.MarshalBinary()
	halt.onerror(err)
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
}

//------------------------------------

type simpleDecDriver[T decReader] struct {
	h *SimpleHandle
	d *decoderBase
	r T

	bdAndBdread
	// bytes bool

	noBuiltInTypes
	// decDriverNoopNumberHelper
	decDriverNoopContainerReader
	decInit2er

	// ds interface{} // must be *decoder[simpleDecDriverM[bytes...]]
}

func (d *simpleDecDriver[T]) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *simpleDecDriver[T]) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == simpleVdNil {
		d.bdRead = false
		return true // null = true
	}
	return
}

func (d *simpleDecDriver[T]) ContainerType() (vt valueType) {
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

func (d *simpleDecDriver[T]) TryNil() bool {
	return d.advanceNil()
}

func (d *simpleDecDriver[T]) decFloat() (f float64, ok bool) {
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

func (d *simpleDecDriver[T]) decInteger() (ui uint64, neg, ok bool) {
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
		// halt.errorf("integer only valid from pos/neg integer1..8. Invalid descriptor: %v", d.bd)
	}
	// DO NOT do this check below, because callers may only want the unsigned value:
	//
	// if ui > math.MaxInt64 {
	// 	halt.errorf("decIntAny: Integer out of range for signed int64: %v", ui)
	//		return
	// }
	return
}

func (d *simpleDecDriver[T]) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

// bool can be decoded from bool only (single byte).
func (d *simpleDecDriver[T]) DecodeBool() (b bool) {
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

func (d *simpleDecDriver[T]) ReadMapStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriver[T]) ReadArrayStart() (length int) {
	if d.advanceNil() {
		return containerLenNil
	}
	d.bdRead = false
	return d.decLen()
}

func (d *simpleDecDriver[T]) uint2Len(ui uint64) int {
	if chkOvf.Uint(ui, intBitsize) {
		halt.errorf("overflow integer: %v", ui)
	}
	return int(ui)
}

func (d *simpleDecDriver[T]) decLen() int {
	switch d.bd & 7 { // d.bd % 8 {
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

func (d *simpleDecDriver[T]) DecodeStringAsBytes() ([]byte, dBytesAttachState) {
	return d.DecodeBytes()
}

func (d *simpleDecDriver[T]) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}
	var cond bool
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
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

func (d *simpleDecDriver[T]) DecodeTime() (t time.Time) {
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

func (d *simpleDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *simpleDecDriver[T]) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *simpleDecDriver[T]) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
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

func (d *simpleDecDriver[T]) DecodeNaked() {
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
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString) //, d.h.ZeroCopy)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(uint(l))
		// MARKER: not necessary to detach for extensions
		// var useBuf bool
		// n.l, useBuf = d.r.readxb(uint(l))
		// n.a = d.d.attachState(useBuf)
		// n.l = d.d.detach2Bytes(n.l, nil, n.a)
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

func (d *simpleDecDriver[T]) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *simpleDecDriver[T]) nextValueBytesBdReadR() {
	c := d.bd

	var length uint

	switch c {
	case simpleVdNil, simpleVdFalse, simpleVdTrue, simpleVdString, simpleVdByteArray:
		// pass
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
		switch c & 7 { // c % 8 {
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
			d.r.readn1() // tag
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

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *simpleEncDriver[T]) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*SimpleHandle)
	d.e = shared
	if shared.bytes {
		fp = simpleFpEncBytes
	} else {
		fp = simpleFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *simpleEncDriver[T]) writeBytesAsis(b []byte) { e.w.writeb(b) }

func (e *simpleEncDriver[T]) writerEnd() { e.w.end() }

func (e *simpleEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *simpleEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *simpleDecDriver[T]) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*SimpleHandle)
	d.d = shared
	if shared.bytes {
		fp = simpleFpDecBytes
	} else {
		fp = simpleFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *simpleDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *simpleDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *simpleDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

// ---- (custom stanza)

func (d *simpleDecDriver[T]) descBd() string {
	return sprintf("%v (%s)", d.bd, simpledesc(d.bd))
}

func (d *simpleDecDriver[T]) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
