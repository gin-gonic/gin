package handshake

import (
	"encoding/binary"

	"github.com/quic-go/quic-go/internal/protocol"
)

func createAEAD(suite *cipherSuite, trafficSecret []byte, v protocol.Version) *xorNonceAEAD {
	keyLabel := hkdfLabelKeyV1
	ivLabel := hkdfLabelIVV1
	if v == protocol.Version2 {
		keyLabel = hkdfLabelKeyV2
		ivLabel = hkdfLabelIVV2
	}
	key := hkdfExpandLabel(suite.Hash, trafficSecret, []byte{}, keyLabel, suite.KeyLen)
	iv := hkdfExpandLabel(suite.Hash, trafficSecret, []byte{}, ivLabel, suite.IVLen())
	return suite.AEAD(key, iv)
}

type longHeaderSealer struct {
	aead            *xorNonceAEAD
	headerProtector headerProtector
	nonceBuf        [8]byte
}

var _ LongHeaderSealer = &longHeaderSealer{}

func newLongHeaderSealer(aead *xorNonceAEAD, headerProtector headerProtector) LongHeaderSealer {
	if aead.NonceSize() != 8 {
		panic("unexpected nonce size")
	}
	return &longHeaderSealer{
		aead:            aead,
		headerProtector: headerProtector,
	}
}

func (s *longHeaderSealer) Seal(dst, src []byte, pn protocol.PacketNumber, ad []byte) []byte {
	binary.BigEndian.PutUint64(s.nonceBuf[:], uint64(pn))
	return s.aead.Seal(dst, s.nonceBuf[:], src, ad)
}

func (s *longHeaderSealer) EncryptHeader(sample []byte, firstByte *byte, pnBytes []byte) {
	s.headerProtector.EncryptHeader(sample, firstByte, pnBytes)
}

func (s *longHeaderSealer) Overhead() int {
	return s.aead.Overhead()
}

type longHeaderOpener struct {
	aead            *xorNonceAEAD
	headerProtector headerProtector
	highestRcvdPN   protocol.PacketNumber // highest packet number received (which could be successfully unprotected)

	// use a single array to avoid allocations
	nonceBuf [8]byte
}

var _ LongHeaderOpener = &longHeaderOpener{}

func newLongHeaderOpener(aead *xorNonceAEAD, headerProtector headerProtector) LongHeaderOpener {
	if aead.NonceSize() != 8 {
		panic("unexpected nonce size")
	}
	return &longHeaderOpener{
		aead:            aead,
		headerProtector: headerProtector,
	}
}

func (o *longHeaderOpener) DecodePacketNumber(wirePN protocol.PacketNumber, wirePNLen protocol.PacketNumberLen) protocol.PacketNumber {
	return protocol.DecodePacketNumber(wirePNLen, o.highestRcvdPN, wirePN)
}

func (o *longHeaderOpener) Open(dst, src []byte, pn protocol.PacketNumber, ad []byte) ([]byte, error) {
	binary.BigEndian.PutUint64(o.nonceBuf[:], uint64(pn))
	dec, err := o.aead.Open(dst, o.nonceBuf[:], src, ad)
	if err == nil {
		o.highestRcvdPN = max(o.highestRcvdPN, pn)
	} else {
		err = ErrDecryptionFailed
	}
	return dec, err
}

func (o *longHeaderOpener) DecryptHeader(sample []byte, firstByte *byte, pnBytes []byte) {
	o.headerProtector.DecryptHeader(sample, firstByte, pnBytes)
}
