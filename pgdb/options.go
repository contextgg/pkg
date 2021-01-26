package pgdb

import "github.com/go-pg/pg/v10"

type PostgresOptions func(*pg.DB)

func AddLogging() PostgresOptions {
	return func(db *pg.DB) {
		db.AddQueryHook(dbLogger{})
	}
}
