package librato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Gauge struct {
	Name        string  `json:"name"`
	MeasureTime int64   `json:"measure_time"`
	Source      string  `json:"source"`
	Count       int64   `json:"count"`
	Sum         float64 `json:"sum"`
	Min         float64 `json:"min",omitempty`
	Max         float64 `json:"max",omitempty`
	SumSquares  float64 `json:"max",omitempty`
}

type Payload struct {
	Gauges []Gauge `json:"gauges"`
}

type Client struct {
	sync.Mutex
	User    string
	Token   string
	payload Payload
}

func (c *Client) AddGauge(g Gauge) {
	c.Lock()
	defer c.Unlock()
	if g.MeasureTime == 0 {
		g.MeasureTime = time.Now().Unix()
	}
	c.payload.Gauges = append(c.payload.Gauges, g)
}

func (c *Client) Flush() error {
	c.Lock()
	defer c.Unlock()

	if len(c.payload.Gauges) == 0 {
		return nil
	}

	b, err := json.Marshal(c.payload)
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
	req.SetBasicAuth(c.User, c.Token)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("status code: %d\n", res.StatusCode)
	}

	c.payload = Payload{}

	return nil
}
