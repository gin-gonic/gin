package quic

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qlogwriter"
)

// ErrTransportClosed is returned by the [Transport]'s Listen or Dial method after it was closed.
var ErrTransportClosed = &errTransportClosed{}

type errTransportClosed struct {
	err error
}

func (e *errTransportClosed) Unwrap() []error { return []error{net.ErrClosed, e.err} }

func (e *errTransportClosed) Error() string {
	if e.err == nil {
		return "quic: transport closed"
	}
	return fmt.Sprintf("quic: transport closed: %s", e.err)
}

func (e *errTransportClosed) Is(target error) bool {
	_, ok := target.(*errTransportClosed)
	return ok
}

var errListenerAlreadySet = errors.New("listener already set")

type closePacket struct {
	payload []byte
	addr    net.Addr
	info    packetInfo
}

// The Transport is the central point to manage incoming and outgoing QUIC connections.
// QUIC demultiplexes connections based on their QUIC Connection IDs, not based on the 4-tuple.
// This means that a single UDP socket can be used for listening for incoming connections, as well as
// for dialing an arbitrary number of outgoing connections.
// A Transport handles a single net.PacketConn, and offers a range of configuration options
// compared to the simple helper functions like [Listen] and [Dial] that this package provides.
type Transport struct {
	// A single net.PacketConn can only be handled by one Transport.
	// Bad things will happen if passed to multiple Transports.
	//
	// A number of optimizations will be enabled if the connections implements the OOBCapablePacketConn interface,
	// as a *net.UDPConn does.
	// 1. It enables the Don't Fragment (DF) bit on the IP header.
	//    This is required to run DPLPMTUD (Path MTU Discovery, RFC 8899).
	// 2. It enables reading of the ECN bits from the IP header.
	//    This allows the remote node to speed up its loss detection and recovery.
	// 3. It uses batched syscalls (recvmmsg) to more efficiently receive packets from the socket.
	// 4. It uses Generic Segmentation Offload (GSO) to efficiently send batches of packets (on Linux).
	//
	// After passing the connection to the Transport, it's invalid to call ReadFrom or WriteTo on the connection.
	Conn net.PacketConn

	// The length of the connection ID in bytes.
	// It can be any value between 1 and 20.
	// Due to the increased risk of collisions, it is not recommended to use connection IDs shorter than 4 bytes.
	// If unset, a 4 byte connection ID will be used.
	ConnectionIDLength int

	// Use for generating new connection IDs.
	// This allows the application to control of the connection IDs used,
	// which allows routing / load balancing based on connection IDs.
	// All Connection IDs returned by the ConnectionIDGenerator MUST
	// have the same length.
	ConnectionIDGenerator ConnectionIDGenerator

	// The StatelessResetKey is used to generate stateless reset tokens.
	// If no key is configured, sending of stateless resets is disabled.
	// It is highly recommended to configure a stateless reset key, as stateless resets
	// allow the peer to quickly recover from crashes and reboots of this node.
	// See section 10.3 of RFC 9000 for details.
	StatelessResetKey *StatelessResetKey

	// The TokenGeneratorKey is used to encrypt session resumption tokens.
	// If no key is configured, a random key will be generated.
	// If multiple servers are authoritative for the same domain, they should use the same key,
	// see section 8.1.3 of RFC 9000 for details.
	TokenGeneratorKey *TokenGeneratorKey

	// MaxTokenAge is the maximum age of the resumption token presented during the handshake.
	// These tokens allow skipping address resumption when resuming a QUIC connection,
	// and are especially useful when using 0-RTT.
	// If not set, it defaults to 24 hours.
	// See section 8.1.3 of RFC 9000 for details.
	MaxTokenAge time.Duration

	// DisableVersionNegotiationPackets disables the sending of Version Negotiation packets.
	// This can be useful if version information is exchanged out-of-band.
	// It has no effect for clients.
	DisableVersionNegotiationPackets bool

	// VerifySourceAddress decides if a connection attempt originating from unvalidated source
	// addresses first needs to go through source address validation using QUIC's Retry mechanism,
	// as described in RFC 9000 section 8.1.2.
	// Note that the address passed to this callback is unvalidated, and might be spoofed in case
	// of an attack.
	// Validating the source address adds one additional network roundtrip to the handshake,
	// and should therefore only be used if a suspiciously high number of incoming connection is recorded.
	// For most use cases, wrapping the Allow function of a rate.Limiter will be a reasonable
	// implementation of this callback (negating its return value).
	VerifySourceAddress func(net.Addr) bool

	// ConnContext is called when the server accepts a new connection. To reject a connection return
	// a non-nil error.
	// The context is closed when the connection is closed, or when the handshake fails for any reason.
	// The context returned from the callback is used to derive every other context used during the
	// lifetime of the connection:
	// * the context passed to crypto/tls (and used on the tls.ClientHelloInfo)
	// * the context used in Config.QlogTrace
	// * the context returned from Conn.Context
	// * the context returned from SendStream.Context
	// It is not used for dialed connections.
	ConnContext func(context.Context, *ClientInfo) (context.Context, error)

	// A Tracer traces events that don't belong to a single QUIC connection.
	// Recorder.Close is called when the transport is closed.
	Tracer qlogwriter.Recorder

	mutex       sync.Mutex
	handlers    map[protocol.ConnectionID]packetHandler
	resetTokens map[protocol.StatelessResetToken]packetHandler

	initOnce sync.Once
	initErr  error

	// If no ConnectionIDGenerator is set, this is the ConnectionIDLength.
	connIDLen int
	// Set in init.
	// If no ConnectionIDGenerator is set, this is set to a default.
	connIDGenerator   ConnectionIDGenerator
	statelessResetter *statelessResetter

	server *baseServer

	conn rawConn

	closeQueue          chan closePacket
	statelessResetQueue chan receivedPacket

	listening   chan struct{} // is closed when listen returns
	closeErr    error
	createdConn bool
	isSingleUse bool // was created for a single server or client, i.e. by calling quic.Listen or quic.Dial

	readingNonQUICPackets atomic.Bool
	nonQUICPackets        chan receivedPacket

	logger utils.Logger
}

// Listen starts listening for incoming QUIC connections.
// There can only be a single listener on any net.PacketConn.
// Listen may only be called again after the current listener was closed.
func (t *Transport) Listen(tlsConf *tls.Config, conf *Config) (*Listener, error) {
	s, err := t.createServer(tlsConf, conf, false)
	if err != nil {
		return nil, err
	}
	return &Listener{baseServer: s}, nil
}

// ListenEarly starts listening for incoming QUIC connections.
// There can only be a single listener on any net.PacketConn.
// ListenEarly may only be called again after the current listener was closed.
func (t *Transport) ListenEarly(tlsConf *tls.Config, conf *Config) (*EarlyListener, error) {
	s, err := t.createServer(tlsConf, conf, true)
	if err != nil {
		return nil, err
	}
	return &EarlyListener{baseServer: s}, nil
}

func (t *Transport) createServer(tlsConf *tls.Config, conf *Config, allow0RTT bool) (*baseServer, error) {
	if tlsConf == nil {
		return nil, errors.New("quic: tls.Config not set")
	}
	if err := validateConfig(conf); err != nil {
		return nil, err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.closeErr != nil {
		return nil, t.closeErr
	}
	if t.server != nil {
		return nil, errListenerAlreadySet
	}
	conf = populateConfig(conf)
	if err := t.init(false); err != nil {
		return nil, err
	}
	maxTokenAge := t.MaxTokenAge
	if maxTokenAge == 0 {
		maxTokenAge = 24 * time.Hour
	}
	s := newServer(
		t.conn,
		(*packetHandlerMap)(t),
		t.connIDGenerator,
		t.statelessResetter,
		t.ConnContext,
		tlsConf,
		conf,
		t.Tracer,
		t.closeServer,
		*t.TokenGeneratorKey,
		maxTokenAge,
		t.VerifySourceAddress,
		t.DisableVersionNegotiationPackets,
		allow0RTT,
	)
	t.server = s
	return s, nil
}

// Dial dials a new connection to a remote host (not using 0-RTT).
func (t *Transport) Dial(ctx context.Context, addr net.Addr, tlsConf *tls.Config, conf *Config) (*Conn, error) {
	return t.dial(ctx, addr, "", tlsConf, conf, false)
}

// DialEarly dials a new connection, attempting to use 0-RTT if possible.
func (t *Transport) DialEarly(ctx context.Context, addr net.Addr, tlsConf *tls.Config, conf *Config) (*Conn, error) {
	return t.dial(ctx, addr, "", tlsConf, conf, true)
}

func (t *Transport) dial(ctx context.Context, addr net.Addr, host string, tlsConf *tls.Config, conf *Config, use0RTT bool) (*Conn, error) {
	if err := t.init(t.isSingleUse); err != nil {
		return nil, err
	}
	if err := validateConfig(conf); err != nil {
		return nil, err
	}
	conf = populateConfig(conf)
	tlsConf = tlsConf.Clone()
	setTLSConfigServerName(tlsConf, addr, host)
	return t.doDial(ctx,
		newSendConn(t.conn, addr, packetInfo{}, utils.DefaultLogger),
		tlsConf,
		conf,
		0,
		false,
		use0RTT,
		conf.Versions[0],
	)
}

func (t *Transport) doDial(
	ctx context.Context,
	sendConn sendConn,
	tlsConf *tls.Config,
	config *Config,
	initialPacketNumber protocol.PacketNumber,
	hasNegotiatedVersion bool,
	use0RTT bool,
	version protocol.Version,
) (*Conn, error) {
	srcConnID, err := t.connIDGenerator.GenerateConnectionID()
	if err != nil {
		return nil, err
	}
	destConnID, err := generateConnectionIDForInitial()
	if err != nil {
		return nil, err
	}

	tracingID := nextConnTracingID()
	ctx = context.WithValue(ctx, ConnectionTracingKey, tracingID)

	t.mutex.Lock()
	if t.closeErr != nil {
		t.mutex.Unlock()
		return nil, t.closeErr
	}

	var qlogTrace qlogwriter.Trace
	if config.Tracer != nil {
		qlogTrace = config.Tracer(ctx, true, destConnID)
	}

	logger := utils.DefaultLogger.WithPrefix("client")
	logger.Infof("Starting new connection to %s (%s -> %s), source connection ID %s, destination connection ID %s, version %s", tlsConf.ServerName, sendConn.LocalAddr(), sendConn.RemoteAddr(), srcConnID, destConnID, version)

	conn := newClientConnection(
		context.WithoutCancel(ctx),
		sendConn,
		(*packetHandlerMap)(t),
		destConnID,
		srcConnID,
		t.connIDGenerator,
		t.statelessResetter,
		config,
		tlsConf,
		initialPacketNumber,
		use0RTT,
		hasNegotiatedVersion,
		qlogTrace,
		logger,
		version,
	)
	t.handlers[srcConnID] = conn
	t.mutex.Unlock()

	// The error channel needs to be buffered, as the run loop will continue running
	// after doDial returns (if the handshake is successful).
	// Similarly, the recreateChan needs to be buffered; in case a different case is selected.
	errChan := make(chan error, 1)
	recreateChan := make(chan errCloseForRecreating, 1)
	go func() {
		err := conn.run()
		var recreateErr *errCloseForRecreating
		if errors.As(err, &recreateErr) {
			recreateChan <- *recreateErr
			return
		}
		if t.isSingleUse {
			t.Close()
		}
		errChan <- err
	}()

	// Only set when we're using 0-RTT.
	// Otherwise, earlyConnChan will be nil. Receiving from a nil chan blocks forever.
	var earlyConnChan <-chan struct{}
	if use0RTT {
		earlyConnChan = conn.earlyConnReady()
	}

	select {
	case <-ctx.Done():
		conn.destroy(nil)
		// wait until the Go routine that called Conn.run() returns
		select {
		case <-errChan:
		case <-recreateChan:
		}
		return nil, context.Cause(ctx)
	case params := <-recreateChan:
		return t.doDial(ctx,
			sendConn,
			tlsConf,
			config,
			params.nextPacketNumber,
			true,
			use0RTT,
			params.nextVersion,
		)
	case err := <-errChan:
		return nil, err
	case <-earlyConnChan:
		// ready to send 0-RTT data
		return conn.Conn, nil
	case <-conn.HandshakeComplete():
		// handshake successfully completed
		return conn.Conn, nil
	}
}

func (t *Transport) init(allowZeroLengthConnIDs bool) error {
	t.initOnce.Do(func() {
		var conn rawConn
		if c, ok := t.Conn.(rawConn); ok {
			conn = c
		} else {
			var err error
			conn, err = wrapConn(t.Conn)
			if err != nil {
				t.initErr = err
				return
			}
		}

		t.logger = utils.DefaultLogger // TODO: make this configurable
		t.conn = conn
		t.handlers = make(map[protocol.ConnectionID]packetHandler)
		t.resetTokens = make(map[protocol.StatelessResetToken]packetHandler)
		t.listening = make(chan struct{})

		t.closeQueue = make(chan closePacket, 4)
		t.statelessResetQueue = make(chan receivedPacket, 4)
		if t.TokenGeneratorKey == nil {
			var key TokenGeneratorKey
			if _, err := rand.Read(key[:]); err != nil {
				t.initErr = err
				return
			}
			t.TokenGeneratorKey = &key
		}

		if t.ConnectionIDGenerator != nil {
			t.connIDGenerator = t.ConnectionIDGenerator
			t.connIDLen = t.ConnectionIDGenerator.ConnectionIDLen()
		} else {
			connIDLen := t.ConnectionIDLength
			if t.ConnectionIDLength == 0 && !allowZeroLengthConnIDs {
				connIDLen = protocol.DefaultConnectionIDLength
			}
			t.connIDLen = connIDLen
			t.connIDGenerator = &protocol.DefaultConnectionIDGenerator{ConnLen: t.connIDLen}
		}
		t.statelessResetter = newStatelessResetter(t.StatelessResetKey)

		go func() {
			defer close(t.listening)
			t.listen(conn)

			if t.createdConn {
				conn.Close()
			}
		}()
		go t.runSendQueue()
	})
	return t.initErr
}

// WriteTo sends a packet on the underlying connection.
func (t *Transport) WriteTo(b []byte, addr net.Addr) (int, error) {
	if err := t.init(false); err != nil {
		return 0, err
	}
	return t.conn.WritePacket(b, addr, nil, 0, protocol.ECNUnsupported)
}

func (t *Transport) runSendQueue() {
	for {
		select {
		case <-t.listening:
			return
		case p := <-t.closeQueue:
			t.conn.WritePacket(p.payload, p.addr, p.info.OOB(), 0, protocol.ECNUnsupported)
		case p := <-t.statelessResetQueue:
			t.sendStatelessReset(p)
		}
	}
}

// Close stops listening for UDP datagrams on the Transport.Conn.
// It abruptly terminates all existing connections, without sending a CONNECTION_CLOSE
// to the peers. It is the application's responsibility to cleanly terminate existing
// connections prior to calling Close.
//
// If a server was started, it will be closed as well.
// It is not possible to start any new server or dial new connections after that.
func (t *Transport) Close() error {
	// avoid race condition if the transport is currently being initialized
	t.init(false)

	t.close(nil)
	if t.createdConn {
		if err := t.Conn.Close(); err != nil {
			return err
		}
	} else if t.conn != nil {
		t.conn.SetReadDeadline(time.Now())
		defer func() { t.conn.SetReadDeadline(time.Time{}) }()
	}
	if t.listening != nil {
		<-t.listening // wait until listening returns
	}
	return nil
}

func (t *Transport) closeServer() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.server = nil
	if t.isSingleUse {
		t.closeErr = ErrServerClosed
	}

	if len(t.handlers) == 0 {
		t.maybeStopListening()
	}
}

func (t *Transport) close(e error) {
	t.mutex.Lock()

	if t.closeErr != nil {
		t.mutex.Unlock()
		return
	}

	e = &errTransportClosed{err: e}
	t.closeErr = e
	server := t.server
	t.server = nil
	if server != nil {
		t.mutex.Unlock()
		server.close(e, true)
		t.mutex.Lock()
	}

	// Close existing connections
	var wg sync.WaitGroup
	for _, handler := range t.handlers {
		wg.Add(1)
		go func(handler packetHandler) {
			handler.destroy(e)
			wg.Done()
		}(handler)
	}
	t.mutex.Unlock() // closing connections requires releasing transport mutex
	wg.Wait()

	if t.Tracer != nil {
		t.Tracer.Close()
	}
}

// only print warnings about the UDP receive buffer size once
var setBufferWarningOnce sync.Once

func (t *Transport) listen(conn rawConn) {
	for {
		p, err := conn.ReadPacket()
		//nolint:staticcheck // SA1019 ignore this!
		// TODO: This code is used to ignore wsa errors on Windows.
		// Since net.Error.Temporary is deprecated as of Go 1.18, we should find a better solution.
		// See https://github.com/quic-go/quic-go/issues/1737 for details.
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			t.mutex.Lock()
			closed := t.closeErr != nil
			t.mutex.Unlock()
			if closed {
				return
			}
			t.logger.Debugf("Temporary error reading from conn: %w", err)
			continue
		}
		if err != nil {
			// Windows returns an error when receiving a UDP datagram that doesn't fit into the provided buffer.
			if isRecvMsgSizeErr(err) {
				continue
			}
			t.close(err)
			return
		}
		t.handlePacket(p)
	}
}

func (t *Transport) maybeStopListening() {
	if t.isSingleUse && t.closeErr != nil {
		t.conn.SetReadDeadline(time.Now())
	}
}

func (t *Transport) handlePacket(p receivedPacket) {
	if len(p.data) == 0 {
		return
	}
	if !wire.IsPotentialQUICPacket(p.data[0]) && !wire.IsLongHeaderPacket(p.data[0]) {
		t.handleNonQUICPacket(p)
		return
	}
	connID, err := wire.ParseConnectionID(p.data, t.connIDLen)
	if err != nil {
		t.logger.Debugf("error parsing connection ID on packet from %s: %s", p.remoteAddr, err)
		if t.Tracer != nil {
			t.Tracer.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropHeaderParseError,
			})
		}
		p.buffer.MaybeRelease()
		return
	}

	// If there's a connection associated with the connection ID, pass the packet there.
	if handler, ok := (*packetHandlerMap)(t).Get(connID); ok {
		handler.handlePacket(p)
		return
	}
	// RFC 9000 section 10.3.1 requires that the stateless reset detection logic is run for both
	// packets that cannot be associated with any connections, and for packets that can't be decrypted.
	// We deviate from the RFC and ignore the latter: If a packet's connection ID is associated with an
	// existing connection, it is dropped there if if it can't be decrypted.
	// Stateless resets use random connection IDs, and at reasonable connection ID lengths collisions are
	// exceedingly rare. In the unlikely event that a stateless reset is misrouted to an existing connection,
	// it is to be expected that the next stateless reset will be correctly detected.
	if isStatelessReset := t.maybeHandleStatelessReset(p.data); isStatelessReset {
		return
	}
	if !wire.IsLongHeaderPacket(p.data[0]) {
		if statelessResetQueued := t.maybeSendStatelessReset(p); !statelessResetQueued {
			if t.Tracer != nil {
				t.Tracer.RecordEvent(qlog.PacketDropped{
					Header:  qlog.PacketHeader{PacketType: qlog.PacketType1RTT},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropUnknownConnectionID,
				})
			}
			p.buffer.Release()
		}
		return
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.server == nil { // no server set
		t.logger.Debugf("received a packet with an unexpected connection ID %s", connID)
		if t.Tracer != nil {
			t.Tracer.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnknownConnectionID,
			})
		}
		p.buffer.MaybeRelease()
		return
	}
	t.server.handlePacket(p)
}

func (t *Transport) maybeSendStatelessReset(p receivedPacket) (statelessResetQueued bool) {
	if t.StatelessResetKey == nil {
		return false
	}

	// Don't send a stateless reset in response to very small packets.
	// This includes packets that could be stateless resets.
	if len(p.data) <= protocol.MinStatelessResetSize {
		return false
	}

	select {
	case t.statelessResetQueue <- p:
		return true
	default:
		// it's fine to not send a stateless reset when we're busy
		return false
	}
}

func (t *Transport) sendStatelessReset(p receivedPacket) {
	defer p.buffer.Release()

	connID, err := wire.ParseConnectionID(p.data, t.connIDLen)
	if err != nil {
		t.logger.Errorf("error parsing connection ID on packet from %s: %s", p.remoteAddr, err)
		return
	}
	token := t.statelessResetter.GetStatelessResetToken(connID)
	t.logger.Debugf("Sending stateless reset to %s (connection ID: %s). Token: %#x", p.remoteAddr, connID, token)
	data := make([]byte, protocol.MinStatelessResetSize-16, protocol.MinStatelessResetSize)
	rand.Read(data)
	data[0] = (data[0] & 0x7f) | 0x40
	data = append(data, token[:]...)
	if _, err := t.conn.WritePacket(data, p.remoteAddr, p.info.OOB(), 0, protocol.ECNUnsupported); err != nil {
		t.logger.Debugf("Error sending Stateless Reset to %s: %s", p.remoteAddr, err)
	}
}

func (t *Transport) maybeHandleStatelessReset(data []byte) bool {
	// stateless resets are always short header packets
	if wire.IsLongHeaderPacket(data[0]) {
		return false
	}
	if len(data) < 17 /* type byte + 16 bytes for the reset token */ {
		return false
	}

	token := protocol.StatelessResetToken(data[len(data)-16:])
	t.mutex.Lock()
	conn, ok := t.resetTokens[token]
	t.mutex.Unlock()

	if ok {
		t.logger.Debugf("Received a stateless reset with token %#x. Closing connection.", token)
		go conn.destroy(&StatelessResetError{})
		return true
	}
	return false
}

func (t *Transport) handleNonQUICPacket(p receivedPacket) {
	// Strictly speaking, this is racy,
	// but we only care about receiving packets at some point after ReadNonQUICPacket has been called.
	if !t.readingNonQUICPackets.Load() {
		return
	}
	select {
	case t.nonQUICPackets <- p:
	default:
		if t.Tracer != nil {
			t.Tracer.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropDOSPrevention,
			})
		}
	}
}

const maxQueuedNonQUICPackets = 32

// ReadNonQUICPacket reads non-QUIC packets received on the underlying connection.
// The detection logic is very simple: Any packet that has the first and second bit of the packet set to 0.
// Note that this is stricter than the detection logic defined in RFC 9443.
func (t *Transport) ReadNonQUICPacket(ctx context.Context, b []byte) (int, net.Addr, error) {
	if err := t.init(false); err != nil {
		return 0, nil, err
	}
	if !t.readingNonQUICPackets.Load() {
		t.nonQUICPackets = make(chan receivedPacket, maxQueuedNonQUICPackets)
		t.readingNonQUICPackets.Store(true)
	}
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	case p := <-t.nonQUICPackets:
		n := copy(b, p.data)
		return n, p.remoteAddr, nil
	case <-t.listening:
		return 0, nil, errors.New("closed")
	}
}

func setTLSConfigServerName(tlsConf *tls.Config, addr net.Addr, host string) {
	// If no ServerName is set, infer the ServerName from the host we're connecting to.
	if tlsConf.ServerName != "" {
		return
	}
	if host == "" {
		if udpAddr, ok := addr.(*net.UDPAddr); ok {
			tlsConf.ServerName = udpAddr.IP.String()
			return
		}
	}
	h, _, err := net.SplitHostPort(host)
	if err != nil { // This happens if the host doesn't contain a port number.
		tlsConf.ServerName = host
		return
	}
	tlsConf.ServerName = h
}

type packetHandlerMap Transport

var _ connRunner = &packetHandlerMap{}

func (h *packetHandlerMap) Add(id protocol.ConnectionID, handler packetHandler) bool /* was added */ {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.handlers[id]; ok {
		h.logger.Debugf("Not adding connection ID %s, as it already exists.", id)
		return false
	}
	h.handlers[id] = handler
	h.logger.Debugf("Adding connection ID %s.", id)
	return true
}

func (h *packetHandlerMap) Get(connID protocol.ConnectionID) (packetHandler, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	handler, ok := h.handlers[connID]
	return handler, ok
}

func (h *packetHandlerMap) AddResetToken(token protocol.StatelessResetToken, handler packetHandler) {
	h.mutex.Lock()
	h.resetTokens[token] = handler
	h.mutex.Unlock()
}

func (h *packetHandlerMap) RemoveResetToken(token protocol.StatelessResetToken) {
	h.mutex.Lock()
	delete(h.resetTokens, token)
	h.mutex.Unlock()
}

func (h *packetHandlerMap) AddWithConnID(clientDestConnID, newConnID protocol.ConnectionID, handler packetHandler) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.handlers[clientDestConnID]; ok {
		h.logger.Debugf("Not adding connection ID %s for a new connection, as it already exists.", clientDestConnID)
		return false
	}
	h.handlers[clientDestConnID] = handler
	h.handlers[newConnID] = handler
	h.logger.Debugf("Adding connection IDs %s and %s for a new connection.", clientDestConnID, newConnID)
	return true
}

func (h *packetHandlerMap) Remove(id protocol.ConnectionID) {
	h.mutex.Lock()
	delete(h.handlers, id)
	h.mutex.Unlock()
	h.logger.Debugf("Removing connection ID %s.", id)
}

// ReplaceWithClosed is called when a connection is closed.
// Depending on which side closed the connection, we need to:
// * remote close: absorb delayed packets
// * local close: retransmit the CONNECTION_CLOSE packet, in case it was lost
func (h *packetHandlerMap) ReplaceWithClosed(ids []protocol.ConnectionID, connClosePacket []byte, expiry time.Duration) {
	var handler packetHandler
	if connClosePacket != nil {
		handler = newClosedLocalConn(
			func(addr net.Addr, info packetInfo) {
				select {
				case h.closeQueue <- closePacket{payload: connClosePacket, addr: addr, info: info}:
				default:
					// We're backlogged.
					// Just drop the packet, sending CONNECTION_CLOSE copies is best effort anyway.
				}
			},
			h.logger,
		)
	} else {
		handler = newClosedRemoteConn()
	}

	h.mutex.Lock()
	for _, id := range ids {
		h.handlers[id] = handler
	}
	h.mutex.Unlock()
	h.logger.Debugf("Replacing connection for connection IDs %s with a closed connection.", ids)

	time.AfterFunc(expiry, func() {
		h.mutex.Lock()
		for _, id := range ids {
			delete(h.handlers, id)
		}
		if len(h.handlers) == 0 {
			t := (*Transport)(h)
			t.maybeStopListening()
		}
		h.mutex.Unlock()
		h.logger.Debugf("Removing connection IDs %s for a closed connection after it has been retired.", ids)
	})
}
