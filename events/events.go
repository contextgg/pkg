package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ErrTypeNotFound when the type isn't found
var ErrTypeNotFound = errors.New("Type not found")

// Metadata is a simple map to store event's metadata
type Metadata = map[string]interface{}

// Event stores the data for every event
type Event struct {
	AggregateNamespace string      `json:"aggregate_namespace"`
	AggregateID        uuid.UUID   `json:"aggregate_id"`
	AggregateType      string      `json:"aggregate_type"`
	Version            int         `json:"version"`
	Type               string      `json:"type"`
	Timestamp          time.Time   `json:"timestamp"`
	Data               interface{} `json:"data"`
	Metadata           Metadata    `json:"metadata"`
}

// String implements the String method of the Event interface.
func (e Event) String() string {
	return fmt.Sprintf("%s@%d", e.Type, e.Version)
}
