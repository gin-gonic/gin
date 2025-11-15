package wire

import (
	"fmt"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/quicvarint"
)

// A ResetStreamFrame is a RESET_STREAM or RESET_STREAM_AT frame in QUIC
type ResetStreamFrame struct {
	StreamID     protocol.StreamID
	ErrorCode    qerr.StreamErrorCode
	FinalSize    protocol.ByteCount
	ReliableSize protocol.ByteCount
}

func parseResetStreamFrame(b []byte, isResetStreamAt bool, _ protocol.Version) (*ResetStreamFrame, int, error) {
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
	finalSize, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]

	var reliableSize uint64
	if isResetStreamAt {
		reliableSize, l, err = quicvarint.Parse(b)
		if err != nil {
			return nil, 0, replaceUnexpectedEOF(err)
		}
		b = b[l:]
	}
	if reliableSize > finalSize {
		return nil, 0, fmt.Errorf("RESET_STREAM_AT: reliable size can't be larger than final size (%d vs %d)", reliableSize, finalSize)
	}

	return &ResetStreamFrame{
		StreamID:     protocol.StreamID(streamID),
		ErrorCode:    qerr.StreamErrorCode(errorCode),
		FinalSize:    protocol.ByteCount(finalSize),
		ReliableSize: protocol.ByteCount(reliableSize),
	}, startLen - len(b), nil
}

func (f *ResetStreamFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	if f.ReliableSize == 0 {
		b = quicvarint.Append(b, uint64(FrameTypeResetStream))
	} else {
		b = quicvarint.Append(b, uint64(FrameTypeResetStreamAt))
	}
	b = quicvarint.Append(b, uint64(f.StreamID))
	b = quicvarint.Append(b, uint64(f.ErrorCode))
	b = quicvarint.Append(b, uint64(f.FinalSize))
	if f.ReliableSize > 0 {
		b = quicvarint.Append(b, uint64(f.ReliableSize))
	}
	return b, nil
}

// Length of a written frame
func (f *ResetStreamFrame) Length(protocol.Version) protocol.ByteCount {
	size := 1 // the frame type for both RESET_STREAM and RESET_STREAM_AT fits into 1 byte
	if f.ReliableSize > 0 {
		size += quicvarint.Len(uint64(f.ReliableSize))
	}
	return protocol.ByteCount(size + quicvarint.Len(uint64(f.StreamID)) + quicvarint.Len(uint64(f.ErrorCode)) + quicvarint.Len(uint64(f.FinalSize)))
}
