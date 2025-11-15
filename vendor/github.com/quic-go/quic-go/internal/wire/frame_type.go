package wire

import "github.com/quic-go/quic-go/internal/protocol"

type FrameType uint64

// These constants correspond to those defined in RFC 9000.
// Stream frame types are not listed explicitly here; use FrameType.IsStreamFrameType() to identify them.
const (
	FrameTypePing        FrameType = 0x1
	FrameTypeAck         FrameType = 0x2
	FrameTypeAckECN      FrameType = 0x3
	FrameTypeResetStream FrameType = 0x4
	FrameTypeStopSending FrameType = 0x5
	FrameTypeCrypto      FrameType = 0x6
	FrameTypeNewToken    FrameType = 0x7

	FrameTypeMaxData            FrameType = 0x10
	FrameTypeMaxStreamData      FrameType = 0x11
	FrameTypeBidiMaxStreams     FrameType = 0x12
	FrameTypeUniMaxStreams      FrameType = 0x13
	FrameTypeDataBlocked        FrameType = 0x14
	FrameTypeStreamDataBlocked  FrameType = 0x15
	FrameTypeBidiStreamBlocked  FrameType = 0x16
	FrameTypeUniStreamBlocked   FrameType = 0x17
	FrameTypeNewConnectionID    FrameType = 0x18
	FrameTypeRetireConnectionID FrameType = 0x19
	FrameTypePathChallenge      FrameType = 0x1a
	FrameTypePathResponse       FrameType = 0x1b
	FrameTypeConnectionClose    FrameType = 0x1c
	FrameTypeApplicationClose   FrameType = 0x1d
	FrameTypeHandshakeDone      FrameType = 0x1e
	// https://datatracker.ietf.org/doc/draft-ietf-quic-reliable-stream-reset/07/
	FrameTypeResetStreamAt FrameType = 0x24
	// https://datatracker.ietf.org/doc/draft-ietf-quic-ack-frequency/11/
	FrameTypeAckFrequency FrameType = 0xaf
	FrameTypeImmediateAck FrameType = 0x1f

	FrameTypeDatagramNoLength   FrameType = 0x30
	FrameTypeDatagramWithLength FrameType = 0x31
)

func (t FrameType) IsStreamFrameType() bool {
	return t >= 0x8 && t <= 0xf
}

func (t FrameType) isValidRFC9000() bool {
	return t <= 0x1e
}

func (t FrameType) IsAckFrameType() bool {
	return t == FrameTypeAck || t == FrameTypeAckECN
}

func (t FrameType) IsDatagramFrameType() bool {
	return t == FrameTypeDatagramNoLength || t == FrameTypeDatagramWithLength
}

func (t FrameType) isAllowedAtEncLevel(encLevel protocol.EncryptionLevel) bool {
	//nolint:exhaustive
	switch encLevel {
	case protocol.EncryptionInitial, protocol.EncryptionHandshake:
		switch t {
		case FrameTypeCrypto, FrameTypeAck, FrameTypeAckECN, FrameTypeConnectionClose, FrameTypePing:
			return true
		default:
			return false
		}
	case protocol.Encryption0RTT:
		switch t {
		case FrameTypeCrypto, FrameTypeAck, FrameTypeAckECN, FrameTypeConnectionClose, FrameTypeNewToken, FrameTypePathResponse, FrameTypeRetireConnectionID:
			return false
		default:
			return true
		}
	case protocol.Encryption1RTT:
		return true
	default:
		panic("unknown encryption level")
	}
}
