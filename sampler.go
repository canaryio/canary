package canary

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Site struct {
	URL  string
	Name string
}

type Sample struct {
	StatusCode int
	T1         time.Time
	T2         time.Time
}

type Sampler interface {
	Sample(Site) (Sample, error)
}

type TransportSampler struct {
	tr http.Transport
}

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

func (s TransportSampler) Sample(site Site) (sample Sample, err error) {
	req, err := http.NewRequest("GET", site.URL, nil)
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
