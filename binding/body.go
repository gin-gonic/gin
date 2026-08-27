// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"io"
	"net/http"
)

// MaxBodyBytes is the maximum number of bytes the BSON and Protobuf binders
// will read from a request body. Zero (the default) means no limit, matching
// historical behavior. Applications that bind untrusted BSON or Protobuf
// input should set this to a positive value during initialization, for
// example 32 << 20.
var MaxBodyBytes int64

func readBody(r io.Reader) ([]byte, error) {
	if MaxBodyBytes <= 0 {
		return io.ReadAll(r)
	}
	body, err := io.ReadAll(io.LimitReader(r, MaxBodyBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > MaxBodyBytes {
		return nil, &http.MaxBytesError{Limit: MaxBodyBytes}
	}
	return body, nil
}
