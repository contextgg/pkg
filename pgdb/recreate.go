package pgdb

import "fmt"

func Recreate(conn, dbName, dbUser, dbPassword string, opts ...PostgresOptions) (DB, error) {
	// connect to the postgres db!
	db, err := SetupPostgres(conn, "postgres", dbUser, dbPassword, opts...)
	if err != nil {
		return nil, err
	}
	defer func() {
		db.Close()
	}()

	// drop all connections
	if err := dropConnections(db, dbName); err != nil {
		return nil, err
	}

	// drop the database
	if err := dropCreate(db, dbName); err != nil {
		return nil, err
	}

	// connect to the db
	return SetupPostgres(conn, dbName, dbUser, dbPassword, opts...)
}

func dropConnections(db DB, name string) error {
	query := `
		select pg_terminate_backend(pg_stat_activity.pid)
		from pg_stat_activity
		where pg_stat_activity.datname = ? and pid <> pg_backend_pid()`
	_, err := db.Exec(query, name)
	return err
}

func dropCreate(db DB, name string) error {
	q1 := fmt.Sprintf(`drop database if exists %s`, name)
	q2 := fmt.Sprintf(`create database %s`, name)

	if _, err := db.Exec(q1); err != nil {
		return err
	}
	if _, err := db.Exec(q2); err != nil {
		return err
	}
	return nil
}
