package pg

import (
	"fmt"

	"github.com/uptrace/bun"
)

type DbFactory func() (*bun.DB, error)

func Recreate(fact DbFactory, dbName string) error {
	// connect to the postgres db!
	db, err := fact()
	if err != nil {
		return err
	}
	defer func() {
		db.Close()
	}()

	// drop all connections
	if err := dropConnections(db, dbName); err != nil {
		return err
	}

	// drop the database
	if err := dropDb(db, dbName); err != nil {
		return err
	}

	if err := createDb(db, dbName); err != nil {
		return err
	}
	return nil
}

func dropConnections(db *bun.DB, name string) error {
	query := `
		select pg_terminate_backend(pg_stat_activity.pid)
		from pg_stat_activity
		where pg_stat_activity.datname = ? and pid <> pg_backend_pid()`
	_, err := db.Exec(query, name)
	return err
}

func dropDb(db *bun.DB, name string) error {
	q1 := fmt.Sprintf(`drop database if exists %s`, name)
	if _, err := db.Exec(q1); err != nil {
		return err
	}
	return nil
}
func createDb(db *bun.DB, name string) error {
	q2 := fmt.Sprintf(`create database %s`, name)
	if _, err := db.Exec(q2); err != nil {
		return err
	}
	return nil
}
