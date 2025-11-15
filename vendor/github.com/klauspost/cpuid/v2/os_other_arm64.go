// Copyright (c) 2020 Klaus Post, released under MIT License. See LICENSE file.

//go:build arm64 && !linux && !darwin
// +build arm64,!linux,!darwin

package cpuid

import "runtime"

func detectOS(c *CPUInfo) bool {
	c.PhysicalCores = runtime.NumCPU()
	// For now assuming 1 thread per core...
	c.ThreadsPerCore = 1
	c.LogicalCores = c.PhysicalCores
	return false
}
