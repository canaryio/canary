package canary

import (
	"fmt"
	"log"
)

// Publisher is the interface that adds the Publish method.
//
// Pubilsh takes a Site, and Sample, and an error, and is
// expected to deliver that data somewhere.
type Publisher interface {
	Publish(Site, Sample, error) error
}

// StdoutPublisher implements the Publisher interface, and emits
// sampler data to STDOUT
type StdoutPublisher struct{}

// Publish takes a Site, a Sample, and an error (can be nil) and sends that data
// to STDOUT.  An error is retured if anything goes wrong.
func (p StdoutPublisher) Publish(site Site, sample Sample, err error) error {
	duration := sample.T2.Sub(sample.T1).Nanoseconds() / 1000 / 1000

	isOK := true
	if err != nil {
		isOK = false
	}

	errMessage := ``
	if err != nil {
		errMessage = fmt.Sprintf("'%s'", err)
	}

	log.Printf("%s %d %d %t %s", site.URL, sample.StatusCode, duration, isOK, errMessage)
	return nil
}
