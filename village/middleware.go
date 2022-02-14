package village

import (
	"fmt"
	"net/http"

	"github.com/contextgg/pkg/identity"
	"github.com/contextgg/pkg/jwks"
	"github.com/contextgg/pkg/ns"
	"github.com/contextgg/pkg/x"
)

func getUserFromToken(cfg *jwks.JWTConfig, r *http.Request) (*identity.User, error) {
	if cfg == nil {
		return nil, jwks.ErrJWTMissing
	}

	token, err := cfg.GetToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwks.UserClaims)
	if !ok {
		return nil, jwks.ErrJWTInvalid
	}

	// convert it.
	user := ToUser(claims)

	return user, nil
}

func NewMiddleware(userCfg *jwks.JWTConfig, apiCfg *jwks.JWTConfig, required bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			api, err := getUserFromToken(apiCfg, r)
			if err != nil && err != jwks.ErrJWTMissing {
				x.WriteError(w, err)
				return
			}
			if api != nil {
				ctx = identity.SetUser(ctx, api)
				ctx = ns.SetNamespace(ctx, api.Audience)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			user, err := getUserFromToken(userCfg, r)
			if err != nil && err != jwks.ErrJWTMissing {
				x.WriteError(w, err)
				return
			}
			if user != nil {
				ctx = identity.SetUser(ctx, user)
				ctx = ns.SetNamespace(ctx, user.Audience)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if required {
				x.WriteError(w, fmt.Errorf("Invalid"))
				return
			}

		})
	}
}
