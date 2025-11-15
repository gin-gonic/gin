package protocol

import "github.com/quic-go/quic-go/quicvarint"

// StreamType encodes if this is a unidirectional or bidirectional stream
type StreamType uint8

const (
	// StreamTypeUni is a unidirectional stream
	StreamTypeUni StreamType = iota
	// StreamTypeBidi is a bidirectional stream
	StreamTypeBidi
)

// InvalidPacketNumber is a stream ID that is invalid.
// The first valid stream ID in QUIC is 0.
const InvalidStreamID StreamID = -1

// StreamNum is the stream number
type StreamNum int64

const (
	// InvalidStreamNum is an invalid stream number.
	InvalidStreamNum = -1
	// MaxStreamCount is the maximum stream count value that can be sent in MAX_STREAMS frames
	// and as the stream count in the transport parameters
	MaxStreamCount StreamNum = 1 << 60
	// MaxStreamID is the maximum stream ID
	MaxStreamID StreamID = quicvarint.Max
)

const (
	// FirstOutgoingBidiStreamClient is the first bidirectional stream opened by the client
	FirstOutgoingBidiStreamClient StreamID = 0
	// FirstOutgoingUniStreamClient is the first unidirectional stream opened by the client
	FirstOutgoingUniStreamClient StreamID = 2
	// FirstOutgoingBidiStreamServer is the first bidirectional stream opened by the server
	FirstOutgoingBidiStreamServer StreamID = 1
	// FirstOutgoingUniStreamServer is the first unidirectional stream opened by the server
	FirstOutgoingUniStreamServer StreamID = 3
)

const (
	// FirstIncomingBidiStreamServer is the first bidirectional stream accepted by the server
	FirstIncomingBidiStreamServer = FirstOutgoingBidiStreamClient
	// FirstIncomingUniStreamServer is the first unidirectional stream accepted by the server
	FirstIncomingUniStreamServer = FirstOutgoingUniStreamClient
	// FirstIncomingBidiStreamClient is the first bidirectional stream accepted by the client
	FirstIncomingBidiStreamClient = FirstOutgoingBidiStreamServer
	// FirstIncomingUniStreamClient is the first unidirectional stream accepted by the client
	FirstIncomingUniStreamClient = FirstOutgoingUniStreamServer
)

// StreamID calculates the stream ID.
func (s StreamNum) StreamID(stype StreamType, pers Perspective) StreamID {
	if s == 0 {
		return InvalidStreamID
	}
	var first StreamID
	switch stype {
	case StreamTypeBidi:
		switch pers {
		case PerspectiveClient:
			first = 0
		case PerspectiveServer:
			first = 1
		}
	case StreamTypeUni:
		switch pers {
		case PerspectiveClient:
			first = 2
		case PerspectiveServer:
			first = 3
		}
	}
	return first + 4*StreamID(s-1)
}

// A StreamID in QUIC
type StreamID int64

// InitiatedBy says if the stream was initiated by the client or by the server
func (s StreamID) InitiatedBy() Perspective {
	if s%2 == 0 {
		return PerspectiveClient
	}
	return PerspectiveServer
}

// Type says if this is a unidirectional or bidirectional stream
func (s StreamID) Type() StreamType {
	if s%4 >= 2 {
		return StreamTypeUni
	}
	return StreamTypeBidi
}

// StreamNum returns how many streams in total are below this
// Example: for stream 9 it returns 3 (i.e. streams 1, 5 and 9)
func (s StreamID) StreamNum() StreamNum {
	return StreamNum(s/4) + 1
}
