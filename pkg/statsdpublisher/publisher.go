package statsdpublisher

import (
	"fmt"
	"io"
	"net"
	"os"

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
	var err error
	source := os.Getenv("SOURCE")
	if source == "" {
		source, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	host := os.Getenv("STATSD_HOST")
	if host == "" {
		return nil, fmt.Errorf("STATSD_HOST not set in ENV")
	}

	prefix := os.Getenv("STATSD_PREFIX")
	if prefix == "" {
		prefix = "canary"
	}

	return New(source, host, prefix)
}

func (p *Publisher) Publish(m sensor.Measurement) error {
	return Write(p.conn, p.Source, p.Prefix, m)
}

func Write(w io.Writer, source, prefix string, m sensor.Measurement) error {
	timeConnect := int(m.Sample.TimeToConnect.Sub(m.Sample.TimeStart).Seconds() * 1000)
	_, err := fmt.Fprintf(w,
		"%s.%s.%s.time_to_connect:%d|ms\n",
		prefix,
		m.Target.Name,
		source,
		timeConnect)

	timeToFirstByte := int(m.Sample.TimeToFirstByte.Sub(m.Sample.TimeStart).Seconds() * 1000)
	_, err = fmt.Fprintf(w,
		"%s.%s.%s.time_to_first_byte:%d|ms\n",
		prefix,
		m.Target.Name,
		source,
		timeToFirstByte)

	timeTotal := int(m.Sample.TimeEnd.Sub(m.Sample.TimeStart).Seconds() * 1000)
	_, err = fmt.Fprintf(w,
		"%s.%s.%s.time_total:%d|ms\n",
		prefix,
		m.Target.Name,
		source,
		timeTotal)

	_, err = fmt.Fprintf(w,
		"%s.%s.%s.http_status.%d:1|c\n",
		prefix,
		m.Target.Name,
		source,
		m.Sample.StatusCode)

	return err
}
