package es

import (
	"context"
	"errors"

	"github.com/contextgg/pkg/events"
)

var ErrNoRows = errors.New("No rows found")

type DataOpts struct {
	Snapshots bool
	Events    bool
	Entities  []Entity
}

func InitializeSnapshots() *DataOpts {
	return &DataOpts{
		Snapshots: true,
	}
}
func InitializeEvents() *DataOpts {
	return &DataOpts{
		Events: true,
	}
}
func InitializeEntities(entities ...Entity) *DataOpts {
	return &DataOpts{
		Entities: entities,
	}
}

// Data for all
type Data interface {
	BeginContext(ctx context.Context) (Transaction, error)
	LoadEntity(ctx context.Context, entity Entity) error
	SaveEntity(ctx context.Context, entity Entity) error
	DeleteEntry(ctx context.Context, entity Entity) error
	LoadSnapshot(ctx context.Context, rev string, agg AggregateSourced) error
	SaveSnapshot(ctx context.Context, rev string, agg AggregateSourced) error
	LoadUniqueEvents(ctx context.Context, aggregateTypeName string) ([]events.Event, error)
	LoadEventsByType(ctx context.Context, aggregateTypeName string, eventTypeNames ...string) ([]events.Event, error)
	LoadAllEvents(ctx context.Context) ([]events.Event, error)
	LoadEvents(ctx context.Context, id string, aggregateTypeName string, fromVersion int) ([]events.Event, error)
	SaveEvents(ctx context.Context, events ...events.Event) error
}

// Transaction for doing things in a transaction
type Transaction interface {
	Data

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
