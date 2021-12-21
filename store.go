package coost

import (
	"encoding/json"
	"github.com/boggydigital/nod"
	"net/url"
	"os"
	"path/filepath"
)

const cookiesFilename = "cookies.json"

func (pj persistentJar) Store() error {

	if pj.directory != "" {
		if _, err := os.Stat(pj.directory); os.IsNotExist(err) {
			nod.Log("making all directories on the path: %s", pj.directory)
			if err := os.MkdirAll(pj.directory, 0755); err != nil {
				return err
			}
		}
	}

	//read all cookies from the file to avoid overwriting with only the hosts for that jar
	hostCookies, err := readHostCookies(pj.directory)
	if !os.IsNotExist(err) {
		return err
	}

	for _, host := range pj.hosts {
		u := &url.URL{
			Scheme: httpsScheme,
			Host:   host,
		}
		hostCookies[host] = dehydrate(pj.Cookies(u))
	}

	cookiesPath := filepath.Join(pj.directory, cookiesFilename)
	cookiesFile, err := os.Create(cookiesPath)
	if err != nil {
		nod.Log("error creating %s: %s", cookiesPath, err.Error())
		return err
	}

	defer cookiesFile.Close()

	return json.NewEncoder(cookiesFile).Encode(hostCookies)
}
