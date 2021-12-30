package auth

import (
	"net/http"
)

func Middleware(mgr Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := mgr.GetAuthUser(r)
			if err != nil {
				ctx = SetAuth(ctx, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
