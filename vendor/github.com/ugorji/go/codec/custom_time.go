// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"time"
)

// EncodeTime encodes a time.Time as a []byte, including
// information on the instant in time and UTC offset.
//
// Format Description
//
//	A timestamp is composed of 3 components:
//
//	- secs: signed integer representing seconds since unix epoch
//	- nsces: unsigned integer representing fractional seconds as a
//	  nanosecond offset within secs, in the range 0 <= nsecs < 1e9
//	- tz: signed integer representing timezone offset in minutes east of UTC,
//	  and a dst (daylight savings time) flag
//
//	When encoding a timestamp, the first byte is the descriptor, which
//	defines which components are encoded and how many bytes are used to
//	encode secs and nsecs components. *If secs/nsecs is 0 or tz is UTC, it
//	is not encoded in the byte array explicitly*.
//
//	    Descriptor 8 bits are of the form `A B C DDD EE`:
//	        A:   Is secs component encoded? 1 = true
//	        B:   Is nsecs component encoded? 1 = true
//	        C:   Is tz component encoded? 1 = true
//	        DDD: Number of extra bytes for secs (range 0-7).
//	             If A = 1, secs encoded in DDD+1 bytes.
//	                 If A = 0, secs is not encoded, and is assumed to be 0.
//	                 If A = 1, then we need at least 1 byte to encode secs.
//	                 DDD says the number of extra bytes beyond that 1.
//	                 E.g. if DDD=0, then secs is represented in 1 byte.
//	                      if DDD=2, then secs is represented in 3 bytes.
//	        EE:  Number of extra bytes for nsecs (range 0-3).
//	             If B = 1, nsecs encoded in EE+1 bytes (similar to secs/DDD above)
//
//	Following the descriptor bytes, subsequent bytes are:
//
//	    secs component encoded in `DDD + 1` bytes (if A == 1)
//	    nsecs component encoded in `EE + 1` bytes (if B == 1)
//	    tz component encoded in 2 bytes (if C == 1)
//
//	secs and nsecs components are integers encoded in a BigEndian
//	2-complement encoding format.
//
//	tz component is encoded as 2 bytes (16 bits). Most significant bit 15 to
//	Least significant bit 0 are described below:
//
//	    Timezone offset has a range of -12:00 to +14:00 (ie -720 to +840 minutes).
//	    Bit 15 = have\_dst: set to 1 if we set the dst flag.
//	    Bit 14 = dst\_on: set to 1 if dst is in effect at the time, or 0 if not.
//	    Bits 13..0 = timezone offset in minutes. It is a signed integer in Big Endian format.
func customEncodeTime(t time.Time) []byte {
	// t := rv2i(rv).(time.Time)
	tsecs, tnsecs := t.Unix(), t.Nanosecond()
	var (
		bd byte
		bs [16]byte
		i  int = 1
	)
	l := t.Location()
	if l == time.UTC {
		l = nil
	}
	if tsecs != 0 {
		bd = bd | 0x80
		btmp := bigen.PutUint64(uint64(tsecs))
		f := pruneSignExt(btmp[:], tsecs >= 0)
		bd = bd | (byte(7-f) << 2)
		copy(bs[i:], btmp[f:])
		i = i + (8 - f)
	}
	if tnsecs != 0 {
		bd = bd | 0x40
		btmp := bigen.PutUint32(uint32(tnsecs))
		f := pruneSignExt(btmp[:4], true)
		bd = bd | byte(3-f)
		copy(bs[i:], btmp[f:4])
		i = i + (4 - f)
	}
	if l != nil {
		bd = bd | 0x20
		// Note that Go Libs do not give access to dst flag.
		_, zoneOffset := t.Zone()
		// zoneName, zoneOffset := t.Zone()
		zoneOffset /= 60
		z := uint16(zoneOffset)
		btmp0, btmp1 := bigen.PutUint16(z)
		// clear dst flags
		bs[i] = btmp0 & 0x3f
		bs[i+1] = btmp1
		i = i + 2
	}
	bs[0] = bd
	return bs[0:i]
}

// customDecodeTime decodes a []byte into a time.Time.
func customDecodeTime(bs []byte) (tt time.Time, err error) {
	bd := bs[0]
	var (
		tsec  int64
		tnsec uint32
		tz    uint16
		i     byte = 1
		i2    byte
		n     byte
	)
	if bd&(1<<7) != 0 {
		var btmp [8]byte
		n = ((bd >> 2) & 0x7) + 1
		i2 = i + n
		copy(btmp[8-n:], bs[i:i2])
		// if first bit of bs[i] is set, then fill btmp[0..8-n] with 0xff (ie sign extend it)
		if bs[i]&(1<<7) != 0 {
			copy(btmp[0:8-n], bsAll0xff)
		}
		i = i2
		tsec = int64(bigen.Uint64(btmp))
	}
	if bd&(1<<6) != 0 {
		var btmp [4]byte
		n = (bd & 0x3) + 1
		i2 = i + n
		copy(btmp[4-n:], bs[i:i2])
		i = i2
		tnsec = bigen.Uint32(btmp)
	}
	if bd&(1<<5) == 0 {
		tt = time.Unix(tsec, int64(tnsec)).UTC()
		return
	}
	// In stdlib time.Parse, when a date is parsed without a zone name, it uses "" as zone name.
	// However, we need name here, so it can be shown when time is printf.d.
	// Zone name is in form: UTC-08:00.
	// Note that Go Libs do not give access to dst flag, so we ignore dst bits

	tz = bigen.Uint16([2]byte{bs[i], bs[i+1]})
	// sign extend sign bit into top 2 MSB (which were dst bits):
	if tz&(1<<13) == 0 { // positive
		tz = tz & 0x3fff //clear 2 MSBs: dst bits
	} else { // negative
		tz = tz | 0xc000 //set 2 MSBs: dst bits
	}
	tzint := int16(tz)
	if tzint == 0 {
		tt = time.Unix(tsec, int64(tnsec)).UTC()
	} else {
		// For Go Time, do not use a descriptive timezone.
		// It's unnecessary, and makes it harder to do a reflect.DeepEqual.
		// The Offset already tells what the offset should be, if not on UTC and unknown zone name.
		// var zoneName = timeLocUTCName(tzint)
		tt = time.Unix(tsec, int64(tnsec)).In(time.FixedZone("", int(tzint)*60))
	}
	return
}

// customEncodeTimeAsNum encodes time.Time exactly as cbor does.
func customEncodeTimeAsNum(t time.Time) (r interface{}) {
	t = t.UTC().Round(time.Microsecond)
	sec, nsec := t.Unix(), uint64(t.Nanosecond())
	if nsec == 0 {
		r = sec
	} else {
		r = float64(sec) + float64(nsec)/1e9
	}
	return r
}

// customDecodeTimeAsNum decodes time.Time exactly as cbor does.
func customDecodeTimeAsNum(v interface{}) (t time.Time) {
	switch vv := v.(type) {
	case int64:
		t = time.Unix(vv, 0)
	case uint64:
		t = time.Unix((int64)(vv), 0)
	case float64:
		f1, f2 := math.Modf(vv)
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		halt.errorf("expect int64/float64 for time.Time ext: got %T", v)
	}
	t = t.UTC().Round(time.Microsecond)
	return
}
