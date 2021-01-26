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

// CommandBus for creating commands
type CommandBus interface {
	CommandRegistry
	CommandHandler
}

// NewCommandBus create a new bus from a registry
func NewCommandBus() CommandBus {
	return &commandBus{
		NewCommandRegistry(),
	}
}

type commandBus struct {
	CommandRegistry
}

func (b *commandBus) HandleCommand(ctx context.Context, cmd Command) error {
	// find the handler!
	handler, err := b.GetHandler(cmd)
	if err != nil {
		return err
	}
	return handler.HandleCommand(ctx, cmd)
}
