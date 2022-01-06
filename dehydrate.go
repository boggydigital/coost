package coost

import (
	"fmt"
	"net/http"
)

const nameValueSep = "="

func dehydrate(cookies []*http.Cookie) []string {
	cnv := make([]string, 0, len(cookies))
	for _, ck := range cookies {
		cnv = append(cnv,
			fmt.Sprintf("%s%s%s", ck.Name, nameValueSep, ck.Value))
	}
	return cnv
}
