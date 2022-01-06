package coost

import (
	"github.com/boggydigital/wits"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
)

type PersistentCookieJar interface {
	http.CookieJar
	Store() error
	NewHttpClient() *http.Client
}

type persistentJar struct {
	jar   http.CookieJar
	dir   string
	hosts []string
}

func NewJar(hosts []string, dir string) (PersistentCookieJar, error) {
	pj := &persistentJar{
		dir:   dir,
		hosts: hosts,
	}

	var err error
	pj.jar, err = cookiejar.New(nil)
	if err != nil {
		return pj, err
	}

	hostCookies, err := wits.Read(filepath.Join(pj.dir, cookiesFilename))
	if err != nil &&
		!os.IsNotExist(err) {
		return pj, err
	}

	for _, host := range hosts {
		if cookies, ok := hostCookies[host]; ok {
			pj.jar.SetCookies(hydrate(host, cookies))
		}
	}

	return pj, nil
}
