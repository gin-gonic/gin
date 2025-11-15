package quic

import (
	"fmt"

	"github.com/quic-go/quic-go/internal/qerr"
)

type (
	// TransportError indicates an error that occurred on the QUIC transport layer.
	// Every transport error other than CONNECTION_REFUSED and APPLICATION_ERROR is
	// likely a bug in the implementation.
	TransportError = qerr.TransportError
	// ApplicationError is an application-defined error.
	ApplicationError = qerr.ApplicationError
	// VersionNegotiationError indicates a failure to negotiate a QUIC version.
	VersionNegotiationError = qerr.VersionNegotiationError
	// StatelessResetError indicates a stateless reset was received.
	// This can happen when the peer reboots, or when packets are misrouted.
	// See section 10.3 of RFC 9000 for details.
	StatelessResetError = qerr.StatelessResetError
	// IdleTimeoutError indicates that the connection timed out because it was inactive for too long.
	IdleTimeoutError = qerr.IdleTimeoutError
	// HandshakeTimeoutError indicates that the connection timed out before completing the handshake.
	HandshakeTimeoutError = qerr.HandshakeTimeoutError
)

type (
	// TransportErrorCode is a QUIC transport error code, see section 20 of RFC 9000.
	TransportErrorCode = qerr.TransportErrorCode
	// ApplicationErrorCode is an QUIC application error code.
	ApplicationErrorCode = qerr.ApplicationErrorCode
	// StreamErrorCode is a QUIC stream error code. The meaning of the value is defined by the application.
	StreamErrorCode = qerr.StreamErrorCode
)

const (
	// NoError is the NO_ERROR transport error code.
	NoError = qerr.NoError
	// InternalError is the INTERNAL_ERROR transport error code.
	InternalError = qerr.InternalError
	// ConnectionRefused is the CONNECTION_REFUSED transport error code.
	ConnectionRefused = qerr.ConnectionRefused
	// FlowControlError is the FLOW_CONTROL_ERROR transport error code.
	FlowControlError = qerr.FlowControlError
	// StreamLimitError is the STREAM_LIMIT_ERROR transport error code.
	StreamLimitError = qerr.StreamLimitError
	// StreamStateError is the STREAM_STATE_ERROR transport error code.
	StreamStateError = qerr.StreamStateError
	// FinalSizeError is the FINAL_SIZE_ERROR transport error code.
	FinalSizeError = qerr.FinalSizeError
	// FrameEncodingError is the FRAME_ENCODING_ERROR transport error code.
	FrameEncodingError = qerr.FrameEncodingError
	// TransportParameterError is the TRANSPORT_PARAMETER_ERROR transport error code.
	TransportParameterError = qerr.TransportParameterError
	// ConnectionIDLimitError is the CONNECTION_ID_LIMIT_ERROR transport error code.
	ConnectionIDLimitError = qerr.ConnectionIDLimitError
	// ProtocolViolation is the PROTOCOL_VIOLATION transport error code.
	ProtocolViolation = qerr.ProtocolViolation
	// InvalidToken is the INVALID_TOKEN transport error code.
	InvalidToken = qerr.InvalidToken
	// ApplicationErrorErrorCode is the APPLICATION_ERROR transport error code.
	ApplicationErrorErrorCode = qerr.ApplicationErrorErrorCode
	// CryptoBufferExceeded is the CRYPTO_BUFFER_EXCEEDED transport error code.
	CryptoBufferExceeded = qerr.CryptoBufferExceeded
	// KeyUpdateError is the KEY_UPDATE_ERROR transport error code.
	KeyUpdateError = qerr.KeyUpdateError
	// AEADLimitReached is the AEAD_LIMIT_REACHED transport error code.
	AEADLimitReached = qerr.AEADLimitReached
	// NoViablePathError is the NO_VIABLE_PATH_ERROR transport error code.
	NoViablePathError = qerr.NoViablePathError
)

// A StreamError is used to signal stream cancellations.
// It is returned from the Read and Write methods of the [ReceiveStream], [SendStream] and [Stream].
type StreamError struct {
	StreamID  StreamID
	ErrorCode StreamErrorCode
	Remote    bool
}

func (e *StreamError) Is(target error) bool {
	t, ok := target.(*StreamError)
	return ok && e.StreamID == t.StreamID && e.ErrorCode == t.ErrorCode && e.Remote == t.Remote
}

func (e *StreamError) Error() string {
	pers := "local"
	if e.Remote {
		pers = "remote"
	}
	return fmt.Sprintf("stream %d canceled by %s with error code %d", e.StreamID, pers, e.ErrorCode)
}

// DatagramTooLargeError is returned from Conn.SendDatagram if the payload is too large to be sent.
type DatagramTooLargeError struct {
	MaxDatagramPayloadSize int64
}

func (e *DatagramTooLargeError) Is(target error) bool {
	t, ok := target.(*DatagramTooLargeError)
	return ok && e.MaxDatagramPayloadSize == t.MaxDatagramPayloadSize
}

func (e *DatagramTooLargeError) Error() string { return "DATAGRAM frame too large" }
