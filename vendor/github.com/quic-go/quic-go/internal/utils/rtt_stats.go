package utils

import (
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
)

const (
	rttAlpha      = 0.125
	oneMinusAlpha = 1 - rttAlpha
	rttBeta       = 0.25
	oneMinusBeta  = 1 - rttBeta
)

// The default RTT used before an RTT sample is taken
const DefaultInitialRTT = 100 * time.Millisecond

// RTTStats provides round-trip statistics
type RTTStats struct {
	hasMeasurement bool

	minRTT        atomic.Int64 // nanoseconds
	latestRTT     atomic.Int64 // nanoseconds
	smoothedRTT   atomic.Int64 // nanoseconds
	meanDeviation atomic.Int64 // nanoseconds

	maxAckDelay atomic.Int64 // nanoseconds
}

func NewRTTStats() *RTTStats {
	var rttStats RTTStats
	rttStats.minRTT.Store(DefaultInitialRTT.Nanoseconds())
	rttStats.latestRTT.Store(DefaultInitialRTT.Nanoseconds())
	rttStats.smoothedRTT.Store(DefaultInitialRTT.Nanoseconds())
	return &rttStats
}

// MinRTT Returns the minRTT for the entire connection.
// May return Zero if no valid updates have occurred.
func (r *RTTStats) MinRTT() time.Duration {
	return time.Duration(r.minRTT.Load())
}

// LatestRTT returns the most recent rtt measurement.
// May return Zero if no valid updates have occurred.
func (r *RTTStats) LatestRTT() time.Duration {
	return time.Duration(r.latestRTT.Load())
}

// SmoothedRTT returns the smoothed RTT for the connection.
// May return Zero if no valid updates have occurred.
func (r *RTTStats) SmoothedRTT() time.Duration {
	return time.Duration(r.smoothedRTT.Load())
}

// MeanDeviation gets the mean deviation
func (r *RTTStats) MeanDeviation() time.Duration {
	return time.Duration(r.meanDeviation.Load())
}

// MaxAckDelay gets the max_ack_delay advertised by the peer
func (r *RTTStats) MaxAckDelay() time.Duration {
	return time.Duration(r.maxAckDelay.Load())
}

// PTO gets the probe timeout duration.
func (r *RTTStats) PTO(includeMaxAckDelay bool) time.Duration {
	if !r.hasMeasurement {
		return 2 * DefaultInitialRTT
	}
	pto := r.SmoothedRTT() + max(4*r.MeanDeviation(), protocol.TimerGranularity)
	if includeMaxAckDelay {
		pto += r.MaxAckDelay()
	}
	return pto
}

// UpdateRTT updates the RTT based on a new sample.
func (r *RTTStats) UpdateRTT(sendDelta, ackDelay time.Duration) {
	if sendDelta <= 0 {
		return
	}

	// Update r.minRTT first. r.minRTT does not use an rttSample corrected for
	// ackDelay but the raw observed sendDelta, since poor clock granularity at
	// the client may cause a high ackDelay to result in underestimation of the
	// r.minRTT.
	minRTT := time.Duration(r.minRTT.Load())
	if !r.hasMeasurement || minRTT > sendDelta {
		minRTT = sendDelta
		r.minRTT.Store(sendDelta.Nanoseconds())
	}

	// Correct for ackDelay if information received from the peer results in a
	// an RTT sample at least as large as minRTT. Otherwise, only use the
	// sendDelta.
	sample := sendDelta
	if sample-minRTT >= ackDelay {
		sample -= ackDelay
	}
	r.latestRTT.Store(sample.Nanoseconds())
	// First time call.
	if !r.hasMeasurement {
		r.hasMeasurement = true
		r.smoothedRTT.Store(sample.Nanoseconds())
		r.meanDeviation.Store(sample.Nanoseconds() / 2)
	} else {
		smoothedRTT := r.SmoothedRTT()
		meanDev := time.Duration(oneMinusBeta*float32(r.MeanDeviation()/time.Microsecond)+rttBeta*float32((smoothedRTT-sample).Abs()/time.Microsecond)) * time.Microsecond
		newSmoothedRTT := time.Duration((float32(smoothedRTT/time.Microsecond)*oneMinusAlpha)+(float32(sample/time.Microsecond)*rttAlpha)) * time.Microsecond
		r.meanDeviation.Store(meanDev.Nanoseconds())
		r.smoothedRTT.Store(newSmoothedRTT.Nanoseconds())
	}
}

func (r *RTTStats) HasMeasurement() bool {
	return r.hasMeasurement
}

// SetMaxAckDelay sets the max_ack_delay
func (r *RTTStats) SetMaxAckDelay(mad time.Duration) {
	r.maxAckDelay.Store(int64(mad))
}

// SetInitialRTT sets the initial RTT.
// It is used during handshake when restoring the RTT stats from the token.
func (r *RTTStats) SetInitialRTT(t time.Duration) {
	// On the server side, by the time we get to process the session ticket,
	// we might already have obtained an RTT measurement.
	// This can happen if we received the ClientHello in multiple pieces, and one of those pieces was lost.
	// Discard the restored value. A fresh measurement is always better.
	if r.hasMeasurement {
		return
	}
	r.smoothedRTT.Store(int64(t))
	r.latestRTT.Store(int64(t))
}

func (r *RTTStats) ResetForPathMigration() {
	r.hasMeasurement = false
	r.minRTT.Store(DefaultInitialRTT.Nanoseconds())
	r.latestRTT.Store(DefaultInitialRTT.Nanoseconds())
	r.smoothedRTT.Store(DefaultInitialRTT.Nanoseconds())
	r.meanDeviation.Store(0)
	// max_ack_delay remains valid
}

func (r *RTTStats) Clone() *RTTStats {
	out := &RTTStats{}
	out.hasMeasurement = r.hasMeasurement
	out.minRTT.Store(r.minRTT.Load())
	out.latestRTT.Store(r.latestRTT.Load())
	out.smoothedRTT.Store(r.smoothedRTT.Load())
	out.meanDeviation.Store(r.meanDeviation.Load())
	out.maxAckDelay.Store(r.maxAckDelay.Load())
	return out
}
