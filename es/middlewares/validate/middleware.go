package validate

import (
	"context"

	"github.com/contextgg/pkg/es"

	"github.com/go-playground/validator/v10"
)

// NewMiddleware returns a new middleware that validate commands with its own
// validation method; `Validate() error`. Commands without the validate method
// will not be validated.
func NewMiddleware(validate *validator.Validate) es.CommandHandlerMiddleware {
	v := validate
	if v == nil {
		v = validator.New()
	}

	return es.CommandHandlerMiddleware(func(h es.CommandHandler) es.CommandHandler {
		return es.CommandHandlerFunc(func(ctx context.Context, cmd es.Command) error {
			if err := v.StructCtx(ctx, cmd); err != nil {
				return err
			}
			// Immediate command execution.
			return h.HandleCommand(ctx, cmd)
		})
	})
}
