package es

import "context"

type Client interface {
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

	return newUnit(c.cfg, data)
}

func NewClient(cfg Config, conn Conn) (Client, error) {
	cli := &client{
		cfg:  cfg,
		conn: conn,
	}
	return cli, nil
}
