//go:build notmono || codec.notmono

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
	"math"
	"reflect"
	"time"
	"unicode/utf8"
)

type bincEncDriver[T encWriter] struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverContainerNoTrackerT
	encInit2er

	h *BincHandle
	e *encoderBase
	w T
	bincEncState
}

func (e *bincEncDriver[T]) EncodeNil() {
	e.w.writen1(bincBdNil)
}

func (e *bincEncDriver[T]) EncodeTime(t time.Time) {
	if t.IsZero() {
		e.EncodeNil()
	} else {
		bs := bincEncodeTime(t)
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	}
}

func (e *bincEncDriver[T]) EncodeBool(b bool) {
	if b {
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	} else {
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	}
}

func (e *bincEncDriver[T]) encSpFloat(f float64) (done bool) {
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

func (e *bincEncDriver[T]) EncodeFloat32(f float32) {
	if !e.encSpFloat(float64(f)) {
		e.w.writen1(bincVdFloat<<4 | bincFlBin32)
		e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
	}
}

func (e *bincEncDriver[T]) EncodeFloat64(f float64) {
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

func (e *bincEncDriver[T]) encIntegerPrune32(bd byte, pos bool, v uint64) {
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

func (e *bincEncDriver[T]) encIntegerPrune64(bd byte, pos bool, v uint64) {
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

func (e *bincEncDriver[T]) EncodeInt(v int64) {
	if v >= 0 {
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	} else if v == -1 {
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	} else {
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	}
}

func (e *bincEncDriver[T]) EncodeUint(v uint64) {
	e.encUint(bincVdPosInt<<4, true, v)
}

func (e *bincEncDriver[T]) encUint(bd byte, pos bool, v uint64) {
	if v == 0 {
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	} else if pos && v >= 1 && v <= 16 {
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	} else if v <= math.MaxUint8 {
		e.w.writen2(bd, byte(v)) // bd|0x0
	} else if v <= math.MaxUint16 {
		e.w.writen1(bd | 0x01)
		e.w.writen2(bigen.PutUint16(uint16(v)))
	} else if v <= math.MaxUint32 {
		e.encIntegerPrune32(bd, pos, v)
	} else {
		e.encIntegerPrune64(bd, pos, v)
	}
}

func (e *bincEncDriver[T]) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *bincEncDriver[T]) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *bincEncDriver[T]) encodeExtPreamble(xtag byte, length int) {
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
}

func (e *bincEncDriver[T]) WriteArrayStart(length int) {
	e.encLen(bincVdArray<<4, uint64(length))
}

func (e *bincEncDriver[T]) WriteMapStart(length int) {
	e.encLen(bincVdMap<<4, uint64(length))
}

func (e *bincEncDriver[T]) WriteArrayEmpty() {
	// e.WriteArrayStart(0) = e.encLen(bincVdArray<<4, 0)
	e.w.writen1(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriver[T]) WriteMapEmpty() {
	// e.WriteMapStart(0) = e.encLen(bincVdMap<<4, 0)
	e.w.writen1(bincVdMap<<4 | uint8(0+4))
}

func (e *bincEncDriver[T]) EncodeSymbol(v string) {
	//symbols only offer benefit when string length > 1.
	//This is because strings with length 1 take only 2 bytes to store
	//(bd with embedded length, and single byte for string val).

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
			// lenprec = 0
		} else if l <= math.MaxUint16 {
			lenprec = 1
		} else if int64(l) <= math.MaxUint32 {
			lenprec = 2
		} else {
			lenprec = 3
		}
		if ui <= math.MaxUint8 {
			e.w.writen2(bincVdSymbol<<4|0x4|lenprec, byte(ui)) // bincVdSymbol<<4|0x0|0x4|lenprec
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

func (e *bincEncDriver[T]) EncodeString(v string) {
	if e.h.StringToRaw {
		e.encLen(bincVdByteArray<<4, uint64(len(v)))
		if len(v) > 0 {
			e.w.writestr(v)
		}
		return
	}
	e.EncodeStringEnc(cUTF8, v)
}

func (e *bincEncDriver[T]) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *bincEncDriver[T]) EncodeStringEnc(c charEncoding, v string) {
	if e.e.c == containerMapKey && c == cUTF8 && (e.h.AsSymbols == 1) {
		e.EncodeSymbol(v)
		return
	}
	e.encLen(bincVdString<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writestr(v)
	}
}

func (e *bincEncDriver[T]) EncodeStringBytesRaw(v []byte) {
	e.encLen(bincVdByteArray<<4, uint64(len(v)))
	if len(v) > 0 {
		e.w.writeb(v)
	}
}

func (e *bincEncDriver[T]) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *bincEncDriver[T]) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = bincBdNil
	}
	e.w.writen1(v)
}

func (e *bincEncDriver[T]) writeNilArray() {
	e.writeNilOr(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriver[T]) writeNilMap() {
	e.writeNilOr(bincVdMap<<4 | uint8(0+4))
}

func (e *bincEncDriver[T]) writeNilBytes() {
	e.writeNilOr(bincVdArray<<4 | uint8(0+4))
}

func (e *bincEncDriver[T]) encBytesLen(c charEncoding, length uint64) {
	// MARKER: we currently only support UTF-8 (string) and RAW (bytearray).
	// We should consider supporting bincUnicodeOther.

	if c == cRAW {
		e.encLen(bincVdByteArray<<4, length)
	} else {
		e.encLen(bincVdString<<4, length)
	}
}

func (e *bincEncDriver[T]) encLen(bd byte, l uint64) {
	if l < 12 {
		e.w.writen1(bd | uint8(l+4))
	} else {
		e.encLenNumber(bd, l)
	}
}

func (e *bincEncDriver[T]) encLenNumber(bd byte, v uint64) {
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

//------------------------------------

type bincDecDriver[T decReader] struct {
	decDriverNoopContainerReader
	// decDriverNoopNumberHelper
	decInit2er
	noBuiltInTypes

	h *BincHandle
	d *decoderBase
	r T

	bincDecState

	// bytes bool
}

func (d *bincDecDriver[T]) readNextBd() {
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
}

func (d *bincDecDriver[T]) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == bincBdNil {
		d.bdRead = false
		return true // null = true
	}
	return
}

func (d *bincDecDriver[T]) TryNil() bool {
	return d.advanceNil()
}

func (d *bincDecDriver[T]) ContainerType() (vt valueType) {
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

func (d *bincDecDriver[T]) DecodeTime() (t time.Time) {
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

func (d *bincDecDriver[T]) decFloatPruned(maxlen uint8) {
	l := d.r.readn1()
	if l > maxlen {
		halt.errorf("cannot read float - at most %v bytes used to represent float - received %v bytes", maxlen, l)
	}
	for i := l; i < maxlen; i++ {
		d.d.b[i] = 0
	}
	d.r.readb(d.d.b[0:l])
}

func (d *bincDecDriver[T]) decFloatPre32() (b [4]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn4()
	} else {
		d.decFloatPruned(4)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriver[T]) decFloatPre64() (b [8]byte) {
	if d.vs&0x8 == 0 {
		b = d.r.readn8()
	} else {
		d.decFloatPruned(8)
		copy(b[:], d.d.b[:])
	}
	return
}

func (d *bincDecDriver[T]) decFloatVal() (f float64) {
	switch d.vs & 0x7 {
	case bincFlBin32:
		f = float64(math.Float32frombits(bigen.Uint32(d.decFloatPre32())))
	case bincFlBin64:
		f = math.Float64frombits(bigen.Uint64(d.decFloatPre64()))
	default:
		// ok = false
		halt.errorf("read float supports only float32/64 - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	return
}

func (d *bincDecDriver[T]) decUint() (v uint64) {
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
		// lim := 7 - d.vs
		// bs := d.d.b[lim:8]
		// d.r.readb(bs)
		// var b [8]byte
		// copy(b[lim:], bs)
		// v = bigen.Uint64(b)
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

func (d *bincDecDriver[T]) uintBytes() (bs []byte) {
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

func (d *bincDecDriver[T]) decInteger() (ui uint64, neg, ok bool) {
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
			// i = 0
		} else if vs == bincSpNegOne {
			neg = true
			ui = 1
		} else {
			ok = false
			// halt.errorf("integer decode has invalid special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	} else {
		ok = false
		// halt.errorf("integer can only be decoded from int/uint. d.bd: 0x%x, d.vd: 0x%x", d.bd, d.vd)
	}
	return
}

func (d *bincDecDriver[T]) decFloat() (f float64, ok bool) {
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
			// halt.errorf("float - invalid special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	} else if vd == bincVdFloat {
		f = d.decFloatVal()
	} else {
		ok = false
	}
	return
}

func (d *bincDecDriver[T]) DecodeInt64() (i int64) {
	if d.advanceNil() {
		return
	}
	v1, v2, v3 := d.decInteger()
	i = decNegintPosintFloatNumberHelper{d}.int64(v1, v2, v3, false)
	d.bdRead = false
	return
}

func (d *bincDecDriver[T]) DecodeUint64() (ui uint64) {
	if d.advanceNil() {
		return
	}
	ui = decNegintPosintFloatNumberHelper{d}.uint64(d.decInteger())
	d.bdRead = false
	return
}

func (d *bincDecDriver[T]) DecodeFloat64() (f float64) {
	if d.advanceNil() {
		return
	}
	v1, v2 := d.decFloat()
	f = decNegintPosintFloatNumberHelper{d}.float64(v1, v2, false)
	d.bdRead = false
	return
}

func (d *bincDecDriver[T]) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == (bincVdSpecial | bincSpFalse) {
		// b = false
	} else if d.bd == (bincVdSpecial | bincSpTrue) {
		b = true
	} else {
		halt.errorf("bool - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	return
}

func (d *bincDecDriver[T]) ReadMapStart() (length int) {
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

func (d *bincDecDriver[T]) ReadArrayStart() (length int) {
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

func (d *bincDecDriver[T]) decLen() int {
	if d.vs > 3 {
		return int(d.vs - 4)
	}
	return int(d.decLenNumber())
}

func (d *bincDecDriver[T]) decLenNumber() (v uint64) {
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

// func (d *bincDecDriver[T]) decStringBytes(bs []byte, zerocopy bool) (bs2 []byte) {
func (d *bincDecDriver[T]) DecodeStringAsBytes() (bs []byte, state dBytesAttachState) {
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
		// zerocopy doesn't apply for symbols,
		// as the values must be stored in a table for later use.
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
			// As we are using symbols, do not store any part of
			// the parameter bs in the map, as it might be a shared buffer.
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

func (d *bincDecDriver[T]) DecodeBytes() (bs []byte, state dBytesAttachState) {
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

func (d *bincDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *bincDecDriver[T]) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *bincDecDriver[T]) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
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
		// zerocopy = d.d.bytes
	} else if d.vd == bincVdByteArray {
		xbs, bstate = d.DecodeBytes()
	} else {
		halt.errorf("ext expects extensions or byte array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	}
	d.bdRead = false
	ok = true
	return
}

func (d *bincDecDriver[T]) DecodeNaked() {
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
			n.u = uint64(0) // int8(0)
		case bincSpNegOne:
			n.v = valueTypeInt
			n.i = int64(-1) // int8(-1)
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt:
		n.v = valueTypeUint
		n.u = uint64(int8(d.vs)) + 1 // int8(d.vs) + 1
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
		d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString) //, d.h.ZeroCopy)
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

func (d *bincDecDriver[T]) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

// func (d *bincDecDriver[T]) nextValueBytesR(v0 []byte) (v []byte) {
// 	d.readNextBd()
// 	v = v0
// 	var h decNextValueBytesHelper
// 	h.append1(&v, d.bytes, d.bd)
// 	return d.nextValueBytesBdReadR(v)
// }

func (d *bincDecDriver[T]) nextValueBytesBdReadR() {
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
		case bincSpNil, bincSpFalse, bincSpTrue, bincSpNan, bincSpPosInf: // pass
		case bincSpNegInf, bincSpZeroFloat, bincSpZero, bincSpNegOne: // pass
		default:
			halt.errorf("cannot infer value - unrecognized special value %x-%x/%s", d.vd, d.vs, bincdesc(d.vd, d.vs))
		}
	case bincVdSmallInt: // pass
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
		d.r.readn1() // tag
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

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *bincEncDriver[T]) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*BincHandle)
	d.e = shared
	if shared.bytes {
		fp = bincFpEncBytes
	} else {
		fp = bincFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *bincEncDriver[T]) writeBytesAsis(b []byte) { e.w.writeb(b) }

// func (e *bincEncDriver[T]) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }

func (e *bincEncDriver[T]) writerEnd() { e.w.end() }

func (e *bincEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *bincEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *bincDecDriver[T]) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*BincHandle)
	d.d = shared
	if shared.bytes {
		fp = bincFpDecBytes
	} else {
		fp = bincFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *bincDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *bincDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *bincDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

// ---- (custom stanza)

func (d *bincDecDriver[T]) descBd() string {
	return sprintf("%v (%s)", d.bd, bincdescbd(d.bd))
}

func (d *bincDecDriver[T]) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
