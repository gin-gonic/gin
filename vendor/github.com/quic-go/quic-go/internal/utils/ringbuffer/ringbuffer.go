package ringbuffer

// A RingBuffer is a ring buffer.
// It acts as a heap that doesn't cause any allocations.
type RingBuffer[T any] struct {
	ring             []T
	headPos, tailPos int
	full             bool
}

// Init preallocates a buffer with a certain size.
func (r *RingBuffer[T]) Init(size int) {
	r.ring = make([]T, size)
}

// Len returns the number of elements in the ring buffer.
func (r *RingBuffer[T]) Len() int {
	if r.full {
		return len(r.ring)
	}
	if r.tailPos >= r.headPos {
		return r.tailPos - r.headPos
	}
	return r.tailPos - r.headPos + len(r.ring)
}

// Empty says if the ring buffer is empty.
func (r *RingBuffer[T]) Empty() bool {
	return !r.full && r.headPos == r.tailPos
}

// PushBack adds a new element.
// If the ring buffer is full, its capacity is increased first.
func (r *RingBuffer[T]) PushBack(t T) {
	if r.full || len(r.ring) == 0 {
		r.grow()
	}
	r.ring[r.tailPos] = t
	r.tailPos++
	if r.tailPos == len(r.ring) {
		r.tailPos = 0
	}
	if r.tailPos == r.headPos {
		r.full = true
	}
}

// PopFront returns the next element.
// It must not be called when the buffer is empty, that means that
// callers might need to check if there are elements in the buffer first.
func (r *RingBuffer[T]) PopFront() T {
	if r.Empty() {
		panic("github.com/quic-go/quic-go/internal/utils/ringbuffer: pop from an empty queue")
	}
	r.full = false
	t := r.ring[r.headPos]
	r.ring[r.headPos] = *new(T)
	r.headPos++
	if r.headPos == len(r.ring) {
		r.headPos = 0
	}
	return t
}

// PeekFront returns the next element.
// It must not be called when the buffer is empty, that means that
// callers might need to check if there are elements in the buffer first.
func (r *RingBuffer[T]) PeekFront() T {
	if r.Empty() {
		panic("github.com/quic-go/quic-go/internal/utils/ringbuffer: peek from an empty queue")
	}
	return r.ring[r.headPos]
}

// Grow the maximum size of the queue.
// This method assume the queue is full.
func (r *RingBuffer[T]) grow() {
	oldRing := r.ring
	newSize := len(oldRing) * 2
	if newSize == 0 {
		newSize = 1
	}
	r.ring = make([]T, newSize)
	headLen := copy(r.ring, oldRing[r.headPos:])
	copy(r.ring[headLen:], oldRing[:r.headPos])
	r.headPos, r.tailPos, r.full = 0, len(oldRing), false
}

// Clear removes all elements.
func (r *RingBuffer[T]) Clear() {
	var zeroValue T
	for i := range r.ring {
		r.ring[i] = zeroValue
	}
	r.headPos, r.tailPos, r.full = 0, 0, false
}
