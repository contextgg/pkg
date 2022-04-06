package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/events"
)

// EventBusError is an async error containing the error returned from a handler
// or observer and the event that it happened on.
type EventBusError struct {
	Err   error
	Ctx   context.Context
	Event *events.Event
}

// Error implements the Error method of the error interface.
func (e EventBusError) Error() string {
	return fmt.Sprintf("%s: (%s)", e.Err, e.Event)
}

// EventBus for handling events
type EventBus interface {
	EventHandler
	AddHandler(EventHandler, EventMatcher)
}

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

type eventBus struct {
	handlers EventHandlers
	uniter   Uniter
}

func (s *eventBus) AddHandler(handler EventHandler, matcher EventMatcher) {
	h := UseEventHandlerMiddleware(handler, EventMatcherMiddleware(matcher))
	s.handlers = append(s.handlers, h)
}

func (b *eventBus) HandleEvent(ctx context.Context, evt events.Event) error {
	exec := func(ctx context.Context) error {
		return b.handlers.HandleEvent(ctx, evt)
	}
	return b.uniter.Run(ctx, exec)
}

// NewEventBus to handle aggregates
func NewEventBus(uniter Uniter) EventBus {
	return &eventBus{
		uniter: uniter,
	}
}
