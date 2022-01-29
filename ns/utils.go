package ns

import (
	"net/http"
	"strings"

	"github.com/contextgg/pkg/jwks"
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

	return hostname
}

type JwtExtractor func(interface{}) string

func JwtValue(r *http.Request, extract JwtExtractor) string {
	c := jwks.ClaimsFromContext(r.Context())
	if c == nil {
		return ""
	}
	return extract(c)
}
