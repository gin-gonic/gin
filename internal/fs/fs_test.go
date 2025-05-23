package fs

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFileSystem struct {
	open func(name string) (http.File, error)
}

func (m *mockFileSystem) Open(name string) (http.File, error) {
	return m.open(name)
}

func TestFileSystem_Open(t *testing.T) {
	var testFile *os.File
	mockFS := &mockFileSystem{
		open: func(name string) (http.File, error) {
			return testFile, nil
		},
	}
	fs := &FileSystem{mockFS}

	file, err := fs.Open("foo")

	require.NoError(t, err)
	assert.Equal(t, testFile, file)
}

func TestFileSystem_Open_err(t *testing.T) {
	testError := errors.New("mock")
	mockFS := &mockFileSystem{
		open: func(_ string) (http.File, error) {
			return nil, testError
		},
	}
	fs := &FileSystem{mockFS}

	file, err := fs.Open("foo")

	require.ErrorIs(t, err, testError)
	assert.Nil(t, file)
}
