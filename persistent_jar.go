package cooja

import (
	"encoding/json"
	"log"
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

	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		return pj, nil
	}

	cookiesFile, err := os.Open(cookiePath)
	if err != nil {
		return pj, err
	}
	defer cookiesFile.Close()

	var hostCookies map[string]map[string]string
	if err := json.NewDecoder(cookiesFile).Decode(&hostCookies); err != nil {
		log.Println(err)
		return pj, nil
	}

	for host, cookies := range hostCookies {
		pj.jar.SetCookies(hydrate(host, cookies))
	}

	return pj, nil
}
