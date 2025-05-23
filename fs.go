// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"os"
)

// OnlyFilesFS implements an http.FileSystem without `Readdir` functionality.
type OnlyFilesFS struct {
	FileSystem http.FileSystem
}

// Open passes `Open` to the upstream implementation without `Readdir` functionality.
func (o OnlyFilesFS) Open(name string) (http.File, error) {
	f, err := o.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return neutralizedReaddirFile{f}, nil
}

// neutralizedReaddirFile wraps http.File with a specific implementation of `Readdir`.
type neutralizedReaddirFile struct {
	http.File
}

// Readdir overrides the http.File default implementation and always returns nil.
func (n neutralizedReaddirFile) Readdir(_ int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}

// Dir returns an http.FileSystem that can be used by http.FileServer().
// It is used internally in router.Static().
// if listDirectory == true, then it works the same as http.Dir(),
// otherwise it returns a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)

	if listDirectory {
		return fs
	}

	return &OnlyFilesFS{FileSystem: fs}
}
