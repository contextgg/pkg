package logit

import (
	"context"
	"fmt"
	"time"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/logger"
)

// NewMiddleware logs stuff about the command
func NewMiddleware(l logger.Logger) es.CommandHandlerMiddleware {
	return es.CommandHandlerMiddleware(func(h es.CommandHandler) es.CommandHandler {
		return es.CommandHandlerFunc(func(ctx context.Context, cmd es.Command) error {
			start := time.Now()
			name := fmt.Sprintf("%T", cmd)
			err := h.HandleCommand(ctx, cmd)
			duration := time.Since(start)

			if err != nil {
				l.Error("Command failed", "err", err, "cmd", name, "start", start, "duration", duration)
				return err
			}
			l.Error("Command completed", "cmd", name, "start", start, "duration", duration)

			return nil
		})
	})
}
