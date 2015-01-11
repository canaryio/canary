package canary

import (
	"fmt"
	"time"
)

// StatusCodeError is an error representing an HTTP Status code
// of 400 or greater.
type StatusCodeError struct {
	StatusCode int
}

func (e StatusCodeError) Error() string {
	return fmt.Sprintf(
		"recieved HTTP status %d",
		e.StatusCode,
	)
}

// Target represents the things that we are measureing.
type Target struct {
	URL  string
	Name string
}

// Sample represents HTTP state from a given point in time.
type Sample struct {
	StatusCode int
	T1         time.Time
	T2         time.Time
}

// Latency returns the amount of milliseconds between T1
// and T2 (start and finish).
func (s Sample) Latency() float64 {
	return s.T2.Sub(s.T1).Seconds() * 1000
}

// A Sampler is an interface that provides the Sample method.
type Sampler interface {
	Sample(Target) (Sample, error)
}
