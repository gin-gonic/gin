//go:build go1.24 && !go1.26 && !goexperiment.swissmap
// +build go1.24,!go1.26,!goexperiment.swissmap

package rt

import (
	"unsafe"
)

type GoMapIterator struct {
	K           unsafe.Pointer
	V           unsafe.Pointer
	T           *GoMapType
	H           unsafe.Pointer
	Buckets     unsafe.Pointer
	Bptr        *unsafe.Pointer
	Overflow    *[]unsafe.Pointer
	OldOverflow *[]unsafe.Pointer
	StartBucket uintptr
	Offset      uint8
	Wrapped     bool
	B           uint8
	I           uint8
	Bucket      uintptr
	CheckBucket uintptr
	// different from go1.23
	ClearSeq    uint64
}
