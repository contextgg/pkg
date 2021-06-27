package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// EventApplier for applying events
type EventApplier interface {
	// ApplyEvent applies an event on the aggregate by setting its values.
	// If there are no errors the version should be incremented by calling
	// IncrementVersion.
	ApplyEvent(context.Context, events.Event) error
}

// AggregateSourced for event stored aggregates
type AggregateSourced interface {
	Entity
	CommandHandler
	EventApplier

	// StoreEventData will create an event and store it
	StoreEventData(context.Context, interface{})

	// GetVersion returns the version of the aggregate.
	GetVersion() int

	// Increment version increments the version of the aggregate. It should be
	// called after an event has been successfully applied.
	IncrementVersion()

	// Events returns all uncommitted events that are not yet saved.
	Events() []events.Event

	// ClearEvents clears all uncommitted events after saving.
	ClearEvents()
}

// BaseAggregateSourced to make our commands smaller
type BaseAggregateSourced struct {
	Namespace string `pg:",pk" json:"-"`
	Id        string `pg:",pk,type:uuid"`
	Version   int    `pg:"-"`

	typeName string
	events   []events.Event
}

// GetID of the aggregate
func (a *BaseAggregateSourced) GetID() string {
	return a.Id
}

// GetTypeName of the aggregate
func (a *BaseAggregateSourced) GetTypeName() string {
	return a.typeName
}

// SetNamespace of the aggregate
func (a *BaseAggregateSourced) SetNamespace(namespace string) {
	a.Namespace = namespace
}

// StoreEventData will add the event to a list which will be persisted later
func (a *BaseAggregateSourced) StoreEventData(ctx context.Context, data interface{}) {
	v := a.Version + len(a.events) + 1

	e := NewEvent(ctx, a, v, data)
	a.events = append(a.events, e)
}

// GetVersion returns the version of the aggregate.
func (a *BaseAggregateSourced) GetVersion() int {
	return a.Version
}

// IncrementVersion ads 1 to the current version
func (a *BaseAggregateSourced) IncrementVersion() {
	a.Version++
}

// Events returns all uncommitted events that are not yet saved.
func (a *BaseAggregateSourced) Events() []events.Event {
	return a.events
}

// ClearEvents clears all uncommitted events after saving.
func (a *BaseAggregateSourced) ClearEvents() {
	a.events = []events.Event{}
}

// NewBaseAggregateSourced create a new base aggregate
func NewBaseAggregateSourced(id string, typeName string) BaseAggregateSourced {
	return BaseAggregateSourced{
		Id:       id,
		typeName: typeName,
	}
}
