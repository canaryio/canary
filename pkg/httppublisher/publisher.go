package httppublisher

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/canaryio/canary/pkg/sensor"
)

type Publisher struct {
	url      string
	interval time.Duration
	buffer   []*sensor.Measurement
	c        chan *sensor.Measurement
}

func New(url string, interval time.Duration) *Publisher {
	p := &Publisher{
		buffer:   make([]*sensor.Measurement, 0),
		c:        make(chan *sensor.Measurement),
		url:      url,
		interval: interval,
	}
	go p.run()
	return p
}

func (p *Publisher) Publish(m sensor.Measurement) error {
	p.c <- &m
	return nil
}

func (p *Publisher) run() {
	t := time.NewTicker(p.interval)
	for {
		select {
		case <-t.C:
			payload, err := json.Marshal(p.buffer)
			if err != nil {
				log.Printf("httppublisher err='%s'", err)
			}

			resp, err := http.Post(p.url, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("httppublisher err='%s'", err)
			}

			if resp.StatusCode != 200 || resp.StatusCode != 201 || resp.StatusCode != 202 {
				log.Printf("httppublisher err='http status' status_code=%d", resp.StatusCode)
			}

			p.buffer = make([]*sensor.Measurement, 0)
		case m := <-p.c:
			p.buffer = append(p.buffer, m)
		}
	}
}
