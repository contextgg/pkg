package es

import (
	"context"
	"fmt"
)

type aggregateHandler struct {
	cfg         *AggregateConfig
	entityStore EntityStore
}

func (a *aggregateHandler) HandleCommand(ctx context.Context, cmd Command) error {
	unit := UnitFromContext(ctx)
	if unit == nil {
		return fmt.Errorf("unit not found")
	}
	replay := IsReplayCommand(cmd)
	id := cmd.GetAggregateId()
	name := a.cfg.Name

	aggregate := a.cfg.Factory(id)
	if err := unit.Load(ctx, id, name, aggregate); 

	aggregate, err := a.entityStore.Load(ctx, a.name, id, DataLoadForce(replay))
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

	return a.entityStore.Save(ctx, aggregate)
}

// NewAggregateHandler creates the commandhandler
func NewAggregateHandler(cfg *AggregateConfig) CommandHandler {
	handler := &aggregateHandler{
		cfg: cfg,
	}
	return handler
}
