package coost

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
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

	return json.NewEncoder(cf).Encode(jar.Cookies(u))
}

func Import(cookieStr string, u *url.URL, path string) error {

	cookieStr = strings.TrimPrefix(cookieStr, cookieHeaderPfx)
	cookieNameValues := strings.Split(cookieStr, cookieNameValuesSep)

	if len(cookieNameValues) == 0 {
		return nil
	}

	cookies := make([]*http.Cookie, 0, len(cookieNameValues))

	for _, cnv := range cookieNameValues {
		cnv = strings.TrimSpace(cnv)
		if name, value, ok := strings.Cut(cnv, cookieNameValueSep); ok {
			cookie := &http.Cookie{
				Name:     name,
				Value:    value,
				Path:     rootPath,
				Domain:   u.Host,
				Expires:  time.Now().Add(30 * 24 * time.Hour),
				Secure:   true,
				HttpOnly: true,
			}
			cookies = append(cookies, cookie)
		}
	}

	jar, err := newCookieJar(u, cookies)
	if err != nil {
		return err
	}

	return Write(jar, u, path)
}
