package ackhandler

import "github.com/quic-go/quic-go/internal/wire"

// IsFrameTypeAckEliciting returns true if the frame is ack-eliciting.
func IsFrameTypeAckEliciting(t wire.FrameType) bool {
	//nolint:exhaustive // The default case catches the rest.
	switch t {
	case wire.FrameTypeAck, wire.FrameTypeAckECN:
		return false
	case wire.FrameTypeConnectionClose, wire.FrameTypeApplicationClose:
		return false
	default:
		return true
	}
}

// IsFrameAckEliciting returns true if the frame is ack-eliciting.
func IsFrameAckEliciting(f wire.Frame) bool {
	_, isAck := f.(*wire.AckFrame)
	_, isConnectionClose := f.(*wire.ConnectionCloseFrame)
	return !isAck && !isConnectionClose
}

// HasAckElicitingFrames returns true if at least one frame is ack-eliciting.
func HasAckElicitingFrames(fs []Frame) bool {
	for _, f := range fs {
		if IsFrameAckEliciting(f.Frame) {
			return true
		}
	}
	return false
}
