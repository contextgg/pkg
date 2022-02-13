package jwks

import (
	"net/http"

	"github.com/contextgg/pkg/x"
)

type key int

const tokenKey key = 1

// NewJWTMiddleware
func NewJWTMiddleware(config JWTConfig, required bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := config.GetToken(r)
			if err != nil && required {
				x.WriteError(w, ErrJWTMissing)
				return
			}
			if err == nil {
				ctx = SetToken(ctx, token)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
