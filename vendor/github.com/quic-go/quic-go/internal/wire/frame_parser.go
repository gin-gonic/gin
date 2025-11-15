package wire

import (
	"errors"
	"fmt"
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/quicvarint"
)

var errUnknownFrameType = errors.New("unknown frame type")

// The FrameParser parses QUIC frames, one by one.
type FrameParser struct {
	ackDelayExponent      uint8
	supportsDatagrams     bool
	supportsResetStreamAt bool
	supportsAckFrequency  bool

	// To avoid allocating when parsing, keep a single ACK frame struct.
	// It is used over and over again.
	ackFrame *AckFrame
}

// NewFrameParser creates a new frame parser.
func NewFrameParser(supportsDatagrams, supportsResetStreamAt, supportsAckFrequency bool) *FrameParser {
	return &FrameParser{
		supportsDatagrams:     supportsDatagrams,
		supportsResetStreamAt: supportsResetStreamAt,
		supportsAckFrequency:  supportsAckFrequency,
		ackFrame:              &AckFrame{},
	}
}

// ParseType parses the frame type of the next frame.
// It skips over PADDING frames.
func (p *FrameParser) ParseType(b []byte, encLevel protocol.EncryptionLevel) (FrameType, int, error) {
	var parsed int
	for len(b) != 0 {
		typ, l, err := quicvarint.Parse(b)
		parsed += l
		if err != nil {
			return 0, parsed, &qerr.TransportError{
				ErrorCode:    qerr.FrameEncodingError,
				ErrorMessage: err.Error(),
			}
		}
		b = b[l:]
		if typ == 0x0 { // skip PADDING frames
			continue
		}
		ft := FrameType(typ)
		valid := ft.isValidRFC9000() ||
			(p.supportsDatagrams && ft.IsDatagramFrameType()) ||
			(p.supportsResetStreamAt && ft == FrameTypeResetStreamAt) ||
			(p.supportsAckFrequency && (ft == FrameTypeAckFrequency || ft == FrameTypeImmediateAck))
		if !valid {
			return 0, parsed, &qerr.TransportError{
				ErrorCode:    qerr.FrameEncodingError,
				FrameType:    typ,
				ErrorMessage: errUnknownFrameType.Error(),
			}
		}
		if !ft.isAllowedAtEncLevel(encLevel) {
			return 0, parsed, &qerr.TransportError{
				ErrorCode:    qerr.FrameEncodingError,
				FrameType:    typ,
				ErrorMessage: fmt.Sprintf("%d not allowed at encryption level %s", ft, encLevel),
			}
		}
		return ft, parsed, nil
	}
	return 0, parsed, io.EOF
}

func (p *FrameParser) ParseStreamFrame(frameType FrameType, data []byte, v protocol.Version) (*StreamFrame, int, error) {
	frame, n, err := ParseStreamFrame(data, frameType, v)
	if err != nil {
		return nil, n, &qerr.TransportError{
			ErrorCode:    qerr.FrameEncodingError,
			FrameType:    uint64(frameType),
			ErrorMessage: err.Error(),
		}
	}
	return frame, n, nil
}

func (p *FrameParser) ParseAckFrame(frameType FrameType, data []byte, encLevel protocol.EncryptionLevel, v protocol.Version) (*AckFrame, int, error) {
	ackDelayExponent := p.ackDelayExponent
	if encLevel != protocol.Encryption1RTT {
		ackDelayExponent = protocol.DefaultAckDelayExponent
	}
	p.ackFrame.Reset()
	l, err := parseAckFrame(p.ackFrame, data, frameType, ackDelayExponent, v)
	if err != nil {
		return nil, l, &qerr.TransportError{
			ErrorCode:    qerr.FrameEncodingError,
			FrameType:    uint64(frameType),
			ErrorMessage: err.Error(),
		}
	}

	return p.ackFrame, l, nil
}

func (p *FrameParser) ParseDatagramFrame(frameType FrameType, data []byte, v protocol.Version) (*DatagramFrame, int, error) {
	f, l, err := parseDatagramFrame(data, frameType, v)
	if err != nil {
		return nil, 0, &qerr.TransportError{
			ErrorCode:    qerr.FrameEncodingError,
			FrameType:    uint64(frameType),
			ErrorMessage: err.Error(),
		}
	}
	return f, l, nil
}

// ParseLessCommonFrame parses everything except STREAM, ACK or DATAGRAM.
// These cases should be handled separately for performance reasons.
func (p *FrameParser) ParseLessCommonFrame(frameType FrameType, data []byte, v protocol.Version) (Frame, int, error) {
	var frame Frame
	var l int
	var err error
	//nolint:exhaustive // Common frames should already be handled.
	switch frameType {
	case FrameTypePing:
		frame = &PingFrame{}
	case FrameTypeResetStream:
		frame, l, err = parseResetStreamFrame(data, false, v)
	case FrameTypeStopSending:
		frame, l, err = parseStopSendingFrame(data, v)
	case FrameTypeCrypto:
		frame, l, err = parseCryptoFrame(data, v)
	case FrameTypeNewToken:
		frame, l, err = parseNewTokenFrame(data, v)
	case FrameTypeMaxData:
		frame, l, err = parseMaxDataFrame(data, v)
	case FrameTypeMaxStreamData:
		frame, l, err = parseMaxStreamDataFrame(data, v)
	case FrameTypeBidiMaxStreams, FrameTypeUniMaxStreams:
		frame, l, err = parseMaxStreamsFrame(data, frameType, v)
	case FrameTypeDataBlocked:
		frame, l, err = parseDataBlockedFrame(data, v)
	case FrameTypeStreamDataBlocked:
		frame, l, err = parseStreamDataBlockedFrame(data, v)
	case FrameTypeBidiStreamBlocked, FrameTypeUniStreamBlocked:
		frame, l, err = parseStreamsBlockedFrame(data, frameType, v)
	case FrameTypeNewConnectionID:
		frame, l, err = parseNewConnectionIDFrame(data, v)
	case FrameTypeRetireConnectionID:
		frame, l, err = parseRetireConnectionIDFrame(data, v)
	case FrameTypePathChallenge:
		frame, l, err = parsePathChallengeFrame(data, v)
	case FrameTypePathResponse:
		frame, l, err = parsePathResponseFrame(data, v)
	case FrameTypeConnectionClose, FrameTypeApplicationClose:
		frame, l, err = parseConnectionCloseFrame(data, frameType, v)
	case FrameTypeHandshakeDone:
		frame = &HandshakeDoneFrame{}
	case FrameTypeResetStreamAt:
		frame, l, err = parseResetStreamFrame(data, true, v)
	case FrameTypeAckFrequency:
		frame, l, err = parseAckFrequencyFrame(data, v)
	case FrameTypeImmediateAck:
		frame = &ImmediateAckFrame{}
	default:
		err = errUnknownFrameType
	}
	if err != nil {
		return frame, l, &qerr.TransportError{
			ErrorCode:    qerr.FrameEncodingError,
			FrameType:    uint64(frameType),
			ErrorMessage: err.Error(),
		}
	}
	return frame, l, err
}

// SetAckDelayExponent sets the acknowledgment delay exponent (sent in the transport parameters).
// This value is used to scale the ACK Delay field in the ACK frame.
func (p *FrameParser) SetAckDelayExponent(exp uint8) {
	p.ackDelayExponent = exp
}

func replaceUnexpectedEOF(e error) error {
	if e == io.ErrUnexpectedEOF {
		return io.EOF
	}
	return e
}
