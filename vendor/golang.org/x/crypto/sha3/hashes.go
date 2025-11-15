// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sha3 implements the SHA-3 hash algorithms and the SHAKE extendable
// output functions defined in FIPS 202.
//
// Most of this package is a wrapper around the crypto/sha3 package in the
// standard library. The only exception is the legacy Keccak hash functions.
package sha3

import (
	"crypto/sha3"
	"hash"
)

// New224 creates a new SHA3-224 hash.
// Its generic security strength is 224 bits against preimage attacks,
// and 112 bits against collision attacks.
//
// It is a wrapper for the [sha3.New224] function in the standard library.
//
//go:fix inline
func New224() hash.Hash {
	return sha3.New224()
}

// New256 creates a new SHA3-256 hash.
// Its generic security strength is 256 bits against preimage attacks,
// and 128 bits against collision attacks.
//
// It is a wrapper for the [sha3.New256] function in the standard library.
//
//go:fix inline
func New256() hash.Hash {
	return sha3.New256()
}

// New384 creates a new SHA3-384 hash.
// Its generic security strength is 384 bits against preimage attacks,
// and 192 bits against collision attacks.
//
// It is a wrapper for the [sha3.New384] function in the standard library.
//
//go:fix inline
func New384() hash.Hash {
	return sha3.New384()
}

// New512 creates a new SHA3-512 hash.
// Its generic security strength is 512 bits against preimage attacks,
// and 256 bits against collision attacks.
//
// It is a wrapper for the [sha3.New512] function in the standard library.
//
//go:fix inline
func New512() hash.Hash {
	return sha3.New512()
}

// Sum224 returns the SHA3-224 digest of the data.
//
// It is a wrapper for the [sha3.Sum224] function in the standard library.
//
//go:fix inline
func Sum224(data []byte) [28]byte {
	return sha3.Sum224(data)
}

// Sum256 returns the SHA3-256 digest of the data.
//
// It is a wrapper for the [sha3.Sum256] function in the standard library.
//
//go:fix inline
func Sum256(data []byte) [32]byte {
	return sha3.Sum256(data)
}

// Sum384 returns the SHA3-384 digest of the data.
//
// It is a wrapper for the [sha3.Sum384] function in the standard library.
//
//go:fix inline
func Sum384(data []byte) [48]byte {
	return sha3.Sum384(data)
}

// Sum512 returns the SHA3-512 digest of the data.
//
// It is a wrapper for the [sha3.Sum512] function in the standard library.
//
//go:fix inline
func Sum512(data []byte) [64]byte {
	return sha3.Sum512(data)
}
