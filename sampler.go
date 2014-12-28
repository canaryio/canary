package canary

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Target represents the things that we are measureing.
type Target struct {
	URL  string
	Name string
}

// Sample represents HTTP state from a given point in time.
type Sample struct {
	StatusCode int
	T1         time.Time
	T2         time.Time
}

// A Sampler is an interface that provides the Sample method.
type Sampler interface {
	Sample(Target) (Sample, error)
}

// TransportSampler implements Sampler, using http.Transport.
type TransportSampler struct {
	tr http.Transport
}

// NewTransportSampler initializes a sane sampler.
func NewTransportSampler() TransportSampler {
	return TransportSampler{
		tr: http.Transport{
			DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 10*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second))
				return c, nil
			},
		},
	}
}

// Sample measures a given target and returns both a Sample and error details.
func (s TransportSampler) Sample(target Target) (sample Sample, err error) {
	req, err := http.NewRequest("GET", target.URL, nil)
	if err != nil {
		return sample, err
	}

	sample.T1 = time.Now()
	defer func() { sample.T2 = time.Now() }()

	resp, err := s.tr.RoundTrip(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	sample.StatusCode = resp.StatusCode
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if sample.StatusCode >= 400 {
		err = fmt.Errorf("http status code is %d", sample.StatusCode)
	}

	return
}
