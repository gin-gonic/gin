// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

// major
const (
	cborMajorUint byte = iota
	cborMajorNegInt
	cborMajorBytes
	cborMajorString
	cborMajorArray
	cborMajorMap
	cborMajorTag
	cborMajorSimpleOrFloat
)

// simple
const (
	cborBdFalse byte = 0xf4 + iota
	cborBdTrue
	cborBdNil
	cborBdUndefined
	cborBdExt
	cborBdFloat16
	cborBdFloat32
	cborBdFloat64
)

// indefinite
const (
	cborBdIndefiniteBytes  byte = 0x5f
	cborBdIndefiniteString byte = 0x7f
	cborBdIndefiniteArray  byte = 0x9f
	cborBdIndefiniteMap    byte = 0xbf
	cborBdBreak            byte = 0xff
)

// These define some in-stream descriptors for
// manual encoding e.g. when doing explicit indefinite-length
const (
	CborStreamBytes  byte = 0x5f
	CborStreamString byte = 0x7f
	CborStreamArray  byte = 0x9f
	CborStreamMap    byte = 0xbf
	CborStreamBreak  byte = 0xff
)

// base values
const (
	cborBaseUint   byte = 0x00
	cborBaseNegInt byte = 0x20
	cborBaseBytes  byte = 0x40
	cborBaseString byte = 0x60
	cborBaseArray  byte = 0x80
	cborBaseMap    byte = 0xa0
	cborBaseTag    byte = 0xc0
	cborBaseSimple byte = 0xe0
)

// const (
// 	cborSelfDesrTag  byte = 0xd9
// 	cborSelfDesrTag2 byte = 0xd9
// 	cborSelfDesrTag3 byte = 0xf7
// )

var (
	cbordescSimpleNames = map[byte]string{
		cborBdNil:     "nil",
		cborBdFalse:   "false",
		cborBdTrue:    "true",
		cborBdFloat16: "float",
		cborBdFloat32: "float",
		cborBdFloat64: "float",
		cborBdBreak:   "break",
	}
	cbordescIndefNames = map[byte]string{
		cborBdIndefiniteBytes:  "bytes*",
		cborBdIndefiniteString: "string*",
		cborBdIndefiniteArray:  "array*",
		cborBdIndefiniteMap:    "map*",
	}
	cbordescMajorNames = map[byte]string{
		cborMajorUint:          "(u)int",
		cborMajorNegInt:        "int",
		cborMajorBytes:         "bytes",
		cborMajorString:        "string",
		cborMajorArray:         "array",
		cborMajorMap:           "map",
		cborMajorTag:           "tag",
		cborMajorSimpleOrFloat: "simple",
	}
)

func cbordesc(bd byte) (s string) {
	bm := bd >> 5
	if bm == cborMajorSimpleOrFloat {
		s = cbordescSimpleNames[bd]
	} else {
		s = cbordescMajorNames[bm]
		if s == "" {
			s = cbordescIndefNames[bd]
		}
	}
	if s == "" {
		s = "unknown"
	}
	return
}

// -------------------------

// CborHandle is a Handle for the CBOR encoding format,
// defined at http://tools.ietf.org/html/rfc7049 and documented further at http://cbor.io .
//
// CBOR is comprehensively supported, including support for:
//   - indefinite-length arrays/maps/bytes/strings
//   - (extension) tags in range 0..0xffff (0 .. 65535)
//   - half, single and double-precision floats
//   - all numbers (1, 2, 4 and 8-byte signed and unsigned integers)
//   - nil, true, false, ...
//   - arrays and maps, bytes and text strings
//
// None of the optional extensions (with tags) defined in the spec are supported out-of-the-box.
// Users can implement them as needed (using SetExt), including spec-documented ones:
//   - timestamp, BigNum, BigFloat, Decimals,
//   - Encoded Text (e.g. URL, regexp, base64, MIME Message), etc.
type CborHandle struct {
	binaryEncodingType
	notJsonType
	// noElemSeparators
	BasicHandle

	// IndefiniteLength=true, means that we encode using indefinitelength
	IndefiniteLength bool

	// TimeRFC3339 says to encode time.Time using RFC3339 format.
	// If unset, we encode time.Time using seconds past epoch.
	TimeRFC3339 bool

	// SkipUnexpectedTags says to skip over any tags for which extensions are
	// not defined. This is in keeping with the cbor spec on "Optional Tagging of Items".
	//
	// Furthermore, this allows the skipping over of the Self Describing Tag 0xd9d9f7.
	SkipUnexpectedTags bool
}

// Name returns the name of the handle: cbor
func (h *CborHandle) Name() string { return "cbor" }

func (h *CborHandle) desc(bd byte) string { return cbordesc(bd) }

// SetInterfaceExt sets an extension
func (h *CborHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}
