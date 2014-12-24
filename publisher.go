package canary

import (
	"fmt"
	"log"
)

type Publisher interface {
	Publish(Site, Sample, error) error
}

type StdoutPublisher struct{}

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
