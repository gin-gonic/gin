package quic

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
)

type connCapabilities struct {
	// This connection has the Don't Fragment (DF) bit set.
	// This means it makes to run DPLPMTUD.
	DF bool
	// GSO (Generic Segmentation Offload) supported
	GSO bool
	// ECN (Explicit Congestion Notifications) supported
	ECN bool
}

// rawConn is a connection that allow reading of a receivedPackeh.
type rawConn interface {
	ReadPacket() (receivedPacket, error)
	// WritePacket writes a packet on the wire.
	// gsoSize is the size of a single packet, or 0 to disable GSO.
	// It is invalid to set gsoSize if capabilities.GSO is not set.
	WritePacket(b []byte, addr net.Addr, packetInfoOOB []byte, gsoSize uint16, ecn protocol.ECN) (int, error)
	LocalAddr() net.Addr
	SetReadDeadline(time.Time) error
	io.Closer

	capabilities() connCapabilities
}

// OOBCapablePacketConn is a connection that allows the reading of ECN bits from the IP header.
// If the PacketConn passed to the [Transport] satisfies this interface, quic-go will use it.
// In this case, ReadMsgUDP() will be used instead of ReadFrom() to read packets.
type OOBCapablePacketConn interface {
	net.PacketConn
	SyscallConn() (syscall.RawConn, error)
	SetReadBuffer(int) error
	ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error)
	WriteMsgUDP(b, oob []byte, addr *net.UDPAddr) (n, oobn int, err error)
}

var _ OOBCapablePacketConn = &net.UDPConn{}

func wrapConn(pc net.PacketConn) (rawConn, error) {
	if err := setReceiveBuffer(pc); err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			setBufferWarningOnce.Do(func() {
				if disable, _ := strconv.ParseBool(os.Getenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING")); disable {
					return
				}
				log.Printf("%s. See https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes for details.", err)
			})
		}
	}
	if err := setSendBuffer(pc); err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			setBufferWarningOnce.Do(func() {
				if disable, _ := strconv.ParseBool(os.Getenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING")); disable {
					return
				}
				log.Printf("%s. See https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes for details.", err)
			})
		}
	}

	conn, ok := pc.(interface {
		SyscallConn() (syscall.RawConn, error)
	})
	var supportsDF bool
	if ok {
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return nil, err
		}

		// only set DF on UDP sockets
		if _, ok := pc.LocalAddr().(*net.UDPAddr); ok {
			var err error
			supportsDF, err = setDF(rawConn)
			if err != nil {
				return nil, err
			}
		}
	}
	c, ok := pc.(OOBCapablePacketConn)
	if !ok {
		utils.DefaultLogger.Infof("PacketConn is not a net.UDPConn. Disabling optimizations possible on UDP connections.")
		return &basicConn{PacketConn: pc, supportsDF: supportsDF}, nil
	}
	return newConn(c, supportsDF)
}

// The basicConn is the most trivial implementation of a rawConn.
// It reads a single packet from the underlying net.PacketConn.
// It is used when
// * the net.PacketConn is not a OOBCapablePacketConn, and
// * when the OS doesn't support OOB.
type basicConn struct {
	net.PacketConn
	supportsDF bool
}

var _ rawConn = &basicConn{}

func (c *basicConn) ReadPacket() (receivedPacket, error) {
	buffer := getPacketBuffer()
	// The packet size should not exceed protocol.MaxPacketBufferSize bytes
	// If it does, we only read a truncated packet, which will then end up undecryptable
	buffer.Data = buffer.Data[:protocol.MaxPacketBufferSize]
	n, addr, err := c.ReadFrom(buffer.Data)
	if err != nil {
		return receivedPacket{}, err
	}
	return receivedPacket{
		remoteAddr: addr,
		rcvTime:    monotime.Now(),
		data:       buffer.Data[:n],
		buffer:     buffer,
	}, nil
}

func (c *basicConn) WritePacket(b []byte, addr net.Addr, _ []byte, gsoSize uint16, ecn protocol.ECN) (n int, err error) {
	if gsoSize != 0 {
		panic("cannot use GSO with a basicConn")
	}
	if ecn != protocol.ECNUnsupported {
		panic("cannot use ECN with a basicConn")
	}
	return c.WriteTo(b, addr)
}

func (c *basicConn) capabilities() connCapabilities { return connCapabilities{DF: c.supportsDF} }
