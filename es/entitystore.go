package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/events"
)

type EntityStore interface {
	Load(ctx context.Context, entityType string, entityId string, opts ...DataLoadOption) (Entity, error)
	Save(ctx context.Context, entities ...Entity) error
}

type entityStore struct {
	eventHandler EventHandler
}

func (e *entityStore) Load(ctx context.Context, entityType string, entityId string, opts ...DataLoadOption) (Entity, error) {
	unit := UnitFromContext(ctx)
	if unit == nil {
		return nil, fmt.Errorf("unit not found")
	}

	store := NewDataStore(unit.Data(), &entityOptions)
	return store.Load(ctx, entityId)
}
func (e *entityStore) Save(ctx context.Context, entities ...Entity) error {
	unit, err := GetUnit(ctx)
	if err != nil {
		return err
	}

	var events []events.Event
	for _, entity := range entities {
		entityOptions, err := e.EntityRegistry.GetOptions(entity.GetTypeName())
		if err != nil {
			return err
		}
		store := NewDataStore(unit.Data(), &entityOptions)
		out, err := store.Save(ctx, entity)
		if err != nil {
			return err
		}
		events = append(events, out...)
	}

	unit.StoreEvents(events...)

	// run the event handler!.
	for _, evt := range events {
		if err := e.eventHandler.HandleEvent(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

func NewEntityStore(registry EntityRegistry, eventHandler EventHandler) EntityStore {
	return &entityStore{
		EntityRegistry: registry,
		eventHandler:   eventHandler,
	}
}
