package quic

import (
	"context"
	"net"
	"os"
	"sync"
	"time"

	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/flowcontrol"
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
)

type deadlineError struct{}

func (deadlineError) Error() string   { return "deadline exceeded" }
func (deadlineError) Temporary() bool { return true }
func (deadlineError) Timeout() bool   { return true }
func (deadlineError) Unwrap() error   { return os.ErrDeadlineExceeded }

var errDeadline net.Error = &deadlineError{}

// The streamSender is notified by the stream about various events.
type streamSender interface {
	onHasConnectionData()
	onHasStreamData(protocol.StreamID, *SendStream)
	onHasStreamControlFrame(protocol.StreamID, streamControlFrameGetter)
	// must be called without holding the mutex that is acquired by closeForShutdown
	onStreamCompleted(protocol.StreamID)
}

// Each of the both stream halves gets its own uniStreamSender.
// This is necessary in order to keep track when both halves have been completed.
type uniStreamSender struct {
	streamSender
	onStreamCompletedImpl       func()
	onHasStreamControlFrameImpl func(protocol.StreamID, streamControlFrameGetter)
}

func (s *uniStreamSender) onHasStreamData(id protocol.StreamID, str *SendStream) {
	s.streamSender.onHasStreamData(id, str)
}
func (s *uniStreamSender) onStreamCompleted(protocol.StreamID) { s.onStreamCompletedImpl() }
func (s *uniStreamSender) onHasStreamControlFrame(id protocol.StreamID, str streamControlFrameGetter) {
	s.onHasStreamControlFrameImpl(id, str)
}

var _ streamSender = &uniStreamSender{}

type Stream struct {
	receiveStr *ReceiveStream
	sendStr    *SendStream

	completedMutex         sync.Mutex
	sender                 streamSender
	receiveStreamCompleted bool
	sendStreamCompleted    bool
}

var (
	_ outgoingStream            = &Stream{}
	_ sendStreamFrameHandler    = &Stream{}
	_ receiveStreamFrameHandler = &Stream{}
)

// newStream creates a new Stream
func newStream(
	ctx context.Context,
	streamID protocol.StreamID,
	sender streamSender,
	flowController flowcontrol.StreamFlowController,
	supportsResetStreamAt bool,
) *Stream {
	s := &Stream{sender: sender}
	senderForSendStream := &uniStreamSender{
		streamSender: sender,
		onStreamCompletedImpl: func() {
			s.completedMutex.Lock()
			s.sendStreamCompleted = true
			s.checkIfCompleted()
			s.completedMutex.Unlock()
		},
		onHasStreamControlFrameImpl: func(id protocol.StreamID, str streamControlFrameGetter) {
			sender.onHasStreamControlFrame(streamID, s)
		},
	}
	s.sendStr = newSendStream(ctx, streamID, senderForSendStream, flowController, supportsResetStreamAt)
	senderForReceiveStream := &uniStreamSender{
		streamSender: sender,
		onStreamCompletedImpl: func() {
			s.completedMutex.Lock()
			s.receiveStreamCompleted = true
			s.checkIfCompleted()
			s.completedMutex.Unlock()
		},
		onHasStreamControlFrameImpl: func(id protocol.StreamID, str streamControlFrameGetter) {
			sender.onHasStreamControlFrame(streamID, s)
		},
	}
	s.receiveStr = newReceiveStream(streamID, senderForReceiveStream, flowController)
	return s
}

// StreamID returns the stream ID.
func (s *Stream) StreamID() protocol.StreamID {
	// the result is same for receiveStream and sendStream
	return s.sendStr.StreamID()
}

// Read reads data from the stream.
// Read can be made to time out using [Stream.SetReadDeadline] and [Stream.SetDeadline].
// If the stream was canceled, the error is a [StreamError].
func (s *Stream) Read(p []byte) (int, error) {
	return s.receiveStr.Read(p)
}

// Write writes data to the stream.
// Write can be made to time out using [Stream.SetWriteDeadline] or [Stream.SetDeadline].
// If the stream was canceled, the error is a [StreamError].
func (s *Stream) Write(p []byte) (int, error) {
	return s.sendStr.Write(p)
}

// CancelWrite aborts sending on this stream.
// See [SendStream.CancelWrite] for more details.
func (s *Stream) CancelWrite(errorCode StreamErrorCode) {
	s.sendStr.CancelWrite(errorCode)
}

// CancelRead aborts receiving on this stream.
// See [ReceiveStream.CancelRead] for more details.
func (s *Stream) CancelRead(errorCode StreamErrorCode) {
	s.receiveStr.CancelRead(errorCode)
}

// The Context is canceled as soon as the write-side of the stream is closed.
// See [SendStream.Context] for more details.
func (s *Stream) Context() context.Context {
	return s.sendStr.Context()
}

// Close closes the send-direction of the stream.
// It does not close the receive-direction of the stream.
func (s *Stream) Close() error {
	return s.sendStr.Close()
}

func (s *Stream) handleResetStreamFrame(frame *wire.ResetStreamFrame, rcvTime monotime.Time) error {
	return s.receiveStr.handleResetStreamFrame(frame, rcvTime)
}

func (s *Stream) handleStreamFrame(frame *wire.StreamFrame, rcvTime monotime.Time) error {
	return s.receiveStr.handleStreamFrame(frame, rcvTime)
}

func (s *Stream) handleStopSendingFrame(frame *wire.StopSendingFrame) {
	s.sendStr.handleStopSendingFrame(frame)
}

func (s *Stream) updateSendWindow(limit protocol.ByteCount) {
	s.sendStr.updateSendWindow(limit)
}

func (s *Stream) enableResetStreamAt() {
	s.sendStr.enableResetStreamAt()
}

func (s *Stream) popStreamFrame(maxBytes protocol.ByteCount, v protocol.Version) (_ ackhandler.StreamFrame, _ *wire.StreamDataBlockedFrame, hasMore bool) {
	return s.sendStr.popStreamFrame(maxBytes, v)
}

func (s *Stream) getControlFrame(now monotime.Time) (_ ackhandler.Frame, ok, hasMore bool) {
	f, ok, _ := s.sendStr.getControlFrame(now)
	if ok {
		return f, true, true
	}
	return s.receiveStr.getControlFrame(now)
}

// SetReadDeadline sets the deadline for future Read calls.
// See [ReceiveStream.SetReadDeadline] for more details.
func (s *Stream) SetReadDeadline(t time.Time) error {
	return s.receiveStr.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
// See [SendStream.SetWriteDeadline] for more details.
func (s *Stream) SetWriteDeadline(t time.Time) error {
	return s.sendStr.SetWriteDeadline(t)
}

// SetDeadline sets the read and write deadlines associated with the stream.
// It is equivalent to calling both SetReadDeadline and SetWriteDeadline.
func (s *Stream) SetDeadline(t time.Time) error {
	_ = s.receiveStr.SetReadDeadline(t) // SetReadDeadline never errors
	_ = s.sendStr.SetWriteDeadline(t)   // SetWriteDeadline never errors
	return nil
}

// CloseForShutdown closes a stream abruptly.
// It makes Read and Write unblock (and return the error) immediately.
// The peer will NOT be informed about this: the stream is closed without sending a FIN or RST.
func (s *Stream) closeForShutdown(err error) {
	s.sendStr.closeForShutdown(err)
	s.receiveStr.closeForShutdown(err)
}

// checkIfCompleted is called from the uniStreamSender, when one of the stream halves is completed.
// It makes sure that the onStreamCompleted callback is only called if both receive and send side have completed.
func (s *Stream) checkIfCompleted() {
	if s.sendStreamCompleted && s.receiveStreamCompleted {
		s.sender.onStreamCompleted(s.StreamID())
	}
}
