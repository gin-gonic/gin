// Copyright 2020 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build gofuzz

package render

import (
	"net/http/httptest"
)

func FuzzRender(data []byte) int {
	w := httptest.NewRecorder()
	(YAML{string(data)}).WriteContentType(w)
	err := (YAML{string(data)}).Render(w)
	if err != nil {
		return 0
	}
	return 1
}
