// Copyright 2017 Bo-Yi Wu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build jsoniter

package json

import "github.com/json-iterator/go"

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
)
