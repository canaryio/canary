package canary

import "time"

type Sensor struct {
	u    string
	l    string
	C    chan *Sample
	quit chan int
}

func NewSensor(u, l string) *Sensor {
	return &Sensor{
		u:    u,
		l:    l,
		quit: make(chan int),
		C:    make(chan *Sample),
	}
}

func (s *Sensor) Start() {
	t := time.NewTicker(time.Second)
	sampler := NewSampler()

	for {
		select {
		case <-t.C:
			s.C <- sampler.Sample(s.u, s.l)
		case <-s.quit:
			break
		}
	}
}

func (s *Sensor) Stop() {
	close(s.quit)
}
