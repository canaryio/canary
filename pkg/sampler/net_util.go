package sampler

import (
	"fmt"
	"net"
	"net/url"
	"crypto/tls"
)

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

