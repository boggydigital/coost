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

// FIXME: replace with strings.Cut when 1.18 releases
func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func hydrate(host string, cnv []string) (*url.URL, []*http.Cookie) {

	//replace cookie-header with extended values
	if len(cnv) == 1 {
		if name, value, ok := cut(cnv[0], nameValueSep); ok {
			if name == cookieHeaderKey {
				cnv = strings.Split(value, keyValuePairsSep)
			}
		}
	}

	cookies := make([]*http.Cookie, 0, len(cnv))

	for _, nv := range cnv {
		if name, value, ok := cut(nv, nameValueSep); ok {
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
