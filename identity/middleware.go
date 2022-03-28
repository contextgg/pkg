package identity

import (
	"net/http"

	"github.com/contextgg/pkg/ns"
)

type Fetch func(r *http.Request) (*User, error)

func Middleware(fn Fetch) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := fn(r)
			if err != nil {
				ctx = SetUser(ctx, user)
				ctx = ns.SetNamespace(ctx, user.Audience)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
