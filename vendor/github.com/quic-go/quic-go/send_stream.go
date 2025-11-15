package quic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/flowcontrol"
	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
)

// A SendStream is a unidirectional Send Stream.
type SendStream struct {
	mutex sync.Mutex

	numOutstandingFrames int64 // outstanding STREAM and RESET_STREAM frames
	retransmissionQueue  []*wire.StreamFrame

	ctx       context.Context
	ctxCancel context.CancelCauseFunc

	streamID protocol.StreamID
	sender   streamSender

	// reliableSize is the portion of the stream that needs to be transmitted reliably,
	// even if the stream is cancelled.
	// This requires the peer to support RESET_STREAM_AT.
	// This value should not be accessed directly, but only through the reliableOffset method.
	// This method returns 0 if the peer doesn't support the RESET_STREAM_AT extension.
	reliableSize protocol.ByteCount
	writeOffset  protocol.ByteCount

	shutdownErr            error
	resetErr               *StreamError
	queuedResetStreamFrame *wire.ResetStreamFrame

	supportsResetStreamAt bool
	finishedWriting       bool // set once Close() is called
	finSent               bool // set when a STREAM_FRAME with FIN bit has been sent
	// Set when the application knows about the cancellation.
	// This can happen because the application called CancelWrite,
	// or because Write returned the error (for remote cancellations).
	cancellationFlagged bool
	completed           bool // set when this stream has been reported to the streamSender as completed

	dataForWriting []byte // during a Write() call, this slice is the part of p that still needs to be sent out
	nextFrame      *wire.StreamFrame

	writeChan chan struct{}
	writeOnce chan struct{}
	deadline  monotime.Time

	flowController flowcontrol.StreamFlowController
}

var (
	_ streamControlFrameGetter = &SendStream{}
	_ outgoingStream           = &SendStream{}
	_ sendStreamFrameHandler   = &SendStream{}
)

func newSendStream(
	ctx context.Context,
	streamID protocol.StreamID,
	sender streamSender,
	flowController flowcontrol.StreamFlowController,
	supportsResetStreamAt bool,
) *SendStream {
	s := &SendStream{
		streamID:              streamID,
		sender:                sender,
		flowController:        flowController,
		writeChan:             make(chan struct{}, 1),
		writeOnce:             make(chan struct{}, 1), // cap: 1, to protect against concurrent use of Write
		supportsResetStreamAt: supportsResetStreamAt,
	}
	s.ctx, s.ctxCancel = context.WithCancelCause(ctx)
	return s
}

// StreamID returns the stream ID.
func (s *SendStream) StreamID() StreamID {
	return s.streamID // same for receiveStream and sendStream
}

// Write writes data to the stream.
// Write can be made to time out using [SendStream.SetWriteDeadline].
// If the stream was canceled, the error is a [StreamError].
func (s *SendStream) Write(p []byte) (int, error) {
	// Concurrent use of Write is not permitted (and doesn't make any sense),
	// but sometimes people do it anyway.
	// Make sure that we only execute one call at any given time to avoid hard to debug failures.
	s.writeOnce <- struct{}{}
	defer func() { <-s.writeOnce }()

	isNewlyCompleted, n, err := s.write(p)
	if isNewlyCompleted {
		s.sender.onStreamCompleted(s.streamID)
	}
	return n, err
}

func (s *SendStream) write(p []byte) (bool /* is newly completed */, int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.resetErr != nil {
		s.cancellationFlagged = true
		return s.isNewlyCompleted(), 0, s.resetErr
	}
	if s.shutdownErr != nil {
		return false, 0, s.shutdownErr
	}
	if s.finishedWriting {
		return false, 0, fmt.Errorf("write on closed stream %d", s.streamID)
	}
	if !s.deadline.IsZero() && !monotime.Now().Before(s.deadline) {
		return false, 0, errDeadline
	}
	if len(p) == 0 {
		return false, 0, nil
	}

	s.dataForWriting = p

	var (
		deadlineTimer  *time.Timer
		bytesWritten   int
		notifiedSender bool
	)
	for {
		var copied bool
		var deadline monotime.Time
		// As soon as dataForWriting becomes smaller than a certain size x, we copy all the data to a STREAM frame (s.nextFrame),
		// which can then be popped the next time we assemble a packet.
		// This allows us to return Write() when all data but x bytes have been sent out.
		// When the user now calls Close(), this is much more likely to happen before we popped that last STREAM frame,
		// allowing us to set the FIN bit on that frame (instead of sending an empty STREAM frame with FIN).
		if s.canBufferStreamFrame() && len(s.dataForWriting) > 0 {
			if s.nextFrame == nil {
				f := wire.GetStreamFrame()
				f.Offset = s.writeOffset
				f.StreamID = s.streamID
				f.DataLenPresent = true
				f.Data = f.Data[:len(s.dataForWriting)]
				copy(f.Data, s.dataForWriting)
				s.nextFrame = f
			} else {
				l := len(s.nextFrame.Data)
				s.nextFrame.Data = s.nextFrame.Data[:l+len(s.dataForWriting)]
				copy(s.nextFrame.Data[l:], s.dataForWriting)
			}
			s.dataForWriting = nil
			bytesWritten = len(p)
			copied = true
		} else {
			bytesWritten = len(p) - len(s.dataForWriting)
			deadline = s.deadline
			if !deadline.IsZero() {
				if !monotime.Now().Before(deadline) {
					s.dataForWriting = nil
					return false, bytesWritten, errDeadline
				}
				if deadlineTimer == nil {
					deadlineTimer = time.NewTimer(monotime.Until(deadline))
					defer deadlineTimer.Stop()
				} else {
					deadlineTimer.Reset(monotime.Until(deadline))
				}
			}
			if s.dataForWriting == nil || s.shutdownErr != nil || s.resetErr != nil {
				break
			}
		}

		s.mutex.Unlock()
		if !notifiedSender {
			s.sender.onHasStreamData(s.streamID, s) // must be called without holding the mutex
			notifiedSender = true
		}
		if copied {
			s.mutex.Lock()
			break
		}
		if deadline.IsZero() {
			<-s.writeChan
		} else {
			select {
			case <-s.writeChan:
			case <-deadlineTimer.C:
			}
		}
		s.mutex.Lock()
	}

	if bytesWritten == len(p) {
		return false, bytesWritten, nil
	}
	if s.shutdownErr != nil {
		return false, bytesWritten, s.shutdownErr
	}
	if s.resetErr != nil {
		s.cancellationFlagged = true
		return s.isNewlyCompleted(), bytesWritten, s.resetErr
	}
	return false, bytesWritten, nil
}

func (s *SendStream) canBufferStreamFrame() bool {
	var l protocol.ByteCount
	if s.nextFrame != nil {
		l = s.nextFrame.DataLen()
	}
	return l+protocol.ByteCount(len(s.dataForWriting)) <= protocol.MaxPacketBufferSize
}

// popStreamFrame returns the next STREAM frame that is supposed to be sent on this stream
// maxBytes is the maximum length this frame (including frame header) will have.
func (s *SendStream) popStreamFrame(maxBytes protocol.ByteCount, v protocol.Version) (_ ackhandler.StreamFrame, _ *wire.StreamDataBlockedFrame, hasMore bool) {
	s.mutex.Lock()
	f, blocked, hasMoreData := s.popNewOrRetransmittedStreamFrame(maxBytes, v)
	if f != nil {
		s.numOutstandingFrames++
	}
	s.mutex.Unlock()

	if f == nil {
		return ackhandler.StreamFrame{}, blocked, hasMoreData
	}
	return ackhandler.StreamFrame{
		Frame:   f,
		Handler: (*sendStreamAckHandler)(s),
	}, blocked, hasMoreData
}

func (s *SendStream) popNewOrRetransmittedStreamFrame(maxBytes protocol.ByteCount, v protocol.Version) (_ *wire.StreamFrame, _ *wire.StreamDataBlockedFrame, hasMoreData bool) {
	if s.shutdownErr != nil {
		return nil, nil, false
	}
	if s.resetErr != nil {
		reliableOffset := s.reliableOffset()
		if reliableOffset == 0 || (s.writeOffset >= reliableOffset && len(s.retransmissionQueue) == 0) {
			return nil, nil, false
		}
	}

	if len(s.retransmissionQueue) > 0 {
		f, hasMoreRetransmissions := s.maybeGetRetransmission(maxBytes, v)
		if f != nil || hasMoreRetransmissions {
			if f == nil {
				return nil, nil, true
			}
			// We always claim that we have more data to send.
			// This might be incorrect, in which case there'll be a spurious call to popStreamFrame in the future.
			return f, nil, true
		}
	}

	if len(s.dataForWriting) == 0 && s.nextFrame == nil {
		if s.finishedWriting && !s.finSent {
			s.finSent = true
			return &wire.StreamFrame{
				StreamID:       s.streamID,
				Offset:         s.writeOffset,
				DataLenPresent: true,
				Fin:            true,
			}, nil, false
		}
		return nil, nil, false
	}

	maxDataLen := s.flowController.SendWindowSize()
	if maxDataLen == 0 {
		return nil, nil, true
	}

	// if the stream is canceled, only data up to the reliable size needs to be sent
	reliableOffset := s.reliableOffset()
	if s.resetErr != nil && reliableOffset > 0 {
		maxDataLen = min(maxDataLen, reliableOffset-s.writeOffset)
	}
	f, hasMoreData := s.popNewStreamFrame(maxBytes, maxDataLen, v)
	if f == nil {
		return nil, nil, hasMoreData
	}
	if f.DataLen() > 0 {
		s.writeOffset += f.DataLen()
		s.flowController.AddBytesSent(f.DataLen())
	}
	if s.resetErr != nil && s.writeOffset >= reliableOffset {
		hasMoreData = false
	}
	var blocked *wire.StreamDataBlockedFrame
	// If the entire send window is used, the stream might have become blocked on stream-level flow control.
	// This is not guaranteed though, because the stream might also have been blocked on connection-level flow control.
	if f.DataLen() == maxDataLen && s.flowController.IsNewlyBlocked() {
		blocked = &wire.StreamDataBlockedFrame{StreamID: s.streamID, MaximumStreamData: s.writeOffset}
	}
	f.Fin = s.finishedWriting && s.dataForWriting == nil && s.nextFrame == nil && !s.finSent
	if f.Fin {
		s.finSent = true
	}
	return f, blocked, hasMoreData
}

// popNewStreamFrame returns a new STREAM frame to send for this stream
// hasMoreData says if there's more data to send, *not* taking into account the reliable size
func (s *SendStream) popNewStreamFrame(maxBytes, maxDataLen protocol.ByteCount, v protocol.Version) (_ *wire.StreamFrame, hasMoreData bool) {
	if s.nextFrame != nil {
		maxDataLen := min(maxDataLen, s.nextFrame.MaxDataLen(maxBytes, v))
		if maxDataLen == 0 {
			return nil, true
		}
		nextFrame := s.nextFrame
		s.nextFrame = nil
		if nextFrame.DataLen() > maxDataLen {
			s.nextFrame = wire.GetStreamFrame()
			s.nextFrame.StreamID = s.streamID
			s.nextFrame.Offset = s.writeOffset + maxDataLen
			s.nextFrame.Data = s.nextFrame.Data[:nextFrame.DataLen()-maxDataLen]
			s.nextFrame.DataLenPresent = true
			copy(s.nextFrame.Data, nextFrame.Data[maxDataLen:])
			nextFrame.Data = nextFrame.Data[:maxDataLen]
		} else {
			s.signalWrite()
		}
		return nextFrame, s.nextFrame != nil || s.dataForWriting != nil
	}

	f := wire.GetStreamFrame()
	f.Fin = false
	f.StreamID = s.streamID
	f.Offset = s.writeOffset
	f.DataLenPresent = true
	f.Data = f.Data[:0]

	hasMoreData = s.popNewStreamFrameWithoutBuffer(f, maxBytes, maxDataLen, v)
	if len(f.Data) == 0 && !f.Fin {
		f.PutBack()
		return nil, hasMoreData
	}
	return f, hasMoreData
}

func (s *SendStream) popNewStreamFrameWithoutBuffer(f *wire.StreamFrame, maxBytes, sendWindow protocol.ByteCount, v protocol.Version) bool {
	maxDataLen := f.MaxDataLen(maxBytes, v)
	if maxDataLen == 0 { // a STREAM frame must have at least one byte of data
		return s.dataForWriting != nil || s.nextFrame != nil || s.finishedWriting
	}
	s.getDataForWriting(f, min(maxDataLen, sendWindow))

	return s.dataForWriting != nil || s.nextFrame != nil || s.finishedWriting
}

func (s *SendStream) maybeGetRetransmission(maxBytes protocol.ByteCount, v protocol.Version) (*wire.StreamFrame, bool /* has more retransmissions */) {
	f := s.retransmissionQueue[0]
	newFrame, needsSplit := f.MaybeSplitOffFrame(maxBytes, v)
	if needsSplit {
		return newFrame, true
	}
	s.retransmissionQueue = s.retransmissionQueue[1:]
	return f, len(s.retransmissionQueue) > 0
}

func (s *SendStream) getDataForWriting(f *wire.StreamFrame, maxBytes protocol.ByteCount) {
	if protocol.ByteCount(len(s.dataForWriting)) <= maxBytes {
		f.Data = f.Data[:len(s.dataForWriting)]
		copy(f.Data, s.dataForWriting)
		s.dataForWriting = nil
		s.signalWrite()
		return
	}
	f.Data = f.Data[:maxBytes]
	copy(f.Data, s.dataForWriting)
	s.dataForWriting = s.dataForWriting[maxBytes:]
	if s.canBufferStreamFrame() {
		s.signalWrite()
	}
}

func (s *SendStream) isNewlyCompleted() bool {
	if s.completed {
		return false
	}
	if s.nextFrame != nil && s.nextFrame.DataLen() > 0 {
		return false
	}
	// We need to keep the stream around until all frames have been sent and acknowledged.
	if s.numOutstandingFrames > 0 || len(s.retransmissionQueue) > 0 || s.queuedResetStreamFrame != nil {
		return false
	}
	// The stream is completed if we sent the FIN.
	if s.finSent {
		s.completed = true
		return true
	}
	// The stream is also completed if:
	// 1. the application called CancelWrite, or
	// 2. we received a STOP_SENDING, and
	// 		* the application consumed the error via Write, or
	//		* the application called Close
	if s.resetErr != nil && (s.cancellationFlagged || s.finishedWriting) {
		s.completed = true
		return true
	}
	return false
}

// Close closes the write-direction of the stream.
// Future calls to Write are not permitted after calling Close.
// It must not be called concurrently with Write.
// It must not be called after calling CancelWrite.
func (s *SendStream) Close() error {
	s.mutex.Lock()
	if s.shutdownErr != nil || s.finishedWriting {
		s.mutex.Unlock()
		return nil
	}
	s.finishedWriting = true
	cancelled := s.resetErr != nil
	if cancelled {
		s.cancellationFlagged = true
	}
	completed := s.isNewlyCompleted()
	s.mutex.Unlock()

	if completed {
		s.sender.onStreamCompleted(s.streamID)
	}
	if cancelled {
		return fmt.Errorf("close called for canceled stream %d", s.streamID)
	}
	s.sender.onHasStreamData(s.streamID, s) // need to send the FIN, must be called without holding the mutex

	s.ctxCancel(nil)
	return nil
}

// SetReliableBoundary marks the data written to this stream so far as reliable.
// It is valid to call this function multiple times, thereby increasing the reliable size.
// It only has an effect if the peer enabled support for the RESET_STREAM_AT extension,
// otherwise, it is a no-op.
func (s *SendStream) SetReliableBoundary() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.reliableSize = s.writeOffset
	if s.nextFrame != nil {
		s.reliableSize += s.nextFrame.DataLen()
	}
}

// returnFramesToPool returns all queued frames to the sync.Pool
func (s *SendStream) returnFramesToPool() {
	for _, f := range s.retransmissionQueue {
		f.PutBack()
	}
	clear(s.retransmissionQueue)
	s.retransmissionQueue = nil
	if s.nextFrame != nil {
		s.nextFrame.PutBack()
		s.nextFrame = nil
	}
}

// CancelWrite aborts sending on this stream.
// Data already written, but not yet delivered to the peer is not guaranteed to be delivered reliably.
// Write will unblock immediately, and future calls to Write will fail.
// When called multiple times it is a no-op.
// When called after Close, it aborts reliable delivery of outstanding stream data.
// Note that there is no guarantee if the peer will receive the FIN or the cancellation error first.
func (s *SendStream) CancelWrite(errorCode StreamErrorCode) {
	s.mutex.Lock()
	if s.shutdownErr != nil {
		s.mutex.Unlock()
		return
	}

	s.cancellationFlagged = true

	if s.resetErr != nil {
		completed := s.isNewlyCompleted()
		s.mutex.Unlock()
		// The user has called CancelWrite. If the previous cancellation was because of a
		// STOP_SENDING, we don't need to flag the error to the user anymore.
		if completed {
			s.sender.onStreamCompleted(s.streamID)
		}
		return
	}
	s.resetErr = &StreamError{StreamID: s.streamID, ErrorCode: errorCode, Remote: false}
	s.ctxCancel(s.resetErr)

	reliableOffset := s.reliableOffset()
	if reliableOffset == 0 {
		s.numOutstandingFrames = 0
		s.returnFramesToPool()
	}
	s.queuedResetStreamFrame = &wire.ResetStreamFrame{
		StreamID:  s.streamID,
		FinalSize: max(s.writeOffset, reliableOffset),
		ErrorCode: errorCode,
		// if the peer doesn't support the extension, the reliable offset will always be 0
		ReliableSize: reliableOffset,
	}
	if reliableOffset > 0 {
		if s.nextFrame != nil {
			if s.nextFrame.Offset >= reliableOffset {
				s.nextFrame.PutBack()
				s.nextFrame = nil
			} else if s.nextFrame.Offset+s.nextFrame.DataLen() > reliableOffset {
				s.nextFrame.Data = s.nextFrame.Data[:reliableOffset-s.nextFrame.Offset]
			}
		}
		if len(s.retransmissionQueue) > 0 {
			retransmissionQueue := make([]*wire.StreamFrame, 0, len(s.retransmissionQueue))
			for _, f := range s.retransmissionQueue {
				if f.Offset >= reliableOffset {
					f.PutBack()
					continue
				}
				if f.Offset+f.DataLen() <= reliableOffset {
					retransmissionQueue = append(retransmissionQueue, f)
				} else {
					f.Data = f.Data[:reliableOffset-f.Offset]
					retransmissionQueue = append(retransmissionQueue, f)
				}
			}
			s.retransmissionQueue = retransmissionQueue
		}
	}
	s.mutex.Unlock()

	s.signalWrite()
	s.sender.onHasStreamControlFrame(s.streamID, s)
}

func (s *SendStream) enableResetStreamAt() {
	s.mutex.Lock()
	s.supportsResetStreamAt = true
	s.mutex.Unlock()
}

func (s *SendStream) updateSendWindow(limit protocol.ByteCount) {
	updated := s.flowController.UpdateSendWindow(limit)
	if !updated { // duplicate or reordered MAX_STREAM_DATA frame
		return
	}
	s.mutex.Lock()
	hasStreamData := s.dataForWriting != nil || s.nextFrame != nil
	s.mutex.Unlock()
	if hasStreamData {
		s.sender.onHasStreamData(s.streamID, s)
	}
}

func (s *SendStream) handleStopSendingFrame(f *wire.StopSendingFrame) {
	s.mutex.Lock()
	if s.shutdownErr != nil {
		s.mutex.Unlock()
		return
	}

	// If the stream was already cancelled (either locally, or due to a previous STOP_SENDING frame),
	// there's nothing else to do.
	if s.resetErr != nil && s.reliableOffset() == 0 {
		s.mutex.Unlock()
		return
	}
	// if the peer stopped reading from the stream, there's no need to transmit any data reliably
	s.reliableSize = 0
	s.numOutstandingFrames = 0
	s.returnFramesToPool()
	if s.resetErr == nil {
		s.resetErr = &StreamError{StreamID: s.streamID, ErrorCode: f.ErrorCode, Remote: true}
		s.ctxCancel(s.resetErr)
	}
	s.queuedResetStreamFrame = &wire.ResetStreamFrame{
		StreamID:  s.streamID,
		FinalSize: s.writeOffset,
		ErrorCode: s.resetErr.ErrorCode,
	}
	s.mutex.Unlock()

	s.signalWrite()
	s.sender.onHasStreamControlFrame(s.streamID, s)
}

func (s *SendStream) getControlFrame(monotime.Time) (_ ackhandler.Frame, ok, hasMore bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.queuedResetStreamFrame == nil {
		return ackhandler.Frame{}, false, false
	}
	s.numOutstandingFrames++
	f := ackhandler.Frame{
		Frame:   s.queuedResetStreamFrame,
		Handler: (*sendStreamResetStreamHandler)(s),
	}
	s.queuedResetStreamFrame = nil
	return f, true, false
}

func (s *SendStream) reliableOffset() protocol.ByteCount {
	if !s.supportsResetStreamAt {
		return 0
	}
	return s.reliableSize
}

// The Context is canceled as soon as the write-side of the stream is closed.
// This happens when Close() or CancelWrite() is called, or when the peer
// cancels the read-side of their stream.
// The cancellation cause is set to the error that caused the stream to
// close, or `context.Canceled` in case the stream is closed without error.
func (s *SendStream) Context() context.Context {
	return s.ctx
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some data was successfully written.
// A zero value for t means Write will not time out.
func (s *SendStream) SetWriteDeadline(t time.Time) error {
	s.mutex.Lock()
	s.deadline = monotime.FromTime(t)
	s.mutex.Unlock()
	s.signalWrite()
	return nil
}

// CloseForShutdown closes a stream abruptly.
// It makes Write unblock (and return the error) immediately.
// The peer will NOT be informed about this: the stream is closed without sending a FIN or RST.
func (s *SendStream) closeForShutdown(err error) {
	s.mutex.Lock()
	if s.shutdownErr == nil && !s.finishedWriting {
		s.shutdownErr = err
		s.returnFramesToPool()
	}
	s.mutex.Unlock()
	s.signalWrite()
}

// signalWrite performs a non-blocking send on the writeChan
func (s *SendStream) signalWrite() {
	select {
	case s.writeChan <- struct{}{}:
	default:
	}
}

type sendStreamAckHandler SendStream

var _ ackhandler.FrameHandler = &sendStreamAckHandler{}

func (s *sendStreamAckHandler) OnAcked(f wire.Frame) {
	sf := f.(*wire.StreamFrame)
	sf.PutBack()

	s.mutex.Lock()
	if s.resetErr != nil && (*SendStream)(s).reliableOffset() == 0 {
		s.mutex.Unlock()
		return
	}
	s.numOutstandingFrames--
	if s.numOutstandingFrames < 0 {
		panic("numOutStandingFrames negative")
	}
	completed := (*SendStream)(s).isNewlyCompleted()
	s.mutex.Unlock()

	if completed {
		s.sender.onStreamCompleted(s.streamID)
	}
}

func (s *sendStreamAckHandler) OnLost(f wire.Frame) {
	sf := f.(*wire.StreamFrame)
	s.mutex.Lock()
	// If the reliable size was 0 when the stream was cancelled,
	// the number of outstanding frames was immediately set to 0, and the retransmission queue was dropped.
	if s.resetErr != nil && (*SendStream)(s).reliableOffset() == 0 {
		// Return the frame to pool since it won't be retransmitted
		sf.PutBack()
		s.mutex.Unlock()
		return
	}
	s.numOutstandingFrames--
	if s.numOutstandingFrames < 0 {
		panic("numOutStandingFrames negative")
	}

	if s.resetErr != nil && (*SendStream)(s).reliableOffset() > 0 {
		// If the stream was reset, and this frame is beyond the reliable offset,
		// it doesn't need to be retransmitted.
		if sf.Offset >= (*SendStream)(s).reliableOffset() {
			sf.PutBack()
			// If this frame was the last one tracked, losing it might cause the stream to be completed.
			completed := (*SendStream)(s).isNewlyCompleted()
			s.mutex.Unlock()
			if completed {
				s.sender.onStreamCompleted(s.streamID)
			}
			return
		}
		// If the payload of the frame extends beyond the reliable size,
		// truncate the frame to the reliable size.
		if sf.Offset+sf.DataLen() > (*SendStream)(s).reliableOffset() {
			sf.Data = sf.Data[:(*SendStream)(s).reliableOffset()-sf.Offset]
		}
	}

	sf.DataLenPresent = true
	s.retransmissionQueue = append(s.retransmissionQueue, sf)
	s.mutex.Unlock()

	s.sender.onHasStreamData(s.streamID, (*SendStream)(s))
}

type sendStreamResetStreamHandler SendStream

var _ ackhandler.FrameHandler = &sendStreamResetStreamHandler{}

func (s *sendStreamResetStreamHandler) OnAcked(f wire.Frame) {
	rsf := f.(*wire.ResetStreamFrame)
	s.mutex.Lock()
	// If the peer sent a STOP_SENDING after we sent a RESET_STREAM_AT frame,
	// we sent 1. reduced the reliable size to 0 and 2. sent a RESET_STREAM frame.
	// In this case, we don't care about the acknowledgment of this frame.
	if rsf.ReliableSize != (*SendStream)(s).reliableOffset() {
		s.mutex.Unlock()
		return
	}
	s.numOutstandingFrames--
	if s.numOutstandingFrames < 0 {
		panic("numOutStandingFrames negative")
	}
	completed := (*SendStream)(s).isNewlyCompleted()
	s.mutex.Unlock()

	if completed {
		s.sender.onStreamCompleted(s.streamID)
	}
}

func (s *sendStreamResetStreamHandler) OnLost(f wire.Frame) {
	rsf := f.(*wire.ResetStreamFrame)
	s.mutex.Lock()
	// If the peer sent a STOP_SENDING after we sent a RESET_STREAM_AT frame,
	// we sent 1. reduced the reliable size to 0 and 2. sent a RESET_STREAM frame.
	// In this case, the loss of the RESET_STREAM_AT frame can be ignored.
	if rsf.ReliableSize != (*SendStream)(s).reliableOffset() {
		s.mutex.Unlock()
		return
	}
	s.queuedResetStreamFrame = rsf
	s.numOutstandingFrames--
	s.mutex.Unlock()
	s.sender.onHasStreamControlFrame(s.streamID, (*SendStream)(s))
}
