package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"

	"github.com/canaryio/canary"

	"github.com/canaryio/canary/pkg/manifest"
	"github.com/canaryio/canary/pkg/libratopublisher"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
)

// builds the app configuration via ENV
func getConfig() (c canary.Config, err error) {
	c.ManifestURL = os.Getenv("MANIFEST_URL")
	if c.ManifestURL == "" {
		err = fmt.Errorf("MANIFEST_URL not defined in ENV")
	}

	interval := os.Getenv("DEFAULT_SAMPLE_INTERVAL")
	// if the variable is unset, an empty string will be returned
	if interval == "" {
		interval = "1"
	}
	c.DefaultSampleInterval, err = strconv.Atoi(interval)
	if err != nil {
		err = fmt.Errorf("DEFAULT_SAMPLE_INTERVAL is not a valid integer")
	}

	// Set RampupSensors if RAMPUP_SENSORS is set to 'yes'
	rampUp := os.Getenv("RAMPUP_SENSORS")
	if rampUp == "yes" {
		c.RampupSensors = true
	} else {
		c.RampupSensors = false
	}

	return
}

func createPublishers() (publishers []canary.Publisher) {
	list := os.Getenv("PUBLISHERS")
	if list == "" {
		list = "stdout"
	}
	publisherList := strings.Split(list, ",")

	for _, publisher := range publisherList {
		switch publisher {
		case "stdout":
			publishers = append(publishers, stdoutpublisher.New())
		case "librato":
			p, err := libratopublisher.NewFromEnv()
			if err != nil {
				log.Fatal(err)
			}
			publishers = append(publishers, p)
		default:
			log.Fatalf("Unknown publisher: %s", publisher)
		}
	}
	
	return
}

func main() {
	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	manifest, err := manifest.GetManifest(conf.ManifestURL)
	if err != nil {
		log.Fatal(err)
	}

	if conf.RampupSensors {
		manifest.GenerateRampupDelays(conf.DefaultSampleInterval)
	}

	c := canary.New(createPublishers())
	c.Config = conf
	c.Manifest = manifest

	// Start canary and block in the signal handler
	c.Run()
	c.SignalHandler()
}
