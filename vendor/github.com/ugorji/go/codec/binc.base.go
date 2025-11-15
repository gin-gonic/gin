// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"time"
)

// Symbol management:
// - symbols are stored in a symbol map during encoding and decoding.
// - the symbols persist until the (En|De)coder ResetXXX method is called.

const bincDoPrune = true

// vd as low 4 bits (there are 16 slots)
const (
	bincVdSpecial byte = iota
	bincVdPosInt
	bincVdNegInt
	bincVdFloat

	bincVdString
	bincVdByteArray
	bincVdArray
	bincVdMap

	bincVdTimestamp
	bincVdSmallInt
	_ // bincVdUnicodeOther
	bincVdSymbol

	_               // bincVdDecimal
	_               // open slot
	_               // open slot
	bincVdCustomExt = 0x0f
)

const (
	bincSpNil byte = iota
	bincSpFalse
	bincSpTrue
	bincSpNan
	bincSpPosInf
	bincSpNegInf
	bincSpZeroFloat
	bincSpZero
	bincSpNegOne
)

const (
	_ byte = iota // bincFlBin16
	bincFlBin32
	_ // bincFlBin32e
	bincFlBin64
	_ // bincFlBin64e
	// others not currently supported
)

const bincBdNil = 0 // bincVdSpecial<<4 | bincSpNil // staticcheck barfs on this (SA4016)

var (
	bincdescSpecialVsNames = map[byte]string{
		bincSpNil:       "nil",
		bincSpFalse:     "false",
		bincSpTrue:      "true",
		bincSpNan:       "float",
		bincSpPosInf:    "float",
		bincSpNegInf:    "float",
		bincSpZeroFloat: "float",
		bincSpZero:      "uint",
		bincSpNegOne:    "int",
	}
	bincdescVdNames = map[byte]string{
		bincVdSpecial:   "special",
		bincVdSmallInt:  "uint",
		bincVdPosInt:    "uint",
		bincVdFloat:     "float",
		bincVdSymbol:    "string",
		bincVdString:    "string",
		bincVdByteArray: "bytes",
		bincVdTimestamp: "time",
		bincVdCustomExt: "ext",
		bincVdArray:     "array",
		bincVdMap:       "map",
	}
)

func bincdescbd(bd byte) (s string) {
	return bincdesc(bd>>4, bd&0x0f)
}

func bincdesc(vd, vs byte) (s string) {
	if vd == bincVdSpecial {
		s = bincdescSpecialVsNames[vs]
	} else {
		s = bincdescVdNames[vd]
	}
	if s == "" {
		s = "unknown"
	}
	return
}

type bincEncState struct {
	m map[string]uint16 // symbols
}

// func (e *bincEncState) restoreState(v interface{}) { e.m = v.(map[string]uint16) }
// func (e bincEncState) captureState() interface{}   { return e.m }
// func (e *bincEncState) resetState()                { e.m = nil }
// func (e *bincEncState) reset()                     { e.resetState() }
func (e *bincEncState) reset() { e.m = nil }

type bincDecState struct {
	bdRead bool
	bd     byte
	vd     byte
	vs     byte

	_ bool
	// MARKER: consider using binary search here instead of a map (ie bincDecSymbol)
	s map[uint16][]byte
}

// func (x bincDecState) captureState() interface{}   { return x }
// func (x *bincDecState) resetState()                { *x = bincDecState{} }
// func (x *bincDecState) reset()                     { x.resetState() }
// func (x *bincDecState) restoreState(v interface{}) { *x = v.(bincDecState) }
func (x *bincDecState) reset() { *x = bincDecState{} }

//------------------------------------

// BincHandle is a Handle for the Binc Schema-Free Encoding Format
// defined at https://github.com/ugorji/binc .
//
// BincHandle currently supports all Binc features with the following EXCEPTIONS:
//   - only integers up to 64 bits of precision are supported.
//     big integers are unsupported.
//   - Only IEEE 754 binary32 and binary64 floats are supported (ie Go float32 and float64 types).
//     extended precision and decimal IEEE 754 floats are unsupported.
//   - Only UTF-8 strings supported.
//     Unicode_Other Binc types (UTF16, UTF32) are currently unsupported.
//
// Note that these EXCEPTIONS are temporary and full support is possible and may happen soon.
type BincHandle struct {
	binaryEncodingType
	notJsonType
	// noElemSeparators
	BasicHandle

	// AsSymbols defines what should be encoded as symbols.
	//
	// Encoding as symbols can reduce the encoded size significantly.
	//
	// However, during decoding, each string to be encoded as a symbol must
	// be checked to see if it has been seen before. Consequently, encoding time
	// will increase if using symbols, because string comparisons has a clear cost.
	//
	// Values:
	// - 0: default: library uses best judgement
	// - 1: use symbols
	// - 2: do not use symbols
	AsSymbols uint8

	// AsSymbols: may later on introduce more options ...
	// - m: map keys
	// - s: struct fields
	// - n: none
	// - a: all: same as m, s, ...

	// _ [7]uint64 // padding (cache-aligned)
}

// Name returns the name of the handle: binc
func (h *BincHandle) Name() string { return "binc" }

func (h *BincHandle) desc(bd byte) string { return bincdesc(bd>>4, bd&0x0f) }

// SetBytesExt sets an extension
func (h *BincHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}

// var timeDigits = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func bincEncodeTime(t time.Time) []byte {
	return customEncodeTime(t)
}

func bincDecodeTime(bs []byte) (tt time.Time, err error) {
	return customDecodeTime(bs)
}
