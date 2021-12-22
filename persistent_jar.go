package coost

import (
	"encoding/json"
	"github.com/boggydigital/nod"
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
	jar       http.CookieJar
	directory string
	hosts     []string
}

func readHostCookies(directory string) (map[string]map[string]string, error) {

	hostCookies := make(map[string]map[string]string)

	cookiePath := filepath.Join(directory, cookiesFilename)
	nod.Log("reading host cookies from %s", cookiePath)

	if _, err := os.Stat(cookiePath); err != nil {
		nod.Log("error getting %s stat: %s", cookiePath, err.Error())
		return hostCookies, nil
	}

	cookiesFile, err := os.Open(cookiePath)
	if err != nil {
		nod.Log("error opening %s: %s", cookiePath, err.Error())
		return hostCookies, err
	}
	defer cookiesFile.Close()

	err = json.NewDecoder(cookiesFile).Decode(&hostCookies)

	return hostCookies, err
}

func NewJar(hosts []string, dir string) (PersistentCookieJar, error) {
	pj := &persistentJar{
		directory: dir,
		hosts:     hosts,
	}

	var err error
	pj.jar, err = cookiejar.New(nil)
	if err != nil {
		return pj, err
	}

	hostCookies, err := readHostCookies(pj.directory)
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
