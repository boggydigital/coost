package coost

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultCookieLifespan = time.Hour * 24 * 30

const (
	httpsScheme      = "https"
	cookieHeaderKey  = "cookie-header"
	keyValuePairsSep = "; "
)

func hydrate(host string, cnv []string) (*url.URL, []*http.Cookie) {

	//replace cookie-header with extended values
	if len(cnv) == 1 {
		if name, value, ok := strings.Cut(cnv[0], nameValueSep); ok {
			if name == cookieHeaderKey {
				cnv = strings.Split(value, keyValuePairsSep)
			}
		}
	}

	cookies := make([]*http.Cookie, 0, len(cnv))

	for _, nv := range cnv {
		if name, value, ok := strings.Cut(nv, nameValueSep); ok {
			ck := &http.Cookie{
				Name:     name,
				Value:    value,
				Path:     "/",
				Domain:   host,
				Expires:  time.Now().Add(defaultCookieLifespan),
				Secure:   true,
				HttpOnly: true,
			}
			cookies = append(cookies, ck)
		}
	}

	u := &url.URL{
		Scheme: httpsScheme,
		Host:   host,
	}

	return u, cookies
}
