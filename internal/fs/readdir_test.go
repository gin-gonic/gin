package fs

import (
	"errors"
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type readDirFile struct {
	http.File
	infos  []fs.FileInfo
	err    error
	closed bool
}

func (f *readDirFile) Readdir(count int) ([]fs.FileInfo, error) { return f.infos, f.err }
func (f *readDirFile) Close() error                             { f.closed = true; return nil }

func TestFileSystemReadDir(t *testing.T) {
	source := fstest.MapFS{"z.html": &fstest.MapFile{}, "a.html": &fstest.MapFile{}}
	z, err := fs.Stat(source, "z.html")
	require.NoError(t, err)
	a, err := fs.Stat(source, "a.html")
	require.NoError(t, err)
	directory := &readDirFile{infos: []fs.FileInfo{z, a}}
	filesystem := FileSystem{&mockFileSystem{open: func(name string) (http.File, error) {
		assert.Equal(t, ".", name)
		return directory, nil
	}}}
	entries, err := fs.ReadDir(filesystem, ".")
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "a.html", entries[0].Name())
	assert.Equal(t, "z.html", entries[1].Name())
	assert.True(t, directory.closed)
	info, err := entries[0].Info()
	require.NoError(t, err)
	assert.Equal(t, a, info)
}

type modernReadDirFile struct {
	readDirFile
	entries []fs.DirEntry
}

func (f *modernReadDirFile) ReadDir(count int) ([]fs.DirEntry, error) { return f.entries, f.err }
func (f *modernReadDirFile) Readdir(count int) ([]fs.FileInfo, error) {
	panic("ReadDir should be preferred")
}

func TestFileSystemReadDirPartialAndModern(t *testing.T) {
	source := fstest.MapFS{"z.html": &fstest.MapFile{}, "a.html": &fstest.MapFile{}}
	z, err := fs.Stat(source, "z.html")
	require.NoError(t, err)
	a, err := fs.Stat(source, "a.html")
	require.NoError(t, err)
	sentinel := errors.New("partial directory read")
	for _, modern := range []bool{false, true} {
		directory := &readDirFile{infos: []fs.FileInfo{z, a}, err: sentinel}
		var file http.File = directory
		if modern {
			file = &modernReadDirFile{readDirFile: *directory, entries: []fs.DirEntry{fs.FileInfoToDirEntry(z), fs.FileInfoToDirEntry(a)}}
		}
		filesystem := FileSystem{&mockFileSystem{open: func(string) (http.File, error) { return file, nil }}}
		entries, err := filesystem.ReadDir(".")
		require.ErrorIs(t, err, sentinel)
		require.Len(t, entries, 2)
		assert.Equal(t, "a.html", entries[0].Name())
		assert.Equal(t, "z.html", entries[1].Name())
		if modern {
			assert.True(t, file.(*modernReadDirFile).closed)
		} else {
			assert.True(t, directory.closed)
		}
	}
}

func TestFileSystemReadDirErrors(t *testing.T) {
	sentinel := errors.New("directory read failed")
	t.Run("open", func(t *testing.T) {
		filesystem := FileSystem{&mockFileSystem{open: func(string) (http.File, error) { return nil, sentinel }}}
		entries, err := filesystem.ReadDir(".")
		require.ErrorIs(t, err, sentinel)
		assert.Nil(t, entries)
	})
	t.Run("read", func(t *testing.T) {
		directory := &readDirFile{err: sentinel}
		filesystem := FileSystem{&mockFileSystem{open: func(string) (http.File, error) { return directory, nil }}}
		entries, err := filesystem.ReadDir(".")
		require.ErrorIs(t, err, sentinel)
		assert.Nil(t, entries)
		assert.True(t, directory.closed)
	})
	t.Run("empty", func(t *testing.T) {
		directory := &readDirFile{}
		filesystem := FileSystem{&mockFileSystem{open: func(string) (http.File, error) { return directory, nil }}}
		entries, err := filesystem.ReadDir(".")
		require.NoError(t, err)
		assert.Empty(t, entries)
		assert.True(t, directory.closed)
	})
}
