package es

import (
	"context"
	"errors"

	"github.com/contextgg/pkg/events"
)

var (
	// ErrInvalidAggregateType is when the aggregate does not implement event.Aggregte.
	ErrInvalidAggregateType = errors.New("Invalid aggregate type")
	// ErrMismatchedEventType occurs when loaded events from ID does not match aggregate type.
	ErrMismatchedEventType = errors.New("mismatched event type and aggregate type")
	// ErrWrongVersion when the version number is wrong
	ErrWrongVersion = errors.New("When we compute the wrong version")
	// ErrCreatingAggregate whoops when creating aggregate
	ErrCreatingAggregate = errors.New("Issue create aggregate")
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

func isReplay(cmd Command) bool {
	// handle the command
	switch cmd.(type) {
	case *ReplayCommand:
		return true
	default:
		return false
	}
}

// NewAggregateSourcedHandler creates the commandhandler
func NewAggregateSourcedHandler(store Store) CommandHandler {
	return CommandHandlerFunc(func(ctx context.Context, cmd Command) error {
		replay := isReplay(cmd)

		id := cmd.GetAggregateId()
		aggregate, err := store.Load(ctx, id, replay)
		if err != nil {
			return err
		}

		// handle the command
		if !replay {
			if handler, ok := aggregate.(CommandHandler); ok {
				if err := handler.HandleCommand(ctx, cmd); err != nil {
					return err
				}
			}
		}

		return store.Save(ctx, aggregate)
	})
}
