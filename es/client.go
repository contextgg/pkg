package es

import (
	"context"
)

type Client interface {
	Initialize(ctx context.Context) error
	Unit(ctx context.Context) (Unit, error)
}

type client struct {
	cfg  Config
	conn Conn
}

func (c *client) Unit(ctx context.Context) (Unit, error) {
	if unit := UnitFromContext(ctx); unit != nil {
		return unit, nil
	}

	return newUnit(c.cfg, c.conn.Db())
}

func (c *client) Initialize(ctx context.Context) error {
	entities := c.cfg.GetEntities()
	db := c.conn.Db()

	if err := MigrateDatabase(
		db,
		InitializeEvents(),
		InitializeSnapshots(),
		InitializeEntities(entities...),
	); err != nil {
		return err
	}

	return nil
}

func NewClient(cfg Config, conn Conn) (Client, error) {
	cli := &client{
		cfg:  cfg,
		conn: conn,
	}

	ctx := context.Background()
	if err := cli.Initialize(ctx); err != nil {
		return nil, err
	}

	return cli, nil
}
