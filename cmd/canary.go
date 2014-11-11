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
	flag.StringVar(&output, "o", "logfmt", "output format")
}

func main() {
	flag.Parse()

	var r canary.Reporter
	// an opportunity to change the output
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

	s := canary.NewSensor(url, source)
	go s.Start()

	for sample := range s.C {
		r.Ingest(sample)
	}
}
