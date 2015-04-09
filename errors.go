// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
)

const (
	ErrorTypeInternal = 1 << iota
	ErrorTypeExternal = 1 << iota
	ErrorTypeAll      = 0xffffffff
)

// Used internally to collect errors that occurred during an http request.
type errorMsg struct {
	Err  string      `json:"error"`
	Type int         `json:"-"`
	Meta interface{} `json:"meta"`
}

type errorMsgs []errorMsg

func (a errorMsgs) ByType(typ int) errorMsgs {
	if len(a) == 0 {
		return a
	}
	result := make(errorMsgs, 0, len(a))
	for _, msg := range a {
		if msg.Type&typ > 0 {
			result = append(result, msg)
		}
	}
	return result
}

func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		text := fmt.Sprintf("Error #%02d: %s\n     Meta: %v\n", (i + 1), msg.Err, msg.Meta)
		buffer.WriteString(text)
	}
	return buffer.String()
}
