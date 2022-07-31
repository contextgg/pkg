package es

import (
	"context"

	"github.com/contextgg/pkg/es/db"
	"github.com/uptrace/bun"
)

type Conn interface {
	Db() bun.IDB
}

type conn struct {
	db *bun.DB
}

func (c *conn) Db() bun.IDB {
	return c.db
}

func NewConn(opts ...db.OptionFunc) (Conn, error) {
	o := db.NewOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.Recreate {
		if err := db.Reset(opts...); err != nil {
			return nil, err
		}
	}

	db, err := db.NewDb(o)
	if err != nil {
		return nil, err
	}

	c := &conn{
		db: db,
	}
	return c, nil
}

func MigrateDatabase(db bun.IDB, options ...DataOption) error {
	opts := dataOptions(options)

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

		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
