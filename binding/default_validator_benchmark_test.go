package binding

import (
	"errors"
	"strconv"
	"testing"
)

func BenchmarkSliceValidateError(b *testing.B) {
	const size int = 100
	for i := 0; i < b.N; i++ {
		e := make(sliceValidateError, size)
		for j := 0; j < size; j++ {
			e[j] = errors.New(strconv.Itoa(j))
		}
		if len(e.Error()) == 0 {
			b.Errorf("error")
		}
	}
}
