package es2

import "context"

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}
