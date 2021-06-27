package x

import (
	"net/http"
)

func UserMiddleware(jwksClient JwksClient) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwksClient.ParseFromRequestWithClaims(r)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}

			user := &User{
				Id: claims.Subject,
			}
			ctx := SetUser(r.Context(), user)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
