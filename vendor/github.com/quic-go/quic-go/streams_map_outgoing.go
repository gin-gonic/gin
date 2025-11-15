package quic

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/wire"
)

type outgoingStream interface {
	updateSendWindow(protocol.ByteCount)
	enableResetStreamAt()
	closeForShutdown(error)
}

type outgoingStreamsMap[T outgoingStream] struct {
	mutex sync.RWMutex

	streamType protocol.StreamType
	streams    map[protocol.StreamID]T

	openQueue []chan struct{}

	nextStream  protocol.StreamID // stream ID of the stream returned by OpenStream(Sync)
	maxStream   protocol.StreamID // the maximum stream ID we're allowed to open
	blockedSent bool              // was a STREAMS_BLOCKED sent for the current maxStream

	newStream            func(protocol.StreamID) T
	queueStreamIDBlocked func(*wire.StreamsBlockedFrame)

	closeErr error
}

func newOutgoingStreamsMap[T outgoingStream](
	streamType protocol.StreamType,
	newStream func(protocol.StreamID) T,
	queueControlFrame func(wire.Frame),
	pers protocol.Perspective,
) *outgoingStreamsMap[T] {
	var nextStream protocol.StreamID
	switch {
	case streamType == protocol.StreamTypeBidi && pers == protocol.PerspectiveServer:
		nextStream = protocol.FirstOutgoingBidiStreamServer
	case streamType == protocol.StreamTypeBidi && pers == protocol.PerspectiveClient:
		nextStream = protocol.FirstOutgoingBidiStreamClient
	case streamType == protocol.StreamTypeUni && pers == protocol.PerspectiveServer:
		nextStream = protocol.FirstOutgoingUniStreamServer
	case streamType == protocol.StreamTypeUni && pers == protocol.PerspectiveClient:
		nextStream = protocol.FirstOutgoingUniStreamClient
	}
	return &outgoingStreamsMap[T]{
		streamType:           streamType,
		streams:              make(map[protocol.StreamID]T),
		maxStream:            protocol.InvalidStreamNum,
		nextStream:           nextStream,
		newStream:            newStream,
		queueStreamIDBlocked: func(f *wire.StreamsBlockedFrame) { queueControlFrame(f) },
	}
}

func (m *outgoingStreamsMap[T]) OpenStream() (T, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.closeErr != nil {
		return *new(T), m.closeErr
	}

	// if there are OpenStreamSync calls waiting, return an error here
	if len(m.openQueue) > 0 || m.nextStream > m.maxStream {
		m.maybeSendBlockedFrame()
		return *new(T), &StreamLimitReachedError{}
	}
	return m.openStream(), nil
}

func (m *outgoingStreamsMap[T]) OpenStreamSync(ctx context.Context) (T, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.closeErr != nil {
		return *new(T), m.closeErr
	}
	if err := ctx.Err(); err != nil {
		return *new(T), err
	}
	if len(m.openQueue) == 0 && m.nextStream <= m.maxStream {
		return m.openStream(), nil
	}

	waitChan := make(chan struct{}, 1)
	m.openQueue = append(m.openQueue, waitChan)
	m.maybeSendBlockedFrame()

	for {
		m.mutex.Unlock()
		select {
		case <-ctx.Done():
			m.mutex.Lock()
			m.openQueue = slices.DeleteFunc(m.openQueue, func(c chan struct{}) bool {
				return c == waitChan
			})
			// If we just received a MAX_STREAMS frame, this might have been the next stream
			// that could be opened. Make sure we unblock the next OpenStreamSync call.
			m.maybeUnblockOpenSync()
			return *new(T), ctx.Err()
		case <-waitChan:
		}

		m.mutex.Lock()
		if m.closeErr != nil {
			return *new(T), m.closeErr
		}
		if m.nextStream > m.maxStream {
			// no stream available. Continue waiting
			continue
		}
		str := m.openStream()
		m.openQueue = m.openQueue[1:]
		m.maybeUnblockOpenSync()
		return str, nil
	}
}

func (m *outgoingStreamsMap[T]) openStream() T {
	s := m.newStream(m.nextStream)
	m.streams[m.nextStream] = s
	m.nextStream += 4
	return s
}

// maybeSendBlockedFrame queues a STREAMS_BLOCKED frame for the current stream offset,
// if we haven't sent one for this offset yet
func (m *outgoingStreamsMap[T]) maybeSendBlockedFrame() {
	if m.blockedSent {
		return
	}

	var streamLimit protocol.StreamNum
	if m.maxStream != protocol.InvalidStreamID {
		streamLimit = m.maxStream.StreamNum()
	}
	m.queueStreamIDBlocked(&wire.StreamsBlockedFrame{
		Type:        m.streamType,
		StreamLimit: streamLimit,
	})
	m.blockedSent = true
}

func (m *outgoingStreamsMap[T]) GetStream(id protocol.StreamID) (T, error) {
	m.mutex.RLock()
	if id >= m.nextStream {
		m.mutex.RUnlock()
		return *new(T), &qerr.TransportError{
			ErrorCode:    qerr.StreamStateError,
			ErrorMessage: fmt.Sprintf("peer attempted to open stream %d", id),
		}
	}
	s := m.streams[id]
	m.mutex.RUnlock()
	return s, nil
}

func (m *outgoingStreamsMap[T]) DeleteStream(id protocol.StreamID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.streams[id]; !ok {
		return &qerr.TransportError{
			ErrorCode:    qerr.StreamStateError,
			ErrorMessage: fmt.Sprintf("tried to delete unknown outgoing stream %d", id),
		}
	}
	delete(m.streams, id)
	return nil
}

func (m *outgoingStreamsMap[T]) SetMaxStream(id protocol.StreamID) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if id <= m.maxStream {
		return
	}
	m.maxStream = id
	m.blockedSent = false
	if m.maxStream < m.nextStream-4+4*protocol.StreamID(len(m.openQueue)) {
		m.maybeSendBlockedFrame()
	}
	m.maybeUnblockOpenSync()
}

// UpdateSendWindow is called when the peer's transport parameters are received.
// Only in the case of a 0-RTT handshake will we have open streams at this point.
// We might need to update the send window, in case the server increased it.
func (m *outgoingStreamsMap[T]) UpdateSendWindow(limit protocol.ByteCount) {
	m.mutex.Lock()
	for _, str := range m.streams {
		str.updateSendWindow(limit)
	}
	m.mutex.Unlock()
}

func (m *outgoingStreamsMap[T]) EnableResetStreamAt() {
	m.mutex.Lock()
	for _, str := range m.streams {
		str.enableResetStreamAt()
	}
	m.mutex.Unlock()
}

// unblockOpenSync unblocks the next OpenStreamSync go-routine to open a new stream
func (m *outgoingStreamsMap[T]) maybeUnblockOpenSync() {
	if len(m.openQueue) == 0 {
		return
	}
	if m.nextStream > m.maxStream {
		return
	}
	// unblockOpenSync is called both from OpenStreamSync and from SetMaxStream.
	// It's sufficient to only unblock OpenStreamSync once.
	select {
	case m.openQueue[0] <- struct{}{}:
	default:
	}
}

func (m *outgoingStreamsMap[T]) CloseWithError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.closeErr = err
	for _, str := range m.streams {
		str.closeForShutdown(err)
	}
	for _, c := range m.openQueue {
		if c != nil {
			close(c)
		}
	}
	m.openQueue = nil
}
