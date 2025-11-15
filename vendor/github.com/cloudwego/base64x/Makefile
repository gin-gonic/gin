.PHONY: all clean

CFLAGS := -mavx
CFLAGS += -mavx2
CFLAGS += -mno-bmi
CFLAGS += -mno-red-zone
CFLAGS += -fno-asynchronous-unwind-tables
CFLAGS += -fno-stack-protector
CFLAGS += -fno-exceptions
CFLAGS += -fno-builtin
CFLAGS += -fno-rtti
CFLAGS += -nostdlib
CFLAGS += -O3

NATIVE_ASM := $(wildcard native/*.S)
NATIVE_SRC := $(wildcard native/*.h)
NATIVE_SRC += $(wildcard native/*.c)

all: native_amd64.s

clean:
	rm -vf native_text_amd64.go native_subr_amd64.go output/*.s

native_amd64.s: ${NATIVE_SRC} ${NATIVE_ASM} native_amd64.go
	mkdir -p output
	clang ${CFLAGS} -S -o output/native.s native/native.c
	python3 tools/asm2asm/asm2asm.py -r native_amd64.go output/native.s ${NATIVE_ASM}
	awk '{gsub(/Text__native_entry__/, "text__native_entry__")}1' native_text_amd64.go > native_text_amd64.go.tmp && mv native_text_amd64.go.tmp native_text_amd64.go
	awk '{gsub(/Funcs/, "funcs")}1' native_subr_amd64.go > native_subr_amd64.go.tmp && mv native_subr_amd64.go.tmp native_subr_amd64.go
