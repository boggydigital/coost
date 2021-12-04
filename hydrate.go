package cooja

import (
	"net/http"
	"net/url"
	"time"
)

var defaultCookieLifespan = time.Hour * 24 * 30

const httpsScheme = "https"

func hydrate(host string, cookieValues map[string]string) (*url.URL, []*http.Cookie) {

	cookies := make([]*http.Cookie, 0, len(cookieValues))

	for cookie, value := range cookieValues {
		cookie := &http.Cookie{
			Name:     cookie,
			Value:    value,
			Path:     "/",
			Domain:   "." + host,
			Expires:  time.Now().Add(defaultCookieLifespan),
			Secure:   true,
			HttpOnly: true,
		}
		cookies = append(cookies, cookie)
	}

	u := &url.URL{
		Scheme: httpsScheme,
		Host:   host,
	}

	return u, cookies
}