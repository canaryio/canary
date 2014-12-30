package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canaryio/canary"
)

type command struct {
	sampler   canary.Sampler
	publisher canary.Publisher
	target    canary.Target
}

func (cmd command) Run() {
	sensor := canary.Sensor{
		Target:  cmd.target,
		C:       make(chan canary.Measurement),
		Sampler: canary.NewTransportSampler(),
	}
	go sensor.Start()

	for m := range sensor.C {
		cmd.publisher.Publish(m.Target, m.Sample, m.Error)
	}
}

// usage prints a useful usage message.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [url]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	cmd := command{
		target: canary.Target{
			URL: args[0],
		},
		sampler:   canary.NewTransportSampler(),
		publisher: canary.StdoutPublisher{},
	}
	cmd.Run()
}
