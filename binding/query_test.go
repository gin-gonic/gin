package binding

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryBindingWithQueryTag(t *testing.T) {
	var s struct {
		Foo string `query:"foo"`
		Bar string `query:"bar"`
	}

	request := &http.Request{URL: &url.URL{RawQuery: "foo=HELLO&bar=WORLD"}}

	err := queryBinding{}.Bind(request, &s)
	require.NoError(t, err)

	assert.Equal(t, "HELLO", s.Foo)
	assert.Equal(t, "WORLD", s.Bar)
}

func TestQueryBindingWithFormTag(t *testing.T) {
	var s struct {
		Foo string `form:"foo"`
		Bar string `form:"bar"`
	}

	request := &http.Request{URL: &url.URL{RawQuery: "foo=HELLO&bar=WORLD"}}

	err := queryBinding{}.Bind(request, &s)
	require.NoError(t, err)

	assert.Equal(t, "HELLO", s.Foo)
	assert.Equal(t, "WORLD", s.Bar)
}

func TestQueryBindingMixedTags(t *testing.T) {
	var s struct {
		Foo string `query:"foo"`
		Bar string `form:"bar"`
	}

	request := &http.Request{URL: &url.URL{RawQuery: "foo=HELLO&bar=WORLD"}}

	err := queryBinding{}.Bind(request, &s)
	require.NoError(t, err)

	assert.Equal(t, "HELLO", s.Foo)
	assert.Equal(t, "WORLD", s.Bar)
}