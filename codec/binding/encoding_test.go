//go:build gin_bind_encoding

package bindingcodec

import (
	"encoding"
	"errors"
	"reflect"
	"testing"
	"time"

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

// mockTextUnmarshaler implements encoding.TextUnmarshaler for testing
type mockTextUnmarshaler struct {
	receivedText []byte
	returnError  error
}

func (m *mockTextUnmarshaler) UnmarshalText(text []byte) error {
	m.receivedText = text
	return m.returnError
}

func TestTrySetByInterface_WithBindUnmarshaler(t *testing.T) {
	api := bindingApi{}
	mock := &mockBindUnmarshaler{}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "test-value"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.NoError(t, err)
	require.Equal(t, inputVal, mock.receivedParam)
}

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

func TestTrySetByInterface_WithTextUnmarshaler(t *testing.T) {
	api := bindingApi{}
	mock := &mockTextUnmarshaler{}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "text-value"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.NoError(t, err)
	require.Equal(t, []byte(inputVal), mock.receivedText)
}

func TestTrySetByInterface_WithTextUnmarshalerError(t *testing.T) {
	api := bindingApi{}
	expectedErr := errors.New("text unmarshal error")
	mock := &mockTextUnmarshaler{returnError: expectedErr}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "text-value"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.Error(t, err)
}

func TestTrySetByInterface_WithTimeTime(t *testing.T) {
	api := bindingApi{}
	now := time.Now()
	value := reflect.ValueOf(&now).Elem()

	inputVal := "2023-01-01T00:00:00Z"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.False(t, isSet)
	require.NoError(t, err)
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
			ptr := reflect.New(reflect.TypeOf(tc.value))
			ptr.Elem().Set(reflect.ValueOf(tc.value))
			value := ptr.Elem()

			isSet, err := api.TrySetByInterface("input", value)
			require.False(t, isSet)
			require.NoError(t, err)
		})
	}
}

// mockBothInterfaces implements both bindUnmarshaler and encoding.TextUnmarshaler
type mockBothInterfaces struct {
	bindCalled bool
	textCalled bool
}

func (m *mockBothInterfaces) UnmarshalParam(param string) error {
	m.bindCalled = true
	return nil
}

func (m *mockBothInterfaces) UnmarshalText(text []byte) error {
	m.textCalled = true
	return nil
}

var _ encoding.TextUnmarshaler = (*mockBothInterfaces)(nil)

func TestTrySetByInterface_PriorityBindUnmarshaler(t *testing.T) {
	api := bindingApi{}
	mock := &mockBothInterfaces{}
	value := reflect.ValueOf(mock).Elem()

	inputVal := "test"
	isSet, err := api.TrySetByInterface(inputVal, value)
	require.True(t, isSet)
	require.NoError(t, err)
	require.True(t, mock.bindCalled)
	require.False(t, mock.textCalled)
}
