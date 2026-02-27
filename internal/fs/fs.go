package fs

import (
	"io/fs"
	"net/http"
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
