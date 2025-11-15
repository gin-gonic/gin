package wire

import (
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// A CryptoFrame is a CRYPTO frame
type CryptoFrame struct {
	Offset protocol.ByteCount
	Data   []byte
}

func parseCryptoFrame(b []byte, _ protocol.Version) (*CryptoFrame, int, error) {
	startLen := len(b)
	frame := &CryptoFrame{}
	offset, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	frame.Offset = protocol.ByteCount(offset)
	dataLen, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, 0, replaceUnexpectedEOF(err)
	}
	b = b[l:]
	if dataLen > uint64(len(b)) {
		return nil, 0, io.EOF
	}
	if dataLen != 0 {
		frame.Data = make([]byte, dataLen)
		copy(frame.Data, b)
	}
	return frame, startLen - len(b) + int(dataLen), nil
}

func (f *CryptoFrame) Append(b []byte, _ protocol.Version) ([]byte, error) {
	b = append(b, byte(FrameTypeCrypto))
	b = quicvarint.Append(b, uint64(f.Offset))
	b = quicvarint.Append(b, uint64(len(f.Data)))
	b = append(b, f.Data...)
	return b, nil
}

// Length of a written frame
func (f *CryptoFrame) Length(_ protocol.Version) protocol.ByteCount {
	return protocol.ByteCount(1 + quicvarint.Len(uint64(f.Offset)) + quicvarint.Len(uint64(len(f.Data))) + len(f.Data))
}

// MaxDataLen returns the maximum data length
func (f *CryptoFrame) MaxDataLen(maxSize protocol.ByteCount) protocol.ByteCount {
	// pretend that the data size will be 1 bytes
	// if it turns out that varint encoding the length will consume 2 bytes, we need to adjust the data length afterwards
	headerLen := protocol.ByteCount(1 + quicvarint.Len(uint64(f.Offset)) + 1)
	if headerLen > maxSize {
		return 0
	}
	maxDataLen := maxSize - headerLen
	if quicvarint.Len(uint64(maxDataLen)) != 1 {
		maxDataLen--
	}
	return maxDataLen
}

// MaybeSplitOffFrame splits a frame such that it is not bigger than n bytes.
// It returns if the frame was actually split.
// The frame might not be split if:
// * the size is large enough to fit the whole frame
// * the size is too small to fit even a 1-byte frame. In that case, the frame returned is nil.
func (f *CryptoFrame) MaybeSplitOffFrame(maxSize protocol.ByteCount, version protocol.Version) (*CryptoFrame, bool /* was splitting required */) {
	if f.Length(version) <= maxSize {
		return nil, false
	}

	n := f.MaxDataLen(maxSize)
	if n == 0 {
		return nil, true
	}

	newLen := protocol.ByteCount(len(f.Data)) - n

	new := &CryptoFrame{}
	new.Offset = f.Offset
	new.Data = make([]byte, newLen)

	// swap the data slices
	new.Data, f.Data = f.Data, new.Data

	copy(f.Data, new.Data[n:])
	new.Data = new.Data[:n]
	f.Offset += n

	return new, true
}
