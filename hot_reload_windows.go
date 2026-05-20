//go:build windows

// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import "errors"

// RunWithHotReload is not supported on Windows because SIGHUP and fd
// inheritance via ExtraFiles are Unix-only primitives. Use an external
// hot-reload tool such as Air (https://github.com/air-verse/air) instead.
func (engine *Engine) RunWithHotReload(addr ...string) error {
	return errors.New("gin: RunWithHotReload is not supported on Windows")
}
