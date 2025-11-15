// Package monotime provides a monotonic time representation that is useful for
// measuring elapsed time.
// It is designed as a memory optimized drop-in replacement for time.Time, with
// a monotime.Time consuming just 8 bytes instead of 24 bytes.
package monotime

import (
	"time"
)

// The absolute value doesn't matter, but it should be in the past,
// so that every timestamp obtained with Now() is non-zero,
// even on systems with low timer resolutions (e.g. Windows).
var start = time.Now().Add(-time.Hour)

// A Time represents an instant in monotonic time.
// Times can be compared using the comparison operators, but the specific
// value is implementation-dependent and should not be relied upon.
// The zero value of Time doesn't have any specific meaning.
type Time int64

// Now returns the current monotonic time.
func Now() Time {
	return Time(time.Since(start).Nanoseconds())
}

// Sub returns the duration t-t2. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t Time) Sub(t2 Time) time.Duration {
	return time.Duration(t - t2)
}

// Add returns the time t+d.
func (t Time) Add(d time.Duration) Time {
	return Time(int64(t) + d.Nanoseconds())
}

// After reports whether the time instant t is after t2.
func (t Time) After(t2 Time) bool {
	return t > t2
}

// Before reports whether the time instant t is before t2.
func (t Time) Before(t2 Time) bool {
	return t < t2
}

// IsZero reports whether t represents the zero time instant.
func (t Time) IsZero() bool {
	return t == 0
}

// Equal reports whether t and t2 represent the same time instant.
func (t Time) Equal(t2 Time) bool {
	return t == t2
}

// ToTime converts the monotonic time to a time.Time value.
// The returned time.Time will have the same instant as the monotonic time,
// but may be subject to clock adjustments.
func (t Time) ToTime() time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	return start.Add(time.Duration(t))
}

// Since returns the time elapsed since t. It is shorthand for Now().Sub(t).
func Since(t Time) time.Duration {
	return Now().Sub(t)
}

// Until returns the duration until t.
// It is shorthand for t.Sub(Now()).
// If t is in the past, the returned duration will be negative.
func Until(t Time) time.Duration {
	return time.Duration(t - Now())
}

// FromTime converts a time.Time to a monotonic Time.
// The conversion is relative to the package's start time and may lose
// precision if the time.Time is far from the start time.
func FromTime(t time.Time) Time {
	if t.IsZero() {
		return 0
	}
	return Time(t.Sub(start).Nanoseconds())
}
