package coost

import (
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"net/url"
	"os"
	"path/filepath"
)

const cookiesFilename = "cookies.txt"

func (pj persistentJar) Store() error {

	if pj.dir != "" {
		if _, err := os.Stat(pj.dir); os.IsNotExist(err) {
			nod.Log("making all directories on the path: %s", pj.dir)
			if err := os.MkdirAll(pj.dir, 0755); err != nil {
				return err
			}
		}
	}

	//read all cookies from the file to avoid overwriting with only the hosts for that jar
	hostCookies, err := wits.Read(filepath.Join(pj.dir, cookiesFilename))
	if err != nil {
		return err
	}

	for _, host := range pj.hosts {
		u := &url.URL{
			Scheme: httpsScheme,
			Host:   host,
		}
		hostCookies[host] = dehydrate(pj.Cookies(u))
	}

	cookiesPath := filepath.Join(pj.dir, cookiesFilename)
	cookiesFile, err := os.Create(cookiesPath)
	if err != nil {
		nod.Log("error creating %s: %s", cookiesPath, err.Error())
		return err
	}

	defer cookiesFile.Close()

	return wits.Write(hostCookies, cookiesPath)
}
