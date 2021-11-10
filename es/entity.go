package es

// EntityType for helping building registry
type EntityType interface{}

// EntityFunc for creating an entity
type EntityFunc func(string) Entity

// Entity for models
type Entity interface {
	// GetId return the ID of the aggregate
	GetId() string

	// GetTypeName return the TypeBame of the aggregate
	GetTypeName() string

	// SetNamespace of the entity
	SetNamespace(namespace string)
}
