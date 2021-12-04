package cooja

import (
	"net"
	"net/http"
	"time"
)

const defaultTimeout = 20 * time.Second

func (lj persistentJar) GetClient() *http.Client {
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
		Jar: lj,
	}
}
