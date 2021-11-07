package sagas

import (
	"context"
	"fmt"
	"log"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"

	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/eventdata"
)

type counter struct {
	store es.Store
}

func (s *counter) Run(ctx context.Context, event events.Event) ([]es.Command, error) {
	switch d := event.Data.(type) {
	case *eventdata.DemoCreated:
		return s.runDemoCreated(ctx, event, d)
	}
	return nil, nil
}

func (s *counter) runDemoCreated(ctx context.Context, event events.Event, data *eventdata.DemoCreated) ([]es.Command, error) {
	e, err := s.store.Load(ctx, event.AggregateId, false)
	if err != nil {
		return nil, err
	}
	agg, ok := e.(*aggregates.Demo)
	if !ok {
		return nil, nil
	}

	log.Printf("Agg %v", agg)

	// default create the invite!
	return nil, fmt.Errorf("Whoops")
}

func NewCounter(store es.Store) es.Saga {
	return &counter{
		store: store,
	}
}
