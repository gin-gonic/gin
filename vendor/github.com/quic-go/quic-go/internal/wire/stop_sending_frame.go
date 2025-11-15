package wire

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/quicvarint"
)

// A StopSendingFrame is a STOP_SENDING frame
type StopSendingFrame struct {
	StreamID  protocol.StreamID
	ErrorCode qerr.StreamErrorCode
}

// parseStopSendingFrame parses a STOP_SENDING frame
func parseStopSendingFrame(b []byte, _ protocol.Version) (*StopSendingFrame, int, error) {
	startLen := len(b)
	streamID, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	errorCode, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]

	return &StopSendingFrame{
		StreamID:  protocol.StreamID(streamID),
		ErrorCode: qerr.StreamErrorCode(errorCode),
	}, startLen - len(b), nil
}

// Length of a written frame
func (f *StopSendingFrame) Length(_ protocol.Version) protocol.ByteCount {
	return 1 + protocol.ByteCount(quicvarint.Len(uint64(f.StreamID))+quicvarint.Len(uint64(f.ErrorCode)))
}

func (f *StopSendingFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, byte(FrameTypeStopSending))
	b = quicvarint.Append(b, uint64(f.StreamID))
	b = quicvarint.Append(b, uint64(f.ErrorCode))
	return b, nil
}
