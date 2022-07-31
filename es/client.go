package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/types"
)

type Client interface {
	Initialize(ctx context.Context) error
	Unit(ctx context.Context) (Unit, error)
	GetCommandHandler(cmd Command) (CommandHandler, error)
}

type client struct {
	cfg  Config
	conn Conn

	entities []Entity
	handlers map[string]CommandHandler
}

func (c *client) Unit(ctx context.Context) (Unit, error) {
	if unit := UnitFromContext(ctx); unit != nil {
		return unit, nil
	}

	return newUnit(c, c.conn.Db()), nil
}

func (c *client) Initialize(ctx context.Context) error {
	aggregates := c.cfg.GetAggregates()
	for _, agg := range aggregates {
		ent := agg.Factory("")

		c.entities = append(c.entities, ent)
	}

	db := c.conn.Db()

	if err := MigrateDatabase(
		db,
		InitializeEvents(),
		InitializeSnapshots(),
		InitializeEntities(c.entities...),
	); err != nil {
		return err
	}
	return nil
}

func (c *client) GetCommandHandler(cmd Command) (CommandHandler, error) {
	name := types.GetTypeName(cmd)
	h, ok := c.handlers[name]
	if !ok {
		return nil, fmt.Errorf("command handler not found: %s", name)
	}
	return h, nil
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
