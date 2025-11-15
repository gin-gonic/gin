package quic

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/flowcontrol"
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/wire"
)

// A ReceiveStream is a unidirectional Receive Stream.
type ReceiveStream struct {
	mutex sync.Mutex

	streamID protocol.StreamID

	sender streamSender

	frameQueue  *frameSorter
	finalOffset protocol.ByteCount

	currentFrame       []byte
	currentFrameDone   func()
	readPosInFrame     int
	currentFrameIsLast bool // is the currentFrame the last frame on this stream

	queuedStopSending   bool
	queuedMaxStreamData bool

	// Set once we read the io.EOF or the cancellation error.
	// Note that for local cancellations, this doesn't necessarily mean that we know the final offset yet.
	errorRead           bool
	completed           bool // set once we've called streamSender.onStreamCompleted
	cancelledRemotely   bool
	cancelledLocally    bool
	cancelErr           *StreamError
	closeForShutdownErr error

	readPos      protocol.ByteCount
	reliableSize protocol.ByteCount

	readChan chan struct{}
	readOnce chan struct{} // cap: 1, to protect against concurrent use of Read
	deadline monotime.Time

	flowController flowcontrol.StreamFlowController
}

var (
	_ streamControlFrameGetter  = &ReceiveStream{}
	_ receiveStreamFrameHandler = &ReceiveStream{}
)

func newReceiveStream(
	streamID protocol.StreamID,
	sender streamSender,
	flowController flowcontrol.StreamFlowController,
) *ReceiveStream {
	return &ReceiveStream{
		streamID:       streamID,
		sender:         sender,
		flowController: flowController,
		frameQueue:     newFrameSorter(),
		readChan:       make(chan struct{}, 1),
		readOnce:       make(chan struct{}, 1),
		finalOffset:    protocol.MaxByteCount,
	}
}

// StreamID returns the stream ID.
func (s *ReceiveStream) StreamID() protocol.StreamID {
	return s.streamID
}

// Read reads data from the stream.
// Read can be made to time out using [ReceiveStream.SetReadDeadline].
// If the stream was canceled, the error is a [StreamError].
func (s *ReceiveStream) Read(p []byte) (int, error) {
	// Concurrent use of Read is not permitted (and doesn't make any sense),
	// but sometimes people do it anyway.
	// Make sure that we only execute one call at any given time to avoid hard to debug failures.
	s.readOnce <- struct{}{}
	defer func() { <-s.readOnce }()

	s.mutex.Lock()
	queuedStreamWindowUpdate, queuedConnWindowUpdate, n, err := s.readImpl(p)
	completed := s.isNewlyCompleted()
	s.mutex.Unlock()

	if completed {
		s.sender.onStreamCompleted(s.streamID)
	}
	if queuedStreamWindowUpdate {
		s.sender.onHasStreamControlFrame(s.streamID, s)
	}
	if queuedConnWindowUpdate {
		s.sender.onHasConnectionData()
	}
	return n, err
}

func (s *ReceiveStream) isNewlyCompleted() bool {
	if s.completed {
		return false
	}
	// We need to know the final offset (either via FIN or RESET_STREAM) for flow control accounting.
	if s.finalOffset == protocol.MaxByteCount {
		return false
	}
	// We're done with the stream if it was cancelled locally...
	if s.cancelledLocally {
		s.completed = true
		return true
	}
	// ... or if the error (either io.EOF or the reset error) was read
	if s.errorRead {
		s.completed = true
		return true
	}
	return false
}

func (s *ReceiveStream) readImpl(p []byte) (hasStreamWindowUpdate bool, hasConnWindowUpdate bool, _ int, _ error) {
	if s.currentFrameIsLast && s.currentFrame == nil {
		s.errorRead = true
		return false, false, 0, io.EOF
	}
	if s.cancelledLocally || (s.cancelledRemotely && s.readPos >= s.reliableSize) {
		s.errorRead = true
		return false, false, 0, s.cancelErr
	}
	if s.closeForShutdownErr != nil {
		return false, false, 0, s.closeForShutdownErr
	}

	var bytesRead int
	var deadlineTimer *time.Timer
	for bytesRead < len(p) {
		if s.currentFrame == nil || s.readPosInFrame >= len(s.currentFrame) {
			s.dequeueNextFrame()
		}
		if s.currentFrame == nil && bytesRead > 0 {
			return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, s.closeForShutdownErr
		}

		for {
			// Stop waiting on errors
			if s.closeForShutdownErr != nil {
				return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, s.closeForShutdownErr
			}
			if s.cancelledLocally || (s.cancelledRemotely && s.readPos >= s.reliableSize) {
				s.errorRead = true
				return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, s.cancelErr
			}

			deadline := s.deadline
			if !deadline.IsZero() {
				if !monotime.Now().Before(deadline) {
					return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, errDeadline
				}
				if deadlineTimer == nil {
					deadlineTimer = time.NewTimer(monotime.Until(deadline))
					defer deadlineTimer.Stop()
				} else {
					deadlineTimer.Reset(monotime.Until(deadline))
				}
			}

			if s.currentFrame != nil || s.currentFrameIsLast {
				break
			}

			s.mutex.Unlock()
			if deadline.IsZero() {
				<-s.readChan
			} else {
				select {
				case <-s.readChan:
				case <-deadlineTimer.C:
				}
			}
			s.mutex.Lock()
			if s.currentFrame == nil {
				s.dequeueNextFrame()
			}
		}

		if bytesRead > len(p) {
			return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, fmt.Errorf("BUG: bytesRead (%d) > len(p) (%d) in stream.Read", bytesRead, len(p))
		}
		if s.readPosInFrame > len(s.currentFrame) {
			return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, fmt.Errorf("BUG: readPosInFrame (%d) > frame.DataLen (%d) in stream.Read", s.readPosInFrame, len(s.currentFrame))
		}
		m := copy(p[bytesRead:], s.currentFrame[s.readPosInFrame:])

		// when a RESET_STREAM was received, the flow controller was already
		// informed about the final offset for this stream
		if !s.cancelledRemotely || s.readPos < s.reliableSize {
			hasStream, hasConn := s.flowController.AddBytesRead(protocol.ByteCount(m))
			if hasStream {
				s.queuedMaxStreamData = true
				hasStreamWindowUpdate = true
			}
			if hasConn {
				hasConnWindowUpdate = true
			}
		}

		s.readPosInFrame += m
		s.readPos += protocol.ByteCount(m)
		bytesRead += m

		if s.cancelledRemotely && s.readPos >= s.reliableSize {
			s.flowController.Abandon()
		}

		if s.readPosInFrame >= len(s.currentFrame) && s.currentFrameIsLast {
			s.currentFrame = nil
			if s.currentFrameDone != nil {
				s.currentFrameDone()
			}
			s.errorRead = true
			return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, io.EOF
		}
	}
	if s.cancelledRemotely && s.readPos >= s.reliableSize {
		s.errorRead = true
		return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, s.cancelErr
	}
	return hasStreamWindowUpdate, hasConnWindowUpdate, bytesRead, nil
}

func (s *ReceiveStream) dequeueNextFrame() {
	var offset protocol.ByteCount
	// We're done with the last frame. Release the buffer.
	if s.currentFrameDone != nil {
		s.currentFrameDone()
	}
	offset, s.currentFrame, s.currentFrameDone = s.frameQueue.Pop()
	s.currentFrameIsLast = offset+protocol.ByteCount(len(s.currentFrame)) >= s.finalOffset && !s.cancelledRemotely
	s.readPosInFrame = 0
}

// CancelRead aborts receiving on this stream.
// It instructs the peer to stop transmitting stream data.
// Read will unblock immediately, and future Read calls will fail.
// When called multiple times or after reading the io.EOF it is a no-op.
func (s *ReceiveStream) CancelRead(errorCode StreamErrorCode) {
	s.mutex.Lock()
	queuedNewControlFrame := s.cancelReadImpl(errorCode)
	completed := s.isNewlyCompleted()
	s.mutex.Unlock()

	if queuedNewControlFrame {
		s.sender.onHasStreamControlFrame(s.streamID, s)
	}
	if completed {
		s.flowController.Abandon()
		s.sender.onStreamCompleted(s.streamID)
	}
}

func (s *ReceiveStream) cancelReadImpl(errorCode qerr.StreamErrorCode) (queuedNewControlFrame bool) {
	if s.cancelledLocally { // duplicate call to CancelRead
		return false
	}
	if s.closeForShutdownErr != nil {
		return false
	}
	s.cancelledLocally = true
	if s.errorRead || s.cancelledRemotely {
		return false
	}
	s.queuedStopSending = true
	s.cancelErr = &StreamError{StreamID: s.streamID, ErrorCode: errorCode, Remote: false}
	s.signalRead()
	return true
}

func (s *ReceiveStream) handleStreamFrame(frame *wire.StreamFrame, now monotime.Time) error {
	s.mutex.Lock()
	err := s.handleStreamFrameImpl(frame, now)
	completed := s.isNewlyCompleted()
	s.mutex.Unlock()

	if completed {
		s.flowController.Abandon()
		s.sender.onStreamCompleted(s.streamID)
	}
	return err
}

func (s *ReceiveStream) handleStreamFrameImpl(frame *wire.StreamFrame, now monotime.Time) error {
	maxOffset := frame.Offset + frame.DataLen()
	if err := s.flowController.UpdateHighestReceived(maxOffset, frame.Fin, now); err != nil {
		return err
	}
	if frame.Fin {
		s.finalOffset = maxOffset
	}
	if s.cancelledLocally {
		return nil
	}
	if err := s.frameQueue.Push(frame.Data, frame.Offset, frame.PutBack); err != nil {
		return err
	}
	s.signalRead()
	return nil
}

func (s *ReceiveStream) handleResetStreamFrame(frame *wire.ResetStreamFrame, now monotime.Time) error {
	s.mutex.Lock()
	err := s.handleResetStreamFrameImpl(frame, now)
	completed := s.isNewlyCompleted()
	s.mutex.Unlock()

	if completed {
		s.sender.onStreamCompleted(s.streamID)
	}
	return err
}

func (s *ReceiveStream) handleResetStreamFrameImpl(frame *wire.ResetStreamFrame, now monotime.Time) error {
	if s.closeForShutdownErr != nil {
		return nil
	}
	if err := s.flowController.UpdateHighestReceived(frame.FinalSize, true, now); err != nil {
		return err
	}
	s.finalOffset = frame.FinalSize

	// senders are allowed to reduce the reliable size, but frames might have been reordered
	if (!s.cancelledRemotely && s.reliableSize == 0) || frame.ReliableSize < s.reliableSize {
		s.reliableSize = frame.ReliableSize
	}
	if s.readPos >= s.reliableSize {
		// calling Abandon multiple times is a no-op
		s.flowController.Abandon()
	}
	// ignore duplicate RESET_STREAM frames for this stream (after checking their final offset)
	if s.cancelledRemotely {
		return nil
	}

	// don't save the error if the RESET_STREAM frames was received after CancelRead was called
	if s.cancelledLocally {
		return nil
	}
	s.cancelledRemotely = true
	s.cancelErr = &StreamError{StreamID: s.streamID, ErrorCode: frame.ErrorCode, Remote: true}
	s.signalRead()
	return nil
}

func (s *ReceiveStream) getControlFrame(now monotime.Time) (_ ackhandler.Frame, ok, hasMore bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.queuedStopSending && !s.queuedMaxStreamData {
		return ackhandler.Frame{}, false, false
	}
	if s.queuedStopSending {
		s.queuedStopSending = false
		return ackhandler.Frame{
			Frame: &wire.StopSendingFrame{StreamID: s.streamID, ErrorCode: s.cancelErr.ErrorCode},
		}, true, s.queuedMaxStreamData
	}

	s.queuedMaxStreamData = false
	return ackhandler.Frame{
		Frame: &wire.MaxStreamDataFrame{
			StreamID:          s.streamID,
			MaximumStreamData: s.flowController.GetWindowUpdate(now),
		},
	}, true, false
}

// SetReadDeadline sets the deadline for future Read calls and
// any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (s *ReceiveStream) SetReadDeadline(t time.Time) error {
	s.mutex.Lock()
	s.deadline = monotime.FromTime(t)
	s.mutex.Unlock()
	s.signalRead()
	return nil
}

// CloseForShutdown closes a stream abruptly.
// It makes Read unblock (and return the error) immediately.
// The peer will NOT be informed about this: the stream is closed without sending a FIN or RESET.
func (s *ReceiveStream) closeForShutdown(err error) {
	s.mutex.Lock()
	s.closeForShutdownErr = err
	s.mutex.Unlock()
	s.signalRead()
}

// signalRead performs a non-blocking send on the readChan
func (s *ReceiveStream) signalRead() {
	select {
	case s.readChan <- struct{}{}:
	default:
	}
}
