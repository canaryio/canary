package stdoutpublisher

import (
	"fmt"
	"time"

	"github.com/canaryio/canary/pkg/sensor"
)

// Publisher implements canary.Publisher, and is our
// gateway for delivering canary.Measurement data to STDOUT.
type Publisher struct{}

// New returns a pointer to a new Publsher.
func New() *Publisher {
	return &Publisher{}
}

// Publish takes a canary.Measurement and emits data to STDOUT.
func (p *Publisher) Publish(m sensor.Measurement) (err error) {
	duration := m.Sample.T2.Sub(m.Sample.T1).Seconds() * 1000

	isOK := true
	if m.Error != nil {
		isOK = false
	}

	errMessage := ``
	if m.Error != nil {
		errMessage = fmt.Sprintf("'%s'", m.Error)
	}

	fmt.Printf(
		"%s %s %d %f %t %s\n",
		m.Sample.T2.Format(time.RFC3339),
		m.Target.URL,
		m.Sample.StatusCode,
		duration,
		isOK,
		errMessage,
	)
	return
}
