package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/events"
)

var isLocal = IsLocal()

// EventPublisher for publishing events
type EventPublisher interface {
	// PublishEvent the event on the bus.
	PublishEvent(context.Context, events.Event) error
}

// EventSubscriber used to listen for events on a bus of some sort
type EventSubscriber interface {
	Errors() <-chan EventBusError
	Start()
}

type EventPublishers []EventPublisher

func (ps EventPublishers) PublishEvents(ctx context.Context, evts []events.Event) error {
	var errs []error
	for _, p := range ps {
		for _, evt := range evts {
			if isLocal(evt) {
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
