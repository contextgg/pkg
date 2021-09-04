package pgdb

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
		where pg_stat_activity.datname = $1 and pid <> pg_backend_pid()`
	_, err := db.Exec(query, name)
	return err
}

func dropCreate(db DB, name string) error {
	query := `
    drop database if exists $1; create database $1
  `
	_, err := db.Exec(query, name)
	return err
}
