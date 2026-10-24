// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderRenderNoHeaders(t *testing.T) {
	content := "test"
	r := Reader{
		ContentLength: int64(len(content)),
		Reader:        strings.NewReader(content),
	}
	err := r.Render(httptest.NewRecorder())
	require.NoError(t, err)
}

func TestReaderRenderDoesNotMutateHeaders(t *testing.T) {
	headers := map[string]string{"Cache-Control": "public"}
	r := Reader{
		ContentType:   "text/plain",
		ContentLength: 5,
		Reader:        strings.NewReader("hello"),
		Headers:       headers,
	}
	w := httptest.NewRecorder()
	require.NoError(t, r.Render(w))
	require.Equal(t, "5", w.Header().Get("Content-Length"))
	require.NotContains(t, headers, "Content-Length", "caller's headers map must stay read-only")
}

func TestReaderRenderNoStaleContentLengthOnMapReuse(t *testing.T) {
	headers := map[string]string{"Cache-Control": "public"}
	r1 := Reader{
		ContentType:   "text/plain",
		ContentLength: 5,
		Reader:        strings.NewReader("hello"),
		Headers:       headers,
	}
	require.NoError(t, r1.Render(httptest.NewRecorder()))

	r2 := Reader{
		ContentType:   "text/plain",
		ContentLength: -1,
		Reader:        strings.NewReader("hello, world!"),
		Headers:       headers,
	}
	w2 := httptest.NewRecorder()
	require.NoError(t, r2.Render(w2))
	require.Empty(t, w2.Header().Get("Content-Length"), "stale Content-Length must not leak into later responses")
}

func TestReaderRenderConcurrentSharedHeaders(t *testing.T) {
	headers := map[string]string{"Cache-Control": "public"}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := Reader{
				ContentType:   "text/plain",
				ContentLength: 5,
				Reader:        strings.NewReader("hello"),
				Headers:       headers,
			}
			_ = r.Render(httptest.NewRecorder())
		}()
	}
	wg.Wait()
	require.NotContains(t, headers, "Content-Length")
}
