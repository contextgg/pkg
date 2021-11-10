package es

import (
	"errors"

	"github.com/contextgg/pkg/events"
)

var (
	// ErrInvalidAggregateType is when the aggregate does not implement event.Aggregate.
	ErrInvalidAggregateType = errors.New("Invalid aggregate type")
	// ErrMismatchedEventType occurs when loaded events from ID does not match aggregate type.
	ErrMismatchedEventType = errors.New("mismatched event type and aggregate type")
	// ErrWrongVersion when the version number is wrong
	ErrWrongVersion = errors.New("When we compute the wrong version")
	// ErrCreatingAggregate whoops when creating aggregate
	ErrCreatingAggregate = errors.New("Issue create aggregate")
	// ErrVersionMismatch when the command's version doesn't match the aggregate
	ErrVersionMismatch = errors.New("Version mismatch")
)

// ApplyEventError is when an event could not be applied. It contains the error
// and the event that caused it.
type ApplyEventError struct {
	// Event is the event that caused the error.
	Event events.Event
	// Err is the error that happened when applying the event.
	Err error
}

// Error implements the Error method of the error interface.
func (a ApplyEventError) Error() string {
	return "failed to apply event " + a.Event.String() + ": " + a.Err.Error()
}
