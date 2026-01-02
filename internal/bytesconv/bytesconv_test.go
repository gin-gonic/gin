// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bytesconv

import (
	"bytes"
	cRand "crypto/rand"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var (
	testString = "Albert Einstein: Logic will get you from A to B. Imagination will take you everywhere."
	testBytes  = []byte(testString)
)

func rawBytesToStr(b []byte) string {
	return string(b)
}

func rawStrToBytes(s string) []byte {
	return []byte(s)
}

// go test -v

func TestBytesToString(t *testing.T) {
	data := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		_, err := cRand.Read(data)
		if err != nil {
			t.Fatal(err)
		}
		if rawBytesToStr(data) != BytesToString(data) {
			t.Fatal("don't match")
		}
	}
}

func TestBytesToStringEmpty(t *testing.T) {
	if got := BytesToString([]byte{}); got != "" {
		t.Fatalf("BytesToString([]byte{}) = %q; want empty string", got)
	}
	if got := BytesToString(nil); got != "" {
		t.Fatalf("BytesToString(nil) = %q; want empty string", got)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func TestStringToBytes(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandStringBytesMaskImprSrcSB(64)
		if !bytes.Equal(rawStrToBytes(s), StringToBytes(s)) {
			t.Fatal("don't match")
		}
	}
}

func TestStringToBytesEmpty(t *testing.T) {
	b := StringToBytes("")
	if len(b) != 0 {
		t.Fatalf(`StringToBytes("") length = %d; want 0`, len(b))
	}
	if !bytes.Equal(b, []byte("")) {
		t.Fatalf(`StringToBytes("") = %v; want []byte("")`, b)
	}
}

// go test -v -run=none -bench=^BenchmarkBytesConv -benchmem=true

func BenchmarkBytesConvBytesToStrRaw(b *testing.B) {
	for b.Loop() {
		rawBytesToStr(testBytes)
	}
}

func BenchmarkBytesConvBytesToStr(b *testing.B) {
	for b.Loop() {
		BytesToString(testBytes)
	}
}

func BenchmarkBytesConvStrToBytesRaw(b *testing.B) {
	for b.Loop() {
		rawStrToBytes(testString)
	}
}

func BenchmarkBytesConvStrToBytes(b *testing.B) {
	for b.Loop() {
		StringToBytes(testString)
	}
}
