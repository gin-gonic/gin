// Copyright (c) 2015 Klaus Post, released under MIT License. See LICENSE file.

//+build amd64,!gccgo,!noasm,!appengine

// func asmCpuid(op uint32) (eax, ebx, ecx, edx uint32)
TEXT ·asmCpuid(SB), 7, $0
	XORQ CX, CX
	MOVL op+0(FP), AX
	CPUID
	MOVL AX, eax+8(FP)
	MOVL BX, ebx+12(FP)
	MOVL CX, ecx+16(FP)
	MOVL DX, edx+20(FP)
	RET

// func asmCpuidex(op, op2 uint32) (eax, ebx, ecx, edx uint32)
TEXT ·asmCpuidex(SB), 7, $0
	MOVL op+0(FP), AX
	MOVL op2+4(FP), CX
	CPUID
	MOVL AX, eax+8(FP)
	MOVL BX, ebx+12(FP)
	MOVL CX, ecx+16(FP)
	MOVL DX, edx+20(FP)
	RET

// func asmXgetbv(index uint32) (eax, edx uint32)
TEXT ·asmXgetbv(SB), 7, $0
	MOVL index+0(FP), CX
	BYTE $0x0f; BYTE $0x01; BYTE $0xd0 // XGETBV
	MOVL AX, eax+8(FP)
	MOVL DX, edx+12(FP)
	RET

// func asmRdtscpAsm() (eax, ebx, ecx, edx uint32)
TEXT ·asmRdtscpAsm(SB), 7, $0
	BYTE $0x0F; BYTE $0x01; BYTE $0xF9 // RDTSCP
	MOVL AX, eax+0(FP)
	MOVL BX, ebx+4(FP)
	MOVL CX, ecx+8(FP)
	MOVL DX, edx+12(FP)
	RET

// From https://go-review.googlesource.com/c/sys/+/285572/
// func asmDarwinHasAVX512() bool
TEXT ·asmDarwinHasAVX512(SB), 7, $0-1
	MOVB $0, ret+0(FP) // default to false

#ifdef GOOS_darwin // return if not darwin
#ifdef GOARCH_amd64 // return if not amd64
// These values from:
// https://github.com/apple/darwin-xnu/blob/xnu-4570.1.46/osfmk/i386/cpu_capabilities.h
#define commpage64_base_address         0x00007fffffe00000
#define commpage64_cpu_capabilities64   (commpage64_base_address+0x010)
#define commpage64_version              (commpage64_base_address+0x01E)
#define hasAVX512F                      0x0000004000000000
	MOVQ $commpage64_version, BX
	MOVW (BX), AX
	CMPW AX, $13                            // versions < 13 do not support AVX512
	JL   no_avx512
	MOVQ $commpage64_cpu_capabilities64, BX
	MOVQ (BX), AX
	MOVQ $hasAVX512F, CX
	ANDQ CX, AX
	JZ   no_avx512
	MOVB $1, ret+0(FP)

no_avx512:
#endif
#endif
	RET

