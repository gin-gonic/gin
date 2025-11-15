package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go/internal/handshake"
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qlogwriter"
)

// ErrServerClosed is returned by the [Listener] or [EarlyListener]'s Accept method after a call to Close.
var ErrServerClosed = errServerClosed{}

type errServerClosed struct{}

func (errServerClosed) Error() string { return "quic: server closed" }
func (errServerClosed) Unwrap() error { return net.ErrClosed }

// packetHandler handles packets
type packetHandler interface {
	handlePacket(receivedPacket)
	destroy(error)
	closeWithTransportError(qerr.TransportErrorCode)
}

type zeroRTTQueue struct {
	packets    []receivedPacket
	expiration monotime.Time
}

type rejectedPacket struct {
	receivedPacket
	hdr *wire.Header
}

// A Listener of QUIC
type baseServer struct {
	tr                        *packetHandlerMap
	disableVersionNegotiation bool
	acceptEarlyConns          bool

	tlsConf *tls.Config
	config  *Config

	conn rawConn

	tokenGenerator *handshake.TokenGenerator
	maxTokenAge    time.Duration

	connIDGenerator   ConnectionIDGenerator
	statelessResetter *statelessResetter
	onClose           func()

	receivedPackets chan receivedPacket

	nextZeroRTTCleanup monotime.Time
	zeroRTTQueues      map[protocol.ConnectionID]*zeroRTTQueue // only initialized if acceptEarlyConns == true

	connContext func(context.Context, *ClientInfo) (context.Context, error)

	// set as a member, so they can be set in the tests
	newConn func(
		context.Context,
		context.CancelCauseFunc,
		sendConn,
		connRunner,
		protocol.ConnectionID, /* original dest connection ID */
		*protocol.ConnectionID, /* retry src connection ID */
		protocol.ConnectionID, /* client dest connection ID */
		protocol.ConnectionID, /* destination connection ID */
		protocol.ConnectionID, /* source connection ID */
		ConnectionIDGenerator,
		*statelessResetter,
		*Config,
		*tls.Config,
		*handshake.TokenGenerator,
		bool, /* client address validated by an address validation token */
		time.Duration,
		qlogwriter.Trace,
		utils.Logger,
		protocol.Version,
	) *wrappedConn

	closeMx sync.Mutex
	// errorChan is closed when Close is called. This has two effects:
	// 1. it cancels handshakes that are still in flight (using CONNECTION_REFUSED) errors
	// 2. it stops handling of packets passed to this server
	errorChan chan struct{}
	// acceptChan is closed when Close returns.
	// This only happens once all handshake in flight have either completed and canceled.
	// Calls to Accept will first drain the queue of connections that have completed the handshake,
	// and then return ErrServerClosed.
	stopAccepting chan struct{}
	closeErr      error
	running       chan struct{} // closed as soon as run() returns

	versionNegotiationQueue chan receivedPacket
	invalidTokenQueue       chan rejectedPacket
	connectionRefusedQueue  chan rejectedPacket
	retryQueue              chan rejectedPacket
	handshakingCount        sync.WaitGroup

	verifySourceAddress func(net.Addr) bool

	connQueue chan *Conn

	qlogger qlogwriter.Recorder

	logger utils.Logger
}

// A Listener listens for incoming QUIC connections.
// It returns connections once the handshake has completed.
type Listener struct {
	baseServer *baseServer
}

// Accept returns new connections. It should be called in a loop.
func (l *Listener) Accept(ctx context.Context) (*Conn, error) {
	return l.baseServer.Accept(ctx)
}

// Close closes the listener.
// Accept will return [ErrServerClosed] as soon as all connections in the accept queue have been accepted.
// QUIC handshakes that are still in flight will be rejected with a CONNECTION_REFUSED error.
// Already established (accepted) connections will be unaffected.
func (l *Listener) Close() error {
	return l.baseServer.Close()
}

// Addr returns the local network address that the server is listening on.
func (l *Listener) Addr() net.Addr {
	return l.baseServer.Addr()
}

// An EarlyListener listens for incoming QUIC connections, and returns them before the handshake completes.
// For connections that don't use 0-RTT, this allows the server to send 0.5-RTT data.
// This data is encrypted with forward-secure keys, however, the client's identity has not yet been verified.
// For connection using 0-RTT, this allows the server to accept and respond to streams that the client opened in the
// 0-RTT data it sent. Note that at this point during the handshake, the live-ness of the
// client has not yet been confirmed, and the 0-RTT data could have been replayed by an attacker.
type EarlyListener struct {
	baseServer *baseServer
}

// Accept returns a new connections. It should be called in a loop.
func (l *EarlyListener) Accept(ctx context.Context) (*Conn, error) {
	conn, err := l.baseServer.accept(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Close closes the listener.
// Accept will return [ErrServerClosed] as soon as all connections in the accept queue have been accepted.
// Early connections that are still in flight will be rejected with a CONNECTION_REFUSED error.
// Already established (accepted) connections will be unaffected.
func (l *EarlyListener) Close() error {
	return l.baseServer.Close()
}

// Addr returns the local network addr that the server is listening on.
func (l *EarlyListener) Addr() net.Addr {
	return l.baseServer.Addr()
}

// ListenAddr creates a QUIC server listening on a given address.
// See [Listen] for more details.
func ListenAddr(addr string, tlsConf *tls.Config, config *Config) (*Listener, error) {
	conn, err := listenUDP(addr)
	if err != nil {
		return nil, err
	}
	return (&Transport{
		Conn:        conn,
		createdConn: true,
		isSingleUse: true,
	}).Listen(tlsConf, config)
}

// ListenAddrEarly works like [ListenAddr], but it returns connections before the handshake completes.
func ListenAddrEarly(addr string, tlsConf *tls.Config, config *Config) (*EarlyListener, error) {
	conn, err := listenUDP(addr)
	if err != nil {
		return nil, err
	}
	return (&Transport{
		Conn:        conn,
		createdConn: true,
		isSingleUse: true,
	}).ListenEarly(tlsConf, config)
}

func listenUDP(addr string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", udpAddr)
}

// Listen listens for QUIC connections on a given net.PacketConn.
// If the PacketConn satisfies the [OOBCapablePacketConn] interface (as a [net.UDPConn] does),
// ECN and packet info support will be enabled. In this case, ReadMsgUDP and WriteMsgUDP
// will be used instead of ReadFrom and WriteTo to read/write packets.
// A single net.PacketConn can only be used for a single call to Listen.
//
// The tls.Config must not be nil and must contain a certificate configuration.
// Furthermore, it must define an application control (using [NextProtos]).
// The quic.Config may be nil, in that case the default values will be used.
//
// This is a convenience function. More advanced use cases should instantiate a [Transport],
// which offers configuration options for a more fine-grained control of the connection establishment,
// including reusing the underlying UDP socket for outgoing QUIC connections.
// When closing a listener created with Listen, all established QUIC connections will be closed immediately.
func Listen(conn net.PacketConn, tlsConf *tls.Config, config *Config) (*Listener, error) {
	tr := &Transport{Conn: conn, isSingleUse: true}
	return tr.Listen(tlsConf, config)
}

// ListenEarly works like [Listen], but it returns connections before the handshake completes.
func ListenEarly(conn net.PacketConn, tlsConf *tls.Config, config *Config) (*EarlyListener, error) {
	tr := &Transport{Conn: conn, isSingleUse: true}
	return tr.ListenEarly(tlsConf, config)
}

func newServer(
	conn rawConn,
	tr *packetHandlerMap,
	connIDGenerator ConnectionIDGenerator,
	statelessResetter *statelessResetter,
	connContext func(context.Context, *ClientInfo) (context.Context, error),
	tlsConf *tls.Config,
	config *Config,
	qlogger qlogwriter.Recorder,
	onClose func(),
	tokenGeneratorKey TokenGeneratorKey,
	maxTokenAge time.Duration,
	verifySourceAddress func(net.Addr) bool,
	disableVersionNegotiation bool,
	acceptEarly bool,
) *baseServer {
	s := &baseServer{
		conn:                      conn,
		connContext:               connContext,
		tr:                        tr,
		tlsConf:                   tlsConf,
		config:                    config,
		tokenGenerator:            handshake.NewTokenGenerator(tokenGeneratorKey),
		maxTokenAge:               maxTokenAge,
		verifySourceAddress:       verifySourceAddress,
		connIDGenerator:           connIDGenerator,
		statelessResetter:         statelessResetter,
		connQueue:                 make(chan *Conn, protocol.MaxAcceptQueueSize),
		errorChan:                 make(chan struct{}),
		stopAccepting:             make(chan struct{}),
		running:                   make(chan struct{}),
		receivedPackets:           make(chan receivedPacket, protocol.MaxServerUnprocessedPackets),
		versionNegotiationQueue:   make(chan receivedPacket, 4),
		invalidTokenQueue:         make(chan rejectedPacket, 4),
		connectionRefusedQueue:    make(chan rejectedPacket, 4),
		retryQueue:                make(chan rejectedPacket, 8),
		newConn:                   newConnection,
		qlogger:                   qlogger,
		logger:                    utils.DefaultLogger.WithPrefix("server"),
		acceptEarlyConns:          acceptEarly,
		disableVersionNegotiation: disableVersionNegotiation,
		onClose:                   onClose,
	}
	if acceptEarly {
		s.zeroRTTQueues = map[protocol.ConnectionID]*zeroRTTQueue{}
	}
	go s.run()
	go s.runSendQueue()
	s.logger.Debugf("Listening for %s connections on %s", conn.LocalAddr().Network(), conn.LocalAddr().String())
	return s
}

func (s *baseServer) run() {
	defer close(s.running)
	for {
		select {
		case <-s.errorChan:
			return
		default:
		}
		select {
		case <-s.errorChan:
			return
		case p := <-s.receivedPackets:
			if bufferStillInUse := s.handlePacketImpl(p); !bufferStillInUse {
				p.buffer.Release()
			}
		}
	}
}

func (s *baseServer) runSendQueue() {
	for {
		select {
		case <-s.running:
			return
		case p := <-s.versionNegotiationQueue:
			s.maybeSendVersionNegotiationPacket(p)
		case p := <-s.invalidTokenQueue:
			s.maybeSendInvalidToken(p)
		case p := <-s.connectionRefusedQueue:
			s.sendConnectionRefused(p)
		case p := <-s.retryQueue:
			s.sendRetry(p)
		}
	}
}

// Accept returns connections that already completed the handshake.
// It is only valid if acceptEarlyConns is false.
func (s *baseServer) Accept(ctx context.Context) (*Conn, error) {
	return s.accept(ctx)
}

func (s *baseServer) accept(ctx context.Context) (*Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case conn := <-s.connQueue:
		return conn, nil
	case <-s.stopAccepting:
		// first drain the queue
		select {
		case conn := <-s.connQueue:
			return conn, nil
		default:
		}
		return nil, s.closeErr
	}
}

func (s *baseServer) Close() error {
	s.close(ErrServerClosed, false)
	return nil
}

// close closes the server. The Transport mutex must not be held while calling this method.
// This method closes any handshaking connections which requires the tranpsort mutex.
func (s *baseServer) close(e error, transportClose bool) {
	s.closeMx.Lock()
	if s.closeErr != nil {
		s.closeMx.Unlock()
		return
	}
	s.closeErr = e
	close(s.errorChan)
	<-s.running
	s.closeMx.Unlock()

	if !transportClose {
		s.onClose()
	}

	// wait until all handshakes in flight have terminated
	s.handshakingCount.Wait()
	close(s.stopAccepting)

	if transportClose {
		// if the transport is closing, drain the connQueue. All connections in the queue
		// will be closed by the transport.
		for {
			select {
			case <-s.connQueue:
			default:
				return
			}
		}
	}
}

// Addr returns the server's network address
func (s *baseServer) Addr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *baseServer) handlePacket(p receivedPacket) {
	select {
	case s.receivedPackets <- p:
	case <-s.errorChan:
		return
	default:
		s.logger.Debugf("Dropping packet from %s (%d bytes). Server receive queue full.", p.remoteAddr, p.Size())
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropDOSPrevention,
			})
		}
	}
}

func (s *baseServer) handlePacketImpl(p receivedPacket) bool /* is the buffer still in use? */ {
	if !s.nextZeroRTTCleanup.IsZero() && p.rcvTime.After(s.nextZeroRTTCleanup) {
		defer s.cleanupZeroRTTQueues(p.rcvTime)
	}

	if wire.IsVersionNegotiationPacket(p.data) {
		s.logger.Debugf("Dropping Version Negotiation packet.")
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header:  qlog.PacketHeader{PacketType: qlog.PacketTypeVersionNegotiation},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return false
	}
	// Short header packets should never end up here in the first place
	if !wire.IsLongHeaderPacket(p.data[0]) {
		panic(fmt.Sprintf("misrouted packet: %#v", p.data))
	}
	v, err := wire.ParseVersion(p.data)
	// drop the packet if we failed to parse the protocol version
	if err != nil {
		s.logger.Debugf("Dropping a packet with an unknown version")
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return false
	}
	// send a Version Negotiation Packet if the client is speaking a different protocol version
	if !protocol.IsSupportedVersion(s.config.Versions, v) {
		if s.disableVersionNegotiation {
			if s.qlogger != nil {
				s.qlogger.RecordEvent(qlog.PacketDropped{
					Header:  qlog.PacketHeader{Version: v},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropUnexpectedVersion,
				})
			}
			return false
		}

		if p.Size() < protocol.MinUnknownVersionPacketSize {
			s.logger.Debugf("Dropping a packet with an unsupported version number %d that is too small (%d bytes)", v, p.Size())
			if s.qlogger != nil {
				s.qlogger.RecordEvent(qlog.PacketDropped{
					Header:  qlog.PacketHeader{Version: v},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropUnexpectedPacket,
				})
			}
			return false
		}
		return s.enqueueVersionNegotiationPacket(p)
	}

	if wire.Is0RTTPacket(p.data) {
		if !s.acceptEarlyConns {
			if s.qlogger != nil {
				s.qlogger.RecordEvent(qlog.PacketDropped{
					Header: qlog.PacketHeader{
						PacketType:   qlog.PacketType0RTT,
						PacketNumber: protocol.InvalidPacketNumber,
					},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropUnexpectedPacket,
				})
			}
			return false
		}
		return s.handle0RTTPacket(p)
	}

	// If we're creating a new connection, the packet will be passed to the connection.
	// The header will then be parsed again.
	hdr, _, _, err := wire.ParsePacket(p.data)
	if err != nil {
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropHeaderParseError,
			})
		}
		s.logger.Debugf("Error parsing packet: %s", err)
		return false
	}
	if hdr.Type == protocol.PacketTypeInitial && p.Size() < protocol.MinInitialPacketSize {
		s.logger.Debugf("Dropping a packet that is too small to be a valid Initial (%d bytes)", p.Size())
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketTypeInitial,
					PacketNumber: protocol.InvalidPacketNumber,
					Version:      v,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return false
	}

	if hdr.Type != protocol.PacketTypeInitial {
		// Drop long header packets.
		// There's little point in sending a Stateless Reset, since the client
		// might not have received the token yet.
		s.logger.Debugf("Dropping long header packet of type %s (%d bytes)", hdr.Type, len(p.data))
		if s.qlogger != nil {
			var pt qlog.PacketType
			switch hdr.Type {
			case protocol.PacketTypeInitial:
				pt = qlog.PacketTypeInitial
			case protocol.PacketTypeHandshake:
				pt = qlog.PacketTypeHandshake
			case protocol.PacketType0RTT:
				pt = qlog.PacketType0RTT
			case protocol.PacketTypeRetry:
				pt = qlog.PacketTypeRetry
			}
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   pt,
					PacketNumber: protocol.InvalidPacketNumber,
					Version:      v,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return false
	}

	s.logger.Debugf("<- Received Initial packet.")

	if err := s.handleInitialImpl(p, hdr); err != nil {
		s.logger.Errorf("Error occurred handling initial packet: %s", err)
	}
	// Don't put the packet buffer back.
	// handleInitialImpl deals with the buffer.
	return true
}

func (s *baseServer) handle0RTTPacket(p receivedPacket) bool {
	connID, err := wire.ParseConnectionID(p.data, 0)
	if err != nil {
		if s.qlogger != nil {
			v, _ := wire.ParseVersion(p.data)
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketType0RTT,
					PacketNumber: protocol.InvalidPacketNumber,
					Version:      v,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropHeaderParseError,
			})
		}
		return false
	}

	// check again if we might have a connection now
	if handler, ok := s.tr.Get(connID); ok {
		handler.handlePacket(p)
		return true
	}

	if q, ok := s.zeroRTTQueues[connID]; ok {
		if len(q.packets) >= protocol.Max0RTTQueueLen {
			if s.qlogger != nil {
				v, _ := wire.ParseVersion(p.data)
				s.qlogger.RecordEvent(qlog.PacketDropped{
					Header: qlog.PacketHeader{
						PacketType:   qlog.PacketType0RTT,
						PacketNumber: protocol.InvalidPacketNumber,
						Version:      v,
					},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropDOSPrevention,
				})
			}
			return false
		}
		q.packets = append(q.packets, p)
		return true
	}

	if len(s.zeroRTTQueues) >= protocol.Max0RTTQueues {
		if s.qlogger != nil {
			v, _ := wire.ParseVersion(p.data)
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketType0RTT,
					PacketNumber: protocol.InvalidPacketNumber,
					Version:      v,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropDOSPrevention,
			})
		}
		return false
	}
	queue := &zeroRTTQueue{packets: make([]receivedPacket, 1, 8)}
	queue.packets[0] = p
	expiration := p.rcvTime.Add(protocol.Max0RTTQueueingDuration)
	queue.expiration = expiration
	if s.nextZeroRTTCleanup.IsZero() || s.nextZeroRTTCleanup.After(expiration) {
		s.nextZeroRTTCleanup = expiration
	}
	s.zeroRTTQueues[connID] = queue
	return true
}

func (s *baseServer) cleanupZeroRTTQueues(now monotime.Time) {
	// Iterate over all queues to find those that are expired.
	// This is ok since we're placing a pretty low limit on the number of queues.
	var nextCleanup monotime.Time
	for connID, q := range s.zeroRTTQueues {
		if q.expiration.After(now) {
			if nextCleanup.IsZero() || nextCleanup.After(q.expiration) {
				nextCleanup = q.expiration
			}
			continue
		}
		for _, p := range q.packets {
			if s.qlogger != nil {
				v, _ := wire.ParseVersion(p.data)
				s.qlogger.RecordEvent(qlog.PacketDropped{
					Header: qlog.PacketHeader{
						PacketType:   qlog.PacketType0RTT,
						PacketNumber: protocol.InvalidPacketNumber,
						Version:      v,
					},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropDOSPrevention,
				})
			}
			p.buffer.Release()
		}
		delete(s.zeroRTTQueues, connID)
		if s.logger.Debug() {
			s.logger.Debugf("Removing 0-RTT queue for %s.", connID)
		}
	}
	s.nextZeroRTTCleanup = nextCleanup
}

// validateToken returns false if:
//   - address is invalid
//   - token is expired
//   - token is null
func (s *baseServer) validateToken(token *handshake.Token, addr net.Addr) bool {
	if token == nil {
		return false
	}
	if !token.ValidateRemoteAddr(addr) {
		return false
	}
	if !token.IsRetryToken && time.Since(token.SentTime) > s.maxTokenAge {
		return false
	}
	if token.IsRetryToken && time.Since(token.SentTime) > s.config.maxRetryTokenAge() {
		return false
	}
	return true
}

func (s *baseServer) handleInitialImpl(p receivedPacket, hdr *wire.Header) error {
	if len(hdr.Token) == 0 && hdr.DestConnectionID.Len() < protocol.MinConnectionIDLenInitial {
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketTypeInitial,
					PacketNumber: protocol.InvalidPacketNumber,
					Version:      hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		p.buffer.Release()
		return errors.New("too short connection ID")
	}

	// The server queues packets for a while, and we might already have established a connection by now.
	// This results in a second check in the connection map.
	// That's ok since it's not the hot path (it's only taken by some Initial and 0-RTT packets).
	if handler, ok := s.tr.Get(hdr.DestConnectionID); ok {
		handler.handlePacket(p)
		return nil
	}

	var (
		token              *handshake.Token
		retrySrcConnID     *protocol.ConnectionID
		clientAddrVerified bool
	)
	origDestConnID := hdr.DestConnectionID
	if len(hdr.Token) > 0 {
		tok, err := s.tokenGenerator.DecodeToken(hdr.Token)
		if err == nil {
			if tok.IsRetryToken {
				origDestConnID = tok.OriginalDestConnectionID
				retrySrcConnID = &tok.RetrySrcConnectionID
			}
			token = tok
		}
	}
	if token != nil {
		clientAddrVerified = s.validateToken(token, p.remoteAddr)
		if !clientAddrVerified {
			// For invalid and expired non-retry tokens, we don't send an INVALID_TOKEN error.
			// We just ignore them, and act as if there was no token on this packet at all.
			// This also means we might send a Retry later.
			if !token.IsRetryToken {
				token = nil
			} else {
				// For Retry tokens, we send an INVALID_ERROR if
				// * the token is too old, or
				// * the token is invalid, in case of a retry token.
				select {
				case s.invalidTokenQueue <- rejectedPacket{receivedPacket: p, hdr: hdr}:
				default:
					// drop packet if we can't send out the  INVALID_TOKEN packets fast enough
					p.buffer.Release()
				}
				return nil
			}
		}
	}

	if token == nil && s.verifySourceAddress != nil && s.verifySourceAddress(p.remoteAddr) {
		// Retry invalidates all 0-RTT packets sent.
		delete(s.zeroRTTQueues, hdr.DestConnectionID)
		select {
		case s.retryQueue <- rejectedPacket{receivedPacket: p, hdr: hdr}:
		default:
			// drop packet if we can't send out Retry packets fast enough
			p.buffer.Release()
		}
		return nil
	}

	// restore RTT from token
	var rtt time.Duration
	if token != nil && !token.IsRetryToken {
		rtt = token.RTT
	}

	config := s.config
	clientInfo := &ClientInfo{
		RemoteAddr:   p.remoteAddr,
		AddrVerified: clientAddrVerified,
	}
	if s.config.GetConfigForClient != nil {
		conf, err := s.config.GetConfigForClient(clientInfo)
		if err != nil {
			s.logger.Debugf("Rejecting new connection due to GetConfigForClient callback")
			s.refuseNewConn(p, hdr)
			return nil
		}
		config = populateConfig(conf)
	}

	var conn *wrappedConn
	var cancel context.CancelCauseFunc
	ctx, cancel1 := context.WithCancelCause(context.Background())
	if s.connContext != nil {
		var err error
		ctx, err = s.connContext(ctx, clientInfo)
		if err != nil {
			cancel1(err)
			s.logger.Debugf("Rejecting new connection due to ConnContext callback: %s", err)
			s.refuseNewConn(p, hdr)
			return nil
		}
		if ctx == nil {
			panic("quic: ConnContext returned nil")
		}
		// There's no guarantee that the application returns a context
		// that's derived from the context we passed into ConnContext.
		// We need to make sure that both contexts are cancelled.
		var cancel2 context.CancelCauseFunc
		ctx, cancel2 = context.WithCancelCause(ctx)
		cancel = func(cause error) {
			cancel1(cause)
			cancel2(cause)
		}
	} else {
		cancel = cancel1
	}
	ctx = context.WithValue(ctx, ConnectionTracingKey, nextConnTracingID())
	var qlogTrace qlogwriter.Trace
	if config.Tracer != nil {
		// Use the same connection ID that is passed to the client's GetLogWriter callback.
		connID := hdr.DestConnectionID
		if origDestConnID.Len() > 0 {
			connID = origDestConnID
		}
		qlogTrace = config.Tracer(ctx, false, connID)
	}
	connID, err := s.connIDGenerator.GenerateConnectionID()
	if err != nil {
		return err
	}
	s.logger.Debugf("Changing connection ID to %s.", connID)
	conn = s.newConn(
		ctx,
		cancel,
		newSendConn(s.conn, p.remoteAddr, p.info, s.logger),
		s.tr,
		origDestConnID,
		retrySrcConnID,
		hdr.DestConnectionID,
		hdr.SrcConnectionID,
		connID,
		s.connIDGenerator,
		s.statelessResetter,
		config,
		s.tlsConf,
		s.tokenGenerator,
		clientAddrVerified,
		rtt,
		qlogTrace,
		s.logger,
		hdr.Version,
	)
	conn.handlePacket(p)
	// Adding the connection will fail if the client's chosen Destination Connection ID is already in use.
	// This is very unlikely: Even if an attacker chooses a connection ID that's already in use,
	// under normal circumstances the packet would just be routed to that connection.
	// The only time this collision will occur if we receive the two Initial packets at the same time.
	if added := s.tr.AddWithConnID(hdr.DestConnectionID, connID, conn); !added {
		delete(s.zeroRTTQueues, hdr.DestConnectionID)
		conn.closeWithTransportError(ConnectionRefused)
		return nil
	}
	// Pass queued 0-RTT to the newly established connection.
	if q, ok := s.zeroRTTQueues[hdr.DestConnectionID]; ok {
		for _, p := range q.packets {
			conn.handlePacket(p)
		}
		delete(s.zeroRTTQueues, hdr.DestConnectionID)
	}

	s.handshakingCount.Add(1)
	go func() {
		defer s.handshakingCount.Done()
		s.handleNewConn(conn)
	}()
	go conn.run()
	return nil
}

func (s *baseServer) refuseNewConn(p receivedPacket, hdr *wire.Header) {
	delete(s.zeroRTTQueues, hdr.DestConnectionID)
	select {
	case s.connectionRefusedQueue <- rejectedPacket{receivedPacket: p, hdr: hdr}:
	default:
		// drop packet if we can't send out the CONNECTION_REFUSED fast enough
		p.buffer.Release()
	}
}

func (s *baseServer) handleNewConn(conn *wrappedConn) {
	if s.acceptEarlyConns {
		// wait until the early connection is ready, the handshake fails, or the server is closed
		select {
		case <-s.errorChan:
			conn.closeWithTransportError(ConnectionRefused)
			return
		case <-conn.Context().Done():
			return
		case <-conn.earlyConnReady():
		}
	} else {
		// wait until the handshake completes, fails, or the server is closed
		select {
		case <-s.errorChan:
			conn.closeWithTransportError(ConnectionRefused)
			return
		case <-conn.Context().Done():
			return
		case <-conn.HandshakeComplete():
		}
	}

	select {
	case s.connQueue <- conn.Conn:
	default:
		conn.closeWithTransportError(ConnectionRefused)
	}
}

func (s *baseServer) sendRetry(p rejectedPacket) {
	if err := s.sendRetryPacket(p); err != nil {
		s.logger.Debugf("Error sending Retry packet: %s", err)
	}
}

func (s *baseServer) sendRetryPacket(p rejectedPacket) error {
	hdr := p.hdr
	// Log the Initial packet now.
	// If no Retry is sent, the packet will be logged by the connection.
	(&wire.ExtendedHeader{Header: *hdr}).Log(s.logger)
	srcConnID, err := s.connIDGenerator.GenerateConnectionID()
	if err != nil {
		return err
	}
	token, err := s.tokenGenerator.NewRetryToken(p.remoteAddr, hdr.DestConnectionID, srcConnID)
	if err != nil {
		return err
	}
	replyHdr := &wire.ExtendedHeader{}
	replyHdr.Type = protocol.PacketTypeRetry
	replyHdr.Version = hdr.Version
	replyHdr.SrcConnectionID = srcConnID
	replyHdr.DestConnectionID = hdr.SrcConnectionID
	replyHdr.Token = token
	if s.logger.Debug() {
		s.logger.Debugf("Changing connection ID to %s.", srcConnID)
		s.logger.Debugf("-> Sending Retry")
		replyHdr.Log(s.logger)
	}

	buf := getPacketBuffer()
	defer buf.Release()
	buf.Data, err = replyHdr.Append(buf.Data, hdr.Version)
	if err != nil {
		return err
	}
	// append the Retry integrity tag
	tag := handshake.GetRetryIntegrityTag(buf.Data, hdr.DestConnectionID, hdr.Version)
	buf.Data = append(buf.Data, tag[:]...)
	if s.qlogger != nil {
		s.qlogger.RecordEvent(qlog.PacketSent{
			Header: qlog.PacketHeader{
				PacketType:       qlog.PacketTypeRetry,
				SrcConnectionID:  replyHdr.SrcConnectionID,
				DestConnectionID: replyHdr.DestConnectionID,
				Version:          replyHdr.Version,
				Token:            &qlog.Token{Raw: token},
			},
			Raw: qlog.RawInfo{
				Length:        len(buf.Data),
				PayloadLength: int(replyHdr.Length),
			},
		})
	}
	_, err = s.conn.WritePacket(buf.Data, p.remoteAddr, p.info.OOB(), 0, protocol.ECNUnsupported)
	return err
}

func (s *baseServer) maybeSendInvalidToken(p rejectedPacket) {
	defer p.buffer.Release()

	// Only send INVALID_TOKEN if we can unprotect the packet.
	// This makes sure that we won't send it for packets that were corrupted.
	hdr := p.hdr
	sealer, opener := handshake.NewInitialAEAD(hdr.DestConnectionID, protocol.PerspectiveServer, hdr.Version)
	data := p.data[:hdr.ParsedLen()+hdr.Length]
	extHdr, err := unpackLongHeader(opener, hdr, data)
	// Only send INVALID_TOKEN if we can unprotect the packet.
	// This makes sure that we won't send it for packets that were corrupted.
	if err != nil {
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketTypeInitial,
					PacketNumber: protocol.InvalidPacketNumber,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropHeaderParseError,
			})
		}
		return
	}
	hdrLen := extHdr.ParsedLen()
	if _, err := opener.Open(data[hdrLen:hdrLen], data[hdrLen:], extHdr.PacketNumber, data[:hdrLen]); err != nil {
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketTypeInitial,
					PacketNumber: protocol.InvalidPacketNumber,
					Version:      hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropPayloadDecryptError,
			})
		}
		return
	}
	if s.logger.Debug() {
		s.logger.Debugf("Client sent an invalid retry token. Sending INVALID_TOKEN to %s.", p.remoteAddr)
	}
	if err := s.sendError(p.remoteAddr, hdr, sealer, InvalidToken, p.info); err != nil {
		s.logger.Debugf("Error sending INVALID_TOKEN error: %s", err)
	}
}

func (s *baseServer) sendConnectionRefused(p rejectedPacket) {
	defer p.buffer.Release()
	sealer, _ := handshake.NewInitialAEAD(p.hdr.DestConnectionID, protocol.PerspectiveServer, p.hdr.Version)
	if err := s.sendError(p.remoteAddr, p.hdr, sealer, ConnectionRefused, p.info); err != nil {
		s.logger.Debugf("Error sending CONNECTION_REFUSED error: %s", err)
	}
}

// sendError sends the error as a response to the packet received with header hdr
func (s *baseServer) sendError(remoteAddr net.Addr, hdr *wire.Header, sealer handshake.LongHeaderSealer, errorCode qerr.TransportErrorCode, info packetInfo) error {
	b := getPacketBuffer()
	defer b.Release()

	ccf := &wire.ConnectionCloseFrame{ErrorCode: uint64(errorCode)}

	replyHdr := &wire.ExtendedHeader{}
	replyHdr.Type = protocol.PacketTypeInitial
	replyHdr.Version = hdr.Version
	replyHdr.SrcConnectionID = hdr.DestConnectionID
	replyHdr.DestConnectionID = hdr.SrcConnectionID
	replyHdr.PacketNumberLen = protocol.PacketNumberLen4
	replyHdr.Length = 4 /* packet number len */ + ccf.Length(hdr.Version) + protocol.ByteCount(sealer.Overhead())
	var err error
	b.Data, err = replyHdr.Append(b.Data, hdr.Version)
	if err != nil {
		return err
	}
	payloadOffset := len(b.Data)

	b.Data, err = ccf.Append(b.Data, hdr.Version)
	if err != nil {
		return err
	}

	_ = sealer.Seal(b.Data[payloadOffset:payloadOffset], b.Data[payloadOffset:], replyHdr.PacketNumber, b.Data[:payloadOffset])
	b.Data = b.Data[0 : len(b.Data)+sealer.Overhead()]

	pnOffset := payloadOffset - int(replyHdr.PacketNumberLen)
	sealer.EncryptHeader(
		b.Data[pnOffset+4:pnOffset+4+16],
		&b.Data[0],
		b.Data[pnOffset:payloadOffset],
	)

	replyHdr.Log(s.logger)
	wire.LogFrame(s.logger, ccf, true)
	if s.qlogger != nil {
		s.qlogger.RecordEvent(qlog.PacketSent{
			Header: qlog.PacketHeader{
				PacketType:       qlog.PacketTypeInitial,
				SrcConnectionID:  replyHdr.SrcConnectionID,
				DestConnectionID: replyHdr.DestConnectionID,
				PacketNumber:     replyHdr.PacketNumber,
				Version:          replyHdr.Version,
			},
			Raw: qlog.RawInfo{
				Length:        len(b.Data),
				PayloadLength: int(replyHdr.Length),
			},
			Frames: []qlog.Frame{{Frame: ccf}},
		})
	}
	_, err = s.conn.WritePacket(b.Data, remoteAddr, info.OOB(), 0, protocol.ECNUnsupported)
	return err
}

func (s *baseServer) enqueueVersionNegotiationPacket(p receivedPacket) (bufferInUse bool) {
	select {
	case s.versionNegotiationQueue <- p:
		return true
	default:
		// it's fine to not send version negotiation packets when we are busy
	}
	return false
}

func (s *baseServer) maybeSendVersionNegotiationPacket(p receivedPacket) {
	defer p.buffer.Release()

	v, err := wire.ParseVersion(p.data)
	if err != nil {
		s.logger.Debugf("failed to parse version for sending version negotiation packet: %s", err)
		return
	}

	_, src, dest, err := wire.ParseArbitraryLenConnectionIDs(p.data)
	if err != nil { // should never happen
		s.logger.Debugf("Dropping a packet with an unknown version for which we failed to parse connection IDs")
		if s.qlogger != nil {
			s.qlogger.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return
	}

	s.logger.Debugf("Client offered version %s, sending Version Negotiation", v)

	data := wire.ComposeVersionNegotiation(dest, src, s.config.Versions)
	if s.qlogger != nil {
		s.qlogger.RecordEvent(qlog.VersionNegotiationSent{
			Header: qlog.PacketHeaderVersionNegotiation{
				SrcConnectionID:  src,
				DestConnectionID: dest,
			},
			SupportedVersions: s.config.Versions,
		})
	}
	if _, err := s.conn.WritePacket(data, p.remoteAddr, p.info.OOB(), 0, protocol.ECNUnsupported); err != nil {
		s.logger.Debugf("Error sending Version Negotiation: %s", err)
	}
}
