package mustns

import (
	"context"
	"errors"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/ns"
)

var ErrNamespaceMismatch = errors.New("command not supported")

// NewMiddleware returns a new middleware that checks the ns matches against one supplied
func NewMiddleware(namespaces ...string) es.CommandHandlerMiddleware {

	return es.CommandHandlerMiddleware(func(h es.CommandHandler) es.CommandHandler {
		return es.CommandHandlerFunc(func(ctx context.Context, cmd es.Command) error {
			namespace := ns.FromContext(ctx)
			for _, item := range namespaces {
				if namespace == item {
					return h.HandleCommand(ctx, cmd)
				}
			}

			return ErrNamespaceMismatch
		})
	})
}
