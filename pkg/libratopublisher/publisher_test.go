package libratopublisher

import (
	"fmt"
	"testing"
	"time"

	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
)

func TestGoodMeasurement(t *testing.T) {
	expectedLatency := 7000.0

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:07Z")

	m := sensor.Measurement{
		Target: sampler.Target{
			Name: "test",
		},
		Sample: sampler.Sample{
			TimeStart:  t1,
			TimeEnd:    t2,
			StatusCode: 200,
		},
	}
	res := mapMeasurement(m)

	if len(res) != 1 {
		t.Fatalf("expected 1 metric to be in this list, found %d", len(res))
	}

	val := res["canary.test.latency"]
	if val != expectedLatency {

		t.Fatalf(
			"expected canary.test.latency to equal %f, but it was %f",
			expectedLatency,
			val,
		)
	}
}

func TestBadHTTPMeasurement(t *testing.T) {
	expectedLatency := 7000.0

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:07Z")

	m := sensor.Measurement{
		Target: sampler.Target{
			Name: "test",
		},
		Sample: sampler.Sample{
			TimeStart:  t1,
			TimeEnd:    t2,
			StatusCode: 502,
		},
		Error: sampler.StatusCodeError{
			StatusCode: 502,
		},
	}
	res := mapMeasurement(m)

	val := res["canary.test.latency"]
	if val != expectedLatency {
		t.Fatalf(
			"expected canary.test.latency to equal %f, but it was %f",
			expectedLatency,
			val,
		)
	}

	val = res["canary.test.errors"]
	if res["canary.test.errors"] != 1.0 {
		t.Errorf(
			"expected canary.test.errors to equal %f, but it was %f",
			1.0,
			val,
		)
	}

	val = res["canary.test.errors.http"]
	if val != 1.0 {
		t.Errorf(
			"expected canary.test.errors.http to equal %f, but it was %f",
			1.0,
			val,
		)
	}
}

func TestBadTransportMeasurement(t *testing.T) {
	expectedLatency := 7000.0

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:07Z")

	m := sensor.Measurement{
		Target: sampler.Target{
			Name: "test",
		},
		Sample: sampler.Sample{
			TimeStart: t1,
			TimeEnd:   t2,
		},
		Error: fmt.Errorf("test error"),
	}
	res := mapMeasurement(m)

	val := res["canary.test.latency"]
	if val != expectedLatency {
		t.Fatalf(
			"expected canary.test.latency to equal %f, but it was %f",
			expectedLatency,
			val,
		)
	}

	val = res["canary.test.errors"]
	if res["canary.test.errors"] != 1.0 {
		t.Errorf(
			"expected canary.test.errors to equal %f, but it was %f",
			1.0,
			val,
		)
	}

	val = res["canary.test.errors.sampler"]
	if val != 1.0 {
		t.Errorf(
			"expected canary.test.errors.sampler to equal %f, but it was %f",
			1.0,
			val,
		)
	}
}
