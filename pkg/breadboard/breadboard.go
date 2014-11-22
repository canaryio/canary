package breadboard

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/libratoreporter"
)

// Breadboards wire everything together
type Breadboard struct {
	sync.Mutex
	sensors   map[string]*canary.Sensor
	reporters map[string]canary.Reporter
	c         chan *canary.Sample
	stop      chan int
	source    string
}

func New(source string) *Breadboard {
	return &Breadboard{
		sensors:   make(map[string]*canary.Sensor),
		reporters: make(map[string]canary.Reporter),
		c:         make(chan *canary.Sample),
		source:    source,
		stop:      make(chan int),
	}
}

func (b *Breadboard) Start() {
	for {
		select {
		case <-b.stop:
			for _, s := range b.sensors {
				s.Stop()
			}

			for _, r := range b.reporters {
				r.Stop()
			}
		case sample := <-b.c:
			for _, r := range b.reporters {
				r.Ingest(sample)
			}
		}
	}
}

func (b *Breadboard) Stop() {
	close(b.stop)
}

func (b *Breadboard) Update(m *canary.Manifest) error {
	b.Lock()
	defer b.Unlock()

	// launch new sensors if they are not already running
	for _, site := range m.Sites {
		sensor := canary.NewSensor(site, b.source)
		sensor.C = b.c
		b.sensors[site.URL] = sensor
		go sensor.Start()
	}

	// TODO: teardown old sensors

	// launch new reporters if they are not already running
	for _, service := range m.Services {
		switch service.Name {
		case "librato":
			config := &libratoreporter.Config{}
			err := json.Unmarshal(service.Config, config)
			if err != nil {
				return err
			}
			r := libratoreporter.New(config)
			b.reporters[service.Name] = r
			go r.Start()
		default:
			log.Printf("unknown service: %s", service.Name)
		}
	}

	// TODO: teardown old reporters

	return nil
}
