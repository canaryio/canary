package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canaryio/canary"
)

var (
	url    string
	source string
	output string
)

func init() {
	flag.StringVar(&url, "u", "http://www.canary.io", "url to monitor")
	flag.StringVar(&source, "s", "unknown", "source / location of this sensor")
	flag.StringVar(&output, "o", "tsv", "output format")
}

func main() {
	flag.Parse()

	// we support various reporters
	// pretty sure we could do something more clever here,
	// but this seems to work well enough
	var r canary.Reporter
	switch output {
	case "tsv":
		r = &canary.TSVReporter{}
	case "logfmt":
		r = &canary.LogfmtReporter{}
	}

	if r == nil {
		fmt.Fprintf(os.Stderr, "-o was set to an invalid option: %s\n", output)
		os.Exit(1)
	}
	go r.Start()

	// fire up a sensor
	s := canary.NewSensor(url, source)
	go s.Start()

	// move samples from the sensor to the reporter
	for sample := range s.C {
		r.Ingest(sample)
	}
}
