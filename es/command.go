package es

import (
	"errors"

	"github.com/google/uuid"
)

func ReturnCommands(cmds ...Command) ([]Command, error) {
	return cmds, nil
}

var (
	ErrNotHandled = errors.New("Command not handled")
)

// Command will find its way to an aggregate
type Command interface {
	GetAggregateId() uuid.UUID
}

// BaseCommand to make it easier to get the ID
type BaseCommand struct {
	AggregateId uuid.UUID `json:"aggregate_id"`
}

// GetAggregateID return the aggregate id
func (c BaseCommand) GetAggregateId() uuid.UUID {
	return c.AggregateId
}

// ReplayCommand a command that load and reply events ontop of an aggregate.
type ReplayCommand struct {
	AggregateId uuid.UUID
}

// GetAggregateID return the aggregate id
func (c ReplayCommand) GetAggregateId() uuid.UUID {
	return c.AggregateId
}
