# cpuid
Package cpuid provides information about the CPU running the current program.

CPU features are detected on startup, and kept for fast access through the life of the application.
Currently x86 / x64 (AMD64/i386) and ARM (ARM64) is supported, and no external C (cgo) code is used, which should make the library very easy to use.

You can access the CPU information by accessing the shared CPU variable of the cpuid library.

Package home: https://github.com/klauspost/cpuid

[![PkgGoDev](https://pkg.go.dev/badge/github.com/klauspost/cpuid)](https://pkg.go.dev/github.com/klauspost/cpuid/v2)
[![Go](https://github.com/klauspost/cpuid/actions/workflows/go.yml/badge.svg)](https://github.com/klauspost/cpuid/actions/workflows/go.yml)

## installing

`go get -u github.com/klauspost/cpuid/v2` using modules.
Drop `v2` for others.

Installing binary:

`go install github.com/klauspost/cpuid/v2/cmd/cpuid@latest`

Or download binaries from release page: https://github.com/klauspost/cpuid/releases

### Homebrew

For macOS/Linux users, you can install via [brew](https://brew.sh/)

```sh
$ brew install cpuid
```

## example

```Go
package main

import (
	"fmt"
	"strings"

	. "github.com/klauspost/cpuid/v2"
)

func main() {
	// Print basic CPU information:
	fmt.Println("Name:", CPU.BrandName)
	fmt.Println("PhysicalCores:", CPU.PhysicalCores)
	fmt.Println("ThreadsPerCore:", CPU.ThreadsPerCore)
	fmt.Println("LogicalCores:", CPU.LogicalCores)
	fmt.Println("Family", CPU.Family, "Model:", CPU.Model, "Vendor ID:", CPU.VendorID)
	fmt.Println("Features:", strings.Join(CPU.FeatureSet(), ","))
	fmt.Println("Cacheline bytes:", CPU.CacheLine)
	fmt.Println("L1 Data Cache:", CPU.Cache.L1D, "bytes")
	fmt.Println("L1 Instruction Cache:", CPU.Cache.L1I, "bytes")
	fmt.Println("L2 Cache:", CPU.Cache.L2, "bytes")
	fmt.Println("L3 Cache:", CPU.Cache.L3, "bytes")
	fmt.Println("Frequency", CPU.Hz, "hz")

	// Test if we have these specific features:
	if CPU.Supports(SSE, SSE2) {
		fmt.Println("We have Streaming SIMD 2 Extensions")
	}
}
```

Sample output:
```
>go run main.go
Name: AMD Ryzen 9 3950X 16-Core Processor
PhysicalCores: 16
ThreadsPerCore: 2
LogicalCores: 32
Family 23 Model: 113 Vendor ID: AMD
Features: ADX,AESNI,AVX,AVX2,BMI1,BMI2,CLMUL,CMOV,CX16,F16C,FMA3,HTT,HYPERVISOR,LZCNT,MMX,MMXEXT,NX,POPCNT,RDRAND,RDSEED,RDTSCP,SHA,SSE,SSE2,SSE3,SSE4,SSE42,SSE4A,SSSE3
Cacheline bytes: 64
L1 Data Cache: 32768 bytes
L1 Instruction Cache: 32768 bytes
L2 Cache: 524288 bytes
L3 Cache: 16777216 bytes
Frequency 0 hz
We have Streaming SIMD 2 Extensions
```

# usage

The `cpuid.CPU` provides access to CPU features. Use `cpuid.CPU.Supports()` to check for CPU features.
A faster `cpuid.CPU.Has()` is provided which will usually be inlined by the gc compiler.  

To test a larger number of features, they can be combined using `f := CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SYSCALL, SSE, SSE2)`, etc.
This can be using with `cpuid.CPU.HasAll(f)` to quickly test if all features are supported.

Note that for some cpu/os combinations some features will not be detected.
`amd64` has rather good support and should work reliably on all platforms.

Note that hypervisors may not pass through all CPU features through to the guest OS,
so even if your host supports a feature it may not be visible on guests.

## arm64 feature detection

Not all operating systems provide ARM features directly 
and there is no safe way to do so for the rest.

Currently `arm64/linux` and `arm64/freebsd` should be quite reliable. 
`arm64/darwin` adds features expected from the M1 processor, but a lot remains undetected.

A `DetectARM()` can be used if you are able to control your deployment,
it will detect CPU features, but may crash if the OS doesn't intercept the calls.
A `-cpu.arm` flag for detecting unsafe ARM features can be added. See below.
 
Note that currently only features are detected on ARM, 
no additional information is currently available. 

## flags

It is possible to add flags that affects cpu detection.

For this the `Flags()` command is provided.

This must be called *before* `flag.Parse()` AND after the flags have been parsed `Detect()` must be called.

This means that any detection used in `init()` functions will not contain these flags.

Example:

```Go
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/klauspost/cpuid/v2"
)

func main() {
	cpuid.Flags()
	flag.Parse()
	cpuid.Detect()

	// Test if we have these specific features:
	if cpuid.CPU.Supports(cpuid.SSE, cpuid.SSE2) {
		fmt.Println("We have Streaming SIMD 2 Extensions")
	}
}
```

## commandline

Download as binary from: https://github.com/klauspost/cpuid/releases

Install from source:

`go install github.com/klauspost/cpuid/v2/cmd/cpuid@latest`

### Example

```
λ cpuid
Name: AMD Ryzen 9 3950X 16-Core Processor
Vendor String: AuthenticAMD
Vendor ID: AMD
PhysicalCores: 16
Threads Per Core: 2
Logical Cores: 32
CPU Family 23 Model: 113
Features: ADX,AESNI,AVX,AVX2,BMI1,BMI2,CLMUL,CLZERO,CMOV,CMPXCHG8,CPBOOST,CX16,F16C,FMA3,FXSR,FXSROPT,HTT,HYPERVISOR,LAHF,LZCNT,MCAOVERFLOW,MMX,MMXEXT,MOVBE,NX,OSXSAVE,POPCNT,RDRAND,RDSEED,RDTSCP,SCE,SHA,SSE,SSE2,SSE3,SSE4,SSE42,SSE4A,SSSE3,SUCCOR,X87,XSAVE
Microarchitecture level: 3
Cacheline bytes: 64
L1 Instruction Cache: 32768 bytes
L1 Data Cache: 32768 bytes
L2 Cache: 524288 bytes
L3 Cache: 16777216 bytes

```
### JSON Output:

```
λ cpuid --json
{
  "BrandName": "AMD Ryzen 9 3950X 16-Core Processor",
  "VendorID": 2,
  "VendorString": "AuthenticAMD",
  "PhysicalCores": 16,
  "ThreadsPerCore": 2,
  "LogicalCores": 32,
  "Family": 23,
  "Model": 113,
  "CacheLine": 64,
  "Hz": 0,
  "BoostFreq": 0,
  "Cache": {
    "L1I": 32768,
    "L1D": 32768,
    "L2": 524288,
    "L3": 16777216
  },
  "SGX": {
    "Available": false,
    "LaunchControl": false,
    "SGX1Supported": false,
    "SGX2Supported": false,
    "MaxEnclaveSizeNot64": 0,
    "MaxEnclaveSize64": 0,
    "EPCSections": null
  },
  "Features": [
    "ADX",
    "AESNI",
    "AVX",
    "AVX2",
    "BMI1",
    "BMI2",
    "CLMUL",
    "CLZERO",
    "CMOV",
    "CMPXCHG8",
    "CPBOOST",
    "CX16",
    "F16C",
    "FMA3",
    "FXSR",
    "FXSROPT",
    "HTT",
    "HYPERVISOR",
    "LAHF",
    "LZCNT",
    "MCAOVERFLOW",
    "MMX",
    "MMXEXT",
    "MOVBE",
    "NX",
    "OSXSAVE",
    "POPCNT",
    "RDRAND",
    "RDSEED",
    "RDTSCP",
    "SCE",
    "SHA",
    "SSE",
    "SSE2",
    "SSE3",
    "SSE4",
    "SSE42",
    "SSE4A",
    "SSSE3",
    "SUCCOR",
    "X87",
    "XSAVE"
  ],
  "X64Level": 3
}
```

### Check CPU microarch level

```
λ cpuid --check-level=3
2022/03/18 17:04:40 AMD Ryzen 9 3950X 16-Core Processor
2022/03/18 17:04:40 Microarchitecture level 3 is supported. Max level is 3.
Exit Code 0

λ cpuid --check-level=4
2022/03/18 17:06:18 AMD Ryzen 9 3950X 16-Core Processor
2022/03/18 17:06:18 Microarchitecture level 4 not supported. Max level is 3.
Exit Code 1
```


## Available flags

### x86 & amd64 

| Feature Flag       | Description                                                                                                                                                                        |
|--------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ADX                | Intel ADX (Multi-Precision Add-Carry Instruction Extensions)                                                                                                                       |
| AESNI              | Advanced Encryption Standard New Instructions                                                                                                                                      |
| AMD3DNOW           | AMD 3DNOW                                                                                                                                                                          |
| AMD3DNOWEXT        | AMD 3DNowExt                                                                                                                                                                       |
| AMXBF16            | Tile computational operations on BFLOAT16 numbers                                                                                                                                  |
| AMXINT8            | Tile computational operations on 8-bit integers                                                                                                                                    |
| AMXFP16            | Tile computational operations on FP16 numbers                                                                                                                                      |
| AMXFP8             | Tile computational operations on FP8 numbers                                                                                                                                       |
| AMXCOMPLEX         | Tile computational operations on complex numbers                                                                                                                                   |
| AMXTILE            | Tile architecture                                                                                                                                                                  |
| AMXTF32            | Matrix Multiplication of TF32 Tiles into Packed Single Precision Tile                                                                                                              |
| AMXTRANSPOSE       | Tile multiply where the first operand is transposed                                                                                                                                |
| APX_F              | Intel APX                                                                                                                                                                          |
| AVX                | AVX functions                                                                                                                                                                      |
| AVX10              | If set the Intel AVX10 Converged Vector ISA is supported                                                                                                                           |
| AVX10_128          | If set indicates that AVX10 128-bit vector support is present                                                                                                                      |
| AVX10_256          | If set indicates that AVX10 256-bit vector support is present                                                                                                                      |
| AVX10_512          | If set indicates that AVX10 512-bit vector support is present                                                                                                                      |
| AVX2               | AVX2 functions                                                                                                                                                                     |
| AVX512BF16         | AVX-512 BFLOAT16 Instructions                                                                                                                                                      |
| AVX512BITALG       | AVX-512 Bit Algorithms                                                                                                                                                             |
| AVX512BW           | AVX-512 Byte and Word Instructions                                                                                                                                                 |
| AVX512CD           | AVX-512 Conflict Detection Instructions                                                                                                                                            |
| AVX512DQ           | AVX-512 Doubleword and Quadword Instructions                                                                                                                                       |
| AVX512ER           | AVX-512 Exponential and Reciprocal Instructions                                                                                                                                    |
| AVX512F            | AVX-512 Foundation                                                                                                                                                                 |
| AVX512FP16         | AVX-512 FP16 Instructions                                                                                                                                                          |
| AVX512IFMA         | AVX-512 Integer Fused Multiply-Add Instructions                                                                                                                                    |
| AVX512PF           | AVX-512 Prefetch Instructions                                                                                                                                                      |
| AVX512VBMI         | AVX-512 Vector Bit Manipulation Instructions                                                                                                                                       |
| AVX512VBMI2        | AVX-512 Vector Bit Manipulation Instructions, Version 2                                                                                                                            |
| AVX512VL           | AVX-512 Vector Length Extensions                                                                                                                                                   |
| AVX512VNNI         | AVX-512 Vector Neural Network Instructions                                                                                                                                         |
| AVX512VP2INTERSECT | AVX-512 Intersect for D/Q                                                                                                                                                          |
| AVX512VPOPCNTDQ    | AVX-512 Vector Population Count Doubleword and Quadword                                                                                                                            |
| AVXIFMA            | AVX-IFMA instructions                                                                                                                                                              |
| AVXNECONVERT       | AVX-NE-CONVERT instructions                                                                                                                                                        |
| AVXSLOW            | Indicates the CPU performs 2 128 bit operations instead of one                                                                                                                     |
| AVXVNNI            | AVX (VEX encoded) VNNI neural network instructions                                                                                                                                 |
| AVXVNNIINT8        | AVX-VNNI-INT8 instructions                                                                                                                                                         |
| AVXVNNIINT16       | AVX-VNNI-INT16 instructions                                                                                                                                                        |
| BHI_CTRL           | Branch History Injection and Intra-mode Branch Target Injection / CVE-2022-0001, CVE-2022-0002 / INTEL-SA-00598                                                                    |
| BMI1               | Bit Manipulation Instruction Set 1                                                                                                                                                 |
| BMI2               | Bit Manipulation Instruction Set 2                                                                                                                                                 |
| CETIBT             | Intel CET Indirect Branch Tracking                                                                                                                                                 |
| CETSS              | Intel CET Shadow Stack                                                                                                                                                             |
| CLDEMOTE           | Cache Line Demote                                                                                                                                                                  |
| CLMUL              | Carry-less Multiplication                                                                                                                                                          |
| CLZERO             | CLZERO instruction supported                                                                                                                                                       |
| CMOV               | i686 CMOV                                                                                                                                                                          |
| CMPCCXADD          | CMPCCXADD instructions                                                                                                                                                             |
| CMPSB_SCADBS_SHORT | Fast short CMPSB and SCASB                                                                                                                                                         |
| CMPXCHG8           | CMPXCHG8 instruction                                                                                                                                                               |
| CPBOOST            | Core Performance Boost                                                                                                                                                             |
| CPPC               | AMD: Collaborative Processor Performance Control                                                                                                                                   |
| CX16               | CMPXCHG16B Instruction                                                                                                                                                             |
| EFER_LMSLE_UNS     | AMD: =Core::X86::Msr::EFER[LMSLE] is not supported, and MBZ                                                                                                                        |
| ENQCMD             | Enqueue Command                                                                                                                                                                    |
| ERMS               | Enhanced REP MOVSB/STOSB                                                                                                                                                           |
| F16C               | Half-precision floating-point conversion                                                                                                                                           |
| FLUSH_L1D          | Flush L1D cache                                                                                                                                                                    |
| FMA3               | Intel FMA 3. Does not imply AVX.                                                                                                                                                   |
| FMA4               | Bulldozer FMA4 functions                                                                                                                                                           |
| FP128              | AMD: When set, the internal FP/SIMD execution datapath is 128-bits wide                                                                                                            |
| FP256              | AMD: When set, the internal FP/SIMD execution datapath is 256-bits wide                                                                                                            |
| FSRM               | Fast Short Rep Mov                                                                                                                                                                 |
| FXSR               | FXSAVE, FXRESTOR instructions, CR4 bit 9                                                                                                                                           |
| FXSROPT            | FXSAVE/FXRSTOR optimizations                                                                                                                                                       |
| GFNI               | Galois Field New Instructions. May require other features (AVX, AVX512VL,AVX512F) based on usage.                                                                                  |
| HLE                | Hardware Lock Elision                                                                                                                                                              |
| HRESET             | If set CPU supports history reset and the IA32_HRESET_ENABLE MSR                                                                                                                   |
| HTT                | Hyperthreading (enabled)                                                                                                                                                           |
| HWA                | Hardware assert supported. Indicates support for MSRC001_10                                                                                                                        |
| HYBRID_CPU         | This part has CPUs of more than one type.                                                                                                                                          |
| HYPERVISOR         | This bit has been reserved by Intel & AMD for use by hypervisors                                                                                                                   |
| IA32_ARCH_CAP      | IA32_ARCH_CAPABILITIES MSR (Intel)                                                                                                                                                 |
| IA32_CORE_CAP      | IA32_CORE_CAPABILITIES MSR                                                                                                                                                         |
| IBPB               | Indirect Branch Restricted Speculation (IBRS) and Indirect Branch Predictor Barrier (IBPB)                                                                                         |
| IBRS               | AMD: Indirect Branch Restricted Speculation                                                                                                                                        |
| IBRS_PREFERRED     | AMD: IBRS is preferred over software solution                                                                                                                                      |
| IBRS_PROVIDES_SMP  | AMD: IBRS provides Same Mode Protection                                                                                                                                            |
| IBS                | Instruction Based Sampling (AMD)                                                                                                                                                   |
| IBSBRNTRGT         | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSFETCHSAM        | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSFFV             | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSOPCNT           | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSOPCNTEXT        | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSOPSAM           | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSRDWROPCNT       | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBSRIPINVALIDCHK   | Instruction Based Sampling Feature (AMD)                                                                                                                                           |
| IBS_FETCH_CTLX     | AMD: IBS fetch control extended MSR supported                                                                                                                                      |
| IBS_OPDATA4        | AMD: IBS op data 4 MSR supported                                                                                                                                                   |
| IBS_OPFUSE         | AMD: Indicates support for IbsOpFuse                                                                                                                                               |
| IBS_PREVENTHOST    | Disallowing IBS use by the host supported                                                                                                                                          |
| IBS_ZEN4           | Fetch and Op IBS support IBS extensions added with Zen4                                                                                                                            |
| IDPRED_CTRL        | IPRED_DIS                                                                                                                                                                          |
| INT_WBINVD         | WBINVD/WBNOINVD are interruptible.                                                                                                                                                 |
| INVLPGB            | NVLPGB and TLBSYNC instruction supported                                                                                                                                           |
| KEYLOCKER          | Key locker                                                                                                                                                                         |
| KEYLOCKERW         | Key locker wide                                                                                                                                                                    |
| LAHF               | LAHF/SAHF in long mode                                                                                                                                                             |
| LAM                | If set, CPU supports Linear Address Masking                                                                                                                                        |
| LBRVIRT            | LBR virtualization                                                                                                                                                                 |
| LZCNT              | LZCNT instruction                                                                                                                                                                  |
| MCAOVERFLOW        | MCA overflow recovery support.                                                                                                                                                     |
| MCDT_NO            | Processor do not exhibit MXCSR Configuration Dependent Timing behavior and do not need to mitigate it.                                                                             |
| MCOMMIT            | MCOMMIT instruction supported                                                                                                                                                      |
| MD_CLEAR           | VERW clears CPU buffers                                                                                                                                                            |
| MMX                | standard MMX                                                                                                                                                                       |
| MMXEXT             | SSE integer functions or AMD MMX ext                                                                                                                                               |
| MOVBE              | MOVBE instruction (big-endian)                                                                                                                                                     |
| MOVDIR64B          | Move 64 Bytes as Direct Store                                                                                                                                                      |
| MOVDIRI            | Move Doubleword as Direct Store                                                                                                                                                    |
| MOVSB_ZL           | Fast Zero-Length MOVSB                                                                                                                                                             |
| MPX                | Intel MPX (Memory Protection Extensions)                                                                                                                                           |
| MOVU               | MOVU SSE instructions are more efficient and should be preferred to SSE	MOVL/MOVH. MOVUPS is more efficient than MOVLPS/MOVHPS. MOVUPD is more efficient than MOVLPD/MOVHPD        |
| MSRIRC             | Instruction Retired Counter MSR available                                                                                                                                          |
| MSRLIST            | Read/Write List of Model Specific Registers                                                                                                                                        |
| MSR_PAGEFLUSH      | Page Flush MSR available                                                                                                                                                           |
| NRIPS              | Indicates support for NRIP save on VMEXIT                                                                                                                                          |
| NX                 | NX (No-Execute) bit                                                                                                                                                                |
| OSXSAVE            | XSAVE enabled by OS                                                                                                                                                                |
| PCONFIG            | PCONFIG for Intel Multi-Key Total Memory Encryption                                                                                                                                |
| POPCNT             | POPCNT instruction                                                                                                                                                                 |
| PPIN               | AMD: Protected Processor Inventory Number support. Indicates that Protected Processor Inventory Number (PPIN) capability can be enabled                                            |
| PREFETCHI          | PREFETCHIT0/1 instructions                                                                                                                                                         |
| PSFD               | Predictive Store Forward Disable                                                                                                                                                   |
| RDPRU              | RDPRU instruction supported                                                                                                                                                        |
| RDRAND             | RDRAND instruction is available                                                                                                                                                    |
| RDSEED             | RDSEED instruction is available                                                                                                                                                    |
| RDTSCP             | RDTSCP Instruction                                                                                                                                                                 |
| RRSBA_CTRL         | Restricted RSB Alternate                                                                                                                                                           |
| RTM                | Restricted Transactional Memory                                                                                                                                                    |
| RTM_ALWAYS_ABORT   | Indicates that the loaded microcode is forcing RTM abort.                                                                                                                          |
| SERIALIZE          | Serialize Instruction Execution                                                                                                                                                    |
| SEV                | AMD Secure Encrypted Virtualization supported                                                                                                                                      |
| SEV_64BIT          | AMD SEV guest execution only allowed from a 64-bit host                                                                                                                            |
| SEV_ALTERNATIVE    | AMD SEV Alternate Injection supported                                                                                                                                              |
| SEV_DEBUGSWAP      | Full debug state swap supported for SEV-ES guests                                                                                                                                  |
| SEV_ES             | AMD SEV Encrypted State supported                                                                                                                                                  |
| SEV_RESTRICTED     | AMD SEV Restricted Injection supported                                                                                                                                             |
| SEV_SNP            | AMD SEV Secure Nested Paging supported                                                                                                                                             |
| SGX                | Software Guard Extensions                                                                                                                                                          |
| SGXLC              | Software Guard Extensions Launch Control                                                                                                                                           |
| SGXPQC             | Software Guard Extensions 256-bit Encryption                                                                                                                                       |
| SHA                | Intel SHA Extensions                                                                                                                                                               |
| SME                | AMD Secure Memory Encryption supported                                                                                                                                             |
| SME_COHERENT       | AMD Hardware cache coherency across encryption domains enforced                                                                                                                    |
| SM3_X86            | SM3 instructions                                                                                                                                                                   |
| SM4_X86            | SM4 instructions                                                                                                                                                                   |
| SPEC_CTRL_SSBD     | Speculative Store Bypass Disable                                                                                                                                                   |
| SRBDS_CTRL         | SRBDS mitigation MSR available                                                                                                                                                     |
| SSE                | SSE functions                                                                                                                                                                      |
| SSE2               | P4 SSE functions                                                                                                                                                                   |
| SSE3               | Prescott SSE3 functions                                                                                                                                                            |
| SSE4               | Penryn SSE4.1 functions                                                                                                                                                            |
| SSE42              | Nehalem SSE4.2 functions                                                                                                                                                           |
| SSE4A              | AMD Barcelona microarchitecture SSE4a instructions                                                                                                                                 |
| SSSE3              | Conroe SSSE3 functions                                                                                                                                                             |
| STIBP              | Single Thread Indirect Branch Predictors                                                                                                                                           |
| STIBP_ALWAYSON     | AMD: Single Thread Indirect Branch Prediction Mode has Enhanced Performance and may be left Always On                                                                              |
| STOSB_SHORT        | Fast short STOSB                                                                                                                                                                   |
| SUCCOR             | Software uncorrectable error containment and recovery capability.                                                                                                                  |
| SVM                | AMD Secure Virtual Machine                                                                                                                                                         |
| SVMDA              | Indicates support for the SVM decode assists.                                                                                                                                      |
| SVMFBASID          | SVM, Indicates that TLB flush events, including CR3 writes and CR4.PGE toggles, flush only the current ASID's TLB entries. Also indicates support for the extended VMCBTLB_Control |
| SVML               | AMD SVM lock. Indicates support for SVM-Lock.                                                                                                                                      |
| SVMNP              | AMD SVM nested paging                                                                                                                                                              |
| SVMPF              | SVM pause intercept filter. Indicates support for the pause intercept filter                                                                                                       |
| SVMPFT             | SVM PAUSE filter threshold. Indicates support for the PAUSE filter cycle count threshold                                                                                           |
| SYSCALL            | System-Call Extension (SCE): SYSCALL and SYSRET instructions.                                                                                                                      |
| SYSEE              | SYSENTER and SYSEXIT instructions                                                                                                                                                  |
| TBM                | AMD Trailing Bit Manipulation                                                                                                                                                      |
| TDX_GUEST          | Intel Trust Domain Extensions Guest                                                                                                                                                |
| TLB_FLUSH_NESTED   | AMD: Flushing includes all the nested translations for guest translations                                                                                                          |
| TME                | Intel Total Memory Encryption. The following MSRs are supported: IA32_TME_CAPABILITY, IA32_TME_ACTIVATE, IA32_TME_EXCLUDE_MASK, and IA32_TME_EXCLUDE_BASE.                         |
| TOPEXT             | TopologyExtensions: topology extensions support. Indicates support for CPUID Fn8000_001D_EAX_x[N:0]-CPUID Fn8000_001E_EDX.                                                         |
| TSA_L1_NO          | AMD only: Not vulnerable to TSA-L1                                                                                                                                                 |
| TSA_SQ_NO          | AMD only: Not vulnerable to TSA-SQ                                                                                                                                                 |
| TSA_VERW_CLEAR     | AMD: If set, the memory form of the VERW instruction may be used to help mitigate TSA                                                                                              |
| TSCRATEMSR         | MSR based TSC rate control. Indicates support for MSR TSC ratio MSRC000_0104                                                                                                       |
| TSXLDTRK           | Intel TSX Suspend Load Address Tracking                                                                                                                                            |
| VAES               | Vector AES. AVX(512) versions requires additional checks.                                                                                                                          |
| VMCBCLEAN          | VMCB clean bits. Indicates support for VMCB clean bits.                                                                                                                            |
| VMPL               | AMD VM Permission Levels supported                                                                                                                                                 |
| VMSA_REGPROT       | AMD VMSA Register Protection supported                                                                                                                                             |
| VMX                | Virtual Machine Extensions                                                                                                                                                         |
| VPCLMULQDQ         | Carry-Less Multiplication Quadword. Requires AVX for 3 register versions.                                                                                                          |
| VTE                | AMD Virtual Transparent Encryption supported                                                                                                                                       |
| WAITPKG            | TPAUSE, UMONITOR, UMWAIT                                                                                                                                                           |
| WBNOINVD           | Write Back and Do Not Invalidate Cache                                                                                                                                             |
| WRMSRNS            | Non-Serializing Write to Model Specific Register                                                                                                                                   |
| X87                | FPU                                                                                                                                                                                |
| XGETBV1            | Supports XGETBV with ECX = 1                                                                                                                                                       |
| XOP                | Bulldozer XOP functions                                                                                                                                                            |
| XSAVE              | XSAVE, XRESTOR, XSETBV, XGETBV                                                                                                                                                     |
| XSAVEC             | Supports XSAVEC and the compacted form of XRSTOR.                                                                                                                                  |
| XSAVEOPT           | XSAVEOPT available                                                                                                                                                                 |
| XSAVES             | Supports XSAVES/XRSTORS and IA32_XSS                                                                                                                                               |

# ARM features:

| Feature Flag | Description                                                      |
|--------------|------------------------------------------------------------------|
| AESARM       | AES instructions                                                 |
| ARMCPUID     | Some CPU ID registers readable at user-level                     |
| ASIMD        | Advanced SIMD                                                    |
| ASIMDDP      | SIMD Dot Product                                                 |
| ASIMDHP      | Advanced SIMD half-precision floating point                      |
| ASIMDRDM     | Rounding Double Multiply Accumulate/Subtract (SQRDMLAH/SQRDMLSH) |
| ATOMICS      | Large System Extensions (LSE)                                    |
| CRC32        | CRC32/CRC32C instructions                                        |
| DCPOP        | Data cache clean to Point of Persistence (DC CVAP)               |
| EVTSTRM      | Generic timer                                                    |
| FCMA         | Floatin point complex number addition and multiplication         |
| FHM          | FMLAL and FMLSL instructions                                     |
| FP           | Single-precision and double-precision floating point             |
| FPHP         | Half-precision floating point                                    |
| GPA          | Generic Pointer Authentication                                   |
| JSCVT        | Javascript-style double->int convert (FJCVTZS)                   |
| LRCPC        | Weaker release consistency (LDAPR, etc)                          |
| PMULL        | Polynomial Multiply instructions (PMULL/PMULL2)                  |
| RNDR         | Random Number instructions                                       |
| TLB          | Outer Shareable and TLB range maintenance instructions           |
| TS           | Flag manipulation instructions                                   |
| SHA1         | SHA-1 instructions (SHA1C, etc)                                  |
| SHA2         | SHA-2 instructions (SHA256H, etc)                                |
| SHA3         | SHA-3 instructions (EOR3, RAXI, XAR, BCAX)                       |
| SHA512       | SHA512 instructions                                              |
| SM3          | SM3 instructions                                                 |
| SM4          | SM4 instructions                                                 |
| SVE          | Scalable Vector Extension                                        |

# license

This code is published under an MIT license. See LICENSE file for more information.
