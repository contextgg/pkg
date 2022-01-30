package session

import (
	"net/http"

	"github.com/contextgg/pkg/x"
)

func SaveSession(sessionManager Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			sess, err := sessionManager.Get(w, r)
			if err != nil {
				x.WriteError(w, err)
				return
			}

			// save it in the session
			ctx = SetSession(ctx, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
