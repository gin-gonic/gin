//go:build gin_bind_encoding

package binding

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// In gin_bind_encoding mode, TextUnmarshaler is used automatically when present, even without an
// explicit parser tag.
func TestMappingUsingBindUnmarshalerAndTextUnmarshalerWhenOnlyTextUnmarshalerDefined_DefaultEncodingUnmarshalText(t *testing.T) {
	var s struct {
		Hex                customUnmarshalTextHex `form:"hex"`
		HexByUnmarshalText customUnmarshalTextHex `form:"hex2,parser=encoding.TextUnmarshaler"`
	}
	err := mappingByPtr(&s, formSource{
		"hex":  {`11`},
		"hex2": {`11`},
	}, "form")
	require.NoError(t, err)

	assert.EqualValues(t, 0x11, s.Hex)
	assert.EqualValues(t, 0x11, s.HexByUnmarshalText)
}

// ==== Automatic TextUnmarshaler binding tests (no parser tag required) ====

func TestMappingTextUnmarshalerAutoBindForm(t *testing.T) {
	var s struct {
		ID objectIDUnmarshalText `form:"id"`
	}
	err := mappingByPtr(&s, formSource{"id": {"664a062ac74a8ad104e0e80f"}}, "form")
	require.NoError(t, err)
	expected, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80f")
	assert.Equal(t, expected, s.ID)
}

func TestMappingTextUnmarshalerAutoBindURI(t *testing.T) {
	var s struct {
		ID objectIDUnmarshalText `uri:"id"`
	}
	err := mappingByPtr(&s, formSource{"id": {"664a062ac74a8ad104e0e80f"}}, "uri")
	require.NoError(t, err)
	expected, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80f")
	assert.Equal(t, expected, s.ID)
}

func TestMappingTextUnmarshalerAutoBindSlice(t *testing.T) {
	var s struct {
		IDs []objectIDUnmarshalText `form:"ids" collection_format:"csv"`
	}
	err := mappingByPtr(&s, formSource{"ids": {"664a062ac74a8ad104e0e80e,664a062ac74a8ad104e0e80f"}}, "form")
	require.NoError(t, err)
	id1, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80e")
	id2, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80f")
	expected := []objectIDUnmarshalText{id1, id2}
	assert.Equal(t, expected, s.IDs)
}

func TestMappingTextUnmarshalerAutoBindMultipleValues(t *testing.T) {
	var s struct {
		IDs []objectIDUnmarshalText `form:"ids"`
	}
	err := mappingByPtr(&s, formSource{"ids": {
		"664a062ac74a8ad104e0e80e",
		"664a062ac74a8ad104e0e80f",
	}}, "form")
	require.NoError(t, err)
	id1, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80e")
	id2, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80f")
	assert.Equal(t, []objectIDUnmarshalText{id1, id2}, s.IDs)
}

func TestMappingTextUnmarshalerAutoBindDefault(t *testing.T) {
	var s struct {
		ID objectIDUnmarshalText `form:"id,default=664a062ac74a8ad104e0e80f"`
	}
	err := mappingByPtr(&s, formSource{}, "form")
	require.NoError(t, err)
	expected, _ := convertToOidUnmarshalText("664a062ac74a8ad104e0e80f")
	assert.Equal(t, expected, s.ID)
}

func TestMappingTextUnmarshalerAutoBindInvalidValue(t *testing.T) {
	var s struct {
		ID objectIDUnmarshalText `form:"id"`
	}
	err := mappingByPtr(&s, formSource{"id": {"not-a-valid-objectid"}}, "form")
	require.Error(t, err)
}

// BindUnmarshaler should take precedence over TextUnmarshaler
type testDualUnmarshaler struct {
	Value string
}

func (d *testDualUnmarshaler) UnmarshalParam(param string) error {
	d.Value = "param:" + param
	return nil
}

func (d *testDualUnmarshaler) UnmarshalText(text []byte) error {
	d.Value = "text:" + string(text)
	return nil
}

var _ BindUnmarshaler = (*testDualUnmarshaler)(nil)
var _ encoding.TextUnmarshaler = (*testDualUnmarshaler)(nil)

func TestMappingBindUnmarshalerTakesPrecedenceOverTextUnmarshaler(t *testing.T) {
	var s struct {
		Field testDualUnmarshaler `form:"field"`
	}
	err := mappingByPtr(&s, formSource{"field": {"hello"}}, "form")
	require.NoError(t, err)
	assert.Equal(t, "param:hello", s.Field.Value) // BindUnmarshaler wins
}
