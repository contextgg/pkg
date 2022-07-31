package es

import (
	"context"
)

// CommandHandler for handling commands
type CommandHandler interface {
	HandleCommand(context.Context, Command) error
}

// CommandHandlerFunc is a function that can be used as a command handler.
type CommandHandlerFunc func(context.Context, Command) error

// HandleCommand implements the HandleCommand method of the CommandHandler.
func (h CommandHandlerFunc) HandleCommand(ctx context.Context, cmd Command) error {
	return h(ctx, cmd)
}
