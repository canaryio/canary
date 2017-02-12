package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/canaryio/canary"

	"github.com/canaryio/canary/pkg/libratopublisher"
	"github.com/canaryio/canary/pkg/manifest"
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

	defaultTimeout := os.Getenv("DEFAULT_MAX_TIMEOUT")
	// if the variable is unset, an empty string will be returned
	if defaultTimeout == "" {
		c.MaxSampleTimeout = 10
	} else {
		c.MaxSampleTimeout, err = strconv.Atoi(defaultTimeout)
		if err != nil {
			err = fmt.Errorf("DEFAULT_MAX_TIMOEUT is not a valid integer")
		}
	}

	if err == nil {
		c.DefaultSampleInterval, err = strconv.Atoi(interval)
		if err != nil {
			err = fmt.Errorf("DEFAULT_SAMPLE_INTERVAL is not a valid integer")
		}
	}

	if err == nil {
		autoReloadInterval := os.Getenv("AUTO_RELOAD_INTERVAL")
		if autoReloadInterval == "" {
			autoReloadInterval = "0"
		}

		duration, err := time.ParseDuration(autoReloadInterval + "s")
		if err != nil {
			err = fmt.Errorf("AUTO_RELOAD_INTERVAL is not a valid value for seconds")
		}
		c.ReloadInterval = duration
	}

	if err == nil {
		maxReloadFailures := os.Getenv("MAX_RELOAD_FAILURES")
		if maxReloadFailures == "" {
			maxReloadFailures = "5"
		}

		c.MaxReloadFailures, err = strconv.Atoi(maxReloadFailures)
		if err != nil {
			err = fmt.Errorf("MAX_RELOAD_FAILURES is not a valid integer")
		}
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

	manifest, err := manifest.Get(conf.ManifestURL, conf.DefaultSampleInterval)
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

	u, _ := time.ParseDuration("0s")
	if c.Config.ReloadInterval != u {
		go c.StartAutoReload(c.Config.ReloadInterval)
	}
	c.SignalHandler()
}
