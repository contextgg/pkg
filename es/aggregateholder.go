package es

import "github.com/contextgg/pkg/events"

// AggregateHolder for event stored aggregates
type AggregateHolder interface {
	Entity

	// EventsToPublish returns events that need to be published
	EventsToPublish() []events.Event

	// ClearEvents clears all uncommitted events after saving.
	ClearEvents()
}

// BaseAggregateHolder to make our commands smaller
type BaseAggregateHolder struct {
	ID       string
	TypeName string

	events []events.Event
}

// Initialize the aggregate with id and type
func (a *BaseAggregateHolder) Initialize(id string, typeName string) {
	a.ID = id
	a.TypeName = typeName
}

// GetID of the aggregate
func (a *BaseAggregateHolder) GetID() string {
	return a.ID
}

// GetTypeName of the aggregate
func (a *BaseAggregateHolder) GetTypeName() string {
	return a.TypeName
}

// PublishEvent registers an event to be published after the aggregate
// has been successfully saved.
func (a *BaseAggregateHolder) PublishEvent(e events.Event) {
	a.events = append(a.events, e)
}

// Events returns all uncommitted events that are not yet saved.
func (a *BaseAggregateHolder) EventsToPublish() []events.Event {
	return a.events
}

// ClearEvents clears all uncommitted events after saving.
func (a *BaseAggregateHolder) ClearEvents() {
	a.events = nil
}
