package es

import (
	"context"
	"errors"
	"fmt"

	"github.com/contextgg/pkg/events"
)

type StoreOpts func(s Store)

func StoreRevision(revision string) StoreOpts {
	return func(s Store) {
		if imp, ok := s.(*store); ok {
			imp.revision = revision
		}
	}
}
func StoreRevisionMin(minVersionDiff int) StoreOpts {
	return func(s Store) {
		if imp, ok := s.(*store); ok {
			imp.minVersionDiff = minVersionDiff
		}
	}
}
func StoreDisableRevision() StoreOpts {
	return func(s Store) {
		if imp, ok := s.(*store); ok {
			imp.minVersionDiff = -1
		}
	}
}
func StoreDisableProject() StoreOpts {
	return func(s Store) {
		if imp, ok := s.(*store); ok {
			imp.project = false
		}
	}
}

func applyEvents(ctx context.Context, aggregate AggregateSourced, originalEvents []events.Event) error {
	aggregateType := aggregate.GetTypeName()

	for _, event := range originalEvents {
		if event.AggregateType != aggregateType {
			return ErrMismatchedEventType
		}

		// lets build the event!
		if err := aggregate.ApplyEvent(ctx, event); err != nil {
			return ApplyEventError{
				Event: event,
				Err:   err,
			}
		}
		aggregate.IncrementVersion()
	}
	return nil
}

type Store interface {
	Load(ctx context.Context, id string, forced bool) (Entity, error)
	Save(ctx context.Context, aggregate Entity) error
	Delete(ctx context.Context, aggregate Entity) error
	RunInTransaction(ctx context.Context, fn func(Store) error) error
}

type store struct {
	data         Data
	eventHandler EventHandler
	factory      EntityFunc

	revision       string
	minVersionDiff int
	project        bool
}

func (s *store) loadSourced(ctx context.Context, aggregate AggregateSourced, forced bool) (Entity, error) {
	// load up the aggregate
	if s.minVersionDiff >= 0 && !forced {
		if err := s.data.LoadSnapshot(ctx, s.revision, aggregate); err != nil {
			return nil, err
		}
	}
	// load up the events from the DB.
	originalEvents, err := s.data.LoadEvents(ctx, aggregate.GetID(), aggregate.GetTypeName(), aggregate.GetVersion())
	if err != nil {
		return nil, err
	}
	if err := applyEvents(ctx, aggregate, originalEvents); err != nil {
		return nil, err
	}
	return aggregate, nil
}
func (s *store) loadEntity(ctx context.Context, entity Entity, forced bool) (Entity, error) {
	if err := s.data.LoadEntity(ctx, entity); err != nil && !errors.Is(err, ErrNoRows) {
		return nil, err
	}
	return entity, nil
}
func (s *store) saveSourced(ctx context.Context, aggregate AggregateSourced) error {
	originalVersion := aggregate.GetVersion()

	// now save it!.
	events := aggregate.Events()
	if len(events) > 0 {
		if err := s.data.SaveEvents(ctx, events...); err != nil {
			return err
		}
		aggregate.ClearEvents()

		// Apply the events so we can save the aggregate
		if err := applyEvents(ctx, aggregate, events); err != nil {
			return err
		}
	}

	// save the snapshot!
	diff := aggregate.GetVersion() - originalVersion
	if diff < 0 {
		return ErrWrongVersion
	}

	if s.minVersionDiff >= 0 && diff >= s.minVersionDiff {
		if err := s.data.SaveSnapshot(ctx, s.revision, aggregate); err != nil {
			return err
		}
	}

	if s.project {
		if err := s.data.SaveEntity(ctx, aggregate); err != nil {
			return err
		}
	}

	if s.eventHandler != nil {
		for _, e := range events {
			if err := s.eventHandler.HandleEvent(ctx, e); err != nil {
				return err
			}
		}
	}

	return nil
}
func (s *store) saveEntity(ctx context.Context, aggregate Entity) error {
	return s.data.SaveEntity(ctx, aggregate)
}
func (s *store) saveAggregateHolder(ctx context.Context, aggregate AggregateHolder) error {
	if err := s.data.SaveEntity(ctx, aggregate); err != nil {
		return err
	}

	events := aggregate.EventsToPublish()
	aggregate.ClearEvents()

	for _, e := range events {
		if err := s.eventHandler.HandleEvent(ctx, e); err != nil {
			return err
		}
	}

	return nil
}

func (s *store) deleteSourced(ctx context.Context, aggregate AggregateSourced) error {
	return fmt.Errorf("Delete sourced %s - %s", aggregate.GetID(), aggregate.GetTypeName())
}
func (s *store) deleteEntity(ctx context.Context, aggregate Entity) error {
	return s.data.DeleteEntry(ctx, aggregate)
}

func (s *store) Load(ctx context.Context, id string, forced bool) (Entity, error) {
	// make the aggregate
	entity := s.factory(id)

	switch agg := entity.(type) {
	case AggregateSourced:
		return s.loadSourced(ctx, agg, forced)
	default:
		return s.loadEntity(ctx, agg, forced)
	}
}
func (s *store) Save(ctx context.Context, entity Entity) error {
	switch agg := entity.(type) {
	case AggregateSourced:
		return s.saveSourced(ctx, agg)
	case AggregateHolder:
		return s.saveAggregateHolder(ctx, agg)
	default:
		return s.saveEntity(ctx, agg)
	}
}
func (s *store) Delete(ctx context.Context, entity Entity) error {
	switch agg := entity.(type) {
	case AggregateSourced:
		return s.deleteSourced(ctx, agg)
	default:
		return s.deleteEntity(ctx, agg)
	}
}

func (s *store) RunInTransaction(ctx context.Context, fn func(Store) error) error {
	tx, err := s.data.BeginContext(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// store ..
	imp := &store{
		data:         tx,
		eventHandler: s.eventHandler,
		factory:      s.factory,

		revision:       s.revision,
		minVersionDiff: s.minVersionDiff,
		project:        s.project,
	}

	if err := fn(imp); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// NewStore for creating stores
func NewStore(data Data, eventHandler EventHandler, factory EntityFunc, opts ...StoreOpts) Store {
	s := &store{
		data:           data,
		eventHandler:   eventHandler,
		factory:        factory,
		revision:       "rev1",
		minVersionDiff: 0,
		project:        true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
