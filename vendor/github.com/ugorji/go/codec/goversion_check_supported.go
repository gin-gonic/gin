// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !go1.21

package codec

import "errors"

// Moving forward, this codec package will support at least the last 4 major Go releases.
//
// As of early summer 2025, codec will support go 1.21, 1.22, 1.23, 1.24 releases of go.
// This allows use of the followin:
//   - stabilized generics
//   - min/max/clear
//   - slice->array conversion

func init() {
	panic(errors.New("codec: supports go 1.21 and above only"))
}
