// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build appengine
// +build appengine

package gin

func init() {
	defaultPlatform = PlatformGoogleAppEngine
}
