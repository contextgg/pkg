package session

import (
	"net/http"

	"github.com/contextgg/pkg/x"
)

func SaveSession(sessionManager Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if _, err := sessionManager.Get(r); err == NotFound {
				if err := sessionManager.Save(w, sessionManager.New()); err != nil {
					x.WriteError(w, err)
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func ValidateCsrfToken(csrfStore CsrfService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// get the header.
			token := r.Header.Get("X-CSRF-Token")
			if len(token) > 0 {
				if err := csrfStore.Verify(r, token); err != nil {
					x.WriteError(w, err)
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
