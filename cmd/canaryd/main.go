package main

import (
	"fmt"
	"log"
	"os"

	"github.com/canaryio/canary"
)

type config struct {
	ManifestURL string
}

// builds the app configuration via ENV
func getConfig() (c config, err error) {
	c.ManifestURL = os.Getenv("MANIFEST_URL")
	if c.ManifestURL == "" {
		err = fmt.Errorf("MANIFEST_URL not defined in ENV")
	}
	return
}

func main() {
	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	manifest, err := canary.GetManifest(conf.ManifestURL)
	if err != nil {
		log.Fatal(err)
	}

	// output chan
	c := make(chan canary.Measurement)

	p := canary.StdoutPublisher{}

	// spinup a sensor for each target
	for _, target := range manifest.Targets {
		sensor := canary.Sensor{
			Target:  target,
			C:       c,
			Sampler: canary.NewTransportSampler(),
		}
		go sensor.Start()
	}

	for m := range c {
		p.Publish(m.Target, m.Sample, m.Error)
	}
}
