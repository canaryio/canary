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
	TimeToResolveIP time.Time
	TimeToConnect   time.Time
	TimeToFirstByte time.Time
	TimeEnd         time.Time
	ResponseHeaders http.Header
	LocalAddr       net.IP
	RemoteAddr      net.IP
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
	// we require four pieces of information to make the request:
	// * hostname to connect to
	// * port to connect to
	// * ip address of the hostname
	// * Host header value
	
	hostname, port, err := hostnameAndPort(&target.URL)
	if err != nil {
		return
	}

	sample.TimeStart = time.Now()
	defer func() { sample.TimeEnd = time.Now() }()

	deadline := sample.TimeStart.Add(time.Duration(timeout) * time.Second)
	
	ip, err := resolveIPAddr(hostname, deadline)
	if err != nil {
		err = fmt.Errorf("resolving IP for %s: %s", hostname, err)
		return
	}

	sample.TimeToResolveIP = time.Now()
	sample.RemoteAddr = ip
	
	conn, err := dial(target.URL.Scheme, ip.String() + ":" + port, hostname, deadline, target.InsecureSkipVerify)
	if err != nil {
		err = fmt.Errorf("connecting: %s", err)
		return
	}
	defer conn.Close()
	conn.SetDeadline(deadline)

	sample.TimeToConnect = time.Now()
	sample.LocalAddr = conn.LocalAddr().(*net.TCPAddr).IP

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

func hostnameAndPort(u *JsonURL) (hostname string, port string, err error) {
	// @todo investigate net.SplitHostPort
	hostname = u.Host
 	
 	if colonInd := strings.LastIndex(hostname, ":"); colonInd > 0 {
 		// there's a port in the hostname
 		port = hostname[colonInd+1:len(hostname)]
 		
 		hostname = hostname[:colonInd]
 	} else {
		switch u.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			err =  fmt.Errorf("unknown URL scheme '%s' and no port provided", u.Scheme)
		}
	}

	return
}
