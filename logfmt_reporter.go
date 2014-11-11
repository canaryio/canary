package canary

import (
	"fmt"
	"time"
)

type LogfmtReporter struct{}

func (l *LogfmtReporter) Ingest(s *Sample) error {
	fmt.Printf("t=%s source=%s url=%s ip=%s curl_status=%d http_status=%d namelookup_time=%f connect_time=%f start_transfer_time=%f total_time=%f\n",
		s.T.Format(time.RFC3339),
		s.Source,
		s.URL,
		s.IP,
		s.CurlStatus,
		s.HTTPStatus,
		s.NameLookupTime,
		s.ConnectTime,
		s.StartTransferTime,
		s.TotalTime,
	)

	return nil
}

func (l *LogfmtReporter) Start() {
}

func (l *LogfmtReporter) Stop() {
}
