package handshake

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/chacha20"

	"github.com/quic-go/quic-go/internal/protocol"
)

type headerProtector interface {
	EncryptHeader(sample []byte, firstByte *byte, hdrBytes []byte)
	DecryptHeader(sample []byte, firstByte *byte, hdrBytes []byte)
}

func hkdfHeaderProtectionLabel(v protocol.Version) string {
	if v == protocol.Version2 {
		return "quicv2 hp"
	}
	return "quic hp"
}

func newHeaderProtector(suite *cipherSuite, trafficSecret []byte, isLongHeader bool, v protocol.Version) headerProtector {
	hkdfLabel := hkdfHeaderProtectionLabel(v)
	switch suite.ID {
	case tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384:
		return newAESHeaderProtector(suite, trafficSecret, isLongHeader, hkdfLabel)
	case tls.TLS_CHACHA20_POLY1305_SHA256:
		return newChaChaHeaderProtector(suite, trafficSecret, isLongHeader, hkdfLabel)
	default:
		panic(fmt.Sprintf("Invalid cipher suite id: %d", suite.ID))
	}
}

type aesHeaderProtector struct {
	mask         [16]byte // AES always has a 16 byte block size
	block        cipher.Block
	isLongHeader bool
}

var _ headerProtector = &aesHeaderProtector{}

func newAESHeaderProtector(suite *cipherSuite, trafficSecret []byte, isLongHeader bool, hkdfLabel string) headerProtector {
	hpKey := hkdfExpandLabel(suite.Hash, trafficSecret, []byte{}, hkdfLabel, suite.KeyLen)
	block, err := aes.NewCipher(hpKey)
	if err != nil {
		panic(fmt.Sprintf("error creating new AES cipher: %s", err))
	}
	return &aesHeaderProtector{
		block:        block,
		isLongHeader: isLongHeader,
	}
}

func (p *aesHeaderProtector) DecryptHeader(sample []byte, firstByte *byte, hdrBytes []byte) {
	p.apply(sample, firstByte, hdrBytes)
}

func (p *aesHeaderProtector) EncryptHeader(sample []byte, firstByte *byte, hdrBytes []byte) {
	p.apply(sample, firstByte, hdrBytes)
}

func (p *aesHeaderProtector) apply(sample []byte, firstByte *byte, hdrBytes []byte) {
	if len(sample) != len(p.mask) {
		panic("invalid sample size")
	}
	p.block.Encrypt(p.mask[:], sample)
	if p.isLongHeader {
		*firstByte ^= p.mask[0] & 0xf
	} else {
		*firstByte ^= p.mask[0] & 0x1f
	}
	for i := range hdrBytes {
		hdrBytes[i] ^= p.mask[i+1]
	}
}

type chachaHeaderProtector struct {
	mask [5]byte

	key          [32]byte
	isLongHeader bool
}

var _ headerProtector = &chachaHeaderProtector{}

func newChaChaHeaderProtector(suite *cipherSuite, trafficSecret []byte, isLongHeader bool, hkdfLabel string) headerProtector {
	hpKey := hkdfExpandLabel(suite.Hash, trafficSecret, []byte{}, hkdfLabel, suite.KeyLen)

	p := &chachaHeaderProtector{
		isLongHeader: isLongHeader,
	}
	copy(p.key[:], hpKey)
	return p
}

func (p *chachaHeaderProtector) DecryptHeader(sample []byte, firstByte *byte, hdrBytes []byte) {
	p.apply(sample, firstByte, hdrBytes)
}

func (p *chachaHeaderProtector) EncryptHeader(sample []byte, firstByte *byte, hdrBytes []byte) {
	p.apply(sample, firstByte, hdrBytes)
}

func (p *chachaHeaderProtector) apply(sample []byte, firstByte *byte, hdrBytes []byte) {
	if len(sample) != 16 {
		panic("invalid sample size")
	}
	for i := 0; i < 5; i++ {
		p.mask[i] = 0
	}
	cipher, err := chacha20.NewUnauthenticatedCipher(p.key[:], sample[4:])
	if err != nil {
		panic(err)
	}
	cipher.SetCounter(binary.LittleEndian.Uint32(sample[:4]))
	cipher.XORKeyStream(p.mask[:], p.mask[:])
	p.applyMask(firstByte, hdrBytes)
}

func (p *chachaHeaderProtector) applyMask(firstByte *byte, hdrBytes []byte) {
	if p.isLongHeader {
		*firstByte ^= p.mask[0] & 0xf
	} else {
		*firstByte ^= p.mask[0] & 0x1f
	}
	for i := range hdrBytes {
		hdrBytes[i] ^= p.mask[i+1]
	}
}
