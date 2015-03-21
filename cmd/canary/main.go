package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/manifest"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
)

// usage prints a useful usage message.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [url]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

// builds the app configuration via ENV
func getConfig() (c canary.Config, url string, err error) {
	flag.Usage = usage
	flag.Parse()

	sampleIntervalString := os.Getenv("SAMPLE_INTERVAL")
	sampleInterval := 1
	if sampleIntervalString != "" {
		sampleInterval, err = strconv.Atoi(sampleIntervalString)
		if err != nil {
			err = fmt.Errorf("SAMPLE_INTERVAL is not a valid integer")
		}
	}
	c.DefaultSampleInterval = sampleInterval

	timeout := 0
	defaultTimeout := os.Getenv("DEFAULT_MAX_TIMEOUT")
	if defaultTimeout == "" {
		timeout = 10
	} else {
		timeout, err = strconv.Atoi(defaultTimeout)
		if err != nil {
			err = fmt.Errorf("DEFAULT_MAX_TIMOEUT is not a valid integer")
		}
	}
	c.MaxSampleTimeout = timeout

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
	url = args[0]

	return
}

func main() {
	conf, url, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	c := canary.New([]canary.Publisher{stdoutpublisher.New()})
	manifest := manifest.Manifest{}

	manifest.StartDelays = []float64{0.0}
	manifest.Targets = []sampler.Target{
		sampler.Target{
			URL:      url,
			Interval: conf.DefaultSampleInterval,
		},
	}

	c.Config = conf
	c.Manifest = manifest

	// Start canary and block in the signal handler
	c.Run()
	c.SignalHandler()
}
