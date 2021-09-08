package es

import (
	"context"
	"errors"

	"github.com/contextgg/pkg/events"
)

var ErrNoRows = errors.New("No rows found")

type DataOption func(d *DataOpts)

type DataOpts struct {
	RecreateTables bool
	TruncateTables bool
	HasEvents      bool
	HasSnapshots   bool
	ExtraModels    []interface{}
	Migrations     []interface{}
}

func InitializeSnapshots() DataOption {
	return func(o *DataOpts) {
		o.HasSnapshots = true
	}
}
func InitializeEvents() DataOption {
	return func(o *DataOpts) {
		o.HasEvents = true
	}
}
func InitializeEntities(entities ...Entity) DataOption {
	return func(o *DataOpts) {
		for _, ent := range entities {
			o.ExtraModels = append(o.ExtraModels, ent)
		}
	}
}
func WithTruncate() DataOption {
	return func(o *DataOpts) {
		o.TruncateTables = true
	}
}
func WithRecreate() DataOption {
	return func(o *DataOpts) {
		o.RecreateTables = true
	}
}
func RegisterMigrations(migration ...interface{}) DataOption {
	return func(o *DataOpts) {
		for _, m := range migration {
			o.Migrations = append(o.Migrations, m)
		}
	}
}

// Data for all
type Data interface {
	LoadEntity(ctx context.Context, namespace string, entity Entity) error
	SaveEntity(ctx context.Context, namespace string, entity Entity) error
	DeleteEntry(ctx context.Context, namespace string, entity Entity) error
	LoadSnapshot(ctx context.Context, namespace string, rev string, agg AggregateSourced) error
	SaveSnapshot(ctx context.Context, namespace string, rev string, agg AggregateSourced) error
	LoadUniqueEvents(ctx context.Context, namespace string, aggregateTypeName string) ([]events.Event, error)
	LoadEventsByType(ctx context.Context, namespace string, aggregateTypeName string, eventTypeNames ...string) ([]events.Event, error)
	LoadAllEvents(ctx context.Context, namespace string) ([]events.Event, error)
	LoadEvent(ctx context.Context, namespace string, id string, aggregateTypeName string, version int) (*events.Event, error)
	LoadEvents(ctx context.Context, namespace string, id string, aggregateTypeName string, fromVersion int) ([]events.Event, error)
	SaveEvents(ctx context.Context, namespace string, events ...events.Event) error
}

// Transaction for doing things in a transaction
type Transaction interface {
	Data

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
