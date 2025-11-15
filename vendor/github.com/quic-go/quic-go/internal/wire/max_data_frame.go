package wire

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// A MaxDataFrame carries flow control information for the connection
type MaxDataFrame struct {
	MaximumData protocol.ByteCount
}

// parseMaxDataFrame parses a MAX_DATA frame
func parseMaxDataFrame(b []byte, _ protocol.Version) (*MaxDataFrame, int, error) {
	frame := &MaxDataFrame{}
	byteOffset, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	frame.MaximumData = protocol.ByteCount(byteOffset)
	return frame, l, nil
}

func (f *MaxDataFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, byte(FrameTypeMaxData))
	b = quicvarint.Append(b, uint64(f.MaximumData))
	return b, nil
}

// Length of a written frame
func (f *MaxDataFrame) Length(_ protocol.Version) protocol.ByteCount {
	return 1 + protocol.ByteCount(quicvarint.Len(uint64(f.MaximumData)))
}
