package es

import (
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
