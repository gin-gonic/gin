// Copyright (c) 2015 Klaus Post, released under MIT License. See LICENSE file.

//go:build arm64 && !gccgo && !noasm && !appengine
// +build arm64,!gccgo,!noasm,!appengine

package cpuid

import "runtime"

func getMidr() (midr uint64)
func getProcFeatures() (procFeatures uint64)
func getInstAttributes() (instAttrReg0, instAttrReg1 uint64)
func getVectorLength() (vl, pl uint64)

func initCPU() {
	cpuid = func(uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
	cpuidex = func(x, y uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
	xgetbv = func(uint32) (a, b uint32) { return 0, 0 }
	rdtscpAsm = func() (a, b, c, d uint32) { return 0, 0, 0, 0 }
}

func addInfo(c *CPUInfo, safe bool) {
	// Seems to be safe to assume on ARM64
	c.CacheLine = 64
	detectOS(c)

	// ARM64 disabled since it may crash if interrupt is not intercepted by OS.
	if safe && !c.Has(ARMCPUID) && runtime.GOOS != "freebsd" {
		return
	}
	midr := getMidr()

	// MIDR_EL1 - Main ID Register
	// https://developer.arm.com/docs/ddi0595/h/aarch64-system-registers/midr_el1
	//  x--------------------------------------------------x
	//  | Name                         |  bits   | visible |
	//  |--------------------------------------------------|
	//  | Implementer                  | [31-24] |    y    |
	//  |--------------------------------------------------|
	//  | Variant                      | [23-20] |    y    |
	//  |--------------------------------------------------|
	//  | Architecture                 | [19-16] |    y    |
	//  |--------------------------------------------------|
	//  | PartNum                      | [15-4]  |    y    |
	//  |--------------------------------------------------|
	//  | Revision                     | [3-0]   |    y    |
	//  x--------------------------------------------------x

	switch (midr >> 24) & 0xff {
	case 0xC0:
		c.VendorString = "Ampere Computing"
		c.VendorID = Ampere
	case 0x41:
		c.VendorString = "Arm Limited"
		c.VendorID = ARM
	case 0x42:
		c.VendorString = "Broadcom Corporation"
		c.VendorID = Broadcom
	case 0x43:
		c.VendorString = "Cavium Inc"
		c.VendorID = Cavium
	case 0x44:
		c.VendorString = "Digital Equipment Corporation"
		c.VendorID = DEC
	case 0x46:
		c.VendorString = "Fujitsu Ltd"
		c.VendorID = Fujitsu
	case 0x49:
		c.VendorString = "Infineon Technologies AG"
		c.VendorID = Infineon
	case 0x4D:
		c.VendorString = "Motorola or Freescale Semiconductor Inc"
		c.VendorID = Motorola
	case 0x4E:
		c.VendorString = "NVIDIA Corporation"
		c.VendorID = NVIDIA
	case 0x50:
		c.VendorString = "Applied Micro Circuits Corporation"
		c.VendorID = AMCC
	case 0x51:
		c.VendorString = "Qualcomm Inc"
		c.VendorID = Qualcomm
	case 0x56:
		c.VendorString = "Marvell International Ltd"
		c.VendorID = Marvell
	case 0x69:
		c.VendorString = "Intel Corporation"
		c.VendorID = Intel
	}

	// Lower 4 bits: Architecture
	// Architecture	Meaning
	// 0b0001		Armv4.
	// 0b0010		Armv4T.
	// 0b0011		Armv5 (obsolete).
	// 0b0100		Armv5T.
	// 0b0101		Armv5TE.
	// 0b0110		Armv5TEJ.
	// 0b0111		Armv6.
	// 0b1111		Architectural features are individually identified in the ID_* registers, see 'ID registers'.
	// Upper 4 bit: Variant
	// An IMPLEMENTATION DEFINED variant number.
	// Typically, this field is used to distinguish between different product variants, or major revisions of a product.
	c.Family = int(midr>>16) & 0xff

	// PartNum, bits [15:4]
	// An IMPLEMENTATION DEFINED primary part number for the device.
	// On processors implemented by Arm, if the top four bits of the primary
	// part number are 0x0 or 0x7, the variant and architecture are encoded differently.
	// Revision, bits [3:0]
	// An IMPLEMENTATION DEFINED revision number for the device.
	c.Model = int(midr) & 0xffff

	procFeatures := getProcFeatures()

	// ID_AA64PFR0_EL1 - Processor Feature Register 0
	// x--------------------------------------------------x
	// | Name                         |  bits   | visible |
	// |--------------------------------------------------|
	// | DIT                          | [51-48] |    y    |
	// |--------------------------------------------------|
	// | SVE                          | [35-32] |    y    |
	// |--------------------------------------------------|
	// | GIC                          | [27-24] |    n    |
	// |--------------------------------------------------|
	// | AdvSIMD                      | [23-20] |    y    |
	// |--------------------------------------------------|
	// | FP                           | [19-16] |    y    |
	// |--------------------------------------------------|
	// | EL3                          | [15-12] |    n    |
	// |--------------------------------------------------|
	// | EL2                          | [11-8]  |    n    |
	// |--------------------------------------------------|
	// | EL1                          | [7-4]   |    n    |
	// |--------------------------------------------------|
	// | EL0                          | [3-0]   |    n    |
	// x--------------------------------------------------x

	var f flagSet
	// if procFeatures&(0xf<<48) != 0 {
	// 	fmt.Println("DIT")
	// }
	f.setIf(procFeatures&(0xf<<32) != 0, SVE)
	if procFeatures&(0xf<<20) != 15<<20 {
		f.set(ASIMD)
		// https://developer.arm.com/docs/ddi0595/b/aarch64-system-registers/id_aa64pfr0_el1
		// 0b0001 --> As for 0b0000, and also includes support for half-precision floating-point arithmetic.
		f.setIf(procFeatures&(0xf<<20) == 1<<20, FPHP, ASIMDHP)
	}
	f.setIf(procFeatures&(0xf<<16) != 0, FP)

	instAttrReg0, instAttrReg1 := getInstAttributes()

	// https://developer.arm.com/docs/ddi0595/b/aarch64-system-registers/id_aa64isar0_el1
	//
	// ID_AA64ISAR0_EL1 - Instruction Set Attribute Register 0
	// x--------------------------------------------------x
	// | Name                         |  bits   | visible |
	// |--------------------------------------------------|
	// | RNDR                         | [63-60] |    y    |
	// |--------------------------------------------------|
	// | TLB                          | [59-56] |    y    |
	// |--------------------------------------------------|
	// | TS                           | [55-52] |    y    |
	// |--------------------------------------------------|
	// | FHM                          | [51-48] |    y    |
	// |--------------------------------------------------|
	// | DP                           | [47-44] |    y    |
	// |--------------------------------------------------|
	// | SM4                          | [43-40] |    y    |
	// |--------------------------------------------------|
	// | SM3                          | [39-36] |    y    |
	// |--------------------------------------------------|
	// | SHA3                         | [35-32] |    y    |
	// |--------------------------------------------------|
	// | RDM                          | [31-28] |    y    |
	// |--------------------------------------------------|
	// | ATOMICS                      | [23-20] |    y    |
	// |--------------------------------------------------|
	// | CRC32                        | [19-16] |    y    |
	// |--------------------------------------------------|
	// | SHA2                         | [15-12] |    y    |
	// |--------------------------------------------------|
	// | SHA1                         | [11-8]  |    y    |
	// |--------------------------------------------------|
	// | AES                          | [7-4]   |    y    |
	// x--------------------------------------------------x

	f.setIf(instAttrReg0&(0xf<<60) != 0, RNDR)
	f.setIf(instAttrReg0&(0xf<<56) != 0, TLB)
	f.setIf(instAttrReg0&(0xf<<52) != 0, TS)
	f.setIf(instAttrReg0&(0xf<<48) != 0, FHM)
	f.setIf(instAttrReg0&(0xf<<44) != 0, ASIMDDP)
	f.setIf(instAttrReg0&(0xf<<40) != 0, SM4)
	f.setIf(instAttrReg0&(0xf<<36) != 0, SM3)
	f.setIf(instAttrReg0&(0xf<<32) != 0, SHA3)
	f.setIf(instAttrReg0&(0xf<<28) != 0, ASIMDRDM)
	f.setIf(instAttrReg0&(0xf<<20) != 0, ATOMICS)
	f.setIf(instAttrReg0&(0xf<<16) != 0, CRC32)
	f.setIf(instAttrReg0&(0xf<<12) != 0, SHA2)
	// https://developer.arm.com/docs/ddi0595/b/aarch64-system-registers/id_aa64isar0_el1
	// 0b0010 --> As 0b0001, plus SHA512H, SHA512H2, SHA512SU0, and SHA512SU1 instructions implemented.
	f.setIf(instAttrReg0&(0xf<<12) == 2<<12, SHA512)
	f.setIf(instAttrReg0&(0xf<<8) != 0, SHA1)
	f.setIf(instAttrReg0&(0xf<<4) != 0, AESARM)
	// https://developer.arm.com/docs/ddi0595/b/aarch64-system-registers/id_aa64isar0_el1
	// 0b0010 --> As for 0b0001, plus PMULL/PMULL2 instructions operating on 64-bit data quantities.
	f.setIf(instAttrReg0&(0xf<<4) == 2<<4, PMULL)

	// https://developer.arm.com/docs/ddi0595/b/aarch64-system-registers/id_aa64isar1_el1
	//
	// ID_AA64ISAR1_EL1 - Instruction set attribute register 1
	// x--------------------------------------------------x
	// | Name                         |  bits   | visible |
	// |--------------------------------------------------|
	// | GPI                          | [31-28] |    y    |
	// |--------------------------------------------------|
	// | GPA                          | [27-24] |    y    |
	// |--------------------------------------------------|
	// | LRCPC                        | [23-20] |    y    |
	// |--------------------------------------------------|
	// | FCMA                         | [19-16] |    y    |
	// |--------------------------------------------------|
	// | JSCVT                        | [15-12] |    y    |
	// |--------------------------------------------------|
	// | API                          | [11-8]  |    y    |
	// |--------------------------------------------------|
	// | APA                          | [7-4]   |    y    |
	// |--------------------------------------------------|
	// | DPB                          | [3-0]   |    y    |
	// x--------------------------------------------------x

	// if instAttrReg1&(0xf<<28) != 0 {
	// 	fmt.Println("GPI")
	// }
	f.setIf(instAttrReg1&(0xf<<28) != 24, GPA)
	f.setIf(instAttrReg1&(0xf<<20) != 0, LRCPC)
	f.setIf(instAttrReg1&(0xf<<16) != 0, FCMA)
	f.setIf(instAttrReg1&(0xf<<12) != 0, JSCVT)
	// if instAttrReg1&(0xf<<8) != 0 {
	// 	fmt.Println("API")
	// }
	// if instAttrReg1&(0xf<<4) != 0 {
	// 	fmt.Println("APA")
	// }
	f.setIf(instAttrReg1&(0xf<<0) != 0, DCPOP)

	// Store
	c.featureSet.or(f)
}
