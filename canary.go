package canary

import (
	"sync"
	"time"

	"github.com/andelf/go-curl"
)

// A Site represents that which is being measured
type Site struct {
	URL  string
	Name string
}

// A Sample contains a point in time measurement made against a given Site.
type Sample struct {
	Site              *Site
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

// A Sampler is capable of generating Samples from a given Site
type Sampler struct {
	sync.Mutex
	easy *curl.CURL
}

// Sample generates a Sample for a given Site and source.
func (s *Sampler) Sample(site *Site, source string) *Sample {
	s.Lock()
	defer s.easy.Reset()
	defer s.Unlock()

	if s.easy == nil {
		s.easy = curl.EasyInit()
	}

	// curl configuration
	s.easy.Setopt(curl.OPT_URL, site.URL)
	noOut := func(buf []byte, userdata interface{}) bool {
		return true
	}
	s.easy.Setopt(curl.OPT_WRITEFUNCTION, noOut)
	s.easy.Setopt(curl.OPT_TIMEOUT, 10)

	sample := &Sample{
		Site:   site,
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
