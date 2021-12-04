package cooja

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
)

type PersistentCookieJar interface {
	http.CookieJar
	Save() error
	GetClient() *http.Client
}

type persistentJar struct {
	jar           http.CookieJar
	tempDirectory string
	hosts         []string
}

func NewJar(hosts []string, tempDirectory string) (PersistentCookieJar, error) {
	pj := &persistentJar{
		tempDirectory: tempDirectory,
		hosts:         hosts,
	}

	var err error
	pj.jar, err = cookiejar.New(nil)
	if err != nil {
		return pj, err
	}

	cookiePath := filepath.Join(pj.tempDirectory, cookiesFilename)
	cookiesFile, err := os.Open(cookiePath)
	if err != nil {
		return pj, err
	}
	defer cookiesFile.Close()

	var hostCookies map[string]map[string]string
	if err := json.NewDecoder(cookiesFile).Decode(&hostCookies); err != nil {
		return pj, err
	}

	for host, cookies := range hostCookies {
		pj.jar.SetCookies(hydrate(host, cookies))
	}

	return &persistentJar{}, nil
}
