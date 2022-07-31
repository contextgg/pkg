package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// NewSagaHandler Saga event handler
func NewSagaHandler(cli Client, saga Saga) EventHandler {
	return &sagaHandler{
		cli:  cli,
		saga: saga,
	}
}

type sagaHandler struct {
	cli  Client
	saga Saga
}

func (s *sagaHandler) HandleEvent(ctx context.Context, evt events.Event) error {
	cmds, err := s.saga.Run(ctx, evt)
	if err != nil {
		return err
	}

	return s.cli.HandleCommands(ctx, cmds...)
}
