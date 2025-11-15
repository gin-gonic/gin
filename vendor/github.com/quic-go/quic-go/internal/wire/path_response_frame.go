package wire

import (
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
)

// A PathResponseFrame is a PATH_RESPONSE frame
type PathResponseFrame struct {
	Data [8]byte
}

func parsePathResponseFrame(b []byte, _ protocol.Version) (*PathResponseFrame, int, error) {
	f := &PathResponseFrame{}
	if len(b) < 8 {
		return nil, 0, io.EOF
	}
	copy(f.Data[:], b)
	return f, 8, nil
}

func (f *PathResponseFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, byte(FrameTypePathResponse))
	b = append(b, f.Data[:]...)
	return b, nil
}

// Length of a written frame
func (f *PathResponseFrame) Length(_ protocol.Version) protocol.ByteCount {
	return 1 + 8
}
