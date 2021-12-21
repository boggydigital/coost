package pesco

import "net/http"

func dehydrate(cookies []*http.Cookie) map[string]string {
	ckv := make(map[string]string, len(cookies))
	for _, ck := range cookies {
		ckv[ck.Name] = ck.Value
	}
	return ckv
}
