package ns

import (
	"net/http"
	"strings"
)

func Slug(r *http.Request, suffixes ...string) string {
	wport := r.Host
	parts := strings.Split(wport, ":")
	hostname := parts[0]

	for _, s := range suffixes {
		if strings.HasSuffix(hostname, s) {
			hostname = hostname[:len(hostname)-len(s)]
			break
		}
	}

	hostname = strings.ReplaceAll(hostname, ".", "")
	return hostname
}
