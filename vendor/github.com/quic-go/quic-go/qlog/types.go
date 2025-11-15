package qlog

import (
	"fmt"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
)

type (
	ConnectionID             = protocol.ConnectionID
	ArbitraryLenConnectionID = protocol.ArbitraryLenConnectionID
	Version                  = protocol.Version
	PacketNumber             = protocol.PacketNumber
	EncryptionLevel          = protocol.EncryptionLevel
	KeyPhaseBit              = protocol.KeyPhaseBit
	KeyPhase                 = protocol.KeyPhase
	StreamID                 = protocol.StreamID
	TransportErrorCode       = qerr.TransportErrorCode
	ApplicationErrorCode     = qerr.ApplicationErrorCode
)

const (
	// KeyPhaseZero is key phase bit 0
	KeyPhaseZero = protocol.KeyPhaseZero
	// KeyPhaseOne is key phase bit 1
	KeyPhaseOne = protocol.KeyPhaseOne
)

// ECN represents the Explicit Congestion Notification value.
type ECN string

const (
	// ECNUnsupported means that no ECN value was set / received
	ECNUnsupported ECN = ""
	// ECTNot is Not-ECT
	ECTNot ECN = "Not-ECT"
	// ECT0 is ECT(0)
	ECT0 ECN = "ECT(0)"
	// ECT1 is ECT(1)
	ECT1 ECN = "ECT(1)"
	// ECNCE is CE
	ECNCE ECN = "CE"
)

type Initiator string

const (
	InitiatorLocal  Initiator = "local"
	InitiatorRemote Initiator = "remote"
)

type streamType protocol.StreamType

func (s streamType) String() string {
	switch protocol.StreamType(s) {
	case protocol.StreamTypeUni:
		return "unidirectional"
	case protocol.StreamTypeBidi:
		return "bidirectional"
	default:
		return "unknown stream type"
	}
}

type version protocol.Version

func (v version) String() string {
	return fmt.Sprintf("%x", uint32(v))
}

func encLevelToPacketNumberSpace(encLevel protocol.EncryptionLevel) string {
	switch encLevel {
	case protocol.EncryptionInitial:
		return "initial"
	case protocol.EncryptionHandshake:
		return "handshake"
	case protocol.Encryption0RTT, protocol.Encryption1RTT:
		return "application_data"
	default:
		return "unknown encryption level"
	}
}

// KeyType represents the type of cryptographic key used in QUIC connections.
type KeyType string

const (
	// KeyTypeServerInitial represents the server's initial secret key.
	KeyTypeServerInitial KeyType = "server_initial_secret"
	// KeyTypeClientInitial represents the client's initial secret key.
	KeyTypeClientInitial KeyType = "client_initial_secret"
	// KeyTypeServerHandshake represents the server's handshake secret key.
	KeyTypeServerHandshake KeyType = "server_handshake_secret"
	// KeyTypeClientHandshake represents the client's handshake secret key.
	KeyTypeClientHandshake KeyType = "client_handshake_secret"
	// KeyTypeServer0RTT represents the server's 0-RTT secret key.
	KeyTypeServer0RTT KeyType = "server_0rtt_secret"
	// KeyTypeClient0RTT represents the client's 0-RTT secret key.
	KeyTypeClient0RTT KeyType = "client_0rtt_secret"
	// KeyTypeServer1RTT represents the server's 1-RTT secret key.
	KeyTypeServer1RTT KeyType = "server_1rtt_secret"
	// KeyTypeClient1RTT represents the client's 1-RTT secret key.
	KeyTypeClient1RTT KeyType = "client_1rtt_secret"
)

// KeyUpdateTrigger describes what caused a key update event.
type KeyUpdateTrigger string

const (
	// KeyUpdateTLS indicates the key update was triggered by TLS.
	KeyUpdateTLS KeyUpdateTrigger = "tls"
	// KeyUpdateRemote indicates the key update was triggered by the remote peer.
	KeyUpdateRemote KeyUpdateTrigger = "remote_update"
	// KeyUpdateLocal indicates the key update was triggered locally.
	KeyUpdateLocal KeyUpdateTrigger = "local_update"
)

type transportError uint64

func (e transportError) String() string {
	switch qerr.TransportErrorCode(e) {
	case qerr.NoError:
		return "no_error"
	case qerr.InternalError:
		return "internal_error"
	case qerr.ConnectionRefused:
		return "connection_refused"
	case qerr.FlowControlError:
		return "flow_control_error"
	case qerr.StreamLimitError:
		return "stream_limit_error"
	case qerr.StreamStateError:
		return "stream_state_error"
	case qerr.FinalSizeError:
		return "final_size_error"
	case qerr.FrameEncodingError:
		return "frame_encoding_error"
	case qerr.TransportParameterError:
		return "transport_parameter_error"
	case qerr.ConnectionIDLimitError:
		return "connection_id_limit_error"
	case qerr.ProtocolViolation:
		return "protocol_violation"
	case qerr.InvalidToken:
		return "invalid_token"
	case qerr.ApplicationErrorErrorCode:
		return "application_error"
	case qerr.CryptoBufferExceeded:
		return "crypto_buffer_exceeded"
	case qerr.KeyUpdateError:
		return "key_update_error"
	case qerr.AEADLimitReached:
		return "aead_limit_reached"
	case qerr.NoViablePathError:
		return "no_viable_path"
	default:
		return ""
	}
}

type PacketType string

const (
	// PacketTypeInitial represents an Initial packet
	PacketTypeInitial PacketType = "initial"
	// PacketTypeHandshake represents a Handshake packet
	PacketTypeHandshake PacketType = "handshake"
	// PacketTypeRetry represents a Retry packet
	PacketTypeRetry PacketType = "retry"
	// PacketType0RTT represents a 0-RTT packet
	PacketType0RTT PacketType = "0RTT"
	// PacketTypeVersionNegotiation represents a Version Negotiation packet
	PacketTypeVersionNegotiation PacketType = "version_negotiation"
	// PacketTypeStatelessReset represents a Stateless Reset packet
	PacketTypeStatelessReset PacketType = "stateless_reset"
	// PacketType1RTT represents a 1-RTT packet
	PacketType1RTT PacketType = "1RTT"
	// // PacketTypeNotDetermined represents a packet type that could not be determined
	// PacketTypeNotDetermined packetType = ""
)

func EncryptionLevelToPacketType(l EncryptionLevel) PacketType {
	switch l {
	case protocol.EncryptionInitial:
		return PacketTypeInitial
	case protocol.EncryptionHandshake:
		return PacketTypeHandshake
	case protocol.Encryption0RTT:
		return PacketType0RTT
	case protocol.Encryption1RTT:
		return PacketType1RTT
	default:
		panic(fmt.Sprintf("unknown encryption level: %d", l))
	}
}

type PacketLossReason string

const (
	// PacketLossReorderingThreshold is used when a packet is declared lost due to reordering threshold
	PacketLossReorderingThreshold PacketLossReason = "reordering_threshold"
	// PacketLossTimeThreshold is used when a packet is declared lost due to time threshold
	PacketLossTimeThreshold PacketLossReason = "time_threshold"
)

type PacketDropReason string

const (
	// PacketDropKeyUnavailable is used when a packet is dropped because keys are unavailable
	PacketDropKeyUnavailable PacketDropReason = "key_unavailable"
	// PacketDropUnknownConnectionID is used when a packet is dropped because the connection ID is unknown
	PacketDropUnknownConnectionID PacketDropReason = "unknown_connection_id"
	// PacketDropHeaderParseError is used when a packet is dropped because header parsing failed
	PacketDropHeaderParseError PacketDropReason = "header_parse_error"
	// PacketDropPayloadDecryptError is used when a packet is dropped because decrypting the payload failed
	PacketDropPayloadDecryptError PacketDropReason = "payload_decrypt_error"
	// PacketDropProtocolViolation is used when a packet is dropped due to a protocol violation
	PacketDropProtocolViolation PacketDropReason = "protocol_violation"
	// PacketDropDOSPrevention is used when a packet is dropped to mitigate a DoS attack
	PacketDropDOSPrevention PacketDropReason = "dos_prevention"
	// PacketDropUnsupportedVersion is used when a packet is dropped because the version is not supported
	PacketDropUnsupportedVersion PacketDropReason = "unsupported_version"
	// PacketDropUnexpectedPacket is used when an unexpected packet is received
	PacketDropUnexpectedPacket PacketDropReason = "unexpected_packet"
	// PacketDropUnexpectedSourceConnectionID is used when a packet with an unexpected source connection ID is received
	PacketDropUnexpectedSourceConnectionID PacketDropReason = "unexpected_source_connection_id"
	// PacketDropUnexpectedVersion is used when a packet with an unexpected version is received
	PacketDropUnexpectedVersion PacketDropReason = "unexpected_version"
	// PacketDropDuplicate is used when a duplicate packet is received
	PacketDropDuplicate PacketDropReason = "duplicate"
)

type LossTimerUpdateType string

const (
	LossTimerUpdateTypeSet       LossTimerUpdateType = "set"
	LossTimerUpdateTypeExpired   LossTimerUpdateType = "expired"
	LossTimerUpdateTypeCancelled LossTimerUpdateType = "cancelled"
)

type TimerType string

const (
	// TimerTypeACK represents an ACK timer
	TimerTypeACK TimerType = "ack"
	// TimerTypePTO represents a PTO (Probe Timeout) timer
	TimerTypePTO TimerType = "pto"
	// TimerTypePathProbe represents a path probe timer
	TimerTypePathProbe TimerType = "path_probe"
)

type CongestionState string

const (
	// CongestionStateSlowStart is the slow start phase of Reno / Cubic
	CongestionStateSlowStart CongestionState = "slow_start"
	// CongestionStateCongestionAvoidance is the congestion avoidance phase of Reno / Cubic
	CongestionStateCongestionAvoidance CongestionState = "congestion_avoidance"
	// CongestionStateRecovery is the recovery phase of Reno / Cubic
	CongestionStateRecovery CongestionState = "recovery"
	// CongestionStateApplicationLimited means that the congestion controller is application limited
	CongestionStateApplicationLimited CongestionState = "application_limited"
)

func (s CongestionState) String() string {
	return string(s)
}

// ECNState is the state of the ECN state machine (see Appendix A.4 of RFC 9000)
type ECNState string

const (
	// ECNStateTesting is the testing state
	ECNStateTesting ECNState = "testing"
	// ECNStateUnknown is the unknown state
	ECNStateUnknown ECNState = "unknown"
	// ECNStateFailed is the failed state
	ECNStateFailed ECNState = "failed"
	// ECNStateCapable is the capable state
	ECNStateCapable ECNState = "capable"
)

type ConnectionCloseTrigger string

const (
	// IdleTimeout indicates the connection was closed due to idle timeout
	ConnectionCloseTriggerIdleTimeout ConnectionCloseTrigger = "idle_timeout"
	// Application indicates the connection was closed by the application
	ConnectionCloseTriggerApplication ConnectionCloseTrigger = "application"
	// VersionMismatch indicates the connection was closed due to a QUIC version mismatch
	ConnectionCloseTriggerVersionMismatch ConnectionCloseTrigger = "version_mismatch"
	// StatelessReset indicates the connection was closed due to receiving a stateless reset from the peer
	ConnectionCloseTriggerStatelessReset ConnectionCloseTrigger = "stateless_reset"
)
