package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/flowcontrol"
	"github.com/quic-go/quic-go/internal/handshake"
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/internal/utils/ringbuffer"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qlogwriter"
)

type unpacker interface {
	UnpackLongHeader(hdr *wire.Header, data []byte) (*unpackedPacket, error)
	UnpackShortHeader(rcvTime monotime.Time, data []byte) (protocol.PacketNumber, protocol.PacketNumberLen, protocol.KeyPhaseBit, []byte, error)
}

type cryptoStreamHandler interface {
	StartHandshake(context.Context) error
	ChangeConnectionID(protocol.ConnectionID)
	SetLargest1RTTAcked(protocol.PacketNumber) error
	SetHandshakeConfirmed()
	GetSessionTicket() ([]byte, error)
	NextEvent() handshake.Event
	DiscardInitialKeys()
	HandleMessage([]byte, protocol.EncryptionLevel) error
	io.Closer
	ConnectionState() handshake.ConnectionState
}

type receivedPacket struct {
	buffer *packetBuffer

	remoteAddr net.Addr
	rcvTime    monotime.Time
	data       []byte

	ecn protocol.ECN

	info packetInfo // only valid if the contained IP address is valid
}

func (p *receivedPacket) Size() protocol.ByteCount { return protocol.ByteCount(len(p.data)) }

func (p *receivedPacket) Clone() *receivedPacket {
	return &receivedPacket{
		remoteAddr: p.remoteAddr,
		rcvTime:    p.rcvTime,
		data:       p.data,
		buffer:     p.buffer,
		ecn:        p.ecn,
		info:       p.info,
	}
}

type connRunner interface {
	Add(protocol.ConnectionID, packetHandler) bool
	Remove(protocol.ConnectionID)
	ReplaceWithClosed([]protocol.ConnectionID, []byte, time.Duration)
	AddResetToken(protocol.StatelessResetToken, packetHandler)
	RemoveResetToken(protocol.StatelessResetToken)
}

type closeError struct {
	err       error
	immediate bool
}

type errCloseForRecreating struct {
	nextPacketNumber protocol.PacketNumber
	nextVersion      protocol.Version
}

func (e *errCloseForRecreating) Error() string {
	return "closing connection in order to recreate it"
}

var deadlineSendImmediately = monotime.Time(42 * time.Millisecond) // any value > time.Time{} and before time.Now() is fine

var connTracingID atomic.Uint64              // to be accessed atomically
func nextConnTracingID() ConnectionTracingID { return ConnectionTracingID(connTracingID.Add(1)) }

type blockMode uint8

const (
	// blockModeNone means that the connection is not blocked.
	blockModeNone blockMode = iota
	// blockModeCongestionLimited means that the connection is congestion limited.
	// In that case, we can still send acknowledgments and PTO probe packets.
	blockModeCongestionLimited
	// blockModeHardBlocked means that no packet can be sent, under no circumstances. This can happen when:
	// * the send queue is full
	// * the SentPacketHandler returns SendNone, e.g. when we are tracking the maximum number of packets
	// In that case, the timer will be set to the idle timeout.
	blockModeHardBlocked
)

// A Conn is a QUIC connection between two peers.
// Calls to the connection (and to streams) can return the following types of errors:
//   - [ApplicationError]: for errors triggered by the application running on top of QUIC
//   - [TransportError]: for errors triggered by the QUIC transport (in many cases a misbehaving peer)
//   - [IdleTimeoutError]: when the peer goes away unexpectedly (this is a [net.Error] timeout error)
//   - [HandshakeTimeoutError]: when the cryptographic handshake takes too long (this is a [net.Error] timeout error)
//   - [StatelessResetError]: when we receive a stateless reset
//   - [VersionNegotiationError]: returned by the client, when there's no version overlap between the peers
type Conn struct {
	// Destination connection ID used during the handshake.
	// Used to check source connection ID on incoming packets.
	handshakeDestConnID protocol.ConnectionID
	// Set for the client. Destination connection ID used on the first Initial sent.
	origDestConnID protocol.ConnectionID
	retrySrcConnID *protocol.ConnectionID // only set for the client (and if a Retry was performed)

	srcConnIDLen int

	perspective protocol.Perspective
	version     protocol.Version
	config      *Config

	conn      sendConn
	sendQueue sender

	// lazily initialzed: most connections never migrate
	pathManager         *pathManager
	largestRcvdAppData  protocol.PacketNumber
	pathManagerOutgoing atomic.Pointer[pathManagerOutgoing]

	streamsMap      *streamsMap
	connIDManager   *connIDManager
	connIDGenerator *connIDGenerator

	rttStats  *utils.RTTStats
	connStats utils.ConnectionStats

	cryptoStreamManager   *cryptoStreamManager
	sentPacketHandler     ackhandler.SentPacketHandler
	receivedPacketHandler ackhandler.ReceivedPacketHandler
	retransmissionQueue   *retransmissionQueue
	framer                *framer
	connFlowController    flowcontrol.ConnectionFlowController
	tokenStoreKey         string                    // only set for the client
	tokenGenerator        *handshake.TokenGenerator // only set for the server

	unpacker      unpacker
	frameParser   wire.FrameParser
	packer        packer
	mtuDiscoverer mtuDiscoverer // initialized when the transport parameters are received

	currentMTUEstimate atomic.Uint32

	initialStream       *initialCryptoStream
	handshakeStream     *cryptoStream
	oneRTTStream        *cryptoStream // only set for the server
	cryptoStreamHandler cryptoStreamHandler

	notifyReceivedPacket chan struct{}
	sendingScheduled     chan struct{}
	receivedPacketMx     sync.Mutex
	receivedPackets      ringbuffer.RingBuffer[receivedPacket]

	// closeChan is used to notify the run loop that it should terminate
	closeChan chan struct{}
	closeErr  atomic.Pointer[closeError]

	ctx                   context.Context
	ctxCancel             context.CancelCauseFunc
	handshakeCompleteChan chan struct{}

	undecryptablePackets          []receivedPacket // undecryptable packets, waiting for a change in encryption level
	undecryptablePacketsToProcess []receivedPacket

	earlyConnReadyChan chan struct{}
	sentFirstPacket    bool
	droppedInitialKeys bool
	handshakeComplete  bool
	handshakeConfirmed bool

	receivedRetry       bool
	versionNegotiated   bool
	receivedFirstPacket bool

	blocked blockMode

	// the minimum of the max_idle_timeout values advertised by both endpoints
	idleTimeout  time.Duration
	creationTime monotime.Time
	// The idle timeout is set based on the max of the time we received the last packet...
	lastPacketReceivedTime monotime.Time
	// ... and the time we sent a new ack-eliciting packet after receiving a packet.
	firstAckElicitingPacketAfterIdleSentTime monotime.Time
	// pacingDeadline is the time when the next packet should be sent
	pacingDeadline monotime.Time

	peerParams *wire.TransportParameters

	timer *time.Timer
	// keepAlivePingSent stores whether a keep alive PING is in flight.
	// It is reset as soon as we receive a packet from the peer.
	keepAlivePingSent bool
	keepAliveInterval time.Duration

	datagramQueue *datagramQueue

	connStateMutex sync.Mutex
	connState      ConnectionState

	logID     string
	qlogTrace qlogwriter.Trace
	qlogger   qlogwriter.Recorder
	logger    utils.Logger
}

var _ streamSender = &Conn{}

type connTestHooks struct {
	run                     func() error
	earlyConnReady          func() <-chan struct{}
	context                 func() context.Context
	handshakeComplete       func() <-chan struct{}
	closeWithTransportError func(TransportErrorCode)
	destroy                 func(error)
	handlePacket            func(receivedPacket)
}

type wrappedConn struct {
	testHooks *connTestHooks
	*Conn
}

var newConnection = func(
	ctx context.Context,
	ctxCancel context.CancelCauseFunc,
	conn sendConn,
	runner connRunner,
	origDestConnID protocol.ConnectionID,
	retrySrcConnID *protocol.ConnectionID,
	clientDestConnID protocol.ConnectionID,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	connIDGenerator ConnectionIDGenerator,
	statelessResetter *statelessResetter,
	conf *Config,
	tlsConf *tls.Config,
	tokenGenerator *handshake.TokenGenerator,
	clientAddressValidated bool,
	rtt time.Duration,
	qlogTrace qlogwriter.Trace,
	logger utils.Logger,
	v protocol.Version,
) *wrappedConn {
	s := &Conn{
		ctx:                 ctx,
		ctxCancel:           ctxCancel,
		conn:                conn,
		config:              conf,
		handshakeDestConnID: destConnID,
		srcConnIDLen:        srcConnID.Len(),
		tokenGenerator:      tokenGenerator,
		oneRTTStream:        newCryptoStream(),
		perspective:         protocol.PerspectiveServer,
		qlogTrace:           qlogTrace,
		logger:              logger,
		version:             v,
	}
	if qlogTrace != nil {
		s.qlogger = qlogTrace.AddProducer()
	}
	if origDestConnID.Len() > 0 {
		s.logID = origDestConnID.String()
	} else {
		s.logID = destConnID.String()
	}
	s.connIDManager = newConnIDManager(
		destConnID,
		func(token protocol.StatelessResetToken) { runner.AddResetToken(token, s) },
		runner.RemoveResetToken,
		s.queueControlFrame,
	)
	s.connIDGenerator = newConnIDGenerator(
		runner,
		srcConnID,
		&clientDestConnID,
		statelessResetter,
		connRunnerCallbacks{
			AddConnectionID:    func(connID protocol.ConnectionID) { runner.Add(connID, s) },
			RemoveConnectionID: runner.Remove,
			ReplaceWithClosed:  runner.ReplaceWithClosed,
		},
		s.queueControlFrame,
		connIDGenerator,
	)
	s.preSetup()
	s.rttStats.SetInitialRTT(rtt)
	s.sentPacketHandler, s.receivedPacketHandler = ackhandler.NewAckHandler(
		0,
		protocol.ByteCount(s.config.InitialPacketSize),
		s.rttStats,
		&s.connStats,
		clientAddressValidated,
		s.conn.capabilities().ECN,
		s.perspective,
		s.qlogger,
		s.logger,
	)
	s.currentMTUEstimate.Store(uint32(estimateMaxPayloadSize(protocol.ByteCount(s.config.InitialPacketSize))))
	statelessResetToken := statelessResetter.GetStatelessResetToken(srcConnID)
	params := &wire.TransportParameters{
		InitialMaxStreamDataBidiLocal:   protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataBidiRemote:  protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataUni:         protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxData:                  protocol.ByteCount(s.config.InitialConnectionReceiveWindow),
		MaxIdleTimeout:                  s.config.MaxIdleTimeout,
		MaxBidiStreamNum:                protocol.StreamNum(s.config.MaxIncomingStreams),
		MaxUniStreamNum:                 protocol.StreamNum(s.config.MaxIncomingUniStreams),
		MaxAckDelay:                     protocol.MaxAckDelayInclGranularity,
		AckDelayExponent:                protocol.AckDelayExponent,
		MaxUDPPayloadSize:               protocol.MaxPacketBufferSize,
		StatelessResetToken:             &statelessResetToken,
		OriginalDestinationConnectionID: origDestConnID,
		// For interoperability with quic-go versions before May 2023, this value must be set to a value
		// different from protocol.DefaultActiveConnectionIDLimit.
		// If set to the default value, it will be omitted from the transport parameters, which will make
		// old quic-go versions interpret it as 0, instead of the default value of 2.
		// See https://github.com/quic-go/quic-go/pull/3806.
		ActiveConnectionIDLimit:   protocol.MaxActiveConnectionIDs,
		InitialSourceConnectionID: srcConnID,
		RetrySourceConnectionID:   retrySrcConnID,
		EnableResetStreamAt:       conf.EnableStreamResetPartialDelivery,
	}
	if s.config.EnableDatagrams {
		params.MaxDatagramFrameSize = wire.MaxDatagramSize
	} else {
		params.MaxDatagramFrameSize = protocol.InvalidByteCount
	}
	if s.qlogger != nil {
		s.qlogTransportParameters(params, protocol.PerspectiveServer, false)
	}
	cs := handshake.NewCryptoSetupServer(
		clientDestConnID,
		conn.LocalAddr(),
		conn.RemoteAddr(),
		params,
		tlsConf,
		conf.Allow0RTT,
		s.rttStats,
		s.qlogger,
		logger,
		s.version,
	)
	s.cryptoStreamHandler = cs
	s.packer = newPacketPacker(srcConnID, s.connIDManager.Get, s.initialStream, s.handshakeStream, s.sentPacketHandler, s.retransmissionQueue, cs, s.framer, s.receivedPacketHandler, s.datagramQueue, s.perspective)
	s.unpacker = newPacketUnpacker(cs, s.srcConnIDLen)
	s.cryptoStreamManager = newCryptoStreamManager(s.initialStream, s.handshakeStream, s.oneRTTStream)
	return &wrappedConn{Conn: s}
}

// declare this as a variable, such that we can it mock it in the tests
var newClientConnection = func(
	ctx context.Context,
	conn sendConn,
	runner connRunner,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	connIDGenerator ConnectionIDGenerator,
	statelessResetter *statelessResetter,
	conf *Config,
	tlsConf *tls.Config,
	initialPacketNumber protocol.PacketNumber,
	enable0RTT bool,
	hasNegotiatedVersion bool,
	qlogTrace qlogwriter.Trace,
	logger utils.Logger,
	v protocol.Version,
) *wrappedConn {
	s := &Conn{
		conn:                conn,
		config:              conf,
		origDestConnID:      destConnID,
		handshakeDestConnID: destConnID,
		srcConnIDLen:        srcConnID.Len(),
		perspective:         protocol.PerspectiveClient,
		logID:               destConnID.String(),
		logger:              logger,
		qlogTrace:           qlogTrace,
		versionNegotiated:   hasNegotiatedVersion,
		version:             v,
	}
	if qlogTrace != nil {
		s.qlogger = qlogTrace.AddProducer()
	}
	if s.qlogger != nil {
		var srcAddr, destAddr *net.UDPAddr
		if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			srcAddr = addr
		}
		if addr, ok := conn.RemoteAddr().(*net.UDPAddr); ok {
			destAddr = addr
		}
		s.qlogger.RecordEvent(startedConnectionEvent(srcAddr, destAddr))
	}
	s.connIDManager = newConnIDManager(
		destConnID,
		func(token protocol.StatelessResetToken) { runner.AddResetToken(token, s) },
		runner.RemoveResetToken,
		s.queueControlFrame,
	)
	s.connIDGenerator = newConnIDGenerator(
		runner,
		srcConnID,
		nil,
		statelessResetter,
		connRunnerCallbacks{
			AddConnectionID:    func(connID protocol.ConnectionID) { runner.Add(connID, s) },
			RemoveConnectionID: runner.Remove,
			ReplaceWithClosed:  runner.ReplaceWithClosed,
		},
		s.queueControlFrame,
		connIDGenerator,
	)
	s.ctx, s.ctxCancel = context.WithCancelCause(ctx)
	s.preSetup()
	s.sentPacketHandler, s.receivedPacketHandler = ackhandler.NewAckHandler(
		initialPacketNumber,
		protocol.ByteCount(s.config.InitialPacketSize),
		s.rttStats,
		&s.connStats,
		false, // has no effect
		s.conn.capabilities().ECN,
		s.perspective,
		s.qlogger,
		s.logger,
	)
	s.currentMTUEstimate.Store(uint32(estimateMaxPayloadSize(protocol.ByteCount(s.config.InitialPacketSize))))
	oneRTTStream := newCryptoStream()
	params := &wire.TransportParameters{
		InitialMaxStreamDataBidiRemote: protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataBidiLocal:  protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataUni:        protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxData:                 protocol.ByteCount(s.config.InitialConnectionReceiveWindow),
		MaxIdleTimeout:                 s.config.MaxIdleTimeout,
		MaxBidiStreamNum:               protocol.StreamNum(s.config.MaxIncomingStreams),
		MaxUniStreamNum:                protocol.StreamNum(s.config.MaxIncomingUniStreams),
		MaxAckDelay:                    protocol.MaxAckDelayInclGranularity,
		MaxUDPPayloadSize:              protocol.MaxPacketBufferSize,
		AckDelayExponent:               protocol.AckDelayExponent,
		// For interoperability with quic-go versions before May 2023, this value must be set to a value
		// different from protocol.DefaultActiveConnectionIDLimit.
		// If set to the default value, it will be omitted from the transport parameters, which will make
		// old quic-go versions interpret it as 0, instead of the default value of 2.
		// See https://github.com/quic-go/quic-go/pull/3806.
		ActiveConnectionIDLimit:   protocol.MaxActiveConnectionIDs,
		InitialSourceConnectionID: srcConnID,
		EnableResetStreamAt:       conf.EnableStreamResetPartialDelivery,
	}
	if s.config.EnableDatagrams {
		params.MaxDatagramFrameSize = wire.MaxDatagramSize
	} else {
		params.MaxDatagramFrameSize = protocol.InvalidByteCount
	}
	if s.qlogger != nil {
		s.qlogTransportParameters(params, protocol.PerspectiveClient, false)
	}
	cs := handshake.NewCryptoSetupClient(
		destConnID,
		params,
		tlsConf,
		enable0RTT,
		s.rttStats,
		s.qlogger,
		logger,
		s.version,
	)
	s.cryptoStreamHandler = cs
	s.cryptoStreamManager = newCryptoStreamManager(s.initialStream, s.handshakeStream, oneRTTStream)
	s.unpacker = newPacketUnpacker(cs, s.srcConnIDLen)
	s.packer = newPacketPacker(srcConnID, s.connIDManager.Get, s.initialStream, s.handshakeStream, s.sentPacketHandler, s.retransmissionQueue, cs, s.framer, s.receivedPacketHandler, s.datagramQueue, s.perspective)
	if len(tlsConf.ServerName) > 0 {
		s.tokenStoreKey = tlsConf.ServerName
	} else {
		s.tokenStoreKey = conn.RemoteAddr().String()
	}
	if s.config.TokenStore != nil {
		if token := s.config.TokenStore.Pop(s.tokenStoreKey); token != nil {
			s.packer.SetToken(token.data)
			s.rttStats.SetInitialRTT(token.rtt)
		}
	}
	return &wrappedConn{Conn: s}
}

func (c *Conn) preSetup() {
	c.largestRcvdAppData = protocol.InvalidPacketNumber
	c.initialStream = newInitialCryptoStream(c.perspective == protocol.PerspectiveClient)
	c.handshakeStream = newCryptoStream()
	c.sendQueue = newSendQueue(c.conn)
	c.retransmissionQueue = newRetransmissionQueue()
	c.frameParser = *wire.NewFrameParser(
		c.config.EnableDatagrams,
		c.config.EnableStreamResetPartialDelivery,
		false, // ACK_FREQUENCY is not supported yet
	)
	c.rttStats = utils.NewRTTStats()
	c.connFlowController = flowcontrol.NewConnectionFlowController(
		protocol.ByteCount(c.config.InitialConnectionReceiveWindow),
		protocol.ByteCount(c.config.MaxConnectionReceiveWindow),
		func(size protocol.ByteCount) bool {
			if c.config.AllowConnectionWindowIncrease == nil {
				return true
			}
			return c.config.AllowConnectionWindowIncrease(c, uint64(size))
		},
		c.rttStats,
		c.logger,
	)
	c.earlyConnReadyChan = make(chan struct{})
	c.streamsMap = newStreamsMap(
		c.ctx,
		c,
		c.queueControlFrame,
		c.newFlowController,
		uint64(c.config.MaxIncomingStreams),
		uint64(c.config.MaxIncomingUniStreams),
		c.perspective,
	)
	c.framer = newFramer(c.connFlowController)
	c.receivedPackets.Init(8)
	c.notifyReceivedPacket = make(chan struct{}, 1)
	c.closeChan = make(chan struct{}, 1)
	c.sendingScheduled = make(chan struct{}, 1)
	c.handshakeCompleteChan = make(chan struct{})

	now := monotime.Now()
	c.lastPacketReceivedTime = now
	c.creationTime = now

	c.datagramQueue = newDatagramQueue(c.scheduleSending, c.logger)
	c.connState.Version = c.version
}

// run the connection main loop
func (c *Conn) run() (err error) {
	defer func() { c.ctxCancel(err) }()

	defer func() {
		// drain queued packets that will never be processed
		c.receivedPacketMx.Lock()
		defer c.receivedPacketMx.Unlock()

		for !c.receivedPackets.Empty() {
			p := c.receivedPackets.PopFront()
			p.buffer.Decrement()
			p.buffer.MaybeRelease()
		}
	}()

	c.timer = time.NewTimer(monotime.Until(c.idleTimeoutStartTime().Add(c.config.HandshakeIdleTimeout)))

	if err := c.cryptoStreamHandler.StartHandshake(c.ctx); err != nil {
		return err
	}
	if err := c.handleHandshakeEvents(monotime.Now()); err != nil {
		return err
	}
	go func() {
		if err := c.sendQueue.Run(); err != nil {
			c.destroyImpl(err)
		}
	}()

	if c.perspective == protocol.PerspectiveClient {
		c.scheduleSending() // so the ClientHello actually gets sent
	}

	var sendQueueAvailable <-chan struct{}

runLoop:
	for {
		if c.framer.QueuedTooManyControlFrames() {
			c.setCloseError(&closeError{err: &qerr.TransportError{ErrorCode: InternalError}})
			break runLoop
		}
		// Close immediately if requested
		select {
		case <-c.closeChan:
			break runLoop
		default:
		}

		// no need to set a timer if we can send packets immediately
		if c.pacingDeadline != deadlineSendImmediately {
			c.maybeResetTimer()
		}

		// 1st: handle undecryptable packets, if any.
		// This can only occur before completion of the handshake.
		if len(c.undecryptablePacketsToProcess) > 0 {
			var processedUndecryptablePacket bool
			queue := c.undecryptablePacketsToProcess
			c.undecryptablePacketsToProcess = nil
			for _, p := range queue {
				processed, err := c.handleOnePacket(p)
				if err != nil {
					c.setCloseError(&closeError{err: err})
					break runLoop
				}
				if processed {
					processedUndecryptablePacket = true
				}
			}
			if processedUndecryptablePacket {
				// if we processed any undecryptable packets, jump to the resetting of the timers directly
				continue
			}
		}

		// 2nd: receive packets.
		processed, err := c.handlePackets() // don't check receivedPackets.Len() in the run loop to avoid locking the mutex
		if err != nil {
			c.setCloseError(&closeError{err: err})
			break runLoop
		}

		// We don't need to wait for new events if:
		// * we processed packets: we probably need to send an ACK, and potentially more data
		// * the pacer allows us to send more packets immediately
		shouldProceedImmediately := sendQueueAvailable == nil && (processed || c.pacingDeadline.Equal(deadlineSendImmediately))
		if !shouldProceedImmediately {
			// 3rd: wait for something to happen:
			// * closing of the connection
			// * timer firing
			// * sending scheduled
			// * send queue available
			// * received packets
			select {
			case <-c.closeChan:
				break runLoop
			case <-c.timer.C:
			case <-c.sendingScheduled:
			case <-sendQueueAvailable:
			case <-c.notifyReceivedPacket:
				wasProcessed, err := c.handlePackets()
				if err != nil {
					c.setCloseError(&closeError{err: err})
					break runLoop
				}
				// if we processed any undecryptable packets, jump to the resetting of the timers directly
				if !wasProcessed {
					continue
				}
			}
		}

		// Check for loss detection timeout.
		// This could cause packets to be declared lost, and retransmissions to be enqueued.
		now := monotime.Now()
		if timeout := c.sentPacketHandler.GetLossDetectionTimeout(); !timeout.IsZero() && !timeout.After(now) {
			if err := c.sentPacketHandler.OnLossDetectionTimeout(now); err != nil {
				c.setCloseError(&closeError{err: err})
				break runLoop
			}
		}

		if keepAliveTime := c.nextKeepAliveTime(); !keepAliveTime.IsZero() && !now.Before(keepAliveTime) {
			// send a PING frame since there is no activity in the connection
			c.logger.Debugf("Sending a keep-alive PING to keep the connection alive.")
			c.framer.QueueControlFrame(&wire.PingFrame{})
			c.keepAlivePingSent = true
		} else if !c.handshakeComplete && now.Sub(c.creationTime) >= c.config.handshakeTimeout() {
			c.destroyImpl(qerr.ErrHandshakeTimeout)
			break runLoop
		} else {
			idleTimeoutStartTime := c.idleTimeoutStartTime()
			if (!c.handshakeComplete && now.Sub(idleTimeoutStartTime) >= c.config.HandshakeIdleTimeout) ||
				(c.handshakeComplete && !now.Before(c.nextIdleTimeoutTime())) {
				c.destroyImpl(qerr.ErrIdleTimeout)
				break runLoop
			}
		}

		c.connIDGenerator.RemoveRetiredConnIDs(now)

		if c.perspective == protocol.PerspectiveClient {
			pm := c.pathManagerOutgoing.Load()
			if pm != nil {
				tr, ok := pm.ShouldSwitchPath()
				if ok {
					c.switchToNewPath(tr, now)
				}
			}
		}

		if c.sendQueue.WouldBlock() {
			// The send queue is still busy sending out packets. Wait until there's space to enqueue new packets.
			sendQueueAvailable = c.sendQueue.Available()
			// Cancel the pacing timer, as we can't send any more packets until the send queue is available again.
			c.pacingDeadline = 0
			c.blocked = blockModeHardBlocked
			continue
		}

		if c.closeErr.Load() != nil {
			break runLoop
		}

		c.blocked = blockModeNone // sending might set it back to true if we're congestion limited
		if err := c.triggerSending(now); err != nil {
			c.setCloseError(&closeError{err: err})
			break runLoop
		}
		if c.sendQueue.WouldBlock() {
			// The send queue is still busy sending out packets. Wait until there's space to enqueue new packets.
			sendQueueAvailable = c.sendQueue.Available()
			// Cancel the pacing timer, as we can't send any more packets until the send queue is available again.
			c.pacingDeadline = 0
			c.blocked = blockModeHardBlocked
		} else {
			sendQueueAvailable = nil
		}
	}

	closeErr := c.closeErr.Load()
	c.cryptoStreamHandler.Close()
	c.sendQueue.Close() // close the send queue before sending the CONNECTION_CLOSE
	c.handleCloseError(closeErr)
	if c.qlogger != nil {
		if e := (&errCloseForRecreating{}); !errors.As(closeErr.err, &e) {
			c.qlogger.Close()
		}
	}
	c.logger.Infof("Connection %s closed.", c.logID)
	c.timer.Stop()
	return closeErr.err
}

// blocks until the early connection can be used
func (c *Conn) earlyConnReady() <-chan struct{} {
	return c.earlyConnReadyChan
}

// Context returns a context that is cancelled when the connection is closed.
// The cancellation cause is set to the error that caused the connection to close.
func (c *Conn) Context() context.Context {
	return c.ctx
}

func (c *Conn) supportsDatagrams() bool {
	return c.peerParams.MaxDatagramFrameSize > 0
}

// ConnectionState returns basic details about the QUIC connection.
func (c *Conn) ConnectionState() ConnectionState {
	c.connStateMutex.Lock()
	defer c.connStateMutex.Unlock()
	cs := c.cryptoStreamHandler.ConnectionState()
	c.connState.TLS = cs.ConnectionState
	c.connState.Used0RTT = cs.Used0RTT
	c.connState.SupportsStreamResetPartialDelivery = c.peerParams.EnableResetStreamAt
	c.connState.GSO = c.conn.capabilities().GSO
	return c.connState
}

// ConnectionStats contains statistics about the QUIC connection
type ConnectionStats struct {
	// MinRTT is the estimate of the minimum RTT observed on the active network
	// path.
	MinRTT time.Duration
	// LatestRTT is the last RTT sample observed on the active network path.
	LatestRTT time.Duration
	// SmoothedRTT is an exponentially weighted moving average of an endpoint's
	// RTT samples. See https://www.rfc-editor.org/rfc/rfc9002#section-5.3
	SmoothedRTT time.Duration
	// MeanDeviation estimates the variation in the RTT samples using a mean
	// variation. See https://www.rfc-editor.org/rfc/rfc9002#section-5.3
	MeanDeviation time.Duration

	// BytesSent is the number of bytes sent on the underlying connection,
	// including retransmissions. Does not include UDP or any other outer
	// framing.
	BytesSent uint64
	// PacketsSent is the number of packets sent on the underlying connection,
	// including those that are determined to have been lost.
	PacketsSent uint64
	// BytesReceived is the number of total bytes received on the underlying
	// connection, including duplicate data for streams. Does not include UDP or
	// any other outer framing.
	BytesReceived uint64
	// PacketsReceived is the number of total packets received on the underlying
	// connection, including packets that were not processable.
	PacketsReceived uint64
	// BytesLost is the number of bytes lost on the underlying connection (does
	// not monotonically increase, because packets that are declared lost can
	// subsequently be received). Does not include UDP or any other outer
	// framing.
	BytesLost uint64
	// PacketsLost is the number of packets lost on the underlying connection
	// (does not monotonically increase, because packets that are declared lost
	// can subsequently be received).
	PacketsLost uint64
}

func (c *Conn) ConnectionStats() ConnectionStats {
	return ConnectionStats{
		MinRTT:        c.rttStats.MinRTT(),
		LatestRTT:     c.rttStats.LatestRTT(),
		SmoothedRTT:   c.rttStats.SmoothedRTT(),
		MeanDeviation: c.rttStats.MeanDeviation(),

		BytesSent:       c.connStats.BytesSent.Load(),
		PacketsSent:     c.connStats.PacketsSent.Load(),
		BytesReceived:   c.connStats.BytesReceived.Load(),
		PacketsReceived: c.connStats.PacketsReceived.Load(),
		BytesLost:       c.connStats.BytesLost.Load(),
		PacketsLost:     c.connStats.PacketsLost.Load(),
	}
}

// Time when the connection should time out
func (c *Conn) nextIdleTimeoutTime() monotime.Time {
	idleTimeout := max(c.idleTimeout, c.rttStats.PTO(true)*3)
	return c.idleTimeoutStartTime().Add(idleTimeout)
}

// Time when the next keep-alive packet should be sent.
// It returns a zero time if no keep-alive should be sent.
func (c *Conn) nextKeepAliveTime() monotime.Time {
	if c.config.KeepAlivePeriod == 0 || c.keepAlivePingSent {
		return 0
	}
	keepAliveInterval := max(c.keepAliveInterval, c.rttStats.PTO(true)*3/2)
	return c.lastPacketReceivedTime.Add(keepAliveInterval)
}

func (c *Conn) maybeResetTimer() {
	var deadline monotime.Time
	if !c.handshakeComplete {
		deadline = c.creationTime.Add(c.config.handshakeTimeout())
		if t := c.idleTimeoutStartTime().Add(c.config.HandshakeIdleTimeout); t.Before(deadline) {
			deadline = t
		}
	} else {
		// A keep-alive packet is ack-eliciting, so it can only be sent if the connection is
		// neither congestion limited nor hard-blocked.
		if c.blocked != blockModeNone {
			deadline = c.nextIdleTimeoutTime()
		} else {
			if keepAliveTime := c.nextKeepAliveTime(); !keepAliveTime.IsZero() {
				deadline = keepAliveTime
			} else {
				deadline = c.nextIdleTimeoutTime()
			}
		}
	}
	// If the connection is hard-blocked, we can't even send acknowledgments,
	// nor can we send PTO probe packets.
	if c.blocked == blockModeHardBlocked {
		c.timer.Reset(monotime.Until(deadline))
		return
	}

	if t := c.receivedPacketHandler.GetAlarmTimeout(); !t.IsZero() && t.Before(deadline) {
		deadline = t
	}
	if t := c.sentPacketHandler.GetLossDetectionTimeout(); !t.IsZero() && t.Before(deadline) {
		deadline = t
	}
	if c.blocked == blockModeCongestionLimited {
		c.timer.Reset(monotime.Until(deadline))
		return
	}

	if t := c.connIDGenerator.NextRetireTime(); !t.IsZero() && t.Before(deadline) {
		deadline = t
	}
	if !c.pacingDeadline.IsZero() && c.pacingDeadline.Before(deadline) {
		deadline = c.pacingDeadline
	}
	c.timer.Reset(monotime.Until(deadline))
}

func (c *Conn) idleTimeoutStartTime() monotime.Time {
	startTime := c.lastPacketReceivedTime
	if t := c.firstAckElicitingPacketAfterIdleSentTime; !t.IsZero() && t.After(startTime) {
		startTime = t
	}
	return startTime
}

func (c *Conn) switchToNewPath(tr *Transport, now monotime.Time) {
	initialPacketSize := protocol.ByteCount(c.config.InitialPacketSize)
	c.sentPacketHandler.MigratedPath(now, initialPacketSize)
	maxPacketSize := protocol.ByteCount(protocol.MaxPacketBufferSize)
	if c.peerParams.MaxUDPPayloadSize > 0 && c.peerParams.MaxUDPPayloadSize < maxPacketSize {
		maxPacketSize = c.peerParams.MaxUDPPayloadSize
	}
	c.mtuDiscoverer.Reset(now, initialPacketSize, maxPacketSize)
	c.conn = newSendConn(tr.conn, c.conn.RemoteAddr(), packetInfo{}, utils.DefaultLogger) // TODO: find a better way
	c.sendQueue.Close()
	c.sendQueue = newSendQueue(c.conn)
	go func() {
		if err := c.sendQueue.Run(); err != nil {
			c.destroyImpl(err)
		}
	}()
}

func (c *Conn) handleHandshakeComplete(now monotime.Time) error {
	defer close(c.handshakeCompleteChan)
	// Once the handshake completes, we have derived 1-RTT keys.
	// There's no point in queueing undecryptable packets for later decryption anymore.
	c.undecryptablePackets = nil

	c.connIDManager.SetHandshakeComplete()
	c.connIDGenerator.SetHandshakeComplete(now.Add(3 * c.rttStats.PTO(false)))

	if c.qlogger != nil {
		c.qlogger.RecordEvent(qlog.ALPNInformation{
			ChosenALPN: c.cryptoStreamHandler.ConnectionState().NegotiatedProtocol,
		})
	}

	// The server applies transport parameters right away, but the client side has to wait for handshake completion.
	// During a 0-RTT connection, the client is only allowed to use the new transport parameters for 1-RTT packets.
	if c.perspective == protocol.PerspectiveClient {
		c.applyTransportParameters()
		return nil
	}

	// All these only apply to the server side.
	if err := c.handleHandshakeConfirmed(now); err != nil {
		return err
	}

	ticket, err := c.cryptoStreamHandler.GetSessionTicket()
	if err != nil {
		return err
	}
	if ticket != nil { // may be nil if session tickets are disabled via tls.Config.SessionTicketsDisabled
		c.oneRTTStream.Write(ticket)
		for c.oneRTTStream.HasData() {
			if cf := c.oneRTTStream.PopCryptoFrame(protocol.MaxPostHandshakeCryptoFrameSize); cf != nil {
				c.queueControlFrame(cf)
			}
		}
	}
	token, err := c.tokenGenerator.NewToken(c.conn.RemoteAddr(), c.rttStats.SmoothedRTT())
	if err != nil {
		return err
	}
	c.queueControlFrame(&wire.NewTokenFrame{Token: token})
	c.queueControlFrame(&wire.HandshakeDoneFrame{})
	return nil
}

func (c *Conn) handleHandshakeConfirmed(now monotime.Time) error {
	// Drop initial keys.
	// On the client side, this should have happened when sending the first Handshake packet,
	// but this is not guaranteed if the server misbehaves.
	// See CVE-2025-59530 for more details.
	if err := c.dropEncryptionLevel(protocol.EncryptionInitial, now); err != nil {
		return err
	}
	if err := c.dropEncryptionLevel(protocol.EncryptionHandshake, now); err != nil {
		return err
	}

	c.handshakeConfirmed = true
	c.cryptoStreamHandler.SetHandshakeConfirmed()

	if !c.config.DisablePathMTUDiscovery && c.conn.capabilities().DF {
		c.mtuDiscoverer.Start(now)
	}
	return nil
}

func (c *Conn) handlePackets() (wasProcessed bool, _ error) {
	// Now process all packets in the receivedPackets channel.
	// Limit the number of packets to the length of the receivedPackets channel,
	// so we eventually get a chance to send out an ACK when receiving a lot of packets.
	c.receivedPacketMx.Lock()
	numPackets := c.receivedPackets.Len()
	if numPackets == 0 {
		c.receivedPacketMx.Unlock()
		return false, nil
	}

	var hasMorePackets bool
	for i := 0; i < numPackets; i++ {
		if i > 0 {
			c.receivedPacketMx.Lock()
		}
		p := c.receivedPackets.PopFront()
		hasMorePackets = !c.receivedPackets.Empty()
		c.receivedPacketMx.Unlock()

		processed, err := c.handleOnePacket(p)
		if err != nil {
			return false, err
		}
		if processed {
			wasProcessed = true
		}
		if !hasMorePackets {
			break
		}
		// only process a single packet at a time before handshake completion
		if !c.handshakeComplete {
			break
		}
	}
	if hasMorePackets {
		select {
		case c.notifyReceivedPacket <- struct{}{}:
		default:
		}
	}
	return wasProcessed, nil
}

func (c *Conn) handleOnePacket(rp receivedPacket) (wasProcessed bool, _ error) {
	c.sentPacketHandler.ReceivedBytes(rp.Size(), rp.rcvTime)

	if wire.IsVersionNegotiationPacket(rp.data) {
		c.handleVersionNegotiationPacket(rp)
		return false, nil
	}

	var counter uint8
	var lastConnID protocol.ConnectionID
	data := rp.data
	p := rp
	for len(data) > 0 {
		if counter > 0 {
			p = *(p.Clone())
			p.data = data

			destConnID, err := wire.ParseConnectionID(p.data, c.srcConnIDLen)
			if err != nil {
				if c.qlogger != nil {
					c.qlogger.RecordEvent(qlog.PacketDropped{
						Raw:     qlog.RawInfo{Length: len(data)},
						Trigger: qlog.PacketDropHeaderParseError,
					})
				}
				c.logger.Debugf("error parsing packet, couldn't parse connection ID: %s", err)
				break
			}
			if destConnID != lastConnID {
				if c.qlogger != nil {
					c.qlogger.RecordEvent(qlog.PacketDropped{
						Header:  qlog.PacketHeader{DestConnectionID: destConnID},
						Raw:     qlog.RawInfo{Length: len(data)},
						Trigger: qlog.PacketDropUnknownConnectionID,
					})
				}
				c.logger.Debugf("coalesced packet has different destination connection ID: %s, expected %s", destConnID, lastConnID)
				break
			}
		}

		if wire.IsLongHeaderPacket(p.data[0]) {
			hdr, packetData, rest, err := wire.ParsePacket(p.data)
			if err != nil {
				if c.qlogger != nil {
					if err == wire.ErrUnsupportedVersion {
						c.qlogger.RecordEvent(qlog.PacketDropped{
							Header:  qlog.PacketHeader{Version: hdr.Version},
							Raw:     qlog.RawInfo{Length: len(data)},
							Trigger: qlog.PacketDropUnsupportedVersion,
						})
					} else {
						c.qlogger.RecordEvent(qlog.PacketDropped{
							Raw:     qlog.RawInfo{Length: len(data)},
							Trigger: qlog.PacketDropHeaderParseError,
						})
					}
				}
				c.logger.Debugf("error parsing packet: %s", err)
				break
			}
			lastConnID = hdr.DestConnectionID

			if hdr.Version != c.version {
				if c.qlogger != nil {
					c.qlogger.RecordEvent(qlog.PacketDropped{
						Raw:     qlog.RawInfo{Length: len(data)},
						Trigger: qlog.PacketDropUnexpectedVersion,
					})
				}
				c.logger.Debugf("Dropping packet with version %x. Expected %x.", hdr.Version, c.version)
				break
			}

			if counter > 0 {
				p.buffer.Split()
			}
			counter++

			// only log if this actually a coalesced packet
			if c.logger.Debug() && (counter > 1 || len(rest) > 0) {
				c.logger.Debugf("Parsed a coalesced packet. Part %d: %d bytes. Remaining: %d bytes.", counter, len(packetData), len(rest))
			}

			p.data = packetData

			processed, err := c.handleLongHeaderPacket(p, hdr)
			if err != nil {
				return false, err
			}
			if processed {
				wasProcessed = true
			}
			data = rest
		} else {
			if counter > 0 {
				p.buffer.Split()
			}
			processed, err := c.handleShortHeaderPacket(p, counter > 0)
			if err != nil {
				return false, err
			}
			if processed {
				wasProcessed = true
			}
			break
		}
	}

	p.buffer.MaybeRelease()
	c.blocked = blockModeNone
	return wasProcessed, nil
}

func (c *Conn) handleShortHeaderPacket(p receivedPacket, isCoalesced bool) (wasProcessed bool, _ error) {
	var wasQueued bool

	defer func() {
		// Put back the packet buffer if the packet wasn't queued for later decryption.
		if !wasQueued {
			p.buffer.Decrement()
		}
	}()

	destConnID, err := wire.ParseConnectionID(p.data, c.srcConnIDLen)
	if err != nil {
		c.qlogger.RecordEvent(qlog.PacketDropped{
			Header: qlog.PacketHeader{
				PacketType:   qlog.PacketType1RTT,
				PacketNumber: protocol.InvalidPacketNumber,
			},
			Raw:     qlog.RawInfo{Length: len(p.data)},
			Trigger: qlog.PacketDropHeaderParseError,
		})
		return false, nil
	}
	pn, pnLen, keyPhase, data, err := c.unpacker.UnpackShortHeader(p.rcvTime, p.data)
	if err != nil {
		// Stateless reset packets (see RFC 9000, section 10.3):
		// * fill the entire UDP datagram (i.e. they cannot be part of a coalesced packet)
		// * are short header packets (first bit is 0)
		// * have the QUIC bit set (second bit is 1)
		// * are at least 21 bytes long
		if !isCoalesced && len(p.data) >= protocol.MinReceivedStatelessResetSize && p.data[0]&0b11000000 == 0b01000000 {
			token := protocol.StatelessResetToken(p.data[len(p.data)-16:])
			if c.connIDManager.IsActiveStatelessResetToken(token) {
				return false, &StatelessResetError{}
			}
		}
		wasQueued, err = c.handleUnpackError(err, p, qlog.PacketType1RTT)
		return false, err
	}
	c.largestRcvdAppData = max(c.largestRcvdAppData, pn)

	if c.logger.Debug() {
		c.logger.Debugf("<- Reading packet %d (%d bytes) for connection %s, 1-RTT", pn, p.Size(), destConnID)
		wire.LogShortHeader(c.logger, destConnID, pn, pnLen, keyPhase)
	}

	if c.receivedPacketHandler.IsPotentiallyDuplicate(pn, protocol.Encryption1RTT) {
		c.logger.Debugf("Dropping (potentially) duplicate packet.")
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketType1RTT,
					PacketNumber: pn,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropDuplicate,
			})
		}
		return false, nil
	}

	var log func([]qlog.Frame)
	if c.qlogger != nil {
		log = func(frames []qlog.Frame) {
			c.qlogger.RecordEvent(qlog.PacketReceived{
				Header: qlog.PacketHeader{
					PacketType:       qlog.PacketType1RTT,
					DestConnectionID: destConnID,
					PacketNumber:     pn,
					KeyPhaseBit:      keyPhase,
				},
				Raw: qlog.RawInfo{
					Length:        int(p.Size()),
					PayloadLength: int(p.Size() - wire.ShortHeaderLen(destConnID, pnLen)),
				},
				Frames: frames,
				ECN:    toQlogECN(p.ecn),
			})
		}
	}
	isNonProbing, pathChallenge, err := c.handleUnpackedShortHeaderPacket(destConnID, pn, data, p.ecn, p.rcvTime, log)
	if err != nil {
		return false, err
	}

	// In RFC 9000, only the client can migrate between paths.
	if c.perspective == protocol.PerspectiveClient {
		return true, nil
	}
	if addrsEqual(p.remoteAddr, c.RemoteAddr()) {
		return true, nil
	}

	var shouldSwitchPath bool
	if c.pathManager == nil {
		c.pathManager = newPathManager(
			c.connIDManager.GetConnIDForPath,
			c.connIDManager.RetireConnIDForPath,
			c.logger,
		)
	}
	destConnID, frames, shouldSwitchPath := c.pathManager.HandlePacket(p.remoteAddr, p.rcvTime, pathChallenge, isNonProbing)
	if len(frames) > 0 {
		probe, buf, err := c.packer.PackPathProbePacket(destConnID, frames, c.version)
		if err != nil {
			return true, err
		}
		c.logger.Debugf("sending path probe packet to %s", p.remoteAddr)
		c.logShortHeaderPacket(probe.DestConnID, probe.Ack, probe.Frames, probe.StreamFrames, probe.PacketNumber, probe.PacketNumberLen, probe.KeyPhase, protocol.ECNNon, buf.Len(), false)
		c.registerPackedShortHeaderPacket(probe, protocol.ECNNon, p.rcvTime)
		c.sendQueue.SendProbe(buf, p.remoteAddr)
	}
	// We only switch paths in response to the highest-numbered non-probing packet,
	// see section 9.3 of RFC 9000.
	if !shouldSwitchPath || pn != c.largestRcvdAppData {
		return true, nil
	}
	c.pathManager.SwitchToPath(p.remoteAddr)
	c.sentPacketHandler.MigratedPath(p.rcvTime, protocol.ByteCount(c.config.InitialPacketSize))
	maxPacketSize := protocol.ByteCount(protocol.MaxPacketBufferSize)
	if c.peerParams.MaxUDPPayloadSize > 0 && c.peerParams.MaxUDPPayloadSize < maxPacketSize {
		maxPacketSize = c.peerParams.MaxUDPPayloadSize
	}
	c.mtuDiscoverer.Reset(
		p.rcvTime,
		protocol.ByteCount(c.config.InitialPacketSize),
		maxPacketSize,
	)
	c.conn.ChangeRemoteAddr(p.remoteAddr, p.info)
	return true, nil
}

func (c *Conn) handleLongHeaderPacket(p receivedPacket, hdr *wire.Header) (wasProcessed bool, _ error) {
	var wasQueued bool

	defer func() {
		// Put back the packet buffer if the packet wasn't queued for later decryption.
		if !wasQueued {
			p.buffer.Decrement()
		}
	}()

	if hdr.Type == protocol.PacketTypeRetry {
		return c.handleRetryPacket(hdr, p.data, p.rcvTime), nil
	}

	// The server can change the source connection ID with the first Handshake packet.
	// After this, all packets with a different source connection have to be ignored.
	if c.receivedFirstPacket && hdr.Type == protocol.PacketTypeInitial && hdr.SrcConnectionID != c.handshakeDestConnID {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketTypeInitial,
					PacketNumber: protocol.InvalidPacketNumber,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnknownConnectionID,
			})
		}
		c.logger.Debugf("Dropping Initial packet (%d bytes) with unexpected source connection ID: %s (expected %s)", p.Size(), hdr.SrcConnectionID, c.handshakeDestConnID)
		return false, nil
	}
	// drop 0-RTT packets, if we are a client
	if c.perspective == protocol.PerspectiveClient && hdr.Type == protocol.PacketType0RTT {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   qlog.PacketType0RTT,
					PacketNumber: protocol.InvalidPacketNumber,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return false, nil
	}

	packet, err := c.unpacker.UnpackLongHeader(hdr, p.data)
	if err != nil {
		wasQueued, err = c.handleUnpackError(err, p, toQlogPacketType(hdr.Type))
		return false, err
	}

	if c.logger.Debug() {
		c.logger.Debugf("<- Reading packet %d (%d bytes) for connection %s, %s", packet.hdr.PacketNumber, p.Size(), hdr.DestConnectionID, packet.encryptionLevel)
		packet.hdr.Log(c.logger)
	}

	if pn := packet.hdr.PacketNumber; c.receivedPacketHandler.IsPotentiallyDuplicate(pn, packet.encryptionLevel) {
		c.logger.Debugf("Dropping (potentially) duplicate packet.")
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       toQlogPacketType(packet.hdr.Type),
					DestConnectionID: hdr.DestConnectionID,
					SrcConnectionID:  hdr.SrcConnectionID,
					PacketNumber:     pn,
					Version:          packet.hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size()), PayloadLength: int(packet.hdr.Length)},
				Trigger: qlog.PacketDropDuplicate,
			})
		}
		return false, nil
	}

	if err := c.handleUnpackedLongHeaderPacket(packet, p.ecn, p.rcvTime, p.Size()); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Conn) handleUnpackError(err error, p receivedPacket, pt qlog.PacketType) (wasQueued bool, _ error) {
	switch err {
	case handshake.ErrKeysDropped:
		if c.qlogger != nil {
			connID, _ := wire.ParseConnectionID(p.data, c.srcConnIDLen)
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       pt,
					DestConnectionID: connID,
					PacketNumber:     protocol.InvalidPacketNumber,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropKeyUnavailable,
			})
		}
		c.logger.Debugf("Dropping %s packet (%d bytes) because we already dropped the keys.", pt, p.Size())
		return false, nil
	case handshake.ErrKeysNotYetAvailable:
		// Sealer for this encryption level not yet available.
		// Try again later.
		c.tryQueueingUndecryptablePacket(p, pt)
		return true, nil
	case wire.ErrInvalidReservedBits:
		return false, &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: err.Error(),
		}
	case handshake.ErrDecryptionFailed:
		// This might be a packet injected by an attacker. Drop it.
		if c.qlogger != nil {
			connID, _ := wire.ParseConnectionID(p.data, c.srcConnIDLen)
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       pt,
					DestConnectionID: connID,
					PacketNumber:     protocol.InvalidPacketNumber,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropPayloadDecryptError,
			})
		}
		c.logger.Debugf("Dropping %s packet (%d bytes) that could not be unpacked. Error: %s", pt, p.Size(), err)
		return false, nil
	default:
		var headerErr *headerParseError
		if errors.As(err, &headerErr) {
			// This might be a packet injected by an attacker. Drop it.
			if c.qlogger != nil {
				connID, _ := wire.ParseConnectionID(p.data, c.srcConnIDLen)
				c.qlogger.RecordEvent(qlog.PacketDropped{
					Header: qlog.PacketHeader{
						PacketType:       pt,
						DestConnectionID: connID,
						PacketNumber:     protocol.InvalidPacketNumber,
					},
					Raw:     qlog.RawInfo{Length: int(p.Size())},
					Trigger: qlog.PacketDropHeaderParseError,
				})
			}
			c.logger.Debugf("Dropping %s packet (%d bytes) for which we couldn't unpack the header. Error: %s", pt, p.Size(), err)
			return false, nil
		}
		// This is an error returned by the AEAD (other than ErrDecryptionFailed).
		// For example, a PROTOCOL_VIOLATION due to key updates.
		return false, err
	}
}

func (c *Conn) handleRetryPacket(hdr *wire.Header, data []byte, rcvTime monotime.Time) bool /* was this a valid Retry */ {
	if c.perspective == protocol.PerspectiveServer {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       qlog.PacketTypeRetry,
					SrcConnectionID:  hdr.SrcConnectionID,
					DestConnectionID: hdr.DestConnectionID,
					Version:          hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: len(data)},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		c.logger.Debugf("Ignoring Retry.")
		return false
	}
	if c.receivedFirstPacket {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       qlog.PacketTypeRetry,
					SrcConnectionID:  hdr.SrcConnectionID,
					DestConnectionID: hdr.DestConnectionID,
					Version:          hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: len(data)},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		c.logger.Debugf("Ignoring Retry, since we already received a packet.")
		return false
	}
	destConnID := c.connIDManager.Get()
	if hdr.SrcConnectionID == destConnID {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       qlog.PacketTypeRetry,
					SrcConnectionID:  hdr.SrcConnectionID,
					DestConnectionID: hdr.DestConnectionID,
					Version:          hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: len(data)},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		c.logger.Debugf("Ignoring Retry, since the server didn't change the Source Connection ID.")
		return false
	}
	// If a token is already set, this means that we already received a Retry from the server.
	// Ignore this Retry packet.
	if c.receivedRetry {
		c.logger.Debugf("Ignoring Retry, since a Retry was already received.")
		return false
	}

	tag := handshake.GetRetryIntegrityTag(data[:len(data)-16], destConnID, hdr.Version)
	if !bytes.Equal(data[len(data)-16:], tag[:]) {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:       qlog.PacketTypeRetry,
					SrcConnectionID:  hdr.SrcConnectionID,
					DestConnectionID: hdr.DestConnectionID,
					Version:          hdr.Version,
				},
				Raw:     qlog.RawInfo{Length: len(data)},
				Trigger: qlog.PacketDropPayloadDecryptError,
			})
		}
		c.logger.Debugf("Ignoring spoofed Retry. Integrity Tag doesn't match.")
		return false
	}

	newDestConnID := hdr.SrcConnectionID
	c.receivedRetry = true
	c.sentPacketHandler.ResetForRetry(rcvTime)
	c.handshakeDestConnID = newDestConnID
	c.retrySrcConnID = &newDestConnID
	c.cryptoStreamHandler.ChangeConnectionID(newDestConnID)
	c.packer.SetToken(hdr.Token)
	c.connIDManager.ChangeInitialConnID(newDestConnID)

	if c.logger.Debug() {
		c.logger.Debugf("<- Received Retry:")
		(&wire.ExtendedHeader{Header: *hdr}).Log(c.logger)
		c.logger.Debugf("Switching destination connection ID to: %s", hdr.SrcConnectionID)
	}
	if c.qlogger != nil {
		c.qlogger.RecordEvent(qlog.PacketReceived{
			Header: qlog.PacketHeader{
				PacketType:       qlog.PacketTypeRetry,
				DestConnectionID: destConnID,
				SrcConnectionID:  newDestConnID,
				Version:          hdr.Version,
				Token:            &qlog.Token{Raw: hdr.Token},
			},
			Raw: qlog.RawInfo{Length: len(data)},
		})
	}

	c.scheduleSending()
	return true
}

func (c *Conn) handleVersionNegotiationPacket(p receivedPacket) {
	if c.perspective == protocol.PerspectiveServer || // servers never receive version negotiation packets
		c.receivedFirstPacket || c.versionNegotiated { // ignore delayed / duplicated version negotiation packets
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header:  qlog.PacketHeader{PacketType: qlog.PacketTypeVersionNegotiation},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedPacket,
			})
		}
		return
	}

	src, dest, supportedVersions, err := wire.ParseVersionNegotiationPacket(p.data)
	if err != nil {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header:  qlog.PacketHeader{PacketType: qlog.PacketTypeVersionNegotiation},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropHeaderParseError,
			})
		}
		c.logger.Debugf("Error parsing Version Negotiation packet: %s", err)
		return
	}

	if slices.Contains(supportedVersions, c.version) {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header:  qlog.PacketHeader{PacketType: qlog.PacketTypeVersionNegotiation},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropUnexpectedVersion,
			})
		}
		// The Version Negotiation packet contains the version that we offered.
		// This might be a packet sent by an attacker, or it was corrupted.
		return
	}

	c.logger.Infof("Received a Version Negotiation packet. Supported Versions: %s", supportedVersions)
	if c.qlogger != nil {
		c.qlogger.RecordEvent(qlog.VersionNegotiationReceived{
			Header: qlog.PacketHeaderVersionNegotiation{
				DestConnectionID: dest,
				SrcConnectionID:  src,
			},
			SupportedVersions: supportedVersions,
		})
	}
	newVersion, ok := protocol.ChooseSupportedVersion(c.config.Versions, supportedVersions)
	if !ok {
		c.destroyImpl(&VersionNegotiationError{
			Ours:   c.config.Versions,
			Theirs: supportedVersions,
		})
		c.logger.Infof("No compatible QUIC version found.")
		return
	}
	if c.qlogger != nil {
		c.qlogger.RecordEvent(qlog.VersionInformation{
			ChosenVersion:  newVersion,
			ClientVersions: c.config.Versions,
			ServerVersions: supportedVersions,
		})
	}

	c.logger.Infof("Switching to QUIC version %s.", newVersion)
	nextPN, _ := c.sentPacketHandler.PeekPacketNumber(protocol.EncryptionInitial)
	c.destroyImpl(&errCloseForRecreating{
		nextPacketNumber: nextPN,
		nextVersion:      newVersion,
	})
}

func (c *Conn) handleUnpackedLongHeaderPacket(
	packet *unpackedPacket,
	ecn protocol.ECN,
	rcvTime monotime.Time,
	packetSize protocol.ByteCount, // only for logging
) error {
	if !c.receivedFirstPacket {
		c.receivedFirstPacket = true
		if !c.versionNegotiated && c.qlogger != nil {
			var clientVersions, serverVersions []Version
			switch c.perspective {
			case protocol.PerspectiveClient:
				clientVersions = c.config.Versions
			case protocol.PerspectiveServer:
				serverVersions = c.config.Versions
			}
			c.qlogger.RecordEvent(qlog.VersionInformation{
				ChosenVersion:  c.version,
				ClientVersions: clientVersions,
				ServerVersions: serverVersions,
			})
		}
		// The server can change the source connection ID with the first Handshake packet.
		if c.perspective == protocol.PerspectiveClient && packet.hdr.SrcConnectionID != c.handshakeDestConnID {
			cid := packet.hdr.SrcConnectionID
			c.logger.Debugf("Received first packet. Switching destination connection ID to: %s", cid)
			c.handshakeDestConnID = cid
			c.connIDManager.ChangeInitialConnID(cid)
		}
		// We create the connection as soon as we receive the first packet from the client.
		// We do that before authenticating the packet.
		// That means that if the source connection ID was corrupted,
		// we might have created a connection with an incorrect source connection ID.
		// Once we authenticate the first packet, we need to update it.
		if c.perspective == protocol.PerspectiveServer {
			if packet.hdr.SrcConnectionID != c.handshakeDestConnID {
				c.handshakeDestConnID = packet.hdr.SrcConnectionID
				c.connIDManager.ChangeInitialConnID(packet.hdr.SrcConnectionID)
			}
			if c.qlogger != nil {
				var srcAddr, destAddr *net.UDPAddr
				if addr, ok := c.conn.LocalAddr().(*net.UDPAddr); ok {
					srcAddr = addr
				}
				if addr, ok := c.conn.RemoteAddr().(*net.UDPAddr); ok {
					destAddr = addr
				}
				c.qlogger.RecordEvent(startedConnectionEvent(srcAddr, destAddr))
			}
		}
	}

	if c.perspective == protocol.PerspectiveServer && packet.encryptionLevel == protocol.EncryptionHandshake &&
		!c.droppedInitialKeys {
		// On the server side, Initial keys are dropped as soon as the first Handshake packet is received.
		// See Section 4.9.1 of RFC 9001.
		if err := c.dropEncryptionLevel(protocol.EncryptionInitial, rcvTime); err != nil {
			return err
		}
	}

	c.lastPacketReceivedTime = rcvTime
	c.firstAckElicitingPacketAfterIdleSentTime = 0
	c.keepAlivePingSent = false

	if packet.hdr.Type == protocol.PacketType0RTT {
		c.largestRcvdAppData = max(c.largestRcvdAppData, packet.hdr.PacketNumber)
	}

	var log func([]qlog.Frame)
	if c.qlogger != nil {
		log = func(frames []qlog.Frame) {
			var token *qlog.Token
			if len(packet.hdr.Token) > 0 {
				token = &qlog.Token{Raw: packet.hdr.Token}
			}
			c.qlogger.RecordEvent(qlog.PacketReceived{
				Header: qlog.PacketHeader{
					PacketType:       toQlogPacketType(packet.hdr.Type),
					DestConnectionID: packet.hdr.DestConnectionID,
					SrcConnectionID:  packet.hdr.SrcConnectionID,
					PacketNumber:     packet.hdr.PacketNumber,
					Version:          packet.hdr.Version,
					Token:            token,
				},
				Raw: qlog.RawInfo{
					Length:        int(packetSize),
					PayloadLength: int(packet.hdr.Length),
				},
				Frames: frames,
				ECN:    toQlogECN(ecn),
			})
		}
	}
	isAckEliciting, _, _, err := c.handleFrames(packet.data, packet.hdr.DestConnectionID, packet.encryptionLevel, log, rcvTime)
	if err != nil {
		return err
	}
	return c.receivedPacketHandler.ReceivedPacket(packet.hdr.PacketNumber, ecn, packet.encryptionLevel, rcvTime, isAckEliciting)
}

func (c *Conn) handleUnpackedShortHeaderPacket(
	destConnID protocol.ConnectionID,
	pn protocol.PacketNumber,
	data []byte,
	ecn protocol.ECN,
	rcvTime monotime.Time,
	log func([]qlog.Frame),
) (isNonProbing bool, pathChallenge *wire.PathChallengeFrame, _ error) {
	c.lastPacketReceivedTime = rcvTime
	c.firstAckElicitingPacketAfterIdleSentTime = 0
	c.keepAlivePingSent = false

	isAckEliciting, isNonProbing, pathChallenge, err := c.handleFrames(data, destConnID, protocol.Encryption1RTT, log, rcvTime)
	if err != nil {
		return false, nil, err
	}
	if err := c.receivedPacketHandler.ReceivedPacket(pn, ecn, protocol.Encryption1RTT, rcvTime, isAckEliciting); err != nil {
		return false, nil, err
	}
	return isNonProbing, pathChallenge, nil
}

// handleFrames parses the frames, one after the other, and handles them.
// It returns the last PATH_CHALLENGE frame contained in the packet, if any.
func (c *Conn) handleFrames(
	data []byte,
	destConnID protocol.ConnectionID,
	encLevel protocol.EncryptionLevel,
	log func([]qlog.Frame),
	rcvTime monotime.Time,
) (isAckEliciting, isNonProbing bool, pathChallenge *wire.PathChallengeFrame, _ error) {
	// Only used for tracing.
	// If we're not tracing, this slice will always remain empty.
	var frames []qlog.Frame
	if log != nil {
		frames = make([]qlog.Frame, 0, 4)
	}
	handshakeWasComplete := c.handshakeComplete
	var handleErr error
	var skipHandling bool

	for len(data) > 0 {
		frameType, l, err := c.frameParser.ParseType(data, encLevel)
		if err != nil {
			// The frame parser skips over PADDING frames, and returns an io.EOF if the PADDING
			// frames were the last frames in this packet.
			if err == io.EOF {
				break
			}
			return false, false, nil, err
		}
		data = data[l:]

		if ackhandler.IsFrameTypeAckEliciting(frameType) {
			isAckEliciting = true
		}
		if !wire.IsProbingFrameType(frameType) {
			isNonProbing = true
		}

		// We're inlining common cases, to avoid using interfaces
		// Fast path: STREAM, DATAGRAM and ACK
		if frameType.IsStreamFrameType() {
			streamFrame, l, err := c.frameParser.ParseStreamFrame(frameType, data, c.version)
			if err != nil {
				return false, false, nil, err
			}
			data = data[l:]

			if log != nil {
				frames = append(frames, toQlogFrame(streamFrame))
			}
			// an error occurred handling a previous frame, don't handle the current frame
			if skipHandling {
				continue
			}
			wire.LogFrame(c.logger, streamFrame, false)
			handleErr = c.streamsMap.HandleStreamFrame(streamFrame, rcvTime)
		} else if frameType.IsAckFrameType() {
			ackFrame, l, err := c.frameParser.ParseAckFrame(frameType, data, encLevel, c.version)
			if err != nil {
				return false, false, nil, err
			}
			data = data[l:]
			if log != nil {
				frames = append(frames, toQlogFrame(ackFrame))
			}
			// an error occurred handling a previous frame, don't handle the current frame
			if skipHandling {
				continue
			}
			wire.LogFrame(c.logger, ackFrame, false)
			handleErr = c.handleAckFrame(ackFrame, encLevel, rcvTime)
		} else if frameType.IsDatagramFrameType() {
			datagramFrame, l, err := c.frameParser.ParseDatagramFrame(frameType, data, c.version)
			if err != nil {
				return false, false, nil, err
			}
			data = data[l:]

			if log != nil {
				frames = append(frames, toQlogFrame(datagramFrame))
			}
			// an error occurred handling a previous frame, don't handle the current frame
			if skipHandling {
				continue
			}
			wire.LogFrame(c.logger, datagramFrame, false)
			handleErr = c.handleDatagramFrame(datagramFrame)
		} else {
			frame, l, err := c.frameParser.ParseLessCommonFrame(frameType, data, c.version)
			if err != nil {
				return false, false, nil, err
			}
			data = data[l:]

			if log != nil {
				frames = append(frames, toQlogFrame(frame))
			}
			// an error occurred handling a previous frame, don't handle the current frame
			if skipHandling {
				continue
			}
			pc, err := c.handleFrame(frame, encLevel, destConnID, rcvTime)
			if pc != nil {
				pathChallenge = pc
			}
			handleErr = err
		}

		if handleErr != nil {
			// if we're logging, we need to keep parsing (but not handling) all frames
			skipHandling = true
			if log == nil {
				return false, false, nil, handleErr
			}
		}
	}

	if log != nil {
		log(frames)
		if handleErr != nil {
			return false, false, nil, handleErr
		}
	}

	// Handle completion of the handshake after processing all the frames.
	// This ensures that we correctly handle the following case on the server side:
	// We receive a Handshake packet that contains the CRYPTO frame that allows us to complete the handshake,
	// and an ACK serialized after that CRYPTO frame. In this case, we still want to process the ACK frame.
	if !handshakeWasComplete && c.handshakeComplete {
		if err := c.handleHandshakeComplete(rcvTime); err != nil {
			return false, false, nil, err
		}
	}
	return
}

func (c *Conn) handleFrame(
	f wire.Frame,
	encLevel protocol.EncryptionLevel,
	destConnID protocol.ConnectionID,
	rcvTime monotime.Time,
) (pathChallenge *wire.PathChallengeFrame, _ error) {
	var err error
	wire.LogFrame(c.logger, f, false)
	switch frame := f.(type) {
	case *wire.CryptoFrame:
		err = c.handleCryptoFrame(frame, encLevel, rcvTime)
	case *wire.ConnectionCloseFrame:
		err = c.handleConnectionCloseFrame(frame)
	case *wire.ResetStreamFrame:
		err = c.streamsMap.HandleResetStreamFrame(frame, rcvTime)
	case *wire.MaxDataFrame:
		c.connFlowController.UpdateSendWindow(frame.MaximumData)
	case *wire.MaxStreamDataFrame:
		err = c.streamsMap.HandleMaxStreamDataFrame(frame)
	case *wire.MaxStreamsFrame:
		c.streamsMap.HandleMaxStreamsFrame(frame)
	case *wire.DataBlockedFrame:
	case *wire.StreamDataBlockedFrame:
		err = c.streamsMap.HandleStreamDataBlockedFrame(frame)
	case *wire.StreamsBlockedFrame:
	case *wire.StopSendingFrame:
		err = c.streamsMap.HandleStopSendingFrame(frame)
	case *wire.PingFrame:
	case *wire.PathChallengeFrame:
		c.handlePathChallengeFrame(frame)
		pathChallenge = frame
	case *wire.PathResponseFrame:
		err = c.handlePathResponseFrame(frame)
	case *wire.NewTokenFrame:
		err = c.handleNewTokenFrame(frame)
	case *wire.NewConnectionIDFrame:
		err = c.connIDManager.Add(frame)
	case *wire.RetireConnectionIDFrame:
		err = c.connIDGenerator.Retire(frame.SequenceNumber, destConnID, rcvTime.Add(3*c.rttStats.PTO(false)))
	case *wire.HandshakeDoneFrame:
		err = c.handleHandshakeDoneFrame(rcvTime)
	default:
		err = fmt.Errorf("unexpected frame type: %s", reflect.ValueOf(&frame).Elem().Type().Name())
	}
	return pathChallenge, err
}

// handlePacket is called by the server with a new packet
func (c *Conn) handlePacket(p receivedPacket) {
	c.receivedPacketMx.Lock()
	// Discard packets once the amount of queued packets is larger than
	// the channel size, protocol.MaxConnUnprocessedPackets
	if c.receivedPackets.Len() >= protocol.MaxConnUnprocessedPackets {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropDOSPrevention,
			})
		}
		c.receivedPacketMx.Unlock()
		return
	}
	c.receivedPackets.PushBack(p)
	c.receivedPacketMx.Unlock()

	select {
	case c.notifyReceivedPacket <- struct{}{}:
	default:
	}
}

func (c *Conn) handleConnectionCloseFrame(frame *wire.ConnectionCloseFrame) error {
	if frame.IsApplicationError {
		return &qerr.ApplicationError{
			Remote:       true,
			ErrorCode:    qerr.ApplicationErrorCode(frame.ErrorCode),
			ErrorMessage: frame.ReasonPhrase,
		}
	}
	return &qerr.TransportError{
		Remote:       true,
		ErrorCode:    qerr.TransportErrorCode(frame.ErrorCode),
		FrameType:    frame.FrameType,
		ErrorMessage: frame.ReasonPhrase,
	}
}

func (c *Conn) handleCryptoFrame(frame *wire.CryptoFrame, encLevel protocol.EncryptionLevel, rcvTime monotime.Time) error {
	if err := c.cryptoStreamManager.HandleCryptoFrame(frame, encLevel); err != nil {
		return err
	}
	for {
		data := c.cryptoStreamManager.GetCryptoData(encLevel)
		if data == nil {
			break
		}
		if err := c.cryptoStreamHandler.HandleMessage(data, encLevel); err != nil {
			return err
		}
	}
	return c.handleHandshakeEvents(rcvTime)
}

func (c *Conn) handleHandshakeEvents(now monotime.Time) error {
	for {
		ev := c.cryptoStreamHandler.NextEvent()
		var err error
		switch ev.Kind {
		case handshake.EventNoEvent:
			return nil
		case handshake.EventHandshakeComplete:
			// Don't call handleHandshakeComplete yet.
			// It's advantageous to process ACK frames that might be serialized after the CRYPTO frame first.
			c.handshakeComplete = true
		case handshake.EventReceivedTransportParameters:
			err = c.handleTransportParameters(ev.TransportParameters)
		case handshake.EventRestoredTransportParameters:
			c.restoreTransportParameters(ev.TransportParameters)
			close(c.earlyConnReadyChan)
		case handshake.EventReceivedReadKeys:
			// queue all previously undecryptable packets
			c.undecryptablePacketsToProcess = append(c.undecryptablePacketsToProcess, c.undecryptablePackets...)
			c.undecryptablePackets = nil
		case handshake.EventDiscard0RTTKeys:
			err = c.dropEncryptionLevel(protocol.Encryption0RTT, now)
		case handshake.EventWriteInitialData:
			_, err = c.initialStream.Write(ev.Data)
		case handshake.EventWriteHandshakeData:
			_, err = c.handshakeStream.Write(ev.Data)
		}
		if err != nil {
			return err
		}
	}
}

func (c *Conn) handlePathChallengeFrame(f *wire.PathChallengeFrame) {
	if c.perspective == protocol.PerspectiveClient {
		c.queueControlFrame(&wire.PathResponseFrame{Data: f.Data})
	}
}

func (c *Conn) handlePathResponseFrame(f *wire.PathResponseFrame) error {
	switch c.perspective {
	case protocol.PerspectiveClient:
		return c.handlePathResponseFrameClient(f)
	case protocol.PerspectiveServer:
		return c.handlePathResponseFrameServer(f)
	default:
		panic("unreachable")
	}
}

func (c *Conn) handlePathResponseFrameClient(f *wire.PathResponseFrame) error {
	pm := c.pathManagerOutgoing.Load()
	if pm == nil {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "unexpected PATH_RESPONSE frame",
		}
	}
	pm.HandlePathResponseFrame(f)
	return nil
}

func (c *Conn) handlePathResponseFrameServer(f *wire.PathResponseFrame) error {
	if c.pathManager == nil {
		// since we didn't send PATH_CHALLENGEs yet, we don't expect PATH_RESPONSEs
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "unexpected PATH_RESPONSE frame",
		}
	}
	c.pathManager.HandlePathResponseFrame(f)
	return nil
}

func (c *Conn) handleNewTokenFrame(frame *wire.NewTokenFrame) error {
	if c.perspective == protocol.PerspectiveServer {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "received NEW_TOKEN frame from the client",
		}
	}
	if c.config.TokenStore != nil {
		c.config.TokenStore.Put(c.tokenStoreKey, &ClientToken{data: frame.Token, rtt: c.rttStats.SmoothedRTT()})
	}
	return nil
}

func (c *Conn) handleHandshakeDoneFrame(rcvTime monotime.Time) error {
	if c.perspective == protocol.PerspectiveServer {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "received a HANDSHAKE_DONE frame",
		}
	}
	if !c.handshakeConfirmed {
		return c.handleHandshakeConfirmed(rcvTime)
	}
	return nil
}

func (c *Conn) handleAckFrame(frame *wire.AckFrame, encLevel protocol.EncryptionLevel, rcvTime monotime.Time) error {
	acked1RTTPacket, err := c.sentPacketHandler.ReceivedAck(frame, encLevel, c.lastPacketReceivedTime)
	if err != nil {
		return err
	}
	if !acked1RTTPacket {
		return nil
	}
	// On the client side: If the packet acknowledged a 1-RTT packet, this confirms the handshake.
	// This is only possible if the ACK was sent in a 1-RTT packet.
	// This is an optimization over simply waiting for a HANDSHAKE_DONE frame, see section 4.1.2 of RFC 9001.
	if c.perspective == protocol.PerspectiveClient && !c.handshakeConfirmed {
		if err := c.handleHandshakeConfirmed(rcvTime); err != nil {
			return err
		}
	}
	// If one of the acknowledged packets was a Path MTU probe packet, this might have increased the Path MTU estimate.
	if c.mtuDiscoverer != nil {
		if mtu := c.mtuDiscoverer.CurrentSize(); mtu > protocol.ByteCount(c.currentMTUEstimate.Load()) {
			c.currentMTUEstimate.Store(uint32(mtu))
			c.sentPacketHandler.SetMaxDatagramSize(mtu)
		}
	}
	return c.cryptoStreamHandler.SetLargest1RTTAcked(frame.LargestAcked())
}

func (c *Conn) handleDatagramFrame(f *wire.DatagramFrame) error {
	if f.Length(c.version) > wire.MaxDatagramSize {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "DATAGRAM frame too large",
		}
	}
	c.datagramQueue.HandleDatagramFrame(f)
	return nil
}

func (c *Conn) setCloseError(e *closeError) {
	c.closeErr.CompareAndSwap(nil, e)
	select {
	case c.closeChan <- struct{}{}:
	default:
	}
}

// closeLocal closes the connection and send a CONNECTION_CLOSE containing the error
func (c *Conn) closeLocal(e error) {
	c.setCloseError(&closeError{err: e, immediate: false})
}

// destroy closes the connection without sending the error on the wire
func (c *Conn) destroy(e error) {
	c.destroyImpl(e)
	<-c.ctx.Done()
}

func (c *Conn) destroyImpl(e error) {
	c.setCloseError(&closeError{err: e, immediate: true})
}

// CloseWithError closes the connection with an error.
// The error string will be sent to the peer.
func (c *Conn) CloseWithError(code ApplicationErrorCode, desc string) error {
	c.closeLocal(&qerr.ApplicationError{
		ErrorCode:    code,
		ErrorMessage: desc,
	})
	<-c.ctx.Done()
	return nil
}

func (c *Conn) closeWithTransportError(code TransportErrorCode) {
	c.closeLocal(&qerr.TransportError{ErrorCode: code})
	<-c.ctx.Done()
}

func (c *Conn) handleCloseError(closeErr *closeError) {
	if closeErr.immediate {
		if nerr, ok := closeErr.err.(net.Error); ok && nerr.Timeout() {
			c.logger.Errorf("Destroying connection: %s", closeErr.err)
		} else {
			c.logger.Errorf("Destroying connection with error: %s", closeErr.err)
		}
	} else {
		if closeErr.err == nil {
			c.logger.Infof("Closing connection.")
		} else {
			c.logger.Errorf("Closing connection with error: %s", closeErr.err)
		}
	}

	e := closeErr.err
	if e == nil {
		e = &qerr.ApplicationError{}
	} else {
		defer func() { closeErr.err = e }()
	}

	var (
		statelessResetErr     *StatelessResetError
		versionNegotiationErr *VersionNegotiationError
		recreateErr           *errCloseForRecreating
		applicationErr        *ApplicationError
		transportErr          *TransportError
	)
	var isRemoteClose bool
	var trigger qlog.ConnectionCloseTrigger
	var reason string
	var transportErrorCode *qlog.TransportErrorCode
	var applicationErrorCode *qlog.ApplicationErrorCode
	switch {
	case errors.Is(e, qerr.ErrIdleTimeout),
		errors.Is(e, qerr.ErrHandshakeTimeout):
		trigger = qlog.ConnectionCloseTriggerIdleTimeout
	case errors.As(e, &statelessResetErr):
		trigger = qlog.ConnectionCloseTriggerStatelessReset
	case errors.As(e, &versionNegotiationErr):
		trigger = qlog.ConnectionCloseTriggerVersionMismatch
	case errors.As(e, &recreateErr):
	case errors.As(e, &applicationErr):
		isRemoteClose = applicationErr.Remote
		reason = applicationErr.ErrorMessage
		applicationErrorCode = &applicationErr.ErrorCode
	case errors.As(e, &transportErr):
		isRemoteClose = transportErr.Remote
		reason = transportErr.ErrorMessage
		transportErrorCode = &transportErr.ErrorCode
	case closeErr.immediate:
		e = closeErr.err
	default:
		te := &qerr.TransportError{
			ErrorCode:    qerr.InternalError,
			ErrorMessage: e.Error(),
		}
		e = te
		reason = te.ErrorMessage
		code := te.ErrorCode
		transportErrorCode = &code
	}

	c.streamsMap.CloseWithError(e)
	if c.datagramQueue != nil {
		c.datagramQueue.CloseWithError(e)
	}

	// In rare instances, the connection ID manager might switch to a new connection ID
	// when sending the CONNECTION_CLOSE frame.
	// The connection ID manager removes the active stateless reset token from the packet
	// handler map when it is closed, so we need to make sure that this happens last.
	defer c.connIDManager.Close()

	if c.qlogger != nil && !errors.As(e, &recreateErr) {
		initiator := qlog.InitiatorLocal
		if isRemoteClose {
			initiator = qlog.InitiatorRemote
		}
		c.qlogger.RecordEvent(qlog.ConnectionClosed{
			Initiator:        initiator,
			ConnectionError:  transportErrorCode,
			ApplicationError: applicationErrorCode,
			Trigger:          trigger,
			Reason:           reason,
		})
	}

	// If this is a remote close we're done here
	if isRemoteClose {
		c.connIDGenerator.ReplaceWithClosed(nil, 3*c.rttStats.PTO(false))
		return
	}
	if closeErr.immediate {
		c.connIDGenerator.RemoveAll()
		return
	}
	// Don't send out any CONNECTION_CLOSE if this is an error that occurred
	// before we even sent out the first packet.
	if c.perspective == protocol.PerspectiveClient && !c.sentFirstPacket {
		c.connIDGenerator.RemoveAll()
		return
	}
	connClosePacket, err := c.sendConnectionClose(e)
	if err != nil {
		c.logger.Debugf("Error sending CONNECTION_CLOSE: %s", err)
	}
	c.connIDGenerator.ReplaceWithClosed(connClosePacket, 3*c.rttStats.PTO(false))
}

func (c *Conn) dropEncryptionLevel(encLevel protocol.EncryptionLevel, now monotime.Time) error {
	c.sentPacketHandler.DropPackets(encLevel, now)
	c.receivedPacketHandler.DropPackets(encLevel)
	//nolint:exhaustive // only Initial and 0-RTT need special treatment
	switch encLevel {
	case protocol.EncryptionInitial:
		c.droppedInitialKeys = true
		c.cryptoStreamHandler.DiscardInitialKeys()
	case protocol.Encryption0RTT:
		c.streamsMap.ResetFor0RTT()
		c.framer.Handle0RTTRejection()
		return c.connFlowController.Reset()
	}
	return c.cryptoStreamManager.Drop(encLevel)
}

// is called for the client, when restoring transport parameters saved for 0-RTT
func (c *Conn) restoreTransportParameters(params *wire.TransportParameters) {
	if c.logger.Debug() {
		c.logger.Debugf("Restoring Transport Parameters: %s", params)
	}
	if c.qlogger != nil {
		c.qlogger.RecordEvent(qlog.ParametersSet{
			Restore:                         true,
			Initiator:                       qlog.InitiatorRemote,
			SentBy:                          c.perspective,
			OriginalDestinationConnectionID: params.OriginalDestinationConnectionID,
			InitialSourceConnectionID:       params.InitialSourceConnectionID,
			RetrySourceConnectionID:         params.RetrySourceConnectionID,
			StatelessResetToken:             params.StatelessResetToken,
			DisableActiveMigration:          params.DisableActiveMigration,
			MaxIdleTimeout:                  params.MaxIdleTimeout,
			MaxUDPPayloadSize:               params.MaxUDPPayloadSize,
			AckDelayExponent:                params.AckDelayExponent,
			MaxAckDelay:                     params.MaxAckDelay,
			ActiveConnectionIDLimit:         params.ActiveConnectionIDLimit,
			InitialMaxData:                  params.InitialMaxData,
			InitialMaxStreamDataBidiLocal:   params.InitialMaxStreamDataBidiLocal,
			InitialMaxStreamDataBidiRemote:  params.InitialMaxStreamDataBidiRemote,
			InitialMaxStreamDataUni:         params.InitialMaxStreamDataUni,
			InitialMaxStreamsBidi:           int64(params.MaxBidiStreamNum),
			InitialMaxStreamsUni:            int64(params.MaxUniStreamNum),
			MaxDatagramFrameSize:            params.MaxDatagramFrameSize,
			EnableResetStreamAt:             params.EnableResetStreamAt,
		})
	}

	c.peerParams = params
	c.connIDGenerator.SetMaxActiveConnIDs(params.ActiveConnectionIDLimit)
	c.connFlowController.UpdateSendWindow(params.InitialMaxData)
	c.streamsMap.HandleTransportParameters(params)
	c.connStateMutex.Lock()
	c.connState.SupportsDatagrams = c.supportsDatagrams()
	c.connStateMutex.Unlock()
}

func (c *Conn) handleTransportParameters(params *wire.TransportParameters) error {
	if c.qlogger != nil {
		c.qlogTransportParameters(params, c.perspective.Opposite(), false)
	}
	if err := c.checkTransportParameters(params); err != nil {
		return &qerr.TransportError{
			ErrorCode:    qerr.TransportParameterError,
			ErrorMessage: err.Error(),
		}
	}

	if c.perspective == protocol.PerspectiveClient && c.peerParams != nil && c.ConnectionState().Used0RTT && !params.ValidForUpdate(c.peerParams) {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "server sent reduced limits after accepting 0-RTT data",
		}
	}

	c.peerParams = params
	// On the client side we have to wait for handshake completion.
	// During a 0-RTT connection, we are only allowed to use the new transport parameters for 1-RTT packets.
	if c.perspective == protocol.PerspectiveServer {
		c.applyTransportParameters()
		// On the server side, the early connection is ready as soon as we processed
		// the client's transport parameters.
		close(c.earlyConnReadyChan)
	}

	c.connStateMutex.Lock()
	c.connState.SupportsDatagrams = c.supportsDatagrams()
	c.connStateMutex.Unlock()
	return nil
}

func (c *Conn) checkTransportParameters(params *wire.TransportParameters) error {
	if c.logger.Debug() {
		c.logger.Debugf("Processed Transport Parameters: %s", params)
	}

	// check the initial_source_connection_id
	if params.InitialSourceConnectionID != c.handshakeDestConnID {
		return fmt.Errorf("expected initial_source_connection_id to equal %s, is %s", c.handshakeDestConnID, params.InitialSourceConnectionID)
	}

	if c.perspective == protocol.PerspectiveServer {
		return nil
	}
	// check the original_destination_connection_id
	if params.OriginalDestinationConnectionID != c.origDestConnID {
		return fmt.Errorf("expected original_destination_connection_id to equal %s, is %s", c.origDestConnID, params.OriginalDestinationConnectionID)
	}
	if c.retrySrcConnID != nil { // a Retry was performed
		if params.RetrySourceConnectionID == nil {
			return errors.New("missing retry_source_connection_id")
		}
		if *params.RetrySourceConnectionID != *c.retrySrcConnID {
			return fmt.Errorf("expected retry_source_connection_id to equal %s, is %s", c.retrySrcConnID, *params.RetrySourceConnectionID)
		}
	} else if params.RetrySourceConnectionID != nil {
		return errors.New("received retry_source_connection_id, although no Retry was performed")
	}
	return nil
}

func (c *Conn) applyTransportParameters() {
	params := c.peerParams
	// Our local idle timeout will always be > 0.
	c.idleTimeout = c.config.MaxIdleTimeout
	// If the peer advertised an idle timeout, take the minimum of the values.
	if params.MaxIdleTimeout > 0 {
		c.idleTimeout = min(c.idleTimeout, params.MaxIdleTimeout)
	}
	c.keepAliveInterval = min(c.config.KeepAlivePeriod, c.idleTimeout/2)
	c.streamsMap.HandleTransportParameters(params)
	c.frameParser.SetAckDelayExponent(params.AckDelayExponent)
	c.connFlowController.UpdateSendWindow(params.InitialMaxData)
	c.rttStats.SetMaxAckDelay(params.MaxAckDelay)
	c.connIDGenerator.SetMaxActiveConnIDs(params.ActiveConnectionIDLimit)
	if params.StatelessResetToken != nil {
		c.connIDManager.SetStatelessResetToken(*params.StatelessResetToken)
	}
	// We don't support connection migration yet, so we don't have any use for the preferred_address.
	if params.PreferredAddress != nil {
		// Retire the connection ID.
		c.connIDManager.AddFromPreferredAddress(params.PreferredAddress.ConnectionID, params.PreferredAddress.StatelessResetToken)
	}
	maxPacketSize := protocol.ByteCount(protocol.MaxPacketBufferSize)
	if params.MaxUDPPayloadSize > 0 && params.MaxUDPPayloadSize < maxPacketSize {
		maxPacketSize = params.MaxUDPPayloadSize
	}
	c.mtuDiscoverer = newMTUDiscoverer(
		c.rttStats,
		protocol.ByteCount(c.config.InitialPacketSize),
		maxPacketSize,
		c.qlogger,
	)
}

func (c *Conn) triggerSending(now monotime.Time) error {
	c.pacingDeadline = 0

	sendMode := c.sentPacketHandler.SendMode(now)
	switch sendMode {
	case ackhandler.SendAny:
		return c.sendPackets(now)
	case ackhandler.SendNone:
		c.blocked = blockModeHardBlocked
		return nil
	case ackhandler.SendPacingLimited:
		deadline := c.sentPacketHandler.TimeUntilSend()
		if deadline.IsZero() {
			deadline = deadlineSendImmediately
		}
		c.pacingDeadline = deadline
		// Allow sending of an ACK if we're pacing limit.
		// This makes sure that a peer that is mostly receiving data (and thus has an inaccurate cwnd estimate)
		// sends enough ACKs to allow its peer to utilize the bandwidth.
		return c.maybeSendAckOnlyPacket(now)
	case ackhandler.SendAck:
		// We can at most send a single ACK only packet.
		// There will only be a new ACK after receiving new packets.
		// SendAck is only returned when we're congestion limited, so we don't need to set the pacing timer.
		c.blocked = blockModeCongestionLimited
		return c.maybeSendAckOnlyPacket(now)
	case ackhandler.SendPTOInitial, ackhandler.SendPTOHandshake, ackhandler.SendPTOAppData:
		if err := c.sendProbePacket(sendMode, now); err != nil {
			return err
		}
		if c.sendQueue.WouldBlock() {
			c.scheduleSending()
			return nil
		}
		return c.triggerSending(now)
	default:
		return fmt.Errorf("BUG: invalid send mode %d", sendMode)
	}
}

func (c *Conn) sendPackets(now monotime.Time) error {
	if c.perspective == protocol.PerspectiveClient && c.handshakeConfirmed {
		if pm := c.pathManagerOutgoing.Load(); pm != nil {
			connID, frame, tr, ok := pm.NextPathToProbe()
			if ok {
				probe, buf, err := c.packer.PackPathProbePacket(connID, []ackhandler.Frame{frame}, c.version)
				if err != nil {
					return err
				}
				c.logger.Debugf("sending path probe packet from %s", c.LocalAddr())
				c.logShortHeaderPacket(probe.DestConnID, probe.Ack, probe.Frames, probe.StreamFrames, probe.PacketNumber, probe.PacketNumberLen, probe.KeyPhase, protocol.ECNNon, buf.Len(), false)
				c.registerPackedShortHeaderPacket(probe, protocol.ECNNon, now)
				tr.WriteTo(buf.Data, c.conn.RemoteAddr())
				// There's (likely) more data to send. Loop around again.
				c.scheduleSending()
				return nil
			}
		}
	}

	// Path MTU Discovery
	// Can't use GSO, since we need to send a single packet that's larger than our current maximum size.
	// Performance-wise, this doesn't matter, since we only send a very small (<10) number of
	// MTU probe packets per connection.
	if c.handshakeConfirmed && c.mtuDiscoverer != nil && c.mtuDiscoverer.ShouldSendProbe(now) {
		ping, size := c.mtuDiscoverer.GetPing(now)
		p, buf, err := c.packer.PackMTUProbePacket(ping, size, c.version)
		if err != nil {
			return err
		}
		ecn := c.sentPacketHandler.ECNMode(true)
		c.logShortHeaderPacket(p.DestConnID, p.Ack, p.Frames, p.StreamFrames, p.PacketNumber, p.PacketNumberLen, p.KeyPhase, ecn, buf.Len(), false)
		c.registerPackedShortHeaderPacket(p, ecn, now)
		c.sendQueue.Send(buf, 0, ecn)
		// There's (likely) more data to send. Loop around again.
		c.scheduleSending()
		return nil
	}

	if offset := c.connFlowController.GetWindowUpdate(now); offset > 0 {
		c.framer.QueueControlFrame(&wire.MaxDataFrame{MaximumData: offset})
	}
	if cf := c.cryptoStreamManager.GetPostHandshakeData(protocol.MaxPostHandshakeCryptoFrameSize); cf != nil {
		c.queueControlFrame(cf)
	}

	if !c.handshakeConfirmed {
		packet, err := c.packer.PackCoalescedPacket(false, c.maxPacketSize(), now, c.version)
		if err != nil || packet == nil {
			return err
		}
		c.sentFirstPacket = true
		if err := c.sendPackedCoalescedPacket(packet, c.sentPacketHandler.ECNMode(packet.IsOnlyShortHeaderPacket()), now); err != nil {
			return err
		}
		//nolint:exhaustive // only need to handle pacing-related events here
		switch c.sentPacketHandler.SendMode(now) {
		case ackhandler.SendPacingLimited:
			c.resetPacingDeadline()
		case ackhandler.SendAny:
			c.pacingDeadline = deadlineSendImmediately
		}
		return nil
	}

	if c.conn.capabilities().GSO {
		return c.sendPacketsWithGSO(now)
	}
	return c.sendPacketsWithoutGSO(now)
}

func (c *Conn) sendPacketsWithoutGSO(now monotime.Time) error {
	for {
		buf := getPacketBuffer()
		ecn := c.sentPacketHandler.ECNMode(true)
		if _, err := c.appendOneShortHeaderPacket(buf, c.maxPacketSize(), ecn, now); err != nil {
			if err == errNothingToPack {
				buf.Release()
				return nil
			}
			return err
		}

		c.sendQueue.Send(buf, 0, ecn)

		if c.sendQueue.WouldBlock() {
			return nil
		}
		sendMode := c.sentPacketHandler.SendMode(now)
		if sendMode == ackhandler.SendPacingLimited {
			c.resetPacingDeadline()
			return nil
		}
		if sendMode != ackhandler.SendAny {
			return nil
		}
		// Prioritize receiving of packets over sending out more packets.
		c.receivedPacketMx.Lock()
		hasPackets := !c.receivedPackets.Empty()
		c.receivedPacketMx.Unlock()
		if hasPackets {
			c.pacingDeadline = deadlineSendImmediately
			return nil
		}
	}
}

func (c *Conn) sendPacketsWithGSO(now monotime.Time) error {
	buf := getLargePacketBuffer()
	maxSize := c.maxPacketSize()

	ecn := c.sentPacketHandler.ECNMode(true)
	for {
		var dontSendMore bool
		size, err := c.appendOneShortHeaderPacket(buf, maxSize, ecn, now)
		if err != nil {
			if err != errNothingToPack {
				return err
			}
			if buf.Len() == 0 {
				buf.Release()
				return nil
			}
			dontSendMore = true
		}

		if !dontSendMore {
			sendMode := c.sentPacketHandler.SendMode(now)
			if sendMode == ackhandler.SendPacingLimited {
				c.resetPacingDeadline()
			}
			if sendMode != ackhandler.SendAny {
				dontSendMore = true
			}
		}

		// Don't send more packets in this batch if they require a different ECN marking than the previous ones.
		nextECN := c.sentPacketHandler.ECNMode(true)

		// Append another packet if
		// 1. The congestion controller and pacer allow sending more
		// 2. The last packet appended was a full-size packet
		// 3. The next packet will have the same ECN marking
		// 4. We still have enough space for another full-size packet in the buffer
		if !dontSendMore && size == maxSize && nextECN == ecn && buf.Len()+maxSize <= buf.Cap() {
			continue
		}

		c.sendQueue.Send(buf, uint16(maxSize), ecn)

		if dontSendMore {
			return nil
		}
		if c.sendQueue.WouldBlock() {
			return nil
		}

		// Prioritize receiving of packets over sending out more packets.
		c.receivedPacketMx.Lock()
		hasPackets := !c.receivedPackets.Empty()
		c.receivedPacketMx.Unlock()
		if hasPackets {
			c.pacingDeadline = deadlineSendImmediately
			return nil
		}

		ecn = nextECN
		buf = getLargePacketBuffer()
	}
}

func (c *Conn) resetPacingDeadline() {
	deadline := c.sentPacketHandler.TimeUntilSend()
	if deadline.IsZero() {
		deadline = deadlineSendImmediately
	}
	c.pacingDeadline = deadline
}

func (c *Conn) maybeSendAckOnlyPacket(now monotime.Time) error {
	if !c.handshakeConfirmed {
		ecn := c.sentPacketHandler.ECNMode(false)
		packet, err := c.packer.PackCoalescedPacket(true, c.maxPacketSize(), now, c.version)
		if err != nil {
			return err
		}
		if packet == nil {
			return nil
		}
		return c.sendPackedCoalescedPacket(packet, ecn, now)
	}

	ecn := c.sentPacketHandler.ECNMode(true)
	p, buf, err := c.packer.PackAckOnlyPacket(c.maxPacketSize(), now, c.version)
	if err != nil {
		if err == errNothingToPack {
			return nil
		}
		return err
	}
	c.logShortHeaderPacket(p.DestConnID, p.Ack, p.Frames, p.StreamFrames, p.PacketNumber, p.PacketNumberLen, p.KeyPhase, ecn, buf.Len(), false)
	c.registerPackedShortHeaderPacket(p, ecn, now)
	c.sendQueue.Send(buf, 0, ecn)
	return nil
}

func (c *Conn) sendProbePacket(sendMode ackhandler.SendMode, now monotime.Time) error {
	var encLevel protocol.EncryptionLevel
	//nolint:exhaustive // We only need to handle the PTO send modes here.
	switch sendMode {
	case ackhandler.SendPTOInitial:
		encLevel = protocol.EncryptionInitial
	case ackhandler.SendPTOHandshake:
		encLevel = protocol.EncryptionHandshake
	case ackhandler.SendPTOAppData:
		encLevel = protocol.Encryption1RTT
	default:
		return fmt.Errorf("connection BUG: unexpected send mode: %d", sendMode)
	}
	// Queue probe packets until we actually send out a packet,
	// or until there are no more packets to queue.
	var packet *coalescedPacket
	for packet == nil {
		if wasQueued := c.sentPacketHandler.QueueProbePacket(encLevel); !wasQueued {
			break
		}
		var err error
		packet, err = c.packer.PackPTOProbePacket(encLevel, c.maxPacketSize(), false, now, c.version)
		if err != nil {
			return err
		}
	}
	if packet == nil {
		var err error
		packet, err = c.packer.PackPTOProbePacket(encLevel, c.maxPacketSize(), true, now, c.version)
		if err != nil {
			return err
		}
	}
	if packet == nil || (len(packet.longHdrPackets) == 0 && packet.shortHdrPacket == nil) {
		return fmt.Errorf("connection BUG: couldn't pack %s probe packet: %v", encLevel, packet)
	}
	return c.sendPackedCoalescedPacket(packet, c.sentPacketHandler.ECNMode(packet.IsOnlyShortHeaderPacket()), now)
}

// appendOneShortHeaderPacket appends a new packet to the given packetBuffer.
// If there was nothing to pack, the returned size is 0.
func (c *Conn) appendOneShortHeaderPacket(buf *packetBuffer, maxSize protocol.ByteCount, ecn protocol.ECN, now monotime.Time) (protocol.ByteCount, error) {
	startLen := buf.Len()
	p, err := c.packer.AppendPacket(buf, maxSize, now, c.version)
	if err != nil {
		return 0, err
	}
	size := buf.Len() - startLen
	c.logShortHeaderPacket(p.DestConnID, p.Ack, p.Frames, p.StreamFrames, p.PacketNumber, p.PacketNumberLen, p.KeyPhase, ecn, size, false)
	c.registerPackedShortHeaderPacket(p, ecn, now)
	return size, nil
}

func (c *Conn) registerPackedShortHeaderPacket(p shortHeaderPacket, ecn protocol.ECN, now monotime.Time) {
	if p.IsPathProbePacket {
		c.sentPacketHandler.SentPacket(
			now,
			p.PacketNumber,
			protocol.InvalidPacketNumber,
			p.StreamFrames,
			p.Frames,
			protocol.Encryption1RTT,
			ecn,
			p.Length,
			p.IsPathMTUProbePacket,
			true,
		)
		return
	}
	if c.firstAckElicitingPacketAfterIdleSentTime.IsZero() && (len(p.StreamFrames) > 0 || ackhandler.HasAckElicitingFrames(p.Frames)) {
		c.firstAckElicitingPacketAfterIdleSentTime = now
	}

	largestAcked := protocol.InvalidPacketNumber
	if p.Ack != nil {
		largestAcked = p.Ack.LargestAcked()
	}
	c.sentPacketHandler.SentPacket(
		now,
		p.PacketNumber,
		largestAcked,
		p.StreamFrames,
		p.Frames,
		protocol.Encryption1RTT,
		ecn,
		p.Length,
		p.IsPathMTUProbePacket,
		false,
	)
	c.connIDManager.SentPacket()
}

func (c *Conn) sendPackedCoalescedPacket(packet *coalescedPacket, ecn protocol.ECN, now monotime.Time) error {
	c.logCoalescedPacket(packet, ecn)
	for _, p := range packet.longHdrPackets {
		if c.firstAckElicitingPacketAfterIdleSentTime.IsZero() && p.IsAckEliciting() {
			c.firstAckElicitingPacketAfterIdleSentTime = now
		}
		largestAcked := protocol.InvalidPacketNumber
		if p.ack != nil {
			largestAcked = p.ack.LargestAcked()
		}
		c.sentPacketHandler.SentPacket(
			now,
			p.header.PacketNumber,
			largestAcked,
			p.streamFrames,
			p.frames,
			p.EncryptionLevel(),
			ecn,
			p.length,
			false,
			false,
		)
		if c.perspective == protocol.PerspectiveClient && p.EncryptionLevel() == protocol.EncryptionHandshake &&
			!c.droppedInitialKeys {
			// On the client side, Initial keys are dropped as soon as the first Handshake packet is sent.
			// See Section 4.9.1 of RFC 9001.
			if err := c.dropEncryptionLevel(protocol.EncryptionInitial, now); err != nil {
				return err
			}
		}
	}
	if p := packet.shortHdrPacket; p != nil {
		if c.firstAckElicitingPacketAfterIdleSentTime.IsZero() && p.IsAckEliciting() {
			c.firstAckElicitingPacketAfterIdleSentTime = now
		}
		largestAcked := protocol.InvalidPacketNumber
		if p.Ack != nil {
			largestAcked = p.Ack.LargestAcked()
		}
		c.sentPacketHandler.SentPacket(
			now,
			p.PacketNumber,
			largestAcked,
			p.StreamFrames,
			p.Frames,
			protocol.Encryption1RTT,
			ecn,
			p.Length,
			p.IsPathMTUProbePacket,
			false,
		)
	}
	c.connIDManager.SentPacket()
	c.sendQueue.Send(packet.buffer, 0, ecn)
	return nil
}

func (c *Conn) sendConnectionClose(e error) ([]byte, error) {
	var packet *coalescedPacket
	var err error
	var transportErr *qerr.TransportError
	var applicationErr *qerr.ApplicationError
	if errors.As(e, &transportErr) {
		packet, err = c.packer.PackConnectionClose(transportErr, c.maxPacketSize(), c.version)
	} else if errors.As(e, &applicationErr) {
		packet, err = c.packer.PackApplicationClose(applicationErr, c.maxPacketSize(), c.version)
	} else {
		packet, err = c.packer.PackConnectionClose(&qerr.TransportError{
			ErrorCode:    qerr.InternalError,
			ErrorMessage: fmt.Sprintf("connection BUG: unspecified error type (msg: %s)", e.Error()),
		}, c.maxPacketSize(), c.version)
	}
	if err != nil {
		return nil, err
	}
	ecn := c.sentPacketHandler.ECNMode(packet.IsOnlyShortHeaderPacket())
	c.logCoalescedPacket(packet, ecn)
	return packet.buffer.Data, c.conn.Write(packet.buffer.Data, 0, ecn)
}

func (c *Conn) maxPacketSize() protocol.ByteCount {
	if c.mtuDiscoverer == nil {
		// Use the configured packet size on the client side.
		// If the server sends a max_udp_payload_size that's smaller than this size, we can ignore this:
		// Apparently the server still processed the (fully padded) Initial packet anyway.
		if c.perspective == protocol.PerspectiveClient {
			return protocol.ByteCount(c.config.InitialPacketSize)
		}
		// On the server side, there's no downside to using 1200 bytes until we received the client's transport
		// parameters:
		// * If the first packet didn't contain the entire ClientHello, all we can do is ACK that packet. We don't
		//   need a lot of bytes for that.
		// * If it did, we will have processed the transport parameters and initialized the MTU discoverer.
		return protocol.MinInitialPacketSize
	}
	return c.mtuDiscoverer.CurrentSize()
}

// AcceptStream returns the next stream opened by the peer, blocking until one is available.
func (c *Conn) AcceptStream(ctx context.Context) (*Stream, error) {
	return c.streamsMap.AcceptStream(ctx)
}

// AcceptUniStream returns the next unidirectional stream opened by the peer, blocking until one is available.
func (c *Conn) AcceptUniStream(ctx context.Context) (*ReceiveStream, error) {
	return c.streamsMap.AcceptUniStream(ctx)
}

// OpenStream opens a new bidirectional QUIC stream.
// There is no signaling to the peer about new streams:
// The peer can only accept the stream after data has been sent on the stream,
// or the stream has been reset or closed.
// When reaching the peer's stream limit, it is not possible to open a new stream until the
// peer raises the stream limit. In that case, a [StreamLimitReachedError] is returned.
func (c *Conn) OpenStream() (*Stream, error) {
	return c.streamsMap.OpenStream()
}

// OpenStreamSync opens a new bidirectional QUIC stream.
// It blocks until a new stream can be opened.
// There is no signaling to the peer about new streams:
// The peer can only accept the stream after data has been sent on the stream,
// or the stream has been reset or closed.
func (c *Conn) OpenStreamSync(ctx context.Context) (*Stream, error) {
	return c.streamsMap.OpenStreamSync(ctx)
}

// OpenUniStream opens a new outgoing unidirectional QUIC stream.
// There is no signaling to the peer about new streams:
// The peer can only accept the stream after data has been sent on the stream,
// or the stream has been reset or closed.
// When reaching the peer's stream limit, it is not possible to open a new stream until the
// peer raises the stream limit. In that case, a [StreamLimitReachedError] is returned.
func (c *Conn) OpenUniStream() (*SendStream, error) {
	return c.streamsMap.OpenUniStream()
}

// OpenUniStreamSync opens a new outgoing unidirectional QUIC stream.
// It blocks until a new stream can be opened.
// There is no signaling to the peer about new streams:
// The peer can only accept the stream after data has been sent on the stream,
// or the stream has been reset or closed.
func (c *Conn) OpenUniStreamSync(ctx context.Context) (*SendStream, error) {
	return c.streamsMap.OpenUniStreamSync(ctx)
}

func (c *Conn) newFlowController(id protocol.StreamID) flowcontrol.StreamFlowController {
	initialSendWindow := c.peerParams.InitialMaxStreamDataUni
	if id.Type() == protocol.StreamTypeBidi {
		if id.InitiatedBy() == c.perspective {
			initialSendWindow = c.peerParams.InitialMaxStreamDataBidiRemote
		} else {
			initialSendWindow = c.peerParams.InitialMaxStreamDataBidiLocal
		}
	}
	return flowcontrol.NewStreamFlowController(
		id,
		c.connFlowController,
		protocol.ByteCount(c.config.InitialStreamReceiveWindow),
		protocol.ByteCount(c.config.MaxStreamReceiveWindow),
		initialSendWindow,
		c.rttStats,
		c.logger,
	)
}

// scheduleSending signals that we have data for sending
func (c *Conn) scheduleSending() {
	select {
	case c.sendingScheduled <- struct{}{}:
	default:
	}
}

// tryQueueingUndecryptablePacket queues a packet for which we're missing the decryption keys.
// The qlogevents.PacketType is only used for logging purposes.
func (c *Conn) tryQueueingUndecryptablePacket(p receivedPacket, pt qlog.PacketType) {
	if c.handshakeComplete {
		panic("shouldn't queue undecryptable packets after handshake completion")
	}
	if len(c.undecryptablePackets)+1 > protocol.MaxUndecryptablePackets {
		if c.qlogger != nil {
			c.qlogger.RecordEvent(qlog.PacketDropped{
				Header: qlog.PacketHeader{
					PacketType:   pt,
					PacketNumber: protocol.InvalidPacketNumber,
				},
				Raw:     qlog.RawInfo{Length: int(p.Size())},
				Trigger: qlog.PacketDropDOSPrevention,
			})
		}
		c.logger.Infof("Dropping undecryptable packet (%d bytes). Undecryptable packet queue full.", p.Size())
		return
	}
	c.logger.Infof("Queueing packet (%d bytes) for later decryption", p.Size())
	if c.qlogger != nil {
		c.qlogger.RecordEvent(qlog.PacketBuffered{
			Header: qlog.PacketHeader{
				PacketType:   pt,
				PacketNumber: protocol.InvalidPacketNumber,
			},
			Raw: qlog.RawInfo{Length: int(p.Size())},
		})
	}
	c.undecryptablePackets = append(c.undecryptablePackets, p)
}

func (c *Conn) queueControlFrame(f wire.Frame) {
	c.framer.QueueControlFrame(f)
	c.scheduleSending()
}

func (c *Conn) onHasConnectionData() { c.scheduleSending() }

func (c *Conn) onHasStreamData(id protocol.StreamID, str *SendStream) {
	c.framer.AddActiveStream(id, str)
	c.scheduleSending()
}

func (c *Conn) onHasStreamControlFrame(id protocol.StreamID, str streamControlFrameGetter) {
	c.framer.AddStreamWithControlFrames(id, str)
	c.scheduleSending()
}

func (c *Conn) onStreamCompleted(id protocol.StreamID) {
	if err := c.streamsMap.DeleteStream(id); err != nil {
		c.closeLocal(err)
	}
	c.framer.RemoveActiveStream(id)
}

// SendDatagram sends a message using a QUIC datagram, as specified in RFC 9221,
// if the peer enabled datagram support.
// There is no delivery guarantee for DATAGRAM frames, they are not retransmitted if lost.
// The payload of the datagram needs to fit into a single QUIC packet.
// In addition, a datagram may be dropped before being sent out if the available packet size suddenly decreases.
// If the payload is too large to be sent at the current time, a DatagramTooLargeError is returned.
func (c *Conn) SendDatagram(p []byte) error {
	if !c.supportsDatagrams() {
		return errors.New("datagram support disabled")
	}

	f := &wire.DatagramFrame{DataLenPresent: true}
	// The payload size estimate is conservative.
	// Under many circumstances we could send a few more bytes.
	maxDataLen := min(
		f.MaxDataLen(c.peerParams.MaxDatagramFrameSize, c.version),
		protocol.ByteCount(c.currentMTUEstimate.Load()),
	)
	if protocol.ByteCount(len(p)) > maxDataLen {
		return &DatagramTooLargeError{MaxDatagramPayloadSize: int64(maxDataLen)}
	}
	f.Data = make([]byte, len(p))
	copy(f.Data, p)
	return c.datagramQueue.Add(f)
}

// ReceiveDatagram gets a message received in a QUIC datagram, as specified in RFC 9221.
func (c *Conn) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	if !c.config.EnableDatagrams {
		return nil, errors.New("datagram support disabled")
	}
	return c.datagramQueue.Receive(ctx)
}

// LocalAddr returns the local address of the QUIC connection.
func (c *Conn) LocalAddr() net.Addr { return c.conn.LocalAddr() }

// RemoteAddr returns the remote address of the QUIC connection.
func (c *Conn) RemoteAddr() net.Addr { return c.conn.RemoteAddr() }

// getPathManager lazily initializes the Conn's pathManagerOutgoing.
// May create multiple pathManagerOutgoing objects if called concurrently.
func (c *Conn) getPathManager() *pathManagerOutgoing {
	old := c.pathManagerOutgoing.Load()
	if old != nil {
		// Path manager is already initialized
		return old
	}

	// Initialize the path manager
	new := newPathManagerOutgoing(
		c.connIDManager.GetConnIDForPath,
		c.connIDManager.RetireConnIDForPath,
		c.scheduleSending,
	)
	if c.pathManagerOutgoing.CompareAndSwap(old, new) {
		return new
	}

	// Swap failed. A concurrent writer wrote first, use their value.
	return c.pathManagerOutgoing.Load()
}

func (c *Conn) AddPath(t *Transport) (*Path, error) {
	if c.perspective == protocol.PerspectiveServer {
		return nil, errors.New("server cannot initiate connection migration")
	}
	if c.peerParams.DisableActiveMigration {
		return nil, errors.New("server disabled connection migration")
	}
	if err := t.init(false); err != nil {
		return nil, err
	}
	return c.getPathManager().NewPath(
		t,
		200*time.Millisecond, // initial RTT estimate
		func() {
			runner := (*packetHandlerMap)(t)
			c.connIDGenerator.AddConnRunner(
				runner,
				connRunnerCallbacks{
					AddConnectionID:    func(connID protocol.ConnectionID) { runner.Add(connID, c) },
					RemoveConnectionID: runner.Remove,
					ReplaceWithClosed:  runner.ReplaceWithClosed,
				},
			)
		},
	), nil
}

// HandshakeComplete blocks until the handshake completes (or fails).
// For the client, data sent before completion of the handshake is encrypted with 0-RTT keys.
// For the server, data sent before completion of the handshake is encrypted with 1-RTT keys,
// however the client's identity is only verified once the handshake completes.
func (c *Conn) HandshakeComplete() <-chan struct{} {
	return c.handshakeCompleteChan
}

// QlogTrace returns the qlog trace of the QUIC connection.
// It is nil if qlog is not enabled.
func (c *Conn) QlogTrace() qlogwriter.Trace {
	return c.qlogTrace
}

func (c *Conn) NextConnection(ctx context.Context) (*Conn, error) {
	// The handshake might fail after the server rejected 0-RTT.
	// This could happen if the Finished message is malformed or never received.
	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	case <-c.Context().Done():
	case <-c.HandshakeComplete():
		c.streamsMap.UseResetMaps()
	}
	return c, nil
}

// estimateMaxPayloadSize estimates the maximum payload size for short header packets.
// It is not very sophisticated: it just subtracts the size of header (assuming the maximum
// connection ID length), and the size of the encryption tag.
func estimateMaxPayloadSize(mtu protocol.ByteCount) protocol.ByteCount {
	return mtu - 1 /* type byte */ - 20 /* maximum connection ID length */ - 16 /* tag size */
}
