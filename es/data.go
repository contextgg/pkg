package es

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/types"
	"github.com/uptrace/bun"
)

type event struct {
	Namespace     string `bun:",pk"`
	AggregateId   string `bun:",pk,type:uuid"`
	AggregateType string `bun:",pk"`
	Version       int    `bun:",pk"`
	Type          string `bun:",notnull"`
	Timestamp     time.Time
	Data          json.RawMessage        `bun:"type:jsonb"`
	Metadata      map[string]interface{} `bun:"type:jsonb"`
}
type snapshot struct {
	Namespace string          `bun:",pk"`
	Id        string          `bun:",pk,type:uuid"`
	Type      string          `bun:",pk"`
	Revision  string          `bun:",pk"`
	Aggregate json.RawMessage `bun:",notnull,type:jsonb"`
}

var ErrNoRows = errors.New("No rows found")

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

func NewData(db bun.IDB) Data {
	return &data{
		db: db,
	}
}

type data struct {
	db     bun.IDB
	legacy bool
}

func (s *data) event(evt event) (events.Event, error) {
	entry, ok := types.GetByName(evt.Type)
	if !ok {
		return events.Event{}, fmt.Errorf("Type %s is not in registry", evt.Type)
	}
	data, err := types.EntryUnmarshal(entry, evt.Data, types.UseLegacyJsonSerializer(s.legacy))
	if err != nil {
		return events.Event{}, err
	}

	return events.Event{
		AggregateId:   evt.AggregateId,
		AggregateType: evt.AggregateType,
		Version:       evt.Version,
		Type:          evt.Type,
		Timestamp:     evt.Timestamp,
		Data:          data,
	}, nil
}
func (s *data) events(evts []event) ([]events.Event, error) {
	out := make([]events.Event, len(evts))
	for i, evt := range evts {
		m, err := s.event(evt)
		if err != nil {
			return nil, err
		}
		out[i] = m
	}
	return out, nil
}

func (s *data) SaveEntity(ctx context.Context, namespace string, entity Entity) error {
	entity.SetNamespace(namespace)

	_, err := s.db.NewInsert().
		Model(entity).
		On("CONFLICT (namespace, id) DO UPDATE").
		Exec(ctx)
	return err
}
func (s *data) DeleteEntry(ctx context.Context, namespace string, entity Entity) error {
	entity.SetNamespace(namespace)

	_, err := s.db.NewDelete().
		Model(entity).
		WherePK().
		Exec(ctx)
	return err
}
func (s *data) SaveSnapshot(ctx context.Context, namespace string, rev string, agg AggregateSourced) error {
	// to json!
	data, err := json.Marshal(agg)
	if err != nil {
		return err
	}

	ss := &snapshot{
		Namespace: namespace,
		Id:        agg.GetId(),
		Type:      agg.GetTypeName(),
		Revision:  rev,
		Aggregate: json.RawMessage(data),
	}

	_, err = s.db.NewInsert().
		Model(ss).
		On("CONFLICT (namespace,id,type,revision) DO UPDATE").
		Exec(ctx)
	return err
}
func (s *data) SaveEvents(ctx context.Context, namespace string, evts ...events.Event) error {
	all := make([]event, len(evts))
	for i, evt := range evts {
		data, err := types.JsonMarshal(evt.Data, types.UseLegacyJsonSerializer(s.legacy))
		if err != nil {
			return err
		}

		m := event{
			Namespace:     namespace,
			AggregateId:   evt.AggregateId,
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
func (s *data) LoadEntity(ctx context.Context, namespace string, entity Entity) error {
	entity.SetNamespace(namespace)

	if err := s.db.NewSelect().
		Model(entity).
		WherePK().
		Scan(ctx); err != nil {
		if sql.ErrNoRows == err {
			return ErrNoRows
		}
		return err
	}

	return nil
}
func (s *data) LoadSnapshot(ctx context.Context, namespace string, rev string, agg AggregateSourced) error {
	ss := &snapshot{
		Namespace: namespace,
		Id:        agg.GetId(),
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

	if err := types.JsonUnmarshal(ss.Aggregate, agg, types.UseLegacyJsonSerializer(s.legacy)); err != nil {
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
			return nil, ErrNoRows
		}
		return nil, err
	}
	return s.events(evts)
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
			return nil, ErrNoRows
		}
		return nil, err
	}
	return s.events(evts)
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
			return nil, ErrNoRows
		}
		return nil, err
	}
	return s.events(evts)
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
			return nil, ErrNoRows
		}
		return nil, err
	}
	return s.events(evts)
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
			return nil, ErrNoRows
		}
		return nil, err
	}

	entry, ok := types.GetByName(evt.Type)
	if !ok {
		return nil, fmt.Errorf("Type %s is not in registry", evt.Type)
	}
	data, err := types.EntryUnmarshal(entry, evt.Data, types.UseLegacyJsonSerializer(s.legacy))
	if err != nil {
		return nil, err
	}

	m := &events.Event{
		AggregateId:   evt.AggregateId,
		AggregateType: evt.AggregateType,
		Version:       evt.Version,
		Type:          evt.Type,
		Timestamp:     evt.Timestamp,
		Data:          data,
	}
	return m, nil
}
