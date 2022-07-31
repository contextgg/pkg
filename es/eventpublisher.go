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

// EventSubscriber used to listen for events on a bus of some sort
type EventSubscriber interface {
	Errors() <-chan EventBusError
	Start()
}

// EventPublisher for publishing events
type EventPublisher interface {
	// PublishEvent the event on the bus.
	PublishEvent(context.Context, events.Event) error
}

type EventPublishers []EventPublisher

func (ps EventPublishers) PublishEvents(ctx context.Context, evts []events.Event) error {
	var errs []error
	for _, p := range ps {
		for _, evt := range evts {
			if !IsPublicEvent(evt.Type) {
				continue
			}

			if err := p.PublishEvent(ctx, evt); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("Issue publishing events %v", errs)
	}
	return nil
}
