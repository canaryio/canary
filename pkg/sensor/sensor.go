package sensor

import (
	"time"

	"github.com/canaryio/canary/pkg/sampler"
)

// Measurement reprents an aggregate of Target, Sample and error.
type Measurement struct {
	Target sampler.Target
	Sample sampler.Sample
	Error  error
}

// Sensor is capable of repeatedly measuring a given Target
// with a specific Sampler, and returns those results over channel C.
type Sensor struct {
	Target   sampler.Target
	C        chan Measurement
	Sampler  sampler.Sampler
	stopChan chan int
}

// take a sample against a target.
func (s *Sensor) measure() Measurement {
	sample, err := s.Sampler.Sample(s.Target)
	return Measurement{
		Target: s.Target,
		Sample: sample,
		Error:  err,
	}
}

// Start is meant to be called within a goroutine, and fires up the main event loop.
func (s *Sensor) Start(interval int) {
	if s.stopChan == nil {
		s.stopChan = make(chan int)
	}
	t := time.NewTicker((time.Second * time.Duration(interval)))

	for {
		select {
		case <-s.stopChan:
			return
		case <-t.C:
			s.C <- s.measure()
		}
	}
}

// Stop halts the event loop.
func (s *Sensor) Stop() {
	close(s.stopChan)
}
