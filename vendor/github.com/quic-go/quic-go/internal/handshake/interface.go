package handshake

import (
	"context"
	"crypto/tls"
	"errors"
	"io"

	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
)

var (
	// ErrKeysNotYetAvailable is returned when an opener or a sealer is requested for an encryption level,
	// but the corresponding opener has not yet been initialized
	// This can happen when packets arrive out of order.
	ErrKeysNotYetAvailable = errors.New("CryptoSetup: keys at this encryption level not yet available")
	// ErrKeysDropped is returned when an opener or a sealer is requested for an encryption level,
	// but the corresponding keys have already been dropped.
	ErrKeysDropped = errors.New("CryptoSetup: keys were already dropped")
	// ErrDecryptionFailed is returned when the AEAD fails to open the packet.
	ErrDecryptionFailed = errors.New("decryption failed")
)

type headerDecryptor interface {
	DecryptHeader(sample []byte, firstByte *byte, pnBytes []byte)
}

// LongHeaderOpener opens a long header packet
type LongHeaderOpener interface {
	headerDecryptor
	DecodePacketNumber(wirePN protocol.PacketNumber, wirePNLen protocol.PacketNumberLen) protocol.PacketNumber
	Open(dst, src []byte, pn protocol.PacketNumber, associatedData []byte) ([]byte, error)
}

// ShortHeaderOpener opens a short header packet
type ShortHeaderOpener interface {
	headerDecryptor
	DecodePacketNumber(wirePN protocol.PacketNumber, wirePNLen protocol.PacketNumberLen) protocol.PacketNumber
	Open(dst, src []byte, rcvTime monotime.Time, pn protocol.PacketNumber, kp protocol.KeyPhaseBit, associatedData []byte) ([]byte, error)
}

// LongHeaderSealer seals a long header packet
type LongHeaderSealer interface {
	Seal(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) []byte
	EncryptHeader(sample []byte, firstByte *byte, pnBytes []byte)
	Overhead() int
}

// ShortHeaderSealer seals a short header packet
type ShortHeaderSealer interface {
	LongHeaderSealer
	KeyPhase() protocol.KeyPhaseBit
}

type ConnectionState struct {
	tls.ConnectionState
	Used0RTT bool
}

// EventKind is the kind of handshake event.
type EventKind uint8

const (
	// EventNoEvent signals that there are no new handshake events
	EventNoEvent EventKind = iota + 1
	// EventWriteInitialData contains new CRYPTO data to send at the Initial encryption level
	EventWriteInitialData
	// EventWriteHandshakeData contains new CRYPTO data to send at the Handshake encryption level
	EventWriteHandshakeData
	// EventReceivedReadKeys signals that new decryption keys are available.
	// It doesn't say which encryption level those keys are for.
	EventReceivedReadKeys
	// EventDiscard0RTTKeys signals that the Handshake keys were discarded.
	EventDiscard0RTTKeys
	// EventReceivedTransportParameters contains the transport parameters sent by the peer.
	EventReceivedTransportParameters
	// EventRestoredTransportParameters contains the transport parameters restored from the session ticket.
	// It is only used for the client.
	EventRestoredTransportParameters
	// EventHandshakeComplete signals that the TLS handshake was completed.
	EventHandshakeComplete
)

func (k EventKind) String() string {
	switch k {
	case EventNoEvent:
		return "EventNoEvent"
	case EventWriteInitialData:
		return "EventWriteInitialData"
	case EventWriteHandshakeData:
		return "EventWriteHandshakeData"
	case EventReceivedReadKeys:
		return "EventReceivedReadKeys"
	case EventDiscard0RTTKeys:
		return "EventDiscard0RTTKeys"
	case EventReceivedTransportParameters:
		return "EventReceivedTransportParameters"
	case EventRestoredTransportParameters:
		return "EventRestoredTransportParameters"
	case EventHandshakeComplete:
		return "EventHandshakeComplete"
	default:
		return "Unknown EventKind"
	}
}

// Event is a handshake event.
type Event struct {
	Kind                EventKind
	Data                []byte
	TransportParameters *wire.TransportParameters
}

// CryptoSetup handles the handshake and protecting / unprotecting packets
type CryptoSetup interface {
	StartHandshake(context.Context) error
	io.Closer
	ChangeConnectionID(protocol.ConnectionID)
	GetSessionTicket() ([]byte, error)

	HandleMessage([]byte, protocol.EncryptionLevel) error
	NextEvent() Event

	SetLargest1RTTAcked(protocol.PacketNumber) error
	DiscardInitialKeys()
	SetHandshakeConfirmed()
	ConnectionState() ConnectionState

	GetInitialOpener() (LongHeaderOpener, error)
	GetHandshakeOpener() (LongHeaderOpener, error)
	Get0RTTOpener() (LongHeaderOpener, error)
	Get1RTTOpener() (ShortHeaderOpener, error)

	GetInitialSealer() (LongHeaderSealer, error)
	GetHandshakeSealer() (LongHeaderSealer, error)
	Get0RTTSealer() (LongHeaderSealer, error)
	Get1RTTSealer() (ShortHeaderSealer, error)
}
