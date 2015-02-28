package libratoaggregator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/canaryio/canary/pkg/canaryversion"
	"log"
	"net/http"
	"time"
)

type gauge struct {
	Name        string  `json:"name"`
	MeasureTime int64   `json:"measure_time"`
	Source      string  `json:"source"`
	Count       int64   `json:"count"`
	Sum         float64 `json:"sum"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	SumSquares  float64 `json:"sum_squares"`
}

type payload struct {
	Gauges []gauge `json:"gauges"`
}

// An Aggregator is a continuous client to Librato.
// It listens for new metric pairs to come in over a channel,
// aggregates them according to name, and submits those metrics
// to the Librato API every 5 seconds.
type Aggregator struct {
	User   string
	Token  string
	Source string
	C      chan map[string]float64
	gauges map[string]gauge
}

// New takes a user, token and source and returns a
// pointer to an Aggregator.
func New(user, token, source string) (a *Aggregator) {
	a = &Aggregator{
		User:   user,
		Token:  token,
		Source: source,
		C:      make(chan map[string]float64),
		gauges: make(map[string]gauge),
	}
	go a.start()
	return
}

// Start begins the event loop for an Aggregator.  It is blocking,
// and should be run in a goroutine.
func (a *Aggregator) start() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case m := <-a.C:
			for k, v := range m {
				a.ingest(k, v)
			}
		case <-t.C:
			if len(a.gauges) > 0 {
				err := a.flush()
				if err != nil {
					log.Printf("librato: error submitting to API: %s", err)
					continue
				}
				a.gauges = make(map[string]gauge)
			}
		}
	}
}

func (a *Aggregator) ingest(k string, v float64) {
	g := a.gauges[k]
	g.Name = k
	g.MeasureTime = time.Now().Unix()
	g.Source = a.Source
	g.Count++
	g.Sum += v
	if v > g.Max {
		g.Max = v
	}
	if v < g.Min || g.Min == 0 {
		g.Min = v
	}
	g.SumSquares += v * v

	a.gauges[k] = g
}

func (a *Aggregator) flush() error {
	// prepare a payload from our collected gauges
	p := payload{
		Gauges: make([]gauge, 0),
	}

	for _, g := range a.gauges {
		p.Gauges = append(p.Gauges, g)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"https://metrics-api.librato.com/v1/metrics",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "canary/"+canaryversion.Version)

	req.SetBasicAuth(a.User, a.Token)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("status code: %d\n", res.StatusCode)
	}
	return nil
}
