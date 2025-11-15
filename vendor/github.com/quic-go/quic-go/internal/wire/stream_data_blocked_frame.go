package wire

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// A StreamDataBlockedFrame is a STREAM_DATA_BLOCKED frame
type StreamDataBlockedFrame struct {
	StreamID          protocol.StreamID
	MaximumStreamData protocol.ByteCount
}

func parseStreamDataBlockedFrame(b []byte, _ protocol.Version) (*StreamDataBlockedFrame, int, error) {
	startLen := len(b)
	sid, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	offset, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}

	return &StreamDataBlockedFrame{
		StreamID:          protocol.StreamID(sid),
		MaximumStreamData: protocol.ByteCount(offset),
	}, startLen - len(b) + l, nil
}

func (f *StreamDataBlockedFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, 0x15)
	b = quicvarint.Append(b, uint64(f.StreamID))
	b = quicvarint.Append(b, uint64(f.MaximumStreamData))
	return b, nil
}

// Length of a written frame
func (f *StreamDataBlockedFrame) Length(protocol.Version) protocol.ByteCount {
	return 1 + protocol.ByteCount(quicvarint.Len(uint64(f.StreamID))+quicvarint.Len(uint64(f.MaximumStreamData)))
}
