package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/canaryio/canary"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [url]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func scheduler(site *canary.Site, source string, c chan *canary.Sample) {
	t := time.NewTicker(time.Second)
	sampler := canary.NewSampler()

	for {
		select {
		case <-t.C:
			c <- sampler.Sample(site, source)
		}
	}
}

func emit_tsv(s *canary.Sample, source string) {
	fmt.Printf("%s %s %s %s %d %d %f %f %f %f\n",
		s.T.Format(time.RFC3339),
		source,
		s.Site.URL,
		s.IP,
		s.CurlStatus,
		s.HTTPStatus,
		s.NameLookupTime,
		s.ConnectTime,
		s.StartTransferTime,
		s.TotalTime,
	)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	site := &canary.Site{
		URL: args[0],
	}
	c := make(chan *canary.Sample)
	source, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	go scheduler(site, source, c)

	// move samples from the sensor to the reporter
	for sample := range c {
		emit_tsv(sample, source)
	}
}
