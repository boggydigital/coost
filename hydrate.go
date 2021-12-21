package cooja

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultCookieLifespan = time.Hour * 24 * 30

const (
	httpsScheme     = "https"
	cookieHeaderKey = "cookie-header"
)

func hydrate(host string, cookieValues map[string]string) (*url.URL, []*http.Cookie) {

	//replace cookie-header with extended values
	if content, ok := cookieValues[cookieHeaderKey]; ok {
		for key, value := range expandCookieHeader(content) {
			cookieValues[key] = value
		}
		delete(cookieValues, cookieHeaderKey)
	}

	if !strings.HasPrefix(host, ".") {
		host = "." + host
	}

	cookies := make([]*http.Cookie, 0, len(cookieValues))

	for cookie, value := range cookieValues {
		ck := &http.Cookie{
			Name:     cookie,
			Value:    value,
			Path:     "/",
			Domain:   host,
			Expires:  time.Now().Add(defaultCookieLifespan),
			Secure:   true,
			HttpOnly: true,
		}
		cookies = append(cookies, ck)
	}

	u := &url.URL{
		Scheme: httpsScheme,
		Host:   host,
	}

	return u, cookies
}
