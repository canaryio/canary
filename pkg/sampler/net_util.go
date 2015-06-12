package sampler

import (
	"fmt"
	"net"
	"crypto/tls"
)

func dial(t Target) (net.Conn, error) {
	host, err := hostString(&t.URL)
	if err != nil {
		return nil, err
	}

	switch t.URL.Scheme {
	case "http":
		return net.Dial("tcp", host)
	case "https":
		return tls.Dial("tcp", host, &tls.Config{
			InsecureSkipVerify: t.InsecureSkipVerify,
		})
	default:
		return nil, fmt.Errorf("unknown scheme '%s'", t.URL.Scheme)
	}
}
