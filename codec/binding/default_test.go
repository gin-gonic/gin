//go:build !gin_bind_encoding

package bindingcodec

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockBindUnmarshaler implements the bindUnmarshaler interface for testing
type mockBindUnmarshaler struct {
	receivedParam string
	returnError   error
}

func (m *mockBindUnmarshaler) UnmarshalParam(param string) error {
	m.receivedParam = param
	return m.returnError
}

// TestTrySetByInterface_WithBindUnmarshaler tests successful binding with bindUnmarshaler
func TestTrySetByInterface_WithBindUnmarshaler(t *testing.T) {
	api := bindingApi{}
	mock := &mockBindUnmarshaler{}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "test-value"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.NoError(t, err)
	require.Equal(t, "test-value", mock.receivedParam)
}

// TestTrySetByInterface_WithBindUnmarshalerError tests error handling from bindUnmarshaler
func TestTrySetByInterface_WithBindUnmarshalerError(t *testing.T) {
	api := bindingApi{}
	expectedErr := errors.New("unmarshal error")
	mock := &mockBindUnmarshaler{returnError: expectedErr}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "test-value"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.Error(t, err)
}

// TestTrySetByInterface_WithoutInterface tests behavior with regular types
func TestTrySetByInterface_WithoutInterface(t *testing.T) {
	api := bindingApi{}

	testCases := []struct {
		name  string
		value any
	}{
		{"string", "test"},
		{"int", 42},
		{"bool", true},
		{"struct", struct{ Field string }{"value"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pointer to make it addressable
			ptr := reflect.New(reflect.TypeOf(tc.value))
			ptr.Elem().Set(reflect.ValueOf(tc.value))
			value := ptr.Elem()

			isSet, err := api.TrySetByInterface("input", value)
			require.False(t, isSet)
			require.NoError(t, err)
		})
	}
}

// TestTrySetByInterface_WithPointer tests that the method works with pointer values
func TestTrySetByInterface_WithPointer(t *testing.T) {
	api := bindingApi{}
	mock := &mockBindUnmarshaler{}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "pointer-test"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.NoError(t, err)
	require.Equal(t, inputVal, mock.receivedParam)
}
