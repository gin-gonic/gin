//go:build darwin

package quic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// for macOS versions, see https://en.wikipedia.org/wiki/Darwin_(operating_system)#Darwin_20_onwards
const (
	macOSVersion11 = 20
	macOSVersion15 = 24
)

func setDF(rawConn syscall.RawConn) (bool, error) {
	// Setting DF bit is only supported from macOS 11.
	// https://github.com/chromium/chromium/blob/117.0.5881.2/net/socket/udp_socket_posix.cc#L555
	version, err := getMacOSVersion()
	if err != nil || version < macOSVersion11 {
		return false, err
	}

	var controlErr error
	var disableDF bool
	if err := rawConn.Control(func(fd uintptr) {
		addr, err := unix.Getsockname(int(fd))
		if err != nil {
			controlErr = fmt.Errorf("getsockname: %w", err)
			return
		}

		// Dual-stack sockets are effectively IPv6 sockets (with IPV6_ONLY set to 0).
		// On macOS, the DF bit on dual-stack sockets is controlled by the IPV6_DONTFRAG option.
		// See https://datatracker.ietf.org/doc/draft-seemann-tsvwg-udp-fragmentation/ for details.
		switch addr.(type) {
		case *unix.SockaddrInet4:
			controlErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IP, unix.IP_DONTFRAG, 1)
		case *unix.SockaddrInet6:
			controlErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_DONTFRAG, 1)

			// Setting the DF bit on dual-stack sockets works since macOS Sequoia.
			// Disable DF on dual-stack sockets before Sequoia.
			if version < macOSVersion15 {
				// check if this is a dual-stack socket by reading the IPV6_V6ONLY flag
				v6only, err := unix.GetsockoptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_V6ONLY)
				if err != nil {
					controlErr = fmt.Errorf("getting IPV6_V6ONLY: %w", err)
					return
				}
				disableDF = v6only == 0
			}
		default:
			controlErr = fmt.Errorf("unknown address type: %T", addr)
		}
	}); err != nil {
		return false, err
	}
	if controlErr != nil {
		return false, controlErr
	}
	return !disableDF, nil
}

func isSendMsgSizeErr(err error) bool {
	return errors.Is(err, unix.EMSGSIZE)
}

func isRecvMsgSizeErr(error) bool { return false }

func getMacOSVersion() (int, error) {
	uname := &unix.Utsname{}
	if err := unix.Uname(uname); err != nil {
		return 0, err
	}

	release := string(uname.Release[:])
	idx := strings.Index(release, ".")
	if idx == -1 {
		return 0, nil
	}
	version, err := strconv.Atoi(release[:idx])
	if err != nil {
		return 0, err
	}
	return version, nil
}
