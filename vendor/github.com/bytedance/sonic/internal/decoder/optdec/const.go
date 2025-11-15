package optdec

import "math"

/*
Copied from sonic-rs
// JSON Value Type
const NULL: u64 = 0;
const BOOL: u64 = 2;
const FALSE: u64 = BOOL;
const TRUE: u64 = (1 << 3) | BOOL;
const NUMBER: u64 = 3;
const UINT: u64 = NUMBER;
const SINT: u64 = (1 << 3) | NUMBER;
const REAL: u64 = (2 << 3) | NUMBER;
const RAWNUMBER: u64 = (3 << 3) | NUMBER;
const STRING: u64 = 4;
const STRING_COMMON: u64 = STRING;
const STRING_HASESCAPED: u64 = (1 << 3) | STRING;
const OBJECT: u64 = 6;
const ARRAY: u64 = 7;

/// JSON Type Mask
const POS_MASK: u64 = (!0) << 32;
const POS_BITS: u64 = 32;
const TYPE_MASK: u64 = 0xFF;
const TYPE_BITS: u64 = 8;

*/

const (
	// BasicType: 3 bits
	KNull   = 0 // xxxxx000
	KBool   = 2 // xxxxx010
	KNumber = 3 // xxxxx011
	KString = 4 // xxxxx100
	KRaw    = 5 // xxxxx101
	KObject = 6 // xxxxx110
	KArray  = 7 // xxxxx111

	// SubType: 2 bits
	KFalse            = (0 << 3) | KBool   // xxx00_010, 2
	KTrue             = (1 << 3) | KBool   // xxx01_010, 10
	KUint             = (0 << 3) | KNumber // xxx00_011, 3
	KSint             = (1 << 3) | KNumber // xxx01_011, 11
	KReal             = (2 << 3) | KNumber // xxx10_011, 19
	KRawNumber        = (3 << 3) | KNumber // xxx11_011, 27
	KStringCommon     = KString            // xxx00_100, 4
	KStringEscaped = (1 << 3) | KString // xxx01_100, 12
)

const (
	PosMask  = math.MaxUint64 << 32
	PosBits  = 32
	TypeMask = 0xFF
	TypeBits = 8

	ConLenMask = uint64(math.MaxUint32)
	ConLenBits = 32
)
