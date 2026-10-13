package render

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRenderAsciiJSONNonBMP is a regression test for AsciiJSON corrupting
// non-BMP Unicode characters (code points above U+FFFF, such as emoji).
//
// It asserts only AsciiJSON's two user-visible contracts: the output must be
// ASCII-only, and it must decode back to the original value. The exact escape
// sequence used (a UTF-16 surrogate pair, a raw rune, etc.) is an
// implementation detail and is intentionally not asserted.
func TestRenderAsciiJSONNonBMP(t *testing.T) {
	const grinningFace = "😀" // U+1F600 GRINNING FACE, a non-BMP code point

	w := httptest.NewRecorder()
	require.NoError(t, (AsciiJSON{map[string]string{"msg": grinningFace}}).Render(w))

	out := w.Body.String()

	// Contract 1: AsciiJSON must emit ASCII-only output.
	for i := 0; i < len(out); i++ {
		require.LessOrEqualf(t, out[i], byte(unicode.MaxASCII),
			"AsciiJSON must emit ASCII-only output, got non-ASCII byte at %d in %q", i, out)
	}

	// Contract 2 (the bug): the rendered output must round-trip back to the
	// original value. The buggy output {"msg":"ὠ0"} is valid JSON but
	// decodes to "ὠ0" instead of "😀".
	var decoded map[string]string
	require.NoErrorf(t, json.Unmarshal([]byte(out), &decoded),
		"AsciiJSON output is not valid JSON: %q", out)
	assert.Equalf(t, grinningFace, decoded["msg"],
		"AsciiJSON corrupted a non-BMP character; rendered output was %q", out)
}
