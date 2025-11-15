// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha3

import (
	"crypto/sha3"
	"hash"
	"io"
)

// ShakeHash defines the interface to hash functions that support
// arbitrary-length output. When used as a plain [hash.Hash], it
// produces minimum-length outputs that provide full-strength generic
// security.
type ShakeHash interface {
	hash.Hash

	// Read reads more output from the hash; reading affects the hash's
	// state. (ShakeHash.Read is thus very different from Hash.Sum.)
	// It never returns an error, but subsequent calls to Write or Sum
	// will panic.
	io.Reader

	// Clone returns a copy of the ShakeHash in its current state.
	Clone() ShakeHash
}

// NewShake128 creates a new SHAKE128 variable-output-length ShakeHash.
// Its generic security strength is 128 bits against all attacks if at
// least 32 bytes of its output are used.
func NewShake128() ShakeHash {
	return &shakeWrapper{sha3.NewSHAKE128(), 32, false, sha3.NewSHAKE128}
}

// NewShake256 creates a new SHAKE256 variable-output-length ShakeHash.
// Its generic security strength is 256 bits against all attacks if
// at least 64 bytes of its output are used.
func NewShake256() ShakeHash {
	return &shakeWrapper{sha3.NewSHAKE256(), 64, false, sha3.NewSHAKE256}
}

// NewCShake128 creates a new instance of cSHAKE128 variable-output-length ShakeHash,
// a customizable variant of SHAKE128.
// N is used to define functions based on cSHAKE, it can be empty when plain cSHAKE is
// desired. S is a customization byte string used for domain separation - two cSHAKE
// computations on same input with different S yield unrelated outputs.
// When N and S are both empty, this is equivalent to NewShake128.
func NewCShake128(N, S []byte) ShakeHash {
	return &shakeWrapper{sha3.NewCSHAKE128(N, S), 32, false, func() *sha3.SHAKE {
		return sha3.NewCSHAKE128(N, S)
	}}
}

// NewCShake256 creates a new instance of cSHAKE256 variable-output-length ShakeHash,
// a customizable variant of SHAKE256.
// N is used to define functions based on cSHAKE, it can be empty when plain cSHAKE is
// desired. S is a customization byte string used for domain separation - two cSHAKE
// computations on same input with different S yield unrelated outputs.
// When N and S are both empty, this is equivalent to NewShake256.
func NewCShake256(N, S []byte) ShakeHash {
	return &shakeWrapper{sha3.NewCSHAKE256(N, S), 64, false, func() *sha3.SHAKE {
		return sha3.NewCSHAKE256(N, S)
	}}
}

// ShakeSum128 writes an arbitrary-length digest of data into hash.
func ShakeSum128(hash, data []byte) {
	h := NewShake128()
	h.Write(data)
	h.Read(hash)
}

// ShakeSum256 writes an arbitrary-length digest of data into hash.
func ShakeSum256(hash, data []byte) {
	h := NewShake256()
	h.Write(data)
	h.Read(hash)
}

// shakeWrapper adds the Size, Sum, and Clone methods to a sha3.SHAKE
// to implement the ShakeHash interface.
type shakeWrapper struct {
	*sha3.SHAKE
	outputLen int
	squeezing bool
	newSHAKE  func() *sha3.SHAKE
}

func (w *shakeWrapper) Read(p []byte) (n int, err error) {
	w.squeezing = true
	return w.SHAKE.Read(p)
}

func (w *shakeWrapper) Clone() ShakeHash {
	s := w.newSHAKE()
	b, err := w.MarshalBinary()
	if err != nil {
		panic(err) // unreachable
	}
	if err := s.UnmarshalBinary(b); err != nil {
		panic(err) // unreachable
	}
	return &shakeWrapper{s, w.outputLen, w.squeezing, w.newSHAKE}
}

func (w *shakeWrapper) Size() int { return w.outputLen }

func (w *shakeWrapper) Sum(b []byte) []byte {
	if w.squeezing {
		panic("sha3: Sum after Read")
	}
	out := make([]byte, w.outputLen)
	// Clone the state so that we don't affect future Write calls.
	s := w.Clone()
	s.Read(out)
	return append(b, out...)
}
