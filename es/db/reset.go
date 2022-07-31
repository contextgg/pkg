package db

import (
	"fmt"
)

func Reset(opts ...OptionFunc) error {
	o := NewOptions()
	for _, opt := range opts {
		opt(o)
	}
	dbName := o.DbName
	WithDbName("postgres")(o)

	db, err := NewDb(o)
	if err != nil {
		return err
	}

	query := `
		select pg_terminate_backend(pg_stat_activity.pid)
		from pg_stat_activity
		where pg_stat_activity.datname = ? and pid <> pg_backend_pid()`
	if _, err := db.Exec(query, dbName); err != nil {
		return err
	}

	q1 := fmt.Sprintf(`drop database if exists %s`, dbName)
	if _, err := db.Exec(q1); err != nil {
		return err
	}

	q2 := fmt.Sprintf(`create database %s`, dbName)
	if _, err := db.Exec(q2); err != nil {
		return err
	}

	return db.Close()
}
