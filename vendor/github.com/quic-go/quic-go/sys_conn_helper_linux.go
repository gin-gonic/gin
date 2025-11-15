//go:build linux

package quic

import (
	"encoding/binary"
	"errors"
	"net/netip"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	msgTypeIPTOS = unix.IP_TOS
	ipv4PKTINFO  = unix.IP_PKTINFO
)

const ecnIPv4DataLen = 1

const batchSize = 8 // needs to smaller than MaxUint8 (otherwise the type of oobConn.readPos has to be changed)

var kernelVersionMajor int

func init() {
	kernelVersionMajor, _ = kernelVersion()
}

func forceSetReceiveBuffer(c syscall.RawConn, bytes int) error {
	var serr error
	if err := c.Control(func(fd uintptr) {
		serr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUFFORCE, bytes)
	}); err != nil {
		return err
	}
	return serr
}

func forceSetSendBuffer(c syscall.RawConn, bytes int) error {
	var serr error
	if err := c.Control(func(fd uintptr) {
		serr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_SNDBUFFORCE, bytes)
	}); err != nil {
		return err
	}
	return serr
}

func parseIPv4PktInfo(body []byte) (ip netip.Addr, ifIndex uint32, ok bool) {
	// struct in_pktinfo {
	// 	unsigned int   ipi_ifindex;  /* Interface index */
	// 	struct in_addr ipi_spec_dst; /* Local address */
	// 	struct in_addr ipi_addr;     /* Header Destination address */
	// };
	if len(body) != 12 {
		return netip.Addr{}, 0, false
	}
	return netip.AddrFrom4(*(*[4]byte)(body[8:12])), binary.NativeEndian.Uint32(body), true
}

// isGSOEnabled tests if the kernel supports GSO.
// Sending with GSO might still fail later on, if the interface doesn't support it (see isGSOError).
func isGSOEnabled(conn syscall.RawConn) bool {
	if kernelVersionMajor < 5 {
		return false
	}
	disabled, err := strconv.ParseBool(os.Getenv("QUIC_GO_DISABLE_GSO"))
	if err == nil && disabled {
		return false
	}
	var serr error
	if err := conn.Control(func(fd uintptr) {
		_, serr = unix.GetsockoptInt(int(fd), unix.IPPROTO_UDP, unix.UDP_SEGMENT)
	}); err != nil {
		return false
	}
	return serr == nil
}

func appendUDPSegmentSizeMsg(b []byte, size uint16) []byte {
	startLen := len(b)
	const dataLen = 2 // payload is a uint16
	b = append(b, make([]byte, unix.CmsgSpace(dataLen))...)
	h := (*unix.Cmsghdr)(unsafe.Pointer(&b[startLen]))
	h.Level = syscall.IPPROTO_UDP
	h.Type = unix.UDP_SEGMENT
	h.SetLen(unix.CmsgLen(dataLen))

	// UnixRights uses the private `data` method, but I *think* this achieves the same goal.
	offset := startLen + unix.CmsgSpace(0)
	*(*uint16)(unsafe.Pointer(&b[offset])) = size
	return b
}

func isGSOError(err error) bool {
	var serr *os.SyscallError
	if errors.As(err, &serr) {
		// EIO is returned by udp_send_skb() if the device driver does not have tx checksums enabled,
		// which is a hard requirement of UDP_SEGMENT. See:
		// https://git.kernel.org/pub/scm/docs/man-pages/man-pages.git/tree/man7/udp.7?id=806eabd74910447f21005160e90957bde4db0183#n228
		// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/net/ipv4/udp.c?h=v6.2&id=c9c3395d5e3dcc6daee66c6908354d47bf98cb0c#n942
		return serr.Err == unix.EIO
	}
	return false
}

// The first sendmsg call on a new UDP socket sometimes errors on Linux.
// It's not clear why this happens.
// See https://github.com/golang/go/issues/63322.
func isPermissionError(err error) bool {
	var serr *os.SyscallError
	if errors.As(err, &serr) {
		return serr.Syscall == "sendmsg" && serr.Err == unix.EPERM
	}
	return false
}

func isECNEnabled() bool {
	return kernelVersionMajor >= 5 && !isECNDisabledUsingEnv()
}

// kernelVersion returns major and minor kernel version numbers, parsed from
// the syscall.Uname's Release field, or 0, 0 if the version can't be obtained
// or parsed.
//
// copied from the standard library's internal/syscall/unix/kernel_version_linux.go
func kernelVersion() (major, minor int) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return
	}

	var (
		values    [2]int
		value, vi int
	)
	for _, c := range uname.Release {
		if '0' <= c && c <= '9' {
			value = (value * 10) + int(c-'0')
		} else {
			// Note that we're assuming N.N.N here.
			// If we see anything else, we are likely to mis-parse it.
			values[vi] = value
			vi++
			if vi >= len(values) {
				break
			}
			value = 0
		}
	}

	return values[0], values[1]
}
