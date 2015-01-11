package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
	"github.com/canaryio/canary/pkg/transportsampler"
)

type command struct {
	sampler   canary.Sampler
	publisher canary.Publisher
	target    canary.Target
}

func (cmd command) Run() {
	scheduler := canary.Scheduler{
		Target:  cmd.target,
		C:       make(chan canary.Measurement),
		Sampler: cmd.sampler,
	}
	go scheduler.Start()

	for m := range scheduler.C {
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
		target: canary.Target{
			URL: args[0],
		},
		sampler:   transportsampler.New(),
		publisher: stdoutpublisher.New(),
	}
	cmd.Run()
}
