package canary

import "time"

// Measurement reprents an aggregate of Target, Sample and error.
type Measurement struct {
	Target Target
	Sample Sample
	Error  error
}

// Scheduler is capable of repeatedly measuring a given Target
// with a specific Sampler, and returns those results over channel C.
type Scheduler struct {
	Target   Target
	C        chan Measurement
	Sampler  Sampler
	stopChan chan int
}

// take a sample against a target.
func (s *Scheduler) measure() Measurement {
	sample, err := s.Sampler.Sample(s.Target)
	return Measurement{
		Target: s.Target,
		Sample: sample,
		Error:  err,
	}
}

// Start is meant to be called within a goroutine, and fires up the main event loop.
func (s *Scheduler) Start() {
	if s.stopChan == nil {
		s.stopChan = make(chan int)
	}
	t := time.NewTicker(time.Second)

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
func (s *Scheduler) Stop() {
	close(s.stopChan)
}
