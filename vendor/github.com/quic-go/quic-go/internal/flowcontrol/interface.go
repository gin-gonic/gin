package flowcontrol

import (
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
)

type flowController interface {
	// for sending
	SendWindowSize() protocol.ByteCount
	UpdateSendWindow(protocol.ByteCount) (updated bool)
	AddBytesSent(protocol.ByteCount)
	// for receiving
	GetWindowUpdate(monotime.Time) protocol.ByteCount // returns 0 if no update is necessary
}

// A StreamFlowController is a flow controller for a QUIC stream.
type StreamFlowController interface {
	flowController
	AddBytesRead(protocol.ByteCount) (hasStreamWindowUpdate, hasConnWindowUpdate bool)
	// UpdateHighestReceived is called when a new highest offset is received
	// final has to be to true if this is the final offset of the stream,
	// as contained in a STREAM frame with FIN bit, and the RESET_STREAM frame
	UpdateHighestReceived(offset protocol.ByteCount, final bool, now monotime.Time) error
	// Abandon is called when reading from the stream is aborted early,
	// and there won't be any further calls to AddBytesRead.
	Abandon()
	IsNewlyBlocked() bool
}

// The ConnectionFlowController is the flow controller for the connection.
type ConnectionFlowController interface {
	flowController
	AddBytesRead(protocol.ByteCount) (hasWindowUpdate bool)
	Reset() error
	IsNewlyBlocked() (bool, protocol.ByteCount)
}

type connectionFlowControllerI interface {
	ConnectionFlowController
	// The following two methods are not supposed to be called from outside this packet, but are needed internally
	// for sending
	EnsureMinimumWindowSize(protocol.ByteCount, monotime.Time)
	// for receiving
	IncrementHighestReceived(protocol.ByteCount, monotime.Time) error
}
