package canary

import (
	"fmt"
	"time"
)

type TSVReporter struct{}

func (l *TSVReporter) Ingest(s *Sample) error {
	fmt.Printf("%s %s %s %s %d %d %f\n",
		s.T.Format(time.RFC3339),
		s.Source,
		s.URL,
		s.IP,
		s.CurlStatus,
		s.HTTPStatus,
		s.TotalTime,
	)

	return nil
}

// noops for the love of the interface
func (l TSVReporter) Start() {}
func (l TSVReporter) Stop()  {}
