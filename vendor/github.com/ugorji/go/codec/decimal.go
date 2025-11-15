// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"strconv"
)

type readFloatResult struct {
	mantissa uint64
	exp      int8
	neg      bool
	trunc    bool
	bad      bool // bad decimal string
	hardexp  bool // exponent is hard to handle (> 2 digits, etc)
	ok       bool
	// sawdot   bool
	// sawexp   bool
	//_ [2]bool // padding
}

// Per go spec, floats are represented in memory as
// IEEE single or double precision floating point values.
//
// We also looked at the source for stdlib math/modf.go,
// reviewed https://github.com/chewxy/math32
// and read wikipedia documents describing the formats.
//
// It became clear that we could easily look at the bits to determine
// whether any fraction exists.

func parseFloat32(b []byte) (f float32, err error) {
	return parseFloat32_custom(b)
}

func parseFloat64(b []byte) (f float64, err error) {
	return parseFloat64_custom(b)
}

func parseFloat32_strconv(b []byte) (f float32, err error) {
	f64, err := strconv.ParseFloat(stringView(b), 32)
	f = float32(f64)
	return
}

func parseFloat64_strconv(b []byte) (f float64, err error) {
	return strconv.ParseFloat(stringView(b), 64)
}

// ------ parseFloat custom below --------

// JSON really supports decimal numbers in base 10 notation, with exponent support.
//
// We assume the following:
//   - a lot of floating point numbers in json files will have defined precision
//     (in terms of number of digits after decimal point), etc.
//   - these (referenced above) can be written in exact format.
//
// strconv.ParseFloat has some unnecessary overhead which we can do without
// for the common case:
//
//    - expensive char-by-char check to see if underscores are in right place
//    - testing for and skipping underscores
//    - check if the string matches ignorecase +/- inf, +/- infinity, nan
//    - support for base 16 (0xFFFF...)
//
// The functions below will try a fast-path for floats which can be decoded
// without any loss of precision, meaning they:
//
//    - fits within the significand bits of the 32-bits or 64-bits
//    - exponent fits within the exponent value
//    - there is no truncation (any extra numbers are all trailing zeros)
//
// To figure out what the values are for maxMantDigits, use this idea below:
//
// 2^23 =                 838 8608 (between 10^ 6 and 10^ 7) (significand bits of uint32)
// 2^32 =             42 9496 7296 (between 10^ 9 and 10^10) (full uint32)
// 2^52 =      4503 5996 2737 0496 (between 10^15 and 10^16) (significand bits of uint64)
// 2^64 = 1844 6744 0737 0955 1616 (between 10^19 and 10^20) (full uint64)
//
// Note: we only allow for up to what can comfortably fit into the significand
// ignoring the exponent, and we only try to parse iff significand fits.

const (
	fMaxMultiplierForExactPow10_64 = 1e15
	fMaxMultiplierForExactPow10_32 = 1e7

	fUint64Cutoff = (1<<64-1)/10 + 1
	// fUint32Cutoff = (1<<32-1)/10 + 1

	fBase = 10
)

const (
	thousand    = 1000
	million     = thousand * thousand
	billion     = thousand * million
	trillion    = thousand * billion
	quadrillion = thousand * trillion
	quintillion = thousand * quadrillion
)

// Exact powers of 10.
var uint64pow10 = [...]uint64{
	1, 10, 100,
	1 * thousand, 10 * thousand, 100 * thousand,
	1 * million, 10 * million, 100 * million,
	1 * billion, 10 * billion, 100 * billion,
	1 * trillion, 10 * trillion, 100 * trillion,
	1 * quadrillion, 10 * quadrillion, 100 * quadrillion,
	1 * quintillion, 10 * quintillion,
}
var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
var float32pow10 = [...]float32{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10,
}

type floatinfo struct {
	mantbits uint8

	// expbits uint8 // (unused)
	// bias    int16 // (unused)
	// is32bit bool // (unused)

	exactPow10 int8 // Exact powers of ten are <= 10^N (32: 10, 64: 22)

	exactInts int8 // Exact integers are <= 10^N (for non-float, set to 0)

	// maxMantDigits int8 // 10^19 fits in uint64, while 10^9 fits in uint32

	mantCutoffIsUint64Cutoff bool

	mantCutoff uint64
}

var fi32 = floatinfo{23, 10, 7, false, 1<<23 - 1}
var fi64 = floatinfo{52, 22, 15, false, 1<<52 - 1}

var fi64u = floatinfo{0, 19, 0, true, fUint64Cutoff}

func noFrac64(fbits uint64) bool {
	if fbits == 0 {
		return true
	}

	exp := uint64(fbits>>52)&0x7FF - 1023 // uint(x>>shift)&mask - bias
	// clear top 12+e bits, the integer part; if the rest is 0, then no fraction.
	return exp < 52 && fbits<<(12+exp) == 0 // means there's no fractional part
}

func noFrac32(fbits uint32) bool {
	if fbits == 0 {
		return true
	}

	exp := uint32(fbits>>23)&0xFF - 127 // uint(x>>shift)&mask - bias
	// clear top 9+e bits, the integer part; if the rest is 0, then no fraction.
	return exp < 23 && fbits<<(9+exp) == 0 // means there's no fractional part
}

func strconvParseErr(b []byte, fn string) error {
	return &strconv.NumError{
		Func: fn,
		Err:  strconv.ErrSyntax,
		Num:  string(b),
	}
}

func parseFloat32_reader(r readFloatResult) (f float32, fail bool) {
	f = float32(r.mantissa)
	if r.exp == 0 {
	} else if r.exp < 0 { // int / 10^k
		f /= float32pow10[uint8(-r.exp)]
	} else { // exp > 0
		if r.exp > fi32.exactPow10 {
			f *= float32pow10[r.exp-fi32.exactPow10]
			if f > fMaxMultiplierForExactPow10_32 { // exponent too large - outside range
				fail = true
				return // ok = false
			}
			f *= float32pow10[fi32.exactPow10]
		} else {
			f *= float32pow10[uint8(r.exp)]
		}
	}
	if r.neg {
		f = -f
	}
	return
}

func parseFloat32_custom(b []byte) (f float32, err error) {
	r := readFloat(b, fi32)
	if r.bad {
		return 0, strconvParseErr(b, "ParseFloat")
	}
	if r.ok {
		f, r.bad = parseFloat32_reader(r)
		if !r.bad {
			return
		}
	}
	return parseFloat32_strconv(b)
}

func parseFloat64_reader(r readFloatResult) (f float64, fail bool) {
	f = float64(r.mantissa)
	if r.exp == 0 {
	} else if r.exp < 0 { // int / 10^k
		f /= float64pow10[-uint8(r.exp)]
	} else { // exp > 0
		if r.exp > fi64.exactPow10 {
			f *= float64pow10[r.exp-fi64.exactPow10]
			if f > fMaxMultiplierForExactPow10_64 { // exponent too large - outside range
				fail = true
				return
			}
			f *= float64pow10[fi64.exactPow10]
		} else {
			f *= float64pow10[uint8(r.exp)]
		}
	}
	if r.neg {
		f = -f
	}
	return
}

func parseFloat64_custom(b []byte) (f float64, err error) {
	r := readFloat(b, fi64)
	if r.bad {
		return 0, strconvParseErr(b, "ParseFloat")
	}
	if r.ok {
		f, r.bad = parseFloat64_reader(r)
		if !r.bad {
			return
		}
	}
	return parseFloat64_strconv(b)
}

func parseUint64_simple(b []byte) (n uint64, ok bool) {
	if len(b) > 1 && b[0] == '0' { // punt on numbers with leading zeros
		return
	}

	var i int
	var n1 uint64
	var c uint8
LOOP:
	if i < len(b) {
		c = b[i]
		// unsigned integers don't overflow well on multiplication, so check cutoff here
		// e.g. (maxUint64-5)*10 doesn't overflow well ...
		// if n >= fUint64Cutoff || !isDigitChar(b[i]) { // if c < '0' || c > '9' {
		if n >= fUint64Cutoff || c < '0' || c > '9' {
			return
		} else if c == '0' {
			n *= fBase
		} else {
			n1 = n
			n = n*fBase + uint64(c-'0')
			if n < n1 {
				return
			}
		}
		i++
		goto LOOP
	}
	ok = true
	return
}

func parseUint64_reader(r readFloatResult) (f uint64, fail bool) {
	f = r.mantissa
	if r.exp == 0 {
	} else if r.exp < 0 { // int / 10^k
		if f%uint64pow10[uint8(-r.exp)] != 0 {
			fail = true
		} else {
			f /= uint64pow10[uint8(-r.exp)]
		}
	} else { // exp > 0
		f *= uint64pow10[uint8(r.exp)]
	}
	return
}

func parseInteger_bytes(b []byte) (u uint64, neg, ok bool) {
	if len(b) == 0 {
		ok = true
		return
	}
	if b[0] == '-' {
		if len(b) == 1 {
			return
		}
		neg = true
		b = b[1:]
	}

	u, ok = parseUint64_simple(b)
	if ok {
		return
	}

	r := readFloat(b, fi64u)
	if r.ok {
		var fail bool
		u, fail = parseUint64_reader(r)
		if fail {
			f, err := parseFloat64(b)
			if err != nil {
				return
			}
			if !noFrac64(math.Float64bits(f)) {
				return
			}
			u = uint64(f)
		}
		ok = true
		return
	}
	return
}

// parseNumber will return an integer if only composed of [-]?[0-9]+
// Else it will return a float.
func parseNumber(b []byte, z *fauxUnion, preferSignedInt bool) (err error) {
	var ok, neg bool
	var f uint64

	if len(b) == 0 {
		return
	}

	if b[0] == '-' {
		neg = true
		f, ok = parseUint64_simple(b[1:])
	} else {
		f, ok = parseUint64_simple(b)
	}

	if ok {
		if neg {
			z.v = valueTypeInt
			if chkOvf.Uint2Int(f, neg) {
				return strconvParseErr(b, "ParseInt")
			}
			z.i = -int64(f)
		} else if preferSignedInt {
			z.v = valueTypeInt
			if chkOvf.Uint2Int(f, neg) {
				return strconvParseErr(b, "ParseInt")
			}
			z.i = int64(f)
		} else {
			z.v = valueTypeUint
			z.u = f
		}
		return
	}

	z.v = valueTypeFloat
	z.f, err = parseFloat64_custom(b)
	return
}

func readFloat(s []byte, y floatinfo) (r readFloatResult) {
	var i uint // uint, so that we eliminate bounds checking
	var slen = uint(len(s))
	if slen == 0 {
		// read an empty string as the zero value
		// r.bad = true
		r.ok = true
		return
	}

	if s[0] == '-' {
		r.neg = true
		i++
	}

	// considered punting early if string has length > maxMantDigits, but doesn't account
	// for trailing 0's e.g. 700000000000000000000 can be encoded exactly as it is 7e20

	var nd, ndMant, dp int8
	var sawdot, sawexp bool
	var xu uint64

	if i+1 < slen && s[i] == '0' {
		switch s[i+1] {
		case '.', 'e', 'E':
			// ok
		default:
			r.bad = true
			return
		}
	}

LOOP:
	for ; i < slen; i++ {
		switch s[i] {
		case '.':
			if sawdot {
				r.bad = true
				return
			}
			sawdot = true
			dp = nd
		case 'e', 'E':
			sawexp = true
			break LOOP
		case '0':
			if nd == 0 {
				dp--
				continue LOOP
			}
			nd++
			if r.mantissa < y.mantCutoff {
				r.mantissa *= fBase
				ndMant++
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			nd++
			if y.mantCutoffIsUint64Cutoff && r.mantissa < fUint64Cutoff {
				r.mantissa *= fBase
				xu = r.mantissa + uint64(s[i]-'0')
				if xu < r.mantissa {
					r.trunc = true
					return
				}
				r.mantissa = xu
			} else if r.mantissa < y.mantCutoff {
				// mantissa = (mantissa << 1) + (mantissa << 3) + uint64(c-'0')
				r.mantissa = r.mantissa*fBase + uint64(s[i]-'0')
			} else {
				r.trunc = true
				return
			}
			ndMant++
		default:
			r.bad = true
			return
		}
	}

	if !sawdot {
		dp = nd
	}

	if sawexp {
		i++
		if i < slen {
			var eneg bool
			if s[i] == '+' {
				i++
			} else if s[i] == '-' {
				i++
				eneg = true
			}
			if i < slen {
				// for exact match, exponent is 1 or 2 digits (float64: -22 to 37, float32: -1 to 17).
				// exit quick if exponent is more than 2 digits.
				if i+2 < slen {
					r.hardexp = true
					return
				}
				var e int8
				if s[i] < '0' || s[i] > '9' { // !isDigitChar(s[i]) { //
					r.bad = true
					return
				}
				e = int8(s[i] - '0')
				i++
				if i < slen {
					if s[i] < '0' || s[i] > '9' { // !isDigitChar(s[i]) { //
						r.bad = true
						return
					}
					e = e*fBase + int8(s[i]-'0') // (e << 1) + (e << 3) + int8(s[i]-'0')
					i++
				}
				if eneg {
					dp -= e
				} else {
					dp += e
				}
			}
		}
	}

	if r.mantissa != 0 {
		r.exp = dp - ndMant
		// do not set ok=true for cases we cannot handle
		if r.exp < -y.exactPow10 ||
			r.exp > y.exactInts+y.exactPow10 ||
			(y.mantbits != 0 && r.mantissa>>y.mantbits != 0) {
			r.hardexp = true
			return
		}
	}

	_ = i // no-op
	r.ok = true
	return
}
