package es

import "github.com/google/uuid"

func IsPQNoRow(err error) bool {
	return err != nil && err.Error() == "pg: no rows in result set"
}

// EntityFunc for creating an entity
type EntityFunc func(uuid.UUID) Entity

// Entity for models
type Entity interface {
	// GetNamespace get the namespace for an entity
	GetNamespace() string

	// GetID return the ID of the aggregate
	GetID() uuid.UUID

	// GetTypeName return the TypeBame of the aggregate
	GetTypeName() string
}
