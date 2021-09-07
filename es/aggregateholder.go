package es

import (
	"github.com/contextgg/pkg/events"
)

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
	Namespace string `bun:",pk"`
	Id        string `bun:",pk,type:uuid"`

	typeName string
	events   []events.Event
}

// GetID of the aggregate
func (a *BaseAggregateHolder) GetID() string {
	return a.Id
}

// GetTypeName of the aggregate
func (a *BaseAggregateHolder) GetTypeName() string {
	return a.typeName
}

// SetNamespace of the aggregate
func (a *BaseAggregateHolder) SetNamespace(namespace string) {
	a.Namespace = namespace
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

func NewBaseAggregateHolder(id string, typeName string) BaseAggregateHolder {
	return BaseAggregateHolder{
		Id:       id,
		typeName: typeName,
	}
}
