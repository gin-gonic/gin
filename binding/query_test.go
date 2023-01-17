// Copyright 2023 Illia Dimura. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryBinding(t *testing.T) {
	var s struct {
		Foo string `query:"foo"`
	}

	request := &http.Request{URL: &url.URL{RawQuery: "foo=BAR"}}

	err := queryBinding{}.Bind(request, &s)
	require.NoError(t, err)

	assert.Equal(t, "BAR", s.Foo)
}