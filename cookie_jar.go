package pesco

import (
	"net/http"
	"net/url"
)

func (pj persistentJar) Cookies(u *url.URL) []*http.Cookie {
	return pj.jar.Cookies(u)
}

func (pj persistentJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	pj.jar.SetCookies(u, cookies)
}
