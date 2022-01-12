package coost

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultCookieLifespan = time.Hour * 24 * 30

const (
	defaultPath      = "/"
	httpsScheme      = "https"
	cookieHeaderKey  = "cookie-header"
	keyValuePairsSep = "; "
	keyValueSep      = "="
)

func expandCookieHeader(cookieHeader string) map[string]string {
	kvm := make(map[string]string)

	kvps := strings.Split(cookieHeader, keyValuePairsSep)
	for _, kvp := range kvps {
		kv := strings.Split(kvp, keyValueSep)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			kvm[key] = val
		}
	}

	return kvm
}

func hydrate(host string, cookies map[string]string) (*url.URL, []*http.Cookie) {

	//replace cookie-header with extended values
	if content, ok := cookies[cookieHeaderKey]; ok {
		for key, value := range expandCookieHeader(content) {
			cookies[key] = value
		}
		delete(cookies, cookieHeaderKey)
	}

	cs := make([]*http.Cookie, 0, len(cookies))

	for name, value := range cookies {
		ck := &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     defaultPath,
			Domain:   host,
			Expires:  time.Now().Add(defaultCookieLifespan),
			Secure:   true,
			HttpOnly: true,
		}
		cs = append(cs, ck)
	}

	u := &url.URL{
		Scheme: httpsScheme,
		Host:   host,
	}

	return u, cs
}
