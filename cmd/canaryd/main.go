package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/libratopublisher"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
	"github.com/canaryio/canary/pkg/transportsampler"
)

type config struct {
	ManifestURL   string
	PublisherList []string
}

// builds the app configuration via ENV
func getConfig() (c config, err error) {
	c.ManifestURL = os.Getenv("MANIFEST_URL")
	if c.ManifestURL == "" {
		err = fmt.Errorf("MANIFEST_URL not defined in ENV")
	}

	list := os.Getenv("PUBLISHERS")
	if list == "" {
		list = "stdout"
	}
	c.PublisherList = strings.Split(list, ",")

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

	var publishers []canary.Publisher

	// spinup publishers
	for _, publisher := range conf.PublisherList {
		switch publisher {
		case "stdout":
			p := stdoutpublisher.New()
			publishers = append(publishers, p)
		case "librato":
			p, err := libratopublisher.NewFromEnv()
			if err != nil {
				log.Fatal(err)
			}
			go p.Start()
			publishers = append(publishers, p)
		default:
			log.Printf("Unknown publisher: %s", publisher)
		}
	}

	// spinup a scheduler for each target
	for _, target := range manifest.Targets {
		scheduler := canary.Scheduler{
			Target:  target,
			C:       c,
			Sampler: transportsampler.New(),
		}
		go scheduler.Start()
	}

	// publish each incoming measurement
	for m := range c {
		for _, p := range publishers {
			p.Publish(m)
		}
	}
}
