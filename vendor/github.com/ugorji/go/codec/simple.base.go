// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
)

const (
	_               uint8 = iota
	simpleVdNil           = 1
	simpleVdFalse         = 2
	simpleVdTrue          = 3
	simpleVdFloat32       = 4
	simpleVdFloat64       = 5

	// each lasts for 4 (ie n, n+1, n+2, n+3)
	simpleVdPosInt = 8
	simpleVdNegInt = 12

	simpleVdTime = 24

	// containers: each lasts for 8 (ie n, n+1, n+2, ... n+7)
	simpleVdString    = 216
	simpleVdByteArray = 224
	simpleVdArray     = 232
	simpleVdMap       = 240
	simpleVdExt       = 248
)

var simpledescNames = map[byte]string{
	simpleVdNil:     "null",
	simpleVdFalse:   "false",
	simpleVdTrue:    "true",
	simpleVdFloat32: "float32",
	simpleVdFloat64: "float64",

	simpleVdPosInt: "+int",
	simpleVdNegInt: "-int",

	simpleVdTime: "time",

	simpleVdString:    "string",
	simpleVdByteArray: "binary",
	simpleVdArray:     "array",
	simpleVdMap:       "map",
	simpleVdExt:       "ext",
}

func simpledesc(bd byte) (s string) {
	s = simpledescNames[bd]
	if s == "" {
		s = "unknown"
	}
	return
}

//------------------------------------

// SimpleHandle is a Handle for a very simple encoding format.
//
// simple is a simplistic codec similar to binc, but not as compact.
//   - Encoding of a value is always preceded by the descriptor byte (bd)
//   - True, false, nil are encoded fully in 1 byte (the descriptor)
//   - Integers (intXXX, uintXXX) are encoded in 1, 2, 4 or 8 bytes (plus a descriptor byte).
//     There are positive (uintXXX and intXXX >= 0) and negative (intXXX < 0) integers.
//   - Floats are encoded in 4 or 8 bytes (plus a descriptor byte)
//   - Length of containers (strings, bytes, array, map, extensions)
//     are encoded in 0, 1, 2, 4 or 8 bytes.
//     Zero-length containers have no length encoded.
//     For others, the number of bytes is given by pow(2, bd%3)
//   - maps are encoded as [bd] [length] [[key][value]]...
//   - arrays are encoded as [bd] [length] [value]...
//   - extensions are encoded as [bd] [length] [tag] [byte]...
//   - strings/bytearrays are encoded as [bd] [length] [byte]...
//   - time.Time are encoded as [bd] [length] [byte]...
//
// The full spec will be published soon.
type SimpleHandle struct {
	binaryEncodingType
	notJsonType
	BasicHandle

	// EncZeroValuesAsNil says to encode zero values for numbers, bool, string, etc as nil
	EncZeroValuesAsNil bool
}

// Name returns the name of the handle: simple
func (h *SimpleHandle) Name() string { return "simple" }

func (h *SimpleHandle) desc(bd byte) string { return simpledesc(bd) }

// SetBytesExt sets an extension
func (h *SimpleHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}
