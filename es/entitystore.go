package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

type EntityStore interface {
	New(ctx context.Context, entityType string, entityId string) (Entity, error)
	Load(ctx context.Context, entityType string, entityId string, opts ...DataLoadOption) (Entity, error)
	Save(ctx context.Context, entities ...Entity) error
}

type entityStore struct {
	EntityRegistry
	eventHandler EventHandler
}

func (e *entityStore) New(ctx context.Context, entityType string, entityId string) (Entity, error) {
	entityOptions, err := e.EntityRegistry.GetOptions(entityType)
	if err != nil {
		return nil, err
	}

	// make the aggregate
	entity := entityOptions.Factory(entityId)
	return entity, nil
}
func (e *entityStore) Load(ctx context.Context, entityType string, entityId string, opts ...DataLoadOption) (Entity, error) {
	entityOptions, err := e.EntityRegistry.GetOptions(entityType)
	if err != nil {
		return nil, err
	}

	unit, err := GetUnit(ctx)
	if err != nil {
		return nil, err
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
