package pgdb

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	_ "github.com/contextgg/pkg/pgdb/timestamp"
)

var (
	// ErrNotInTransaction is returned when using Commit
	// outside of a transaction.
	ErrNotInTransaction = errors.New("not in transaction")
)

type DB interface {
	orm.DB

	BeginContext(ctx context.Context) (DB, error)
	CommitContext(ctx context.Context) error
	RollbackContext(ctx context.Context) error
	Close() error
}

func IsPQNoRow(err error) bool {
	return err != nil && err.Error() == "pg: no rows in result set"
}

type database struct {
	orm.DB

	db     *pg.DB
	tx     *pg.Tx
	nested bool
}

func (n database) BeginContext(ctx context.Context) (DB, error) {
	var err error

	switch {
	case n.tx == nil:
		// new actual transaction
		n.tx, err = n.db.BeginContext(ctx)
		n.DB = n.tx
	default:
		// already in a transaction: reusing current transaction
		n.nested = true
	}

	if err != nil {
		return nil, err
	}

	return &n, nil
}
func (n *database) CommitContext(ctx context.Context) error {
	if n.tx == nil {
		return nil
	}

	var err error

	if !n.nested {
		err = n.tx.Commit()
	}

	if err != nil {
		return err
	}

	n.tx = nil
	return nil
}
func (n *database) RollbackContext(ctx context.Context) error {
	if n.tx == nil {
		return nil
	}

	var err error

	if !n.nested {
		err = n.tx.Rollback()
	}

	if err != nil {
		return err
	}

	n.tx = nil
	return nil
}
func (n *database) Close() error {
	return n.db.Close()
}

// SetupPostgres for connecting to a database
func SetupPostgres(conn, dbName, dbUser, dbPassword string, opts ...PostgresOptions) (DB, error) {
	o, err := pg.ParseURL(conn)
	if err != nil {
		return nil, err
	}

	if len(dbName) > 0 {
		o.Database = dbName
	}
	if len(dbUser) > 0 {
		o.User = dbUser
	}
	if len(dbPassword) > 0 {
		o.Password = dbPassword
	}

	db := pg.Connect(o)
	for _, item := range opts {
		item(db)
	}
	return &database{
		DB: db,
		db: db,
	}, nil
}
