// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"os"

	"github.com/gin-gonic/gin/binding"
	"github.com/mattn/go-colorable"
)

const ENV_GIN_MODE = "GIN_MODE"

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestMode    string = "test"
)
const (
	debugCode   = iota
	releaseCode = iota
	testCode    = iota
)

var DefaultWriter io.Writer = colorable.NewColorableStdout()
var ginMode int = debugCode
var modeName string = DebugMode

func init() {
	mode := os.Getenv(ENV_GIN_MODE)
	if len(mode) == 0 {
		SetMode(DebugMode)
	} else {
		SetMode(mode)
	}
}

func SetMode(value string) {
	switch value {
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("gin mode unknown: " + value)
	}
	modeName = value
}

func DisableBindValidation() {
	binding.Validator = nil
}

func Mode() string {
	return modeName
}
