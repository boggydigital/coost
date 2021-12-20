package cooja

import "strings"

const (
	keyValuePairsSep = ";"
	keyValueSep      = "="
)

func expandCookieHeader(cookieHeader string) map[string]string {
	kvm := make(map[string]string)

	kvps := strings.Split(cookieHeader, keyValuePairsSep)
	for _, kvp := range kvps {
		kv := strings.Split(kvp, keyValueSep)
		if len(kv) == 2 {
			key := strings.Trim(kv[0], " ")
			value := strings.Trim(kv[1], " ")
			kvm[key] = value
		}
	}

	return kvm
}
