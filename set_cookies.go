package cooja

import (
	"net/http"
	"net/url"
)

func (lj persistentJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	lj.jar.SetCookies(u, cookies)
}
