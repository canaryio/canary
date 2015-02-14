package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canaryio/canary/pkg/libratopublisher"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
	"github.com/canaryio/canary/pkg/stdoutpublisher"
)

// the public version of canary
const Version string = "v3"

// Publisher is the interface that adds the Publish method.
//
// Pubilsh takes a Target, and Sample, and an error, and is
// expected to deliver that data somewhere.
type Publisher interface {
	Publish(sensor.Measurement) error
}

type config struct {
	PublisherList []string
}

// builds the app configuration via ENV
func getConfig() (c config, err error) {
	list := os.Getenv("PUBLISHERS")
	if list == "" {
		list = "stdout"
	}
	c.PublisherList = strings.Split(list, ",")

	return
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

	target := sampler.Target{
		URL: args[0],
	}

	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	// output chan
	c := make(chan sensor.Measurement)

	var publishers []Publisher

	// spinup publishers
	for _, publisher := range conf.PublisherList {
		switch publisher {
		case "stdout":
			p := stdoutpublisher.New()
			publishers = append(publishers, p)
		case "librato":
			p, err := libratopublisher.NewFromEnv()
			if err != nil {
				log.Fatal(err)
			}
			publishers = append(publishers, p)
		default:
			log.Printf("Unknown publisher: %s", publisher)
		}
	}

	// spinup a sensor for our target
	sensor := sensor.Sensor{
		Target:  target,
		C:       make(chan sensor.Measurement),
		Sampler: sampler.New(),
	}
	go sensor.Start()

	// publish each incoming measurement
	for m := range c {
		for _, p := range publishers {
			p.Publish(m)
		}
	}
}
