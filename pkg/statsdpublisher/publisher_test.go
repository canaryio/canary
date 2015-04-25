package statsdpublisher

import (
	"os"
	"time"

	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
)

func ExampleStatsdpublisher_Write() {
	target := sampler.Target{
		Name: "github",
		URL:  "https://github.com",
	}

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:01Z")

	sample := sampler.Sample{
		TimeStart:       t1,
		TimeToConnect:   t2,
		TimeToFirstByte: t2,
		TimeEnd:         t2,
		StatusCode:      200,
	}

	Write(os.Stdout, "test", "canary", sensor.Measurement{
		Target: target,
		Sample: sample,
	})
	// Output:
	// canary.github.test.time_to_connect:1000|ms
	// canary.github.test.time_to_first_byte:1000|ms
	// canary.github.test.time_total:1000|ms
	// canary.github.test.http_status.200:1|c
}
