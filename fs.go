package gin

import (
	"net/http"
	"os"
)

type (
	onlyfilesFS struct {
		fs http.FileSystem
	}
	neuteredReaddirFile struct {
		http.File
	}
)

func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	} else {
		return &onlyfilesFS{fs}
	}
}

func (fs onlyfilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
