package es

import (
	"context"
	"errors"

	"github.com/contextgg/pkg/events"
)

// Projector takes a events and may return new commands
type Projector interface {
	Project(ctx context.Context, evt events.Event, entity Entity) (Entity, error)
}

// NewProjectorHandler turns an
func NewProjectorHandler(
	factory EntityFunc,
	data Data,
	projector Projector,
) EventHandler {
	return &projectorHandler{factory, data, projector}
}

type projectorHandler struct {
	factory   EntityFunc
	data      Data
	projector Projector
}

func (p *projectorHandler) HandleEvent(ctx context.Context, evt events.Event) error {
	// load up the entity
	entity := p.factory(evt.AggregateID)
	if err := p.data.LoadEntity(ctx, entity); err != nil && !errors.Is(err, ErrNoRows) {
		return err
	}

	n, err := p.projector.Project(ctx, evt, entity)
	if err != nil {
		return err
	}

	return p.data.SaveEntity(ctx, n)
}
