package handshake

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qlogwriter"
	"github.com/quic-go/quic-go/quicvarint"
)

type quicVersionContextKey struct{}

var QUICVersionContextKey = &quicVersionContextKey{}

const clientSessionStateRevision = 5

type cryptoSetup struct {
	tlsConf *tls.Config
	conn    *tls.QUICConn

	events []Event

	version protocol.Version

	ourParams  *wire.TransportParameters
	peerParams *wire.TransportParameters

	zeroRTTParameters *wire.TransportParameters
	allow0RTT         bool

	rttStats *utils.RTTStats

	qlogger qlogwriter.Recorder
	logger  utils.Logger

	perspective protocol.Perspective

	handshakeCompleteTime time.Time

	zeroRTTOpener LongHeaderOpener // only set for the server
	zeroRTTSealer LongHeaderSealer // only set for the client

	initialOpener LongHeaderOpener
	initialSealer LongHeaderSealer

	handshakeOpener LongHeaderOpener
	handshakeSealer LongHeaderSealer

	used0RTT atomic.Bool

	aead          *updatableAEAD
	has1RTTSealer bool
	has1RTTOpener bool
}

var _ CryptoSetup = &cryptoSetup{}

// NewCryptoSetupClient creates a new crypto setup for the client
func NewCryptoSetupClient(
	connID protocol.ConnectionID,
	tp *wire.TransportParameters,
	tlsConf *tls.Config,
	enable0RTT bool,
	rttStats *utils.RTTStats,
	qlogger qlogwriter.Recorder,
	logger utils.Logger,
	version protocol.Version,
) CryptoSetup {
	cs := newCryptoSetup(
		connID,
		tp,
		rttStats,
		qlogger,
		logger,
		protocol.PerspectiveClient,
		version,
	)

	tlsConf = tlsConf.Clone()
	tlsConf.MinVersion = tls.VersionTLS13
	cs.tlsConf = tlsConf
	cs.allow0RTT = enable0RTT

	cs.conn = tls.QUICClient(&tls.QUICConfig{
		TLSConfig:           tlsConf,
		EnableSessionEvents: true,
	})
	cs.conn.SetTransportParameters(cs.ourParams.Marshal(protocol.PerspectiveClient))

	return cs
}

// NewCryptoSetupServer creates a new crypto setup for the server
func NewCryptoSetupServer(
	connID protocol.ConnectionID,
	localAddr, remoteAddr net.Addr,
	tp *wire.TransportParameters,
	tlsConf *tls.Config,
	allow0RTT bool,
	rttStats *utils.RTTStats,
	qlogger qlogwriter.Recorder,
	logger utils.Logger,
	version protocol.Version,
) CryptoSetup {
	cs := newCryptoSetup(
		connID,
		tp,
		rttStats,
		qlogger,
		logger,
		protocol.PerspectiveServer,
		version,
	)
	cs.allow0RTT = allow0RTT

	tlsConf = setupConfigForServer(tlsConf, localAddr, remoteAddr)

	cs.tlsConf = tlsConf
	cs.conn = tls.QUICServer(&tls.QUICConfig{
		TLSConfig:           tlsConf,
		EnableSessionEvents: true,
	})
	return cs
}

func newCryptoSetup(
	connID protocol.ConnectionID,
	tp *wire.TransportParameters,
	rttStats *utils.RTTStats,
	qlogger qlogwriter.Recorder,
	logger utils.Logger,
	perspective protocol.Perspective,
	version protocol.Version,
) *cryptoSetup {
	initialSealer, initialOpener := NewInitialAEAD(connID, perspective, version)
	if qlogger != nil {
		qlogger.RecordEvent(qlog.KeyUpdated{
			Trigger: qlog.KeyUpdateTLS,
			KeyType: encLevelToKeyType(protocol.EncryptionInitial, protocol.PerspectiveClient),
		})
		qlogger.RecordEvent(qlog.KeyUpdated{
			Trigger: qlog.KeyUpdateTLS,
			KeyType: encLevelToKeyType(protocol.EncryptionInitial, protocol.PerspectiveServer),
		})
	}
	return &cryptoSetup{
		initialSealer: initialSealer,
		initialOpener: initialOpener,
		aead:          newUpdatableAEAD(rttStats, qlogger, logger, version),
		events:        make([]Event, 0, 16),
		ourParams:     tp,
		rttStats:      rttStats,
		qlogger:       qlogger,
		logger:        logger,
		perspective:   perspective,
		version:       version,
	}
}

func (h *cryptoSetup) ChangeConnectionID(id protocol.ConnectionID) {
	initialSealer, initialOpener := NewInitialAEAD(id, h.perspective, h.version)
	h.initialSealer = initialSealer
	h.initialOpener = initialOpener
	if h.qlogger != nil {
		h.qlogger.RecordEvent(qlog.KeyUpdated{
			Trigger: qlog.KeyUpdateTLS,
			KeyType: encLevelToKeyType(protocol.EncryptionInitial, protocol.PerspectiveClient),
		})
		h.qlogger.RecordEvent(qlog.KeyUpdated{
			Trigger: qlog.KeyUpdateTLS,
			KeyType: encLevelToKeyType(protocol.EncryptionInitial, protocol.PerspectiveServer),
		})
	}
}

func (h *cryptoSetup) SetLargest1RTTAcked(pn protocol.PacketNumber) error {
	return h.aead.SetLargestAcked(pn)
}

func (h *cryptoSetup) StartHandshake(ctx context.Context) error {
	err := h.conn.Start(context.WithValue(ctx, QUICVersionContextKey, h.version))
	if err != nil {
		return wrapError(err)
	}
	for {
		ev := h.conn.NextEvent()
		if err := h.handleEvent(ev); err != nil {
			return wrapError(err)
		}
		if ev.Kind == tls.QUICNoEvent {
			break
		}
	}
	if h.perspective == protocol.PerspectiveClient {
		if h.zeroRTTSealer != nil && h.zeroRTTParameters != nil {
			h.logger.Debugf("Doing 0-RTT.")
			h.events = append(h.events, Event{Kind: EventRestoredTransportParameters, TransportParameters: h.zeroRTTParameters})
		} else {
			h.logger.Debugf("Not doing 0-RTT. Has sealer: %t, has params: %t", h.zeroRTTSealer != nil, h.zeroRTTParameters != nil)
		}
	}
	return nil
}

// Close closes the crypto setup.
// It aborts the handshake, if it is still running.
func (h *cryptoSetup) Close() error {
	return h.conn.Close()
}

// HandleMessage handles a TLS handshake message.
// It is called by the crypto streams when a new message is available.
func (h *cryptoSetup) HandleMessage(data []byte, encLevel protocol.EncryptionLevel) error {
	if err := h.handleMessage(data, encLevel); err != nil {
		return wrapError(err)
	}
	return nil
}

func (h *cryptoSetup) handleMessage(data []byte, encLevel protocol.EncryptionLevel) error {
	if err := h.conn.HandleData(encLevel.ToTLSEncryptionLevel(), data); err != nil {
		return err
	}
	for {
		ev := h.conn.NextEvent()
		if err := h.handleEvent(ev); err != nil {
			return err
		}
		if ev.Kind == tls.QUICNoEvent {
			return nil
		}
	}
}

func (h *cryptoSetup) handleEvent(ev tls.QUICEvent) (err error) {
	switch ev.Kind {
	case tls.QUICNoEvent:
		return nil
	case tls.QUICSetReadSecret:
		h.setReadKey(ev.Level, ev.Suite, ev.Data)
		return nil
	case tls.QUICSetWriteSecret:
		h.setWriteKey(ev.Level, ev.Suite, ev.Data)
		return nil
	case tls.QUICTransportParameters:
		return h.handleTransportParameters(ev.Data)
	case tls.QUICTransportParametersRequired:
		h.conn.SetTransportParameters(h.ourParams.Marshal(h.perspective))
		return nil
	case tls.QUICRejectedEarlyData:
		h.rejected0RTT()
		return nil
	case tls.QUICWriteData:
		h.writeRecord(ev.Level, ev.Data)
		return nil
	case tls.QUICHandshakeDone:
		h.handshakeComplete()
		return nil
	case tls.QUICStoreSession:
		if h.perspective == protocol.PerspectiveServer {
			panic("cryptoSetup BUG: unexpected QUICStoreSession event for the server")
		}
		ev.SessionState.Extra = append(
			ev.SessionState.Extra,
			addSessionStateExtraPrefix(h.marshalDataForSessionState(ev.SessionState.EarlyData)),
		)
		return h.conn.StoreSession(ev.SessionState)
	case tls.QUICResumeSession:
		var allowEarlyData bool
		switch h.perspective {
		case protocol.PerspectiveClient:
			// for clients, this event occurs when a session ticket is selected
			allowEarlyData = h.handleDataFromSessionState(
				findSessionStateExtraData(ev.SessionState.Extra),
				ev.SessionState.EarlyData,
			)
		case protocol.PerspectiveServer:
			// for servers, this event occurs when receiving the client's session ticket
			allowEarlyData = h.handleSessionTicket(
				findSessionStateExtraData(ev.SessionState.Extra),
				ev.SessionState.EarlyData,
			)
		}
		if ev.SessionState.EarlyData {
			ev.SessionState.EarlyData = allowEarlyData
		}
		return nil
	default:
		// Unknown events should be ignored.
		// crypto/tls will ensure that this is safe to do.
		// See the discussion following https://github.com/golang/go/issues/68124#issuecomment-2187042510 for details.
		return nil
	}
}

func (h *cryptoSetup) NextEvent() Event {
	if len(h.events) == 0 {
		return Event{Kind: EventNoEvent}
	}
	ev := h.events[0]
	h.events = h.events[1:]
	return ev
}

func (h *cryptoSetup) handleTransportParameters(data []byte) error {
	var tp wire.TransportParameters
	if err := tp.Unmarshal(data, h.perspective.Opposite()); err != nil {
		return err
	}
	h.peerParams = &tp
	h.events = append(h.events, Event{Kind: EventReceivedTransportParameters, TransportParameters: h.peerParams})
	return nil
}

// must be called after receiving the transport parameters
func (h *cryptoSetup) marshalDataForSessionState(earlyData bool) []byte {
	b := make([]byte, 0, 256)
	b = quicvarint.Append(b, clientSessionStateRevision)
	if earlyData {
		// only save the transport parameters for 0-RTT enabled session tickets
		return h.peerParams.MarshalForSessionTicket(b)
	}
	return b
}

func (h *cryptoSetup) handleDataFromSessionState(data []byte, earlyData bool) (allowEarlyData bool) {
	tp, err := decodeDataFromSessionState(data, earlyData)
	if err != nil {
		h.logger.Debugf("Restoring of transport parameters from session ticket failed: %s", err.Error())
		return
	}
	// The session ticket might have been saved from a connection that allowed 0-RTT,
	// and therefore contain transport parameters.
	// Only use them if 0-RTT is actually used on the new connection.
	if tp != nil && h.allow0RTT {
		h.zeroRTTParameters = tp
		return true
	}
	return false
}

func decodeDataFromSessionState(b []byte, earlyData bool) (*wire.TransportParameters, error) {
	ver, l, err := quicvarint.Parse(b)
	if err != nil {
		return nil, err
	}
	b = b[l:]
	if ver != clientSessionStateRevision {
		return nil, fmt.Errorf("mismatching version. Got %d, expected %d", ver, clientSessionStateRevision)
	}
	if !earlyData {
		return nil, nil
	}
	var tp wire.TransportParameters
	if err := tp.UnmarshalFromSessionTicket(b); err != nil {
		return nil, err
	}
	return &tp, nil
}

func (h *cryptoSetup) getDataForSessionTicket() []byte {
	return (&sessionTicket{
		Parameters: h.ourParams,
	}).Marshal()
}

// GetSessionTicket generates a new session ticket.
// Due to limitations in crypto/tls, it's only possible to generate a single session ticket per connection.
// It is only valid for the server.
func (h *cryptoSetup) GetSessionTicket() ([]byte, error) {
	if err := h.conn.SendSessionTicket(tls.QUICSessionTicketOptions{
		EarlyData: h.allow0RTT,
		Extra:     [][]byte{addSessionStateExtraPrefix(h.getDataForSessionTicket())},
	}); err != nil {
		// Session tickets might be disabled by tls.Config.SessionTicketsDisabled.
		// We can't check h.tlsConfig here, since the actual config might have been obtained from
		// the GetConfigForClient callback.
		// See https://github.com/golang/go/issues/62032.
		// Once that issue is resolved, this error assertion can be removed.
		if strings.Contains(err.Error(), "session ticket keys unavailable") {
			return nil, nil
		}
		return nil, err
	}
	ev := h.conn.NextEvent()
	if ev.Kind != tls.QUICWriteData || ev.Level != tls.QUICEncryptionLevelApplication {
		panic("crypto/tls bug: where's my session ticket?")
	}
	ticket := ev.Data
	if ev := h.conn.NextEvent(); ev.Kind != tls.QUICNoEvent {
		panic("crypto/tls bug: why more than one ticket?")
	}
	return ticket, nil
}

// handleSessionTicket is called for the server when receiving the client's session ticket.
// It reads parameters from the session ticket and checks whether to accept 0-RTT if the session ticket enabled 0-RTT.
// Note that the fact that the session ticket allows 0-RTT doesn't mean that the actual TLS handshake enables 0-RTT:
// A client may use a 0-RTT enabled session to resume a TLS session without using 0-RTT.
func (h *cryptoSetup) handleSessionTicket(data []byte, using0RTT bool) (allowEarlyData bool) {
	var t sessionTicket
	if err := t.Unmarshal(data); err != nil {
		h.logger.Debugf("Unmarshalling session ticket failed: %s", err.Error())
		return false
	}
	if !using0RTT {
		return false
	}
	valid := h.ourParams.ValidFor0RTT(t.Parameters)
	if !valid {
		h.logger.Debugf("Transport parameters changed. Rejecting 0-RTT.")
		return false
	}
	if !h.allow0RTT {
		h.logger.Debugf("0-RTT not allowed. Rejecting 0-RTT.")
		return false
	}
	return true
}

// rejected0RTT is called for the client when the server rejects 0-RTT.
func (h *cryptoSetup) rejected0RTT() {
	h.logger.Debugf("0-RTT was rejected. Dropping 0-RTT keys.")

	had0RTTKeys := h.zeroRTTSealer != nil
	h.zeroRTTSealer = nil

	if had0RTTKeys {
		h.events = append(h.events, Event{Kind: EventDiscard0RTTKeys})
	}
}

func (h *cryptoSetup) setReadKey(el tls.QUICEncryptionLevel, suiteID uint16, trafficSecret []byte) {
	suite := getCipherSuite(suiteID)
	//nolint:exhaustive // The TLS stack doesn't export Initial keys.
	switch el {
	case tls.QUICEncryptionLevelEarly:
		if h.perspective == protocol.PerspectiveClient {
			panic("Received 0-RTT read key for the client")
		}
		h.zeroRTTOpener = newLongHeaderOpener(
			createAEAD(suite, trafficSecret, h.version),
			newHeaderProtector(suite, trafficSecret, true, h.version),
		)
		h.used0RTT.Store(true)
		if h.logger.Debug() {
			h.logger.Debugf("Installed 0-RTT Read keys (using %s)", tls.CipherSuiteName(suite.ID))
		}
	case tls.QUICEncryptionLevelHandshake:
		h.handshakeOpener = newLongHeaderOpener(
			createAEAD(suite, trafficSecret, h.version),
			newHeaderProtector(suite, trafficSecret, true, h.version),
		)
		if h.logger.Debug() {
			h.logger.Debugf("Installed Handshake Read keys (using %s)", tls.CipherSuiteName(suite.ID))
		}
	case tls.QUICEncryptionLevelApplication:
		h.aead.SetReadKey(suite, trafficSecret)
		h.has1RTTOpener = true
		if h.logger.Debug() {
			h.logger.Debugf("Installed 1-RTT Read keys (using %s)", tls.CipherSuiteName(suite.ID))
		}
	default:
		panic("unexpected read encryption level")
	}
	h.events = append(h.events, Event{Kind: EventReceivedReadKeys})
	if h.qlogger != nil {
		h.qlogger.RecordEvent(qlog.KeyUpdated{
			Trigger: qlog.KeyUpdateTLS,
			KeyType: encLevelToKeyType(protocol.FromTLSEncryptionLevel(el), h.perspective.Opposite()),
		})
	}
}

func (h *cryptoSetup) setWriteKey(el tls.QUICEncryptionLevel, suiteID uint16, trafficSecret []byte) {
	suite := getCipherSuite(suiteID)
	//nolint:exhaustive // The TLS stack doesn't export Initial keys.
	switch el {
	case tls.QUICEncryptionLevelEarly:
		if h.perspective == protocol.PerspectiveServer {
			panic("Received 0-RTT write key for the server")
		}
		h.zeroRTTSealer = newLongHeaderSealer(
			createAEAD(suite, trafficSecret, h.version),
			newHeaderProtector(suite, trafficSecret, true, h.version),
		)
		if h.logger.Debug() {
			h.logger.Debugf("Installed 0-RTT Write keys (using %s)", tls.CipherSuiteName(suite.ID))
		}
		if h.qlogger != nil {
			h.qlogger.RecordEvent(qlog.KeyUpdated{
				Trigger: qlog.KeyUpdateTLS,
				KeyType: encLevelToKeyType(protocol.Encryption0RTT, h.perspective),
			})
		}
		// don't set used0RTT here. 0-RTT might still get rejected.
		return
	case tls.QUICEncryptionLevelHandshake:
		h.handshakeSealer = newLongHeaderSealer(
			createAEAD(suite, trafficSecret, h.version),
			newHeaderProtector(suite, trafficSecret, true, h.version),
		)
		if h.logger.Debug() {
			h.logger.Debugf("Installed Handshake Write keys (using %s)", tls.CipherSuiteName(suite.ID))
		}
	case tls.QUICEncryptionLevelApplication:
		h.aead.SetWriteKey(suite, trafficSecret)
		h.has1RTTSealer = true
		if h.logger.Debug() {
			h.logger.Debugf("Installed 1-RTT Write keys (using %s)", tls.CipherSuiteName(suite.ID))
		}
		if h.zeroRTTSealer != nil {
			// Once we receive handshake keys, we know that 0-RTT was not rejected.
			h.used0RTT.Store(true)
			h.zeroRTTSealer = nil
			h.logger.Debugf("Dropping 0-RTT keys.")
			if h.qlogger != nil {
				h.qlogger.RecordEvent(qlog.KeyDiscarded{KeyType: qlog.KeyTypeClient0RTT})
			}
		}
	default:
		panic("unexpected write encryption level")
	}
	if h.qlogger != nil {
		h.qlogger.RecordEvent(qlog.KeyUpdated{
			Trigger: qlog.KeyUpdateTLS,
			KeyType: encLevelToKeyType(protocol.FromTLSEncryptionLevel(el), h.perspective),
		})
	}
}

// writeRecord is called when TLS writes data
func (h *cryptoSetup) writeRecord(encLevel tls.QUICEncryptionLevel, p []byte) {
	//nolint:exhaustive // handshake records can only be written for Initial and Handshake.
	switch encLevel {
	case tls.QUICEncryptionLevelInitial:
		h.events = append(h.events, Event{Kind: EventWriteInitialData, Data: p})
	case tls.QUICEncryptionLevelHandshake:
		h.events = append(h.events, Event{Kind: EventWriteHandshakeData, Data: p})
	case tls.QUICEncryptionLevelApplication:
		panic("unexpected write")
	default:
		panic(fmt.Sprintf("unexpected write encryption level: %s", encLevel))
	}
}

func (h *cryptoSetup) DiscardInitialKeys() {
	dropped := h.initialOpener != nil
	h.initialOpener = nil
	h.initialSealer = nil
	if dropped {
		h.logger.Debugf("Dropping Initial keys.")
		if h.qlogger != nil {
			h.qlogger.RecordEvent(qlog.KeyDiscarded{KeyType: qlog.KeyTypeClientInitial})
			h.qlogger.RecordEvent(qlog.KeyDiscarded{KeyType: qlog.KeyTypeServerInitial})
		}
	}
}

func (h *cryptoSetup) handshakeComplete() {
	h.handshakeCompleteTime = time.Now()
	h.events = append(h.events, Event{Kind: EventHandshakeComplete})
}

func (h *cryptoSetup) SetHandshakeConfirmed() {
	h.aead.SetHandshakeConfirmed()
	// drop Handshake keys
	var dropped bool
	if h.handshakeOpener != nil {
		h.handshakeOpener = nil
		h.handshakeSealer = nil
		dropped = true
	}
	if dropped {
		h.logger.Debugf("Dropping Handshake keys.")
		if h.qlogger != nil {
			h.qlogger.RecordEvent(qlog.KeyDiscarded{KeyType: qlog.KeyTypeClientHandshake})
			h.qlogger.RecordEvent(qlog.KeyDiscarded{KeyType: qlog.KeyTypeServerHandshake})
		}
	}
}

func (h *cryptoSetup) GetInitialSealer() (LongHeaderSealer, error) {
	if h.initialSealer == nil {
		return nil, ErrKeysDropped
	}
	return h.initialSealer, nil
}

func (h *cryptoSetup) Get0RTTSealer() (LongHeaderSealer, error) {
	if h.zeroRTTSealer == nil {
		return nil, ErrKeysDropped
	}
	return h.zeroRTTSealer, nil
}

func (h *cryptoSetup) GetHandshakeSealer() (LongHeaderSealer, error) {
	if h.handshakeSealer == nil {
		if h.initialSealer == nil {
			return nil, ErrKeysDropped
		}
		return nil, ErrKeysNotYetAvailable
	}
	return h.handshakeSealer, nil
}

func (h *cryptoSetup) Get1RTTSealer() (ShortHeaderSealer, error) {
	if !h.has1RTTSealer {
		return nil, ErrKeysNotYetAvailable
	}
	return h.aead, nil
}

func (h *cryptoSetup) GetInitialOpener() (LongHeaderOpener, error) {
	if h.initialOpener == nil {
		return nil, ErrKeysDropped
	}
	return h.initialOpener, nil
}

func (h *cryptoSetup) Get0RTTOpener() (LongHeaderOpener, error) {
	if h.zeroRTTOpener == nil {
		if h.initialOpener != nil {
			return nil, ErrKeysNotYetAvailable
		}
		// if the initial opener is also not available, the keys were already dropped
		return nil, ErrKeysDropped
	}
	return h.zeroRTTOpener, nil
}

func (h *cryptoSetup) GetHandshakeOpener() (LongHeaderOpener, error) {
	if h.handshakeOpener == nil {
		if h.initialOpener != nil {
			return nil, ErrKeysNotYetAvailable
		}
		// if the initial opener is also not available, the keys were already dropped
		return nil, ErrKeysDropped
	}
	return h.handshakeOpener, nil
}

func (h *cryptoSetup) Get1RTTOpener() (ShortHeaderOpener, error) {
	if h.zeroRTTOpener != nil && time.Since(h.handshakeCompleteTime) > 3*h.rttStats.PTO(true) {
		h.zeroRTTOpener = nil
		h.logger.Debugf("Dropping 0-RTT keys.")
		if h.qlogger != nil {
			h.qlogger.RecordEvent(qlog.KeyDiscarded{KeyType: qlog.KeyTypeClient0RTT})
		}
	}

	if !h.has1RTTOpener {
		return nil, ErrKeysNotYetAvailable
	}
	return h.aead, nil
}

func (h *cryptoSetup) ConnectionState() ConnectionState {
	return ConnectionState{
		ConnectionState: h.conn.ConnectionState(),
		Used0RTT:        h.used0RTT.Load(),
	}
}

func wrapError(err error) error {
	if alertErr := tls.AlertError(0); errors.As(err, &alertErr) {
		return qerr.NewLocalCryptoError(uint8(alertErr), err)
	}
	return &qerr.TransportError{ErrorCode: qerr.InternalError, ErrorMessage: err.Error()}
}

func encLevelToKeyType(encLevel protocol.EncryptionLevel, pers protocol.Perspective) qlog.KeyType {
	if pers == protocol.PerspectiveServer {
		switch encLevel {
		case protocol.EncryptionInitial:
			return qlog.KeyTypeServerInitial
		case protocol.EncryptionHandshake:
			return qlog.KeyTypeServerHandshake
		case protocol.Encryption0RTT:
			return qlog.KeyTypeServer0RTT
		case protocol.Encryption1RTT:
			return qlog.KeyTypeServer1RTT
		default:
			return ""
		}
	}
	switch encLevel {
	case protocol.EncryptionInitial:
		return qlog.KeyTypeClientInitial
	case protocol.EncryptionHandshake:
		return qlog.KeyTypeClientHandshake
	case protocol.Encryption0RTT:
		return qlog.KeyTypeClient0RTT
	case protocol.Encryption1RTT:
		return qlog.KeyTypeClient1RTT
	default:
		return ""
	}
}
