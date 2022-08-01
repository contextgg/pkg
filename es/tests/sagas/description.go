package sagas

import (
	"context"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"

	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
)

type description struct {
}

func (s *description) Run(ctx context.Context, event events.Event) ([]es.Command, error) {
	switch d := event.Data.(type) {
	case *eventdata.DemoCreated:
		return s.runDemoCreated(ctx, event, d)
	}
	return nil, nil
}

func (s *description) runDemoCreated(ctx context.Context, event events.Event, data *eventdata.DemoCreated) ([]es.Command, error) {
	return es.ReturnCommands(&commands.AddDescription{
		BaseCommand: es.BaseCommand{
			AggregateId: event.AggregateId,
		},
		Description: "Done!",
	})
}

func NewDescription() es.Saga {
	return &description{}
}
