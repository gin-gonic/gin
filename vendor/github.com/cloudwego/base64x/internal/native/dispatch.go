package native

import (
	"unsafe"
    
    `github.com/klauspost/cpuid/v2`
	"github.com/cloudwego/base64x/internal/rt"
	"github.com/cloudwego/base64x/internal/native/avx2"
	"github.com/cloudwego/base64x/internal/native/sse"
)

var (
    hasAVX2 = cpuid.CPU.Has(cpuid.AVX2)
    hasSSE = cpuid.CPU.Has(cpuid.SSE)
)

var (
	S_b64decode uintptr
	S_b64encode uintptr
)

var (
	F_b64decode func(out unsafe.Pointer, src unsafe.Pointer, len int, mod int) (ret int)
	F_b64encode func(out unsafe.Pointer, src unsafe.Pointer, mod int)
)

func useAVX2() {
	avx2.Use()
	S_b64decode = avx2.S_b64decode
	S_b64encode = avx2.S_b64encode

	F_b64decode = avx2.F_b64decode
	F_b64encode = avx2.F_b64encode
}

func useSSE() {
	sse.Use()
	S_b64decode = sse.S_b64decode
	S_b64encode = sse.S_b64encode

	F_b64decode = sse.F_b64decode
	F_b64encode = sse.F_b64encode
}

//go:nosplit
func B64Decode(out *[]byte, src unsafe.Pointer, len int, mod int) (ret int) {
    return F_b64decode(rt.NoEscape(unsafe.Pointer(out)), rt.NoEscape(unsafe.Pointer(src)), len, mod)
}

//go:nosplit
func B64Encode(out *[]byte, src *[]byte, mod int) {
	F_b64encode(rt.NoEscape(unsafe.Pointer(out)), rt.NoEscape(unsafe.Pointer(src)), mod)
}

func init() {
	if hasAVX2 {
		useAVX2()
	} else if hasSSE {
		useSSE()
	} else {
		panic("Unsupported CPU, lacks of AVX2 or SSE CPUID Flag. maybe it's too old to run Sonic.")
	}
}
