package coost

import (
	"github.com/boggydigital/wits"
	"maps"
	"net/http"
	"net/http/cookiejar"
	"os"
	"slices"
)

type PersistentCookieJar interface {
	http.CookieJar
	Store(string) error
	NewHttpClient() *http.Client
}

type persistentJar struct {
	jar   http.CookieJar
	hosts []string
}

func NewJar(path string) (PersistentCookieJar, error) {
	pj := &persistentJar{}

	var err error
	pj.jar, err = cookiejar.New(nil)
	if err != nil {
		return pj, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hostCookies, err := wits.ReadSectionKeyValue(file)
	if err != nil &&
		!os.IsNotExist(err) {
		return pj, err
	}

	pj.hosts = slices.Collect(maps.Keys(hostCookies))

	for host, cookies := range hostCookies {
		pj.jar.SetCookies(hydrate(host, cookies))
	}

	return pj, nil
}
