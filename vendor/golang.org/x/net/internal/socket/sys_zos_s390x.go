// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"net"
	"syscall"
	"unsafe"
)

func syscall_syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
func syscall_syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

func probeProtocolStack() int {
	return 4 // sizeof(int) on GOOS=zos GOARCH=s390x
}

func getsockopt(s uintptr, level, name int, b []byte) (int, error) {
	l := uint32(len(b))
	_, _, errno := syscall_syscall6(syscall.SYS_GETSOCKOPT, s, uintptr(level), uintptr(name), uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&l)), 0)
	return int(l), errnoErr(errno)
}

func setsockopt(s uintptr, level, name int, b []byte) error {
	_, _, errno := syscall_syscall6(syscall.SYS_SETSOCKOPT, s, uintptr(level), uintptr(name), uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), 0)
	return errnoErr(errno)
}

func recvmsg(s uintptr, buffers [][]byte, oob []byte, flags int, network string) (n, oobn int, recvflags int, from net.Addr, err error) {
	var h msghdr
	vs := make([]iovec, len(buffers))
	var sa []byte
	if network != "tcp" {
		sa = make([]byte, sizeofSockaddrInet6)
	}
	h.pack(vs, buffers, oob, sa)
	sn, _, errno := syscall_syscall(syscall.SYS___RECVMSG_A, s, uintptr(unsafe.Pointer(&h)), uintptr(flags))
	n = int(sn)
	oobn = h.controllen()
	recvflags = h.flags()
	err = errnoErr(errno)
	if network != "tcp" {
		var err2 error
		from, err2 = parseInetAddr(sa, network)
		if err2 != nil && err == nil {
			err = err2
		}
	}
	return
}

func sendmsg(s uintptr, buffers [][]byte, oob []byte, to net.Addr, flags int) (int, error) {
	var h msghdr
	vs := make([]iovec, len(buffers))
	var sa []byte
	if to != nil {
		var a [sizeofSockaddrInet6]byte
		n := marshalInetAddr(to, a[:])
		sa = a[:n]
	}
	h.pack(vs, buffers, oob, sa)
	n, _, errno := syscall_syscall(syscall.SYS___SENDMSG_A, s, uintptr(unsafe.Pointer(&h)), uintptr(flags))
	return int(n), errnoErr(errno)
}
