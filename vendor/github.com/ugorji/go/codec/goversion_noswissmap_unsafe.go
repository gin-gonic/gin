// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build !safe && !codec.safe && !appengine && !go1.24

package codec

import "unsafe"

// retrofited from hIter struct

type unsafeMapIterPadding struct {
	_ [6]unsafe.Pointer // padding: *maptype, *hmap, buckets, *bmap, overflow, oldoverflow,
	_ [4]uintptr        // padding: uintptr, uint8, bool fields
	_ uintptr           // padding: wasted (try to fill cache-line at multiple of 4)
}
