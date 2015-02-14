package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
)

type command struct {
	sampler   sampler.Sampler
	publisher canary.Publisher
	target    sampler.Target
}

func (cmd command) Run() {
	sensor := sensor.Sensor{
		Target:  cmd.target,
		C:       make(chan sensor.Measurement),
		Sampler: cmd.sampler,
	}
	go sensor.Start()

	for m := range sensor.C {
		cmd.publisher.Publish(m)
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
		target: sampler.Target{
			URL: args[0],
		},
		sampler:   sampler.New(),
		publisher: stdoutpublisher.New(),
	}
	cmd.Run()
}
