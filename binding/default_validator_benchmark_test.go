// Copyright 2022 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"strconv"
	"testing"
)

func BenchmarkSliceValidationError(b *testing.B) {
	const size int = 100
	e := make(SliceValidationError, size)
	for j := 0; j < size; j++ {
		e[j] = errors.New(strconv.Itoa(j))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if len(e.Error()) == 0 {
			b.Errorf("error")
		}
	}
}
