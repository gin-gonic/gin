// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build appengine
//go:build appengine

package gin

func init() {
	defaultAppEngine = true
}
