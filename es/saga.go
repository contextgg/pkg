package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// Saga takes a events and may return new commands
type Saga interface {
	Run(ctx context.Context, evt events.Event) ([]Command, error)
}

// NewSagaHandler turns an
func NewSagaHandler(handler CommandBus, saga Saga) EventHandler {
	return &sagaHandler{handler, saga}
}

type sagaHandler struct {
	handler CommandHandler
	saga    Saga
}

func (s *sagaHandler) HandleEvent(ctx context.Context, evt events.Event) error {
	cmds, err := s.saga.Run(ctx, evt)
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		if err := s.handler.HandleCommand(ctx, cmd); err != nil {
			return err
		}
	}

	return nil
}
