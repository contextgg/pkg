package es

import (
	"context"
	"sync"

	"github.com/contextgg/pkg/events"
	"github.com/uptrace/bun"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Unit interface {
	// Db access to the transaction
	Db() bun.IDB

	// Create a new TX
	Begin(ctx context.Context) (Tx, error)

	// Dispatch will dispatch the events to the event publishers
	Dispatch(ctx context.Context, cmds ...Command) error

	// Load will load the entity from the database.
	Load(ctx context.Context, entity Entity, id string, dataOptions ...DataLoadOption) (Entity, error)

	// Save will save the entity to the database.
	Save(ctx context.Context, entity Entity) error
}

type unit struct {
	sync.RWMutex
	cli Client
	db  bun.IDB
	tx  *bun.Tx

	events []events.Event
}

func (u *unit) Db() bun.IDB {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}

func (u *unit) Begin(ctx context.Context) (Tx, error) {
	u.Lock()
	defer u.Unlock()

	if u.tx == nil {
		tx, err := u.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		u.tx = &tx
	}
	return u, nil
}

func (u *unit) Commit(ctx context.Context) error {
	u.Lock()
	defer u.Unlock()

	if u.tx == nil {
		return nil
	}
	err := u.tx.Commit()
	if err != nil {
		return err
	}
	u.tx = nil

	// send over the
	if err := u.cli.PublishEvents(ctx, u.events...); err != nil {
		// TODO log this!!!
		return err
	}
	u.events = nil
	return nil
}

func (u *unit) Rollback(ctx context.Context) error {
	u.Lock()
	defer u.Unlock()

	if u.tx == nil {
		return nil
	}
	err := u.tx.Rollback()
	if err != nil {
		return err
	}
	u.tx = nil
	return nil
}

func (u *unit) Dispatch(ctx context.Context, cmds ...Command) error {
	ctx = SetUnit(ctx, u)
	return u.cli.HandleCommands(ctx, cmds...)
}

func (u *unit) Load(ctx context.Context, entity Entity, id string, dataOptions ...DataLoadOption) (Entity, error) {
	entityOptions, err := u.cli.GetEntityOptions(entity)
	if err != nil {
		return nil, err
	}

	db := u.Db()
	data := NewData(db)
	dataStore := NewDataStore(data, entityOptions)

	return dataStore.Load(ctx, id, dataOptions...)
}

func (u *unit) Save(ctx context.Context, entity Entity) error {
	entityOptions, err := u.cli.GetEntityOptions(entity)
	if err != nil {
		return err
	}

	db := u.Db()
	data := NewData(db)
	dataStore := NewDataStore(data, entityOptions)

	events, err := dataStore.Save(ctx, entity)
	if err != nil {
		return err
	}

	u.events = append(u.events, events...)
	return u.cli.HandleEvents(ctx, events...)
}

func newUnit(cli Client, db bun.IDB) Unit {
	return &unit{
		cli: cli,
		db:  db,
	}
}
