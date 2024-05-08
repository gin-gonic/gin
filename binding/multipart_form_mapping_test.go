// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormMultipartBindingBindOneFile(t *testing.T) {
	var s struct {
		FileValue   multipart.FileHeader     `form:"file"`
		FilePtr     *multipart.FileHeader    `form:"file"`
		SliceValues []multipart.FileHeader   `form:"file"`
		SlicePtrs   []*multipart.FileHeader  `form:"file"`
		ArrayValues [1]multipart.FileHeader  `form:"file"`
		ArrayPtrs   [1]*multipart.FileHeader `form:"file"`
	}
	file := testFile{"file", "file1", []byte("hello")}

	req := createRequestMultipartFiles(t, file)
	err := FormMultipart.Bind(req, &s)
	assert.NoError(t, err)

	assertMultipartFileHeader(t, &s.FileValue, file)
	assertMultipartFileHeader(t, s.FilePtr, file)
	assert.Len(t, s.SliceValues, 1)
	assertMultipartFileHeader(t, &s.SliceValues[0], file)
	assert.Len(t, s.SlicePtrs, 1)
	assertMultipartFileHeader(t, s.SlicePtrs[0], file)
	assertMultipartFileHeader(t, &s.ArrayValues[0], file)
	assertMultipartFileHeader(t, s.ArrayPtrs[0], file)
}

func TestFormMultipartBindingBindTwoFiles(t *testing.T) {
	var s struct {
		SliceValues []multipart.FileHeader   `form:"file"`
		SlicePtrs   []*multipart.FileHeader  `form:"file"`
		ArrayValues [2]multipart.FileHeader  `form:"file"`
		ArrayPtrs   [2]*multipart.FileHeader `form:"file"`
	}
	files := []testFile{
		{"file", "file1", []byte("hello")},
		{"file", "file2", []byte("world")},
	}

	req := createRequestMultipartFiles(t, files...)
	err := FormMultipart.Bind(req, &s)
	assert.NoError(t, err)

	assert.Len(t, s.SliceValues, len(files))
	assert.Len(t, s.SlicePtrs, len(files))
	assert.Len(t, s.ArrayValues, len(files))
	assert.Len(t, s.ArrayPtrs, len(files))

	for i, file := range files {
		assertMultipartFileHeader(t, &s.SliceValues[i], file)
		assertMultipartFileHeader(t, s.SlicePtrs[i], file)
		assertMultipartFileHeader(t, &s.ArrayValues[i], file)
		assertMultipartFileHeader(t, s.ArrayPtrs[i], file)
	}
}

func TestFormMultipartBindingBindError(t *testing.T) {
	files := []testFile{
		{"file", "file1", []byte("hello")},
		{"file", "file2", []byte("world")},
	}

	for _, tt := range []struct {
		name string
		s    any
	}{
		{"wrong type", &struct {
			Files int `form:"file"`
		}{}},
		{"wrong array size", &struct {
			Files [1]*multipart.FileHeader `form:"file"`
		}{}},
		{"wrong slice type", &struct {
			Files []int `form:"file"`
		}{}},
	} {
		req := createRequestMultipartFiles(t, files...)
		err := FormMultipart.Bind(req, tt.s)
		assert.Error(t, err)
	}
}

type testFile struct {
	Fieldname string
	Filename  string
	Content   []byte
}

func createRequestMultipartFiles(t *testing.T, files ...testFile) *http.Request {
	var body bytes.Buffer

	mw := multipart.NewWriter(&body)
	for _, file := range files {
		fw, err := mw.CreateFormFile(file.Fieldname, file.Filename)
		assert.NoError(t, err)

		n, err := fw.Write(file.Content)
		assert.NoError(t, err)
		assert.Equal(t, len(file.Content), n)
	}
	err := mw.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/", &body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())
	return req
}

func assertMultipartFileHeader(t *testing.T, fh *multipart.FileHeader, file testFile) {
	assert.Equal(t, file.Filename, fh.Filename)
	assert.Equal(t, int64(len(file.Content)), fh.Size)

	fl, err := fh.Open()
	assert.NoError(t, err)

	body, err := io.ReadAll(fl)
	assert.NoError(t, err)
	assert.Equal(t, string(file.Content), string(body))

	err = fl.Close()
	assert.NoError(t, err)
}
