package sampler

import (
	"net/http"
	"net/url"
	"bufio"
	"strings"
	"fmt"
	"strconv"
)

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
