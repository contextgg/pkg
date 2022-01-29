package ns

import (
	"net/http"
)

func MiddlewareSubdomain(domains ...string) func(next http.Handler) http.Handler {
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

func MiddlewareHeader() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			namespace := r.Header.Get("X-Namespace")
			if len(namespace) > 0 {
				ctx = SetNamespace(ctx, namespace)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func MiddlewareJwtAud(extractor JwtExtractor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			namespace := JwtValue(r, extractor)
			if len(namespace) > 0 {
				ctx = SetNamespace(ctx, namespace)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
