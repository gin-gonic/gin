package quicvarint

import (
	"encoding/binary"
	"fmt"
	"io"
)

// taken from the QUIC draft
const (
	// Min is the minimum value allowed for a QUIC varint.
	Min = 0

	// Max is the maximum allowed value for a QUIC varint (2^62-1).
	Max = maxVarInt8

	maxVarInt1 = 63
	maxVarInt2 = 16383
	maxVarInt4 = 1073741823
	maxVarInt8 = 4611686018427387903
)

type varintLengthError struct {
	Num uint64
}

func (e *varintLengthError) Error() string {
	return fmt.Sprintf("value doesn't fit into 62 bits: %d", e.Num)
}

// Read reads a number in the QUIC varint format from r.
func Read(r io.ByteReader) (uint64, error) {
	firstByte, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	// the first two bits of the first byte encode the length
	l := 1 << ((firstByte & 0xc0) >> 6)
	b1 := firstByte & (0xff - 0xc0)
	if l == 1 {
		return uint64(b1), nil
	}
	b2, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	if l == 2 {
		return uint64(b2) + uint64(b1)<<8, nil
	}
	b3, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	b4, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	if l == 4 {
		return uint64(b4) + uint64(b3)<<8 + uint64(b2)<<16 + uint64(b1)<<24, nil
	}
	b5, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	b6, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	b7, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	b8, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint64(b8) + uint64(b7)<<8 + uint64(b6)<<16 + uint64(b5)<<24 + uint64(b4)<<32 + uint64(b3)<<40 + uint64(b2)<<48 + uint64(b1)<<56, nil
}

// Parse reads a number in the QUIC varint format.
// It returns the number of bytes consumed.
func Parse(b []byte) (uint64 /* value */, int /* bytes consumed */, error) {
	if len(b) == 0 {
		return 0, 0, io.EOF
	}

	first := b[0]
	switch first >> 6 {
	case 0: // 1-byte encoding: 00xxxxxx
		return uint64(first & 0b00111111), 1, nil
	case 1: // 2-byte encoding: 01xxxxxx
		if len(b) < 2 {
			return 0, 0, io.ErrUnexpectedEOF
		}
		return uint64(b[1]) | uint64(first&0b00111111)<<8, 2, nil
	case 2: // 4-byte encoding: 10xxxxxx
		if len(b) < 4 {
			return 0, 0, io.ErrUnexpectedEOF
		}
		return uint64(b[3]) | uint64(b[2])<<8 | uint64(b[1])<<16 | uint64(first&0b00111111)<<24, 4, nil
	case 3: // 8-byte encoding: 00xxxxxx
		if len(b) < 8 {
			return 0, 0, io.ErrUnexpectedEOF
		}
		// binary.BigEndian.Uint64 only reads the first 8 bytes. Passing the full slice avoids slicing overhead.
		return binary.BigEndian.Uint64(b) & 0x3fffffffffffffff, 8, nil
	}

	panic("unreachable")
}

// Append appends i in the QUIC varint format.
func Append(b []byte, i uint64) []byte {
	if i <= maxVarInt1 {
		return append(b, uint8(i))
	}
	if i <= maxVarInt2 {
		return append(b, []byte{uint8(i>>8) | 0x40, uint8(i)}...)
	}
	if i <= maxVarInt4 {
		return append(b, []byte{uint8(i>>24) | 0x80, uint8(i >> 16), uint8(i >> 8), uint8(i)}...)
	}
	if i <= maxVarInt8 {
		return append(b, []byte{
			uint8(i>>56) | 0xc0, uint8(i >> 48), uint8(i >> 40), uint8(i >> 32),
			uint8(i >> 24), uint8(i >> 16), uint8(i >> 8), uint8(i),
		}...)
	}
	panic(&varintLengthError{Num: i})
}

// AppendWithLen append i in the QUIC varint format with the desired length.
func AppendWithLen(b []byte, i uint64, length int) []byte {
	if length != 1 && length != 2 && length != 4 && length != 8 {
		panic("invalid varint length")
	}
	l := Len(i)
	if l == length {
		return Append(b, i)
	}
	if l > length {
		panic(fmt.Sprintf("cannot encode %d in %d bytes", i, length))
	}
	switch length {
	case 2:
		b = append(b, 0b01000000)
	case 4:
		b = append(b, 0b10000000)
	case 8:
		b = append(b, 0b11000000)
	}
	for range length - l - 1 {
		b = append(b, 0)
	}
	for j := range l {
		b = append(b, uint8(i>>(8*(l-1-j))))
	}
	return b
}

// Len determines the number of bytes that will be needed to write the number i.
//
//gcassert:inline
func Len(i uint64) int {
	if i <= maxVarInt1 {
		return 1
	}
	if i <= maxVarInt2 {
		return 2
	}
	if i <= maxVarInt4 {
		return 4
	}
	if i <= maxVarInt8 {
		return 8
	}
	// Don't use a fmt.Sprintf here to format the error message.
	// The function would then exceed the inlining budget.
	panic(&varintLengthError{Num: i})
}
