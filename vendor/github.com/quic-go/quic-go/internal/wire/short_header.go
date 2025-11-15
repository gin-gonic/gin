package wire

import (
	"errors"
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
)

// ParseShortHeader parses a short header packet.
// It must be called after header protection was removed.
// Otherwise, the check for the reserved bits will (most likely) fail.
func ParseShortHeader(data []byte, connIDLen int) (length int, _ protocol.PacketNumber, _ protocol.PacketNumberLen, _ protocol.KeyPhaseBit, _ error) {
	if len(data) == 0 {
		return 0, 0, 0, 0, io.EOF
	}
	if data[0]&0x80 > 0 {
		return 0, 0, 0, 0, errors.New("not a short header packet")
	}
	if data[0]&0x40 == 0 {
		return 0, 0, 0, 0, errors.New("not a QUIC packet")
	}
	pnLen := protocol.PacketNumberLen(data[0]&0b11) + 1
	if len(data) < 1+int(pnLen)+connIDLen {
		return 0, 0, 0, 0, io.EOF
	}

	pos := 1 + connIDLen
	pn, err := readPacketNumber(data[pos:], pnLen)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	kp := protocol.KeyPhaseZero
	if data[0]&0b100 > 0 {
		kp = protocol.KeyPhaseOne
	}

	if data[0]&0x18 != 0 {
		err = ErrInvalidReservedBits
	}
	return 1 + connIDLen + int(pnLen), pn, pnLen, kp, err
}

// AppendShortHeader writes a short header.
func AppendShortHeader(b []byte, connID protocol.ConnectionID, pn protocol.PacketNumber, pnLen protocol.PacketNumberLen, kp protocol.KeyPhaseBit) ([]byte, error) {
	typeByte := 0x40 | uint8(pnLen-1)
	if kp == protocol.KeyPhaseOne {
		typeByte |= byte(1 << 2)
	}
	b = append(b, typeByte)
	b = append(b, connID.Bytes()...)
	return appendPacketNumber(b, pn, pnLen)
}

func ShortHeaderLen(dest protocol.ConnectionID, pnLen protocol.PacketNumberLen) protocol.ByteCount {
	return 1 + protocol.ByteCount(dest.Len()) + protocol.ByteCount(pnLen)
}

func LogShortHeader(logger utils.Logger, dest protocol.ConnectionID, pn protocol.PacketNumber, pnLen protocol.PacketNumberLen, kp protocol.KeyPhaseBit) {
	logger.Debugf("\tShort Header{DestConnectionID: %s, PacketNumber: %d, PacketNumberLen: %d, KeyPhase: %s}", dest, pn, pnLen, kp)
}
