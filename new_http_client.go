package coost

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

const defaultTimeout = 20 * time.Second

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
			ExpectContinueTimeout: defaultTimeout,
			ResponseHeaderTimeout: defaultTimeout,
			IdleConnTimeout:       defaultTimeout,
		},
		Jar: pj,
	}
}

func NewHttpClientFromFile(path string, hosts ...string) (*http.Client, error) {

	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		} else {
			pj := &persistentJar{}
			pj.jar, err = cookiejar.New(nil)
			return pj.NewHttpClient(), err
		}
	}

	cookieFile, err := os.Open(path)
	defer cookieFile.Close()
	if err != nil {
		return nil, err
	}

	pj, err := NewJar(cookieFile, hosts...)
	if err != nil {
		return nil, err
	}

	return pj.NewHttpClient(), nil
}
