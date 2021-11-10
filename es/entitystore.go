package es

import "context"

type EntityStore interface {
	Load(ctx context.Context, entityType string, entityId string, opts ...DataLoadOption) (Entity, error)
	Save(ctx context.Context, entity Entity) error
}

type entityStore struct {
	EntityRegistry
	eventHandler EventHandler
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
func (e *entityStore) Save(ctx context.Context, entity Entity) error {
	entityOptions, err := e.EntityRegistry.GetOptions(entity.GetTypeName())
	if err != nil {
		return err
	}

	unit, err := GetUnit(ctx)
	if err != nil {
		return err
	}

	store := NewDataStore(unit.Data(), &entityOptions)
	events, err := store.Save(ctx, entity)
	if err != nil {
		return err
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
