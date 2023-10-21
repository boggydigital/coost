package coost

import (
	"github.com/boggydigital/wits"
	"net/url"
	"os"
)

func (pj persistentJar) Store(path string) error {

	hostCookies := make(wits.SectionKeyValue)

	for _, host := range pj.hosts {
		u := &url.URL{
			Scheme: httpsScheme,
			Host:   host,
		}
		hostCookies[host] = dehydrate(pj.Cookies(u))
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return hostCookies.Write(file)
}
