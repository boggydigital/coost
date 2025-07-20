package coost

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

const defaultTimeout = 30 * time.Second

func (pj persistentJar) NewHttpClient() *http.Client {
	if pj.jar == nil {
		var err error
		if pj.jar, err = cookiejar.New(nil); err != nil {
			return nil
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   defaultTimeout,
				KeepAlive: defaultTimeout,
			}).DialContext,
			TLSHandshakeTimeout:   defaultTimeout,
			IdleConnTimeout:       defaultTimeout,
			ResponseHeaderTimeout: defaultTimeout,
			ExpectContinueTimeout: defaultTimeout,
		},
		Jar: pj,
	}
}

func NewHttpClientFromFile(path string) (*http.Client, error) {

	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		} else {
			pj := &persistentJar{}
			pj.jar, err = cookiejar.New(nil)
			return pj.NewHttpClient(), err
		}
	}

	pj, err := NewJar(path)
	if err != nil {
		return nil, err
	}

	return pj.NewHttpClient(), nil
}
