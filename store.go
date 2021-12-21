package pesco

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
)

const cookiesFilename = "cookies.json"

func (pj persistentJar) Store() error {

	if pj.directory != "" {
		if _, err := os.Stat(pj.directory); os.IsNotExist(err) {
			if err := os.MkdirAll(pj.directory, 0755); err != nil {
				return err
			}
		}
	}

	cookiesPath := filepath.Join(pj.directory, cookiesFilename)
	cookiesFile, err := os.Create(cookiesPath)
	if err != nil {
		return err
	}

	defer cookiesFile.Close()

	hostCookies := make(map[string]map[string]string)
	for _, host := range pj.hosts {
		u := &url.URL{
			Scheme: httpsScheme,
			Host:   host,
		}
		hostCookies[host] = dehydrate(pj.Cookies(u))
	}

	return json.NewEncoder(cookiesFile).Encode(hostCookies)
}
