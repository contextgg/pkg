package es2

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type CommandGroup interface {
	Command | []Command
}

type CommandBus interface {
	Dispatch(ctx context.Context, cmds ...CommandGroup) error
}

type commandBus struct {
	CommandRegistry
}

func (b *commandBus) handle(ctx context.Context, cmds ...Command) error {
	for _, cmd := range cmds {
		ctx, span := tracer.Start(ctx, "Handle")
		defer span.End()

		// find the handler!
		handler, err := b.GetHandler(cmd)
		if err != nil {
			return err
		}

		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (b *commandBus) Dispatch(ctx context.Context, cmds ...CommandGroup) error {
	ctx, span := tracer.Start(ctx, "Dispatch")
	defer span.End()

	errs, ctx := errgroup.WithContext(ctx)
	for _, cmd := range cmds {
		var all []Command

		switch c := cmd.(type) {
		case Command:
			all = []Command{c}
		case []Command:
			all = c
		}

		errs.Go(func() error {
			return b.handle(ctx, all...)
		})
	}

	return errs.Wait()
}

func NewCommandBus(commandRegistry CommandRegistry) CommandBus {
	return &commandBus{
		CommandRegistry: commandRegistry,
	}
}
