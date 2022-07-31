package es

import (
	"context"
	"fmt"
)

type aggregateHandler struct {
	cfg *AggregateConfig
}

func (a *aggregateHandler) HandleCommand(ctx context.Context, cmd Command) error {
	unit := UnitFromContext(ctx)
	if unit == nil {
		return fmt.Errorf("unit not found")
	}

	id := cmd.GetAggregateId()
	replay := IsReplayCommand(cmd)

	aggregate, err := unit.Load(ctx, &a.cfg.EntityOptions, id, DataLoadForce(replay))
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
func NewAggregateHandler(cfg *AggregateConfig) CommandHandler {
	handler := &aggregateHandler{
		cfg: cfg,
	}
	return handler
}
