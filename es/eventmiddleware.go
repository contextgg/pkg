package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

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

func EventUniterMiddleware(uniter Uniter) EventHandlerMiddleware {
	return func(h EventHandler) EventHandler {
		return EventHandlerFunc(func(ctx context.Context, evt events.Event) error {
			exec := func(ctx context.Context) error {
				return h.HandleEvent(ctx, evt)
			}
			return uniter.Run(ctx, exec)
		})
	}
}
