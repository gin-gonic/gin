// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bytesconv

// StringToBytes converts string to byte slice
func StringToBytes(s string) []byte {
	return []byte(s)
}

// BytesToString converts byte slice to string
func BytesToString(b []byte) string {
	return string(b)
}
