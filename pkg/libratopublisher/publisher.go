package libratopublisher

import (
	"fmt"
	"os"

	"github.com/canaryio/canary/pkg/libratoaggregator"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
)

// Publisher implements the canary.Publisher interface and
// is our means of ingesting canary.Measurements and converting
// them to Librato metrics.
type Publisher struct {
	aggregator *libratoaggregator.Aggregator
}

// New takes a user, token and source and return a pointer
// to a Publisher.
func New(user, token, source string) (p *Publisher) {
	p = &Publisher{
		aggregator: libratoaggregator.New(user, token, source),
	}
	return
}

// NewFromEnv is a convenience func that wraps New,
// and populates the required arguments via environment variables.
// If required variables cannot be found, errors are returned.
func NewFromEnv() (*Publisher, error) {
	user := os.Getenv("LIBRATO_USER")
	if user == "" {
		return nil, fmt.Errorf("LIBRATO_USER not set in ENV")
	}

	token := os.Getenv("LIBRATO_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("LIBRATO_TOKEN not set in ENV")
	}

	var err error
	source := os.Getenv("SOURCE")
	if source == "" {
		source, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	return New(user, token, source), nil
}

// Publish takes a canary.Measurement and delivers it to the aggregator.
func (p *Publisher) Publish(m sensor.Measurement) (err error) {
	// convert our measurement into a map of metrics
	// send the map on to the librato aggregator
	p.aggregator.C <- mapMeasurement(m)
	return
}

// mapMeasurments takes a canary.Measurement and returns a map with all of the appropriate metrics
func mapMeasurement(m sensor.Measurement) map[string]float64 {
	metrics := make(map[string]float64)
	// latency
	metrics["canary."+m.Target.Name+".latency"] = m.Sample.Latency()
	if m.Error != nil {
		// increment a general error metric
		metrics["canary."+m.Target.Name+".errors"] = 1

		// increment a specific error metric
		switch m.Error.(type) {
		case sampler.StatusCodeError:
			metrics["canary."+m.Target.Name+".errors.http"] = 1
		default:
			metrics["canary."+m.Target.Name+".errors.sampler"] = 1
		}
	}

	return metrics
}
