package identity

import (
	"context"
	"net/http"
)

type Fetch func(ctx context.Context) (*User, error)

func Middleware(fn Fetch) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := fn(ctx)
			if err != nil {
				ctx = SetUser(ctx, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}