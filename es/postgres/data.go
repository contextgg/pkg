package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/types"
)

type event struct {
	Namespace     string `bun:",pk"`
	AggregateID   string `bun:",pk,type:uuid"`
	AggregateType string `bun:",pk"`
	Version       int    `bun:",pk"`
	Type          string `bun:",notnull"`
	Timestamp     time.Time
	Data          json.RawMessage
	Metadata      map[string]interface{}
}
type snapshot struct {
	Namespace string          `bun:",pk"`
	ID        string          `bun:",pk,type:uuid"`
	Type      string          `bun:",pk"`
	Revision  string          `bun:",pk"`
	Aggregate json.RawMessage `bun:",notnull"`
}

func run(db *bun.DB, opts *es.DataOpts) error {
	var models []interface{}

	if opts.HasEvents {
		models = append(models, &event{})
	}
	if opts.HasSnapshots {
		models = append(models, &snapshot{})
	}
	for _, model := range opts.ExtraModels {
		models = append(models, model)
	}

	ctx := context.Background()
	for _, model := range models {
		if opts.TruncateTables {
			_, err := db.NewTruncateTable().Model(model).Exec(ctx)
			if err != nil {
				return err
			}
		}

		if opts.RecreateTables {
			_, err := db.NewDropTable().Model(model).IfExists().Exec(ctx)
			if err != nil {
				return err
			}
		}

		_, err := db.NewCreateTable().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}

	migrations := migrate.NewMigrations()
	// TODO register migrations!

	migrator := migrate.NewMigrator(db, migrations)
	if _, err := migrator.Migrate(ctx); err != nil {
		if err.Error() != "migrate: there are no any migrations" {
			return err
		}
	}
	return nil
}

func NewPostgresData(db *bun.DB, opts ...es.DataOption) es.Data {
	all := &es.DataOpts{}
	for _, o := range opts {
		o(all)
	}

	// run the migrations?
	if err := run(db, all); err != nil {
		panic(err)
	}

	return &data{
		db: db,
	}
}

type data struct {
	db     bun.IDB
	legacy bool
}

func (s *data) SaveEntity(ctx context.Context, namespace string, entity es.Entity) error {
	entity.SetNamespace(namespace)

	_, err := s.db.NewInsert().
		Model(entity).
		On("CONFLICT (namespace, id) DO UPDATE").
		Exec(ctx)
	return err
}
func (s *data) DeleteEntry(ctx context.Context, namespace string, entity es.Entity) error {
	entity.SetNamespace(namespace)

	_, err := s.db.NewDelete().
		Model(entity).
		WherePK().
		Exec(ctx)
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

	_, err = s.db.NewInsert().
		Model(ss).
		On("CONFLICT (namespace,id,type,revision) DO UPDATE").
		Exec(ctx)
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
	if _, err := s.db.NewInsert().
		Model(&all).
		On("CONFLICT (namespace, aggregate_id, aggregate_type, version) DO UPDATE").
		Exec(ctx); err != nil {
		return err
	}
	return nil
}
func (s *data) LoadEntity(ctx context.Context, namespace string, entity es.Entity) error {
	entity.SetNamespace(namespace)

	if err := s.db.NewSelect().
		Model(entity).
		WherePK().
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
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

	if err := s.db.NewSelect().
		Model(ss).
		WherePK().
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
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

	if err := s.db.NewSelect().
		Model(&evts).
		Where("namespace = ?", namespace).
		Where("aggregate_type = ?", aggregateTypeName).
		Where("type IN (?)", bun.In(eventTypeNames)).
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
			return nil, es.ErrNoRows
		}
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
	if err := s.db.NewSelect().
		Model(&evts).
		Where("namespace = ?", namespace).
		Where("aggregate_type = ?", typeName).
		Where("version = ?", 1).
		Order("version").
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
			return nil, es.ErrNoRows
		}
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
	if err := s.db.NewSelect().
		Model(&evts).
		Where("namespace = ?", namespace).
		Order("aggregate_type", "version").
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
			return nil, es.ErrNoRows
		}
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
	if err := s.db.NewSelect().
		Model(&evt).
		Where("namespace = ?", namespace).
		Where("aggregate_id = ?", id).
		Where("aggregate_type = ?", typeName).
		Where("version = ?", version).
		Order("version").
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
			return nil, es.ErrNoRows
		}
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
	if err := s.db.NewSelect().
		Model(&evts).
		Where("namespace = ?", namespace).
		Where("aggregate_id = ?", id).
		Where("aggregate_type = ?", typeName).
		Where("version > ?", fromVersion).
		Order("version").
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
			return nil, es.ErrNoRows
		}
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
