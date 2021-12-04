package cooja

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
)

const cookiesFilename = "cookies.json"

func (lj persistentJar) Save() error {

	if _, err := os.Stat(lj.tempDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(lj.tempDirectory, 0755); err != nil {
			return err
		}
	}

	cookiesPath := filepath.Join(lj.tempDirectory, cookiesFilename)
	cookiesFile, err := os.Create(cookiesPath)
	if err != nil {
		return err
	}

	defer cookiesFile.Close()

	hostCookies := make(map[string]map[string]string)
	for _, host := range lj.hosts {
		u := &url.URL{
			Scheme: httpsScheme,
			Host:   host,
		}
		hostCookies[host] = dehydrate(lj.Cookies(u))
	}

	return json.NewEncoder(cookiesFile).Encode(hostCookies)
}
