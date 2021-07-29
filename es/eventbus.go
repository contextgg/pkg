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
	AddPublisher(string, EventPublisher)
	Close()
}

// EventHandler for handling commands
type EventHandler interface {
	HandleEvent(context.Context, events.Event) error
}

// EventPublisher for publishing events
type EventPublisher interface {
	// PublishEvent the event on the bus.
	PublishEvent(context.Context, events.Event) error
	Errors() <-chan EventBusError
	Start()
	Close()
}

// EventHandlerFunc is a function that can be used as a event handler.
type EventHandlerFunc func(context.Context, events.Event) error

// HandleEvent implements the HandleEvent method of the EventHandler.
func (h EventHandlerFunc) HandleEvent(ctx context.Context, evt events.Event) error {
	return h(ctx, evt)
}

// EventHandlerMiddleware is a function that middlewares can implement to be able to chain.
type EventHandlerMiddleware func(EventHandler) EventHandler

// UseEventHandlerMiddleware wraps a EventHandler in one or more middleware.
func UseEventHandlerMiddleware(h EventHandler, middleware ...EventHandlerMiddleware) EventHandler {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}
	return h
}

func EventMatcherMiddleware(matcher EventMatcher) EventHandlerMiddleware {
	return func(h EventHandler) EventHandler {
		return EventHandlerFunc(func(ctx context.Context, evt events.Event) error {
			if !matcher(evt) {
				return nil
			}
			return h.HandleEvent(ctx, evt)
		})
	}
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

// NewEventBus to handle aggregates
func NewEventBus() EventBus {
	return &eventBus{
		publishers: make(map[string]EventPublisher),
	}
}

type eventBus struct {
	handlers   EventHandlers
	publishers map[string]EventPublisher
}

func (s *eventBus) AddHandler(handler EventHandler, matcher EventMatcher) {
	h := UseEventHandlerMiddleware(handler, EventMatcherMiddleware(matcher))
	s.handlers = append(s.handlers, h)
}

func (b *eventBus) AddPublisher(key string, publisher EventPublisher) {
	b.publishers[key] = publisher
}

func (b *eventBus) HandleEvent(ctx context.Context, evt events.Event) error {
	if err := b.handlers.HandleEvent(ctx, evt); err != nil {
		return err
	}

	notLocal := MatchNotLocal()
	notPublisher := MatchNotPublisher()
	if !notLocal(evt) || !notPublisher(evt) {
		return nil
	}

	for _, p := range b.publishers {
		if err := p.PublishEvent(ctx, evt); err != nil {
			return err
		}
	}

	return nil
}

// Close underlying connection
func (b *eventBus) Close() {
	for _, p := range b.publishers {
		p.Close()
	}
}
