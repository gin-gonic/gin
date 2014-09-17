// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"os"
)

const GIN_MODE = "GIN_MODE"

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

var gin_mode int = debugCode

func SetMode(value string) {
	switch value {
	case DebugMode:
		gin_mode = debugCode
	case ReleaseMode:
		gin_mode = releaseCode
	case TestMode:
		gin_mode = testCode
	default:
		panic("gin mode unknown, the allowed modes are: " + DebugMode + " and " + ReleaseMode)
	}
}

func init() {
	value := os.Getenv(GIN_MODE)
	if len(value) == 0 {
		SetMode(DebugMode)
	} else {
		SetMode(value)
	}
}
