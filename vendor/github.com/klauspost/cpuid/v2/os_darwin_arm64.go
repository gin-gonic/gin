// Copyright (c) 2020 Klaus Post, released under MIT License. See LICENSE file.

package cpuid

import (
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

func detectOS(c *CPUInfo) bool {
	if runtime.GOOS != "ios" {
		tryToFillCPUInfoFomSysctl(c)
	}
	// There are no hw.optional sysctl values for the below features on Mac OS 11.0
	// to detect their supported state dynamically. Assume the CPU features that
	// Apple Silicon M1 supports to be available as a minimal set of features
	// to all Go programs running on darwin/arm64.
	// TODO: Add more if we know them.
	c.featureSet.setIf(runtime.GOOS != "ios", AESARM, PMULL, SHA1, SHA2)

	return true
}

func sysctlGetBool(name string) bool {
	value, err := unix.SysctlUint32(name)
	if err != nil {
		return false
	}
	return value != 0
}

func sysctlGetString(name string) string {
	value, err := unix.Sysctl(name)
	if err != nil {
		return ""
	}
	return value
}

func sysctlGetInt(unknown int, names ...string) int {
	for _, name := range names {
		value, err := unix.SysctlUint32(name)
		if err != nil {
			continue
		}
		if value != 0 {
			return int(value)
		}
	}
	return unknown
}

func sysctlGetInt64(unknown int, names ...string) int {
	for _, name := range names {
		value64, err := unix.SysctlUint64(name)
		if err != nil {
			continue
		}
		if int(value64) != unknown {
			return int(value64)
		}
	}
	return unknown
}

func setFeature(c *CPUInfo, feature FeatureID, aliases ...string) {
	for _, alias := range aliases {
		set := sysctlGetBool(alias)
		c.featureSet.setIf(set, feature)
		if set {
			break
		}
	}
}

func tryToFillCPUInfoFomSysctl(c *CPUInfo) {
	c.BrandName = sysctlGetString("machdep.cpu.brand_string")

	if len(c.BrandName) != 0 {
		c.VendorString = strings.Fields(c.BrandName)[0]
	}

	c.PhysicalCores = sysctlGetInt(runtime.NumCPU(), "hw.physicalcpu")
	c.ThreadsPerCore = sysctlGetInt(1, "machdep.cpu.thread_count", "kern.num_threads") /
		sysctlGetInt(1, "hw.physicalcpu")
	c.LogicalCores = sysctlGetInt(runtime.NumCPU(), "machdep.cpu.core_count")
	c.Family = sysctlGetInt(0, "machdep.cpu.family", "hw.cpufamily")
	c.Model = sysctlGetInt(0, "machdep.cpu.model")
	c.CacheLine = sysctlGetInt64(0, "hw.cachelinesize")
	c.Cache.L1I = sysctlGetInt64(-1, "hw.l1icachesize")
	c.Cache.L1D = sysctlGetInt64(-1, "hw.l1dcachesize")
	c.Cache.L2 = sysctlGetInt64(-1, "hw.l2cachesize")
	c.Cache.L3 = sysctlGetInt64(-1, "hw.l3cachesize")

	// ARM features:
	//
	// Note: On some Apple Silicon system, some feats have aliases. See:
	// https://developer.apple.com/documentation/kernel/1387446-sysctlbyname/determining_instruction_set_characteristics
	// When so, we look at all aliases and consider a feature available when at least one identifier matches.
	setFeature(c, AESARM, "hw.optional.arm.FEAT_AES")                                   // AES instructions
	setFeature(c, ASIMD, "hw.optional.arm.AdvSIMD", "hw.optional.neon")                 // Advanced SIMD
	setFeature(c, ASIMDDP, "hw.optional.arm.FEAT_DotProd")                              // SIMD Dot Product
	setFeature(c, ASIMDHP, "hw.optional.arm.AdvSIMD_HPFPCvt", "hw.optional.neon_hpfp")  // Advanced SIMD half-precision floating point
	setFeature(c, ASIMDRDM, "hw.optional.arm.FEAT_RDM")                                 // Rounding Double Multiply Accumulate/Subtract
	setFeature(c, ATOMICS, "hw.optional.arm.FEAT_LSE", "hw.optional.armv8_1_atomics")   // Large System Extensions (LSE)
	setFeature(c, CRC32, "hw.optional.arm.FEAT_CRC32", "hw.optional.armv8_crc32")       // CRC32/CRC32C instructions
	setFeature(c, DCPOP, "hw.optional.arm.FEAT_DPB")                                    // Data cache clean to Point of Persistence (DC CVAP)
	setFeature(c, EVTSTRM, "hw.optional.arm.FEAT_ECV")                                  // Generic timer
	setFeature(c, FCMA, "hw.optional.arm.FEAT_FCMA", "hw.optional.armv8_3_compnum")     // Floating point complex number addition and multiplication
	setFeature(c, FHM, "hw.optional.armv8_2_fhm", "hw.optional.arm.FEAT_FHM")           // FMLAL and FMLSL instructions
	setFeature(c, FP, "hw.optional.floatingpoint")                                      // Single-precision and double-precision floating point
	setFeature(c, FPHP, "hw.optional.arm.FEAT_FP16", "hw.optional.neon_fp16")           // Half-precision floating point
	setFeature(c, GPA, "hw.optional.arm.FEAT_PAuth")                                    // Generic Pointer Authentication
	setFeature(c, JSCVT, "hw.optional.arm.FEAT_JSCVT")                                  // Javascript-style double->int convert (FJCVTZS)
	setFeature(c, LRCPC, "hw.optional.arm.FEAT_LRCPC")                                  // Weaker release consistency (LDAPR, etc)
	setFeature(c, PMULL, "hw.optional.arm.FEAT_PMULL")                                  // Polynomial Multiply instructions (PMULL/PMULL2)
	setFeature(c, RNDR, "hw.optional.arm.FEAT_RNG")                                     // Random Number instructions
	setFeature(c, TLB, "hw.optional.arm.FEAT_TLBIOS", "hw.optional.arm.FEAT_TLBIRANGE") // Outer Shareable and TLB range maintenance instructions
	setFeature(c, TS, "hw.optional.arm.FEAT_FlagM", "hw.optional.arm.FEAT_FlagM2")      // Flag manipulation instructions
	setFeature(c, SHA1, "hw.optional.arm.FEAT_SHA1")                                    // SHA-1 instructions (SHA1C, etc)
	setFeature(c, SHA2, "hw.optional.arm.FEAT_SHA256")                                  // SHA-2 instructions (SHA256H, etc)
	setFeature(c, SHA3, "hw.optional.arm.FEAT_SHA3")                                    // SHA-3 instructions (EOR3, RAXI, XAR, BCAX)
	setFeature(c, SHA512, "hw.optional.arm.FEAT_SHA512")                                // SHA512 instructions
	setFeature(c, SM3, "hw.optional.arm.FEAT_SM3")                                      // SM3 instructions
	setFeature(c, SM4, "hw.optional.arm.FEAT_SM4")                                      // SM4 instructions
	setFeature(c, SVE, "hw.optional.arm.FEAT_SVE")                                      // Scalable Vector Extension
}
