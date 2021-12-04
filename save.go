package cooja

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
)

const cookiesFilename = "cookies.json"

func (lj persistentJar) Save() error {
	cookiesFile, err := os.Create(filepath.Join(lj.tempDirectory, cookiesFilename))
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
