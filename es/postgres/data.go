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
	Namespace     string `pg:",pk"`
	AggregateID   string `pg:",pk,type:uuid"`
	AggregateType string `pg:",pk"`
	Version       int    `pg:",pk"`
	Type          string `pg:",notnull"`
	Timestamp     time.Time
	Data          json.RawMessage
	Metadata      map[string]interface{}
}
type snapshot struct {
	Namespace string          `pg:",pk"`
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
		if err := run(db, o); err != nil {
			panic(err)
		}
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

func (s *data) SaveEntity(ctx context.Context, namespace string, entity es.Entity) error {
	entity.SetNamespace(namespace)
	_, err := s.db.ModelContext(ctx, entity).
		OnConflict("(namespace, id) DO UPDATE").
		Insert()
	return err
}
func (s *data) DeleteEntry(ctx context.Context, namespace string, entity es.Entity) error {
	entity.SetNamespace(namespace)
	_, err := s.db.ModelContext(ctx, entity).
		WherePK().
		Delete()
	return err
}
func (s *data) SaveSnapshot(ctx context.Context, namespace string, rev string, agg es.AggregateSourced) error {
	// to json!
	data, err := json.Marshal(agg)
	if err != nil {
		return err
	}

	ss := &snapshot{
		Namespace: namespace,
		ID:        agg.GetID(),
		Type:      agg.GetTypeName(),
		Revision:  rev,
		Aggregate: json.RawMessage(data),
	}
	if _, err := s.db.ModelContext(ctx, ss).
		OnConflict("(namespace,id,type,revision) DO UPDATE").
		Insert(); err != nil {
		return err
	}
	return nil
}
func (s *data) SaveEvents(ctx context.Context, namespace string, evts ...events.Event) error {
	all := make([]event, len(evts))
	for i, evt := range evts {
		data, err := types.Marshal(evt.Data, s.legacy)
		if err != nil {
			return err
		}

		m := event{
			Namespace:     namespace,
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
		OnConflict("(namespace, aggregate_id, aggregate_type, version) DO UPDATE").
		Insert(); err != nil {
		return err
	}
	return nil
}
func (s *data) LoadEntity(ctx context.Context, namespace string, entity es.Entity) error {
	entity.SetNamespace(namespace)
	if err := s.db.ModelContext(ctx, entity).WherePK().Select(); err != nil {
		if isNoRow(err) {
			return es.ErrNoRows
		}
		return err
	}
	return nil
}
func (s *data) LoadSnapshot(ctx context.Context, namespace string, rev string, agg es.AggregateSourced) error {
	ss := &snapshot{
		Namespace: namespace,
		ID:        agg.GetID(),
		Type:      agg.GetTypeName(),
		Revision:  rev,
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
func (s *data) LoadEventsByType(ctx context.Context, namespace string, aggregateTypeName string, eventTypeNames ...string) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("namespace = ?", namespace).
		Where("aggregate_type = ?", aggregateTypeName).
		WhereIn("type IN (?)", eventTypeNames).
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
func (s *data) LoadUniqueEvents(ctx context.Context, namespace string, typeName string) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("namespace = ?", namespace).
		Where("aggregate_type = ?", typeName).
		Where("version = ?", 1).
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
func (s *data) LoadAllEvents(ctx context.Context, namespace string) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("namespace = ?", namespace).
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
func (s *data) LoadEvent(ctx context.Context, namespace string, id string, typeName string, version int) (*events.Event, error) {
	// Select all users.
	var evt event
	if err := s.db.
		ModelContext(ctx, &evt).
		Where("namespace = ?", namespace).
		Where("aggregate_id = ?", id).
		Where("aggregate_type = ?", typeName).
		Where("version = ?", version).
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
func (s *data) LoadEvents(ctx context.Context, namespace string, id string, typeName string, fromVersion int) ([]events.Event, error) {
	// Select all users.
	var evts []event
	if err := s.db.
		ModelContext(ctx, &evts).
		Where("namespace = ?", namespace).
		Where("aggregate_id = ?", id).
		Where("aggregate_type = ?", typeName).
		Where("version > ?", fromVersion).
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
