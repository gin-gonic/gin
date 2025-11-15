// Copyright (c) 2015 Klaus Post, released under MIT License. See LICENSE file.

//go:build (!amd64 && !386 && !arm64) || gccgo || noasm || appengine
// +build !amd64,!386,!arm64 gccgo noasm appengine

package cpuid

func initCPU() {
	cpuid = func(uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
	cpuidex = func(x, y uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
	xgetbv = func(uint32) (a, b uint32) { return 0, 0 }
	rdtscpAsm = func() (a, b, c, d uint32) { return 0, 0, 0, 0 }

}

func addInfo(info *CPUInfo, safe bool) {}
func getVectorLength() (vl, pl uint64) { return 0, 0 }
