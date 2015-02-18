package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
)

type command struct {
	sampler   sampler.Sampler
	publisher canary.Publisher
	target    sampler.Target
	interval  int
}

func (cmd command) Run() {
	sensor := sensor.Sensor{
		Target:  cmd.target,
		C:       make(chan sensor.Measurement),
		Sampler: cmd.sampler,
	}
	go sensor.Start(cmd.interval)

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

	cmd := command{
		target: sampler.Target{
			URL: args[0],
		},
		sampler:   sampler.New(),
		publisher: stdoutpublisher.New(),
		interval: sample_interval,
	}
	cmd.Run()
}
