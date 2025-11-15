package wire

import (
	"errors"
	"fmt"
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// A NewConnectionIDFrame is a NEW_CONNECTION_ID frame
type NewConnectionIDFrame struct {
	SequenceNumber      uint64
	RetirePriorTo       uint64
	ConnectionID        protocol.ConnectionID
	StatelessResetToken protocol.StatelessResetToken
}

func parseNewConnectionIDFrame(b []byte, _ protocol.Version) (*NewConnectionIDFrame, int, error) {
	startLen := len(b)
	seq, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	ret, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	if ret > seq {
		//nolint:staticcheck // SA1021: Retire Prior To is the name of the field
		return nil, 0, fmt.Errorf("Retire Prior To value (%d) larger than Sequence Number (%d)", ret, seq)
	}
	if len(b) == 0 {
		return nil, 0, io.EOF
	}
	connIDLen := int(b[0])
	b = b[1:]
	if connIDLen == 0 {
		return nil, 0, errors.New("invalid zero-length connection ID")
	}
	if connIDLen > protocol.MaxConnIDLen {
		return nil, 0, protocol.ErrInvalidConnectionIDLen
	}
	if len(b) < connIDLen {
		return nil, 0, io.EOF
	}
	frame := &NewConnectionIDFrame{
		SequenceNumber: seq,
		RetirePriorTo:  ret,
		ConnectionID:   protocol.ParseConnectionID(b[:connIDLen]),
	}
	b = b[connIDLen:]
	if len(b) < len(frame.StatelessResetToken) {
		return nil, 0, io.EOF
	}
	copy(frame.StatelessResetToken[:], b)
	return frame, startLen - len(b) + len(frame.StatelessResetToken), nil
}

func (f *NewConnectionIDFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, byte(FrameTypeNewConnectionID))
	b = quicvarint.Append(b, f.SequenceNumber)
	b = quicvarint.Append(b, f.RetirePriorTo)
	connIDLen := f.ConnectionID.Len()
	if connIDLen > protocol.MaxConnIDLen {
		return nil, fmt.Errorf("invalid connection ID length: %d", connIDLen)
	}
	b = append(b, uint8(connIDLen))
	b = append(b, f.ConnectionID.Bytes()...)
	b = append(b, f.StatelessResetToken[:]...)
	return b, nil
}

// Length of a written frame
func (f *NewConnectionIDFrame) Length(protocol.Version) protocol.ByteCount {
	return 1 + protocol.ByteCount(quicvarint.Len(f.SequenceNumber)+quicvarint.Len(f.RetirePriorTo)+1 /* connection ID length */ +f.ConnectionID.Len()) + 16
}
