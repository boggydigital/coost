package coost

import (
	"net/http"
)

func dehydrate(cookies []*http.Cookie) map[string]string {
	cnv := make(map[string]string, len(cookies))
	for _, ck := range cookies {
		cnv[ck.Name] = ck.Value
	}
	return cnv
}
