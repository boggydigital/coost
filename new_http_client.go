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
			Proxy: nil,
			DialContext: (&net.Dialer{
				Timeout:   defaultTimeout,
				KeepAlive: defaultTimeout,
			}).DialContext,
			Dial:                   nil,
			DialTLSContext:         nil,
			DialTLS:                nil,
			TLSClientConfig:        nil,
			TLSHandshakeTimeout:    defaultTimeout,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    0,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        defaultTimeout,
			ResponseHeaderTimeout:  defaultTimeout,
			ExpectContinueTimeout:  defaultTimeout,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			GetProxyConnectHeader:  nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
			ForceAttemptHTTP2:      false,
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
