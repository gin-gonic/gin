package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetCookieDataPreservesInput(t *testing.T) {
	cookie := &http.Cookie{Name: "session", Value: "value", SameSite: http.SameSiteDefaultMode}
	original := *cookie
	writer := httptest.NewRecorder()
	context, _ := CreateTestContext(writer)
	context.SetSameSite(http.SameSiteStrictMode)
	context.SetCookieData(cookie)

	assert.Equal(t, original, *cookie)
	require.Len(t, writer.Result().Cookies(), 1)
	result := writer.Result().Cookies()[0]
	assert.Equal(t, "/", result.Path)
	assert.Equal(t, http.SameSiteStrictMode, result.SameSite)
}

func TestSetCookieDataReusesContextDefaults(t *testing.T) {
	cookie := &http.Cookie{Name: "session", Value: "value", SameSite: http.SameSiteDefaultMode}
	for _, mode := range []http.SameSite{http.SameSiteStrictMode, http.SameSiteLaxMode} {
		writer := httptest.NewRecorder()
		context, _ := CreateTestContext(writer)
		context.SetSameSite(mode)
		context.SetCookieData(cookie)
		require.Len(t, writer.Result().Cookies(), 1)
		assert.Equal(t, mode, writer.Result().Cookies()[0].SameSite)
	}
}
