// Copyright (c) 2015 Klaus Post, released under MIT License. See LICENSE file.

// Package cpuid provides information about the CPU running the current program.
//
// CPU features are detected on startup, and kept for fast access through the life of the application.
// Currently x86 / x64 (AMD64) as well as arm64 is supported.
//
// You can access the CPU information by accessing the shared CPU variable of the cpuid library.
//
// Package home: https://github.com/klauspost/cpuid
package cpuid

import (
	"flag"
	"fmt"
	"math"
	"math/bits"
	"os"
	"runtime"
	"strings"
)

// AMD refererence: https://www.amd.com/system/files/TechDocs/25481.pdf
// and Processor Programming Reference (PPR)

// Vendor is a representation of a CPU vendor.
type Vendor int

const (
	VendorUnknown Vendor = iota
	Intel
	AMD
	VIA
	Transmeta
	NSC
	KVM  // Kernel-based Virtual Machine
	MSVM // Microsoft Hyper-V or Windows Virtual PC
	VMware
	XenHVM
	Bhyve
	Hygon
	SiS
	RDC

	Ampere
	ARM
	Broadcom
	Cavium
	DEC
	Fujitsu
	Infineon
	Motorola
	NVIDIA
	AMCC
	Qualcomm
	Marvell

	QEMU
	QNX
	ACRN
	SRE
	Apple

	lastVendor
)

//go:generate stringer -type=FeatureID,Vendor

// FeatureID is the ID of a specific cpu feature.
type FeatureID int

const (
	// Keep index -1 as unknown
	UNKNOWN = -1

	// x86 features
	ADX                 FeatureID = iota // Intel ADX (Multi-Precision Add-Carry Instruction Extensions)
	AESNI                                // Advanced Encryption Standard New Instructions
	AMD3DNOW                             // AMD 3DNOW
	AMD3DNOWEXT                          // AMD 3DNowExt
	AMXBF16                              // Tile computational operations on BFLOAT16 numbers
	AMXFP16                              // Tile computational operations on FP16 numbers
	AMXINT8                              // Tile computational operations on 8-bit integers
	AMXFP8                               // Tile computational operations on FP8 numbers
	AMXTILE                              // Tile architecture
	AMXTF32                              // Tile architecture
	AMXCOMPLEX                           // Matrix Multiplication of TF32 Tiles into Packed Single Precision Tile
	AMXTRANSPOSE                         // Tile multiply where the first operand is transposed
	APX_F                                // Intel APX
	AVX                                  // AVX functions
	AVX10                                // If set the Intel AVX10 Converged Vector ISA is supported
	AVX10_128                            // If set indicates that AVX10 128-bit vector support is present
	AVX10_256                            // If set indicates that AVX10 256-bit vector support is present
	AVX10_512                            // If set indicates that AVX10 512-bit vector support is present
	AVX2                                 // AVX2 functions
	AVX512BF16                           // AVX-512 BFLOAT16 Instructions
	AVX512BITALG                         // AVX-512 Bit Algorithms
	AVX512BW                             // AVX-512 Byte and Word Instructions
	AVX512CD                             // AVX-512 Conflict Detection Instructions
	AVX512DQ                             // AVX-512 Doubleword and Quadword Instructions
	AVX512ER                             // AVX-512 Exponential and Reciprocal Instructions
	AVX512F                              // AVX-512 Foundation
	AVX512FP16                           // AVX-512 FP16 Instructions
	AVX512IFMA                           // AVX-512 Integer Fused Multiply-Add Instructions
	AVX512PF                             // AVX-512 Prefetch Instructions
	AVX512VBMI                           // AVX-512 Vector Bit Manipulation Instructions
	AVX512VBMI2                          // AVX-512 Vector Bit Manipulation Instructions, Version 2
	AVX512VL                             // AVX-512 Vector Length Extensions
	AVX512VNNI                           // AVX-512 Vector Neural Network Instructions
	AVX512VP2INTERSECT                   // AVX-512 Intersect for D/Q
	AVX512VPOPCNTDQ                      // AVX-512 Vector Population Count Doubleword and Quadword
	AVXIFMA                              // AVX-IFMA instructions
	AVXNECONVERT                         // AVX-NE-CONVERT instructions
	AVXSLOW                              // Indicates the CPU performs 2 128 bit operations instead of one
	AVXVNNI                              // AVX (VEX encoded) VNNI neural network instructions
	AVXVNNIINT8                          // AVX-VNNI-INT8 instructions
	AVXVNNIINT16                         // AVX-VNNI-INT16 instructions
	BHI_CTRL                             // Branch History Injection and Intra-mode Branch Target Injection / CVE-2022-0001, CVE-2022-0002 / INTEL-SA-00598
	BMI1                                 // Bit Manipulation Instruction Set 1
	BMI2                                 // Bit Manipulation Instruction Set 2
	CETIBT                               // Intel CET Indirect Branch Tracking
	CETSS                                // Intel CET Shadow Stack
	CLDEMOTE                             // Cache Line Demote
	CLMUL                                // Carry-less Multiplication
	CLZERO                               // CLZERO instruction supported
	CMOV                                 // i686 CMOV
	CMPCCXADD                            // CMPCCXADD instructions
	CMPSB_SCADBS_SHORT                   // Fast short CMPSB and SCASB
	CMPXCHG8                             // CMPXCHG8 instruction
	CPBOOST                              // Core Performance Boost
	CPPC                                 // AMD: Collaborative Processor Performance Control
	CX16                                 // CMPXCHG16B Instruction
	EFER_LMSLE_UNS                       // AMD: =Core::X86::Msr::EFER[LMSLE] is not supported, and MBZ
	ENQCMD                               // Enqueue Command
	ERMS                                 // Enhanced REP MOVSB/STOSB
	F16C                                 // Half-precision floating-point conversion
	FLUSH_L1D                            // Flush L1D cache
	FMA3                                 // Intel FMA 3. Does not imply AVX.
	FMA4                                 // Bulldozer FMA4 functions
	FP128                                // AMD: When set, the internal FP/SIMD execution datapath is no more than 128-bits wide
	FP256                                // AMD: When set, the internal FP/SIMD execution datapath is no more than 256-bits wide
	FSRM                                 // Fast Short Rep Mov
	FXSR                                 // FXSAVE, FXRESTOR instructions, CR4 bit 9
	FXSROPT                              // FXSAVE/FXRSTOR optimizations
	GFNI                                 // Galois Field New Instructions. May require other features (AVX, AVX512VL,AVX512F) based on usage.
	HLE                                  // Hardware Lock Elision
	HRESET                               // If set CPU supports history reset and the IA32_HRESET_ENABLE MSR
	HTT                                  // Hyperthreading (enabled)
	HWA                                  // Hardware assert supported. Indicates support for MSRC001_10
	HYBRID_CPU                           // This part has CPUs of more than one type.
	HYPERVISOR                           // This bit has been reserved by Intel & AMD for use by hypervisors
	IA32_ARCH_CAP                        // IA32_ARCH_CAPABILITIES MSR (Intel)
	IA32_CORE_CAP                        // IA32_CORE_CAPABILITIES MSR
	IBPB                                 // Indirect Branch Restricted Speculation (IBRS) and Indirect Branch Predictor Barrier (IBPB)
	IBPB_BRTYPE                          // Indicates that MSR 49h (PRED_CMD) bit 0 (IBPB) flushes	all branch type predictions from the CPU branch predictor
	IBRS                                 // AMD: Indirect Branch Restricted Speculation
	IBRS_PREFERRED                       // AMD: IBRS is preferred over software solution
	IBRS_PROVIDES_SMP                    // AMD: IBRS provides Same Mode Protection
	IBS                                  // Instruction Based Sampling (AMD)
	IBSBRNTRGT                           // Instruction Based Sampling Feature (AMD)
	IBSFETCHSAM                          // Instruction Based Sampling Feature (AMD)
	IBSFFV                               // Instruction Based Sampling Feature (AMD)
	IBSOPCNT                             // Instruction Based Sampling Feature (AMD)
	IBSOPCNTEXT                          // Instruction Based Sampling Feature (AMD)
	IBSOPSAM                             // Instruction Based Sampling Feature (AMD)
	IBSRDWROPCNT                         // Instruction Based Sampling Feature (AMD)
	IBSRIPINVALIDCHK                     // Instruction Based Sampling Feature (AMD)
	IBS_FETCH_CTLX                       // AMD: IBS fetch control extended MSR supported
	IBS_OPDATA4                          // AMD: IBS op data 4 MSR supported
	IBS_OPFUSE                           // AMD: Indicates support for IbsOpFuse
	IBS_PREVENTHOST                      // Disallowing IBS use by the host supported
	IBS_ZEN4                             // AMD: Fetch and Op IBS support IBS extensions added with Zen4
	IDPRED_CTRL                          // IPRED_DIS
	INT_WBINVD                           // WBINVD/WBNOINVD are interruptible.
	INVLPGB                              // NVLPGB and TLBSYNC instruction supported
	KEYLOCKER                            // Key locker
	KEYLOCKERW                           // Key locker wide
	LAHF                                 // LAHF/SAHF in long mode
	LAM                                  // If set, CPU supports Linear Address Masking
	LBRVIRT                              // LBR virtualization
	LZCNT                                // LZCNT instruction
	MCAOVERFLOW                          // MCA overflow recovery support.
	MCDT_NO                              // Processor do not exhibit MXCSR Configuration Dependent Timing behavior and do not need to mitigate it.
	MCOMMIT                              // MCOMMIT instruction supported
	MD_CLEAR                             // VERW clears CPU buffers
	MMX                                  // standard MMX
	MMXEXT                               // SSE integer functions or AMD MMX ext
	MOVBE                                // MOVBE instruction (big-endian)
	MOVDIR64B                            // Move 64 Bytes as Direct Store
	MOVDIRI                              // Move Doubleword as Direct Store
	MOVSB_ZL                             // Fast Zero-Length MOVSB
	MOVU                                 // AMD: MOVU SSE instructions are more efficient and should be preferred to SSE	MOVL/MOVH. MOVUPS is more efficient than MOVLPS/MOVHPS. MOVUPD is more efficient than MOVLPD/MOVHPD
	MPX                                  // Intel MPX (Memory Protection Extensions)
	MSRIRC                               // Instruction Retired Counter MSR available
	MSRLIST                              // Read/Write List of Model Specific Registers
	MSR_PAGEFLUSH                        // Page Flush MSR available
	NRIPS                                // Indicates support for NRIP save on VMEXIT
	NX                                   // NX (No-Execute) bit
	OSXSAVE                              // XSAVE enabled by OS
	PCONFIG                              // PCONFIG for Intel Multi-Key Total Memory Encryption
	POPCNT                               // POPCNT instruction
	PPIN                                 // AMD: Protected Processor Inventory Number support. Indicates that Protected Processor Inventory Number (PPIN) capability can be enabled
	PREFETCHI                            // PREFETCHIT0/1 instructions
	PSFD                                 // Predictive Store Forward Disable
	RDPRU                                // RDPRU instruction supported
	RDRAND                               // RDRAND instruction is available
	RDSEED                               // RDSEED instruction is available
	RDTSCP                               // RDTSCP Instruction
	RRSBA_CTRL                           // Restricted RSB Alternate
	RTM                                  // Restricted Transactional Memory
	RTM_ALWAYS_ABORT                     // Indicates that the loaded microcode is forcing RTM abort.
	SBPB                                 // Indicates support for the Selective Branch Predictor Barrier
	SERIALIZE                            // Serialize Instruction Execution
	SEV                                  // AMD Secure Encrypted Virtualization supported
	SEV_64BIT                            // AMD SEV guest execution only allowed from a 64-bit host
	SEV_ALTERNATIVE                      // AMD SEV Alternate Injection supported
	SEV_DEBUGSWAP                        // Full debug state swap supported for SEV-ES guests
	SEV_ES                               // AMD SEV Encrypted State supported
	SEV_RESTRICTED                       // AMD SEV Restricted Injection supported
	SEV_SNP                              // AMD SEV Secure Nested Paging supported
	SGX                                  // Software Guard Extensions
	SGXLC                                // Software Guard Extensions Launch Control
	SGXPQC                               // Software Guard Extensions 256-bit Encryption
	SHA                                  // Intel SHA Extensions
	SME                                  // AMD Secure Memory Encryption supported
	SME_COHERENT                         // AMD Hardware cache coherency across encryption domains enforced
	SM3_X86                              // SM3 instructions
	SM4_X86                              // SM4 instructions
	SPEC_CTRL_SSBD                       // Speculative Store Bypass Disable
	SRBDS_CTRL                           // SRBDS mitigation MSR available
	SRSO_MSR_FIX                         // Indicates that software may use MSR BP_CFG[BpSpecReduce] to mitigate SRSO.
	SRSO_NO                              // Indicates the CPU is not subject to the SRSO vulnerability
	SRSO_USER_KERNEL_NO                  // Indicates the CPU is not subject to the SRSO vulnerability across user/kernel boundaries
	SSE                                  // SSE functions
	SSE2                                 // P4 SSE functions
	SSE3                                 // Prescott SSE3 functions
	SSE4                                 // Penryn SSE4.1 functions
	SSE42                                // Nehalem SSE4.2 functions
	SSE4A                                // AMD Barcelona microarchitecture SSE4a instructions
	SSSE3                                // Conroe SSSE3 functions
	STIBP                                // Single Thread Indirect Branch Predictors
	STIBP_ALWAYSON                       // AMD: Single Thread Indirect Branch Prediction Mode has Enhanced Performance and may be left Always On
	STOSB_SHORT                          // Fast short STOSB
	SUCCOR                               // Software uncorrectable error containment and recovery capability.
	SVM                                  // AMD Secure Virtual Machine
	SVMDA                                // Indicates support for the SVM decode assists.
	SVMFBASID                            // SVM, Indicates that TLB flush events, including CR3 writes and CR4.PGE toggles, flush only the current ASID's TLB entries. Also indicates support for the extended VMCBTLB_Control
	SVML                                 // AMD SVM lock. Indicates support for SVM-Lock.
	SVMNP                                // AMD SVM nested paging
	SVMPF                                // SVM pause intercept filter. Indicates support for the pause intercept filter
	SVMPFT                               // SVM PAUSE filter threshold. Indicates support for the PAUSE filter cycle count threshold
	SYSCALL                              // System-Call Extension (SCE): SYSCALL and SYSRET instructions.
	SYSEE                                // SYSENTER and SYSEXIT instructions
	TBM                                  // AMD Trailing Bit Manipulation
	TDX_GUEST                            // Intel Trust Domain Extensions Guest
	TLB_FLUSH_NESTED                     // AMD: Flushing includes all the nested translations for guest translations
	TME                                  // Intel Total Memory Encryption. The following MSRs are supported: IA32_TME_CAPABILITY, IA32_TME_ACTIVATE, IA32_TME_EXCLUDE_MASK, and IA32_TME_EXCLUDE_BASE.
	TOPEXT                               // TopologyExtensions: topology extensions support. Indicates support for CPUID Fn8000_001D_EAX_x[N:0]-CPUID Fn8000_001E_EDX.
	TSA_L1_NO                            // AMD only: Not vulnerable to TSA-L1
	TSA_SQ_NO                            // AM onlyD: Not vulnerable to TSA-SQ
	TSA_VERW_CLEAR                       // If set, the memory form of the VERW instruction may be used to help mitigate TSA
	TSCRATEMSR                           // MSR based TSC rate control. Indicates support for MSR TSC ratio MSRC000_0104
	TSXLDTRK                             // Intel TSX Suspend Load Address Tracking
	VAES                                 // Vector AES. AVX(512) versions requires additional checks.
	VMCBCLEAN                            // VMCB clean bits. Indicates support for VMCB clean bits.
	VMPL                                 // AMD VM Permission Levels supported
	VMSA_REGPROT                         // AMD VMSA Register Protection supported
	VMX                                  // Virtual Machine Extensions
	VPCLMULQDQ                           // Carry-Less Multiplication Quadword. Requires AVX for 3 register versions.
	VTE                                  // AMD Virtual Transparent Encryption supported
	WAITPKG                              // TPAUSE, UMONITOR, UMWAIT
	WBNOINVD                             // Write Back and Do Not Invalidate Cache
	WRMSRNS                              // Non-Serializing Write to Model Specific Register
	X87                                  // FPU
	XGETBV1                              // Supports XGETBV with ECX = 1
	XOP                                  // Bulldozer XOP functions
	XSAVE                                // XSAVE, XRESTOR, XSETBV, XGETBV
	XSAVEC                               // Supports XSAVEC and the compacted form of XRSTOR.
	XSAVEOPT                             // XSAVEOPT available
	XSAVES                               // Supports XSAVES/XRSTORS and IA32_XSS

	// ARM features:
	AESARM   // AES instructions
	ARMCPUID // Some CPU ID registers readable at user-level
	ASIMD    // Advanced SIMD
	ASIMDDP  // SIMD Dot Product
	ASIMDHP  // Advanced SIMD half-precision floating point
	ASIMDRDM // Rounding Double Multiply Accumulate/Subtract (SQRDMLAH/SQRDMLSH)
	ATOMICS  // Large System Extensions (LSE)
	CRC32    // CRC32/CRC32C instructions
	DCPOP    // Data cache clean to Point of Persistence (DC CVAP)
	EVTSTRM  // Generic timer
	FCMA     // Floating point complex number addition and multiplication
	FHM      // FMLAL and FMLSL instructions
	FP       // Single-precision and double-precision floating point
	FPHP     // Half-precision floating point
	GPA      // Generic Pointer Authentication
	JSCVT    // Javascript-style double->int convert (FJCVTZS)
	LRCPC    // Weaker release consistency (LDAPR, etc)
	PMULL    // Polynomial Multiply instructions (PMULL/PMULL2)
	RNDR     // Random Number instructions
	TLB      // Outer Shareable and TLB range maintenance instructions
	TS       // Flag manipulation instructions
	SHA1     // SHA-1 instructions (SHA1C, etc)
	SHA2     // SHA-2 instructions (SHA256H, etc)
	SHA3     // SHA-3 instructions (EOR3, RAXI, XAR, BCAX)
	SHA512   // SHA512 instructions
	SM3      // SM3 instructions
	SM4      // SM4 instructions
	SVE      // Scalable Vector Extension

	// PMU
	PMU_FIXEDCOUNTER_CYCLES
	PMU_FIXEDCOUNTER_REFCYCLES
	PMU_FIXEDCOUNTER_INSTRUCTIONS
	PMU_FIXEDCOUNTER_TOPDOWN_SLOTS

	// Keep it last. It automatically defines the size of []flagSet
	lastID

	firstID FeatureID = UNKNOWN + 1
)

// CPUInfo contains information about the detected system CPU.
type CPUInfo struct {
	BrandName              string  // Brand name reported by the CPU
	VendorID               Vendor  // Comparable CPU vendor ID
	VendorString           string  // Raw vendor string.
	HypervisorVendorID     Vendor  // Hypervisor vendor
	HypervisorVendorString string  // Raw hypervisor vendor string
	featureSet             flagSet // Features of the CPU
	PhysicalCores          int     // Number of physical processor cores in your CPU. Will be 0 if undetectable.
	ThreadsPerCore         int     // Number of threads per physical core. Will be 1 if undetectable.
	LogicalCores           int     // Number of physical cores times threads that can run on each core through the use of hyperthreading. Will be 0 if undetectable.
	Family                 int     // CPU family number
	Model                  int     // CPU model number
	Stepping               int     // CPU stepping info
	CacheLine              int     // Cache line size in bytes. Will be 0 if undetectable.
	Hz                     int64   // Clock speed, if known, 0 otherwise. Will attempt to contain base clock speed.
	BoostFreq              int64   // Max clock speed, if known, 0 otherwise
	Cache                  struct {
		L1I int // L1 Instruction Cache (per core or shared). Will be -1 if undetected
		L1D int // L1 Data Cache (per core or shared). Will be -1 if undetected
		L2  int // L2 Cache (per core or shared). Will be -1 if undetected
		L3  int // L3 Cache (per core, per ccx or shared). Will be -1 if undetected
	}
	SGX              SGXSupport
	AMDMemEncryption AMDMemEncryptionSupport
	AVX10Level       uint8
	PMU              PerformanceMonitoringInfo //  holds information about the PMU

	maxFunc   uint32
	maxExFunc uint32
}

// PerformanceMonitoringInfo holds information about CPU performance monitoring capabilities.
// This is primarily populated from CPUID leaf 0xAh on x86
type PerformanceMonitoringInfo struct {
	// VersionID (x86 only): Version ID of architectural performance monitoring.
	// A value of 0 means architectural performance monitoring is not supported or information is unavailable.
	VersionID uint8
	// NumGPPMC: Number of General-Purpose Performance Monitoring Counters per logical processor.
	// On ARM, this is derived from PMCR_EL0.N (number of event counters).
	NumGPCounters uint8
	// GPPMCWidth: Bit width of General-Purpose Performance Monitoring Counters.
	// On ARM, typically 64 for PMU event counters.
	GPPMCWidth uint8
	// NumFixedPMC: Number of Fixed-Function Performance Counters.
	// Valid on x86 if VersionID > 1. On ARM, this typically includes at least the cycle counter (PMCCNTR_EL0).
	NumFixedPMC uint8
	// FixedPMCWidth: Bit width of Fixed-Function Performance Counters.
	// Valid on x86 if VersionID > 1. On ARM, the cycle counter (PMCCNTR_EL0) is 64-bit.
	FixedPMCWidth uint8
	// Raw register output from CPUID leaf 0xAh.
	RawEBX uint32
	RawEAX uint32
	RawEDX uint32
}

var cpuid func(op uint32) (eax, ebx, ecx, edx uint32)
var cpuidex func(op, op2 uint32) (eax, ebx, ecx, edx uint32)
var xgetbv func(index uint32) (eax, edx uint32)
var rdtscpAsm func() (eax, ebx, ecx, edx uint32)
var darwinHasAVX512 = func() bool { return false }

// CPU contains information about the CPU as detected on startup,
// or when Detect last was called.
//
// Use this as the primary entry point to you data.
var CPU CPUInfo

func init() {
	initCPU()
	Detect()
}

// Detect will re-detect current CPU info.
// This will replace the content of the exported CPU variable.
//
// Unless you expect the CPU to change while you are running your program
// you should not need to call this function.
// If you call this, you must ensure that no other goroutine is accessing the
// exported CPU variable.
func Detect() {
	// Set defaults
	CPU.ThreadsPerCore = 1
	CPU.Cache.L1I = -1
	CPU.Cache.L1D = -1
	CPU.Cache.L2 = -1
	CPU.Cache.L3 = -1
	safe := true
	if detectArmFlag != nil {
		safe = !*detectArmFlag
	}
	addInfo(&CPU, safe)
	if displayFeats != nil && *displayFeats {
		fmt.Println("cpu features:", strings.Join(CPU.FeatureSet(), ","))
		// Exit with non-zero so tests will print value.
		os.Exit(1)
	}
	if disableFlag != nil {
		s := strings.Split(*disableFlag, ",")
		for _, feat := range s {
			feat := ParseFeature(strings.TrimSpace(feat))
			if feat != UNKNOWN {
				CPU.featureSet.unset(feat)
			}
		}
	}
}

// DetectARM will detect ARM64 features.
// This is NOT done automatically since it can potentially crash
// if the OS does not handle the command.
// If in the future this can be done safely this function may not
// do anything.
func DetectARM() {
	addInfo(&CPU, false)
}

var detectArmFlag *bool
var displayFeats *bool
var disableFlag *string

// Flags will enable flags.
// This must be called *before* flag.Parse AND
// Detect must be called after the flags have been parsed.
// Note that this means that any detection used in init() functions
// will not contain these flags.
func Flags() {
	disableFlag = flag.String("cpu.disable", "", "disable cpu features; comma separated list")
	displayFeats = flag.Bool("cpu.features", false, "lists cpu features and exits")
	detectArmFlag = flag.Bool("cpu.arm", false, "allow ARM features to be detected; can potentially crash")
}

// Supports returns whether the CPU supports all of the requested features.
func (c CPUInfo) Supports(ids ...FeatureID) bool {
	for _, id := range ids {
		if !c.featureSet.inSet(id) {
			return false
		}
	}
	return true
}

// Has allows for checking a single feature.
// Should be inlined by the compiler.
func (c *CPUInfo) Has(id FeatureID) bool {
	return c.featureSet.inSet(id)
}

// AnyOf returns whether the CPU supports one or more of the requested features.
func (c CPUInfo) AnyOf(ids ...FeatureID) bool {
	for _, id := range ids {
		if c.featureSet.inSet(id) {
			return true
		}
	}
	return false
}

// Features contains several features combined for a fast check using
// CpuInfo.HasAll
type Features *flagSet

// CombineFeatures allows to combine several features for a close to constant time lookup.
func CombineFeatures(ids ...FeatureID) Features {
	var v flagSet
	for _, id := range ids {
		v.set(id)
	}
	return &v
}

func (c *CPUInfo) HasAll(f Features) bool {
	return c.featureSet.hasSetP(f)
}

// https://en.wikipedia.org/wiki/X86-64#Microarchitecture_levels
var oneOfLevel = CombineFeatures(SYSEE, SYSCALL)
var level1Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2)
var level2Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2, CX16, LAHF, POPCNT, SSE3, SSE4, SSE42, SSSE3)
var level3Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2, CX16, LAHF, POPCNT, SSE3, SSE4, SSE42, SSSE3, AVX, AVX2, BMI1, BMI2, F16C, FMA3, LZCNT, MOVBE, OSXSAVE)
var level4Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2, CX16, LAHF, POPCNT, SSE3, SSE4, SSE42, SSSE3, AVX, AVX2, BMI1, BMI2, F16C, FMA3, LZCNT, MOVBE, OSXSAVE, AVX512F, AVX512BW, AVX512CD, AVX512DQ, AVX512VL)

// X64Level returns the microarchitecture level detected on the CPU.
// If features are lacking or non x64 mode, 0 is returned.
// See https://en.wikipedia.org/wiki/X86-64#Microarchitecture_levels
func (c CPUInfo) X64Level() int {
	if !c.featureSet.hasOneOf(oneOfLevel) {
		return 0
	}
	if c.featureSet.hasSetP(level4Features) {
		return 4
	}
	if c.featureSet.hasSetP(level3Features) {
		return 3
	}
	if c.featureSet.hasSetP(level2Features) {
		return 2
	}
	if c.featureSet.hasSetP(level1Features) {
		return 1
	}
	return 0
}

// Disable will disable one or several features.
func (c *CPUInfo) Disable(ids ...FeatureID) bool {
	for _, id := range ids {
		c.featureSet.unset(id)
	}
	return true
}

// Enable will disable one or several features even if they were undetected.
// This is of course not recommended for obvious reasons.
func (c *CPUInfo) Enable(ids ...FeatureID) bool {
	for _, id := range ids {
		c.featureSet.set(id)
	}
	return true
}

// IsVendor returns true if vendor is recognized as Intel
func (c CPUInfo) IsVendor(v Vendor) bool {
	return c.VendorID == v
}

// FeatureSet returns all available features as strings.
func (c CPUInfo) FeatureSet() []string {
	s := make([]string, 0, c.featureSet.nEnabled())
	s = append(s, c.featureSet.Strings()...)
	return s
}

// RTCounter returns the 64-bit time-stamp counter
// Uses the RDTSCP instruction. The value 0 is returned
// if the CPU does not support the instruction.
func (c CPUInfo) RTCounter() uint64 {
	if !c.Has(RDTSCP) {
		return 0
	}
	a, _, _, d := rdtscpAsm()
	return uint64(a) | (uint64(d) << 32)
}

// Ia32TscAux returns the IA32_TSC_AUX part of the RDTSCP.
// This variable is OS dependent, but on Linux contains information
// about the current cpu/core the code is running on.
// If the RDTSCP instruction isn't supported on the CPU, the value 0 is returned.
func (c CPUInfo) Ia32TscAux() uint32 {
	if !c.Has(RDTSCP) {
		return 0
	}
	_, _, ecx, _ := rdtscpAsm()
	return ecx
}

// SveLengths returns arm SVE vector and predicate lengths in bits.
// Will return 0, 0 if SVE is not enabled or otherwise unable to detect.
func (c CPUInfo) SveLengths() (vl, pl uint64) {
	if !c.Has(SVE) {
		return 0, 0
	}
	return getVectorLength()
}

// LogicalCPU will return the Logical CPU the code is currently executing on.
// This is likely to change when the OS re-schedules the running thread
// to another CPU.
// If the current core cannot be detected, -1 will be returned.
func (c CPUInfo) LogicalCPU() int {
	if c.maxFunc < 1 {
		return -1
	}
	_, ebx, _, _ := cpuid(1)
	return int(ebx >> 24)
}

// frequencies tries to compute the clock speed of the CPU. If leaf 15 is
// supported, use it, otherwise parse the brand string. Yes, really.
func (c *CPUInfo) frequencies() {
	c.Hz, c.BoostFreq = 0, 0
	mfi := maxFunctionID()
	if mfi >= 0x15 {
		eax, ebx, ecx, _ := cpuid(0x15)
		if eax != 0 && ebx != 0 && ecx != 0 {
			c.Hz = (int64(ecx) * int64(ebx)) / int64(eax)
		}
	}
	if mfi >= 0x16 {
		a, b, _, _ := cpuid(0x16)
		// Base...
		if a&0xffff > 0 {
			c.Hz = int64(a&0xffff) * 1_000_000
		}
		// Boost...
		if b&0xffff > 0 {
			c.BoostFreq = int64(b&0xffff) * 1_000_000
		}
	}
	if c.Hz > 0 {
		return
	}

	// computeHz determines the official rated speed of a CPU from its brand
	// string. This insanity is *actually the official documented way to do
	// this according to Intel*, prior to leaf 0x15 existing. The official
	// documentation only shows this working for exactly `x.xx` or `xxxx`
	// cases, e.g., `2.50GHz` or `1300MHz`; this parser will accept other
	// sizes.
	model := c.BrandName
	hz := strings.LastIndex(model, "Hz")
	if hz < 3 {
		return
	}
	var multiplier int64
	switch model[hz-1] {
	case 'M':
		multiplier = 1000 * 1000
	case 'G':
		multiplier = 1000 * 1000 * 1000
	case 'T':
		multiplier = 1000 * 1000 * 1000 * 1000
	}
	if multiplier == 0 {
		return
	}
	freq := int64(0)
	divisor := int64(0)
	decimalShift := int64(1)
	var i int
	for i = hz - 2; i >= 0 && model[i] != ' '; i-- {
		if model[i] >= '0' && model[i] <= '9' {
			freq += int64(model[i]-'0') * decimalShift
			decimalShift *= 10
		} else if model[i] == '.' {
			if divisor != 0 {
				return
			}
			divisor = decimalShift
		} else {
			return
		}
	}
	// we didn't find a space
	if i < 0 {
		return
	}
	if divisor != 0 {
		c.Hz = (freq * multiplier) / divisor
		return
	}
	c.Hz = freq * multiplier
}

// VM Will return true if the cpu id indicates we are in
// a virtual machine.
func (c CPUInfo) VM() bool {
	return CPU.featureSet.inSet(HYPERVISOR)
}

// flags contains detected cpu features and characteristics
type flags uint64

// log2(bits_in_uint64)
const flagBitsLog2 = 6
const flagBits = 1 << flagBitsLog2
const flagMask = flagBits - 1

// flagSet contains detected cpu features and characteristics in an array of flags
type flagSet [(lastID + flagMask) / flagBits]flags

func (s *flagSet) inSet(feat FeatureID) bool {
	return s[feat>>flagBitsLog2]&(1<<(feat&flagMask)) != 0
}

func (s *flagSet) set(feat FeatureID) {
	s[feat>>flagBitsLog2] |= 1 << (feat & flagMask)
}

// setIf will set a feature if boolean is true.
func (s *flagSet) setIf(cond bool, features ...FeatureID) {
	if cond {
		for _, offset := range features {
			s[offset>>flagBitsLog2] |= 1 << (offset & flagMask)
		}
	}
}

func (s *flagSet) unset(offset FeatureID) {
	bit := flags(1 << (offset & flagMask))
	s[offset>>flagBitsLog2] = s[offset>>flagBitsLog2] & ^bit
}

// or with another flagset.
func (s *flagSet) or(other flagSet) {
	for i, v := range other[:] {
		s[i] |= v
	}
}

// hasSet returns whether all features are present.
func (s *flagSet) hasSet(other flagSet) bool {
	for i, v := range other[:] {
		if s[i]&v != v {
			return false
		}
	}
	return true
}

// hasSet returns whether all features are present.
func (s *flagSet) hasSetP(other *flagSet) bool {
	for i, v := range other[:] {
		if s[i]&v != v {
			return false
		}
	}
	return true
}

// hasOneOf returns whether one or more features are present.
func (s *flagSet) hasOneOf(other *flagSet) bool {
	for i, v := range other[:] {
		if s[i]&v != 0 {
			return true
		}
	}
	return false
}

// nEnabled will return the number of enabled flags.
func (s *flagSet) nEnabled() (n int) {
	for _, v := range s[:] {
		n += bits.OnesCount64(uint64(v))
	}
	return n
}

func flagSetWith(feat ...FeatureID) flagSet {
	var res flagSet
	for _, f := range feat {
		res.set(f)
	}
	return res
}

// ParseFeature will parse the string and return the ID of the matching feature.
// Will return UNKNOWN if not found.
func ParseFeature(s string) FeatureID {
	s = strings.ToUpper(s)
	for i := firstID; i < lastID; i++ {
		if i.String() == s {
			return i
		}
	}
	return UNKNOWN
}

// Strings returns an array of the detected features for FlagsSet.
func (s flagSet) Strings() []string {
	if len(s) == 0 {
		return []string{""}
	}
	r := make([]string, 0)
	for i := firstID; i < lastID; i++ {
		if s.inSet(i) {
			r = append(r, i.String())
		}
	}
	return r
}

func maxExtendedFunction() uint32 {
	eax, _, _, _ := cpuid(0x80000000)
	return eax
}

func maxFunctionID() uint32 {
	a, _, _, _ := cpuid(0)
	return a
}

func brandName() string {
	if maxExtendedFunction() >= 0x80000004 {
		v := make([]uint32, 0, 48)
		for i := uint32(0); i < 3; i++ {
			a, b, c, d := cpuid(0x80000002 + i)
			v = append(v, a, b, c, d)
		}
		return strings.Trim(string(valAsString(v...)), " ")
	}
	return "unknown"
}

func threadsPerCore() int {
	mfi := maxFunctionID()
	vend, _ := vendorID()

	if mfi < 0x4 || (vend != Intel && vend != AMD) {
		return 1
	}

	if mfi < 0xb {
		if vend != Intel {
			return 1
		}
		_, b, _, d := cpuid(1)
		if (d & (1 << 28)) != 0 {
			// v will contain logical core count
			v := (b >> 16) & 255
			if v > 1 {
				a4, _, _, _ := cpuid(4)
				// physical cores
				v2 := (a4 >> 26) + 1
				if v2 > 0 {
					return int(v) / int(v2)
				}
			}
		}
		return 1
	}
	_, b, _, _ := cpuidex(0xb, 0)
	if b&0xffff == 0 {
		if vend == AMD {
			// if >= Zen 2 0x8000001e EBX 15-8 bits means threads per core.
			// The number of threads per core is ThreadsPerCore+1
			// See PPR for AMD Family 17h Models 00h-0Fh (page 82)
			fam, _, _ := familyModel()
			_, _, _, d := cpuid(1)
			if (d&(1<<28)) != 0 && fam >= 23 {
				if maxExtendedFunction() >= 0x8000001e {
					_, b, _, _ := cpuid(0x8000001e)
					return int((b>>8)&0xff) + 1
				}
				return 2
			}
		}
		return 1
	}
	return int(b & 0xffff)
}

func logicalCores() int {
	mfi := maxFunctionID()
	v, _ := vendorID()
	switch v {
	case Intel:
		// Use this on old Intel processors
		if mfi < 0xb {
			if mfi < 1 {
				return 0
			}
			// CPUID.1:EBX[23:16] represents the maximum number of addressable IDs (initial APIC ID)
			// that can be assigned to logical processors in a physical package.
			// The value may not be the same as the number of logical processors that are present in the hardware of a physical package.
			_, ebx, _, _ := cpuid(1)
			logical := (ebx >> 16) & 0xff
			return int(logical)
		}
		_, b, _, _ := cpuidex(0xb, 1)
		return int(b & 0xffff)
	case AMD, Hygon:
		_, b, _, _ := cpuid(1)
		return int((b >> 16) & 0xff)
	default:
		return 0
	}
}

func familyModel() (family, model, stepping int) {
	if maxFunctionID() < 0x1 {
		return 0, 0, 0
	}
	eax, _, _, _ := cpuid(1)
	// If BaseFamily[3:0] is less than Fh then ExtendedFamily[7:0] is reserved and Family is equal to BaseFamily[3:0].
	family = int((eax >> 8) & 0xf)
	extFam := family == 0x6 // Intel is 0x6, needs extended model.
	if family == 0xf {
		// Add ExtFamily
		family += int((eax >> 20) & 0xff)
		extFam = true
	}
	// If BaseFamily[3:0] is less than 0Fh then ExtendedModel[3:0] is reserved and Model is equal to BaseModel[3:0].
	model = int((eax >> 4) & 0xf)
	if extFam {
		// Add ExtModel
		model += int((eax >> 12) & 0xf0)
	}
	stepping = int(eax & 0xf)
	return family, model, stepping
}

func physicalCores() int {
	v, _ := vendorID()
	switch v {
	case Intel:
		lc := logicalCores()
		tpc := threadsPerCore()
		if lc > 0 && tpc > 0 {
			return lc / tpc
		}
		return 0
	case AMD, Hygon:
		lc := logicalCores()
		tpc := threadsPerCore()
		if lc > 0 && tpc > 0 {
			return lc / tpc
		}

		// The following is inaccurate on AMD EPYC 7742 64-Core Processor
		if maxExtendedFunction() >= 0x80000008 {
			_, _, c, _ := cpuid(0x80000008)
			if c&0xff > 0 {
				return int(c&0xff) + 1
			}
		}
	}
	return 0
}

// Except from http://en.wikipedia.org/wiki/CPUID#EAX.3D0:_Get_vendor_ID
var vendorMapping = map[string]Vendor{
	"AMDisbetter!": AMD,
	"AuthenticAMD": AMD,
	"CentaurHauls": VIA,
	"GenuineIntel": Intel,
	"TransmetaCPU": Transmeta,
	"GenuineTMx86": Transmeta,
	"Geode by NSC": NSC,
	"VIA VIA VIA ": VIA,
	"KVMKVMKVM":    KVM,
	"Linux KVM Hv": KVM,
	"TCGTCGTCGTCG": QEMU,
	"Microsoft Hv": MSVM,
	"VMwareVMware": VMware,
	"XenVMMXenVMM": XenHVM,
	"bhyve bhyve ": Bhyve,
	"HygonGenuine": Hygon,
	"Vortex86 SoC": SiS,
	"SiS SiS SiS ": SiS,
	"RiseRiseRise": SiS,
	"Genuine  RDC": RDC,
	"QNXQVMBSQG":   QNX,
	"ACRNACRNACRN": ACRN,
	"SRESRESRESRE": SRE,
	"Apple VZ":     Apple,
}

func vendorID() (Vendor, string) {
	_, b, c, d := cpuid(0)
	v := string(valAsString(b, d, c))
	vend, ok := vendorMapping[v]
	if !ok {
		return VendorUnknown, v
	}
	return vend, v
}

func hypervisorVendorID() (Vendor, string) {
	// https://lwn.net/Articles/301888/
	_, b, c, d := cpuid(0x40000000)
	v := string(valAsString(b, c, d))
	vend, ok := vendorMapping[v]
	if !ok {
		return VendorUnknown, v
	}
	return vend, v
}

func cacheLine() int {
	if maxFunctionID() < 0x1 {
		return 0
	}

	_, ebx, _, _ := cpuid(1)
	cache := (ebx & 0xff00) >> 5 // cflush size
	if cache == 0 && maxExtendedFunction() >= 0x80000006 {
		_, _, ecx, _ := cpuid(0x80000006)
		cache = ecx & 0xff // cacheline size
	}
	// TODO: Read from Cache and TLB Information
	return int(cache)
}

func (c *CPUInfo) cacheSize() {
	c.Cache.L1D = -1
	c.Cache.L1I = -1
	c.Cache.L2 = -1
	c.Cache.L3 = -1
	vendor, _ := vendorID()
	switch vendor {
	case Intel:
		if maxFunctionID() < 4 {
			return
		}
		c.Cache.L1I, c.Cache.L1D, c.Cache.L2, c.Cache.L3 = 0, 0, 0, 0
		for i := uint32(0); ; i++ {
			eax, ebx, ecx, _ := cpuidex(4, i)
			cacheType := eax & 15
			if cacheType == 0 {
				break
			}
			cacheLevel := (eax >> 5) & 7
			coherency := int(ebx&0xfff) + 1
			partitions := int((ebx>>12)&0x3ff) + 1
			associativity := int((ebx>>22)&0x3ff) + 1
			sets := int(ecx) + 1
			size := associativity * partitions * coherency * sets
			switch cacheLevel {
			case 1:
				if cacheType == 1 {
					// 1 = Data Cache
					c.Cache.L1D = size
				} else if cacheType == 2 {
					// 2 = Instruction Cache
					c.Cache.L1I = size
				} else {
					if c.Cache.L1D < 0 {
						c.Cache.L1I = size
					}
					if c.Cache.L1I < 0 {
						c.Cache.L1I = size
					}
				}
			case 2:
				c.Cache.L2 = size
			case 3:
				c.Cache.L3 = size
			}
		}
	case AMD, Hygon:
		// Untested.
		if maxExtendedFunction() < 0x80000005 {
			return
		}
		_, _, ecx, edx := cpuid(0x80000005)
		c.Cache.L1D = int(((ecx >> 24) & 0xFF) * 1024)
		c.Cache.L1I = int(((edx >> 24) & 0xFF) * 1024)

		if maxExtendedFunction() < 0x80000006 {
			return
		}
		_, _, ecx, _ = cpuid(0x80000006)
		c.Cache.L2 = int(((ecx >> 16) & 0xFFFF) * 1024)

		// CPUID Fn8000_001D_EAX_x[N:0] Cache Properties
		if maxExtendedFunction() < 0x8000001D || !c.Has(TOPEXT) {
			return
		}

		// Xen Hypervisor is buggy and returns the same entry no matter ECX value.
		// Hack: When we encounter the same entry 100 times we break.
		nSame := 0
		var last uint32
		for i := uint32(0); i < math.MaxUint32; i++ {
			eax, ebx, ecx, _ := cpuidex(0x8000001D, i)

			level := (eax >> 5) & 7
			cacheNumSets := ecx + 1
			cacheLineSize := 1 + (ebx & 2047)
			cachePhysPartitions := 1 + ((ebx >> 12) & 511)
			cacheNumWays := 1 + ((ebx >> 22) & 511)

			typ := eax & 15
			size := int(cacheNumSets * cacheLineSize * cachePhysPartitions * cacheNumWays)
			if typ == 0 {
				return
			}

			// Check for the same value repeated.
			comb := eax ^ ebx ^ ecx
			if comb == last {
				nSame++
				if nSame == 100 {
					return
				}
			}
			last = comb

			switch level {
			case 1:
				switch typ {
				case 1:
					// Data cache
					c.Cache.L1D = size
				case 2:
					// Inst cache
					c.Cache.L1I = size
				default:
					if c.Cache.L1D < 0 {
						c.Cache.L1I = size
					}
					if c.Cache.L1I < 0 {
						c.Cache.L1I = size
					}
				}
			case 2:
				c.Cache.L2 = size
			case 3:
				c.Cache.L3 = size
			}
		}
	}
}

type SGXEPCSection struct {
	BaseAddress uint64
	EPCSize     uint64
}

type SGXSupport struct {
	Available           bool
	LaunchControl       bool
	SGX1Supported       bool
	SGX2Supported       bool
	MaxEnclaveSizeNot64 int64
	MaxEnclaveSize64    int64
	EPCSections         []SGXEPCSection
}

func hasSGX(available, lc bool) (rval SGXSupport) {
	rval.Available = available

	if !available {
		return
	}

	rval.LaunchControl = lc

	a, _, _, d := cpuidex(0x12, 0)
	rval.SGX1Supported = a&0x01 != 0
	rval.SGX2Supported = a&0x02 != 0
	rval.MaxEnclaveSizeNot64 = 1 << (d & 0xFF)     // pow 2
	rval.MaxEnclaveSize64 = 1 << ((d >> 8) & 0xFF) // pow 2
	rval.EPCSections = make([]SGXEPCSection, 0)

	for subleaf := uint32(2); subleaf < 2+8; subleaf++ {
		eax, ebx, ecx, edx := cpuidex(0x12, subleaf)
		leafType := eax & 0xf

		if leafType == 0 {
			// Invalid subleaf, stop iterating
			break
		} else if leafType == 1 {
			// EPC Section subleaf
			baseAddress := uint64(eax&0xfffff000) + (uint64(ebx&0x000fffff) << 32)
			size := uint64(ecx&0xfffff000) + (uint64(edx&0x000fffff) << 32)

			section := SGXEPCSection{BaseAddress: baseAddress, EPCSize: size}
			rval.EPCSections = append(rval.EPCSections, section)
		}
	}

	return
}

type AMDMemEncryptionSupport struct {
	Available          bool
	CBitPossition      uint32
	NumVMPL            uint32
	PhysAddrReduction  uint32
	NumEntryptedGuests uint32
	MinSevNoEsAsid     uint32
}

func hasAMDMemEncryption(available bool) (rval AMDMemEncryptionSupport) {
	rval.Available = available
	if !available {
		return
	}

	_, b, c, d := cpuidex(0x8000001f, 0)

	rval.CBitPossition = b & 0x3f
	rval.PhysAddrReduction = (b >> 6) & 0x3F
	rval.NumVMPL = (b >> 12) & 0xf
	rval.NumEntryptedGuests = c
	rval.MinSevNoEsAsid = d

	return
}

func support() flagSet {
	var fs flagSet
	mfi := maxFunctionID()
	vend, _ := vendorID()
	if mfi < 0x1 {
		return fs
	}
	family, model, _ := familyModel()

	_, _, c, d := cpuid(1)
	fs.setIf((d&(1<<0)) != 0, X87)
	fs.setIf((d&(1<<8)) != 0, CMPXCHG8)
	fs.setIf((d&(1<<11)) != 0, SYSEE)
	fs.setIf((d&(1<<15)) != 0, CMOV)
	fs.setIf((d&(1<<23)) != 0, MMX)
	fs.setIf((d&(1<<24)) != 0, FXSR)
	fs.setIf((d&(1<<25)) != 0, FXSROPT)
	fs.setIf((d&(1<<25)) != 0, SSE)
	fs.setIf((d&(1<<26)) != 0, SSE2)
	fs.setIf((c&1) != 0, SSE3)
	fs.setIf((c&(1<<5)) != 0, VMX)
	fs.setIf((c&(1<<9)) != 0, SSSE3)
	fs.setIf((c&(1<<19)) != 0, SSE4)
	fs.setIf((c&(1<<20)) != 0, SSE42)
	fs.setIf((c&(1<<25)) != 0, AESNI)
	fs.setIf((c&(1<<1)) != 0, CLMUL)
	fs.setIf(c&(1<<22) != 0, MOVBE)
	fs.setIf(c&(1<<23) != 0, POPCNT)
	fs.setIf(c&(1<<30) != 0, RDRAND)

	// This bit has been reserved by Intel & AMD for use by hypervisors,
	// and indicates the presence of a hypervisor.
	fs.setIf(c&(1<<31) != 0, HYPERVISOR)
	fs.setIf(c&(1<<29) != 0, F16C)
	fs.setIf(c&(1<<13) != 0, CX16)

	if vend == Intel && (d&(1<<28)) != 0 && mfi >= 4 {
		fs.setIf(threadsPerCore() > 1, HTT)
	}
	if vend == AMD && (d&(1<<28)) != 0 && mfi >= 4 {
		fs.setIf(threadsPerCore() > 1, HTT)
	}
	fs.setIf(c&1<<26 != 0, XSAVE)
	fs.setIf(c&1<<27 != 0, OSXSAVE)
	// Check XGETBV/XSAVE (26), OXSAVE (27) and AVX (28) bits
	const avxCheck = 1<<26 | 1<<27 | 1<<28
	if c&avxCheck == avxCheck {
		// Check for OS support
		eax, _ := xgetbv(0)
		if (eax & 0x6) == 0x6 {
			fs.set(AVX)
			switch vend {
			case Intel:
				// Older than Haswell.
				fs.setIf(family == 6 && model < 60, AVXSLOW)
			case AMD:
				// Older than Zen 2
				fs.setIf(family < 23 || (family == 23 && model < 49), AVXSLOW)
			}
		}
	}
	// FMA3 can be used with SSE registers, so no OS support is strictly needed.
	// fma3 and OSXSAVE needed.
	const fma3Check = 1<<12 | 1<<27
	fs.setIf(c&fma3Check == fma3Check, FMA3)

	// Check AVX2, AVX2 requires OS support, but BMI1/2 don't.
	if mfi >= 7 {
		_, ebx, ecx, edx := cpuidex(7, 0)
		if fs.inSet(AVX) && (ebx&0x00000020) != 0 {
			fs.set(AVX2)
		}
		// CPUID.(EAX=7, ECX=0).EBX
		if (ebx & 0x00000008) != 0 {
			fs.set(BMI1)
			fs.setIf((ebx&0x00000100) != 0, BMI2)
		}
		fs.setIf(ebx&(1<<2) != 0, SGX)
		fs.setIf(ebx&(1<<4) != 0, HLE)
		fs.setIf(ebx&(1<<9) != 0, ERMS)
		fs.setIf(ebx&(1<<11) != 0, RTM)
		fs.setIf(ebx&(1<<14) != 0, MPX)
		fs.setIf(ebx&(1<<18) != 0, RDSEED)
		fs.setIf(ebx&(1<<19) != 0, ADX)
		fs.setIf(ebx&(1<<29) != 0, SHA)

		// CPUID.(EAX=7, ECX=0).ECX
		fs.setIf(ecx&(1<<5) != 0, WAITPKG)
		fs.setIf(ecx&(1<<7) != 0, CETSS)
		fs.setIf(ecx&(1<<8) != 0, GFNI)
		fs.setIf(ecx&(1<<9) != 0, VAES)
		fs.setIf(ecx&(1<<10) != 0, VPCLMULQDQ)
		fs.setIf(ecx&(1<<13) != 0, TME)
		fs.setIf(ecx&(1<<25) != 0, CLDEMOTE)
		fs.setIf(ecx&(1<<23) != 0, KEYLOCKER)
		fs.setIf(ecx&(1<<27) != 0, MOVDIRI)
		fs.setIf(ecx&(1<<28) != 0, MOVDIR64B)
		fs.setIf(ecx&(1<<29) != 0, ENQCMD)
		fs.setIf(ecx&(1<<30) != 0, SGXLC)

		// CPUID.(EAX=7, ECX=0).EDX
		fs.setIf(edx&(1<<4) != 0, FSRM)
		fs.setIf(edx&(1<<9) != 0, SRBDS_CTRL)
		fs.setIf(edx&(1<<10) != 0, MD_CLEAR)
		fs.setIf(edx&(1<<11) != 0, RTM_ALWAYS_ABORT)
		fs.setIf(edx&(1<<14) != 0, SERIALIZE)
		fs.setIf(edx&(1<<15) != 0, HYBRID_CPU)
		fs.setIf(edx&(1<<16) != 0, TSXLDTRK)
		fs.setIf(edx&(1<<18) != 0, PCONFIG)
		fs.setIf(edx&(1<<20) != 0, CETIBT)
		fs.setIf(edx&(1<<26) != 0, IBPB)
		fs.setIf(edx&(1<<27) != 0, STIBP)
		fs.setIf(edx&(1<<28) != 0, FLUSH_L1D)
		fs.setIf(edx&(1<<29) != 0, IA32_ARCH_CAP)
		fs.setIf(edx&(1<<30) != 0, IA32_CORE_CAP)
		fs.setIf(edx&(1<<31) != 0, SPEC_CTRL_SSBD)

		// CPUID.(EAX=7, ECX=1).EAX
		eax1, _, _, edx1 := cpuidex(7, 1)
		fs.setIf(fs.inSet(AVX) && eax1&(1<<4) != 0, AVXVNNI)
		fs.setIf(eax1&(1<<1) != 0, SM3_X86)
		fs.setIf(eax1&(1<<2) != 0, SM4_X86)
		fs.setIf(eax1&(1<<7) != 0, CMPCCXADD)
		fs.setIf(eax1&(1<<10) != 0, MOVSB_ZL)
		fs.setIf(eax1&(1<<11) != 0, STOSB_SHORT)
		fs.setIf(eax1&(1<<12) != 0, CMPSB_SCADBS_SHORT)
		fs.setIf(eax1&(1<<22) != 0, HRESET)
		fs.setIf(eax1&(1<<23) != 0, AVXIFMA)
		fs.setIf(eax1&(1<<26) != 0, LAM)

		// CPUID.(EAX=7, ECX=1).EDX
		fs.setIf(edx1&(1<<4) != 0, AVXVNNIINT8)
		fs.setIf(edx1&(1<<5) != 0, AVXNECONVERT)
		fs.setIf(edx1&(1<<6) != 0, AMXTRANSPOSE)
		fs.setIf(edx1&(1<<7) != 0, AMXTF32)
		fs.setIf(edx1&(1<<8) != 0, AMXCOMPLEX)
		fs.setIf(edx1&(1<<10) != 0, AVXVNNIINT16)
		fs.setIf(edx1&(1<<14) != 0, PREFETCHI)
		fs.setIf(edx1&(1<<19) != 0, AVX10)
		fs.setIf(edx1&(1<<21) != 0, APX_F)

		// Only detect AVX-512 features if XGETBV is supported
		if c&((1<<26)|(1<<27)) == (1<<26)|(1<<27) {
			// Check for OS support
			eax, _ := xgetbv(0)

			// Verify that XCR0[7:5] = ‘111b’ (OPMASK state, upper 256-bit of ZMM0-ZMM15 and
			// ZMM16-ZMM31 state are enabled by OS)
			/// and that XCR0[2:1] = ‘11b’ (XMM state and YMM state are enabled by OS).
			hasAVX512 := (eax>>5)&7 == 7 && (eax>>1)&3 == 3
			if runtime.GOOS == "darwin" {
				hasAVX512 = fs.inSet(AVX) && darwinHasAVX512()
			}
			if hasAVX512 {
				fs.setIf(ebx&(1<<16) != 0, AVX512F)
				fs.setIf(ebx&(1<<17) != 0, AVX512DQ)
				fs.setIf(ebx&(1<<21) != 0, AVX512IFMA)
				fs.setIf(ebx&(1<<26) != 0, AVX512PF)
				fs.setIf(ebx&(1<<27) != 0, AVX512ER)
				fs.setIf(ebx&(1<<28) != 0, AVX512CD)
				fs.setIf(ebx&(1<<30) != 0, AVX512BW)
				fs.setIf(ebx&(1<<31) != 0, AVX512VL)
				// ecx
				fs.setIf(ecx&(1<<1) != 0, AVX512VBMI)
				fs.setIf(ecx&(1<<3) != 0, AMXFP8)
				fs.setIf(ecx&(1<<6) != 0, AVX512VBMI2)
				fs.setIf(ecx&(1<<11) != 0, AVX512VNNI)
				fs.setIf(ecx&(1<<12) != 0, AVX512BITALG)
				fs.setIf(ecx&(1<<14) != 0, AVX512VPOPCNTDQ)
				// edx
				fs.setIf(edx&(1<<8) != 0, AVX512VP2INTERSECT)
				fs.setIf(edx&(1<<22) != 0, AMXBF16)
				fs.setIf(edx&(1<<23) != 0, AVX512FP16)
				fs.setIf(edx&(1<<24) != 0, AMXTILE)
				fs.setIf(edx&(1<<25) != 0, AMXINT8)
				// eax1 = CPUID.(EAX=7, ECX=1).EAX
				fs.setIf(eax1&(1<<5) != 0, AVX512BF16)
				fs.setIf(eax1&(1<<19) != 0, WRMSRNS)
				fs.setIf(eax1&(1<<21) != 0, AMXFP16)
				fs.setIf(eax1&(1<<27) != 0, MSRLIST)
			}
		}

		// CPUID.(EAX=7, ECX=2)
		_, _, _, edx = cpuidex(7, 2)
		fs.setIf(edx&(1<<0) != 0, PSFD)
		fs.setIf(edx&(1<<1) != 0, IDPRED_CTRL)
		fs.setIf(edx&(1<<2) != 0, RRSBA_CTRL)
		fs.setIf(edx&(1<<4) != 0, BHI_CTRL)
		fs.setIf(edx&(1<<5) != 0, MCDT_NO)

		if fs.inSet(SGX) {
			eax, _, _, _ := cpuidex(0x12, 0)
			fs.setIf(eax&(1<<12) != 0, SGXPQC)
		}

		// Add keylocker features.
		if fs.inSet(KEYLOCKER) && mfi >= 0x19 {
			_, ebx, _, _ := cpuidex(0x19, 0)
			fs.setIf(ebx&5 == 5, KEYLOCKERW) // Bit 0 and 2 (1+4)
		}

		// Add AVX10 features.
		if fs.inSet(AVX10) && mfi >= 0x24 {
			_, ebx, _, _ := cpuidex(0x24, 0)
			fs.setIf(ebx&(1<<16) != 0, AVX10_128)
			fs.setIf(ebx&(1<<17) != 0, AVX10_256)
			fs.setIf(ebx&(1<<18) != 0, AVX10_512)
		}

	}

	// Processor Extended State Enumeration Sub-leaf (EAX = 0DH, ECX = 1)
	// EAX
	// Bit 00: XSAVEOPT is available.
	// Bit 01: Supports XSAVEC and the compacted form of XRSTOR if set.
	// Bit 02: Supports XGETBV with ECX = 1 if set.
	// Bit 03: Supports XSAVES/XRSTORS and IA32_XSS if set.
	// Bits 31 - 04: Reserved.
	// EBX
	// Bits 31 - 00: The size in bytes of the XSAVE area containing all states enabled by XCRO | IA32_XSS.
	// ECX
	// Bits 31 - 00: Reports the supported bits of the lower 32 bits of the IA32_XSS MSR. IA32_XSS[n] can be set to 1 only if ECX[n] is 1.
	// EDX?
	// Bits 07 - 00: Used for XCR0. Bit 08: PT state. Bit 09: Used for XCR0. Bits 12 - 10: Reserved. Bit 13: HWP state. Bits 31 - 14: Reserved.
	if mfi >= 0xd {
		if fs.inSet(XSAVE) {
			eax, _, _, _ := cpuidex(0xd, 1)
			fs.setIf(eax&(1<<0) != 0, XSAVEOPT)
			fs.setIf(eax&(1<<1) != 0, XSAVEC)
			fs.setIf(eax&(1<<2) != 0, XGETBV1)
			fs.setIf(eax&(1<<3) != 0, XSAVES)
		}
	}
	if maxExtendedFunction() >= 0x80000001 {
		_, _, c, d := cpuid(0x80000001)
		if (c & (1 << 5)) != 0 {
			fs.set(LZCNT)
			fs.set(POPCNT)
		}
		// ECX
		fs.setIf((c&(1<<0)) != 0, LAHF)
		fs.setIf((c&(1<<2)) != 0, SVM)
		fs.setIf((c&(1<<6)) != 0, SSE4A)
		fs.setIf((c&(1<<10)) != 0, IBS)
		fs.setIf((c&(1<<22)) != 0, TOPEXT)

		// EDX
		fs.setIf(d&(1<<11) != 0, SYSCALL)
		fs.setIf(d&(1<<20) != 0, NX)
		fs.setIf(d&(1<<22) != 0, MMXEXT)
		fs.setIf(d&(1<<23) != 0, MMX)
		fs.setIf(d&(1<<24) != 0, FXSR)
		fs.setIf(d&(1<<25) != 0, FXSROPT)
		fs.setIf(d&(1<<27) != 0, RDTSCP)
		fs.setIf(d&(1<<30) != 0, AMD3DNOWEXT)
		fs.setIf(d&(1<<31) != 0, AMD3DNOW)

		/* XOP and FMA4 use the AVX instruction coding scheme, so they can't be
		 * used unless the OS has AVX support. */
		if fs.inSet(AVX) {
			fs.setIf((c&(1<<11)) != 0, XOP)
			fs.setIf((c&(1<<16)) != 0, FMA4)
		}

	}
	if maxExtendedFunction() >= 0x80000007 {
		_, b, _, d := cpuid(0x80000007)
		fs.setIf((b&(1<<0)) != 0, MCAOVERFLOW)
		fs.setIf((b&(1<<1)) != 0, SUCCOR)
		fs.setIf((b&(1<<2)) != 0, HWA)
		fs.setIf((d&(1<<9)) != 0, CPBOOST)
	}

	if maxExtendedFunction() >= 0x80000008 {
		_, b, _, _ := cpuid(0x80000008)
		fs.setIf(b&(1<<28) != 0, PSFD)
		fs.setIf(b&(1<<27) != 0, CPPC)
		fs.setIf(b&(1<<24) != 0, SPEC_CTRL_SSBD)
		fs.setIf(b&(1<<23) != 0, PPIN)
		fs.setIf(b&(1<<21) != 0, TLB_FLUSH_NESTED)
		fs.setIf(b&(1<<20) != 0, EFER_LMSLE_UNS)
		fs.setIf(b&(1<<19) != 0, IBRS_PROVIDES_SMP)
		fs.setIf(b&(1<<18) != 0, IBRS_PREFERRED)
		fs.setIf(b&(1<<17) != 0, STIBP_ALWAYSON)
		fs.setIf(b&(1<<15) != 0, STIBP)
		fs.setIf(b&(1<<14) != 0, IBRS)
		fs.setIf((b&(1<<13)) != 0, INT_WBINVD)
		fs.setIf(b&(1<<12) != 0, IBPB)
		fs.setIf((b&(1<<9)) != 0, WBNOINVD)
		fs.setIf((b&(1<<8)) != 0, MCOMMIT)
		fs.setIf((b&(1<<4)) != 0, RDPRU)
		fs.setIf((b&(1<<3)) != 0, INVLPGB)
		fs.setIf((b&(1<<1)) != 0, MSRIRC)
		fs.setIf((b&(1<<0)) != 0, CLZERO)
	}

	if fs.inSet(SVM) && maxExtendedFunction() >= 0x8000000A {
		_, _, _, edx := cpuid(0x8000000A)
		fs.setIf((edx>>0)&1 == 1, SVMNP)
		fs.setIf((edx>>1)&1 == 1, LBRVIRT)
		fs.setIf((edx>>2)&1 == 1, SVML)
		fs.setIf((edx>>3)&1 == 1, NRIPS)
		fs.setIf((edx>>4)&1 == 1, TSCRATEMSR)
		fs.setIf((edx>>5)&1 == 1, VMCBCLEAN)
		fs.setIf((edx>>6)&1 == 1, SVMFBASID)
		fs.setIf((edx>>7)&1 == 1, SVMDA)
		fs.setIf((edx>>10)&1 == 1, SVMPF)
		fs.setIf((edx>>12)&1 == 1, SVMPFT)
	}

	if maxExtendedFunction() >= 0x8000001a {
		eax, _, _, _ := cpuid(0x8000001a)
		fs.setIf((eax>>0)&1 == 1, FP128)
		fs.setIf((eax>>1)&1 == 1, MOVU)
		fs.setIf((eax>>2)&1 == 1, FP256)
	}

	if maxExtendedFunction() >= 0x8000001b && fs.inSet(IBS) {
		eax, _, _, _ := cpuid(0x8000001b)
		fs.setIf((eax>>0)&1 == 1, IBSFFV)
		fs.setIf((eax>>1)&1 == 1, IBSFETCHSAM)
		fs.setIf((eax>>2)&1 == 1, IBSOPSAM)
		fs.setIf((eax>>3)&1 == 1, IBSRDWROPCNT)
		fs.setIf((eax>>4)&1 == 1, IBSOPCNT)
		fs.setIf((eax>>5)&1 == 1, IBSBRNTRGT)
		fs.setIf((eax>>6)&1 == 1, IBSOPCNTEXT)
		fs.setIf((eax>>7)&1 == 1, IBSRIPINVALIDCHK)
		fs.setIf((eax>>8)&1 == 1, IBS_OPFUSE)
		fs.setIf((eax>>9)&1 == 1, IBS_FETCH_CTLX)
		fs.setIf((eax>>10)&1 == 1, IBS_OPDATA4) // Doc says "Fixed,0. IBS op data 4 MSR supported", but assuming they mean 1.
		fs.setIf((eax>>11)&1 == 1, IBS_ZEN4)
	}

	if maxExtendedFunction() >= 0x8000001f && vend == AMD {
		a, _, _, _ := cpuid(0x8000001f)
		fs.setIf((a>>0)&1 == 1, SME)
		fs.setIf((a>>1)&1 == 1, SEV)
		fs.setIf((a>>2)&1 == 1, MSR_PAGEFLUSH)
		fs.setIf((a>>3)&1 == 1, SEV_ES)
		fs.setIf((a>>4)&1 == 1, SEV_SNP)
		fs.setIf((a>>5)&1 == 1, VMPL)
		fs.setIf((a>>10)&1 == 1, SME_COHERENT)
		fs.setIf((a>>11)&1 == 1, SEV_64BIT)
		fs.setIf((a>>12)&1 == 1, SEV_RESTRICTED)
		fs.setIf((a>>13)&1 == 1, SEV_ALTERNATIVE)
		fs.setIf((a>>14)&1 == 1, SEV_DEBUGSWAP)
		fs.setIf((a>>15)&1 == 1, IBS_PREVENTHOST)
		fs.setIf((a>>16)&1 == 1, VTE)
		fs.setIf((a>>24)&1 == 1, VMSA_REGPROT)
	}

	if maxExtendedFunction() >= 0x80000021 && vend == AMD {
		a, _, c, _ := cpuid(0x80000021)
		fs.setIf((a>>31)&1 == 1, SRSO_MSR_FIX)
		fs.setIf((a>>30)&1 == 1, SRSO_USER_KERNEL_NO)
		fs.setIf((a>>29)&1 == 1, SRSO_NO)
		fs.setIf((a>>28)&1 == 1, IBPB_BRTYPE)
		fs.setIf((a>>27)&1 == 1, SBPB)
		fs.setIf((c>>1)&1 == 1, TSA_L1_NO)
		fs.setIf((c>>2)&1 == 1, TSA_SQ_NO)
		fs.setIf((a>>5)&1 == 1, TSA_VERW_CLEAR)
	}
	if vend == AMD {
		if family < 0x19 {
			// AMD CPUs that are older than Family 19h are not vulnerable to TSA but do not set TSA_L1_NO or TSA_SQ_NO.
			// Source: https://www.amd.com/content/dam/amd/en/documents/resources/bulletin/technical-guidance-for-mitigating-transient-scheduler-attacks.pdf
			fs.set(TSA_L1_NO)
			fs.set(TSA_SQ_NO)
		} else if family == 0x1a {
			// AMD Family 1Ah models 00h-4Fh and 60h-7Fh are also not vulnerable to TSA but do not set TSA_L1_NO or TSA_SQ_NO.
			// Future AMD CPUs will set these CPUID bits if appropriate. CPUs will be designed to set these CPUID bits if appropriate.
			notVuln := model <= 0x4f || (model >= 0x60 && model <= 0x7f)
			fs.setIf(notVuln, TSA_L1_NO, TSA_SQ_NO)
		}
	}

	if mfi >= 0x20 {
		// Microsoft has decided to purposefully hide the information
		// of the guest TEE when VMs are being created using Hyper-V.
		//
		// This leads us to check for the Hyper-V cpuid features
		// (0x4000000C), and then for the `ebx` value set.
		//
		// For Intel TDX, `ebx` is set as `0xbe3`, being 3 the part
		// we're mostly interested about,according to:
		// https://github.com/torvalds/linux/blob/d2f51b3516dade79269ff45eae2a7668ae711b25/arch/x86/include/asm/hyperv-tlfs.h#L169-L174
		_, ebx, _, _ := cpuid(0x4000000C)
		fs.setIf(ebx == 0xbe3, TDX_GUEST)
	}

	if mfi >= 0x21 {
		// Intel Trusted Domain Extensions Guests have their own cpuid leaf (0x21).
		_, ebx, ecx, edx := cpuid(0x21)
		identity := string(valAsString(ebx, edx, ecx))
		fs.setIf(identity == "IntelTDX    ", TDX_GUEST)
	}

	return fs
}

func (c *CPUInfo) supportAVX10() uint8 {
	if c.maxFunc >= 0x24 && c.featureSet.inSet(AVX10) {
		_, ebx, _, _ := cpuidex(0x24, 0)
		return uint8(ebx)
	}
	return 0
}

func valAsString(values ...uint32) []byte {
	r := make([]byte, 4*len(values))
	for i, v := range values {
		dst := r[i*4:]
		dst[0] = byte(v & 0xff)
		dst[1] = byte((v >> 8) & 0xff)
		dst[2] = byte((v >> 16) & 0xff)
		dst[3] = byte((v >> 24) & 0xff)
		switch {
		case dst[0] == 0:
			return r[:i*4]
		case dst[1] == 0:
			return r[:i*4+1]
		case dst[2] == 0:
			return r[:i*4+2]
		case dst[3] == 0:
			return r[:i*4+3]
		}
	}
	return r
}

func parseLeaf0AH(c *CPUInfo, eax, ebx, edx uint32) (info PerformanceMonitoringInfo) {
	info.VersionID = uint8(eax & 0xFF)
	info.NumGPCounters = uint8((eax >> 8) & 0xFF)
	info.GPPMCWidth = uint8((eax >> 16) & 0xFF)

	info.RawEBX = ebx
	info.RawEAX = eax
	info.RawEDX = edx

	if info.VersionID > 1 { // This information is only valid if VersionID > 1
		info.NumFixedPMC = uint8(edx & 0x1F)          // Bits 4:0
		info.FixedPMCWidth = uint8((edx >> 5) & 0xFF) // Bits 12:5
	}
	if info.VersionID > 0 {
		// first 4 fixed events are always instructions retired, cycles, ref cycles and topdown slots
		if ebx == 0x0 && info.NumFixedPMC == 3 {
			c.featureSet.set(PMU_FIXEDCOUNTER_INSTRUCTIONS)
			c.featureSet.set(PMU_FIXEDCOUNTER_CYCLES)
			c.featureSet.set(PMU_FIXEDCOUNTER_REFCYCLES)
		}
		if ebx == 0x0 && info.NumFixedPMC == 4 {
			c.featureSet.set(PMU_FIXEDCOUNTER_INSTRUCTIONS)
			c.featureSet.set(PMU_FIXEDCOUNTER_CYCLES)
			c.featureSet.set(PMU_FIXEDCOUNTER_REFCYCLES)
			c.featureSet.set(PMU_FIXEDCOUNTER_TOPDOWN_SLOTS)
		}
		if ebx != 0x0 {
			if ((ebx >> 0) & 1) == 0 {
				c.featureSet.set(PMU_FIXEDCOUNTER_INSTRUCTIONS)
			}
			if ((ebx >> 1) & 1) == 0 {
				c.featureSet.set(PMU_FIXEDCOUNTER_CYCLES)
			}
			if ((ebx >> 2) & 1) == 0 {
				c.featureSet.set(PMU_FIXEDCOUNTER_REFCYCLES)
			}
			if ((ebx >> 3) & 1) == 0 {
				c.featureSet.set(PMU_FIXEDCOUNTER_TOPDOWN_SLOTS)
			}
		}
	}
	return info
}
