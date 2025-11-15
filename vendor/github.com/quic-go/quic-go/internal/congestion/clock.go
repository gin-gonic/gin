package congestion

import (
	"github.com/quic-go/quic-go/internal/monotime"
)

// A Clock returns the current time
type Clock interface {
	Now() monotime.Time
}

// DefaultClock implements the Clock interface using the Go stdlib clock.
type DefaultClock struct{}

var _ Clock = DefaultClock{}

// Now gets the current time
func (DefaultClock) Now() monotime.Time {
	return monotime.Now()
}
