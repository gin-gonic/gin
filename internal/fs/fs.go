package fs

import (
	"io/fs"
	"net/http"
	"slices"
	"strings"
)

// FileSystem implements an [fs.FS].
type FileSystem struct {
	http.FileSystem
}

// Open passes `Open` to the upstream implementation and return an [fs.File].
func (o FileSystem) Open(name string) (fs.File, error) {
	f, err := o.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return fs.File(f), nil
}

// ReadDir adapts http.File.Readdir for globbing through io/fs.
func (o FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := o.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []fs.DirEntry
	if dir, ok := f.(fs.ReadDirFile); ok {
		entries, err = dir.ReadDir(-1)
	} else {
		var infos []fs.FileInfo
		infos, err = f.Readdir(-1)
		if infos != nil {
			entries = make([]fs.DirEntry, len(infos))
		}
		for i, info := range infos {
			entries[i] = fs.FileInfoToDirEntry(info)
		}
	}
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return entries, err
}
