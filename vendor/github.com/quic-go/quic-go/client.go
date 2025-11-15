package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"net"

	"github.com/quic-go/quic-go/internal/protocol"
)

// make it possible to mock connection ID for initial generation in the tests
var generateConnectionIDForInitial = protocol.GenerateConnectionIDForInitial

// DialAddr establishes a new QUIC connection to a server.
// It resolves the address, and then creates a new UDP connection to dial the QUIC server.
// When the QUIC connection is closed, this UDP connection is closed.
// See [Dial] for more details.
func DialAddr(ctx context.Context, addr string, tlsConf *tls.Config, conf *Config) (*Conn, error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	tr, err := setupTransport(udpConn, tlsConf, true)
	if err != nil {
		return nil, err
	}
	conn, err := tr.dial(ctx, udpAddr, addr, tlsConf, conf, false)
	if err != nil {
		tr.Close()
		return nil, err
	}
	return conn, nil
}

// DialAddrEarly establishes a new 0-RTT QUIC connection to a server.
// See [DialAddr] for more details.
func DialAddrEarly(ctx context.Context, addr string, tlsConf *tls.Config, conf *Config) (*Conn, error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	tr, err := setupTransport(udpConn, tlsConf, true)
	if err != nil {
		return nil, err
	}
	conn, err := tr.dial(ctx, udpAddr, addr, tlsConf, conf, true)
	if err != nil {
		tr.Close()
		return nil, err
	}
	return conn, nil
}

// DialEarly establishes a new 0-RTT QUIC connection to a server using a net.PacketConn.
// See [Dial] for more details.
func DialEarly(ctx context.Context, c net.PacketConn, addr net.Addr, tlsConf *tls.Config, conf *Config) (*Conn, error) {
	dl, err := setupTransport(c, tlsConf, false)
	if err != nil {
		return nil, err
	}
	conn, err := dl.DialEarly(ctx, addr, tlsConf, conf)
	if err != nil {
		dl.Close()
		return nil, err
	}
	return conn, nil
}

// Dial establishes a new QUIC connection to a server using a net.PacketConn.
// If the PacketConn satisfies the [OOBCapablePacketConn] interface (as a [net.UDPConn] does),
// ECN and packet info support will be enabled. In this case, ReadMsgUDP and WriteMsgUDP
// will be used instead of ReadFrom and WriteTo to read/write packets.
// The [tls.Config] must define an application protocol (using tls.Config.NextProtos).
//
// This is a convenience function. More advanced use cases should instantiate a [Transport],
// which offers configuration options for a more fine-grained control of the connection establishment,
// including reusing the underlying UDP socket for multiple QUIC connections.
func Dial(ctx context.Context, c net.PacketConn, addr net.Addr, tlsConf *tls.Config, conf *Config) (*Conn, error) {
	dl, err := setupTransport(c, tlsConf, false)
	if err != nil {
		return nil, err
	}
	conn, err := dl.Dial(ctx, addr, tlsConf, conf)
	if err != nil {
		dl.Close()
		return nil, err
	}
	return conn, nil
}

func setupTransport(c net.PacketConn, tlsConf *tls.Config, createdPacketConn bool) (*Transport, error) {
	if tlsConf == nil {
		return nil, errors.New("quic: tls.Config not set")
	}
	return &Transport{
		Conn:        c,
		createdConn: createdPacketConn,
		isSingleUse: true,
	}, nil
}
