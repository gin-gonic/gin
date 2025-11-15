// Copyright (c) 2015 Klaus Post, released under MIT License. See LICENSE file.

//go:build (386 && !gccgo && !noasm && !appengine) || (amd64 && !gccgo && !noasm && !appengine)
// +build 386,!gccgo,!noasm,!appengine amd64,!gccgo,!noasm,!appengine

package cpuid

func asmCpuid(op uint32) (eax, ebx, ecx, edx uint32)
func asmCpuidex(op, op2 uint32) (eax, ebx, ecx, edx uint32)
func asmXgetbv(index uint32) (eax, edx uint32)
func asmRdtscpAsm() (eax, ebx, ecx, edx uint32)
func asmDarwinHasAVX512() bool

func initCPU() {
	cpuid = asmCpuid
	cpuidex = asmCpuidex
	xgetbv = asmXgetbv
	rdtscpAsm = asmRdtscpAsm
	darwinHasAVX512 = asmDarwinHasAVX512
}

func addInfo(c *CPUInfo, safe bool) {
	c.maxFunc = maxFunctionID()
	c.maxExFunc = maxExtendedFunction()
	c.BrandName = brandName()
	c.CacheLine = cacheLine()
	c.Family, c.Model, c.Stepping = familyModel()
	c.featureSet = support()
	c.SGX = hasSGX(c.featureSet.inSet(SGX), c.featureSet.inSet(SGXLC))
	c.AMDMemEncryption = hasAMDMemEncryption(c.featureSet.inSet(SME) || c.featureSet.inSet(SEV))
	c.ThreadsPerCore = threadsPerCore()
	c.LogicalCores = logicalCores()
	c.PhysicalCores = physicalCores()
	c.VendorID, c.VendorString = vendorID()
	c.HypervisorVendorID, c.HypervisorVendorString = hypervisorVendorID()
	c.AVX10Level = c.supportAVX10()
	c.cacheSize()
	c.frequencies()
	if c.maxFunc >= 0x0A {
		eax, ebx, _, edx := cpuid(0x0A)
		c.PMU = parseLeaf0AH(c, eax, ebx, edx)
	}
}

func getVectorLength() (vl, pl uint64) { return 0, 0 }
