package canary

import "time"

type Sensor struct {
	site *Site
	l    string
	C    chan *Sample
	quit chan int
}

func NewSensor(site *Site, l string) *Sensor {
	return &Sensor{
		site: site,
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
			s.C <- sampler.Sample(s.site, s.l)
		case <-s.quit:
			break
		}
	}
}

func (s *Sensor) Stop() {
	close(s.quit)
}
