package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// NewSagaHandler Saga event handler
func NewSagaHandler(handler CommandHandler, saga Saga) EventHandler {
	return &sagaHandler{
		handler: handler,
		saga:    saga,
	}
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
