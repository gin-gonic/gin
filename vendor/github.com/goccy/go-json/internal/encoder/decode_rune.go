package encoder

import "unicode/utf8"

const (
	// The default lowest and highest continuation byte.
	locb = 128 //0b10000000
	hicb = 191 //0b10111111

	// These names of these constants are chosen to give nice alignment in the
	// table below. The first nibble is an index into acceptRanges or F for
	// special one-byte cases. The second nibble is the Rune length or the
	// Status for the special one-byte case.
	xx = 0xF1 // invalid: size 1
	as = 0xF0 // ASCII: size 1
	s1 = 0x02 // accept 0, size 2
	s2 = 0x13 // accept 1, size 3
	s3 = 0x03 // accept 0, size 3
	s4 = 0x23 // accept 2, size 3
	s5 = 0x34 // accept 3, size 4
	s6 = 0x04 // accept 0, size 4
	s7 = 0x44 // accept 4, size 4
)

// first is information about the first byte in a UTF-8 sequence.
var first = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x00-0x0F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x10-0x1F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x20-0x2F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x30-0x3F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x40-0x4F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x50-0x5F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x60-0x6F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x80-0x8F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x90-0x9F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xA0-0xAF
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xB0-0xBF
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xC0-0xCF
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xD0-0xDF
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, // 0xE0-0xEF
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xF0-0xFF
}

const (
	lineSep      = byte(168) //'\u2028'
	paragraphSep = byte(169) //'\u2029'
)

type decodeRuneState int

const (
	validUTF8State decodeRuneState = iota
	runeErrorState
	lineSepState
	paragraphSepState
)

func decodeRuneInString(s string) (decodeRuneState, int) {
	n := len(s)
	s0 := s[0]
	x := first[s0]
	if x >= as {
		// The following code simulates an additional check for x == xx and
		// handling the ASCII and invalid cases accordingly. This mask-and-or
		// approach prevents an additional branch.
		mask := rune(x) << 31 >> 31 // Create 0x0000 or 0xFFFF.
		if rune(s[0])&^mask|utf8.RuneError&mask == utf8.RuneError {
			return runeErrorState, 1
		}
		return validUTF8State, 1
	}
	sz := int(x & 7)
	if n < sz {
		return runeErrorState, 1
	}
	s1 := s[1]
	switch x >> 4 {
	case 0:
		if s1 < locb || hicb < s1 {
			return runeErrorState, 1
		}
	case 1:
		if s1 < 0xA0 || hicb < s1 {
			return runeErrorState, 1
		}
	case 2:
		if s1 < locb || 0x9F < s1 {
			return runeErrorState, 1
		}
	case 3:
		if s1 < 0x90 || hicb < s1 {
			return runeErrorState, 1
		}
	case 4:
		if s1 < locb || 0x8F < s1 {
			return runeErrorState, 1
		}
	}
	if sz <= 2 {
		return validUTF8State, 2
	}
	s2 := s[2]
	if s2 < locb || hicb < s2 {
		return runeErrorState, 1
	}
	if sz <= 3 {
		// separator character prefixes: [2]byte{226, 128}
		if s0 == 226 && s1 == 128 {
			switch s2 {
			case lineSep:
				return lineSepState, 3
			case paragraphSep:
				return paragraphSepState, 3
			}
		}
		return validUTF8State, 3
	}
	s3 := s[3]
	if s3 < locb || hicb < s3 {
		return runeErrorState, 1
	}
	return validUTF8State, 4
}
