package coost

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	cookieHeaderPfx     = "Cookie:"
	cookieNameValuesSep = ";"
	cookieNameValueSep  = "="
	rootPath            = "/"
)

func newCookieJar(u *url.URL, cookies []*http.Cookie) (http.CookieJar, error) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	jar.SetCookies(u, cookies)

	return jar, nil
}

func Read(u *url.URL, path string) (http.CookieJar, error) {

	cf, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer cf.Close()

	var cookies []*http.Cookie

	if err = json.NewDecoder(cf).Decode(&cookies); err != nil {
		return nil, err
	}

	return newCookieJar(u, cookies)
}

func Write(jar http.CookieJar, u *url.URL, path string) error {

	dir, _ := filepath.Split(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0777); err != nil {
			return err
		}
	}

	cf, err := os.Create(path)
	if err != nil {
		return err
	}
	defer cf.Close()

	cookies := jar.Cookies(u)
	writeableCookies := make([]*http.Cookie, 0, len(cookies))

	for _, cookie := range cookies {
		writeableCookies = append(writeableCookies, newCookie(cookie.Name, cookie.Value, u))
	}

	return json.NewEncoder(cf).Encode(writeableCookies)
}

func newCookie(name, value string, u *url.URL) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     rootPath,
		Domain:   u.Host,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
}

func Import(cookieStr string, u *url.URL, path string, filter ...string) error {

	cookieStr = strings.TrimPrefix(cookieStr, cookieHeaderPfx)
	cookieNameValues := strings.Split(cookieStr, cookieNameValuesSep)

	if len(cookieNameValues) == 0 {
		return nil
	}

	cookies := make([]*http.Cookie, 0, len(cookieNameValues))

	for _, cnv := range cookieNameValues {
		cnv = strings.TrimSpace(cnv)
		if name, value, ok := strings.Cut(cnv, cookieNameValueSep); ok {
			if len(filter) > 0 && !slices.Contains(filter, name) {
				continue
			}
			cookies = append(cookies, newCookie(name, value, u))
		}
	}

	jar, err := newCookieJar(u, cookies)
	if err != nil {
		return err
	}

	return Write(jar, u, path)
}
