package statsdpublisher

import (
	"os"
	"time"

	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
)

func ExampleStatsdpublisher_Write() {
	target := sampler.Target{
		Name: "foo",
		URL:  "http://www.canary.io",
	}

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:01Z")

	sample := sampler.Sample{
		TimeStart:  t1,
		TimeEnd:    t2,
		StatusCode: 200,
	}

	Write(os.Stdout, "test", "canary", sensor.Measurement{
		Target:     target,
		Sample:     sample,
		IsOK:       true,
		StateCount: 2,
	})
	// Output:
	// canary.foo.test.time_total:1000|ms
}
