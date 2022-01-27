package coost

import (
	"github.com/boggydigital/wits"
	"io"
	"net/url"
)

func (pj persistentJar) Store(cookieWriter io.Writer) error {

	hostCookies := make(wits.SectionKeyValue)

	for _, host := range pj.hosts {
		u := &url.URL{
			Scheme: httpsScheme,
			Host:   host,
		}
		hostCookies[host] = dehydrate(pj.Cookies(u))
	}

	return hostCookies.Write(cookieWriter)
}
