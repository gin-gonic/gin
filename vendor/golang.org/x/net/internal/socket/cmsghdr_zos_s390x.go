// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

func (h *cmsghdr) set(l, lvl, typ int) {
	h.Len = int32(l)
	h.Level = int32(lvl)
	h.Type = int32(typ)
}
