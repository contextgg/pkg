package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-pg/pg/v10/orm"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/pgdb"
	"github.com/contextgg/pkg/types"
)

type event struct {
	AggregateID   string `pg:",pk,type:uuid"`
	AggregateType string `pg:",pk"`
	Version       int    `pg:",pk"`
	Type          string `pg:",notnull"`
	Timestamp     time.Time
	Data          json.RawMessage
	Metadata      map[string]interface{}
}
type snapshot struct {
	ID        string          `pg:",pk,type:uuid"`
	Type      string          `pg:",pk"`
	Revision  string          `pg:",pk"`
	Aggregate json.RawMessage `pg:",notnull"`
}

func isNoRow(err error) bool {
	return err.Error() == "pg: no rows in result set"
}

func run(db pgdb.DB, opts *es.DataOpts) error {
	tableOpts := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	ctx := context.Background()
	if opts.Events {
		if err := db.ModelContext(ctx, &event{}).CreateTable(tableOpts); err != nil {
			return err
		}
	}
	if opts.Snapshots {
		if err := db.ModelContext(ctx, &snapshot{}).CreateTable(tableOpts); err != nil {
			return err
		}
	}
	for _, model := range opts.Entities {
		if err := db.ModelContext(ctx, model).CreateTable(tableOpts); err != nil {
			return err
		}
	}
	return nil
}

func NewPostgresData(db pgdb.DB, opts ...*es.DataOpts) es.Data {
	return NewPostgresDataWithLegacy(db, false, opts...)
}
func NewPostgresDataWithLegacy(db pgdb.DB, legacy bool, opts ...*es.DataOpts) es.Data {
	for _, o := range opts {
		run(db, o)
	}

	return &data{
		db:     db,
		legacy: legacy,
	}
}

type data struct {
	db     pgdb.DB
	legacy bool
}

func (s *data) BeginContext(ctx context.Context) (es.Transaction, error) {
	_, err := s.db.BeginContext(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *data) Commit(ctx context.Context) error {
	return s.db.CommitContext(ctx)
}
func (s *data) Rollback(ctx context.Context) error {
	return s.db.RollbackContext(ctx)
}

func (s *data) SaveEntity(ctx context.Context, entity es.Entity) error {
	_, err := s.db.ModelContext(ctx, entity).
		OnConflict("(id) DO UPDATE").
		Insert()
	return err
}
func (s *data) DeleteEntry(ctx context.Context, entity es.Entity) error {
	_, err := s.db.ModelContext(ctx, entity).
		WherePK().
		Delete()
	return err
}

func (s *data) SaveSnapshot(ctx context.Context, rev string, agg es.AggregateSourced) error {
	// to json!
	data, err := json.Marshal(agg)
	if err != nil {
		return err
	}

	ss := &snapshot{
		ID:        agg.GetID(),
		Type:      agg.GetTypeName(),
		Revision:  rev,
		Aggregate: json.RawMessage(data),
	}
	if _, err := s.db.ModelContext(ctx, ss).
		OnConflict("(id,type,revision) DO UPDATE").
		Insert(); err != nil {
		return err
	}
	return nil
}
func (s *data) SaveEvents(ctx context.Context, evts ...events.Event) error {
	all := make([]event, len(evts))
	for i, evt := range evts {
		data, err := types.Marshal(evt.Data, s.legacy)
		if err != nil {
			return err
		}

		m := event{
			AggregateID:   evt.AggregateID,
			AggregateType: evt.AggregateType,
			Version:       evt.Version,
			Type:          evt.Type,
			Timestamp:     evt.Timestamp,
			Data:          data,
			Metadata:      evt.Metadata,
		}
		all[i] = m
	}

	// save em
	if _, err := s.db.ModelContext(ctx, &all).
		OnConflict("(aggregate_id, aggregate_type, version) DO UPDATE").
		Insert(); err != nil {
		return err
	}
	return nil
}
func (s *data) LoadEntity(ctx context.Context, entity es.Entity) error {
	if err := s.db.ModelContext(ctx, entity).WherePK().Select(); err != nil {
		if isNoRow(err) {
			return es.ErrNoRows
		}
		return err
	}
	return nil
}
func (s *data) LoadSnapshot(ctx context.Context, rev string, agg es.AggregateSourced) error {
	ss := &snapshot{
		ID:       agg.GetID(),
		Type:     agg.GetTypeName(),
		Revision: rev,
	}
	if err := s.db.ModelContext(ctx, ss).WherePK().Select(); err != nil {
		if isNoRow(err) {
			return nil
		}
		return err
	}

	if _, err := types.Unmarshal(agg, ss.Aggregate, true); err != nil {
		return err
	}
	return nil
}
func (s *data) LoadEventsByType(ctx context.Context, aggregateTypeName string, eventTypeNames ...string) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("event.aggregate_type = ?", aggregateTypeName).
		WhereIn("event.type IN (?)", eventTypeNames).
		Select(); err != nil {
		return nil, err
	}

	out := make([]events.Event, len(evts))
	for i, evt := range evts {
		data, err := types.UnmarshalByName(evt.Type, evt.Data, s.legacy)
		if err != nil {
			return nil, err
		}

		m := events.Event{
			AggregateID:   evt.AggregateID,
			AggregateType: evt.AggregateType,
			Version:       evt.Version,
			Type:          evt.Type,
			Timestamp:     evt.Timestamp,
			Data:          data,
		}
		out[i] = m
	}
	return out, nil
}
func (s *data) LoadUniqueEvents(ctx context.Context, typeName string) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("event.aggregate_type = ? AND event.version = ?", typeName, 1).
		Order("version").
		Select(); err != nil {
		return nil, err
	}

	out := make([]events.Event, len(evts))
	for i, evt := range evts {
		data, err := types.UnmarshalByName(evt.Type, evt.Data, s.legacy)
		if err != nil {
			return nil, err
		}

		m := events.Event{
			AggregateID:   evt.AggregateID,
			AggregateType: evt.AggregateType,
			Version:       evt.Version,
			Type:          evt.Type,
			Timestamp:     evt.Timestamp,
			Data:          data,
		}
		out[i] = m
	}
	return out, nil
}
func (s *data) LoadAllEvents(ctx context.Context) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Order("aggregate_type", "version").
		Select(); err != nil {
		return nil, err
	}

	out := make([]events.Event, len(evts))
	for i, evt := range evts {
		data, err := types.UnmarshalByName(evt.Type, evt.Data, s.legacy)
		if err != nil {
			return nil, err
		}

		m := events.Event{
			AggregateID:   evt.AggregateID,
			AggregateType: evt.AggregateType,
			Version:       evt.Version,
			Type:          evt.Type,
			Timestamp:     evt.Timestamp,
			Data:          data,
		}
		out[i] = m
	}
	return out, nil
}
func (s *data) LoadEvent(ctx context.Context, id string, typeName string, version int) (*events.Event, error) {
	// Select all users.
	var evt event
	if err := s.db.
		ModelContext(ctx, &evt).
		Where("event.aggregate_id = ? AND event.aggregate_type = ? AND event.version = ?", id, typeName, version).
		Order("version").
		Select(); err != nil {
		return nil, err
	}

	data, err := types.UnmarshalByName(evt.Type, evt.Data, s.legacy)
	if err != nil {
		return nil, err
	}

	m := &events.Event{
		AggregateID:   evt.AggregateID,
		AggregateType: evt.AggregateType,
		Version:       evt.Version,
		Type:          evt.Type,
		Timestamp:     evt.Timestamp,
		Data:          data,
	}
	return m, nil
}
func (s *data) LoadEvents(ctx context.Context, id string, typeName string, fromVersion int) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("event.aggregate_id = ? AND event.aggregate_type = ? AND event.version > ?", id, typeName, fromVersion).
		Order("version").
		Select(); err != nil {
		return nil, err
	}

	out := make([]events.Event, len(evts))
	for i, evt := range evts {
		data, err := types.UnmarshalByName(evt.Type, evt.Data, s.legacy)
		if err != nil {
			return nil, err
		}

		m := events.Event{
			AggregateID:   evt.AggregateID,
			AggregateType: evt.AggregateType,
			Version:       evt.Version,
			Type:          evt.Type,
			Timestamp:     evt.Timestamp,
			Data:          data,
		}
		out[i] = m
	}
	return out, nil
}
