package pesco

import (
	"net"
	"net/http"
	"time"
)

const defaultTimeout = 20 * time.Second

func (pj persistentJar) NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   defaultTimeout,
				KeepAlive: defaultTimeout,
			}).DialContext,
			TLSHandshakeTimeout:   defaultTimeout,
			ExpectContinueTimeout: defaultTimeout,
			ResponseHeaderTimeout: defaultTimeout,
			IdleConnTimeout:       defaultTimeout,
		},
		Jar: pj,
	}
}
