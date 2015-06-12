package sampler

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Sample struct {
	StatusCode      int
	TimeStart       time.Time
	TimeToConnect   time.Time
	TimeToFirstByte time.Time
	TimeEnd         time.Time
	ResponseHeaders http.Header
	LocalAddr       net.Addr
	RemoteAddr      net.Addr
}

// StatusCodeError is an error representing an HTTP Status code
// of 400 or greater.
type StatusCodeError struct {
	StatusCode int
}

func (e StatusCodeError) Error() string {
	return fmt.Sprintf(
		"recieved HTTP status %d",
		e.StatusCode,
	)
}

// Ping measures a given URL and returns a Sample
func Ping(target Target, timeout int) (sample Sample, err error) {
	sample.TimeStart = time.Now()
	defer func() { sample.TimeEnd = time.Now() }()

	conn, err := dial(target)
	if err != nil {
		err = fmt.Errorf("connecting: %s", err)
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	sample.TimeToConnect = time.Now()

	sample.LocalAddr = conn.LocalAddr()
	sample.RemoteAddr = conn.RemoteAddr()

	req, err := genRequest(target)
	if err != nil {
		return
	}
	fmt.Fprint(
		conn,
		req)

	r := bufio.NewReader(conn)

	sample.StatusCode, err = parseStatus(r)
	if err != nil {
		err = fmt.Errorf("parsing status: %s", err)
		return
	}

	sample.TimeToFirstByte = time.Now()

	sample.ResponseHeaders, err = parseHeaders(r)
	if err != nil {
		err = fmt.Errorf("parsing headers: %s", err)
		return
	}

	// if we have a Content-Length, go ahead and read the body
	val := sample.ResponseHeaders.Get("Content-Length")
	if val != "" {
		contentLength, err := strconv.Atoi(val)
		if err != nil {
			err = fmt.Errorf("parsing Content-Length: %s", err)
			return sample, err
		}
		n := 0
		buf := make([]byte, contentLength)
		for i := contentLength; i > 0; i = i - n {
			n, err = r.Read(buf)
			if err != nil {
				return sample, err
			}
		}
	}

	if sample.StatusCode >= 400 {
		err = &StatusCodeError{
			StatusCode: sample.StatusCode,
		}
	}

	return
}

func hostString(u *JsonURL) (string, error) {
	// if our Host already has a port, bail out
	if strings.Contains(u.Host, ":") {
		return u.Host, nil
	}

	switch u.Scheme {
	case "http":
		return u.Host + ":80", nil
	case "https":
		return u.Host + ":443", nil
	default:
		return "", fmt.Errorf("unknown URL scheme '%s'", u.Scheme)
	}
}

