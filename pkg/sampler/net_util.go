package sampler

import (
	"net"
	"time"
	"fmt"
	"errors"
	"crypto/tls"
)

type lookupResult struct {
	IP net.IP
	Err error
}

func resolveIPAddr(host string, deadline time.Time) (net.IP, error) {
	resultChan := make(chan *lookupResult)
	
	go func() {
		var ip net.IP
		
		ips, err := net.LookupIP(host)
		if err == nil {
			ip = ips[0]
		}
		
		resultChan <- &lookupResult{ip, err}
	}()
	
	select {
		case result := <- resultChan:
			return result.IP, result.Err
		case <- time.After(deadline.Sub(time.Now())):
			return nil, errors.New("dns timeout")
	}
}

func dial(scheme string, addr string, serverName string, deadline time.Time, insecure bool) (net.Conn, error) {
	dialer := &net.Dialer{
		Deadline: deadline,
	}

	switch scheme {
	case "http":
		return dialer.Dial("tcp", addr)
	case "https":
		return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			ServerName:         serverName,
			InsecureSkipVerify: insecure,
		})
	default:
		return nil, fmt.Errorf("unknown scheme '%s'", scheme)
	}
}
