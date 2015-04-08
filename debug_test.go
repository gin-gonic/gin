// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDebugging(t *testing.T) {
	SetMode(DebugMode)
	assert.True(t, IsDebugging())
	SetMode(ReleaseMode)
	assert.False(t, IsDebugging())
	SetMode(TestMode)
	assert.False(t, IsDebugging())
}

// TODO
// func TestDebugPrint(t *testing.T) {
// 	buffer := bytes.NewBufferString("")
// 	debugLogger.
// 	log.SetOutput(buffer)

// 	SetMode(ReleaseMode)
// 	debugPrint("This is a example")
// 	assert.Equal(t, buffer.Len(), 0)

// 	SetMode(DebugMode)
// 	debugPrint("This is %s", "a example")
// 	assert.Equal(t, buffer.String(), "[GIN-debug] This is a example")

// 	SetMode(TestMode)
// 	log.SetOutput(os.Stdout)
// }
