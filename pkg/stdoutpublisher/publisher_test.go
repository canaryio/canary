package stdoutpublisher

import (
	"time"

	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
)

func ExamplePublisher_Publish() {
	target := sampler.Target{
		URL: "http://www.canary.io",
	}

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:01Z")

	sample := sampler.Sample{
		T1:         t1,
		T2:         t2,
		StatusCode: 200,
	}

	p := New()
	p.Publish(sensor.Measurement{
		Target:     target,
		Sample:     sample,
		IsOK:       true,
		StateCount: 2,
	})
	// Output:
	// 2014-12-28T00:00:01Z http://www.canary.io 200 1000.000000 true 2
}
