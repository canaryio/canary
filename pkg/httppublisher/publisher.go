package httppublisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/canaryio/canary/pkg/sensor"
)

const MAX_BUFFER = 1000

type Publisher struct {
	url      string
	interval time.Duration
	buffer   []*sensor.Measurement
	c        chan *sensor.Measurement
	i        int
}

func New(url string, interval time.Duration) *Publisher {
	p := &Publisher{
		buffer:   make([]*sensor.Measurement, MAX_BUFFER),
		c:        make(chan *sensor.Measurement),
		url:      url,
		interval: interval,
		i:        0,
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
			// flush when our ticker fires
			if err := p.flush(); err != nil {
				log.Print(err)
			}
		case m := <-p.c:
			p.buffer[p.i] = m
			p.i++

			// flush if we've exceeded the bounds of our buffer
			if p.i >= MAX_BUFFER {
				if err := p.flush(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (p *Publisher) flush() error {
	payload, err := json.Marshal(p.buffer)
	defer func() { p.i = 0 }()

	if err != nil {
		return fmt.Errorf("httppublisher err='%s'", err)
	}

	resp, err := http.Post(p.url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("httppublisher err='%s'", err)
	}

	if resp.StatusCode != 200 || resp.StatusCode != 201 || resp.StatusCode != 202 {
		return fmt.Errorf("httppublisher err='http status' status_code=%d", resp.StatusCode)
	}

	return nil
}
