package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/ns"
)

type ReplayService interface {
	All(ctx context.Context, aggregateType string) (int, error)
	One(ctx context.Context, aggregateType string, aggregateId string) error
}

type replayService struct {
	handlers map[string]CommandHandler
}

func (r *replayService) All(ctx context.Context, aggregateType string) (int, error) {
	handler, ok := r.handlers[aggregateType]
	if !ok {
		return 0, fmt.Errorf("no handler for %s", aggregateType)
	}

	unit := UnitFromContext(ctx)
	if unit == nil {
		return 0, fmt.Errorf("no unit in context")
	}

	namespace := ns.FromContext(ctx)
	evts, err := unit.Data().LoadUniqueEvents(ctx, namespace, aggregateType)
	if err != nil {
		return 0, err
	}

	for _, evt := range evts {
		cmd := &ReplayCommand{
			AggregateId: evt.AggregateId,
		}
		if err := handler.HandleCommand(ctx, cmd); err != nil {
			return 0, err
		}
	}

	return len(evts), nil
}

func (r *replayService) One(ctx context.Context, aggregateType string, id string) error {
	handler, ok := r.handlers[aggregateType]
	if !ok {
		return fmt.Errorf("no handler for %s", aggregateType)
	}

	cmd := &ReplayCommand{
		AggregateId: id,
	}
	if err := handler.HandleCommand(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func NewReplayService(handlers map[string]CommandHandler) ReplayService {
	return &replayService{
		handlers: handlers,
	}
}
