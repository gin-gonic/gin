// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package socket

import (
	"net"
	"unsafe"

	"golang.org/x/sys/unix"
)

//go:linkname syscall_getsockopt syscall.getsockopt
func syscall_getsockopt(s, level, name int, val unsafe.Pointer, vallen *uint32) error

//go:linkname syscall_setsockopt syscall.setsockopt
func syscall_setsockopt(s, level, name int, val unsafe.Pointer, vallen uintptr) error

func getsockopt(s uintptr, level, name int, b []byte) (int, error) {
	l := uint32(len(b))
	err := syscall_getsockopt(int(s), level, name, unsafe.Pointer(&b[0]), &l)
	return int(l), err
}

func setsockopt(s uintptr, level, name int, b []byte) error {
	return syscall_setsockopt(int(s), level, name, unsafe.Pointer(&b[0]), uintptr(len(b)))
}

func recvmsg(s uintptr, buffers [][]byte, oob []byte, flags int, network string) (n, oobn int, recvflags int, from net.Addr, err error) {
	var unixFrom unix.Sockaddr
	n, oobn, recvflags, unixFrom, err = unix.RecvmsgBuffers(int(s), buffers, oob, flags)
	if unixFrom != nil {
		from = sockaddrToAddr(unixFrom, network)
	}
	return
}

func sendmsg(s uintptr, buffers [][]byte, oob []byte, to net.Addr, flags int) (int, error) {
	var unixTo unix.Sockaddr
	if to != nil {
		unixTo = addrToSockaddr(to)
	}
	return unix.SendmsgBuffers(int(s), buffers, oob, unixTo, flags)
}

// addrToSockaddr converts a net.Addr to a unix.Sockaddr.
func addrToSockaddr(a net.Addr) unix.Sockaddr {
	var (
		ip   net.IP
		port int
		zone string
	)
	switch a := a.(type) {
	case *net.TCPAddr:
		ip = a.IP
		port = a.Port
		zone = a.Zone
	case *net.UDPAddr:
		ip = a.IP
		port = a.Port
		zone = a.Zone
	case *net.IPAddr:
		ip = a.IP
		zone = a.Zone
	default:
		return nil
	}

	if ip4 := ip.To4(); ip4 != nil {
		sa := unix.SockaddrInet4{Port: port}
		copy(sa.Addr[:], ip4)
		return &sa
	}

	if ip6 := ip.To16(); ip6 != nil && ip.To4() == nil {
		sa := unix.SockaddrInet6{Port: port}
		copy(sa.Addr[:], ip6)
		if zone != "" {
			sa.ZoneId = uint32(zoneCache.index(zone))
		}
		return &sa
	}

	return nil
}

// sockaddrToAddr converts a unix.Sockaddr to a net.Addr.
func sockaddrToAddr(sa unix.Sockaddr, network string) net.Addr {
	var (
		ip   net.IP
		port int
		zone string
	)
	switch sa := sa.(type) {
	case *unix.SockaddrInet4:
		ip = make(net.IP, net.IPv4len)
		copy(ip, sa.Addr[:])
		port = sa.Port
	case *unix.SockaddrInet6:
		ip = make(net.IP, net.IPv6len)
		copy(ip, sa.Addr[:])
		port = sa.Port
		if sa.ZoneId > 0 {
			zone = zoneCache.name(int(sa.ZoneId))
		}
	default:
		return nil
	}

	switch network {
	case "tcp", "tcp4", "tcp6":
		return &net.TCPAddr{IP: ip, Port: port, Zone: zone}
	case "udp", "udp4", "udp6":
		return &net.UDPAddr{IP: ip, Port: port, Zone: zone}
	default:
		return &net.IPAddr{IP: ip, Zone: zone}
	}
}
