package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// Projector takes a events and may return new commands
type Projector interface {
	Project(ctx context.Context, evt events.Event, entity Entity) (Entity, error)
}

// NewProjectorHandler turns an
func NewProjectorHandler(
	store Store,
	projector Projector,
) EventHandler {
	return &projectorHandler{store, projector}
}

type projectorHandler struct {
	store     Store
	projector Projector
}

func (p *projectorHandler) HandleEvent(ctx context.Context, evt events.Event) error {
	entity, err := p.store.Load(ctx, evt.AggregateID, false)
	if err != nil {
		return err
	}

	n, err := p.projector.Project(ctx, evt, entity)
	if err != nil {
		return err
	}

	return p.store.Save(ctx, n)
}
