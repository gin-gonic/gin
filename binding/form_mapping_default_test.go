//go:build !gin_bind_encoding

package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If someone does not specify parser=TextUnmarshaler even when it's defined for the type, gin should ignore the
// UnmarshalText logic and continue using its default binding logic. (This ensures gin does not break backwards
// compatibility)
//
// Note: TestMappingUsingBindUnmarshalerAndTextUnmarshalerWhenOnlyTextUnmarshalerDefined works differently when:
// - form_mapping_encoding_test.go (with gin_bind_encoding build tag enabled)
func TestMappingUsingBindUnmarshalerAndTextUnmarshalerWhenOnlyTextUnmarshalerDefined(t *testing.T) {
	var s struct {
		Hex                customUnmarshalTextHex `form:"hex"`
		HexByUnmarshalText customUnmarshalTextHex `form:"hex2,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&s, formSource{
		"hex":  {`11`},
		"hex2": {`11`},
	}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, 11, s.Hex)                  // this is using default int binding, not our custom hex binding. 0x11 should be 17 in decimal
	assert.EqualValues(t, 0x11, s.HexByUnmarshalText) // correct expected value for normal hex binding
}
