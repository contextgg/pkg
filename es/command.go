package es

import (
	"errors"
)

func ReturnCommands(cmds ...Command) ([]Command, error) {
	return cmds, nil
}

var (
	ErrNotHandled = errors.New("Command not handled")
)

// Command will find its way to an aggregate
type Command interface {
	GetAggregateId() string
}

// BaseCommand to make it easier to get the ID
type BaseCommand struct {
	AggregateId string `json:"aggregate_id"`
}

// GetAggregateId return the aggregate id
func (c BaseCommand) GetAggregateId() string {
	return c.AggregateId
}

// ReplayCommand a command that load and reply events ontop of an aggregate.
type ReplayCommand struct {
	AggregateId string
}

// GetAggregateId return the aggregate id
func (c ReplayCommand) GetAggregateId() string {
	return c.AggregateId
}

func IsReplayCommand(cmd Command) bool {
	// handle the command
	switch cmd.(type) {
	case *ReplayCommand:
		return true
	default:
		return false
	}
}
