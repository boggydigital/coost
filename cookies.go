package cooja

import (
	"net/http"
	"net/url"
)

func (lj persistentJar) Cookies(u *url.URL) []*http.Cookie {
	return lj.jar.Cookies(u)
}
