package es

import (
	"context"

	"github.com/contextgg/pkg/types"
)

type aggregateCommandHandler struct {
	entityStore EntityStore
	name        string
}

func (a *aggregateCommandHandler) HandleCommand(ctx context.Context, cmd Command) error {
	replay := IsReplayCommand(cmd)
	id := cmd.GetAggregateId()

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

// NewAggregateSourcedHandler creates the commandhandler
func NewAggregateSourcedHandler(entityStore EntityStore, entityType EntityType) CommandHandler {
	_, name := types.GetTypeName(entityType)
	handler := &aggregateCommandHandler{
		entityStore: entityStore,
		name:        name,
	}
	return handler
}
