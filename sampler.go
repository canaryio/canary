package canary

import (
	"time"

	"github.com/andelf/go-curl"
)

type Sample struct {
	URL               string
	Name              string
	T                 time.Time
	Source            string
	IP                string
	CurlStatus        int
	HTTPStatus        int
	NameLookupTime    float64
	ConnectTime       float64
	StartTransferTime float64
	TotalTime         float64
}

type Sampler struct {
	easy *curl.CURL
}

func NewSampler() *Sampler {
	return &Sampler{
		easy: curl.EasyInit(),
	}
}

func (s *Sampler) Sample(site *Site, source string) *Sample {
	defer s.easy.Reset()

	// curl configuration
	s.easy.Setopt(curl.OPT_URL, site.URL)
	noOut := func(buf []byte, userdata interface{}) bool {
		return true
	}
	s.easy.Setopt(curl.OPT_WRITEFUNCTION, noOut)
	s.easy.Setopt(curl.OPT_TIMEOUT, 10)

	sample := &Sample{
		URL:    site.URL,
		Name:   site.Name,
		Source: source,
		T:      time.Now(),
		IP:     "n/a",
	}

	if err := s.easy.Perform(); err != nil {
		if e, ok := err.(curl.CurlError); ok {
			sample.CurlStatus = (int(e))
			return sample
		}
		return sample
	}

	httpStatus, _ := s.easy.Getinfo(curl.INFO_RESPONSE_CODE)
	sample.HTTPStatus = httpStatus.(int)

	ip, _ := s.easy.Getinfo(curl.INFO_PRIMARY_IP)
	sample.IP = ip.(string)

	namelookupTime, _ := s.easy.Getinfo(curl.INFO_NAMELOOKUP_TIME)
	sample.NameLookupTime = namelookupTime.(float64) * 1000
	connectTime, _ := s.easy.Getinfo(curl.INFO_CONNECT_TIME)
	sample.ConnectTime = connectTime.(float64) * 1000

	starttransferTime, _ := s.easy.Getinfo(curl.INFO_STARTTRANSFER_TIME)
	sample.StartTransferTime = starttransferTime.(float64) * 1000

	totalTime, _ := s.easy.Getinfo(curl.INFO_TOTAL_TIME)
	sample.TotalTime = totalTime.(float64) * 1000

	return sample
}
