// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http/httptest"
	"strconv"
	"strings"
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

func TestReaderRenderWithHeaders(t *testing.T) {
	content := "test"
	r := Reader{
		ContentLength: int64(len(content)),
		Reader:        strings.NewReader(content),
		Headers: map[string]string{
			"Test-Content": "test/content",
		},
	}
	recorder := httptest.NewRecorder()
	err := r.Render(recorder)
	require.NoError(t, err)

	require.Contains(t, recorder.Header()["Content-Length"], strconv.FormatInt(r.ContentLength, 10))

	require.Contains(t, recorder.Header()["Test-Content"], "test/content")
}
