package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/manifest"
)

// usage prints a useful usage message.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [url]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	interval_string := os.Getenv("SAMPLE_INTERVAL")
	if interval_string == "" {
		interval_string = "1"
	}
	sample_interval, err := strconv.Atoi(interval_string)
	if err != nil {
		err = fmt.Errorf("SAMPLE_INTERVAL is not a valid integer")
	}

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	c := canary.New()
	conf := canary.Config{}
	manifest := manifest.Manifest{}

	conf.DefaultSampleInterval = sample_interval
	conf.PublisherList = []string{"stdout"}
	manifest.StartDelays = []float64{0.0}
	manifest.Targets = []sampler.Target{ sampler.Target{URL: args[0]} }

	c.Config = conf
	c.Manifest = manifest

	// Start canary and block in the signal handler
	c.Run()
	c.SignalHandler()
}