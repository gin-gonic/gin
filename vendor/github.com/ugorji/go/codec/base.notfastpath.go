// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build notfastpath || codec.notfastpath

package codec

import "reflect"

const fastpathEnabled = false

// The generated fast-path code is very large, and adds a few seconds to the build time.
// This causes test execution, execution of small tools which use codec, etc
// to take a long time.
//
// To mitigate, we now support the notfastpath tag.
// This tag disables fastpath during build, allowing for faster build, test execution,
// short-program runs, etc.

// func fastpathEncodeTypeSwitchSlice(iv interface{}, e *Encoder) bool { return false }
// func fastpathEncodeTypeSwitchMap(iv interface{}, e *Encoder) bool   { return false }

func fastpathDecodeSetZeroTypeSwitch(iv interface{}) bool { return false }

func fastpathAvIndex(rtid uintptr) (uint, bool) { return 0, false }

type fastpathRtRtid struct {
	rtid uintptr
	rt   reflect.Type
}

type fastpathARtRtid [0]fastpathRtRtid

var fastpathAvRtRtid fastpathARtRtid
