//go:build !darwin && !linux && !freebsd && !windows

package quic

import (
	"net"
	"net/netip"
)

func newConn(c net.PacketConn, supportsDF bool) (*basicConn, error) {
	return &basicConn{PacketConn: c, supportsDF: supportsDF}, nil
}

func inspectReadBuffer(any) (int, error)  { return 0, nil }
func inspectWriteBuffer(any) (int, error) { return 0, nil }

type packetInfo struct {
	addr netip.Addr
}

func (i *packetInfo) OOB() []byte { return nil }
