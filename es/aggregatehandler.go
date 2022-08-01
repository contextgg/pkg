package es

import (
	"context"
	"fmt"
)

type aggregateHandler struct {
	fn EntityFunc
}

func (a *aggregateHandler) HandleCommand(ctx context.Context, cmd Command) error {
	unit := UnitFromContext(ctx)
	if unit == nil {
		return fmt.Errorf("unit not found")
	}

	id := cmd.GetAggregateId()
	replay := IsReplayCommand(cmd)
	entity := a.fn(id)

	aggregate, err := unit.Load(ctx, entity, id, DataLoadForce(replay))
	if err != nil {
		return err
	}

	cv, okCV := cmd.(CommandVersion)
	versioned, okVersioned := aggregate.(Versioned)
	if okCV && okVersioned && cv.GetVersion() != versioned.GetVersion() {
		return ErrVersionMismatch
	}

	// handle the command
	if !replay {
		if handler, ok := aggregate.(CommandHandler); ok {
			if err := handler.HandleCommand(ctx, cmd); err != nil {
				return err
			}
		}
	}

	return unit.Save(ctx, aggregate)
}

// NewAggregateHandler creates the commandhandler
func NewAggregateHandler(fn EntityFunc) CommandHandler {
	handler := &aggregateHandler{
		fn: fn,
	}
	return handler
}
