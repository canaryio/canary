package transportsampler

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/canaryio/canary"
)

// TransportSampler implements Sampler, using http.Transport.
type TransportSampler struct {
	tr http.Transport
}

// New initializes a sane sampler.
func New() TransportSampler {
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
func (s TransportSampler) Sample(target canary.Target) (sample canary.Sample, err error) {
	req, err := http.NewRequest("GET", target.URL, nil)
	if err != nil {
		return sample, err
	}

	req.Header.Add("User-Agent", "canary/"+canary.Version)

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
		err = &canary.StatusCodeError{
			StatusCode: sample.StatusCode,
		}
	}

	return
}
