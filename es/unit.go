package es

import (
	"github.com/contextgg/pkg/events"
	"github.com/uptrace/bun"
)

type Unit interface {
	// Error will return an error if there's one!
	Error() error
	// SetError from the error
	SetError(err error)
	// Parent unit
	Parent() Unit
	// Data get the for usage
	Data() Data
	// DB access to the transaction
	DB() *bun.Tx
	// StoreEvents store all the events!
	StoreEvents(evts ...events.Event)
	// Events returns all uncommitted events that are not yet saved.
	Events() []events.Event
	// ClearEvents clears all uncommitted events after saving.
	ClearEvents()
}

type unit struct {
	parent Unit
	db     *bun.DB
	tx     *bun.Tx
	err    error
	events []events.Event
}

func (u *unit) Error() error {
	return u.err
}
func (u *unit) SetError(err error) {
	u.err = err
}
func (u *unit) Parent() Unit {
	return u.parent
}
func (u *unit) DB() *bun.Tx {
	return u.tx
}
func (u *unit) Data() Data {
	return NewData(u.tx)
}
func (u *unit) StoreEvents(evts ...events.Event) {
	u.events = append(u.events, evts...)
}
func (u *unit) Events() []events.Event {
	return u.events
}
func (u *unit) ClearEvents() {
	u.events = nil
}
