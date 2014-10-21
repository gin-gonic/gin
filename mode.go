// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
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
var mode_name string = DebugMode

func init() {
	value := os.Getenv(GIN_MODE)
	if len(value) == 0 {
		SetMode(DebugMode)
	} else {
		SetMode(value)
	}
}

func SetMode(value string) {
	switch value {
	case DebugMode:
		gin_mode = debugCode
	case ReleaseMode:
		gin_mode = releaseCode
	case TestMode:
		gin_mode = testCode
	default:
		panic("gin mode unknown: " + value)
	}
	mode_name = value
}

func Mode() string {
	return mode_name
}

func IsDebugging() bool {
	return gin_mode == debugCode
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		fmt.Printf("[GIN-debug] "+format, values...)
	}
}
