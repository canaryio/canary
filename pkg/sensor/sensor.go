package sensor

import (
	"time"

	"github.com/canaryio/canary/pkg/sampler"
)

// Measurement reprents an aggregate of Target, Sample and error.
type Measurement struct {
	Target     sampler.Target
	Sample     sampler.Sample
	IsOK       bool
	StateCount int
	Error      error
}

// Sensor is capable of repeatedly measuring a given Target
// with a specific Sampler, and returns those results over channel C.
type Sensor struct {
	Target         sampler.Target
	C              chan Measurement
	StateCounter   int
	StopChan       chan int
	StopNotifyChan chan bool
	IsOK           bool
	Timeout        int // timeout in secs
}

// take a sample against a target.
func (s *Sensor) measure() Measurement {
	sample, err := sampler.Ping(s.Target, s.Timeout)
	m := Measurement{
		Target: s.Target,
		Sample: sample,
		Error:  err,
	}

	// Record the pass/fail for this measurement
	m.IsOK = (m.Error == nil)

	// Update the Sensors value for IsOK and counter for said state.
	if s.IsOK != m.IsOK {
		s.IsOK = m.IsOK
		s.StateCounter = 0
	}
	s.StateCounter++
	m.StateCount = s.StateCounter

	return m
}

// Start is meant to be called within a goroutine, and fires up the main event loop.
// interval is number of seconds. delay is number of ms.
func (s *Sensor) Start(delay float64) {
	// Delay for loop start offset.
	time.Sleep((time.Millisecond * time.Duration(delay)))

	// Start the ticker for this sensors interval
	t := time.NewTicker((time.Second * time.Duration(s.Target.Interval)))

	// Measure, then wait for ticker interval
	s.C <- s.measure()

	isStopped := false
	for ! isStopped {
		<-t.C
		select {
		case <-s.StopChan:
			isStopped = true
			s.StopNotifyChan <- true
			return
		default:
			s.C <- s.measure()
		}
	}
}

// Stop halts the event loop.
func (s *Sensor) Stop() {
	s.StopChan <- 1
}
