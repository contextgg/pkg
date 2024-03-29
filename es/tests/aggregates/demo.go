package aggregates

import (
	"context"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"

	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/eventdata"
)

type Demo struct {
	es.BaseAggregateSourced

	Name        string
	Description string
}

func (a *Demo) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch c := cmd.(type) {
	case *commands.NewDemo:
		return a.handleNewDemo(ctx, c)
	case *commands.AddDescription:
		return a.handleAddDescription(ctx, c)
	}
	return es.ErrNotHandled
}

func (a *Demo) handleNewDemo(ctx context.Context, cmd *commands.NewDemo) error {
	a.StoreEventData(ctx, &eventdata.DemoCreated{
		Name: cmd.Name,
	})
	return nil
}
func (a *Demo) handleAddDescription(ctx context.Context, cmd *commands.AddDescription) error {
	a.StoreEventData(ctx, &eventdata.DemoDescriptionAdded{
		Description: cmd.Description,
	})
	return nil
}

func (a *Demo) ApplyEvent(ctx context.Context, event events.Event) error {
	switch d := event.Data.(type) {
	case *eventdata.DemoCreated:
		return a.applyDemoCreated(ctx, event, d)
	case *eventdata.DemoDescriptionAdded:
		return a.applyDemoDescriptionAdded(ctx, event, d)
	}
	return nil
}
func (a *Demo) applyDemoCreated(ctx context.Context, event events.Event, data *eventdata.DemoCreated) error {
	a.Name = data.Name
	return nil
}
func (a *Demo) applyDemoDescriptionAdded(ctx context.Context, event events.Event, data *eventdata.DemoDescriptionAdded) error {
	a.Description = data.Description
	return nil
}

func NewDemo(id string) es.Entity {
	return &Demo{
		BaseAggregateSourced: es.NewBaseAggregateSourced(id, "Demo"),
	}
}
