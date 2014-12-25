package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/canaryio/canary"
)

type command struct {
	sampler   canary.Sampler
	publisher canary.Publisher
	target    canary.Target
}

func (cmd command) Run() {
	c := make(chan measurement)
	go scheduler(cmd.target, cmd.sampler, c)

	for m := range c {
		cmd.publisher.Publish(m.Target, m.Sample, m.Error)
	}
}

type measurement struct {
	Target canary.Target
	Sample canary.Sample
	Error  error
}

// usage prints a useful usage message.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [url]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

// schedule repeatedly produces samples of a given canary.Target and reports
// the samples over a channel.
func scheduler(target canary.Target, sampler canary.Sampler, c chan measurement) {
	t := time.NewTicker(time.Second)

	for {
		select {
		case <-t.C:
			sample, err := sampler.Sample(target)
			m := measurement{
				Target: target,
				Sample: sample,
				Error:  err,
			}
			c <- m
		}
	}
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
