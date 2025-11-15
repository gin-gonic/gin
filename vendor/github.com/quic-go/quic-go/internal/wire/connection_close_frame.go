package wire

import (
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// A ConnectionCloseFrame is a CONNECTION_CLOSE frame
type ConnectionCloseFrame struct {
	IsApplicationError bool
	ErrorCode          uint64
	FrameType          uint64
	ReasonPhrase       string
}

func parseConnectionCloseFrame(b []byte, typ FrameType, _ protocol.Version) (*ConnectionCloseFrame, int, error) {
	startLen := len(b)
	f := &ConnectionCloseFrame{IsApplicationError: typ == FrameTypeApplicationClose}
	ec, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	f.ErrorCode = ec
	// read the Frame Type, if this is not an application error
	if !f.IsApplicationError {
		ft, l, err := quicvarint.Parse(b)
		if err != nil {
			return nil, 0, replaceUnexpectedEOF(err)
		}
		b = b[l:]
		f.FrameType = ft
	}
	var reasonPhraseLen uint64
	reasonPhraseLen, l, err = quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	if int(reasonPhraseLen) > len(b) {
		return nil, 0, io.EOF
	}

	reasonPhrase := make([]byte, reasonPhraseLen)
	copy(reasonPhrase, b)
	f.ReasonPhrase = string(reasonPhrase)
	return f, startLen - len(b) + int(reasonPhraseLen), nil
}

// Length of a written frame
func (f *ConnectionCloseFrame) Length(protocol.Version) protocol.ByteCount {
	length := 1 + protocol.ByteCount(quicvarint.Len(f.ErrorCode)+quicvarint.Len(uint64(len(f.ReasonPhrase)))) + protocol.ByteCount(len(f.ReasonPhrase))
	if !f.IsApplicationError {
		length += protocol.ByteCount(quicvarint.Len(f.FrameType)) // for the frame type
	}
	return length
}

func (f *ConnectionCloseFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	if f.IsApplicationError {
		b = append(b, byte(FrameTypeApplicationClose))
	} else {
		b = append(b, byte(FrameTypeConnectionClose))
	}

	b = quicvarint.Append(b, f.ErrorCode)
	if !f.IsApplicationError {
		b = quicvarint.Append(b, f.FrameType)
	}
	b = quicvarint.Append(b, uint64(len(f.ReasonPhrase)))
	b = append(b, []byte(f.ReasonPhrase)...)
	return b, nil
}
