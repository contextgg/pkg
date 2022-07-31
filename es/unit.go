package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/ns"
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

func (u *unit) Dispatch(ctx context.Context, cmds ...Command) error {
	ctx = SetUnit(ctx, u)

	for _, cmd := range cmds {
		h, err := u.cli.GetCommandHandler(cmd)
		if err != nil {
			return err
		}

		if err := h.HandleCommand(ctx, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (u *unit) Load(ctx context.Context, id string, aggregateName string, out interface{}) error {
	namespace := ns.FromContext(ctx)
	data := NewData(u.DB())

	// return data.Load(ctx, u.serviceName, aggregateName, namespace, id, out)
	return fmt.Errorf("not implemented")
}

func newUnit(cli Client, db bun.IDB) Unit {
	return &unit{
		cli: cli,
		db:  db,
	}
}
