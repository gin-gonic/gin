package handshake

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// TokenProtectorKey is the key used to encrypt both Retry and session resumption tokens.
type TokenProtectorKey [32]byte

const tokenNonceSize = 32

// tokenProtector is used to create and verify a token
type tokenProtector struct {
	key TokenProtectorKey
}

// newTokenProtector creates a source for source address tokens
func newTokenProtector(key TokenProtectorKey) *tokenProtector {
	return &tokenProtector{key: key}
}

// NewToken encodes data into a new token.
func (s *tokenProtector) NewToken(data []byte) ([]byte, error) {
	var nonce [tokenNonceSize]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	aead, aeadNonce, err := s.createAEAD(nonce[:])
	if err != nil {
		return nil, err
	}
	return append(nonce[:], aead.Seal(nil, aeadNonce, data, nil)...), nil
}

// DecodeToken decodes a token.
func (s *tokenProtector) DecodeToken(p []byte) ([]byte, error) {
	if len(p) < tokenNonceSize {
		return nil, fmt.Errorf("token too short: %d", len(p))
	}
	nonce := p[:tokenNonceSize]
	aead, aeadNonce, err := s.createAEAD(nonce)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, aeadNonce, p[tokenNonceSize:], nil)
}

func (s *tokenProtector) createAEAD(nonce []byte) (cipher.AEAD, []byte, error) {
	h := hkdf.New(sha256.New, s.key[:], nonce, []byte("quic-go token source"))
	key := make([]byte, 32) // use a 32 byte key, in order to select AES-256
	if _, err := io.ReadFull(h, key); err != nil {
		return nil, nil, err
	}
	aeadNonce := make([]byte, 12)
	if _, err := io.ReadFull(h, aeadNonce); err != nil {
		return nil, nil, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aead, err := cipher.NewGCM(c)
	if err != nil {
		return nil, nil, err
	}
	return aead, aeadNonce, nil
}
