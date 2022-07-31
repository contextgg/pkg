package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// EventHandler for handling commands
type EventHandler interface {
	HandleEvent(context.Context, events.Event) error
}

// EventHandlerFunc is a function that can be used as a event handler.
type EventHandlerFunc func(context.Context, events.Event) error

// HandleEvent implements the HandleEvent method of the EventHandler.
func (h EventHandlerFunc) HandleEvent(ctx context.Context, evt events.Event) error {
	return h(ctx, evt)
}

type EventHandlers []EventHandler

func (h EventHandlers) HandleEvent(ctx context.Context, evt events.Event) error {
	for _, eh := range h {
		if err := eh.HandleEvent(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}
