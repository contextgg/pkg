package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDb(opts *Options) (*bun.DB, error) {
	dsn := opts.DSN()
	conn := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithDatabase(opts.DbName),
		pgdriver.WithUser(opts.User),
		pgdriver.WithPassword(opts.Password),
	)
	sqldb := sql.OpenDB(conn)
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(opts.Debug)))

	return db, nil
}
