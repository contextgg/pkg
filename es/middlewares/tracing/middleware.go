package tracing

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/types"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// NewMiddleware returns a new command handler middleware that adds tracing spans.
func NewMiddleware() es.CommandHandlerMiddleware {
	return es.CommandHandlerMiddleware(func(h es.CommandHandler) es.CommandHandler {
		return es.CommandHandlerFunc(func(ctx context.Context, cmd es.Command) error {
			name := types.GetTypeName(cmd)
			opName := fmt.Sprintf("Command(%s)", name)
			sp, ctx := opentracing.StartSpanFromContext(ctx, opName)

			err := h.HandleCommand(ctx, cmd)

			sp.SetTag("es.aggregate_id", cmd.GetAggregateId())
			if err != nil {
				ext.LogError(sp, err)
			}
			sp.Finish()

			return err
		})
	})
}
