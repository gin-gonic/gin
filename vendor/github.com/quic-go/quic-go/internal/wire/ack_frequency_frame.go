package wire

import (
	"math"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

type AckFrequencyFrame struct {
	SequenceNumber        uint64
	AckElicitingThreshold uint64
	RequestMaxAckDelay    time.Duration
	ReorderingThreshold   protocol.PacketNumber
}

func parseAckFrequencyFrame(b []byte, _ protocol.Version) (*AckFrequencyFrame, int, error) {
	startLen := len(b)
	seq, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	aeth, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	mad, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	// prevents overflows if the peer sends a very large value
	maxAckDelay := time.Duration(mad) * time.Microsecond
	if maxAckDelay < 0 {
		maxAckDelay = math.MaxInt64
	}
	b = b[l:]
	rth, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]

	return &AckFrequencyFrame{
		SequenceNumber:        seq,
		AckElicitingThreshold: aeth,
		RequestMaxAckDelay:    maxAckDelay,
		ReorderingThreshold:   protocol.PacketNumber(rth),
	}, startLen - len(b), nil
}

func (f *AckFrequencyFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = quicvarint.Append(b, uint64(FrameTypeAckFrequency))
	b = quicvarint.Append(b, f.SequenceNumber)
	b = quicvarint.Append(b, f.AckElicitingThreshold)
	b = quicvarint.Append(b, uint64(f.RequestMaxAckDelay/time.Microsecond))
	return quicvarint.Append(b, uint64(f.ReorderingThreshold)), nil
}

func (f *AckFrequencyFrame) Length(_ protocol.Version) protocol.ByteCount {
	return protocol.ByteCount(2 + quicvarint.Len(f.SequenceNumber) + quicvarint.Len(f.AckElicitingThreshold) +
		quicvarint.Len(uint64(f.RequestMaxAckDelay/time.Microsecond)) + quicvarint.Len(uint64(f.ReorderingThreshold)))
}
