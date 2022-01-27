package coost

import (
	"github.com/boggydigital/wits"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
)

type PersistentCookieJar interface {
	http.CookieJar
	Store(io.Writer) error
	NewHttpClient() *http.Client
}

type persistentJar struct {
	jar   http.CookieJar
	hosts []string
}

func NewJar(cookiesReader io.Reader, hosts ...string) (PersistentCookieJar, error) {
	pj := &persistentJar{
		hosts: hosts,
	}

	var err error
	pj.jar, err = cookiejar.New(nil)
	if err != nil {
		return pj, err
	}

	hostCookies, err := wits.ReadSectionKeyValue(cookiesReader)
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
