package ns

import (
	"net/http"
)

func Middleware(domains ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			slug := Slug(r, domains...)

			ctx := r.Context()
			ctx = SetNamespace(ctx, slug)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
