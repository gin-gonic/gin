// Copyright 2026 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"mime"
	"net/http"
	"net/url"
)

func parseMultipartForm(req *http.Request, maxMemory int64) error {
	err := req.ParseMultipartForm(maxMemory)
	if err == nil || !errors.Is(err, http.ErrNotMultipart) {
		return err
	}

	mediaType, _, parseErr := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if parseErr != nil || mediaType != MIMEMultipartMixed {
		return err
	}

	reader, readerErr := req.MultipartReader()
	if readerErr != nil {
		return readerErr
	}

	form, readErr := reader.ReadForm(maxMemory)
	if readErr != nil {
		return readErr
	}
	req.MultipartForm = form

	if req.PostForm == nil {
		req.PostForm = make(url.Values)
	}
	for key, values := range form.Value {
		req.PostForm[key] = append(req.PostForm[key], values...)
	}

	if req.Form == nil {
		req.Form = make(url.Values, len(req.PostForm))
		for key, values := range req.PostForm {
			req.Form[key] = append(req.Form[key], values...)
		}
	}
	if req.URL != nil {
		for key, values := range req.URL.Query() {
			req.Form[key] = append(req.Form[key], values...)
		}
	}

	return nil
}
