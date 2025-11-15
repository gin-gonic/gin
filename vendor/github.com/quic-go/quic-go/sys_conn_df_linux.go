//go:build linux

package quic

import (
	"errors"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/quic-go/quic-go/internal/utils"
)

func setDF(rawConn syscall.RawConn) (bool, error) {
	// Enabling IP_MTU_DISCOVER will force the kernel to return "sendto: message too long"
	// and the datagram will not be fragmented
	var errDFIPv4, errDFIPv6 error
	if err := rawConn.Control(func(fd uintptr) {
		errDFIPv4 = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_MTU_DISCOVER, unix.IP_PMTUDISC_PROBE)
		errDFIPv6 = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_MTU_DISCOVER, unix.IPV6_PMTUDISC_PROBE)
	}); err != nil {
		return false, err
	}
	switch {
	case errDFIPv4 == nil && errDFIPv6 == nil:
		utils.DefaultLogger.Debugf("Setting DF for IPv4 and IPv6.")
	case errDFIPv4 == nil && errDFIPv6 != nil:
		utils.DefaultLogger.Debugf("Setting DF for IPv4.")
	case errDFIPv4 != nil && errDFIPv6 == nil:
		utils.DefaultLogger.Debugf("Setting DF for IPv6.")
	case errDFIPv4 != nil && errDFIPv6 != nil:
		return false, errors.New("setting DF failed for both IPv4 and IPv6")
	}
	return true, nil
}

func isSendMsgSizeErr(err error) bool {
	// https://man7.org/linux/man-pages/man7/udp.7.html
	return errors.Is(err, unix.EMSGSIZE)
}

func isRecvMsgSizeErr(error) bool { return false }
