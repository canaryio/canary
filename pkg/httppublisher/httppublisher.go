package httppublisher

import "github.com/canaryio/canary/pkg/sensor"

type Publisher struct{}

func New() *Publisher {
	return &Publisher
}

func (p *Publisher) Publish(m sensor.Measurement) (err error) {
	return
}
