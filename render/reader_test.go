// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestReaderRenderDoesNotModifyHeaders(t *testing.T) {
	headers := map[string]string{"Cache-Control": "public"}

	w := httptest.NewRecorder()
	r := Reader{
		ContentLength: 5,
		Reader:        strings.NewReader("hello"),
		Headers:       headers,
	}
	require.NoError(t, r.Render(w))
	require.Equal(t, map[string]string{"Cache-Control": "public"}, headers)
	require.Equal(t, "5", w.Header().Get("Content-Length"))

	// Reusing the same map for a response of unknown length must not
	// carry over the Content-Length of the previous response.
	w = httptest.NewRecorder()
	r = Reader{
		ContentLength: -1,
		Reader:        strings.NewReader("hello, world!"),
		Headers:       headers,
	}
	require.NoError(t, r.Render(w))
	require.Empty(t, w.Header().Get("Content-Length"))
	require.Equal(t, "public", w.Header().Get("Cache-Control"))
	require.Equal(t, "hello, world!", w.Body.String())
}

// Run with -race: rendering concurrently with a shared Headers map must not
// write to the map.
func TestReaderRenderConcurrentSharedHeaders(t *testing.T) {
	headers := map[string]string{"Cache-Control": "public"}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			r := Reader{
				ContentLength: 5,
				Reader:        strings.NewReader("hello"),
				Headers:       headers,
			}
			assert.NoError(t, r.Render(httptest.NewRecorder()))
		}()
	}
	wg.Wait()

	require.Equal(t, map[string]string{"Cache-Control": "public"}, headers)
}
