package wire

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// ParseConnectionID parses the destination connection ID of a packet.
func ParseConnectionID(data []byte, shortHeaderConnIDLen int) (protocol.ConnectionID, error) {
	if len(data) == 0 {
		return protocol.ConnectionID{}, io.EOF
	}
	if !IsLongHeaderPacket(data[0]) {
		if len(data) < shortHeaderConnIDLen+1 {
			return protocol.ConnectionID{}, io.EOF
		}
		return protocol.ParseConnectionID(data[1 : 1+shortHeaderConnIDLen]), nil
	}
	if len(data) < 6 {
		return protocol.ConnectionID{}, io.EOF
	}
	destConnIDLen := int(data[5])
	if destConnIDLen > protocol.MaxConnIDLen {
		return protocol.ConnectionID{}, protocol.ErrInvalidConnectionIDLen
	}
	if len(data) < 6+destConnIDLen {
		return protocol.ConnectionID{}, io.EOF
	}
	return protocol.ParseConnectionID(data[6 : 6+destConnIDLen]), nil
}

// ParseArbitraryLenConnectionIDs parses the most general form of a Long Header packet,
// using only the version-independent packet format as described in Section 5.1 of RFC 8999:
// https://datatracker.ietf.org/doc/html/rfc8999#section-5.1.
// This function should only be called on Long Header packets for which we don't support the version.
func ParseArbitraryLenConnectionIDs(data []byte) (bytesParsed int, dest, src protocol.ArbitraryLenConnectionID, _ error) {
	startLen := len(data)
	if len(data) < 6 {
		return 0, nil, nil, io.EOF
	}
	data = data[5:] // skip first byte and version field
	destConnIDLen := data[0]
	data = data[1:]
	destConnID := make(protocol.ArbitraryLenConnectionID, destConnIDLen)
	if len(data) < int(destConnIDLen)+1 {
		return 0, nil, nil, io.EOF
	}
	copy(destConnID, data)
	data = data[destConnIDLen:]
	srcConnIDLen := data[0]
	data = data[1:]
	if len(data) < int(srcConnIDLen) {
		return 0, nil, nil, io.EOF
	}
	srcConnID := make(protocol.ArbitraryLenConnectionID, srcConnIDLen)
	copy(srcConnID, data)
	return startLen - len(data) + int(srcConnIDLen), destConnID, srcConnID, nil
}

func IsPotentialQUICPacket(firstByte byte) bool {
	return firstByte&0x40 > 0
}

// IsLongHeaderPacket says if this is a Long Header packet
func IsLongHeaderPacket(firstByte byte) bool {
	return firstByte&0x80 > 0
}

// ParseVersion parses the QUIC version.
// It should only be called for Long Header packets (Short Header packets don't contain a version number).
func ParseVersion(data []byte) (protocol.Version, error) {
	if len(data) < 5 {
		return 0, io.EOF
	}
	return protocol.Version(binary.BigEndian.Uint32(data[1:5])), nil
}

// IsVersionNegotiationPacket says if this is a version negotiation packet
func IsVersionNegotiationPacket(b []byte) bool {
	if len(b) < 5 {
		return false
	}
	return IsLongHeaderPacket(b[0]) && b[1] == 0 && b[2] == 0 && b[3] == 0 && b[4] == 0
}

// Is0RTTPacket says if this is a 0-RTT packet.
// A packet sent with a version we don't understand can never be a 0-RTT packet.
func Is0RTTPacket(b []byte) bool {
	if len(b) < 5 {
		return false
	}
	if !IsLongHeaderPacket(b[0]) {
		return false
	}
	version := protocol.Version(binary.BigEndian.Uint32(b[1:5]))
	//nolint:exhaustive // We only need to test QUIC versions that we support.
	switch version {
	case protocol.Version1:
		return b[0]>>4&0b11 == 0b01
	case protocol.Version2:
		return b[0]>>4&0b11 == 0b10
	default:
		return false
	}
}

var ErrUnsupportedVersion = errors.New("unsupported version")

// The Header is the version independent part of the header
type Header struct {
	typeByte byte
	Type     protocol.PacketType

	Version          protocol.Version
	SrcConnectionID  protocol.ConnectionID
	DestConnectionID protocol.ConnectionID

	Length protocol.ByteCount

	Token []byte

	parsedLen protocol.ByteCount // how many bytes were read while parsing this header
}

// ParsePacket parses a long header packet.
// The packet is cut according to the length field.
// If we understand the version, the packet is parsed up unto the packet number.
// Otherwise, only the invariant part of the header is parsed.
func ParsePacket(data []byte) (*Header, []byte, []byte, error) {
	if len(data) == 0 || !IsLongHeaderPacket(data[0]) {
		return nil, nil, nil, errors.New("not a long header packet")
	}
	hdr, err := parseHeader(data)
	if err != nil {
		if errors.Is(err, ErrUnsupportedVersion) {
			return hdr, nil, nil, err
		}
		return nil, nil, nil, err
	}
	if protocol.ByteCount(len(data)) < hdr.ParsedLen()+hdr.Length {
		return nil, nil, nil, fmt.Errorf("packet length (%d bytes) is smaller than the expected length (%d bytes)", len(data)-int(hdr.ParsedLen()), hdr.Length)
	}
	packetLen := int(hdr.ParsedLen() + hdr.Length)
	return hdr, data[:packetLen], data[packetLen:], nil
}

// ParseHeader parses the header:
// * if we understand the version: up to the packet number
// * if not, only the invariant part of the header
func parseHeader(b []byte) (*Header, error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	typeByte := b[0]

	h := &Header{typeByte: typeByte}
	l, err := h.parseLongHeader(b[1:])
	h.parsedLen = protocol.ByteCount(l) + 1
	return h, err
}

func (h *Header) parseLongHeader(b []byte) (int, error) {
	startLen := len(b)
	if len(b) < 5 {
		return 0, io.EOF
	}
	h.Version = protocol.Version(binary.BigEndian.Uint32(b[:4]))
	if h.Version != 0 && h.typeByte&0x40 == 0 {
		return startLen - len(b), errors.New("not a QUIC packet")
	}
	destConnIDLen := int(b[4])
	if destConnIDLen > protocol.MaxConnIDLen {
		return startLen - len(b), protocol.ErrInvalidConnectionIDLen
	}
	b = b[5:]
	if len(b) < destConnIDLen+1 {
		return startLen - len(b), io.EOF
	}
	h.DestConnectionID = protocol.ParseConnectionID(b[:destConnIDLen])
	srcConnIDLen := int(b[destConnIDLen])
	if srcConnIDLen > protocol.MaxConnIDLen {
		return startLen - len(b), protocol.ErrInvalidConnectionIDLen
	}
	b = b[destConnIDLen+1:]
	if len(b) < srcConnIDLen {
		return startLen - len(b), io.EOF
	}
	h.SrcConnectionID = protocol.ParseConnectionID(b[:srcConnIDLen])
	b = b[srcConnIDLen:]
	if h.Version == 0 { // version negotiation packet
		return startLen - len(b), nil
	}
	// If we don't understand the version, we have no idea how to interpret the rest of the bytes
	if !protocol.IsSupportedVersion(protocol.SupportedVersions, h.Version) {
		return startLen - len(b), ErrUnsupportedVersion
	}

	if h.Version == protocol.Version2 {
		switch h.typeByte >> 4 & 0b11 {
		case 0b00:
			h.Type = protocol.PacketTypeRetry
		case 0b01:
			h.Type = protocol.PacketTypeInitial
		case 0b10:
			h.Type = protocol.PacketType0RTT
		case 0b11:
			h.Type = protocol.PacketTypeHandshake
		}
	} else {
		switch h.typeByte >> 4 & 0b11 {
		case 0b00:
			h.Type = protocol.PacketTypeInitial
		case 0b01:
			h.Type = protocol.PacketType0RTT
		case 0b10:
			h.Type = protocol.PacketTypeHandshake
		case 0b11:
			h.Type = protocol.PacketTypeRetry
		}
	}

	if h.Type == protocol.PacketTypeRetry {
		tokenLen := len(b) - 16
		if tokenLen <= 0 {
			return startLen - len(b), io.EOF
		}
		h.Token = make([]byte, tokenLen)
		copy(h.Token, b[:tokenLen])
		return startLen - len(b) + tokenLen + 16, nil
	}

	if h.Type == protocol.PacketTypeInitial {
		tokenLen, n, err := quicvarint.Parse(b)
		if err != nil {
			return startLen - len(b), err
		}
		b = b[n:]
		if tokenLen > uint64(len(b)) {
			return startLen - len(b), io.EOF
		}
		h.Token = make([]byte, tokenLen)
		copy(h.Token, b[:tokenLen])
		b = b[tokenLen:]
	}

	pl, n, err := quicvarint.Parse(b)
	if err != nil {
		return 0, err
	}
	h.Length = protocol.ByteCount(pl)
	return startLen - len(b) + n, nil
}

// ParsedLen returns the number of bytes that were consumed when parsing the header
func (h *Header) ParsedLen() protocol.ByteCount {
	return h.parsedLen
}

// ParseExtended parses the version dependent part of the header.
// The Reader has to be set such that it points to the first byte of the header.
func (h *Header) ParseExtended(data []byte) (*ExtendedHeader, error) {
	extHdr := h.toExtendedHeader()
	reservedBitsValid, err := extHdr.parse(data)
	if err != nil {
		return nil, err
	}
	if !reservedBitsValid {
		return extHdr, ErrInvalidReservedBits
	}
	return extHdr, nil
}

func (h *Header) toExtendedHeader() *ExtendedHeader {
	return &ExtendedHeader{Header: *h}
}

// PacketType is the type of the packet, for logging purposes
func (h *Header) PacketType() string {
	return h.Type.String()
}

func readPacketNumber(data []byte, pnLen protocol.PacketNumberLen) (protocol.PacketNumber, error) {
	var pn protocol.PacketNumber
	switch pnLen {
	case protocol.PacketNumberLen1:
		pn = protocol.PacketNumber(data[0])
	case protocol.PacketNumberLen2:
		pn = protocol.PacketNumber(binary.BigEndian.Uint16(data[:2]))
	case protocol.PacketNumberLen3:
		pn = protocol.PacketNumber(uint32(data[2]) + uint32(data[1])<<8 + uint32(data[0])<<16)
	case protocol.PacketNumberLen4:
		pn = protocol.PacketNumber(binary.BigEndian.Uint32(data[:4]))
	default:
		return 0, fmt.Errorf("invalid packet number length: %d", pnLen)
	}
	return pn, nil
}
