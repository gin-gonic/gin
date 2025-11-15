//go:build notmono || codec.notmono

// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

/*
Msgpack-c implementation powers the c, c++, python, ruby, etc libraries.
We need to maintain compatibility with it and how it encodes integer values
without caring about the type.

For compatibility with behaviour of msgpack-c reference implementation:
  - Go intX (>0) and uintX
       IS ENCODED AS
    msgpack +ve fixnum, unsigned
  - Go intX (<0)
       IS ENCODED AS
    msgpack -ve fixnum, signed
*/

package codec

import (
	"io"
	"math"
	"reflect"
	"time"
	"unicode/utf8"
)

//---------------------------------------------

type msgpackEncDriver[T encWriter] struct {
	noBuiltInTypes
	encDriverNoopContainerWriter
	encDriverNoState
	encDriverContainerNoTrackerT
	encInit2er

	h *MsgpackHandle
	e *encoderBase
	w T
	// x [8]byte
}

func (e *msgpackEncDriver[T]) EncodeNil() {
	e.w.writen1(mpNil)
}

func (e *msgpackEncDriver[T]) EncodeInt(i int64) {
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

func (e *msgpackEncDriver[T]) EncodeUint(i uint64) {
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

func (e *msgpackEncDriver[T]) EncodeBool(b bool) {
	if b {
		e.w.writen1(mpTrue)
	} else {
		e.w.writen1(mpFalse)
	}
}

func (e *msgpackEncDriver[T]) EncodeFloat32(f float32) {
	e.w.writen1(mpFloat)
	e.w.writen4(bigen.PutUint32(math.Float32bits(f)))
}

func (e *msgpackEncDriver[T]) EncodeFloat64(f float64) {
	e.w.writen1(mpDouble)
	e.w.writen8(bigen.PutUint64(math.Float64bits(f)))
}

func (e *msgpackEncDriver[T]) EncodeTime(t time.Time) {
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

func (e *msgpackEncDriver[T]) EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (e *msgpackEncDriver[T]) EncodeRawExt(re *RawExt) {
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
}

func (e *msgpackEncDriver[T]) encodeExtPreamble(xtag byte, l int) {
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

func (e *msgpackEncDriver[T]) WriteArrayStart(length int) {
	e.writeContainerLen(msgpackContainerList, length)
}

func (e *msgpackEncDriver[T]) WriteMapStart(length int) {
	e.writeContainerLen(msgpackContainerMap, length)
}

func (e *msgpackEncDriver[T]) WriteArrayEmpty() {
	// e.WriteArrayStart(0) = e.writeContainerLen(msgpackContainerList, 0)
	e.w.writen1(mpFixArrayMin)
}

func (e *msgpackEncDriver[T]) WriteMapEmpty() {
	// e.WriteMapStart(0) = e.writeContainerLen(msgpackContainerMap, 0)
	e.w.writen1(mpFixMapMin)
}

func (e *msgpackEncDriver[T]) EncodeString(s string) {
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

func (e *msgpackEncDriver[T]) EncodeStringNoEscape4Json(v string) { e.EncodeString(v) }

func (e *msgpackEncDriver[T]) EncodeStringBytesRaw(bs []byte) {
	if e.h.WriteExt {
		e.writeContainerLen(msgpackContainerBin, len(bs))
	} else {
		e.writeContainerLen(msgpackContainerRawLegacy, len(bs))
	}
	if len(bs) > 0 {
		e.w.writeb(bs)
	}
}

func (e *msgpackEncDriver[T]) EncodeBytes(v []byte) {
	if v == nil {
		e.writeNilBytes()
		return
	}
	e.EncodeStringBytesRaw(v)
}

func (e *msgpackEncDriver[T]) writeNilOr(v byte) {
	if !e.h.NilCollectionToZeroLength {
		v = mpNil
	}
	e.w.writen1(v)
}

func (e *msgpackEncDriver[T]) writeNilArray() {
	e.writeNilOr(mpFixArrayMin)
}

func (e *msgpackEncDriver[T]) writeNilMap() {
	e.writeNilOr(mpFixMapMin)
}

func (e *msgpackEncDriver[T]) writeNilBytes() {
	e.writeNilOr(mpFixStrMin)
}

func (e *msgpackEncDriver[T]) writeContainerLen(ct msgpackContainerType, l int) {
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

//---------------------------------------------

type msgpackDecDriver[T decReader] struct {
	decDriverNoopContainerReader
	decDriverNoopNumberHelper
	decInit2er

	h *MsgpackHandle
	d *decoderBase
	r T

	bdAndBdread
	// bytes bool
	noBuiltInTypes
}

// Note: This returns either a primitive (int, bool, etc) for non-containers,
// or a containerType, or a specific type denoting nil or extension.
// It is called when a nil interface{} is passed, leaving it up to the DecDriver
// to introspect the stream and decide how best to decode.
// It deciphers the value by looking at the stream first.
func (d *msgpackDecDriver[T]) DecodeNaked() {
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
			// positive fixnum (always signed)
			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
			// negative fixnum
			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd == mpStr8, bd == mpStr16, bd == mpStr32, bd >= mpFixStrMin && bd <= mpFixStrMax:
			d.d.fauxUnionReadRawBytes(d, d.h.WriteExt, d.h.RawToString) //, d.h.ZeroCopy)
			// if d.h.WriteExt || d.h.RawToString {
			// 	n.v = valueTypeString
			// 	n.s = d.d.stringZC(d.DecodeStringAsBytes())
			// } else {
			// 	n.v = valueTypeBytes
			// 	n.l = d.DecodeBytes([]byte{})
			// }
		case bd == mpBin8, bd == mpBin16, bd == mpBin32:
			d.d.fauxUnionReadRawBytes(d, false, d.h.RawToString) //, d.h.ZeroCopy)
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

func (d *msgpackDecDriver[T]) nextValueBytes() (v []byte) {
	if !d.bdRead {
		d.readNextBd()
	}
	d.r.startRecording()
	d.nextValueBytesBdReadR()
	v = d.r.stopRecording()
	d.bdRead = false
	return
}

func (d *msgpackDecDriver[T]) nextValueBytesBdReadR() {
	bd := d.bd

	var clen uint

	switch bd {
	case mpNil, mpFalse, mpTrue: // pass
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
		d.r.readn1() // tag
		d.r.readn1()
	case mpFixExt2:
		d.r.readn1() // tag
		d.r.skip(2)
	case mpFixExt4:
		d.r.readn1() // tag
		d.r.skip(4)
	case mpFixExt8:
		d.r.readn1() // tag
		d.r.skip(8)
	case mpFixExt16:
		d.r.readn1() // tag
		d.r.skip(16)
	case mpExt8:
		clen = uint(d.r.readn1())
		d.r.readn1() // tag
		d.r.skip(clen)
	case mpExt16:
		x := d.r.readn2()
		clen = uint(bigen.Uint16(x))
		d.r.readn1() // tag
		d.r.skip(clen)
	case mpExt32:
		x := d.r.readn4()
		clen = uint(bigen.Uint32(x))
		d.r.readn1() // tag
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
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax: // pass
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax: // pass
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

func (d *msgpackDecDriver[T]) decFloat4Int32() (f float32) {
	fbits := bigen.Uint32(d.r.readn4())
	f = math.Float32frombits(fbits)
	if !noFrac32(fbits) {
		halt.errorf("assigning integer value from float32 with a fraction: %v", f)
	}
	return
}

func (d *msgpackDecDriver[T]) decFloat4Int64() (f float64) {
	fbits := bigen.Uint64(d.r.readn8())
	f = math.Float64frombits(fbits)
	if !noFrac64(fbits) {
		halt.errorf("assigning integer value from float64 with a fraction: %v", f)
	}
	return
}

// int can be decoded from msgpack type: intXXX or uintXXX
func (d *msgpackDecDriver[T]) DecodeInt64() (i int64) {
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

// uint can be decoded from msgpack type: intXXX or uintXXX
func (d *msgpackDecDriver[T]) DecodeUint64() (ui uint64) {
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

// float can either be decoded from msgpack type: float, double or intX
func (d *msgpackDecDriver[T]) DecodeFloat64() (f float64) {
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

// bool can be decoded from bool, fixnum 0 or 1.
func (d *msgpackDecDriver[T]) DecodeBool() (b bool) {
	if d.advanceNil() {
		return
	}
	if d.bd == mpFalse || d.bd == 0 {
		// b = false
	} else if d.bd == mpTrue || d.bd == 1 {
		b = true
	} else {
		halt.errorf("cannot decode bool: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
	}
	d.bdRead = false
	return
}

func (d *msgpackDecDriver[T]) DecodeBytes() (bs []byte, state dBytesAttachState) {
	if d.advanceNil() {
		return
	}

	var cond bool
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		clen = d.readContainerLen(msgpackContainerBin) // binary
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		clen = d.readContainerLen(msgpackContainerStr) // string/raw
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

func (d *msgpackDecDriver[T]) DecodeStringAsBytes() (out []byte, state dBytesAttachState) {
	out, state = d.DecodeBytes()
	if d.h.ValidateUnicode && !utf8.Valid(out) {
		halt.errorf("DecodeStringAsBytes: invalid UTF-8: %s", out)
	}
	return
}

func (d *msgpackDecDriver[T]) readNextBd() {
	d.bd = d.r.readn1()
	d.bdRead = true
}

func (d *msgpackDecDriver[T]) advanceNil() (null bool) {
	if !d.bdRead {
		d.readNextBd()
	}
	if d.bd == mpNil {
		d.bdRead = false
		return true // null = true
	}
	return
}

func (d *msgpackDecDriver[T]) TryNil() (v bool) {
	return d.advanceNil()
}

func (d *msgpackDecDriver[T]) ContainerType() (vt valueType) {
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
		if d.h.WriteExt || d.h.RawToString { // UTF-8 string (new spec)
			return valueTypeString
		}
		return valueTypeBytes // raw (old spec)
	} else if bd == mpArray16 || bd == mpArray32 || (bd >= mpFixArrayMin && bd <= mpFixArrayMax) {
		return valueTypeArray
	} else if bd == mpMap16 || bd == mpMap32 || (bd >= mpFixMapMin && bd <= mpFixMapMax) {
		return valueTypeMap
	}
	return valueTypeUnset
}

func (d *msgpackDecDriver[T]) readContainerLen(ct msgpackContainerType) (clen int) {
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

func (d *msgpackDecDriver[T]) ReadMapStart() int {
	if d.advanceNil() {
		return containerLenNil
	}
	return d.readContainerLen(msgpackContainerMap)
}

func (d *msgpackDecDriver[T]) ReadArrayStart() int {
	if d.advanceNil() {
		return containerLenNil
	}
	return d.readContainerLen(msgpackContainerList)
}

func (d *msgpackDecDriver[T]) readExtLen() (clen int) {
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

func (d *msgpackDecDriver[T]) DecodeTime() (t time.Time) {
	// decode time from string bytes or ext
	if d.advanceNil() {
		return
	}
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 {
		clen = d.readContainerLen(msgpackContainerBin) // binary
	} else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) {
		clen = d.readContainerLen(msgpackContainerStr) // string/raw
	} else {
		// expect to see mpFixExt4,-1 OR mpFixExt8,-1 OR mpExt8,12,-1
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

func (d *msgpackDecDriver[T]) decodeTime(clen int) (t time.Time) {
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

func (d *msgpackDecDriver[T]) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
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

func (d *msgpackDecDriver[T]) DecodeRawExt(re *RawExt) {
	xbs, realxtag, state, ok := d.decodeExtV(false, 0)
	if !ok {
		return
	}
	re.Tag = uint64(realxtag)
	re.setData(xbs, state >= dBytesAttachViewZerocopy)
}

func (d *msgpackDecDriver[T]) decodeExtV(verifyTag bool, xtagIn uint64) (xbs []byte, xtag byte, bstate dBytesAttachState, ok bool) {
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
		// zerocopy = d.d.bytes
	}
	d.bdRead = false
	ok = true
	return
}

// ----
//
// The following below are similar across all format files (except for the format name).
//
// We keep them together here, so that we can easily copy and compare.

// ----

func (d *msgpackEncDriver[T]) init(hh Handle, shared *encoderBase, enc encoderI) (fp interface{}) {
	callMake(&d.w)
	d.h = hh.(*MsgpackHandle)
	d.e = shared
	if shared.bytes {
		fp = msgpackFpEncBytes
	} else {
		fp = msgpackFpEncIO
	}
	// d.w.init()
	d.init2(enc)
	return
}

func (e *msgpackEncDriver[T]) writeBytesAsis(b []byte) { e.w.writeb(b) }

// func (e *msgpackEncDriver[T]) writeStringAsisDblQuoted(v string) { e.w.writeqstr(v) }

func (e *msgpackEncDriver[T]) writerEnd() { e.w.end() }

func (e *msgpackEncDriver[T]) resetOutBytes(out *[]byte) {
	e.w.resetBytes(*out, out)
}

func (e *msgpackEncDriver[T]) resetOutIO(out io.Writer) {
	e.w.resetIO(out, e.h.WriterBufferSize, &e.e.blist)
}

// ----

func (d *msgpackDecDriver[T]) init(hh Handle, shared *decoderBase, dec decoderI) (fp interface{}) {
	callMake(&d.r)
	d.h = hh.(*MsgpackHandle)
	d.d = shared
	if shared.bytes {
		fp = msgpackFpDecBytes
	} else {
		fp = msgpackFpDecIO
	}
	// d.r.init()
	d.init2(dec)
	return
}

func (d *msgpackDecDriver[T]) NumBytesRead() int {
	return int(d.r.numread())
}

func (d *msgpackDecDriver[T]) resetInBytes(in []byte) {
	d.r.resetBytes(in)
}

func (d *msgpackDecDriver[T]) resetInIO(r io.Reader) {
	d.r.resetIO(r, d.h.ReaderBufferSize, d.h.MaxInitLen, &d.d.blist)
}

// ---- (custom stanza)

func (d *msgpackDecDriver[T]) descBd() string {
	return sprintf("%v (%s)", d.bd, mpdesc(d.bd))
}

func (d *msgpackDecDriver[T]) DecodeFloat32() (f float32) {
	return float32(chkOvf.Float32V(d.DecodeFloat64()))
}
