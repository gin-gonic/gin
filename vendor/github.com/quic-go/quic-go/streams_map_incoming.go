package quic

import (
	"context"
	"fmt"
	"sync"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/wire"
)

type incomingStream interface {
	closeForShutdown(error)
}

// When a stream is deleted before it was accepted, we can't delete it from the map immediately.
// We need to wait until the application accepts it, and delete it then.
type incomingStreamEntry[T incomingStream] struct {
	stream       T
	shouldDelete bool
}

type incomingStreamsMap[T incomingStream] struct {
	mutex         sync.RWMutex
	newStreamChan chan struct{}

	streamType protocol.StreamType
	streams    map[protocol.StreamID]incomingStreamEntry[T]

	nextStreamToAccept protocol.StreamID // the next stream that will be returned by AcceptStream()
	nextStreamToOpen   protocol.StreamID // the highest stream that the peer opened
	maxStream          protocol.StreamID // the highest stream that the peer is allowed to open
	maxNumStreams      uint64            // maximum number of streams

	newStream        func(protocol.StreamID) T
	queueMaxStreamID func(*wire.MaxStreamsFrame)

	closeErr error
}

func newIncomingStreamsMap[T incomingStream](
	streamType protocol.StreamType,
	newStream func(protocol.StreamID) T,
	maxStreams uint64,
	queueControlFrame func(wire.Frame),
	pers protocol.Perspective,
) *incomingStreamsMap[T] {
	var nextStreamToAccept protocol.StreamID
	switch {
	case streamType == protocol.StreamTypeBidi && pers == protocol.PerspectiveServer:
		nextStreamToAccept = protocol.FirstIncomingBidiStreamServer
	case streamType == protocol.StreamTypeBidi && pers == protocol.PerspectiveClient:
		nextStreamToAccept = protocol.FirstIncomingBidiStreamClient
	case streamType == protocol.StreamTypeUni && pers == protocol.PerspectiveServer:
		nextStreamToAccept = protocol.FirstIncomingUniStreamServer
	case streamType == protocol.StreamTypeUni && pers == protocol.PerspectiveClient:
		nextStreamToAccept = protocol.FirstIncomingUniStreamClient
	}
	return &incomingStreamsMap[T]{
		newStreamChan:      make(chan struct{}, 1),
		streamType:         streamType,
		streams:            make(map[protocol.StreamID]incomingStreamEntry[T]),
		maxStream:          protocol.StreamNum(maxStreams).StreamID(streamType, pers.Opposite()),
		maxNumStreams:      maxStreams,
		newStream:          newStream,
		nextStreamToOpen:   nextStreamToAccept,
		nextStreamToAccept: nextStreamToAccept,
		queueMaxStreamID:   func(f *wire.MaxStreamsFrame) { queueControlFrame(f) },
	}
}

func (m *incomingStreamsMap[T]) AcceptStream(ctx context.Context) (T, error) {
	// drain the newStreamChan, so we don't check the map twice if the stream doesn't exist
	select {
	case <-m.newStreamChan:
	default:
	}

	m.mutex.Lock()

	var id protocol.StreamID
	var entry incomingStreamEntry[T]
	for {
		id = m.nextStreamToAccept
		if m.closeErr != nil {
			m.mutex.Unlock()
			return *new(T), m.closeErr
		}
		var ok bool
		entry, ok = m.streams[id]
		if ok {
			break
		}
		m.mutex.Unlock()
		select {
		case <-ctx.Done():
			return *new(T), ctx.Err()
		case <-m.newStreamChan:
		}
		m.mutex.Lock()
	}
	m.nextStreamToAccept += 4
	// If this stream was completed before being accepted, we can delete it now.
	if entry.shouldDelete {
		if err := m.deleteStream(id); err != nil {
			m.mutex.Unlock()
			return *new(T), err
		}
	}
	m.mutex.Unlock()
	return entry.stream, nil
}

func (m *incomingStreamsMap[T]) GetOrOpenStream(id protocol.StreamID) (T, error) {
	m.mutex.RLock()
	if id > m.maxStream {
		m.mutex.RUnlock()
		return *new(T), &qerr.TransportError{
			ErrorCode:    qerr.StreamLimitError,
			ErrorMessage: fmt.Sprintf("peer tried to open stream %d (current limit: %d)", id, m.maxStream),
		}
	}
	// if the num is smaller than the highest we accepted
	// * this stream exists in the map, and we can return it, or
	// * this stream was already closed, then we can return the nil
	if id < m.nextStreamToOpen {
		var s T
		// If the stream was already queued for deletion, and is just waiting to be accepted, don't return it.
		if entry, ok := m.streams[id]; ok && !entry.shouldDelete {
			s = entry.stream
		}
		m.mutex.RUnlock()
		return s, nil
	}
	m.mutex.RUnlock()

	m.mutex.Lock()
	// no need to check the two error conditions from above again
	// * maxStream can only increase, so if the id was valid before, it definitely is valid now
	// * highestStream is only modified by this function
	for newNum := m.nextStreamToOpen; newNum <= id; newNum += 4 {
		m.streams[newNum] = incomingStreamEntry[T]{stream: m.newStream(newNum)}
		select {
		case m.newStreamChan <- struct{}{}:
		default:
		}
	}
	m.nextStreamToOpen = id + 4
	entry := m.streams[id]
	m.mutex.Unlock()
	return entry.stream, nil
}

func (m *incomingStreamsMap[T]) DeleteStream(id protocol.StreamID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.deleteStream(id); err != nil {
		return &qerr.TransportError{
			ErrorCode:    qerr.StreamStateError,
			ErrorMessage: err.Error(),
		}
	}
	return nil
}

func (m *incomingStreamsMap[T]) deleteStream(id protocol.StreamID) error {
	if _, ok := m.streams[id]; !ok {
		return fmt.Errorf("tried to delete unknown incoming stream %d", id)
	}

	// Don't delete this stream yet, if it was not yet accepted.
	// Just save it to streamsToDelete map, to make sure it is deleted as soon as it gets accepted.
	if id >= m.nextStreamToAccept {
		entry, ok := m.streams[id]
		if ok && entry.shouldDelete {
			return fmt.Errorf("tried to delete incoming stream %d multiple times", id)
		}
		entry.shouldDelete = true
		m.streams[id] = entry // can't assign to struct in map, so we need to reassign
		return nil
	}

	delete(m.streams, id)
	// queue a MAX_STREAM_ID frame, giving the peer the option to open a new stream
	if m.maxNumStreams > uint64(len(m.streams)) {
		maxStream := m.nextStreamToOpen + 4*protocol.StreamID(m.maxNumStreams-uint64(len(m.streams))-1)
		// never send a value larger than the maximum value for a stream number
		if maxStream <= protocol.MaxStreamID {
			m.maxStream = maxStream
			m.queueMaxStreamID(&wire.MaxStreamsFrame{
				Type:         m.streamType,
				MaxStreamNum: m.maxStream.StreamNum(),
			})
		}
	}
	return nil
}

func (m *incomingStreamsMap[T]) CloseWithError(err error) {
	m.mutex.Lock()
	m.closeErr = err
	for _, entry := range m.streams {
		entry.stream.closeForShutdown(err)
	}
	m.mutex.Unlock()
	close(m.newStreamChan)
}
