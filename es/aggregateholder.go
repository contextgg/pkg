package es

import (
	"github.com/contextgg/pkg/events"
	"github.com/google/uuid"
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
	Namespace string    `pg:",pk"`
	Id        uuid.UUID `pg:",pk,type:uuid"`

	typeName string
	events   []events.Event
}

// GetNamespace of the aggregate
func (a *BaseAggregateHolder) GetNamespace() string {
	return a.Namespace
}

// GetID of the aggregate
func (a *BaseAggregateHolder) GetID() uuid.UUID {
	return a.Id
}

// GetTypeName of the aggregate
func (a *BaseAggregateHolder) GetTypeName() string {
	return a.typeName
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

func NewBaseAggregateHolder(namespace string, id uuid.UUID, typeName string) BaseAggregateHolder {
	return BaseAggregateHolder{
		Namespace: namespace,
		Id:        id,
		typeName:  typeName,
	}
}
