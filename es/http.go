package es

import (
	"context"
	"net/http"

	"github.com/uptrace/bun"
)

func Middleware(db *bun.DB) func(next http.Handler) http.Handler {
	uniter := NewUniter(db)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			uniter.Run(r.Context(), func(ctx context.Context) error {
				next.ServeHTTP(w, r.WithContext(ctx))
				return nil
			})
		}
		return http.HandlerFunc(fn)
	}
}
