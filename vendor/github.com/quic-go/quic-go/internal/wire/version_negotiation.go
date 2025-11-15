package wire

import (
	"crypto/rand"
	"encoding/binary"
	"errors"

	"github.com/quic-go/quic-go/internal/protocol"
)

// ParseVersionNegotiationPacket parses a Version Negotiation packet.
func ParseVersionNegotiationPacket(b []byte) (dest, src protocol.ArbitraryLenConnectionID, _ []protocol.Version, _ error) {
	n, dest, src, err := ParseArbitraryLenConnectionIDs(b)
	if err != nil {
		return nil, nil, nil, err
	}
	b = b[n:]
	if len(b) == 0 {
		//nolint:staticcheck // SA1021: the packet is called Version Negotiation packet
		return nil, nil, nil, errors.New("Version Negotiation packet has empty version list")
	}
	if len(b)%4 != 0 {
		//nolint:staticcheck // SA1021: the packet is called Version Negotiation packet
		return nil, nil, nil, errors.New("Version Negotiation packet has a version list with an invalid length")
	}
	versions := make([]protocol.Version, len(b)/4)
	for i := 0; len(b) > 0; i++ {
		versions[i] = protocol.Version(binary.BigEndian.Uint32(b[:4]))
		b = b[4:]
	}
	return dest, src, versions, nil
}

// ComposeVersionNegotiation composes a Version Negotiation
func ComposeVersionNegotiation(destConnID, srcConnID protocol.ArbitraryLenConnectionID, versions []protocol.Version) []byte {
	greasedVersions := protocol.GetGreasedVersions(versions)
	expectedLen := 1 /* type byte */ + 4 /* version field */ + 1 /* dest connection ID length field */ + destConnID.Len() + 1 /* src connection ID length field */ + srcConnID.Len() + len(greasedVersions)*4
	buf := make([]byte, 1+4 /* type byte and version field */, expectedLen)
	_, _ = rand.Read(buf[:1]) // ignore the error here. It is not critical to have perfect random here.
	// Setting the "QUIC bit" (0x40) is not required by the RFC,
	// but it allows clients to demultiplex QUIC with a long list of other protocols.
	// See RFC 9443 and https://mailarchive.ietf.org/arch/msg/quic/oR4kxGKY6mjtPC1CZegY1ED4beg/ for details.
	buf[0] |= 0xc0
	// The next 4 bytes are left at 0 (version number).
	buf = append(buf, uint8(destConnID.Len()))
	buf = append(buf, destConnID.Bytes()...)
	buf = append(buf, uint8(srcConnID.Len()))
	buf = append(buf, srcConnID.Bytes()...)
	for _, v := range greasedVersions {
		buf = binary.BigEndian.AppendUint32(buf, uint32(v))
	}
	return buf
}
