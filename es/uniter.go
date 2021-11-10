package es

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type key int

const unitKey key = 0

var ErrUnitNotFound = fmt.Errorf("Unit not found")

type UniterExec = func(ctx context.Context) error

type Uniter interface {
	Run(ctx context.Context, exec UniterExec) error
}

type uniter struct {
	db              *bun.DB
	eventPublishers EventPublishers
}

func (u *uniter) create(ctx context.Context, parent Unit) (Unit, error) {
	raw, _ := parent.(*unit)

	var tx *bun.Tx
	var err error
	if raw != nil {
		tx = raw.tx
		err = raw.err
	}

	if tx == nil {
		// create the root!
		// create the tx!
		begin, err := u.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		tx = &begin
	}

	return &unit{
		parent: parent,
		db:     u.db,
		tx:     tx,
		err:    err,
	}, nil
}

func (u *uniter) rollback(ctx context.Context, unit Unit) error {
	if unit.Error() != nil {
		unit.SetError(fmt.Errorf("Unit rolled back"))
	}

	// do transaction stuff!
	tx := unit.DB()
	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}
func (u *uniter) commit(ctx context.Context, unit Unit) error {
	evts := unit.Events()
	unit.ClearEvents()

	parent := unit.Parent()
	if parent != nil {
		parent.StoreEvents(evts...)
		parent.SetError(unit.Error())
		return unit.Error()
	}

	if unit.Error() != nil {
		return unit.Error()
	}

	// do transaction stuff!
	tx := unit.DB()
	if err := tx.Commit(); err != nil {
		unit.SetError(err)
		return err
	}

	if err := u.eventPublishers.PublishEvents(ctx, evts); err != nil {
		return err
	}
	return nil
}

func (u *uniter) Run(ctx context.Context, exec UniterExec) error {
	// setup the new ctx!
	parent, _ := GetUnit(ctx)
	unit, err := u.create(ctx, parent)
	if err != nil {
		return err
	}

	ctxn := context.WithValue(ctx, unitKey, unit)
	if err := exec(ctxn); err != nil {
		// rollback!
		u.rollback(ctx, unit)
		return err
	}

	// commit!
	u.commit(ctx, unit)
	return nil
}

func GetUnit(ctx context.Context) (Unit, error) {
	current, ok := ctx.Value(unitKey).(Unit)
	if !ok {
		return nil, ErrUnitNotFound
	}
	return current, nil
}

func NewUniter(db *bun.DB, eventPublishers ...EventPublisher) Uniter {
	u := &uniter{
		db:              db,
		eventPublishers: eventPublishers,
	}
	return u
}
