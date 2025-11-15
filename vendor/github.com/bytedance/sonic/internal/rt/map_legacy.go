// +build !go1.24

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
}
