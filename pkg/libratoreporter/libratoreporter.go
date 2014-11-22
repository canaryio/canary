package libratoreporter

import (
	"log"
	"strings"
	"time"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/librato"
)

type Reporter struct {
	ingest chan *canary.Sample
	stop   chan int
	gauges map[string]*librato.Gauge
	client *librato.Client
}

type Config struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

func New(c *Config) *Reporter {
	return &Reporter{
		ingest: make(chan *canary.Sample),
		stop:   make(chan int),
		gauges: make(map[string]*librato.Gauge),
		client: &librato.Client{
			User:  c.User,
			Token: c.Token,
		},
	}
}

func (r *Reporter) Start() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-r.stop:
			break
		case <-t.C:
			r.flush()
		case sample := <-r.ingest:
			r.push(sample)
		}
	}
}

func (r *Reporter) Stop() {
	close(r.stop)
}

func (r *Reporter) Ingest(s *canary.Sample) error {
	r.ingest <- s
	return nil
}

func (r *Reporter) push(s *canary.Sample) {
	r.merge(LatencyGauge(s))
	r.merge(HTTPErrorGauge(s))
	r.merge(TransportErrorGauge(s))
}

func (r *Reporter) merge(g *librato.Gauge) {
	og := r.gauges[g.Name]
	if og == nil {
		r.gauges[g.Name] = g
		return
	}

	og.Count++
	og.Sum += g.Sum
	og.SumSquares += g.Sum * g.Sum

	if g.Sum < og.Min {
		og.Min = g.Sum
	}

	if g.Sum < og.Max {
		og.Max = g.Sum
	}
}

func (r *Reporter) flush() error {
	for _, g := range r.gauges {
		r.client.AddGauge(*g)
	}
	err := r.client.Flush()
	if err != nil {
		log.Fatal(err)
	}

	r.gauges = make(map[string]*librato.Gauge)
	return nil
}

func LatencyGauge(s *canary.Sample) *librato.Gauge {
	name := strings.Join([]string{s.Name, "latency"}, ".")
	return &librato.Gauge{
		Name:   name,
		Source: s.Source,
		Count:  1,
		Sum:    s.TotalTime,
	}
}

func HTTPErrorGauge(s *canary.Sample) *librato.Gauge {
	name := strings.Join([]string{s.Name, "errors", "http"}, ".")
	count := 0
	if s.HTTPStatus > 399 {
		count = 1
	}

	return &librato.Gauge{
		Name:   name,
		Source: s.Source,
		Count:  int64(count),
		Sum:    float64(count),
	}
}

func TransportErrorGauge(s *canary.Sample) *librato.Gauge {
	name := strings.Join([]string{s.Name, "errors", "transport"}, ".")
	count := 0
	if s.CurlStatus > 0 {
		count = 1
	}

	return &librato.Gauge{
		Name:   name,
		Source: s.Source,
		Count:  int64(count),
		Sum:    float64(count),
	}
}
