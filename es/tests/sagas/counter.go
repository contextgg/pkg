package sagas

import (
	"context"
	"log"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/types"

	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
)

type counter struct {
	entityStore es.EntityStore
	name        string
}

func (s *counter) Run(ctx context.Context, event events.Event) ([]es.Command, error) {
	switch d := event.Data.(type) {
	case *eventdata.DemoCreated:
		return s.runDemoCreated(ctx, event, d)
	}
	return nil, nil
}

func (s *counter) runDemoCreated(ctx context.Context, event events.Event, data *eventdata.DemoCreated) ([]es.Command, error) {
	agg, err := s.entityStore.Load(ctx, s.name, event.AggregateId)
	if err != nil {
		return nil, err
	}

	log.Printf("Agg %v", agg)

	return es.ReturnCommands(&commands.AddDescription{
		BaseCommand: es.BaseCommand{
			AggregateId: event.AggregateId,
		},
		Description: "Done!",
	})
}

func NewCounter(entityStore es.EntityStore) es.Saga {
	name := types.GetTypeName(&aggregates.Demo{})

	return &counter{
		name:        name,
		entityStore: entityStore,
	}
}
