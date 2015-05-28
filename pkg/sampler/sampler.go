package sampler

import (
	"bufio"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Target struct {
	URL      string
	Name     string
	Interval int
	// metadata
	Tags               []string
	Attributes         map[string]string
	Hash               string
	RequestHeaders     map[string]string
	InsecureSkipVerify bool
}

func (t *Target) SetHash() {
	jsonTarget, _ := json.Marshal(t)
	hasher := md5.New()
	hasher.Write(jsonTarget)
	t.Hash = hex.EncodeToString(hasher.Sum(nil))
}

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
		return
	}

	sample.TimeToFirstByte = time.Now()

	sample.ResponseHeaders, err = parseHeaders(r)
	if err != nil {
		return
	}

	// if we have a Content-Length, go ahead and read the body
	val := sample.ResponseHeaders.Get("Content-Length")
	if val != "" {
		contentLength, err := strconv.Atoi(val)
		if err != nil {
			return sample, err
		}
		n := 0
		buf := make([]byte, contentLength)
		for i := contentLength; i > 0; i = i - n {
			n, err = r.Read(buf)
		}
	}

	if sample.StatusCode >= 400 {
		err = &StatusCodeError{
			StatusCode: sample.StatusCode,
		}
	}

	return
}

func dial(t Target) (net.Conn, error) {
	u, err := url.Parse(t.URL)
	if err != nil {
		return nil, err
	}

	host, err := hostString(u)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http":
		return net.Dial("tcp", host)
	case "https":
		return tls.Dial("tcp", host, &tls.Config{
			InsecureSkipVerify: t.InsecureSkipVerify,
		})
	default:
		return nil, fmt.Errorf("unknown scheme '%s'", u.Scheme)
	}
}

func hostString(u *url.URL) (string, error) {
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

func parseStatus(r *bufio.Reader) (int, error) {
	statusLine, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}

	parts := strings.Split(statusLine, " ")
	if len(parts) < 3 {
		return 0, fmt.Errorf("'%s' is an invalid HTTP status response", strings.TrimSpace(statusLine))
	}

	status, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return status, nil
}

func parseHeaders(r *bufio.Reader) (http.Header, error) {
	headers := make(http.Header)
	for {
		line, err := r.ReadString('\n')

		if err != nil {
			return headers, nil
		}

		cleanLine := strings.TrimSpace(line)
		if cleanLine == "" {
			// end of headers
			break
		}

		parts := strings.SplitN(cleanLine, ": ", 2)
		if len(parts) == 2 {
			headers.Add(parts[0], parts[1])
		}
	}
	return headers, nil
}

func genRequest(t Target) (string, error) {
	u, err := url.Parse(t.URL)
	if err != nil {
		return "", err
	}

	// allow Host header to be set via t.RequestHeaders
	// otherwise, use the host of the URL
	hostHeader := t.RequestHeaders["Host"]
	if hostHeader == "" {
		hostHeader = u.Host
	}

	// our standard request
	req := fmt.Sprintf("GET %s HTTP/1.1\r\n", u.RequestURI())
	req += fmt.Sprintf("Host: %s\r\n", hostHeader)

	for k, v := range t.RequestHeaders {
		if k != "Host" {
			req += fmt.Sprintf("%s: %s\r\n", k, v)
		}
	}

	// trailing newline
	req += "\r\n"

	return req, nil
}
