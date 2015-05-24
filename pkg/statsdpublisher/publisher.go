package statsdpublisher

import (
	"fmt"
	"io"
	"net"

	"github.com/canaryio/canary/pkg/sensor"
)

// Publisher implements canary.Publisher
type Publisher struct {
	conn   net.Conn
	Source string
	Prefix string
}

func New(source, host, prefix string) (*Publisher, error) {
	conn, err := net.Dial("udp", host)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn:   conn,
		Source: source,
		Prefix: prefix,
	}, nil
}

func NewFromEnv() (*Publisher, error) {
	return New("test", "45.55.133.179:8125", "canary")
}

func (p *Publisher) Publish(m sensor.Measurement) error {
	return Write(p.conn, p.Source, p.Prefix, m)
}

func Write(w io.Writer, source, prefix string, m sensor.Measurement) error {
	timeTotal := int(m.Sample.TimeEnd.Sub(m.Sample.TimeStart).Seconds() * 1000)
	_, err := fmt.Fprintf(w,
		"%s.%s.%s.time_total:%d|ms\n",
		prefix,
		m.Target.Name,
		source,
		timeTotal)
	return err
}
