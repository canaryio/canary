package canary

import "github.com/canaryio/canary/pkg/sensor"

// Publisher is the interface that adds the Publish method.
//
// Publish takes a Target, and Sample, and an error, and is
// expected to deliver that data somewhere.
type Publisher interface {
	Publish(sensor.Measurement) error
}
