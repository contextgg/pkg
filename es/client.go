package es

import (
	"context"
	"fmt"
	"reflect"

	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/types"
)

type Client interface {
	Initialize(ctx context.Context) error
	Unit(ctx context.Context) (Unit, error)
	GetEntityOptions(entity Entity) (*EntityOptions, error)
	HandleCommands(ctx context.Context, cmds ...Command) error
	HandleEvents(ctx context.Context, events ...events.Event) error
	PublishEvents(ctx context.Context, events ...events.Event) error
}

type client struct {
	cfg  Config
	conn Conn

	entities        []Entity
	entityOptions   map[string]*EntityOptions
	commandHandlers map[reflect.Type]CommandHandler
	eventHandlers   map[reflect.Type]EventHandler
}

func (c *client) Unit(ctx context.Context) (Unit, error) {
	if unit := UnitFromContext(ctx); unit != nil {
		return unit, nil
	}
	return newUnit(c, c.conn.Db()), nil
}

func (c *client) Initialize(ctx context.Context) error {
	sagas := c.cfg.GetSagas()
	for _, saga := range sagas {
		handler := NewSagaHandler(c, saga)
		for _, evt := range saga.events {
			t := types.GetElemType(evt)
			c.eventHandlers[t] = handler
		}
	}

	aggregates := c.cfg.GetAggregates()
	for _, agg := range aggregates {
		ent := agg.Factory("")

		c.entities = append(c.entities, ent)
		c.entityOptions[ent.GetTypeName()] = &agg.EntityOptions

		handler := NewAggregateHandler(agg)
		for _, cmd := range agg.commands {
			t := types.GetElemType(cmd)
			c.commandHandlers[t] = handler
		}
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

func (c *client) GetEntityOptions(entity Entity) (*EntityOptions, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}

	name := entity.GetTypeName()
	if opts, ok := c.entityOptions[name]; ok {
		return opts, nil
	}

	return nil, fmt.Errorf("entity options not found: %s", name)
}

func (c *client) HandleCommands(ctx context.Context, cmds ...Command) error {
	for _, cmd := range cmds {
		t := types.GetElemType(cmd)
		h, ok := c.commandHandlers[t]
		if !ok {
			return fmt.Errorf("command handler not found: %v", t)
		}
		if err := h.HandleCommand(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) HandleEvents(ctx context.Context, evts ...events.Event) error {
	for _, evt := range evts {
		t := types.GetElemType(evt.Data)
		h, ok := c.eventHandlers[t]
		if !ok {
			return nil
		}
		if err := h.HandleEvent(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) PublishEvents(ctx context.Context, evts ...events.Event) error {
	publishers := c.cfg.GetPublishers()

	for _, p := range publishers {
		for _, evt := range evts {
			if err := p.PublishEvent(ctx, evt); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewClient(cfg Config, conn Conn) (Client, error) {
	cli := &client{
		cfg:             cfg,
		conn:            conn,
		entityOptions:   map[string]*EntityOptions{},
		commandHandlers: map[reflect.Type]CommandHandler{},
		eventHandlers:   map[reflect.Type]EventHandler{},
	}

	ctx := context.Background()
	if err := cli.Initialize(ctx); err != nil {
		return nil, err
	}

	return cli, nil
}
