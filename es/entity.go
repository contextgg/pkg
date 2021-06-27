package es

func IsPQNoRow(err error) bool {
	return err != nil && err.Error() == "pg: no rows in result set"
}

// EntityFunc for creating an entity
type EntityFunc func(string) Entity

// Entity for models
type Entity interface {
	// GetID return the ID of the aggregate
	GetID() string

	// GetTypeName return the TypeBame of the aggregate
	GetTypeName() string

	// SetNamespace of the entity
	SetNamespace(namespace string)
}
