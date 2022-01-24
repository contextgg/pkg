package mustns

import (
	"context"
	"errors"

	"github.com/contextgg/pkg/auth"
	"github.com/contextgg/pkg/es"
)

var ErrNoAuth = errors.New("no auth found")

// NewMiddleware returns a new middleware that checks the ns matches against one supplied
func NewMiddleware() es.CommandHandlerMiddleware {
	return es.CommandHandlerMiddleware(func(h es.CommandHandler) es.CommandHandler {
		return es.CommandHandlerFunc(func(ctx context.Context, cmd es.Command) error {
			_, ok := auth.FromContext(ctx)
			if !ok {
				return ErrNoAuth
			}

			return h.HandleCommand(ctx, cmd)
		})
	})
}
