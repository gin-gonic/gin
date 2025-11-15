package wire

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// A DataBlockedFrame is a DATA_BLOCKED frame
type DataBlockedFrame struct {
	MaximumData protocol.ByteCount
}

func parseDataBlockedFrame(b []byte, _ protocol.Version) (*DataBlockedFrame, int, error) {
	offset, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	return &DataBlockedFrame{MaximumData: protocol.ByteCount(offset)}, l, nil
}

func (f *DataBlockedFrame) Append(b []byte, version protocol.Version) ([]byte, error) {
	b = append(b, byte(FrameTypeDataBlocked))
	return quicvarint.Append(b, uint64(f.MaximumData)), nil
}

// Length of a written frame
func (f *DataBlockedFrame) Length(version protocol.Version) protocol.ByteCount {
	return 1 + protocol.ByteCount(quicvarint.Len(uint64(f.MaximumData)))
}
