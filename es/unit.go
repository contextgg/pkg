package es

import (
	"context"

	"github.com/contextgg/pkg/ns"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Unit interface {
	// DB access to the transaction
	DB() bun.IDB
	// Dispatch will dispatch the events to the event publishers
	Dispatch(ctx context.Context, cmds ...Command) error
	// Load will load the aggregate from the database.
	Load(ctx context.Context, id string, aggregateName string, out interface{}) error
}

type unit struct {
	cli Client
	db  bun.IDB
	tx  *bun.Tx
}

func (u *unit) DB() bun.IDB {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}

func (u *unit) Load(ctx context.Context, id uuid.UUID, aggregateName string, out interface{}) error {
	namespace := ns.FromContext(ctx)
	data := NewData(u.DB())

	return data.Load(ctx, u.serviceName, aggregateName, namespace, id, out)
}

func newUnit(cli Client, db bun.IDB) Unit {
	return &unit{
		cli: cli,
		db:  db,
	}
}
